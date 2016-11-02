package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Clever/wag/samples/gen-go/models"
	"github.com/go-errors/errors"
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

// statusCodeForGetBooks returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetBooks(obj interface{}) int {

	switch obj.(type) {

	case []models.Book:
		return 200

	case models.DefaultBadRequest:
		return 400
	case models.DefaultInternalError:
		return 500
	default:
		return -1
	}
}

func (h handler) GetBooksHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newGetBooksInput(r)
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

	resp, err := h.GetBooks(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		}
		statusCode := statusCodeForGetBooks(err)
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
	w.WriteHeader(statusCodeForGetBooks(resp))
	w.Write(respBytes)

}

// newGetBooksInput takes in an http.Request an returns the input struct.
func newGetBooksInput(r *http.Request) (*models.GetBooksInput, error) {
	var input models.GetBooksInput

	var err error
	_ = err

	authorsStr := r.URL.Query().Get("authors")
	if len(authorsStr) != 0 {
		var authorsTmp []string
		authorsTmp, err = swag.SplitByFormat(authorsStr, ""), error(nil)
		if err != nil {
			return nil, err
		}
		input.Authors = authorsTmp

	}
	availableStr := r.URL.Query().Get("available")
	if len(availableStr) == 0 {
		// Use the default value
		availableStr = "true"
	}
	if len(availableStr) != 0 {
		var availableTmp bool
		availableTmp, err = strconv.ParseBool(availableStr)
		if err != nil {
			return nil, err
		}
		input.Available = &availableTmp

	}
	stateStr := r.URL.Query().Get("state")
	if len(stateStr) == 0 {
		// Use the default value
		stateStr = "finished"
	}
	if len(stateStr) != 0 {
		var stateTmp string
		stateTmp, err = stateStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.State = &stateTmp

	}
	publishedStr := r.URL.Query().Get("published")
	if len(publishedStr) != 0 {
		var publishedTmp strfmt.Date
		publishedTmp, err = convertDate(publishedStr)
		if err != nil {
			return nil, err
		}
		input.Published = &publishedTmp

	}
	snakeCaseStr := r.URL.Query().Get("snake_case")
	if len(snakeCaseStr) != 0 {
		var snakeCaseTmp string
		snakeCaseTmp, err = snakeCaseStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.SnakeCase = &snakeCaseTmp

	}
	completedStr := r.URL.Query().Get("completed")
	if len(completedStr) != 0 {
		var completedTmp strfmt.DateTime
		completedTmp, err = convertDateTime(completedStr)
		if err != nil {
			return nil, err
		}
		input.Completed = &completedTmp

	}
	maxPagesStr := r.URL.Query().Get("maxPages")
	if len(maxPagesStr) == 0 {
		// Use the default value
		maxPagesStr = "5.005E+02"
	}
	if len(maxPagesStr) != 0 {
		var maxPagesTmp float64
		maxPagesTmp, err = swag.ConvertFloat64(maxPagesStr)
		if err != nil {
			return nil, err
		}
		input.MaxPages = &maxPagesTmp

	}
	minPagesStr := r.URL.Query().Get("min_pages")
	if len(minPagesStr) == 0 {
		// Use the default value
		minPagesStr = "5"
	}
	if len(minPagesStr) != 0 {
		var minPagesTmp int32
		minPagesTmp, err = swag.ConvertInt32(minPagesStr)
		if err != nil {
			return nil, err
		}
		input.MinPages = &minPagesTmp

	}
	pagesToTimeStr := r.URL.Query().Get("pagesToTime")
	if len(pagesToTimeStr) != 0 {
		var pagesToTimeTmp float32
		pagesToTimeTmp, err = swag.ConvertFloat32(pagesToTimeStr)
		if err != nil {
			return nil, err
		}
		input.PagesToTime = &pagesToTimeTmp

	}

	return &input, nil
}

// statusCodeForCreateBook returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForCreateBook(obj interface{}) int {

	switch obj.(type) {

	case *models.Book:
		return 200

	case models.Book:
		return 200

	case models.DefaultBadRequest:
		return 400
	case models.DefaultInternalError:
		return 500
	default:
		return -1
	}
}

func (h handler) CreateBookHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newCreateBookInput(r)
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

	resp, err := h.CreateBook(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		}
		statusCode := statusCodeForCreateBook(err)
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
	w.WriteHeader(statusCodeForCreateBook(resp))
	w.Write(respBytes)

}

