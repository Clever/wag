package jsclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Clever/go-utils/stringset"
	"github.com/Clever/wag/v10/swagger"
	"github.com/Clever/wag/v10/templates"
	"github.com/Clever/wag/v10/utils"
	"github.com/go-openapi/spec"
)

// Generate generates a client
func Generate(modulePath string, s spec.Swagger) error {
	pkgName, ok := s.Info.Extensions.GetString("x-npm-package")
	if !ok {
		return errors.New("must provide 'x-npm-package' in the 'info' section of the swagger.yml")
	}
	unwrappedFetch, _ := s.Info.Extensions.GetBool("x-unwrapped-fetch")

	indexJSTmpl, packageJSONTmpl := wrappedIndexJSTmplStr, wrappedPackageJSONTmplStr
	if unwrappedFetch {
		indexJSTmpl, packageJSONTmpl = unwrappedIndexJSTmplStr, unwrappedPackageJSONTmplStr
	}

	tmplInfo := clientCodeTemplate{
		ClassName:   utils.CamelCase(s.Info.InfoProps.Title, true),
		PackageName: pkgName,
		ServiceName: s.Info.InfoProps.Title,
		Version:     s.Info.InfoProps.Version,
		Description: s.Info.InfoProps.Description,
	}

	for _, path := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		pathItem := s.Paths.Paths[path]
		pathItemOps := swagger.PathItemOperations(pathItem)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]
			if op.Deprecated {
				continue
			}
			methodCode, err := methodCode(s, op, method, path, unwrappedFetch)
			if err != nil {
				return err
			}
			tmplInfo.Methods = append(tmplInfo.Methods, methodCode)
		}
	}

	indexDTS, err := generateTypescriptTypes(s, unwrappedFetch)
	if err != nil {
		return err
	}

	errorsJS, err := generateErrorsFile(s)
	if err != nil {
		return err
	}

	indexJS, err := templates.WriteTemplate(indexJSTmpl, tmplInfo)
	if err != nil {
		return err
	}

	packageJSON, err := templates.WriteTemplate(packageJSONTmpl, tmplInfo)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(modulePath, "types.js"), []byte(errorsJS), 0644); err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(modulePath, "index.js"), []byte(indexJS), 0644); err != nil {
		return err
	}

	if ioutil.WriteFile(filepath.Join(modulePath, "package.json"), []byte(packageJSON), 0644); err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(modulePath, "index.d.ts"), []byte(indexDTS), 0644)
}

type clientCodeTemplate struct {
	PackageName string
	ClassName   string
	ServiceName string
	Version     string
	Description string
	Methods     []string
}

const singleParamMethodDefinitionTemplateString = `/**{{if .Description}}
   * {{.Description}}{{end}}{{range $param := .Params}}
   * @param {{if $param.JSDocType}}{{.JSDocType}} {{end}}{{$param.JSName}}{{if $param.Default}}={{$param.Default}}{{end}}{{if $param.Description}} - {{.Description}}{{end}}{{end}}
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:{{.ServiceName}}.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   {{- if .IterMethod}}
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   * @returns {function} iter.forEachAsync - takes in an async function, applies it to each resource
   {{- else}}
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill{{if .JSDocSuccessReturnType}} {{.JSDocSuccessReturnType}}{{else}} {*}{{end}}{{$ServiceName := .ServiceName}}{{range $response := .Responses}}{{if $response.IsError}}
   * @reject {module:{{$ServiceName}}.Errors.{{$response.Name}}}{{end}}{{end}}
   * @reject {Error}{{end}}
   */
  {{- if .IterMethod}}
  {{.MethodName}}({{range $param := .Params}}{{$param.JSName}}, {{end}}options) {
  {{- else}}
  {{.MethodName}}({{range $param := .Params}}{{$param.JSName}}, {{end}}options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._{{.MethodName}}, arguments), callback);
  }

  {{if .UnwrappedFetch}}async {{end}}_{{.MethodName}}({{range $param := .Params}}{{$param.JSName}}, {{end}}options, cb) {
  {{- end}}
    const params = {};{{range $param := .Params}}
    params["{{$param.JSName}}"] = {{$param.JSName}};{{end}}
`

