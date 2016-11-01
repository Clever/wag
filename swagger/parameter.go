package swagger

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Clever/wag/utils"
	"github.com/go-openapi/spec"
)

// This code defines all the operations on parameter objects. The swagger spec for parameters
// is defined here: http://swagger.io/specification/#parameterObject. Note that currently we don't
// support array, file, or string.binary data types and that the schema logic isn't currently defined
// in this file.
//
// There are four common operations on parameter objects and we have one function for each:
// 1. Param -> Go Type
// 2. Go Type -> string
// 3. String -> Go type
// 4. Validation logic

// BaseParamToStringCode returns the code for a function that joins a string
// array by a known format (csv, ssv, tsv, pipes).
func BaseParamToStringCode() string {
	return `
// JoinByFormat joins a string array by a known format:
//	 csv: comma separated value (default)
//	 ssv: space separated value
//	 tsv: tab separated value
//	 pipes: pipe (|) separated value
func JoinByFormat(data []string, format string) string {
	if len(data) == 0 {
		return ""
	}
	var sep string
	switch format {
	case "ssv":
		sep = " "
	case "tsv":
		sep = "\t"
	case "pipes":
		sep = "|"
	default:
		sep = ","
	}
	return strings.Join(data, sep)
}
`
}

// ParamToStringCode returns a function that converts a Parameter into the code to convert
// it into a string (for serialization). For example, a integer named 'Size' becomes
// `strconv.FormatInt(i.Size, 10)`
func ParamToStringCode(param spec.Parameter) string {
	valToSet := fmt.Sprintf("i.%s", StructParamName(param))
	if !param.Required && param.Type != "array" {
		valToSet = "*" + valToSet
	}
	switch param.Type {
	case "string":
		switch param.Format {
		case "byte":
			return fmt.Sprintf("string(%s)", valToSet)
		case "date":
			return fmt.Sprintf("(%s).String()", valToSet)
		case "date-time":
			return fmt.Sprintf("(%s).String()", valToSet)
		case "mongo-id":
			return fmt.Sprintf("%s", valToSet)
		case "":
			return fmt.Sprintf("%s", valToSet)
		default:
			panic(fmt.Errorf("unsupported string format %s", param.Format))
		}
	case "integer":
		if param.Format == "int32" {
			return fmt.Sprintf("strconv.FormatInt(int64(%s), 10)", valToSet)
		}
		return fmt.Sprintf("strconv.FormatInt(%s, 10)", valToSet)
	case "number":
		if param.Format == "float" {
			return fmt.Sprintf("strconv.FormatFloat(float64(%s), 'E', -1, 32)", valToSet)
		}
		return fmt.Sprintf("strconv.FormatFloat(%s, 'E', -1, 64)", valToSet)
	case "boolean":
		return fmt.Sprintf("strconv.FormatBool(%s)", valToSet)
	case "array":
		if param.Items.Type != "string" {
			panic(fmt.Errorf("Array parameters must have string sub-types"))
		}
		return fmt.Sprintf("JoinByFormat(%s, \"%s\")", valToSet, param.Format)
	default:
		// Theoretically should have validated before getting here
		panic(fmt.Errorf("unsupported parameter type %s", param.Type))
	}
}

// ParamToType converts a parameter into its Go type. If withPointer is true
// it will make the variable a pointer if it's not required and not an array type
func ParamToType(param spec.Parameter, withPointer bool) (string, error) {
	var typeName string
	switch param.Type {
	case "string":
		switch param.Format {
		case "byte":
			typeName = "strfmt.Base64"
		case "date":
			typeName = "strfmt.Date"
		case "date-time":
			typeName = "strfmt.DateTime"
		case "mongo-id":
			typeName = "string"
		case "":
			typeName = "string"
		default:
			return "", fmt.Errorf("unsupported string format \"%s\"", param.Format)
		}
	case "integer":
		if param.Format == "int32" {
			typeName = "int32"
		} else {
			typeName = "int64"
		}
	case "boolean":
		typeName = "bool"
	case "number":
		if param.Format == "float" {
			typeName = "float32"
		} else {
			typeName = "float64"
		}
	case "array":
		if param.Items.Type != "string" {
			return "", fmt.Errorf("array parameters must have string sub-types")
		}
		typeName = "[]string"
	default:
		// Note. We don't support 'array' or 'file' types even though they're in the
		// Swagger spec.
		return "", fmt.Errorf("unsupported param type: \"%s\"", param.Type)
	}
	if withPointer && !param.Required && param.Type != "array" {
		typeName = "*" + typeName
	}
	return typeName, nil
}

