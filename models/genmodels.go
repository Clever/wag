package models

import (
	"bytes"
	"fmt"

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
		"github.com/go-openapi/strfmt"
)

// These imports may not be used depending on the input parameters
var _ = json.Marshal
var _ = strconv.FormatInt
var _ = validate.Maximum
var _ = strfmt.NewFormats
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
	capOpID := swagger.Capitalize(op.ID)
	g.Printf("// %sInput holds the input parameters for a %s operation.\n", capOpID, op.ID)
	g.Printf("type %sInput struct {\n", capOpID)

	for _, param := range op.Parameters {
		if param.In == "formData" {
			return fmt.Errorf("Input parameters with 'In' formData are not supported")
		}

		var typeName string
		var err error
		if param.In != "body" {
			typeName, err = swagger.ParamToType(param, true)
			if err != nil {
				return err
			}
		} else {
			typeName, err = swagger.TypeFromSchema(param.Schema, false)
			if err != nil {
				return err
			}
			// All schema types are pointers
			typeName = "*" + typeName
		}

		g.Printf("\t%s %s\n", swagger.StructParamName(param), typeName)
	}
	g.Printf("}\n\n")

	return nil
}

func printInputValidation(g *swagger.Generator, op *spec.Operation) error {
	capOpID := swagger.Capitalize(op.ID)
	g.Printf("// Validate returns an error if any of the %sInput parameters don't satisfy the\n",
		capOpID)
	g.Printf("// requirements from the swagger yml file.\n")
	g.Printf("func (i %sInput) Validate() error{\n", capOpID)

	for _, param := range op.Parameters {
		if param.In == "body" {
			g.Printf("\tif err := i.%s.Validate(nil); err != nil {\n", swagger.StructParamName(param))
			g.Printf("\t\treturn err\n")
			g.Printf("\t}\n\n")
		}

		validations, err := swagger.ParamToValidationCode(param)
		if err != nil {
			return err
		}
		for _, validation := range validations {
			if param.Required {
				g.Printf(errCheck(validation))
			} else {
				g.Printf("\tif i.%s != nil {\n", swagger.StructParamName(param))
				g.Printf(errCheck(validation))
				g.Printf("\t}\n")
			}
		}
	}
	g.Printf("\treturn nil\n")
	g.Printf("}\n\n")

	return nil
}

// errCheck returns an if err := ifCondition; err != nil { return err } function
func errCheck(ifCondition string) string {
	return fmt.Sprintf(
		`	if err := %s; err != nil {
		return err
	}
`, ifCondition)
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
			// 3. Default - this is defined as a 500
			if err := validateStatusCodes(op.Responses.StatusCodeResponses); err != nil {
				return err
			}
			successTypes, err := generateSuccessTypes(capOpID, op.Responses.StatusCodeResponses)
			if err != nil {
				return err
			}
			g.Printf(successTypes)
			errorTypes, err := generateErrorTypes(capOpID, op.Responses.StatusCodeResponses)
			if err != nil {
				return err
			}
			g.Printf(errorTypes)
		}
	}
	return g.WriteFile("models/outputs.go")
}

func validateStatusCodes(responses map[int]spec.Response) error {
	for _, statusCode := range swagger.SortedStatusCodeKeys(responses) {
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
	}
	return nil
}

