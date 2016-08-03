package models

import (
	"encoding/json"
	"strconv"
)

var _ = json.Marshal

var _ = strconv.FormatInt

type GetBookByIDInput struct {
	BookID int64
	Authorization string
}

func (i GetBookByIDInput) Validate() error{
	return nil
}

type CreateBookInput struct {
	NewBook Book
}

func (i CreateBookInput) Validate() error{
	if err := i.NewBook.Validate(nil); err != nil {
		return err
	}

	return nil
}

type GetBooksInput struct {
	Author string
	Available bool
	MaxPages float64
}

func (i GetBooksInput) Validate() error{
	return nil
}

