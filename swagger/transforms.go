package swagger

import (
	"errors"
	"fmt"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
)

// TransformErrors - TODO: Add a nice comment! This is a somewhat non-trivial function.
// TODO: Figure out the exact interface here... Know that it modifies the input...
// TODO: Add some unit tests...
func TransformErrors(s spec.Swagger) error {

	// Confirm that we have the global error types we're expecting
	global400 := false
	global500 := false
	for name, resp := range s.Responses {

		if resp.Schema == nil {
			return fmt.Errorf("%s response must have schema", name)
		}

		if err := validReference(resp.Schema.Ref, s); err != nil {
			return fmt.Errorf("%s response is invalid: %s", name, err)
		}
		if name == "BadRequest" {
			global400 = true
		} else if name == "InternalError" {
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

				// Do I need to do any more checking???

				if code == 400 || code == 500 {

					// Two cases to support here. One is a reference to a
					// response object. The other is a schema with a reference to
					// a type
					refToCheck := resp.Ref
					if resp.Schema != nil {
						refToCheck = resp.Schema.Ref
					}

					if err := validReference(refToCheck, s); err != nil {
						// TODO: clean up this message...
						return fmt.Errorf("invalid 400 response: %s", err)
					}
					if code == 400 {
						has400 = true
					} else {
						has500 = true
					}
				}
			}

			if !has400 {
				refResponse, err := createRefResponse("BadRequest", "#/responses/BadRequest")
				if err != nil {
					return err
				}
				op.Responses.StatusCodeResponses[400] = *refResponse
			}
			if !has500 {
				refResponse, err := createRefResponse("InternalError", "#/responses/InternalError")
				if err != nil {
					return err
				}
				op.Responses.StatusCodeResponses[500] = *refResponse
			}
		}
	}

	return nil
}

// createRefResponse returns a pointer to a spec.Response object
func createRefResponse(description, ref string) (*spec.Response, error) {
	jsonref, err := jsonreference.New(ref)
	if err != nil {
		return nil, err
	}

	return &spec.Response{
		Refable: spec.Refable{Ref: spec.Ref{Ref: jsonref}},
	}, nil
}

// TODO: Add a nice comment
// Make sure we have a Msg field and no other fields are required
func validReference(ref spec.Ref, s spec.Swagger) error {

	refObj, _, err := ref.GetPointer().Get(s)
	if err != nil {
		return fmt.Errorf("invalid schema reference: %s", err)
	}

	schema, ok := refObj.(spec.Schema)
	if !ok {

		// TODO: clean this up... maybe move it into the upper layer?
		r, ok := refObj.(spec.Response)
		if !ok {
			return errors.New("invalid schema reference")
		}
		return validReference(r.Schema.Ref, s)
	}

	msgField, ok := schema.Properties["msg"]
	if !ok {
		return fmt.Errorf("schema must have a 'msg' field: %s", ref.String())
	}

	if len(msgField.Type) != 1 || msgField.Type[0] != "string" {
		return fmt.Errorf("msg field must be of type 'string': %s", ref.String())
	}

	// TODO: Check for required fields...

	return nil
}
