package client

// Using Alpha version of WAG Yay!
import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Clever/wag/samples/gen-go-nils/models/v9"

	discovery "github.com/Clever/discovery-go"
	wcl "github.com/Clever/wag/logging/wagclientlogger"

	"github.com/afex/hystrix-go/hystrix"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// Version of the client.
const Version = "9.0.0"

// VersionHeader is sent with every request.
const VersionHeader = "X-Client-Version"

// WagClient is used to make requests to the nil-test service.
type WagClient struct {
	basePath    string
	requestDoer doer
	client      *http.Client
	// Keep the retry doer around so that we can set the number of retries
	retryDoer *retryDoer
	// Keep the circuit doer around so that we can turn it on / off
	circuitDoer    *circuitBreakerDoer
	defaultTimeout time.Duration
	logger         wcl.WagClientLogger
}

var _ Client = (*WagClient)(nil)

//This pattern is used instead of using closures for greater transparency and the ability to implement additional interfaces.
type options struct {
	transport    http.RoundTripper
	logger       wcl.WagClientLogger
	instrumentor Instrumentor
	exporter     sdktrace.SpanExporter
}

type Option interface {
	apply(*options)
}

//WithLogger sets client logger option.
func WithLogger(log wcl.WagClientLogger) Option {
	return loggerOption{Log: log}
}

type loggerOption struct {
	Log wcl.WagClientLogger
}

func (l loggerOption) apply(opts *options) {
	opts.logger = l.Log
}

type roundTripperOption struct {
	rt http.RoundTripper
}

func (t roundTripperOption) apply(opts *options) {
	opts.transport = t.rt
}

// WithRoundTripper allows you to pass in intrumented/custom roundtrippers which will then wrap the
// transport roundtripper
func WithRoundTripper(t http.RoundTripper) Option {
	return roundTripperOption{rt: t}
}

// Instrumentor is a function that creates an instrumented round tripper
type Instrumentor func(baseTransport http.RoundTripper, spanNameCtxValue interface{}, tp sdktrace.TracerProvider) http.RoundTripper

// WithInstrumentor sets a instrumenting function that will be used to wrap the roundTripper for tracing.
// For standard instrumentation with tracing use tracing.InstrumentedTransport, default is non-instrumented.

func WithInstrumentor(fn Instrumentor) Option {
	return instrumentorOption{instrumentor: fn}
}

type instrumentorOption struct {
	instrumentor Instrumentor
}

func (i instrumentorOption) apply(opts *options) {
	opts.instrumentor = i.instrumentor
}

// WithExporter sets client span exporter option.
func WithExporter(se sdktrace.SpanExporter) Option {
	return exporterOption{exporter: se}
}

type exporterOption struct {
	exporter sdktrace.SpanExporter
}

func (se exporterOption) apply(opts *options) {
	opts.exporter = se.exporter
}

//----------------------BEGIN LOGGING RELATED FUNCTIONS----------------------

// NewLogger creates a logger for id that produces logs at and below the indicated level.
// level here indicates the level at and below which logs are created.
func NewLogger(id string, level wcl.LogLevel) PrintlnLogger {
	return PrintlnLogger{id: id, level: level}
}

type PrintlnLogger struct {
	level wcl.LogLevel
	id    string
}

func (w PrintlnLogger) Log(level wcl.LogLevel, message string, m map[string]interface{}) {

	if level >= level {
		m["id"] = w.id
		jsonLog, err := json.Marshal(m)
		if err != nil {
			jsonLog, err = json.Marshal(map[string]interface{}{"Error Marshalling Log": err})
		}
		fmt.Println(string(jsonLog))
	}
}

//----------------------END LOGGING RELATED FUNCTIONS------------------------

//----------------------BEGIN TRACING RELATED FUNCTIONS----------------------

// newResource returns a resource describing this application.
// Used for setting up tracer provider
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("nil-test"),
			semconv.ServiceVersionKey.String("9.0.0"),
		),
	)
	return r
}

