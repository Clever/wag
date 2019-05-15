package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Clever/wag/samples/gen-go-nils/models"
	"github.com/go-errors/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gorilla/mux"
	"golang.org/x/xerrors"
	"gopkg.in/Clever/kayvee-go.v6/logger"
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
	bytes, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		// This should never happen
		return ""
	}
	return string(bytes)
}

// statusCodeForNilCheck returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForNilCheck(obj interface{}) int {

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

func (h handler) NilCheckHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newNilCheckInput(r)
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

	err = h.NilCheck(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForNilCheck(err)
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

// newNilCheckInput takes in an http.Request an returns the input struct.
func newNilCheckInput(r *http.Request) (*models.NilCheckInput, error) {
	var input models.NilCheckInput

	var err error
	_ = err

	idStr := mux.Vars(r)["id"]
	if len(idStr) == 0 {
		return nil, errors.New("path parameter 'id' must be specified")
	}
	idStrs := []string{idStr}

	if len(idStrs) > 0 {
		var idTmp string
		idStr := idStrs[0]
		idTmp, err = idStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.ID = idTmp
	}

	queryStrs := r.URL.Query()["query"]

	if len(queryStrs) > 0 {
		var queryTmp string
		queryStr := queryStrs[0]
		queryTmp, err = queryStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Query = &queryTmp
	}

	headerStrs := r.Header.Get("header")

	if len(headerStrs) > 0 {
		var headerTmp string
		headerTmp = headerStrs
		input.Header = headerTmp
	}
	if array, ok := r.URL.Query()["array"]; ok {
		input.Array = array
	}

	data, err := ioutil.ReadAll(r.Body)

	if len(data) > 0 {

		input.Body = &models.NilFields{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.Body); err != nil {
			return nil, err
		}

	}

	return &input, nil
}
