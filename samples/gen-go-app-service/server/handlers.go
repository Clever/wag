package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Clever/wag/samples/gen-go-app-service/models"
	"github.com/go-errors/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
var _ = log.String

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

// statusCodeForHealthCheck returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForHealthCheck(obj interface{}) int {

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

func (h handler) HealthCheckHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	err := h.HealthCheck(ctx)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForHealthCheck(err)
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

// newHealthCheckInput takes in an http.Request an returns the input struct.
func newHealthCheckInput(r *http.Request) (*models.HealthCheckInput, error) {
	var input models.HealthCheckInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	return &input, nil
}

// statusCodeForGetAdmins returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAdmins(obj interface{}) int {

	switch obj.(type) {

	case *[]models.Admin:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case []models.Admin:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) GetAdminsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newGetAdminsInput(r)
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

	resp, err := h.GetAdmins(ctx, input)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.Admin{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAdmins(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAdmins(resp))
	w.Write(respBytes)

}

// newGetAdminsInput takes in an http.Request an returns the input struct.
func newGetAdminsInput(r *http.Request) (*models.GetAdminsInput, error) {
	var input models.GetAdminsInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	emailStrs := r.URL.Query()["email"]

	if len(emailStrs) > 0 {
		var emailTmp string
		emailStr := emailStrs[0]
		emailTmp, err = emailStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Email = &emailTmp
	}

	passwordStrs := r.URL.Query()["password"]

	if len(passwordStrs) > 0 {
		var passwordTmp string
		passwordStr := passwordStrs[0]
		passwordTmp, err = passwordStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Password = &passwordTmp
	}

	return &input, nil
}

// statusCodeForDeleteAdmin returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForDeleteAdmin(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) DeleteAdminHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	adminID, err := newDeleteAdminInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateDeleteAdminInput(adminID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = h.DeleteAdmin(ctx, adminID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForDeleteAdmin(err)
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

// newDeleteAdminInput takes in an http.Request an returns the adminID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newDeleteAdminInput(r *http.Request) (string, error) {
	adminID := mux.Vars(r)["adminID"]
	if len(adminID) == 0 {
		return "", errors.New("Parameter adminID must be specified")
	}
	return adminID, nil
}

// statusCodeForGetAdminByID returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAdminByID(obj interface{}) int {

	switch obj.(type) {

	case *models.Admin:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.Admin:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetAdminByIDHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	adminID, err := newGetAdminByIDInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetAdminByIDInput(adminID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetAdminByID(ctx, adminID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAdminByID(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAdminByID(resp))
	w.Write(respBytes)

}

// newGetAdminByIDInput takes in an http.Request an returns the adminID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetAdminByIDInput(r *http.Request) (string, error) {
	adminID := mux.Vars(r)["adminID"]
	if len(adminID) == 0 {
		return "", errors.New("Parameter adminID must be specified")
	}
	return adminID, nil
}

// statusCodeForUpdateAdmin returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForUpdateAdmin(obj interface{}) int {

	switch obj.(type) {

	case *models.Admin:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.Admin:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) UpdateAdminHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newUpdateAdminInput(r)
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

	resp, err := h.UpdateAdmin(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForUpdateAdmin(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForUpdateAdmin(resp))
	w.Write(respBytes)

}

// newUpdateAdminInput takes in an http.Request an returns the input struct.
func newUpdateAdminInput(r *http.Request) (*models.UpdateAdminInput, error) {
	var input models.UpdateAdminInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	adminIDStr := mux.Vars(r)["adminID"]
	if len(adminIDStr) == 0 {
		return nil, errors.New("path parameter 'adminID' must be specified")
	}
	adminIDStrs := []string{adminIDStr}

	if len(adminIDStrs) > 0 {
		var adminIDTmp string
		adminIDStr := adminIDStrs[0]
		adminIDTmp, err = adminIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AdminID = adminIDTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.Admin = &models.PatchAdminRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.Admin); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForCreateAdmin returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForCreateAdmin(obj interface{}) int {

	switch obj.(type) {

	case *models.Admin:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case models.Admin:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) CreateAdminHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newCreateAdminInput(r)
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

	resp, err := h.CreateAdmin(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForCreateAdmin(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForCreateAdmin(resp))
	w.Write(respBytes)

}

// newCreateAdminInput takes in an http.Request an returns the input struct.
func newCreateAdminInput(r *http.Request) (*models.CreateAdminInput, error) {
	var input models.CreateAdminInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.CreateAdmin = &models.CreateAdminRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.CreateAdmin); err != nil {
			return nil, err
		}

	}

	adminIDStr := mux.Vars(r)["adminID"]
	if len(adminIDStr) == 0 {
		return nil, errors.New("path parameter 'adminID' must be specified")
	}
	adminIDStrs := []string{adminIDStr}

	if len(adminIDStrs) > 0 {
		var adminIDTmp string
		adminIDStr := adminIDStrs[0]
		adminIDTmp, err = adminIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AdminID = adminIDTmp
	}

	return &input, nil
}

// statusCodeForGetAppsForAdminDeprecated returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAppsForAdminDeprecated(obj interface{}) int {

	switch obj.(type) {

	case *[]models.App:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case []models.App:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetAppsForAdminDeprecatedHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	adminID, err := newGetAppsForAdminDeprecatedInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetAppsForAdminDeprecatedInput(adminID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetAppsForAdminDeprecated(ctx, adminID)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.App{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAppsForAdminDeprecated(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAppsForAdminDeprecated(resp))
	w.Write(respBytes)

}

// newGetAppsForAdminDeprecatedInput takes in an http.Request an returns the adminID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetAppsForAdminDeprecatedInput(r *http.Request) (string, error) {
	adminID := mux.Vars(r)["adminID"]
	if len(adminID) == 0 {
		return "", errors.New("Parameter adminID must be specified")
	}
	return adminID, nil
}

// statusCodeForVerifyCode returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForVerifyCode(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) VerifyCodeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newVerifyCodeInput(r)
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

	err = h.VerifyCode(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForVerifyCode(err)
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

// newVerifyCodeInput takes in an http.Request an returns the input struct.
func newVerifyCodeInput(r *http.Request) (*models.VerifyCodeInput, error) {
	var input models.VerifyCodeInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	codeStrs := r.URL.Query()["code"]
	if len(codeStrs) == 0 {
		return nil, errors.New("query parameter 'code' must be specified")
	}

	if len(codeStrs) > 0 {
		var codeTmp string
		codeStr := codeStrs[0]
		codeTmp, err = codeStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Code = codeTmp
	}

	invalidateStrs := r.URL.Query()["invalidate"]

	if len(invalidateStrs) > 0 {
		var invalidateTmp bool
		invalidateStr := invalidateStrs[0]
		invalidateTmp, err = strconv.ParseBool(invalidateStr)
		if err != nil {
			return nil, err
		}
		input.Invalidate = &invalidateTmp
	}

	adminIDStr := mux.Vars(r)["adminID"]
	if len(adminIDStr) == 0 {
		return nil, errors.New("path parameter 'adminID' must be specified")
	}
	adminIDStrs := []string{adminIDStr}

	if len(adminIDStrs) > 0 {
		var adminIDTmp string
		adminIDStr := adminIDStrs[0]
		adminIDTmp, err = adminIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AdminID = adminIDTmp
	}

	return &input, nil
}

// statusCodeForCreateVerificationCode returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForCreateVerificationCode(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.VerificationCodeResponse:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.VerificationCodeResponse:
		return 200

	default:
		return -1
	}
}

func (h handler) CreateVerificationCodeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newCreateVerificationCodeInput(r)
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

	resp, err := h.CreateVerificationCode(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForCreateVerificationCode(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForCreateVerificationCode(resp))
	w.Write(respBytes)

}

// newCreateVerificationCodeInput takes in an http.Request an returns the input struct.
func newCreateVerificationCodeInput(r *http.Request) (*models.CreateVerificationCodeInput, error) {
	var input models.CreateVerificationCodeInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	durationStrs := r.URL.Query()["duration"]
	if len(durationStrs) == 0 {
		return nil, errors.New("query parameter 'duration' must be specified")
	}

	if len(durationStrs) > 0 {
		var durationTmp int32
		durationStr := durationStrs[0]
		durationTmp, err = swag.ConvertInt32(durationStr)
		if err != nil {
			return nil, err
		}
		input.Duration = durationTmp
	}

	adminIDStr := mux.Vars(r)["adminID"]
	if len(adminIDStr) == 0 {
		return nil, errors.New("path parameter 'adminID' must be specified")
	}
	adminIDStrs := []string{adminIDStr}

	if len(adminIDStrs) > 0 {
		var adminIDTmp string
		adminIDStr := adminIDStrs[0]
		adminIDTmp, err = adminIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AdminID = adminIDTmp
	}

	return &input, nil
}

// statusCodeForVerifyAdminEmail returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForVerifyAdminEmail(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) VerifyAdminEmailHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newVerifyAdminEmailInput(r)
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

	err = h.VerifyAdminEmail(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForVerifyAdminEmail(err)
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

// newVerifyAdminEmailInput takes in an http.Request an returns the input struct.
func newVerifyAdminEmailInput(r *http.Request) (*models.VerifyAdminEmailInput, error) {
	var input models.VerifyAdminEmailInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	adminIDStr := mux.Vars(r)["adminID"]
	if len(adminIDStr) == 0 {
		return nil, errors.New("path parameter 'adminID' must be specified")
	}
	adminIDStrs := []string{adminIDStr}

	if len(adminIDStrs) > 0 {
		var adminIDTmp string
		adminIDStr := adminIDStrs[0]
		adminIDTmp, err = adminIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AdminID = adminIDTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.Request = &models.VerifyAdminEmailRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.Request); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForGetAllAnalyticsApps returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAllAnalyticsApps(obj interface{}) int {

	switch obj.(type) {

	case *models.AnalyticsApps:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.AnalyticsApps:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetAllAnalyticsAppsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	resp, err := h.GetAllAnalyticsApps(ctx)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAllAnalyticsApps(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAllAnalyticsApps(resp))
	w.Write(respBytes)

}

// newGetAllAnalyticsAppsInput takes in an http.Request an returns the input struct.
func newGetAllAnalyticsAppsInput(r *http.Request) (*models.GetAllAnalyticsAppsInput, error) {
	var input models.GetAllAnalyticsAppsInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	return &input, nil
}

// statusCodeForGetAnalyticsAppByShortname returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAnalyticsAppByShortname(obj interface{}) int {

	switch obj.(type) {

	case *models.AnalyticsApp:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.AnalyticsApp:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetAnalyticsAppByShortnameHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	shortname, err := newGetAnalyticsAppByShortnameInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetAnalyticsAppByShortnameInput(shortname)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetAnalyticsAppByShortname(ctx, shortname)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAnalyticsAppByShortname(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAnalyticsAppByShortname(resp))
	w.Write(respBytes)

}

// newGetAnalyticsAppByShortnameInput takes in an http.Request an returns the shortname parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetAnalyticsAppByShortnameInput(r *http.Request) (string, error) {
	shortname := mux.Vars(r)["shortname"]
	if len(shortname) == 0 {
		return "", errors.New("Parameter shortname must be specified")
	}
	return shortname, nil
}

// statusCodeForGetAllTrackableApps returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAllTrackableApps(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.TrackableApps:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.TrackableApps:
		return 200

	default:
		return -1
	}
}

func (h handler) GetAllTrackableAppsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	resp, err := h.GetAllTrackableApps(ctx)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAllTrackableApps(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAllTrackableApps(resp))
	w.Write(respBytes)

}

