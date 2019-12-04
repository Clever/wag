package jsclient

import (
	"fmt"
	"sort"

	"github.com/go-openapi/spec"
)

// Given a schema, generates type declarations for each property
func generatePropertyDeclarations(schema *spec.Schema, refPrefix string) (
	[]string,
	error,
) {
	keys := []string{}
	typeDeclarations := []string{}
	requiredFields := extractRequiredFields(schema.Required)

	for key := range schema.Properties {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		propertySchema := schema.Properties[key]
		required := requiredFields[key]

		typeForKey, err := asJSType(&propertySchema, refPrefix)
		if err != nil {
			return []string{}, fmt.Errorf("Error getting type for key '%s': %w", key, err)
		}

		declaration := generatePropertyDeclaration(key, typeForKey, required)
		typeDeclarations = append(typeDeclarations, declaration)
	}

	return typeDeclarations, nil
}

// Given a property key, type, and requirement generate a property declaration
func generatePropertyDeclaration(key string, typeForKey JSType, required bool) string {
	if required {
		return fmt.Sprintf("%s: %s;", key, typeForKey)
	}

	return fmt.Sprintf("%s?: %s;", key, typeForKey)
}

// Given a spec.Required generate a map of required fields
func extractRequiredFields(required []string) map[string]bool {
	requiredFields := map[string]bool{}

	for _, key := range required {
		requiredFields[key] = true
	}

	return requiredFields
}
