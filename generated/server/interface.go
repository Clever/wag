package server

import (
	"github.com/Clever/wag/generated/models"
	"golang.org/x/net/context"
)

type Controller interface {
	GetBooks(ctx context.Context, input *models.GetBooksInput) (models.GetBooksOutput, error)
	GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (models.GetBookByIDOutput, error)
	CreateBook(ctx context.Context, input *models.CreateBookInput) (models.CreateBookOutput, error)
}
