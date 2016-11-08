package models

func (o BadRequest) Error() string {
	return o.Message
}

func (o Error) Error() string {
	return o.Message
}

func (o InternalError) Error() string {
	return o.Message
}

func (o Unathorized) Error() string {
	return o.Message
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
