package models

// DefaultInternalError represents a generic 500 response.
type DefaultInternalError struct {
	Msg string `json:"msg"`
}

func (d DefaultInternalError) Error() string {
	return d.Msg
}

type DefaultBadRequest struct {
	Msg string `json:"msg"`
}

func (d DefaultBadRequest) Error() string {
	return d.Msg
}

type GetBooksError interface {
	error // Extend the error interface
	GetBooksStatusCode() int
}

type GetBookByIDOutput interface {
	GetBookByIDStatus() int
}

type GetBookByID200Output Book

func (o GetBookByID200Output) GetBookByIDStatus() int {
	return 200
}

type GetBookByID204Output string

func (o GetBookByID204Output) GetBookByIDStatus() int {
	return 204
}

type GetBookByIDError interface {
	error // Extend the error interface
	GetBookByIDStatusCode() int
}

type GetBookByID401Output string

func (o GetBookByID401Output) GetBookByIDData() interface{} {
	return o
}

func (o GetBookByID401Output) Error() string {
	// We implement this to satisfy the error interface. This has a generic error message.
	// If the user wants a more details error message they should put it in the output type
	return "Status Code: 401"
}

func (o GetBookByID401Output) GetBookByIDStatusCode() int {
	return 401
}

type GetBookByID404Output Error

func (o GetBookByID404Output) GetBookByIDData() interface{} {
	return o
}

func (o GetBookByID404Output) Error() string {
	// We implement this to satisfy the error interface. This has a generic error message.
	// If the user wants a more details error message they should put it in the output type
	return "Status Code: 404"
}

func (o GetBookByID404Output) GetBookByIDStatusCode() int {
	return 404
}

type CreateBookError interface {
	error // Extend the error interface
	CreateBookStatusCode() int
}

type HealthCheckError interface {
	error // Extend the error interface
	HealthCheckStatusCode() int
}
