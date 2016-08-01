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
	BookID string
	Authorization string
}

func NewGetBookByIDInput(r *http.Request) (*GetBookByIDInput, error) {
	var input GetBookByIDInput

	input.BookID = mux.Vars(r)["bookID"]
	input.Authorization = r.Header.Get("authorization")

	return &input, nil
}

func (i GetBookByIDInput) Validate() error{
	return nil
}

type CreateBookInput struct {
	NewBook models.Book
}

func NewCreateBookInput(r *http.Request) (*CreateBookInput, error) {
	var input CreateBookInput

	json.NewDecoder(r.Body).Decode(&input.NewBook)

	return &input, nil
}

func (i CreateBookInput) Validate() error{
	if err := i.NewBook.Validate(nil); err != nil {
		return err
	}

	return nil
}

type GetBooksInput struct {
	Author string
}

func NewGetBooksInput(r *http.Request) (*GetBooksInput, error) {
	var input GetBooksInput

	input.Author = r.URL.Query().Get("author")

	return &input, nil
}

func (i GetBooksInput) Validate() error{
	return nil
}

type Controller interface {
	CreateBook(ctx context.Context, input *CreateBookInput) (CreateBookOutput, error)
	GetBookByID(ctx context.Context, input *GetBookByIDInput) (GetBookByIDOutput, error)
	GetBooks(ctx context.Context, input *GetBooksInput) (GetBooksOutput, error)
}
