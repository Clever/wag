package test

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"net/http/httptrace"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Clever/wag/samples/gen-go/client"
	"github.com/Clever/wag/samples/gen-go/models"
	"github.com/Clever/wag/samples/gen-go/server"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ClientContextTest struct {
	getCount      int
	getTimes      []time.Time
	getErrorCount int
	postCount     int
}

func (c *ClientContextTest) GetBooks(ctx context.Context, input *models.GetBooksInput) ([]models.Book, error) {
	c.getCount++
	c.getTimes = append(c.getTimes, time.Now())
	if c.getCount <= c.getErrorCount {
		return nil, fmt.Errorf("Error count: %d", c.getCount)
	}
	return []models.Book{}, nil
}
func (c *ClientContextTest) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (*models.Book, error) {
	return nil, nil
}
func (c *ClientContextTest) GetBookByID2(ctx context.Context, id string) (*models.Book, error) {
	return nil, nil
}
func (c *ClientContextTest) CreateBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	c.postCount++
	if c.postCount == 1 {
		return nil, fmt.Errorf("Error count: %d", c.postCount)
	}
	return &models.Book{}, nil
}

func (c *ClientContextTest) HealthCheck(ctx context.Context) error {
	return nil
}

type ClientCircuitTest struct {
	down bool
}

func (c *ClientCircuitTest) GetBooks(ctx context.Context, input *models.GetBooksInput) ([]models.Book, error) {
	if c.down {
		return nil, errors.New("fail")
	}
	return []models.Book{}, nil
}
func (c *ClientCircuitTest) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (*models.Book, error) {
	if c.down {
		return nil, errors.New("fail")
	}
	return nil, nil
}
func (c *ClientCircuitTest) GetBookByID2(ctx context.Context, id string) (*models.Book, error) {
	if c.down {
		return nil, errors.New("fail")
	}
	return nil, nil
}
func (c *ClientCircuitTest) CreateBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	if c.down {
		return nil, errors.New("fail")
	}
	return &models.Book{}, nil
}

func (c *ClientCircuitTest) HealthCheck(ctx context.Context) error {
	if c.down {
		return errors.New("fail")
	}
	return nil
}

func TestDefaultClientRetries(t *testing.T) {
	controller := ClientContextTest{getErrorCount: 2}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()
	c := client.New(testServer.URL)
	_, err := c.GetBooks(context.Background(), &models.GetBooksInput{})
	require.NoError(t, err)
	require.Equal(t, len(controller.getTimes), 3, "expected three requests")
	assert.WithinDuration(t, controller.getTimes[1], controller.getTimes[0].Add(100*time.Millisecond), 10*time.Millisecond,
		"expected first backoff to be about 100ms")
	assert.WithinDuration(t, controller.getTimes[2], controller.getTimes[1].Add(200*time.Millisecond), 20*time.Millisecond,
		"expected first backoff to be about 200ms")
}

func TestCustomClientRetries(t *testing.T) {
	controller := ClientContextTest{getErrorCount: 1}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()

	// Should fail if no retries
	c := client.New(testServer.URL).WithRetryPolicy(client.NoRetryPolicy{})
	_, err := c.GetBooks(context.Background(), &models.GetBooksInput{})
	assert.Error(t, err)
	assert.Equal(t, 1, controller.getCount)
}

func TestCustomContextRetries(t *testing.T) {
	controller := ClientContextTest{getErrorCount: 1}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()

	// Should fail if no retries
	c := client.New(testServer.URL)
	_, err := c.GetBooks(client.WithRetryPolicy(context.Background(), client.NoRetryPolicy{}), &models.GetBooksInput{})
	assert.Error(t, err)
	assert.Equal(t, 1, controller.getCount)
}

func TestNonGetRetries(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()
	c := client.New(testServer.URL)
	_, err := c.CreateBook(context.Background(), &models.Book{})
	assert.Error(t, err)
	assert.Equal(t, 1, controller.postCount)
}

func TestNetworkErrorRetries(t *testing.T) {
	c := client.New("https://thisshouldnotresolve1234567890.com/")
	_, err := c.CreateBook(context.Background(), &models.Book{})
	assert.Error(t, err)
}