// BaseStringToTypeCode is the helper code from base string to type
func BaseStringToTypeCode() string {
	return `

var formats = strfmt.Default
var _ = formats

// convertBase64 takes in a string and returns a strfmt.Base64 if the input
// is valid base64 and an error otherwise.
func convertBase64(input string) (strfmt.Base64, error) {
	temp, err := formats.Parse("byte", input)
	if err != nil {
		return strfmt.Base64{}, err
	}
	return *temp.(*strfmt.Base64), nil
}

// convertDateTime takes in a string and returns a strfmt.DateTime if the input
// is a valid DateTime and an error otherwise.
func convertDateTime(input string) (strfmt.DateTime, error) {
	temp, err := formats.Parse("date-time", input)
	if err != nil {
		return strfmt.DateTime{}, err
	}
	return *temp.(*strfmt.DateTime), nil
}

// convertDate takes in a string and returns a strfmt.Date if the input
// is a valid Date and an error otherwise.
func convertDate(input string) (strfmt.Date, error) {
	temp, err := formats.Parse("date", input)
	if err != nil {
		return strfmt.Date{}, err
	}
	return *temp.(*strfmt.Date), nil
}
`
}

// StringToTypeCode returns the code to convert a string into the type. For example,
// for an int64 it generates swag.ConvertFloat64(sizeStr). It returns an array because sometimes
// this takes multiple lines of code
func StringToTypeCode(strField string, param spec.Parameter) (string, error) {

	switch param.Type {
	case "integer":
		if param.Format == "int32" {
			return fmt.Sprintf("swag.ConvertInt32(%s)", strField), nil
		}
		return fmt.Sprintf("swag.ConvertInt64(%s)", strField), nil
	case "number":
		if param.Format == "float" {
			return fmt.Sprintf("swag.ConvertFloat32(%s)", strField), nil
		}
		return fmt.Sprintf("swag.ConvertFloat64(%s)", strField), nil
	case "boolean":
		return fmt.Sprintf("strconv.ParseBool(%s)", strField), nil
	case "string":
		switch param.Format {
		case "byte":
			return fmt.Sprintf("convertBase64(%s)", strField), nil
		case "mongo-id":
			return fmt.Sprintf("%s, error(nil)", strField), nil
		case "":
			return fmt.Sprintf("%s, error(nil)", strField), nil
		case "date":
			return fmt.Sprintf("convertDate(%s)", strField), nil
		case "date-time":
			return fmt.Sprintf("convertDateTime(%s)", strField), nil
		default:
			return "", fmt.Errorf("param format %s not supported", param.Format)
		}
	case "array":
		if param.Items.Type != "string" {
			return "", fmt.Errorf("array parameters must have string sub-types")
		}
		return fmt.Sprintf("swag.SplitByFormat(%s, \"%s\"), error(nil)", strField, param.Format), nil
	default:
		return "", fmt.Errorf("param type %s not supported", param.Type)
	}
}

