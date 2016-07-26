package main

import (
	"net/http"

	"gopkg.in/Clever/kayvee-go.v3/logger"
	kvMiddleware "gopkg.in/Clever/kayvee-go.v3/middleware"

	"github.com/gorilla/mux"

	"golang.org/x/net/context"
)

// TODO: Need a way to build custom middleware...
// TODO: This we can just copy and let people play around with? Or should we just re-generate this???
func withMiddleware(router *mux.Router) http.Handler {

	// I guess the context handlers should go earlier

	// TODO: We should have some good standard logging approach...
	handler := kvMiddleware.New(router, logger.New("TODO: CHANGE THIS"))
	ctxHandler := convertToContextHandler(handler)
	ctxHandler = tracingHandler(ctxHandler)
	return convertFromContextHandler(ctxHandler)
}

type ContextHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request)

func convertToContextHandler(h http.Handler) ContextHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		// In this case we know that the context is empty...
		h.ServeHTTP(w, r)
	}
}

func convertFromContextHandler(c ContextHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a context object that we will pass through the rest of the handlers
		// TODO: This should really get created in a middleware itself...
		var ctx context.Context
		c(ctx, w, r)
	})
}

func tracingHandler(c ContextHandler) ContextHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(ctx, "trace-id", "tempID")
		c(newCtx, w, r)
	}
}
