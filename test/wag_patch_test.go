package test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Clever/wag/samples/gen-wag-patch/client"
	"github.com/Clever/wag/samples/gen-wag-patch/models"
	"github.com/Clever/wag/samples/gen-wag-patch/server"
	"github.com/stretchr/testify/assert"
)

type WagPatchController struct {
	Data *models.Data
}

func (w *WagPatchController) Wagpatch(ctx context.Context, i *models.PatchData) (*models.Data, error) {

	if i.ID != nil {
		w.Data.ID = *i.ID
	}
	if i.ArrayField != nil {
		w.Data.ArrayField = i.ArrayField
	}
	if i.Num != nil {
		w.Data.Num = *i.Num
	}
	if i.Nested != nil {
		w.Data.Nested = i.Nested
	}
	return w.Data, nil
}

func TestWagPatch(t *testing.T) {
	s := server.New(&WagPatchController{Data: &models.Data{}}, "")
	testServer := httptest.NewServer(s.Handler)
	c := client.New(testServer.URL)

	// Patch to starting values
	// TODO: Add nested field
	id := "id"
	num := int64(4)
	out, err := c.Wagpatch(context.Background(), &models.PatchData{
		ID:         &id,
		ArrayField: []string{"start"},
		Num:        &num,
		Nested:     &models.NestedData{Name: "test"}})
	assert.NoError(t, err)
	assert.Equal(t, "id", out.ID)
	assert.Equal(t, int64(4), out.Num)
	assert.Equal(t, 1, len(out.ArrayField))
	assert.Equal(t, "test", out.Nested.Name)

	// Setting the values to nil shouldn't do anything
	out, err = c.Wagpatch(context.Background(), &models.PatchData{
		ID:         nil,
		ArrayField: nil,
		Num:        nil,
		Nested:     nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, "id", out.ID)
	assert.Equal(t, int64(4), out.Num)
	assert.Equal(t, 1, len(out.ArrayField))
	assert.Equal(t, "test", out.Nested.Name)

	id = ""
	num = int64(0)
	out, err = c.Wagpatch(context.Background(), &models.PatchData{
		ID: &id,
		//ArrayField: []string{},
		Num:    &num,
		Nested: &models.NestedData{}})
	assert.NoError(t, err)
	assert.Equal(t, "", out.ID)
	assert.Equal(t, int64(0), out.Num)
	// TODO: Fix this in go-swagger
	//assert.Equal(t, 0, len(out.ArrayField))
	assert.Equal(t, "", out.Nested.Name)
}
