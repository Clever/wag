package main

import (
	"net/http"
	"golang.org/x/net/context"
	"github.com/gorilla/mux"
	"encoding/json"
)

type GetBookByIDInput struct {
	Author string
	BookID string
	Authorization string
	TestBook Book
}
func NewGetBookByIDInput(r *http.Request) (*GetBookByIDInput, error) {
	var input GetBookByIDInput

	input.Author = r.URL.Query().Get("author")
	input.BookID = mux.Vars(r)["bookID"]
	input.Authorization = r.Header.Get("authorization")
	if err := json.NewDecoder(r.Body).Decode(&input.TestBook); err != nil{
		return nil, err
	}

	return &input, nil
}
func (i GetBookByIDInput) Validate() error{
	if err := i.TestBook.Validate(); err != nil {
		return err
	}

	return nil
}

type Controller interface {
	GetBookByID(ctx context.Context, input *GetBookByIDInput) (GetBookByIDOutput, error)
}
