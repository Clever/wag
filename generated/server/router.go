package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	gContext "github.com/gorilla/context"
	"golang.org/x/net/context"

	"gopkg.in/tylerb/graceful.v1"
)

type contextKey struct{}

type Server struct {
	Handler http.Handler
	port    int
}

func (s Server) Serve() error {
	// Give the sever 30 seconds to shut down
	graceful.Run(":"+string(s.port), 30*time.Second, s.Handler)

	// This should never return
	return errors.New("This should never happen")
}

func New(c Controller, port int) Server {
	controller = c // TODO: get rid of global variable?
	r := mux.NewRouter()

	r.Methods("GET").Path("/v1/books").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gContext.Get(r, contextKey{}).(context.Context)
		GetBooksHandler(ctx, w, r)
	})

	r.Methods("GET").Path("/v1/books/{bookID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gContext.Get(r, contextKey{}).(context.Context)
		GetBookByIDHandler(ctx, w, r)
	})

	r.Methods("POST").Path("/v1/books/{bookID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gContext.Get(r, contextKey{}).(context.Context)
		CreateBookHandler(ctx, w, r)
	})
	handler := withMiddleware(r)
	return Server{Handler: handler, port: port}
}
