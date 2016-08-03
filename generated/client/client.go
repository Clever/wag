package client

import "net/http"
import "net/url"
import "encoding/json"
import "strings"
import "golang.org/x/net/context"
import "bytes"
import "fmt"
import "strconv"
import "github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"

import opentracing "github.com/opentracing/opentracing-go"

var _ = json.Marshal
var _ = strings.Replace

var _ = strconv.FormatInt

func GetBookByID(ctx context.Context, i *models.GetBookByIDInput) (models.GetBookByIDOutput, error) {
	path := "http://localhost:8080" + "/v1/books/{bookID}"
	urlVals := url.Values{}
	var body []byte

	path = strings.Replace(path, "{bookID}", strconv.FormatInt(i.BookID, 10), -1)
	path = path + "?" + urlVals.Encode()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", path, bytes.NewBuffer(body))
	req.Header.Set("authorization", i.Authorization)

	// Inject tracing headers
	opName := "GetBookByID"
	var sp opentracing.Span
	// TODO: add tags relating to input data?
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		sp = opentracing.StartSpan(opName, opentracing.ChildOf(parentSpan.Context()))
	} else {
		sp = opentracing.StartSpan(opName)
	}
	if err := sp.Tracer().Inject(sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
		return nil, fmt.Errorf("couldn't inject tracing headers (%v)", err)
	}

	resp, _ := client.Do(req)

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

func CreateBook(ctx context.Context, i *models.CreateBookInput) (models.CreateBookOutput, error) {
	path := "http://localhost:8080" + "/v1/books/{bookID}"
	urlVals := url.Values{}
	var body []byte

	path = path + "?" + urlVals.Encode()

	body, _ = json.Marshal(i.NewBook)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(body))

	// Inject tracing headers
	opName := "CreateBook"
	var sp opentracing.Span
	// TODO: add tags relating to input data?
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		sp = opentracing.StartSpan(opName, opentracing.ChildOf(parentSpan.Context()))
	} else {
		sp = opentracing.StartSpan(opName)
	}
	if err := sp.Tracer().Inject(sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
		return nil, fmt.Errorf("couldn't inject tracing headers (%v)", err)
	}

	resp, _ := client.Do(req)

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

func GetBooks(ctx context.Context, i *models.GetBooksInput) (models.GetBooksOutput, error) {
	path := "http://localhost:8080" + "/v1/books"
	urlVals := url.Values{}
	var body []byte

	urlVals.Add("author", i.Author)
	urlVals.Add("available", strconv.FormatBool(i.Available))
	urlVals.Add("maxPages", strconv.FormatFloat(i.MaxPages, 'E', -1, 64))
	path = path + "?" + urlVals.Encode()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", path, bytes.NewBuffer(body))

	// Inject tracing headers
	opName := "GetBooks"
	var sp opentracing.Span
	// TODO: add tags relating to input data?
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		sp = opentracing.StartSpan(opName, opentracing.ChildOf(parentSpan.Context()))
	} else {
		sp = opentracing.StartSpan(opName)
	}
	if err := sp.Tracer().Inject(sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
		return nil, fmt.Errorf("couldn't inject tracing headers (%v)", err)
	}

	resp, _ := client.Do(req)

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

