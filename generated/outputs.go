package main

import "github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"

type GetBookByIDOutput interface {
	GetBookByIDStatus() int
	// Data is the underlying model object. We know it is json serializable
	GetBookByIDData() interface{}
}

type GetBookByIDError interface {
	error // Extend the error interface
	GetBookByIDStatusCode() int
}

type GetBookByID200Output struct {
	Data models.Book
}

func (o GetBookByID200Output) GetBookByIDData() interface{} {
	return o.Data
}

func (o GetBookByID200Output) GetBookByIDStatus() int {
	return 200
}

type GetBookByID404Output struct{}

func (o GetBookByID404Output) Error() string {
	return "Status Code: " + "404"
}

func (o GetBookByID404Output) GetBookByIDStatusCode() int {
	return 404}

type GetBookByIDDefaultOutput struct{}

func (o GetBookByIDDefaultOutput) Error() string {
	return "Status Code: " + "500"
}

func (o GetBookByIDDefaultOutput) GetBookByIDStatusCode() int {
	return 500}

