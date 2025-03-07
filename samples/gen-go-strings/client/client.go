package client

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Clever/wag/samples/gen-go-strings/models/v9"

	discovery "github.com/Clever/discovery-go"
	wcl "github.com/Clever/wag/logging/wagclientlogger"
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
	retryDoer      *retryDoer
	defaultTimeout time.Duration
	logger         wcl.WagClientLogger
}

var _ Client = (*WagClient)(nil)

// New creates a new client. The base path, logger, and http transport are configurable.
// The logger provided should be specifically created for this wag client. If tracing is required,
// provide an instrumented transport using the wag clientconfig module. If no tracing is required, pass nil to use
// the default transport.
func New(basePath string, logger wcl.WagClientLogger, transport *http.RoundTripper) *WagClient {

	t := http.DefaultTransport
	if transport != nil {
		t = *transport
	}

	basePath = strings.TrimSuffix(basePath, "/")
	base := baseDoer{}

	// Don't use the default retry policy since its 5 retries can 5X the traffic
	retry := retryDoer{d: base, retryPolicy: SingleRetryPolicy{}}

	client := &WagClient{
		basePath:    basePath,
		requestDoer: &retry,
		client: &http.Client{
			Transport: t,
		},
		retryDoer:      &retry,
		defaultTimeout: 5 * time.Second,
		logger:         logger,
	}
	return client
}

// NewFromDiscovery creates a client from the discovery environment variables. This method requires
// the three env vars: SERVICE_NIL_TEST_HTTP_(HOST/PORT/PROTO) to be set. Otherwise it returns an error.
// The logger provided should be specifically created for this wag client. If tracing is required,
// provide an instrumented transport using the wag clientconfig module. If no tracing is required, pass nil to use
// the default transport.
func NewFromDiscovery(logger wcl.WagClientLogger, transport *http.RoundTripper) (*WagClient, error) {
	url, err := discovery.URL("nil-test", "default")
	if err != nil {
		url, err = discovery.URL("nil-test", "http") // Added fallback to maintain reverse compatibility
		if err != nil {
			return nil, err
		}
	}
	return New(url, logger, transport), nil
}

// SetRetryPolicy sets a the given retry policy for all requests.
func (c *WagClient) SetRetryPolicy(retryPolicy RetryPolicy) {
	c.retryDoer.retryPolicy = retryPolicy
}

// SetLogger allows for setting a custom logger
func (c *WagClient) SetLogger(l wcl.WagClientLogger) {
	c.logger = l
}

// SetTimeout sets a timeout on all operations for the client. To make a single request with a shorter timeout
// than the default on the client, use context.WithTimeout as described here: https://godoc.org/golang.org/x/net/context#WithTimeout.
func (c *WagClient) SetTimeout(timeout time.Duration) {
	c.defaultTimeout = timeout
}

// GetDistricts makes a POST request to /check
//
// 200: nil
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) GetDistricts(ctx context.Context, i *models.GetDistrictsInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	if i.Where != nil {

		var err error
		body, err = json.Marshal(i.Where)

		if err != nil {
			return err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doGetDistrictsRequest(ctx, req, headers)
}

func (c *WagClient) doGetDistrictsRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "getDistricts")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getDistricts")
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
	if err == nil && retCode > 399 && retCode < 500 {
		logData["message"] = resp.Status
		c.logger.Log(wcl.Warning, "client-request-finished", logData)
	}
	if err == nil && retCode > 499 {
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
