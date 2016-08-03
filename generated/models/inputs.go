package models

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
)

var _ = json.Marshal

var _ = strconv.FormatInt

type CreateBookInput struct {
	NewBook Book
}

func NewCreateBookInput(r *http.Request) (*CreateBookInput, error) {
	var input CreateBookInput


	return &input, nil
}

func (i CreateBookInput) Validate() error{
	if err := i.NewBook.Validate(nil); err != nil {
		return err
	}

	return nil
}

type GetBookByIDInput struct {
	BookID int64
	Authorization string
}

func NewGetBookByIDInput(r *http.Request) (*GetBookByIDInput, error) {
	var input GetBookByIDInput

	bookIDStr := mux.Vars(r)["bookID"]
	var err error
	input.BookID, err = strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil {
		return nil, err
	}
	authorizationStr := r.Header.Get("authorization")
	input.Authorization = authorizationStr

	return &input, nil
}

func (i GetBookByIDInput) Validate() error{
	return nil
}

type GetBooksInput struct {
	Author string
	Available bool
	MaxPages float64
}

func NewGetBooksInput(r *http.Request) (*GetBooksInput, error) {
	var input GetBooksInput

	authorStr := r.URL.Query().Get("author")
	input.Author = authorStr
	availableStr := r.URL.Query().Get("available")
	var err error
	input.Available, err = strconv.ParseBool(availableStr)
	if err != nil {
		return nil, err
	}
	maxPagesStr := r.URL.Query().Get("maxPages")
	input.MaxPages, err = strconv.ParseFloat(maxPagesStr, 64)
	if err != nil {
		return nil, err
	}

	return &input, nil
}

func (i GetBooksInput) Validate() error{
	return nil
}

