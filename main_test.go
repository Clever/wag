package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_config_validate(t *testing.T) {
	tests := []struct {
		name    string
		input   config
		output  config
		wantErr bool
	}{
		{
			name: "with server",
			input: config{
				outputPath:    stringPtr("output-path"),
				goPackageName: stringPtr("github.com/Clever/wag/v6/output-path"),
				jsModulePath:  stringPtr("jsModulePath"),
			},
			output: config{
				outputPath:       stringPtr("output-path"),
				goPackageName:    stringPtr("github.com/Clever/wag/v6/output-path"),
				jsModulePath:     stringPtr("jsModulePath"),
				goPackagePath:    "github.com/Clever/wag/output-path",
				generateServer:   true,
				generateGoClient: true,
				generateGoModels: true,
				generateJSClient: true,
			},
		},
		{
			name: "client only",
			input: config{
				clientOnly:    boolPtr(true),
				outputPath:    stringPtr("output-path"),
				goPackageName: stringPtr("goPackageName"),
				jsModulePath:  stringPtr("jsModulePath"),
			},
			output: config{
				clientOnly:       boolPtr(true),
				outputPath:       stringPtr("output-path"),
				goPackageName:    stringPtr("github.com/Clever/wag/v6/output-path"),
				jsModulePath:     stringPtr("jsModulePath"),
				goPackagePath:    "github.com/Clever/wag/output-path",
				generateServer:   false,
				generateGoClient: true,
				generateGoModels: true,
				generateJSClient: true,
			},
		},
		{
			name: "client only go",
			input: config{
				clientOnly:     boolPtr(true),
				clientLanguage: stringPtr("go"),
				outputPath:     stringPtr("output-path"),
				goPackageName:  stringPtr("github.com/Clever/wag/v6/output-path"),
			},
			output: config{
				clientOnly:       boolPtr(true),
				clientLanguage:   stringPtr("go"),
				outputPath:       stringPtr("output-path"),
				goPackageName:    stringPtr("github.com/Clever/wag/v6/output-path"),
				goPackagePath:    "github.com/Clever/wag/output-path",
				generateServer:   false,
				generateGoClient: true,
				generateGoModels: true,
				generateJSClient: false,
			},
		},
		{
			name: "client only js",
			input: config{
				clientOnly:     boolPtr(true),
				clientLanguage: stringPtr("js"),
				outputPath:     stringPtr("output-path"),
				goPackageName:  stringPtr("github.com/Clever/wag/v6/output-path"),
				jsModulePath:   stringPtr("jsModulePath"),
			},
			output: config{
				clientOnly:       boolPtr(true),
				clientLanguage:   stringPtr("js"),
				outputPath:       stringPtr("output-path"),
				goPackageName:    stringPtr("github.com/Clever/wag/v6/output-path"),
				jsModulePath:     stringPtr("jsModulePath"),
				goPackagePath:    "github.com/Clever/wag/output-path",
				generateServer:   false,
				generateGoClient: false,
				generateGoModels: false,
				generateJSClient: true,
			},
		},
		{
			name: "client only invalid language",
			input: config{
				clientOnly:     boolPtr(true),
				clientLanguage: stringPtr("invalid"),
				outputPath:     stringPtr("output-path"),
				goPackageName:  stringPtr("goPackageName"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("config.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// clear paths so they are not diffed
				tt.input.modelsPath = ""
				tt.input.serverPath = ""
				tt.input.goClientPath = ""
				tt.input.jsClientPath = ""

				assert.Equal(t, tt.output, tt.input)
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}
