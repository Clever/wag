package test

import (
	"context"
	"net/http/httptest"

	"github.com/Clever/wag/gen-go/models"
	"github.com/Clever/wag/gen-go/server"
)

type ControllerImpl struct {
	books map[int64]*models.Book
}

func (c *ControllerImpl) GetBooks(ctx context.Context, input *models.GetBooksInput) ([]models.Book, error) {
	bookList := make([]models.Book, 0)
	for _, book := range c.books {
		bookList = append(bookList, *book)
	}
	return bookList, nil
}
func (c *ControllerImpl) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (models.GetBookByIDOutput, error) {
	book, ok := c.books[input.BookID]
	if !ok {
		return nil, models.GetBookByID404Output{}
	}
	if input.BookID%4 == 2 {
		return models.GetBookByID204Output{}, nil
	} else {
		return models.GetBookByID200Output(*book), nil
	}

}
func (c *ControllerImpl) CreateBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	c.books[input.ID] = input
	return input, nil
}
func (c *ControllerImpl) HealthCheck(ctx context.Context) error {
	return nil
}

func setupServer() *httptest.Server {
	controller := ControllerImpl{books: make(map[int64]*models.Book)}

	s := server.New(&controller, ":8080")

	return httptest.NewServer(s.Handler)
}
