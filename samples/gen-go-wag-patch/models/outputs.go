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
