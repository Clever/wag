package swagger

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"
)

// Interface returns the interface for an operation
func Interface(op *spec.Operation) string {
	successCodes := make([]int, 0)
	for statusCode, _ := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			successCodes = append(successCodes, statusCode)
		}
	}

	capOpID := Capitalize(op.ID)
	if NoSuccessType(op) {
		return fmt.Sprintf("%s(ctx context.Context, i *models.%sInput) error", capOpID, capOpID)
	}

	singleSchema := singleSuccessOutputType(op)
	successType := ""
	if singleSchema != nil {
		var err error
		successType, err = TypeFromSchema(op.Responses.StatusCodeResponses[successCodes[0]].Schema, true)
		if err != nil {
			panic(fmt.Errorf("Could not convert operation to type %s", err))
		}
		// Make non-arrays pointers
		if singleSchema.Ref.String() != "" {
			successType = "*" + successType
		}
	} else {
		successType = fmt.Sprintf("models.%sOutput", capOpID)
	}

	return fmt.Sprintf("%s(ctx context.Context, i *models.%sInput) (%s, error)",
		capOpID, capOpID, successType)
}

// OutputType returns the output type for a given status code of an operation
func OutputType(op *spec.Operation, statusCode int) string {
	singleSuccessSchema := singleSuccessOutputType(op)
	if singleSuccessSchema != nil && statusCode < 400 {
		successType, err := TypeFromSchema(singleSuccessSchema, true)
		if err != nil {
			panic(fmt.Errorf("Could not convert operation to type %s", err))
		}
		return successType
	}
	return fmt.Sprintf("models.%s%dOutput", Capitalize(op.ID), statusCode)
}

// singleSuccessOutputType returns nil if there is more than one success output type for an
// operation. If there is only one, then it returns its name as a string pointer.
func singleSuccessOutputType(op *spec.Operation) *spec.Schema {
	successCodes := make([]int, 0)
	for statusCode, _ := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			successCodes = append(successCodes, statusCode)
		}
	}

	if len(successCodes) == 1 {
		return op.Responses.StatusCodeResponses[successCodes[0]].Schema
	} else {
		return nil
	}
}

// TODO: Add a nice comment!!!
func NoSuccessType(op *spec.Operation) bool {
	successCodes := make([]int, 0)
	for statusCode, _ := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			successCodes = append(successCodes, statusCode)
		}
	}
	if len(successCodes) > 1 {
		return false
	}
	if len(successCodes) == 0 {
		return true
	}
	return op.Responses.StatusCodeResponses[successCodes[0]].Schema == nil
}

// TypeFromSchema returns the type of a Swagger schema as a string. If includeModels is true
// then it returns models.TYPE
func TypeFromSchema(schema *spec.Schema, includeModels bool) (string, error) {
	// We support three types of schemas
	// 1. An empty schema, which we represent by an empty string by default
	// 2. A schema with one element, the $ref key
	// 3. A schema with two elements. One a type with value 'array' and another items field
	// referencing the $ref
	if schema == nil {
		return "string", nil
	} else if schema.Ref.String() != "" {
		ref := schema.Ref.String()
		if !strings.HasPrefix(ref, "#/definitions/") {
			return "", fmt.Errorf("schema.$ref has undefined reference type. Must start with #/definitions")
		}
		def := ref[len("#/definitions/"):]
		if includeModels {
			def = "models." + def
		}
		return def, nil
	} else {
		schemaType := schema.Type
		if len(schemaType) != 1 || schemaType[0] != "array" {
			return "", fmt.Errorf("Two element schemas must have a 'type' field with the value 'array'")
		}
		items := schema.Items
		if items == nil || items.Schema == nil {
			return "", fmt.Errorf("Two element schemas must have a '$ref' field in the 'items' descriptions")
		}
		ref := items.Schema.Ref.String()
		if !strings.HasPrefix(ref, "#/definitions/") {
			return "", fmt.Errorf("schema.$ref has undefined reference type. Must start with #/definitions")
		}
		def := ref[len("#/definitions/"):]
		if includeModels {
			def = "models." + def
		}
		return "[]" + def, nil
	}
}
