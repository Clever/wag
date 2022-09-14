package client

import (
	"context"

	"github.com/Clever/wag/samples/v9/gen-go-deprecated/gen-go/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package=client

// Client defines the methods available to clients of the swagger-test service.
type Client interface {
}
