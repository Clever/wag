package main

import "golang.org/x/net/context"
import "errors"
import "github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"

type ControllerImpl struct{
}
func (c ControllerImpl) CreateBook(ctx context.Context, input *CreateBookInput) (CreateBookOutput, error) {
	// TODO: Implement me!
	return nil, errors.New("Not implemented")
}
func (c ControllerImpl) GetBookByID(ctx context.Context, input *GetBookByIDInput) (GetBookByIDOutput, error) {
	// TODO: Implement me!
	return nil, errors.New("Not implemented")
}
func (c ControllerImpl) GetBooks(ctx context.Context, input *GetBooksInput) (GetBooksOutput, error) {
	// TODO: Implement me!
	// return nil, errors.New("Not implemented")
        return GetBooks200Output{Data: []models.Book{models.Book{Name: "Test1"}, models.Book{Name: "Test2"}}}, nil
}
