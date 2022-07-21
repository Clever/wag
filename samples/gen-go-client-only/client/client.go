package client

// Using Alpha version of WAG Yay!

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Clever/wag/loggers/printlogger"
	waglogger "github.com/Clever/wag/loggers/waglogger"

	"github.com/Clever/wag/samples/v8/gen-go-client-only/models"

	discovery "github.com/Clever/discovery-go"

	"github.com/afex/hystrix-go/hystrix"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// Version of the client.
const Version = "0.1.0"

// VersionHeader is sent with every request.
const VersionHeader = "X-Client-Version"

// WagClient is used to make requests to the swagger-test service.
type WagClient struct {
	basePath    string
	requestDoer doer
	client      *http.Client
	// Keep the retry doer around so that we can set the number of retries
	retryDoer *retryDoer
	// Keep the circuit doer around so that we can turn it on / off
	circuitDoer    *circuitBreakerDoer
	defaultTimeout time.Duration
	logger         waglogger.WagClientLogger
}

var _ Client = (*WagClient)(nil)

//This pattern is used instead of using closures for greater transparency and the ability to implement additional interfaces.
type options struct {
	transport    http.RoundTripper
	logger       waglogger.WagClientLogger
	instrumentor Instrumentor
	exporter     sdktrace.SpanExporter
}

type Option interface {
	apply(*options)
}

//Logger

//WithLogger sets client logger option.
func WithLogger(log waglogger.WagClientLogger) Option {
	return loggerOption{Log: log}
}

type loggerOption struct {
	Log waglogger.WagClientLogger
}

func (l loggerOption) apply(opts *options) {
	opts.logger = l.Log
}

//RoundTripper

type roundTripperOption struct {
	rt http.RoundTripper
}

func (t roundTripperOption) apply(opts *options) {
	opts.transport = t.rt
}

//WithRoundTripper allows you to pass in intrumented/custom roundtrippers which will then wrap the
//transport roundtripper
func WithRoundTripper(t http.RoundTripper) Option {
	return roundTripperOption{rt: t}
}

//Instrumentor is a function that creates an instrumented round tripper
type Instrumentor func(baseTransport http.RoundTripper, spanNameCtxValue interface{}, tp sdktrace.TracerProvider) http.RoundTripper

//WithInstrumentor sets a instrumenting function that will be used to wrap the roundTripper for tracing.
func WithInstrumentor(fn Instrumentor) Option {
	return instrumentorOption{instrumentor: fn}
}

type instrumentorOption struct {
	instrumentor Instrumentor
}

func (i instrumentorOption) apply(opts *options) {
	opts.instrumentor = i.instrumentor
}

//WithExporter sets client span exporter option.
func WithExporter(se sdktrace.SpanExporter) Option {
	return exporterOption{exporter: se}
}

type exporterOption struct {
	exporter sdktrace.SpanExporter
}

func (se exporterOption) apply(opts *options) {
	opts.exporter = se.exporter
}

//----------------------BEGIN TRACING RELATED FUNCTIONS----------------------

// newResource returns a resource describing this application.
// Used for setting up tracer provider
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("dapple"),
			semconv.ServiceVersionKey.String("1.11.0"),
		),
	)
	return r
}

func newTracerProvider(exporter sdktrace.SpanExporter, samplingProbability float64) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		// We use the default ID generator. In order for sampling to work (at least with this sampler)
		// the ID generator must generate trace IDs uniformly at random from the entire space of uint64.
		// For example, the default x-ray ID generator does not do this.
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(samplingProbability))),
		// These maximums are to guard against something going wrong and sending a ton of data unexpectedly
		sdktrace.WithSpanLimits(sdktrace.SpanLimits{
			AttributeCountLimit: 100,
			EventCountLimit:     100,
			LinkCountLimit:      100,
		}),
		//Batcher is more efficient, switch to it after testing
		sdktrace.WithSyncer(exporter),
		//sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource()),
	)
}
func doNothing(baseTransport http.RoundTripper, spanNameCtxValue interface{}, tp sdktrace.TracerProvider) http.RoundTripper {
	return baseTransport
}
func determineSampling() (samplingProbability float64, err error) {

	// 	// If we're running locally, then turn off sampling. Otherwise sample
	// 	// 1%!o(MISSING)r whatever TRACING_SAMPLING_PROBABILITY specifies.
	samplingProbability = 0.01
	isLocal := os.Getenv("_IS_LOCAL") == "true"
	if isLocal {
		fmt.Println("Set to Local")
		samplingProbability = 1.0
	} else if v := os.Getenv("TRACING_SAMPLING_PROBABILITY"); v != "" {
		samplingProbabilityFromEnv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("could not parse '%!s(MISSING)' to float", v)
		}
		samplingProbability = samplingProbabilityFromEnv
	}
	return
}

