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
	g := swagger.Generator{PackageName: packageName}
	g.Printf("package server\n\n")
	g.Printf(swagger.ImportStatements([]string{"golang.org/x/net/context", packageName + "/models"}))
	g.Printf("type Controller interface {\n")

	for _, pathKey := range swagger.SortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, opKey := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			capOpID := swagger.Capitalize(op.ID)
			definition := fmt.Sprintf("%s(ctx context.Context, input *models.%sInput) (models.%sOutput, error)",
				capOpID, capOpID, capOpID)
			g.Printf("\t%s\n", definition)
		}
	}
	g.Printf("}\n")

	return g.WriteFile("server/interface.go")
}

func printNewInput(g *swagger.Generator, op *spec.Operation) error {
	capOpID := swagger.Capitalize(op.ID)
	g.Printf("func New%sInput(r *http.Request) (*models.%sInput, error) {\n",
		capOpID, capOpID)

	g.Printf("\tvar input models.%sInput\n\n", capOpID)
	g.Printf("\tvar err error\n")
	g.Printf("\t_ = err\n\n")
	g.Printf("\tformats := strfmt.Default\n")
	g.Printf("\t_ = formats\n\n")

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

			if param.Required {
				g.Printf("\tif len(%sStr) == 0{\n", param.Name)
				g.Printf("\t\treturn nil, errors.New(\"Parameter must be specified\")\n")
				g.Printf("\t}\n")
			}

			switch param.Type {
			case "integer":
				if param.Format == "int32" {
					g.Printf("\t%sTmp, err := swag.ConvertInt32(%sStr)\n", param.Name, param.Name)
				} else {
					g.Printf("\t%sTmp, err := swag.ConvertInt64(%sStr)\n", param.Name, param.Name)
				}
			case "number":
				if param.Format == "float" {
					g.Printf("\t%sTmp, err := swag.ConvertFloat32(%sStr)\n", param.Name, param.Name)
				} else {
					g.Printf("\t%sTmp, err := swag.ConvertFloat64(%sStr)\n", param.Name, param.Name)
				}

			case "boolean":
				g.Printf("\t%sTmp, err := strconv.ParseBool(%sStr)\n",
					param.Name, param.Name)
			case "string":
				switch param.Format {
				case "byte":
					g.Printf("\t%sTmpInterface, err := formats.Parse(\"byte\", %sStr)\n", param.Name, param.Name)
					g.Printf("\t%sTmp := %sTmpInterface.([]byte)\n", param.Name, param.Name)
				case "":
					g.Printf("\t%sTmp := %sStr\n", param.Name, param.Name)
				case "date":
					g.Printf("\t%sTmpInterface, err := formats.Parse(\"date\", %sStr)\n", param.Name, param.Name)
					g.Printf("\t%sTmp := %sTmpInterface.(strfmt.Date)\n", param.Name, param.Name)
				case "date-time":
					g.Printf("\t%sTmpInterface, err := formats.Parse(\"date-time\", %sStr)\n", param.Name, param.Name)
					g.Printf("\t%sTmp := %sTmpInterface.(strfmt.DateTime)\n", param.Name, param.Name)
				default:
					return fmt.Errorf("Param format %s not supported", param.Format)
				}
			default:
				return fmt.Errorf("Param type %s not supported", param.Type)
			}

			if param.Required {
				g.Printf("\tinput.%s = %sTmp\n\n", capParamName, param.Name)
			} else {
				g.Printf("\tinput.%s = &%sTmp\n\n", capParamName, param.Name)
			}

		} else {
			if param.Schema == nil {
				return fmt.Errorf("Body parameters must have a schema defined")
			}
			g.Printf("\terr = json.NewDecoder(r.Body).Decode(input.%s)\n",
				capParamName)

		}
		// TODO: Figure out how to handle schemas here...
		if param.Required || param.Schema != nil {
			g.Printf("\tif err != nil {\n")
		} else {
			g.Printf("\tif err != nil && len(%sStr) != 0 {\n", param.Name)
		}
		// TODO: Make this an error of the right type
		g.Printf("\t\treturn nil, err\n")
		g.Printf("\t}\n")
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
		"net/http", "strconv", "encoding/json", "strconv", packageName + "/models", "errors",
		"github.com/go-openapi/strfmt", "github.com/go-openapi/swag"}))

	g.Printf("var _ = strconv.ParseInt\n")
	g.Printf("var _ = strfmt.Default\n")
	g.Printf("var _ = swag.ConvertInt32\n\n")

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
