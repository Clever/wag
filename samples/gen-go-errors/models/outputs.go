package models

func (o ExtendedError) Error() string {
	return o.Message
}

func (o InternalError) Error() string {
	return o.Message
}

func (o NotFound) Error() string {
	return o.Message
}
