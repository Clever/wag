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

// GetDistrictsInput holds the input parameters for a getDistricts operation.
type GetDistrictsInput struct {
	Where         *WhereQueryString
	StartingAfter *string
	EndingBefore  *string
	PageSize      *int64
}

// Validate returns an error if any of the GetDistrictsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetDistrictsInput) Validate() error {

	if i.Where != nil {
		if err := i.Where.Validate(nil); err != nil {
			return err
		}
	}

	if i.StartingAfter != nil {
		if err := validate.FormatOf("starting_after", "query", "mongo-id", *i.StartingAfter, strfmt.Default); err != nil {
			return err
		}
	}

	if i.EndingBefore != nil {
		if err := validate.FormatOf("ending_before", "query", "mongo-id", *i.EndingBefore, strfmt.Default); err != nil {
			return err
		}
	}

	if i.PageSize != nil {
		if err := validate.MaximumInt("page_size", "query", *i.PageSize, int64(10000), false); err != nil {
			return err
		}
	}
	if i.PageSize != nil {
		if err := validate.MinimumInt("page_size", "query", *i.PageSize, int64(1), false); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the URI path for the input.
func (i GetDistrictsInput) Path() (string, error) {
	path := "/v1/check"
	urlVals := url.Values{}

	if i.StartingAfter != nil {
		urlVals.Add("starting_after", *i.StartingAfter)
	}

	if i.EndingBefore != nil {
		urlVals.Add("ending_before", *i.EndingBefore)
	}

	if i.PageSize != nil {
		urlVals.Add("page_size", strconv.FormatInt(*i.PageSize, 10))
	}

	return path + "?" + urlVals.Encode(), nil
}
