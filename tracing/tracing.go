package tracing

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

//propagator to use
var propagator propagation.TextMapPropagator = propagation.TraceContext{} // traceparent header
type Option interface {
	apply(*options)
}
type options struct {
	address string
}

//WithAddress takes an address in the form of Host:Port
func WithAddress(addr string) Option {
	return addrOption{address: addr}
}

type addrOption struct {
	address string
}

func (o addrOption) apply(opts *options) {
	opts.address = o.address
}

//OtlpGrpcExporter uses the otlptracegrpc modules and the otlptrace module to produce a new exporter at our default addr
func OtlpGrpcExporter(ctx context.Context, opts ...Option) sdktrace.SpanExporter {

	DefaultCollectorHost := "localhost"
	var defaultCollectorPort uint16 = 4317
	addr := fmt.Sprintf("%s:%d", DefaultCollectorHost, defaultCollectorPort)

	options := options{
		address: addr,
	}

	for _, o := range opts {
		o.apply(&options)
	}

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
		return nil
	}
	return spanExporter

}

func JaegerExporter() (spanExporter sdktrace.SpanExporter) {
	fmt.Println("Creating Jaeger Exporter")
	spanExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		log.Fatal("Error creating Jaeger Exporter")
	}
	return
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
		otelhttp.WithTracerProvider(otel.GetTracerProvider())
		// otelhttp.WithTracerProvider(&rt.tp),
		// otelhttp.WithTracerProvider(tracer),
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
	fmt.Println("Extracting TraceID")
	s := trace.SpanFromContext(r.Context())
	if s.SpanContext().HasTraceID() {
		return s.SpanContext().TraceID().String(), s.SpanContext().SpanID().String()
	}
	sc := trace.SpanContextFromContext(propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header)))
	return sc.SpanID().String(), sc.TraceID().String()
}
