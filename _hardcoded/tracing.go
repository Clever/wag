package tracing

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagation"
	sdkexporttrace "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/export/trace/tracetest"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// propagator to use.
var propagator propagation.TextMapPropagator = propagation.NewCompositeTextMapPropagator(
	propagation.TraceContext{}, // traceparent header
	xray.Propagator{},          // x-amzn-trace-id header
)

// defaultCollectorPort was changed from 55860 in November and the Go library
// hasn't been updated when it is updated we can use otlp.DefaultCollectorPort
var defaultCollectorPort uint16 = 4317

// SetupGlobalTraceProviderAndExporter sets up an exporter to export,
// as well as the opentelemetry global trace provider for trace generators to use.
// The exporter and provider are returned in order for the caller to defer shutdown.
func SetupGlobalTraceProviderAndExporter() (sdkexporttrace.SpanExporter, *sdktrace.TracerProvider, error) {
	// 1. set up exporter
	// 2. set up tracer provider
	// 3. assign global tracer provider

	// If we're running locally, then turn off sampling. Otherwise sample
	// 1% or whatever TRACING_SAMPLING_PROBABILITY specifies.
	samplingProbability := 0.01
	isLocal := os.Getenv("_IS_LOCAL") == "true"
	if isLocal {
		samplingProbability = 1.0
	} else if v := os.Getenv("TRACING_SAMPLING_PROBABILITY"); v != "" {
		samplingProbabilityFromEnv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("could not parse '%s' to integer", v)
		}
		samplingProbability = samplingProbabilityFromEnv
	}

	// Every 15 seconds we'll try to connect to opentelemetry collector at
	// the default location of localhost:4317
	// When running in production this is a sidecar, and when running
	// locally this is a locally running opetelemetry-collector.
	exporter, err := otlp.NewExporter(
		context.Background(),
		otlp.WithAddress(fmt.Sprintf("%s:%s", otlp.DefaultCollectorHost, defaultCollectorPort)),
		otlp.WithReconnectionPeriod(15*time.Second),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating exporter: %v", err)
	}

	tp := newTracerProvider(exporter, samplingProbability)
	otel.SetTracerProvider(tp)
	return exporter, tp, nil
}

// SetupGlobalTraceProviderAndExporterForTest is meant to be used in unit testing,
// and mirrors the setup above for outside of unit testing. It returns an in-memory
// exporter for examining generated spans.
func SetupGlobalTraceProviderAndExporterForTest() (*tracetest.InMemoryExporter, *sdktrace.TracerProvider, error) {
	exporter := tracetest.NewInMemoryExporter()
	tp := newTracerProvider(exporter, 1.0)
	otel.SetTracerProvider(tp)
	return exporter, tp, nil
}

func newTracerProvider(exporter sdkexporttrace.SpanExporter, samplingProbability float64) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{
			IDGenerator:          xray.NewIDGenerator(),
			DefaultSampler:       sdktrace.ParentBased(sdktrace.TraceIDRatioBased(samplingProbability)),
			MaxEventsPerSpan:     100,
			MaxAttributesPerSpan: 100,
			MaxLinksPerSpan:      100,
		}),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			label.KeyValue{Key: "app_name", Value: label.StringValue(os.Getenv("_APP_NAME"))},
			label.KeyValue{Key: "build_id", Value: label.StringValue(os.Getenv("_BUILD_ID"))},
			label.KeyValue{Key: "deploy_env", Value: label.StringValue(os.Getenv("_DEPLOY_ENV"))},
			label.KeyValue{Key: "team_owner", Value: label.StringValue(os.Getenv("_TEAM_OWNER"))},
			label.KeyValue{Key: "pod_id", Value: label.StringValue(os.Getenv("_POD_ID"))},
			label.KeyValue{Key: "pod_shortname", Value: label.StringValue(os.Getenv("_POD_SHORTNAME"))},
			label.KeyValue{Key: "pod_account", Value: label.StringValue(os.Getenv("_POD_ACCOUNT"))},
			label.KeyValue{Key: "pod_region", Value: label.StringValue(os.Getenv("_POD_REGION"))},
		)),
	)
}

// MuxServerMiddelware returns middleware that should be attached to a gorilla/mux server.
func MuxServerMiddleware(serviceName string) func(http.Handler) http.Handler {
	return otelmux.Middleware(serviceName, otelmux.WithPropagators(propagator))
}

// NewTransport returns the transport to use in client requests.
// It takes in a transport to wrap, e.g. http.DefaultTransport, and the request
// context value to pull the span name out from.
// The exporter is pulled from the global one on each request, so tracing won't
// begin until that is initialized (e.g, in in server startup).
func NewTransport(baseTransport http.RoundTripper, spanNameCtxValue interface{}) http.RoundTripper {
	return roundTripper{baseTransport: baseTransport, spanNameCtxValue: spanNameCtxValue}
}

type roundTripper struct {
	baseTransport    http.RoundTripper
	spanNameCtxValue interface{}
}

func (rt roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return otelhttp.NewTransport(
		rt.baseTransport,
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(propagator),
		otelhttp.WithSpanNameFormatter(func(method string, r *http.Request) string {
			v, ok := r.Context().Value(rt.spanNameCtxValue).(string)
			if ok {
				return v
			}
			return r.Method // same as otelhttp's default span naming
		}),
	).RoundTrip(r)
}

// ExtractSpanAndTraceID extracts span and trace IDs from an http request header.
func ExtractSpanAndTraceID(r *http.Request) (traceID, spanID string) {
	sc := trace.RemoteSpanContextFromContext(propagator.Extract(r.Context(), r.Header))
	return sc.SpanID.String(), sc.TraceID.String()
}
