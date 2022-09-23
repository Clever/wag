package client

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-basic/models/v9"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package client --build_flags=--mod=mod

// Client defines the methods available to clients of the wag/samples service.
type Client interface {

	// GetAuthors makes a GET request to /authors
	// Gets authors
	// 200: *models.AuthorsResponse
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAuthors(ctx context.Context, i *models.GetAuthorsInput) (*models.AuthorsResponse, error)

	NewGetAuthorsIter(ctx context.Context, i *models.GetAuthorsInput) (GetAuthorsIter, error)

	// GetAuthorsWithPut makes a PUT request to /authors
	// Gets authors, but needs to use the body so it's a PUT
	// 200: *models.AuthorsResponse
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetAuthorsWithPut(ctx context.Context, i *models.GetAuthorsWithPutInput) (*models.AuthorsResponse, error)

	NewGetAuthorsWithPutIter(ctx context.Context, i *models.GetAuthorsWithPutInput) (GetAuthorsWithPutIter, error)

	// GetBooks makes a GET request to /books
	// Returns a list of books
	// 200: []models.Book
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetBooks(ctx context.Context, i *models.GetBooksInput) ([]models.Book, error)

	NewGetBooksIter(ctx context.Context, i *models.GetBooksInput) (GetBooksIter, error)

	// CreateBook makes a POST request to /books
	// Creates a book
	// 200: *models.Book
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateBook(ctx context.Context, i *models.Book) (*models.Book, error)

	// PutBook makes a PUT request to /books
	// Puts a book
	// 200: *models.Book
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	PutBook(ctx context.Context, i *models.Book) (*models.Book, error)

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

// GetAuthorsIter defines the methods available on GetAuthors iterators.
type GetAuthorsIter interface {
	Next(*models.Author) bool
	Err() error
}

// GetAuthorsWithPutIter defines the methods available on GetAuthorsWithPut iterators.
type GetAuthorsWithPutIter interface {
	Next(*models.Author) bool
	Err() error
}

// GetBooksIter defines the methods available on GetBooks iterators.
type GetBooksIter interface {
	Next(*models.Book) bool
	Err() error
}