func generateSuccessTypes(capOpID string, responses map[int]spec.Response) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("// %sOutput defines the success output interface for %s.\n",
		capOpID, capOpID))
	buf.WriteString(fmt.Sprintf("type %sOutput interface {\n", capOpID))
	buf.WriteString(fmt.Sprintf("\t%sStatusCode() int\n", capOpID))
	buf.WriteString(fmt.Sprintf("}\n\n"))

	successStatusCodes := make([]int, 0)
	for _, statusCode := range swagger.SortedStatusCodeKeys(responses) {
		if statusCode >= 400 {
			continue
		}
		successStatusCodes = append(successStatusCodes, statusCode)
	}

	// We don't need to generate any success types if there is one or less success responses. In that
	// case we can just use the raw type
	if len(successStatusCodes) < 2 {
		return "", nil
	}

	for _, statusCode := range successStatusCodes {
		response := responses[statusCode]
		outputName := fmt.Sprintf("%s%dOutput", capOpID, statusCode)
		typeName, err := swagger.TypeFromSchema(response.Schema, false)
		if err != nil {
			return "", err
		}
		buf.WriteString(fmt.Sprintf("// %s defines the %d status code response for %s.\n",
			outputName, statusCode, capOpID))
		buf.WriteString(fmt.Sprintf("type %s %s\n\n", outputName, typeName))

		buf.WriteString(fmt.Sprintf("// %sStatusCode returns the status code for the operation.\n",
			capOpID))
		buf.WriteString(fmt.Sprintf("func (o %s) %sStatusCode() int {\n", outputName, capOpID))
		buf.WriteString(fmt.Sprintf("\treturn %d\n", statusCode))
		buf.WriteString(fmt.Sprintf("}\n\n"))
	}
	return buf.String(), nil
}

func generateErrorTypes(capOpID string, responses map[int]spec.Response) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("// %sError defines the error interface for %s.\n",
		capOpID, capOpID))
	buf.WriteString(fmt.Sprintf("type %sError interface {\n", capOpID))
	buf.WriteString(fmt.Sprintf("\terror // Extend the error interface\n"))
	buf.WriteString(fmt.Sprintf("\t%sStatusCode() int\n", capOpID))
	buf.WriteString(fmt.Sprintf("}\n\n"))

	for _, statusCode := range swagger.SortedStatusCodeKeys(responses) {
		if statusCode < 400 {
			continue
		}
		response := responses[statusCode]

		outputName := fmt.Sprintf("%s%dOutput", capOpID, statusCode)
		typeName, err := swagger.TypeFromSchema(response.Schema, false)
		if err != nil {
			return "", err
		}

		buf.WriteString(fmt.Sprintf("// %s defines the %d status code response for %s.\n",
			outputName, statusCode, capOpID))
		buf.WriteString(fmt.Sprintf("type %s %s\n\n", outputName, typeName))

		buf.WriteString(fmt.Sprintf("// Error returns `Status Code: X`. We implemeted it to satisfy\n"))
		buf.WriteString(fmt.Sprintf("// the error interface. For a more descriptive error message see\n"))
		buf.WriteString(fmt.Sprintf("// the output type.\n"))
		buf.WriteString(fmt.Sprintf("func (o %s) Error() string {\n", outputName))
		buf.WriteString(fmt.Sprintf("\treturn \"Status Code: %d\"\n", statusCode))
		buf.WriteString(fmt.Sprintf("}\n\n"))

		buf.WriteString(fmt.Sprintf("// %sStatusCode returns the status code for the operation.\n",
			capOpID))
		buf.WriteString(fmt.Sprintf("func (o %s) %sStatusCode() int {\n", outputName, capOpID))
		buf.WriteString(fmt.Sprintf("\treturn %d\n", statusCode))
		buf.WriteString(fmt.Sprintf("}\n\n"))
	}
	return buf.String(), nil
}

// defaultOutputTypes returns the string defining the default output type
func defaultOutputTypes() string {
	return fmt.Sprintf(`
// DefaultInternalError represents a generic 500 response.
type DefaultInternalError struct {
	Msg string %s
}

// Error returns the internal error that caused the 500.
func (d DefaultInternalError) Error() string {
	return d.Msg
}

// DefaultBadRequest represents a generic 400 response. It used internally by Swagger as the
// response when a request fails the validation defined in the Swagger yml file.
type DefaultBadRequest struct {
	Msg string %s
}

// Error returns the validation error that caused the 400.
func (d DefaultBadRequest) Error() string {
	return d.Msg
}

`, "`json:\"msg\"`", "`json:\"msg\"`")
}
