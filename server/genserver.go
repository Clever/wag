package server

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/swagger"

	"text/template"
)

func Generate(packageName string, s spec.Swagger) error {

	if err := generateRouter(packageName, s, s.Paths); err != nil {
		return err
	}
	if err := generateInterface(packageName, s.Info.InfoProps.Title, s.Paths); err != nil {
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
	"gopkg.in/Clever/kayvee-go.v4/logger"
	"gopkg.in/tylerb/graceful.v1"
)

type contextKey struct{}

// Server defines a HTTP server that implements the Controller interface.
type Server struct {
	// Handler should generally not be changed. It exposed to make testing easier.
	Handler http.Handler
	addr string
}

// Serve starts the server. It will return if an error occurs.
func (s Server) Serve() error {
	// Give the sever 30 seconds to shut down
	return graceful.RunWithErr(s.addr,30*time.Second,s.Handler)
}

type handler struct {
	Controller
}

// New returns a Server that implements the Controller interface. It will start when "Serve" is called.
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
			err = tmpl.Execute(&tmpBuf, routerTemplate{
				Method:      method,
				Path:        s.BasePath + path,
				HandlerName: swagger.Capitalize(op.ID),
				OpID:        op.ID,
			})
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
	OpID        string
}

var routerFunctionTemplate = `	r.Methods("{{.Method}}").Path("{{.Path}}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "{{.OpID}}")
		h.{{.HandlerName}}Handler(r.Context(), w, r)
	})
`

func generateInterface(packageName string, serviceName string, paths *spec.Paths) error {
	g := swagger.Generator{PackageName: packageName}
	g.Printf("package server\n\n")
	g.Printf(swagger.ImportStatements([]string{"context", packageName + "/models"}))
	g.Printf("//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_controller.go -package=server\n\n")
	g.Printf("// Controller defines the interface for the %s service.\n", serviceName)
	g.Printf("type Controller interface {\n")

	for _, pathKey := range swagger.SortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			g.Printf("\t%s\n", swagger.InterfaceComment(method, pathKey, pathItemOps[method]))
			g.Printf("\t%s\n\n", swagger.Interface(pathItemOps[method]))
		}
	}
	g.Printf("}\n")

	return g.WriteFile("server/interface.go")
}

func lowercase(input string) string {
	return strings.ToLower(input[0:1]) + input[1:]
}

