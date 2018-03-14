package gendb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/awslabs/goformation/cloudformation"

	"github.com/Clever/wag/swagger"
	"github.com/go-openapi/spec"
	"github.com/go-swagger/go-swagger/generator"
)

//go:generate $PWD/bin/go-bindata -ignore .*\.go$ -pkg gendb -prefix $PWD/server/gendb/ $PWD/server/gendb/

const xdbExtensionKey = "x-db"

// XDBConfig is the configuration that exists in swagger.yml for auto-generated database code.
type XDBConfig struct {
	// AllowOverwrites sets whether saving an object that already exists should fail.
	AllowOverwrites bool

	// DynamoDB configuration.
	DynamoDB AWSDynamoDBTable

	// Schema and SchemaName that the config was contained within.
	Schema     spec.Schema
	SchemaName string
}

// AWSDynamoDBTable is a subset of clouformation.AWSDynamoDBTable. Currently supported fields:
// -.DynamoDB.KeySchema: configures primary key
// future/todo:
// - GlobalSecondaryIndexes
// - TableName (if you want something other than pascalized model name)
type AWSDynamoDBTable struct {
	KeySchema []cloudformation.AWSDynamoDBTable_KeySchema `json:"KeySchema,omitempty"`
}

// DecodeConfig extracts a db configuration from the schema definition, if one exists.
func DecodeConfig(schemaName string, schema spec.Schema) (*XDBConfig, error) {
	var config *XDBConfig
	for k, v := range schema.VendorExtensible.Extensions {
		switch k {
		case xdbExtensionKey:
			bs, _ := json.Marshal(v)
			if err := json.Unmarshal(bs, &config); err != nil {
				return nil, err
			}
			break
		}
	}
	if config != nil {
		config.SchemaName = schemaName
		config.Schema = schema
		if config.DynamoDB.KeySchema == nil || len(config.DynamoDB.KeySchema) == 0 {
			return nil, fmt.Errorf("x-db DynamoDB config must contain.DynamoDB.KeySchema: %s", schemaName)
		}
	}
	return config, nil
}

var primaryKeyUsesDateTime = func(config XDBConfig) bool {
	for _, ks := range config.DynamoDB.KeySchema {
		if config.Schema.Properties[ks.AttributeName].Format == "date-time" {
			return true
		}
	}
	return false
}

// funcMap contains useful functiosn to use in templates
var funcMap = template.FuncMap(map[string]interface{}{
	"primaryKeyUsesDateTime": primaryKeyUsesDateTime,
	"anyPrimaryKeyUsesDateTime": func(configs []XDBConfig) bool {
		for _, config := range configs {
			if primaryKeyUsesDateTime(config) {
				return true
			}
		}
		return false
	},
	"goType": func(propertySchema spec.Schema) string {
		if propertySchema.Format == "date-time" {
			return "strfmt.DateTime"
		} else if len(propertySchema.Type) > 0 {
			if propertySchema.Type[0] == "string" {
				return "string"
			} else if propertySchema.Type[0] == "integer" {
				return "int64"
			}
		}
		return "unknownType"
	},
	"exampleValue1": func(propertySchema spec.Schema) string {
		if propertySchema.Format == "date-time" {
			return "strfmt.DateTime(time.Unix(1522279646, 0))"
		} else if len(propertySchema.Type) > 0 {
			if propertySchema.Type[0] == "string" {
				return `"string1"`
			} else if propertySchema.Type[0] == "integer" {
				return "1"
			}
		}
		return "unknownType"
	},
	"exampleValue2": func(propertySchema spec.Schema) string {
		if propertySchema.Format == "date-time" {
			return "strfmt.DateTime(time.Unix(2522279646, 0))"
		} else if len(propertySchema.Type) > 0 {
			if propertySchema.Type[0] == "string" {
				return `"string2"`
			} else if propertySchema.Type[0] == "integer" {
				return "2"
			}
		}
		return "unknownType"
	},
})

// GenerateDB generates DB code for schemas annotated with the x-db extension.
func GenerateDB(packageName string, s *spec.Swagger, serviceName string, paths *spec.Paths) error {
	var xdbConfigs []XDBConfig
	for schemaName, schema := range s.Definitions {
		if config, err := DecodeConfig(schemaName, schema); err != nil {
			return err
		} else if config != nil {
			xdbConfigs = append(xdbConfigs, *config)
		}
	}
	if len(xdbConfigs) == 0 {
		return nil
	}
	sort.Slice(xdbConfigs, func(i, j int) bool { return xdbConfigs[i].SchemaName < xdbConfigs[j].SchemaName })

	writeTemplate := func(tmplFilename, outputFilename string, data interface{}) error {
		tmpl, err := template.New("test").
			Funcs(generator.FuncMap).
			Funcs(funcMap).
			Parse(string(MustAsset(tmplFilename)))
		if err != nil {
			return err
		}

		var tmpBuf bytes.Buffer
		err = tmpl.Execute(&tmpBuf, data)
		if err != nil {
			return err
		}

		g := swagger.Generator{PackageName: packageName}
		g.Printf(tmpBuf.String())
		return g.WriteFile(outputFilename)
	}

	type writeTemplateInput struct {
		tmplFilename   string
		outputFilename string
		data           interface{}
	}
	wtis := []writeTemplateInput{
		{
			tmplFilename:   "dynamodb-local.sh.tmpl",
			outputFilename: "server/db/dynamodb/dynamodb-local.sh",
			data:           nil,
		},
		{
			tmplFilename:   "dynamodb.go.tmpl",
			outputFilename: "server/db/dynamodb/dynamodb.go",
			data: map[string]interface{}{
				"PackageName": packageName,
				"XDBConfigs":  xdbConfigs,
			},
		},
		{
			tmplFilename:   "dynamodb_test.go.tmpl",
			outputFilename: "server/db/dynamodb/dynamodb_test.go",
			data: map[string]interface{}{
				"PackageName": packageName,
			},
		},
		{
			tmplFilename:   "interface.go.tmpl",
			outputFilename: "server/db/interface.go",
			data: map[string]interface{}{
				"PackageName": packageName,
				"ServiceName": serviceName,
				"XDBConfigs":  xdbConfigs,
			},
		},
		{
			tmplFilename:   "tests.go.tmpl",
			outputFilename: "server/db/tests/tests.go",
			data: map[string]interface{}{
				"PackageName": packageName,
				"XDBConfigs":  xdbConfigs,
			},
		},
	}
	for _, xdbConfig := range xdbConfigs {
		wtis = append(wtis, writeTemplateInput{
			tmplFilename:   "table.go.tmpl",
			outputFilename: "server/db/dynamodb/" + strings.ToLower(xdbConfig.SchemaName) + ".go",
			data: map[string]interface{}{
				"PackageName": packageName,
				"XDBConfig":   xdbConfig,
			},
		})
	}

	for _, wti := range wtis {
		if err := writeTemplate(wti.tmplFilename, wti.outputFilename, wti.data); err != nil {
			return err
		}
	}

	return nil
}
