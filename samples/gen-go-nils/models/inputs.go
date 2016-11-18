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

// NilCheckInput holds the input parameters for a nilCheck operation.
type NilCheckInput struct {
	ID     string
	Query  *string
	Header string
	Body   *NilFields
}

// Validate returns an error if any of the NilCheckInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i NilCheckInput) Validate() error {
	if err := i.Body.Validate(nil); err != nil {
		return err
	}

	return nil
}
