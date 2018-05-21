package db

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-db/models"
	"github.com/go-openapi/strfmt"
)

//go:generate $GOPATH/bin/mockgen -source=$GOFILE -destination=mock_db.go -package=db

// Interface for interacting with the swagger-test database.
type Interface interface {
	// SaveSimpleThing saves a SimpleThing to the database.
	SaveSimpleThing(ctx context.Context, m models.SimpleThing) error
	// GetSimpleThing retrieves a SimpleThing from the database.
	GetSimpleThing(ctx context.Context, name string) (*models.SimpleThing, error)
	// DeleteSimpleThing deletes a SimpleThing from the database.
	DeleteSimpleThing(ctx context.Context, name string) error

	// SaveThing saves a Thing to the database.
	SaveThing(ctx context.Context, m models.Thing) error
	// GetThing retrieves a Thing from the database.
	GetThing(ctx context.Context, name string, version int64) (*models.Thing, error)
	// GetThingsByNameAndVersion retrieves a list of Things from the database.
	GetThingsByNameAndVersion(ctx context.Context, input GetThingsByNameAndVersionInput) ([]models.Thing, error)
	// DeleteThing deletes a Thing from the database.
	DeleteThing(ctx context.Context, name string, version int64) error
	// GetThingByID retrieves a Thing from the database.
	GetThingByID(ctx context.Context, id string) (*models.Thing, error)
	// GetThingsByNameAndCreatedAt retrieves a list of Things from the database.
	GetThingsByNameAndCreatedAt(ctx context.Context, input GetThingsByNameAndCreatedAtInput) ([]models.Thing, error)

	// SaveThingWithDateRange saves a ThingWithDateRange to the database.
	SaveThingWithDateRange(ctx context.Context, m models.ThingWithDateRange) error
	// GetThingWithDateRange retrieves a ThingWithDateRange from the database.
	GetThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) (*models.ThingWithDateRange, error)
	// GetThingWithDateRangesByNameAndDate retrieves a list of ThingWithDateRanges from the database.
	GetThingWithDateRangesByNameAndDate(ctx context.Context, input GetThingWithDateRangesByNameAndDateInput) ([]models.ThingWithDateRange, error)
	// DeleteThingWithDateRange deletes a ThingWithDateRange from the database.
	DeleteThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) error

	// SaveThingWithUnderscores saves a ThingWithUnderscores to the database.
	SaveThingWithUnderscores(ctx context.Context, m models.ThingWithUnderscores) error
	// GetThingWithUnderscores retrieves a ThingWithUnderscores from the database.
	GetThingWithUnderscores(ctx context.Context, idApp string) (*models.ThingWithUnderscores, error)
	// DeleteThingWithUnderscores deletes a ThingWithUnderscores from the database.
	DeleteThingWithUnderscores(ctx context.Context, idApp string) error
}

// Int64 returns a pointer to the int64 value passed in.
func Int64(i int64) *int64 { return &i }

// String returns a pointer to the string value passed in.
func String(s string) *string { return &s }

// DateTime returns a pointer to the strfmt.DateTime value passed in.
func DateTime(d strfmt.DateTime) *strfmt.DateTime { return &d }

// ErrSimpleThingNotFound is returned when the database fails to find a SimpleThing.
type ErrSimpleThingNotFound struct {
	Name string
}

var _ error = ErrSimpleThingNotFound{}

// Error returns a description of the error.
func (e ErrSimpleThingNotFound) Error() string {
	return "could not find SimpleThing"
}

// ErrSimpleThingAlreadyExists is returned when trying to overwrite a SimpleThing.
type ErrSimpleThingAlreadyExists struct {
	Name string
}

var _ error = ErrSimpleThingAlreadyExists{}

// Error returns a description of the error.
func (e ErrSimpleThingAlreadyExists) Error() string {
	return "SimpleThing already exists"
}

// GetThingsByNameAndVersionInput is the query input to GetThingsByNameAndVersion.
type GetThingsByNameAndVersionInput struct {
	Name                  string
	VersionStartingAt     *int64
	Descending            bool
	DisableConsistentRead bool
}

// ErrThingNotFound is returned when the database fails to find a Thing.
type ErrThingNotFound struct {
	Name    string
	Version int64
}

var _ error = ErrThingNotFound{}

// Error returns a description of the error.
func (e ErrThingNotFound) Error() string {
	return "could not find Thing"
}

// ErrThingByIDNotFound is returned when the database fails to find a Thing.
type ErrThingByIDNotFound struct {
	ID string
}

var _ error = ErrThingByIDNotFound{}

// Error returns a description of the error.
func (e ErrThingByIDNotFound) Error() string {
	return "could not find Thing"
}

// GetThingsByNameAndCreatedAtInput is the query input to GetThingsByNameAndCreatedAt.
type GetThingsByNameAndCreatedAtInput struct {
	Name                string
	CreatedAtStartingAt *strfmt.DateTime
	Descending          bool
}

// ErrThingByNameAndCreatedAtNotFound is returned when the database fails to find a Thing.
type ErrThingByNameAndCreatedAtNotFound struct {
	Name      string
	CreatedAt strfmt.DateTime
}

var _ error = ErrThingByNameAndCreatedAtNotFound{}

// Error returns a description of the error.
func (e ErrThingByNameAndCreatedAtNotFound) Error() string {
	return "could not find Thing"
}

// ErrThingAlreadyExists is returned when trying to overwrite a Thing.
type ErrThingAlreadyExists struct {
	Name    string
	Version int64
}

var _ error = ErrThingAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingAlreadyExists) Error() string {
	return "Thing already exists"
}

// GetThingWithDateRangesByNameAndDateInput is the query input to GetThingWithDateRangesByNameAndDate.
type GetThingWithDateRangesByNameAndDateInput struct {
	Name                  string
	DateStartingAt        *strfmt.DateTime
	Descending            bool
	DisableConsistentRead bool
}

// ErrThingWithDateRangeNotFound is returned when the database fails to find a ThingWithDateRange.
type ErrThingWithDateRangeNotFound struct {
	Name string
	Date strfmt.DateTime
}

var _ error = ErrThingWithDateRangeNotFound{}

// Error returns a description of the error.
func (e ErrThingWithDateRangeNotFound) Error() string {
	return "could not find ThingWithDateRange"
}

// ErrThingWithUnderscoresNotFound is returned when the database fails to find a ThingWithUnderscores.
type ErrThingWithUnderscoresNotFound struct {
	IDApp string
}

var _ error = ErrThingWithUnderscoresNotFound{}

// Error returns a description of the error.
func (e ErrThingWithUnderscoresNotFound) Error() string {
	return "could not find ThingWithUnderscores"
}
