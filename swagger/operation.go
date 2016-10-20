package swagger

import (
	"fmt"
	"strings"

	"github.com/go-openapi/spec"
)

// Interface returns the interface for an operation
func Interface(op *spec.Operation) string {
	capOpID := Capitalize(op.ID)

	// Don't add the input parameter argument unless there are some arguments.
	// If a method has a single input parameter, and it's a schema, make the
	// generated type for that schema the input of the method.
	// TODO: If a method has a single input parameter, and it's a simple type (int/string),
	// make that the input of the method.
	// If a method has multiple input parameters, wrap them in a struct.
	input := ""
	if ssbp, opModel := SingleSchemaedBodyParameter(op); ssbp {
		input = fmt.Sprintf("i *models.%s", opModel)
	} else if len(op.Parameters) > 0 {
		input = fmt.Sprintf("i *models.%sInput", capOpID)
	}

	if NoSuccessType(op) {
		return fmt.Sprintf("%s(ctx context.Context, %s) error", capOpID, input)
	}
	successType, makePointer := OutputType(op, -1)
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
func OutputType(op *spec.Operation, statusCode int) (string, bool) {
	// If there is no success type and this is a success status code return the empty
	// string to indicate no type
	if NoSuccessType(op) && statusCode < 400 {
		return "", false
	}
	successCodes := successStatusCodes(op)
	if len(successCodes) == 1 && statusCode < 400 {
		singleSchema := op.Responses.StatusCodeResponses[successCodes[0]].Schema
		var err error
		successType, err := TypeFromSchema(singleSchema, true)
		if err != nil {
			panic(fmt.Errorf("Could not convert operation to type for %s, %s", op.ID, err))
		}
		return successType, singleSchema != nil && singleSchema.Ref.String() != ""
	}
	// This magic number is only used internally in this file. I will clean it up soon.
	if statusCode == -1 {
		return fmt.Sprintf("models.%sOutput", Capitalize(op.ID)), false
	}
	return fmt.Sprintf("models.%s%dOutput", Capitalize(op.ID), statusCode), true
}

// NoSuccessType returns true if the operation has no-success response type. This includes
// either no 200-399 response code or a 200-399 response code without a schema.
func NoSuccessType(op *spec.Operation) bool {
	successCodes := successStatusCodes(op)
	if len(successCodes) > 1 {
		return false
	}
	if len(successCodes) == 0 {
		return true
	}
	return op.Responses.StatusCodeResponses[successCodes[0]].Schema == nil
}

// CodeToTypeMap returns a map from return status code to its corresponding type
func CodeToTypeMap(op *spec.Operation) map[int]string {
	resp := make(map[int]string)
	for _, statusCode := range SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
		outputType, makePointer := OutputType(op, statusCode)
		if makePointer {
			outputType = "*" + outputType
		}
		resp[statusCode] = outputType
	}
	return resp
}

// successStatusCodes returns a slice of all the success status codes for an operation
func successStatusCodes(op *spec.Operation) []int {
	var successCodes []int
	for statusCode := range op.Responses.StatusCodeResponses {
		if statusCode < 400 {
			successCodes = append(successCodes, statusCode)
		}
	}
	return successCodes
}

// SingleSchemaedBodyParameter returns true if the operation has a single,
// schema'd body input. It also returns the name of the model (without "models.").
func SingleSchemaedBodyParameter(op *spec.Operation) (bool, string) {
	if len(op.Parameters) == 1 && op.Parameters[0].ParamProps.Schema != nil {
		typ, err := TypeFromSchema(op.Parameters[0].ParamProps.Schema, false)
		if err != nil {
			panic(err)
		}
		// Assuming that we're in a patch here (should we check that somehow???)
		if t, _ := wagPatchType(op); t != "" {
			typ = "Patch" + t
		}
		return true, typ
	}
	return false, ""
}

func wagPatchType(op *spec.Operation) (string, error) {
	for _, param := range op.Parameters {
		wagPatch, ok := param.Extensions.GetBool("x-wag-patch")
		if wagPatch && ok {
			if param.Schema == nil {
				// TODO: Should move these to validation before this code?
				return "", fmt.Errorf("TODO: nice error message")
			}
			ref := param.Schema.Ref.String()
			if !strings.HasPrefix(ref, "#/definitions/") {
				return "", fmt.Errorf("TODO: add a nice error message")
			}
			return ref[len("#/definitions/"):], nil
		}
	}
	return "", nil
}

// WagPatchDataTypes returns a set of all the data types that will have Patch{{Name}}
// structs created. This the input parameters marked with `x-wag-patch`
func WagPatchDataTypes(paths map[string]spec.PathItem) (map[string]struct{}, error) {

	// TODO: Better name from definitionsToExtend
	definitionsToExtend := make(map[string]struct{})

	for _, path := range paths {
		if path.Patch != nil {
			typeStr, err := wagPatchType(path.Patch)
			if err != nil {
				return nil, err
			}
			definitionsToExtend[typeStr] = struct{}{}
		}
	}
	return definitionsToExtend, nil
}
