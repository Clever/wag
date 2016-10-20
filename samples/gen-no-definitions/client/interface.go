package client

import (
	"context"

	"github.com/Clever/wag/samples/gen-no-definitions/models"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_client.go -package=client

// Client defines the methods available to clients of the swagger-test service.
type Client interface {

	// DeleteBook makes a DELETE request to /books/{id}.
	DeleteBook(ctx context.Context, i *models.DeleteBookInput) error
}
