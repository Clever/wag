package test

import (
	"context"
	"net/http/httptest"
	"strconv"

	"github.com/Clever/wag/v6/samples/gen-go/models"
	"github.com/Clever/wag/v6/samples/gen-go/server"
)

// ControllerImpl implements the test server controller interface.
type ControllerImpl struct {
	books           map[int64]*models.Book
	maxID           int64
	pageSize        int
	authors         []*models.Author
	nilPutBookCount int
}

// GetBooks returns a list of books.
func (c *ControllerImpl) GetBooks(ctx context.Context, input *models.GetBooksInput) ([]models.Book, int64, error) {
	var bookList []models.Book

	var nextPage int64 // default to no next page

	idx := int64(0)
	if input.StartingAfter != nil {
		idx = *input.StartingAfter + 1
	}

	// loop through all indices to get predictable ordering
	for ; idx <= c.maxID; idx++ {
		if book, ok := c.books[idx]; ok {
			bookList = append(bookList, *book)
			if len(bookList) == c.pageSize {
				nextPage = idx
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

	book, ok := c.books[input.BookID]
	if !ok {
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

	book, ok := c.books[int64(i)]
	if !ok {
		return nil, models.Error{}
	}
	return book, nil
}

// CreateBook creates a book.
func (c *ControllerImpl) CreateBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	if input.ID > c.maxID {
		c.maxID = input.ID
	}
	c.books[input.ID] = input
	return input, nil
}

// PutBook creates a book.
func (c *ControllerImpl) PutBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	if input == nil {
		c.nilPutBookCount++
		return nil, nil
	}
	if input.ID > c.maxID {
		c.maxID = input.ID
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

// GetAuthorsWithPut gets authors with a PUT (because it needs a body).
func (c *ControllerImpl) GetAuthorsWithPut(ctx context.Context, i *models.GetAuthorsWithPutInput) (
	*models.AuthorsResponse, string, error,
) {
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
	controller := ControllerImpl{books: make(map[int64]*models.Book), pageSize: 100}

	s := server.New(&controller, ":8080")

	return httptest.NewServer(s.Handler), &controller
}