func TestNewWithDiscovery(t *testing.T) {
	controller := ClientContextTest{getErrorCount: 1}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)

	// Should be an err if env vars aren't set
	_, err := client.NewFromDiscovery()
	assert.Error(t, err)

	splitURL := strings.Split(testServer.URL, ":")
	assert.Equal(t, 3, len(splitURL))

	os.Setenv("SERVICE_SWAGGER_TEST_DEFAULT_PROTO", "http")
	os.Setenv("SERVICE_SWAGGER_TEST_DEFAULT_PORT", splitURL[2])
	os.Setenv("SERVICE_SWAGGER_TEST_DEFAULT_HOST", splitURL[1][2:])

	c, err := client.NewFromDiscovery()
	assert.NoError(t, err)
	_, err = c.GetBooks(context.Background(), &models.GetBooksInput{})
	assert.NoError(t, err)
	assert.Equal(t, 2, controller.getCount)

	// Testing fallback
	os.Unsetenv("SERVICE_SWAGGER_TEST_DEFAULT_PROTO")
	os.Unsetenv("SERVICE_SWAGGER_TEST_DEFAULT_PORT")
	os.Unsetenv("SERVICE_SWAGGER_TEST_DEFAULT_HOST")
	os.Setenv("SERVICE_SWAGGER_TEST_HTTP_PROTO", "http")
	os.Setenv("SERVICE_SWAGGER_TEST_HTTP_PORT", splitURL[2])
	os.Setenv("SERVICE_SWAGGER_TEST_HTTP_HOST", splitURL[1][2:])

	c, err = client.NewFromDiscovery()
	assert.NoError(t, err)
	_, err = c.GetBooks(context.Background(), &models.GetBooksInput{})
	assert.NoError(t, err)
	assert.Equal(t, 3, controller.getCount)
}

func TestCircuitBreaker(t *testing.T) {
	controller := ClientCircuitTest{}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()
	hystrix.Flush()
	c := client.New(testServer.URL)
	c.SetCircuitBreakerDebug(false)
	c.SetCircuitBreakerSettings(client.CircuitBreakerSettings{
		MaxConcurrentRequests:  client.DefaultCircuitBreakerSettings.MaxConcurrentRequests,
		RequestVolumeThreshold: 1,
		SleepWindow:            2000,
		ErrorPercentThreshold:  client.DefaultCircuitBreakerSettings.ErrorPercentThreshold,
	})

	// the circuit should open after one or two failed attempts, since the volume
	// threshold (set above) is 1.
	controller.down = true
	var connAttempts int64
	ctx := httptrace.WithClientTrace(context.Background(),
		&httptrace.ClientTrace{
			GetConn: func(hostPort string) {
				atomic.AddInt64(&connAttempts, 1)
			},
		})

	_, err := c.CreateBook(ctx, &models.Book{})
	assert.Error(t, err)
	_, err = c.CreateBook(ctx, &models.Book{})
	assert.Error(t, err)
	_, err = c.CreateBook(ctx, &models.Book{})
	assert.Error(t, err)
	assert.Equal(t, true, connAttempts <= int64(2), "circuit should have opened, saw too many connection attempts: %d", connAttempts)

	// we should see an attempt go through after two seconds (this is the
	// sleep window configured above).
	circuitOpened := time.Now()
	for _ = range time.Tick(100 * time.Millisecond) {
		_, err := c.CreateBook(ctx, &models.Book{})
		assert.Error(t, err)
		if connAttempts == 2 {
			assert.WithinDuration(t, time.Now(), circuitOpened,
				2*time.Second+500*time.Millisecond)
			break
		}
		if time.Now().Sub(circuitOpened) > 10*time.Second {
			t.Fatal("circuit should let through a 2nd attempt by now")
		}
	}

	// bring the server back up, and we should see successes after another
	// two seconds, for a total of 4 seconds.
	controller.down = false
	for _ = range time.Tick(100 * time.Millisecond) {
		_, err := c.CreateBook(ctx, &models.Book{})
		if err == nil {
			assert.WithinDuration(t, time.Now(), circuitOpened,
				4*time.Second+500*time.Millisecond)
			break
		}
		if time.Now().Sub(circuitOpened) > 10*time.Second {
			t.Fatal("circuit should have closed by now")
		}
	}
}
