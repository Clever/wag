package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-openapi/spec"
	"github.com/go-swagger/go-swagger/generator"

	"text/template"
)

type SwaggerOperation struct {
	ID          string                     `yaml:"operationId"`
	Description string                     `yaml:"description"`
	Responses   map[string]SwaggerResponse `yaml:"responses"`
	// TODO: Parameters can also be a reference type. We should disallow that given that
	// we don't support parameters to be defined anywhere else.
	Parameters []SwaggerParameter `yaml:"parameters"`

	// Not supported
	Consumes []string               `yaml:"consumers"`
	Produces []string               `yaml:"produces"`
	Schemes  []string               `yaml:"schemes"`
	Security map[string]interface{} `yaml:"security"`
	Tags     []string               `yaml:"tags"`
}

// A regex requiring the field to be start with a letter and be alphanumeric
var alphaNumRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")

// Validate checks if the swagger operation has any fields we don't support
func ValidateOp(s *spec.Operation) error {
	if len(s.Consumes) != 0 {
		fmt.Errorf("Consumes not supported in swagger operations")
	}
	if len(s.Produces) != 0 {
		fmt.Errorf("Produces not supported in swagger operations")
	}
	if len(s.Schemes) != 0 {
		fmt.Errorf("Schemes not supported in swagger operations")
	}
	if len(s.Security) != 0 {
		fmt.Errorf("Security not supported in swagger operations")
	}

	// TODO: A test for this?
	/*
		for path, pathObj := range paths {
		for method, op := range pathObj {
	*/
	if !alphaNumRegex.MatchString(s.ID) {
		// We need this because we build function / struct names with the operationID.
		// We could strip out the special characters, but it seems clearly to just enforce
		// this.
		return fmt.Errorf("OperationIDs must be alphanumeric and start with a letter")
	}

	return nil
}

type SwaggerParameter struct {
	Name string `yaml:"name"`
	In   string `yaml:"in"`
	Type string `yaml:"type"`
	// The schema here has to be {$ref: "..."}. We don't support defining your own
	// schema here.
	Schema map[string]interface{} `yaml:"schema"`
}

type SwaggerResponse struct {
	Description string                 `yaml:"description"`
	Schema      map[string]interface{} `yaml:"schema"`

	// Fields we don't support
	Header map[string]interface{} `yaml:"headers"`
}

type Swagger struct {
	BasePath string `yaml:"basePath"`

	// Partially implemented
	Paths map[string]map[string]SwaggerOperation `yaml:"paths"`

	// We rely on the go-swagger code to generate all the definitions
	Definitions map[string]interface{} `yaml:"definitions"`

	// Fields we support, but only with a certain set of values
	Swagger  string   `yaml:"swagger"`
	Schemes  []string `yaml:"schemes"`
	Consumes []string `yaml:"consumes"`
	Produces []string `yaml:"produces"`

	// Fields we don't support
	Host                string                 `yaml:"host"`
	Parameters          map[string]interface{} `yaml:"parameters"`
	Responses           map[string]interface{} `yaml:"responses"`
	SecurityDefinitions map[string]interface{} `yaml:"securityDefinitions"`
	Security            map[string]interface{} `yaml:"security"`
	Tags                []string               `yaml:"tags"`
}

