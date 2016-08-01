package generated

import "net/http"
import "net/url"
import "encoding/json"
import "strings"
import "errors"
import "golang.org/x/net/context"
import "bytes"
import "fmt"
import opentracing "github.com/opentracing/opentracing-go"

var _ = json.Marshal
var _ = strings.Replace

func GetBookByID(ctx context.Context, i *GetBookByIDInput) (GetBookByIDOutput, error) {
	path := "http://localhost:8080" + "/books/{bookID}"
	urlVals := url.Values{}
	var body []byte

	path = strings.Replace(path, "{bookID}", i.BookID, -1)
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
	case 404:
		return nil, GetBookByID404Output{}
	case 200:

		var output GetBookByID200Output
		if err := json.NewDecoder(resp.Body).Decode(&output.Data); err != nil {
			return nil, err
		}
		return output, nil
	default:
		return nil, errors.New("Unknown response")
	}
}

func CreateBook(ctx context.Context, i *CreateBookInput) (CreateBookOutput, error) {
	path := "http://localhost:8080" + "/books/{bookID}"
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

		var output CreateBook200Output
		if err := json.NewDecoder(resp.Body).Decode(&output.Data); err != nil {
			return nil, err
		}
		return output, nil
	default:
		return nil, errors.New("Unknown response")
	}
}

func GetBooks(ctx context.Context, i *GetBooksInput) (GetBooksOutput, error) {
	path := "http://localhost:8080" + "/books"
	urlVals := url.Values{}
	var body []byte

	urlVals.Add("author", i.Author)
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

		var output GetBooks200Output
		if err := json.NewDecoder(resp.Body).Decode(&output.Data); err != nil {
			return nil, err
		}
		return output, nil
	default:
		return nil, errors.New("Unknown response")
	}
}

