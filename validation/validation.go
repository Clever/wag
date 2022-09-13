package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Clever/wag/v9/swagger"
	swaggererrors "github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// A regex requiring the field to be start with a letter and be alphanumeric
var alphaNumRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")

// Validate checks if the swagger operation has any fields we don't support
func validateOp(s *spec.Swagger, path, method string, op *spec.Operation) error {
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
		// We could strip out the special characters, but it seems clearer to just enforce
		// this.
		return fmt.Errorf("The operationId for %s %s must be alphanumeric and start with a letter",
			method, path)
	}

	if err := validateResponses(path, method, op); err != nil {
		return err
	}

	if err := validateParams(path, method, op); err != nil {
		return err
	}

	return validatePaging(s, path, method, op)
}

func validateResponses(path, method string, op *spec.Operation) error {

	for _, statusCode := range swagger.SortedStatusCodeKeys(op.Responses.StatusCodeResponses) {
		if statusCode < 200 || statusCode > 599 {
			return fmt.Errorf("response map key must be an integer between 200 and 599 or "+
				"the string 'default'. Was %d for %s %s", statusCode, path, method)
		}
		response := op.Responses.StatusCodeResponses[statusCode]
		refStr := response.Ref.String()
		if refStr != "" && strings.HasPrefix(refStr, "#/definitions") {
			return fmt.Errorf("responses with references should nest the ref in a schema. "+
				"responses.%d.$ref = '%s' for %s %s should be "+
				"responses.%d.schema.$ref = '%s'",
				statusCode, refStr, path, method, statusCode, refStr)
		}

		_, err := swagger.TypeFromSchema(response.Schema, false)
		if err != nil {
			return fmt.Errorf("responses.%d for %s %s: %s", statusCode, method, path, err.Error())
		}
	}
	return nil
}

func validateParams(path, method string, op *spec.Operation) error {

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
			if strings.ToLower(param.Name) == "x-next-page-path" {
				return fmt.Errorf("%s %s uses reserved header parameter name X-Next-Page-Path",
					method, path)
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

func validatePaging(s *spec.Swagger, path, method string, op *spec.Operation) error {
	pagingConfig, ok := op.Extensions["x-paging"].(map[string]interface{})
	if !ok {
		return nil
	}

	upperMethod := strings.ToUpper(method)
	// Technically we could support paging for other methods as well, but in most cases it doesn't
	// make sense (for example, on DELETEs)
	if upperMethod != "GET" && upperMethod != "PUT" {
		return fmt.Errorf("%s %s cannot use x-paging. WAG only supports "+
			"paging on GET endpoints", method, path)
	}

	if singleString, _ := swagger.SingleStringPathParameter(op); singleString {
		return fmt.Errorf("%s %s cannot use x-paging. WAG doesn't support "+
			"paging on endpoints with a single string path parameter", method, path)
	}

	pagingParamName, ok := pagingConfig["pageParameter"].(string)
	if !ok {
		return fmt.Errorf("%s %s has invalid x-paging section. x-paging must include "+
			"a `pageParameter` field of type string set to the name of the parameter"+
			"that specifies the page ID for this operation", method, path)
	}

	var pagingParam *spec.Parameter
	for _, p := range op.Parameters {
		if p.Name == pagingParamName {
			pagingParam = &p
			break
		}
	}
	if pagingParam == nil {
		return fmt.Errorf("%s %s has invalid x-paging.pageParameter. Parameter '%s' does not exist",
			method, path, pagingParamName)
	}

	if resourcePathIntf, ok := pagingConfig["resourcePath"]; ok {
		_, isString := resourcePathIntf.(string)
		if !isString {
			return fmt.Errorf("%s %s has invalid x-paging.resourcePath. Field must be a string or undefined",
				method, path)
		}
	}

	if _, _, err := swagger.PagingResourceType(s, op); err != nil {
		return fmt.Errorf("%s %s has invalid paging resource (defaults to successful return type): %s",
			method, path, err.Error())
	}

	return nil
}

func validateDefinitions(definitions map[string]spec.Schema) error {
	for name, def := range definitions {
		for _, subDef := range def.Properties {
			if len(subDef.Type) == 1 && subDef.Type[0] == "object" {
				// We throw an error here because nested objects generate a compiler error in
				// the go-swagger model code.
				return fmt.Errorf("%s cannot have nested object types", name)
			}
		}
	}
	return nil
}

// Validate returns an error if the swagger file is invalid or uses fields
// we don't support. Note that this isn't a comprehensive check for all things
// we don't support, so this may not return an error, but the Swagger file might
// have values we don't support
func Validate(d loads.Document, generateJSClient bool) error {
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
	if !ok && generateJSClient {
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
			if err := validateOp(s, path, op, opItem); err != nil {
				return err
			}
		}
	}

	return validateDefinitions(s.Definitions)
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