// Validates returns an error if the swagger file is invalid or uses fields
// we don't support. Note that this isn't a comprehensive check for all things
// we don't support, so this may not return an error, but the Swagger file might
// have values we don't support
func Validate(s spec.Swagger) error {
	if s.Swagger != "2.0" {
		return fmt.Errorf("Unsupported Swagger version %s", s.Swagger)
	}

	if len(s.Schemes) != 1 || s.Schemes[0] != "http" {
		return fmt.Errorf("Schemes only supports 'http', not %s")
	}

	if len(s.Consumes) != 1 || s.Consumes[0] != "application/json" {
		return fmt.Errorf("Consumes only supports 'application/json'")
	}

	if len(s.Produces) != 1 || s.Produces[0] != "application/json" {
		return fmt.Errorf("Produces only supports 'application/json'")
	}

	if s.Host != "" {
		return fmt.Errorf("Host parameter is not supported")
	}

	if len(s.Parameters) != 0 {
		return fmt.Errorf("Global parameters definitions are not supported. Define parameters on a per request basis.")
	}

	if len(s.Responses) != 0 {
		return fmt.Errorf("Global response definitions are not supported. Define responses on a per request basis")
	}

	if len(s.SecurityDefinitions) != 0 {
		return fmt.Errorf("Security definitions definition not supported")
	}

	if len(s.Security) != 0 {
		return fmt.Errorf("Security field not supported")
	}

	// Validate paths
	for fieldName, pathItem := range s.Paths.Paths {
		// Note that $ref and parameters are not valid as of now
		if !sliceContains([]string{"get", "put", "post", "delete", "options", "head", "patch"}, fieldName) {
			fmt.Errorf("Invalid path field name: %s", fieldName)
		}

		for _, op := range pathItemOperations(pathItem) {
			if err := ValidateOp(op); err != nil {
				return err
			}
		}
	}

	return nil
}

func pathItemOperations(item spec.PathItem) map[string]*spec.Operation {
	ops := make(map[string]*spec.Operation)
	if item.Get != nil {
		ops["GET"] = item.Get
	}
	if item.Put != nil {
		ops["PUT"] = item.Put
	}
	if item.Post != nil {
		ops["POST"] = item.Post
	}
	if item.Delete != nil {
		ops["DELETE"] = item.Delete
	}
	if item.Options != nil {
		ops["OPTIONS"] = item.Options
	}
	if item.Head != nil {
		ops["HEAD"] = item.Head
	}
	if item.Patch != nil {
		ops["PATCH"] = item.Patch
	}
	return ops
}

func sliceContains(slice []string, key string) bool {
	for _, val := range slice {
		if val == key {
			return true
		}
	}
	return false
}

// TODO: Should this be a function on the swagger object directly?
func capitalize(input string) string {
	return strings.ToUpper(input[0:1]) + input[1:]
}

func main() {

	swaggerFile := flag.String("file", "swagger.yml", "the spec file to use")
	packageName := flag.String("package", "", "package of the generated code")
	flag.Parse()
	if *packageName == "" {
		log.Fatal("package is required")
	}

	// generate models with go-swagger
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
	if err := generator.GenerateServer("", []string{}, []string{}, generator.GenOpts{
		Spec:           *swaggerFile,
		ModelPackage:   "models",
		Target:         "./generated/",
		IncludeModel:   true,
		IncludeHandler: false,
		IncludeSupport: false,
	}); err != nil {
		log.Fatal(err)
	}

	doc, err := loads.Spec(*swaggerFile)
	if err != nil {
		log.Fatalf("Error loading swagger file: %s", err)
	}
	swagger := *doc.Spec()

	if err := Validate(swagger); err != nil {
		log.Fatalf("Swagger file not valid: %s", err)
	}

	for _, dir := range []string{"server", "client", "models"} {
		dirName := os.Getenv("GOPATH") + "/src/" + *packageName + "/" + dir
		if err := os.Mkdir(dirName, 0700); err != nil {
			if !os.IsExist(err.(*os.PathError)) {
				log.Fatalf("Could not create directory: %s, error: %s", dirName, err)
			}
		}
	}

	if err := generateRouter(swagger.BasePath, swagger.Paths); err != nil {
		panic(err)
	}
	// TODO: Is this really the way I want to do this???
	if err := generateContextsAndControllers(*packageName, swagger.Paths); err != nil {
		log.Fatalf("Failed creating context and controllers: %s", err)
	}
	if err := generateHandlers(*packageName, swagger.Paths); err != nil {
		log.Fatalf("Failed building handlers: %s", err)
	}
	if err := generateOutputs(*packageName, swagger.Paths); err != nil {
		log.Fatalf("Failed building outputs %s", err)
	}
	if err := generateClients(*packageName, swagger); err != nil {
		log.Fatalf("Failed generating clients %s", err)
	}

}

