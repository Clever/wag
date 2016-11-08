package swagger

import (
	"errors"
	"fmt"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
)

// ValidateResponses checks the responses from swagger operations, and in a few cases
// tweaks the swagger spec to make them valid. This means a few things:
// - Verifying that the spec only has one success type
// - Verifying the required global error types exist and that they have the required fields
// - Adding 400 / 500 responses to any operation that doesn't have them defined
func ValidateResponses(s spec.Swagger) error {

	// Confirm that we have the global error types we're expecting
	global400 := false
	global500 := false
	for name, resp := range s.Responses {

		if resp.Schema == nil {
			return fmt.Errorf("%s response must have schema", name)
		}
		if err := refHasMessageField(resp.Schema.Ref, s); err != nil {
			return fmt.Errorf("%s response is invalid: %s", name, err)
		}
		if name == "BadRequest" {
			global400 = true
		} else if name == "InternalError" {
			global500 = true
		}
	}
	if !global400 || !global500 {
		return errors.New("must specify global 'BadRequest' response type and global " +
			"'InternalError' response type")
	}

	for _, pathKey := range SortedPathItemKeys(s.Paths.Paths) {
		path := s.Paths.Paths[pathKey]
		pathItemOps := PathItemOperations(path)
		for _, opKey := range SortedOperationsKeys(pathItemOps) {
			op := pathItemOps[opKey]

			has400 := false
			has500 := false

			successResponseCount := 0
			for code, resp := range op.Responses.StatusCodeResponses {

				// Any 400+ responses must have `message` field so that
				// they can have a generated Error message
				if code < 300 {
					successResponseCount++
				} else if code < 400 {
					return fmt.Errorf("cannot define 3XX status codes: %s", op.ID)
				} else {
					if err := responseHasMessageField(resp, s); err != nil {
						return fmt.Errorf("invalid %d response: %s", code, err)
					}
				}

				if code == 400 {
					has400 = true
				} else if code == 500 {
					has500 = true
				}
			}

			if successResponseCount > 1 {
				return fmt.Errorf("can only define one success type (statusCode < 400) for %s", op.ID)
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

			// Confirm that the operation has a one-to-one map from status code -> type.
			_, err := TypeToCodeMap(&s, op)
			if err != nil {
				return err
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

// responseHasMessageField checks that a response points to a type with
// a message field. This should be used by responses defined in an operation
// (i.e. not global response type). Responses in an operation can either
// have a reference to a global response type or they have have a schema.
func responseHasMessageField(r spec.Response, s spec.Swagger) error {
	refToCheck := r.Ref
	if r.Schema != nil {
		refToCheck = r.Schema.Ref
	}

	return refHasMessageField(refToCheck, s)
}

// refHasMessageField ensures that the reference points to a schema with
// a `message` field and no other required fields.
func refHasMessageField(ref spec.Ref, s spec.Swagger) error {

	refObj, _, err := ref.GetPointer().Get(s)
	if err != nil {
		return fmt.Errorf("invalid schema reference: %s", err)
	}

	// The reference can point directly to a schema, or it can
	// point to a global response type which can then point to a
	// schema.
	r, ok := refObj.(spec.Response)
	if ok {
		return refHasMessageField(r.Schema.Ref, s)
	}
	schema, ok := refObj.(spec.Schema)
	if !ok {
		return errors.New("invalid schema reference")

	}

	messageField, ok := schema.Properties["message"]
	if !ok {
		return fmt.Errorf("schema must have a 'message' field: %s", ref.String())
	}

	if len(messageField.Type) != 1 || messageField.Type[0] != "string" {
		return fmt.Errorf("'message' field must be of type 'string': %s", ref.String())
	}

	// Don't allow any required fields. We need this because Wag won't know what those
	// fields should be when it generates the default 400 + 500 responses. We don't even
	// allow `message` to be required because go-swagger would make it a pointer which
	// complicates things.
	//
	// Note that we enforce this on all global response types, not just the 400 / 500s.
	// For now we do it because it makes the code simpler, but we could relax the
	// restriction if it limits users.
	if len(schema.Required) > 0 {
		return fmt.Errorf("%s cannot have required fields", ref.String())
	}

	return nil
}
