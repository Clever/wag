package test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Clever/wag/v6/samples/gen-go-errors/client"
	"github.com/Clever/wag/v6/samples/gen-go-errors/models"
	"github.com/Clever/wag/v6/samples/gen-go-errors/server"

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
	c := client.New(testServer.URL)

	err := c.GetBook(context.Background(), &models.GetBookInput{ID: 404})
	require.Error(t, err)
	assert.IsType(t, &models.NotFound{}, err)
}

func TestOverridenBadRequest(t *testing.T) {
	s := server.New(&ErrorsController{}, "")
	testServer := httptest.NewServer(s.Handler)
	c := client.New(testServer.URL)

	err := c.GetBook(context.Background(), &models.GetBookInput{ID: 50000})
	require.Error(t, err)
	assert.IsType(t, &models.ExtendedError{}, err)
}