const pluralParamMethodDefinitionTemplateString = `/**{{if .Description}}
   * {{.Description}}{{end}}
   * @param {Object} params{{range $param := .Params}}
   * @param {{if $param.JSDocType}}{{.JSDocType}} {{end}}{{if not $param.Required}}[{{end}}params.{{$param.JSName}}{{if $param.Default}}={{$param.Default}}{{end}}{{if not $param.Required}}]{{end}}{{if $param.Description}} - {{.Description}}{{end}}{{end}}
   * @param {object} [options]
   * @param {number} [options.timeout] - A request specific timeout
   * @param {Map<string, string | number>} [options.baggage] - A request-specific baggage to be propagated
   * @param {module:{{.ServiceName}}.RetryPolicies} [options.retryPolicy] - A request specific retryPolicy
   * @param {Object.<string, string>} [options.headers] - Additional headers to send with the request
   {{- if .IterMethod}}
   * @returns {Object} iter
   * @returns {function} iter.map - takes in a function, applies it to each resource, and returns a promise to the result as an array
   * @returns {function} iter.toArray - returns a promise to the resources as an array
   * @returns {function} iter.forEach - takes in a function, applies it to each resource
   * @returns {function} iter.forEachAsync - takes in an async function, applies it to each resource
   {{- else}}
   * @param {function} [cb]
   * @returns {Promise}
   * @fulfill{{if .JSDocSuccessReturnType}} {{.JSDocSuccessReturnType}}{{else}} {*}{{end}}{{$ServiceName := .ServiceName}}{{range $response := .Responses}}{{if $response.IsError}}
   * @reject {module:{{$ServiceName}}.Errors.{{$response.Name}}}{{end}}{{end}}
   * @reject {Error}{{end}}
   */
  {{- if .IterMethod}}
  {{.MethodName}}(params, options) {
  {{- else}}
  {{.MethodName}}(params, options, cb) {
    let callback = cb;
    if (!cb && typeof options === "function") {
      callback = options;
    }
    return applyCallback(this._hystrixCommand.execute(this._{{.MethodName}}, arguments), callback);
  }

  {{if .UnwrappedFetch}}async {{end}}_{{.MethodName}}(params, options, cb) {
  {{- end}}`

type paramMapping struct {
	JSName      string
	WagName     string
	Required    bool
	JSDocType   string
	Default     interface{}
	Description string
}

type responseMapping struct {
	StatusCode int
	Name       string
	IsError    bool
	IsNoData   bool
}

type methodTemplate struct {
	UnwrappedFetch           bool
	ServiceName              string
	Operation                string
	MethodName               string
	IterMethod               bool
	IterResourceAccessString string
	Description              string
	MethodDefinition         string
	Params                   []paramMapping
	Method                   string
	PathCode                 string
	Path                     string
	HeaderParams             []paramMapping
	PathParams               []paramMapping
	QueryParams              []paramMapping
	BodyParam                string
	Responses                []responseMapping
	JSDocSuccessReturnType   string
}

// This function takes in a swagger path such as "/path/goes/to/{location}/and/to/{other_Location}"
// and returns a string of javacript code such as "/path/goes/to/" + location + "/and/to/" + otherLocation.
func fillOutPath(path string) string {
	return replacePathParams(path, func(jsParam string) string { return jsParam })
}

// fillOutEncodedPath is like fillOutPath, but URI-encodes each path parameter, e.g.
// "/path/goes/to/" + encodeURIComponent(location) + "/and/to/" + encodeURIComponent(otherLocation).
func fillOutEncodedPath(path string) string {
	return replacePathParams(path, func(jsParam string) string { return "encodeURIComponent(" + jsParam + ")" })
}

