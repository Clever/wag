package server

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-errors/models/v9"
)

//go:generate mockgen -source=$GOFILE -destination=mock_controller.go -package server --build_flags=--mod=mod -imports=models=github.com/Clever/wag/samples/gen-go-errors/models/v9

// Controller defines the interface for the swagger-test service.
type Controller interface {

	// GetBook handles GET requests to /books/{id}
	//
	// 200: nil
	// 400: *models.ExtendedError
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetBook(ctx context.Context, i *models.GetBookInput) error
}
