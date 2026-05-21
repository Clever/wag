package gendb

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/require"
)

// compositeConfig builds an XDBConfig with a composite attribute "name_id"
// composed of string "name" and string "id" with separator "@".
func compositeConfig() XDBConfig {
	return XDBConfig{
		SchemaName: "Thing",
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Properties: map[string]spec.Schema{
					"name": {SchemaProps: spec.SchemaProps{Type: []string{"string"}}},
					"id":   {SchemaProps: spec.SchemaProps{Type: []string{"string"}}},
					"date": {SchemaProps: spec.SchemaProps{Type: []string{"string"}, Format: "date-time"}},
				},
			},
		},
		CompositeAttributes: []CompositeAttribute{
			{AttributeName: "name_id", Properties: []string{"name", "id"}, Separator: "@"},
		},
	}
}

func TestCompositeValueFromArray(t *testing.T) {
	config := compositeConfig()
	tests := []struct {
		name            string
		attributeName   string
		sliceIdentifier string
		want            string
	}{
		{
			name:            "string properties",
			attributeName:   "name_id",
			sliceIdentifier: "ms[i]",
			want:            `fmt.Sprintf("%s@%s",ms[i].Name, ms[i].ID)`,
		},
		{
			name:            "non-existent attribute",
			attributeName:   "does_not_exist",
			sliceIdentifier: "ms[i]",
			want:            "not-a-composite-attribute",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compositeValueFromArrayFunc(config, tt.attributeName, tt.sliceIdentifier)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCompositeValue(t *testing.T) {
	config := compositeConfig()
	tests := []struct {
		name         string
		attributeName string
		modelVarName string
		want         string
	}{
		{
			name:          "m prefix",
			attributeName: "name_id",
			modelVarName:  "m",
			want:          `fmt.Sprintf("%s@%s",m.Name, m.ID)`,
		},
		{
			name:          "empty model var produces varname style",
			attributeName: "name_id",
			modelVarName:  "",
			want:          `fmt.Sprintf("%s@%s",name, id)`,
		},
		{
			name:          "custom model var",
			attributeName: "name_id",
			modelVarName:  "item",
			want:          `fmt.Sprintf("%s@%s",item.Name, item.ID)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compositeValueFunc(config, tt.attributeName, tt.modelVarName)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestAttributeToModelValueNotPtr(t *testing.T) {
	config := compositeConfig()
	tests := []struct {
		name          string
		attributeName string
		prefix        string
		want          string
	}{
		{name: "slice prefix", attributeName: "name", prefix: "ms[i].", want: "ms[i].Name"},
		{name: "m prefix", attributeName: "name", prefix: "m.", want: "m.Name"},
		{name: "no prefix returns varname", attributeName: "name", prefix: "", want: "name"},
		{name: "no prefix acronym", attributeName: "id", prefix: "", want: "id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := attributeToModelValueNotPtr(config, tt.attributeName, tt.prefix)
			require.Equal(t, tt.want, got)
		})
	}
}

// compositeValueFromArrayFunc extracts the closure from funcMap for direct testing.
func compositeValueFromArrayFunc(config XDBConfig, attributeName, sliceIdentifier string) string {
	fn := funcMap["compositeValueFromSlice"].(func(XDBConfig, string, string) string)
	return fn(config, attributeName, sliceIdentifier)
}

func compositeValueFunc(config XDBConfig, attributeName, modelVarName string) string {
	fn := funcMap["compositeValue"].(func(XDBConfig, string, string) string)
	return fn(config, attributeName, modelVarName)
}
