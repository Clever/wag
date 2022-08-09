package models

import "fmt"

func (o BadRequest) Error() string {
	return o.Message
}

func (o InternalError) Error() string {
	return o.Message
}

func (o NotFound) Error() string {
	return o.Message
}

func (u UnknownResponse) Error() string {
	return fmt.Sprintf("unknown response with status: %d body: %s", u.StatusCode, u.Body)
}
