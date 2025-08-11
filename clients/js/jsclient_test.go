package jsclient

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestFillOutPath(t *testing.T) {
	testSpecs := []struct {
		i string
		o string
	}{
		{"/url", "/url"},
		{"/url/{param}", `/url/" + params.param + "`},
		{"/url/{Param}", `/url/" + params.Param + "`},
		{"/url/{param}/other/{longParam}", `/url/" + params.param + "/other/" + params.longParam + "`},
		{"/url/{param_1}/other/{param_2}", `/url/" + params.param1 + "/other/" + params.param2 + "`},
	}
	for _, spec := range testSpecs {
		assert.Equal(t, spec.o, fillOutPath(spec.i))
	}
}

func TestGetDefaultTimeout(t *testing.T) {
	tests := []struct {
		name     string
		swagger  spec.Swagger
		expected int
	}{
		{
			name: "no extension",
			swagger: spec.Swagger{
				SwaggerProps: spec.SwaggerProps{
					Info: &spec.Info{},
				},
			},
			expected: 5000,
		},
		{
			name: "float64 timeout",
			swagger: spec.Swagger{
				SwaggerProps: spec.SwaggerProps{
					Info: &spec.Info{
						VendorExtensible: spec.VendorExtensible{
							Extensions: map[string]interface{}{
								"x-client-timeout": 15000.0,
							},
						},
					},
				},
			},
			expected: 15000,
		},
		{
			name: "int64 timeout",
			swagger: spec.Swagger{
				SwaggerProps: spec.SwaggerProps{
					Info: &spec.Info{
						VendorExtensible: spec.VendorExtensible{
							Extensions: map[string]interface{}{
								"x-client-timeout": int64(30000),
							},
						},
					},
				},
			},
			expected: 30000,
		},
		{
			name: "int timeout",
			swagger: spec.Swagger{
				SwaggerProps: spec.SwaggerProps{
					Info: &spec.Info{
						VendorExtensible: spec.VendorExtensible{
							Extensions: map[string]interface{}{
								"x-client-timeout": 10000,
							},
						},
					},
				},
			},
			expected: 10000,
		},
		{
			name: "zero timeout falls back to default",
			swagger: spec.Swagger{
				SwaggerProps: spec.SwaggerProps{
					Info: &spec.Info{
						VendorExtensible: spec.VendorExtensible{
							Extensions: map[string]interface{}{
								"x-client-timeout": 0,
							},
						},
					},
				},
			},
			expected: 5000,
		},
		{
			name: "negative timeout falls back to default",
			swagger: spec.Swagger{
				SwaggerProps: spec.SwaggerProps{
					Info: &spec.Info{
						VendorExtensible: spec.VendorExtensible{
							Extensions: map[string]interface{}{
								"x-client-timeout": -1000,
							},
						},
					},
				},
			},
			expected: 5000,
		},
		{
			name: "string timeout falls back to default",
			swagger: spec.Swagger{
				SwaggerProps: spec.SwaggerProps{
					Info: &spec.Info{
						VendorExtensible: spec.VendorExtensible{
							Extensions: map[string]interface{}{
								"x-client-timeout": "15000",
							},
						},
					},
				},
			},
			expected: 5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDefaultTimeout(tt.swagger)
			assert.Equal(t, tt.expected, result)
		})
	}
}
