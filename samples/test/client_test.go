package test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Clever/wag/logging/wagclientlogger"
	"github.com/Clever/wag/samples/gen-go-basic/client/v9"
	"github.com/Clever/wag/samples/gen-go-basic/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-basic/server"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type NoopLogger struct{}

// Log is a no-op implementation of the Log method.
func (n *NoopLogger) Log(level wagclientlogger.LogLevel, message string, pairs map[string]interface{}) {
	// No operation performed
}

var wcl wagclientlogger.WagClientLogger = &NoopLogger{}

type ClientContextTest struct {
	getCount      int
	getTimes      []time.Time
	getErrorCount int
	postCount     int
	putCount      int
}

func (c *ClientContextTest) GetBooks(ctx context.Context, input *models.GetBooksInput) ([]models.Book, int64, error) {
	c.getCount++
	c.getTimes = append(c.getTimes, time.Now())
	if c.getCount <= c.getErrorCount {
		return nil, int64(0), fmt.Errorf("Error count: %d", c.getCount)
	}
	return []models.Book{}, int64(0), nil
}
func (c *ClientContextTest) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (*models.Book, error) {
	return nil, nil
}
func (c *ClientContextTest) GetBookByID2(ctx context.Context, id string) (*models.Book, error) {
	return nil, nil
}
func (c *ClientContextTest) PutBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	c.putCount++
	if c.putCount == 1 {
		return nil, fmt.Errorf("Error count: %d", c.putCount)
	}
	return input, nil
}
func (c *ClientContextTest) CreateBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	c.postCount++
	if c.postCount == 1 {
		return nil, fmt.Errorf("Error count: %d", c.postCount)
	}
	return &models.Book{}, nil
}
func (c *ClientContextTest) GetAuthors(ctx context.Context, i *models.GetAuthorsInput) (*models.AuthorsResponse, string, error) {
	return nil, "", nil
}
func (c *ClientContextTest) GetAuthorsWithPut(ctx context.Context, i *models.GetAuthorsWithPutInput) (*models.AuthorsResponse, string, error) {
	return nil, "", nil
}
func (c *ClientContextTest) HealthCheck(ctx context.Context) error {
	return nil
}

func (c *ClientContextTest) LowercaseModelsTest(ctx context.Context, i *models.LowercaseModelsTestInput) error {
	return nil
}

type ClientPutPagingTest struct {
	pageToReturn        string
	t                   *testing.T
	expectedRequestBody *models.Book
	timesPutCalled      int
}

func (c *ClientPutPagingTest) GetBooks(ctx context.Context, input *models.GetBooksInput) ([]models.Book, int64, error) {
	return nil, int64(0), nil
}
func (c *ClientPutPagingTest) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (*models.Book, error) {
	return nil, nil
}
func (c *ClientPutPagingTest) GetBookByID2(ctx context.Context, id string) (*models.Book, error) {
	return nil, nil
}
func (c *ClientPutPagingTest) CreateBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	return nil, nil
}
func (c *ClientPutPagingTest) PutBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	return nil, nil
}
func (c *ClientPutPagingTest) GetAuthors(ctx context.Context, i *models.GetAuthorsInput) (*models.AuthorsResponse, string, error) {
	return nil, "", nil
}

func (c *ClientPutPagingTest) LowercaseModelsTest(ctx context.Context, i *models.LowercaseModelsTestInput) error {
	return nil
}

func (c *ClientPutPagingTest) GetAuthorsWithPut(ctx context.Context, i *models.GetAuthorsWithPutInput) (*models.AuthorsResponse, string, error) {
	assert.Equal(c.t, c.expectedRequestBody, i.FavoriteBooks)
	c.timesPutCalled++
	return &models.AuthorsResponse{
		AuthorSet: &models.AuthorSet{
			Results: models.AuthorArray{
				&models.Author{
					ID:   "123",
					Name: "Mary Shelley",
				},
			},
		},
	}, c.pageToReturn, nil
}
func (c *ClientPutPagingTest) HealthCheck(ctx context.Context) error {
	return nil
}