// newGetAllTrackableAppsInput takes in an http.Request an returns the input struct.
func newGetAllTrackableAppsInput(r *http.Request) (*models.GetAllTrackableAppsInput, error) {
	var input models.GetAllTrackableAppsInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	return &input, nil
}

// statusCodeForGetAnalyticsUsageUrls returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAnalyticsUsageUrls(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.UsageUrls:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.UsageUrls:
		return 200

	default:
		return -1
	}
}

func (h handler) GetAnalyticsUsageUrlsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	resp, err := h.GetAnalyticsUsageUrls(ctx)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAnalyticsUsageUrls(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAnalyticsUsageUrls(resp))
	w.Write(respBytes)

}

// newGetAnalyticsUsageUrlsInput takes in an http.Request an returns the input struct.
func newGetAnalyticsUsageUrlsInput(r *http.Request) (*models.GetAnalyticsUsageUrlsInput, error) {
	var input models.GetAnalyticsUsageUrlsInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	return &input, nil
}

// statusCodeForGetAllUsageUrls returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAllUsageUrls(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.UsageUrls:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.UsageUrls:
		return 200

	default:
		return -1
	}
}

func (h handler) GetAllUsageUrlsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	resp, err := h.GetAllUsageUrls(ctx)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAllUsageUrls(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAllUsageUrls(resp))
	w.Write(respBytes)

}

// newGetAllUsageUrlsInput takes in an http.Request an returns the input struct.
func newGetAllUsageUrlsInput(r *http.Request) (*models.GetAllUsageUrlsInput, error) {
	var input models.GetAllUsageUrlsInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	return &input, nil
}

// statusCodeForGetApps returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetApps(obj interface{}) int {

	switch obj.(type) {

	case *[]models.App:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case []models.App:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) GetAppsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newGetAppsInput(r)
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

	resp, err := h.GetApps(ctx, input)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.App{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetApps(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetApps(resp))
	w.Write(respBytes)

}

// newGetAppsInput takes in an http.Request an returns the input struct.
func newGetAppsInput(r *http.Request) (*models.GetAppsInput, error) {
	var input models.GetAppsInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err
	if ids, ok := r.URL.Query()["ids"]; ok {
		input.Ids = ids
	}

	clientIdStrs := r.URL.Query()["clientId"]

	if len(clientIdStrs) > 0 {
		var clientIdTmp string
		clientIdStr := clientIdStrs[0]
		clientIdTmp, err = clientIdStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.ClientId = &clientIdTmp
	}

	clientSecretStrs := r.URL.Query()["clientSecret"]

	if len(clientSecretStrs) > 0 {
		var clientSecretTmp string
		clientSecretStr := clientSecretStrs[0]
		clientSecretTmp, err = clientSecretStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.ClientSecret = &clientSecretTmp
	}

	shortnameStrs := r.URL.Query()["shortname"]

	if len(shortnameStrs) > 0 {
		var shortnameTmp string
		shortnameStr := shortnameStrs[0]
		shortnameTmp, err = shortnameStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Shortname = &shortnameTmp
	}

	businessTokenStrs := r.URL.Query()["businessToken"]

	if len(businessTokenStrs) > 0 {
		var businessTokenTmp string
		businessTokenStr := businessTokenStrs[0]
		businessTokenTmp, err = businessTokenStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.BusinessToken = &businessTokenTmp
	}
	if tags, ok := r.URL.Query()["tags"]; ok {
		input.Tags = tags
	}
	if skipTags, ok := r.URL.Query()["skipTags"]; ok {
		input.SkipTags = skipTags
	}

	return &input, nil
}

