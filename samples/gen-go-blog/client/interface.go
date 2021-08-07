package client

import (
	"context"

	"github.com/Clever/wag/samples/v8/gen-go-blog/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package=client

// Client defines the methods available to clients of the blog service.
type Client interface {

	// PostGradeFileForStudent makes a POST request to /students/{student_id}/gradeFile
	// Posts the grade file for the specified student
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	PostGradeFileForStudent(ctx context.Context, i *models.PostGradeFileForStudentInput) error

	// GetSectionsForStudent makes a GET request to /students/{student_id}/sections
	// Gets the sections for the specified student
	// 200: []models.Section
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetSectionsForStudent(ctx context.Context, studentID string) ([]models.Section, error)

	// PostSectionsForStudent makes a POST request to /students/{student_id}/sections
	// Posts the sections for the specified student
	// 200: []models.Section
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	PostSectionsForStudent(ctx context.Context, i *models.PostSectionsForStudentInput) ([]models.Section, error)
}
