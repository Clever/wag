package server

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/swagger"
	"github.com/Clever/wag/templates"
	"github.com/Clever/wag/utils"
)

// Generate server package for a swagger spec.
func Generate(packageName string, s spec.Swagger) error {

	if err := generateRouter(packageName, s, s.Paths); err != nil {
		return err
	}
	if err := generateInterface(packageName, &s, s.Info.InfoProps.Title, s.Paths); err != nil {
		return err
	}
	if err := generateHandlers(packageName, &s, s.Paths); err != nil {
		return err
	}
	return nil
}

type routerFunction struct {
	Method      string
	Path        string
	HandlerName string
	OpID        string
}

type routerTemplate struct {
	Title     string
	Functions []routerFunction
}

var routerTemplateStr = `
package server

// Code auto-generated. Do not edit.

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
	"path"

	"github.com/gorilla/mux"
	lightstep "github.com/lightstep/lightstep-tracer-go"
	opentracing "github.com/opentracing/opentracing-go"
	"gopkg.in/Clever/kayvee-go.v6/logger"
	kvMiddleware "gopkg.in/Clever/kayvee-go.v6/middleware"
	"gopkg.in/tylerb/graceful.v1"
	"github.com/Clever/go-process-metrics/metrics"
	"github.com/kardianos/osext"
)

type contextKey struct{}

// Server defines a HTTP server that implements the Controller interface.
type Server struct {
	// Handler should generally not be changed. It exposed to make testing easier.
	Handler http.Handler
	addr string
	l logger.KayveeLogger
}

// Serve starts the server. It will return if an error occurs.
func (s *Server) Serve() error {

	go func() {
		metrics.Log("{{.Title}}", 1*time.Minute)
	}()

	go func() {
		// This should never return. Listen on the pprof port
		log.Printf("PProf server crashed: %%s", http.ListenAndServe(":6060", nil))
	}()

	dir, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.SetGlobalRouting(path.Join(dir, "kvconfig.yml")); err != nil {
		s.l.Info("please provide a kvconfig.yml file to enable app log routing")
	}

	if lightstepToken := os.Getenv("LIGHTSTEP_ACCESS_TOKEN"); lightstepToken != "" {
		tags := make(map[string]interface{})
		tags[lightstep.ComponentNameKey] = "{{.Title}}"
		lightstepTracer := lightstep.NewTracer(lightstep.Options{
			AccessToken: lightstepToken,
			Tags:        tags,
			UseGRPC:     true,
		})
		defer lightstep.FlushLightStepTracer(lightstepTracer)
		opentracing.InitGlobalTracer(lightstepTracer)
	} else {
		s.l.Error("please set LIGHTSTEP_ACCESS_TOKEN to enable tracing")
	}

	s.l.Counter("server-started")

	// Give the sever 30 seconds to shut down
	return graceful.RunWithErr(s.addr,30*time.Second,s.Handler)
}

type handler struct {
	Controller
}

func withMiddleware(serviceName string, router http.Handler, m []func(http.Handler) http.Handler) http.Handler {
	handler := router

	// Wrap the middleware in the opposite order specified so that when called then run
	// in the order specified
	for i := len(m) - 1; i >= 0; i-- {
		handler = m[i](handler)
	}
	handler = TracingMiddleware(handler)
	handler = PanicMiddleware(handler)
	// Logging middleware comes last, i.e. will be run first.
	// This makes it so that other middleware has access to the logger
	// that kvMiddleware injects into the request context.
	handler = kvMiddleware.New(handler, serviceName)
	return handler
}


// New returns a Server that implements the Controller interface. It will start when "Serve" is called.
func New(c Controller, addr string) *Server {
	return NewWithMiddleware(c, addr, []func(http.Handler) http.Handler{})
}

// NewWithMiddleware returns a Server that implemenets the Controller interface. It runs the
// middleware after the built-in middleware (e.g. logging), but before the controller methods.
// The middleware is executed in the order specified. The server will start when "Serve" is called.
func NewWithMiddleware(c Controller, addr string, m []func(http.Handler) http.Handler) *Server {
	router := mux.NewRouter()
	h := handler{Controller: c}

	l := logger.New("{{.Title}}")

	{{range $index, $val := .Functions}}
	router.Methods("{{$val.Method}}").Path("{{$val.Path}}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.FromContext(r.Context()).AddContext("op", "{{$val.OpID}}")
		h.{{$val.HandlerName}}Handler(r.Context(), w, r)
		ctx := WithTracingOpName(r.Context(), "{{$val.OpID}}")
		r = r.WithContext(ctx)
	})
	{{end}}

	handler := withMiddleware("{{.Title}}", router, m)
	return &Server{Handler: handler, addr: addr, l: l}
}
`

