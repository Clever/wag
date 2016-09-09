package test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Clever/wag/gen-go/client"
	"github.com/Clever/wag/gen-go/models"
	"github.com/Clever/wag/gen-go/server"
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
	c := client.New(testServer.URL)
	_, err := c.GetBooks(context.Background(), &models.GetBooksInput{})
	assert.NoError(t, err)
	assert.Equal(t, 2, controller.getCount)
}

func TestCustomClientRetries(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)

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
	c := client.New(testServer.URL)
	_, err := c.CreateBook(context.Background(), &models.Book{})
	assert.Error(t, err)
	assert.Equal(t, 1, controller.postCount)
}
