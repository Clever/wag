package models

// DefaultInternalError represents a generic 500 response.
type DefaultInternalError struct {
	Msg string `json:"msg"`
}

// Error returns the internal error that caused the 500.
func (d DefaultInternalError) Error() string {
	return d.Msg
}

// DefaultBadRequest represents a generic 400 response. It used internally by Swagger as the
// response when a request fails the validation defined in the Swagger yml file.
type DefaultBadRequest struct {
	Msg string `json:"msg"`
}

// Error returns the validation error that caused the 400.
func (d DefaultBadRequest) Error() string {
	return d.Msg
}

// GetBooksError defines the error interface for GetBooks.
type GetBooksError interface {
	error // Extend the error interface
	GetBooksStatusCode() int
}

// GetBookByIDOutput defines the success output interface for GetBookByID.
type GetBookByIDOutput interface {
	GetBookByIDStatusCode() int
}

// GetBookByID200Output defines the 200 status code response for GetBookByID.
type GetBookByID200Output Book

// GetBookByIDStatusCode returns the status code for the operation.
func (o GetBookByID200Output) GetBookByIDStatusCode() int {
	return 200
}

// GetBookByID204Output defines the 204 status code response for GetBookByID.
type GetBookByID204Output struct{}

// GetBookByIDStatusCode returns the status code for the operation.
func (o GetBookByID204Output) GetBookByIDStatusCode() int {
	return 204
}

// GetBookByIDError defines the error interface for GetBookByID.
type GetBookByIDError interface {
	error // Extend the error interface
	GetBookByIDStatusCode() int
}

// GetBookByID401Output defines the 401 status code response for GetBookByID.
type GetBookByID401Output struct{}

// Error returns `Status Code: X`. We implemeted it to satisfy
// the error interface. For a more descriptive error message see
// the output type.
func (o GetBookByID401Output) Error() string {
	return "Status Code: 401"
}

// GetBookByIDStatusCode returns the status code for the operation.
func (o GetBookByID401Output) GetBookByIDStatusCode() int {
	return 401
}

// GetBookByID404Output defines the 404 status code response for GetBookByID.
type GetBookByID404Output Error

// Error returns `Status Code: X`. We implemeted it to satisfy
// the error interface. For a more descriptive error message see
// the output type.
func (o GetBookByID404Output) Error() string {
	return "Status Code: 404"
}

// GetBookByIDStatusCode returns the status code for the operation.
func (o GetBookByID404Output) GetBookByIDStatusCode() int {
	return 404
}

// CreateBookError defines the error interface for CreateBook.
type CreateBookError interface {
	error // Extend the error interface
	CreateBookStatusCode() int
}

// HealthCheckError defines the error interface for HealthCheck.
type HealthCheckError interface {
	error // Extend the error interface
	HealthCheckStatusCode() int
}