func generateRouter(packageName string, s spec.Swagger, paths *spec.Paths) error {

	var template routerTemplate
	template.Title = s.Info.Title
	for _, path := range swagger.SortedPathItemKeys(paths.Paths) {
		pathItem := paths.Paths[path]
		pathItemOps := swagger.PathItemOperations(pathItem)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]

			template.Functions = append(template.Functions, routerFunction{
				Method:      method,
				Path:        s.BasePath + gorillaPathFor(op, path),
				HandlerName: swagger.Capitalize(op.ID),
				OpID:        op.ID,
			})
		}
	}

	routerCode, err := templates.WriteTemplate(routerTemplateStr, template)
	if err != nil {
		return err
	}
	g := swagger.Generator{PackageName: packageName}
	g.Printf(routerCode)
	return g.WriteFile("server/router.go")
}

func gorillaPathFor(op *spec.Operation, path string) string {
	for _, param := range op.Parameters {
		if param.In != "path" || param.Pattern == "" {
			continue
		}
		plainName := fmt.Sprintf("{%s}", param.Name)
		nameWithRegex := fmt.Sprintf("{%s:%s}", param.Name, param.Pattern)
		path = strings.Replace(path, plainName, nameWithRegex, 1)
	}
	return path
}

type interfaceTemplate struct {
	Comment    string
	Definition string
}

type interfaceFileTemplate struct {
	ImportStatements string
	ServiceName      string
	Interfaces       []interfaceTemplate
}

var interfaceTemplateStr = `
package server

{{.ImportStatements}}

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_controller.go -package=server

// Controller defines the interface for the {{.ServiceName}} service.
type Controller interface {

	{{range $interface := .Interfaces}}
		{{$interface.Comment}}
		{{$interface.Definition}}
	{{end}}
}
`

func generateInterface(packageName string, s *spec.Swagger, serviceName string, paths *spec.Paths) error {

	tmpl := interfaceFileTemplate{
		ImportStatements: swagger.ImportStatements([]string{"context", packageName + "/models"}),
		ServiceName:      serviceName,
	}

	for _, pathKey := range swagger.SortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			interfaceComment, err := swagger.InterfaceComment(method, pathKey, false, s, pathItemOps[method])
			if err != nil {
				return err
			}
			tmpl.Interfaces = append(tmpl.Interfaces, interfaceTemplate{
				Comment:    interfaceComment,
				Definition: swagger.Interface(s, pathItemOps[method]),
			})
		}
	}

	interfaceCode, err := templates.WriteTemplate(interfaceTemplateStr, tmpl)
	if err != nil {
		return err
	}
	g := swagger.Generator{PackageName: packageName}
	g.Printf(interfaceCode)
	return g.WriteFile("server/interface.go")
}

func lowercase(input string) string {
	return strings.ToLower(input[0:1]) + input[1:]
}

type handlerFileTemplate struct {
	ImportStatements string
	// TODO: Think about possibly factoring this out...
	BaseStringToTypeCode string
	Handlers             []string
}

var handlerFileTemplateStr = `
package server

{{.ImportStatements}}

var _ = strconv.ParseInt
var _ = strfmt.Default
var _ = swag.ConvertInt32
var _ = errors.New
var _ = mux.Vars
var _ = bytes.Compare
var _ = ioutil.ReadAll

{{.BaseStringToTypeCode}}

func jsonMarshalNoError(i interface{}) string {
	bytes, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		// This should never happen
		return ""
	}
	return string(bytes)
}

{{ range $handler := .Handlers }}
	{{ $handler }}
{{end}}
`

