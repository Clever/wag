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
	"github.com/Clever/wag/samples/v8/gen-go-errors/models"
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

// statusCodeForGetBook returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetBook(obj interface{}) int {

	switch obj.(type) {

	case *models.ExtendedError:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.ExtendedError:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetBookHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newGetBookInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.ExtendedError{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = input.Validate()

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.ExtendedError{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = h.GetBook(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetBook(err)
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

// newGetBookInput takes in an http.Request an returns the input struct.
func newGetBookInput(r *http.Request) (*models.GetBookInput, error) {
	var input models.GetBookInput

	var err error
	_ = err

	idStr := mux.Vars(r)["id"]
	if len(idStr) == 0 {
		return nil, errors.New("path parameter 'id' must be specified")
	}
	idStrs := []string{idStr}

	if len(idStrs) > 0 {
		var idTmp int64
		idStr := idStrs[0]
		idTmp, err = swag.ConvertInt64(idStr)
		if err != nil {
			return nil, err
		}
		input.ID = idTmp
	}

	return &input, nil
}
