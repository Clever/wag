package gendb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/awslabs/goformation/v2/cloudformation/resources"
	"github.com/go-openapi/spec"
	"github.com/go-swagger/go-swagger/generator"

	"github.com/Clever/wag/v9/swagger"
	"github.com/Clever/wag/v9/utils"
)

//go:generate go-bindata -nometadata -ignore .*\.go$ -pkg gendb -prefix ../../server/gendb/ ../../server/gendb/
//go:generate gofmt -w bindata.go

const xdbExtensionKey = "x-db"

// XDBConfig is the configuration that exists in swagger.yml for auto-generated database code.
type XDBConfig struct {
	// AllowOverwrites sets whether saving an object that already exists should fail.
	AllowOverwrites bool

	// CompositeAttributes encodes attributes that are composed of multiple properties in the schema.
	CompositeAttributes []CompositeAttribute

	// AllowBatchWrites determines whether a batch write method should be generated for the table.
	AllowBatchWrites bool

	// AllowPrimaryIndexScan determines whether methods should be generated that scan the primary index.
	AllowPrimaryIndexScan bool

	// AllowSecondaryIndexScan determines whether methods should be generated that scan each of the secondary indexes.
	AllowSecondaryIndexScan []string

	// DynamoDB configuration.
	DynamoDB AWSDynamoDBTable

	// EnableTransactions determines which schemas this schema will be able to perform transactions with. It only needs to be set for one per pair.
	EnableTransactions []string

	// SwaggerSpec, Schema and SchemaName that the config was contained within.
	SwaggerSpec spec.Swagger
	Schema      spec.Schema
	SchemaName  string
}

// CompositeAttribute is an attribute that is composed of multiple properties in the object's schema.
type CompositeAttribute struct {
	AttributeName string
	Properties    []string
	Separator     string
}

// Validate checks that the user enter a valid x-db config.
func (config XDBConfig) Validate(schemaNames []string) error {
	// check that all attribute names show up in the schema or in composite attribute defs.
	for _, ks := range config.DynamoDB.KeySchema {
		if err := config.attributeNameIsDefined(ks.AttributeName); err != nil {
			return err
		}
	}
	for _, gsi := range config.DynamoDB.GlobalSecondaryIndexes {
		for _, ks := range gsi.KeySchema {
			if err := config.attributeNameIsDefined(ks.AttributeName); err != nil {
				return err
			}
		}
	}
	// check that the transaction config is valid i.e. each schema name is valid
	for _, t := range config.EnableTransactions {
		if !contains(t, schemaNames) {
			return fmt.Errorf("invalid transaction config for %s: no matching schema %s", config.SchemaName, t)
		}
	}
	return nil
}

// attributeNameIsDefined checks whether a user has provided an AttributeName that
// is either contained as a property in the swagger schema or defined as a composite
// attribute.
func (config XDBConfig) attributeNameIsDefined(attributeName string) error {
	if _, ok := config.Schema.SchemaProps.Properties[attributeName]; ok {
		return nil
	} else if ca := findCompositeAttribute(config, attributeName); ca != nil {
		return nil
	}
	return fmt.Errorf("unrecognized attribute: '%s'. AttributeNames must match schema properties or be defined as composite attributes", attributeName)
}

// AWSDynamoDBTable is a subset of clouformation.AWSDynamoDBTable. Currently supported fields:
// -.DynamoDB.KeySchema: configures primary key
// future/todo:
// - GlobalSecondaryIndexes
// - TableName (if you want something other than pascalized model name)
type AWSDynamoDBTable struct {
	KeySchema              []resources.AWSDynamoDBTable_KeySchema            `json:"KeySchema,omitempty"`
	GlobalSecondaryIndexes []resources.AWSDynamoDBTable_GlobalSecondaryIndex `json:"GlobalSecondaryIndexes,omitempty"`
	AttributesDefinitions  []resources.AWSDynamoDBTable_AttributeDefinition  `json:"AttributeDefinitions,omitempty"`
}

// DecodeConfig extracts a db configuration from the schema definition, if one exists.
func DecodeConfig(schemaName string, schema spec.Schema, swaggerSpec spec.Swagger) (*XDBConfig, error) {
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
		config.SwaggerSpec = swaggerSpec
		if config.DynamoDB.KeySchema == nil || len(config.DynamoDB.KeySchema) == 0 {
			return nil, fmt.Errorf("x-db DynamoDB config must contain.DynamoDB.KeySchema: %s", schemaName)
		}
	}
	return config, nil
}

