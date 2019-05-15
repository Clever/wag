package client

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-errors/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package=client

// Client defines the methods available to clients of the swagger-test service.
type Client interface {

	// GetBook makes a GET request to /books/{id}
	//
	// 200: nil
	// 400: *models.ExtendedError
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetBook(ctx context.Context, i *models.GetBookInput) error
}
