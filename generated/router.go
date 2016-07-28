package main

import "github.com/gorilla/mux"

func withRoutes(r *mux.Router) *mux.Router {
	r.Methods("GET").Path("/books/{id}").HandlerFunc(GetBookByIDHandler)
	return r
}
