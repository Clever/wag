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
	"github.com/Clever/wag/samples/gen-go-app-service/models"
	"github.com/afex/hystrix-go/hystrix"
	logger "gopkg.in/Clever/kayvee-go.v6/logger"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// Version of the client.
const Version = "12.15.1"

// VersionHeader is sent with every request.
const VersionHeader = "X-Client-Version"

// WagClient is used to make requests to the app-service service.
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
	logger := logger.New("app-service-wagclient")
	circuit := &circuitBreakerDoer{
		d:     &retry,
		debug: true,
		// one circuit for each service + url pair
		circuitName: fmt.Sprintf("app-service-%s", shortHash(basePath)),
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
// the three env vars: SERVICE_APP_SERVICE_HTTP_(HOST/PORT/PROTO) to be set. Otherwise it returns an error.
func NewFromDiscovery() (*WagClient, error) {
	url, err := discovery.URL("app-service", "default")
	if err != nil {
		url, err = discovery.URL("app-service", "http") // Added fallback to maintain reverse compatibility
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

// HealthCheck makes a GET request to /_health
// Checks if the service is healthy
// 200: nil
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) HealthCheck(ctx context.Context) error {
	headers := make(map[string]string)

	var body []byte
	path := c.basePath + "/_health"

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
	logData := logger.M{
		"backend":     "app-service",
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// GetAdmins makes a GET request to /v1/admins
//
// 200: []models.Admin
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAdmins(ctx context.Context, i *models.GetAdminsInput) ([]models.Admin, error) {
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

	return c.doGetAdminsRequest(ctx, req, headers)
}

func (c *WagClient) doGetAdminsRequest(ctx context.Context, req *http.Request, headers map[string]string) ([]models.Admin, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAdmins")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAdmins")
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
		"backend":     "app-service",
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

		var output []models.Admin
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// DeleteAdmin makes a DELETE request to /v1/admins/{adminID}
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) DeleteAdmin(ctx context.Context, adminID string) error {
	headers := make(map[string]string)

	var body []byte
	path, err := models.DeleteAdminInputPath(adminID)

	if err != nil {
		return err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "DELETE", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doDeleteAdminRequest(ctx, req, headers)
}

func (c *WagClient) doDeleteAdminRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "deleteAdmin")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "deleteAdmin")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// GetAdminByID makes a GET request to /v1/admins/{adminID}
//
// 200: *models.Admin
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAdminByID(ctx context.Context, adminID string) (*models.Admin, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetAdminByIDInputPath(adminID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetAdminByIDRequest(ctx, req, headers)
}

func (c *WagClient) doGetAdminByIDRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Admin, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAdminByID")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAdminByID")
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
		"backend":     "app-service",
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

		var output models.Admin
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// UpdateAdmin makes a PATCH request to /v1/admins/{adminID}
//
// 200: *models.Admin
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) UpdateAdmin(ctx context.Context, i *models.UpdateAdminInput) (*models.Admin, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.Admin != nil {

		var err error
		body, err = json.Marshal(i.Admin)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doUpdateAdminRequest(ctx, req, headers)
}

func (c *WagClient) doUpdateAdminRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Admin, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "updateAdmin")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "updateAdmin")
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
		"backend":     "app-service",
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

		var output models.Admin
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// CreateAdmin makes a PUT request to /v1/admins/{adminID}
//
// 200: *models.Admin
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) CreateAdmin(ctx context.Context, i *models.CreateAdminInput) (*models.Admin, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.CreateAdmin != nil {

		var err error
		body, err = json.Marshal(i.CreateAdmin)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doCreateAdminRequest(ctx, req, headers)
}

func (c *WagClient) doCreateAdminRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Admin, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "createAdmin")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "createAdmin")
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
		"backend":     "app-service",
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

		var output models.Admin
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// VerifyCode makes a POST request to /v1/admins/{adminID}/confirmation_code
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) VerifyCode(ctx context.Context, i *models.VerifyCodeInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doVerifyCodeRequest(ctx, req, headers)
}

func (c *WagClient) doVerifyCodeRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "verifyCode")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "verifyCode")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// CreateVerificationCode makes a PUT request to /v1/admins/{adminID}/confirmation_code
//
// 200: *models.VerificationCodeResponse
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) CreateVerificationCode(ctx context.Context, i *models.CreateVerificationCodeInput) (*models.VerificationCodeResponse, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doCreateVerificationCodeRequest(ctx, req, headers)
}

func (c *WagClient) doCreateVerificationCodeRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.VerificationCodeResponse, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "createVerificationCode")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "createVerificationCode")
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
		"backend":     "app-service",
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

		var output models.VerificationCodeResponse
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// VerifyAdminEmail makes a POST request to /v1/admins/{adminID}/verify_email
// set the verified email of an admin
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) VerifyAdminEmail(ctx context.Context, i *models.VerifyAdminEmailInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	if i.Request != nil {

		var err error
		body, err = json.Marshal(i.Request)

		if err != nil {
			return err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doVerifyAdminEmailRequest(ctx, req, headers)
}

func (c *WagClient) doVerifyAdminEmailRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "verifyAdminEmail")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "verifyAdminEmail")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// GetAllAnalyticsApps makes a GET request to /v1/analytics/apps
//
// 200: *models.AnalyticsApps
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAllAnalyticsApps(ctx context.Context) (*models.AnalyticsApps, error) {
	headers := make(map[string]string)

	var body []byte
	path := c.basePath + "/v1/analytics/apps"

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetAllAnalyticsAppsRequest(ctx, req, headers)
}

func (c *WagClient) doGetAllAnalyticsAppsRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.AnalyticsApps, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAllAnalyticsApps")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAllAnalyticsApps")
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
		"backend":     "app-service",
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

		var output models.AnalyticsApps
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetAnalyticsAppByShortname makes a GET request to /v1/analytics/apps/{shortname}
//
// 200: *models.AnalyticsApp
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAnalyticsAppByShortname(ctx context.Context, shortname string) (*models.AnalyticsApp, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetAnalyticsAppByShortnameInputPath(shortname)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetAnalyticsAppByShortnameRequest(ctx, req, headers)
}

func (c *WagClient) doGetAnalyticsAppByShortnameRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.AnalyticsApp, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAnalyticsAppByShortname")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAnalyticsAppByShortname")
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
		"backend":     "app-service",
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

		var output models.AnalyticsApp
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetAllTrackableApps makes a GET request to /v1/analytics/trackable_apps
//
// 200: *models.TrackableApps
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAllTrackableApps(ctx context.Context) (*models.TrackableApps, error) {
	headers := make(map[string]string)

	var body []byte
	path := c.basePath + "/v1/analytics/trackable_apps"

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetAllTrackableAppsRequest(ctx, req, headers)
}

func (c *WagClient) doGetAllTrackableAppsRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.TrackableApps, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAllTrackableApps")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAllTrackableApps")
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
		"backend":     "app-service",
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

		var output models.TrackableApps
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetAnalyticsUsageUrls makes a GET request to /v1/analytics/usageUrls
//
// 200: *models.UsageUrls
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAnalyticsUsageUrls(ctx context.Context) (*models.UsageUrls, error) {
	headers := make(map[string]string)

	var body []byte
	path := c.basePath + "/v1/analytics/usageUrls"

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetAnalyticsUsageUrlsRequest(ctx, req, headers)
}

func (c *WagClient) doGetAnalyticsUsageUrlsRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.UsageUrls, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAnalyticsUsageUrls")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAnalyticsUsageUrls")
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
		"backend":     "app-service",
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

		var output models.UsageUrls
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetAllUsageUrls makes a GET request to /v1/appUniverse/usageUrls
//
// 200: *models.UsageUrls
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAllUsageUrls(ctx context.Context) (*models.UsageUrls, error) {
	headers := make(map[string]string)

	var body []byte
	path := c.basePath + "/v1/appUniverse/usageUrls"

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetAllUsageUrlsRequest(ctx, req, headers)
}

func (c *WagClient) doGetAllUsageUrlsRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.UsageUrls, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAllUsageUrls")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAllUsageUrls")
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
		"backend":     "app-service",
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

		var output models.UsageUrls
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetApps makes a GET request to /v1/apps
// The server takes in the intersection of input parameters
// 200: []models.App
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetApps(ctx context.Context, i *models.GetAppsInput) ([]models.App, error) {
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

	return c.doGetAppsRequest(ctx, req, headers)
}

func (c *WagClient) doGetAppsRequest(ctx context.Context, req *http.Request, headers map[string]string) ([]models.App, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getApps")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getApps")
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
		"backend":     "app-service",
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

		var output []models.App
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// DeleteApp makes a DELETE request to /v1/apps/{appID}
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) DeleteApp(ctx context.Context, appID string) error {
	headers := make(map[string]string)

	var body []byte
	path, err := models.DeleteAppInputPath(appID)

	if err != nil {
		return err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "DELETE", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doDeleteAppRequest(ctx, req, headers)
}

func (c *WagClient) doDeleteAppRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "deleteApp")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "deleteApp")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// GetAppByID makes a GET request to /v1/apps/{appID}
//
// 200: *models.App
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAppByID(ctx context.Context, appID string) (*models.App, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetAppByIDInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetAppByIDRequest(ctx, req, headers)
}

