package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Clever/wag/generated/models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// Client makes requests to Swagger Test
type Client struct {
	basePath    string
	requestDoer doer
	transport   *http.Transport
	// Keep the retry doer around so that we can set the number of retries
	retryDoer *retryDoer
}

// New creates a new client. The base path and http transport are configurable
func New(basePath string) Client {
	base := baseDoer{}
	tracing := tracingDoer{d: base}
	retry := retryDoer{d: tracing, defaultRetries: 1}

	return Client{requestDoer: &retry, retryDoer: &retry, transport: &http.Transport{}, basePath: basePath}
}

// WithRetries returns a new client that retries all GET operations until they either succeed or fail the
// number of times specified.
func (c Client) WithRetries(retries int) Client {
	c.retryDoer.defaultRetries = retries
	return c
}

// JoinByFormat joins a string array by a known format:
//	 csv: comma separated value (default)
//   ssv: space separated value
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

// GetBooks returns...
func (c Client) GetBooks(ctx context.Context, i *models.GetBooksInput) ([]models.Book, error) {
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
	resp, err := c.requestDoer.Do(client, req.WithContext(ctx))

	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}

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

// GetBookByID returns...
func (c Client) GetBookByID(ctx context.Context, i *models.GetBookByIDInput) (models.GetBookByIDOutput, error) {
<<<<<<< 8840822c0bcfcfe7eb5289400cca5cd013ae9222
	path := c.BasePath + "/v1/books/{book_id}"
=======
	path := c.basePath + "/v1/books/{bookID}"
>>>>>>> Fix client linting
	urlVals := url.Values{}
	var body []byte

	path = strings.Replace(path, "{book_id}", strconv.FormatInt(i.BookID, 10), -1)
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
	resp, err := c.requestDoer.Do(client, req.WithContext(ctx))

	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}

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

// CreateBook returns...
func (c Client) CreateBook(ctx context.Context, i *models.CreateBookInput) (*models.Book, error) {
<<<<<<< 8840822c0bcfcfe7eb5289400cca5cd013ae9222
	path := c.BasePath + "/v1/books/{book_id}"
=======
	path := c.basePath + "/v1/books/{bookID}"
>>>>>>> Fix client linting
	urlVals := url.Values{}
	var body []byte

	path = path + "?" + urlVals.Encode()

	if i.NewBook != nil {

		var err error
		body, err = json.Marshal(i.NewBook)

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
	resp, err := c.requestDoer.Do(client, req.WithContext(ctx))

	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}

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

<<<<<<< 8840822c0bcfcfe7eb5289400cca5cd013ae9222
func (c Client) HealthCheck(ctx context.Context) error {
	path := c.BasePath + "/v1/health/check"
=======
// HealthCheck returns...
func (c Client) HealthCheck(ctx context.Context, i *models.HealthCheckInput) error {
	path := c.basePath + "/v1/health/check"
>>>>>>> Fix client linting
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
	resp, err := c.requestDoer.Do(client, req.WithContext(ctx))

	if err != nil {
		return models.DefaultInternalError{Msg: err.Error()}
	}

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