func findCompositeAttribute(config XDBConfig, attributeName string) *CompositeAttribute {
	for _, compositeAttr := range config.CompositeAttributes {
		if compositeAttr.AttributeName == attributeName {
			return &compositeAttr
		}
	}
	return nil
}

// GenerateDB generates DB code for schemas annotated with the x-db extension.
func GenerateDB(packageName, basePath, goOutputPath string, withTests bool, s *spec.Swagger, outputPath string) error {
	goOutputPath = strings.TrimPrefix(goOutputPath, ".")
	moduleName, versionSuffix := utils.ExtractModuleNameAndVersionSuffix(packageName, goOutputPath)
	var schemaNames []string
	for schemaName := range s.Definitions {
		schemaNames = append(schemaNames, schemaName)
	}

	var xdbConfigs []XDBConfig
	for schemaName, schema := range s.Definitions {
		if config, err := DecodeConfig(schemaName, schema, *s); err != nil {
			return err
		} else if config != nil {
			if err := config.Validate(schemaNames); err != nil {
				return err
			}
			xdbConfigs = append(xdbConfigs, *config)
		}
	}
	if len(xdbConfigs) == 0 {
		return nil
	}
	sort.Slice(xdbConfigs, func(i, j int) bool { return xdbConfigs[i].SchemaName < xdbConfigs[j].SchemaName })

	writeTemplate := func(tmplFilename, outputFilename string, data interface{}) error {
		tmpl, err := template.New(tmplFilename).
			Funcs(generator.FuncMapFunc(generator.DefaultLanguageFunc())).
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

		g := swagger.Generator{BasePath: basePath}
		g.Print(tmpBuf.String())
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
			outputFilename: path.Join(outputPath, "dynamodb/dynamodb-local.sh"),
			data:           nil,
		},
		{
			tmplFilename:   "dynamodb.go.tmpl",
			outputFilename: path.Join(outputPath, "dynamodb/dynamodb.go"),
			data: map[string]interface{}{
				"PackageName":   packageName,
				"XDBConfigs":    xdbConfigs,
				"OutputPath":    outputPath,
				"ModuleName":    moduleName,
				"VersionSuffix": versionSuffix,
				"GoOutputPath":  goOutputPath,
			},
		},

		{
			tmplFilename:   "interface.go.tmpl",
			outputFilename: path.Join(outputPath, "interface.go"),
			data: map[string]interface{}{
				"PackageName":   packageName,
				"ServiceName":   s.Info.InfoProps.Title,
				"XDBConfigs":    xdbConfigs,
				"ModuleName":    moduleName,
				"VersionSuffix": versionSuffix,
				"GoOutputPath":  goOutputPath,
			},
		},
	}
	for _, xdbConfig := range xdbConfigs {
		wtis = append(wtis, writeTemplateInput{
			tmplFilename:   "table.go.tmpl",
			outputFilename: path.Join(outputPath, "dynamodb", fmt.Sprintf("%v.go", strings.ToLower(xdbConfig.SchemaName))),
			data: map[string]interface{}{
				"PackageName":   packageName,
				"XDBConfig":     xdbConfig,
				"OutputPath":    outputPath,
				"ModuleName":    moduleName,
				"VersionSuffix": versionSuffix,
				"GoOutputPath":  goOutputPath,
			},
		})
	}
	if withTests {
		wtis = append(wtis,
			writeTemplateInput{
				tmplFilename:   "tests.go.tmpl",
				outputFilename: path.Join(outputPath, "tests/tests.go"),
				data: map[string]interface{}{
					"PackageName":   packageName,
					"XDBConfigs":    xdbConfigs,
					"OutputPath":    outputPath,
					"ModuleName":    moduleName,
					"VersionSuffix": versionSuffix,
					"GoOutputPath":  goOutputPath,
				},
			},
			writeTemplateInput{
				tmplFilename:   "dynamodb_test.go.tmpl",
				outputFilename: path.Join(outputPath, "dynamodb/dynamodb_test.go"),
				data: map[string]interface{}{
					"PackageName":   packageName,
					"XDBConfigs":    xdbConfigs,
					"OutputPath":    outputPath,
					"ModuleName":    moduleName,
					"VersionSuffix": versionSuffix,
					"GoOutputPath":  goOutputPath,
				},
			},
		)
	}

	for _, wti := range wtis {
		if err := writeTemplate(wti.tmplFilename, wti.outputFilename, wti.data); err != nil {
			return err
		}
	}

	return nil
}
