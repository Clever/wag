package client

import "net/http"
import "net/url"
import "encoding/json"
import "strings"
import "golang.org/x/net/context"
import "bytes"
import "strconv"
import "github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"

var _ = json.Marshal
var _ = strings.Replace

var _ = strconv.FormatInt

type Client struct {
	BasePath    string
	requestDoer doer
	transport   *http.Transport
	// Keep the retry doer around so that we can set the number of retries
	retryDoer retryDoer
}

// NewClient creates a new client. The base path and http transport are configurable
func NewClient(basePath string) Client {
	base := baseDoer{}
	tracing := tracingDoer{d: base}
	retry := retryDoer{d: tracing, defaultRetries: 1}

	return Client{requestDoer: retry, retryDoer: retry, transport: nil, BasePath: basePath}
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

	// Inject tracing headers
	ctx = context.WithValue(ctx, opNameCtx{}, "GetBooks")
	resp, err := c.requestDoer.Do(ctx, client, req)
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
