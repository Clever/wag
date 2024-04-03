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
	"github.com/Clever/wag/samples/gen-go-basic/models/v9"
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

// statusCodeForGetAuthors returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAuthors(obj interface{}) int {

	switch obj.(type) {

	case *models.AuthorsResponse:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case models.AuthorsResponse:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) GetAuthorsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newGetAuthorsInput(r)
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

	resp, nextPageID, err := h.GetAuthors(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAuthors(err)
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

	if !swag.IsZero(nextPageID) {
		input.StartingAfter = &nextPageID
		path, err := input.Path()
		if err != nil {
			logger.FromContext(ctx).AddContext("error", err.Error())
			http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Next-Page-Path", path)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAuthors(resp))
	w.Write(respBytes)

}

// newGetAuthorsInput takes in an http.Request an returns the input struct.
func newGetAuthorsInput(r *http.Request) (*models.GetAuthorsInput, error) {
	var input models.GetAuthorsInput

	var err error
	_ = err

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

	startingAfterStrs := r.URL.Query()["startingAfter"]

	if len(startingAfterStrs) > 0 {
		var startingAfterTmp string
		startingAfterStr := startingAfterStrs[0]
		startingAfterTmp, err = startingAfterStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.StartingAfter = &startingAfterTmp
	}

	return &input, nil
}

// statusCodeForGetAuthorsWithPut returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetAuthorsWithPut(obj interface{}) int {

	switch obj.(type) {

	case *models.AuthorsResponse:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case models.AuthorsResponse:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) GetAuthorsWithPutHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newGetAuthorsWithPutInput(r)
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

	resp, nextPageID, err := h.GetAuthorsWithPut(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetAuthorsWithPut(err)
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

	if !swag.IsZero(nextPageID) {
		input.StartingAfter = &nextPageID
		path, err := input.Path()
		if err != nil {
			logger.FromContext(ctx).AddContext("error", err.Error())
			http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Next-Page-Path", path)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCodeForGetAuthorsWithPut(resp))
	w.Write(respBytes)

}

// newGetAuthorsWithPutInput takes in an http.Request an returns the input struct.
func newGetAuthorsWithPutInput(r *http.Request) (*models.GetAuthorsWithPutInput, error) {
	var input models.GetAuthorsWithPutInput

	var err error
	_ = err

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

	startingAfterStrs := r.URL.Query()["startingAfter"]

	if len(startingAfterStrs) > 0 {
		var startingAfterTmp string
		startingAfterStr := startingAfterStrs[0]
		startingAfterTmp, err = startingAfterStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.StartingAfter = &startingAfterTmp
	}

	data, err := ioutil.ReadAll(r.Body)

	if len(data) > 0 {
		input.FavoriteBooks = &models.Book{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.FavoriteBooks); err != nil {
			return nil, err
		}
	}

	return &input, nil
}

// statusCodeForGetBooks returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetBooks(obj interface{}) int {

	switch obj.(type) {

	case *[]models.Book:
		return 200

	case *models.BadRequest:
		return 400

	case *models.InternalError:
		return 500

	case []models.Book:
		return 200

	case models.BadRequest:
		return 400

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) GetBooksHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newGetBooksInput(r)
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

	resp, nextPageID, err := h.GetBooks(ctx, input)

	// Success types that return an array should never return nil so let's make this easier
	// for consumers by converting nil arrays to empty arrays
	if resp == nil {
		resp = []models.Book{}
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetBooks(err)
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

	if !swag.IsZero(nextPageID) {
		input.StartingAfter = &nextPageID
		path, err := input.Path()
		if err != nil {
			logger.FromContext(ctx).AddContext("error", err.Error())
			http.Error(w, jsonMarshalNoError(models.InternalError{Message: err.Error()}), http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Next-Page-Path", path)
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
	if authors, ok := r.URL.Query()["authors"]; ok {
		input.Authors = authors
	}

	availableStrs := r.URL.Query()["available"]

	if len(availableStrs) == 0 {
		availableStrs = []string{"true"}
	}
	if len(availableStrs) > 0 {
		var availableTmp bool
		availableStr := availableStrs[0]
		availableTmp, err = strconv.ParseBool(availableStr)
		if err != nil {
			return nil, err
		}
		input.Available = &availableTmp
	}

	stateStrs := r.URL.Query()["state"]

	if len(stateStrs) == 0 {
		stateStrs = []string{"finished"}
	}
	if len(stateStrs) > 0 {
		var stateTmp string
		stateStr := stateStrs[0]
		stateTmp, err = stateStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.State = &stateTmp
	}

	publishedStrs := r.URL.Query()["published"]

	if len(publishedStrs) > 0 {
		var publishedTmp strfmt.Date
		publishedStr := publishedStrs[0]
		publishedTmp, err = convertDate(publishedStr)
		if err != nil {
			return nil, err
		}
		input.Published = &publishedTmp
	}

	snakeCaseStrs := r.URL.Query()["snake_case"]

	if len(snakeCaseStrs) > 0 {
		var snakeCaseTmp string
		snakeCaseStr := snakeCaseStrs[0]
		snakeCaseTmp, err = snakeCaseStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.SnakeCase = &snakeCaseTmp
	}

	completedStrs := r.URL.Query()["completed"]

	if len(completedStrs) > 0 {
		var completedTmp strfmt.DateTime
		completedStr := completedStrs[0]
		completedTmp, err = convertDateTime(completedStr)
		if err != nil {
			return nil, err
		}
		input.Completed = &completedTmp
	}

	maxPagesStrs := r.URL.Query()["maxPages"]

	if len(maxPagesStrs) == 0 {
		maxPagesStrs = []string{"5.005E+02"}
	}
	if len(maxPagesStrs) > 0 {
		var maxPagesTmp float64
		maxPagesStr := maxPagesStrs[0]
		maxPagesTmp, err = swag.ConvertFloat64(maxPagesStr)
		if err != nil {
			return nil, err
		}
		input.MaxPages = &maxPagesTmp
	}

	minPagesStrs := r.URL.Query()["min_pages"]

	if len(minPagesStrs) == 0 {
		minPagesStrs = []string{"5"}
	}
	if len(minPagesStrs) > 0 {
		var minPagesTmp int32
		minPagesStr := minPagesStrs[0]
		minPagesTmp, err = swag.ConvertInt32(minPagesStr)
		if err != nil {
			return nil, err
		}
		input.MinPages = &minPagesTmp
	}

	pagesToTimeStrs := r.URL.Query()["pagesToTime"]

	if len(pagesToTimeStrs) > 0 {
		var pagesToTimeTmp float32
		pagesToTimeStr := pagesToTimeStrs[0]
		pagesToTimeTmp, err = swag.ConvertFloat32(pagesToTimeStr)
		if err != nil {
			return nil, err
		}
		input.PagesToTime = &pagesToTimeTmp
	}

	authorizationStrs := r.Header.Get("authorization")

	if len(authorizationStrs) > 0 {
		var authorizationTmp string
		authorizationTmp = authorizationStrs
		input.Authorization = authorizationTmp
	}

	startingAfterStrs := r.URL.Query()["startingAfter"]

	if len(startingAfterStrs) > 0 {
		var startingAfterTmp int64
		startingAfterStr := startingAfterStrs[0]
		startingAfterTmp, err = swag.ConvertInt64(startingAfterStr)
		if err != nil {
			return nil, err
		}
		input.StartingAfter = &startingAfterTmp
	}

	return &input, nil
}

// statusCodeForCreateBook returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForCreateBook(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Book:
		return 200

	case *models.InternalError:
		return 500

	case models.BadRequest:
		return 400

	case models.Book:
		return 200

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) CreateBookHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newCreateBookInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	if input != nil {
		err = input.Validate(nil)
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.CreateBook(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForCreateBook(err)
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
	w.WriteHeader(statusCodeForCreateBook(resp))
	w.Write(respBytes)

}

// newCreateBookInput takes in an http.Request an returns the input struct.
func newCreateBookInput(r *http.Request) (*models.Book, error) {
	var err error
	_ = err

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	if len(data) > 0 {
		var input models.Book
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(&input); err != nil {
			return nil, err
		}
		return &input, nil
	}

	return nil, nil
}

// statusCodeForPutBook returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForPutBook(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Book:
		return 200

	case *models.InternalError:
		return 500

	case models.BadRequest:
		return 400

	case models.Book:
		return 200

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) PutBookHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newPutBookInput(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	if input != nil {
		err = input.Validate(nil)
	}

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.PutBook(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForPutBook(err)
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
	w.WriteHeader(statusCodeForPutBook(resp))
	w.Write(respBytes)

}

// newPutBookInput takes in an http.Request an returns the input struct.
func newPutBookInput(r *http.Request) (*models.Book, error) {
	var err error
	_ = err

	data, err := ioutil.ReadAll(r.Body)

	if len(data) > 0 {
		var input models.Book
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(&input); err != nil {
			return nil, err
		}
		return &input, nil
	}

	return nil, nil
}

// statusCodeForGetBookByID returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForGetBookByID(obj interface{}) int {

	switch obj.(type) {

	case *models.BadRequest:
		return 400

	case *models.Book:
		return 200

	case *models.Error:
		return 404

	case *models.InternalError:
		return 500

	case *models.Unathorized:
		return 401

	case models.BadRequest:
		return 400

	case models.Book:
		return 200

	case models.Error:
		return 404

	case models.InternalError:
		return 500

	case models.Unathorized:
		return 401

	default:
		return -1
	}
}

func (h handler) GetBookByIDHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newGetBookByIDInput(r)
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

	resp, err := h.GetBookByID(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetBookByID(err)
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
		return nil, errors.New("path parameter 'book_id' must be specified")
	}
	bookIDStrs := []string{bookIDStr}

	if len(bookIDStrs) > 0 {
		var bookIDTmp int64
		bookIDStr := bookIDStrs[0]
		bookIDTmp, err = swag.ConvertInt64(bookIDStr)
		if err != nil {
			return nil, err
		}
		input.BookID = bookIDTmp
	}

	authorIDStrs := r.URL.Query()["authorID"]

	if len(authorIDStrs) > 0 {
		var authorIDTmp string
		authorIDStr := authorIDStrs[0]
		authorIDTmp, err = authorIDStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.AuthorID = &authorIDTmp
	}

	authorizationStrs := r.Header.Get("authorization")

	if len(authorizationStrs) > 0 {
		var authorizationTmp string
		authorizationTmp = authorizationStrs
		input.Authorization = authorizationTmp
	}

	xDontRateLimitMeBroStrs := r.Header.Get("X-Dont-Rate-Limit-Me-Bro")

	if len(xDontRateLimitMeBroStrs) > 0 {
		var xDontRateLimitMeBroTmp string
		xDontRateLimitMeBroTmp = xDontRateLimitMeBroStrs
		input.XDontRateLimitMeBro = xDontRateLimitMeBroTmp
	}

	randomBytesStrs := r.URL.Query()["randomBytes"]

	if len(randomBytesStrs) > 0 {
		var randomBytesTmp strfmt.Base64
		randomBytesStr := randomBytesStrs[0]
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

	case *models.BadRequest:
		return 400

	case *models.Book:
		return 200

	case *models.Error:
		return 404

	case *models.InternalError:
		return 500

	case models.BadRequest:
		return 400

	case models.Book:
		return 200

	case models.Error:
		return 404

	case models.InternalError:
		return 500

	default:
		return -1
	}
}

func (h handler) GetBookByID2Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	id, err := newGetBookByID2Input(r)
	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	err = models.ValidateGetBookByID2Input(id)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		http.Error(w, jsonMarshalNoError(models.BadRequest{Message: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetBookByID2(ctx, id)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForGetBookByID2(err)
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
	w.WriteHeader(statusCodeForGetBookByID2(resp))
	w.Write(respBytes)

}

// newGetBookByID2Input takes in an http.Request an returns the id parameter
// that it contains. It returns an error if the request doesn't contain the parameter.
func newGetBookByID2Input(r *http.Request) (string, error) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		return "", errors.New("Parameter id must be specified")
	}
	return id, nil
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

	var err error
	_ = err

	return &input, nil
}

// statusCodeForLowercaseModelsTest returns the status code corresponding to the returned
// object. It returns -1 if the type doesn't correspond to anything.
func statusCodeForLowercaseModelsTest(obj interface{}) int {

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

func (h handler) LowercaseModelsTestHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := newLowercaseModelsTestInput(r)
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

	err = h.LowercaseModelsTest(ctx, input)

	if err != nil {
		logger.FromContext(ctx).AddContext("error", err.Error())
		if btErr, ok := err.(*errors.Error); ok {
			logger.FromContext(ctx).AddContext("stacktrace", string(btErr.Stack()))
		} else if xerr, ok := err.(xerrors.Formatter); ok {
			logger.FromContext(ctx).AddContext("frames", fmt.Sprintf("%+v", xerr))
		}
		statusCode := statusCodeForLowercaseModelsTest(err)
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

// newLowercaseModelsTestInput takes in an http.Request an returns the input struct.
func newLowercaseModelsTestInput(r *http.Request) (*models.LowercaseModelsTestInput, error) {
	var input models.LowercaseModelsTestInput

	var err error
	_ = err

	data, err := ioutil.ReadAll(r.Body)
	if len(data) == 0 {
		return nil, errors.New("request body is required, but was empty")
	}
	if len(data) > 0 {
		input.Lowercase = &models.lowercase{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.Lowercase); err != nil {
			return nil, err
		}
	}

	pathParamStr := mux.Vars(r)["pathParam"]
	if len(pathParamStr) == 0 {
		return nil, errors.New("path parameter 'pathParam' must be specified")
	}
	pathParamStrs := []string{pathParamStr}

	if len(pathParamStrs) > 0 {
		var pathParamTmp string
		pathParamStr := pathParamStrs[0]
		pathParamTmp, err = pathParamStr, error(nil)
		if err != nil {
			return nil, err
		}
		input.PathParam = pathParamTmp
	}

	return &input, nil
}
