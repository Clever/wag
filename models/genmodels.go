package models

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/go-utils/stringset"
	goClient "github.com/Clever/wag/v9/clients/go"
	"github.com/Clever/wag/v9/swagger"
	"github.com/Clever/wag/v9/templates"
	"github.com/Clever/wag/v9/utils"

	"github.com/go-swagger/go-swagger/generator"
)

// Generate writes the files to the client directories
func Generate(packageName, basePath, outputPath string, s spec.Swagger) error {
	tmpFile, err := swagger.WriteToFile(&s)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	genopts := generator.GenOpts{
		Spec:           tmpFile,
		ModelPackage:   "models",
		Target:         basePath,
		IncludeModel:   true,
		IncludeHandler: false,
		IncludeSupport: false,
	}

	// The zero-values for many fields in GenOpts are not good defaults; this call
	// sets them to actually good defaults.
	// Setting GenOpts.FlattenOpts is particularly important.
	genopts.EnsureDefaults()

	// generate models with go-swagger
	if err := generator.GenerateServer("", []string{}, []string{}, &genopts); err != nil {
		return fmt.Errorf("error generating go-swagger models: %s", err)
	}

	if err := generateOutputs(basePath, s); err != nil {
		return fmt.Errorf("error generating outputs: %s", err)
	}
	if err := generateInputs(basePath, s); err != nil {
		return fmt.Errorf("error generating inputs: %s", err)
	}
	if err := CreateModFile("models/go.mod", basePath, packageName, outputPath); err != nil {
		return fmt.Errorf("error creating go.mod file: %s", err)
	}
	return nil
}

// CreateModFile creates a go.mod file for the client module.
func CreateModFile(path string, basePath, packageName, outputPath string) error {
	outputPath = strings.TrimPrefix(outputPath, ".")
	moduleName, versionSuffix := utils.ExtractModuleNameAndVersionSuffix(packageName, outputPath)

	absPath := basePath + "/" + path
	f, err := os.Create(absPath)
	if err != nil {
		return err
	}

	defer f.Close()
	modFileString := `
module ` + moduleName + outputPath + `/models` + versionSuffix + `


go 1.24

require (
	github.com/go-openapi/errors v0.20.2
	github.com/go-openapi/strfmt v0.21.2
	github.com/go-openapi/swag v0.21.1
	github.com/go-openapi/validate v0.22.0
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/uuid v1.1.2 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
)

	`
	_, err = f.WriteString(modFileString)
	if err != nil {
		return err
	}

	return nil
}

func generateInputs(basePath string, s spec.Swagger) error {
	g := swagger.Generator{BasePath: basePath}

	g.Printf(`
package models

import(
		"encoding/json"
		"fmt"
		"net/url"
		"strconv"
		"strings"

		"github.com/go-openapi/validate"
		"github.com/go-openapi/strfmt"
)

// These imports may not be used depending on the input parameters
var _ = json.Marshal
var _ = fmt.Sprintf
var _ = url.QueryEscape
var _ = strconv.FormatInt
var _ = strings.Replace
var _ = validate.Maximum
var _ = strfmt.NewFormats
`)

	paths := s.Paths
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
			if err := printInputValidation(&g, op, goClient.IsBinaryBody(op, s.Definitions)); err != nil {
				return err
			}
			if err := printInputSerializer(&g, op, s.BasePath, pathKey); err != nil {
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

func printInputValidation(g *swagger.Generator, op *spec.Operation, binaryBody bool) error {
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
			BinaryBody:   binaryBody,
		}
		if single, _ := swagger.SingleStringPathParameter(op); single {
			t.AccessString = param.Name
		}
		str, err := templates.WriteTemplate(validationStr, t)
		if err != nil {
			return err
		}
		g.Print(str)
	}
	g.Print("\treturn nil\n")
	g.Print("}\n\n")

	return nil
}

type validateTemplate struct {
	Type         string
	AccessString string
	Pointer      bool
	Validations  []string
	BinaryBody   bool
}

//

var validationStr = `
	{{if eq .Type "body" -}}
	{{if not .BinaryBody -}}
	if {{.AccessString}} != nil {
		if err := {{.AccessString}}.Validate(nil); err != nil {
			return err
		}
	}
	{{- end -}}
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

func printInputSerializer(g *swagger.Generator, op *spec.Operation, basePath, methodPath string) error {
	singleStringPathParameter, paramName := swagger.SingleStringPathParameter(op)
	if singleStringPathParameter {
		capOpID := swagger.Capitalize(op.ID)
		g.Printf("// %sInputPath returns the URI path for the input.\n", capOpID)
		g.Printf("func %sInputPath(%s string) (string, error) {\n", capOpID, paramName)
	} else {
		capOpID := swagger.Capitalize(op.ID)
		g.Printf("// Path returns the URI path for the input.\n")
		g.Printf("func (i %sInput) Path() (string, error) {\n", capOpID)
	}

	g.Printf("\tpath := \"%s\"\n", basePath+methodPath)
	g.Printf("\turlVals := url.Values{}\n")

	for _, param := range op.Parameters {
		t := swagger.ParamToTemplate(&param, op)
		if param.In == "path" {
			pt := pathParamTemplate{
				Name:         param.Name,
				PathName:     "{" + t.Name + "}",
				ToStringCode: t.ToStringCode,
			}
			str, err := templates.WriteTemplate(pathParamStr, pt)
			if err != nil {
				panic(fmt.Errorf("unexpected error: %s", err))
			}
			g.Print(str)
		} else if param.In == "query" {
			str, err := templates.WriteTemplate(queryParamStr, t)
			if err != nil {
				panic(fmt.Errorf("unexpected error: %s", err))
			}
			g.Print(str)
		}
	}

	g.Print("\n\treturn path + \"?\" + urlVals.Encode(), nil\n")
	g.Print("}\n\n")

	return nil
}

// pathParamTemplate is a seperate template because I couldn't find a non-annoying
// way to encode {".Name"} in the template
type pathParamTemplate struct {
	Name         string
	PathName     string
	ToStringCode string
}

var pathParamStr = `
	path{{.Name}} := {{.ToStringCode}}
	if path{{.Name}} == "" {
		err := fmt.Errorf("{{.Name}} cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{{.PathName}}", path{{.Name}}, -1)
`

var queryParamStr = `
	{{if .Pointer -}}
	if {{.AccessString}} != nil {
	{{end}}
	{{- if eq .Type "array" -}}
	for _, v := range {{.AccessString}} {
		urlVals.Add("{{.Name}}", v)
	}
	{{- else -}}
	urlVals.Add("{{.Name}}", {{.ToStringCode}})
	{{- end -}}
	{{if .Pointer}}
	}
	{{end}}
`

func generateOutputs(basePath string, s spec.Swagger) error {
	g := swagger.Generator{BasePath: basePath}

	g.Printf("package models\n\nimport \"fmt\"\n")

	// It's a bit wonky that we're writing these into output.go instead of the file
	// defining each of the types, but I think that's okay for now. We can clean this
	// up if it becomes confusing.
	errorMethodCode, err := generateErrorMethods(&s)
	if err != nil {
		return err
	}
	g.Print(errorMethodCode)
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

	// default unknown response type is also an error type
	buf.WriteString(`
	func (u UnknownResponse) Error() string {
		return fmt.Sprintf("unknown response with status: %d body: %s", u.StatusCode, u.Body)
	}
	
`)

	return buf.String(), nil
}
