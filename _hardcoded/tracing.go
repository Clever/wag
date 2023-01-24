package servertracing

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Clever/kayvee-go/v7/logger"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

var defaultCollectorHost string = "localhost"
var defaultCollectorPort uint16 = 4317

func SetupGlobalTraceProviderAndExporter(ctx context.Context) (sdktrace.SpanExporter, *sdktrace.TracerProvider, error) {

	addr := fmt.Sprintf("%s:%d", defaultCollectorHost, defaultCollectorPort)

	// Every 15 seconds we'll try to connect to opentelemetry collector at
	// the default location of localhost:4317
	// When running in production this is a sidecar, and when running
	// locally this is a locally running opetelemetry-collector.
	otlpClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithReconnectionPeriod(15*time.Second),
		otlptracegrpc.WithEndpoint(addr),
		otlptracegrpc.WithInsecure(),
	)
	spanExporter, err := otlptrace.New(ctx, otlpClient)
	// spanExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))

	if err != nil {
		return nil, nil, fmt.Errorf("error creating exporter: %v", err)
	}

	tp := newTracerProvider(spanExporter, newResource())
	otel.SetTracerProvider(tp)

	logger.FromContext(ctx).InfoD("starting-tracer", logger.M{
		"address": addr,
	})
	return spanExporter, tp, nil
}

func newTracerProvider(exporter sdktrace.SpanExporter, resource *resource.Resource) *sdktrace.TracerProvider {
	samplingProbability := 0.05
	isLocal := os.Getenv("_IS_LOCAL") == "true"
	if isLocal {
		samplingProbability = 1.0
	} else if v := os.Getenv("TRACING_SAMPLING_PROBABILITY"); v != "" {
		samplingProbabilityFromEnv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			samplingProbabilityFromEnv = 1
		}
		samplingProbability = samplingProbabilityFromEnv
	}

	tp := sdktrace.NewTracerProvider(
		// We use the default ID generator. In order for sampling to work (at least with this sampler)
		// the ID generator must generate trace IDs uniformly at random from the entire space of uint64.
		// For example, the default x-ray ID generator does not do this.
		// sdktrace.WithSampler(sdktrace.TraceIDRatioBased()),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(samplingProbability))),
		// These maximums are to guard against something going wrong and sending a ton of data unexpectedly
		sdktrace.WithSpanLimits(sdktrace.SpanLimits{
			AttributeCountLimit: 100,
			EventCountLimit:     100,
			LinkCountLimit:      100,
		}),

		// Batcher is more efficient, switch to it after testing
		// sdktrace.WithSyncer(exporter),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp
}

// SetupGlobalTraceProviderAndExporterForTest is meant to be used in unit testing,
// and mirrors the setup above for outside of unit testing. It returns an in-memory
// exporter for examining generated spans.
func SetupGlobalTraceProviderAndExporterForTest() (*tracetest.InMemoryExporter, *sdktrace.TracerProvider, error) {
	exporter := tracetest.NewInMemoryExporter()
	tp := newTracerProvider(exporter, newResource())
	otel.SetTracerProvider(tp)
	return exporter, tp, nil
}

// MuxServerMiddleware returns middleware that should be attached to a gorilla/mux server.
// It does two things: starts spans, and adds span/trace info to the request-specific logger.
// Right now we only support logging IDs in the format that Datadog expects.
func MuxServerMiddleware(serviceName string) func(http.Handler) http.Handler {
	otlmux := otelmux.Middleware(serviceName, otelmux.WithPropagators(otel.GetTextMapPropagator()))
	// fmt.Println("Adding mux server middleware")
	return func(h http.Handler) http.Handler {
		return otlmux(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var rid string
			var crid string

			s := trace.SpanFromContext(r.Context())
			bag := baggage.FromContext(r.Context())

			// Prefer to grab values from baggage
			crid = bag.Member("clever-request-id").Value()
			rid = bag.Member("X-Request-ID").Value()

			// If the values aren't set in baggage grab the traceid from otel, and the x-request-id
			// from the headers (set by envoy)
			if crid == "" {
				crid = s.SpanContext().TraceID().String()
			}

			if rid == "" {
				rid = r.Header.Get("X-Request-ID")
			}

			s.SetAttributes(attribute.String("clever-request-id", crid))

			cridMember, err := baggage.NewMember("clever-request-id", crid)
			if err != nil {
				s.RecordError(err)
			}

			bag, err = bag.SetMember(cridMember)
			if err != nil {
				s.RecordError(err)
			}

			logger.FromContext(r.Context()).AddContext("clever-request-id", crid)
			rw.Header().Add("clever-request-id", crid)

			if rid != "" {
				s.SetAttributes(attribute.String("X-Request-ID", rid))
				ridMember, err := baggage.NewMember("X-Request-ID", rid)
				if err != nil {
					s.RecordError(err)
				}
				bag, err = bag.SetMember(ridMember)
				if err != nil {
					s.RecordError(err)
				}

				// Envoy logs store this as request_id so lets match it for easier filtering.
				logger.FromContext(r.Context()).AddContext("request_id", rid)

				rw.Header().Add("X-Request-ID", rid)
			}

			if sc := s.SpanContext(); sc.HasTraceID() {
				spanID, traceID := sc.SpanID().String(), sc.TraceID().String()
				// datadog converts hex strings to uint64 IDs, so log those so that correlating logs and traces works
				if len(traceID) == 32 && len(spanID) == 16 { // opentelemetry format: 16 byte (32-char hex), 8 byte (16-char hex) trace and span ids

					traceIDBs, _ := hex.DecodeString(traceID)
					logger.FromContext(r.Context()).AddContext("dd.trace_id",
						fmt.Sprintf("%d", binary.BigEndian.Uint64(traceIDBs[8:])))
					spanIDBs, _ := hex.DecodeString(spanID)
					logger.FromContext(r.Context()).AddContext("dd.span_id",
						fmt.Sprintf("%d", binary.BigEndian.Uint64(spanIDBs)))
				}
			}

			ctx := baggage.ContextWithBaggage(r.Context(), bag)
			r = r.WithContext(ctx)

			h.ServeHTTP(rw, r)
		}))
	}
}

// newResource returns a resource describing this application.
// Used for setting up tracer provider
func newResource() *resource.Resource {
	var appName string
	if os.Getenv("_APP_NAME") != "" {
		appName = os.Getenv("_APP_NAME")
	} else if os.Getenv("APP_NAME") != "" {
		appName = os.Getenv("APP_NAME")
	}
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		),
	)
	return r
}
