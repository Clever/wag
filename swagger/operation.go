package swagger

import (
	"fmt"
	"strings"

	"github.com/Clever/wag/utils"
	"github.com/go-openapi/spec"
)

// Interface returns the interface for an operation
func Interface(s *spec.Swagger, op *spec.Operation) string {
	capOpID := Capitalize(op.ID)

	// Don't add the input parameter argument unless there are some arguments.
	// If a method has a single input parameter, and it's a schema, make the
	// generated type for that schema the input of the method.
	// If a method has a single input parameter, and it's a simple type (string, TODO: int),
	// make that the input of the method.
	// If a method has multiple input parameters, wrap them in a struct.
	input := ""
	if singleSchemaedBodyParameter, opModel := SingleSchemaedBodyParameter(op); singleSchemaedBodyParameter {
		input = fmt.Sprintf("i *models.%s", opModel)
	} else if singleStringPathParameter, inputName := SingleStringPathParameter(op); singleStringPathParameter {
		input = fmt.Sprintf("%s string", inputName)
	} else if len(op.Parameters) > 0 {
		input = fmt.Sprintf("i *models.%sInput", capOpID)
	}

	if NoSuccessType(op) {
		return fmt.Sprintf("%s(ctx context.Context, %s) error", capOpID, input)
	}
	successCode := successStatusCode(op)
	successType, makePointer := OutputType(s, op, *successCode)
	if makePointer {
		successType = "*" + successType
	}

	return fmt.Sprintf("%s(ctx context.Context, %s) (%s, error)",
		capOpID, input, successType)
}

// InterfaceComment returns the comment for the interface for the operation
func InterfaceComment(method, path string, op *spec.Operation) string {

	capOpID := Capitalize(op.ID)
	comment := fmt.Sprintf("// %s makes a %s request to %s.", capOpID, method, path)
	if op.Description != "" {
		comment += "\n// " + op.Description
	}
	return comment
}

// OutputType returns the output type for a given status code of an operation and whether it
// is a pointer in the interface.
func OutputType(s *spec.Swagger, op *spec.Operation, statusCode int) (string, bool) {
	// If there is no success type and this is a success status code return the empty
	// string to indicate no type
	if NoSuccessType(op) && statusCode < 400 {
		return "", false
	}

	// This magic number is only used internally in this file. I will clean it up at some point.
	// It is used if there are multiple success types
	successCodes := successStatusCodes(op)
	if len(successCodes) > 1 {
		if statusCode == -1 {
			return fmt.Sprintf("models.%sOutput", Capitalize(op.ID)), false
		} else if statusCode < 400 {
			return fmt.Sprintf("models.%s%dOutput", Capitalize(op.ID), statusCode), true
		}
	} else if statusCode == -1 {
		statusCode = successCodes[0]
	}

	resp := op.Responses.StatusCodeResponses[statusCode]
	schema := resp.Schema
	if strings.HasPrefix(resp.Ref.String(), "#/responses") {
		// Follow the pointer to the response and then to the schema
		refObj, _, err := resp.Ref.GetPointer().Get(s)
		if err != nil {
			panic("bad response schema reference")
		}
		r, ok := refObj.(spec.Response)
		if !ok {
			panic("bad response schema reference")
		}
		schema = r.Schema
	}

	successType, err := TypeFromSchema(schema, true)
	if err != nil {
		panic(fmt.Errorf("could not convert operation to type for %s, %s", op.ID, err))
	}
	return successType, schema != nil && schema.Ref.String() != ""
}

// NoSuccessType returns true if the operation has no-success response type. This includes
// either no 200-399 response code or a 200-399 response code without a schema.
func NoSuccessType(op *spec.Operation) bool {
	successCode := successStatusCodes(op)
	if successCode == nil {
		return true
	}
	return op.Responses.StatusCodeResponses[*successCode].Schema == nil
}

// CodeToTypeMap returns a map from return status code to its corresponding type
func CodeToTypeMap(s *spec.Swagger, op *spec.Operation) map[int]string {
	resp := make(map[int]string)
	for _, statusCode := range SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
		outputType, _ := OutputType(s, op, statusCode)
		resp[statusCode] = outputType
	}
	return resp
}

// TypeToCodeMap returns a map from the type to its corresponding status code. It returns
// an error if mutiple status codes map to the same type
func TypeToCodeMap(s *spec.Swagger, op *spec.Operation) (map[string]int, error) {
	typeToCode := make(map[string]int)
	for code, typeStr := range CodeToTypeMap(s, op) {
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

// successStatusCode returns the success status code. If there is no success status code
// then it returns nil.
func successStatusCode(op *spec.Operation) *int {
	for statusCode := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			return &statusCode
		}
	}
	return nil
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
