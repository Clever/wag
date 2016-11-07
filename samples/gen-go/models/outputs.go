package models

// BadRequest defines a response type.
// Bad Request
type BadRequest Error

// Error returns the message encoded in the error type
func (o BadRequest) Error() string {
	return o.Msg
}

// InternalError defines a response type.
// Internal Error
type InternalError Error

// Error returns the message encoded in the error type
func (o InternalError) Error() string {
	return o.Msg
}

// GetBooks400Output defines the 400 status code response for GetBooks.
type GetBooks400Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBooks400Output) Error() string {
	return "Status Code: 400"
}

// GetBooks500Output defines the 500 status code response for GetBooks.
type GetBooks500Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBooks500Output) Error() string {
	return "Status Code: 500"
}

// CreateBook400Output defines the 400 status code response for CreateBook.
type CreateBook400Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o CreateBook400Output) Error() string {
	return "Status Code: 400"
}

// CreateBook500Output defines the 500 status code response for CreateBook.
type CreateBook500Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o CreateBook500Output) Error() string {
	return "Status Code: 500"
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

// GetBookByID400Output defines the 400 status code response for GetBookByID.
type GetBookByID400Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBookByID400Output) Error() string {
	return "Status Code: 400"
}

// GetBookByID401Output defines the 401 status code response for GetBookByID.
type GetBookByID401Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBookByID401Output) Error() string {
	return "Status Code: 401"
}

// GetBookByID404Output defines the 404 status code response for GetBookByID.
type GetBookByID404Output Error

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBookByID404Output) Error() string {
	return "Status Code: 404"
}

// GetBookByID500Output defines the 500 status code response for GetBookByID.
type GetBookByID500Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBookByID500Output) Error() string {
	return "Status Code: 500"
}

// GetBookByID2400Output defines the 400 status code response for GetBookByID2.
type GetBookByID2400Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBookByID2400Output) Error() string {
	return "Status Code: 400"
}

// GetBookByID2404Output defines the 404 status code response for GetBookByID2.
type GetBookByID2404Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBookByID2404Output) Error() string {
	return "Status Code: 404"
}

// GetBookByID2500Output defines the 500 status code response for GetBookByID2.
type GetBookByID2500Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBookByID2500Output) Error() string {
	return "Status Code: 500"
}

// HealthCheck400Output defines the 400 status code response for HealthCheck.
type HealthCheck400Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o HealthCheck400Output) Error() string {
	return "Status Code: 400"
}

// HealthCheck500Output defines the 500 status code response for HealthCheck.
type HealthCheck500Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o HealthCheck500Output) Error() string {
	return "Status Code: 500"
}
