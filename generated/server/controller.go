package server

import "golang.org/x/net/context"
import "errors"

import "github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"

type ControllerImpl struct{
}
func (c ControllerImpl) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (models.GetBookByIDOutput, error) {
	// TODO: Implement me!
	return nil, errors.New("Not implemented")
}
func (c ControllerImpl) CreateBook(ctx context.Context, input *models.CreateBookInput) (models.CreateBookOutput, error) {
	// TODO: Implement me!
	return nil, errors.New("Not implemented")
}
func (c ControllerImpl) GetBooks(ctx context.Context, input *models.GetBooksInput) (models.GetBooksOutput, error) {
	// TODO: Implement me!
	return nil, errors.New("Not implemented")
}
