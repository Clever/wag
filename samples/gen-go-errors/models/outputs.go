package models

// NotFound defines a response type.
// Not found
type NotFound Error

// Error returns the message encoded in the error type
func (o NotFound) Error() string {
	return o.Msg
}

// BadRequest defines a response type.
// Bad Request
type BadRequest Error

// Error returns the message encoded in the error type
func (o BadRequest) Error() string {
	return o.Msg
}

// InternalError defines a response type.
// Internal Error
type InternalError ExtendedError

// Error returns the message encoded in the error type
func (o InternalError) Error() string {
	return o.Msg
}

// GetBook400Output defines the 400 status code response for GetBook.
type GetBook400Output ExtendedError

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBook400Output) Error() string {
	return "Status Code: 400"
}

// GetBook404Output defines the 404 status code response for GetBook.
type GetBook404Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBook404Output) Error() string {
	return "Status Code: 404"
}

// GetBook500Output defines the 500 status code response for GetBook.
type GetBook500Output struct{}

// Error returns "Status Code: X". We implemented in to satisfy the error
// interface. For a more descriptive error message see the output type.
func (o GetBook500Output) Error() string {
	return "Status Code: 500"
}
