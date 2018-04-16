package gendb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/Clever/wag/swagger"
	"github.com/awslabs/goformation/cloudformation"
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
	KeySchema              []cloudformation.AWSDynamoDBTable_KeySchema            `json:"KeySchema,omitempty"`
	GlobalSecondaryIndexes []cloudformation.AWSDynamoDBTable_GlobalSecondaryIndex `json:"GlobalSecondaryIndexes,omitempty"`
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

var tableUsesDateTime = func(config XDBConfig) bool {
	keySchemas := config.DynamoDB.KeySchema
	for _, gsi := range config.DynamoDB.GlobalSecondaryIndexes {
		keySchemas = append(keySchemas, gsi.KeySchema...)
	}
	for _, ks := range keySchemas {
		if config.Schema.Properties[ks.AttributeName].Format == "date-time" {
			return true
		}
	}
	return false
}

// funcMap contains useful functiosn to use in templates
var funcMap = template.FuncMap(map[string]interface{}{
	"tableUsesDateTime": tableUsesDateTime,
	"anyTableUsesDateTime": func(configs []XDBConfig) bool {
		for _, config := range configs {
			if tableUsesDateTime(config) {
				return true
			}
		}
		return false
	},
	"indexHasRangeKey": func(index []cloudformation.AWSDynamoDBTable_KeySchema) bool {
		return len(index) == 2 && index[1].KeyType == "RANGE"
	},
	"indexes": func(config XDBConfig) [][]cloudformation.AWSDynamoDBTable_KeySchema {
		indexes := [][]cloudformation.AWSDynamoDBTable_KeySchema{config.DynamoDB.KeySchema}
		for _, gsi := range config.DynamoDB.GlobalSecondaryIndexes {
			indexes = append(indexes, gsi.KeySchema)
		}
		return indexes
	},
	"unionKeySchemas": func(a, b []cloudformation.AWSDynamoDBTable_KeySchema) []cloudformation.AWSDynamoDBTable_KeySchema {
		ret := []cloudformation.AWSDynamoDBTable_KeySchema{}
		seen := map[string]struct{}{}
		for _, ks := range append(a, b...) {
			if _, ok := seen[ks.AttributeName]; ok {
				continue
			}
			seen[ks.AttributeName] = struct{}{}
			cpy := ks
			ret = append(ret, cpy)
		}
		return ret
	},
	"differenceKeySchemas": func(a, b []cloudformation.AWSDynamoDBTable_KeySchema) []cloudformation.AWSDynamoDBTable_KeySchema {
		ret := []cloudformation.AWSDynamoDBTable_KeySchema{}
		inB := map[string]struct{}{}
		for _, ks := range b {
			inB[ks.AttributeName] = struct{}{}
		}
		for _, ks := range a {
			if _, ok := inB[ks.AttributeName]; ok {
				continue
			}
			cpy := ks
			ret = append(ret, cpy)
		}
		return ret
	},
	"indexName": func(index []cloudformation.AWSDynamoDBTable_KeySchema) string {
		pascalize := generator.FuncMap["pascalize"].(func(string) string)
		if len(index) == 1 {
			return pascalize(index[0].AttributeName)
		} else if len(index) == 2 {
			return fmt.Sprintf("%sAnd%s",
				pascalize(index[0].AttributeName),
				pascalize(index[1].AttributeName),
			)
		} else {
			return ""
		}
	},
	"attributeNames": func(table AWSDynamoDBTable) []string {
		attrnames := map[string]struct{}{}
		for _, ks := range table.KeySchema {
			attrnames[ks.AttributeName] = struct{}{}
		}
		for _, gsi := range table.GlobalSecondaryIndexes {
			for _, ks := range gsi.KeySchema {
				attrnames[ks.AttributeName] = struct{}{}
			}
		}
		attrs := []string{}
		for k := range attrnames {
			attrs = append(attrs, k)
		}
		sort.Strings(attrs)
		return attrs
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
	"exampleValue": func(propertySchema spec.Schema, i int) string {
		if propertySchema.Format == "date-time" {
			return fmt.Sprintf(`mustTime("2018-03-11T15:04:0%d+07:00")`, i)
		} else if len(propertySchema.Type) > 0 {
			if propertySchema.Type[0] == "string" {
				return fmt.Sprintf(`"string%d"`, i)
			} else if propertySchema.Type[0] == "integer" {
				return fmt.Sprintf("%d", i)
			}
		}
		return "unknownType"
	},
	"exampleValuePtr": func(propertySchema spec.Schema, i int) string {
		if propertySchema.Format == "date-time" {
			return fmt.Sprintf(`DateTime(mustTime("2018-03-11T15:04:0%d+07:00"))`, i)
		} else if len(propertySchema.Type) > 0 {
			if propertySchema.Type[0] == "string" {
				return fmt.Sprintf(`String("string%d")`, i)
			} else if propertySchema.Type[0] == "integer" {
				return fmt.Sprintf("Int64(%d)", i)
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
