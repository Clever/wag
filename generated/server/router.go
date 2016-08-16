package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"gopkg.in/tylerb/graceful.v1"
)

type contextKey struct{}

type Server struct {
	Handler http.Handler
	addr    string
}

func (s Server) Serve() error {
	// Give the sever 30 seconds to shut down
	return graceful.RunWithErr(s.addr, 30*time.Second, s.Handler)
}

func New(c Controller, addr string) Server {
	controller = c // TODO: get rid of global variable?
	r := mux.NewRouter()

	r.Methods("GET").Path("/v1/books").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetBooksHandler(r.Context(), w, r)
	})

	r.Methods("GET").Path("/v1/books/{bookID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetBookByIDHandler(r.Context(), w, r)
	})

	r.Methods("POST").Path("/v1/books/{bookID}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateBookHandler(r.Context(), w, r)
	})
	handler := withMiddleware(r)
	return Server{Handler: handler, addr: addr}
}
