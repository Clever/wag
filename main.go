package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-swagger/go-swagger/generator"
	"gopkg.in/yaml.v2"

	"text/template"
)

type SwaggerOperation struct {
	// TODO: Implement
	Tags []string `yaml:"tags"`

	OperationID string                     `yaml:"operationId"`
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
}

// Validate checks if the swagger operation has any fields we don't support
func (s SwaggerOperation) Validate() error {
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

	return nil
}

type SwaggerParameter struct {
	Name string `yaml:"name"`
	In   string `yaml:"in"`
	// TODO: Support more types
	Type string `yaml:"type"`
	// The schema here has to be {$ref: "..."}. We don't support defining your own
	// schema here.
	Schema map[string]string `yaml:"schema"`
}

type SwaggerResponse struct {
	Description string `yaml:"description"`
	// Schema only supports schemas of the form {$ref: "..."}
	// TODO: Support {type: 'array', items: {$ref: "..."}}
	Schema map[string]string `yaml:"schema"`

	// TODO: Implement at least some part of the Header Object spec
	Header map[string]interface{} `yaml:"headers"`
}

type Swagger struct {
	// TODO: Field to implement or error on if invalid
	BasePath string                 `yaml:"basePath"`
	Tags     map[string]interface{} `yaml:"tags"`

	// Partially implemented
	Paths map[string]map[string]SwaggerOperation `yaml:"paths"`

	// We rely on the go-swagger code to generate all the definitions
	Definitions map[string]interface{} `yaml:"definitions"`

	// Fields we support, but only with a certain set of values
	Swagger  string   `yaml:"swagger"`
	Schemes  []string `yaml:"schemes"`
	Consumes []string `yaml:"consumers"`
	Produces []string `yaml:"produces"`

	// Fields we don't support
	Host                string                 `yaml:"host"`
	Parameters          map[string]interface{} `yaml:"parameters"`
	Responses           map[string]interface{} `yaml:"responses"`
	SecurityDefinitions map[string]interface{} `yaml:"securityDefinitions"`
	Security            map[string]interface{} `yaml:"security"`
}

