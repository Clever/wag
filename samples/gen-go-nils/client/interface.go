package client

import (
	"context"

	"github.com/Clever/wag/v5/samples/gen-go-nils/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package=client

// Client defines the methods available to clients of the nil-test service.
type Client interface {

	// NilCheck makes a POST request to /check/{id}
	// Nil check tests
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	NilCheck(ctx context.Context, i *models.NilCheckInput) error
}
