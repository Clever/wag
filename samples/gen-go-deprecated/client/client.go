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

	discovery "github.com/Clever/discovery-go"
	"github.com/Clever/wag/samples/v8/gen-go-deprecated/models"
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

func NewLogger(id string, level int) WagClientLogger {
	return WagClientPrintlnLogger{id: id, level: level}
}

type WagClientPrintlnLogger struct {
	level int
	id    string
}

func (w WagClientPrintlnLogger) Log(level int, message string, m map[string]interface{}) {
	if w.level >= level {
		fmt.Print(w.id, ": ")
		fmt.Print(message)
		for k, v := range m {
			fmt.Print(" ", k, " : ", v)
		}
		fmt.Println()
	}
}

type WagClientLogger interface {
	Log(level int, message string, pairs map[string]interface{})
}

var CRITICALD int = 0
var ERRORD int = 1
var WARND int = 2
var INFOD int = 3
var DEBUGD int = 4

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
	logger         WagClientLogger
}

var _ Client = (*WagClient)(nil)

// New creates a new client. The base path and http transport are configurable.
func New(basePath string) *WagClient {
	basePath = strings.TrimSuffix(basePath, "/")
	base := baseDoer{}
	// For the short-term don't use the default retry policy since its 5 retries can 5X
	// the traffic. Once we've enabled circuit breakers by default we can turn it on.
	retry := retryDoer{d: base, retryPolicy: SingleRetryPolicy{}}
	logger := NewLogger("swagger-test-wagclient", 3)
	circuit := &circuitBreakerDoer{
		d: &retry,
		// TODO: INFRANG-4404 allow passing circuitBreakerOptions
		debug: true,
		// one circuit for each service + url pair
		circuitName: fmt.Sprintf("swagger-test-%s", shortHash(basePath)),
		logger:      logger,
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
		logger:         logger,
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
func (c *WagClient) SetLogger(logger WagClientLogger) {
	c.logger = logger
	c.circuitDoer.logger = logger
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

func shortHash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))[0:6]
}
