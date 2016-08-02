package generated

import "golang.org/x/net/context"
import "github.com/Clever/inter-service-api-testing/codegen-poc/generated/models"

type Controller interface {
	GetBookByID(ctx context.Context, input *models.GetBookByIDInput) (models.GetBookByIDOutput, error)
	CreateBook(ctx context.Context, input *models.CreateBookInput) (models.CreateBookOutput, error)
	GetBooks(ctx context.Context, input *models.GetBooksInput) (models.GetBooksOutput, error)
}