//----------------------END TRACING RELATEDFUNCTIONS----------------------

// New creates a new client. The base path and http transport are configurable.
func New(basePath string, opts ...Option) *WagClient {

	defaultTransport := http.DefaultTransport
	defaultLogger := printlogger.NewLogger("swagger-test-wagclient", "info")
	defaultExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		fmt.Println(err)
	}
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

	samplingProbability := 1.0 // TODO: Put back logic to set this to 1 for local, 0.1 otherwise etc.
	// samplingProbability := determineSampling()

	tp := newTracerProvider(options.exporter, samplingProbability)
	options.transport = options.instrumentor(options.transport, context.TODO(), *tp)

	circuit := &circuitBreakerDoer{
		d: &retry,
		// TODO: INFRANG-4404 allow passing circuitBreakerOptions
		debug: true,
		// one circuit for each service + url pair
		circuitName: fmt.Sprintf("swagger-test-%s", shortHash(basePath)),
		logger:      options.logger,
	}
	circuit.init()
	client := &WagClient{
		basePath:    basePath,
		requestDoer: circuit,
		client: &http.Client{
			Transport: defaultTransport,
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
// the three env vars: SERVICE_SWAGGER_TEST_HTTP_(HOST/PORT/PROTO) to be set. Otherwise it returns an error.
func NewFromDiscovery() (*WagClient, error) {
	url, err := discovery.URL("swagger-test", "default")
	if err != nil {
		url, err = discovery.URL("swagger-test", "http") // Added fallback to maintain reverse compatibility
		if err != nil {
			return nil, err
		}
	}
	return New(url), nil
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
func (c *WagClient) SetLogger(l waglogger.WagClientLogger) {
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

// SetTransport sets the http transport used by the client.
func (c *WagClient) SetTransport(t http.RoundTripper) {
	// c.client.Transport = tracing.NewTransport(t, opNameCtx{})
}

// GetAuthors makes a GET request to /authors
// Gets authors
// 200: *models.AuthorsResponse
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAuthors(ctx context.Context, i *models.GetAuthorsInput) (*models.AuthorsResponse, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	resp, _, err := c.doGetAuthorsRequest(ctx, req, headers)
	return resp, err
}

type getAuthorsIterImpl struct {
	c            *WagClient
	ctx          context.Context
	lastResponse []*models.Author
	index        int
	err          error
	nextURL      string
	headers      map[string]string
	body         []byte
}

// NewgetAuthorsIter constructs an iterator that makes calls to getAuthors for
// each page.
func (c *WagClient) NewGetAuthorsIter(ctx context.Context, i *models.GetAuthorsInput) (GetAuthorsIter, error) {
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	headers := make(map[string]string)

	var body []byte

	return &getAuthorsIterImpl{
		c:            c,
		ctx:          ctx,
		lastResponse: []*models.Author{},
		nextURL:      path,
		headers:      headers,
		body:         body,
	}, nil
}

func (i *getAuthorsIterImpl) refresh() error {
	req, err := http.NewRequestWithContext(i.ctx, "GET", i.nextURL, bytes.NewBuffer(i.body))

	if err != nil {
		i.err = err
		return err
	}

	resp, nextPage, err := i.c.doGetAuthorsRequest(i.ctx, req, i.headers)
	if err != nil {
		i.err = err
		return err
	}

	i.lastResponse = resp.AuthorSet.Results
	i.index = 0
	if nextPage != "" {
		i.nextURL = i.c.basePath + nextPage
	} else {
		i.nextURL = ""
	}
	return nil
}

// Next retrieves the next resource from the iterator and assigns it to the
// provided pointer, fetching a new page if necessary. Returns true if it
// successfully retrieves a new resource.
func (i *getAuthorsIterImpl) Next(v *models.Author) bool {
	if i.err != nil {
		return false
	} else if i.index < len(i.lastResponse) {
		*v = *i.lastResponse[i.index]
		i.index++
		return true
	} else if i.nextURL == "" {
		return false
	}

	if err := i.refresh(); err != nil {
		return false
	}
	return i.Next(v)
}

// Err returns an error if one occurred when .Next was called.
func (i *getAuthorsIterImpl) Err() error {
	return i.err
}

func (c *WagClient) doGetAuthorsRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.AuthorsResponse, string, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAuthors")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAuthors")
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
		"backend":     "swagger-test",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.Log("error", "client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.Log("error", "client-request-finished", logData)
		return nil, "", err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {

	case 200:

		var output models.AuthorsResponse
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, "", err
		}

		return &output, resp.Header.Get("X-Next-Page-Path"), nil

	case 400:

		var output models.BadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, "", err
		}
		return nil, "", &output

	case 500:

		var output models.InternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, "", err
		}
		return nil, "", &output

	default:
		return nil, "", &models.InternalError{Message: fmt.Sprintf("Unknown status code %v", resp.StatusCode)}
	}
}

