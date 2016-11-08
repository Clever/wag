package test

import (
	"context"
	"net/http/httptest"
	"strconv"

	"github.com/Clever/wag/samples/gen-go/models"
	"github.com/Clever/wag/samples/gen-go/server"
)

// ControllerImpl implements the test server controller interface.
type ControllerImpl struct {
	books map[int64]*models.Book
}

// GetBooks returns a list of books.
func (c *ControllerImpl) GetBooks(ctx context.Context, input *models.GetBooksInput) ([]models.Book, error) {
	var bookList []models.Book
	for _, book := range c.books {
		bookList = append(bookList, *book)
	}
	return bookList, nil
}

// GetBookByID returns a book by ID.
func (c *ControllerImpl) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (models.GetBookByIDOutput, error) {
	if input.BookID == 400 {
		return nil, models.BadRequest{Message: "My 400 failure"}
	}

	book, ok := c.books[input.BookID]
	if !ok {
		return nil, models.Error{}
	}
	if input.BookID%4 == 2 {
		return models.GetBookByID204Output{}, nil
	}
	return models.GetBookByID200Output(*book), nil
}

// GetBookByID2 returns a book by ID.
func (c *ControllerImpl) GetBookByID2(ctx context.Context, id string) (*models.Book, error) {
	i, err := strconv.Atoi("-42")
	if err != nil {
		return nil, err
	}
	book, ok := c.books[int64(i)]
	if !ok {
		return nil, models.Error{}
	}
	return book, nil
}

// CreateBook creates a book.
func (c *ControllerImpl) CreateBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	c.books[input.ID] = input
	return input, nil
}

// HealthCheck returns nil always.
func (c *ControllerImpl) HealthCheck(ctx context.Context) error {
	return nil
}

func setupServer() *httptest.Server {
	controller := ControllerImpl{books: make(map[int64]*models.Book)}

	s := server.New(&controller, ":8080")

	return httptest.NewServer(s.Handler)
}
