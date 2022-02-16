package gendb

import (
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/awslabs/goformation/v2/cloudformation/resources"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/go-swagger/go-swagger/generator"
)

// funcMap contains useful functions to use in templates
var funcMap = template.FuncMap(map[string]interface{}{
	"anyTableUsesDateTime": func(configs []XDBConfig) bool {
		for _, config := range configs {
			if tableUsesDateTime(config) {
				return true
			}
		}
		return false
	},
	"tableAllowsScans": tableAllowsScans,
	"anyTableAllowsScans": func(configs []XDBConfig) bool {
		for _, config := range configs {
			if tableAllowsScans(config) {
				return true
			}
		}
		return false
	},
	"indexAllowsScans": func(config XDBConfig, index string) bool {
		for _, ind := range config.AllowSecondaryIndexScan {
			if ind == index {
				return true
			}
		}
		return false
	},
	"indexHasRangeKey": func(index []resources.AWSDynamoDBTable_KeySchema) bool {
		return len(index) == 2 && index[1].KeyType == "RANGE"
	},
	"indexes": func(config XDBConfig) [][]resources.AWSDynamoDBTable_KeySchema {
		indexes := [][]resources.AWSDynamoDBTable_KeySchema{config.DynamoDB.KeySchema}
		for _, gsi := range config.DynamoDB.GlobalSecondaryIndexes {
			indexes = append(indexes, gsi.KeySchema)
		}
		return indexes
	},
	"secondaryIndexes": func(config XDBConfig) [][]resources.AWSDynamoDBTable_KeySchema {
		indexes := [][]resources.AWSDynamoDBTable_KeySchema{}
		for _, gsi := range config.DynamoDB.GlobalSecondaryIndexes {
			indexes = append(indexes, gsi.KeySchema)
		}
		return indexes
	},
	"projectedIndexesWithCompositeAttributes": func(config XDBConfig) [][]resources.AWSDynamoDBTable_KeySchema {
		indexes := [][]resources.AWSDynamoDBTable_KeySchema{}
		for _, gsi := range config.DynamoDB.GlobalSecondaryIndexes {
			if gsi.Projection == nil || gsi.Projection.ProjectionType == "ALL" ||
				!indexContainsCompositeAttribute(config, gsi.KeySchema) {
				continue
			}
			indexes = append(indexes, gsi.KeySchema)
		}
		return indexes
	},
	"unionKeySchemas": func(a, b []resources.AWSDynamoDBTable_KeySchema) []resources.AWSDynamoDBTable_KeySchema {
		ret := []resources.AWSDynamoDBTable_KeySchema{}
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
	"differenceKeySchemas": func(a, b []resources.AWSDynamoDBTable_KeySchema) []resources.AWSDynamoDBTable_KeySchema {
		ret := []resources.AWSDynamoDBTable_KeySchema{}
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
	"indexName": func(index []resources.AWSDynamoDBTable_KeySchema) string {
		generatorFuncs := generator.FuncMapFunc(generator.DefaultLanguageFunc())
		pascalize := generatorFuncs["pascalize"].(func(string) string)
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
	"transactionsEnabled": tableHasTransactionsEnabled,
	"anyTableHasTransactionsEnabled": func(configs []XDBConfig) bool {
		for _, config := range configs {
			if tableHasTransactionsEnabled(config) {
				return true
			}
		}
		return false
	},
	"findCompositeAttribute":                              findCompositeAttribute,
	"indexContainsCompositeAttribute":                     indexContainsCompositeAttribute,
	"indexContainsNonCompositeAttribute":                  indexContainsNonCompositeAttribute,
	"anyIndexRangeKeyContainsSpecifiedCompositeAttribute": anyIndexRangeKeyContainsSpecifiedCompositeAttribute,
	"isComposite": isComposite,
	"stringPropertiesInComposites": func(config XDBConfig) map[string][]string {
		sepToProps := map[string][]string{}
		for _, attr := range config.CompositeAttributes {
			for _, prop := range attr.Properties {
				if goTypeForAttribute(config, prop) == "string" {
					sepToProps[attr.Separator] = append(sepToProps[attr.Separator], prop)
				}
			}
		}
		for k, v := range sepToProps {
			sepToProps[k] = uniq(v)
		}
		return sepToProps
	},
	"getCompositeAttributeAndCompositeAttributeIndexProperties": func(config XDBConfig) []string {
		props := []string{}
		for _, attr := range config.CompositeAttributes {
			for _, prop := range attr.Properties {
				props = append(props, prop)
			}
		}
		for _, gsi := range config.DynamoDB.GlobalSecondaryIndexes {
			if indexContainsCompositeAttribute(config, gsi.KeySchema) {
				for _, ks := range gsi.KeySchema {
					if !isComposite(config, ks.AttributeName) {
						props = append(props, ks.AttributeName)
					}
				}
			}
		}
		props = uniq(props)

		return props
	},
	"sortedKeys": func(m map[string][]string) []string {
		keys := []string{}
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	},
	"compositeValue": func(config XDBConfig, attributeName string, modelVarName string) string {
		ca := findCompositeAttribute(config, attributeName)
		if ca == nil {
			return "not-a-composite-attribute"
		}
		value := `fmt.Sprintf("`
		for i, prop := range ca.Properties {
			goTyp := goTypeForAttribute(config, prop)
			if goTyp == "int64" {
				value += `%%d`
			} else {
				value += `%%s`
			}
			if i != len(ca.Properties)-1 {
				value += ca.Separator
			}
		}
		value += `",`
		for i, prop := range ca.Properties {
			if modelVarName == "m" {
				// hackyaf: usually "m." signifies it could be a pointer
				value += attributeToModelValue(config, prop, modelVarName+".")
			} else if modelVarName == "" {
				value += attributeToModelValueNotPtr(config, prop, "")
			} else {
				value += attributeToModelValueNotPtr(config, prop, modelVarName+".")
			}
			if i != len(ca.Properties)-1 {
				value += `, `
			}
		}
		value += `)`
		return value
	},
	"compositeValuePage": func(config XDBConfig, attributeName string, modelVarName string) string {
		ca := findCompositeAttribute(config, attributeName)
		if ca == nil {
			return "not-a-composite-attribute"
		}
		value := `fmt.Sprintf("`
		for i, prop := range ca.Properties {
			goTyp := goTypeForAttribute(config, prop)
			if goTyp == "int64" {
				value += `%%d`
			} else {
				value += `%%s`
			}
			if i != len(ca.Properties)-1 {
				value += ca.Separator
			}
		}
		value += `",`
		for i, prop := range ca.Properties {
			if modelVarName == "m" {
				// hackyaf: usually "m." signifies it could be a pointer
				value += attributeToModelValue(config, prop, modelVarName+".")
			} else if modelVarName == "" {
				value += attributeToModelValueNotPtr(config, prop, "")
			} else if attributeIsPointer(config, prop) {
				value += attributeToModelValuePtr(config, prop, modelVarName+".")
			} else {
				value += attributeToModelValueNotPtr(config, prop, modelVarName+".")
			}
			if i != len(ca.Properties)-1 {
				value += `, `
			}
		}
		value += `)`
		return value
	},
	"compositeValueFromArray": func(config XDBConfig, attributeName string, modelVarName string, sliceIdentifier string) string {
		ca := findCompositeAttribute(config, attributeName)
		if ca == nil {
			return "not-a-composite-attribute"
		}
		value := `fmt.Sprintf("`
		for i, prop := range ca.Properties {
			goTyp := goTypeForAttribute(config, prop)
			if goTyp == "int64" {
				value += `%%d`
			} else {
				value += `%%s`
			}
			if i != len(ca.Properties)-1 {
				value += ca.Separator
			}
		}
		value += `",`
		for i, prop := range ca.Properties {
			value += sliceIdentifier + "."
			if modelVarName == "m" {
				// hackyaf: usually "m." signifies it could be a pointer
				value += strings.Title(attributeToModelValue(config, prop, modelVarName+"."))
			} else if modelVarName == "" {
				value += strings.Title(attributeToModelValueNotPtr(config, prop, ""))
			} else {
				value += strings.Title(attributeToModelValueNotPtr(config, prop, modelVarName+"."))
			}
			if i != len(ca.Properties)-1 {
				value += `, `
			}
		}
		value += `)`
		return value
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
	"modelAttributeNames": func(config XDBConfig) []string {
		table := config.DynamoDB
		attrnames := map[string]struct{}{}
		for _, ks := range table.KeySchema {
			attrnames[ks.AttributeName] = struct{}{}
		}
		for _, gsi := range table.GlobalSecondaryIndexes {
			for _, ks := range gsi.KeySchema {
				attrnames[ks.AttributeName] = struct{}{}
			}
		}
		for k := range attrnames {
			if ca := findCompositeAttribute(config, k); ca != nil {
				delete(attrnames, k)
				for _, prop := range ca.Properties {
					attrnames[prop] = struct{}{}
				}
			}
		}

		attrs := []string{}
		for k := range attrnames {
			attrs = append(attrs, k)
		}
		sort.Strings(attrs)
		return attrs
	},
	"modelAttributeNamesForIndex": modelAttributeNamesForIndex,
	"modelAttributeNamesForKeyType": func(config XDBConfig, keySchema []resources.AWSDynamoDBTable_KeySchema, keyType string) []string {
		attributeNames := []string{}
		for _, ks := range keySchema {
			if ks.KeyType != keyType {
				continue
			}
			if _, ok := config.Schema.Properties[ks.AttributeName]; ok {
				attributeNames = append(attributeNames, ks.AttributeName)
			} else if ca := findCompositeAttribute(config, ks.AttributeName); ca != nil {
				attributeNames = append(attributeNames, ca.Properties...)
			} else {
				attributeNames = append(attributeNames, "unknownAttributeName")
			}
		}
		return attributeNames
	},
	"configForSchemaName": func(configs []XDBConfig, schemaName string) XDBConfig {
		var config XDBConfig
		for _, c := range configs {
			if c.SchemaName == schemaName {
				config = c
				break
			}
		}
		return config
	},
	"nonPKSecondaryStringProperties": func(config XDBConfig) []string {
		// find attributes in non-primary indexes that are strings.
		// these must be specified when saving a model
		secondaryStringAttributes := []string{}
		for _, gsi := range config.DynamoDB.GlobalSecondaryIndexes {
			if gsi.Projection != nil && gsi.Projection.ProjectionType == "ALL" {
				continue
			}
			for _, attrName := range modelAttributeNamesForIndex(config, gsi.KeySchema) {
				if dynamoDBTypeForAttribute(config, attrName) != "S" {
					// only care about string properties
					continue
				}
				secondaryStringAttributes = append(secondaryStringAttributes, attrName)
			}
		}
		pkAttributes := modelAttributeNamesForIndex(config, config.DynamoDB.KeySchema)
		return difference(secondaryStringAttributes, pkAttributes)
	},
	"goTypeForAttribute":             goTypeForAttribute,
	"dynamoDBTypeForAttribute":       dynamoDBTypeForAttribute,
	"exampleValueForAttribute":       exampleValueForAttribute,
	"exampleValueNotPtrForAttribute": exampleValueNotPtrForAttribute,
	"exampleValuePtrForAttribute":    exampleValuePtrForAttribute,
	"difference":                     difference,
	"pascalizeAndJoin": func(s []string) string {
		ret := ""
		for _, el := range s {
			ret += swag.ToGoName(el)
		}
		return ret
	},
	"attributeIsPointer":          attributeIsPointer,
	"attributeIsEnum":             attributeIsEnum,
	"attributeToModelValue":       attributeToModelValue,
	"attributeToModelValueNotPtr": attributeToModelValueNotPtr,
	"attributeToModelValuePtr":    attributeToModelValuePtr,
	"stringsEqual": func(s1, s2 string) bool {
		return s1 == s2
	},
	"union": union,
})

func contains(el string, arr []string) bool {
	for _, val := range arr {
		if el == val {
			return true
		}
	}
	return false
}

func goTypeForAttribute(config XDBConfig, attributeName string) string {
	if propertySchema, ok := config.Schema.Properties[attributeName]; ok {
		if isDateTimeFormat(propertySchema.Format) {
			return "strfmt.DateTime"
		} else if propertySchema.Format == "byte" {
			return "[]byte"
		} else if len(propertySchema.Type) > 0 {
			if propertySchema.Type[0] == "string" {
				return "string"
			} else if propertySchema.Type[0] == "integer" {
				return "int64"
			}
		} else if ref := propertySchema.Ref.String(); ref != "" {
			refSplit := strings.Split(ref, "/")
			return "models." + swag.ToGoName(refSplit[len(refSplit)-1])
		}
	} else if ca := findCompositeAttribute(config, attributeName); ca != nil {
		// composite attributes must be strings, since they are
		// a concatenation of values
		return "string"
	}
	return "unknownType"
}

func indexContainsCompositeAttribute(config XDBConfig, keySchema []resources.AWSDynamoDBTable_KeySchema) bool {
	for _, ks := range keySchema {
		if ca := findCompositeAttribute(config, ks.AttributeName); ca != nil {
			return true
		}
	}
	return false
}

func indexContainsNonCompositeAttribute(config XDBConfig, keySchema []resources.AWSDynamoDBTable_KeySchema) bool {
	for _, ks := range keySchema {
		if !isComposite(config, ks.AttributeName) {
			return true
		}
	}
	return false
}

func anyIndexRangeKeyContainsSpecifiedCompositeAttribute(config XDBConfig, compAttr string) bool {
	for _, ks := range config.DynamoDB.KeySchema {
		if compAttr == ks.AttributeName && ks.KeyType == "RANGE" {
			return true
		}
	}
	for _, gsi := range config.DynamoDB.GlobalSecondaryIndexes {
		for _, ks := range gsi.KeySchema {
			if compAttr == ks.AttributeName && ks.KeyType == "RANGE" {
				return true
			}
		}
	}

	return false
}

func difference(a, b []string) []string {
	diff := []string{}
	for _, el := range a {
		if !contains(el, b) {
			diff = append(diff, el)
		}
	}
	return diff
}

func dynamoDBTypeForAttribute(config XDBConfig, attributeName string) string {
	if propertySchema, ok := config.Schema.Properties[attributeName]; ok {
		if len(propertySchema.Type) > 0 {
			if propertySchema.Type[0] == "string" {
				if propertySchema.Format == "byte" {
					return "B"
				}
				return "S"
			} else if propertySchema.Type[0] == "integer" {
				return "N"
			}
		} else if propertySchema.Ref.String() != "" {
			return "S"
		}
	} else if ca := findCompositeAttribute(config, attributeName); ca != nil {
		// composite attributes must be strings, since they are
		// a concatenation of values
		return "S"
	}
	return "unknownType"
}

func modelAttributeNamesForIndex(config XDBConfig, keySchema []resources.AWSDynamoDBTable_KeySchema) []string {
	attributeNames := []string{}
	for _, ks := range keySchema {
		if _, ok := config.Schema.Properties[ks.AttributeName]; ok {
			attributeNames = append(attributeNames, ks.AttributeName)
		} else if ca := findCompositeAttribute(config, ks.AttributeName); ca != nil {
			attributeNames = append(attributeNames, ca.Properties...)
		} else {
			attributeNames = append(attributeNames, "unknownAttributeName")
		}
	}
	return attributeNames
}

func isComposite(config XDBConfig, attributeName string) bool {
	if ca := findCompositeAttribute(config, attributeName); ca != nil {
		return true
	}
	return false
}

func union(a, b []string) []string {
	return uniq(append(a, b...))
}

func uniq(arr []string) []string {
	u := map[string]struct{}{}
	for _, el := range arr {
		u[el] = struct{}{}
	}
	unique := []string{}
	for k := range u {
		unique = append(unique, k)
	}
	sort.Strings(unique)
	return unique
}

func attributeIsPointer(config XDBConfig, attributeName string) bool {
	if propertySchema, ok := config.Schema.Properties[attributeName]; ok {
		if propertySchema.Ref.String() != "" {
			// most ref types aren't pointers, for now treat all as non pointers
			return false
		}
	}
	return contains(attributeName, config.Schema.Required) || attributeIsXNullable(config, attributeName)
}

func attributeIsXNullable(config XDBConfig, attributeName string) bool {
	propertySchema, ok := config.Schema.Properties[attributeName]
	if !ok {
		return false
	}
	res, ok := propertySchema.Extensions["x-nullable"]
	if !ok {
		return false
	}

	xNullable, ok := res.(bool)
	if !ok {
		return false
	}
	return xNullable
}

func attributeIsEnum(config XDBConfig, attributeName string) bool {
	if propertySchema, ok := config.Schema.Properties[attributeName]; ok {
		if propertySchema.Ref.String() != "" {
			schemaI, _, err := propertySchema.Ref.GetPointer().Get(config.SwaggerSpec)
			if err != nil {
				panic(err)
			}
			schema, ok := schemaI.(spec.Schema)
			if !ok {
				panic("expected spec.Schema")
			}
			if len(schema.SchemaProps.Enum) > 0 {
				return true
			}
		}
	}
	return false
}

func exampleValueNotPtrForAttribute(config XDBConfig, attributeName string, i int) string {
	if propertySchema, ok := config.Schema.Properties[attributeName]; ok {
		if isDateTimeFormat(propertySchema.Format) {
			return fmt.Sprintf(`mustTime("2018-03-11T15:04:0%d+07:00")`, i)
		} else if len(propertySchema.Type) > 0 {
			if propertySchema.Type[0] == "string" {
				if propertySchema.Format == "byte" {
					return fmt.Sprintf(`[]byte("string%d")`, i)
				}
				return fmt.Sprintf(`"string%d"`, i)
			} else if propertySchema.Type[0] == "integer" {
				return fmt.Sprintf("%d", i)
			}
		} else if propertySchema.Ref.String() != "" {
			schemaI, _, err := propertySchema.Ref.GetPointer().Get(config.SwaggerSpec)
			if err != nil {
				panic(err)
			}
			schema, ok := schemaI.(spec.Schema)
			if !ok {
				panic("expected spec.Schema")
			}
			enumVals := schema.SchemaProps.Enum
			if enumLength := len(enumVals); enumLength > 0 {
				return fmt.Sprintf(`models.Branch%s`, swag.ToGoName(enumVals[(i-1)%enumLength].(string)))
			}
		}
	} else if ca := findCompositeAttribute(config, attributeName); ca != nil {
		// composite attributes must be strings, since they are
		// a concatenation of values
		return fmt.Sprintf(`"string%d"`, i)
	}
	return "unknownType"
}

func exampleValueForAttribute(config XDBConfig, attributeName string, i int) string {
	if attributeIsPointer(config, attributeName) {
		return exampleValuePtrForAttribute(config, attributeName, i)
	}
	return exampleValueNotPtrForAttribute(config, attributeName, i)
}

func exampleValuePtrForAttribute(config XDBConfig, attributeName string, i int) string {
	if propertySchema, ok := config.Schema.Properties[attributeName]; ok {
		if isDateTimeFormat(propertySchema.Format) {
			return fmt.Sprintf(`db.DateTime(mustTime("2018-03-11T15:04:0%d+07:00"))`, i)
		} else if len(propertySchema.Type) > 0 {
			if propertySchema.Type[0] == "string" {
				if propertySchema.Format == "byte" {
					return fmt.Sprintf(`[]byte("string%d")`, i)
				}
				return fmt.Sprintf(`db.String("string%d")`, i)
			} else if propertySchema.Type[0] == "integer" {
				return fmt.Sprintf("db.Int64(%d)", i)
			}
		}
	} else if ca := findCompositeAttribute(config, attributeName); ca != nil {
		// composite attributes must be strings, since they are
		// a concatenation of values
		return fmt.Sprintf(`db.String("string%d")`, i)

	}
	return "unknownType"
}

func attributeToModelValue(config XDBConfig, attributeName string, prefix string) string {
	if attributeIsPointer(config, attributeName) {
		return attributeToModelValuePtr(config, attributeName, prefix)
	}
	return attributeToModelValueNotPtr(config, attributeName, prefix)
}

func attributeToModelValuePtr(config XDBConfig, attributeName string, prefix string) string {
	return "*" + attributeToModelValueNotPtr(config, attributeName, prefix)
}

func attributeToModelValueNotPtr(config XDBConfig, attributeName string, prefix string) string {
	goName := swag.ToGoName(attributeName)
	if prefix == "" {
		goName = strings.ToLower(goName)[0:1] + goName[1:]
		generatorFuncs := generator.FuncMapFunc(generator.DefaultLanguageFunc())
		goName = generatorFuncs["varname"].(func(string) string)(goName)
	}
	return prefix + goName
}

func isDateTimeFormat(s string) bool {
	return contains(s, []string{"datetime", "date-time"})
}

func tableUsesDateTime(config XDBConfig) bool {
	keySchemas := config.DynamoDB.KeySchema
	for _, gsi := range config.DynamoDB.GlobalSecondaryIndexes {
		keySchemas = append(keySchemas, gsi.KeySchema...)
	}
	schemaProperties := []string{}
	for _, ks := range keySchemas {
		if _, ok := config.Schema.Properties[ks.AttributeName]; ok {
			schemaProperties = append(schemaProperties, ks.AttributeName)
		} else if ca := findCompositeAttribute(config, ks.AttributeName); ca != nil {
			schemaProperties = append(schemaProperties, ca.Properties...)
		}
	}
	for _, schemaProp := range schemaProperties {
		if isDateTimeFormat(config.Schema.Properties[schemaProp].Format) {
			return true
		}
	}
	return false
}

func tableAllowsScans(config XDBConfig) bool {
	return config.AllowPrimaryIndexScan || (len(config.AllowSecondaryIndexScan) > 0)
}

func tableHasTransactionsEnabled(config XDBConfig) bool {
	return len(config.EnableTransactions) > 0
}
