package validation

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestValidateOpID(t *testing.T) {
	op := spec.Operation{}
	op.ID = "with spaces"
	err := validateOp("/books", "GET", &op)
	assert.Error(t, err)
	assert.Equal(t, "The operationId for GET /books must be alphanumeric and start with a letter", err.Error())
}

func TestValidatePathParams(t *testing.T) {
	op := spec.Operation{}
	op.ID = "op"
	param := spec.Parameter{}
	param.Name = "paramName"
	param.In = "path"
	param.Required = false
	op.Parameters = []spec.Parameter{param}
	err := validateOp("/books", "GET", &op)
	assert.Error(t, err)
	assert.Equal(t, "paramName for GET /books is a path parameter so it must be required", err.Error())
}
