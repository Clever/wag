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

func ConvertBase64(input string) (strfmt.Base64, error) {
	temp, err := formats.Parse("byte", input)
	if err != nil {
		return strfmt.Base64{}, err
	}
	return *temp.(*strfmt.Base64), nil
}

func ConvertDateTime(input string) (strfmt.DateTime, error) {
	temp, err := formats.Parse("date-time", input)
	if err != nil {
		return strfmt.DateTime{}, err
	}
	return *temp.(*strfmt.DateTime), nil
}

func ConvertDate(input string) (strfmt.Date, error) {
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

	input, err := NewGetBooksInput(r)
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
		} else {
			http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(200)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)

}
func NewGetBooksInput(r *http.Request) (*models.GetBooksInput, error) {
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
		publishedTmp, err = ConvertDate(publishedStr)
		if err != nil {
			return nil, err
		}
		input.Published = &publishedTmp

	}
	completedStr := r.URL.Query().Get("completed")
	if len(completedStr) != 0 {
		var completedTmp strfmt.DateTime
		completedTmp, err = ConvertDateTime(completedStr)
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

	input, err := NewGetBookByIDInput(r)
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
		} else {
			http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(resp.GetBookByIDStatus())

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)

}
func NewGetBookByIDInput(r *http.Request) (*models.GetBookByIDInput, error) {
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
		randomBytesTmp, err = ConvertBase64(randomBytesStr)
		if err != nil {
			return nil, err
		}
		input.RandomBytes = &randomBytesTmp

	}

	return &input, nil
}

func (h handler) CreateBookHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	input, err := NewCreateBookInput(r)
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
		} else {
			http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(200)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)

}
func NewCreateBookInput(r *http.Request) (*models.CreateBookInput, error) {
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

	err := h.HealthCheck(ctx)

	if err != nil {
		if respErr, ok := err.(models.HealthCheckError); ok {
			http.Error(w, respErr.Error(), respErr.HealthCheckStatusCode())
			return
		} else {
			http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(200)

	w.Write([]byte(""))

}
func NewHealthCheckInput(r *http.Request) (*models.HealthCheckInput, error) {
	var input models.HealthCheckInput

	var err error
	_ = err

	return &input, nil
}
