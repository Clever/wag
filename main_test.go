package main

import (
	"testing"

	"github.com/go-openapi/swag"
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
				outputPath:    swag.String("output-path"),
				goPackageName: swag.String("github.com/Clever/wag/v9/output-path"),
				jsModulePath:  swag.String("jsModulePath"),
			},
			output: config{
				outputPath:         swag.String("output-path"),
				goPackageName:      swag.String("github.com/Clever/wag/v9/output-path"),
				jsModulePath:       swag.String("jsModulePath"),
				relativeDynamoPath: swag.String("server/db"),
				generateDynamo:     true,
				generateServer:     true,
				generateTracing:    true,
				generateGoClient:   true,
				generateGoModels:   true,
				generateJSClient:   true,
			},
		},
		{
			name: "client only",
			input: config{
				clientOnly:    swag.Bool(true),
				outputPath:    swag.String("output-path"),
				goPackageName: swag.String("goPackageName"),
				jsModulePath:  swag.String("jsModulePath"),
			},
			output: config{
				clientOnly:       swag.Bool(true),
				outputPath:       swag.String("output-path"),
				goPackageName:    swag.String("github.com/Clever/wag/v9/output-path"),
				jsModulePath:     swag.String("jsModulePath"),
				generateServer:   false,
				generateDynamo:   false,
				generateTracing:  true,
				generateGoClient: true,
				generateGoModels: true,
				generateJSClient: true,
			},
		},
		{
			name: "client only go",
			input: config{
				clientOnly:     swag.Bool(true),
				clientLanguage: swag.String("go"),
				outputPath:     swag.String("output-path"),
				goPackageName:  swag.String("github.com/Clever/wag/v9/output-path"),
			},
			output: config{
				clientOnly:       swag.Bool(true),
				clientLanguage:   swag.String("go"),
				outputPath:       swag.String("output-path"),
				goPackageName:    swag.String("github.com/Clever/wag/v9/output-path"),
				generateServer:   false,
				generateDynamo:   false,
				generateTracing:  true,
				generateGoClient: true,
				generateGoModels: true,
				generateJSClient: false,
			},
		},
		{
			name: "client only js",
			input: config{
				clientOnly:     swag.Bool(true),
				clientLanguage: swag.String("js"),
				outputPath:     swag.String("output-path"),
				goPackageName:  swag.String("github.com/Clever/wag/v9/output-path"),
				jsModulePath:   swag.String("jsModulePath"),
			},
			output: config{
				clientOnly:       swag.Bool(true),
				clientLanguage:   swag.String("js"),
				outputPath:       swag.String("output-path"),
				goPackageName:    swag.String("github.com/Clever/wag/v9/output-path"),
				jsModulePath:     swag.String("jsModulePath"),
				generateServer:   false,
				generateDynamo:   false,
				generateTracing:  false,
				generateGoClient: false,
				generateGoModels: false,
				generateJSClient: true,
			},
		},
		{
			name: "server with js client",
			input: config{
				clientLanguage: swag.String("js"),
				outputPath:     swag.String("output-path"),
				goPackageName:  swag.String("github.com/Clever/wag/v9/output-path"),
				jsModulePath:   swag.String("jsModulePath"),
			},
			output: config{
				clientLanguage:     swag.String("js"),
				outputPath:         swag.String("output-path"),
				goPackageName:      swag.String("github.com/Clever/wag/v9/output-path"),
				jsModulePath:       swag.String("jsModulePath"),
				relativeDynamoPath: swag.String("server/db"),
				generateServer:     true,
				generateDynamo:     true,
				generateTracing:    true,
				generateGoClient:   false,
				generateGoModels:   true,
				generateJSClient:   true,
			},
		},
		{
			name: "server with go client",
			input: config{
				clientLanguage: swag.String("go"),
				outputPath:     swag.String("output-path"),
				goPackageName:  swag.String("github.com/Clever/wag/v9/output-path"),
			},
			output: config{
				clientLanguage:     swag.String("go"),
				outputPath:         swag.String("output-path"),
				goPackageName:      swag.String("github.com/Clever/wag/v9/output-path"),
				relativeDynamoPath: swag.String("server/db"),
				generateServer:     true,
				generateTracing:    true,
				generateGoClient:   true,
				generateDynamo:     true,
				generateGoModels:   true,
				generateJSClient:   false,
			},
		},
		{
			name: "js client no jsModulePath",
			input: config{
				clientLanguage: swag.String("js"),
				outputPath:     swag.String("output-path"),
				goPackageName:  swag.String("github.com/Clever/wag/v9/output-path"),
			},
			wantErr: true,
		},
		{
			name: "client only invalid language",
			input: config{
				clientOnly:     swag.Bool(true),
				clientLanguage: swag.String("invalid"),
				outputPath:     swag.String("output-path"),
				goPackageName:  swag.String("goPackageName"),
			},
			wantErr: true,
		},
		{
			name: "dynamo only custom path",
			input: config{
				dynamoOnly:         swag.Bool(true),
				relativeDynamoPath: swag.String("gen-db/db"),
				outputPath:         swag.String("output-path"),
				goPackageName:      swag.String("github.com/Clever/wag/v9/output-path"),
			},
			output: config{
				outputPath:         swag.String("output-path"),
				goPackageName:      swag.String("github.com/Clever/wag/v9/output-path"),
				relativeDynamoPath: swag.String("gen-db/db"),
				dynamoOnly:         swag.Bool(true),
				generateServer:     false,
				generateTracing:    false,
				generateDynamo:     true,
				generateGoClient:   false,
				generateGoModels:   true,
				generateJSClient:   false,
			},
		},
		{
			name: "dynamo only default path",
			input: config{
				dynamoOnly:    swag.Bool(true),
				outputPath:    swag.String("output-path"),
				goPackageName: swag.String("github.com/Clever/wag/v9/output-path"),
			},
			output: config{
				outputPath:         swag.String("output-path"),
				goPackageName:      swag.String("github.com/Clever/wag/v9/output-path"),
				relativeDynamoPath: swag.String("db"),
				dynamoOnly:         swag.Bool(true),
				generateServer:     false,
				generateTracing:    false,
				generateDynamo:     true,
				generateGoClient:   false,
				generateGoModels:   true,
				generateJSClient:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("config.validate() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				// clear paths so they are not diffed
				tt.input.jsClientPath = ""
				tt.input.dynamoPath = ""
				tt.input.goAbsolutePackagePath = ""

				assert.Equal(t, tt.output, tt.input)
			}
		})
	}
}
