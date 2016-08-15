package models

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/swagger"

	"github.com/go-swagger/go-swagger/generator"
)

// Generate writes the files to the client directories
func Generate(packageName, swaggerFile string, swagger spec.Swagger) error {

	// generate models with go-swagger
	if err := generator.GenerateServer("", []string{}, []string{}, generator.GenOpts{
		Spec:           swaggerFile,
		ModelPackage:   "models",
		Target:         "./generated/",
		IncludeModel:   true,
		IncludeHandler: false,
		IncludeSupport: false,
	}); err != nil {
		return err
	}

	if err := generateOutputs(packageName, swagger.Paths); err != nil {
		return err
	}
	return generateInputs(packageName, swagger.Paths)
}

func generateInputs(packageName string, paths *spec.Paths) error {

	g := swagger.Generator{PackageName: packageName}

	g.Printf(`
package models

import(
		"encoding/json"
		"strconv"
		"github.com/go-openapi/validate"
)

// These imports may not be used depending on the input parameters
var _ = json.Marshal
var _ = strconv.FormatInt
var _ = validate.Maximum
`)

	for _, pathKey := range swagger.SortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, opKey := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			if err := printInputStruct(&g, op); err != nil {
				return err
			}
			if err := printInputValidation(&g, op); err != nil {
				return err
			}
		}
	}

	return g.WriteFile("models/inputs.go")
}

func printInputStruct(g *swagger.Generator, op *spec.Operation) error {
	g.Printf("type %sInput struct {\n", swagger.Capitalize(op.ID))

	for _, param := range op.Parameters {
		if param.In == "formData" {
			return fmt.Errorf("Input parameters with 'In' formData are not supported")
		}

		var typeName string
		if param.In != "body" {
			switch param.Type {
			case "string":
				typeName = "string"
			case "integer":
				typeName = "int64"
			case "boolean":
				typeName = "bool"
			case "number":
				typeName = "float64"
			default:
				// Note. We don't support 'array' or 'file' types even though they're in the
				// Swagger spec.
				return fmt.Errorf("Unsupported param type")
			}
		} else {
			var err error
			typeName, err = typeFromSchema(param.Schema)
			if err != nil {
				return err
			}
		}
		pointerStr := ""
		if !param.Required {
			pointerStr = "*"
		}

		g.Printf("\t%s %s%s\n", swagger.Capitalize(param.Name), pointerStr, typeName)
	}
	g.Printf("}\n\n")

	return nil
}

// oneTabErrorCheck returns an if err := ifCondition; err != nil { return err } function
func oneTabErrorCheck(ifCondition string) string {
	return fmt.Sprintf(`
	if err := %s; err != nil {
		return err
	}
`, ifCondition)
}

func printInputValidation(g *swagger.Generator, op *spec.Operation) error {
	g.Printf("func (i %sInput) Validate() error{\n", swagger.Capitalize(op.ID))

	for _, param := range op.Parameters {
		if param.In == "body" {
			g.Printf("\tif err := i.%s.Validate(nil); err != nil {\n", swagger.Capitalize(param.Name))
			g.Printf("\t\treturn err\n")
			g.Printf("\t}\n\n")
		}

		if param.Type == "string" {
			if param.MaxLength != nil {
				g.Printf(oneTabErrorCheck(fmt.Sprintf("validate.MaxLength(\"%s\", \"%s\", %s, %d)",
					param.Name, param.In, accessString(param), *param.MaxLength)))
			}
			if param.MinLength != nil {
				g.Printf(oneTabErrorCheck(fmt.Sprintf("validate.MinLength(\"%s\", \"%s\", %s, %d)",
					param.Name, param.In, accessString(param), *param.MaxLength)))
			}
			if param.Pattern != "" {
				g.Printf(oneTabErrorCheck(fmt.Sprintf("validate.Pattern(\"%s\", \"%s\", %s, \"%s\")",
					param.Name, param.In, accessString(param), param.Pattern)))
			}
			if len(param.Enum) != 0 {
				strEnum := []string{}
				for _, enum := range param.Enum {
					strEnum = append(strEnum, enum.(string))
				}
				g.Printf(oneTabErrorCheck(fmt.Sprintf("validate.Enum(\"%s\", \"%s\", %s, %s)",
					param.Name, param.In, accessString(param), fmt.Sprintf("[]interface{}{\"%s\"}", strings.Join(strEnum, "\",\"")))))
			}
		} else if param.Type == "integer" {
			if param.Maximum != nil {
				g.Printf(oneTabErrorCheck(fmt.Sprintf("validate.MaximumInt(\"%s\", \"%s\", %s, %d, %t)",
					param.Name, param.In, accessString(param), int64(*param.Maximum), param.ExclusiveMaximum)))
			}
			if param.Minimum != nil {
				g.Printf(oneTabErrorCheck(fmt.Sprintf("validate.MinimumInt(\"%s\", \"%s\", %s, %d, %t)",
					param.Name, param.In, accessString(param), int64(*param.Minimum), param.ExclusiveMinimum)))
			}
			if param.MultipleOf != nil {
				g.Printf(oneTabErrorCheck(fmt.Sprintf("validate.MultipleOf(\"%s\", \"%s\", float64(%s), %f)",
					param.Name, param.In, accessString(param), *param.MultipleOf)))
			}
		} else if param.Type == "number" {
			if param.Maximum != nil {
				g.Printf(oneTabErrorCheck(fmt.Sprintf("validate.Maximum(\"%s\", \"%s\",  %s, %f, %t)",
					param.Name, param.In, accessString(param), *param.Maximum, param.ExclusiveMaximum)))
			}
			if param.Minimum != nil {
				g.Printf(oneTabErrorCheck(fmt.Sprintf("validate.Minimum(\"%s\", \"%s\", %s, %f, %t)",
					param.Name, param.In, accessString(param), *param.Minimum, param.ExclusiveMinimum)))
			}
			if param.MultipleOf != nil {
				g.Printf(oneTabErrorCheck(fmt.Sprintf("validate.MultipleOf(\"%s\", \"%s\", %s, %f)",
					param.Name, param.In, accessString(param), *param.MultipleOf)))
			}
		}
	}
	g.Printf("\treturn nil\n")
	g.Printf("}\n\n")

	return nil
}

