package db

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-db/models"
	"github.com/go-openapi/strfmt"
	"golang.org/x/time/rate"
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

	// SaveTeacherSharingRule saves a TeacherSharingRule to the database.
	SaveTeacherSharingRule(ctx context.Context, m models.TeacherSharingRule) error
	// GetTeacherSharingRule retrieves a TeacherSharingRule from the database.
	GetTeacherSharingRule(ctx context.Context, teacher string, school string, app string) (*models.TeacherSharingRule, error)
	// GetTeacherSharingRulesByTeacherAndSchoolApp retrieves a page of TeacherSharingRules from the database.
	GetTeacherSharingRulesByTeacherAndSchoolApp(ctx context.Context, input GetTeacherSharingRulesByTeacherAndSchoolAppInput, fn func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool) error
	// DeleteTeacherSharingRule deletes a TeacherSharingRule from the database.
	DeleteTeacherSharingRule(ctx context.Context, teacher string, school string, app string) error
	// GetTeacherSharingRulesByDistrictAndSchoolTeacherApp retrieves a page of TeacherSharingRules from the database.
	GetTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx context.Context, input GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput, fn func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool) error

	// SaveThing saves a Thing to the database.
	SaveThing(ctx context.Context, m models.Thing) error
	// GetThing retrieves a Thing from the database.
	GetThing(ctx context.Context, name string, version int64) (*models.Thing, error)
	// ScanThings runs a scan on the Things table.
	ScanThings(ctx context.Context, input ScanThingsInput, fn func(m *models.Thing, lastThing bool) bool) error
	// GetThingsByNameAndVersion retrieves a page of Things from the database.
	GetThingsByNameAndVersion(ctx context.Context, input GetThingsByNameAndVersionInput, fn func(m *models.Thing, lastThing bool) bool) error
	// DeleteThing deletes a Thing from the database.
	DeleteThing(ctx context.Context, name string, version int64) error
	// GetThingByID retrieves a Thing from the database.
	GetThingByID(ctx context.Context, id string) (*models.Thing, error)
	// GetThingsByNameAndCreatedAt retrieves a page of Things from the database.
	GetThingsByNameAndCreatedAt(ctx context.Context, input GetThingsByNameAndCreatedAtInput, fn func(m *models.Thing, lastThing bool) bool) error

	// SaveThingWithCompositeAttributes saves a ThingWithCompositeAttributes to the database.
	SaveThingWithCompositeAttributes(ctx context.Context, m models.ThingWithCompositeAttributes) error
	// GetThingWithCompositeAttributes retrieves a ThingWithCompositeAttributes from the database.
	GetThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) (*models.ThingWithCompositeAttributes, error)
	// GetThingWithCompositeAttributessByNameBranchAndDate retrieves a page of ThingWithCompositeAttributess from the database.
	GetThingWithCompositeAttributessByNameBranchAndDate(ctx context.Context, input GetThingWithCompositeAttributessByNameBranchAndDateInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error
	// DeleteThingWithCompositeAttributes deletes a ThingWithCompositeAttributes from the database.
	DeleteThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) error
	// GetThingWithCompositeAttributessByNameVersionAndDate retrieves a page of ThingWithCompositeAttributess from the database.
	GetThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input GetThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error

	// SaveThingWithCompositeEnumAttributes saves a ThingWithCompositeEnumAttributes to the database.
	SaveThingWithCompositeEnumAttributes(ctx context.Context, m models.ThingWithCompositeEnumAttributes) error
	// GetThingWithCompositeEnumAttributes retrieves a ThingWithCompositeEnumAttributes from the database.
	GetThingWithCompositeEnumAttributes(ctx context.Context, name string, branchID models.Branch, date strfmt.DateTime) (*models.ThingWithCompositeEnumAttributes, error)
	// GetThingWithCompositeEnumAttributessByNameBranchAndDate retrieves a page of ThingWithCompositeEnumAttributess from the database.
	GetThingWithCompositeEnumAttributessByNameBranchAndDate(ctx context.Context, input GetThingWithCompositeEnumAttributessByNameBranchAndDateInput, fn func(m *models.ThingWithCompositeEnumAttributes, lastThingWithCompositeEnumAttributes bool) bool) error
	// DeleteThingWithCompositeEnumAttributes deletes a ThingWithCompositeEnumAttributes from the database.
	DeleteThingWithCompositeEnumAttributes(ctx context.Context, name string, branchID models.Branch, date strfmt.DateTime) error

	// SaveThingWithDateRange saves a ThingWithDateRange to the database.
	SaveThingWithDateRange(ctx context.Context, m models.ThingWithDateRange) error
	// GetThingWithDateRange retrieves a ThingWithDateRange from the database.
	GetThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) (*models.ThingWithDateRange, error)
	// GetThingWithDateRangesByNameAndDate retrieves a page of ThingWithDateRanges from the database.
	GetThingWithDateRangesByNameAndDate(ctx context.Context, input GetThingWithDateRangesByNameAndDateInput, fn func(m *models.ThingWithDateRange, lastThingWithDateRange bool) bool) error
	// DeleteThingWithDateRange deletes a ThingWithDateRange from the database.
	DeleteThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) error

	// SaveThingWithDateTimeComposite saves a ThingWithDateTimeComposite to the database.
	SaveThingWithDateTimeComposite(ctx context.Context, m models.ThingWithDateTimeComposite) error
	// GetThingWithDateTimeComposite retrieves a ThingWithDateTimeComposite from the database.
	GetThingWithDateTimeComposite(ctx context.Context, typeVar string, id string, created strfmt.DateTime, resource string) (*models.ThingWithDateTimeComposite, error)
	// GetThingWithDateTimeCompositesByTypeIDAndCreatedResource retrieves a page of ThingWithDateTimeComposites from the database.
	GetThingWithDateTimeCompositesByTypeIDAndCreatedResource(ctx context.Context, input GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput, fn func(m *models.ThingWithDateTimeComposite, lastThingWithDateTimeComposite bool) bool) error
	// DeleteThingWithDateTimeComposite deletes a ThingWithDateTimeComposite from the database.
	DeleteThingWithDateTimeComposite(ctx context.Context, typeVar string, id string, created strfmt.DateTime, resource string) error

	// SaveThingWithRequiredFields saves a ThingWithRequiredFields to the database.
	SaveThingWithRequiredFields(ctx context.Context, m models.ThingWithRequiredFields) error
	// GetThingWithRequiredFields retrieves a ThingWithRequiredFields from the database.
	GetThingWithRequiredFields(ctx context.Context, name string) (*models.ThingWithRequiredFields, error)
	// DeleteThingWithRequiredFields deletes a ThingWithRequiredFields from the database.
	DeleteThingWithRequiredFields(ctx context.Context, name string) error

	// SaveThingWithUnderscores saves a ThingWithUnderscores to the database.
	SaveThingWithUnderscores(ctx context.Context, m models.ThingWithUnderscores) error
	// GetThingWithUnderscores retrieves a ThingWithUnderscores from the database.
	GetThingWithUnderscores(ctx context.Context, iDApp string) (*models.ThingWithUnderscores, error)
	// DeleteThingWithUnderscores deletes a ThingWithUnderscores from the database.
	DeleteThingWithUnderscores(ctx context.Context, iDApp string) error
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