// GetAuthorsWithPut makes a PUT request to /authors
// Gets authors, but needs to use the body so it's a PUT
// 200: *models.AuthorsResponse
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAuthorsWithPut(ctx context.Context, i *models.GetAuthorsWithPutInput) (*models.AuthorsResponse, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.FavoriteBooks != nil {

		var err error
		body, err = json.Marshal(i.FavoriteBooks)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	resp, _, err := c.doGetAuthorsWithPutRequest(ctx, req, headers)
	return resp, err
}

type getAuthorsWithPutIterImpl struct {
	c            *WagClient
	ctx          context.Context
	lastResponse []*models.Author
	index        int
	err          error
	nextURL      string
	headers      map[string]string
	body         []byte
}

// NewgetAuthorsWithPutIter constructs an iterator that makes calls to getAuthorsWithPut for
// each page.
func (c *WagClient) NewGetAuthorsWithPutIter(ctx context.Context, i *models.GetAuthorsWithPutInput) (GetAuthorsWithPutIter, error) {
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	headers := make(map[string]string)

	var body []byte

	if i.FavoriteBooks != nil {

		var err error
		body, err = json.Marshal(i.FavoriteBooks)

		if err != nil {
			return nil, err
		}

	}

	return &getAuthorsWithPutIterImpl{
		c:            c,
		ctx:          ctx,
		lastResponse: []*models.Author{},
		nextURL:      path,
		headers:      headers,
		body:         body,
	}, nil
}

func (i *getAuthorsWithPutIterImpl) refresh() error {
	req, err := http.NewRequestWithContext(i.ctx, "PUT", i.nextURL, bytes.NewBuffer(i.body))

	if err != nil {
		i.err = err
		return err
	}

	resp, nextPage, err := i.c.doGetAuthorsWithPutRequest(i.ctx, req, i.headers)
	if err != nil {
		i.err = err
		return err
	}

	i.lastResponse = resp.AuthorSet.Results
	i.index = 0
	if nextPage != "" {
		i.nextURL = i.c.basePath + nextPage
	} else {
		i.nextURL = ""
	}
	return nil
}

// Next retrieves the next resource from the iterator and assigns it to the
// provided pointer, fetching a new page if necessary. Returns true if it
// successfully retrieves a new resource.
func (i *getAuthorsWithPutIterImpl) Next(v *models.Author) bool {
	if i.err != nil {
		return false
	} else if i.index < len(i.lastResponse) {
		*v = *i.lastResponse[i.index]
		i.index++
		return true
	} else if i.nextURL == "" {
		return false
	}

	if err := i.refresh(); err != nil {
		return false
	}
	return i.Next(v)
}

// Err returns an error if one occurred when .Next was called.
func (i *getAuthorsWithPutIterImpl) Err() error {
	return i.err
}

func (c *WagClient) doGetAuthorsWithPutRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.AuthorsResponse, string, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAuthorsWithPut")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAuthorsWithPut")
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
		"backend":     "swagger-test",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.Log("error", "client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.Log("error", "client-request-finished", logData)
		return nil, "", err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {

	case 200:

		var output models.AuthorsResponse
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, "", err
		}

		return &output, resp.Header.Get("X-Next-Page-Path"), nil

	case 400:

		var output models.BadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, "", err
		}
		return nil, "", &output

	case 500:

		var output models.InternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, "", err
		}
		return nil, "", &output

	default:
		return nil, "", &models.InternalError{Message: fmt.Sprintf("Unknown status code %v", resp.StatusCode)}
	}
}