// Validates returns an error if the swagger file is invalid or uses fields
// we don't support. Note that this isn't a comprehensive check for all things
// we don't support, so this may not return an error, but the Swagger file might
// have values we don't support
func (s Swagger) Validate() error {
	if s.Swagger != "2.0" {
		return fmt.Errorf("Unsupported Swagger version %s", s.Swagger)
	}

	if len(s.Schemes) != 1 || s.Schemes[0] != "http" {
		return fmt.Errorf("Schemes only supports 'http'")
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
	for fieldName, _ := range s.Paths {
		// Note that $ref and parameters are not valid as of now
		if !sliceContains([]string{"get", "put", "post", "delete", "options", "head", "patch"}, fieldName) {
			fmt.Errorf("Invalid path field name: %s", fieldName)
		}
	}

	// TODO: Verify operationID is unique

	return nil
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

	// generate models with go-swagger
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
	if err := generator.GenerateServer("", []string{}, []string{}, generator.GenOpts{
		Spec:           "test.yml",
		ModelPackage:   "models",
		Target:         "./generated/",
		IncludeModel:   true,
		IncludeHandler: false,
		IncludeSupport: false,
	}); err != nil {
		log.Fatal(err)
	}

	// TODO: Make this configurable
	bytes, err := ioutil.ReadFile("test.yml")
	if err != nil {
		panic(err)
	}

	var swagger Swagger
	if err := yaml.Unmarshal(bytes, &swagger); err != nil {
		panic(err)
	}

	fmt.Printf("Swagger: %+v\n", swagger)

	if err := buildRouter(swagger.Paths); err != nil {
		panic(err)
	}
	// TODO: Is this really the way I want to do this???
	if err := buildContextsAndControllers(swagger.Paths); err != nil {
		panic(err)
	}
	if err := buildHandlers(swagger.Paths); err != nil {
		panic(err)
	}
	if err := buildOutputs(swagger.Paths); err != nil {
		panic(err)
	}

	if err := generateClients(swagger); err != nil {
		panic(err)
	}

}

type Generator struct {
	buf bytes.Buffer
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func buildRouter(paths map[string]map[string]SwaggerOperation) error {
	var g Generator

	// TODO: Add something to all these about being auto-generated

	g.Printf(
		`package main

import (
	"net/http"

	"github.com/gorilla/mux"

	gContext "github.com/gorilla/context"
	"golang.org/x/net/context"	
)

type contextKey struct{}

func withRoutes(r *mux.Router) *mux.Router {`)

	for path, pathObj := range paths {
		for method, op := range pathObj {
			// TODO: Note the coupling for the handler name here and in the handler function. Does that mean these should be
			// together? Probably...

			g.Printf("\n")
			tmpl, err := template.New("routerFunction").Parse(routerFunctionTemplate)
			if err != nil {
				return err
			}
			err = tmpl.Execute(&g.buf, routerTemplate{Method: method, Path: path,
				HandlerName: capitalize(op.OperationID)})
			if err != nil {
				return err
			}
		}
	}
	// TODO: It's a bit weird that this returns a pointer that it modifies...
	g.Printf("\treturn r\n")
	g.Printf("}\n")

	return ioutil.WriteFile("generated/router.go", g.buf.Bytes(), 0644)
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

func buildContextsAndControllers(paths map[string]map[string]SwaggerOperation) error {

	// Only create the controller the first time. After that leave it as is
	// TODO: This isn't convenient when developing, so maybe have a flag...?
	// TODO: We still want to build the contexts.go file, so need to figure that out...
	// if _, err := os.Stat("generated/controller.go"); err == nil {
	//	return nil
	//}

	var g Generator

	g.Printf("package main\n\n")

	g.Printf("import (\n")
	g.Printf("\t\"net/http\"\n")
	g.Printf("\t\"golang.org/x/net/context\"\n")
	g.Printf("\t\"github.com/gorilla/mux\"\n")

	g.Printf("\t\"encoding/json\"\n")
	g.Printf("\t\"github.com/Clever/inter-service-api-testing/codegen-poc/generated/models\"\n")
	g.Printf(")\n\n")

	// These two imports are only used if we have body parameters, so if we don't have these the
	// compiler will complain about unused imports
	g.Printf("var _ = json.Marshal\n\n")
	// TODO: Not as easy to do this for models since we have to find some type. Will be easier if
	// we move it into the same package

	for _, path := range paths {
		for _, op := range path {
			if err := printInputStruct(&g, op); err != nil {
				return err
			}

			if err := printNewInput(&g, op); err != nil {
				return err
			}

			if err := printInputValidation(&g, op); err != nil {
				return err
			}
		}
	}

	// TODO: How should I name these things??? Should they be on a per-tag basis???
	g.Printf("type Controller interface {\n")

	var controllerGenerator Generator
	controllerGenerator.Printf("package main\n\n")
	controllerGenerator.Printf("import \"golang.org/x/net/context\"\n")
	controllerGenerator.Printf("import \"errors\"\n\n")
	// TODO: Better name for this... very java-y. Also shouldn't necessarily be controller
	// TODO: Should we plug this in more nicely??
	controllerGenerator.Printf("type ControllerImpl struct{\n")
	controllerGenerator.Printf("}\n")

	for _, path := range paths {
		for _, op := range path {
			definition := fmt.Sprintf("%s(ctx context.Context, input *%sInput) (%sOutput, error)",
				capitalize(op.OperationID), capitalize(op.OperationID), capitalize(op.OperationID))
			g.Printf("\t%s\n", definition)

			// TODO: We could add a nice comment here...
			controllerGenerator.Printf("func (c ControllerImpl) %s {\n", definition)
			controllerGenerator.Printf("\t// TODO: Implement me!\n")
			controllerGenerator.Printf("\treturn nil, errors.New(\"Not implemented\")\n")
			controllerGenerator.Printf("}\n")
		}
	}
	g.Printf("}\n")

	if err := ioutil.WriteFile("generated/contexts.go", g.buf.Bytes(), 0644); err != nil {
		return err
	}
	return ioutil.WriteFile("generated/controller.go", controllerGenerator.buf.Bytes(), 0644)
}

func printInputStruct(g *Generator, op SwaggerOperation) error {
	g.Printf("type %sInput struct {\n", capitalize(op.OperationID))

	// TODO: Explicitly disallow formData param type

	fmt.Printf("Parameters %+v\n", op.Parameters)

	for _, param := range op.Parameters {
		if param.In == "formData" {
			return fmt.Errorf("Input parameters with 'In' formData are not supported")
		}

		var typeName string
		if param.Type == "string" && param.In != "body" {
			typeName = "string"
		} else if param.In == "body" {
			var err error
			typeName, err = typeFromSchema(param.Schema)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Unsupported param types, at least not yet")
		}
		g.Printf("\t%s %s\n", capitalize(param.Name), typeName)
	}
	g.Printf("}\n\n")

	return nil
}

func printInputValidation(g *Generator, op SwaggerOperation) error {
	g.Printf("func (i %sInput) Validate() error{\n", capitalize(op.OperationID))

	// TODO: Right now we only support validation on complex types (schemas)
	for _, param := range op.Parameters {
		if param.In == "body" {
			g.Printf("\tif err := i.%s.Validate(nil); err != nil {\n", capitalize(param.Name))
			g.Printf("\t\treturn err\n")
			g.Printf("\t}\n\n")
		}

	}
	// TODO: Add in any validation...
	g.Printf("\treturn nil\n")
	g.Printf("}\n\n")

	return nil
}

func printNewInput(g *Generator, op SwaggerOperation) error {
	g.Printf("func New%sInput(r *http.Request) (*%sInput, error) {\n",
		capitalize(op.OperationID), capitalize(op.OperationID))

	g.Printf("\tvar input %sInput\n\n", capitalize(op.OperationID))

	for _, param := range op.Parameters {
		// TODO: Handle non-string types. This probably means extracting the
		// code that pulls the variable out

		// TODO: Switch statement
		if param.In == "query" {
			g.Printf("\tinput.%s = r.URL.Query().Get(\"%s\")\n",
				capitalize(param.Name), param.Name)
		} else if param.In == "path" {
			g.Printf("\tinput.%s = mux.Vars(r)[\"%s\"]\n",
				capitalize(param.Name), param.Name)
		} else if param.In == "header" {
			g.Printf("\tinput.%s = r.Header.Get(\"%s\")\n",
				capitalize(param.Name), param.Name)
		} else if param.In == "body" {
			// g.Printf("\tif err := json.NewDecoder(r.Body).Decode(&input.%s); err != nil{\n",
			// 	capitalize(param.Name))
			// g.Printf("\t\treturn nil, err\n") // TODO: This should probably return a 400 or something
			// g.Printf("\t}\n")
			// TODO: Make it like the above, and only do this is required = false
			g.Printf("\tjson.NewDecoder(r.Body).Decode(&input.%s)\n", capitalize(param.Name))
		} else {
			fmt.Errorf("Unsupported param type %s", param)
		}
	}
	g.Printf("\n")

	g.Printf("\treturn &input, nil\n")
	g.Printf("}\n\n")

	return nil
}

func buildOutputs(paths map[string]map[string]SwaggerOperation) error {
	var g Generator

	g.Printf("package main\n\n")
	// TODO: We're going to have to be smarter about these imports
	g.Printf("import \"github.com/Clever/inter-service-api-testing/codegen-poc/generated/models\"\n\n")

	for _, path := range paths {
		for _, op := range path {

			// We classify response keys into three types:
			// 1. 200-399 - these are "success" responses and implement the Output interface
			// 	defined above
			// 2. 400-599 - these are "failure" responses and implement the error interface
			// 3. Default - this is defined as a 500 (TODO: decide if we want to keep this...)

			// Define the success interface
			g.Printf("type %sOutput interface {\n", capitalize(op.OperationID))
			g.Printf("\t%sStatus() int\n", capitalize(op.OperationID))
			g.Printf("\t// Data is the underlying model object. We know it is json serializable\n")
			g.Printf("\t%sData() interface{}\n", capitalize(op.OperationID))
			g.Printf("}\n\n")

			// Define the error interface
			g.Printf("type %sError interface {\n", capitalize(op.OperationID))
			g.Printf("\terror // Extend the error interface\n")
			g.Printf("\t%sStatusCode() int\n", capitalize(op.OperationID))
			g.Printf("}\n\n")

			for key, response := range op.Responses {

				var statusCode int
				if key == "default" {
					statusCode = 500
				} else {
					statusCode64, err := strconv.ParseInt(key, 10, 32)
					if err != nil || statusCode64 < 200 || statusCode64 > 599 {
						// TODO: Write a test for this...
						return fmt.Errorf("Response map key must be an integer between 200 and 599 or "+
							"the string 'default'. Was %s", key)
					}
					statusCode = int(statusCode64)
				}

				outputName := fmt.Sprintf("%s%sOutput", capitalize(op.OperationID), capitalize(key))
				if statusCode < 400 {

					// TODO: Be a bit smarter with this...
					typeName, err := typeFromSchema(response.Schema)
					if err != nil {
						return err
					}
					g.Printf("type %s struct {\n", outputName) //, typeName)
					g.Printf("\tData %s\n", typeName)
					g.Printf("}\n\n")

					g.Printf("func (o %s) %sData() interface{} {\n", outputName, capitalize(op.OperationID))
					g.Printf("\treturn o.Data\n")
					g.Printf("}\n\n")

					// TODO: Do we really want to have that as part of the interface?
					g.Printf("func (o %s) %sStatus() int {\n", outputName, capitalize(op.OperationID))
					// TODO: Use the right status code...
					g.Printf("\treturn %d\n", statusCode)
					g.Printf("}\n\n")

				} else {

					// TODO: Should these by named operationID + "Error" (instead of + "Output")
					// Would need to change the client too
					g.Printf("type %s struct{}\n\n", outputName)
					g.Printf("func (o %s) Error() string {\n", outputName)
					// TODO: Would it make sense to give this a constructor so we can have a more detailed
					// error message?
					g.Printf("\treturn \"Status Code: \" + \"%d\"\n", statusCode)
					g.Printf("}\n\n")
					g.Printf("func (o %s) %sStatusCode() int {\n", outputName, capitalize(op.OperationID))
					g.Printf("\treturn %d", statusCode)
					g.Printf("}\n\n")
				}

			}
		}
	}

	return ioutil.WriteFile("generated/outputs.go", g.buf.Bytes(), 0644)
}

func typeFromSchema(schema map[string]string) (string, error) {
	if _, ok := schema["$ref"]; ok && len(schema) == 1 {
		ref, _ := schema["$ref"]

		// TODO: Handle references outside of '#/definitions'
		if !strings.HasPrefix(ref, "#/definitions/") {
			return "", fmt.Errorf("schema.$ref has undefined reference type. Must be #/definitions")
		}
		return "models." + ref[len("#/definitions/"):], nil
	} else {
		// TODO: Support more?
		return "", fmt.Errorf("Other ref types not supported")
	}
}

func buildHandlers(paths map[string]map[string]SwaggerOperation) error {
	var g Generator

	g.Printf("package main\n\n")
	g.Printf("import (\n")
	g.Printf("\t\"net/http\"\n")
	g.Printf("\t\"golang.org/x/net/context\"\n")
	g.Printf("\t\"encoding/json\"\n")
	g.Printf(")\n\n")

	// TODO: Make this not be a global variable
	g.Printf("var controller Controller\n\n")

	for _, path := range paths {
		for _, op := range path {

			tmpl, err := template.New("test").Parse(handlerTemplate)
			if err != nil {
				return err
			}
			err = tmpl.Execute(&g.buf, handlerOp{Op: capitalize(op.OperationID)})
			if err != nil {
				return err
			}
		}
	}

	return ioutil.WriteFile("generated/handlers.go", g.buf.Bytes(), 0644)
}

type handlerOp struct {
	Op string
}

var handlerTemplate = `func {{.Op}}Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	input, err := New{{.Op}}Input(r)
	if err != nil {
		// TODO: Think about this whether this is usually an internal error or it could
		// be from a bad request format...
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := controller.{{.Op}}(ctx, input)
	if err != nil {
		if respErr, ok := err.({{.Op}}Error); ok {
			http.Error(w, respErr.Error(), respErr.{{.Op}}StatusCode())
			return
		} else {
			// This is the default case
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	respBytes, err := json.Marshal(resp.{{.Op}}Data())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
`