// statusCodeForDeleteApp returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForDeleteApp(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) DeleteAppHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	appID, err := newDeleteAppInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateDeleteAppInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = h.DeleteApp(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForDeleteApp(err)
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

// newDeleteAppInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newDeleteAppInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForGetAppByID returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAppByID(obj interface{}) int {

	switch obj.(type) {

	case *models.App:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.App:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetAppByIDHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGetAppByIDInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetAppByIDInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetAppByID(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAppByID(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAppByID(resp))
	w.Write(respBytes)

}

// newGetAppByIDInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetAppByIDInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForUpdateApp returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForUpdateApp(obj interface{}) int {

	switch obj.(type) {

	case *models.App:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.App:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) UpdateAppHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newUpdateAppInput(r)
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

	resp, err := h.UpdateApp(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForUpdateApp(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForUpdateApp(resp))
	w.Write(respBytes)

}

// newUpdateAppInput takes in an http.Request an returns the input struct.
func newUpdateAppInput(r *http.Request) (*models.UpdateAppInput, error) {
	var input models.UpdateAppInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	withSchemaPropagationStrs := r.URL.Query()["withSchemaPropagation"]

	if len(withSchemaPropagationStrs) > 0 {
		var withSchemaPropagationTmp bool
		withSchemaPropagationStr := withSchemaPropagationStrs[0]
		withSchemaPropagationTmp, err = strconv.ParseBool(withSchemaPropagationStr)
		if err != nil {
			return nil, err
		}
		input.WithSchemaPropagation = &withSchemaPropagationTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.App = &models.PatchAppRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.App); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForCreateApp returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForCreateApp(obj interface{}) int {

	switch obj.(type) {

	case *models.App:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case models.App:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) CreateAppHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newCreateAppInput(r)
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

	resp, err := h.CreateApp(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForCreateApp(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForCreateApp(resp))
	w.Write(respBytes)

}

// newCreateAppInput takes in an http.Request an returns the input struct.
func newCreateAppInput(r *http.Request) (*models.CreateAppInput, error) {
	var input models.CreateAppInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	data, err := ioutil.ReadAll(r.Body)

	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.App = &models.App{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.App); err != nil {
			return nil, err
		}

	}

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	return &input, nil
}

// statusCodeForGetAdminsForApp returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAdminsForApp(obj interface{}) int {

	switch obj.(type) {

	case *[]models.AppAdminResponse:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case []models.AppAdminResponse:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetAdminsForAppHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGetAdminsForAppInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetAdminsForAppInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetAdminsForApp(ctx, appID)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.AppAdminResponse{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAdminsForApp(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAdminsForApp(resp))
	w.Write(respBytes)

}

// newGetAdminsForAppInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetAdminsForAppInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForUnlinkAppAdmin returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForUnlinkAppAdmin(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Forbidden:
		return 403

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.Forbidden:
		return 403

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) UnlinkAppAdminHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newUnlinkAppAdminInput(r)
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

	err = h.UnlinkAppAdmin(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForUnlinkAppAdmin(err)
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

// newUnlinkAppAdminInput takes in an http.Request an returns the input struct.
func newUnlinkAppAdminInput(r *http.Request) (*models.UnlinkAppAdminInput, error) {
	var input models.UnlinkAppAdminInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	adminIDStr := mux.Vars(r)["adminID"]
	if len(adminIDStr) == 0 {
		return nil, errors.New("path parameter 'adminID' must be specified")
	}
	adminIDStrs := []string{adminIDStr}

	if len(adminIDStrs) > 0 {
		var adminIDTmp string
		adminIDStr := adminIDStrs[0]
		adminIDTmp, err = adminIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AdminID = adminIDTmp
	}

	return &input, nil
}

// statusCodeForLinkAppAdmin returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForLinkAppAdmin(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Forbidden:
		return 403

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.Forbidden:
		return 403

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) LinkAppAdminHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newLinkAppAdminInput(r)
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

	err = h.LinkAppAdmin(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForLinkAppAdmin(err)
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

// newLinkAppAdminInput takes in an http.Request an returns the input struct.
func newLinkAppAdminInput(r *http.Request) (*models.LinkAppAdminInput, error) {
	var input models.LinkAppAdminInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	adminIDStr := mux.Vars(r)["adminID"]
	if len(adminIDStr) == 0 {
		return nil, errors.New("path parameter 'adminID' must be specified")
	}
	adminIDStrs := []string{adminIDStr}

	if len(adminIDStrs) > 0 {
		var adminIDTmp string
		adminIDStr := adminIDStrs[0]
		adminIDTmp, err = adminIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AdminID = adminIDTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.Permissions = &models.PermissionList{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.Permissions); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForGetGuideConfig returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetGuideConfig(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Forbidden:
		return 403

	case *models.GuideConfig:
		return 200

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.Forbidden:
		return 403

	case models.GuideConfig:
		return 200

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetGuideConfigHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newGetGuideConfigInput(r)
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

	resp, err := h.GetGuideConfig(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetGuideConfig(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetGuideConfig(resp))
	w.Write(respBytes)

}

// newGetGuideConfigInput takes in an http.Request an returns the input struct.
func newGetGuideConfigInput(r *http.Request) (*models.GetGuideConfigInput, error) {
	var input models.GetGuideConfigInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	adminIDStr := mux.Vars(r)["adminID"]
	if len(adminIDStr) == 0 {
		return nil, errors.New("path parameter 'adminID' must be specified")
	}
	adminIDStrs := []string{adminIDStr}

	if len(adminIDStrs) > 0 {
		var adminIDTmp string
		adminIDStr := adminIDStrs[0]
		adminIDTmp, err = adminIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AdminID = adminIDTmp
	}

	guideIDStr := mux.Vars(r)["guideID"]
	if len(guideIDStr) == 0 {
		return nil, errors.New("path parameter 'guideID' must be specified")
	}
	guideIDStrs := []string{guideIDStr}

	if len(guideIDStrs) > 0 {
		var guideIDTmp string
		guideIDStr := guideIDStrs[0]
		guideIDTmp, err = guideIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.GuideID = guideIDTmp
	}

	return &input, nil
}

// statusCodeForSetGuideConfig returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForSetGuideConfig(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Forbidden:
		return 403

	case *models.GuideConfig:
		return 200

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.Forbidden:
		return 403

	case models.GuideConfig:
		return 200

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) SetGuideConfigHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newSetGuideConfigInput(r)
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

	resp, err := h.SetGuideConfig(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForSetGuideConfig(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForSetGuideConfig(resp))
	w.Write(respBytes)

}

// newSetGuideConfigInput takes in an http.Request an returns the input struct.
func newSetGuideConfigInput(r *http.Request) (*models.SetGuideConfigInput, error) {
	var input models.SetGuideConfigInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	adminIDStr := mux.Vars(r)["adminID"]
	if len(adminIDStr) == 0 {
		return nil, errors.New("path parameter 'adminID' must be specified")
	}
	adminIDStrs := []string{adminIDStr}

	if len(adminIDStrs) > 0 {
		var adminIDTmp string
		adminIDStr := adminIDStrs[0]
		adminIDTmp, err = adminIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AdminID = adminIDTmp
	}

	guideIDStr := mux.Vars(r)["guideID"]
	if len(guideIDStr) == 0 {
		return nil, errors.New("path parameter 'guideID' must be specified")
	}
	guideIDStrs := []string{guideIDStr}

	if len(guideIDStrs) > 0 {
		var guideIDTmp string
		guideIDStr := guideIDStrs[0]
		guideIDTmp, err = guideIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.GuideID = guideIDTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.GuideConfig = &models.GuideConfig{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.GuideConfig); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForGetPermissionsForAdmin returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetPermissionsForAdmin(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.PermissionList:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.PermissionList:
		return 200

	default:
		return -1
	}
}

func (h handler) GetPermissionsForAdminHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newGetPermissionsForAdminInput(r)
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

	resp, err := h.GetPermissionsForAdmin(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetPermissionsForAdmin(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetPermissionsForAdmin(resp))
	w.Write(respBytes)

}

// newGetPermissionsForAdminInput takes in an http.Request an returns the input struct.
func newGetPermissionsForAdminInput(r *http.Request) (*models.GetPermissionsForAdminInput, error) {
	var input models.GetPermissionsForAdminInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	adminIDStr := mux.Vars(r)["adminID"]
	if len(adminIDStr) == 0 {
		return nil, errors.New("path parameter 'adminID' must be specified")
	}
	adminIDStrs := []string{adminIDStr}

	if len(adminIDStrs) > 0 {
		var adminIDTmp string
		adminIDStr := adminIDStrs[0]
		adminIDTmp, err = adminIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AdminID = adminIDTmp
	}

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	return &input, nil
}

// statusCodeForVerifyAppAdmin returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForVerifyAppAdmin(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Forbidden:
		return 403

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.Forbidden:
		return 403

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) VerifyAppAdminHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newVerifyAppAdminInput(r)
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

	err = h.VerifyAppAdmin(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForVerifyAppAdmin(err)
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

// newVerifyAppAdminInput takes in an http.Request an returns the input struct.
func newVerifyAppAdminInput(r *http.Request) (*models.VerifyAppAdminInput, error) {
	var input models.VerifyAppAdminInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	adminIDStr := mux.Vars(r)["adminID"]
	if len(adminIDStr) == 0 {
		return nil, errors.New("path parameter 'adminID' must be specified")
	}
	adminIDStrs := []string{adminIDStr}

	if len(adminIDStrs) > 0 {
		var adminIDTmp string
		adminIDStr := adminIDStrs[0]
		adminIDTmp, err = adminIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AdminID = adminIDTmp
	}

	verifiedStrs := r.URL.Query()["verified"]
	if len(verifiedStrs) == 0 {
		return nil, errors.New("query parameter 'verified' must be specified")
	}

	if len(verifiedStrs) > 0 {
		var verifiedTmp bool
		verifiedStr := verifiedStrs[0]
		verifiedTmp, err = strconv.ParseBool(verifiedStr)
		if err != nil {
			return nil, err
		}
		input.Verified = verifiedTmp
	}

	return &input, nil
}

// statusCodeForGenerateNewBusinessToken returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGenerateNewBusinessToken(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.SecretConfig:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.SecretConfig:
		return 200

	default:
		return -1
	}
}

func (h handler) GenerateNewBusinessTokenHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGenerateNewBusinessTokenInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGenerateNewBusinessTokenInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GenerateNewBusinessToken(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGenerateNewBusinessToken(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGenerateNewBusinessToken(resp))
	w.Write(respBytes)

}

// newGenerateNewBusinessTokenInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGenerateNewBusinessTokenInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForGetCertifications returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetCertifications(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Certifications:
		return 200

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.Certifications:
		return 200

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetCertificationsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newGetCertificationsInput(r)
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

	resp, err := h.GetCertifications(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetCertifications(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetCertifications(resp))
	w.Write(respBytes)

}

// newGetCertificationsInput takes in an http.Request an returns the input struct.
func newGetCertificationsInput(r *http.Request) (*models.GetCertificationsInput, error) {
	var input models.GetCertificationsInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	schoolYearStartStr := mux.Vars(r)["schoolYearStart"]
	if len(schoolYearStartStr) == 0 {
		return nil, errors.New("path parameter 'schoolYearStart' must be specified")
	}
	schoolYearStartStrs := []string{schoolYearStartStr}

	if len(schoolYearStartStrs) > 0 {
		var schoolYearStartTmp int32
		schoolYearStartStr := schoolYearStartStrs[0]
		schoolYearStartTmp, err = swag.ConvertInt32(schoolYearStartStr)
		if err != nil {
			return nil, err
		}
		input.SchoolYearStart = schoolYearStartTmp
	}

	return &input, nil
}

// statusCodeForSetCertifications returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForSetCertifications(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Certifications:
		return 200

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.Certifications:
		return 200

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) SetCertificationsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newSetCertificationsInput(r)
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

	resp, err := h.SetCertifications(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForSetCertifications(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForSetCertifications(resp))
	w.Write(respBytes)

}

// newSetCertificationsInput takes in an http.Request an returns the input struct.
func newSetCertificationsInput(r *http.Request) (*models.SetCertificationsInput, error) {
	var input models.SetCertificationsInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	schoolYearStartStr := mux.Vars(r)["schoolYearStart"]
	if len(schoolYearStartStr) == 0 {
		return nil, errors.New("path parameter 'schoolYearStart' must be specified")
	}
	schoolYearStartStrs := []string{schoolYearStartStr}

	if len(schoolYearStartStrs) > 0 {
		var schoolYearStartTmp int32
		schoolYearStartStr := schoolYearStartStrs[0]
		schoolYearStartTmp, err = swag.ConvertInt32(schoolYearStartStr)
		if err != nil {
			return nil, err
		}
		input.SchoolYearStart = schoolYearStartTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.Certifications = &models.SetCertificationsRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.Certifications); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForGetSetupStep returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetSetupStep(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.SetupStep:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.SetupStep:
		return 200

	default:
		return -1
	}
}

func (h handler) GetSetupStepHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGetSetupStepInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetSetupStepInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetSetupStep(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetSetupStep(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetSetupStep(resp))
	w.Write(respBytes)

}

// newGetSetupStepInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetSetupStepInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForCreateSetupStep returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForCreateSetupStep(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) CreateSetupStepHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newCreateSetupStepInput(r)
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

	err = h.CreateSetupStep(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForCreateSetupStep(err)
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

// newCreateSetupStepInput takes in an http.Request an returns the input struct.
func newCreateSetupStepInput(r *http.Request) (*models.CreateSetupStepInput, error) {
	var input models.CreateSetupStepInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	data, err := ioutil.ReadAll(r.Body)

	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.SetupStep = &models.SetupStep{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.SetupStep); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForGetDataRules returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetDataRules(obj interface{}) int {

	switch obj.(type) {

	case *[]models.DataRule:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case []models.DataRule:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetDataRulesHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGetDataRulesInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetDataRulesInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetDataRules(ctx, appID)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.DataRule{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetDataRules(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetDataRules(resp))
	w.Write(respBytes)

}

// newGetDataRulesInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetDataRulesInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForSetDataRules returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForSetDataRules(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) SetDataRulesHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newSetDataRulesInput(r)
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

	err = h.SetDataRules(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForSetDataRules(err)
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

// newSetDataRulesInput takes in an http.Request an returns the input struct.
func newSetDataRulesInput(r *http.Request) (*models.SetDataRulesInput, error) {
	var input models.SetDataRulesInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	data, err := ioutil.ReadAll(r.Body)

	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.Rules = &models.SetDataRulesRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.Rules); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForGetManagers returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetManagers(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.Managers:
		return 200

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.Managers:
		return 200

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetManagersHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGetManagersInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetManagersInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetManagers(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetManagers(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetManagers(resp))
	w.Write(respBytes)

}

// newGetManagersInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetManagersInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForGetOnboarding returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetOnboarding(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.Onboarding:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.Onboarding:
		return 200

	default:
		return -1
	}
}

func (h handler) GetOnboardingHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGetOnboardingInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetOnboardingInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetOnboarding(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetOnboarding(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetOnboarding(resp))
	w.Write(respBytes)

}

// newGetOnboardingInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetOnboardingInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForUpdateOnboarding returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForUpdateOnboarding(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) UpdateOnboardingHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newUpdateOnboardingInput(r)
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

	err = h.UpdateOnboarding(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForUpdateOnboarding(err)
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

// newUpdateOnboardingInput takes in an http.Request an returns the input struct.
func newUpdateOnboardingInput(r *http.Request) (*models.UpdateOnboardingInput, error) {
	var input models.UpdateOnboardingInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.Update = &models.UpdateOnboardingRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.Update); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForInitializeOnboarding returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForInitializeOnboarding(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) InitializeOnboardingHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	appID, err := newInitializeOnboardingInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateInitializeOnboardingInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = h.InitializeOnboarding(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForInitializeOnboarding(err)
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

// newInitializeOnboardingInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newInitializeOnboardingInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForDeletePlatform returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForDeletePlatform(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) DeletePlatformHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newDeletePlatformInput(r)
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

	err = h.DeletePlatform(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForDeletePlatform(err)
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

// newDeletePlatformInput takes in an http.Request an returns the input struct.
func newDeletePlatformInput(r *http.Request) (*models.DeletePlatformInput, error) {
	var input models.DeletePlatformInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	clientIDStr := mux.Vars(r)["clientID"]
	if len(clientIDStr) == 0 {
		return nil, errors.New("path parameter 'clientID' must be specified")
	}
	clientIDStrs := []string{clientIDStr}

	if len(clientIDStrs) > 0 {
		var clientIDTmp string
		clientIDStr := clientIDStrs[0]
		clientIDTmp, err = clientIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.ClientID = clientIDTmp
	}

	return &input, nil
}

// statusCodeForUpdatePlatform returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForUpdatePlatform(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.Platform:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.Platform:
		return 200

	default:
		return -1
	}
}

func (h handler) UpdatePlatformHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newUpdatePlatformInput(r)
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

	resp, err := h.UpdatePlatform(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForUpdatePlatform(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForUpdatePlatform(resp))
	w.Write(respBytes)

}

// newUpdatePlatformInput takes in an http.Request an returns the input struct.
func newUpdatePlatformInput(r *http.Request) (*models.UpdatePlatformInput, error) {
	var input models.UpdatePlatformInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	clientIDStr := mux.Vars(r)["clientID"]
	if len(clientIDStr) == 0 {
		return nil, errors.New("path parameter 'clientID' must be specified")
	}
	clientIDStrs := []string{clientIDStr}

	if len(clientIDStrs) > 0 {
		var clientIDTmp string
		clientIDStr := clientIDStrs[0]
		clientIDTmp, err = clientIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.ClientID = clientIDTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.Request = &models.PatchPlatformRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.Request); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForGetPlatformsByAppID returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetPlatformsByAppID(obj interface{}) int {

	switch obj.(type) {

	case *[]models.Platform:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case []models.Platform:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetPlatformsByAppIDHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGetPlatformsByAppIDInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetPlatformsByAppIDInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetPlatformsByAppID(ctx, appID)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.Platform{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetPlatformsByAppID(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetPlatformsByAppID(resp))
	w.Write(respBytes)

}

// newGetPlatformsByAppIDInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetPlatformsByAppIDInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForCreatePlatform returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForCreatePlatform(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.Platform:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.Platform:
		return 200

	default:
		return -1
	}
}

func (h handler) CreatePlatformHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newCreatePlatformInput(r)
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

	resp, err := h.CreatePlatform(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForCreatePlatform(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForCreatePlatform(resp))
	w.Write(respBytes)

}

// newCreatePlatformInput takes in an http.Request an returns the input struct.
func newCreatePlatformInput(r *http.Request) (*models.CreatePlatformInput, error) {
	var input models.CreatePlatformInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.Request = &models.CreatePlatformRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.Request); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForDeleteAppSchema returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForDeleteAppSchema(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) DeleteAppSchemaHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newDeleteAppSchemaInput(r)
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

	err = h.DeleteAppSchema(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForDeleteAppSchema(err)
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

// newDeleteAppSchemaInput takes in an http.Request an returns the input struct.
func newDeleteAppSchemaInput(r *http.Request) (*models.DeleteAppSchemaInput, error) {
	var input models.DeleteAppSchemaInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	deleteDataRulesStrs := r.URL.Query()["deleteDataRules"]

	if len(deleteDataRulesStrs) > 0 {
		var deleteDataRulesTmp bool
		deleteDataRulesStr := deleteDataRulesStrs[0]
		deleteDataRulesTmp, err = strconv.ParseBool(deleteDataRulesStr)
		if err != nil {
			return nil, err
		}
		input.DeleteDataRules = &deleteDataRulesTmp
	}

	return &input, nil
}

// statusCodeForGetAppSchema returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAppSchema(obj interface{}) int {

	switch obj.(type) {

	case *models.AppSchema:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.AppSchema:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetAppSchemaHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGetAppSchemaInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetAppSchemaInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetAppSchema(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAppSchema(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAppSchema(resp))
	w.Write(respBytes)

}

// newGetAppSchemaInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetAppSchemaInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForCreateAppSchema returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForCreateAppSchema(obj interface{}) int {

	switch obj.(type) {

	case *models.AppSchema:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.AppSchema:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) CreateAppSchemaHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newCreateAppSchemaInput(r)
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

	resp, err := h.CreateAppSchema(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForCreateAppSchema(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForCreateAppSchema(resp))
	w.Write(respBytes)

}

// newCreateAppSchemaInput takes in an http.Request an returns the input struct.
func newCreateAppSchemaInput(r *http.Request) (*models.CreateAppSchemaInput, error) {
	var input models.CreateAppSchemaInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	skipPropagationStrs := r.URL.Query()["skipPropagation"]

	if len(skipPropagationStrs) > 0 {
		var skipPropagationTmp bool
		skipPropagationStr := skipPropagationStrs[0]
		skipPropagationTmp, err = strconv.ParseBool(skipPropagationStr)
		if err != nil {
			return nil, err
		}
		input.SkipPropagation = &skipPropagationTmp
	}

	updateDataRulesStrs := r.URL.Query()["updateDataRules"]

	if len(updateDataRulesStrs) > 0 {
		var updateDataRulesTmp bool
		updateDataRulesStr := updateDataRulesStrs[0]
		updateDataRulesTmp, err = strconv.ParseBool(updateDataRulesStr)
		if err != nil {
			return nil, err
		}
		input.UpdateDataRules = &updateDataRulesTmp
	}

	return &input, nil
}

// statusCodeForSetAppSchema returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForSetAppSchema(obj interface{}) int {

	switch obj.(type) {

	case *models.AppSchema:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.AppSchema:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) SetAppSchemaHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newSetAppSchemaInput(r)
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

	resp, err := h.SetAppSchema(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForSetAppSchema(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForSetAppSchema(resp))
	w.Write(respBytes)

}

// newSetAppSchemaInput takes in an http.Request an returns the input struct.
func newSetAppSchemaInput(r *http.Request) (*models.SetAppSchemaInput, error) {
	var input models.SetAppSchemaInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	skipPropagationStrs := r.URL.Query()["skipPropagation"]

	if len(skipPropagationStrs) > 0 {
		var skipPropagationTmp bool
		skipPropagationStr := skipPropagationStrs[0]
		skipPropagationTmp, err = strconv.ParseBool(skipPropagationStr)
		if err != nil {
			return nil, err
		}
		input.SkipPropagation = &skipPropagationTmp
	}

	updateDataRulesStrs := r.URL.Query()["updateDataRules"]

	if len(updateDataRulesStrs) > 0 {
		var updateDataRulesTmp bool
		updateDataRulesStr := updateDataRulesStrs[0]
		updateDataRulesTmp, err = strconv.ParseBool(updateDataRulesStr)
		if err != nil {
			return nil, err
		}
		input.UpdateDataRules = &updateDataRulesTmp
	}

	data, err := ioutil.ReadAll(r.Body)

	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.AppSchema = &models.AppSchema{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.AppSchema); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForGetSecrets returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetSecrets(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.SecretConfig:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.SecretConfig:
		return 200

	default:
		return -1
	}
}

func (h handler) GetSecretsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGetSecretsInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetSecretsInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetSecrets(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetSecrets(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetSecrets(resp))
	w.Write(respBytes)

}

// newGetSecretsInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetSecretsInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForRevokeOldClientSecret returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForRevokeOldClientSecret(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.SecretConfig:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.SecretConfig:
		return 200

	default:
		return -1
	}
}

func (h handler) RevokeOldClientSecretHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newRevokeOldClientSecretInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateRevokeOldClientSecretInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.RevokeOldClientSecret(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForRevokeOldClientSecret(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForRevokeOldClientSecret(resp))
	w.Write(respBytes)

}

// newRevokeOldClientSecretInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newRevokeOldClientSecretInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForGenerateNewClientSecret returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGenerateNewClientSecret(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.SecretConfig:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.SecretConfig:
		return 200

	default:
		return -1
	}
}

func (h handler) GenerateNewClientSecretHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGenerateNewClientSecretInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGenerateNewClientSecretInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GenerateNewClientSecret(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGenerateNewClientSecret(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGenerateNewClientSecret(resp))
	w.Write(respBytes)

}

// newGenerateNewClientSecretInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGenerateNewClientSecretInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForResetClientSecret returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForResetClientSecret(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.SecretConfig:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.SecretConfig:
		return 200

	default:
		return -1
	}
}

func (h handler) ResetClientSecretHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newResetClientSecretInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateResetClientSecretInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.ResetClientSecret(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForResetClientSecret(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForResetClientSecret(resp))
	w.Write(respBytes)

}

// newResetClientSecretInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newResetClientSecretInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForGetRecommendedSharing returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetRecommendedSharing(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.SharingRecommendations:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.SharingRecommendations:
		return 200

	default:
		return -1
	}
}

func (h handler) GetRecommendedSharingHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	appID, err := newGetRecommendedSharingInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetRecommendedSharingInput(appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetRecommendedSharing(ctx, appID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetRecommendedSharing(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetRecommendedSharing(resp))
	w.Write(respBytes)

}

// newGetRecommendedSharingInput takes in an http.Request an returns the appID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetRecommendedSharingInput(r *http.Request) (string, error) {
	appID := mux.Vars(r)["appID"]
	if len(appID) == 0 {
		return "", errors.New("Parameter appID must be specified")
	}
	return appID, nil
}

// statusCodeForSetRecommendedSharing returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForSetRecommendedSharing(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) SetRecommendedSharingHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newSetRecommendedSharingInput(r)
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

	err = h.SetRecommendedSharing(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForSetRecommendedSharing(err)
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

// newSetRecommendedSharingInput takes in an http.Request an returns the input struct.
func newSetRecommendedSharingInput(r *http.Request) (*models.SetRecommendedSharingInput, error) {
	var input models.SetRecommendedSharingInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	data, err := ioutil.ReadAll(r.Body)

	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.Recommendations = &models.SharingRecommendations{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.Recommendations); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForUpdateAppIcon returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForUpdateAppIcon(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Image:
		return 200

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.UnprocessableEntity:
		return 422

	case models.BadRequest:
		return 400

	case models.Image:
		return 200

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.UnprocessableEntity:
		return 422

	default:
		return -1
	}
}

func (h handler) UpdateAppIconHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newUpdateAppIconInput(r)
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

	resp, err := h.UpdateAppIcon(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForUpdateAppIcon(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForUpdateAppIcon(resp))
	w.Write(respBytes)

}

// newUpdateAppIconInput takes in an http.Request an returns the input struct.
func newUpdateAppIconInput(r *http.Request) (*models.UpdateAppIconInput, error) {
	var input models.UpdateAppIconInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	appIDStr := mux.Vars(r)["appID"]
	if len(appIDStr) == 0 {
		return nil, errors.New("path parameter 'appID' must be specified")
	}
	appIDStrs := []string{appIDStr}

	if len(appIDStrs) > 0 {
		var appIDTmp string
		appIDStr := appIDStrs[0]
		appIDTmp, err = appIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AppID = appIDTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.App = &models.UpdateAppIconRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.App); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForGetAllCategories returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAllCategories(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Categories:
		return 200

	case *models.InternalError:
		return 500

	case models.BadRequest:
		return 400

	case models.Categories:
		return 200

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) GetAllCategoriesHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	resp, err := h.GetAllCategories(ctx)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAllCategories(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAllCategories(resp))
	w.Write(respBytes)

}

// newGetAllCategoriesInput takes in an http.Request an returns the input struct.
func newGetAllCategoriesInput(r *http.Request) (*models.GetAllCategoriesInput, error) {
	var input models.GetAllCategoriesInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	return &input, nil
}

// statusCodeForGetKnownHosts returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetKnownHosts(obj interface{}) int {

	switch obj.(type) {

	case *[]models.KnownHost:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case []models.KnownHost:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) GetKnownHostsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	resp, err := h.GetKnownHosts(ctx)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.KnownHost{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetKnownHosts(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetKnownHosts(resp))
	w.Write(respBytes)

}

// newGetKnownHostsInput takes in an http.Request an returns the input struct.
func newGetKnownHostsInput(r *http.Request) (*models.GetKnownHostsInput, error) {
	var input models.GetKnownHostsInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	return &input, nil
}

// statusCodeForGetAllLibraryResources returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAllLibraryResources(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.LibraryResources:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.LibraryResources:
		return 200

	default:
		return -1
	}
}

func (h handler) GetAllLibraryResourcesHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newGetAllLibraryResourcesInput(r)
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

	resp, err := h.GetAllLibraryResources(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAllLibraryResources(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAllLibraryResources(resp))
	w.Write(respBytes)

}

// newGetAllLibraryResourcesInput takes in an http.Request an returns the input struct.
func newGetAllLibraryResourcesInput(r *http.Request) (*models.GetAllLibraryResourcesInput, error) {
	var input models.GetAllLibraryResourcesInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	categoryStrs := r.URL.Query()["category"]

	if len(categoryStrs) > 0 {
		var categoryTmp string
		categoryStr := categoryStrs[0]
		categoryTmp, err = categoryStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Category = &categoryTmp
	}

	includeDevAppsStrs := r.URL.Query()["includeDevApps"]

	if len(includeDevAppsStrs) > 0 {
		var includeDevAppsTmp bool
		includeDevAppsStr := includeDevAppsStrs[0]
		includeDevAppsTmp, err = strconv.ParseBool(includeDevAppsStr)
		if err != nil {
			return nil, err
		}
		input.IncludeDevApps = &includeDevAppsTmp
	}

	includeLinksStrs := r.URL.Query()["includeLinks"]

	if len(includeLinksStrs) > 0 {
		var includeLinksTmp bool
		includeLinksStr := includeLinksStrs[0]
		includeLinksTmp, err = strconv.ParseBool(includeLinksStr)
		if err != nil {
			return nil, err
		}
		input.IncludeLinks = &includeLinksTmp
	}

	return &input, nil
}

// statusCodeForSearchLibraryResource returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForSearchLibraryResource(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.LibraryResources:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.LibraryResources:
		return 200

	default:
		return -1
	}
}

