package client

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	discovery "github.com/Clever/discovery-go"
	"github.com/Clever/wag/samples/gen-go-errors/models"
	"github.com/afex/hystrix-go/hystrix"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// WagClient is used to make requests to the swagger-test service.
type WagClient struct {
	basePath    string
	requestDoer doer
	transport   *http.Transport
	timeout     time.Duration
	// Keep the retry doer around so that we can set the number of retries
	retryDoer *retryDoer
	// Keep the circuit doer around so that we can turn it on / off
	circuitDoer    *circuitBreakerDoer
	defaultTimeout time.Duration
}

var _ Client = (*WagClient)(nil)

// New creates a new client. The base path and http transport are configurable.
func New(basePath string) *WagClient {
	base := baseDoer{}
	tracing := tracingDoer{d: base}
	retry := retryDoer{d: tracing, retryPolicy: DefaultRetryPolicy{}}
	circuit := &circuitBreakerDoer{
		d:     &retry,
		debug: true,
		// one circuit for each service + url pair
		circuitName: fmt.Sprintf("swagger-test-%s", shortHash(basePath)),
	}
	circuit.init()
	client := &WagClient{requestDoer: circuit, retryDoer: &retry, circuitDoer: circuit, defaultTimeout: 10 * time.Second,
		transport: &http.Transport{}, basePath: basePath}
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

// WithRetryPolicy returns a new client that will use the given retry policy for
// all requests.
func (c *WagClient) WithRetryPolicy(retryPolicy RetryPolicy) *WagClient {
	c.retryDoer.retryPolicy = retryPolicy
	return c
}

// SetCircuitBreakerDebug puts the circuit
func (c *WagClient) SetCircuitBreakerDebug(b bool) {
	c.circuitDoer.debug = b
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

// WithTimeout returns a new client that has the specified timeout on all operations. To make a single request
// have a timeout use context.WithTimeout as described here: https://godoc.org/golang.org/x/net/context#WithTimeout.
func (c *WagClient) WithTimeout(timeout time.Duration) *WagClient {
	c.defaultTimeout = timeout
	return c
}

// JoinByFormat joins a string array by a known format:
//	 csv: comma separated value (default)
//	 ssv: space separated value
//	 tsv: tab separated value
//	 pipes: pipe (|) separated value
func JoinByFormat(data []string, format string) string {
	if len(data) == 0 {
		return ""
	}
	var sep string
	switch format {
	case "ssv":
		sep = " "
	case "tsv":
		sep = "\t"
	case "pipes":
		sep = "|"
	default:
		sep = ","
	}
	return strings.Join(data, sep)
}

// GetBook makes a GET request to /books/{id}.
func (c *WagClient) GetBook(ctx context.Context, i *models.GetBookInput) error {
	path := c.basePath + "/v1/books/{id}"
	urlVals := url.Values{}
	var body []byte

	path = strings.Replace(path, "{id}", strconv.FormatInt(i.ID, 10), -1)
	path = path + "?" + urlVals.Encode()

	client := &http.Client{Transport: c.transport}
	req, err := http.NewRequest("GET", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getBook")
	req = req.WithContext(ctx)
	// Don't add the timeout in a "doer" because we don't want to call "defer.cancel()"
	// until we've finished all the processing of the request object. Otherwise we'll cancel
	// our own request before we've finished it.
	if c.defaultTimeout != 0 {
		ctx, cancel := context.WithTimeout(req.Context(), c.defaultTimeout)
		defer cancel()
		req = req.WithContext(ctx)
	}
	resp, err := c.requestDoer.Do(client, req)

	if err != nil {
		return &models.InternalError{Msg: err.Error()}
	}

	defer resp.Body.Close()
	switch resp.StatusCode {

	case 200:

		return nil

	case 400:

		var output models.ExtendedError
		// Any errors other than EOF should result in an error. EOF is acceptable for empty
		// types.
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil && err != io.EOF {
			return &models.InternalError{Msg: err.Error()}
		}
		return &output

	case 404:

		var output models.NotFoundError
		// Any errors other than EOF should result in an error. EOF is acceptable for empty
		// types.
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil && err != io.EOF {
			return &models.InternalError{Msg: err.Error()}
		}
		return &output

	case 500:

		var output models.InternalError
		// Any errors other than EOF should result in an error. EOF is acceptable for empty
		// types.
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil && err != io.EOF {
			return &models.InternalError{Msg: err.Error()}
		}
		return &output

	default:
		return &models.InternalError{Msg: "Unknown response"}
	}
}

func shortHash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))[0:6]
}
