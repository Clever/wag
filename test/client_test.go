package test

import (
	"fmt"

	"github.com/stretchr/testify/assert"

	"github.com/Clever/wag/generated/client"
	"github.com/Clever/wag/generated/models"
	"github.com/Clever/wag/generated/server"

	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"
)

type ClientContextTest struct {
	getCount  int
	postCount int
}

func (c *ClientContextTest) GetBooks(ctx context.Context, input *models.GetBooksInput) (*[]models.Book, error) {
	c.getCount++
	if c.getCount == 1 {
		return nil, fmt.Errorf("Error count: %d", c.getCount)
	}
	return &[]models.Book{}, nil
}
func (c *ClientContextTest) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (models.GetBookByIDOutput, error) {
	return nil, nil
}
func (c *ClientContextTest) CreateBook(ctx context.Context, input *models.CreateBookInput) (*models.Book, error) {
	c.postCount++
	if c.postCount == 1 {
		return nil, fmt.Errorf("Error count: %d", c.postCount)
	}
	return &models.Book{}, nil
}

func TestDefaultClientRetries(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, ":8080")
	testServer := httptest.NewServer(s.Handler)
	c := client.New(testServer.URL)
	_, err := c.GetBooks(context.Background(), &models.GetBooksInput{})
	assert.NoError(t, err)
	assert.Equal(t, 2, controller.getCount)
}

func TestCustomClientRetries(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, ":8080")
	testServer := httptest.NewServer(s.Handler)

	// Should fail if no retries
	c := client.New(testServer.URL).WithRetries(0)
	_, err := c.GetBooks(context.Background(), &models.GetBooksInput{})
	assert.Error(t, err)
	assert.Equal(t, 1, controller.getCount)
}

func TestCustomContextRetries(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, ":8080")
	testServer := httptest.NewServer(s.Handler)

	// Should fail if no retries
	c := client.New(testServer.URL)
	_, err := c.GetBooks(client.WithRetry(context.Background(), 0), &models.GetBooksInput{})
	assert.Error(t, err)
	assert.Equal(t, 1, controller.getCount)
}

func TestNonGetRetries(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, ":8080")
	testServer := httptest.NewServer(s.Handler)
	c := client.New(testServer.URL)
	_, err := c.CreateBook(context.Background(), &models.CreateBookInput{})
	assert.Error(t, err)
	assert.Equal(t, 1, controller.postCount)
}