func (h handler) SearchLibraryResourceHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newSearchLibraryResourceInput(r)
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

	resp, err := h.SearchLibraryResource(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForSearchLibraryResource(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForSearchLibraryResource(resp))
	w.Write(respBytes)

}

// newSearchLibraryResourceInput takes in an http.Request an returns the input struct.
func newSearchLibraryResourceInput(r *http.Request) (*models.SearchLibraryResourceInput, error) {
	var input models.SearchLibraryResourceInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	searchTermStrs := r.URL.Query()["searchTerm"]
	if len(searchTermStrs) == 0 {
		return nil, errors.New("query parameter 'searchTerm' must be specified")
	}

	if len(searchTermStrs) > 0 {
		var searchTermTmp string
		searchTermStr := searchTermStrs[0]
		searchTermTmp, err = searchTermStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.SearchTerm = searchTermTmp
	}

	showInLibraryOnlyStrs := r.URL.Query()["showInLibraryOnly"]

	if len(showInLibraryOnlyStrs) > 0 {
		var showInLibraryOnlyTmp bool
		showInLibraryOnlyStr := showInLibraryOnlyStrs[0]
		showInLibraryOnlyTmp, err = strconv.ParseBool(showInLibraryOnlyStr)
		if err != nil {
			return nil, err
		}
		input.ShowInLibraryOnly = &showInLibraryOnlyTmp
	}

	includeLinksStrs := r.URL.Query()["includeLinks"]

	if len(includeLinksStrs) > 0 {
		var includeLinksTmp bool
		includeLinksStr := includeLinksStrs[0]
		includeLinksTmp, err = strconv.ParseBool(includeLinksStr)
		if err != nil {
			return nil, err
		}
		input.IncludeLinks = &includeLinksTmp
	}

	return &input, nil
}