type Generator struct {
	buf bytes.Buffer
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) WriteFile(path string) error {
	fileBytes, err := format.Source(g.buf.Bytes())
	if err != nil {
		// This will error if the code isn't valid so let's print it to make it
		// easier to debug
		fmt.Printf("BAD CODE\n%s\n", string(g.buf.Bytes()))
		return err
	}
	return ioutil.WriteFile(path, fileBytes, 0644)
}

func generateRouter(basePath string, paths *spec.Paths) error {
	var g Generator

	// TODO: Add something to all these about being auto-generated

	g.Printf(
		`package server

import (
	"net/http"

	"github.com/gorilla/mux"

	gContext "github.com/gorilla/context"
	"golang.org/x/net/context"
)

type contextKey struct{}

func SetupServer(r *mux.Router, c Controller) http.Handler {
	controller = c // TODO: get rid of global variable?`)

	for path, pathItem := range paths.Paths {
		for method, op := range pathItemOperations(pathItem) {
			// TODO: Note the coupling for the handler name here and in the handler function. Does that mean these should be
			// together? Probably...

			g.Printf("\n")
			tmpl, err := template.New("routerFunction").Parse(routerFunctionTemplate)
			if err != nil {
				return err
			}
			err = tmpl.Execute(&g.buf, routerTemplate{Method: method, Path: basePath + path,
				HandlerName: capitalize(op.ID)})
			if err != nil {
				return err
			}
		}
	}
	// TODO: It's a bit weird that this returns a pointer that it modifies...
	g.Printf("\treturn withMiddleware(r)\n")
	g.Printf("}\n")

	return g.WriteFile("generated/server/router.go")
}

type routerTemplate struct {
	Method      string
	Path        string
	HandlerName string
}

var routerFunctionTemplate = `	r.Methods("{{.Method}}").Path("{{.Path}}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gContext.Get(r, contextKey{}).(context.Context)
		{{.HandlerName}}Handler(ctx, w, r)
	})
`

func generateContextsAndControllers(packageName string, paths *spec.Paths) error {

	// Only create the controller the first time. After that leave it as is
	// TODO: This isn't convenient when developing, so maybe have a flag...?
	// TODO: We still want to build the contexts.go file, so need to figure that out...
	// if _, err := os.Stat("generated/controller.go"); err == nil {
	//	return nil
	//}

	var g Generator

	g.Printf("package models\n\n")

	g.Printf("import (\n")

	g.Printf("\t\"encoding/json\"\n")
	g.Printf("\t\"strconv\"\n")
	g.Printf(")\n\n")

	// These two imports are only used if we have body parameters, so if we don't have these the
	// compiler will complain about unused imports
	g.Printf("var _ = json.Marshal\n\n")
	g.Printf("var _ = strconv.FormatInt\n\n")
	// TODO: Not as easy to do this for models since we have to find some type. Will be easier if
	// we move it into the same package

	for _, path := range paths.Paths {
		for _, op := range pathItemOperations(path) {

			if err := printInputStruct(&g, op); err != nil {
				return err
			}
			if err := printInputValidation(&g, op); err != nil {
				return err
			}
		}
	}

	var interfaceGenerator Generator
	interfaceGenerator.Printf("package server\n\n")
	interfaceGenerator.Printf(importStatements([]string{"golang.org/x/net/context", packageName + "/models"}))
	interfaceGenerator.Printf("type Controller interface {\n")

	var controllerGenerator Generator
	controllerGenerator.Printf("package server\n\n")
	controllerGenerator.Printf(importStatements([]string{"golang.org/x/net/context",
		"errors", packageName + "/models"}))

	// TODO: Better name for this... very java-y. Also shouldn't necessarily be controller
	// TODO: Should we plug this in more nicely??
	controllerGenerator.Printf("type ControllerImpl struct{\n")
	controllerGenerator.Printf("}\n")

	for _, path := range paths.Paths {
		for _, op := range pathItemOperations(path) {
			definition := fmt.Sprintf("%s(ctx context.Context, input *models.%sInput) (models.%sOutput, error)",
				capitalize(op.ID), capitalize(op.ID), capitalize(op.ID))
			interfaceGenerator.Printf("\t%s\n", definition)

			controllerGenerator.Printf("func (c ControllerImpl) %s {\n", definition)
			controllerGenerator.Printf("\t// TODO: Implement me!\n")
			controllerGenerator.Printf("\treturn nil, errors.New(\"Not implemented\")\n")
			controllerGenerator.Printf("}\n")
		}
	}
	interfaceGenerator.Printf("}\n")

	if err := g.WriteFile("generated/models/inputs.go"); err != nil {
		return err
	}
	if err := interfaceGenerator.WriteFile("generated/server/interface.go"); err != nil {
		return err
	}
	return controllerGenerator.WriteFile("generated/server/controller.go")
}

