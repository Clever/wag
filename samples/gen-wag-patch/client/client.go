package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Clever/wag/samples/gen-wag-patch/models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	discovery "github.com/Clever/discovery-go"
)

var _ = json.Marshal
var _ = strings.Replace
var _ = strconv.FormatInt
var _ = bytes.Compare

// WagClient is used to make requests to the wag-patch service.
type WagClient struct {
	basePath    string
	requestDoer doer
	transport   *http.Transport
	timeout     time.Duration
	// Keep the retry doer around so that we can set the number of retries
	retryDoer      *retryDoer
	defaultTimeout time.Duration
}

var _ Client = (*WagClient)(nil)

// New creates a new client. The base path and http transport are configurable.
func New(basePath string) *WagClient {
	base := baseDoer{}
	tracing := tracingDoer{d: base}
	retry := retryDoer{d: tracing, defaultRetries: 1}

	return &WagClient{requestDoer: &retry, retryDoer: &retry, defaultTimeout: 10 * time.Second,
		transport: &http.Transport{}, basePath: basePath}
}

// NewFromDiscovery creates a client from the discovery environment variables. This method requires
// the three env vars: SERVICE_WAG_PATCH_HTTP_(HOST/PORT/PROTO) to be set. Otherwise it returns an error.
func NewFromDiscovery() (*WagClient, error) {
	url, err := discovery.URL("wag-patch", "default")
	if err != nil {
		url, err = discovery.URL("wag-patch", "http") // Added fallback to maintain reverse compatibility
		if err != nil {
			return nil, err
		}
	}
	return New(url), nil
}

// WithRetries returns a new client that retries all GET operations until they either succeed or fail the
// number of times specified.
func (c *WagClient) WithRetries(retries int) *WagClient {
	c.retryDoer.defaultRetries = retries
	return c
}

// WithTimeout returns a new client that has the specified timeout on all operations. To make a single request
// have a timeout use context.WithTimeout as described here: https://godoc.org/golang.org/x/net/context#WithTimeout.
func (c *WagClient) WithTimeout(timeout time.Duration) *WagClient {
	c.defaultTimeout = timeout
	return c
}

// JoinByFormat joins a string array by a known format:
//	 csv: comma separated value (default)
//	 ssv: space separated value
//	 tsv: tab separated value
//	 pipes: pipe (|) separated value
func JoinByFormat(data []string, format string) string {
	if len(data) == 0 {
		return ""
	}
	var sep string
	switch format {
	case "ssv":
		sep = " "
	case "tsv":
		sep = "\t"
	case "pipes":
		sep = "|"
	default:
		sep = ","
	}
	return strings.Join(data, sep)
}

// Wagpatch makes a PATCH request to /wagpatch.
// Special wag patch type
func (c *WagClient) Wagpatch(ctx context.Context, i *models.Data) (*models.Data, error) {
	path := c.basePath + "/wagpatch"
	urlVals := url.Values{}
	var body []byte

	path = path + "?" + urlVals.Encode()

	if i != nil {

		var err error
		body, err = json.Marshal(i)

		if err != nil {
			return nil, err
		}

	}

	client := &http.Client{Transport: c.transport}
	req, err := http.NewRequest("PATCH", path, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	// Add the opname for doers like tracing
	ctx = context.WithValue(ctx, opNameCtx{}, "wagpatch")
	req = req.WithContext(ctx)
	// Don't add the timeout in a "doer" because we don't want to call "defer.cancel()"
	// until we've finished all the processing of the request object. Otherwise we'll cancel
	// our own request before we've finished it.
	if c.defaultTimeout != 0 {
		ctx, cancel := context.WithTimeout(req.Context(), c.defaultTimeout)
		defer cancel()
		req = req.WithContext(ctx)
	}
	resp, err := c.requestDoer.Do(client, req)

	if err != nil {
		return nil, models.DefaultInternalError{Msg: err.Error()}
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		var output models.Data

		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}

		return &output, nil
	case 400:
		var output models.DefaultBadRequest

		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}

		return nil, output
	case 500:
		var output models.DefaultInternalError

		if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
			return nil, models.DefaultInternalError{Msg: err.Error()}
		}

		return nil, output

	default:
		return nil, models.DefaultInternalError{Msg: "Unknown response"}
	}
}