// ParamToValidationCode takes in a param and returns a list of parameter validation
// functions, each of which have a single return value, error
func ParamToValidationCode(param spec.Parameter) ([]string, error) {
	var validations []string
	if param.Type == "string" {
		if param.MaxLength != nil {
			validations = append(validations, fmt.Sprintf("validate.MaxLength(\"%s\", \"%s\", string(%s), %d)",
				param.Name, param.In, accessString(param), *param.MaxLength))
		}
		if param.MinLength != nil {
			validations = append(validations, fmt.Sprintf("validate.MinLength(\"%s\", \"%s\", string(%s), %d)",
				param.Name, param.In, accessString(param), *param.MaxLength))
		}
		if param.Pattern != "" {
			validations = append(validations, fmt.Sprintf("validate.Pattern(\"%s\", \"%s\", string(%s), \"%s\")",
				param.Name, param.In, accessString(param), param.Pattern))
		}
		if param.Format != "" {
			validations = append(validations, fmt.Sprintf("validate.FormatOf(\"%s\", \"%s\", \"%s\", %s, strfmt.Default)",
				param.Name, param.In, param.Format, ParamToStringCode(param)))
		}
		if len(param.Enum) != 0 {
			strEnum := []string{}
			for _, enum := range param.Enum {
				strEnum = append(strEnum, enum.(string))
			}
			validations = append(validations, fmt.Sprintf("validate.Enum(\"%s\", \"%s\", %s, %s)",
				param.Name, param.In, accessString(param), fmt.Sprintf("[]interface{}{\"%s\"}", strings.Join(strEnum, "\",\""))))
		}
	} else if param.Type == "integer" {
		if param.Maximum != nil {
			validations = append(validations, fmt.Sprintf("validate.MaximumInt(\"%s\", \"%s\", %s, int64(%d), %t)",
				param.Name, param.In, accessString(param), int64(*param.Maximum), param.ExclusiveMaximum))
		}
		if param.Minimum != nil {
			validations = append(validations, fmt.Sprintf("validate.MinimumInt(\"%s\", \"%s\", %s, int64(%d), %t)",
				param.Name, param.In, accessString(param), int64(*param.Minimum), param.ExclusiveMinimum))
		}
		if param.MultipleOf != nil {
			validations = append(validations, fmt.Sprintf("validate.MultipleOf(\"%s\", \"%s\", float64(%s), %f)",
				param.Name, param.In, accessString(param), *param.MultipleOf))
		}
	} else if param.Type == "number" {
		if param.Maximum != nil {
			validations = append(validations, fmt.Sprintf("validate.Maximum(\"%s\", \"%s\",  float64(%s), %f, %t)",
				param.Name, param.In, accessString(param), *param.Maximum, param.ExclusiveMaximum))
		}
		if param.Minimum != nil {
			validations = append(validations, fmt.Sprintf("validate.Minimum(\"%s\", \"%s\", float64(%s), %f, %t)",
				param.Name, param.In, accessString(param), *param.Minimum, param.ExclusiveMinimum))
		}
		if param.MultipleOf != nil {
			validations = append(validations, fmt.Sprintf("validate.MultipleOf(\"%s\", \"%s\", float64(%s), %f)",
				param.Name, param.In, accessString(param), *param.MultipleOf))
		}
	} else if param.Type == "array" {
		if param.MaxItems != nil {
			validations = append(validations, fmt.Sprintf("validate.MaxItems(\"%s\", \"%s\", int64(len(%s)), %d)",
				param.Name, param.In, accessString(param), *param.MaxItems))
		}
		if param.MinItems != nil {
			validations = append(validations, fmt.Sprintf("validate.MinItems(\"%s\", \"%s\", int64(len(%s)), %d)",
				param.Name, param.In, accessString(param), *param.MinItems))
		}
		if param.UniqueItems {
			validations = append(validations, fmt.Sprintf("validate.UniqueItems(\"%s\", \"%s\", %s)",
				param.Name, param.In, accessString(param)))
		}
	}
	return validations, nil
}

// accessString returns a string with the access of a variable in the struct named 'i'. For example:
// *i.Length
func accessString(param spec.Parameter) string {
	pointer := ""
	if !param.Required && param.Type != "array" {
		pointer = "*"
	}
	return fmt.Sprintf("%si.%s", pointer, utils.CamelCase(param.Name, true))
}

// StructParamName returns the name of the struct as used in the model struct
func StructParamName(param spec.Parameter) string {
	return utils.CamelCase(param.Name, true)
}

// DefaultAsString returns the default value as a string. We convert it into a string so it's easier to insert
// into the generated code and it doesn't make this logic really any different.
func DefaultAsString(param spec.Parameter) string {
	switch param.Default.(type) {
	case string:
		return param.Default.(string)
	case float64:
		if param.Type == "integer" {
			return strconv.FormatInt(int64(param.Default.(float64)), 10)
		}
		return strconv.FormatFloat(param.Default.(float64), 'E', -1, 64)
	case bool:
		return strconv.FormatBool(param.Default.(bool))
	default:
		panic(fmt.Errorf("unknown param type: %T", param))
	}
}
