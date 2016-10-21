package server

import (
	"context"

	"github.com/Clever/wag/samples/gen-deprecated/models"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_controller.go -package=server

// Controller defines the interface for the swagger-test service.
type Controller interface {

	// GetBook makes a GET request to /books/{id}.
	GetBook(ctx context.Context, i *models.GetBookInput) error
}