// TODO: Add a nice comment!!!
func accessString(param spec.Parameter) string {
	pointer := ""
	if !param.Required {
		pointer = "*"
	}
	return fmt.Sprintf("%si.%s", pointer, swagger.Capitalize(param.Name))
}

func generateOutputs(packageName string, paths *spec.Paths) error {
	g := swagger.Generator{PackageName: packageName}

	g.Printf("package models\n\n")

	g.Printf(defaultOutputTypes())

	for _, pathKey := range swagger.SortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, opKey := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			capOpID := swagger.Capitalize(op.ID)
			// We classify response keys into three types:
			// 1. 200-399 - these are "success" responses and implement the Output interface
			// 	defined above
			// 2. 400-599 - these are "failure" responses and implement the error interface
			// 3. Default - this is defined as a 500 (TODO: decide if we want to keep this...)

			// Define the success interface
			g.Printf("type %sOutput interface {\n", capOpID)
			g.Printf("\t%sStatus() int\n", capOpID)
			g.Printf("\t// Data is the underlying model object. We know it is json serializable\n")
			g.Printf("\t%sData() interface{}\n", capOpID)
			g.Printf("}\n\n")

			// Define the error interface
			g.Printf("type %sError interface {\n", capOpID)
			g.Printf("\terror // Extend the error interface\n")
			g.Printf("\t%sStatusCode() int\n", capOpID)
			g.Printf("}\n\n")

			for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
				response := op.Responses.StatusCodeResponses[statusCode]

				if statusCode < 200 || statusCode > 599 {
					// TODO: Write a test for this...
					return fmt.Errorf("Response map key must be an integer between 200 and 599 or "+
						"the string 'default'. Was %d", statusCode)
				}
				if statusCode == 400 {
					return fmt.Errorf("Use the pre-defined default 400 response 'DefaultBadRequest' " +
						"instead of defining your own")
				} else if statusCode == 500 {
					return fmt.Errorf("Use the pre-defined default 500 response `DefaultInternalError` " +
						"instead of defining your own")
				}

				outputName := fmt.Sprintf("%s%dOutput", capOpID, statusCode)
				typeName, err := typeFromSchema(response.Schema)
				if err != nil {
					return err
				}

				g.Printf("type %s %s\n\n", outputName, typeName)

				g.Printf("func (o %s) %sData() interface{} {\n", outputName, capOpID)
				g.Printf("\treturn o\n")
				g.Printf("}\n\n")

				if statusCode < 400 {

					// TODO: Do we really want to have that as part of the interface?
					g.Printf("func (o %s) %sStatus() int {\n", outputName, capOpID)
					// TODO: Use the right status code...
					g.Printf("\treturn %d\n", statusCode)
					g.Printf("}\n\n")

				} else {

					g.Printf("func (o %s) Error() string {\n", outputName)
					// TODO: Would it make sense to give this a constructor so we can have a more detailed
					// error message?
					g.Printf("\treturn \"Status Code: \" + \"%d\"\n", statusCode)
					g.Printf("}\n\n")

					g.Printf("func (o %s) %sStatusCode() int {\n", outputName, capOpID)
					g.Printf("\treturn %d\n", statusCode)
					g.Printf("}\n\n")
				}
			}
		}
	}
	return g.WriteFile("models/outputs.go")
}

// defaultOutputTypes returns the string defining the default output type
func defaultOutputTypes() string {
	return fmt.Sprintf(`
// DefaultInternalError represents a generic 500 response.
type DefaultInternalError struct {
	Msg string %s
}

func (d DefaultInternalError) Error() string {
	return d.Msg
}

type DefaultBadRequest struct {
	Msg string %s
}

func (d DefaultBadRequest) Error() string {
	return d.Msg
}

`, "`json:\"msg\"`", "`json:\"msg\"`")
}

func typeFromSchema(schema *spec.Schema) (string, error) {
	// We support three types of schemas
	// 1. An empty schema
	// 2. A schema with one element, the $ref key
	// 3. A schema with two elements. One a type with value 'array' and another items field
	// referencing the $ref
	// TODO: The error messages here aren't great...
	if schema == nil {
		// represent this as a string, which is empty by default
		return "string", nil
	} else if schema.Ref.String() != "" {
		ref := schema.Ref.String()
		if !strings.HasPrefix(ref, "#/definitions/") {
			return "", fmt.Errorf("schema.$ref has undefined reference type. Must be #/definitions")
		}
		return ref[len("#/definitions/"):], nil
	} else {
		schemaType := schema.Type
		if len(schemaType) != 1 || schemaType[0] != "array" {
			return "", fmt.Errorf("Two element schemas must have a 'type' field with the value 'array'")
		}
		items := schema.Items
		if items == nil || items.Schema == nil {
			return "", fmt.Errorf("Two element schemas must have an '$ref' field in the 'items' descriptions")
		}
		ref := items.Schema.Ref.String()
		if !strings.HasPrefix(ref, "#/definitions/") {
			return "", fmt.Errorf("schema.$ref has undefined reference type. Must be #/definitions")
		}
		return "[]" + ref[len("#/definitions/"):], nil
	}
}
