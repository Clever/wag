package swagger

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"
)

// TypeFromSchema returns the type of a Swagger schema as a string. If includeModels is true
// then it returns models.TYPE
func TypeFromSchema(schema *spec.Schema, includeModels bool) (string, error) {
	// We support one of two schemas:
	// 1. A schema with one element, the $ref key
	// 2. A schema with two elements. One a type with value 'array' and another items field
	// referencing the $ref
	if schema == nil {
		return "", nil
	} else if schema.Ref.String() != "" {
		def, err := DefFromRef(schema.Ref.String())
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
			return "", fmt.Errorf("Cannot define complex data types inline. They must be defined in " +
				"the #/definitions section of the swagger yaml.")
		}
		items := schema.Items
		if items == nil || items.Schema == nil || items.Schema.Ref.String() == "" {
			return "", fmt.Errorf("Cannot define complex data types inline. They must be defined in " +
				"the #/definitions section of the swagger yaml.")
		}
		def, err := DefFromRef(items.Schema.Ref.String())
		if err != nil {
			return "", err
		}
		if includeModels {
			def = "models." + def
		}
		return "[]" + def, nil
	}
}

// DefFromRef returns the schema definition given the reference
func DefFromRef(ref string) (string, error) {
	if strings.HasPrefix(ref, "#/definitions/") {
		return ref[len("#/definitions/"):], nil
	}
	return "", fmt.Errorf("schema.$ref has undefined reference type \"%s\". "+
		"Must start with #/definitions or #/responses.", ref)
}