func replacePathParams(path string, wrap func(jsParam string) string) string {
	paramRegex := regexp.MustCompile("({.+?})")
	paramNameRegex := regexp.MustCompile("{(.+?)}")
	return paramRegex.ReplaceAllStringFunc(path, func(param string) string {
		return paramNameRegex.ReplaceAllStringFunc(param, func(paramName string) string {
			return "\" + " + wrap("params."+utils.CamelCase(paramName, false)) + " + \""
		})
	})
}

func methodCode(s spec.Swagger, op *spec.Operation, method, path string, unwrappedFetch bool) (string, error) {
	basePath := s.BasePath
	tmplStr, pathCode := wrappedMethodTmplStr, fillOutPath(path)
	if unwrappedFetch {
		tmplStr, pathCode = unwrappedMethodTmplStr, fillOutEncodedPath(path)
	}
	tmplInfo := methodTemplate{
		UnwrappedFetch: unwrappedFetch,
		ServiceName:    s.Info.InfoProps.Title,
		Operation:      op.ID,
		MethodName:     op.ID, // might mutate to op.ID + "Iter" for paging methods
		Description:    op.Description,
		Method:         method,
		PathCode:       basePath + pathCode,
		Path:           basePath + path,
	}

	var successResponse *spec.Response
	for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
		if successResponse == nil && statusCode >= 200 && statusCode < 400 {
			r := op.Responses.StatusCodeResponses[statusCode]
			successResponse = &r
		}
		response := responseMapping{
			StatusCode: statusCode,
			IsError:    statusCode >= 400,
		}
		typeName, _ := swagger.OutputType(&s, op, statusCode)
		if strings.HasPrefix(typeName, "models.") {
			typeName = typeName[7:] // models.ResponseType -> ResponseType
		}
		response.Name = typeName
		if typeName == "" {
			response.IsNoData = true
		}
		tmplInfo.Responses = append(tmplInfo.Responses, response)
	}
	tmplInfo.JSDocSuccessReturnType = responseToJSDocReturnType(successResponse)

	for _, wagParam := range op.Parameters {
		param := paramMapping{
			JSName:      utils.CamelCase(wagParam.Name, false),
			WagName:     wagParam.Name,
			Required:    wagParam.Required,
			JSDocType:   paramToJSDocType(wagParam),
			Description: wagParam.Description,
			Default:     wagParam.Default,
		}

		tmplInfo.Params = append(tmplInfo.Params, param)
		switch wagParam.In {
		case "path":
			tmplInfo.PathParams = append(tmplInfo.PathParams, param)
		case "header":
			tmplInfo.HeaderParams = append(tmplInfo.HeaderParams, param)
		case "body": // Will only ever be a single bodyParam so we can just set here
			tmplInfo.BodyParam = param.JSName
		case "query":
			tmplInfo.QueryParams = append(tmplInfo.QueryParams, param)
		}
	}

	if err := fillMethodDefinition(op, &tmplInfo); err != nil {
		return "", err
	}

	res, err := templates.WriteTemplate(tmplStr, tmplInfo)
	if err != nil {
		return "", err
	}

	if _, hasPaging := swagger.PagingParam(op); hasPaging {
		tmplInfo.IterMethod = true
		tmplInfo.MethodName += "Iter"

		if err := fillMethodDefinition(op, &tmplInfo); err != nil {
			return "", err
		}

		resourcePath := swagger.PagingResourcePath(op)
		if len(resourcePath) > 0 {
			tmplInfo.IterResourceAccessString = "." + strings.Join(resourcePath, ".")
		}

		iterMethodCode, err := templates.WriteTemplate(tmplStr, tmplInfo)
		if err != nil {
			return "", err
		}
		res += "\n" + iterMethodCode
	}

	return res, nil
}

