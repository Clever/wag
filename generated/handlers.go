package main

import (
	"net/http"
	"golang.org/x/net/context"
	"encoding/json"
)

var controller Controller

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