func (c *WagClient) doGetAppByIDRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.App, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAppByID")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAppByID")
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
		"backend":     "app-service",
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

		var output models.App
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// UpdateApp makes a PATCH request to /v1/apps/{appID}
//
// 200: *models.App
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) UpdateApp(ctx context.Context, i *models.UpdateAppInput) (*models.App, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.App != nil {

		var err error
		body, err = json.Marshal(i.App)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doUpdateAppRequest(ctx, req, headers)
}

func (c *WagClient) doUpdateAppRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.App, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "updateApp")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "updateApp")
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
		"backend":     "app-service",
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

		var output models.App
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// CreateApp makes a PUT request to /v1/apps/{appID}
//
// 200: *models.App
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) CreateApp(ctx context.Context, i *models.CreateAppInput) (*models.App, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.App != nil {

		var err error
		body, err = json.Marshal(i.App)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doCreateAppRequest(ctx, req, headers)
}

func (c *WagClient) doCreateAppRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.App, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "createApp")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "createApp")
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
		"backend":     "app-service",
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

		var output models.App
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetAdminsForApp makes a GET request to /v1/apps/{appID}/admins
//
// 200: []models.AppAdminResponse
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAdminsForApp(ctx context.Context, appID string) ([]models.AppAdminResponse, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetAdminsForAppInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetAdminsForAppRequest(ctx, req, headers)
}

func (c *WagClient) doGetAdminsForAppRequest(ctx context.Context, req *http.Request, headers map[string]string) ([]models.AppAdminResponse, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAdminsForApp")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAdminsForApp")
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
		"backend":     "app-service",
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

		var output []models.AppAdminResponse
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

	case 404:

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// UnlinkAppAdmin makes a DELETE request to /v1/apps/{appID}/admins/{adminID}
//
// 200: nil
// 400: *models.BadRequest
// 403: *models.Forbidden
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) UnlinkAppAdmin(ctx context.Context, i *models.UnlinkAppAdminInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "DELETE", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doUnlinkAppAdminRequest(ctx, req, headers)
}

func (c *WagClient) doUnlinkAppAdminRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "unlinkAppAdmin")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "unlinkAppAdmin")
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
		"backend":     "app-service",
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

	case 403:

		var output models.Forbidden
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return err
		}
		return &output

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// LinkAppAdmin makes a PUT request to /v1/apps/{appID}/admins/{adminID}
//
// 200: nil
// 400: *models.BadRequest
// 403: *models.Forbidden
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) LinkAppAdmin(ctx context.Context, i *models.LinkAppAdminInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	if i.Permissions != nil {

		var err error
		body, err = json.Marshal(i.Permissions)

		if err != nil {
			return err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doLinkAppAdminRequest(ctx, req, headers)
}

func (c *WagClient) doLinkAppAdminRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "linkAppAdmin")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "linkAppAdmin")
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
		"backend":     "app-service",
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

	case 403:

		var output models.Forbidden
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return err
		}
		return &output

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// GetGuideConfig makes a GET request to /v1/apps/{appID}/admins/{adminID}/guides/{guideID}
//
// 200: *models.GuideConfig
// 400: *models.BadRequest
// 403: *models.Forbidden
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetGuideConfig(ctx context.Context, i *models.GetGuideConfigInput) (*models.GuideConfig, error) {
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

	return c.doGetGuideConfigRequest(ctx, req, headers)
}