// statusCodeForGetLibraryResourceByShortname returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetLibraryResourceByShortname(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.LibraryResource:
		return 200

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.LibraryResource:
		return 200

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetLibraryResourceByShortnameHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newGetLibraryResourceByShortnameInput(r)
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

	resp, err := h.GetLibraryResourceByShortname(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetLibraryResourceByShortname(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetLibraryResourceByShortname(resp))
	w.Write(respBytes)

}

// newGetLibraryResourceByShortnameInput takes in an http.Request an returns the input struct.
func newGetLibraryResourceByShortnameInput(r *http.Request) (*models.GetLibraryResourceByShortnameInput, error) {
	var input models.GetLibraryResourceByShortnameInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	shortnameStr := mux.Vars(r)["shortname"]
	if len(shortnameStr) == 0 {
		return nil, errors.New("path parameter 'shortname' must be specified")
	}
	shortnameStrs := []string{shortnameStr}

	if len(shortnameStrs) > 0 {
		var shortnameTmp string
		shortnameStr := shortnameStrs[0]
		shortnameTmp, err = shortnameStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Shortname = shortnameTmp
	}

	includeDevAppsStrs := r.URL.Query()["includeDevApps"]

	if len(includeDevAppsStrs) > 0 {
		var includeDevAppsTmp bool
		includeDevAppsStr := includeDevAppsStrs[0]
		includeDevAppsTmp, err = strconv.ParseBool(includeDevAppsStr)
		if err != nil {
			return nil, err
		}
		input.IncludeDevApps = &includeDevAppsTmp
	}

	includeLinksStrs := r.URL.Query()["includeLinks"]

	if len(includeLinksStrs) > 0 {
		var includeLinksTmp bool
		includeLinksStr := includeLinksStrs[0]
		includeLinksTmp, err = strconv.ParseBool(includeLinksStr)
		if err != nil {
			return nil, err
		}
		input.IncludeLinks = &includeLinksTmp
	}

	return &input, nil
}

