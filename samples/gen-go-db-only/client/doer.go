package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// doer is an interface for "doing" http requests possibly with wrapping
type doer interface {
	Do(c *http.Client, r *http.Request) (*http.Response, error)
}

type opNameCtx struct{}

// baseRequestHandler performs the base http request
type baseDoer struct{}

func (d baseDoer) Do(c *http.Client, r *http.Request) (*http.Response, error) {
	return c.Do(r)
}

// retryHandler retries 50X http requests
type retryDoer struct {
	d           doer
	retryPolicy RetryPolicy
}

// RetryPolicy defines a retry policy.
type RetryPolicy interface {
	// Backoffs returns the number and timing of retry attempts.
	Backoffs() []time.Duration
	// Retry receives the http request, as well as the result of
	// net/http.Client's `Do` method.
	Retry(*http.Request, *http.Response, error) bool
}

// SingleRetryPolicy defines a retry that retries a request once
type SingleRetryPolicy struct{}

// Backoffs returns that you should retry the request 1second after it fails.
func (SingleRetryPolicy) Backoffs() []time.Duration {
	return []time.Duration{1 * time.Second}
}

// Retry will retry non-POST, non-PATCH requests that 5XX.
// TODO: It does not currently retry any errors returned by net/http.Client's `Do`.
func (SingleRetryPolicy) Retry(req *http.Request, resp *http.Response, err error) bool {
	if err != nil || req.Method == "POST" || req.Method == "PATCH" ||
		resp.StatusCode < 500 {
		return false
	}
	return true
}

// ExponentialRetryPolicy defines an exponential retry policy
type ExponentialRetryPolicy struct{}

// Backoffs returns five backoffs with exponentially increasing wait times
// between requests: 100, 200, 400, 800, and 1600 milliseconds +/- up to 5% jitter.
func (ExponentialRetryPolicy) Backoffs() []time.Duration {
	ret := make([]time.Duration, 5)
	next := 100 * time.Millisecond
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	e := 0.05 // +/- 5 percent jitter
	for i := range ret {
		ret[i] = next + time.Duration(((rnd.Float64()*2)-1)*e*float64(next))
		next *= 2
	}
	return ret
}

// Retry will retry non-POST, non-PATCH requests that 5XX.
// TODO: It does not currently retry any errors returned by net/http.Client's `Do`.
func (ExponentialRetryPolicy) Retry(req *http.Request, resp *http.Response, err error) bool {
	if err != nil || req.Method == "POST" || req.Method == "PATCH" ||
		resp.StatusCode < 500 {
		return false
	}
	return true
}

// NoRetryPolicy defines a policy of never retrying a request.
type NoRetryPolicy struct{}

// Backoffs returns an empty slice.
func (NoRetryPolicy) Backoffs() []time.Duration {
	return []time.Duration{}
}

// Retry always returns false.
func (NoRetryPolicy) Retry(*http.Request, *http.Response, error) bool {
	return false
}

type retryContext struct{}

// WithRetryPolicy returns a new context that overrides the client object's
// retry policy.
func WithRetryPolicy(ctx context.Context, retryPolicy RetryPolicy) context.Context {
	return context.WithValue(ctx, retryContext{}, retryPolicy)
}

func (d *retryDoer) Do(c *http.Client, r *http.Request) (*http.Response, error) {
	retryPolicy, ok := r.Context().Value(retryContext{}).(RetryPolicy)
	if !ok {
		retryPolicy = d.retryPolicy
	}
	backoffs := retryPolicy.Backoffs()
	var resp *http.Response
	var err error

	// Save the request body in case we have to retry. Otherwise we will have already read
	// the buffer on retry and the request will fail. See
	// http://stackoverflow.com/questions/23070876/reading-body-of-http-request-without-modifying-request-state
	var buf []byte
	if r.Body != nil {
		var err error
		buf, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
	}

	for retries := 0; true; retries++ {
		if r.Body != nil {
			rdr := ioutil.NopCloser(bytes.NewBuffer(buf))
			r.Body = rdr
		}
		resp, err = d.d.Do(c, r)
		if retries == len(backoffs) || !retryPolicy.Retry(r, resp, err) {
			break
		}
		// Close the response body if response is not nil
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(backoffs[retries])
	}
	return resp, err
}
