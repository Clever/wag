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

var _ = json.Marshal
var _ = strings.Replace
var _ = fmt.Printf

type GetBookByIDInput struct {
	Author string
	BookID string
	Authorization string
	TestBook map[string]string
}

type GetBookByID404Output struct{}

func (g GetBookByID404Output) Error() string {
        return "got a 404"
}

type GetBookByIDDefaultOutput struct{}

type GetBookByID200Output struct {
}

type GetBookByIDOutput interface {}

func GetBookByID(ctx context.Context, i *GetBookByIDInput) (GetBookByIDOutput, error) {
	path := "http://localhost:8080" + "/books/{bookID}"
	urlVals := url.Values{}
	var body []byte
	urlVals.Add("author", i.Author)
	path = strings.Replace(path, "{bookID}", i.BookID, -1)
	body, _ = json.Marshal(i.TestBook)
	path = path + "?" + urlVals.Encode()
	client := &http.Client{}
	req, _ := http.NewRequest("get", path, bytes.NewBuffer(body))
	req.Header.Set("authorization", i.Authorization)
	resp, _ := client.Do(req)

	switch resp.StatusCode {
	case 200:
		return GetBookByID200Output{}, nil
	case 404:
		return nil, GetBookByID404Output{}
	default:
		return nil, errors.New("Unknown response")
	}
}


func main() {

        output, err := GetBookByID(nil, &GetBookByIDInput{Author: "Kyle", BookID: "1234"})
        fmt.Printf("Output: %+v\n", output)
        fmt.Printf("Error: %+v\n", err)
}
