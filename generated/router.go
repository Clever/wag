package main

import (
	"net/http"

	"github.com/gorilla/mux"

	gContext "github.com/gorilla/context"
	"golang.org/x/net/context"
)

// TODO: Is this the way to do
type contextKey struct{}

func withRoutes(r *mux.Router) *mux.Router {
	r.Methods("GET").Path("/books/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gContext.Get(r, contextKey{}).(context.Context)
		GetBookByIDHandler(ctx, w, r)
	})
	return r
}
