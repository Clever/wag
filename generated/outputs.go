package generated

import "github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"

type GetBooksOutput interface {
	GetBooksStatus() int
	// Data is the underlying model object. We know it is json serializable
	GetBooksData() interface{}
}

type GetBooksError interface {
	error // Extend the error interface
	GetBooksStatusCode() int
}

<<<<<<< a0005dbc4f1c70496590515a9cdd8b111d9685bb
type GetBookByIDDefaultOutput struct{}

func (o GetBookByIDDefaultOutput) Error() string {
	return "Status Code: " + "500"
}

func (o GetBookByIDDefaultOutput) GetBookByIDStatusCode() int {
	return 500}

type GetBookByID200Output struct {
	Data models.Book
=======
type GetBooks200Output struct {
	Data []models.Book
>>>>>>> Add other input types to code generation
}

func (o GetBooks200Output) GetBooksData() interface{} {
	return o.Data
}

func (o GetBooks200Output) GetBooksStatus() int {
	return 200
}

type GetBooksDefaultOutput struct{}

func (o GetBooksDefaultOutput) Error() string {
	return "Status Code: " + "500"
}

func (o GetBooksDefaultOutput) GetBooksStatusCode() int {
	return 500}

type GetBookByIDOutput interface {
	GetBookByIDStatus() int
	// Data is the underlying model object. We know it is json serializable
	GetBookByIDData() interface{}
}

type GetBookByIDError interface {
	error // Extend the error interface
	GetBookByIDStatusCode() int
}

type GetBookByID404Output struct{}

func (o GetBookByID404Output) Error() string {
	return "Status Code: " + "404"
}

func (o GetBookByID404Output) GetBookByIDStatusCode() int {
	return 404}

<<<<<<< a0005dbc4f1c70496590515a9cdd8b111d9685bb
=======
type GetBookByIDDefaultOutput struct{}

func (o GetBookByIDDefaultOutput) Error() string {
	return "Status Code: " + "500"
}

func (o GetBookByIDDefaultOutput) GetBookByIDStatusCode() int {
	return 500}

type GetBookByID200Output struct {
	Data models.Book
}

func (o GetBookByID200Output) GetBookByIDData() interface{} {
	return o.Data
}

func (o GetBookByID200Output) GetBookByIDStatus() int {
	return 200
}

>>>>>>> Add other input types to code generation
type CreateBookOutput interface {
	CreateBookStatus() int
	// Data is the underlying model object. We know it is json serializable
	CreateBookData() interface{}
}

type CreateBookError interface {
	error // Extend the error interface
	CreateBookStatusCode() int
}

type CreateBook200Output struct {
	Data models.Book
}

func (o CreateBook200Output) CreateBookData() interface{} {
	return o.Data
}

func (o CreateBook200Output) CreateBookStatus() int {
	return 200
}

type CreateBookDefaultOutput struct{}

func (o CreateBookDefaultOutput) Error() string {
	return "Status Code: " + "500"
}

func (o CreateBookDefaultOutput) CreateBookStatusCode() int {
	return 500}

