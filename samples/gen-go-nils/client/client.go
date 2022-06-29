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
	"github.com/Clever/wag/loggers/waglogger"

	discovery "github.com/Clever/discovery-go"
	"github.com/Clever/wag/samples/v8/gen-go-nils/models"
	"github.com/afex/hystrix-go/hystrix"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// Version of the client.
const Version = "0.1.0"

// VersionHeader is sent with every request.
const VersionHeader = "X-Client-Version"

var CRITICALD int = 0
var ERRORD int = 1
var WARND int = 2
var INFOD int = 3
var DEBUGD int = 4

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
	logger         waglogger.WagClientLogger
}

var _ Client = (*WagClient)(nil)

//This pattern is used instead of using closures for greater transparency and the ability to implement additional interfaces.
type options struct {
	tracing interface{}
	logger  waglogger.WagClientLogger
}

type Option interface {
	apply(*options)
}

//WithLogger sets client logger option.
func WithLogger(log waglogger.WagClientLogger) Option {
	return loggerOption{Log: log}
}

type loggerOption struct {
	Log *waglogger.WagClientLogger
}

func (l loggerOption) apply(opts *options) {
	opts.logger = l.Log
}

type tracingOption interface{}

func (t tracingOption) apply(opts *options) {
	opts.tracing = t //This is where I would normally have a : tracing.NewTransport(http.DefaultTransport, opNameCtx{})
}

//WithTracingProvider defines the tracing provider for
func WithTracingProvider(t interface{}) Option {
	return tracingOption(t)
}

// New creates a new client. The base path and http transport are configurable.
func New(basePath string) *WagClient {
	basePath = strings.TrimSuffix(basePath, "/")
	base := baseDoer{}
	// For the short-term don't use the default retry policy since its 5 retries can 5X
	// the traffic. Once we've enabled circuit breakers by default we can turn it on.
	retry := retryDoer{d: base, retryPolicy: SingleRetryPolicy{}}
	options := options{
		tracing: defaultTracing,
		logger:  printlogger.NewLogger("nil-test-wagclient", 3),
	}

	for _, o := range opts {
		o.apply(&options)
	}

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
		client:      &http.Client{
			// Transport: tracing.NewTransport(http.DefaultTransport, opNameCtx{}),
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
func NewFromDiscovery() (*WagClient, error) {
	url, err := discovery.URL("nil-test", "default")
	if err != nil {
		url, err = discovery.URL("nil-test", "http") // Added fallback to maintain reverse compatibility
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
func (c *WagClient) SetLogger(l logger.WagClientLogger) {
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