func (c *WagClient) doGetGuideConfigRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.GuideConfig, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getGuideConfig")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getGuideConfig")
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
		"backend":     "app-service",
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

		var output models.GuideConfig
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

	case 403:

		var output models.Forbidden
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	case 404:

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// SetGuideConfig makes a PUT request to /v1/apps/{appID}/admins/{adminID}/guides/{guideID}
//
// 200: *models.GuideConfig
// 400: *models.BadRequest
// 403: *models.Forbidden
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) SetGuideConfig(ctx context.Context, i *models.SetGuideConfigInput) (*models.GuideConfig, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.GuideConfig != nil {

		var err error
		body, err = json.Marshal(i.GuideConfig)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doSetGuideConfigRequest(ctx, req, headers)
}

func (c *WagClient) doSetGuideConfigRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.GuideConfig, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "setGuideConfig")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "setGuideConfig")
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
		"backend":     "app-service",
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

		var output models.GuideConfig
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

	case 403:

		var output models.Forbidden
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	case 404:

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetPermissionsForAdmin makes a GET request to /v1/apps/{appID}/admins/{adminID}/permissions
//
// 200: *models.PermissionList
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetPermissionsForAdmin(ctx context.Context, i *models.GetPermissionsForAdminInput) (*models.PermissionList, error) {
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

	return c.doGetPermissionsForAdminRequest(ctx, req, headers)
}

func (c *WagClient) doGetPermissionsForAdminRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.PermissionList, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getPermissionsForAdmin")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getPermissionsForAdmin")
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
		"backend":     "app-service",
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

		var output models.PermissionList
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// VerifyAppAdmin makes a POST request to /v1/apps/{appID}/admins/{adminID}/verify
//
// 200: nil
// 400: *models.BadRequest
// 403: *models.Forbidden
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) VerifyAppAdmin(ctx context.Context, i *models.VerifyAppAdminInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doVerifyAppAdminRequest(ctx, req, headers)
}

func (c *WagClient) doVerifyAppAdminRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "verifyAppAdmin")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "verifyAppAdmin")
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
		"backend":     "app-service",
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

	case 403:

		var output models.Forbidden
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return err
		}
		return &output

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// GenerateNewBusinessToken makes a POST request to /v1/apps/{appID}/business_token
//
// 200: *models.SecretConfig
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GenerateNewBusinessToken(ctx context.Context, appID string) (*models.SecretConfig, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GenerateNewBusinessTokenInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGenerateNewBusinessTokenRequest(ctx, req, headers)
}

func (c *WagClient) doGenerateNewBusinessTokenRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.SecretConfig, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "generateNewBusinessToken")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "generateNewBusinessToken")
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
		"backend":     "app-service",
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

		var output models.SecretConfig
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetCertifications makes a GET request to /v1/apps/{appID}/certifications/{schoolYearStart}
//
// 200: *models.Certifications
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetCertifications(ctx context.Context, i *models.GetCertificationsInput) (*models.Certifications, error) {
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

	return c.doGetCertificationsRequest(ctx, req, headers)
}

func (c *WagClient) doGetCertificationsRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Certifications, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getCertifications")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getCertifications")
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
		"backend":     "app-service",
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

		var output models.Certifications
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// SetCertifications makes a POST request to /v1/apps/{appID}/certifications/{schoolYearStart}
//
// 200: *models.Certifications
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) SetCertifications(ctx context.Context, i *models.SetCertificationsInput) (*models.Certifications, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.Certifications != nil {

		var err error
		body, err = json.Marshal(i.Certifications)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doSetCertificationsRequest(ctx, req, headers)
}

func (c *WagClient) doSetCertificationsRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Certifications, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "setCertifications")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "setCertifications")
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
		"backend":     "app-service",
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

		var output models.Certifications
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetSetupStep makes a GET request to /v1/apps/{appID}/customStep
//
// 200: *models.SetupStep
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetSetupStep(ctx context.Context, appID string) (*models.SetupStep, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetSetupStepInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetSetupStepRequest(ctx, req, headers)
}

func (c *WagClient) doGetSetupStepRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.SetupStep, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getSetupStep")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getSetupStep")
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
		"backend":     "app-service",
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

		var output models.SetupStep
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// CreateSetupStep makes a PATCH request to /v1/apps/{appID}/customStep
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) CreateSetupStep(ctx context.Context, i *models.CreateSetupStepInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	if i.SetupStep != nil {

		var err error
		body, err = json.Marshal(i.SetupStep)

		if err != nil {
			return err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doCreateSetupStepRequest(ctx, req, headers)
}

func (c *WagClient) doCreateSetupStepRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "createSetupStep")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "createSetupStep")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// GetDataRules makes a GET request to /v1/apps/{appID}/data_rules
//
// 200: []models.DataRule
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetDataRules(ctx context.Context, appID string) ([]models.DataRule, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetDataRulesInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetDataRulesRequest(ctx, req, headers)
}

