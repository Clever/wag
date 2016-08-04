package client

import (
	"fmt"
	"net/http"
)
import "golang.org/x/net/context"

import opentracing "github.com/opentracing/opentracing-go"

// doer is an interface for "doing" http requests possibly with wrapping
type doer interface {
	Do(ctx context.Context, c *http.Client, r *http.Request) (*http.Response, error)
}

type opNameCtx struct{}

// baseRequestHandler performs the base http request
type baseDoer struct{}

func (d baseDoer) Do(ctx context.Context, c *http.Client, r *http.Request) (*http.Response, error) {
	// TODO: Add a timeout handler, probably based on https://godoc.org/golang.org/x/net/context/ctxhttp#Do
	return c.Do(r)
}

// tracingDoer adds tracing to http requests
type tracingDoer struct {
	d doer
}

func (d tracingDoer) Do(
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
	return d.d.Do(ctx, c, r)
}

// retryHandler retries 50X http requests
type retryDoer struct {
	d              doer
	defaultRetries int
}

// WithRetry returns a new context that overrides the number of retries to do for a particular
// request.
func WithRetry(ctx context.Context, retries int) context.Context {
	return context.WithValue(ctx, retryContext{}, retries)
}

// retryContext is the key the retry configuration. For demonstration purposes it's just a count
// of the number of retries right now.
type retryContext struct{}

func (d retryDoer) Do(ctx context.Context, c *http.Client, r *http.Request) (*http.Response, error) {

	resp, err := d.d.Do(ctx, c, r)
	if err != nil {
		return resp, err
	}

	// If the request can't be retried then just return immediately. For this proof of concept only
	// GETs can be retried
	if r.Method != "GET" {
		return resp, err
	}

	var retries int
	retries, ok := ctx.Value(retryContext{}).(int)
	if !ok {
		retries = d.defaultRetries
	}

	for i := 0; i < retries; i++ {
		if resp.StatusCode < 500 {
			break
		}
		resp, err = d.d.Do(ctx, c, r)
	}
	return resp, err
}
