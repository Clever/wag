package models

func (o BadRequest) Error() string {
	return o.Message
}

func (o InternalError) Error() string {
	return o.Message
}
