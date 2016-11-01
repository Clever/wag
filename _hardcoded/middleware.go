package server

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"gopkg.in/Clever/kayvee-go.v5/logger"
	kvMiddleware "gopkg.in/Clever/kayvee-go.v5/middleware"

	opentracing "github.com/opentracing/opentracing-go"
)

func withMiddleware(serviceName string, router http.Handler) http.Handler {
	handler := tracingMiddleware(router)
	handler = panicMiddleware(handler)
	// Logging middleware comes last, i.e. will be run first.
	// This makes it so that other middleware has access to the logger
	// that kvMiddleware injects into the request context.
	handler = kvMiddleware.New(handler, serviceName)
	return handler
}

// panicMiddleware logs any panics. For now, we're continue throwing the panic up
// the stack so this may crash the process.
func panicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			panicErr := recover()
			if panicErr == nil {
				return
			}
			var err error

			switch panicErr := panicErr.(type) {
			case string:
				err = fmt.Errorf(panicErr)
			case error:
				err = panicErr
			default:
				err = fmt.Errorf("unknown panic %#v of type %T", panicErr, panicErr)
			}

			logger.FromContext(r.Context()).ErrorD("panic",
				logger.M{"err": err, "stacktrace": string(debug.Stack())})
			panic(panicErr)
		}()
		h.ServeHTTP(w, r)
	})
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

		// inject span ID into logs to aid in request debugging
		t := make(map[string]string)
		if err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap,
			opentracing.TextMapCarrier(t)); err == nil {
			if spanid, ok := t["ot-tracer-spanid"]; ok {
				logger.FromContext(r.Context()).AddContext("ot-tracer-spanid", spanid)
			}
		}

		sp.LogEvent("request_received")
		defer func() {
			sp.LogEvent("request_finished")
		}()
		newCtx := opentracing.ContextWithSpan(r.Context(), sp)

		h.ServeHTTP(w, r.WithContext(newCtx))
	})
}
