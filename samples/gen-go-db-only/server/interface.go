package server

import (
	"context"
)

//go:generate mockgen -source=$GOFILE -destination=mock_controller.go -package server --build_flags=--mod=mod -imports=models=github.com/Clever/wag/samples/gen-go-db-only/models/v9

// Controller defines the interface for the swagger-test service.
type Controller interface {

	// HealthCheck handles GET requests to /health/check
	//
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	HealthCheck(ctx context.Context) error
}