// newCreateBookInput takes in an http.Request an returns the input struct.
func newCreateBookInput(r *http.Request) (*models.Book, error) {
	var input models.Book

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

// statusCodeForGetBookByID returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetBookByID(obj interface{}) int {

	switch obj.(type) {

	case *models.GetBookByID200Output:
		return 200

	case *models.GetBookByID204Output:
		return 204

	case *models.GetBookByID401Output:
		return 401

	case *models.GetBookByID404Output:
		return 404

	case models.GetBookByID200Output:
		return 200

	case models.GetBookByID204Output:
		return 204

	case models.GetBookByID401Output:
		return 401

	case models.GetBookByID404Output:
		return 404

	case models.DefaultBadRequest:
		return 400
	case models.DefaultInternalError:
		return 500
	default:
		return -1
	}
}

func (h handler) GetBookByIDHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newGetBookByIDInput(r)
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

	resp, err := h.GetBookByID(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		}
		statusCode := statusCodeForGetBookByID(err)
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
	w.WriteHeader(statusCodeForGetBookByID(resp))
	w.Write(respBytes)

}

// newGetBookByIDInput takes in an http.Request an returns the input struct.
func newGetBookByIDInput(r *http.Request) (*models.GetBookByIDInput, error) {
	var input models.GetBookByIDInput

	var err error
	_ = err

	bookIDStr := mux.Vars(r)["book_id"]
	if len(bookIDStr) == 0 {
		return nil, errors.New("Parameter must be specified")
	}
	if len(bookIDStr) != 0 {
		var bookIDTmp int64
		bookIDTmp, err = swag.ConvertInt64(bookIDStr)
		if err != nil {
			return nil, err
		}
		input.BookID = bookIDTmp

	}
	authorIDStr := r.URL.Query().Get("authorID")
	if len(authorIDStr) != 0 {
		var authorIDTmp string
		authorIDTmp, err = authorIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AuthorID = &authorIDTmp

	}
	authorizationStr := r.Header.Get("authorization")
	if len(authorizationStr) != 0 {
		var authorizationTmp string
		authorizationTmp, err = authorizationStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.Authorization = &authorizationTmp

	}
	randomBytesStr := r.URL.Query().Get("randomBytes")
	if len(randomBytesStr) != 0 {
		var randomBytesTmp strfmt.Base64
		randomBytesTmp, err = convertBase64(randomBytesStr)
		if err != nil {
			return nil, err
		}
		input.RandomBytes = &randomBytesTmp

	}

	return &input, nil
}

// statusCodeForGetBookByID2 returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetBookByID2(obj interface{}) int {

	switch obj.(type) {

	case *models.Book:
		return 200

	case *models.GetBookByID2404Output:
		return 404

	case models.Book:
		return 200

	case models.GetBookByID2404Output:
		return 404

	case models.DefaultBadRequest:
		return 400
	case models.DefaultInternalError:
		return 500
	default:
		return -1
	}
}

func (h handler) GetBookByID2Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newGetBookByID2Input(r)
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

	resp, err := h.GetBookByID2(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		}
		statusCode := statusCodeForGetBookByID2(err)
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
	w.WriteHeader(statusCodeForGetBookByID2(resp))
	w.Write(respBytes)

}

// newGetBookByID2Input takes in an http.Request an returns the input struct.
func newGetBookByID2Input(r *http.Request) (*models.GetBookByID2Input, error) {
	var input models.GetBookByID2Input

	var err error
	_ = err

	iDStr := mux.Vars(r)["id"]
	if len(iDStr) == 0 {
		return nil, errors.New("Parameter must be specified")
	}
	if len(iDStr) != 0 {
		var iDTmp string
		iDTmp, err = iDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.ID = iDTmp

	}

	return &input, nil
}

// statusCodeForHealthCheck returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForHealthCheck(obj interface{}) int {

	switch obj.(type) {

	case models.DefaultBadRequest:
		return 400
	case models.DefaultInternalError:
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
		}
		statusCode := statusCodeForHealthCheck(err)
		if statusCode != -1 {
			http.Error(w, err.Error(), statusCode)
		} else {
			http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(""))

}

// newHealthCheckInput takes in an http.Request an returns the input struct.
func newHealthCheckInput(r *http.Request) (*models.HealthCheckInput, error) {
	var input models.HealthCheckInput

	var err error
	_ = err

	return &input, nil
}
