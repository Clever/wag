package server

import (
	"github.com/Clever/wag/generated/models"
	"golang.org/x/net/context"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_controller.go -package=server

type Controller interface {
	GetBooks(ctx context.Context, i *models.GetBooksInput) (*[]models.Book, error)
	GetBookByID(ctx context.Context, i *models.GetBookByIDInput) (models.GetBookByIDOutput, error)
	CreateBook(ctx context.Context, i *models.CreateBookInput) (*models.Book, error)
}