func generateHandlers(packageName string, s *spec.Swagger, paths *spec.Paths) error {

	tmpl := handlerFileTemplate{
		ImportStatements: swagger.ImportStatements([]string{"context", "github.com/gorilla/mux", "gopkg.in/Clever/kayvee-go.v6/logger",
			"net/http", "strconv", "encoding/json", "strconv", packageName + "/models",
			"github.com/go-openapi/strfmt", "github.com/go-openapi/swag", "io/ioutil", "bytes",
			"github.com/go-errors/errors"}),
		BaseStringToTypeCode: swagger.BaseStringToTypeCode(),
	}

	for _, pathKey := range swagger.SortedPathItemKeys(paths.Paths) {
		path := paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, opKey := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]

			operationHandler, err := generateOperationHandler(s, op)
			if err != nil {
				return err
			}
			tmpl.Handlers = append(tmpl.Handlers, operationHandler)
		}
	}

	handlerCode, err := templates.WriteTemplate(handlerFileTemplateStr, tmpl)
	if err != nil {
		return err
	}
	g := swagger.Generator{PackageName: packageName}
	g.Printf(handlerCode)
	return g.WriteFile("server/handlers.go")
}

var jsonMarshalString = `

`

// generateOperationHandler generates the handler code for a single handler
func generateOperationHandler(s *spec.Swagger, op *spec.Operation) (string, error) {
	typeToCode := make(map[string]int)
	emptyResponseCode := 200
	codeToType := swagger.CodeToTypeMap(s, op, false)
	typeToCode, err := swagger.TypeToCodeMap(s, op)
	if err != nil {
		return "", err
	}
	if empty, ok := typeToCode[""]; ok {
		emptyResponseCode = empty
		delete(typeToCode, "")
	}

	singleSchemaedBodyParameter, _ := swagger.SingleSchemaedBodyParameter(op)
	singleStringPathParameter, singleStringPathParameterVarName := swagger.SingleStringPathParameter(op)
	successType := swagger.SuccessType(s, op)
	arraySuccessType := ""
	if successType != nil && strings.HasPrefix(*successType, "[]") {
		arraySuccessType = *successType
	}
	pagingParam, hasPaging := swagger.PagingParam(op)
	var pagingParamPointer bool
	if hasPaging {
		_, pagingParamPointer, err = swagger.ParamToType(pagingParam)
		if err != nil {
			return "", err
		}
	}
	inputVarName := "input"
	if singleStringPathParameter {
		inputVarName = singleStringPathParameterVarName
	}

	handlerOp := handlerOp{
		Op:                               swagger.Capitalize(op.ID),
		SuccessReturnType:                successType != nil,
		ArraySuccessType:                 arraySuccessType,
		HasParams:                        len(op.Parameters) != 0,
		InputVarName:                     inputVarName,
		SingleSchemaedBodyParameter:      singleSchemaedBodyParameter,
		HasPaging:                        hasPaging,
		PagingParamField:                 swagger.StructParamName(pagingParam),
		PagingParamPointer:               pagingParamPointer,
		EmptyStatusCode:                  emptyResponseCode,
		TypesToStatusCodes:               typeToCode,
		SingleStringPathParameter:        singleStringPathParameter,
		SingleStringPathParameterVarName: singleStringPathParameterVarName,
		StatusCodeToType:                 codeToType,
	}
	handlerCode, err := templates.WriteTemplate(handlerTemplate, handlerOp)
	if err != nil {
		return "", err
	}

	newInputCode, err := generateNewInput(op)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.WriteString(handlerCode)
	buf.WriteString(newInputCode)

	return buf.String(), nil
}

// handlerOp contains the template variables for the handlerTemplate
type handlerOp struct {
	Op                               string
	SuccessReturnType                bool
	ArraySuccessType                 string
	HasParams                        bool
	InputVarName                     string
	HasPaging                        bool
	PagingParamField                 string
	PagingParamPointer               bool
	SingleSchemaedBodyParameter      bool
	EmptyStatusCode                  int
	TypesToStatusCodes               map[string]int
	SingleStringPathParameter        bool
	SingleStringPathParameterVarName string
	StatusCodeToType                 map[int]string
}

