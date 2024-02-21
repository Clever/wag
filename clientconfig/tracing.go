package clientconfig

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

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
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
		otelhttp.WithSpanNameFormatter(func(method string, r *http.Request) string {
			// I passed in appName because ideally I'd also include it here in the span name.
			// Currently left out to avoid other possible issues.
			v, ok := r.Context().Value("otelSpanName").(string)
			if ok {
				return v
			}
			// return fmt.Sprintf("%s-wagclient %s %s", rt.appName, r.Method, r.URL.Path)

			return r.Method
		}),
	).RoundTrip(r)
}