func fillMethodDefinition(op *spec.Operation, tmplInfo *methodTemplate) error {
	var err error
	var methodDefinition string
	if len(op.Parameters) <= 1 {
		methodDefinition, err = templates.WriteTemplate(singleParamMethodDefinitionTemplateString, tmplInfo)
	} else {
		methodDefinition, err = templates.WriteTemplate(pluralParamMethodDefinitionTemplateString, tmplInfo)
	}
	if err != nil {
		return err
	}
	tmplInfo.MethodDefinition = methodDefinition
	return err
}

// paramToJSDocType returns the JSDoc type to assign a parameter.
// It returns empty string if it cannot discern a type for the parameter.
func paramToJSDocType(param spec.Parameter) string {
	if param.Type == "string" || param.Type == "number" || param.Type == "boolean" {
		return "{" + param.Type + "}"
	} else if param.Type == "integer" {
		return "{number}"
	} else if param.Type == "array" && param.Items != nil &&
		(param.Items.Type == "string" || param.Items.Type == "number" || param.Items.Type == "boolean") {
		return fmt.Sprintf("{%s[]}", param.Items.Type)
	}
	log.Printf("TODO: unhandled param name=%s. Documentation will be incomplete for this parameter.", param.Name)
	return ""
}

// schemaToJSDocType returns the JSDoc type to assign a schema.
// It returns empty string if it cannot discern a type for the schema.
func schemaToJSDocType(schema *spec.Schema) string {
	if len(schema.Type) == 1 {
		typ := schema.Type[0]
		if typ == "string" {
			return "{string}"
		} else if typ == "integer" {
			return "{number}"
		}
	}
	log.Printf("TODO: unhandled schema type %v.", schema.Type)
	return ""
}

// responseToJSDocReturnType returns the JS Doc type for a response.
// It returns empty string if it cannot determine the type.
func responseToJSDocReturnType(r *spec.Response) string {
	if r.Schema == nil {
		return "{undefined}"
	} else if r.Schema.Type != nil && len(r.Schema.Type) == 1 && r.Schema.Type[0] == "array" &&
		r.Schema.Items != nil && r.Schema.Items.Schema != nil {
		return "{Object[]}" // in the future, a more specific type would be nice
	} else if r.Schema.Ref.String() != "" {
		return "{Object}" // in the future, a more specific type would be nice
	}
	log.Printf("TODO: unhandled response: %#v. Documentation will be incomplete for this response type.", *r)
	return ""
}

var typeTmplString = `module.exports.Errors = {};
{{$ServiceName := .ServiceName}}
{{range .ErrorTypes}}/**
 * {{.Name}}
 * @extends Error
 * @memberof module:{{$ServiceName}}
 * @alias module:{{$ServiceName}}.Errors.{{.Name}}{{range .JSDocProperties}}
 * @property {{.Type}} {{.Name}}{{end}}
 */
module.exports.Errors.{{ .Name }} = class extends Error {
  constructor(body) {
    super(body.message);
    for (const k of Object.keys(body)) {
      this[k] = body[k];
    }
  }
};

{{ end }}`

type typesTemplate struct {
	ServiceName string
	ErrorTypes  []errorType
}

type errorType struct {
	StatusCode      int
	Name            string
	JSDocProperties []jsDocProperty
}

type jsDocProperty struct {
	Type string
	Name string
}

func jsDocPropertyFromSchema(name string, schema *spec.Schema) jsDocProperty {
	return jsDocProperty{
		Name: name,
		Type: schemaToJSDocType(schema),
	}
}

