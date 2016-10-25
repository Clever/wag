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

// Wagpatch2Input holds the input parameters for a wagpatch2 operation.
type Wagpatch2Input struct {
	Other        *string
	SpecialPatch *PatchData
}

// Validate returns an error if any of the Wagpatch2Input parameters don't satisfy the
// requirements from the swagger yml file.
func (i Wagpatch2Input) Validate() error {
	if err := i.SpecialPatch.Validate(nil); err != nil {
		return err
	}

	return nil
}
