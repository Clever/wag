package validation

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestValidateOpID(t *testing.T) {
	op := spec.Operation{}
	op.ID = "with spaces"
	err := validateOp(&op)
	assert.Error(t, err)
	assert.Equal(t, "OperationIDs must be alphanumeric and start with a letter", err.Error())
}

func TestValidatePathParams(t *testing.T) {
	op := spec.Operation{}
	param := spec.Parameter{}
	param.In = "path"
	param.Required = false
	op.Parameters = []spec.Parameter{param}
	err := validateOp(&op)
	assert.Error(t, err)
	assert.Equal(t, "OperationIDs must be alphanumeric and start with a letter", err.Error())
}
