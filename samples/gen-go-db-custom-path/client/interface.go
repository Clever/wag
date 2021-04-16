package client

import (
	"context"

	"github.com/Clever/wag/v7/samples/gen-go-db-custom-path/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package=client

// Client defines the methods available to clients of the swagger-test service.
type Client interface {

	// HealthCheck makes a GET request to /health/check
	//
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	HealthCheck(ctx context.Context) error
}
