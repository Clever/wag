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

// GetBooksInput holds the input parameters for a getBooks operation
type GetBooksInput struct {
	Authors     []string
	Available   *bool
	State       *string
	Published   *strfmt.Date
	SnakeCase   *string
	Completed   *strfmt.DateTime
	MaxPages    *float64
	MinPages    *int32
	PagesToTime *float32
}

// Validate returns an error if any of the GetBooksInput parameters don't satisfy the
// requirements from the swagger yml file
func (i GetBooksInput) Validate() error {
	if i.Authors != nil {
		if err := validate.MaxItems("authors", "query", int64(len(i.Authors)), 2); err != nil {
			return err
		}
	}
	if i.Authors != nil {
		if err := validate.MinItems("authors", "query", int64(len(i.Authors)), 1); err != nil {
			return err
		}
	}
	if i.Authors != nil {
		if err := validate.UniqueItems("authors", "query", i.Authors); err != nil {
			return err
		}
	}
	if i.State != nil {
		if err := validate.Enum("state", "query", *i.State, []interface{}{"finished", "inprogress"}); err != nil {
			return err
		}
	}
	if i.Published != nil {
		if err := validate.FormatOf("published", "query", "date", (*i.Published).String(), strfmt.Default); err != nil {
			return err
		}
	}
	if i.SnakeCase != nil {
		if err := validate.MaxLength("snake_case", "query", string(*i.SnakeCase), 5); err != nil {
			return err
		}
	}
	if i.Completed != nil {
		if err := validate.FormatOf("completed", "query", "date-time", (*i.Completed).String(), strfmt.Default); err != nil {
			return err
		}
	}
	if i.MaxPages != nil {
		if err := validate.Maximum("maxPages", "query", float64(*i.MaxPages), 1000.000000, false); err != nil {
			return err
		}
	}
	if i.MaxPages != nil {
		if err := validate.Minimum("maxPages", "query", float64(*i.MaxPages), -5.000000, false); err != nil {
			return err
		}
	}
	if i.MaxPages != nil {
		if err := validate.MultipleOf("maxPages", "query", float64(*i.MaxPages), 0.500000); err != nil {
			return err
		}
	}
	return nil
}

// GetBookByIDInput holds the input parameters for a getBookByID operation
type GetBookByIDInput struct {
	BookID        int64
	Authorization *string
	RandomBytes   *strfmt.Base64
}

// Validate returns an error if any of the GetBookByIDInput parameters don't satisfy the
// requirements from the swagger yml file
func (i GetBookByIDInput) Validate() error {
	if err := validate.MaximumInt("book_id", "path", i.BookID, int64(10000000), false); err != nil {
		return err
	}
	if err := validate.MinimumInt("book_id", "path", i.BookID, int64(2), false); err != nil {
		return err
	}
	if err := validate.MultipleOf("book_id", "path", float64(i.BookID), 2.000000); err != nil {
		return err
	}
	if i.Authorization != nil {
		if err := validate.MaxLength("authorization", "header", string(*i.Authorization), 24); err != nil {
			return err
		}
	}
	if i.Authorization != nil {
		if err := validate.MinLength("authorization", "header", string(*i.Authorization), 24); err != nil {
			return err
		}
	}
	if i.Authorization != nil {
		if err := validate.Pattern("authorization", "header", string(*i.Authorization), "[0-9a-f]+"); err != nil {
			return err
		}
	}
	if i.RandomBytes != nil {
		if err := validate.FormatOf("randomBytes", "query", "byte", string(*i.RandomBytes), strfmt.Default); err != nil {
			return err
		}
	}
	return nil
}

// CreateBookInput holds the input parameters for a createBook operation
type CreateBookInput struct {
	NewBook *Book
}

// Validate returns an error if any of the CreateBookInput parameters don't satisfy the
// requirements from the swagger yml file
func (i CreateBookInput) Validate() error {
	if err := i.NewBook.Validate(nil); err != nil {
		return err
	}

	return nil
}

// HealthCheckInput holds the input parameters for a healthCheck operation
type HealthCheckInput struct {
}

// Validate returns an error if any of the HealthCheckInput parameters don't satisfy the
// requirements from the swagger yml file
func (i HealthCheckInput) Validate() error {
	return nil
}