// GetTeacherSharingRulesByTeacherAndSchoolAppInput is the query input to GetTeacherSharingRulesByTeacherAndSchoolApp.
type GetTeacherSharingRulesByTeacherAndSchoolAppInput struct {
	// StartingAt is a required specification of an (exclusive) starting point.
	StartingAt *models.TeacherSharingRule
	// Exclusive toggles whether results include the start point
	Exclusive bool
	// Limit is a required limit on how many items to evaluate.
	Limit *int64
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	Descending            bool
}

// SchoolApp struct.
type SchoolApp struct {
	School string
	App    string
}

// ErrTeacherSharingRuleNotFound is returned when the database fails to find a TeacherSharingRule.
type ErrTeacherSharingRuleNotFound struct {
	Teacher string
	School  string
	App     string
}

var _ error = ErrTeacherSharingRuleNotFound{}

// Error returns a description of the error.
func (e ErrTeacherSharingRuleNotFound) Error() string {
	return "could not find TeacherSharingRule"
}

// GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput is the query input to GetTeacherSharingRulesByDistrictAndSchoolTeacherApp.
type GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput struct {
	// StartingAt is a required specification of an (exclusive) starting point.
	StartingAt *models.TeacherSharingRule
	// Exclusive toggles whether results include the start point
	Exclusive bool
	// Limit is a required limit on how many items to evaluate.
	Limit      *int64
	Descending bool
}

