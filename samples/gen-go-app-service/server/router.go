package server

// Code auto-generated. Do not edit.

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
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
	config  serverConfig
}

type serverConfig struct {
	compressionLevel int
}

func CompressionLevel(level int) func(*serverConfig) {
	return func(c *serverConfig) {
		c.compressionLevel = level
	}
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
	server := &http.Server{
		Addr:        s.addr,
		Handler:     s.Handler,
		IdleTimeout: 3 * time.Minute,
	}
	server.SetKeepAlivesEnabled(true)

	// Give the server 30 seconds to shut down gracefully after it receives a signal
	shutdown := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Signal(syscall.SIGTERM))
		sig := <-c
		s.l.CriticalD("shutdown-initiated", logger.M{"signal": sig.String()})
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		defer close(shutdown)
		if err := server.Shutdown(ctx); err != nil {
			s.l.CriticalD("error-during-shutdown", logger.M{"error": err.Error()})
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	// ensure we wait for graceful shutdown
	<-shutdown

	return nil
}

type handler struct {
	Controller
}

func startLoggingProcessMetrics() {
	metrics.Log("app-service", 1*time.Minute)
}

func withMiddleware(serviceName string, router http.Handler, m []func(http.Handler) http.Handler, config serverConfig) http.Handler {
	handler := router

	// compress everything
	handler = handlers.CompressHandlerLevel(handler, config.compressionLevel)

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
func New(c Controller, addr string, options ...func(*serverConfig)) *Server {
	return NewWithMiddleware(c, addr, []func(http.Handler) http.Handler{}, options...)
}

// NewRouter returns a mux.Router with no middleware. This is so we can attach additional routes to the
// router if necessary
func NewRouter(c Controller) *mux.Router {
	return newRouter(c)
}

func newRouter(c Controller) *mux.Router {
	router := mux.NewRouter()
	h := handler{Controller: c}

	router.Methods("GET").Path("/_health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "healthCheck")
		h.HealthCheckHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "healthCheck")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/admins").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAdmins")
		h.GetAdminsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAdmins")
		r = r.WithContext(ctx)
	})

	router.Methods("DELETE").Path("/v1/admins/{adminID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "deleteAdmin")
		h.DeleteAdminHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "deleteAdmin")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/admins/{adminID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAdminByID")
		h.GetAdminByIDHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAdminByID")
		r = r.WithContext(ctx)
	})

	router.Methods("PATCH").Path("/v1/admins/{adminID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "updateAdmin")
		h.UpdateAdminHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "updateAdmin")
		r = r.WithContext(ctx)
	})

	router.Methods("PUT").Path("/v1/admins/{adminID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "createAdmin")
		h.CreateAdminHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "createAdmin")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/admins/{adminID}/apps").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAppsForAdminDeprecated")
		h.GetAppsForAdminDeprecatedHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAppsForAdminDeprecated")
		r = r.WithContext(ctx)
	})

	router.Methods("POST").Path("/v1/admins/{adminID}/confirmation_code").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "verifyCode")
		h.VerifyCodeHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "verifyCode")
		r = r.WithContext(ctx)
	})

	router.Methods("PUT").Path("/v1/admins/{adminID}/confirmation_code").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "createVerificationCode")
		h.CreateVerificationCodeHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "createVerificationCode")
		r = r.WithContext(ctx)
	})

	router.Methods("POST").Path("/v1/admins/{adminID}/verify_email").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "verifyAdminEmail")
		h.VerifyAdminEmailHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "verifyAdminEmail")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/analytics/apps").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAllAnalyticsApps")
		h.GetAllAnalyticsAppsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAllAnalyticsApps")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/analytics/apps/{shortname}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAnalyticsAppByShortname")
		h.GetAnalyticsAppByShortnameHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAnalyticsAppByShortname")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/analytics/trackable_apps").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAllTrackableApps")
		h.GetAllTrackableAppsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAllTrackableApps")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/analytics/usageUrls").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAnalyticsUsageUrls")
		h.GetAnalyticsUsageUrlsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAnalyticsUsageUrls")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/appUniverse/usageUrls").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAllUsageUrls")
		h.GetAllUsageUrlsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAllUsageUrls")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getApps")
		h.GetAppsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getApps")
		r = r.WithContext(ctx)
	})

	router.Methods("DELETE").Path("/v1/apps/{appID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "deleteApp")
		h.DeleteAppHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "deleteApp")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAppByID")
		h.GetAppByIDHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAppByID")
		r = r.WithContext(ctx)
	})

	router.Methods("PATCH").Path("/v1/apps/{appID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "updateApp")
		h.UpdateAppHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "updateApp")
		r = r.WithContext(ctx)
	})

	router.Methods("PUT").Path("/v1/apps/{appID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "createApp")
		h.CreateAppHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "createApp")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/admins").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAdminsForApp")
		h.GetAdminsForAppHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAdminsForApp")
		r = r.WithContext(ctx)
	})

	router.Methods("DELETE").Path("/v1/apps/{appID}/admins/{adminID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "unlinkAppAdmin")
		h.UnlinkAppAdminHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "unlinkAppAdmin")
		r = r.WithContext(ctx)
	})

	router.Methods("PUT").Path("/v1/apps/{appID}/admins/{adminID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "linkAppAdmin")
		h.LinkAppAdminHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "linkAppAdmin")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/admins/{adminID}/guides/{guideID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getGuideConfig")
		h.GetGuideConfigHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getGuideConfig")
		r = r.WithContext(ctx)
	})

	router.Methods("PUT").Path("/v1/apps/{appID}/admins/{adminID}/guides/{guideID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "setGuideConfig")
		h.SetGuideConfigHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "setGuideConfig")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/admins/{adminID}/permissions").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getPermissionsForAdmin")
		h.GetPermissionsForAdminHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getPermissionsForAdmin")
		r = r.WithContext(ctx)
	})

	router.Methods("POST").Path("/v1/apps/{appID}/admins/{adminID}/verify").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "verifyAppAdmin")
		h.VerifyAppAdminHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "verifyAppAdmin")
		r = r.WithContext(ctx)
	})

	router.Methods("POST").Path("/v1/apps/{appID}/business_token").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "generateNewBusinessToken")
		h.GenerateNewBusinessTokenHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "generateNewBusinessToken")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/certifications/{schoolYearStart}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getCertifications")
		h.GetCertificationsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getCertifications")
		r = r.WithContext(ctx)
	})

	router.Methods("POST").Path("/v1/apps/{appID}/certifications/{schoolYearStart}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "setCertifications")
		h.SetCertificationsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "setCertifications")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/customStep").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getSetupStep")
		h.GetSetupStepHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getSetupStep")
		r = r.WithContext(ctx)
	})

	router.Methods("PATCH").Path("/v1/apps/{appID}/customStep").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "createSetupStep")
		h.CreateSetupStepHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "createSetupStep")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/data_rules").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getDataRules")
		h.GetDataRulesHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getDataRules")
		r = r.WithContext(ctx)
	})

	router.Methods("PUT").Path("/v1/apps/{appID}/data_rules").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "setDataRules")
		h.SetDataRulesHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "setDataRules")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/managers").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getManagers")
		h.GetManagersHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getManagers")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/onboarding").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getOnboarding")
		h.GetOnboardingHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getOnboarding")
		r = r.WithContext(ctx)
	})

	router.Methods("PATCH").Path("/v1/apps/{appID}/onboarding").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "updateOnboarding")
		h.UpdateOnboardingHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "updateOnboarding")
		r = r.WithContext(ctx)
	})

	router.Methods("PUT").Path("/v1/apps/{appID}/onboarding").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "initializeOnboarding")
		h.InitializeOnboardingHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "initializeOnboarding")
		r = r.WithContext(ctx)
	})

	router.Methods("DELETE").Path("/v1/apps/{appID}/platform/{clientID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "deletePlatform")
		h.DeletePlatformHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "deletePlatform")
		r = r.WithContext(ctx)
	})

	router.Methods("PATCH").Path("/v1/apps/{appID}/platform/{clientID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "updatePlatform")
		h.UpdatePlatformHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "updatePlatform")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/platforms").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getPlatformsByAppID")
		h.GetPlatformsByAppIDHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getPlatformsByAppID")
		r = r.WithContext(ctx)
	})

	router.Methods("PUT").Path("/v1/apps/{appID}/platforms").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "createPlatform")
		h.CreatePlatformHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "createPlatform")
		r = r.WithContext(ctx)
	})

	router.Methods("DELETE").Path("/v1/apps/{appID}/schema").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "deleteAppSchema")
		h.DeleteAppSchemaHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "deleteAppSchema")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/schema").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAppSchema")
		h.GetAppSchemaHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAppSchema")
		r = r.WithContext(ctx)
	})

	router.Methods("POST").Path("/v1/apps/{appID}/schema").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "createAppSchema")
		h.CreateAppSchemaHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "createAppSchema")
		r = r.WithContext(ctx)
	})

	router.Methods("PUT").Path("/v1/apps/{appID}/schema").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "setAppSchema")
		h.SetAppSchemaHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "setAppSchema")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/secrets").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getSecrets")
		h.GetSecretsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getSecrets")
		r = r.WithContext(ctx)
	})

	router.Methods("PATCH").Path("/v1/apps/{appID}/secrets").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "revokeOldClientSecret")
		h.RevokeOldClientSecretHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "revokeOldClientSecret")
		r = r.WithContext(ctx)
	})

	router.Methods("POST").Path("/v1/apps/{appID}/secrets").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "generateNewClientSecret")
		h.GenerateNewClientSecretHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "generateNewClientSecret")
		r = r.WithContext(ctx)
	})

	router.Methods("PUT").Path("/v1/apps/{appID}/secrets").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "resetClientSecret")
		h.ResetClientSecretHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "resetClientSecret")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/apps/{appID}/sharing").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getRecommendedSharing")
		h.GetRecommendedSharingHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getRecommendedSharing")
		r = r.WithContext(ctx)
	})

	router.Methods("PUT").Path("/v1/apps/{appID}/sharing").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "setRecommendedSharing")
		h.SetRecommendedSharingHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "setRecommendedSharing")
		r = r.WithContext(ctx)
	})

	router.Methods("POST").Path("/v1/apps/{appID}/update_icon").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "updateAppIcon")
		h.UpdateAppIconHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "updateAppIcon")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/categories").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAllCategories")
		h.GetAllCategoriesHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAllCategories")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/knownhosts").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getKnownHosts")
		h.GetKnownHostsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getKnownHosts")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/libraryResources").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAllLibraryResources")
		h.GetAllLibraryResourcesHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAllLibraryResources")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/libraryResources/search").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "searchLibraryResource")
		h.SearchLibraryResourceHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "searchLibraryResource")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/libraryResources/{shortname}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getLibraryResourceByShortname")
		h.GetLibraryResourceByShortnameHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getLibraryResourceByShortname")
		r = r.WithContext(ctx)
	})

	router.Methods("PATCH").Path("/v1/libraryResources/{shortname}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "updateLibraryResourceByShortname")
		h.UpdateLibraryResourceByShortnameHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "updateLibraryResourceByShortname")
		r = r.WithContext(ctx)
	})

	router.Methods("POST").Path("/v1/libraryResources/{shortname}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "createLibraryResource")
		h.CreateLibraryResourceHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "createLibraryResource")
		r = r.WithContext(ctx)
	})

	router.Methods("DELETE").Path("/v1/libraryResources/{shortname}/link").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "deleteLibraryResourceLink")
		h.DeleteLibraryResourceLinkHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "deleteLibraryResourceLink")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/permissions").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getValidPermissions")
		h.GetValidPermissionsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getValidPermissions")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/platforms").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getPlatforms")
		h.GetPlatformsHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getPlatforms")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v1/platforms/{clientID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getPlatformByClientID")
		h.GetPlatformByClientIDHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getPlatformByClientID")
		r = r.WithContext(ctx)
	})

	router.Methods("GET").Path("/v2/admins/{adminID}/apps").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "getAppsForAdmin")
		h.GetAppsForAdminHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "getAppsForAdmin")
		r = r.WithContext(ctx)
	})

	router.Methods("POST").Path("/v2/apps/{srcAppID}/override-config/{destAppID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "overrideConfig")
		h.OverrideConfigHandler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "overrideConfig")
		r = r.WithContext(ctx)
	})

	return router
}

