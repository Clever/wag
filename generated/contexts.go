package main

import (
	"net/http"
	"golang.org/x/net/context"
	"github.com/gorilla/mux"
	"encoding/json"
	"github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"
)

var _ = json.Marshal
type GetBookByIDInput struct {
	Author string
	BookID string
	Authorization string
	TestBook models.Book
}
func NewGetBookByIDInput(r *http.Request) (*GetBookByIDInput, error) {
	var input GetBookByIDInput

	input.Author = r.URL.Query().Get("author")
	input.BookID = mux.Vars(r)["bookID"]
	input.Authorization = r.Header.Get("authorization")
json.NewDecoder(r.Body).Decode(&input.TestBook)

	return &input, nil
}
func (i GetBookByIDInput) Validate() error{
	if err := i.TestBook.Validate(nil); err != nil {
		return err
	}

	return nil
}

type Controller interface {
	GetBookByID(ctx context.Context, input *GetBookByIDInput) (GetBookByIDOutput, error)
}