// GetBooks makes a GET request to /books
// Returns a list of books
// 200: []models.Book
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetBooks(ctx context.Context, i *models.GetBooksInput) ([]models.Book, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	headers["authorization"] = i.Authorization

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	resp, _, err := c.doGetBooksRequest(ctx, req, headers)
	return resp, err
}

type getBooksIterImpl struct {
	c            *WagClient
	ctx          context.Context
	lastResponse []models.Book
	index        int
	err          error
	nextURL      string
	headers      map[string]string
	body         []byte
}

// NewgetBooksIter constructs an iterator that makes calls to getBooks for
// each page.
func (c *WagClient) NewGetBooksIter(ctx context.Context, i *models.GetBooksInput) (GetBooksIter, error) {
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	headers := make(map[string]string)

	headers["authorization"] = i.Authorization

	var body []byte

	return &getBooksIterImpl{
		c:            c,
		ctx:          ctx,
		lastResponse: []models.Book{},
		nextURL:      path,
		headers:      headers,
		body:         body,
	}, nil
}

func (i *getBooksIterImpl) refresh() error {
	req, err := http.NewRequestWithContext(i.ctx, "GET", i.nextURL, bytes.NewBuffer(i.body))

	if err != nil {
		i.err = err
		return err
	}

	resp, nextPage, err := i.c.doGetBooksRequest(i.ctx, req, i.headers)
	if err != nil {
		i.err = err
		return err
	}

	i.lastResponse = resp
	i.index = 0
	if nextPage != "" {
		i.nextURL = i.c.basePath + nextPage
	} else {
		i.nextURL = ""
	}
	return nil
}

// Next retrieves the next resource from the iterator and assigns it to the
// provided pointer, fetching a new page if necessary. Returns true if it
// successfully retrieves a new resource.
func (i *getBooksIterImpl) Next(v *models.Book) bool {
	if i.err != nil {
		return false
	} else if i.index < len(i.lastResponse) {
		*v = i.lastResponse[i.index]
		i.index++
		return true
	} else if i.nextURL == "" {
		return false
	}

	if err := i.refresh(); err != nil {
		return false
	}
	return i.Next(v)
}

// Err returns an error if one occurred when .Next was called.
func (i *getBooksIterImpl) Err() error {
	return i.err
}

func (c *WagClient) doGetBooksRequest(ctx context.Context, req *http.Request, headers map[string]string) ([]models.Book, string, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getBooks")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getBooks")
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
		"backend":     "swagger-test",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.Log("error", "client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.Log("error", "client-request-finished", logData)
		return nil, "", err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {

	case 200:

		var output []models.Book
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, "", err
		}

		return output, resp.Header.Get("X-Next-Page-Path"), nil

	case 400:

		var output models.BadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, "", err
		}
		return nil, "", &output

	case 500:

		var output models.InternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, "", err
		}
		return nil, "", &output

	default:
		return nil, "", &models.InternalError{Message: fmt.Sprintf("Unknown status code %v", resp.StatusCode)}
	}
}

