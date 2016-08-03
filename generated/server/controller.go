package server

import (
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"errors"
	"github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"
)

type ControllerImpl struct{
}
func NewGetBooksInput(r *http.Request) (*models.GetBooksInput, error) {
	var input models.GetBooksInput

	authorStr := r.URL.Query().Get("author")
	input.Author = authorStr
	availableStr := r.URL.Query().Get("available")
	var err error
	input.Available, err = strconv.ParseBool(availableStr)
	if err != nil {
		return nil, err
	}
	maxPagesStr := r.URL.Query().Get("maxPages")
	input.MaxPages, err = strconv.ParseFloat(maxPagesStr, 64)
	if err != nil {
		return nil, err
	}

	return &input, nil
}

func (c ControllerImpl) GetBooks(ctx context.Context, input *models.GetBooksInput) (models.GetBooksOutput, error) {
	// TODO: Implement me!
	return nil, errors.New("Not implemented")
}
func NewGetBookByIDInput(r *http.Request) (*models.GetBookByIDInput, error) {
	var input models.GetBookByIDInput

	bookIDStr := mux.Vars(r)["bookID"]
	var err error
	input.BookID, err = strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil {
		return nil, err
	}
	authorizationStr := r.Header.Get("authorization")
	input.Authorization = authorizationStr

	return &input, nil
}

func (c ControllerImpl) GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (models.GetBookByIDOutput, error) {
	// TODO: Implement me!
	return nil, errors.New("Not implemented")
}
func NewCreateBookInput(r *http.Request) (*models.CreateBookInput, error) {
	var input models.CreateBookInput


	return &input, nil
}

func (c ControllerImpl) CreateBook(ctx context.Context, input *models.CreateBookInput) (models.CreateBookOutput, error) {
	// TODO: Implement me!
	return nil, errors.New("Not implemented")
}
