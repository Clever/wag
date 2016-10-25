package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Clever/wag/samples/gen-wag-patch/models"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gorilla/mux"
	"gopkg.in/Clever/kayvee-go.v5/logger"
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

// statusCodeForWagpatch returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForWagpatch(obj interface{}) int {

	switch obj.(type) {

	case *models.Data:
		return 200

	case models.Data:
		return 200

	case models.DefaultBadRequest:
		return 400
	case models.DefaultInternalError:
		return 500
	default:
		return -1
	}
}

func (h handler) WagpatchHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newWagpatchInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	err = input.Validate(nil)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.Wagpatch(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		statusCode := statusCodeForWagpatch(err)
		if statusCode != -1 {
			http.Error(w, err.Error(), statusCode)
		} else {
			http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		}
		return
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForWagpatch(resp))
	w.Write(respBytes)

}

// newWagpatchInput takes in an http.Request an returns the input struct.
func newWagpatchInput(r *http.Request) (*models.PatchData, error) {
	var input models.PatchData

	var err error
	_ = err

	data, err := ioutil.ReadAll(r.Body)
	if len(data) > 0 {
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(&input); err != nil {
			return nil, err
		}
	}

	return &input, nil
}

// statusCodeForWagpatch2 returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForWagpatch2(obj interface{}) int {

	switch obj.(type) {

	case *models.Data:
		return 200

	case models.Data:
		return 200

	case models.DefaultBadRequest:
		return 400
	case models.DefaultInternalError:
		return 500
	default:
		return -1
	}
}

func (h handler) Wagpatch2Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newWagpatch2Input(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	err = input.Validate()
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.Wagpatch2(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		statusCode := statusCodeForWagpatch2(err)
		if statusCode != -1 {
			http.Error(w, err.Error(), statusCode)
		} else {
			http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		}
		return
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForWagpatch2(resp))
	w.Write(respBytes)

}

// newWagpatch2Input takes in an http.Request an returns the input struct.
func newWagpatch2Input(r *http.Request) (*models.Wagpatch2Input, error) {
	var input models.Wagpatch2Input

	var err error
	_ = err

	otherStr := r.URL.Query().Get("other")
	if len(otherStr) != 0 {
		var otherTmp string
		otherTmp, err = otherStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Other = &otherTmp

	}
	data, err := ioutil.ReadAll(r.Body)
	if len(data) > 0 {
		input.SpecialPatch = &models.PatchData{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.SpecialPatch); err != nil {
			return nil, err
		}
	}

	return &input, nil
}
