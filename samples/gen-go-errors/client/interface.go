package client

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-errors/models"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_client.go -package=client

// Client defines the methods available to clients of the swagger-test service.
type Client interface {

	// GetBook makes a GET request to /books/{id}.
	GetBook(ctx context.Context, i *models.GetBookInput) error
}