func generateErrorsFile(s spec.Swagger) (string, error) {
	typesTmpl := typesTemplate{
		ServiceName: s.Info.InfoProps.Title,
	}

	typeNames := stringset.New()

	for _, pathKey := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		path := s.Paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, opKey := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
				if statusCode < 400 {
					continue
				}
				typeName, _ := swagger.OutputType(&s, op, statusCode)
				if strings.HasPrefix(typeName, "models.") {
					typeName = typeName[7:]
				}
				if typeNames.Contains(typeName) {
					continue
				}
				typeNames.Add(typeName)

				etype := errorType{
					StatusCode: statusCode,
					Name:       typeName,
				}

				if schema, ok := s.Definitions[typeName]; !ok {
					log.Printf("TODO: could not find schema for %s, JS documentation will be incomplete", typeName)
				} else if len(schema.Properties) > 0 {
					for _, name := range swagger.SortedSchemaProperties(schema) {
						propertySchema := schema.Properties[name]
						etype.JSDocProperties = append(etype.JSDocProperties, jsDocPropertyFromSchema(name, &propertySchema))
					}
				}

				typesTmpl.ErrorTypes = append(typesTmpl.ErrorTypes, etype)
			}
		}
	}

	return templates.WriteTemplate(typeTmplString, typesTmpl)
}

// JSType is a wrapper type for things that represent a TypeScript type
type JSType string

// JSTypeMap is a representation of a TypeScript object
type JSTypeMap map[string]JSType

type typescriptTypes struct {
	UnwrappedFetch bool
	ServiceName    string
	IncludedTypes  []string
	MethodDecls    []string
	ErrorTypes     []string
}

var isDefaultIncludedType = map[string]bool{
	"BadRequest":    true,
	"InternalError": true,
	"NotFound":      true,
}

var primitiveTypes = map[string]string{
	"string":  "string",
	"integer": "number",
	"number":  "number",
	"boolean": "boolean",
}

func generateTypescriptTypes(s spec.Swagger, unwrappedFetch bool) (string, error) {
	tt := typescriptTypes{
		UnwrappedFetch: unwrappedFetch,
		ServiceName:    utils.CamelCase(s.Info.InfoProps.Title, true),
		IncludedTypes:  []string{},
		MethodDecls:    []string{},
	}
	includedTypeMap := JSTypeMap{}

	for _, path := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		pathItem := s.Paths.Paths[path]
		pathItemOps := swagger.PathItemOperations(pathItem)
		for _, method := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[method]
			if op.Deprecated {
				continue
			}
			methodDecl, err := methodDecl(s, op, path, method)
			if err != nil {
				return "", err
			}
			tt.MethodDecls = append(tt.MethodDecls, methodDecl)
			err = addInputType(&includedTypeMap, op)
			if err != nil {
				return "", err
			}
		}
	}

	for name, schema := range s.Definitions {
		if !isDefaultIncludedType[name] {
			theType, err := asJSType(&schema, "")
			if err != nil {
				return "", err
			}
			includedTypeMap[name] = theType
		}
	}

	var keys []string
	for k := range includedTypeMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	includedTypes := []string{}
	for _, typeName := range keys {
		includedTypes = append(includedTypes, fmt.Sprintf("type %s = %s;", typeName, includedTypeMap[typeName]))
	}
	tt.IncludedTypes = includedTypes

	errorTypes, err := getErrorTypes(s)
	if err != nil {
		return "", err
	}
	tt.ErrorTypes = errorTypes

	types, err := templates.WriteTemplate(typescriptTmplStr, tt)
	if err != nil {
		return "", err
	}

	return types, nil
}

func getErrorTypes(s spec.Swagger) ([]string, error) {
	typeNames := map[string]struct{}{}
	errorTypes := []string{}

	for _, pathKey := range swagger.SortedPathItemKeys(s.Paths.Paths) {
		path := s.Paths.Paths[pathKey]
		pathItemOps := swagger.PathItemOperations(path)
		for _, opKey := range swagger.SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]
			for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
				if statusCode < 400 {
					continue
				}
				typeName, _ := swagger.OutputType(&s, op, statusCode)
				if strings.HasPrefix(typeName, "models.") {
					typeName = typeName[7:]
				}

				if _, exists := typeNames[typeName]; exists {
					continue
				}
				typeNames[typeName] = struct{}{}

				if schema, ok := s.Definitions[typeName]; !ok {
					errorTypes = append(errorTypes, fmt.Sprintf("class %s {}", typeName))
				} else if len(schema.Properties) > 0 {
					declaration, err := generateErrorDeclaration(&schema, typeName, "models.")
					if err != nil {
						return errorTypes, err
					}
					errorTypes = append(errorTypes, declaration)
				}
			}
		}
	}
	return errorTypes, nil
}

