package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/donovanhide/eventsource"
	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context/ctxhttp"
	logger "gopkg.in/Clever/kayvee-go.v5/logger"
)

// doer is an interface for "doing" http requests possibly with wrapping
type doer interface {
	Do(c *http.Client, r *http.Request) (*http.Response, error)
}

type opNameCtx struct{}

// baseRequestHandler performs the base http request
type baseDoer struct{}

func (d baseDoer) Do(c *http.Client, r *http.Request) (*http.Response, error) {
	return ctxhttp.Do(r.Context(), c, r)
}

// tracingDoer adds tracing to http requests
type tracingDoer struct {
	d doer
}

func (d tracingDoer) Do(c *http.Client, r *http.Request) (*http.Response, error) {

	ctx := r.Context()
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
	return d.d.Do(c, r)
}

// retryHandler retries 50X http requests
type retryDoer struct {
	d              doer
	defaultRetries int
}

// WithRetries returns a new context that overrides the number of retries to do for a particular
// request.
func WithRetries(ctx context.Context, retries int) context.Context {
	return context.WithValue(ctx, retryContext{}, retries)
}

// retryContext is the key the retry configuration. For demonstration purposes it's just a count
// of the number of retries right now.
type retryContext struct{}

func (d *retryDoer) Do(c *http.Client, r *http.Request) (*http.Response, error) {

	resp, err := d.d.Do(c, r)
	if err != nil {
		return resp, err
	}

	// If the request can't be retried then just return immediately. We retry all idempotent
	// http requests. The only two that can't be retried are post and patch
	if r.Method == "POST" || r.Method == "PATCH" {
		return resp, err
	}

	var retries int
	retries, ok := r.Context().Value(retryContext{}).(int)
	if !ok {
		retries = d.defaultRetries
	}

	for i := 0; i < retries; i++ {
		if resp.StatusCode < 500 {
			break
		}
		resp.Body.Close()
		resp, err = d.d.Do(c, r)
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

func logEvent(l *logger.Logger, e HystrixSSEEvent) {
	l.InfoD(e.Name, map[string]interface{}{
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
			logger := logger.New("wag")
			logFrequency := 30 * time.Second
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
					lastSeen, ok := lastEventSeen[e.Name]
					lastEventSeen[e.Name] = e

					// log the circuit state if it's either
					// (1) the first event we've seen for the circuit or
					// (2) the circuit open state has changed or
					// (3) 30 seconds have passed since we logged something for the circuit
					if !ok {
						lastEventLogTime[e.Name] = time.Now()
						logEvent(logger, e)
						continue
					}
					if lastSeen.IsCircuitBreakerOpen != e.IsCircuitBreakerOpen {
						lastEventLogTime[e.Name] = time.Now()
						logEvent(logger, e)
					} else if time.Now().Sub(lastEventLogTime[e.Name]) > logFrequency {
						lastEventLogTime[e.Name] = time.Now()
						logEvent(logger, e)
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
