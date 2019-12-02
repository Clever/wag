package server

// Code auto-generated. Do not edit.

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	// register pprof listener
	_ "net/http/pprof"

	"github.com/Clever/go-process-metrics/metrics"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kardianos/osext"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
	"gopkg.in/Clever/kayvee-go.v6/logger"
	kvMiddleware "gopkg.in/Clever/kayvee-go.v6/middleware"
	"gopkg.in/tylerb/graceful.v1"
)

const (
	// lowerBoundRateLimiter determines the lower bound interval that we sample every operation.
	// https://godoc.org/github.com/uber/jaeger-client-go#GuaranteedThroughputProbabilisticSampler
	lowerBoundRateLimiter = 1.0 / 60 // 1 request/minute/operation
)

type contextKey struct{}

// Server defines a HTTP server that implements the Controller interface.
type Server struct {
	// Handler should generally not be changed. It exposed to make testing easier.
	Handler http.Handler
	addr    string
	l       logger.KayveeLogger
}

// Serve starts the server. It will return if an error occurs.
func (s *Server) Serve() error {
	tracingToken := os.Getenv("TRACING_ACCESS_TOKEN")
	ingestURL := os.Getenv("TRACING_INGEST_URL")
	isLocal := os.Getenv("_IS_LOCAL") == "true"

	if !isLocal {
		go startLoggingProcessMetrics()
	}

	go func() {
		// This should never return. Listen on the pprof port
		log.Printf("PProf server crashed: %s", http.ListenAndServe("localhost:6060", nil))
	}()

	dir, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.SetGlobalRouting(path.Join(dir, "kvconfig.yml")); err != nil {
		s.l.Info("please provide a kvconfig.yml file to enable app log routing")
	}

	if (tracingToken != "" && ingestURL != "") || isLocal {
		samplingRate := .01 // 1% of requests

		if samplingRateStr := os.Getenv("TRACING_SAMPLING_RATE_PERCENT"); samplingRateStr != "" {
			samplingRateP, err := strconv.ParseFloat(samplingRateStr, 64)
			if err != nil {
				s.l.ErrorD("tracing-sampling-override-failed", logger.M{
					"msg": fmt.Sprintf("could not parse '%s' to integer", samplingRateStr),
				})
			} else {
				samplingRate = samplingRateP
			}

			s.l.InfoD("tracing-sampling-rate", logger.M{
				"msg": fmt.Sprintf("sampling rate will be %.3f", samplingRate),
			})
		}

		sampler, err := jaeger.NewGuaranteedThroughputProbabilisticSampler(lowerBoundRateLimiter, samplingRate)
		if err != nil {
			return fmt.Errorf("failed to build jaeger sampler: %s", err)
		}

		cfg := &jaegercfg.Configuration{
			ServiceName: os.Getenv("_APP_NAME"),
			Tags: []opentracing.Tag{
				opentracing.Tag{Key: "app_name", Value: os.Getenv("_APP_NAME")},
				opentracing.Tag{Key: "build_id", Value: os.Getenv("_BUILD_ID")},
				opentracing.Tag{Key: "deploy_env", Value: os.Getenv("_DEPLOY_ENV")},
				opentracing.Tag{Key: "team_owner", Value: os.Getenv("_TEAM_OWNER")},
				opentracing.Tag{Key: "pod_id", Value: os.Getenv("_POD_ID")},
				opentracing.Tag{Key: "pod_shortname", Value: os.Getenv("_POD_SHORTNAME")},
				opentracing.Tag{Key: "pod_account", Value: os.Getenv("_POD_ACCOUNT")},
				opentracing.Tag{Key: "pod_region", Value: os.Getenv("_POD_REGION")},
			},
		}

		var tracer opentracing.Tracer
		var closer io.Closer
		if isLocal {
			// when local, send everything and use the default params for the Jaeger collector
			cfg.Sampler = &jaegercfg.SamplerConfig{
				Type:  "const",
				Param: 1.0,
			}
			tracer, closer, err = cfg.NewTracer()
			s.l.InfoD("local-tracing", logger.M{"msg": "sending traces to default localhost jaeger address"})
		} else {
			// Create a Jaeger HTTP Thrift transport
			transport := transport.NewHTTPTransport(ingestURL, transport.HTTPBasicAuth("auth", tracingToken))
			tracer, closer, err = cfg.NewTracer(
				jaegercfg.Reporter(jaeger.NewRemoteReporter(transport)),
				jaegercfg.Sampler(sampler))
		}
		if err != nil {
			log.Fatalf("Could not initialize jaeger tracer: %s", err)
		}
		defer closer.Close()

		opentracing.SetGlobalTracer(tracer)
	} else {
		s.l.Error("please set TRACING_ACCESS_TOKEN & TRACING_INGEST_URL to enable tracing")
	}

	s.l.Counter("server-started")

	// Give the sever 30 seconds to shut down
	return graceful.RunWithErr(s.addr, 30*time.Second, s.Handler)
}

type handler struct {
	Controller
}

func startLoggingProcessMetrics() {
	metrics.Log("swagger-test", 1*time.Minute)
}

func withMiddleware(serviceName string, router http.Handler, m []func(http.Handler) http.Handler) http.Handler {
	handler := router

	// compress everything
	handler = handlers.CompressHandler(handler)

	// Wrap the middleware in the opposite order specified so that when called then run
	// in the order specified
	for i := len(m) - 1; i >= 0; i-- {
		handler = m[i](handler)
	}
	handler = TracingMiddleware(handler)
	handler = PanicMiddleware(handler)
	// Logging middleware comes last, i.e. will be run first.
	// This makes it so that other middleware has access to the logger
	// that kvMiddleware injects into the request context.
	handler = kvMiddleware.New(handler, serviceName)
	return handler
}

// New returns a Server that implements the Controller interface. It will start when "Serve" is called.
func New(c Controller, addr string) *Server {
	return NewWithMiddleware(c, addr, []func(http.Handler) http.Handler{})
}

// NewRouter returns a mux.Router with no middleware. This is so we can attach additional routes to the
// router if necessary
func NewRouter(c Controller) *mux.Router {
	return newRouter(c)
}

func newRouter(c Controller) *mux.Router {
	router := mux.NewRouter()
	h := handler{Controller: c}

	router.Methods("GET").Path("/v1/books/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getBook")
		h.GetBookHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getBook")
		r = r.WithContext(ctx)
	})

	return router
}

// NewWithMiddleware returns a Server that implemenets the Controller interface. It runs the
// middleware after the built-in middleware (e.g. logging), but before the controller methods.
// The middleware is executed in the order specified. The server will start when "Serve" is called.
func NewWithMiddleware(c Controller, addr string, m []func(http.Handler) http.Handler) *Server {
	router := newRouter(c)

	return AttachMiddleware(router, addr, m)
}

// AttachMiddleware attaches the given middleware to the router; this is to be used in conjunction with
// NewServer. It attaches custom middleware passed as arguments as well as the built-in middleware for
// logging, tracing, and handling panics. It should be noted that the built-in middleware executes first
// followed by the passed in middleware (in the order specified).
func AttachMiddleware(router *mux.Router, addr string, m []func(http.Handler) http.Handler) *Server {
	l := logger.New("swagger-test")

	handler := withMiddleware("swagger-test", router, m)
	return &Server{Handler: handler, addr: addr, l: l}
}
