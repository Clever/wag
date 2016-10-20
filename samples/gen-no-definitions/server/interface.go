package server

import (
	"context"

	"github.com/Clever/wag/samples/gen-no-definitions/models"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_controller.go -package=server

// Controller defines the interface for the swagger-test service.
type Controller interface {

	// DeleteBook makes a DELETE request to /books/{id}.
	DeleteBook(ctx context.Context, i *models.DeleteBookInput) error
}
