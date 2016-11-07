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

// GetBookInput holds the input parameters for a getBook operation.
type GetBookInput struct {
	ID int64
}

// Validate returns an error if any of the GetBookInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetBookInput) Validate() error {
	if err := validate.MaximumInt("id", "path", i.ID, int64(4000), false); err != nil {
		return err
	}
	return nil
}
