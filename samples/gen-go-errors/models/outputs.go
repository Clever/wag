package models

func (o ExtendedError) Error() string {
	return o.Msg
}

func (o InternalError) Error() string {
	return o.Msg
}

func (o NotFoundError) Error() string {
	return o.Msg
}
