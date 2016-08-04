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

var clientHandler RequestHandler

type Client struct {
	BasePath  string
	Handler   RequestHandler
	Transport *http.Transport
}

// NewClient creates a new client. The base path and http transport are configurable
func NewClient(basePath string) Client {
	var handler RequestHandler
	handler = baseRequestHandler{}
	handler = tracingHandler{handler: clientHandler}

	// We could add a default timeout in the http transport here...
	return Client{Handler: handler, Transport: nil, BasePath: basePath}
}

func (c Client) GetBooks(ctx context.Context, i *models.GetBooksInput) (models.GetBooksOutput, error) {
	path := c.BasePath + "/v1/books"
	urlVals := url.Values{}
	var body []byte

	urlVals.Add("author", i.Author)
	urlVals.Add("available", strconv.FormatBool(i.Available))
	urlVals.Add("maxPages", strconv.FormatFloat(i.MaxPages, 'E', -1, 64))
	path = path + "?" + urlVals.Encode()

	client := &http.Client{Transport: c.Transport}
	req, _ := http.NewRequest("GET", path, bytes.NewBuffer(body))

	// Inject tracing headers
	ctx = context.WithValue(ctx, opNameCtx{}, "GetBooks")
	resp, err := c.Handler.HandleRequest(ctx, client, req)
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
