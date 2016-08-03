package server

import (
	"net/http"

	"github.com/gorilla/mux"

	gContext "github.com/gorilla/context"
	"golang.org/x/net/context"
)

type contextKey struct{}

func SetupServer(r *mux.Router, c Controller) http.Handler {
	controller = c // TODO: get rid of global variable?
	r.Methods("get").Path("/v1/books/{bookID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gContext.Get(r, contextKey{}).(context.Context)
		GetBookByIDHandler(ctx, w, r)
	})

	r.Methods("post").Path("/v1/books/{bookID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gContext.Get(r, contextKey{}).(context.Context)
		CreateBookHandler(ctx, w, r)
	})

	r.Methods("get").Path("/v1/books").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gContext.Get(r, contextKey{}).(context.Context)
		GetBooksHandler(ctx, w, r)
	})
	return withMiddleware(r)
}
