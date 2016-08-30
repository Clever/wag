package server

import (
	"context"
	"github.com/Clever/wag/generated/models"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_controller.go -package=server

// Controller defines the interface for Swagger Test
type Controller interface {
	// GetBooks returns...
	GetBooks(ctx context.Context, i *models.GetBooksInput) ([]models.Book, error)
	// GetBookByID returns...
	GetBookByID(ctx context.Context, i *models.GetBookByIDInput) (models.GetBookByIDOutput, error)
	// CreateBook returns...
	CreateBook(ctx context.Context, i *models.CreateBookInput) (*models.Book, error)
<<<<<<< 60d0f23f90d41093fff0b1b68ecf1dfd124d5d2e
	HealthCheck(ctx context.Context) error
=======
	// HealthCheck returns...
	HealthCheck(ctx context.Context, i *models.HealthCheckInput) error
>>>>>>> Server comments and linting
}