var handlerTemplate = `
// statusCodeFor{{.Op}} returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeFor{{.Op}}(obj interface{}) int {

	switch obj.(type) {
	{{ range $type, $code := .TypesToStatusCodes }}
   	case {{$type}}:
   		return {{$code}}
	{{ end }}
	default:
		return -1
	}
}

func (h handler) {{.Op}}Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
{{if .HasParams}}
	{{.InputVarName}}, err := new{{.Op}}Input(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError({{index .StatusCodeToType 400}}{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	{{if .SingleStringPathParameter}}
		err = models.Validate{{.Op}}Input({{.SingleStringPathParameterVarName}})
	{{else}}
		err = input.Validate({{if .SingleSchemaedBodyParameter}}nil{{end}})
	{{end}}
		if err != nil {
			logger.FromContext(ctx).AddContext("error", err.Error())
			http.Error(w, jsonMarshalNoError({{index .StatusCodeToType 400}}{Message: err.Error()}), http.StatusBadRequest)
			return
		}
{{end}}
{{if .SuccessReturnType}}
	{{if .HasParams}}
		resp,{{if .HasPaging}} nextPageID,{{end}} err := h.{{.Op}}(ctx, {{.InputVarName}})
	{{else}}
		resp, err := h.{{.Op}}(ctx)
	{{end}}
	{{if gt (len .ArraySuccessType) 0}}
		// Success types that return an array should never return nil so let's make this easier
		// for consumers by converting nil arrays to empty arrays
		if resp == nil {
			resp = {{.ArraySuccessType}}{}
		}
	{{end}}
{{else}}
	{{if .HasParams}}
		err = h.{{.Op}}(ctx, {{.InputVarName}})
	{{else}}
		err := h.{{.Op}}(ctx)
	{{end}}
{{end}}
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		}
		statusCode := statusCodeFor{{.Op}}(err)
		if statusCode == -1 {
			err = {{index .StatusCodeToType 500}}{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

{{if .SuccessReturnType}}
	respBytes, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError({{index .StatusCodeToType 500}}{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	{{if .HasPaging}}
		if !swag.IsZero(nextPageID) {
			{{.InputVarName}}.{{.PagingParamField}} = {{if .PagingParamPointer}}&{{end}}nextPageID
			path, err := {{.InputVarName}}.Path()
			if err != nil {
				logger.FromContext(ctx).AddContext("error", err.Error())
				http.Error(w, jsonMarshalNoError({{index .StatusCodeToType 500}}{Message: err.Error()}), http.StatusInternalServerError)
				return
			}
			w.Header().Set("X-Next-Page-Path", path)
		}
	{{end}}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeFor{{.Op}}(resp))
	w.Write(respBytes)
{{else}}
	w.WriteHeader({{.EmptyStatusCode}})
	w.Write([]byte(""))
{{end}}
}
`

type singleStringPathParameterTemplateData struct {
	Op           string
	ParamName    string
	ParamVarName string
}

var singleStringPathParameterTemplate = `
// new{{.Op}}Input takes in an http.Request an returns the {{.ParamName}} parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func new{{.Op}}Input(r *http.Request) (string, error) {
	{{.ParamVarName}} := mux.Vars(r)["{{.ParamName}}"]
	if len({{.ParamVarName}}) == 0 {
		return "", errors.New("Parameter {{.ParamName}} must be specified")
	}
	return {{.ParamVarName}}, nil
}
`

func generateNewInput(op *spec.Operation) (string, error) {
	capOpID := swagger.Capitalize(op.ID)

	singleStringPathParameter, paramVarName := swagger.SingleStringPathParameter(op)
	if singleStringPathParameter {
		return templates.WriteTemplate(singleStringPathParameterTemplate,
			singleStringPathParameterTemplateData{
				Op:           capOpID,
				ParamName:    op.Parameters[0].Name,
				ParamVarName: paramVarName,
			})
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("// new%sInput takes in an http.Request an returns the input struct.\n", capOpID))
	singleSchemaedBodyParameter, opModel := swagger.SingleSchemaedBodyParameter(op)
	if singleSchemaedBodyParameter {
		buf.WriteString(fmt.Sprintf("func new%sInput(r *http.Request) (*models.%s, error) {\n",
			capOpID, opModel))
		buf.WriteString(fmt.Sprintf("\tvar input models.%s\n\n", opModel))
	} else {
		buf.WriteString(fmt.Sprintf("func new%sInput(r *http.Request) (*models.%sInput, error) {\n",
			capOpID, capOpID))
		buf.WriteString(fmt.Sprintf("\tvar input models.%sInput\n\n", capOpID))
	}

	buf.WriteString(fmt.Sprintf("\tvar err error\n"))
	buf.WriteString(fmt.Sprintf("\t_ = err\n\n"))

	for _, param := range op.Parameters {

		structFieldName := swagger.StructParamName(param)
		paramVarName := lowercase(utils.CamelCase(param.Name, false))

		typeName, pointer, err := swagger.ParamToType(param)
		if err != nil {
			return "", err
		}

		if param.In != "body" {
			if param.Type == "array" && param.In == "query" {
				buf.WriteString(fmt.Sprintf("\tif %s, ok := r.URL.Query()[\"%s\"]; ok {\n\t\tinput.%s = %s\n\t}\n",
					paramVarName, param.Name, structFieldName, paramVarName))
			} else {
				typeCode, err := swagger.StringToTypeCode(fmt.Sprintf("%sStr", paramVarName), param, op)
				if err != nil {
					return "", err
				}
				defaultVal := ""
				if param.Default != nil {
					defaultVal = swagger.DefaultAsString(param)
				}
				str, err := templates.WriteTemplate(paramTemplateStr, paramTemplate{
					Required:        param.Required,
					ParamType:       param.In,
					VarName:         paramVarName,
					ParamName:       param.Name,
					CapParamName:    structFieldName,
					TypeName:        typeName,
					TypeCode:        typeCode,
					DefaultValue:    defaultVal,
					PointerInStruct: pointer,
				})
				if err != nil {
					return "", err
				}
				buf.WriteString(str)
			}
		} else {
			if param.Schema == nil {
				return "", fmt.Errorf("body parameters must have a schema defined")
			}
			paramField := structFieldName
			if singleSchemaedBodyParameter {
				paramField = ""
			}
			str, err := templates.WriteTemplate(bodyParamTemplateStr, bodyParamTemplate{
				Required:   param.Required,
				ParamField: paramField,
				TypeName:   typeName,
			})
			if err != nil {
				return "", err
			}
			buf.WriteString(str)
		}
	}
	buf.WriteString(fmt.Sprintf("\n"))

	buf.WriteString(fmt.Sprintf("\treturn &input, nil\n"))
	buf.WriteString(fmt.Sprintf("}\n\n"))

	return buf.String(), nil
}

