package main

import "golang.org/x/net/context"
import "errors"

type ControllerImpl struct{
}
func (c ControllerImpl) GetBook(ctx context.Context, input *GetBookInput) (GetBookOutput, error) {
	// TODO: Implement me!
	return nil, errors.New("Not implemented")
}