func newTracerProvider(exporter sdktrace.SpanExporter, samplingProbability float64) *sdktrace.TracerProvider {

	tp := sdktrace.NewTracerProvider(
		// We use the default ID generator. In order for sampling to work (at least with this sampler)
		// the ID generator must generate trace IDs uniformly at random from the entire space of uint64.
		// For example, the default x-ray ID generator does not do this.
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// These maximums are to guard against something going wrong and sending a ton of data unexpectedly
		sdktrace.WithSpanLimits(sdktrace.SpanLimits{
			AttributeCountLimit: 100,
			EventCountLimit:     100,
			LinkCountLimit:      100,
		}),
		//Batcher is more efficient, switch to it after testing
		// sdktrace.WithSyncer(exporter),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func doNothing(baseTransport http.RoundTripper, spanNameCtxValue interface{}, tp sdktrace.TracerProvider) http.RoundTripper {
	return baseTransport
}

func determineSampling() (samplingProbability float64, err error) {

	// If we're running locally, then turn off sampling. Otherwise sample
	// 1% or whatever TRACING_SAMPLING_PROBABILITY specifies.
	samplingProbability = 0.01
	isLocal := os.Getenv("_IS_LOCAL") == "true"
	if isLocal {
		fmt.Println("Set to Local")
		samplingProbability = 1.0
	} else if v := os.Getenv("TRACING_SAMPLING_PROBABILITY"); v != "" {
		samplingProbabilityFromEnv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("could not parse '%s' to float", v)
		}
		samplingProbability = samplingProbabilityFromEnv
	}
	return
}

//----------------------END TRACING RELATEDFUNCTIONS----------------------