// importStatements takes a list of import strings and converts them to a formatted imports
func importStatements(imports []string) string {
	if len(imports) == 0 {
		return ""
	}
	output := "import (\n"
	for _, importStr := range imports {
		output += fmt.Sprintf("\t\"%s\"\n", importStr)
	}
	output += ")\n\n"
	return output
}

func printInputStruct(g *Generator, op *spec.Operation) error {
	g.Printf("type %sInput struct {\n", capitalize(op.ID))

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
		g.Printf("\t%s %s\n", capitalize(param.Name), typeName)
	}
	g.Printf("}\n\n")

	return nil
}

func printInputValidation(g *Generator, op *spec.Operation) error {
	g.Printf("func (i %sInput) Validate() error{\n", capitalize(op.ID))

	// TODO: Right now we only support validation on complex types (schemas)
	for _, param := range op.Parameters {
		if param.In == "body" {
			g.Printf("\tif err := i.%s.Validate(nil); err != nil {\n", capitalize(param.Name))
			g.Printf("\t\treturn err\n")
			g.Printf("\t}\n\n")
		}
	}
	g.Printf("\treturn nil\n")
	g.Printf("}\n\n")

	return nil
}

func printNewInput(g *Generator, op *spec.Operation) error {
	g.Printf("func New%sInput(r *http.Request) (*models.%sInput, error) {\n",
		capitalize(op.ID), capitalize(op.ID))

	g.Printf("\tvar input models.%sInput\n\n", capitalize(op.ID))

	printedError := false

	for _, param := range op.Parameters {

		if param.In != "body" {
			extractCode := ""
			switch param.In {
			case "query":
				extractCode = fmt.Sprintf("r.URL.Query().Get(\"%s\")", param.Name)
			case "path":
				extractCode = fmt.Sprintf("mux.Vars(r)[\"%s\"]", param.Name)
			case "header":
				extractCode = fmt.Sprintf("r.Header.Get(\"%s\")", param.Name)
			}

			g.Printf("\t%sStr := %s\n", param.Name, extractCode)

			if param.Type != "string" {
				if !printedError {
					g.Printf("\tvar err error\n")
					printedError = true
				}

				switch param.Type {
				case "integer":
					g.Printf("\tinput.%s, err = strconv.ParseInt(%sStr, 10, 64)\n",
						capitalize(param.Name), param.Name)
				case "number":
					g.Printf("\tinput.%s, err = strconv.ParseFloat(%sStr, 64)\n",
						capitalize(param.Name), param.Name)
				case "boolean":
					g.Printf("\tinput.%s, err = strconv.ParseBool(%sStr)\n",
						capitalize(param.Name), param.Name)
				default:
					return fmt.Errorf("Param type %s not supported", param.Type)
				}
				// TODO: These error message aren't great. We should probalby clean up...
				g.Printf("\tif err != nil {\n")
				g.Printf("\t\treturn nil, err\n")
				g.Printf("\t}\n")

			} else {
				g.Printf("\tinput.%s = %sStr\n", capitalize(param.Name), param.Name)
			}
		}
	}
	g.Printf("\n")

	g.Printf("\treturn &input, nil\n")
	g.Printf("}\n\n")

	return nil
}

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