// SchoolTeacherApp struct.
type SchoolTeacherApp struct {
	School  string
	Teacher string
	App     string
}

// ErrTeacherSharingRuleByDistrictAndSchoolTeacherAppNotFound is returned when the database fails to find a TeacherSharingRule.
type ErrTeacherSharingRuleByDistrictAndSchoolTeacherAppNotFound struct {
	District string
	School   string
	Teacher  string
	App      string
}

var _ error = ErrTeacherSharingRuleByDistrictAndSchoolTeacherAppNotFound{}

// Error returns a description of the error.
func (e ErrTeacherSharingRuleByDistrictAndSchoolTeacherAppNotFound) Error() string {
	return "could not find TeacherSharingRule"
}

// ScanThingsInput is the input to the ScanThings method.
type ScanThingsInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.Thing
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingsByNameAndVersionInput is the query input to GetThingsByNameAndVersion.
type GetThingsByNameAndVersionInput struct {
	// StartingAt is a required specification of an (exclusive) starting point.
	StartingAt *models.Thing
	// Exclusive toggles whether results include the start point
	Exclusive bool
	// Limit is a required limit on how many items to evaluate.
	Limit *int64
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	Descending            bool
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
	// StartingAt is a required specification of an (exclusive) starting point.
	StartingAt *models.Thing
	// Exclusive toggles whether results include the start point
	Exclusive bool
	// Limit is a required limit on how many items to evaluate.
	Limit      *int64
	Descending bool
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

// GetThingWithCompositeAttributessByNameBranchAndDateInput is the query input to GetThingWithCompositeAttributessByNameBranchAndDate.
type GetThingWithCompositeAttributessByNameBranchAndDateInput struct {
	// StartingAt is a required specification of an (exclusive) starting point.
	StartingAt *models.ThingWithCompositeAttributes
	// Exclusive toggles whether results include the start point
	Exclusive bool
	// Limit is a required limit on how many items to evaluate.
	Limit *int64
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	Descending            bool
}

// ErrThingWithCompositeAttributesNotFound is returned when the database fails to find a ThingWithCompositeAttributes.
type ErrThingWithCompositeAttributesNotFound struct {
	Name   string
	Branch string
	Date   strfmt.DateTime
}

var _ error = ErrThingWithCompositeAttributesNotFound{}

// Error returns a description of the error.
func (e ErrThingWithCompositeAttributesNotFound) Error() string {
	return "could not find ThingWithCompositeAttributes"
}

// GetThingWithCompositeAttributessByNameVersionAndDateInput is the query input to GetThingWithCompositeAttributessByNameVersionAndDate.
type GetThingWithCompositeAttributessByNameVersionAndDateInput struct {
	// StartingAt is a required specification of an (exclusive) starting point.
	StartingAt *models.ThingWithCompositeAttributes
	// Exclusive toggles whether results include the start point
	Exclusive bool
	// Limit is a required limit on how many items to evaluate.
	Limit      *int64
	Descending bool
}

// ErrThingWithCompositeAttributesByNameVersionAndDateNotFound is returned when the database fails to find a ThingWithCompositeAttributes.
type ErrThingWithCompositeAttributesByNameVersionAndDateNotFound struct {
	Name    string
	Version int64
	Date    strfmt.DateTime
}

var _ error = ErrThingWithCompositeAttributesByNameVersionAndDateNotFound{}

// Error returns a description of the error.
func (e ErrThingWithCompositeAttributesByNameVersionAndDateNotFound) Error() string {
	return "could not find ThingWithCompositeAttributes"
}

// ErrThingWithCompositeAttributesAlreadyExists is returned when trying to overwrite a ThingWithCompositeAttributes.
type ErrThingWithCompositeAttributesAlreadyExists struct {
	NameBranch string
	Date       strfmt.DateTime
}

var _ error = ErrThingWithCompositeAttributesAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithCompositeAttributesAlreadyExists) Error() string {
	return "ThingWithCompositeAttributes already exists"
}