// NewWithMiddleware returns a Server that implemenets the Controller interface. It runs the
// middleware after the built-in middleware (e.g. logging), but before the controller methods.
// The middleware is executed in the order specified. The server will start when "Serve" is called.
func NewWithMiddleware(c Controller, addr string, m []func(http.Handler) http.Handler, options ...func(*serverConfig)) *Server {
	router := newRouter(c)

	return AttachMiddleware(router, addr, m, options...)
}

// AttachMiddleware attaches the given middleware to the router; this is to be used in conjunction with
// NewServer. It attaches custom middleware passed as arguments as well as the built-in middleware for
// logging, tracing, and handling panics. It should be noted that the built-in middleware executes first
// followed by the passed in middleware (in the order specified).
func AttachMiddleware(router *mux.Router, addr string, m []func(http.Handler) http.Handler, options ...func(*serverConfig)) *Server {
	// Set sane defaults, to be overriden by the varargs functions.
	// This would probably be better done in NewWithMiddleware, but there are services that call
	// AttachMiddleWare directly instead.
	config := serverConfig{
		compressionLevel: gzip.DefaultCompression,
	}
	for _, option := range options {
		option(&config)
	}

	l := logger.New("app-service")

	handler := withMiddleware("app-service", router, m, config)
	return &Server{Handler: handler, addr: addr, l: l, config: config}
}
