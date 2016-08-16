package client

import (
	"bytes"
	"encoding/json"
	"github.com/Clever/wag/generated/models"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var _ = json.Marshal
var _ = strings.Replace

var _ = strconv.FormatInt

type Client struct {
	BasePath    string
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

	return Client{requestDoer: &retry, retryDoer: &retry, transport: &http.Transport{}, BasePath: basePath}
}

func (c Client) WithRetries(retries int) Client {
	c.retryDoer.defaultRetries = retries
	return c
}

func (c Client) GetBooks(ctx context.Context, i *models.GetBooksInput) (models.GetBooksOutput, error) {
	path := c.BasePath + "/v1/books"
	urlVals := url.Values{}
	var body []byte

	urlVals.Add("author", i.Author)
	urlVals.Add("available", strconv.FormatBool(i.Available))
	urlVals.Add("maxPages", strconv.FormatFloat(i.MaxPages, 'E', -1, 64))
	path = path + "?" + urlVals.Encode()

	client := &http.Client{Transport: c.transport}
	req, _ := http.NewRequest("GET", path, bytes.NewBuffer(body))

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "getBooks")
	resp, err := c.requestDoer.Do(client, req.WithContext(ctx))
	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}
	switch resp.StatusCode {
	case 200:

		var output models.GetBooks200Output
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

func (c Client) GetBookByID(ctx context.Context, i *models.GetBookByIDInput) (models.GetBookByIDOutput, error) {
	path := c.BasePath + "/v1/books/{bookID}"
	urlVals := url.Values{}
	var body []byte

	path = strings.Replace(path, "{bookID}", strconv.FormatInt(i.BookID, 10), -1)
	path = path + "?" + urlVals.Encode()

	client := &http.Client{Transport: c.transport}
	req, _ := http.NewRequest("GET", path, bytes.NewBuffer(body))
	req.Header.Set("authorization", i.Authorization)

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
		return output, nil
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

func (c Client) CreateBook(ctx context.Context, i *models.CreateBookInput) (models.CreateBookOutput, error) {
	path := c.BasePath + "/v1/books/{bookID}"
	urlVals := url.Values{}
	var body []byte

	path = path + "?" + urlVals.Encode()

	body, _ = json.Marshal(i.NewBook)

	client := &http.Client{Transport: c.transport}
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(body))

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "createBook")
	resp, err := c.requestDoer.Do(client, req.WithContext(ctx))
	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}
	switch resp.StatusCode {
	case 200:

		var output models.CreateBook200Output
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
