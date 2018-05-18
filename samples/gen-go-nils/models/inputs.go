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

// NilCheckInput holds the input parameters for a nilCheck operation.
type NilCheckInput struct {
	ID     string
	Query  *string
	Header string
	Array  []string
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

// Path returns the URI path for the input.
func (i NilCheckInput) Path() (string, error) {
	path := "/v1/check/{id}"
	urlVals := url.Values{}

	pathid := i.ID
	if pathid == "" {
		err := fmt.Errorf("id cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{id}", pathid, -1)

	if i.Query != nil {
		urlVals.Add("query", *i.Query)
	}

	for _, v := range i.Array {
		urlVals.Add("array", v)
	}

	return path + "?" + urlVals.Encode(), nil
}
