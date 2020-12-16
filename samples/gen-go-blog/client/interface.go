package client

import (
	"context"

	"github.com/Clever/wag/v5/samples/gen-go-blog/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package=client

// Client defines the methods available to clients of the blog service.
type Client interface {

	// GetSectionsForStudent makes a GET request to /students/{student_id}/sections
	// Gets the sections for the specified student
	// 200: []models.Section
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetSectionsForStudent(ctx context.Context, studentID string) ([]models.Section, error)
}
