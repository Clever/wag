package main

import (
	"fmt"
	"regexp"

	"github.com/go-openapi/spec"
)

// A regex requiring the field to be start with a letter and be alphanumeric
var alphaNumRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")

// Validate checks if the swagger operation has any fields we don't support
func validateOp(s *spec.Operation) error {
	if len(s.Consumes) != 0 {
		fmt.Errorf("Consumes not supported in swagger operations")
	}
	if len(s.Produces) != 0 {
		fmt.Errorf("Produces not supported in swagger operations")
	}
	if len(s.Schemes) != 0 {
		fmt.Errorf("Schemes not supported in swagger operations")
	}
	if len(s.Security) != 0 {
		fmt.Errorf("Security not supported in swagger operations")
	}

	// TODO: A test for this?
	/*
		for path, pathObj := range paths {
		for method, op := range pathObj {
	*/
	if !alphaNumRegex.MatchString(s.ID) {
		// We need this because we build function / struct names with the operationID.
		// We could strip out the special characters, but it seems clearly to just enforce
		// this.
		return fmt.Errorf("OperationIDs must be alphanumeric and start with a letter")
	}

	return nil
}

// validates returns an error if the swagger file is invalid or uses fields
// we don't support. Note that this isn't a comprehensive check for all things
// we don't support, so this may not return an error, but the Swagger file might
// have values we don't support
func validate(s spec.Swagger) error {
	if s.Swagger != "2.0" {
		return fmt.Errorf("Unsupported Swagger version %s", s.Swagger)
	}

	if len(s.Schemes) != 1 || s.Schemes[0] != "http" {
		return fmt.Errorf("Schemes only supports 'http', not %s")
	}

	if len(s.Consumes) != 1 || s.Consumes[0] != "application/json" {
		return fmt.Errorf("Consumes only supports 'application/json'")
	}

	if len(s.Produces) != 1 || s.Produces[0] != "application/json" {
		return fmt.Errorf("Produces only supports 'application/json'")
	}

	if s.Host != "" {
		return fmt.Errorf("Host parameter is not supported")
	}

	if len(s.Parameters) != 0 {
		return fmt.Errorf("Global parameters definitions are not supported. Define parameters on a per request basis.")
	}

	if len(s.Responses) != 0 {
		return fmt.Errorf("Global response definitions are not supported. Define responses on a per request basis")
	}

	if len(s.SecurityDefinitions) != 0 {
		return fmt.Errorf("Security definitions definition not supported")
	}

	if len(s.Security) != 0 {
		return fmt.Errorf("Security field not supported")
	}

	// Validate paths
	for fieldName, pathItem := range s.Paths.Paths {
		// Note that $ref and parameters are not valid as of now
		if !sliceContains([]string{"get", "put", "post", "delete", "options", "head", "patch"}, fieldName) {
			fmt.Errorf("Invalid path field name: %s", fieldName)
		}

		for _, op := range pathItemOperations(pathItem) {
			if err := validateOp(op); err != nil {
				return err
			}
		}
	}

	return nil
}
