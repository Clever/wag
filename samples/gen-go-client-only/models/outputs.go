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
