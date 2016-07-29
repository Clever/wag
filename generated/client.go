package main

import "net/http"
import "net/url"
import "encoding/json"
import "strings"
import "errors"
import "golang.org/x/net/context"
import "bytes"

var _ = json.Marshal
var _ = strings.Replace

func GetBookByID(ctx context.Context, i *GetBookByIDInput) (GetBookByIDOutput, error) {
	path := "http://localhost:8080" + "/books/{bookID}"
	urlVals := url.Values{}
	var body []byte

	urlVals.Add("author", i.Author)
	path = strings.Replace(path, "{bookID}", i.BookID, -1)
	path = path + "?" + urlVals.Encode()

	body, _ = json.Marshal(i.TestBook)

	client := &http.Client{}
	req, _ := http.NewRequest("get", path, bytes.NewBuffer(body))
	req.Header.Set("authorization", i.Authorization)
	resp, _ := client.Do(req)

	switch resp.StatusCode {
	case 200:
		// TODO: Actually read the body and set the data on the Output object correctly
		return GetBookByID200Output{}, nil
	case 404:
		return nil, GetBookByID404Output{}
	default:
		return nil, errors.New("Unknown response")
	}
}

