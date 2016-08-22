package server

import (
	"github.com/Clever/wag/generated/models"
	"golang.org/x/net/context"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_controller.go -package=server

type Controller interface {
	GetBooks(ctx context.Context, input *models.GetBooksInput) (models.GetBooksOutput, error)
	GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (models.GetBookByIDOutput, error)
	CreateBook(ctx context.Context, input *models.CreateBookInput) (models.CreateBookOutput, error)
}
