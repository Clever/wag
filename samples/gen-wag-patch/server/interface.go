package server

import (
	"context"

	"github.com/Clever/wag/samples/gen-wag-patch/models"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_controller.go -package=server

// Controller defines the interface for the wag-patch service.
type Controller interface {

	// Wagpatch makes a PATCH request to /wagpatch.
	// Special wag patch type
	Wagpatch(ctx context.Context, i *models.PatchData) (*models.Data, error)
}
