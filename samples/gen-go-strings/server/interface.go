package server

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-strings/models/v9"
)

//go:generate mockgen -source=$GOFILE -destination=mock_controller.go -package server --build_flags=--mod=mod -imports=models=github.com/Clever/wag/samples/gen-go-strings/models/v9

// Controller defines the interface for the nil-test service.
type Controller interface {

	// GetDistricts handles POST requests to /check
	//
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetDistricts(ctx context.Context, i *models.GetDistrictsInput) error
}
