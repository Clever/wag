package jsclient

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/go-openapi/spec"
)

// TypeScriptErrorDeclaration holds attributes for an Error Type Decclaration
type TypeScriptErrorDeclaration struct {
	Name       string
	Properties []string
}

const errorDeclarationTemplate = `class {{.Name}} {
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
		Name:       typeName,
		Properties: properties,
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

		declaration := generatePropertyDeclaration(key, typeForKey, required)
		typeDeclarations = append(typeDeclarations, declaration)
	}

	return typeDeclarations, nil
}

// JS identifiers must begin with $, _, or a letter, and subsequent characters can also be numbers.
// Unfortunately,"letter" here is somewhat complicated and defined by a list of ranges within Unicode.
// For our case, we'll be simplify and be a bit cautious and quote anything that isn't ASCII.
var simpleJSIdentifierRegexp = regexp.MustCompile(`^[$_a-zA-Z][$_a-zA-Z0-9]*$`)

// Given a property key, type, and requirement generate a property declaration
func generatePropertyDeclaration(key string, typeForKey JSType, required bool) string {
	if !simpleJSIdentifierRegexp.MatchString(key) {
		key = `"` + key + `"`
	}

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
