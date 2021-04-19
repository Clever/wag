package swagger

import (
	"fmt"
	"strings"

	"github.com/Clever/wag/v7/templates"
	"github.com/Clever/wag/v7/utils"
	"github.com/go-openapi/spec"
)

// Interface returns the interface for the server-side handler of an operation
func Interface(s *spec.Swagger, op *spec.Operation) string {
	return opInterface(s, op, true)
}

// ClientInterface returns the client-facing interface for an operation
func ClientInterface(s *spec.Swagger, op *spec.Operation) string {
	return opInterface(s, op, false)
}

// ClientIterInterface returns the client-facing interface for the iterator
// builder of an operation
func ClientIterInterface(s *spec.Swagger, op *spec.Operation) string {
	capOpID := Capitalize(op.ID)
	input := OperationInput(op)
	return fmt.Sprintf(
		"New%sIter(ctx context.Context, %s) (%sIter, error)",
		capOpID,
		input,
		capOpID,
	)
}

// OperationInput returns the input to an operation
func OperationInput(op *spec.Operation) string {
	// Don't add the input parameter argument unless there are some arguments.
	// If a method has a single input parameter, and it's a schema, make the
	// generated type for that schema the input of the method.
	// If a method has a single input parameter, and it's a simple type (string, TODO: int),
	// make that the input of the method.
	// If a method has multiple input parameters, wrap them in a struct.
	capOpID := Capitalize(op.ID)
	input := ""
	if singleSchemaedBodyParameter, opModel := SingleSchemaedBodyParameter(op); singleSchemaedBodyParameter {
		input = fmt.Sprintf("i *models.%s", opModel)
	} else if singleStringPathParameter, inputName := SingleStringPathParameter(op); singleStringPathParameter {
		input = fmt.Sprintf("%s string", inputName)
	} else if len(op.Parameters) > 0 {
		input = fmt.Sprintf("i *models.%sInput", capOpID)
	}
	return input
}

// generateInterface returns the interface for an operation
func opInterface(s *spec.Swagger, op *spec.Operation, includePaging bool) string {
	capOpID := Capitalize(op.ID)

	input := OperationInput(op)

	returnTypes := []string{}
	if successType := SuccessType(s, op); successType != nil {
		returnTypes = append(returnTypes, *successType)
	}
	if pagingParam, ok := PagingParam(op); includePaging && ok {
		pagingParamType, _, err := ParamToType(pagingParam)
		if err != nil {
			panic(fmt.Errorf("could not convert paging parameter to type for %s: %s", op.ID, err))
		}
		returnTypes = append(returnTypes, pagingParamType)
	}
	returnTypes = append(returnTypes, "error")

	var output string
	if len(returnTypes) == 1 {
		output = returnTypes[0]
	} else {
		output = fmt.Sprintf("(%s)", strings.Join(returnTypes, ", "))
	}

	return fmt.Sprintf("%s(ctx context.Context, %s) %s", capOpID, input, output)
}

// InterfaceComment returns the comment for the interface for the operation. If the client
// flag is set then it generates the client version. Otherwise it generates the server version.
func InterfaceComment(method, path string, client bool, s *spec.Swagger, op *spec.Operation) (string, error) {

	statusCodeToType := CodeToTypeMap(s, op, true)
	for code, typ := range statusCodeToType {
		if typ == "" {
			statusCodeToType[code] = "nil"
		}
	}
	_, hasPaging := PagingParam(op)
	tmpl := struct {
		OpID             string
		Method           string
		Path             string
		Client           bool
		HasPaging        bool
		Description      string
		StatusCodeToType map[int]string
	}{
		OpID:             Capitalize(op.ID),
		Method:           method,
		Path:             path,
		Client:           client,
		HasPaging:        hasPaging,
		Description:      op.Description,
		StatusCodeToType: statusCodeToType,
	}
	return templates.WriteTemplate(interfaceCommentTmplStr, tmpl)
}

var interfaceCommentTmplStr = `
{{if .Client -}}
// {{.OpID}} makes a {{.Method}} request to {{.Path}}
{{- else -}}
// {{.OpID}} handles {{.Method}} requests to {{.Path}}
{{- if .HasPaging}}
// Returns response object and the ID of the next page
{{- end}}
{{- end}}
// {{.Description}} {{ range $code, $type := .StatusCodeToType }}
// {{$code}}: {{$type}} {{end}}
// default: client side HTTP errors, for example: context.DeadlineExceeded.`

// OutputSchema returns the Swagger schema for an operation and statusCode
func OutputSchema(s *spec.Swagger, op *spec.Operation, statusCode int) *spec.Schema {
	resp := op.Responses.StatusCodeResponses[statusCode]
	if !strings.HasPrefix(resp.Ref.String(), "#/responses") {
		return resp.Schema
	}

	// Follow the pointer to the response and then to the schema
	refObj, _, err := resp.Ref.GetPointer().Get(s)
	if err != nil {
		panic("bad response schema reference")
	}
	r, ok := refObj.(spec.Response)
	if !ok {
		panic("bad response schema reference")
	}
	return r.Schema
}

// OutputType returns the output type for a given status code of an operation and whether it
// is a pointer in the interface.
func OutputType(s *spec.Swagger, op *spec.Operation, statusCode int) (string, bool) {
	schema := OutputSchema(s, op, statusCode)

	successType, err := TypeFromSchema(schema, true)
	if err != nil {
		panic(fmt.Errorf("could not convert operation to type for %s, %s", op.ID, err))
	}
	return successType, schema != nil && schema.Ref.String() != ""
}

