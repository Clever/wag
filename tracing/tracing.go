package tracing

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// propagator to use
var propagator propagation.TextMapPropagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}) // traceparent header
type Option interface {
	apply(*options)
}
type options struct {
	address string
}

// WithAddress takes an address in the form of Host:Port
func WithAddress(addr string) Option {
	return addrOption{address: addr}
}

type addrOption struct {
	address string
}

func (o addrOption) apply(opts *options) {
	opts.address = o.address
}

// OtlpGrpcExporter uses the otlptracegrpc modules and the otlptrace module to produce a new exporter at our default addr
func OtlpGrpcExporter(ctx context.Context, opts ...Option) (sdktrace.SpanExporter, error) {

	const DefaultCollectorHost = "localhost"
	const defaultCollectorPort uint16 = 4317
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
		return nil, err
	}
	return spanExporter, nil

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
