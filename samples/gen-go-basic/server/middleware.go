package server

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Clever/kayvee-go/v7/logger"
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

// VersionRange decides whether to accept a version.
type VersionRange func(version string) bool

// ClientVersionCheckMiddleware checks the client version.
func ClientVersionCheckMiddleware(h http.Handler, rng VersionRange) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := r.Header.Get("X-Client-Version")
		logger.FromContext(r.Context()).AddContext("client-version", version)
		if !rng(version) {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf(`{"message": "client version '%s' not accepted, please upgrade"}`, version)))
			return
		}
		h.ServeHTTP(w, r)
	})
}
