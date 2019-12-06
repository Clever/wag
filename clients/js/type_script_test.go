package jsclient

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestGenerateErrorDeclaration(t *testing.T) {
	schema := &spec.Schema{}
	prefix := ""

	t.Run("It generates a type declaration", func(t *testing.T) {
		actual, err := generateErrorDeclaration(schema, "Foo", prefix)

		expected := `class Foo {

  constructor(body: ErrorBody);
}`

		assert.Nil(t, err, "No error occurred")
		assert.Equal(t, expected, actual)
	})

	t.Run("Given some additional properties", func(t *testing.T) {
		properties := map[string]spec.Schema{
			"bar": spec.Schema{},
			"baz": spec.Schema{},
		}
		schema = schema.WithProperties(properties)
		actual, err := generateErrorDeclaration(schema, "Foo", prefix)

		expected := `class Foo {
  bar?: any;
  baz?: any;

  constructor(body: ErrorBody);
}`

		assert.Nil(t, err, "No error occurred")
		assert.Equal(t, expected, actual)

		t.Run("When some of the properties are required", func(t *testing.T) {
			schema = schema.WithRequired("bar")
			actual, err := generateErrorDeclaration(schema, "Foo", prefix)

			expected := `class Foo {
  bar: any;
  baz?: any;

  constructor(body: ErrorBody);
}`

			assert.Nil(t, err, "No error occurred")
			assert.Equal(t, expected, actual)
		})
	})
}