func (c *WagClient) doGetDataRulesRequest(ctx context.Context, req *http.Request, headers map[string]string) ([]models.DataRule, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getDataRules")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getDataRules")
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
		"backend":     "app-service",
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

		var output []models.DataRule
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

	case 404:

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// SetDataRules makes a PUT request to /v1/apps/{appID}/data_rules
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) SetDataRules(ctx context.Context, i *models.SetDataRulesInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	if i.Rules != nil {

		var err error
		body, err = json.Marshal(i.Rules)

		if err != nil {
			return err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doSetDataRulesRequest(ctx, req, headers)
}

func (c *WagClient) doSetDataRulesRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "setDataRules")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "setDataRules")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// GetManagers makes a GET request to /v1/apps/{appID}/managers
//
// 200: *models.Managers
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetManagers(ctx context.Context, appID string) (*models.Managers, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetManagersInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetManagersRequest(ctx, req, headers)
}

func (c *WagClient) doGetManagersRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Managers, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getManagers")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getManagers")
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
		"backend":     "app-service",
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

		var output models.Managers
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetOnboarding makes a GET request to /v1/apps/{appID}/onboarding
//
// 200: *models.Onboarding
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetOnboarding(ctx context.Context, appID string) (*models.Onboarding, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetOnboardingInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetOnboardingRequest(ctx, req, headers)
}

func (c *WagClient) doGetOnboardingRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Onboarding, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getOnboarding")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getOnboarding")
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
		"backend":     "app-service",
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

		var output models.Onboarding
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// UpdateOnboarding makes a PATCH request to /v1/apps/{appID}/onboarding
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) UpdateOnboarding(ctx context.Context, i *models.UpdateOnboardingInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	if i.Update != nil {

		var err error
		body, err = json.Marshal(i.Update)

		if err != nil {
			return err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doUpdateOnboardingRequest(ctx, req, headers)
}

func (c *WagClient) doUpdateOnboardingRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "updateOnboarding")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "updateOnboarding")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// InitializeOnboarding makes a PUT request to /v1/apps/{appID}/onboarding
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) InitializeOnboarding(ctx context.Context, appID string) error {
	headers := make(map[string]string)

	var body []byte
	path, err := models.InitializeOnboardingInputPath(appID)

	if err != nil {
		return err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doInitializeOnboardingRequest(ctx, req, headers)
}

func (c *WagClient) doInitializeOnboardingRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "initializeOnboarding")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "initializeOnboarding")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// DeletePlatform makes a DELETE request to /v1/apps/{appID}/platform/{clientID}
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) DeletePlatform(ctx context.Context, i *models.DeletePlatformInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "DELETE", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doDeletePlatformRequest(ctx, req, headers)
}

func (c *WagClient) doDeletePlatformRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "deletePlatform")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "deletePlatform")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// UpdatePlatform makes a PATCH request to /v1/apps/{appID}/platform/{clientID}
//
// 200: *models.Platform
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) UpdatePlatform(ctx context.Context, i *models.UpdatePlatformInput) (*models.Platform, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.Request != nil {

		var err error
		body, err = json.Marshal(i.Request)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doUpdatePlatformRequest(ctx, req, headers)
}

func (c *WagClient) doUpdatePlatformRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Platform, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "updatePlatform")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "updatePlatform")
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
		"backend":     "app-service",
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

		var output models.Platform
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetPlatformsByAppID makes a GET request to /v1/apps/{appID}/platforms
//
// 200: []models.Platform
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetPlatformsByAppID(ctx context.Context, appID string) ([]models.Platform, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetPlatformsByAppIDInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetPlatformsByAppIDRequest(ctx, req, headers)
}

func (c *WagClient) doGetPlatformsByAppIDRequest(ctx context.Context, req *http.Request, headers map[string]string) ([]models.Platform, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getPlatformsByAppID")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getPlatformsByAppID")
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
		"backend":     "app-service",
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

		var output []models.Platform
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

	case 404:

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// CreatePlatform makes a PUT request to /v1/apps/{appID}/platforms
//
// 200: *models.Platform
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) CreatePlatform(ctx context.Context, i *models.CreatePlatformInput) (*models.Platform, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.Request != nil {

		var err error
		body, err = json.Marshal(i.Request)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doCreatePlatformRequest(ctx, req, headers)
}

func (c *WagClient) doCreatePlatformRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Platform, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "createPlatform")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "createPlatform")
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
		"backend":     "app-service",
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

		var output models.Platform
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// DeleteAppSchema makes a DELETE request to /v1/apps/{appID}/schema
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) DeleteAppSchema(ctx context.Context, i *models.DeleteAppSchemaInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "DELETE", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doDeleteAppSchemaRequest(ctx, req, headers)
}

func (c *WagClient) doDeleteAppSchemaRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "deleteAppSchema")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "deleteAppSchema")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// GetAppSchema makes a GET request to /v1/apps/{appID}/schema
//
// 200: *models.AppSchema
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAppSchema(ctx context.Context, appID string) (*models.AppSchema, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetAppSchemaInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetAppSchemaRequest(ctx, req, headers)
}

func (c *WagClient) doGetAppSchemaRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.AppSchema, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAppSchema")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAppSchema")
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
		"backend":     "app-service",
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

		var output models.AppSchema
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// CreateAppSchema makes a POST request to /v1/apps/{appID}/schema
//
// 200: *models.AppSchema
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) CreateAppSchema(ctx context.Context, i *models.CreateAppSchemaInput) (*models.AppSchema, error) {
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

	return c.doCreateAppSchemaRequest(ctx, req, headers)
}

func (c *WagClient) doCreateAppSchemaRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.AppSchema, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "createAppSchema")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "createAppSchema")
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
		"backend":     "app-service",
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

		var output models.AppSchema
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// SetAppSchema makes a PUT request to /v1/apps/{appID}/schema
//
// 200: *models.AppSchema
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) SetAppSchema(ctx context.Context, i *models.SetAppSchemaInput) (*models.AppSchema, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.AppSchema != nil {

		var err error
		body, err = json.Marshal(i.AppSchema)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doSetAppSchemaRequest(ctx, req, headers)
}

func (c *WagClient) doSetAppSchemaRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.AppSchema, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "setAppSchema")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "setAppSchema")
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
		"backend":     "app-service",
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

		var output models.AppSchema
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetSecrets makes a GET request to /v1/apps/{appID}/secrets
//
// 200: *models.SecretConfig
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetSecrets(ctx context.Context, appID string) (*models.SecretConfig, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetSecretsInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetSecretsRequest(ctx, req, headers)
}

func (c *WagClient) doGetSecretsRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.SecretConfig, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getSecrets")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getSecrets")
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
		"backend":     "app-service",
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

		var output models.SecretConfig
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// RevokeOldClientSecret makes a PATCH request to /v1/apps/{appID}/secrets
//
// 200: *models.SecretConfig
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) RevokeOldClientSecret(ctx context.Context, appID string) (*models.SecretConfig, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.RevokeOldClientSecretInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "PATCH", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doRevokeOldClientSecretRequest(ctx, req, headers)
}

func (c *WagClient) doRevokeOldClientSecretRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.SecretConfig, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "revokeOldClientSecret")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "revokeOldClientSecret")
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
		"backend":     "app-service",
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

		var output models.SecretConfig
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GenerateNewClientSecret makes a POST request to /v1/apps/{appID}/secrets
//
// 200: *models.SecretConfig
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GenerateNewClientSecret(ctx context.Context, appID string) (*models.SecretConfig, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GenerateNewClientSecretInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGenerateNewClientSecretRequest(ctx, req, headers)
}

func (c *WagClient) doGenerateNewClientSecretRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.SecretConfig, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "generateNewClientSecret")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "generateNewClientSecret")
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
		"backend":     "app-service",
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

		var output models.SecretConfig
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// ResetClientSecret makes a PUT request to /v1/apps/{appID}/secrets
//
// 200: *models.SecretConfig
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) ResetClientSecret(ctx context.Context, appID string) (*models.SecretConfig, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.ResetClientSecretInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doResetClientSecretRequest(ctx, req, headers)
}

