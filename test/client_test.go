package test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/http/httptrace"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Clever/wag/gen-go/client"
	"github.com/Clever/wag/gen-go/models"
	"github.com/Clever/wag/gen-go/server"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
)

type ClientContextTest struct {
	getCount  int
	postCount int
}

func (c *ClientContextTest) GetBooks(ctx context.Context, input *models.GetBooksInput) ([]models.Book, error) {
	c.getCount++
	if c.getCount == 1 {
		return nil, fmt.Errorf("Error count: %d", c.getCount)
	}
	return []models.Book{}, nil
}
func (c *ClientContextTest) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (models.GetBookByIDOutput, error) {
	return nil, nil
}
func (c *ClientContextTest) GetBookByID2(ctx context.Context, input *models.GetBookByID2Input) (*models.Book, error) {
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

func TestDefaultClientRetries(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()
	c := client.New(testServer.URL)
	_, err := c.GetBooks(context.Background(), &models.GetBooksInput{})
	assert.NoError(t, err)
	assert.Equal(t, 2, controller.getCount)
}

func TestCustomClientRetries(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()

	// Should fail if no retries
	c := client.New(testServer.URL).WithRetries(0)
	_, err := c.GetBooks(context.Background(), &models.GetBooksInput{})
	assert.Error(t, err)
	assert.Equal(t, 1, controller.getCount)
}

func TestCustomContextRetries(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()

	// Should fail if no retries
	c := client.New(testServer.URL)
	_, err := c.GetBooks(client.WithRetries(context.Background(), 0), &models.GetBooksInput{})
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

func TestNewWithDiscovery(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)

	// Should be an err if env vars aren't set
	_, err := client.NewFromDiscovery()
	assert.Error(t, err)

	splitURL := strings.Split(testServer.URL, ":")
	assert.Equal(t, 3, len(splitURL))

	os.Setenv("SERVICE_SWAGGER_TEST_HTTP_PROTO", "http")
	os.Setenv("SERVICE_SWAGGER_TEST_HTTP_PORT", splitURL[2])
	os.Setenv("SERVICE_SWAGGER_TEST_HTTP_HOST", splitURL[1][2:])

	c, err := client.NewFromDiscovery()
	assert.NoError(t, err)
	_, err = c.GetBooks(context.Background(), &models.GetBooksInput{})
	assert.NoError(t, err)
	assert.Equal(t, 2, controller.getCount)
}

func TestCircuitBreaker(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, "")
	testServerDown := httptest.NewServer(s.Handler)
	testServerDown.Close()
	testServerUp := httptest.NewServer(s.Handler)
	defer testServerUp.Close()
	hystrix.Flush()

	// the circuit should open after 20 consecutive failed attempts (this is
	// the default volume threshold, after which health can open the circuit)
	c := client.New(testServerDown.URL)
	var connAttempts int64
	ctx := httptrace.WithClientTrace(context.Background(),
		&httptrace.ClientTrace{
			GetConn: func(hostPort string) {
				atomic.AddInt64(&connAttempts, 1)
			},
		})
	for count := 0; count < 100; count++ {
		_, err := c.CreateBook(ctx, &models.Book{})
		assert.Error(t, err)
	}
	assert.Equal(t, int64(20), connAttempts)

	// we should see an attempts go through after five seconds
	//c = c.WithBasePath(testServerUp.URL)
	circuitOpened := time.Now()
	for _ = range time.Tick(100 * time.Millisecond) {
		_, err := c.CreateBook(ctx, &models.Book{})
		assert.Error(t, err)
		if connAttempts == 21 {
			assert.WithinDuration(t, time.Now(), circuitOpened,
				5*time.Second+500*time.Millisecond)
			break
		}
	}

	// bring the server back up, and we should see successes after another 5s
	c = c.WithBasePath(testServerUp.URL)
	for _ = range time.Tick(100 * time.Millisecond) {
		_, err := c.CreateBook(ctx, &models.Book{})
		if err == nil {
			assert.WithinDuration(t, time.Now(), circuitOpened,
				10*time.Second+500*time.Millisecond)
			break
		}
	}
}
