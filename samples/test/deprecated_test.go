package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Clever/wag/samples/gen-go-deprecated/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-deprecated/server"

	"github.com/stretchr/testify/assert"
)

type DeprecatedController struct{}

func (c *DeprecatedController) Health(ctx context.Context, i *models.HealthInput) error {
	return nil
}

func TestRawHttp(t *testing.T) {
	// This method isn't in the client, but we should still be able to hit it with raw-http
	// for backwards compatability
	s := server.New(&DeprecatedController{}, "")
	testServer := httptest.NewServer(s.Handler)

	// With invalid params
	resp, err := http.Get(testServer.URL + "/v1/health")
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// With valid params
	resp, err = http.Get(testServer.URL + "/v1/health?section=1")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
