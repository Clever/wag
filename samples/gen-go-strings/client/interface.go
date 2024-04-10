package client

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-strings/models/v9"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package client --build_flags=--mod=mod -imports=models=github.com/Clever/wag/samples/gen-go-strings/models/v9

// Client defines the methods available to clients of the nil-test service.
type Client interface {

	// GetDistricts makes a POST request to /check
	//
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetDistricts(ctx context.Context, i *models.GetDistrictsInput) error
}
