package main

// This is just a shameless copy and paste job to sanity check things..

import "net/http"
import "net/url"
import "fmt"
import "encoding/json"
import "strings"
import "golang.org/x/net/context"
import "bytes"
import "errors"
import opentracing "github.com/opentracing/opentracing-go"
import "github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"
import "strconv"

var _ = json.Marshal
var _ = strings.Replace
var _ = fmt.Printf

type GetBookByIDInput struct {
	BookID int64
	Authorization string
}

type GetBookByID404Output struct{}

func (g GetBookByID404Output) Error() string {
        return "got a 404"
}

type GetBookByIDDefaultOutput struct{}

type GetBookByID200Output struct {
	Data models.Book
}

type GetBookByIDOutput interface {}

func GetBookByID(ctx context.Context, i *GetBookByIDInput) (GetBookByIDOutput, error) {
	path := "http://localhost:8080" + "/books/{bookID}"
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


func main() {

        output, err := GetBookByID(context.Background(), &GetBookByIDInput{BookID: 1234})
        fmt.Printf("Output: %+v\n", output)
        fmt.Printf("Error: %+v\n", err)
}
