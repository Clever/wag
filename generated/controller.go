package main

import (
	"fmt"

	"golang.org/x/net/context"
)
import "errors"

type ControllerImpl struct {
}

func (c ControllerImpl) GetBookByID(ctx context.Context, input *GetBookByIDInput) (GetBookByIDOutput, error) {

	fmt.Printf("Context value: %s\n", ctx.Value("addedKey"))

	// TODO: Implement me!
	return nil, errors.New("Not implemented")
}