// statusCodeForUpdateLibraryResourceByShortname returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForUpdateLibraryResourceByShortname(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.LibraryResource:
		return 200

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.LibraryResource:
		return 200

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) UpdateLibraryResourceByShortnameHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newUpdateLibraryResourceByShortnameInput(r)
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

	resp, err := h.UpdateLibraryResourceByShortname(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForUpdateLibraryResourceByShortname(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForUpdateLibraryResourceByShortname(resp))
	w.Write(respBytes)

}

// newUpdateLibraryResourceByShortnameInput takes in an http.Request an returns the input struct.
func newUpdateLibraryResourceByShortnameInput(r *http.Request) (*models.UpdateLibraryResourceByShortnameInput, error) {
	var input models.UpdateLibraryResourceByShortnameInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	shortnameStr := mux.Vars(r)["shortname"]
	if len(shortnameStr) == 0 {
		return nil, errors.New("path parameter 'shortname' must be specified")
	}
	shortnameStrs := []string{shortnameStr}

	if len(shortnameStrs) > 0 {
		var shortnameTmp string
		shortnameStr := shortnameStrs[0]
		shortnameTmp, err = shortnameStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Shortname = shortnameTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.LibraryResource = &models.PatchLibraryResourceRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.LibraryResource); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForCreateLibraryResource returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForCreateLibraryResource(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.LibraryResource:
		return 200

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.LibraryResource:
		return 200

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) CreateLibraryResourceHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newCreateLibraryResourceInput(r)
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

	resp, err := h.CreateLibraryResource(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForCreateLibraryResource(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForCreateLibraryResource(resp))
	w.Write(respBytes)

}

// newCreateLibraryResourceInput takes in an http.Request an returns the input struct.
func newCreateLibraryResourceInput(r *http.Request) (*models.CreateLibraryResourceInput, error) {
	var input models.CreateLibraryResourceInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	shortnameStr := mux.Vars(r)["shortname"]
	if len(shortnameStr) == 0 {
		return nil, errors.New("path parameter 'shortname' must be specified")
	}
	shortnameStrs := []string{shortnameStr}

	if len(shortnameStrs) > 0 {
		var shortnameTmp string
		shortnameStr := shortnameStrs[0]
		shortnameTmp, err = shortnameStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Shortname = shortnameTmp
	}

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	sp.LogFields(log.Int("request-size-bytes", len(data)))

	if len(data) > 0 {
		jsonSpan, _ := opentracing.StartSpanFromContext(r.Context(), "json-request-marshaling")
		defer jsonSpan.Finish()

		input.LibraryResource = &models.CreateLibraryResourceRequest{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.LibraryResource); err != nil {
			return nil, err
		}

	}

	return &input, nil
}

// statusCodeForDeleteLibraryResourceLink returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForDeleteLibraryResourceLink(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) DeleteLibraryResourceLinkHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	shortname, err := newDeleteLibraryResourceLinkInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateDeleteLibraryResourceLinkInput(shortname)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = h.DeleteLibraryResourceLink(ctx, shortname)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForDeleteLibraryResourceLink(err)
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

// newDeleteLibraryResourceLinkInput takes in an http.Request an returns the shortname parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newDeleteLibraryResourceLinkInput(r *http.Request) (string, error) {
	shortname := mux.Vars(r)["shortname"]
	if len(shortname) == 0 {
		return "", errors.New("Parameter shortname must be specified")
	}
	return shortname, nil
}

