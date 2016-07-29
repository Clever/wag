package main

import "golang.org/x/net/context"

type ControllerImpl struct{
}
func (c ControllerImpl) GetBookByID(ctx context.Context, input *GetBookByIDInput) (GetBookByIDOutput, error) {
	// TODO: Implement me!
	// return nil, errors.New("Not implemented")
        return nil, GetBookByID404Output{}
}
