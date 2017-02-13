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
	books    []*models.Book
	pageSize int
	authors  []*models.Author
}

// GetBooks returns a list of books.
func (c *ControllerImpl) GetBooks(ctx context.Context, input *models.GetBooksInput) ([]models.Book, int64, error) {
	var bookList []models.Book

	begin := 0
	if input.StartingAfter != nil {
		begin = int(*input.StartingAfter) + 1
	}

	nextPage := int64(0) // default to no next page
	for idx, book := range c.books[begin:] {
		if book != nil {
			bookList = append(bookList, *book)
			if len(bookList) == c.pageSize {
				nextPage = int64(idx)
				break
			}
		}
	}

	return bookList, nextPage, nil
}

// GetBookByID returns a book by ID.
func (c *ControllerImpl) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (*models.Book, error) {
	if input.BookID == 400 {
		return nil, models.BadRequest{Message: "My 400 failure"}
	}

	if int(input.BookID) >= len(c.books) {
		return nil, models.Error{}
	}
	book := c.books[input.BookID]
	if book == nil {
		return nil, models.Error{}
	}
	return book, nil
}

// GetBookByID2 returns a book by ID.
func (c *ControllerImpl) GetBookByID2(ctx context.Context, id string) (*models.Book, error) {
	i, err := strconv.Atoi("-42")
	if err != nil {
		return nil, err
	}
	if i >= len(c.books) {
		return nil, models.Error{}
	}
	book := c.books[i]
	if book == nil {
		return nil, models.Error{}
	}
	return book, nil
}

// CreateBook creates a book.
func (c *ControllerImpl) CreateBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	if int(input.ID) >= len(c.books) {
		extension := make([]*models.Book, int(input.ID)-len(c.books)+1)
		c.books = append(c.books, extension...)
	}
	c.books[input.ID] = input
	return input, nil
}

// GetAuthors gets authors.
func (c *ControllerImpl) GetAuthors(ctx context.Context, i *models.GetAuthorsInput) (*models.AuthorsResponse, string, error) {
	return &models.AuthorsResponse{
		AuthorSet: &models.AuthorSet{
			RandomProp: int64(3),
			Results:    models.AuthorArray(c.authors),
		},
	}, "", nil
}

// HealthCheck returns nil always.
func (c *ControllerImpl) HealthCheck(ctx context.Context) error {
	return nil
}

func setupServer() (*httptest.Server, *ControllerImpl) {
	controller := ControllerImpl{pageSize: 100}

	s := server.New(&controller, ":8080")

	return httptest.NewServer(s.Handler), &controller
}
