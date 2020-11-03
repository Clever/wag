package db

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-app-service/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_db.go -package=db

// Interface for interacting with the app-service database.
type Interface interface {
	// SaveSetupStep saves a SetupStep to the database.
	SaveSetupStep(ctx context.Context, m models.SetupStep) error
	// GetSetupStep retrieves a SetupStep from the database.
	GetSetupStep(ctx context.Context, appID string, id string) (*models.SetupStep, error)
	// GetSetupStepsByAppIDAndID retrieves a page of SetupSteps from the database.
	GetSetupStepsByAppIDAndID(ctx context.Context, input GetSetupStepsByAppIDAndIDInput, fn func(m *models.SetupStep, lastSetupStep bool) bool) error
	// DeleteSetupStep deletes a SetupStep from the database.
	DeleteSetupStep(ctx context.Context, appID string, id string) error
}

// Int64 returns a pointer to the int64 value passed in.
func Int64(i int64) *int64 { return &i }

// String returns a pointer to the string value passed in.
func String(s string) *string { return &s }

// GetSetupStepsByAppIDAndIDInput is the query input to GetSetupStepsByAppIDAndID.
type GetSetupStepsByAppIDAndIDInput struct {
	// AppID is required
	AppID        string
	IDStartingAt *string
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.SetupStep
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int64
}

// ErrSetupStepNotFound is returned when the database fails to find a SetupStep.
type ErrSetupStepNotFound struct {
	AppID string
	ID    string
}

var _ error = ErrSetupStepNotFound{}

// Error returns a description of the error.
func (e ErrSetupStepNotFound) Error() string {
	return "could not find SetupStep"
}
