package models

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/go-utils/stringset"
	"github.com/Clever/wag/swagger"
	"github.com/Clever/wag/templates"

	"github.com/go-swagger/go-swagger/generator"
)

// Generate writes the files to the client directories
func Generate(packageName string, s spec.Swagger) error {

	tmpFile, err := swagger.WriteToFile(&s)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	// generate models with go-swagger
	if err := generator.GenerateServer("", []string{}, []string{}, &generator.GenOpts{
		Spec:           tmpFile,
		ModelPackage:   "models",
		Target:         fmt.Sprintf("%s/src/%s/", os.Getenv("GOPATH"), packageName),
		IncludeModel:   true,
		IncludeHandler: false,
		IncludeSupport: false,
	}); err != nil {
		return fmt.Errorf("error generating go-swagger models: %s", err)
	}

	if err := generateOutputs(packageName, s); err != nil {
		return fmt.Errorf("error generating outputs: %s", err)
	}
	if err := generateInputs(packageName, s.Paths); err != nil {
		return fmt.Errorf("error generating inputs: %s", err)
	}
	return nil
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
			// Do not generate an input struct + validation for an
			// operation that has a single, schema'd input.
			// The input to these will be the model generated for
			// the schema.
			if singleSchemaedBodyParameter, _ := swagger.SingleSchemaedBodyParameter(op); singleSchemaedBodyParameter {
				continue
			}
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
			return fmt.Errorf("input parameters with 'In' formData are not supported")
		}

		typeName, pointer, err := swagger.ParamToType(param)
		if err != nil {
			return err
		}
		if pointer {
			typeName = "*" + typeName
		}

		g.Printf("\t%s %s\n", swagger.StructParamName(param), typeName)
	}
	g.Printf("}\n\n")

	return nil
}

func printInputValidation(g *swagger.Generator, op *spec.Operation) error {
	singleStringPathParameter, paramName := swagger.SingleStringPathParameter(op)
	if singleStringPathParameter {
		capOpID := swagger.Capitalize(op.ID)
		g.Printf("// Validate%sInput returns an error if the input parameter doesn't\n",
			capOpID)
		g.Printf("// satisfy the requirements in the swagger yml file.\n")
		g.Printf("func Validate%sInput(%s string) error{\n", capOpID, paramName)
	} else {
		capOpID := swagger.Capitalize(op.ID)
		g.Printf("// Validate returns an error if any of the %sInput parameters don't satisfy the\n",
			capOpID)
		g.Printf("// requirements from the swagger yml file.\n")
		g.Printf("func (i %sInput) Validate() error{\n", capOpID)
	}

	for _, param := range op.Parameters {
		_, pointer, err := swagger.ParamToType(param)
		if err != nil {
			return err
		}
		validations, err := swagger.ParamToValidationCode(param, op)
		if err != nil {
			return err
		}
		t := validateTemplate{
			Type:         param.In,
			AccessString: "i." + swagger.StructParamName(param),
			Pointer:      pointer || param.Type == "array",
			Validations:  validations,
		}
		if single, _ := swagger.SingleStringPathParameter(op); single {
			t.AccessString = param.Name
		}
		str, err := templates.WriteTemplate(validationStr, t)
		if err != nil {
			return err
		}
		g.Printf(str)
	}
	g.Printf("\treturn nil\n")
	g.Printf("}\n\n")

	return nil
}

type validateTemplate struct {
	Type         string
	AccessString string
	Pointer      bool
	Validations  []string
}

//

var validationStr = `
	{{if eq .Type "body" -}}
	if err := {{.AccessString}}.Validate(nil); err != nil {
		return err
	}
	{{- end -}}
	{{- $type := .Type -}}
	{{- $accessString := .AccessString -}}
	{{- $pointer := .Pointer -}}
	{{- range $i, $validation := .Validations -}}
		{{- if eq $type "header" -}}
		if len({{$accessString}}) > 0 {
			if err := {{$validation}}; err != nil {
				return err
			}
		}
		{{else -}}
		{{- if $pointer}}
		if {{$accessString}} != nil {
		{{- end -}}
		if err := {{$validation}}; err != nil {
			return err
		}
		{{if $pointer -}}
		}
		{{- end -}}
		{{end -}}
	{{- end }}
`

func generateOutputs(packageName string, s spec.Swagger) error {
	g := swagger.Generator{PackageName: packageName}

	g.Printf("package models\n\n")

	// It's a bit wonky that we're writing these into output.go instead of the file
	// defining each of the types, but I think that's okay for now. We can clean this
	// up if it becomes confusing.
	errorMethodCode, err := generateErrorMethods(&s)
	if err != nil {
		return err
	}
	g.Printf(errorMethodCode)
	return g.WriteFile("models/outputs.go")
}

// generateErrorMethods finds all responses all error responses and generates an error
// method for them.
func generateErrorMethods(s *spec.Swagger) (string, error) {
	errorTypes := stringset.New()

	for _, pathKey := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		path := s.Paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, opKey := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {

				if statusCode < 400 {
					continue
				}

				typeName, _ := swagger.OutputType(s, op, statusCode)
				if strings.HasPrefix(typeName, "models.") {
					typeName = typeName[7:]
				}
				errorTypes.Add(typeName)
			}
		}
	}

	sortedErrors := errorTypes.ToList()
	sort.Strings(sortedErrors)

	var buf bytes.Buffer
	for _, errorType := range sortedErrors {
		buf.WriteString(fmt.Sprintf(`
func (o %s) Error() string {
	return o.Message
}

`, errorType))
	}

	return buf.String(), nil

}
