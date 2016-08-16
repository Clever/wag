package server

import (
	"encoding/json"
	"errors"
	"github.com/Clever/wag/generated/models"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

var _ = strconv.ParseInt
var _ = strfmt.Default
var _ = swag.ConvertInt32

var controller Controller

func jsonMarshalNoError(i interface{}) string {
	bytes, err := json.Marshal(i)
	if err != nil {
		// This should never happen
		return ""
	}
	return string(bytes)
}
func GetBooksHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	resp, err := controller.GetBooks(ctx, input)
	if err != nil {
		if respErr, ok := err.(models.GetBooksError); ok {
			http.Error(w, respErr.Error(), respErr.GetBooksStatusCode())
			return
		} else {
			http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
			return
		}
	}

	respBytes, err := json.Marshal(resp.GetBooksData())
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

	formats := strfmt.Default
	_ = formats

	authorStr := r.URL.Query().Get("author")
	authorTmp := authorStr
	input.Author = &authorTmp

	if err != nil && len(authorStr) != 0 {
		return nil, err
	}
	availableStr := r.URL.Query().Get("available")
	availableTmp, err := strconv.ParseBool(availableStr)
	input.Available = &availableTmp

	if err != nil && len(availableStr) != 0 {
		return nil, err
	}
	stateStr := r.URL.Query().Get("state")
	stateTmp := stateStr
	input.State = &stateTmp

	if err != nil && len(stateStr) != 0 {
		return nil, err
	}
	publishedStr := r.URL.Query().Get("published")
	publishedTmpInterface, err := formats.Parse("date", publishedStr)
	publishedTmp := publishedTmpInterface.(strfmt.Date)
	input.Published = &publishedTmp

	if err != nil && len(publishedStr) != 0 {
		return nil, err
	}
	completedStr := r.URL.Query().Get("completed")
	completedTmpInterface, err := formats.Parse("date-time", completedStr)
	completedTmp := completedTmpInterface.(strfmt.DateTime)
	input.Completed = &completedTmp

	if err != nil && len(completedStr) != 0 {
		return nil, err
	}
	maxPagesStr := r.URL.Query().Get("maxPages")
	maxPagesTmp, err := swag.ConvertFloat64(maxPagesStr)
	input.MaxPages = &maxPagesTmp

	if err != nil && len(maxPagesStr) != 0 {
		return nil, err
	}
	minPagesStr := r.URL.Query().Get("minPages")
	minPagesTmp, err := swag.ConvertInt32(minPagesStr)
	input.MinPages = &minPagesTmp

	if err != nil && len(minPagesStr) != 0 {
		return nil, err
	}
	pagesToTimeStr := r.URL.Query().Get("pagesToTime")
	pagesToTimeTmp, err := swag.ConvertFloat32(pagesToTimeStr)
	input.PagesToTime = &pagesToTimeTmp

	if err != nil && len(pagesToTimeStr) != 0 {
		return nil, err
	}

	return &input, nil
}

func GetBookByIDHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	resp, err := controller.GetBookByID(ctx, input)
	if err != nil {
		if respErr, ok := err.(models.GetBookByIDError); ok {
			http.Error(w, respErr.Error(), respErr.GetBookByIDStatusCode())
			return
		} else {
			http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
			return
		}
	}

	respBytes, err := json.Marshal(resp.GetBookByIDData())
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

	formats := strfmt.Default
	_ = formats

	bookIDStr := mux.Vars(r)["bookID"]
	if len(bookIDStr) == 0 {
		return nil, errors.New("Parameter must be specified")
	}
	bookIDTmp, err := swag.ConvertInt64(bookIDStr)
	input.BookID = bookIDTmp

	if err != nil {
		return nil, err
	}
	authorizationStr := r.Header.Get("authorization")
	authorizationTmp := authorizationStr
	input.Authorization = &authorizationTmp

	if err != nil && len(authorizationStr) != 0 {
		return nil, err
	}
	randomBytesStr := r.URL.Query().Get("randomBytes")
	randomBytesTmpInterface, err := formats.Parse("byte", randomBytesStr)
	randomBytesTmp := randomBytesTmpInterface.([]byte)
	input.RandomBytes = &randomBytesTmp

	if err != nil && len(randomBytesStr) != 0 {
		return nil, err
	}

	return &input, nil
}

func CreateBookHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	resp, err := controller.CreateBook(ctx, input)
	if err != nil {
		if respErr, ok := err.(models.CreateBookError); ok {
			http.Error(w, respErr.Error(), respErr.CreateBookStatusCode())
			return
		} else {
			http.Error(w, jsonMarshalNoError(models.DefaultInternalError{Msg: err.Error()}), http.StatusInternalServerError)
			return
		}
	}

	respBytes, err := json.Marshal(resp.CreateBookData())
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

	formats := strfmt.Default
	_ = formats

	err = json.NewDecoder(r.Body).Decode(input.NewBook)
	if err != nil {
		return nil, err
	}

	return &input, nil
}
