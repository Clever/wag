package main

import (
	"net/http"

	"gopkg.in/Clever/kayvee-go.v3/logger"
	kvMiddleware "gopkg.in/Clever/kayvee-go.v3/middleware"

	"golang.org/x/net/context"

	gContext "github.com/gorilla/context"
)

// TODO: Need a way to build custom middleware...
// TODO: This we can just copy and let people play around with? Or should we just re-generate this???
func withMiddleware(router http.Handler) http.Handler {

	// We layer middleware as follows
	// 1. Request received
	// 2. Add non-context middleware (for backwards compatibility)
	// 3. Convert to a context aware middleware
	// 4. Add context aware middleware
	// 5. Conver back

	// TODO: Add Add standard, non-context aware middleware
	handler := kvMiddleware.New(router, logger.New("TODO: CHANGE THIS"))

	// We shouldn't have to do all the conversion sillyness in go 1.7 as it appears
	// it will add a context object to the http.Request object
	// (https://github.com/golang/go/issues/14660)

	// TODO: We should have some good standard logging approach...

	ctxHandler := convertToContextHandler(handler)

	// This is a bit annoying...
	ctxHandler = addContextToGorilla(ctxHandler)

	ctxHandler = modifyContextExample(ctxHandler)
	// Add in other handlers...

	return ctxHandler
}

type ContextHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request)

func (c ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	c(ctx, w, r)
}

func convertToContextHandler(h http.Handler) ContextHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

func modifyContextExample(c ContextHandler) ContextHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = context.WithValue(ctx, "addedKey", "addedValue")
		c(ctx, w, r)
	}
}

func addContextToGorilla(c ContextHandler) ContextHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		gContext.Set(r, contextKey{}, ctx)
		c(ctx, w, r)
	}
}
