package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Clever/wag/samples/gen-go-errors/client/v9"
	"github.com/Clever/wag/samples/gen-go-errors/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-errors/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ErrorsController struct{}

func (e *ErrorsController) GetBook(ctx context.Context, i *models.GetBookInput) error {
	if i.ID == 404 {
		return models.NotFound{}
	}
	return nil
}

func TestGlobal404(t *testing.T) {
	s := server.New(&ErrorsController{}, "")
	testServer := httptest.NewServer(s.Handler)
	c := client.New(testServer.URL, wcl, &http.DefaultTransport)

	err := c.GetBook(context.Background(), &models.GetBookInput{ID: 404})
	require.Error(t, err)
	assert.IsType(t, &models.NotFound{}, err)
}

func TestOverridenBadRequest(t *testing.T) {
	s := server.New(&ErrorsController{}, "")
	testServer := httptest.NewServer(s.Handler)
	c := client.New(testServer.URL, wcl, &http.DefaultTransport)

	err := c.GetBook(context.Background(), &models.GetBookInput{ID: 50000})
	require.Error(t, err)
	assert.IsType(t, &models.ExtendedError{}, err)
}
