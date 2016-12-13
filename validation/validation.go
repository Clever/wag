package validation

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Clever/wag/swagger"
	swaggererrors "github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// A regex requiring the field to be start with a letter and be alphanumeric
var alphaNumRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")

// Validate checks if the swagger operation has any fields we don't support
func validateOp(path, method string, op *spec.Operation) error {
	if len(op.Consumes) != 0 {
		return fmt.Errorf("%s %s cannot have a consumes field. WAG does not support the consumes field "+
			"on operations", method, path)
	}
	if len(op.Produces) != 0 {
		return fmt.Errorf("%s %s cannot have a produces field. WAG does not support the produces field "+
			"on operations", method, path)
	}
	if len(op.Schemes) != 0 {
		return fmt.Errorf("%s %s cannot have a schemes field. WAG does not support the schemes field "+
			"on operations", method, path)
	}
	if len(op.Security) != 0 {
		return fmt.Errorf("%s %s cannot have a security field. WAG does not support the security field "+
			"on operations", method, path)
	}

	if op.ID == "" {
		return fmt.Errorf("%s %s must have an operationId field, "+
			"see http://swagger.io/specification/#operationObject", method, path)
	}

	if !alphaNumRegex.MatchString(op.ID) {
		// We need this because we build function / struct names with the operationID.
		// We could strip out the special characters, but it seems clearly to just enforce
		// this.
		return fmt.Errorf("The operationId for %s %s must be alphanumeric and start with a letter",
			method, path)
	}

	for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
		if statusCode < 200 || statusCode > 599 {
			return fmt.Errorf("Response map key must be an integer between 200 and 599 or "+
				"the string 'default'. Was %d", statusCode)
		}
		_, err := swagger.TypeFromSchema(op.Responses.StatusCodeResponses[statusCode].Schema, false)
		if err != nil {
			return err
		}
	}

	for _, param := range op.Parameters {

		switch param.In {
		case "path":
			if !param.Required {
				return fmt.Errorf("%s for %s %s is a path parameter so it must be required",
					param.Name, method, path)
			}
			if param.Type == "array" || param.Type == "object" {
				return fmt.Errorf("%s for %s %s is a path parameter so it must have a primitive type",
					param.Name, method, path)
			}
		case "body":
			if param.Default != nil {
				return fmt.Errorf("%s for %s %s is a body parameter so it must not have a default",
					param.Name, method, path)
			}
			if param.Schema == nil || param.Schema.Ref.String() == "" {
				return fmt.Errorf("%s for %s %s is a body parameter so it must reference a schema",
					param.Name, method, path)
			}
		case "query":
			if param.Type == "object" {
				return fmt.Errorf("%s for %s %s is a query param so it can't have the type 'object'",
					param.Name, method, path)
			}
			if param.Type == "array" && param.Items.Type != "string" {
				return fmt.Errorf("array parameters must have string sub-types")
			}
		case "header":
			if param.Type == "array" || param.Type == "object" {
				return fmt.Errorf("%s for %s %s is a path parameter so it must have a primitive type",
					param.Name, method, path)
			}
		default:
			return fmt.Errorf("unsupported param type: %s", param.In)
		}

		if param.Type == "string" && param.Format != "" {
			if param.MaxLength != nil || param.MinLength != nil || param.Pattern != "" || len(param.Enum) > 0 {
				return fmt.Errorf("%s for %s %s cannot have min/max length, pattern, or enum fields. "+
					"Only string type parameters without a format can have additional validation.",
					param.Name, method, path)
			}
		}
	}

	return nil
}

// Validate returns an error if the swagger file is invalid or uses fields
// we don't support. Note that this isn't a comprehensive check for all things
// we don't support, so this may not return an error, but the Swagger file might
// have values we don't support
func Validate(d loads.Document) error {
	s := d.Spec()

	goSwaggerError := validate.Spec(&d, strfmt.Default)
	if goSwaggerError != nil {
		str := ""
		for _, desc := range goSwaggerError.(*swaggererrors.CompositeError).Errors {
			str += fmt.Sprintf("- %s\n", desc)
		}
		return errors.New(str)
	}

	if s.Swagger != "2.0" {
		return fmt.Errorf("unsupported Swagger version %s", s.Swagger)
	}

	if len(s.Schemes) != 1 || s.Schemes[0] != "http" {
		return fmt.Errorf("wag only supports the scheme 'http'")
	}

	if len(s.Consumes) > 1 || (len(s.Produces) == 0 && s.Consumes[0] != "application/json") {
		return fmt.Errorf("wag only supports the consumes option: 'application/json'")
	}

	if len(s.Produces) > 1 || (len(s.Produces) == 0 && s.Produces[0] != "application/json") {
		return fmt.Errorf("wag only support the consumes option: 'application/json'")
	}

	if s.Host != "" {
		return fmt.Errorf("wag does not support the host field")
	}

	if len(s.Parameters) != 0 {
		return fmt.Errorf("wag does not support global parameters definitions. Define parameters on a per request basis")
	}

	if len(s.SecurityDefinitions) != 0 {
		return fmt.Errorf("wag does not support the security definitions field")
	}

	if len(s.Security) != 0 {
		return fmt.Errorf("wag does not support the security field")
	}

	_, ok := s.Info.Extensions.GetString("x-npm-package")
	if !ok {
		return fmt.Errorf("must provide 'x-npm-package' in the 'info' section of the swagger.yml")
	}

	for path, pathItem := range s.Paths.Paths {
		if pathItem.Ref.String() != "" {
			return fmt.Errorf("wag does not support paths with $ref fields. Define the references on " +
				"a per operation basis")
		}
		if len(pathItem.Parameters) != 0 {
			return fmt.Errorf("parameters cannot be defined for an entire path. " +
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
