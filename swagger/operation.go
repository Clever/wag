package swagger

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"
)

// TOOD: Add a nice comment!!!
func Interface(op *spec.Operation) string {
	successCodes := make([]int, 0)
	for statusCode, _ := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			successCodes = append(successCodes, statusCode)
		}
	}

	capOpID := Capitalize(op.ID)
	singleType := singleSuccessOutputType(op)
	successType := ""
	if singleType != nil {
		successType = "*" + *singleType
	} else {
		successType = fmt.Sprintf("models.%sOutput", capOpID)
	}

	return fmt.Sprintf("%s(ctx context.Context, i *models.%sInput) (%s, error)",
		capOpID, capOpID, successType)
}

// TODO: Add a nice comment!
func OutputType(op *spec.Operation, statusCode int) string {
	singleSuccessType := singleSuccessOutputType(op)
	if singleSuccessType != nil && statusCode < 400 {
		return *singleSuccessType
	}
	return fmt.Sprintf("models.%s%dOutput", Capitalize(op.ID), statusCode)
}

// TODO: Nice comment. Returns non-nil if only one success output
func singleSuccessOutputType(op *spec.Operation) *string {
	successCodes := make([]int, 0)
	for statusCode, _ := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			successCodes = append(successCodes, statusCode)
		}
	}

	successType := ""
	if len(successCodes) == 1 {
		var err error
		successType, err = TypeFromSchema(op.Responses.StatusCodeResponses[successCodes[0]].Schema, true)
		if err != nil {
			panic(fmt.Errorf("Could not convert operation to type %s", err))
		}
		return &successType
	} else {
		return nil
	}
}

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
