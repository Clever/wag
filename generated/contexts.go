package main

import (
	"net/http"
	"golang.org/x/net/context"
)

type GetBookByIDInput struct {
	Author string
	BookID string
	Authorization string
	TestBook Book
}
func NewGetBookByIDInput(r *http.Request) (*GetBookByIDInput, error) {
	return &GetBookByIDInput{}, nil
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