func TestPutIterator(t *testing.T) {
	controller := ClientPutPagingTest{"", t, nil, 0}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()
	hystrix.Flush()
	c := client.New(testServer.URL, wcl, &http.DefaultTransport)

	requestBody := &models.Book{
		ID:   int64(123),
		Name: "Lord of the Flies",
	}

	controller.expectedRequestBody = requestBody

	iter, err := c.NewGetAuthorsWithPutIter(context.Background(), &models.GetAuthorsWithPutInput{
		FavoriteBooks: requestBody,
	})
	require.NoError(t, err)

	var author models.Author

	// Normally iter.Next would be called in a loop but it's easier to do it this
	// way for testing.
	// Additional assertions on the request body happen in the mock handler.
	controller.pageToReturn = "nextID"
	ok := iter.Next(&author)
	require.True(t, ok)

	controller.pageToReturn = ""
	ok = iter.Next(&author)
	require.True(t, ok)

	ok = iter.Next(&author)
	assert.False(t, ok)
	assert.NoError(t, iter.Err())
	assert.Equal(t, 2, controller.timesPutCalled)
}

func TestExponentialClientRetries(t *testing.T) {
	controller := ClientContextTest{getErrorCount: 2}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()
	c := client.New(testServer.URL, wcl, &http.DefaultTransport)
	c.SetRetryPolicy(client.ExponentialRetryPolicy{})
	_, err := c.GetBooks(context.Background(), &models.GetBooksInput{})
	require.NoError(t, err)
	require.Equal(t, len(controller.getTimes), 3, "expected three requests")
	assert.WithinDuration(t, controller.getTimes[1], controller.getTimes[0].Add(100*time.Millisecond), 20*time.Millisecond,
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

	c := client.New(testServer.URL, wcl, &http.DefaultTransport)
	c.SetRetryPolicy(client.NoRetryPolicy{})
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
	c := client.New(testServer.URL, wcl, &http.DefaultTransport)
	_, err := c.GetBooks(client.WithRetryPolicy(context.Background(), client.NoRetryPolicy{}), &models.GetBooksInput{})
	assert.Error(t, err)
	assert.Equal(t, 1, controller.getCount)
}

func TestNonGetRetries(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()
	c := client.New(testServer.URL, wcl, &http.DefaultTransport)
	_, err := c.CreateBook(context.Background(), &models.Book{})
	assert.Error(t, err)
	assert.Equal(t, 1, controller.postCount)
}

func TestPutRetries(t *testing.T) {
	controller := ClientContextTest{}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()
	c := client.New(testServer.URL, wcl, &http.DefaultTransport)
	returnedBook, err := c.PutBook(context.Background(), &models.Book{Name: "test"})
	require.NoError(t, err)
	assert.Equal(t, 2, controller.putCount)
	assert.Equal(t, "test", returnedBook.Name)
}

func TestErrorOnMissingPathParams(t *testing.T) {
	// Should fail client side
	c := client.New("badUrl", wcl, &http.DefaultTransport)
	_, err := c.GetBookByID2(context.Background(), "")
	require.Error(t, err)
	assert.Equal(t, "id cannot be empty because it's a path parameter", err.Error())
}

func TestNetworkErrorRetries(t *testing.T) {
	c := client.New("https://thisshouldnotresolve1234567890.com/", wcl, &http.DefaultTransport)
	_, err := c.CreateBook(context.Background(), &models.Book{})
	assert.Error(t, err)
}

func TestNewWithDiscovery(t *testing.T) {
	controller := ClientContextTest{getErrorCount: 1}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)

	// Should be an err if env vars aren't set
	_, err := client.NewFromDiscovery(wcl, &http.DefaultTransport)
	assert.Error(t, err)

	splitURL := strings.Split(testServer.URL, ":")
	assert.Equal(t, 3, len(splitURL))

	os.Setenv("SERVICE_SWAGGER_TEST_DEFAULT_PROTO", "http")
	os.Setenv("SERVICE_SWAGGER_TEST_DEFAULT_PORT", splitURL[2])
	os.Setenv("SERVICE_SWAGGER_TEST_DEFAULT_HOST", splitURL[1][2:])

	c, err := client.NewFromDiscovery(wcl, &http.DefaultTransport)

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

	c, err = client.NewFromDiscovery(wcl, &http.DefaultTransport)
	assert.NoError(t, err)
	_, err = c.GetBooks(context.Background(), &models.GetBooksInput{})
	assert.NoError(t, err)
	assert.Equal(t, 3, controller.getCount)
}

