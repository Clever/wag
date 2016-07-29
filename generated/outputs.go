package main

	import "github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"

type GetBookByIDOutput interface {
	GetBookByIDStatus() int
}

type GetBookByIDDefaultOutput models.Error

func (o GetBookByIDDefaultOutput) GetBookByIDStatus() int {
	return 200
}

type GetBookByID200Output models.Book

func (o GetBookByID200Output) GetBookByIDStatus() int {
	return 200
}

