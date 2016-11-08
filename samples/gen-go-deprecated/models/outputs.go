package models

func (o BadRequest) Error() string {
	return o.Msg
}

func (o InternalError) Error() string {
	return o.Msg
}

func (o NotFound) Error() string {
	return o.Msg
}