func TestIterator(t *testing.T) {
	// we have 2 pages and 3 books so that we have to request a new page midway
	// through and we need to loop through more than one item on a single page
	s, controller := setupServer()
	controller.pageSize = 2
	c := client.New(s.URL, wcl, &http.DefaultTransport)

	book1ID := int64(1)
	book1Name := "First"
	book2ID := int64(2)
	book2Name := "Second"
	book3ID := int64(3)
	book3Name := "Third"
	_, err := c.CreateBook(context.Background(), &models.Book{
		ID: book1ID, Name: book1Name,
	})
	require.NoError(t, err)
	_, err = c.CreateBook(context.Background(), &models.Book{
		ID: book2ID, Name: book2Name,
	})
	require.NoError(t, err)
	_, err = c.CreateBook(context.Background(), &models.Book{
		ID: book3ID, Name: book3Name,
	})
	require.NoError(t, err)

	iter, err := c.NewGetBooksIter(context.Background(), &models.GetBooksInput{})
	require.NoError(t, err)

	var book models.Book

	// normally iter.Next would be called in a loop but it's easier to do it this
	// way for testing
	ok := iter.Next(&book)
	require.True(t, ok)
	assert.Equal(t, book1ID, book.ID)
	assert.Equal(t, book1Name, book.Name)

	ok = iter.Next(&book)
	require.True(t, ok)
	assert.Equal(t, book2ID, book.ID)
	assert.Equal(t, book2Name, book.Name)

	ok = iter.Next(&book)
	require.True(t, ok)
	assert.Equal(t, book3ID, book.ID)
	assert.Equal(t, book3Name, book.Name)

	ok = iter.Next(&book)
	assert.False(t, ok)
	assert.NoError(t, iter.Err())
}

// TestIteratorWithResourcePath makes sure the client works when
// x-paging.resourcePath is specified
func TestIteratorWithResourcePath(t *testing.T) {
	s, controller := setupServer()
	controller.authors = []*models.Author{
		&models.Author{
			ID:   "abc",
			Name: "Joe",
		},
		&models.Author{
			ID:   "def",
			Name: "Jenny",
		},
	}
	c := client.New(s.URL, wcl, &http.DefaultTransport)

	iter, err := c.NewGetAuthorsIter(context.Background(), &models.GetAuthorsInput{})
	require.NoError(t, err)

	var author models.Author

	// normally iter.Next would be called in a loop but it's easier to do it this
	// way for testing
	ok := iter.Next(&author)
	require.True(t, ok)
	assert.Equal(t, "abc", author.ID)
	assert.Equal(t, "Joe", author.Name)

	ok = iter.Next(&author)
	require.True(t, ok)
	assert.Equal(t, "def", author.ID)
	assert.Equal(t, "Jenny", author.Name)

	ok = iter.Next(&author)
	assert.False(t, ok)
	assert.NoError(t, iter.Err())
}

type IterFailTest struct {
	sampleController *ControllerImpl
	fail             bool
}

func (c *IterFailTest) GetBooks(ctx context.Context, input *models.GetBooksInput) ([]models.Book, int64, error) {
	if c.fail {
		return nil, int64(0), errors.New("fail")
	}
	return c.sampleController.GetBooks(ctx, input)
}
func (c *IterFailTest) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (*models.Book, error) {
	return c.sampleController.GetBookByID(ctx, input)
}
func (c *IterFailTest) GetBookByID2(ctx context.Context, id string) (*models.Book, error) {
	return c.sampleController.GetBookByID2(ctx, id)
}
func (c *IterFailTest) CreateBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	return c.sampleController.CreateBook(ctx, input)
}
func (c *IterFailTest) PutBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	return c.sampleController.CreateBook(ctx, input)
}
func (c *IterFailTest) GetAuthors(ctx context.Context, input *models.GetAuthorsInput) (*models.AuthorsResponse, string, error) {
	return c.sampleController.GetAuthors(ctx, input)
}
func (c *IterFailTest) GetAuthorsWithPut(ctx context.Context, i *models.GetAuthorsWithPutInput) (*models.AuthorsResponse, string, error) {
	return nil, "", nil
}
func (c *IterFailTest) HealthCheck(ctx context.Context) error {
	return nil
}

func (c *IterFailTest) LowercaseModelsTest(ctx context.Context, i *models.LowercaseModelsTestInput) error {
	return nil
}

