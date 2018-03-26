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
	path := "/v1/students/{student_id}/sections"
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
