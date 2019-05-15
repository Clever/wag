package server

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-blog/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_controller.go -package=server

// Controller defines the interface for the blog service.
type Controller interface {

	// GetSectionsForStudent handles GET requests to /students/{student_id}/sections
	// Gets the sections for the specified student
	// 200: []models.Section
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetSectionsForStudent(ctx context.Context, studentID string) ([]models.Section, error)
}