func methodDecl(s spec.Swagger, op *spec.Operation, path, method string) (string, error) {
	returnType, err := ReturnType(s, op)
	if returnType != "void" && returnType != "never" {
		returnType = JSType(fmt.Sprintf("models.%s", returnType))
	}
	if err != nil {
		return "", err
	}
	methodName := op.ID
	var params string
	var methodDecl string
	if len(op.Parameters) == 0 {
		params = ""
	} else if len(op.Parameters) == 1 {
		paramName := utils.CamelCase(op.Parameters[0].Name, false)
		var paramType JSType
		if op.Parameters[0].ParamProps.Schema != nil {
			paramType, err = asJSType(op.Parameters[0].ParamProps.Schema, "")
			paramType = JSType(fmt.Sprintf("models.%s", paramType))
		} else {
			paramType, err = asJSTypeSimple(op.Parameters[0].SimpleSchema)
		}
		if err != nil {
			return "", err
		}
		if op.Parameters[0].Required {
			params = fmt.Sprintf("%s: %s, ", paramName, paramType)
		} else {
			params = fmt.Sprintf("%s?: %s, ", paramName, paramType)
		}
	} else {
		paramType := fmt.Sprintf("%sParams", utils.CamelCase(methodName, true))
		params = fmt.Sprintf("params: models.%s, ", paramType)
	}
	methodDecl = fmt.Sprintf("%s(%soptions?: RequestOptions, cb?: Callback<%s>): Promise<%s>",
		methodName, params, returnType, returnType)

	if _, hasPaging := swagger.PagingParam(op); hasPaging {
		resourcePath := swagger.PagingResourcePath(op)
		for _, ident := range resourcePath {
			returnType += JSType(fmt.Sprintf("[\"%s\"]", ident))
		}
		returnType = "ArrayInner<" + returnType + ">"
		methodDecl += fmt.Sprintf("\n  %s(%soptions?: RequestOptions): IterResult<%s>", methodName+"Iter", params, returnType)
	}
	return methodDecl, nil
}

type paramDeclTmpl struct {
	TypeName string
	Fields   string
}

func addInputType(jsTypeMap *JSTypeMap, op *spec.Operation) error {
	if len(op.Parameters) <= 1 {
		return nil
	}
	typeName := utils.CamelCase(op.ID+"Params", true)
	paramNames := []string{}
	fields := JSTypeMap{}
	for _, param := range op.Parameters {
		if param.In == "formData" {
			return fmt.Errorf("input parameters with 'In' formData are not supported")
		}
		paramType, err := paramToJSType(param)
		if err != nil {
			return err
		}
		paramType = ": " + paramType
		if !param.Required {
			paramType = "?" + paramType
		}
		paramNames = append(paramNames, param.Name)
		fields[param.Name] = paramType
	}
	fieldsStrings := []string{}
	for _, paramName := range paramNames {
		jsName := utils.CamelCase(paramName, false)
		fieldsStrings = append(fieldsStrings, fmt.Sprintf("%s%s;", jsName, fields[paramName]))
	}
	(*jsTypeMap)[typeName] = JSType(fmt.Sprintf("{\n  %s\n}", strings.Join(fieldsStrings, "\n  ")))
	return nil
}

const paramDeclTmlpStr = `interface {{.TypeName}} {
  {{.Fields}}
}`

const methodDeclTmplStr = `{{.Name}}({{range $index, $param := .Params}}{{if $index}}, {{end}}{{$param}}{{end}}): {{.ReturnType}}`

