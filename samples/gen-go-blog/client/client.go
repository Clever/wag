package client

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
	"github.com/Clever/wag/v6/samples/gen-go-blog/models"
	"github.com/afex/hystrix-go/hystrix"
	logger "gopkg.in/Clever/kayvee-go.v6/logger"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// Version of the client.
const Version = "0.1.0"

// VersionHeader is sent with every request.
const VersionHeader = "X-Client-Version"

// WagClient is used to make requests to the blog service.
type WagClient struct {
	basePath    string
	requestDoer doer
	client      *http.Client
	timeout     time.Duration
	// Keep the retry doer around so that we can set the number of retries
	retryDoer *retryDoer
	// Keep the circuit doer around so that we can turn it on / off
	circuitDoer    *circuitBreakerDoer
	defaultTimeout time.Duration
	logger         logger.KayveeLogger
}

var _ Client = (*WagClient)(nil)

// New creates a new client. The base path and http transport are configurable.
func New(basePath string) *WagClient {
	basePath = strings.TrimSuffix(basePath, "/")
	base := baseDoer{}
	tracing := tracingDoer{d: base}
	// For the short-term don't use the default retry policy since its 5 retries can 5X
	// the traffic. Once we've enabled circuit breakers by default we can turn it on.
	retry := retryDoer{d: tracing, retryPolicy: SingleRetryPolicy{}}
	logger := logger.New("blog-wagclient")
	circuit := &circuitBreakerDoer{
		d:     &retry,
		debug: true,
		// one circuit for each service + url pair
		circuitName: fmt.Sprintf("blog-%s", shortHash(basePath)),
		logger:      logger,
	}
	circuit.init()
	client := &WagClient{
		basePath:       basePath,
		requestDoer:    circuit,
		client:         &http.Client{Transport: http.DefaultTransport},
		retryDoer:      &retry,
		circuitDoer:    circuit,
		defaultTimeout: 5 * time.Second,
		logger:         logger,
	}
	client.SetCircuitBreakerSettings(DefaultCircuitBreakerSettings)
	return client
}

// NewFromDiscovery creates a client from the discovery environment variables. This method requires
// the three env vars: SERVICE_BLOG_HTTP_(HOST/PORT/PROTO) to be set. Otherwise it returns an error.
func NewFromDiscovery() (*WagClient, error) {
	url, err := discovery.URL("blog", "default")
	if err != nil {
		url, err = discovery.URL("blog", "http") // Added fallback to maintain reverse compatibility
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
func (c *WagClient) SetLogger(logger logger.KayveeLogger) {
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

// SetTimeout sets a timeout on all operations for the client. To make a single request
// with a timeout use context.WithTimeout as described here: https://godoc.org/golang.org/x/net/context#WithTimeout.
func (c *WagClient) SetTimeout(timeout time.Duration) {
	c.defaultTimeout = timeout
}

// SetTransport sets the http transport used by the client.
func (c *WagClient) SetTransport(t http.RoundTripper) {
	c.client.Transport = t
}

// PostGradeFileForStudent makes a POST request to /students/{student_id}/gradeFile
// Posts the grade file for the specified student
// 200: nil
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) PostGradeFileForStudent(ctx context.Context, i *models.PostGradeFileForStudentInput) error {
	headers := make(map[string]string)

	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "POST", path, *i.File)

	if err != nil {
		return err
	}

	return c.doPostGradeFileForStudentRequest(ctx, req, headers)
}

func (c *WagClient) doPostGradeFileForStudentRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "postGradeFileForStudent")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "postGradeFileForStudent")
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
	logData := logger.M{
		"backend":     "blog",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.ErrorD("client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.ErrorD("client-request-finished", logData)
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

// GetSectionsForStudent makes a GET request to /students/{student_id}/sections
// Gets the sections for the specified student
// 200: []models.Section
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetSectionsForStudent(ctx context.Context, studentID string) ([]models.Section, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetSectionsForStudentInputPath(studentID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetSectionsForStudentRequest(ctx, req, headers)
}

func (c *WagClient) doGetSectionsForStudentRequest(ctx context.Context, req *http.Request, headers map[string]string) ([]models.Section, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getSectionsForStudent")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getSectionsForStudent")
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
	logData := logger.M{
		"backend":     "blog",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.ErrorD("client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.ErrorD("client-request-finished", logData)
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {

	case 200:

		var output []models.Section
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}

		return output, nil

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

// PostSectionsForStudent makes a POST request to /students/{student_id}/sections
// Posts the sections for the specified student
// 200: []models.Section
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) PostSectionsForStudent(ctx context.Context, i *models.PostSectionsForStudentInput) ([]models.Section, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doPostSectionsForStudentRequest(ctx, req, headers)
}

func (c *WagClient) doPostSectionsForStudentRequest(ctx context.Context, req *http.Request, headers map[string]string) ([]models.Section, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "postSectionsForStudent")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "postSectionsForStudent")
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
	logData := logger.M{
		"backend":     "blog",
		"method":      req.Method,
		"uri":         req.URL,
		"status_code": retCode,
	}
	if err == nil && retCode > 399 {
		logData["message"] = resp.Status
		c.logger.ErrorD("client-request-finished", logData)
	}
	if err != nil {
		logData["message"] = err.Error()
		c.logger.ErrorD("client-request-finished", logData)
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {

	case 200:

		var output []models.Section
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}

		return output, nil

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

func shortHash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))[0:6]
}
