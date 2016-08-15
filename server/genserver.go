package server

import (
	"bytes"
	"fmt"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/swagger"

	"text/template"
)

func Generate(packageName string, swagger spec.Swagger) error {

	if err := generateRouter(packageName, swagger.BasePath, swagger.Paths); err != nil {
		return err
	}

	// TODO: Is this really the way I want to do this???
	if err := generateContextsAndControllers(packageName, swagger.Paths); err != nil {
		return err
	}
	if err := generateHandlers(packageName, swagger.Paths); err != nil {
		return err
	}
	return nil
}

func generateRouter(packageName string, basePath string, paths *spec.Paths) error {
	g := swagger.Generator{PackageName: packageName}

	// TODO: Add something to all these about being auto-generated

	g.Printf(
		`package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	gContext "github.com/gorilla/context"
	"golang.org/x/net/context"

	"gopkg.in/tylerb/graceful.v1"
)

type contextKey struct{}

type Server struct {
	Handler http.Handler
	addr string
}

func (s Server) Serve() error {
	// Give the sever 30 seconds to shut down
	return graceful.RunWithErr(s.addr,30*time.Second,s.Handler)
}

func New(c Controller, addr string) Server {
	controller = c // TODO: get rid of global variable?
	r := mux.NewRouter()
`)

	for _, path := range swagger.SortedPathItemKeys(paths.Paths) {
		pathItem := paths.Paths[path]
		pathItemOps := swagger.PathItemOperations(pathItem)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]
			// TODO: Note the coupling for the handler name here and in the handler function. Does that mean these should be
			// together? Probably...

			g.Printf("\n")
			tmpl, err := template.New("routerFunction").Parse(routerFunctionTemplate)
			if err != nil {
				return err
			}
			var tmpBuf bytes.Buffer
			err = tmpl.Execute(&tmpBuf, routerTemplate{Method: method, Path: basePath + path,
				HandlerName: swagger.Capitalize(op.ID)})
			if err != nil {
				return err
			}
			g.Printf(tmpBuf.String())
		}
	}
	// TODO: It's a bit weird that this returns a pointer that it modifies...
	g.Printf("\thandler := withMiddleware(r)\n")
	g.Printf("\treturn Server{Handler: handler, addr: addr}\n")
	g.Printf("}\n")

	return g.WriteFile("server/router.go")
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

	interfaceGenerator := swagger.Generator{PackageName: packageName}
	interfaceGenerator.Printf("package server\n\n")
	interfaceGenerator.Printf(swagger.ImportStatements([]string{"golang.org/x/net/context", packageName + "/models"}))
	interfaceGenerator.Printf("type Controller interface {\n")

	controllerGenerator := swagger.Generator{PackageName: packageName}
	controllerGenerator.Printf("package server\n\n")
	controllerGenerator.Printf(swagger.ImportStatements([]string{"golang.org/x/net/context",
		"errors", packageName + "/models"}))

	// TODO: Better name for this... very java-y. Also shouldn't necessarily be controller
	// TODO: Should we plug this in more nicely??
	controllerGenerator.Printf("type ControllerImpl struct{\n")
	controllerGenerator.Printf("}\n")

	for _, pathKey := range swagger.SortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, opKey := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			capOpID := swagger.Capitalize(op.ID)
			definition := fmt.Sprintf("%s(ctx context.Context, input *models.%sInput) (models.%sOutput, error)",
				capOpID, capOpID, capOpID)
			interfaceGenerator.Printf("\t%s\n", definition)

			controllerGenerator.Printf("func (c ControllerImpl) %s {\n", definition)
			controllerGenerator.Printf("\t// TODO: Implement me!\n")
			controllerGenerator.Printf("\treturn nil, errors.New(\"Not implemented\")\n")
			controllerGenerator.Printf("}\n")
		}
	}
	interfaceGenerator.Printf("}\n")

	if err := interfaceGenerator.WriteFile("server/interface.go"); err != nil {
		return err
	}
	return controllerGenerator.WriteFile("server/controller.go")
}

func printNewInput(g *swagger.Generator, op *spec.Operation) error {
	capOpID := swagger.Capitalize(op.ID)
	g.Printf("func New%sInput(r *http.Request) (*models.%sInput, error) {\n",
		capOpID, capOpID)

	g.Printf("\tvar input models.%sInput\n\n", capOpID)

	printedError := false

	for _, param := range op.Parameters {

		capParamName := swagger.Capitalize(param.Name)
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
						capParamName, param.Name)
				case "number":
					g.Printf("\tinput.%s, err = strconv.ParseFloat(%sStr, 64)\n",
						capParamName, param.Name)
				case "boolean":
					g.Printf("\tinput.%s, err = strconv.ParseBool(%sStr)\n",
						capParamName, param.Name)
				default:
					return fmt.Errorf("Param type %s not supported", param.Type)
				}
				// TODO: These error message aren't great. We should probalby clean up...
				g.Printf("\tif err != nil {\n")
				g.Printf("\t\treturn nil, err\n")
				g.Printf("\t}\n")

			} else {
				g.Printf("\tinput.%s = %sStr\n", capParamName, param.Name)
			}
		} else {
			if param.Schema == nil {
				return fmt.Errorf("Body parameters must have a schema defined")
			}
			g.Printf("\tif err := json.NewDecoder(r.Body).Decode(&input.%s); err != nil {\n",
				capParamName)
			// TODO: Make this an error of the right type
			g.Printf("\t\treturn nil, err\n")
			g.Printf("\t}\n")
		}
	}
	g.Printf("\n")

	g.Printf("\treturn &input, nil\n")
	g.Printf("}\n\n")

	return nil
}

func generateHandlers(packageName string, paths *spec.Paths) error {
	g := swagger.Generator{PackageName: packageName}

	g.Printf("package server\n\n")
	g.Printf(swagger.ImportStatements([]string{"golang.org/x/net/context", "github.com/gorilla/mux",
		"net/http", "strconv", "encoding/json", "strconv", packageName + "/models"}))

	g.Printf("var _ = strconv.ParseInt\n\n")

	// TODO: Make this not be a global variable
	g.Printf("var controller Controller\n\n")

	g.Printf(jsonMarshalString)

	for _, pathKey := range swagger.SortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, opKey := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			tmpl, err := template.New("test").Parse(handlerTemplate)
			if err != nil {
				return err
			}
			var tmpBuf bytes.Buffer
			err = tmpl.Execute(&tmpBuf, handlerOp{Op: swagger.Capitalize(op.ID)})
			if err != nil {
				return err
			}
			g.Printf(tmpBuf.String())

			if err := printNewInput(&g, op); err != nil {
				return err
			}
		}
	}

	return g.WriteFile("server/handlers.go")
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
