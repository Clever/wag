package server

// Code auto-generated. Do not edit.

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"gopkg.in/tylerb/graceful.v1"
)

type contextKey struct{}

// Server defines a HTTP server that implements the Controller interface.
type Server struct {
	// Handler should generally not be changed. It exposed to make testing easier.
	Handler http.Handler
	addr    string
}

// Serve starts the server. It will return if an error occurs.
func (s Server) Serve() error {
	// Give the sever 30 seconds to shut down
	return graceful.RunWithErr(s.addr, 30*time.Second, s.Handler)
}

type handler struct {
	Controller
}

// New returns a Server that implements the Controller interface. It will start when "Serve" is called.
func New(c Controller, addr string) Server {
	r := mux.NewRouter()
	h := handler{Controller: c}

	r.Methods("GET").Path("/v1/books").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.GetBooksHandler(r.Context(), w, r)
	})

	r.Methods("GET").Path("/v1/books/{book_id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.GetBookByIDHandler(r.Context(), w, r)
	})

	r.Methods("POST").Path("/v1/books/{book_id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.CreateBookHandler(r.Context(), w, r)
	})

	r.Methods("GET").Path("/v1/health/check").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.HealthCheckHandler(r.Context(), w, r)
	})
	handler := withMiddleware("Swagger Test", r)
	return Server{Handler: handler, addr: addr}
}