func (c *WagClient) doResetClientSecretRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.SecretConfig, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "resetClientSecret")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "resetClientSecret")
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
		"backend":     "app-service",
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

		var output models.SecretConfig
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetRecommendedSharing makes a GET request to /v1/apps/{appID}/sharing
//
// 200: *models.SharingRecommendations
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetRecommendedSharing(ctx context.Context, appID string) (*models.SharingRecommendations, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetRecommendedSharingInputPath(appID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetRecommendedSharingRequest(ctx, req, headers)
}

func (c *WagClient) doGetRecommendedSharingRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.SharingRecommendations, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getRecommendedSharing")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getRecommendedSharing")
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
		"backend":     "app-service",
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

		var output models.SharingRecommendations
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// SetRecommendedSharing makes a PUT request to /v1/apps/{appID}/sharing
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) SetRecommendedSharing(ctx context.Context, i *models.SetRecommendedSharingInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	if i.Recommendations != nil {

		var err error
		body, err = json.Marshal(i.Recommendations)

		if err != nil {
			return err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PUT", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doSetRecommendedSharingRequest(ctx, req, headers)
}

func (c *WagClient) doSetRecommendedSharingRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "setRecommendedSharing")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "setRecommendedSharing")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// UpdateAppIcon makes a POST request to /v1/apps/{appID}/update_icon
//
// 200: *models.Image
// 400: *models.BadRequest
// 404: *models.NotFound
// 422: *models.UnprocessableEntity
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) UpdateAppIcon(ctx context.Context, i *models.UpdateAppIconInput) (*models.Image, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.App != nil {

		var err error
		body, err = json.Marshal(i.App)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doUpdateAppIconRequest(ctx, req, headers)
}

func (c *WagClient) doUpdateAppIconRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Image, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "updateAppIcon")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "updateAppIcon")
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
		"backend":     "app-service",
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

		var output models.Image
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

		var output models.NotFound
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, err
		}
		return nil, &output

	case 422:

		var output models.UnprocessableEntity
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetAllCategories makes a GET request to /v1/categories
//
// 200: *models.Categories
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAllCategories(ctx context.Context) (*models.Categories, error) {
	headers := make(map[string]string)

	var body []byte
	path := c.basePath + "/v1/categories"

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetAllCategoriesRequest(ctx, req, headers)
}

func (c *WagClient) doGetAllCategoriesRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Categories, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAllCategories")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAllCategories")
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
		"backend":     "app-service",
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

		var output models.Categories
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetKnownHosts makes a GET request to /v1/knownhosts
//
// 200: []models.KnownHost
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetKnownHosts(ctx context.Context) ([]models.KnownHost, error) {
	headers := make(map[string]string)

	var body []byte
	path := c.basePath + "/v1/knownhosts"

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetKnownHostsRequest(ctx, req, headers)
}

func (c *WagClient) doGetKnownHostsRequest(ctx context.Context, req *http.Request, headers map[string]string) ([]models.KnownHost, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getKnownHosts")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getKnownHosts")
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
		"backend":     "app-service",
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

		var output []models.KnownHost
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetAllLibraryResources makes a GET request to /v1/libraryResources
//
// 200: *models.LibraryResources
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAllLibraryResources(ctx context.Context, i *models.GetAllLibraryResourcesInput) (*models.LibraryResources, error) {
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

	return c.doGetAllLibraryResourcesRequest(ctx, req, headers)
}

func (c *WagClient) doGetAllLibraryResourcesRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.LibraryResources, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAllLibraryResources")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAllLibraryResources")
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
		"backend":     "app-service",
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

		var output models.LibraryResources
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// SearchLibraryResource makes a GET request to /v1/libraryResources/search
//
// 200: *models.LibraryResources
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) SearchLibraryResource(ctx context.Context, i *models.SearchLibraryResourceInput) (*models.LibraryResources, error) {
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

	return c.doSearchLibraryResourceRequest(ctx, req, headers)
}

func (c *WagClient) doSearchLibraryResourceRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.LibraryResources, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "searchLibraryResource")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "searchLibraryResource")
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
		"backend":     "app-service",
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

		var output models.LibraryResources
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetLibraryResourceByShortname makes a GET request to /v1/libraryResources/{shortname}
//
// 200: *models.LibraryResource
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetLibraryResourceByShortname(ctx context.Context, i *models.GetLibraryResourceByShortnameInput) (*models.LibraryResource, error) {
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

	return c.doGetLibraryResourceByShortnameRequest(ctx, req, headers)
}

