package main

import (
	"net/http"
	"golang.org/x/net/context"
)

type GetBookInput struct {
}
func NewGetBookInput(r *http.Request) (*GetBookInput, error) {
	return &GetBookInput{}, nil
}
func (i GetBookInput) Validate() error{
	return nil
}

type Controller interface {
	GetBook(ctx context.Context, input *GetBookInput) (GetBookOutput, error)
}
