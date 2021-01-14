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

// GetSectionsForStudentInput holds the input parameters for a getSectionsForStudent operation.
type GetSectionsForStudentInput struct {
	StudentID string
}

// ValidateGetSectionsForStudentInput returns an error if the input parameter doesn't
// satisfy the requirements in the swagger yml file.
func ValidateGetSectionsForStudentInput(studentID string) error {

	return nil
}

// GetSectionsForStudentInputPath returns the URI path for the input.
func GetSectionsForStudentInputPath(studentID string) (string, error) {
	path := "/students/{student_id}/sections"
	urlVals := url.Values{}

	pathstudent_id := studentID
	if pathstudent_id == "" {
		err := fmt.Errorf("student_id cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{student_id}", pathstudent_id, -1)

	return path + "?" + urlVals.Encode(), nil
}

// PostSectionsForStudentInput holds the input parameters for a postSectionsForStudent operation.
type PostSectionsForStudentInput struct {
	StudentID string
	Sections  string
	UserType  string
}

// Validate returns an error if any of the PostSectionsForStudentInput parameters don't satisfy the
// requirements from the swagger yml file.
func (i PostSectionsForStudentInput) Validate() error {

	if err := validate.Enum("userType", "query", i.UserType, []interface{}{"math", "science", "reading"}); err != nil {
		return err
	}

	return nil
}

// Path returns the URI path for the input.
func (i PostSectionsForStudentInput) Path() (string, error) {
	path := "/students/{student_id}/sections"
	urlVals := url.Values{}

	pathstudent_id := i.StudentID
	if pathstudent_id == "" {
		err := fmt.Errorf("student_id cannot be empty because it's a path parameter")
		if err != nil {
			return "", err
		}
	}
	path = strings.Replace(path, "{student_id}", pathstudent_id, -1)

	urlVals.Add("sections", i.Sections)

	urlVals.Add("userType", i.UserType)

	return path + "?" + urlVals.Encode(), nil
}