func TestIteratorFail(t *testing.T) {
	controller := IterFailTest{sampleController: &ControllerImpl{
		books:    make(map[int64]*models.Book),
		pageSize: 1,
	}}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()
	c := client.New(testServer.URL, wcl, &http.DefaultTransport)

	book1ID := int64(1)
	book1Name := "Test"
	book2ID := int64(2)
	book2Name := "Second"
	_, err := c.CreateBook(context.Background(), &models.Book{
		ID: book1ID, Name: book1Name,
	})
	require.NoError(t, err)
	_, err = c.CreateBook(context.Background(), &models.Book{
		ID: book2ID, Name: book2Name,
	})
	require.NoError(t, err)

	iter, err := c.NewGetBooksIter(context.Background(), &models.GetBooksInput{})
	require.NoError(t, err)

	var book models.Book

	// normally iter.Next would be called in a loop but it's easier to do it this
	// way for testing
	ok := iter.Next(&book)
	require.True(t, ok)
	assert.Equal(t, book1ID, book.ID)
	assert.Equal(t, book1Name, book.Name)

	controller.fail = true

	ok = iter.Next(&book)
	assert.False(t, ok)
	require.Error(t, iter.Err())
	assert.IsType(t, &models.InternalError{}, iter.Err())
	assert.Equal(t, "fail", iter.Err().Error())
}

type IterHeadersTest struct {
	sampleController *ControllerImpl
	t                *testing.T
}

func (c *IterHeadersTest) GetBooks(ctx context.Context, input *models.GetBooksInput) ([]models.Book, int64, error) {
	assert.Equal(c.t, "x-let-me-in-bro", input.Authorization)
	return c.sampleController.GetBooks(ctx, input)
}
func (c *IterHeadersTest) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (*models.Book, error) {
	return c.sampleController.GetBookByID(ctx, input)
}
func (c *IterHeadersTest) GetBookByID2(ctx context.Context, id string) (*models.Book, error) {
	return c.sampleController.GetBookByID2(ctx, id)
}
func (c *IterHeadersTest) CreateBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	return c.sampleController.CreateBook(ctx, input)
}
func (c *IterHeadersTest) PutBook(ctx context.Context, input *models.Book) (*models.Book, error) {
	return c.sampleController.CreateBook(ctx, input)
}
func (c *IterHeadersTest) GetAuthors(ctx context.Context, input *models.GetAuthorsInput) (*models.AuthorsResponse, string, error) {
	return c.sampleController.GetAuthors(ctx, input)
}
func (c *IterHeadersTest) GetAuthorsWithPut(ctx context.Context, i *models.GetAuthorsWithPutInput) (*models.AuthorsResponse, string, error) {
	return nil, "", nil
}
func (c *IterHeadersTest) HealthCheck(ctx context.Context) error {
	return nil
}

func (c *IterHeadersTest) LowercaseModelsTest(ctx context.Context, i *models.LowercaseModelsTestInput) error {
	return nil
}

func TestIteratorHeaders(t *testing.T) {
	controller := IterHeadersTest{
		t:                t,
		sampleController: &ControllerImpl{pageSize: 1, books: make(map[int64]*models.Book)},
	}
	s := server.New(&controller, "")
	testServer := httptest.NewServer(s.Handler)
	defer testServer.Close()
	c := client.New(testServer.URL, wcl, &http.DefaultTransport)

	t.Log("Ensure client.SetLogger works")
	c.SetLogger(wcl)

	book1ID := int64(1)
	book1Name := "Test"
	book2ID := int64(2)
	book2Name := "Second"
	_, err := c.CreateBook(context.Background(), &models.Book{
		ID: book1ID, Name: book1Name,
	})
	require.NoError(t, err)
	_, err = c.CreateBook(context.Background(), &models.Book{
		ID: book2ID, Name: book2Name,
	})
	require.NoError(t, err)

	iter, err := c.NewGetBooksIter(context.Background(), &models.GetBooksInput{
		Authorization: "x-let-me-in-bro",
	})
	require.NoError(t, err)

	count := 0
	var book models.Book
	for iter.Next(&book) {
		count++
	}
	assert.NoError(t, iter.Err())
	assert.Equal(t, 2, count)
}

func TestVersionHeader(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(client.VersionHeader) != client.Version {
			t.Error("did not receive version header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()
	c := client.New(testServer.URL, wcl, &http.DefaultTransport)
	c.HealthCheck(context.Background())
}

func TestUnknownResponseBody(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(420)
		w.Write([]byte(`{"enhance": "your calm"}`))
	}))
	defer testServer.Close()
	c := client.New(testServer.URL, wcl, &http.DefaultTransport)
	_, err := c.GetBookByID(context.Background(), &models.GetBookByIDInput{BookID: 420})
	assert.Error(t, err)
	assert.IsType(t, models.UnknownResponse{}, err)
	assert.Equal(t, err.Error(), `unknown response with status: 420 body: {"enhance": "your calm"}`)
}
