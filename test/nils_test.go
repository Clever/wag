package test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Clever/wag/samples/gen-go-nils/client"
	"github.com/Clever/wag/samples/gen-go-nils/models"
	"github.com/Clever/wag/samples/gen-go-nils/server"
	"github.com/go-openapi/swag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type NilsController struct {
	t *testing.T
}

func (c *NilsController) NilCheck(ctx context.Context, i *models.NilCheckInput) error {
	assert.Nil(c.t, i.Body)
	assert.Equal(c.t, "", i.Header)
	assert.Nil(c.t, i.Query)
	assert.Nil(c.t, i.Array)
	return nil
}

func TestNils(t *testing.T) {
	s := server.New(&NilsController{t: t}, "")
	testServer := httptest.NewServer(s.Handler)
	c := client.New(testServer.URL)

	require.NoError(t, c.NilCheck(context.Background(), &models.NilCheckInput{ID: "a"}))
}

type EmptyController struct {
	t *testing.T
}

func (c *EmptyController) NilCheck(ctx context.Context, i *models.NilCheckInput) error {
	require.NotNil(c.t, i.Body)
	assert.Equal(c.t, "", i.Body.ID)
	assert.Nil(c.t, i.Body.Optional)
	assert.Equal(c.t, "", i.Header)
	require.NotNil(c.t, i.Query)
	assert.Equal(c.t, "", *i.Query)
	// In query params can't distinguish between an empty and nil array
	assert.Nil(c.t, i.Array)
	return nil
}

func TestEmptyStringsAndFields(t *testing.T) {
	s := server.New(&EmptyController{t: t}, "")
	testServer := httptest.NewServer(s.Handler)
	c := client.New(testServer.URL)

	require.NoError(t, c.NilCheck(context.Background(), &models.NilCheckInput{
		ID:     "a",
		Body:   &models.NilFields{ID: ""},
		Header: "",
		Query:  swag.String(""),
		Array:  []string{},
	}))
}
