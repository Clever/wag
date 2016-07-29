package main

	import "github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"

type GetBookByIDOutput interface {
	GetBookByIDStatus() int
	// Data is the underlying model object. We know it is json serializable
	GetBookByIDData() interface{}
}

type GetBookByID200Output struct {
	Data models.Book
}

func (o GetBookByID200Output) GetBookByIDData() interface{} {
	return o.Data
}

func (o GetBookByID200Output) GetBookByIDStatus() int {
	return 0
}

type GetBookByIDDefaultOutput struct{}

func (o GetBookByIDDefaultOutput) Error() string {
	return "Status Code: " + "500"
}

