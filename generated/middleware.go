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

	// Add all non-context aware middleware first. Right now this is just kayvee
	// middleware
	handler := kvMiddleware.New(router, logger.New("TODO: CHANGE THIS"))

	// Add the context.Context to gorilla's context. This is a bit annoying, but I'm not sure
	// how else to get Gorilla's routing since the routing handler functions don't take in a
	// context. This isn't the worst thing in the world since it's localized weirdness, and it
	// should be cleaned up in go 1.7 (https://github.com/golang/go/issues/14660)
	ctxHandler := addContextToGorilla(handler)

	// 2. Add context aware middlware
	// ctxHandler = modifyContextExample(ctxHandler)

	return ContextWrapper{handler: ctxHandler}
}

// ContextWrapper is a struct that converts from the http.Handler to the ContextHandler
// one. It does this by creating a new context when ServeHTTP is called and passes that down
// the stack.
type ContextWrapper struct {
	handler ContextHandler
}

func (c ContextWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	c.handler.ServeHTTPContext(ctx, w, r)
}

// ContextHandler is an interface for handlers / middleware that extends the base Go
// handler interface with a context.Context object
type ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

type ContextHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

func (c ContextHandlerFunc) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c(ctx, w, r)
}

// modifyContextExample is sample middleware that modifies the context. We can remove it when
// we have real middleware that uses the context.
func modifyContextExample(c ContextHandler) ContextHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = context.WithValue(ctx, "addedKey", "addedValue")
		c.ServeHTTPContext(ctx, w, r)
	}
}

// addContextToGorilla adds the context.Context object to the Gorilla context since there doesn't
// seem to be a easier way to get Gorilla routing and contexts in the handlers
func addContextToGorilla(h http.Handler) ContextHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		gContext.Set(r, contextKey{}, ctx)
		h.ServeHTTP(w, r)
	}
}
