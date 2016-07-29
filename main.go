package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-swagger/go-swagger/generator"
	"gopkg.in/yaml.v2"

	"text/template"
)

type SwaggerDefinition struct {
	Type string `yaml:"type"`
	// TODO: Properties should probably be more tightly typed
	Properties map[string]map[string]string `yaml:"properties"`
}

// TODO: I don't like that we have to pass a Generator in here...
// TODO: Should this even be object oriented?
func (d SwaggerDefinition) Printf(g *Generator, name string) {
	// TODO: Handle different types...

	g.Printf("type %s struct {\n", name)
	for key, value := range d.Properties {
		// TODO: Error handling? Or verify that earlier??? Also have switch?
		propType := value["type"]
		structType := ""
		if propType == "string" {
			structType = "string"

		} else if propType == "integer" {
			// TODO: Distinguish between int32 and other possibilities...
			structType = "int"
		} else {
			panic(fmt.Sprintf("Type %s not supported", propType))
		}
		// TODO: Upper case
		g.Printf("\t%s %s `json:\"%s\"`\n", capitalize(key), structType, key)
	}
	g.Printf("}\n\n")

	g.Printf("func (v %s) Validate() error { return nil }\n", name)
}

type SwaggerOperation struct {
	OperationID string                     `yaml:"operationId"`
	Description string                     `yaml:"description"`
	Responses   map[string]SwaggerResponse `yaml:"responses"`
	Parameters  []SwaggerParameter         `yaml:"parameters"`
	// TODO: Will need to add parameters...
}

type SwaggerParameter struct {
	Name   string            `yaml:"name"`
	In     string            `yaml:"in"`
	Type   string            `yaml:"type"`
	Schema map[string]string `yaml:"schema"`
}

type SwaggerResponse struct {
	Description string `yaml:"description"`
	// TODO: Add more types to schema???
	Schema map[string]string `yaml:"schema"`
}

type Swagger struct {
	Definitions map[string]SwaggerDefinition           `yaml:"definitions"`
	Paths       map[string]map[string]SwaggerOperation `yaml:"paths"`
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

	if err := buildTypes(swagger.Definitions); err != nil {
		panic(err)
	}
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

}

type Generator struct {
	buf bytes.Buffer
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// TODO: Add a nice comment!
// TODO: Make this write out to a file...
func buildTypes(definitions map[string]SwaggerDefinition) error {

	// TODO: Verify that the types are correct. In particular make sure they have the right references...

	var g Generator
	g.Printf("package main\n\n")
	for name, definition := range definitions {
		definition.Printf(&g, name)
	}

	return ioutil.WriteFile("generated/types.go", g.buf.Bytes(), 0644)
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
			// TODO: Validate the method
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

	// This includes the interfaces...
	var g Generator

	g.Printf("package main\n\n")

	g.Printf("import (\n")
	g.Printf("\t\"net/http\"\n")
	g.Printf("\t\"golang.org/x/net/context\"\n")
	g.Printf("\t\"github.com/gorilla/mux\"\n")
	g.Printf("\t\"encoding/json\"\n")
	g.Printf(")\n\n")

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
	g.Printf("\ntype Controller interface {\n")

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
	g.Printf("}\n")

	return nil
}

func printInputValidation(g *Generator, op SwaggerOperation) error {
	g.Printf("func (i %sInput) Validate() error{\n", capitalize(op.OperationID))

	// TODO: Right now we only support validation on complex types (schemas)
	for _, param := range op.Parameters {
		if param.In == "body" {
			g.Printf("\tif err := i.%s.Validate(); err != nil {\n", capitalize(param.Name))
			g.Printf("\t\treturn err\n")
			g.Printf("\t}\n\n")
		}

	}
	// TODO: Add in any validation...
	g.Printf("\treturn nil\n")
	g.Printf("}\n")

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
			// var postBody PostRequest
			// if err := json.NewDecoder(r.Body).Decode(&postBody); err != nil {
			// 	return nil, httpwrapper.HTTPErrorf(http.StatusBadRequest, err.Error())
			//}
			g.Printf("\tif err := json.NewDecoder(r.Body).Decode(&input.%s); err != nil{\n",
				capitalize(param.Name))
			g.Printf("\t\treturn nil, err\n") // TODO: This should probably return a 400 or something
			g.Printf("\t}\n")
		} else {
			fmt.Errorf("Unsupported param type %s", param)
		}
	}
	g.Printf("\n")

	g.Printf("\treturn &input, nil\n")
	g.Printf("}\n")

	return nil
}

func buildOutputs(paths map[string]map[string]SwaggerOperation) error {
	var g Generator

	g.Printf("package main\n\n")

	for _, path := range paths {
		for _, op := range path {
			g.Printf("type %sOutput interface {\n", capitalize(op.OperationID))
			g.Printf("\t%sStatus() int\n", capitalize(op.OperationID))
			g.Printf("}\n\n")

			for key, response := range op.Responses {
				outputName := fmt.Sprintf("%s%sOutput", capitalize(op.OperationID), capitalize(key))

				typeName, err := typeFromSchema(response.Schema)
				if err != nil {
					return err
				}
				g.Printf("type %s %s\n\n", outputName, typeName)
				g.Printf("func (o %s) %sStatus() int {\n", outputName, capitalize(op.OperationID))
				// TODO: Use the right status code...
				g.Printf("\treturn 200\n")
				g.Printf("}\n\n")
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
		return ref[len("#/definitions/"):], nil
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
`
