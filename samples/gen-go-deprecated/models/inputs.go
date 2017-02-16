package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// These imports may not be used depending on the input parameters
var _ = json.Marshal
var _ = fmt.Sprintf
var _ = url.QueryEscape
var _ = strconv.FormatInt
var _ = strings.Replace
var _ = validate.Maximum
var _ = strfmt.NewFormats

// HealthInput holds the input parameters for a health operation.
type HealthInput struct {
	Section int64
}

// Validate returns an error if any of the HealthInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i HealthInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i HealthInput) Path() (string, error) {
	path := "/v1/health"
	urlVals := url.Values{}

	urlVals.Add("section", strconv.FormatInt(i.Section, 10))

	return path + "?" + urlVals.Encode(), nil
}