func generateOutputs(packageName string, paths *spec.Paths) error {
	var g Generator

	g.Printf("package models\n\n")

	g.Printf(defaultOutputTypes())

	for _, path := range paths.Paths {
		for _, op := range pathItemOperations(path) {

			// We classify response keys into three types:
			// 1. 200-399 - these are "success" responses and implement the Output interface
			// 	defined above
			// 2. 400-599 - these are "failure" responses and implement the error interface
			// 3. Default - this is defined as a 500 (TODO: decide if we want to keep this...)

			// Define the success interface
			g.Printf("type %sOutput interface {\n", capitalize(op.ID))
			g.Printf("\t%sStatus() int\n", capitalize(op.ID))
			g.Printf("\t// Data is the underlying model object. We know it is json serializable\n")
			g.Printf("\t%sData() interface{}\n", capitalize(op.ID))
			g.Printf("}\n\n")

			// Define the error interface
			g.Printf("type %sError interface {\n", capitalize(op.ID))
			g.Printf("\terror // Extend the error interface\n")
			g.Printf("\t%sStatusCode() int\n", capitalize(op.ID))
			g.Printf("}\n\n")

			for statusCode, response := range op.Responses.StatusCodeResponses {

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

				outputName := fmt.Sprintf("%s%dOutput", capitalize(op.ID), statusCode)
				typeName, err := typeFromSchema(response.Schema)
				if err != nil {
					return err
				}

				g.Printf("type %s %s\n\n", outputName, typeName)

				g.Printf("func (o %s) %sData() interface{} {\n", outputName, capitalize(op.ID))
				g.Printf("\treturn o\n")
				g.Printf("}\n\n")

				if statusCode < 400 {

					// TODO: Do we really want to have that as part of the interface?
					g.Printf("func (o %s) %sStatus() int {\n", outputName, capitalize(op.ID))
					// TODO: Use the right status code...
					g.Printf("\treturn %d\n", statusCode)
					g.Printf("}\n\n")

				} else {

					g.Printf("func (o %s) Error() string {\n", outputName)
					// TODO: Would it make sense to give this a constructor so we can have a more detailed
					// error message?
					g.Printf("\treturn \"Status Code: \" + \"%d\"\n", statusCode)
					g.Printf("}\n\n")

					g.Printf("func (o %s) %sStatusCode() int {\n", outputName, capitalize(op.ID))
					g.Printf("\treturn %d\n", statusCode)
					g.Printf("}\n\n")
				}
			}
		}
	}
	return g.WriteFile("generated/models/outputs.go")
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

func generateHandlers(packageName string, paths *spec.Paths) error {
	var g Generator

	g.Printf("package server\n\n")
	g.Printf(importStatements([]string{"golang.org/x/net/context", "github.com/gorilla/mux",
		"net/http", "strconv", "encoding/json", "strconv", packageName + "/models"}))

	// TODO: Make this not be a global variable
	g.Printf("var controller Controller\n\n")

	g.Printf(jsonMarshalString)

	for _, path := range paths.Paths {
		for _, op := range pathItemOperations(path) {

			tmpl, err := template.New("test").Parse(handlerTemplate)
			if err != nil {
				return err
			}
			err = tmpl.Execute(&g.buf, handlerOp{Op: capitalize(op.ID)})
			if err != nil {
				return err
			}

			if err := printNewInput(&g, op); err != nil {
				return err
			}
		}
	}

	return g.WriteFile("generated/server/handlers.go")
}

type handlerOp struct {
	Op string
}

var jsonMarshalString = `
func jsonMarshalNoError(i interface{}) string {
	bytes, err := json.Marshal(i)
	if err != nil {
		// This should never happen
		return ""
	}
	return string(bytes)
}
`

var handlerTemplate = `func {{.Op}}Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	input, err := New{{.Op}}Input(r)
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := controller.{{.Op}}(ctx, input)
	if err != nil {
		if respErr, ok := err.(models.{{.Op}}Error); ok {
			http.Error(w, respErr.Error(), respErr.{{.Op}}StatusCode())
			return
		} else {
			http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
			return
		}
	}

	respBytes, err := json.Marshal(resp.{{.Op}}Data())
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
`
