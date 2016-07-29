package main

import "golang.org/x/net/context"
// import "errors"
import "github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"
import "fmt"

type ControllerImpl struct{
}

// This is used for easily testing this curl request `curl -H "authorization:test" localhost:8080/books/1234?author=kyle`
func (c ControllerImpl) GetBookByID(ctx context.Context, input *GetBookByIDInput) (GetBookByIDOutput, error) {

        fmt.Printf("Context: %s\n", ctx.Value("addedValue"))

        fmt.Printf("Author: %s\n", input.Author)
        fmt.Printf("BookID: %s\n", input.BookID)
        fmt.Printf("Authorization: %s\n", input.Authorization)

	// TODO: Implement me!
	// return nil, errors.New("Not implemented")
        return &GetBookByID200Output{models.Book{Name: "Test"}}, nil
}