// SuccessType returns the success type for the operation. If there is no success-type then
// it returns nil
func SuccessType(s *spec.Swagger, op *spec.Operation) *string {
	for statusCode := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			successType, makePointer := OutputType(s, op, statusCode)
			if successType == "" {
				return nil
			}
			if makePointer {
				successType = "*" + successType
			}
			return &successType
		}
	}
	return nil
}

// PagingParam returns the parameter that specifies the page ID for this
// operation, if paging is configured. If paging is not configured, the second
// return value is `false`.
func PagingParam(op *spec.Operation) (spec.Parameter, bool) {
	pagingConfig, ok := op.Extensions["x-paging"].(map[string]interface{})
	if !ok {
		return spec.Parameter{}, false
	}
	paramName, ok := pagingConfig["pageParameter"].(string)
	if !ok {
		panic(fmt.Errorf("x-paging does not include pageParameter for op %s", op.ID))
	}
	for _, p := range op.Parameters {
		if p.Name == paramName {
			return p, true
		}
	}
	panic(fmt.Errorf("x-paging.pageParameter specifies nonexistent parameter %s for op %s",
		paramName, op.ID))
}

// PagingResourcePath returns the path to the array to page over for a
// paging-enabled endpoint (specified by x-paging.resourcePath). Panics if
// paging is not enabled.
func PagingResourcePath(op *spec.Operation) []string {
	pagingConfig, ok := op.Extensions["x-paging"].(map[string]interface{})
	if !ok {
		panic(fmt.Errorf("no paging for operation %s", op.ID))
	}
	resourcePath, ok := pagingConfig["resourcePath"].(string)
	if !ok {
		return []string{}
	}
	return strings.Split(resourcePath, ".")
}

// resolveSchema dereferences the $ref in this schema, if there is one
func resolveSchema(s *spec.Swagger, schema *spec.Schema) *spec.Schema {
	if schema.Ref.String() == "" {
		return schema
	}
	refObj, _, err := schema.Ref.GetPointer().Get(s)
	if err != nil {
		panic(fmt.Errorf("bad schema reference %s", schema.Ref.String()))
	}
	r, ok := refObj.(spec.Schema)
	if !ok {
		panic(fmt.Errorf("bad schema reference %s", schema.Ref.String()))
	}
	return &r
}

// PagingResourceType returns the type of the items of the array set to page
// over for this operation and whether the array should consist of pointers to
// that type. Panics if paging is not enabled.
func PagingResourceType(s *spec.Swagger, op *spec.Operation) (string, bool, error) {
	var schema *spec.Schema
	for statusCode := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			schema = OutputSchema(s, op, statusCode)
			break
		}
	}
	if schema == nil {
		return "", false, fmt.Errorf("operation has no success type")
	}

	path := PagingResourcePath(op)
	for idx, pathComponent := range path {
		schema = resolveSchema(s, schema)
		nextSchema, ok := schema.Properties[pathComponent]
		if !ok {
			return "", false, fmt.Errorf("could not resolve x-paging.resourcePath: "+
				"%s has no field %s", strings.Join(path[:idx], "."), pathComponent)
		}
		schema = &nextSchema
	}
	schema = resolveSchema(s, schema)

	if len(schema.Type) != 1 || schema.Type[0] != "array" {
		return "", false, fmt.Errorf("paging resource type is not an array")
	}
	items := schema.Items
	if items == nil || items.Schema == nil {
		return "", false, fmt.Errorf("type of paging resource type has invalid `items`")
	}
	schema = items.Schema

	resourceType, err := TypeFromSchema(schema, true)
	if err != nil {
		return "", false, fmt.Errorf("could not convert paging resource to type: %s", err)
	}
	return resourceType, len(path) > 0, nil
}

// CodeToTypeMap returns a map from return status code to its corresponding type
func CodeToTypeMap(s *spec.Swagger, op *spec.Operation, makePointers bool) map[int]string {
	resp := make(map[int]string)
	for _, statusCode := range SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
		outputType, makeTypePtr := OutputType(s, op, statusCode)
		if makeTypePtr && makePointers {
			outputType = "*" + outputType
		}
		resp[statusCode] = outputType
	}
	return resp
}

// TypeToCodeMap returns a map from the type to its corresponding status code. It returns
// an error if mutiple status codes map to the same type
func TypeToCodeMap(s *spec.Swagger, op *spec.Operation) (map[string]int, error) {
	typeToCode := make(map[string]int)
	for code, typeStr := range CodeToTypeMap(s, op, false) {
		if typeStr != "" {
			if _, ok := typeToCode[typeStr]; ok {
				return nil, fmt.Errorf("duplicate response types %s, %s", typeStr, op.ID)
			}
			typeToCode[typeStr] = code
			typeToCode["*"+typeStr] = code
		} else {
			typeToCode[""] = code
		}
	}
	return typeToCode, nil
}

// SingleSchemaedBodyParameter returns true if the operation has a single,
// schema'd body input. It also returns the name of the model (without "models.").
func SingleSchemaedBodyParameter(op *spec.Operation) (bool, string) {
	if len(op.Parameters) == 1 && op.Parameters[0].ParamProps.Schema != nil {
		typ, err := TypeFromSchema(op.Parameters[0].ParamProps.Schema, false)
		if err != nil {
			panic(err)
		}
		return true, typ
	}
	return false, ""
}

// SingleStringPathParameter returns true if the operation has a single, required
// string input in the URL path. It also returns the name of the parameter.
func SingleStringPathParameter(op *spec.Operation) (bool, string) {
	if len(op.Parameters) != 1 {
		return false, ""
	}
	param := op.Parameters[0]
	if param.ParamProps.In == "path" && param.SimpleSchema.Type == "string" &&
		param.ParamProps.Required {
		return true, utils.CamelCase(param.Name, false)
	}
	return false, ""
}
