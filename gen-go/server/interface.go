package server

import (
	"context"

	"github.com/Clever/wag/gen-go/models"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_controller.go -package=server

// Controller defines the interface for the Swagger Test service.
type Controller interface {
	// GetBooks makes a GET request to /books.
	// Returns a list of books
	GetBooks(ctx context.Context, i *models.GetBooksInput) ([]models.Book, error)

	// GetBookByID makes a GET request to /books/{book_id}.
	// Returns a book
	GetBookByID(ctx context.Context, i *models.GetBookByIDInput) (models.GetBookByIDOutput, error)

	// CreateBook makes a POST request to /books/{book_id}.
	// Creates a book
	CreateBook(ctx context.Context, i *models.Book) (*models.Book, error)

	// GetBookByID2 makes a GET request to /books/{id}.
	// Retrieve a book
	GetBookByID2(ctx context.Context, i *models.GetBookByID2Input) (*models.Book, error)

	// HealthCheck makes a GET request to /health/check.
	// Checks if the service is healthy
	HealthCheck(ctx context.Context) error
}
