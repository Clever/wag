package models

import (
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"strconv"
)

// These imports may not be used depending on the input parameters
var _ = json.Marshal
var _ = strconv.FormatInt
var _ = validate.Maximum
var _ = strfmt.NewFormats

type GetBooksInput struct {
	Author      *string
	Available   *bool
	State       *string
	Published   *strfmt.Date
	Completed   *strfmt.DateTime
	MaxPages    *float64
	MinPages    *int32
	PagesToTime *float32
}

func (i GetBooksInput) Validate() error {

	if err := validate.Enum("state", "query", *i.State, []interface{}{"finished", "inprogress"}); err != nil {
		return err
	}

	if err := validate.FormatOf("published", "query", "date", (*i.Published).String(), strfmt.Default); err != nil {
		return err
	}

	if err := validate.FormatOf("completed", "query", "date-time", (*i.Completed).String(), strfmt.Default); err != nil {
		return err
	}

	if err := validate.Maximum("maxPages", "query", float64(*i.MaxPages), 1.000000, false); err != nil {
		return err
	}

	if err := validate.Minimum("maxPages", "query", float64(*i.MaxPages), -5.000000, false); err != nil {
		return err
	}

	if err := validate.MultipleOf("maxPages", "query", float64(*i.MaxPages), 0.500000); err != nil {
		return err
	}
	return nil
}

type GetBookByIDInput struct {
	BookID        int64
	Authorization *string
	RandomBytes   *[]byte
}

func (i GetBookByIDInput) Validate() error {

	if err := validate.MaximumInt("bookID", "path", i.BookID, int64(10000000), false); err != nil {
		return err
	}

	if err := validate.MinimumInt("bookID", "path", i.BookID, int64(2), false); err != nil {
		return err
	}

	if err := validate.MultipleOf("bookID", "path", float64(i.BookID), 2.000000); err != nil {
		return err
	}

	if err := validate.MaxLength("authorization", "header", string(*i.Authorization), 24); err != nil {
		return err
	}

	if err := validate.MinLength("authorization", "header", string(*i.Authorization), 24); err != nil {
		return err
	}

	if err := validate.Pattern("authorization", "header", string(*i.Authorization), "[0-9a-f]+"); err != nil {
		return err
	}

	if err := validate.FormatOf("randomBytes", "query", "byte", string(*i.RandomBytes), strfmt.Default); err != nil {
		return err
	}
	return nil
}

type CreateBookInput struct {
	NewBook *Book
}

func (i CreateBookInput) Validate() error {
	if err := i.NewBook.Validate(nil); err != nil {
		return err
	}

	return nil
}