// New creates a new client. The base path and http transport are configurable.
func New(ctx context.Context, basePath string, opts ...Option) *WagClient {

	defaultTransport := http.DefaultTransport
	defaultLogger := NewLogger("nil-test-wagclient", wcl.Info)
	defaultExporter := tracetest.NewNoopExporter()
	defaultInstrumentor := doNothing

	basePath = strings.TrimSuffix(basePath, "/")
	base := baseDoer{}
	// For the short-term don't use the default retry policy since its 5 retries can 5X
	// the traffic. Once we've enabled circuit breakers by default we can turn it on.
	retry := retryDoer{d: base, retryPolicy: SingleRetryPolicy{}}
	options := options{
		transport:    defaultTransport,
		logger:       defaultLogger,
		exporter:     defaultExporter,
		instrumentor: defaultInstrumentor,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	samplingProbability := 1.0 // Hard setting this to one for now, because right now
	// it is essentially ignored as the sidecar is determining the sample rate it forwards on to DD.
	// Thus the prefered approach is to sample locally with the sidecar.

	tp := newTracerProvider(options.exporter, samplingProbability)
	options.transport = options.instrumentor(options.transport, ctx, *tp)

	circuit := &circuitBreakerDoer{
		d: &retry,
		// TODO: INFRANG-4404 allow passing circuitBreakerOptions
		debug: true,
		// one circuit for each service + url pair
		circuitName: fmt.Sprintf("nil-test-%s", shortHash(basePath)),
		logger:      options.logger,
	}
	circuit.init()
	client := &WagClient{
		basePath:    basePath,
		requestDoer: circuit,
		client: &http.Client{
			Transport: options.transport,
		},
		retryDoer:      &retry,
		circuitDoer:    circuit,
		defaultTimeout: 5 * time.Second,
		logger:         options.logger,
	}
	client.SetCircuitBreakerSettings(DefaultCircuitBreakerSettings)
	return client
}

// NewFromDiscovery creates a client from the discovery environment variables. This method requires
// the three env vars: SERVICE_NIL_TEST_HTTP_(HOST/PORT/PROTO) to be set. Otherwise it returns an error.
func NewFromDiscovery(opts ...Option) (*WagClient, error) {
	url, err := discovery.URL("nil-test", "default")
	if err != nil {
		url, err = discovery.URL("nil-test", "http") // Added fallback to maintain reverse compatibility
		if err != nil {
			return nil, err
		}
	}
	return New(context.Background(), url, opts...), nil
}

// SetRetryPolicy sets a the given retry policy for all requests.
func (c *WagClient) SetRetryPolicy(retryPolicy RetryPolicy) {
	c.retryDoer.retryPolicy = retryPolicy
}

// SetCircuitBreakerDebug puts the circuit
func (c *WagClient) SetCircuitBreakerDebug(b bool) {
	c.circuitDoer.debug = b
}

// SetLogger allows for setting a custom logger
func (c *WagClient) SetLogger(l wcl.WagClientLogger) {
	c.logger = l
	c.circuitDoer.logger = l
}

// CircuitBreakerSettings are the parameters that govern the client's circuit breaker.
type CircuitBreakerSettings struct {
	// MaxConcurrentRequests is the maximum number of concurrent requests
	// the client can make at the same time. Default: 100.
	MaxConcurrentRequests int
	// RequestVolumeThreshold is the minimum number of requests needed
	// before a circuit can be tripped due to health. Default: 20.
	RequestVolumeThreshold int
	// SleepWindow how long, in milliseconds, to wait after a circuit opens
	// before testing for recovery. Default: 5000.
	SleepWindow int
	// ErrorPercentThreshold is the threshold to place on the rolling error
	// rate. Once the error rate exceeds this percentage, the circuit opens.
	// Default: 90.
	ErrorPercentThreshold int
}

// DefaultCircuitBreakerSettings describes the default circuit parameters.
var DefaultCircuitBreakerSettings = CircuitBreakerSettings{
	MaxConcurrentRequests:  100,
	RequestVolumeThreshold: 20,
	SleepWindow:            5000,
	ErrorPercentThreshold:  90,
}

// SetCircuitBreakerSettings sets parameters on the circuit breaker. It must be
// called on application startup.
func (c *WagClient) SetCircuitBreakerSettings(settings CircuitBreakerSettings) {
	hystrix.ConfigureCommand(c.circuitDoer.circuitName, hystrix.CommandConfig{
		// redundant, with the timeout we set on the context, so set
		// this to something high and irrelevant
		Timeout:                100 * 1000,
		MaxConcurrentRequests:  settings.MaxConcurrentRequests,
		RequestVolumeThreshold: settings.RequestVolumeThreshold,
		SleepWindow:            settings.SleepWindow,
		ErrorPercentThreshold:  settings.ErrorPercentThreshold,
	})
}

// SetTimeout sets a timeout on all operations for the client. To make a single request with a shorter timeout
// than the default on the client, use context.WithTimeout as described here: https://godoc.org/golang.org/x/net/context#WithTimeout.
func (c *WagClient) SetTimeout(timeout time.Duration) {
	c.defaultTimeout = timeout
}

// NilCheck makes a POST request to /check/{id}
// Nil check tests
// 200: nil
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) NilCheck(ctx context.Context, i *models.NilCheckInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	headers["header"] = i.Header

	if i.Body != nil {

		var err error
		body, err = json.Marshal(i.Body)

		if err != nil {
			return err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doNilCheckRequest(ctx, req, headers)
}

func (c *WagClient) doNilCheckRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "nilCheck")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "nilCheck")
	req = req.WithContext(ctx)
	// Don't add the timeout in a "doer" because we don't want to call "defer.cancel()"
	// until we've finished all the processing of the request object. Otherwise we'll cancel
	// our own request before we've finished it.
	if c.defaultTimeout != 0 {
		ctx, cancel := context.WithTimeout(req.Context(), c.defaultTimeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	resp, err := c.requestDoer.Do(c.client, req)
	retCode := 0
	if resp != nil {
		retCode = resp.StatusCode
	}

	// log all client failures and non-successful HT
	logData := map[string]interface{}{
		"backend":     "nil-test",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.Log(wcl.Error, "client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.Log(wcl.Error, "client-request-finished", logData)
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {

	case 200:

		return nil

	case 400:

		var output models.BadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return err
		}
		return &output

	case 500:

		var output models.InternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return err
		}
		return &output

	default:
		bs, _ := ioutil.ReadAll(resp.Body)
		return models.UnknownResponse{StatusCode: int64(resp.StatusCode), Body: string(bs)}
	}
}

func shortHash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))[0:6]
}
