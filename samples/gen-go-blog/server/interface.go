package server

import (
	"context"

	"github.com/Clever/wag/samples/v9/gen-go-blog/gen-go/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_controller.go -package=server

// Controller defines the interface for the blog service.
type Controller interface {

	// PostGradeFileForStudent handles POST requests to /students/{student_id}/gradeFile
	// Posts the grade file for the specified student
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	PostGradeFileForStudent(ctx context.Context, i *models.PostGradeFileForStudentInput) error

	// GetSectionsForStudent handles GET requests to /students/{student_id}/sections
	// Gets the sections for the specified student
	// 200: []models.Section
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetSectionsForStudent(ctx context.Context, studentID string) ([]models.Section, error)

	// PostSectionsForStudent handles POST requests to /students/{student_id}/sections
	// Posts the sections for the specified student
	// 200: []models.Section
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	PostSectionsForStudent(ctx context.Context, i *models.PostSectionsForStudentInput) ([]models.Section, error)
}