// GetThingWithCompositeEnumAttributessByNameBranchAndDateInput is the query input to GetThingWithCompositeEnumAttributessByNameBranchAndDate.
type GetThingWithCompositeEnumAttributessByNameBranchAndDateInput struct {
	// StartingAt is a required specification of an (exclusive) starting point.
	StartingAt *models.ThingWithCompositeEnumAttributes
	// Exclusive toggles whether results include the start point
	Exclusive bool
	// Limit is a required limit on how many items to evaluate.
	Limit *int64
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	Descending            bool
}

// ErrThingWithCompositeEnumAttributesNotFound is returned when the database fails to find a ThingWithCompositeEnumAttributes.
type ErrThingWithCompositeEnumAttributesNotFound struct {
	Name     string
	BranchID models.Branch
	Date     strfmt.DateTime
}

var _ error = ErrThingWithCompositeEnumAttributesNotFound{}

// Error returns a description of the error.
func (e ErrThingWithCompositeEnumAttributesNotFound) Error() string {
	return "could not find ThingWithCompositeEnumAttributes"
}

// ErrThingWithCompositeEnumAttributesAlreadyExists is returned when trying to overwrite a ThingWithCompositeEnumAttributes.
type ErrThingWithCompositeEnumAttributesAlreadyExists struct {
	NameBranch string
	Date       strfmt.DateTime
}

var _ error = ErrThingWithCompositeEnumAttributesAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithCompositeEnumAttributesAlreadyExists) Error() string {
	return "ThingWithCompositeEnumAttributes already exists"
}

// GetThingWithDateRangesByNameAndDateInput is the query input to GetThingWithDateRangesByNameAndDate.
type GetThingWithDateRangesByNameAndDateInput struct {
	// StartingAt is a required specification of an (exclusive) starting point.
	StartingAt *models.ThingWithDateRange
	// Exclusive toggles whether results include the start point
	Exclusive bool
	// Limit is a required limit on how many items to evaluate.
	Limit *int64
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	Descending            bool
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

// GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput is the query input to GetThingWithDateTimeCompositesByTypeIDAndCreatedResource.
type GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput struct {
	// StartingAt is a required specification of an (exclusive) starting point.
	StartingAt *models.ThingWithDateTimeComposite
	// Exclusive toggles whether results include the start point
	Exclusive bool
	// Limit is a required limit on how many items to evaluate.
	Limit *int64
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	Descending            bool
}

// CreatedResource struct.
type CreatedResource struct {
	Created  strfmt.DateTime
	Resource string
}

// ErrThingWithDateTimeCompositeNotFound is returned when the database fails to find a ThingWithDateTimeComposite.
type ErrThingWithDateTimeCompositeNotFound struct {
	Type     string
	ID       string
	Created  strfmt.DateTime
	Resource string
}

var _ error = ErrThingWithDateTimeCompositeNotFound{}

// Error returns a description of the error.
func (e ErrThingWithDateTimeCompositeNotFound) Error() string {
	return "could not find ThingWithDateTimeComposite"
}

// ErrThingWithRequiredFieldsNotFound is returned when the database fails to find a ThingWithRequiredFields.
type ErrThingWithRequiredFieldsNotFound struct {
	Name string
}

var _ error = ErrThingWithRequiredFieldsNotFound{}

// Error returns a description of the error.
func (e ErrThingWithRequiredFieldsNotFound) Error() string {
	return "could not find ThingWithRequiredFields"
}

// ErrThingWithRequiredFieldsAlreadyExists is returned when trying to overwrite a ThingWithRequiredFields.
type ErrThingWithRequiredFieldsAlreadyExists struct {
	Name string
}

var _ error = ErrThingWithRequiredFieldsAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithRequiredFieldsAlreadyExists) Error() string {
	return "ThingWithRequiredFields already exists"
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
