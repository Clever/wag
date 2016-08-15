package models

import (
	"encoding/json"
	"github.com/go-openapi/validate"
	"strconv"
)

// These imports may not be used depending on the input parameters
var _ = json.Marshal
var _ = strconv.FormatInt
var _ = validate.Maximum

type GetBooksInput struct {
	Author    *string
	Available *bool
	MaxPages  *float64
}

func (i GetBooksInput) Validate() error {

	if err := validate.Maximum("maxPages", "query", *i.MaxPages, 1.000000, false); err != nil {
		return err
	}

	if err := validate.Minimum("maxPages", "query", *i.MaxPages, -5.000000, false); err != nil {
		return err
	}

	if err := validate.MultipleOf("maxPages", "query", *i.MaxPages, 0.500000); err != nil {
		return err
	}
	return nil
}

type GetBookByIDInput struct {
	BookID        int64
	Authorization *string
}

func (i GetBookByIDInput) Validate() error {

	if err := validate.MaximumInt("bookID", "path", i.BookID, 10000000, false); err != nil {
		return err
	}

	if err := validate.MinimumInt("bookID", "path", i.BookID, 2, false); err != nil {
		return err
	}

	if err := validate.MultipleOf("bookID", "path", float64(i.BookID), 2.000000); err != nil {
		return err
	}

	if err := validate.MaxLength("authorization", "header", *i.Authorization, 24); err != nil {
		return err
	}

	if err := validate.MinLength("authorization", "header", *i.Authorization, 24); err != nil {
		return err
	}

	if err := validate.Pattern("authorization", "header", *i.Authorization, "[0-9a-f]+"); err != nil {
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
