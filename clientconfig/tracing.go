package clientconfig

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
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
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

// DefaultInstrumentor returns the transport to use in client requests.
// It takes in a transport to wrap, e.g. http.DefaultTransport, and the request
// context value to pull the span name out from.
// 99% sure this is wrapping a wrapped thing and totally redundant. Fix later.
func DefaultInstrumentor(baseTransport http.RoundTripper, tp sdktrace.TracerProvider) http.RoundTripper {
	return roundTripperWithTracing{baseTransport: baseTransport, tp: tp}
}

type roundTripperWithTracing struct {
	baseTransport http.RoundTripper
	tp            sdktrace.TracerProvider
}

func (rt roundTripperWithTracing) RoundTrip(r *http.Request) (*http.Response, error) {
	return otelhttp.NewTransport(
		rt.baseTransport,
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(propagator),
		otelhttp.WithSpanNameFormatter(func(method string, r *http.Request) string {
			v, ok := r.Context().Value("otelSpanName").(string)
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

// newResource returns a resource describing this application.
// Used for setting up tracer provider
func newResource(appName string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		),
	)
	return r
}

func newTracerProvider(exporter sdktrace.SpanExporter, appName string) *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider(
		// We use the default ID generator. In order for sampling to work (at least with this sampler)
		// the ID generator must generate trace IDs uniformly at random from the entire space of uint64.
		// For example, the default x-ray ID generator does not do this.
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// These maximums are to guard against something going wrong and sending a ton of data unexpectedly
		sdktrace.WithRawSpanLimits(sdktrace.SpanLimits{
			AttributeCountLimit: 100,
			EventCountLimit:     100,
			LinkCountLimit:      100,
		}),
		//Batcher is more efficient, switch to it after testing
		// sdktrace.WithSyncer(exporter),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource(appName)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
