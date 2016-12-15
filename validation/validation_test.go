package validation

import (
	"testing"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	op.Responses = &spec.Responses{}
	err := validateOp("/books", "GET", &op)
	assert.Error(t, err)
	assert.Equal(t, "paramName for GET /books is a path parameter so it must be required", err.Error())
}

func TestValidateRawRef(t *testing.T) {
	op := spec.Operation{}
	op.Responses = &spec.Responses{}
	op.Responses.StatusCodeResponses = make(map[int]spec.Response)

	response := spec.Response{}
	jsonref, err := jsonreference.New("#/definitions/testref")
	require.NoError(t, err)
	response.Ref = spec.Ref{Ref: jsonref}

	op.Responses.StatusCodeResponses[200] = response
	err = validateResponses("path", "method", &op)
	require.Error(t, err)
	assert.Equal(t, "responses with references should nest the ref in a schema. "+
		"responses.200.$ref = '#/definitions/testref' for path method should be "+
		"responses.200.schema.$ref = '#/definitions/testref'", err.Error())
}
