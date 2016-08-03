package server

import (
	"golang.org/x/net/context"
	"github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"
)

type Controller interface {
	GetBooks(ctx context.Context, input *models.GetBooksInput) (models.GetBooksOutput, error)
	GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (models.GetBookByIDOutput, error)
	CreateBook(ctx context.Context, input *models.CreateBookInput) (models.CreateBookOutput, error)
}
