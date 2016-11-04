package swagger

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"
)

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
		def, err := defFromRef(schema.Ref.String())
		if err != nil {
			return "", err
		}
		if includeModels {
			def = "models." + def
		}
		return def, nil
	} else {
		schemaType := schema.Type
		if len(schemaType) != 1 || schemaType[0] != "array" {
			return "", fmt.Errorf("Cannot define complex data types inline. They must be defined in" +
				"the #/definitions section of the swagger yaml")
		}
		items := schema.Items
		if items == nil || items.Schema == nil {
			return "", fmt.Errorf("Cannot define complex data types inline. They must be defined in" +
				"the #/definitions section of the swagger yaml")
		}
		def, err := defFromRef(items.Schema.Ref.String())
		if err != nil {
			return "", err
		}
		if includeModels {
			def = "models." + def
		}
		return "[]" + def, nil
	}
}

func defFromRef(ref string) (string, error) {
	if strings.HasPrefix(ref, "#/definitions/") {
		return ref[len("#/definitions/"):], nil
		// TODO: this isn't right...
	} else if strings.HasPrefix(ref, "#/responses") {
		return ref[len("#/responses/"):], nil
	}
	return "", fmt.Errorf("schema.$ref has undefined reference type \"%s\". "+
		"Must start with #/definitions or #/responses.", ref)
}
