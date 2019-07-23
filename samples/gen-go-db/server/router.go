package server

// Code auto-generated. Do not edit.

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	// register pprof listener
	_ "net/http/pprof"

	"github.com/Clever/go-process-metrics/metrics"
	"github.com/gorilla/mux"
	"github.com/kardianos/osext"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
	"gopkg.in/Clever/kayvee-go.v6/logger"
	kvMiddleware "gopkg.in/Clever/kayvee-go.v6/middleware"
	"gopkg.in/tylerb/graceful.v1"
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

	go func() {
		metrics.Log("swagger-test", 1*time.Minute)
	}()

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

	tracingToken := os.Getenv("TRACING_ACCESS_TOKEN")
	ingestURL := os.Getenv("TRACING_INGEST_URL")
	isLocal := os.Getenv("_IS_LOCAL") == "true"
	if (tracingToken != "" && ingestURL != "") || isLocal {
		// Add rate limited sampling. We will only sample [Param] requests per second
		// and [MaxOperations] different endpoints. Any endpoint above the [MaxOperations]
		// limit will be probabilistically sampled.
		cfgSampler := &jaegercfg.SamplerConfig{
			Type:          jaeger.SamplerTypeRateLimiting,
			Param:         5,
			MaxOperations: 100,
		}
		cfgTags := []opentracing.Tag{
			opentracing.Tag{Key: "app_name", Value: os.Getenv("_APP_NAME")},
			opentracing.Tag{Key: "build_id", Value: os.Getenv("_BUILD_ID")},
			opentracing.Tag{Key: "deploy_env", Value: os.Getenv("_DEPLOY_ENV")},
			opentracing.Tag{Key: "team_owner", Value: os.Getenv("_TEAM_OWNER")},
		}

		cfg := &jaegercfg.Configuration{
			ServiceName: "swagger-test",
			Sampler:     cfgSampler,
			Tags:        cfgTags,
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
			tracer, closer, err = cfg.NewTracer(jaegercfg.Reporter(jaeger.NewRemoteReporter(transport)))
		}
		if err != nil {
			log.Fatalf("Could not initialize jaeger tracer: %!s(MISSING)", err)
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

func withMiddleware(serviceName string, router http.Handler, m []func(http.Handler) http.Handler) http.Handler {
	handler := router

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

// NewWithMiddleware returns a Server that implemenets the Controller interface. It runs the
// middleware after the built-in middleware (e.g. logging), but before the controller methods.
// The middleware is executed in the order specified. The server will start when "Serve" is called.
func NewWithMiddleware(c Controller, addr string, m []func(http.Handler) http.Handler) *Server {
	router := mux.NewRouter()
	h := handler{Controller: c}

	l := logger.New("swagger-test")

	router.Methods("GET").Path("/v1/health/check").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "healthCheck")
		h.HealthCheckHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "healthCheck")
		r = r.WithContext(ctx)
	})

	handler := withMiddleware("swagger-test", router, m)
	return &Server{Handler: handler, addr: addr, l: l}
}
