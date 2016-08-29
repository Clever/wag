package swagger

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"
)

// Interface returns the interface for an operation
func Interface(op *spec.Operation) string {
	capOpID := Capitalize(op.ID)
	if NoSuccessType(op) {
		return fmt.Sprintf("%s(ctx context.Context, i *models.%sInput) error", capOpID, capOpID)
	}

	successCodes := SuccessStatusCodes(op)
	successType := ""

	if len(successCodes) == 1 {
		singleSchema := op.Responses.StatusCodeResponses[successCodes[0]].Schema
		var err error
		successType, err = TypeFromSchema(singleSchema, true)
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
	successCodes := SuccessStatusCodes(op)
	if len(successCodes) == 1 && statusCode < 400 {
		successType, err := TypeFromSchema(op.Responses.StatusCodeResponses[successCodes[0]].Schema, true)
		if err != nil {
			panic(fmt.Errorf("Could not convert operation to type %s", err))
		}
		return successType
	}
	return fmt.Sprintf("models.%s%dOutput", Capitalize(op.ID), statusCode)
}

// SUccessStatusCodes returns a slice of all the success status codes for an operation
func SuccessStatusCodes(op *spec.Operation) []int {
	successCodes := make([]int, 0)
	for statusCode, _ := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			successCodes = append(successCodes, statusCode)
		}
	}
	return successCodes
}

// NoSuccessType returns true if the operation has no-success response type. This includes
// either no 200-399 response code or a 200-399 response code without a schema.
func NoSuccessType(op *spec.Operation) bool {
	successCodes := SuccessStatusCodes(op)
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
	// 1. An empty schema, which we represent by the empty struct
	// 2. A schema with one element, the $ref key
	// 3. A schema with two elements. One a type with value 'array' and another items field
	// referencing the $ref
	if schema == nil {
		return "struct{}", nil
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
