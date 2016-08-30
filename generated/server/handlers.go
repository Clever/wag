package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/Clever/wag/generated/models"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
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
func (h handler) GetBooksHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
<<<<<<< 60d0f23f90d41093fff0b1b68ecf1dfd124d5d2e

	input, err := NewGetBooksInput(r)
=======
	input, err := newGetBooksInput(r)
>>>>>>> Server comments and linting
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetBooks(ctx, input)

	if err != nil {
		if respErr, ok := err.(models.GetBooksError); ok {
			http.Error(w, respErr.Error(), respErr.GetBooksStatusCode())
			return
		}
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respBytes)

}

// newGetBooksInput takes in an http.Request an returns the input struct
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

func (h handler) GetBookByIDHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
<<<<<<< 60d0f23f90d41093fff0b1b68ecf1dfd124d5d2e

	input, err := NewGetBookByIDInput(r)
=======
	input, err := newGetBookByIDInput(r)
>>>>>>> Server comments and linting
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.GetBookByID(ctx, input)

	if err != nil {
		if respErr, ok := err.(models.GetBookByIDError); ok {
			http.Error(w, respErr.Error(), respErr.GetBookByIDStatusCode())
			return
		}
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.GetBookByIDStatus())
	w.Write(respBytes)

}

// newGetBookByIDInput takes in an http.Request an returns the input struct
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

func (h handler) CreateBookHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
<<<<<<< 60d0f23f90d41093fff0b1b68ecf1dfd124d5d2e

	input, err := NewCreateBookInput(r)
=======
	input, err := newCreateBookInput(r)
>>>>>>> Server comments and linting
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	resp, err := h.CreateBook(ctx, input)

	if err != nil {
		if respErr, ok := err.(models.CreateBookError); ok {
			http.Error(w, respErr.Error(), respErr.CreateBookStatusCode())
			return
		}
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respBytes)

}

// newCreateBookInput takes in an http.Request an returns the input struct
func newCreateBookInput(r *http.Request) (*models.CreateBookInput, error) {
	var input models.CreateBookInput

	var err error
	_ = err

	data, err := ioutil.ReadAll(r.Body)
	if len(data) > 0 {
		input.NewBook = &models.Book{}
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(input.NewBook); err != nil {
			return nil, err
		}
	}

	return &input, nil
}

func (h handler) HealthCheckHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
<<<<<<< 60d0f23f90d41093fff0b1b68ecf1dfd124d5d2e
=======
	input, err := newHealthCheckInput(r)
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultBadRequest{Msg: err.Error()}), http.StatusBadRequest)
		return
	}
>>>>>>> Server comments and linting

	err := h.HealthCheck(ctx)

	if err != nil {
		if respErr, ok := err.(models.HealthCheckError); ok {
			http.Error(w, respErr.Error(), respErr.HealthCheckStatusCode())
			return
		}
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(""))

}

// newHealthCheckInput takes in an http.Request an returns the input struct
func newHealthCheckInput(r *http.Request) (*models.HealthCheckInput, error) {
	var input models.HealthCheckInput

	var err error
	_ = err

	return &input, nil
}