// ReturnType returns the methods return type
func ReturnType(s spec.Swagger, op *spec.Operation) (JSType, error) {
	successCodes := []int{}
	for statusCode := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			successCodes = append(successCodes, statusCode)
		}
	}
	if len(successCodes) == 0 {
		return "never", nil
	} else if len(successCodes) == 1 {
		return typeOf(s, op, successCodes[0])
	}
	return "", fmt.Errorf("Operation %s has more than one possible success return type", op.ID)
}

func typeOf(s spec.Swagger, op *spec.Operation, statusCode int) (JSType, error) {
	schema := swagger.OutputSchema(&s, op, statusCode)
	if schema == nil {
		return "void", nil
	}
	return asJSType(schema, "")
}

func paramToJSType(param spec.Parameter) (JSType, error) {
	if param.In == "body" {
		typeName, err := asJSType(param.Schema, "")
		if err != nil {
			return "", err
		}
		return typeName, nil
	}

	if len(param.Enum) > 0 {
		return typeFromEnum(param.Enum)
	}

	var typeName string
	switch param.Type {
	case "string":
		typeName = "string"
	case "integer", "number":
		typeName = "number"
	case "boolean":
		typeName = "boolean"
	case "array":
		if param.Items.Type != "string" {
			return JSType(""), fmt.Errorf("array parameters must have string sub-types")
		}
		typeName = "string[]"
	default:
		// Note. We don't support 'array' or 'file' types even though they're in the
		// Swagger spec.
		return JSType(""), fmt.Errorf("unsupported param type: \"%s\"", param.Type)
	}
	return JSType(typeName), nil
}

func asJSTypeSimple(simpleSchema spec.SimpleSchema) (JSType, error) {
	if jsType, ok := primitiveTypes[simpleSchema.Type]; ok {
		return JSType(jsType), nil
	}

	if simpleSchema.Type == "array" {
		return JSType("any[]"), nil
	}

	if simpleSchema.Type == "object" {
		return JSType("any"), nil
	}

	return JSType(""), fmt.Errorf("Unknown type '%v'", simpleSchema.Type)
}

func asJSType(schema *spec.Schema, refPrefix string) (JSType, error) {
	if schema == nil {
		return JSType(""), fmt.Errorf("No schema")
	}

	if len(schema.AllOf) > 0 {
		subTypes := []string{}
		for _, subSchema := range schema.AllOf {
			subType, err := asJSType(&subSchema, refPrefix)
			if err != nil {
				return JSType(""), err
			}
			subTypes = append(subTypes, string(subType))
		}
		return JSType(strings.Join(subTypes, " & ")), nil
	}

	if schema.Ref.String() != "" {
		def, err := defFromRef(schema.Ref.String())
		if err != nil {
			return JSType(""), err
		}
		if refPrefix != "" {
			return JSType(fmt.Sprintf("%s%s", refPrefix, def)), nil
		}
		return JSType(def), nil
	}

	if len(schema.Type) == 0 {
		if schema.AdditionalProperties != nil {
			if schema.AdditionalProperties.Schema != nil {
				innerType, err := asJSType(schema.AdditionalProperties.Schema, refPrefix)
				if err != nil {
					return JSType(""), err
				}
				return "{ [key: string]: " + innerType + " }", nil
			}
			if schema.AdditionalProperties.Allows {
				return "{ [key: string]: any }", nil
			}
		}
		return JSType("any"), nil
	}

	if len(schema.Enum) > 0 {
		return typeFromEnum(schema.Enum)
	}

	if len(schema.Type) > 1 {
		return JSType(""), fmt.Errorf("Having mltiple types in schema is not supported")
	}
	if jsType, ok := primitiveTypes[schema.Type[0]]; ok {
		return JSType(jsType), nil
	}

	if schema.Type[0] == "array" {
		innerType, err := asJSType(schema.Items.Schema, refPrefix)
		return innerType + "[]", err
	}

	if schema.Type[0] == "object" || len(schema.Properties) > 0 {
		fieldsStrings, err := generatePropertyDeclarations(schema, refPrefix)
		if err != nil {
			return JSType(""), err
		}

		if schema.AdditionalProperties != nil {
			if schema.AdditionalProperties.Schema != nil {
				innerType, err := asJSType(schema.AdditionalProperties.Schema, refPrefix)
				if err != nil {
					return JSType(""), err
				}
				fieldsStrings = append(fieldsStrings, "[key: string]: "+string(innerType)+";")
			} else if schema.AdditionalProperties.Allows {
				fieldsStrings = append(fieldsStrings, "[key: string]: any;")
			}
		}
		return JSType(fmt.Sprintf("{\n  %s\n}", strings.Join(fieldsStrings, "\n  "))), nil
	}

	return JSType(""), fmt.Errorf("Unknown type '%v'", schema.Type[0])
}

