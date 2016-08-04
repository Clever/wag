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

type RequestHandler interface {
	// If the response returns an error, everything above it should stay the same...
	HandleRequest(ctx context.Context, c *http.Client, r *http.Request) (*http.Response, error)
}

var clientHandler RequestHandler

// At some point client will be a struct we can modify...
func initClient() {
	clientHandler = baseRequestHandler{}
	clientHandler = tracingHandler{handler: clientHandler}
}

type opNameCtx struct{}

// TODO: Add something like what http has for this for handlers
type baseRequestHandler struct{}

func (b baseRequestHandler) HandleRequest(ctx context.Context, c *http.Client, r *http.Request) (*http.Response, error) {
	return c.Do(r)
}

type tracingHandler struct {
	handler RequestHandler
}

func (t tracingHandler) HandleRequest(ctx context.Context, c *http.Client, r *http.Request) (*http.Response, error) {
	opName := ctx.Value(opNameCtx{}).(string)
	var sp opentracing.Span
	// TODO: add tags relating to input data?
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		sp = opentracing.StartSpan(opName, opentracing.ChildOf(parentSpan.Context()))
	} else {
		sp = opentracing.StartSpan(opName)
	}
	if err := sp.Tracer().Inject(sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header)); err != nil {
		return nil, fmt.Errorf("couldn't inject tracing headers (%v)", err)
	}
	return t.handler.HandleRequest(ctx, c, r)
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
	ctx = context.WithValue(ctx, opNameCtx{}, "GetBooks")
	resp, err := clientHandler.HandleRequest(ctx, client, req)
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
	default:
		return nil, models.DefaultInternalError{Msg: "Unknown response"}
	}
}
