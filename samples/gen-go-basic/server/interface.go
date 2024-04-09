package server

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-basic/models/v9"
)

//go:generate mockgen -source=$GOFILE -destination=mock_controller.go -package server --build_flags=--mod=mod -imports=models=github.com/Clever/wag/samples/gen-go-basic/models/v9

// Controller defines the interface for the swagger-test service.
type Controller interface {

	// GetAuthors handles GET requests to /authors
	// Returns response object and the ID of the next page
	// Gets authors
	// 200: *models.AuthorsResponse
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAuthors(ctx context.Context, i *models.GetAuthorsInput) (*models.AuthorsResponse, string, error)

	// GetAuthorsWithPut handles PUT requests to /authors
	// Returns response object and the ID of the next page
	// Gets authors, but needs to use the body so it's a PUT
	// 200: *models.AuthorsResponse
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAuthorsWithPut(ctx context.Context, i *models.GetAuthorsWithPutInput) (*models.AuthorsResponse, string, error)

	// GetBooks handles GET requests to /books
	// Returns response object and the ID of the next page
	// Returns a list of books
	// 200: []models.Book
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetBooks(ctx context.Context, i *models.GetBooksInput) ([]models.Book, int64, error)

	// CreateBook handles POST requests to /books
	// Creates a book
	// 200: *models.Book
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateBook(ctx context.Context, i *models.Book) (*models.Book, error)

	// PutBook handles PUT requests to /books
	// Puts a book
	// 200: *models.Book
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	PutBook(ctx context.Context, i *models.Book) (*models.Book, error)

	// GetBookByID handles GET requests to /books/{book_id}
	// Returns a book
	// 200: *models.Book
	// 400: *models.BadRequest
	// 401: *models.Unathorized
	// 404: *models.Error
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetBookByID(ctx context.Context, i *models.GetBookByIDInput) (*models.Book, error)

	// GetBookByID2 handles GET requests to /books2/{id}
	// Retrieve a book
	// 200: *models.Book
	// 400: *models.BadRequest
	// 404: *models.Error
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetBookByID2(ctx context.Context, id string) (*models.Book, error)

	// HealthCheck handles GET requests to /health/check
	//
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	HealthCheck(ctx context.Context) error

	// LowercaseModelsTest handles POST requests to /lowercaseModelsTest/{pathParam}
	// testing that we can use a lowercase name for a model
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	LowercaseModelsTest(ctx context.Context, i *models.LowercaseModelsTestInput) error
}