func typeFromEnum(enum []interface{}) (JSType, error) {
	enums := []string{}
	for _, enum := range enum {
		e, err := json.Marshal(enum)
		if err != nil {
			return JSType(""), err
		}
		enums = append(enums, string(e))
	}
	return JSType("(" + strings.Join(enums, " | ") + ")"), nil
}

func defFromRef(ref string) (string, error) {
	if strings.HasPrefix(ref, "#/definitions/") {
		return ref[len("#/definitions/"):], nil
	}
	return "", fmt.Errorf("schema.$ref has undefined reference type \"%s\". "+
		"Must start with #/definitions or #/responses.", ref)
}

const typescriptTmplStr = `import { Logger } from "kayvee";

type Callback<R> = (err: Error, result: R) => void;
type ArrayInner<R> = R extends (infer T)[] ? T : never;

interface RetryPolicy {
  backoffs(): number[];
{{- if .UnwrappedFetch}}
  retry(request: {method: string, url: string}, err: Error | undefined, res: Response | undefined, body: string | undefined): boolean;
{{- else}}
  retry(requestOptions: {method: string}, err: Error, res: {statusCode: number}): boolean;
{{- end}}
}

interface RequestOptions {
  timeout?: number;
  baggage?: Map<string, string | number>;
  retryPolicy?: RetryPolicy;
  headers?: { [key: string]: string };
}

interface IterResult<R> {
  map<T>(f: (r: R) => T, cb?: Callback<T[]>): Promise<T[]>;
  toArray(cb?: Callback<R[]>): Promise<R[]>;
  forEach(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
  forEachAsync(f: (r: R) => void, cb?: Callback<void>): Promise<void>;
}

interface CircuitOptions {
  forceClosed?: boolean;
  maxConcurrentRequests?: number;
  requestVolumeThreshold?: number;
  sleepWindow?: number;
  errorPercentThreshold?: number;
}

interface GenericOptions {
  timeout?: number;
  baggage?: Map<string, string | number>;
{{- if not .UnwrappedFetch}}
  keepalive?: boolean;
{{- end}}
  retryPolicy?: RetryPolicy;
  logger?: Logger;
  circuit?: CircuitOptions;
  serviceName?: string;
  asynclocalstore?: object;
}

interface DiscoveryOptions {
  discovery: true;
  address?: undefined;
}

interface AddressOptions {
  discovery?: false;
  address: string;
}

type {{.ServiceName}}Options = (DiscoveryOptions | AddressOptions) & GenericOptions;

import models = {{.ServiceName}}.Models

declare class {{.ServiceName}} {
  constructor(options: {{.ServiceName}}Options);

  close(): void;
  {{range .MethodDecls}}
  {{.}}
  {{end}}
}

declare namespace {{.ServiceName}} {
  const RetryPolicies: {
    Single: RetryPolicy;
    Exponential: RetryPolicy;
    None: RetryPolicy;
  }

  const DefaultCircuitOptions: CircuitOptions;

  namespace Errors {
    interface ErrorBody {
      message: string;
      [key: string]: any;
    }

    {{range .ErrorTypes}}
    {{.}}
    {{end}}
  }

  namespace Models {
    {{range .IncludedTypes}}
    {{.}}
    {{end}}
  }
}

export = {{.ServiceName}};
`
