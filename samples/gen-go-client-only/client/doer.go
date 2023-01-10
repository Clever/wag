package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	wcl "github.com/Clever/wag/logging/wagclientlogger"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/donovanhide/eventsource"
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

// circuitBreakerDoer implements the circuit breaker pattern.
// Set debug to true to operate the circuit in a mode where it will only log
// circuit state, and will not block any requests from going through.
type circuitBreakerDoer struct {
	d           doer
	debug       bool
	circuitName string
	logger      wcl.WagClientLogger
}

var circuitSSEOnce sync.Once

// HystrixSSEEvent is emitted by hystrix-go via server-sent events. It describes
// the state of a circuit.
type HystrixSSEEvent struct {
	Type                            string `json:"type"`
	Name                            string `json:"name"`
	RequestCount                    int    `json:"requestCount"`
	ErrorCount                      int    `json:"errorCount"`
	ErrorPercentage                 int    `json:"errorPercentage"`
	IsCircuitBreakerOpen            bool   `json:"isCircuitBreakerOpen"`
	RollingCountFailure             int    `json:"rollingCountFailure"`
	RollingCountFallbackFailure     int    `json:"rollingCountFallbackFailure"`
	RollingCountFallbackSuccess     int    `json:"rollingCountFallbackSuccess"`
	RollingCountShortCircuited      int    `json:"rollingCountShortCircuited"`
	RollingCountSuccess             int    `json:"rollingCountSuccess"`
	RollingCountThreadPoolRejected  int    `json:"rollingCountThreadPoolRejected"`
	RollingCountTimeout             int    `json:"rollingCountTimeout"`
	CurrentConcurrentExecutionCount int    `json:"currentConcurrentExecutionCount"`
	LatencyTotalMean                int    `json:"latencyTotal_mean"`
}

func logEvent(l wcl.WagClientLogger, e HystrixSSEEvent) {
	l.Log(wcl.Info, "", map[string]interface{}{
		"requestCount":                    e.RequestCount,
		"errorCount":                      e.ErrorCount,
		"errorPercentage":                 e.ErrorPercentage,
		"isCircuitBreakerOpen":            e.IsCircuitBreakerOpen,
		"rollingCountFailure":             e.RollingCountFailure,
		"rollingCountFallbackFailure":     e.RollingCountFallbackFailure,
		"rollingCountFallbackSuccess":     e.RollingCountFallbackSuccess,
		"rollingCountShortCircuited":      e.RollingCountShortCircuited,
		"rollingCountSuccess":             e.RollingCountSuccess,
		"rollingCountThreadPoolRejected":  e.RollingCountThreadPoolRejected,
		"rollingCountTimeout":             e.RollingCountTimeout,
		"currentConcurrentExecutionCount": e.CurrentConcurrentExecutionCount,
		"latencyTotalMean":                e.LatencyTotalMean,
	})
}

func (d *circuitBreakerDoer) init() {
	// Periodically log internal circuit state to assist in setting
	// circuit thresholds and understanding application behavior.
	// Unfortunately, hystrix-go doesn't have a great way to expose this
	// data, so we resort to turning on its support for broadcasting
	// circuit metrics via http server-sent events (SSE).
	// See https://github.com/afex/hystrix-go/issues/56.
	circuitSSEOnce.Do(func() {
		hystrixStreamHandler := hystrix.NewStreamHandler()
		hystrixStreamHandler.Start()
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			panic(err)
		}

		go http.Serve(listener, hystrixStreamHandler)
		go func() {
			// We log every 10 seconds because the hysterix metrics are on a 10-second
			// rolling window. This let's us sum up the hysterix metrics to get fairly
			// accurate total numbers (note that this isn't perfect because our timing
			// doesn't perfectly match hysterix's, and we also log when the circuit breaker
			// status changes)
			logFrequency := 10 * time.Second
			lastEventSeen := map[string]HystrixSSEEvent{}
			lastEventLogTime := map[string]time.Time{}

			for _ = range time.Tick(1 * time.Second) { // retry indefinitely
				url := "http://" + listener.Addr().String()
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					log.Printf("error connecting to circuit SSE stream: %s", err)
					continue
				}
				stream, err := eventsource.SubscribeWith("", &http.Client{}, req)
				if err != nil {
					log.Printf("error connecting to circuit SSE stream: %s", err)
					continue
				}
				for ev := range stream.Events {
					var e HystrixSSEEvent
					if err := json.Unmarshal([]byte(ev.Data()), &e); err != nil {
						continue
					}
					if e.Type != "HystrixCommand" {
						continue
					}

					// Today we are creating a stream for every client so lets only log events for this
					// circuit. In an ideal world we only create a single stream and pass it to the client.
					// Lets worry about doing this when we implement passing circuitBreakerOptions
					// to disable debug mode
					if e.Name != d.circuitName {
						continue
					}

					lastSeen, ok := lastEventSeen[e.Name]
					lastEventSeen[e.Name] = e

					// log the circuit state if it's either
					// (1) the first event we've seen for the circuit or
					// (2) the circuit open state has changed or
					// (3) 10 seconds have passed since we logged something for the circuit
					if !ok {
						lastEventLogTime[e.Name] = time.Now()
						logEvent(d.logger, e)
						continue
					}
					if lastSeen.IsCircuitBreakerOpen != e.IsCircuitBreakerOpen {
						lastEventLogTime[e.Name] = time.Now()
						logEvent(d.logger, e)
					} else if time.Now().Sub(lastEventLogTime[e.Name]) > logFrequency {
						lastEventLogTime[e.Name] = time.Now()
						logEvent(d.logger, e)
					}
				}
			}
		}()
	})
}

func (d *circuitBreakerDoer) Do(c *http.Client, r *http.Request) (*http.Response, error) {
	if d.debug {
		resp, err := d.d.Do(c, r)
		hystrix.Do(d.circuitName, func() error {
			if err != nil {
				return err
			}
			if resp.StatusCode >= 500 {
				// need to return an error to trigger circuit opening
				return errors.New("5XX")
			}
			return err
		}, nil)
		return resp, err
	}

	var resp *http.Response
	var err error
	err = hystrix.Do(d.circuitName, func() error {
		resp, err = d.d.Do(c, r)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 500 {
			// need to return an error to trigger circuit opening
			return errors.New("5XX")
		}
		return err
	}, nil /* no fallback function, yet */)
	return resp, err
}
