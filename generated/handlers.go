package generated

import (
	"net/http"
	"golang.org/x/net/context"
	"encoding/json"
)

var controller Controller

func GetBooksHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	input, err := NewGetBooksInput(r)
	if err != nil {
		// TODO: Think about this whether this is usually an internal error or it could
		// be from a bad request format...
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := controller.GetBooks(ctx, input)
	if err != nil {
		if respErr, ok := err.(GetBooksError); ok {
			http.Error(w, respErr.Error(), respErr.GetBooksStatusCode())
			return
		} else {
			// This is the default case
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	respBytes, err := json.Marshal(resp.GetBooksData())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
func GetBookByIDHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	input, err := NewGetBookByIDInput(r)
	if err != nil {
		// TODO: Think about this whether this is usually an internal error or it could
		// be from a bad request format...
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := controller.GetBookByID(ctx, input)
	if err != nil {
		if respErr, ok := err.(GetBookByIDError); ok {
			http.Error(w, respErr.Error(), respErr.GetBookByIDStatusCode())
			return
		} else {
			// This is the default case
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	respBytes, err := json.Marshal(resp.GetBookByIDData())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
<<<<<<< a0005dbc4f1c70496590515a9cdd8b111d9685bb
func CreateBookHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	input, err := NewCreateBookInput(r)
=======
func GetBookByIDHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	input, err := NewGetBookByIDInput(r)
>>>>>>> Add other input types to code generation
	if err != nil {
		// TODO: Think about this whether this is usually an internal error or it could
		// be from a bad request format...
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = input.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

<<<<<<< a0005dbc4f1c70496590515a9cdd8b111d9685bb
	resp, err := controller.CreateBook(ctx, input)
	if err != nil {
		if respErr, ok := err.(CreateBookError); ok {
			http.Error(w, respErr.Error(), respErr.CreateBookStatusCode())
=======
	resp, err := controller.GetBookByID(ctx, input)
	if err != nil {
		if respErr, ok := err.(GetBookByIDError); ok {
			http.Error(w, respErr.Error(), respErr.GetBookByIDStatusCode())
>>>>>>> Add other input types to code generation
			return
		} else {
			// This is the default case
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

<<<<<<< a0005dbc4f1c70496590515a9cdd8b111d9685bb
	respBytes, err := json.Marshal(resp.CreateBookData())
=======
	respBytes, err := json.Marshal(resp.GetBookByIDData())
>>>>>>> Add other input types to code generation
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
