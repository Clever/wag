package client

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-wag-patch/models"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_client.go -package=client

// Client defines the methods available to clients of the wag-patch service.
type Client interface {

	// Wagpatch makes a PATCH request to /wagpatch.
	// Special wag patch type
	Wagpatch(ctx context.Context, i *models.PatchData) (*models.Data, error)

	// Wagpatch2 makes a PATCH request to /wagpatch2.
	// Wag patch with another argument
	Wagpatch2(ctx context.Context, i *models.Wagpatch2Input) (*models.Data, error)
}
