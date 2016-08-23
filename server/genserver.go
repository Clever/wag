package server

import (
	"bytes"
	"fmt"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/models"
	"github.com/Clever/wag/swagger"

	"text/template"
)

func Generate(packageName string, s spec.Swagger) error {

	if err := generateRouter(packageName, s, s.Paths); err != nil {
		return err
	}
	if err := generateInterface(packageName, s.Paths); err != nil {
		return err
	}
	if err := generateHandlers(packageName, s.Paths); err != nil {
		return err
	}
	return nil
}

func generateRouter(packageName string, s spec.Swagger, paths *spec.Paths) error {
	g := swagger.Generator{PackageName: packageName}

	g.Printf(
		`package server

// Code auto-generated. Do not edit.

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

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

type handler struct {
	Controller
}

func New(c Controller, addr string) Server {
	r := mux.NewRouter()
	h := handler{Controller: c}
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
			err = tmpl.Execute(&tmpBuf, routerTemplate{Method: method, Path: s.BasePath + path,
				HandlerName: swagger.Capitalize(op.ID)})
			if err != nil {
				return err
			}
			g.Printf(tmpBuf.String())
		}
	}
	g.Printf("\thandler := withMiddleware(\"%s\", r)\n", s.Info.Title)
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
		h.{{.HandlerName}}Handler(r.Context(), w, r)
	})
`

func generateInterface(packageName string, paths *spec.Paths) error {
	g := swagger.Generator{PackageName: packageName}
	g.Printf("package server\n\n")
	g.Printf(swagger.ImportStatements([]string{"golang.org/x/net/context", packageName + "/models"}))
	g.Printf("//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_controller.go -package=server\n\n")
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
			} else if param.Default != nil {
				g.Printf("\tif len(%sStr) == 0 {\n", param.Name)
				g.Printf("\t\t// Use the default value\n")
				g.Printf("\t\t%sStr = \"%s\"\n", param.Name, swagger.DefaultAsString(param))
				g.Printf("\t}\n")
			}

			g.Printf("\tif len(%sStr) != 0 {\n", param.Name)

			typeName, err := swagger.ParamToType(param)
			if err != nil {
				return err
			}
			typeCode, err := swagger.StringToTypeCode(fmt.Sprintf("%sStr", param.Name), param)
			if err != nil {
				return err
			}
			g.Printf("\t\tvar %sTmp %s\n", param.Name, typeName)
			g.Printf("\t\t%sTmp, err = %s\n", param.Name, typeCode)
			g.Printf("\t\tif err != nil {\n")
			g.Printf("\t\t\treturn nil, err\n")
			g.Printf("\t\t}\n")

			if param.Required {
				g.Printf("\t\tinput.%s = %sTmp\n\n", capParamName, param.Name)
			} else {
				g.Printf("\t\tinput.%s = &%sTmp\n\n", capParamName, param.Name)
			}

			g.Printf("\t}\n")

		} else {
			if param.Schema == nil {
				return fmt.Errorf("Body parameters must have a schema defined")
			}
			typeName, err := models.TypeFromSchema(param.Schema)
			if err != nil {
				return err
			}

			g.Printf("\tdata, err := ioutil.ReadAll(r.Body)\n")

			if param.Required {
				g.Printf("\tif len(data) == 0 {\n")
				g.Printf("\t\treturn nil, errors.New(\"Parameter must be specified\")\n")
				g.Printf("\t}\n")
			}

			g.Printf("\tif len(data) > 0 {")
			// Initialize the pointer in the object
			g.Printf("\t\tinput.%s = &models.%s{}\n", capParamName, typeName)
			g.Printf("\t\tif err := json.NewDecoder(bytes.NewReader(data)).Decode(input.%s); err != nil {\n", capParamName)
			g.Printf("\t\t\treturn nil, err\n")
			g.Printf("\t\t}\n")
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
		"net/http", "strconv", "encoding/json", "strconv", packageName + "/models", "errors",
		"github.com/go-openapi/strfmt", "github.com/go-openapi/swag", "io/ioutil", "bytes"}))

	g.Printf("var _ = strconv.ParseInt\n")
	g.Printf("var _ = strfmt.Default\n")
	g.Printf("var _ = swag.ConvertInt32\n")
	g.Printf("var _ = errors.New\n")
	g.Printf("var _ = ioutil.ReadAll\n\n")

	g.Printf(swagger.BaseStringToTypeCode())
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

var handlerTemplate = `func (h handler) {{.Op}}Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.{{.Op}}(ctx, input)
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
