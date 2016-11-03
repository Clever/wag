package swagger

import (
	"errors"
	"fmt"

	"github.com/go-openapi/spec"
)

// TransformErrors - TODO: Add a nice comment! This is a somewhat non-trivial function.
// TODO: Figure out the exact interface here...
// Probably worth testing this directly...
func TransformErrors(s spec.Swagger) error {

	// Confirm that we have the global error types we're expecting
	global400 := false
	global500 := false
	for name, resp := range s.Responses {
		if name == "BadRequest" {
			if err := validErrorResponse(resp, s); err != nil {
				return fmt.Errorf("invalid bad request defined: %s", err)
			}
			global400 = true
		} else if name == "InternalError" {
			if err := validErrorResponse(resp, s); err != nil {
				return fmt.Errorf("invalid bad request defined: %s", err)
			}
			global500 = true
		}
	}
	if !global400 || !global500 {
		// TODO: add these to the template-wag
		// TODO: should I reference something to make it more clear what I mean...
		// probably the readme
		return errors.New("must specify global 'BadRequest' response type and global " +
			"'InternalResponse' response type")
	}

	for _, pathKey := range SortedPathItemKeys(s.Paths.Paths) {
		path := s.Paths.Paths[pathKey]
		pathItemOps := PathItemOperations(path)
		for _, opKey := range SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]

			has400 := false
			has500 := false

			for code, resp := range op.Responses.StatusCodeResponses {
				if code == 400 {
					if err := validErrorResponse(resp, s); err != nil {
						return fmt.Errorf("invalid 400 response: %s", err)
					}
					has400 = true
				} else if code == 500 {
					if err := validErrorResponse(resp, s); err != nil {
						return fmt.Errorf("invalid 500 response: %s", err)
					}
					has500 = true
				}
			}

			if !has400 {
				// TODO: build the default response
				op.Responses.StatusCodeResponses[400] = spec.Response{}
			}
			if !has500 {
				// TODO: build the default response
				op.Responses.StatusCodeResponses[500] = spec.Response{}
			}
		}
	}

	return nil
}

// TODO: Add a nice comment
// Make sure we have a Msg field and no other fields are required
func validErrorResponse(r spec.Response, s spec.Swagger) error {
	if r.Schema == nil {
		return errors.New("response must have schema")
	}

	// We support either object types or types that ref another object
	if r.Schema.Type[0] == "object" {
		// handle this

		// TODO: this isn't right. Need to learn how schema references work better...
		// } else if r.Schema.Ref.Ref.GetPointer().Get(s) {

		return nil

	}
	return errors.New("response schema must be an object or reference an object")
}