func (c *WagClient) doGetLibraryResourceByShortnameRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.LibraryResource, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getLibraryResourceByShortname")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getLibraryResourceByShortname")
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
		"backend":     "app-service",
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

		var output models.LibraryResource
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// UpdateLibraryResourceByShortname makes a PATCH request to /v1/libraryResources/{shortname}
//
// 200: *models.LibraryResource
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) UpdateLibraryResourceByShortname(ctx context.Context, i *models.UpdateLibraryResourceByShortnameInput) (*models.LibraryResource, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.LibraryResource != nil {

		var err error
		body, err = json.Marshal(i.LibraryResource)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doUpdateLibraryResourceByShortnameRequest(ctx, req, headers)
}

func (c *WagClient) doUpdateLibraryResourceByShortnameRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.LibraryResource, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "updateLibraryResourceByShortname")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "updateLibraryResourceByShortname")
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
		"backend":     "app-service",
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

		var output models.LibraryResource
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// CreateLibraryResource makes a POST request to /v1/libraryResources/{shortname}
//
// 200: *models.LibraryResource
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) CreateLibraryResource(ctx context.Context, i *models.CreateLibraryResourceInput) (*models.LibraryResource, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	if i.LibraryResource != nil {

		var err error
		body, err = json.Marshal(i.LibraryResource)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doCreateLibraryResourceRequest(ctx, req, headers)
}

func (c *WagClient) doCreateLibraryResourceRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.LibraryResource, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "createLibraryResource")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "createLibraryResource")
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
		"backend":     "app-service",
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

		var output models.LibraryResource
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// DeleteLibraryResourceLink makes a DELETE request to /v1/libraryResources/{shortname}/link
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) DeleteLibraryResourceLink(ctx context.Context, shortname string) error {
	headers := make(map[string]string)

	var body []byte
	path, err := models.DeleteLibraryResourceLinkInputPath(shortname)

	if err != nil {
		return err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "DELETE", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doDeleteLibraryResourceLinkRequest(ctx, req, headers)
}

func (c *WagClient) doDeleteLibraryResourceLinkRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "deleteLibraryResourceLink")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "deleteLibraryResourceLink")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

// GetValidPermissions makes a GET request to /v1/permissions
//
// 200: *models.GetValidPermissionsResponse
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetValidPermissions(ctx context.Context) (*models.GetValidPermissionsResponse, error) {
	headers := make(map[string]string)

	var body []byte
	path := c.basePath + "/v1/permissions"

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetValidPermissionsRequest(ctx, req, headers)
}

func (c *WagClient) doGetValidPermissionsRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.GetValidPermissionsResponse, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getValidPermissions")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getValidPermissions")
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
		"backend":     "app-service",
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

		var output models.GetValidPermissionsResponse
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetPlatforms makes a GET request to /v1/platforms
// The server takes in the intersection of input parameters
// 200: []models.Platform
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetPlatforms(ctx context.Context, i *models.GetPlatformsInput) ([]models.Platform, error) {
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

	return c.doGetPlatformsRequest(ctx, req, headers)
}

func (c *WagClient) doGetPlatformsRequest(ctx context.Context, req *http.Request, headers map[string]string) ([]models.Platform, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getPlatforms")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getPlatforms")
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
		"backend":     "app-service",
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

		var output []models.Platform
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetPlatformByClientID makes a GET request to /v1/platforms/{clientID}
//
// 200: *models.Platform
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetPlatformByClientID(ctx context.Context, clientID string) (*models.Platform, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetPlatformByClientIDInputPath(clientID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetPlatformByClientIDRequest(ctx, req, headers)
}

func (c *WagClient) doGetPlatformByClientIDRequest(ctx context.Context, req *http.Request, headers map[string]string) (*models.Platform, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getPlatformByClientID")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getPlatformByClientID")
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
		"backend":     "app-service",
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

		var output models.Platform
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

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// GetAppsForAdmin makes a GET request to /v2/admins/{adminID}/apps
//
// 200: []models.AppForAdminResponse
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetAppsForAdmin(ctx context.Context, adminID string) ([]models.AppForAdminResponse, error) {
	headers := make(map[string]string)

	var body []byte
	path, err := models.GetAppsForAdminInputPath(adminID)

	if err != nil {
		return nil, err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return c.doGetAppsForAdminRequest(ctx, req, headers)
}

func (c *WagClient) doGetAppsForAdminRequest(ctx context.Context, req *http.Request, headers map[string]string) ([]models.AppForAdminResponse, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getAppsForAdmin")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getAppsForAdmin")
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
		"backend":     "app-service",
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

		var output []models.AppForAdminResponse
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

	case 404:

		var output models.NotFound
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
		return nil, &models.InternalError{Message: "Unknown response"}
	}
}

// OverrideConfig makes a POST request to /v2/apps/{srcAppID}/override-config/{destAppID}
//
// 200: nil
// 400: *models.BadRequest
// 404: *models.NotFound
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) OverrideConfig(ctx context.Context, i *models.OverrideConfigInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doOverrideConfigRequest(ctx, req, headers)
}

func (c *WagClient) doOverrideConfigRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "overrideConfig")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "overrideConfig")
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
		"backend":     "app-service",
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

	case 404:

		var output models.NotFound
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
		return &models.InternalError{Message: "Unknown response"}
	}
}

func shortHash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))[0:6]
}
