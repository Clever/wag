package main

import (
	"net/http"

	"gopkg.in/Clever/kayvee-go.v3/logger"
	kvMiddleware "gopkg.in/Clever/kayvee-go.v3/middleware"

	"github.com/gorilla/mux"
)

// TODO: Need a way to build custom middleware...
// TODO: This we can just copy and let people play around with? Or should we just re-generate this???
func withMiddleware(router *mux.Router) http.Handler {

	// TODO: We should have some good standard logging approach...
	return kvMiddleware.New(router, logger.New("TODO: CHANGE THIS"))
}