// statusCodeForGetValidPermissions returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetValidPermissions(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.GetValidPermissionsResponse:
		return 200

	case *models.InternalError:
		return 500

	case models.BadRequest:
		return 400

	case models.GetValidPermissionsResponse:
		return 200

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) GetValidPermissionsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	resp, err := h.GetValidPermissions(ctx)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetValidPermissions(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetValidPermissions(resp))
	w.Write(respBytes)

}

// newGetValidPermissionsInput takes in an http.Request an returns the input struct.
func newGetValidPermissionsInput(r *http.Request) (*models.GetValidPermissionsInput, error) {
	var input models.GetValidPermissionsInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	return &input, nil
}

// statusCodeForGetPlatforms returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetPlatforms(obj interface{}) int {

	switch obj.(type) {

	case *[]models.Platform:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case []models.Platform:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) GetPlatformsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	input, err := newGetPlatformsInput(r)
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

	resp, err := h.GetPlatforms(ctx, input)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.Platform{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetPlatforms(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetPlatforms(resp))
	w.Write(respBytes)

}

// newGetPlatformsInput takes in an http.Request an returns the input struct.
func newGetPlatformsInput(r *http.Request) (*models.GetPlatformsInput, error) {
	var input models.GetPlatformsInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err
	if appIds, ok := r.URL.Query()["appIds"]; ok {
		input.AppIds = appIds
	}

	nameStrs := r.URL.Query()["name"]

	if len(nameStrs) > 0 {
		var nameTmp string
		nameStr := nameStrs[0]
		nameTmp, err = nameStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Name = &nameTmp
	}

	return &input, nil
}

