package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Clever/wag/gen-go/models"

	discovery "github.com/Clever/discovery-go"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// Client is used to make requests to the swagger-test service.
type Client struct {
	basePath    string
	requestDoer doer
	transport   *http.Transport
	timeout     time.Duration
	// Keep the retry doer around so that we can set the number of retries
	retryDoer      *retryDoer
	defaultTimeout time.Duration
}

// New creates a new client. The base path and http transport are configurable.
func New(basePath string) *Client {
	base := baseDoer{}
	tracing := tracingDoer{d: base}
	retry := retryDoer{d: tracing, defaultRetries: 1}

	return &Client{requestDoer: &retry, retryDoer: &retry, defaultTimeout: 10 * time.Second,
		transport: &http.Transport{}, basePath: basePath}
}

// NewFromDiscovery creates a client from the discovery environment variables. There must be three
// env vars with the format: SERVICE_SWAGGER_TEST_HTTP_(HOST/PORT/PROTO)
func NewFromDiscovery() (*Client, error) {
	url, err := discovery.URL("swagger-test", "http")
	if err != nil {
		return nil, err
	}
	return New(url), nil
}

// WithRetries returns a new client that retries all GET operations until they either succeed or fail the
// number of times specified.
func (c *Client) WithRetries(retries int) *Client {
	c.retryDoer.defaultRetries = retries
	return c
}

// WithTimeout returns a new client that has the specified timeout on all operations. To make a single request
// have a timeout use context.WithTimeout as described here: https://godoc.org/golang.org/x/net/context#WithTimeout.
func (c *Client) WithTimeout(timeout time.Duration) *Client {
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

// GetBooks makes a GET request to /books.
// Returns a list of books
func (c *Client) GetBooks(ctx context.Context, i *models.GetBooksInput) ([]models.Book, error) {
	path := c.basePath + "/v1/books"
	urlVals := url.Values{}
	var body []byte

	if i.Authors != nil {
		urlVals.Add("authors", JoinByFormat(i.Authors, ""))
	}
	if i.Available != nil {
		urlVals.Add("available", strconv.FormatBool(*i.Available))
	}
	if i.State != nil {
		urlVals.Add("state", *i.State)
	}
	if i.Published != nil {
		urlVals.Add("published", (*i.Published).String())
	}
	if i.SnakeCase != nil {
		urlVals.Add("snake_case", *i.SnakeCase)
	}
	if i.Completed != nil {
		urlVals.Add("completed", (*i.Completed).String())
	}
	if i.MaxPages != nil {
		urlVals.Add("maxPages", strconv.FormatFloat(*i.MaxPages, 'E', -1, 64))
	}
	if i.MinPages != nil {
		urlVals.Add("min_pages", strconv.FormatInt(int64(*i.MinPages), 10))
	}
	if i.PagesToTime != nil {
		urlVals.Add("pagesToTime", strconv.FormatFloat(float64(*i.PagesToTime), 'E', -1, 32))
	}
	path = path + "?" + urlVals.Encode()

	client := &http.Client{Transport: c.transport}
	req, err := http.NewRequest("GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
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
	resp, err := c.requestDoer.Do(client, req)

	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:

		var output []models.Book
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return output, nil

	case 400:
		var output models.DefaultBadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output

	case 500:
		var output models.DefaultInternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output

	default:
		return nil, models.DefaultInternalError{Msg: "Unknown response"}
	}
}

// GetBookByID makes a GET request to /books/{book_id}.
// Returns a book
func (c *Client) GetBookByID(ctx context.Context, i *models.GetBookByIDInput) (models.GetBookByIDOutput, error) {
	path := c.basePath + "/v1/books/{book_id}"
	urlVals := url.Values{}
	var body []byte

	path = strings.Replace(path, "{book_id}", strconv.FormatInt(i.BookID, 10), -1)
	if i.AuthorID != nil {
		urlVals.Add("authorID", *i.AuthorID)
	}
	if i.RandomBytes != nil {
		urlVals.Add("randomBytes", string(*i.RandomBytes))
	}
	path = path + "?" + urlVals.Encode()

	client := &http.Client{Transport: c.transport}
	req, err := http.NewRequest("GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	if i.Authorization != nil {
		req.Header.Set("authorization", *i.Authorization)
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
	resp, err := c.requestDoer.Do(client, req)

	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:

		var output models.GetBookByID200Output
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return &output, nil
	case 204:
		var output models.GetBookByID204Output
		return output, nil
	case 401:
		var output models.GetBookByID401Output
		return nil, output
	case 404:
		return nil, models.GetBookByID404Output{}

	case 400:
		var output models.DefaultBadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output

	case 500:
		var output models.DefaultInternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output

	default:
		return nil, models.DefaultInternalError{Msg: "Unknown response"}
	}
}

// CreateBook makes a POST request to /books/{book_id}.
// Creates a book
func (c *Client) CreateBook(ctx context.Context, i *models.Book) (*models.Book, error) {
	path := c.basePath + "/v1/books/{book_id}"
	urlVals := url.Values{}
	var body []byte

	path = path + "?" + urlVals.Encode()

	if i != nil {

		var err error
		body, err = json.Marshal(i)

		if err != nil {
			return nil, err
		}

	}

	client := &http.Client{Transport: c.transport}
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
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
	resp, err := c.requestDoer.Do(client, req)

	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:

		var output models.Book
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return &output, nil

	case 400:
		var output models.DefaultBadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output

	case 500:
		var output models.DefaultInternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output

	default:
		return nil, models.DefaultInternalError{Msg: "Unknown response"}
	}
}

// GetBookByID2 makes a GET request to /books/{id}.
// Retrieve a book
func (c *Client) GetBookByID2(ctx context.Context, i *models.GetBookByID2Input) (*models.Book, error) {
	path := c.basePath + "/v1/books/{id}"
	urlVals := url.Values{}
	var body []byte

	path = strings.Replace(path, "{id}", i.ID, -1)
	path = path + "?" + urlVals.Encode()

	client := &http.Client{Transport: c.transport}
	req, err := http.NewRequest("GET", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
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
	resp, err := c.requestDoer.Do(client, req)

	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:

		var output models.Book
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return &output, nil
	case 404:
		var output models.GetBookByID2404Output
		return nil, output

	case 400:
		var output models.DefaultBadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output

	case 500:
		var output models.DefaultInternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}
		return nil, output

	default:
		return nil, models.DefaultInternalError{Msg: "Unknown response"}
	}
}

// HealthCheck makes a GET request to /health/check.
// Checks if the service is healthy
func (c *Client) HealthCheck(ctx context.Context) error {
	path := c.basePath + "/v1/health/check"
	urlVals := url.Values{}
	var body []byte

	path = path + "?" + urlVals.Encode()

	client := &http.Client{Transport: c.transport}
	req, err := http.NewRequest("GET", path, bytes.NewBuffer(body))

	if err != nil {
		return err
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
	resp, err := c.requestDoer.Do(client, req)

	if err != nil {
		return models.DefaultInternalError{Msg: err.Error()}
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return nil

	case 400:
		var output models.DefaultBadRequest
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return models.DefaultInternalError{Msg: err.Error()}
		}
		return output

	case 500:
		var output models.DefaultInternalError
		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return models.DefaultInternalError{Msg: err.Error()}
		}
		return output

	default:
		return models.DefaultInternalError{Msg: "Unknown response"}
	}
}
