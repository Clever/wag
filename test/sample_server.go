package test

import (
	"context"
	"net/http/httptest"

	"github.com/Clever/wag/generated/models"
	"github.com/Clever/wag/generated/server"
)

type ControllerImpl struct {
	books map[int64]*models.Book
}

func (c *ControllerImpl) GetBooks(ctx context.Context, input *models.GetBooksInput) (models.GetBooksOutput, error) {
	bookList := make([]models.Book, 0)
	for _, book := range c.books {
		bookList = append(bookList, *book)
	}
	return models.GetBooks200Output(bookList), nil
}
func (c *ControllerImpl) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (models.GetBookByIDOutput, error) {
	book, ok := c.books[input.BookID]
	if !ok {
		return nil, models.GetBookByID404Output{}
	}
	return models.GetBookByID200Output(*book), nil
}
func (c *ControllerImpl) CreateBook(ctx context.Context, input *models.CreateBookInput) (models.CreateBookOutput, error) {
	c.books[input.NewBook.ID] = input.NewBook
	return models.CreateBook200Output(*input.NewBook), nil
}

func setupServer() *httptest.Server {
	controller := ControllerImpl{books: make(map[int64]*models.Book)}

	s := server.New(&controller, ":8080")

	return httptest.NewServer(s.Handler)
}