// statusCodeForGetPlatformByClientID returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetPlatformByClientID(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case *models.Platform:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	case models.Platform:
		return 200

	default:
		return -1
	}
}

func (h handler) GetPlatformByClientIDHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	clientID, err := newGetPlatformByClientIDInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetPlatformByClientIDInput(clientID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetPlatformByClientID(ctx, clientID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetPlatformByClientID(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetPlatformByClientID(resp))
	w.Write(respBytes)

}

// newGetPlatformByClientIDInput takes in an http.Request an returns the clientID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetPlatformByClientIDInput(r *http.Request) (string, error) {
	clientID := mux.Vars(r)["clientID"]
	if len(clientID) == 0 {
		return "", errors.New("Parameter clientID must be specified")
	}
	return clientID, nil
}

// statusCodeForGetAppsForAdmin returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAppsForAdmin(obj interface{}) int {

	switch obj.(type) {

	case *[]models.AppForAdminResponse:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case []models.AppForAdminResponse:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) GetAppsForAdminHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	sp := opentracing.SpanFromContext(ctx)

	adminID, err := newGetAppsForAdminInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetAppsForAdminInput(adminID)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetAppsForAdmin(ctx, adminID)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.AppForAdminResponse{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAppsForAdmin(err)
		if statusCode == -1 {
			err = models.InternalError{Message: err.Error()}
			statusCode = 500
		}
		http.Error(w, jsonMarshalNoError(err), statusCode)
		return
	}

	jsonSpan, _ := opentracing.StartSpanFromContext(ctx, "json-response-marshaling")
	defer jsonSpan.Finish()

	respBytes, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
		return
	}

	sp.LogFields(log.Int("response-size-bytes", len(respBytes)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAppsForAdmin(resp))
	w.Write(respBytes)

}

// newGetAppsForAdminInput takes in an http.Request an returns the adminID parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetAppsForAdminInput(r *http.Request) (string, error) {
	adminID := mux.Vars(r)["adminID"]
	if len(adminID) == 0 {
		return "", errors.New("Parameter adminID must be specified")
	}
	return adminID, nil
}

// statusCodeForOverrideConfig returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForOverrideConfig(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case *models.NotFound:
		return 404

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	case models.NotFound:
		return 404

	default:
		return -1
	}
}

func (h handler) OverrideConfigHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newOverrideConfigInput(r)
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

	err = h.OverrideConfig(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForOverrideConfig(err)
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

// newOverrideConfigInput takes in an http.Request an returns the input struct.
func newOverrideConfigInput(r *http.Request) (*models.OverrideConfigInput, error) {
	var input models.OverrideConfigInput

	sp := opentracing.SpanFromContext(r.Context())
	_ = sp

	var err error
	_ = err

	srcAppIDStr := mux.Vars(r)["srcAppID"]
	if len(srcAppIDStr) == 0 {
		return nil, errors.New("path parameter 'srcAppID' must be specified")
	}
	srcAppIDStrs := []string{srcAppIDStr}

	if len(srcAppIDStrs) > 0 {
		var srcAppIDTmp string
		srcAppIDStr := srcAppIDStrs[0]
		srcAppIDTmp, err = srcAppIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.SrcAppID = srcAppIDTmp
	}

	destAppIDStr := mux.Vars(r)["destAppID"]
	if len(destAppIDStr) == 0 {
		return nil, errors.New("path parameter 'destAppID' must be specified")
	}
	destAppIDStrs := []string{destAppIDStr}

	if len(destAppIDStrs) > 0 {
		var destAppIDTmp string
		destAppIDStr := destAppIDStrs[0]
		destAppIDTmp, err = destAppIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.DestAppID = destAppIDTmp
	}

	return &input, nil
}
