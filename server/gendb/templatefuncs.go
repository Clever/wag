package gendb

import (
	"fmt"
	"sort"
	"text/template"

	"github.com/awslabs/goformation/cloudformation"
	"github.com/go-openapi/swag"
	"github.com/go-swagger/go-swagger/generator"
)

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
	"secondaryIndexes": func(config XDBConfig) [][]cloudformation.AWSDynamoDBTable_KeySchema {
		indexes := [][]cloudformation.AWSDynamoDBTable_KeySchema{}
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
	"findCompositeAttribute": findCompositeAttribute,
	"indexContainsCompositeAttribute": func(config XDBConfig, keySchema []cloudformation.AWSDynamoDBTable_KeySchema) bool {
		for _, ks := range keySchema {
			if ca := findCompositeAttribute(config, ks.AttributeName); ca != nil {
				return true
			}
		}
		return false
	},
	"isComposite": func(config XDBConfig, attributeName string) bool {
		if ca := findCompositeAttribute(config, attributeName); ca != nil {
			return true
		}
		return false
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
			if modelVarName != "" {
				value += fmt.Sprintf("%s.%s", modelVarName, swag.ToGoName(prop))
			} else {
				value += prop
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
	"modelAttributeNamesForIndex": func(config XDBConfig, keySchema []cloudformation.AWSDynamoDBTable_KeySchema) []string {
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
	},
	"modelAttributeNamesForKeyType": func(config XDBConfig, keySchema []cloudformation.AWSDynamoDBTable_KeySchema, keyType string) []string {
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
	"goTypeForAttribute": goTypeForAttribute,
	"dynamoDBTypeForAttribute": func(config XDBConfig, attributeName string) string {
		if propertySchema, ok := config.Schema.Properties[attributeName]; ok {
			if len(propertySchema.Type) > 0 {
				if propertySchema.Type[0] == "string" {
					return "S"
				} else if propertySchema.Type[0] == "integer" {
					return "N"
				}
			}
		} else if ca := findCompositeAttribute(config, attributeName); ca != nil {
			// composite attributes must be strings, since they are
			// a concatenation of values
			return "S"
		}
		return "unknownType"
	},
	"exampleValueForAttribute": func(config XDBConfig, attributeName string, i int) string {
		if propertySchema, ok := config.Schema.Properties[attributeName]; ok {
			if propertySchema.Format == "date-time" {
				return fmt.Sprintf(`mustTime("2018-03-11T15:04:0%d+07:00")`, i)
			} else if len(propertySchema.Type) > 0 {
				if propertySchema.Type[0] == "string" {
					return fmt.Sprintf(`"string%d"`, i)
				} else if propertySchema.Type[0] == "integer" {
					return fmt.Sprintf("%d", i)
				}
			}
		} else if ca := findCompositeAttribute(config, attributeName); ca != nil {
			// composite attributes must be strings, since they are
			// a concatenation of values
			return fmt.Sprintf(`"string%d"`, i)
		}
		return "unknownType"
	},
	"exampleValuePtrForAttribute": func(config XDBConfig, attributeName string, i int) string {
		if propertySchema, ok := config.Schema.Properties[attributeName]; ok {
			if propertySchema.Format == "date-time" {
				return fmt.Sprintf(`DateTime(mustTime("2018-03-11T15:04:0%d+07:00"))`, i)
			} else if len(propertySchema.Type) > 0 {
				if propertySchema.Type[0] == "string" {
					return fmt.Sprintf(`String("string%d")`, i)
				} else if propertySchema.Type[0] == "integer" {
					return fmt.Sprintf("Int64(%d)", i)
				}
			}
		} else if ca := findCompositeAttribute(config, attributeName); ca != nil {
			// composite attributes must be strings, since they are
			// a concatenation of values
			return fmt.Sprintf(`String("string%d")`, i)

		}
		return "unknownType"
	},
	"difference": func(a, b []string) []string {
		diff := []string{}
		for _, el := range a {
			if !contains(el, b) {
				diff = append(diff, el)
			}
		}
		return diff
	},
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
		if propertySchema.Format == "date-time" {
			return "strfmt.DateTime"
		} else if len(propertySchema.Type) > 0 {
			if propertySchema.Type[0] == "string" {
				return "string"
			} else if propertySchema.Type[0] == "integer" {
				return "int64"
			}
		}
	} else if ca := findCompositeAttribute(config, attributeName); ca != nil {
		// composite attributes must be strings, since they are
		// a concatenation of values
		return "string"
	}
	return "unknownType"
}
