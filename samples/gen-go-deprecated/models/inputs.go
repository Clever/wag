package models

import (
	"encoding/json"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// These imports may not be used depending on the input parameters
var _ = json.Marshal
var _ = strconv.FormatInt
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
