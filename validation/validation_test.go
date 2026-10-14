package validation

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateOpID(t *testing.T) {
	s := spec.Swagger{}
	op := spec.Operation{}
	op.ID = "with spaces"
	err := validateOp(&s, "/books", "GET", &op)
	assert.Error(t, err)
	assert.Equal(t, "The operationId for GET /books must be alphanumeric and start with a letter", err.Error())
}

func TestValidatePathParams(t *testing.T) {
	s := spec.Swagger{}
	op := spec.Operation{}
	op.ID = "op"
	param := spec.Parameter{}
	param.Name = "paramName"
	param.In = "path"
	param.Required = false
	op.Parameters = []spec.Parameter{param}
	op.Responses = &spec.Responses{}
	err := validateOp(&s, "/books", "GET", &op)
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

func TestNestedTypes(t *testing.T) {
	nestedSchema := spec.Schema{}
	nestedSchema.Type = spec.StringOrArray([]string{"object"})
	schema := spec.Schema{}
	schema.Properties = make(map[string]spec.Schema)
	schema.Properties["nested"] = nestedSchema

	definitions := spec.Definitions{}
	definitions["badNested"] = schema

	err := validateDefinitions(definitions)
	require.Error(t, err)
	assert.Equal(t, "badNested cannot have nested object types", err.Error())
}

func TestValidateUnwrappedFetch(t *testing.T) {
	for _, tc := range []struct {
		value   string
		wantErr bool
	}{
		{value: "true"},
		{value: "false"},
		{value: `"true"`, wantErr: true},
	} {
		t.Run(tc.value, func(t *testing.T) {
			d, err := loads.Analyzed(json.RawMessage(fmt.Sprintf(`{
				"swagger": "2.0",
				"info": {"title": "test", "version": "1.0.0", "x-npm-package": "test", "x-unwrapped-fetch": %s},
				"schemes": ["http"],
				"consumes": ["application/json"],
				"produces": ["application/json"],
				"paths": {}
			}`, tc.value)), "")
			require.NoError(t, err)
			err = Validate(*d, true)
			if tc.wantErr {
				require.Error(t, err)
				assert.Equal(t, "'x-unwrapped-fetch' in the 'info' section of the swagger.yml must be true or false", err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
