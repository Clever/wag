package main

import (
	"net/http"

	"github.com/gorilla/mux"

	gContext "github.com/gorilla/context"
	"golang.org/x/net/context"
)

type contextKey struct{}

func withRoutes(r *mux.Router) *mux.Router {
	r.Methods("get").Path("/books/{bookID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gContext.Get(r, contextKey{}).(context.Context)
		GetBookByIDHandler(ctx, w, r)
	})
	return r
}
