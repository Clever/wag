package server

import (
	"encoding/json"
	"errors"
	"github.com/Clever/wag/generated/models"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

var _ = strconv.ParseInt

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

	authorStr := r.URL.Query().Get("author")
	authorTmp := authorStr
	input.Author = &authorTmp

	availableStr := r.URL.Query().Get("available")
	availableTmp, err := strconv.ParseBool(availableStr)
	input.Available = &availableTmp

	stateStr := r.URL.Query().Get("state")
	stateTmp := stateStr
	input.State = &stateTmp

	maxPagesStr := r.URL.Query().Get("maxPages")
	maxPagesTmp, err := strconv.ParseFloat(maxPagesStr, 64)
	input.MaxPages = &maxPagesTmp

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

	bookIDStr := mux.Vars(r)["bookID"]
	if len(bookIDStr) == 0 {
		return nil, errors.New("Parameter must be specified")
	}
	bookIDTmp, err := strconv.ParseInt(bookIDStr, 10, 64)
	input.BookID = bookIDTmp

	if err != nil {
		return nil, err
	}
	authorizationStr := r.Header.Get("authorization")
	authorizationTmp := authorizationStr
	input.Authorization = &authorizationTmp

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

	err = json.NewDecoder(r.Body).Decode(input.NewBook)

	return &input, nil
}
