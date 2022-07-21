package tracing

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/Clever/kayvee-go.v6/logger"
)

const (
	instrumentationName    = "github.com/Clever/wag/instrumentation"
	instrumentationVersion = "v0.1.0"
)

//propagator to use
var propagator propagation.TextMapPropagator = propagation.TraceContext{} // traceparent header
type tracerProviderCreator func(*otlptrace.Exporter, float64) *sdktrace.TracerProvider

//OtlpGrpcExporter uses the otlptracegrpc modules and the otlptrace module to produce a new exporter at our default addr
func OtlpGrpcExporter(ctx context.Context) sdktrace.SpanExporter {
	DefaultCollectorHost := "localhost"
	var defaultCollectorPort uint16 = 4317

	addr := fmt.Sprintf("%s:%d", DefaultCollectorHost, defaultCollectorPort)

	otlpClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(addr), //Not strictly needed as we use the defaults
		otlptracegrpc.WithReconnectionPeriod(15*time.Second),
		otlptracegrpc.WithInsecure(),
	)

	spanExporter, err := otlptrace.New(ctx, otlpClient)
	if err != nil {
		log.Fatal(err)
		//Is doing a fatal error here too risky? No easy way to bubble up errors from here to the app using this.
		//without making each of the WithXOption() takes an error as an arg as well.

		fmt.Println(err)
		return nil
	}
	return spanExporter

}

func RoundTripperInstrumentor(tp sdktrace.TracerProvider, baseRT http.RoundTripper, ctx context.Context) (http.RoundTripper, error) {
	return InstrumentedTransport(baseRT, ctx, tp), nil
	// DefaultCollectorHost := "localhost"

	// var defaultCollectorPort uint16 = 4317
	// addr := fmt.Sprintf("%s:%d", DefaultCollectorHost, defaultCollectorPort)
}

//OurDefaultRoundTripper will return an instrumented round tripper with default tracing and logging
//These defaults are overridden with WithExporter() and WithLogger() at time of client creation (wagclient.New())
func OurDefaultRoundTripper(ctx context.Context, resource *resource.Resource) (http.RoundTripper, error) {
	DefaultCollectorHost := "localhost"
	// defaultCollectorPort was changed from 55860 in November and the Go library
	// hasn't been updated when it is updated we can use otlp.DefaultCollectorPort
	var defaultCollectorPort uint16 = 4317
	// I want to see if this is still true or if I can use otlp.DefaultCollectorPort now.
	// fmt.Println("Are these the same now?: ", otel.defaultCollectorPort, defaultCollectorPort)

	addr := fmt.Sprintf("%s:%d", DefaultCollectorHost, defaultCollectorPort)

	samplingProbability, err := determineSampling()

	samplingProbability = 1 //Temp so I can run without ark start -l and still get sampling

	if err != nil {
		return nil, fmt.Errorf("error determining sampling: %s", err)
	}

	otlpClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(addr), //Not strictly needed as we use the defaults
		otlptracegrpc.WithReconnectionPeriod(15*time.Second),
		otlptracegrpc.WithInsecure(),
	)

	spanExporter, err := otlptrace.New(ctx, otlpClient)

	tracerProvider := newTracerProvider(spanExporter, samplingProbability, resource)

	otel.SetTracerProvider(tracerProvider)

	return http.DefaultTransport, nil
	// return NewTransport(http.DefaultTransport, ctx), nil

	//...
}

func determineSampling() (samplingProbability float64, err error) {

	// If we're running locally, then turn off sampling. Otherwise sample
	// 1% or whatever TRACING_SAMPLING_PROBABILITY specifies.
	samplingProbability = 0.01
	isLocal := os.Getenv("_IS_LOCAL") == "true"
	if isLocal {
		fmt.Println("Set to Local")
		samplingProbability = 1.0
	} else if v := os.Getenv("TRACING_SAMPLING_PROBABILITY"); v != "" {
		samplingProbabilityFromEnv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("could not parse '%s' to float", v)
		}
		samplingProbability = samplingProbabilityFromEnv
	}
	return
}

// My thought is this should be created in either the client or the app using the client and then passed in. Haven't decided yet.
//Perhaps since the client module relies on a "generated" gen-go/tracing.go this could live in there.

// func newResource() *resource.Resource {
// 	r, _ := resource.Merge(
// 		resource.Default(),
// 		resource.NewWithAttributes(
// 			semconv.SchemaURL,
// 			semconv.ServiceNameKey.String("service-name-goes-here"),
// 			semconv.ServiceVersionKey.String("service-name-version-goes-here"),
// 			attribute.String("environment", "demo"),
// 		),
// 	)
// 	return r
// }

func newTracerProvider(exporter sdktrace.SpanExporter, samplingProbability float64, resource *resource.Resource) *sdktrace.TracerProvider {
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
		//Batcher is more efficient, switch to it after testing
		sdktrace.WithSyncer(exporter),
		//sdktrace.WithBatcher(exporter),
		//Have to figure out how I'm going to generate this resource first.
		sdktrace.WithResource(resource),
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

// InstrumentedTransport returns the transport to use in client requests.
// It takes in a transport to wrap, e.g. http.DefaultTransport, and the request
// context value to pull the span name out from.
// 99% sure this is wrapping a wrapped thing and totally redundant. Fix later.
func InstrumentedTransport(baseTransport http.RoundTripper, spanNameCtxValue interface{}, tp sdktrace.TracerProvider) http.RoundTripper {
	return roundTripperWithTracing{baseTransport: baseTransport, spanNameCtxValue: spanNameCtxValue, tp: tp}
}

type roundTripperWithTracing struct {
	baseTransport    http.RoundTripper
	spanNameCtxValue interface{}
	tp               sdktrace.TracerProvider
}

func (rt roundTripperWithTracing) RoundTrip(r *http.Request) (*http.Response, error) {

	return otelhttp.NewTransport(
		rt.baseTransport,
		otelhttp.WithTracerProvider(&rt.tp),
		// otelhttp.WithTracerProvider(tracer),
		otelhttp.WithPropagators(propagator),
		otelhttp.WithSpanNameFormatter(func(method string, r *http.Request) string {
			v, ok := r.Context().Value(rt.spanNameCtxValue).(string)
			if ok {
				fmt.Println("---v---")
				spew.Dump(v)
				return v
			}
			return r.Method // same as otelhttp's default span naming
		}),
	).RoundTrip(r)
}

// ExtractSpanAndTraceID extracts span and trace IDs from an http request header.
func ExtractSpanAndTraceID(r *http.Request) (traceID, spanID string) {
	fmt.Println("Extracting TraceID")
	s := trace.SpanFromContext(r.Context())
	if s.SpanContext().HasTraceID() {
		return s.SpanContext().TraceID().String(), s.SpanContext().SpanID().String()
	}
	sc := trace.SpanContextFromContext(propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header)))
	return sc.SpanID().String(), sc.TraceID().String()
}
