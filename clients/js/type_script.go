package jsclient

import (
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/go-openapi/spec"
)

// TypeScriptErrorDeclaration holds attributes for an Error Type Decclaration
type TypeScriptErrorDeclaration struct {
	Name          string
	Properties    []string
	Documentation string
}

const errorDeclarationTemplate = `{{.Documentation}}class {{.Name}} {
  {{- range $property := .Properties}}
  {{$property}}
  {{- end}}

  constructor(body: ErrorBody);
}`

var t *template.Template

func init() {
	t = template.Must(template.New("ErrorTypeDeclaration").Parse(errorDeclarationTemplate))
}

// Given an error schema, generate a TypeScript type declaration for the error
func generateErrorDeclaration(schema *spec.Schema, typeName, refPrefix string) (
	string,
	error,
) {
	var builder strings.Builder

	properties, err := generatePropertyDeclarations(schema, refPrefix)
	if err != nil {
		return "", fmt.Errorf("Error generating property declarations: %w", err)
	}

	declaration := TypeScriptErrorDeclaration{
		Name:          typeName,
		Properties:    properties,
		Documentation: documentationForType(*schema),
	}

	err = t.ExecuteTemplate(&builder, "ErrorTypeDeclaration", declaration)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

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

		declaration := documentationForProperty(propertySchema) + generatePropertyDeclaration(key, typeForKey, required)
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

func formatDocumentation(documentation string) string {
	if documentation == "" {
		return ""
	}

	linesOfDocumentation := strings.Split(documentation, "\n")
	if len(linesOfDocumentation) == 1 {
		return fmt.Sprintf(`/** %s */`, linesOfDocumentation[0]) + "\n"
	}

	return fmt.Sprintf(`/**
	%s
	*/
	`, documentation)
}

func documentationForProperty(propertySchema spec.Schema) string {
	return formatDocumentation(propertySchema.Description)
}

func documentationForType(typeSchema spec.Schema) string {
	return formatDocumentation(typeSchema.Description)
}

func documentationForMethod(s spec.Swagger, op *spec.Operation) string {
	documentation := ""
	if op.Summary != "" && op.Description != "" {
		documentation += fmt.Sprintf(`  %s
		
		%s`, op.Summary, op.Description)
	} else {
		documentation += "  " + op.Summary + op.Description
	}
	_, responses := responsesForOperation(s, op)
	for _, response := range responses {
		if response.IsError {
			documentation += fmt.Sprintf("\n@throws %s", response.Name)
		}
	}
	return formatDocumentation(documentation)
}

func documentationForInputParam(param spec.Parameter) string {
	return formatDocumentation(param.Description)
}