func capitalize(str string) string {
	return swagger.Capitalize(str)
}

type paramTemplate struct {
	Required        bool
	ParamType       string
	VarName         string
	ParamName       string
	CapParamName    string
	TypeName        string
	TypeCode        string
	DefaultValue    string
	PointerInStruct bool
}

var paramTemplateStr = `
	{{if eq .ParamType "query" -}}
		{{.VarName}}Strs := r.URL.Query()["{{.ParamName}}"]
		{{if .Required -}}
			if len({{.VarName}}Strs) == 0 {
				return nil, errors.New("parameter must be specified")
			}
		{{- end -}}
	{{- else if eq .ParamType "path" -}}
		{{.VarName}}Str := mux.Vars(r)["{{.ParamName}}"]
		if len({{.VarName}}Str) == 0 {
			return nil, errors.New("parameter must be specified")
		}
		{{.VarName}}Strs := []string{ {{.VarName}}Str }
	{{- else if eq .ParamType "header" -}}
		{{.VarName}}Strs := r.Header.Get("{{.ParamName}}")
		{{if .Required -}}
			if len({{.VarName}}Strs) == 0 {
				return nil, errors.New("parameter must be specified")
			}
		{{- end -}}
	{{- end}}
	{{if gt (len .DefaultValue) 0 -}}
		if len({{.VarName}}Strs) == 0 {
			{{.VarName}}Strs = []string{"{{.DefaultValue}}"}
		}
	{{- end}}
	if len({{.VarName}}Strs) > 0 {
			var {{.VarName}}Tmp {{.TypeName}}
		{{if eq .ParamType "header" -}}
			{{.VarName}}Tmp = {{.VarName}}Strs
		{{- else -}}
			{{.VarName}}Str := {{.VarName}}Strs[0]
			{{.VarName}}Tmp, err = {{.TypeCode}}
			if err != nil {
				return nil, err
			}
		{{- end}}
		{{if .PointerInStruct -}}
			input.{{.CapParamName}} = &{{.VarName}}Tmp
		{{- else -}}
			input.{{.CapParamName}} = {{.VarName}}Tmp
		{{- end}}
	}
`

type bodyParamTemplate struct {
	Required   bool
	ParamField string
	TypeName   string
}

var bodyParamTemplateStr = `
	data, err := ioutil.ReadAll(r.Body)
	{{if .Required}} if len(data) == 0 {
		return nil, errors.New("parameter must be specified")
	}{{end}}

	if len(data) > 0 { {{if eq (len .ParamField) 0}}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(&input); err != nil {
			return nil, err
		}{{else}}
		input.{{.ParamField}} = &models.{{.TypeName}}{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.{{.ParamField}}); err != nil {
			return nil, err
		}{{end}}
	}
`
