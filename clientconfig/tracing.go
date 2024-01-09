package clientconfig

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var propagator propagation.TextMapPropagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}) // traceparent header

type Option interface {
	apply(*options)
}
type options struct {
	address string
}

// WithAddress takes an address in the form of Host:Port
// Leaving this here as an example of how to add a new option.
// This was removed because we shouldn't be adjusting the exporter for a client connection.

// func WithAddress(addr string) Option {
// 	return addrOption{address: addr}
// }

// type addrOption struct {
// 	address string
// }

// func (o addrOption) apply(opts *options) {
// 	opts.address = o.address
// }

// DefaultInstrumentor returns the transport to use in client requests.
// It takes in a transport to wrap, e.g. http.DefaultTransport, and the request
// context value to pull the span name out from.
// 99% sure this is wrapping a wrapped thing and totally redundant. Fix later.
func DefaultInstrumentor(baseTransport http.RoundTripper, appName string) http.RoundTripper {
	return roundTripperWithTracing{baseTransport: baseTransport, appName: appName}
}

type roundTripperWithTracing struct {
	baseTransport http.RoundTripper
	appName       string
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
			return r.Method
		}),
	).RoundTrip(r)
}
