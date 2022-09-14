package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Clever/kayvee-go/v7/logger"
	"github.com/Clever/wag/samples/v9/gen-go-blog/gen-go/models"
	"github.com/go-errors/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gorilla/mux"
	"golang.org/x/xerrors"
)

var _ = strconv.ParseInt
var _ = strfmt.Default
var _ = swag.ConvertInt32
var _ = errors.New
var _ = mux.Vars
var _ = bytes.Compare
var _ = ioutil.ReadAll

var formats = strfmt.Default
var _ = formats

// convertBase64 takes in a string and returns a strfmt.Base64 if the input
// is valid base64 and an error otherwise.
func convertBase64(input string) (strfmt.Base64, error) {
	temp, err := formats.Parse("byte", input)
	if err != nil {
		return strfmt.Base64{}, err
	}
	return *temp.(*strfmt.Base64), nil
}

// convertDateTime takes in a string and returns a strfmt.DateTime if the input
// is a valid DateTime and an error otherwise.
func convertDateTime(input string) (strfmt.DateTime, error) {
	temp, err := formats.Parse("date-time", input)
	if err != nil {
		return strfmt.DateTime{}, err
	}
	return *temp.(*strfmt.DateTime), nil
}

// convertDate takes in a string and returns a strfmt.Date if the input
// is a valid Date and an error otherwise.
func convertDate(input string) (strfmt.Date, error) {
	temp, err := formats.Parse("date", input)
	if err != nil {
		return strfmt.Date{}, err
	}
	return *temp.(*strfmt.Date), nil
}

func jsonMarshalNoError(i interface{}) string {
	bytes, err := json.Marshal(i)
	if err != nil {
		// This should never happen
		return ""
	}
	return string(bytes)
}

// statusCodeForPostGradeFileForStudent returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForPostGradeFileForStudent(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) PostGradeFileForStudentHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newPostGradeFileForStudentInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = input.Validate()

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = h.PostGradeFileForStudent(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForPostGradeFileForStudent(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(""))

}

// newPostGradeFileForStudentInput takes in an http.Request an returns the input struct.
func newPostGradeFileForStudentInput(r *http.Request) (*models.PostGradeFileForStudentInput, error) {
	var input models.PostGradeFileForStudentInput

	var err error
	_ = err

	studentIDStr := mux.Vars(r)["student_id"]
	if len(studentIDStr) == 0 {
		return nil, errors.New("path parameter 'student_id' must be specified")
	}
	studentIDStrs := []string{studentIDStr}

	if len(studentIDStrs) > 0 {
		var studentIDTmp string
		studentIDStr := studentIDStrs[0]
		studentIDTmp, err = studentIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.StudentID = studentIDTmp
	}

	input.File = (*models.GradeFile)(&r.Body)

	return &input, nil
}

// statusCodeForGetSectionsForStudent returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetSectionsForStudent(obj interface{}) int {

	switch obj.(type) {

	case *[]models.Section:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case []models.Section:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) GetSectionsForStudentHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	studentID, err := newGetSectionsForStudentInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetSectionsForStudentInput(studentID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetSectionsForStudent(ctx, studentID)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.Section{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetSectionsForStudent(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetSectionsForStudent(resp))
	w.Write(respBytes)

}

// newGetSectionsForStudentInput takes in an http.Request an returns the student_id parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetSectionsForStudentInput(r *http.Request) (string, error) {
	studentID := mux.Vars(r)["student_id"]
	if len(studentID) == 0 {
		return "", errors.New("Parameter student_id must be specified")
	}
	return studentID, nil
}

// statusCodeForPostSectionsForStudent returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForPostSectionsForStudent(obj interface{}) int {

	switch obj.(type) {

	case *[]models.Section:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case []models.Section:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) PostSectionsForStudentHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newPostSectionsForStudentInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = input.Validate()

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.PostSectionsForStudent(ctx, input)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.Section{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForPostSectionsForStudent(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForPostSectionsForStudent(resp))
	w.Write(respBytes)

}

// newPostSectionsForStudentInput takes in an http.Request an returns the input struct.
func newPostSectionsForStudentInput(r *http.Request) (*models.PostSectionsForStudentInput, error) {
	var input models.PostSectionsForStudentInput

	var err error
	_ = err

	studentIDStr := mux.Vars(r)["student_id"]
	if len(studentIDStr) == 0 {
		return nil, errors.New("path parameter 'student_id' must be specified")
	}
	studentIDStrs := []string{studentIDStr}

	if len(studentIDStrs) > 0 {
		var studentIDTmp string
		studentIDStr := studentIDStrs[0]
		studentIDTmp, err = studentIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.StudentID = studentIDTmp
	}

	sectionsStrs := r.URL.Query()["sections"]
	if len(sectionsStrs) == 0 {
		return nil, errors.New("query parameter 'sections' must be specified")
	}

	if len(sectionsStrs) > 0 {
		var sectionsTmp string
		sectionsStr := sectionsStrs[0]
		sectionsTmp, err = sectionsStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Sections = sectionsTmp
	}

	userTypeStrs := r.URL.Query()["userType"]
	if len(userTypeStrs) == 0 {
		return nil, errors.New("query parameter 'userType' must be specified")
	}

	if len(userTypeStrs) > 0 {
		var userTypeTmp string
		userTypeStr := userTypeStrs[0]
		userTypeTmp, err = userTypeStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.UserType = userTypeTmp
	}

	return &input, nil
}
