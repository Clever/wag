package main

import "net/http"
import "net/url"
import "encoding/json"
import "strings"
import "golang.org/x/net/context"
import "bytes"

var _ = json.Marshal
var _ = strings.Replace

func GetBookByID(ctx context.Context, i *GetBookByIDInput) {
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
	client.Do(req)

}

