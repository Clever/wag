package test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Clever/wag/samples/gen-no-definitions/client"
	"github.com/Clever/wag/samples/gen-no-definitions/models"
	"github.com/Clever/wag/samples/gen-no-definitions/server"
	"github.com/stretchr/testify/assert"
)

type NoDefinitionsController struct{}

func (m *NoDefinitionsController) DeleteBook(ctx context.Context, i *models.DeleteBookInput) error {
	if i.ID == 404 {
		return models.DeleteBook404Output{}
	}
	return nil
}

func TestNoDefinitions(t *testing.T) {
	s := server.New(&NoDefinitionsController{}, "")
	testServer := httptest.NewServer(s.Handler)
	c := client.New(testServer.URL)

	// Success
	err := c.DeleteBook(context.Background(), &models.DeleteBookInput{})
	assert.NoError(t, err)

	// Failure
	err = c.DeleteBook(context.Background(), &models.DeleteBookInput{ID: 404})
	assert.Error(t, err)
	assert.IsType(t, models.DeleteBook404Output{}, err)
}
