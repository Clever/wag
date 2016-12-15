package server

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	opentracing "github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
	"gopkg.in/Clever/kayvee-go.v5/logger"
)

// PanicMiddleware logs any panics. For now, we're continue throwing the panic up
// the stack so this may crash the process.
func PanicMiddleware(h http.Handler) http.Handler {
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

// statusResponseWriter wraps a response writer
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusResponseWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

type tracingOpName struct{}

// WithTracingOpName adds the op name to a context for use by the tracing library. It uses
// a pointer because it's called below in the stack and the only way to pass the info up
// is to have it a set a pointer. Even though it doesn't change the context we still have
// this return a context to maintain the illusion.
func WithTracingOpName(ctx context.Context, opName string) context.Context {
	strPtr := ctx.Value(tracingOpName{}).(*string)
	if strPtr != nil {
		*strPtr = opName
	}
	return ctx
}

// Tracing creates a new span named after the URL path of the request.
// It places this span in the request context, for use by other handlers via opentracing.SpanFromContext()
// If a span exists in request headers, the span created by this middleware will be a child of that span.
func TracingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attempt to join a span by getting trace info from the headers.
		// To start with use the URL as the opName since we haven't gotten to the router yet and
		// the router knows about opNames
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
		// Use a string pointer so layers below can modify it
		strPtr := ""
		newCtx = context.WithValue(newCtx, tracingOpName{}, &strPtr)

		srw := &statusResponseWriter{
			status:         200,
			ResponseWriter: w,
		}

		tags.HTTPMethod.Set(sp, r.Method)
		tags.SpanKind.Set(sp, tags.SpanKindRPCServerEnum)
		tags.HTTPUrl.Set(sp, r.URL.Path)

		defer func() {
			tags.HTTPStatusCode.Set(sp, uint16(srw.status))
			if srw.status >= 500 {
				tags.Error.Set(sp, true)
			}
			// Now that we have the opName let's try setting it
			opName, ok := newCtx.Value(tracingOpName{}).(*string)
			if ok && opName != nil {
				sp.SetOperationName(*opName)
			}
		}()

		h.ServeHTTP(srw, r.WithContext(newCtx))
	})
}
