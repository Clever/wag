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

// Path returns the URI path for the input.
func (i GetBookInput) Path() (string, error) {
	path := "/v1/books/{id}"
	urlVals := url.Values{}

	pathid := strconv.FormatInt(i.ID, 10)
	if pathid == "" {
		err := fmt.Errorf("id cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{id}", pathid, -1)

	return path + "?" + urlVals.Encode(), nil
}
