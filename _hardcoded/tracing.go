package servertracing

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/Clever/kayvee-go/v7/logger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// propagator to use.
var propagator propagation.TextMapPropagator = propagation.TraceContext{} // traceparent header

// MuxServerMiddleware returns middleware that should be attached to a gorilla/mux server.
// It does two things: starts spans, and adds span/trace info to the request-specific logger.
// Right now we only support logging IDs in the format that Datadog expects.
func MuxServerMiddleware(serviceName string) func(http.Handler) http.Handler {
	otlmux := otelmux.Middleware(serviceName, otelmux.WithPropagators(propagator))
	return func(h http.Handler) http.Handler {
		return otlmux(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// otelmux has extracted the span. now put it into the ctx-specific logger
			s := trace.SpanFromContext(r.Context())
			fmt.Println(s)
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
