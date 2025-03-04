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

	"github.com/Clever/wag/samples/gen-go-client-only/models/v9"

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

// WagClient is used to make requests to the swagger-test service.
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
		requestDoer: &base,
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
// the three env vars: SERVICE_SWAGGER_TEST_HTTP_(HOST/PORT/PROTO) to be set. Otherwise it returns an error.
// The logger provided should be specifically created for this wag client. If tracing is required,
// provide an instrumented transport using the wag clientconfig module. If no tracing is required, pass nil to use
// the default transport.
func NewFromDiscovery(logger wcl.WagClientLogger, transport *http.RoundTripper) (*WagClient, error) {
	url, err := discovery.URL("swagger-test", "default")
	if err != nil {
		url, err = discovery.URL("swagger-test", "http") // Added fallback to maintain reverse compatibility
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
		bs, _ := ioutil.ReadAll(resp.Body)
		return nil, "", models.UnknownResponse{StatusCode: int64(resp.StatusCode), Body: string(bs)}
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
		bs, _ := ioutil.ReadAll(resp.Body)
		return nil, "", models.UnknownResponse{StatusCode: int64(resp.StatusCode), Body: string(bs)}
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
		bs, _ := ioutil.ReadAll(resp.Body)
		return nil, "", models.UnknownResponse{StatusCode: int64(resp.StatusCode), Body: string(bs)}
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
		bs, _ := ioutil.ReadAll(resp.Body)
		return nil, models.UnknownResponse{StatusCode: int64(resp.StatusCode), Body: string(bs)}
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
		bs, _ := ioutil.ReadAll(resp.Body)
		return nil, models.UnknownResponse{StatusCode: int64(resp.StatusCode), Body: string(bs)}
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
		bs, _ := ioutil.ReadAll(resp.Body)
		return nil, models.UnknownResponse{StatusCode: int64(resp.StatusCode), Body: string(bs)}
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
		bs, _ := ioutil.ReadAll(resp.Body)
		return nil, models.UnknownResponse{StatusCode: int64(resp.StatusCode), Body: string(bs)}
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

// LowercaseModelsTest makes a POST request to /lowercaseModelsTest/{pathParam}
// testing that we can use a lowercase name for a model
// 200: nil
// 400: *models.BadRequest
// 500: *models.InternalError
// default: client side HTTP errors, for example: context.DeadlineExceeded.
func (c *WagClient) LowercaseModelsTest(ctx context.Context, i *models.LowercaseModelsTestInput) error {
	headers := make(map[string]string)

	var body []byte
	path, err := i.Path()

	if err != nil {
		return err
	}

	path = c.basePath + path

	if i.Lowercase != nil {

		var err error
		body, err = json.Marshal(i.Lowercase)

		if err != nil {
			return err
		}

	}

	req, err := http.NewRequestWithContext(ctx, "POST", path, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	return c.doLowercaseModelsTestRequest(ctx, req, headers)
}

func (c *WagClient) doLowercaseModelsTestRequest(ctx context.Context, req *http.Request, headers map[string]string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Canonical-Resource", "lowercaseModelsTest")
	req.Header.Set(VersionHeader, Version)

	for field, value := range headers {
		req.Header.Set(field, value)
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "lowercaseModelsTest")
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
