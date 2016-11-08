package models

func (o ExtendedError) Error() string {
	return o.Msg
}

func (o InternalError) Error() string {
	return o.Msg
}

func (o NotFound) Error() string {
	return o.Msg
}
