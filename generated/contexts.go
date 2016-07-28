package main

import (
	"net/http"
	"golang.org/x/net/context"
)

type GetBookByIDInput struct {
}
func NewGetBookByIDInput(r *http.Request) (*GetBookByIDInput, error) {
	return &GetBookByIDInput{}, nil
}
func (i GetBookByIDInput) Validate() error{
	return nil
}

type Controller interface {
	GetBookByID(ctx context.Context, input *GetBookByIDInput) (GetBookByIDOutput, error)
}
