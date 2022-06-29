package clienttracing

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/Clever/kayvee-go.v6/logger"
)

// propagator to use.
var propagator propagation.TextMapPropagator = propagation.TraceContext{} // traceparent header

// defaultCollectorPort was changed from 55860 in November and the Go library
// hasn't been updated when it is updated we can use otlp.DefaultCollectorPort
var defaultCollectorPort uint16 = 4317

// SetupGlobalTraceProviderAndExporter sets up an exporter to export,
// as well as the opentelemetry global trace provider for trace generators to use.
// The exporter and provider are returned in order for the caller to defer shutdown.
func SetupGlobalTraceProviderAndExporter(ctx context.Context) (sdktrace.SpanExporter, *sdktrace.TracerProvider, error) {
	// 1. set up exporter
	// 2. set up tracer provider
	// 3. assign global tracer provider

	// If we're running locally, then turn off sampling. Otherwise sample
	// 1% or whatever TRACING_SAMPLING_PROBABILITY specifies.
	samplingProbability := 0.01
	isLocal := os.Getenv("_IS_LOCAL") == "true"
	if isLocal {
		samplingProbability = 1.0
	} else if v := os.Getenv("TRACING_SAMPLING_PROBABILITY"); v != "" {
		samplingProbabilityFromEnv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("could not parse '%s' to float", v)
		}
		samplingProbability = samplingProbabilityFromEnv
	}

	addr := fmt.Sprintf("%s:%d", otlp.DefaultCollectorHost, defaultCollectorPort)

	// Every 15 seconds we'll try to connect to opentelemetry collector at
	// the default location of localhost:4317
	// When running in production this is a sidecar, and when running
	// locally this is a locally running opetelemetry-collector.
	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithReconnectionPeriod(15*time.Second),
		otlpgrpc.WithEndpoint(addr),
		otlpgrpc.WithInsecure(),
	)
	exporter, err := otlp.NewExporter(
		ctx,
		driver,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating exporter: %v", err)
	}

	tp := newTracerProvider(exporter, samplingProbability)
	otel.SetTracerProvider(tp)
	logger.FromContext(ctx).InfoD("starting-tracer", logger.M{
		"address":       addr,
		"sampling-rate": samplingProbability,
	})
	return exporter, tp, nil
}

// SetupGlobalTraceProviderAndExporterForTest is meant to be used in unit testing,
// and mirrors the setup above for outside of unit testing. It returns an in-memory
// exporter for examining generated spans.
func SetupGlobalTraceProviderAndExporterForTest() (*tracetest.InMemoryExporter, *sdktrace.TracerProvider, error) {
	exporter := tracetest.NewInMemoryExporter()
	tp := newTracerProvider(exporter, 1.0)
	otel.SetTracerProvider(tp)
	return exporter, tp, nil
}

func newTracerProvider(exporter sdktrace.SpanExporter, samplingProbability float64) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		// We use the default ID generator. In order for sampling to work (at least with this sampler)
		// the ID generator must generate trace IDs uniformly at random from the entire space of uint64.
		// For example, the default x-ray ID generator does not do this.
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(samplingProbability))),
		// These maximums are to guard against something going wrong and sending a ton of data unexpectedly
		sdktrace.WithSpanLimits(sdktrace.SpanLimits{
			AttributeCountLimit: 100,
			EventCountLimit:     100,
			LinkCountLimit:      100,
		}),
		sdktrace.WithSyncer(exporter),
	)
}

// MuxServerMiddleware returns middleware that should be attached to a gorilla/mux server.
// It does two things: starts spans, and adds span/trace info to the request-specific logger.
// Right now we only support logging IDs in the format that Datadog expects.
func MuxServerMiddleware(serviceName string) func(http.Handler) http.Handler {
	otlmux := otelmux.Middleware(serviceName, otelmux.WithPropagators(propagator))
	return func(h http.Handler) http.Handler {
		return otlmux(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// otelmux has extracted the span. now put it into the ctx-specific logger
			s := trace.SpanFromContext(r.Context())
			if sc := s.SpanContext(); sc.HasTraceID() {
				spanID, traceID := sc.SpanID().String(), sc.TraceID().String()
				// datadog converts hex strings to uint64 IDs, so log those so that correlating logs and traces works
				if len(traceID) == 32 && len(spanID) == 16 { // opentelemetry format: 16 byte (32-char hex), 8 byte (16-char hex) trace and span ids
					traceIDBs, _ := hex.DecodeString(traceID)
					logger.FromContext(r.Context()).AddContext("trace_id",
						fmt.Sprintf("%d", binary.BigEndian.Uint64(traceIDBs[8:])))
					spanIDBs, _ := hex.DecodeString(spanID)
					logger.FromContext(r.Context()).AddContext("span_id",
						fmt.Sprintf("%d", binary.BigEndian.Uint64(spanIDBs)))
				}
			}
			h.ServeHTTP(rw, r)
		}))
	}
}

// NewTransport returns the transport to use in client requests.
// It takes in a transport to wrap, e.g. http.DefaultTransport, and the request
// context value to pull the span name out from.
// The exporter is pulled from the global one on each request, so tracing won't
// begin until that is initialized (e.g, in in server startup).
func NewTransport(baseTransport http.RoundTripper, spanNameCtxValue interface{}) http.RoundTripper {
	return roundTripper{baseTransport: baseTransport, spanNameCtxValue: spanNameCtxValue}
}

type roundTripper struct {
	baseTransport    http.RoundTripper
	spanNameCtxValue interface{}
}

func (rt roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return otelhttp.NewTransport(
		rt.baseTransport,
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(propagator),
		otelhttp.WithSpanNameFormatter(func(method string, r *http.Request) string {
			v, ok := r.Context().Value(rt.spanNameCtxValue).(string)
			if ok {
				return v
			}
			return r.Method // same as otelhttp's default span naming
		}),
	).RoundTrip(r)
}

// ExtractSpanAndTraceID extracts span and trace IDs from an http request header.
func ExtractSpanAndTraceID(r *http.Request) (traceID, spanID string) {
	s := trace.SpanFromContext(r.Context())
	if s.SpanContext().HasTraceID() {
		return s.SpanContext().TraceID().String(), s.SpanContext().SpanID().String()
	}
	sc := trace.SpanContextFromContext(propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header)))
	return sc.SpanID().String(), sc.TraceID().String()
}
