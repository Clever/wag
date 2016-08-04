package client

import (
	"fmt"
	"net/http"
)
import "golang.org/x/net/context"

import opentracing "github.com/opentracing/opentracing-go"

type RequestHandler interface {
	// If the response returns an error, everything above it should stay the same...
	HandleRequest(ctx context.Context, c *http.Client, r *http.Request) (*http.Response, error)
}

// TODO: Add a timeout handler, probably based on https://godoc.org/golang.org/x/net/context/ctxhttp#Do

type opNameCtx struct{}

// baseRequestHandler performs the base http request
type baseRequestHandler struct{}

func (h baseRequestHandler) HandleRequest(
	ctx context.Context, c *http.Client, r *http.Request) (*http.Response, error) {

	return c.Do(r)
}

// tracingHandler adds tracing to http requests
type tracingHandler struct {
	handler RequestHandler
}

func (h tracingHandler) HandleRequest(
	ctx context.Context, c *http.Client, r *http.Request) (*http.Response, error) {

	opName := ctx.Value(opNameCtx{}).(string)
	var sp opentracing.Span
	// TODO: add tags relating to input data?
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		sp = opentracing.StartSpan(opName, opentracing.ChildOf(parentSpan.Context()))
	} else {
		sp = opentracing.StartSpan(opName)
	}
	if err := sp.Tracer().Inject(sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header)); err != nil {
		return nil, fmt.Errorf("couldn't inject tracing headers (%v)", err)
	}
	return h.handler.HandleRequest(ctx, c, r)
}

// retryHandler retries 50X http requests
type retryHandler struct {
	handler RequestHandler
}

func (h retryHandler) HandleRequest(
	ctx context.Context, c *http.Client, r *http.Request) (*http.Response, error) {

	resp, err := h.handler.HandleRequest(ctx, c, r)
	if err != nil {
		return resp, err
	}

	// If request is retryable then retry it. For this proof of concept let's just retry GETs
	if r.Method != "GET" || resp.StatusCode < 500 {
		return resp, err
	}

	// Otherwise retry the request. For now just do it once without waiting...
	return h.handler.HandleRequest(ctx, c, r)
}
