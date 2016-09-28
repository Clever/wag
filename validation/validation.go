package validation

import (
	"fmt"
	"regexp"

	"github.com/go-openapi/spec"
)

// A regex requiring the field to be start with a letter and be alphanumeric
var alphaNumRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")

// Validate checks if the swagger operation has any fields we don't support
func validateOp(path, op string, s *spec.Operation) error {
	if len(s.Consumes) != 0 {
		return fmt.Errorf("%s %s cannot have a consumes field. WAG does not support the consumes field "+
			"on operations", op, path)
	}
	if len(s.Produces) != 0 {
		return fmt.Errorf("%s %s cannot have a produces field. WAG does not support the produces field "+
			"on operations", op, path)
	}
	if len(s.Schemes) != 0 {
		return fmt.Errorf("%s %s cannot have a schemes field. WAG does not support the schemes field "+
			"on operations", op, path)
	}
	if len(s.Security) != 0 {
		return fmt.Errorf("%s %s cannot have a security field. WAG does not support the security field "+
			"on operations", op, path)
	}

	if s.ID == "" {
		return fmt.Errorf("%s %s must have an operationId field, "+
			"see http://swagger.io/specification/#operationObject", op, path)
	}

	if !alphaNumRegex.MatchString(s.ID) {
		// We need this because we build function / struct names with the operationID.
		// We could strip out the special characters, but it seems clearly to just enforce
		// this.
		return fmt.Errorf("The operationId for %s %s must be alphanumeric and start with a letter",
			op, path)
	}

	for _, param := range s.Parameters {
		if param.In == "path" && !param.Required {
			return fmt.Errorf("%s for %s %s is a path parameter so it must be required",
				param.Name, op, path)
		}

		if param.Type == "string" && param.Format != "" {
			if param.MaxLength != nil || param.MinLength != nil || param.Pattern != "" || len(param.Enum) > 0 {
				return fmt.Errorf("%s for %s %s cannot have min/max length, pattern, or enum fields. "+
					"Only string type parameters without a format can have additional validation.",
					param.Name, op, path)
			}
		}
	}

	return nil
}

// Validate returns an error if the swagger file is invalid or uses fields
// we don't support. Note that this isn't a comprehensive check for all things
// we don't support, so this may not return an error, but the Swagger file might
// have values we don't support
func Validate(s spec.Swagger) error {
	if s.Swagger != "2.0" {
		return fmt.Errorf("Unsupported Swagger version %s", s.Swagger)
	}

	if len(s.Schemes) != 1 || s.Schemes[0] != "http" {
		return fmt.Errorf("WAG only supports the scheme 'http'")
	}

	if len(s.Consumes) > 1 || (len(s.Produces) == 0 && s.Consumes[0] != "application/json") {
		return fmt.Errorf("WAG only supports the consumes option: 'application/json'")
	}

	if len(s.Produces) > 1 || (len(s.Produces) == 0 && s.Produces[0] != "application/json") {
		return fmt.Errorf("WAG only support the consumes option: 'application/json'")
	}

	if s.Host != "" {
		return fmt.Errorf("WAG does not support the host field")
	}

	if len(s.Parameters) != 0 {
		return fmt.Errorf("WAG does not support global parameters definitions. Define parameters on a per request basis")
	}

	if len(s.Responses) != 0 {
		return fmt.Errorf("WAG does not support global response definitions.  Define responses on a per request basis")
	}

	if len(s.SecurityDefinitions) != 0 {
		return fmt.Errorf("WAG does not support the security definitions field")
	}

	if len(s.Security) != 0 {
		return fmt.Errorf("WAG does not support the security field")
	}

	for path, pathItem := range s.Paths.Paths {
		if pathItem.Ref.String() != "" {
			return fmt.Errorf("WAG does not support paths with $ref fields. Define the references on " +
				"a per operation basis")
		}
		if len(pathItem.Parameters) != 0 {
			return fmt.Errorf("Parameters cannot be defined for an entire path. " +
				"They must be defined on the individual method level.")
		}
		for op, opItem := range pathItemOperations(pathItem) {
			if err := validateOp(path, op, opItem); err != nil {
				return err
			}
		}
	}

	return nil
}

func sliceContains(slice []string, key string) bool {
	for _, val := range slice {
		if val == key {
			return true
		}
	}
	return false
}

func pathItemOperations(item spec.PathItem) map[string]*spec.Operation {
	ops := make(map[string]*spec.Operation)
	if item.Get != nil {
		ops["GET"] = item.Get
	}
	if item.Put != nil {
		ops["PUT"] = item.Put
	}
	if item.Post != nil {
		ops["POST"] = item.Post
	}
	if item.Delete != nil {
		ops["DELETE"] = item.Delete
	}
	if item.Options != nil {
		ops["OPTIONS"] = item.Options
	}
	if item.Head != nil {
		ops["HEAD"] = item.Head
	}
	if item.Patch != nil {
		ops["PATCH"] = item.Patch
	}
	return ops
}