// CreateBook makes a POST request to /books
// Creates a book
// 200: *models.Book
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) CreateBook(ctx context.Context, i *models.Book) (*models.Book, error) {
	headers := make(map[string]string)

	var body []byte
	path := c.basePath + "/v1/books"

	if i != nil {

		var err error
		body, err = json.Marshal(i)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doCreateBookRequest(ctx, req, headers)
}

func (c *WagClient) doCreateBookRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Book, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "createBook")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "createBook")
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
		"backend":     "swagger-test",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.Log("error", "client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.Log("error", "client-request-finished", logData)
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {

	case 200:

		var output models.Book
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}

		return &output, nil

	case 400:

		var output models.BadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	case 500:

		var output models.InternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	default:
		return nil, &models.InternalError{Message: fmt.Sprintf("Unknown status code %v", resp.StatusCode)}
	}
}

// PutBook makes a PUT request to /books
// Puts a book
// 200: *models.Book
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) PutBook(ctx context.Context, i *models.Book) (*models.Book, error) {
	headers := make(map[string]string)

	var body []byte
	path := c.basePath + "/v1/books"

	if i != nil {

		var err error
		body, err = json.Marshal(i)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doPutBookRequest(ctx, req, headers)
}

func (c *WagClient) doPutBookRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Book, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "putBook")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "putBook")
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
		"backend":     "swagger-test",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.Log("error", "client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.Log("error", "client-request-finished", logData)
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {

	case 200:

		var output models.Book
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}

		return &output, nil

	case 400:

		var output models.BadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	case 500:

		var output models.InternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	default:
		return nil, &models.InternalError{Message: fmt.Sprintf("Unknown status code %v", resp.StatusCode)}
	}
}

// GetBookByID makes a GET request to /books/{book_id}
// Returns a book
// 200: *models.Book
// 400: *models.BadRequest
// 401: *models.Unathorized
// 404: *models.Error
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetBookByID(ctx context.Context, i *models.GetBookByIDInput) (*models.Book, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	headers["authorization"] = i.Authorization

	headers["X-Dont-Rate-Limit-Me-Bro"] = i.XDontRateLimitMeBro

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetBookByIDRequest(ctx, req, headers)
}

func (c *WagClient) doGetBookByIDRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Book, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getBookByID")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getBookByID")
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
		"backend":     "swagger-test",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.Log("error", "client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.Log("error", "client-request-finished", logData)
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {

	case 200:

		var output models.Book
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}

		return &output, nil

	case 400:

		var output models.BadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	case 401:

		var output models.Unathorized
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	case 404:

		var output models.Error
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	case 500:

		var output models.InternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	default:
		return nil, &models.InternalError{Message: fmt.Sprintf("Unknown status code %v", resp.StatusCode)}
	}
}

// GetBookByID2 makes a GET request to /books2/{id}
// Retrieve a book
// 200: *models.Book
// 400: *models.BadRequest
// 404: *models.Error
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetBookByID2(ctx context.Context, id string) (*models.Book, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetBookByID2InputPath(id)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetBookByID2Request(ctx, req, headers)
}

func (c *WagClient) doGetBookByID2Request(ctx context.Context, req *http.Request, headers map[string]string) (*models.Book, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getBookByID2")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getBookByID2")
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
		"backend":     "swagger-test",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.Log("error", "client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.Log("error", "client-request-finished", logData)
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {

	case 200:

		var output models.Book
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}

		return &output, nil

	case 400:

		var output models.BadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	case 404:

		var output models.Error
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	case 500:

		var output models.InternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	default:
		return nil, &models.InternalError{Message: fmt.Sprintf("Unknown status code %v", resp.StatusCode)}
	}
}

// HealthCheck makes a GET request to /health/check
//
// 200: nil
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) HealthCheck(ctx context.Context) error {
	headers := make(map[string]string)

	var body []byte
	path := c.basePath + "/v1/health/check"

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doHealthCheckRequest(ctx, req, headers)
}

func (c *WagClient) doHealthCheckRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "healthCheck")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "healthCheck")
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
		"backend":     "swagger-test",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.Log("error", "client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.Log("error", "client-request-finished", logData)
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
		return &models.InternalError{Message: fmt.Sprintf("Unknown status code %v", resp.StatusCode)}
	}
}

func shortHash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))[0:6]
}
