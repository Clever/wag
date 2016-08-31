package server

import (
	"net/http"

	"gopkg.in/Clever/kayvee-go.v4/logger"
	kvMiddleware "gopkg.in/Clever/kayvee-go.v4/middleware"

	opentracing "github.com/opentracing/opentracing-go"
)

func withMiddleware(serviceName string, router http.Handler) http.Handler {

	handler := kvMiddleware.New(router, logger.New(serviceName))
	handler = tracingMiddleware(handler)
	return handler
}

// tracingMiddleware creates a new span named after the URL path of the request.
// It places this span in the request context, for use by other handlers via opentracing.SpanFromContext()
// If a span exists in request headers, the span created by this middleware will be a child of that span.
func tracingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attempt to join a span by getting trace info from the headers.
		opName := r.URL.Path
		var sp opentracing.Span
		if sc, err := opentracing.GlobalTracer().
			Extract(opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header)); err != nil {
			sp = opentracing.StartSpan(opName)
		} else {
			sp = opentracing.StartSpan(opName, opentracing.ChildOf(sc))
		}
		defer sp.Finish()
		sp.LogEvent("request_received")
		defer func() {
			sp.LogEvent("request_finished")
		}()
		newCtx := opentracing.ContextWithSpan(r.Context(), sp)

		h.ServeHTTP(w, r.WithContext(newCtx))
	})
}
