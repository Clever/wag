package client

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-deprecated/models/v9"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package client --build_flags=--mod=mod

// Client defines the methods available to clients of the swagger-test service.
type Client interface {
}
