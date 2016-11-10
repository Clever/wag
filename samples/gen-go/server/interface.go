package server

import (
	"context"

	"github.com/Clever/wag/samples/gen-go/models"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_controller.go -package=server

// Controller defines the interface for the swagger-test service.
type Controller interface {

	// GetBooks makes a GET request to /books
	// Returns a list of books
	// 200: []models.Book
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetBooks(ctx context.Context, i *models.GetBooksInput) ([]models.Book, error)

	// CreateBook makes a POST request to /books
	// Creates a book
	// 200: *models.Book
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateBook(ctx context.Context, i *models.Book) (*models.Book, error)

	// GetBookByID makes a GET request to /books/{book_id}
	// Returns a book
	// 200: *models.Book
	// 400: *models.BadRequest
	// 401: *models.Unathorized
	// 404: *models.Error
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetBookByID(ctx context.Context, i *models.GetBookByIDInput) (*models.Book, error)

	// GetBookByID2 makes a GET request to /books2/{id}
	// Retrieve a book
	// 200: *models.Book
	// 400: *models.BadRequest
	// 404: *models.Error
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetBookByID2(ctx context.Context, id string) (*models.Book, error)

	// HealthCheck makes a GET request to /health/check
	//
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	HealthCheck(ctx context.Context) error
}
