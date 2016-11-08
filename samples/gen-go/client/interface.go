package client

import (
	"context"

	"github.com/Clever/wag/samples/gen-go/models"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_client.go -package=client

// Client defines the methods available to clients of the swagger-test service.
type Client interface {

	// GetBooks makes a GET request to /books.
	// Returns a list of books
	GetBooks(ctx context.Context, i *models.GetBooksInput) ([]models.Book, error)

	// CreateBook makes a POST request to /books.
	// Creates a book
	CreateBook(ctx context.Context, i *models.Book) (*models.Book, error)

	// GetBookByID makes a GET request to /books/{book_id}.
	// Returns a book
	GetBookByID(ctx context.Context, i *models.GetBookByIDInput) (*models.Book, error)

	// GetBookByID2 makes a GET request to /books2/{id}.
	// Retrieve a book
	GetBookByID2(ctx context.Context, id string) (*models.Book, error)

	// HealthCheck makes a GET request to /health/check.
	HealthCheck(ctx context.Context) error
}
