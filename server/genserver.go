package server

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-openapi/spec"

	"github.com/Clever/wag/v5/server/gendb"
	"github.com/Clever/wag/v5/swagger"
	"github.com/Clever/wag/v5/templates"
	"github.com/Clever/wag/v5/utils"
)

// Generate server package for a swagger spec.
func Generate(packageName, packagePath string, s spec.Swagger) error {

	if err := generateRouter(packageName, packagePath, s, s.Paths); err != nil {
		return err
	}
	if err := generateInterface(packageName, packagePath, &s, s.Info.InfoProps.Title, s.Paths); err != nil {
		return err
	}
	if err := generateHandlers(packageName, packagePath, &s, s.Paths); err != nil {
		return err
	}
	return gendb.GenerateDB(packageName, packagePath, &s, s.Info.InfoProps.Title, s.Paths)
}

type routerFunction struct {
	Method      string
	Path        string
	HandlerName string
	OpID        string
}

type routerTemplate struct {
	ImportStatements string
	Title            string
	Functions        []routerFunction
}

func generateRouter(packageName, packagePath string, s spec.Swagger, paths *spec.Paths) error {

	var template routerTemplate
	template.Title = s.Info.Title
	for _, path := range swagger.SortedPathItemKeys(paths.Paths) {
		pathItem := paths.Paths[path]
		pathItemOps := swagger.PathItemOperations(pathItem)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]

			template.Functions = append(template.Functions, routerFunction{
				Method:      method,
				Path:        s.BasePath + path,
				HandlerName: swagger.Capitalize(op.ID),
				OpID:        op.ID,
			})
		}
	}
	template.ImportStatements = swagger.ImportStatements([]string{
		"compress/gzip",
		"context",
		"log",
		"net/http",
		`_ "net/http/pprof"`,
		"os",
		"os/signal",
		"path",
		"syscall",
		"time",
		"github.com/Clever/go-process-metrics/metrics",
		packageName + "/tracing",
		"github.com/gorilla/handlers",
		"github.com/gorilla/mux",
		"github.com/kardianos/osext",
		"gopkg.in/Clever/kayvee-go.v6/logger",
		`kvMiddleware "gopkg.in/Clever/kayvee-go.v6/middleware"`,
	})
	routerCode, err := templates.WriteTemplate(routerTemplateStr, template)
	if err != nil {
		return err
	}
	g := swagger.Generator{PackagePath: packagePath}
	g.Printf(routerCode)
	return g.WriteFile("server/router.go")
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

//go:generate mockgen -source=$GOFILE -destination=mock_controller.go -package=server

// Controller defines the interface for the {{.ServiceName}} service.
type Controller interface {

	{{range $interface := .Interfaces}}
		{{$interface.Comment}}
		{{$interface.Definition}}
	{{end}}
}
`

func generateInterface(packageName, packagePath string, s *spec.Swagger, serviceName string, paths *spec.Paths) error {

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
	g := swagger.Generator{PackagePath: packagePath}
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
	bytes, err := json.Marshal(i)
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

func generateHandlers(packageName, packagePath string, s *spec.Swagger, paths *spec.Paths) error {

	tmpl := handlerFileTemplate{
		ImportStatements: swagger.ImportStatements([]string{"context", "github.com/gorilla/mux",
			"gopkg.in/Clever/kayvee-go.v6/logger",
			"net/http", "strconv", "encoding/json", "strconv", "fmt", packageName + "/models",
			"github.com/go-openapi/strfmt", "github.com/go-openapi/swag", "io/ioutil", "bytes",
			"github.com/go-errors/errors", "golang.org/x/xerrors",
		}),
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
	g := swagger.Generator{PackagePath: packagePath}
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
		{{if .SingleSchemaedBodyParameter}}
			if input != nil {
				err = input.Validate(nil)
			}
		{{else}}
			err = input.Validate()
		{{end}}
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
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%%+v", xerr))
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
	respBytes, err := json.Marshal(resp)
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
	} else {
		buf.WriteString(fmt.Sprintf("func new%sInput(r *http.Request) (*models.%sInput, error) {\n",
			capOpID, capOpID))
		buf.WriteString(fmt.Sprintf("\tvar input models.%sInput\n\n", capOpID))
	}

	buf.WriteString(fmt.Sprintf("\tvar err error\n"))
	buf.WriteString(fmt.Sprintf("\t_ = err\n"))

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

	if singleSchemaedBodyParameter {
		buf.WriteString(fmt.Sprintf("\treturn nil, nil\n"))
	} else {
		buf.WriteString(fmt.Sprintf("\treturn &input, nil\n"))
	}
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
				return nil, errors.New("query parameter '{{.ParamName}}' must be specified")
			}
		{{- end -}}
	{{- else if eq .ParamType "path" -}}
		{{.VarName}}Str := mux.Vars(r)["{{.ParamName}}"]
		if len({{.VarName}}Str) == 0 {
			return nil, errors.New("path parameter '{{.ParamName}}' must be specified")
		}
		{{.VarName}}Strs := []string{ {{.VarName}}Str }
	{{- else if eq .ParamType "header" -}}
		{{.VarName}}Strs := r.Header.Get("{{.ParamName}}")
		{{if .Required -}}
			if len({{.VarName}}Strs) == 0 {
				return nil, errors.New("request header '{{.ParamName}}' must be specified")
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
		return nil, errors.New("request body is required, but was empty")
	}{{end}}
	if len(data) > 0 {
		{{- if eq (len .ParamField) 0}}
			var input models.{{.TypeName}}
			if err := json.NewDecoder(bytes.NewReader(data)).Decode(&input); err != nil {
				return nil, err
			}
			return &input, nil
		{{- else}}
			input.{{.ParamField}} = &models.{{.TypeName}}{}
			if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.{{.ParamField}}); err != nil {
				return nil, err
			}
		{{- end}}
	}
`
