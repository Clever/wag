package client

import (
	"context"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package client --build_flags=--mod=mod -imports=models=github.com/Clever/wag/samples/gen-go-db/models/v9

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