func printNewInput(g *swagger.Generator, op *spec.Operation) error {
	capOpID := swagger.Capitalize(op.ID)
	g.Printf("// new%sInput takes in an http.Request an returns the input struct.\n", capOpID)
	singleSchemaedBodyParameter, opModel := swagger.SingleSchemaedBodyParameter(op)
	if singleSchemaedBodyParameter {
		g.Printf("func new%sInput(r *http.Request) (*models.%s, error) {\n",
			capOpID, opModel)
		g.Printf("\tvar input models.%s\n\n", opModel)
	} else {
		g.Printf("func new%sInput(r *http.Request) (*models.%sInput, error) {\n",
			capOpID, capOpID)
		g.Printf("\tvar input models.%sInput\n\n", capOpID)
	}

	g.Printf("\tvar err error\n")
	g.Printf("\t_ = err\n\n")

	for _, param := range op.Parameters {

		camelParamName := swagger.StructParamName(param)
		paramVarName := lowercase(camelParamName)

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
			g.Printf("\t%sStr := %s\n", paramVarName, extractCode)

			if param.Required {
				g.Printf("\tif len(%sStr) == 0{\n", paramVarName)
				g.Printf("\t\treturn nil, errors.New(\"Parameter must be specified\")\n")
				g.Printf("\t}\n")
			} else if param.Default != nil {
				g.Printf("\tif len(%sStr) == 0 {\n", paramVarName)
				g.Printf("\t\t// Use the default value\n")
				g.Printf("\t\t%sStr = \"%s\"\n", paramVarName, swagger.DefaultAsString(param))
				g.Printf("\t}\n")
			}

			g.Printf("\tif len(%sStr) != 0 {\n", paramVarName)

			typeName, err := swagger.ParamToType(param, false)
			if err != nil {
				return err
			}
			typeCode, err := swagger.StringToTypeCode(fmt.Sprintf("%sStr", paramVarName), param)
			if err != nil {
				return err
			}
			g.Printf("\t\tvar %sTmp %s\n", paramVarName, typeName)
			g.Printf("\t\t%sTmp, err = %s\n", paramVarName, typeCode)
			g.Printf("\t\tif err != nil {\n")
			g.Printf("\t\t\treturn nil, err\n")
			g.Printf("\t\t}\n")

			// TODO: Factor this out...
			if param.Required || param.Type == "array" {
				g.Printf("\t\tinput.%s = %sTmp\n\n", camelParamName, paramVarName)
			} else {
				g.Printf("\t\tinput.%s = &%sTmp\n\n", camelParamName, paramVarName)
			}

			g.Printf("\t}\n")

		} else {
			if param.Schema == nil {
				return fmt.Errorf("Body parameters must have a schema defined")
			}
			typeName, err := swagger.TypeFromSchema(param.Schema, true)
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

			if singleSchemaedBodyParameter {
				g.Printf("\t\tif err := json.NewDecoder(bytes.NewReader(data)).Decode(&input); err != nil {\n")
			} else {
				// Initialize the pointer in the object
				g.Printf("\t\tinput.%s = &%s{}\n", camelParamName, typeName)
				g.Printf("\t\tif err := json.NewDecoder(bytes.NewReader(data)).Decode(input.%s); err != nil {\n", camelParamName)
			}
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
	g.Printf(swagger.ImportStatements([]string{"context", "github.com/gorilla/mux", "gopkg.in/Clever/kayvee-go.v4/logger",
		"net/http", "strconv", "encoding/json", "strconv", packageName + "/models", "errors",
		"github.com/go-openapi/strfmt", "github.com/go-openapi/swag", "io/ioutil", "bytes"}))

	g.Printf("var _ = strconv.ParseInt\n")
	g.Printf("var _ = strfmt.Default\n")
	g.Printf("var _ = swag.ConvertInt32\n")
	g.Printf("var _ = errors.New\n")
	g.Printf("var _ = mux.Vars\n")
	g.Printf("var _ = bytes.Compare\n")
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

			statusCodeLogic := ""
			successCodes := swagger.SuccessStatusCodes(op)
			if len(successCodes) == 1 {
				statusCodeLogic = fmt.Sprintf("%d", successCodes[0])
			} else {
				statusCodeLogic = fmt.Sprintf("resp.%sStatusCode()", swagger.Capitalize(op.ID))
			}

			var tmpBuf bytes.Buffer
			singleInputOp, _ := swagger.SingleSchemaedBodyParameter(op)
			err = tmpl.Execute(&tmpBuf, handlerOp{
				Op:                swagger.Capitalize(op.ID),
				SuccessReturnType: !swagger.NoSuccessType(op),
				StatusCode:        statusCodeLogic,
				HasParams:         len(op.Parameters) != 0,
				SingleInputOp:     singleInputOp,
			})
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

// handlerOp contains the template variables for the handlerTemplate
type handlerOp struct {
	Op                string
	SuccessReturnType bool
	StatusCode        string
	HasParams         bool
	SingleInputOp     bool
}

var handlerTemplate = `func (h handler) {{.Op}}Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
{{if .HasParams}}
	input, err := new{{.Op}}Input(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	err = input.Validate({{if .SingleInputOp}}nil{{end}})
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

{{if .SuccessReturnType}}
	resp, err := h.{{.Op}}(ctx, input)
{{else}}
	err = h.{{.Op}}(ctx, input)
{{end}}
{{else}}
{{if .SuccessReturnType}}
	resp, err := h.{{.Op}}(ctx)
{{else}}
	err := h.{{.Op}}(ctx)
{{end}}
{{end}}
	if err != nil {
		if respErr, ok := err.(models.{{.Op}}Error); ok {
			http.Error(w, respErr.Error(), respErr.{{.Op}}StatusCode())
			return
		}
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

{{if .SuccessReturnType}}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader({{.StatusCode}})
	w.Write(respBytes)
{{else}}
	w.WriteHeader({{.StatusCode}})
	w.Write([]byte(""))
{{end}}
}
`
