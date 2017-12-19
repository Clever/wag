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

// GetAuthorsInput holds the input parameters for a getAuthors operation.
type GetAuthorsInput struct {
	Name          *string
	StartingAfter *string
}

// Validate returns an error if any of the GetAuthorsInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetAuthorsInput) Validate() error {

	return nil
}

// Path returns the URI path for the input.
func (i GetAuthorsInput) Path() (string, error) {
	path := "/v1/authors"
	urlVals := url.Values{}

	if i.Name != nil {
		urlVals.Add("name", *i.Name)
	}

	if i.StartingAfter != nil {
		urlVals.Add("startingAfter", *i.StartingAfter)
	}

	return path + "?" + urlVals.Encode(), nil
}

// GetAuthorsWithPutInput holds the input parameters for a getAuthorsWithPut operation.
type GetAuthorsWithPutInput struct {
	Name          *string
	StartingAfter *string
	FavoriteBooks *Book
}

// Validate returns an error if any of the GetAuthorsWithPutInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i GetAuthorsWithPutInput) Validate() error {

	if err := i.FavoriteBooks.Validate(nil); err != nil {
		return err
	}
	return nil
}

// Path returns the URI path for the input.
func (i GetAuthorsWithPutInput) Path() (string, error) {
	path := "/v1/authors"
	urlVals := url.Values{}

	if i.Name != nil {
		urlVals.Add("name", *i.Name)
	}

	if i.StartingAfter != nil {
		urlVals.Add("startingAfter", *i.StartingAfter)
	}

	return path + "?" + urlVals.Encode(), nil
}

// GetBooksInput holds the input parameters for a getBooks operation.
type GetBooksInput struct {
	Authors       []string
	Available     *bool
	State         *string
	Published     *strfmt.Date
	SnakeCase     *string
	Completed     *strfmt.DateTime
	MaxPages      *float64
	MinPages      *int32
	PagesToTime   *float32
	Authorization string
	StartingAfter *int64
}

// Validate returns an error if any of the GetBooksInput parameters don't satisfy the
// requirements from the swagger yml file.
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

// Path returns the URI path for the input.
func (i GetBooksInput) Path() (string, error) {
	path := "/v1/books"
	urlVals := url.Values{}

	for _, v := range i.Authors {
		urlVals.Add("authors", v)
	}

	if i.Available != nil {
		urlVals.Add("available", strconv.FormatBool(*i.Available))
	}

	if i.State != nil {
		urlVals.Add("state", *i.State)
	}

	if i.Published != nil {
		urlVals.Add("published", (*i.Published).String())
	}

	if i.SnakeCase != nil {
		urlVals.Add("snake_case", *i.SnakeCase)
	}

	if i.Completed != nil {
		urlVals.Add("completed", (*i.Completed).String())
	}

	if i.MaxPages != nil {
		urlVals.Add("maxPages", strconv.FormatFloat(*i.MaxPages, 'E', -1, 64))
	}

	if i.MinPages != nil {
		urlVals.Add("min_pages", strconv.FormatInt(int64(*i.MinPages), 10))
	}

	if i.PagesToTime != nil {
		urlVals.Add("pagesToTime", strconv.FormatFloat(float64(*i.PagesToTime), 'E', -1, 32))
	}

	if i.StartingAfter != nil {
		urlVals.Add("startingAfter", strconv.FormatInt(*i.StartingAfter, 10))
	}

	return path + "?" + urlVals.Encode(), nil
}

// GetBookByIDInput holds the input parameters for a getBookByID operation.
type GetBookByIDInput struct {
	BookID              int64
	AuthorID            *string
	Authorization       string
	XDontRateLimitMeBro string
	RandomBytes         *strfmt.Base64
}

// Validate returns an error if any of the GetBookByIDInput parameters don't satisfy the
// requirements from the swagger yml file.
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

	if i.AuthorID != nil {
		if err := validate.FormatOf("authorID", "query", "mongo-id", *i.AuthorID, strfmt.Default); err != nil {
			return err
		}
	}

	if len(i.Authorization) > 0 {
		if err := validate.MaxLength("authorization", "header", string(i.Authorization), 24); err != nil {
			return err
		}
	}
	if len(i.Authorization) > 0 {
		if err := validate.MinLength("authorization", "header", string(i.Authorization), 24); err != nil {
			return err
		}
	}
	if len(i.Authorization) > 0 {
		if err := validate.Pattern("authorization", "header", string(i.Authorization), "[0-9a-f]+"); err != nil {
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

// Path returns the URI path for the input.
func (i GetBookByIDInput) Path() (string, error) {
	path := "/v1/books/{book_id}"
	urlVals := url.Values{}

	pathbook_id := strconv.FormatInt(i.BookID, 10)
	if pathbook_id == "" {
		err := fmt.Errorf("book_id cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{book_id}", pathbook_id, -1)

	if i.AuthorID != nil {
		urlVals.Add("authorID", *i.AuthorID)
	}

	if i.RandomBytes != nil {
		urlVals.Add("randomBytes", string(*i.RandomBytes))
	}

	return path + "?" + urlVals.Encode(), nil
}

// GetBookByID2Input holds the input parameters for a getBookByID2 operation.
type GetBookByID2Input struct {
	ID string
}

// ValidateGetBookByID2Input returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetBookByID2Input(id string) error {

	if err := validate.Pattern("id", "path", string(id), "^[0-9a-f]{24}$"); err != nil {
		return err
	}

	return nil
}

// GetBookByID2InputPath returns the URI path for the input.
func GetBookByID2InputPath(id string) (string, error) {
	path := "/v1/books2/{id}"
	urlVals := url.Values{}

	pathid := id
	if pathid == "" {
		err := fmt.Errorf("id cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{id}", pathid, -1)

	return path + "?" + urlVals.Encode(), nil
}

// GetBookByIDCachedInput holds the input parameters for a getBookByIDCached operation.
type GetBookByIDCachedInput struct {
	ID string
}

// ValidateGetBookByIDCachedInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetBookByIDCachedInput(id string) error {

	if err := validate.Pattern("id", "path", string(id), "^[0-9a-f]{24}$"); err != nil {
		return err
	}

	return nil
}

// GetBookByIDCachedInputPath returns the URI path for the input.
func GetBookByIDCachedInputPath(id string) (string, error) {
	path := "/v1/bookscached/{id}"
	urlVals := url.Values{}

	pathid := id
	if pathid == "" {
		err := fmt.Errorf("id cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{id}", pathid, -1)

	return path + "?" + urlVals.Encode(), nil
}

// HealthCheckInput holds the input parameters for a healthCheck operation.
type HealthCheckInput struct {
}

// Validate returns an error if any of the HealthCheckInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i HealthCheckInput) Validate() error {
	return nil
}

// Path returns the URI path for the input.
func (i HealthCheckInput) Path() (string, error) {
	path := "/v1/health/check"
	urlVals := url.Values{}

	return path + "?" + urlVals.Encode(), nil
}
