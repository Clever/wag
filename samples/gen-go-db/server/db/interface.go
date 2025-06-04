package db

import (
	"context"

	"github.com/Clever/wag/samples/gen-go-db/models/v9"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/go-openapi/strfmt"
	"golang.org/x/time/rate"
)

//go:generate mockgen -source=$GOFILE -destination=mock_db.go -package db --build_flags=--mod=mod -imports=models=github.com/Clever/wag/samples/gen-go-db/models/v9

// Interface for interacting with the swagger-test database.
type Interface interface {
	// SaveDeployment saves a Deployment to the database.
	SaveDeployment(ctx context.Context, m models.Deployment) error
	// GetDeployment retrieves a Deployment from the database.
	GetDeployment(ctx context.Context, environment string, application string, version string) (*models.Deployment, error)
	// ScanDeployments runs a scan on the Deployments table.
	ScanDeployments(ctx context.Context, input ScanDeploymentsInput, fn func(m *models.Deployment, lastDeployment bool) bool) error
	// GetDeploymentsByEnvAppAndVersion retrieves a page of Deployments from the database.
	GetDeploymentsByEnvAppAndVersion(ctx context.Context, input GetDeploymentsByEnvAppAndVersionInput, fn func(m *models.Deployment, lastDeployment bool) bool) error
	// DeleteDeployment deletes a Deployment from the database.
	DeleteDeployment(ctx context.Context, environment string, application string, version string) error
	// GetDeploymentsByEnvAppAndDate retrieves a page of Deployments from the database.
	GetDeploymentsByEnvAppAndDate(ctx context.Context, input GetDeploymentsByEnvAppAndDateInput, fn func(m *models.Deployment, lastDeployment bool) bool) error
	// ScanDeploymentsByEnvAppAndDate runs a scan on the EnvAppAndDate index.
	ScanDeploymentsByEnvAppAndDate(ctx context.Context, input ScanDeploymentsByEnvAppAndDateInput, fn func(m *models.Deployment, lastDeployment bool) bool) error
	// GetDeploymentsByEnvironmentAndDate retrieves a page of Deployments from the database.
	GetDeploymentsByEnvironmentAndDate(ctx context.Context, input GetDeploymentsByEnvironmentAndDateInput, fn func(m *models.Deployment, lastDeployment bool) bool) error
	// GetDeploymentByVersion retrieves a Deployment from the database.
	GetDeploymentByVersion(ctx context.Context, version string) (*models.Deployment, error)
	// ScanDeploymentsByVersion runs a scan on the Version index.
	ScanDeploymentsByVersion(ctx context.Context, input ScanDeploymentsByVersionInput, fn func(m *models.Deployment, lastDeployment bool) bool) error

	// SaveEvent saves a Event to the database.
	SaveEvent(ctx context.Context, m models.Event) error
	// GetEvent retrieves a Event from the database.
	GetEvent(ctx context.Context, pk string, sk string) (*models.Event, error)
	// ScanEvents runs a scan on the Events table.
	ScanEvents(ctx context.Context, input ScanEventsInput, fn func(m *models.Event, lastEvent bool) bool) error
	// GetEventsByPkAndSk retrieves a page of Events from the database.
	GetEventsByPkAndSk(ctx context.Context, input GetEventsByPkAndSkInput, fn func(m *models.Event, lastEvent bool) bool) error
	// DeleteEvent deletes a Event from the database.
	DeleteEvent(ctx context.Context, pk string, sk string) error
	// GetEventsBySkAndData retrieves a page of Events from the database.
	GetEventsBySkAndData(ctx context.Context, input GetEventsBySkAndDataInput, fn func(m *models.Event, lastEvent bool) bool) error
	// ScanEventsBySkAndData runs a scan on the SkAndData index.
	ScanEventsBySkAndData(ctx context.Context, input ScanEventsBySkAndDataInput, fn func(m *models.Event, lastEvent bool) bool) error

	// SaveNoRangeThingWithCompositeAttributes saves a NoRangeThingWithCompositeAttributes to the database.
	SaveNoRangeThingWithCompositeAttributes(ctx context.Context, m models.NoRangeThingWithCompositeAttributes) error
	// GetNoRangeThingWithCompositeAttributes retrieves a NoRangeThingWithCompositeAttributes from the database.
	GetNoRangeThingWithCompositeAttributes(ctx context.Context, name string, branch string) (*models.NoRangeThingWithCompositeAttributes, error)
	// ScanNoRangeThingWithCompositeAttributess runs a scan on the NoRangeThingWithCompositeAttributess table.
	ScanNoRangeThingWithCompositeAttributess(ctx context.Context, input ScanNoRangeThingWithCompositeAttributessInput, fn func(m *models.NoRangeThingWithCompositeAttributes, lastNoRangeThingWithCompositeAttributes bool) bool) error
	// DeleteNoRangeThingWithCompositeAttributes deletes a NoRangeThingWithCompositeAttributes from the database.
	DeleteNoRangeThingWithCompositeAttributes(ctx context.Context, name string, branch string) error
	// GetNoRangeThingWithCompositeAttributessByNameVersionAndDate retrieves a page of NoRangeThingWithCompositeAttributess from the database.
	GetNoRangeThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.NoRangeThingWithCompositeAttributes, lastNoRangeThingWithCompositeAttributes bool) bool) error
	// ScanNoRangeThingWithCompositeAttributessByNameVersionAndDate runs a scan on the NameVersionAndDate index.
	ScanNoRangeThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input ScanNoRangeThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.NoRangeThingWithCompositeAttributes, lastNoRangeThingWithCompositeAttributes bool) bool) error
	// GetNoRangeThingWithCompositeAttributesByNameBranchCommit retrieves a NoRangeThingWithCompositeAttributes from the database.
	GetNoRangeThingWithCompositeAttributesByNameBranchCommit(ctx context.Context, name string, branch string, commit string) (*models.NoRangeThingWithCompositeAttributes, error)

	// SaveSimpleThing saves a SimpleThing to the database.
	SaveSimpleThing(ctx context.Context, m models.SimpleThing) error
	// GetSimpleThing retrieves a SimpleThing from the database.
	GetSimpleThing(ctx context.Context, name string) (*models.SimpleThing, error)
	// ScanSimpleThings runs a scan on the SimpleThings table.
	ScanSimpleThings(ctx context.Context, input ScanSimpleThingsInput, fn func(m *models.SimpleThing, lastSimpleThing bool) bool) error
	// DeleteSimpleThing deletes a SimpleThing from the database.
	DeleteSimpleThing(ctx context.Context, name string) error

	// SaveTeacherSharingRule saves a TeacherSharingRule to the database.
	SaveTeacherSharingRule(ctx context.Context, m models.TeacherSharingRule) error
	// GetTeacherSharingRule retrieves a TeacherSharingRule from the database.
	GetTeacherSharingRule(ctx context.Context, teacher string, school string, app string) (*models.TeacherSharingRule, error)
	// ScanTeacherSharingRules runs a scan on the TeacherSharingRules table.
	ScanTeacherSharingRules(ctx context.Context, input ScanTeacherSharingRulesInput, fn func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool) error
	// GetTeacherSharingRulesByTeacherAndSchoolApp retrieves a page of TeacherSharingRules from the database.
	GetTeacherSharingRulesByTeacherAndSchoolApp(ctx context.Context, input GetTeacherSharingRulesByTeacherAndSchoolAppInput, fn func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool) error
	// DeleteTeacherSharingRule deletes a TeacherSharingRule from the database.
	DeleteTeacherSharingRule(ctx context.Context, teacher string, school string, app string) error
	// GetTeacherSharingRulesByDistrictAndSchoolTeacherApp retrieves a page of TeacherSharingRules from the database.
	GetTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx context.Context, input GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput, fn func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool) error
	// ScanTeacherSharingRulesByDistrictAndSchoolTeacherApp runs a scan on the DistrictAndSchoolTeacherApp index.
	ScanTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx context.Context, input ScanTeacherSharingRulesByDistrictAndSchoolTeacherAppInput, fn func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool) error

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
	// ScanThingsByID runs a scan on the ID index.
	ScanThingsByID(ctx context.Context, input ScanThingsByIDInput, fn func(m *models.Thing, lastThing bool) bool) error
	// GetThingsByNameAndCreatedAt retrieves a page of Things from the database.
	GetThingsByNameAndCreatedAt(ctx context.Context, input GetThingsByNameAndCreatedAtInput, fn func(m *models.Thing, lastThing bool) bool) error
	// ScanThingsByNameAndCreatedAt runs a scan on the NameAndCreatedAt index.
	ScanThingsByNameAndCreatedAt(ctx context.Context, input ScanThingsByNameAndCreatedAtInput, fn func(m *models.Thing, lastThing bool) bool) error
	// GetThingsByNameAndRangeNullable retrieves a page of Things from the database.
	GetThingsByNameAndRangeNullable(ctx context.Context, input GetThingsByNameAndRangeNullableInput, fn func(m *models.Thing, lastThing bool) bool) error
	// ScanThingsByNameAndRangeNullable runs a scan on the NameAndRangeNullable index.
	ScanThingsByNameAndRangeNullable(ctx context.Context, input ScanThingsByNameAndRangeNullableInput, fn func(m *models.Thing, lastThing bool) bool) error
	// GetThingsByHashNullableAndName retrieves a page of Things from the database.
	GetThingsByHashNullableAndName(ctx context.Context, input GetThingsByHashNullableAndNameInput, fn func(m *models.Thing, lastThing bool) bool) error

	// SaveThingAllowingBatchWrites saves a ThingAllowingBatchWrites to the database.
	SaveThingAllowingBatchWrites(ctx context.Context, m models.ThingAllowingBatchWrites) error
	// SaveArrayOfThingAllowingBatchWrites batch saves all items in []ThingAllowingBatchWrites to the database.
	SaveArrayOfThingAllowingBatchWrites(ctx context.Context, ms []models.ThingAllowingBatchWrites) error
	// DeleteArrayOfThingAllowingBatchWrites batch deletes all items in []ThingAllowingBatchWrites in the database.
	DeleteArrayOfThingAllowingBatchWrites(ctx context.Context, ms []models.ThingAllowingBatchWrites) error
	// GetThingAllowingBatchWrites retrieves a ThingAllowingBatchWrites from the database.
	GetThingAllowingBatchWrites(ctx context.Context, name string, version int64) (*models.ThingAllowingBatchWrites, error)
	// ScanThingAllowingBatchWritess runs a scan on the ThingAllowingBatchWritess table.
	ScanThingAllowingBatchWritess(ctx context.Context, input ScanThingAllowingBatchWritessInput, fn func(m *models.ThingAllowingBatchWrites, lastThingAllowingBatchWrites bool) bool) error
	// GetThingAllowingBatchWritessByNameAndVersion retrieves a page of ThingAllowingBatchWritess from the database.
	GetThingAllowingBatchWritessByNameAndVersion(ctx context.Context, input GetThingAllowingBatchWritessByNameAndVersionInput, fn func(m *models.ThingAllowingBatchWrites, lastThingAllowingBatchWrites bool) bool) error
	// DeleteThingAllowingBatchWrites deletes a ThingAllowingBatchWrites from the database.
	DeleteThingAllowingBatchWrites(ctx context.Context, name string, version int64) error

	// SaveThingAllowingBatchWritesWithCompositeAttributes saves a ThingAllowingBatchWritesWithCompositeAttributes to the database.
	SaveThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, m models.ThingAllowingBatchWritesWithCompositeAttributes) error
	// SaveArrayOfThingAllowingBatchWritesWithCompositeAttributes batch saves all items in []ThingAllowingBatchWritesWithCompositeAttributes to the database.
	SaveArrayOfThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, ms []models.ThingAllowingBatchWritesWithCompositeAttributes) error
	// DeleteArrayOfThingAllowingBatchWritesWithCompositeAttributes batch deletes all items in []ThingAllowingBatchWritesWithCompositeAttributes in the database.
	DeleteArrayOfThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, ms []models.ThingAllowingBatchWritesWithCompositeAttributes) error
	// GetThingAllowingBatchWritesWithCompositeAttributes retrieves a ThingAllowingBatchWritesWithCompositeAttributes from the database.
	GetThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, name string, id string, date strfmt.DateTime) (*models.ThingAllowingBatchWritesWithCompositeAttributes, error)
	// ScanThingAllowingBatchWritesWithCompositeAttributess runs a scan on the ThingAllowingBatchWritesWithCompositeAttributess table.
	ScanThingAllowingBatchWritesWithCompositeAttributess(ctx context.Context, input ScanThingAllowingBatchWritesWithCompositeAttributessInput, fn func(m *models.ThingAllowingBatchWritesWithCompositeAttributes, lastThingAllowingBatchWritesWithCompositeAttributes bool) bool) error
	// GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDate retrieves a page of ThingAllowingBatchWritesWithCompositeAttributess from the database.
	GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDate(ctx context.Context, input GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput, fn func(m *models.ThingAllowingBatchWritesWithCompositeAttributes, lastThingAllowingBatchWritesWithCompositeAttributes bool) bool) error
	// DeleteThingAllowingBatchWritesWithCompositeAttributes deletes a ThingAllowingBatchWritesWithCompositeAttributes from the database.
	DeleteThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, name string, id string, date strfmt.DateTime) error

	// SaveThingWithAdditionalAttributes saves a ThingWithAdditionalAttributes to the database.
	SaveThingWithAdditionalAttributes(ctx context.Context, m models.ThingWithAdditionalAttributes) error
	// GetThingWithAdditionalAttributes retrieves a ThingWithAdditionalAttributes from the database.
	GetThingWithAdditionalAttributes(ctx context.Context, name string, version int64) (*models.ThingWithAdditionalAttributes, error)
	// ScanThingWithAdditionalAttributess runs a scan on the ThingWithAdditionalAttributess table.
	ScanThingWithAdditionalAttributess(ctx context.Context, input ScanThingWithAdditionalAttributessInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error
	// GetThingWithAdditionalAttributessByNameAndVersion retrieves a page of ThingWithAdditionalAttributess from the database.
	GetThingWithAdditionalAttributessByNameAndVersion(ctx context.Context, input GetThingWithAdditionalAttributessByNameAndVersionInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error
	// DeleteThingWithAdditionalAttributes deletes a ThingWithAdditionalAttributes from the database.
	DeleteThingWithAdditionalAttributes(ctx context.Context, name string, version int64) error
	// GetThingWithAdditionalAttributesByID retrieves a ThingWithAdditionalAttributes from the database.
	GetThingWithAdditionalAttributesByID(ctx context.Context, id string) (*models.ThingWithAdditionalAttributes, error)
	// ScanThingWithAdditionalAttributessByID runs a scan on the ID index.
	ScanThingWithAdditionalAttributessByID(ctx context.Context, input ScanThingWithAdditionalAttributessByIDInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error
	// GetThingWithAdditionalAttributessByNameAndCreatedAt retrieves a page of ThingWithAdditionalAttributess from the database.
	GetThingWithAdditionalAttributessByNameAndCreatedAt(ctx context.Context, input GetThingWithAdditionalAttributessByNameAndCreatedAtInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error
	// ScanThingWithAdditionalAttributessByNameAndCreatedAt runs a scan on the NameAndCreatedAt index.
	ScanThingWithAdditionalAttributessByNameAndCreatedAt(ctx context.Context, input ScanThingWithAdditionalAttributessByNameAndCreatedAtInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error
	// GetThingWithAdditionalAttributessByNameAndRangeNullable retrieves a page of ThingWithAdditionalAttributess from the database.
	GetThingWithAdditionalAttributessByNameAndRangeNullable(ctx context.Context, input GetThingWithAdditionalAttributessByNameAndRangeNullableInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error
	// ScanThingWithAdditionalAttributessByNameAndRangeNullable runs a scan on the NameAndRangeNullable index.
	ScanThingWithAdditionalAttributessByNameAndRangeNullable(ctx context.Context, input ScanThingWithAdditionalAttributessByNameAndRangeNullableInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error
	// GetThingWithAdditionalAttributessByHashNullableAndName retrieves a page of ThingWithAdditionalAttributess from the database.
	GetThingWithAdditionalAttributessByHashNullableAndName(ctx context.Context, input GetThingWithAdditionalAttributessByHashNullableAndNameInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error

	// SaveThingWithCompositeAttributes saves a ThingWithCompositeAttributes to the database.
	SaveThingWithCompositeAttributes(ctx context.Context, m models.ThingWithCompositeAttributes) error
	// GetThingWithCompositeAttributes retrieves a ThingWithCompositeAttributes from the database.
	GetThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) (*models.ThingWithCompositeAttributes, error)
	// ScanThingWithCompositeAttributess runs a scan on the ThingWithCompositeAttributess table.
	ScanThingWithCompositeAttributess(ctx context.Context, input ScanThingWithCompositeAttributessInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error
	// GetThingWithCompositeAttributessByNameBranchAndDate retrieves a page of ThingWithCompositeAttributess from the database.
	GetThingWithCompositeAttributessByNameBranchAndDate(ctx context.Context, input GetThingWithCompositeAttributessByNameBranchAndDateInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error
	// DeleteThingWithCompositeAttributes deletes a ThingWithCompositeAttributes from the database.
	DeleteThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) error
	// GetThingWithCompositeAttributessByNameVersionAndDate retrieves a page of ThingWithCompositeAttributess from the database.
	GetThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input GetThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error
	// ScanThingWithCompositeAttributessByNameVersionAndDate runs a scan on the NameVersionAndDate index.
	ScanThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input ScanThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error

	// SaveThingWithCompositeEnumAttributes saves a ThingWithCompositeEnumAttributes to the database.
	SaveThingWithCompositeEnumAttributes(ctx context.Context, m models.ThingWithCompositeEnumAttributes) error
	// GetThingWithCompositeEnumAttributes retrieves a ThingWithCompositeEnumAttributes from the database.
	GetThingWithCompositeEnumAttributes(ctx context.Context, name string, branchID models.Branch, date strfmt.DateTime) (*models.ThingWithCompositeEnumAttributes, error)
	// ScanThingWithCompositeEnumAttributess runs a scan on the ThingWithCompositeEnumAttributess table.
	ScanThingWithCompositeEnumAttributess(ctx context.Context, input ScanThingWithCompositeEnumAttributessInput, fn func(m *models.ThingWithCompositeEnumAttributes, lastThingWithCompositeEnumAttributes bool) bool) error
	// GetThingWithCompositeEnumAttributessByNameBranchAndDate retrieves a page of ThingWithCompositeEnumAttributess from the database.
	GetThingWithCompositeEnumAttributessByNameBranchAndDate(ctx context.Context, input GetThingWithCompositeEnumAttributessByNameBranchAndDateInput, fn func(m *models.ThingWithCompositeEnumAttributes, lastThingWithCompositeEnumAttributes bool) bool) error
	// DeleteThingWithCompositeEnumAttributes deletes a ThingWithCompositeEnumAttributes from the database.
	DeleteThingWithCompositeEnumAttributes(ctx context.Context, name string, branchID models.Branch, date strfmt.DateTime) error

	// SaveThingWithDateGSI saves a ThingWithDateGSI to the database.
	SaveThingWithDateGSI(ctx context.Context, m models.ThingWithDateGSI) error
	// GetThingWithDateGSI retrieves a ThingWithDateGSI from the database.
	GetThingWithDateGSI(ctx context.Context, dateH strfmt.Date) (*models.ThingWithDateGSI, error)
	// ScanThingWithDateGSIs runs a scan on the ThingWithDateGSIs table.
	ScanThingWithDateGSIs(ctx context.Context, input ScanThingWithDateGSIsInput, fn func(m *models.ThingWithDateGSI, lastThingWithDateGSI bool) bool) error
	// DeleteThingWithDateGSI deletes a ThingWithDateGSI from the database.
	DeleteThingWithDateGSI(ctx context.Context, dateH strfmt.Date) error
	// GetThingWithDateGSIsByIDAndDateR retrieves a page of ThingWithDateGSIs from the database.
	GetThingWithDateGSIsByIDAndDateR(ctx context.Context, input GetThingWithDateGSIsByIDAndDateRInput, fn func(m *models.ThingWithDateGSI, lastThingWithDateGSI bool) bool) error
	// GetThingWithDateGSIsByDateHAndID retrieves a page of ThingWithDateGSIs from the database.
	GetThingWithDateGSIsByDateHAndID(ctx context.Context, input GetThingWithDateGSIsByDateHAndIDInput, fn func(m *models.ThingWithDateGSI, lastThingWithDateGSI bool) bool) error

	// SaveThingWithDateRange saves a ThingWithDateRange to the database.
	SaveThingWithDateRange(ctx context.Context, m models.ThingWithDateRange) error
	// GetThingWithDateRange retrieves a ThingWithDateRange from the database.
	GetThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) (*models.ThingWithDateRange, error)
	// ScanThingWithDateRanges runs a scan on the ThingWithDateRanges table.
	ScanThingWithDateRanges(ctx context.Context, input ScanThingWithDateRangesInput, fn func(m *models.ThingWithDateRange, lastThingWithDateRange bool) bool) error
	// GetThingWithDateRangesByNameAndDate retrieves a page of ThingWithDateRanges from the database.
	GetThingWithDateRangesByNameAndDate(ctx context.Context, input GetThingWithDateRangesByNameAndDateInput, fn func(m *models.ThingWithDateRange, lastThingWithDateRange bool) bool) error
	// DeleteThingWithDateRange deletes a ThingWithDateRange from the database.
	DeleteThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) error

	// SaveThingWithDateRangeKey saves a ThingWithDateRangeKey to the database.
	SaveThingWithDateRangeKey(ctx context.Context, m models.ThingWithDateRangeKey) error
	// GetThingWithDateRangeKey retrieves a ThingWithDateRangeKey from the database.
	GetThingWithDateRangeKey(ctx context.Context, id string, date strfmt.Date) (*models.ThingWithDateRangeKey, error)
	// ScanThingWithDateRangeKeys runs a scan on the ThingWithDateRangeKeys table.
	ScanThingWithDateRangeKeys(ctx context.Context, input ScanThingWithDateRangeKeysInput, fn func(m *models.ThingWithDateRangeKey, lastThingWithDateRangeKey bool) bool) error
	// GetThingWithDateRangeKeysByIDAndDate retrieves a page of ThingWithDateRangeKeys from the database.
	GetThingWithDateRangeKeysByIDAndDate(ctx context.Context, input GetThingWithDateRangeKeysByIDAndDateInput, fn func(m *models.ThingWithDateRangeKey, lastThingWithDateRangeKey bool) bool) error
	// DeleteThingWithDateRangeKey deletes a ThingWithDateRangeKey from the database.
	DeleteThingWithDateRangeKey(ctx context.Context, id string, date strfmt.Date) error

	// SaveThingWithDateTimeComposite saves a ThingWithDateTimeComposite to the database.
	SaveThingWithDateTimeComposite(ctx context.Context, m models.ThingWithDateTimeComposite) error
	// GetThingWithDateTimeComposite retrieves a ThingWithDateTimeComposite from the database.
	GetThingWithDateTimeComposite(ctx context.Context, typeVar string, id string, created strfmt.DateTime, resource string) (*models.ThingWithDateTimeComposite, error)
	// ScanThingWithDateTimeComposites runs a scan on the ThingWithDateTimeComposites table.
	ScanThingWithDateTimeComposites(ctx context.Context, input ScanThingWithDateTimeCompositesInput, fn func(m *models.ThingWithDateTimeComposite, lastThingWithDateTimeComposite bool) bool) error
	// GetThingWithDateTimeCompositesByTypeIDAndCreatedResource retrieves a page of ThingWithDateTimeComposites from the database.
	GetThingWithDateTimeCompositesByTypeIDAndCreatedResource(ctx context.Context, input GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput, fn func(m *models.ThingWithDateTimeComposite, lastThingWithDateTimeComposite bool) bool) error
	// DeleteThingWithDateTimeComposite deletes a ThingWithDateTimeComposite from the database.
	DeleteThingWithDateTimeComposite(ctx context.Context, typeVar string, id string, created strfmt.DateTime, resource string) error

	// SaveThingWithDatetimeGSI saves a ThingWithDatetimeGSI to the database.
	SaveThingWithDatetimeGSI(ctx context.Context, m models.ThingWithDatetimeGSI) error
	// GetThingWithDatetimeGSI retrieves a ThingWithDatetimeGSI from the database.
	GetThingWithDatetimeGSI(ctx context.Context, id string) (*models.ThingWithDatetimeGSI, error)
	// ScanThingWithDatetimeGSIs runs a scan on the ThingWithDatetimeGSIs table.
	ScanThingWithDatetimeGSIs(ctx context.Context, input ScanThingWithDatetimeGSIsInput, fn func(m *models.ThingWithDatetimeGSI, lastThingWithDatetimeGSI bool) bool) error
	// DeleteThingWithDatetimeGSI deletes a ThingWithDatetimeGSI from the database.
	DeleteThingWithDatetimeGSI(ctx context.Context, id string) error
	// GetThingWithDatetimeGSIsByDatetimeAndID retrieves a page of ThingWithDatetimeGSIs from the database.
	GetThingWithDatetimeGSIsByDatetimeAndID(ctx context.Context, input GetThingWithDatetimeGSIsByDatetimeAndIDInput, fn func(m *models.ThingWithDatetimeGSI, lastThingWithDatetimeGSI bool) bool) error
	// ScanThingWithDatetimeGSIsByDatetimeAndID runs a scan on the DatetimeAndID index.
	ScanThingWithDatetimeGSIsByDatetimeAndID(ctx context.Context, input ScanThingWithDatetimeGSIsByDatetimeAndIDInput, fn func(m *models.ThingWithDatetimeGSI, lastThingWithDatetimeGSI bool) bool) error

	// SaveThingWithEnumHashKey saves a ThingWithEnumHashKey to the database.
	SaveThingWithEnumHashKey(ctx context.Context, m models.ThingWithEnumHashKey) error
	// GetThingWithEnumHashKey retrieves a ThingWithEnumHashKey from the database.
	GetThingWithEnumHashKey(ctx context.Context, branch models.Branch, date strfmt.DateTime) (*models.ThingWithEnumHashKey, error)
	// ScanThingWithEnumHashKeys runs a scan on the ThingWithEnumHashKeys table.
	ScanThingWithEnumHashKeys(ctx context.Context, input ScanThingWithEnumHashKeysInput, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error
	// GetThingWithEnumHashKeysByBranchAndDate retrieves a page of ThingWithEnumHashKeys from the database.
	GetThingWithEnumHashKeysByBranchAndDate(ctx context.Context, input GetThingWithEnumHashKeysByBranchAndDateInput, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error
	// DeleteThingWithEnumHashKey deletes a ThingWithEnumHashKey from the database.
	DeleteThingWithEnumHashKey(ctx context.Context, branch models.Branch, date strfmt.DateTime) error
	// GetThingWithEnumHashKeysByBranchAndDate2 retrieves a page of ThingWithEnumHashKeys from the database.
	GetThingWithEnumHashKeysByBranchAndDate2(ctx context.Context, input GetThingWithEnumHashKeysByBranchAndDate2Input, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error
	// ScanThingWithEnumHashKeysByBranchAndDate2 runs a scan on the BranchAndDate2 index.
	ScanThingWithEnumHashKeysByBranchAndDate2(ctx context.Context, input ScanThingWithEnumHashKeysByBranchAndDate2Input, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error

	// SaveThingWithMatchingKeys saves a ThingWithMatchingKeys to the database.
	SaveThingWithMatchingKeys(ctx context.Context, m models.ThingWithMatchingKeys) error
	// GetThingWithMatchingKeys retrieves a ThingWithMatchingKeys from the database.
	GetThingWithMatchingKeys(ctx context.Context, bear string, assocType string, assocID string) (*models.ThingWithMatchingKeys, error)
	// ScanThingWithMatchingKeyss runs a scan on the ThingWithMatchingKeyss table.
	ScanThingWithMatchingKeyss(ctx context.Context, input ScanThingWithMatchingKeyssInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error
	// GetThingWithMatchingKeyssByBearAndAssocTypeID retrieves a page of ThingWithMatchingKeyss from the database.
	GetThingWithMatchingKeyssByBearAndAssocTypeID(ctx context.Context, input GetThingWithMatchingKeyssByBearAndAssocTypeIDInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error
	// DeleteThingWithMatchingKeys deletes a ThingWithMatchingKeys from the database.
	DeleteThingWithMatchingKeys(ctx context.Context, bear string, assocType string, assocID string) error
	// GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBear retrieves a page of ThingWithMatchingKeyss from the database.
	GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx context.Context, input GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error
	// ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear runs a scan on the AssocTypeIDAndCreatedBear index.
	ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx context.Context, input ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error

	// SaveThingWithMultiUseCompositeAttribute saves a ThingWithMultiUseCompositeAttribute to the database.
	SaveThingWithMultiUseCompositeAttribute(ctx context.Context, m models.ThingWithMultiUseCompositeAttribute) error
	// GetThingWithMultiUseCompositeAttribute retrieves a ThingWithMultiUseCompositeAttribute from the database.
	GetThingWithMultiUseCompositeAttribute(ctx context.Context, one string) (*models.ThingWithMultiUseCompositeAttribute, error)
	// ScanThingWithMultiUseCompositeAttributes runs a scan on the ThingWithMultiUseCompositeAttributes table.
	ScanThingWithMultiUseCompositeAttributes(ctx context.Context, input ScanThingWithMultiUseCompositeAttributesInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error
	// DeleteThingWithMultiUseCompositeAttribute deletes a ThingWithMultiUseCompositeAttribute from the database.
	DeleteThingWithMultiUseCompositeAttribute(ctx context.Context, one string) error
	// GetThingWithMultiUseCompositeAttributesByThreeAndOneTwo retrieves a page of ThingWithMultiUseCompositeAttributes from the database.
	GetThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx context.Context, input GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error
	// ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwo runs a scan on the ThreeAndOneTwo index.
	ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx context.Context, input ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error
	// GetThingWithMultiUseCompositeAttributesByFourAndOneTwo retrieves a page of ThingWithMultiUseCompositeAttributes from the database.
	GetThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx context.Context, input GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error
	// ScanThingWithMultiUseCompositeAttributesByFourAndOneTwo runs a scan on the FourAndOneTwo index.
	ScanThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx context.Context, input ScanThingWithMultiUseCompositeAttributesByFourAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error

	// SaveThingWithRequiredCompositePropertiesAndKeysOnly saves a ThingWithRequiredCompositePropertiesAndKeysOnly to the database.
	SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, m models.ThingWithRequiredCompositePropertiesAndKeysOnly) error
	// GetThingWithRequiredCompositePropertiesAndKeysOnly retrieves a ThingWithRequiredCompositePropertiesAndKeysOnly from the database.
	GetThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, propertyThree string) (*models.ThingWithRequiredCompositePropertiesAndKeysOnly, error)
	// ScanThingWithRequiredCompositePropertiesAndKeysOnlys runs a scan on the ThingWithRequiredCompositePropertiesAndKeysOnlys table.
	ScanThingWithRequiredCompositePropertiesAndKeysOnlys(ctx context.Context, input ScanThingWithRequiredCompositePropertiesAndKeysOnlysInput, fn func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, lastThingWithRequiredCompositePropertiesAndKeysOnly bool) bool) error
	// DeleteThingWithRequiredCompositePropertiesAndKeysOnly deletes a ThingWithRequiredCompositePropertiesAndKeysOnly from the database.
	DeleteThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, propertyThree string) error
	// GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree retrieves a page of ThingWithRequiredCompositePropertiesAndKeysOnlys from the database.
	GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx context.Context, input GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput, fn func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, lastThingWithRequiredCompositePropertiesAndKeysOnly bool) bool) error
	// ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree runs a scan on the PropertyOneAndTwoAndPropertyThree index.
	ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx context.Context, input ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput, fn func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, lastThingWithRequiredCompositePropertiesAndKeysOnly bool) bool) error

	// SaveThingWithRequiredFields saves a ThingWithRequiredFields to the database.
	SaveThingWithRequiredFields(ctx context.Context, m models.ThingWithRequiredFields) error
	// GetThingWithRequiredFields retrieves a ThingWithRequiredFields from the database.
	GetThingWithRequiredFields(ctx context.Context, name string) (*models.ThingWithRequiredFields, error)
	// ScanThingWithRequiredFieldss runs a scan on the ThingWithRequiredFieldss table.
	ScanThingWithRequiredFieldss(ctx context.Context, input ScanThingWithRequiredFieldssInput, fn func(m *models.ThingWithRequiredFields, lastThingWithRequiredFields bool) bool) error
	// DeleteThingWithRequiredFields deletes a ThingWithRequiredFields from the database.
	DeleteThingWithRequiredFields(ctx context.Context, name string) error

	// SaveThingWithRequiredFields2 saves a ThingWithRequiredFields2 to the database.
	SaveThingWithRequiredFields2(ctx context.Context, m models.ThingWithRequiredFields2) error
	// GetThingWithRequiredFields2 retrieves a ThingWithRequiredFields2 from the database.
	GetThingWithRequiredFields2(ctx context.Context, name string, id string) (*models.ThingWithRequiredFields2, error)
	// ScanThingWithRequiredFields2s runs a scan on the ThingWithRequiredFields2s table.
	ScanThingWithRequiredFields2s(ctx context.Context, input ScanThingWithRequiredFields2sInput, fn func(m *models.ThingWithRequiredFields2, lastThingWithRequiredFields2 bool) bool) error
	// GetThingWithRequiredFields2sByNameAndID retrieves a page of ThingWithRequiredFields2s from the database.
	GetThingWithRequiredFields2sByNameAndID(ctx context.Context, input GetThingWithRequiredFields2sByNameAndIDInput, fn func(m *models.ThingWithRequiredFields2, lastThingWithRequiredFields2 bool) bool) error
	// DeleteThingWithRequiredFields2 deletes a ThingWithRequiredFields2 from the database.
	DeleteThingWithRequiredFields2(ctx context.Context, name string, id string) error

	// SaveThingWithTransactMultipleGSI saves a ThingWithTransactMultipleGSI to the database.
	SaveThingWithTransactMultipleGSI(ctx context.Context, m models.ThingWithTransactMultipleGSI) error
	// GetThingWithTransactMultipleGSI retrieves a ThingWithTransactMultipleGSI from the database.
	GetThingWithTransactMultipleGSI(ctx context.Context, dateH strfmt.Date) (*models.ThingWithTransactMultipleGSI, error)
	// ScanThingWithTransactMultipleGSIs runs a scan on the ThingWithTransactMultipleGSIs table.
	ScanThingWithTransactMultipleGSIs(ctx context.Context, input ScanThingWithTransactMultipleGSIsInput, fn func(m *models.ThingWithTransactMultipleGSI, lastThingWithTransactMultipleGSI bool) bool) error
	// DeleteThingWithTransactMultipleGSI deletes a ThingWithTransactMultipleGSI from the database.
	DeleteThingWithTransactMultipleGSI(ctx context.Context, dateH strfmt.Date) error
	// GetThingWithTransactMultipleGSIsByIDAndDateR retrieves a page of ThingWithTransactMultipleGSIs from the database.
	GetThingWithTransactMultipleGSIsByIDAndDateR(ctx context.Context, input GetThingWithTransactMultipleGSIsByIDAndDateRInput, fn func(m *models.ThingWithTransactMultipleGSI, lastThingWithTransactMultipleGSI bool) bool) error
	// GetThingWithTransactMultipleGSIsByDateHAndID retrieves a page of ThingWithTransactMultipleGSIs from the database.
	GetThingWithTransactMultipleGSIsByDateHAndID(ctx context.Context, input GetThingWithTransactMultipleGSIsByDateHAndIDInput, fn func(m *models.ThingWithTransactMultipleGSI, lastThingWithTransactMultipleGSI bool) bool) error
	// TransactSaveThingWithTransactMultipleGSIAndThing saves ThingWithTransactMultipleGSI and Thing as an atomic transaction.
	// Use the optional condition parameters to require pre-transaction conditions for each put
	TransactSaveThingWithTransactMultipleGSIAndThing(ctx context.Context, m1 models.ThingWithTransactMultipleGSI, m1Conditions *expression.ConditionBuilder, m2 models.Thing, m2Conditions *expression.ConditionBuilder) error

	// SaveThingWithTransaction saves a ThingWithTransaction to the database.
	SaveThingWithTransaction(ctx context.Context, m models.ThingWithTransaction) error
	// GetThingWithTransaction retrieves a ThingWithTransaction from the database.
	GetThingWithTransaction(ctx context.Context, name string) (*models.ThingWithTransaction, error)
	// ScanThingWithTransactions runs a scan on the ThingWithTransactions table.
	ScanThingWithTransactions(ctx context.Context, input ScanThingWithTransactionsInput, fn func(m *models.ThingWithTransaction, lastThingWithTransaction bool) bool) error
	// DeleteThingWithTransaction deletes a ThingWithTransaction from the database.
	DeleteThingWithTransaction(ctx context.Context, name string) error
	// TransactSaveThingWithTransactionAndThing saves ThingWithTransaction and Thing as an atomic transaction.
	// Use the optional condition parameters to require pre-transaction conditions for each put
	TransactSaveThingWithTransactionAndThing(ctx context.Context, m1 models.ThingWithTransaction, m1Conditions *expression.ConditionBuilder, m2 models.Thing, m2Conditions *expression.ConditionBuilder) error

	// SaveThingWithTransactionWithSimpleThing saves a ThingWithTransactionWithSimpleThing to the database.
	SaveThingWithTransactionWithSimpleThing(ctx context.Context, m models.ThingWithTransactionWithSimpleThing) error
	// GetThingWithTransactionWithSimpleThing retrieves a ThingWithTransactionWithSimpleThing from the database.
	GetThingWithTransactionWithSimpleThing(ctx context.Context, name string) (*models.ThingWithTransactionWithSimpleThing, error)
	// ScanThingWithTransactionWithSimpleThings runs a scan on the ThingWithTransactionWithSimpleThings table.
	ScanThingWithTransactionWithSimpleThings(ctx context.Context, input ScanThingWithTransactionWithSimpleThingsInput, fn func(m *models.ThingWithTransactionWithSimpleThing, lastThingWithTransactionWithSimpleThing bool) bool) error
	// DeleteThingWithTransactionWithSimpleThing deletes a ThingWithTransactionWithSimpleThing from the database.
	DeleteThingWithTransactionWithSimpleThing(ctx context.Context, name string) error
	// TransactSaveThingWithTransactionWithSimpleThingAndSimpleThing saves ThingWithTransactionWithSimpleThing and SimpleThing as an atomic transaction.
	// Use the optional condition parameters to require pre-transaction conditions for each put
	TransactSaveThingWithTransactionWithSimpleThingAndSimpleThing(ctx context.Context, m1 models.ThingWithTransactionWithSimpleThing, m1Conditions *expression.ConditionBuilder, m2 models.SimpleThing, m2Conditions *expression.ConditionBuilder) error

	// SaveThingWithUnderscores saves a ThingWithUnderscores to the database.
	SaveThingWithUnderscores(ctx context.Context, m models.ThingWithUnderscores) error
	// GetThingWithUnderscores retrieves a ThingWithUnderscores from the database.
	GetThingWithUnderscores(ctx context.Context, iDApp string) (*models.ThingWithUnderscores, error)
	// DeleteThingWithUnderscores deletes a ThingWithUnderscores from the database.
	DeleteThingWithUnderscores(ctx context.Context, iDApp string) error
}

// Int64 returns a pointer to the int64 value passed in.
func Int64(i int64) *int64 { return &i }

// Int32 returns a pointer to the int32 value passed in.
func Int32(i int32) *int32 { return &i }

// String returns a pointer to the string value passed in.
func String(s string) *string { return &s }

// DateTime returns a pointer to the strfmt.DateTime value passed in.
func DateTime(d strfmt.DateTime) *strfmt.DateTime { return &d }

// Date returns a pointer to the strfmt.Date value passed in.
func Date(d strfmt.Date) *strfmt.Date { return &d }

// ScanDeploymentsInput is the input to the ScanDeployments method.
type ScanDeploymentsInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.Deployment
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// DeploymentByEnvAppAndVersionFilterableAttribute represents the fields we can apply filters to for queries on this index
type DeploymentByEnvAppAndVersionFilterableAttribute string

const DeploymentDate DeploymentByEnvAppAndVersionFilterableAttribute = "date"

// DeploymentByEnvAppAndVersionFilterValues represents a filter on a particular field to be included in the query
type DeploymentByEnvAppAndVersionFilterValues struct {
	// AttributeName is the attibute we are attempting to apply the filter to
	AttributeName DeploymentByEnvAppAndVersionFilterableAttribute
	// AttributeValues is an optional parameter to be used when we want to compare the attibute to a single value or multiple values
	AttributeValues []interface{}
}

// GetDeploymentsByEnvAppAndVersionInput is the query input to GetDeploymentsByEnvAppAndVersion.
type GetDeploymentsByEnvAppAndVersionInput struct {
	// Environment is required
	Environment string
	// Application is required
	Application       string
	VersionStartingAt *string
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.Deployment
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// FilterValues is an optional array of filters to apply on various table attributes
	FilterValues []DeploymentByEnvAppAndVersionFilterValues
	// FilterExpression is the filter expression to be applied to our fitlered attributes
	// when referencing an attribute use #ATTRIBUTE_NAME
	// ex: if the attribute is called "created_at" in its wag definition use #CREATED_AT
	// when referencing one of the given values use :{attribute_name}_value0, :{attribute_name}_value1, etc.
	// ex: if the attribute is called "created_at" in its wag definition use :created_at_value0, created_at_value1, etc.
	// see https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.KeyConditionExpressions
	// for guidance on building expressions
	FilterExpression string
}

// ErrDeploymentNotFound is returned when the database fails to find a Deployment.
type ErrDeploymentNotFound struct {
	Environment string
	Application string
	Version     string
}

var _ error = ErrDeploymentNotFound{}

// Error returns a description of the error.
func (e ErrDeploymentNotFound) Error() string {
	return "could not find Deployment"
}

// GetDeploymentsByEnvAppAndDateInput is the query input to GetDeploymentsByEnvAppAndDate.
type GetDeploymentsByEnvAppAndDateInput struct {
	// Environment is required
	Environment string
	// Application is required
	Application    string
	DateStartingAt *strfmt.DateTime
	StartingAfter  *models.Deployment
	Descending     bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrDeploymentByEnvAppAndDateNotFound is returned when the database fails to find a Deployment.
type ErrDeploymentByEnvAppAndDateNotFound struct {
	Environment string
	Application string
	Date        strfmt.DateTime
}

var _ error = ErrDeploymentByEnvAppAndDateNotFound{}

// Error returns a description of the error.
func (e ErrDeploymentByEnvAppAndDateNotFound) Error() string {
	return "could not find Deployment"
}

// ScanDeploymentsByEnvAppAndDateInput is the input to the ScanDeploymentsByEnvAppAndDate method.
type ScanDeploymentsByEnvAppAndDateInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.Deployment
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetDeploymentsByEnvironmentAndDateInput is the query input to GetDeploymentsByEnvironmentAndDate.
type GetDeploymentsByEnvironmentAndDateInput struct {
	// Environment is required
	Environment    string
	DateStartingAt *strfmt.DateTime
	StartingAfter  *models.Deployment
	Descending     bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrDeploymentByEnvironmentAndDateNotFound is returned when the database fails to find a Deployment.
type ErrDeploymentByEnvironmentAndDateNotFound struct {
	Environment string
	Date        strfmt.DateTime
}

var _ error = ErrDeploymentByEnvironmentAndDateNotFound{}

// Error returns a description of the error.
func (e ErrDeploymentByEnvironmentAndDateNotFound) Error() string {
	return "could not find Deployment"
}

// ErrDeploymentByVersionNotFound is returned when the database fails to find a Deployment.
type ErrDeploymentByVersionNotFound struct {
	Version string
}

var _ error = ErrDeploymentByVersionNotFound{}

// Error returns a description of the error.
func (e ErrDeploymentByVersionNotFound) Error() string {
	return "could not find Deployment"
}

// ScanDeploymentsByVersionInput is the input to the ScanDeploymentsByVersion method.
type ScanDeploymentsByVersionInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.Deployment
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ScanEventsInput is the input to the ScanEvents method.
type ScanEventsInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.Event
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// EventByPkAndSkFilterableAttribute represents the fields we can apply filters to for queries on this index
type EventByPkAndSkFilterableAttribute string

const EventData EventByPkAndSkFilterableAttribute = "data"

// EventByPkAndSkFilterValues represents a filter on a particular field to be included in the query
type EventByPkAndSkFilterValues struct {
	// AttributeName is the attibute we are attempting to apply the filter to
	AttributeName EventByPkAndSkFilterableAttribute
	// AttributeValues is an optional parameter to be used when we want to compare the attibute to a single value or multiple values
	AttributeValues []interface{}
}

// GetEventsByPkAndSkInput is the query input to GetEventsByPkAndSk.
type GetEventsByPkAndSkInput struct {
	// Pk is required
	Pk           string
	SkStartingAt *string
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.Event
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// FilterValues is an optional array of filters to apply on various table attributes
	FilterValues []EventByPkAndSkFilterValues
	// FilterExpression is the filter expression to be applied to our fitlered attributes
	// when referencing an attribute use #ATTRIBUTE_NAME
	// ex: if the attribute is called "created_at" in its wag definition use #CREATED_AT
	// when referencing one of the given values use :{attribute_name}_value0, :{attribute_name}_value1, etc.
	// ex: if the attribute is called "created_at" in its wag definition use :created_at_value0, created_at_value1, etc.
	// see https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.KeyConditionExpressions
	// for guidance on building expressions
	FilterExpression string
}

// ErrEventNotFound is returned when the database fails to find a Event.
type ErrEventNotFound struct {
	Pk string
	Sk string
}

var _ error = ErrEventNotFound{}

// Error returns a description of the error.
func (e ErrEventNotFound) Error() string {
	return "could not find Event"
}

// GetEventsBySkAndDataInput is the query input to GetEventsBySkAndData.
type GetEventsBySkAndDataInput struct {
	// Sk is required
	Sk             string
	DataStartingAt []byte
	StartingAfter  *models.Event
	Descending     bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrEventBySkAndDataNotFound is returned when the database fails to find a Event.
type ErrEventBySkAndDataNotFound struct {
	Sk   string
	Data []byte
}

var _ error = ErrEventBySkAndDataNotFound{}

// Error returns a description of the error.
func (e ErrEventBySkAndDataNotFound) Error() string {
	return "could not find Event"
}

// ScanEventsBySkAndDataInput is the input to the ScanEventsBySkAndData method.
type ScanEventsBySkAndDataInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.Event
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ScanNoRangeThingWithCompositeAttributessInput is the input to the ScanNoRangeThingWithCompositeAttributess method.
type ScanNoRangeThingWithCompositeAttributessInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.NoRangeThingWithCompositeAttributes
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ErrNoRangeThingWithCompositeAttributesNotFound is returned when the database fails to find a NoRangeThingWithCompositeAttributes.
type ErrNoRangeThingWithCompositeAttributesNotFound struct {
	Name   string
	Branch string
}

var _ error = ErrNoRangeThingWithCompositeAttributesNotFound{}

// Error returns a description of the error.
func (e ErrNoRangeThingWithCompositeAttributesNotFound) Error() string {
	return "could not find NoRangeThingWithCompositeAttributes"
}

// GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput is the query input to GetNoRangeThingWithCompositeAttributessByNameVersionAndDate.
type GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput struct {
	// Name is required
	Name string
	// Version is required
	Version        int64
	DateStartingAt *strfmt.DateTime
	StartingAfter  *models.NoRangeThingWithCompositeAttributes
	Descending     bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrNoRangeThingWithCompositeAttributesByNameVersionAndDateNotFound is returned when the database fails to find a NoRangeThingWithCompositeAttributes.
type ErrNoRangeThingWithCompositeAttributesByNameVersionAndDateNotFound struct {
	Name    string
	Version int64
	Date    strfmt.DateTime
}

var _ error = ErrNoRangeThingWithCompositeAttributesByNameVersionAndDateNotFound{}

// Error returns a description of the error.
func (e ErrNoRangeThingWithCompositeAttributesByNameVersionAndDateNotFound) Error() string {
	return "could not find NoRangeThingWithCompositeAttributes"
}

// ScanNoRangeThingWithCompositeAttributessByNameVersionAndDateInput is the input to the ScanNoRangeThingWithCompositeAttributessByNameVersionAndDate method.
type ScanNoRangeThingWithCompositeAttributessByNameVersionAndDateInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.NoRangeThingWithCompositeAttributes
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ErrNoRangeThingWithCompositeAttributesByNameBranchCommitNotFound is returned when the database fails to find a NoRangeThingWithCompositeAttributes.
type ErrNoRangeThingWithCompositeAttributesByNameBranchCommitNotFound struct {
	Name   string
	Branch string
	Commit string
}

var _ error = ErrNoRangeThingWithCompositeAttributesByNameBranchCommitNotFound{}

// Error returns a description of the error.
func (e ErrNoRangeThingWithCompositeAttributesByNameBranchCommitNotFound) Error() string {
	return "could not find NoRangeThingWithCompositeAttributes"
}

// ErrNoRangeThingWithCompositeAttributesAlreadyExists is returned when trying to overwrite a NoRangeThingWithCompositeAttributes.
type ErrNoRangeThingWithCompositeAttributesAlreadyExists struct {
	NameBranch string
}

var _ error = ErrNoRangeThingWithCompositeAttributesAlreadyExists{}

// Error returns a description of the error.
func (e ErrNoRangeThingWithCompositeAttributesAlreadyExists) Error() string {
	return "NoRangeThingWithCompositeAttributes already exists"
}

// ScanSimpleThingsInput is the input to the ScanSimpleThings method.
type ScanSimpleThingsInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.SimpleThing
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

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

// ScanTeacherSharingRulesInput is the input to the ScanTeacherSharingRules method.
type ScanTeacherSharingRulesInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.TeacherSharingRule
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// TeacherSharingRuleByTeacherAndSchoolAppFilterableAttribute represents the fields we can apply filters to for queries on this index
type TeacherSharingRuleByTeacherAndSchoolAppFilterableAttribute string

const TeacherSharingRuleDistrict TeacherSharingRuleByTeacherAndSchoolAppFilterableAttribute = "district"

// TeacherSharingRuleByTeacherAndSchoolAppFilterValues represents a filter on a particular field to be included in the query
type TeacherSharingRuleByTeacherAndSchoolAppFilterValues struct {
	// AttributeName is the attibute we are attempting to apply the filter to
	AttributeName TeacherSharingRuleByTeacherAndSchoolAppFilterableAttribute
	// AttributeValues is an optional parameter to be used when we want to compare the attibute to a single value or multiple values
	AttributeValues []interface{}
}

// GetTeacherSharingRulesByTeacherAndSchoolAppInput is the query input to GetTeacherSharingRulesByTeacherAndSchoolApp.
type GetTeacherSharingRulesByTeacherAndSchoolAppInput struct {
	// Teacher is required
	Teacher    string
	StartingAt *SchoolApp
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.TeacherSharingRule
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// FilterValues is an optional array of filters to apply on various table attributes
	FilterValues []TeacherSharingRuleByTeacherAndSchoolAppFilterValues
	// FilterExpression is the filter expression to be applied to our fitlered attributes
	// when referencing an attribute use #ATTRIBUTE_NAME
	// ex: if the attribute is called "created_at" in its wag definition use #CREATED_AT
	// when referencing one of the given values use :{attribute_name}_value0, :{attribute_name}_value1, etc.
	// ex: if the attribute is called "created_at" in its wag definition use :created_at_value0, created_at_value1, etc.
	// see https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.KeyConditionExpressions
	// for guidance on building expressions
	FilterExpression string
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
	// District is required
	District      string
	StartingAt    *SchoolTeacherApp
	StartingAfter *models.TeacherSharingRule
	Descending    bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
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

// ScanTeacherSharingRulesByDistrictAndSchoolTeacherAppInput is the input to the ScanTeacherSharingRulesByDistrictAndSchoolTeacherApp method.
type ScanTeacherSharingRulesByDistrictAndSchoolTeacherAppInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.TeacherSharingRule
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// SchoolApp struct.
type SchoolApp struct {
	School string
	App    string
}

// SchoolTeacherApp struct.
type SchoolTeacherApp struct {
	School  string
	Teacher string
	App     string
}

// ScanThingsInput is the input to the ScanThings method.
type ScanThingsInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.Thing
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ThingByNameAndVersionFilterableAttribute represents the fields we can apply filters to for queries on this index
type ThingByNameAndVersionFilterableAttribute string

const ThingCreatedAt ThingByNameAndVersionFilterableAttribute = "createdAt"
const ThingHashNullable ThingByNameAndVersionFilterableAttribute = "hashNullable"
const ThingID ThingByNameAndVersionFilterableAttribute = "id"
const ThingRangeNullable ThingByNameAndVersionFilterableAttribute = "rangeNullable"

// ThingByNameAndVersionFilterValues represents a filter on a particular field to be included in the query
type ThingByNameAndVersionFilterValues struct {
	// AttributeName is the attibute we are attempting to apply the filter to
	AttributeName ThingByNameAndVersionFilterableAttribute
	// AttributeValues is an optional parameter to be used when we want to compare the attibute to a single value or multiple values
	AttributeValues []interface{}
}

// GetThingsByNameAndVersionInput is the query input to GetThingsByNameAndVersion.
type GetThingsByNameAndVersionInput struct {
	// Name is required
	Name              string
	VersionStartingAt *int64
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.Thing
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// FilterValues is an optional array of filters to apply on various table attributes
	FilterValues []ThingByNameAndVersionFilterValues
	// FilterExpression is the filter expression to be applied to our fitlered attributes
	// when referencing an attribute use #ATTRIBUTE_NAME
	// ex: if the attribute is called "created_at" in its wag definition use #CREATED_AT
	// when referencing one of the given values use :{attribute_name}_value0, :{attribute_name}_value1, etc.
	// ex: if the attribute is called "created_at" in its wag definition use :created_at_value0, created_at_value1, etc.
	// see https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.KeyConditionExpressions
	// for guidance on building expressions
	FilterExpression string
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

// ScanThingsByIDInput is the input to the ScanThingsByID method.
type ScanThingsByIDInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.Thing
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingsByNameAndCreatedAtInput is the query input to GetThingsByNameAndCreatedAt.
type GetThingsByNameAndCreatedAtInput struct {
	// Name is required
	Name                string
	CreatedAtStartingAt *strfmt.DateTime
	StartingAfter       *models.Thing
	Descending          bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
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

// ScanThingsByNameAndCreatedAtInput is the input to the ScanThingsByNameAndCreatedAt method.
type ScanThingsByNameAndCreatedAtInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.Thing
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingsByNameAndRangeNullableInput is the query input to GetThingsByNameAndRangeNullable.
type GetThingsByNameAndRangeNullableInput struct {
	// Name is required
	Name                    string
	RangeNullableStartingAt *strfmt.DateTime
	StartingAfter           *models.Thing
	Descending              bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingByNameAndRangeNullableNotFound is returned when the database fails to find a Thing.
type ErrThingByNameAndRangeNullableNotFound struct {
	Name          string
	RangeNullable strfmt.DateTime
}

var _ error = ErrThingByNameAndRangeNullableNotFound{}

// Error returns a description of the error.
func (e ErrThingByNameAndRangeNullableNotFound) Error() string {
	return "could not find Thing"
}

// ScanThingsByNameAndRangeNullableInput is the input to the ScanThingsByNameAndRangeNullable method.
type ScanThingsByNameAndRangeNullableInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.Thing
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingsByHashNullableAndNameInput is the query input to GetThingsByHashNullableAndName.
type GetThingsByHashNullableAndNameInput struct {
	// HashNullable is required
	HashNullable   string
	NameStartingAt *string
	StartingAfter  *models.Thing
	Descending     bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingByHashNullableAndNameNotFound is returned when the database fails to find a Thing.
type ErrThingByHashNullableAndNameNotFound struct {
	HashNullable string
	Name         string
}

var _ error = ErrThingByHashNullableAndNameNotFound{}

// Error returns a description of the error.
func (e ErrThingByHashNullableAndNameNotFound) Error() string {
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

// ScanThingAllowingBatchWritessInput is the input to the ScanThingAllowingBatchWritess method.
type ScanThingAllowingBatchWritessInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingAllowingBatchWrites
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingAllowingBatchWritessByNameAndVersionInput is the query input to GetThingAllowingBatchWritessByNameAndVersion.
type GetThingAllowingBatchWritessByNameAndVersionInput struct {
	// Name is required
	Name              string
	VersionStartingAt *int64
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.ThingAllowingBatchWrites
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingAllowingBatchWritesNotFound is returned when the database fails to find a ThingAllowingBatchWrites.
type ErrThingAllowingBatchWritesNotFound struct {
	Name    string
	Version int64
}

var _ error = ErrThingAllowingBatchWritesNotFound{}

// Error returns a description of the error.
func (e ErrThingAllowingBatchWritesNotFound) Error() string {
	return "could not find ThingAllowingBatchWrites"
}

// ErrThingAllowingBatchWritesAlreadyExists is returned when trying to overwrite a ThingAllowingBatchWrites.
type ErrThingAllowingBatchWritesAlreadyExists struct {
	Name    string
	Version int64
}

var _ error = ErrThingAllowingBatchWritesAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingAllowingBatchWritesAlreadyExists) Error() string {
	return "ThingAllowingBatchWrites already exists"
}

// ScanThingAllowingBatchWritesWithCompositeAttributessInput is the input to the ScanThingAllowingBatchWritesWithCompositeAttributess method.
type ScanThingAllowingBatchWritesWithCompositeAttributessInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingAllowingBatchWritesWithCompositeAttributes
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput is the query input to GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDate.
type GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput struct {
	// Name is required
	Name string
	// ID is required
	ID             string
	DateStartingAt *strfmt.DateTime
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.ThingAllowingBatchWritesWithCompositeAttributes
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingAllowingBatchWritesWithCompositeAttributesNotFound is returned when the database fails to find a ThingAllowingBatchWritesWithCompositeAttributes.
type ErrThingAllowingBatchWritesWithCompositeAttributesNotFound struct {
	Name string
	ID   string
	Date strfmt.DateTime
}

var _ error = ErrThingAllowingBatchWritesWithCompositeAttributesNotFound{}

// Error returns a description of the error.
func (e ErrThingAllowingBatchWritesWithCompositeAttributesNotFound) Error() string {
	return "could not find ThingAllowingBatchWritesWithCompositeAttributes"
}

// ErrThingAllowingBatchWritesWithCompositeAttributesAlreadyExists is returned when trying to overwrite a ThingAllowingBatchWritesWithCompositeAttributes.
type ErrThingAllowingBatchWritesWithCompositeAttributesAlreadyExists struct {
	NameID string
	Date   strfmt.DateTime
}

var _ error = ErrThingAllowingBatchWritesWithCompositeAttributesAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingAllowingBatchWritesWithCompositeAttributesAlreadyExists) Error() string {
	return "ThingAllowingBatchWritesWithCompositeAttributes already exists"
}

// ScanThingWithAdditionalAttributessInput is the input to the ScanThingWithAdditionalAttributess method.
type ScanThingWithAdditionalAttributessInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithAdditionalAttributes
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ThingWithAdditionalAttributesByNameAndVersionFilterableAttribute represents the fields we can apply filters to for queries on this index
type ThingWithAdditionalAttributesByNameAndVersionFilterableAttribute string

const ThingWithAdditionalAttributesAdditionalBAttribute ThingWithAdditionalAttributesByNameAndVersionFilterableAttribute = "additionalBAttribute"
const ThingWithAdditionalAttributesAdditionalNAttribute ThingWithAdditionalAttributesByNameAndVersionFilterableAttribute = "additionalNAttribute"
const ThingWithAdditionalAttributesAdditionalSAttribute ThingWithAdditionalAttributesByNameAndVersionFilterableAttribute = "additionalSAttribute"
const ThingWithAdditionalAttributesCreatedAt ThingWithAdditionalAttributesByNameAndVersionFilterableAttribute = "createdAt"
const ThingWithAdditionalAttributesHashNullable ThingWithAdditionalAttributesByNameAndVersionFilterableAttribute = "hashNullable"
const ThingWithAdditionalAttributesID ThingWithAdditionalAttributesByNameAndVersionFilterableAttribute = "id"
const ThingWithAdditionalAttributesRangeNullable ThingWithAdditionalAttributesByNameAndVersionFilterableAttribute = "rangeNullable"

// ThingWithAdditionalAttributesByNameAndVersionFilterValues represents a filter on a particular field to be included in the query
type ThingWithAdditionalAttributesByNameAndVersionFilterValues struct {
	// AttributeName is the attibute we are attempting to apply the filter to
	AttributeName ThingWithAdditionalAttributesByNameAndVersionFilterableAttribute
	// AttributeValues is an optional parameter to be used when we want to compare the attibute to a single value or multiple values
	AttributeValues []interface{}
}

// GetThingWithAdditionalAttributessByNameAndVersionInput is the query input to GetThingWithAdditionalAttributessByNameAndVersion.
type GetThingWithAdditionalAttributessByNameAndVersionInput struct {
	// Name is required
	Name              string
	VersionStartingAt *int64
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.ThingWithAdditionalAttributes
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// FilterValues is an optional array of filters to apply on various table attributes
	FilterValues []ThingWithAdditionalAttributesByNameAndVersionFilterValues
	// FilterExpression is the filter expression to be applied to our fitlered attributes
	// when referencing an attribute use #ATTRIBUTE_NAME
	// ex: if the attribute is called "created_at" in its wag definition use #CREATED_AT
	// when referencing one of the given values use :{attribute_name}_value0, :{attribute_name}_value1, etc.
	// ex: if the attribute is called "created_at" in its wag definition use :created_at_value0, created_at_value1, etc.
	// see https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.KeyConditionExpressions
	// for guidance on building expressions
	FilterExpression string
}

// ErrThingWithAdditionalAttributesNotFound is returned when the database fails to find a ThingWithAdditionalAttributes.
type ErrThingWithAdditionalAttributesNotFound struct {
	Name    string
	Version int64
}

var _ error = ErrThingWithAdditionalAttributesNotFound{}

// Error returns a description of the error.
func (e ErrThingWithAdditionalAttributesNotFound) Error() string {
	return "could not find ThingWithAdditionalAttributes"
}

// ErrThingWithAdditionalAttributesByIDNotFound is returned when the database fails to find a ThingWithAdditionalAttributes.
type ErrThingWithAdditionalAttributesByIDNotFound struct {
	ID string
}

var _ error = ErrThingWithAdditionalAttributesByIDNotFound{}

// Error returns a description of the error.
func (e ErrThingWithAdditionalAttributesByIDNotFound) Error() string {
	return "could not find ThingWithAdditionalAttributes"
}

// ScanThingWithAdditionalAttributessByIDInput is the input to the ScanThingWithAdditionalAttributessByID method.
type ScanThingWithAdditionalAttributessByIDInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithAdditionalAttributes
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingWithAdditionalAttributessByNameAndCreatedAtInput is the query input to GetThingWithAdditionalAttributessByNameAndCreatedAt.
type GetThingWithAdditionalAttributessByNameAndCreatedAtInput struct {
	// Name is required
	Name                string
	CreatedAtStartingAt *strfmt.DateTime
	StartingAfter       *models.ThingWithAdditionalAttributes
	Descending          bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithAdditionalAttributesByNameAndCreatedAtNotFound is returned when the database fails to find a ThingWithAdditionalAttributes.
type ErrThingWithAdditionalAttributesByNameAndCreatedAtNotFound struct {
	Name      string
	CreatedAt strfmt.DateTime
}

var _ error = ErrThingWithAdditionalAttributesByNameAndCreatedAtNotFound{}

// Error returns a description of the error.
func (e ErrThingWithAdditionalAttributesByNameAndCreatedAtNotFound) Error() string {
	return "could not find ThingWithAdditionalAttributes"
}

// ScanThingWithAdditionalAttributessByNameAndCreatedAtInput is the input to the ScanThingWithAdditionalAttributessByNameAndCreatedAt method.
type ScanThingWithAdditionalAttributessByNameAndCreatedAtInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithAdditionalAttributes
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingWithAdditionalAttributessByNameAndRangeNullableInput is the query input to GetThingWithAdditionalAttributessByNameAndRangeNullable.
type GetThingWithAdditionalAttributessByNameAndRangeNullableInput struct {
	// Name is required
	Name                    string
	RangeNullableStartingAt *strfmt.DateTime
	StartingAfter           *models.ThingWithAdditionalAttributes
	Descending              bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithAdditionalAttributesByNameAndRangeNullableNotFound is returned when the database fails to find a ThingWithAdditionalAttributes.
type ErrThingWithAdditionalAttributesByNameAndRangeNullableNotFound struct {
	Name          string
	RangeNullable strfmt.DateTime
}

var _ error = ErrThingWithAdditionalAttributesByNameAndRangeNullableNotFound{}

// Error returns a description of the error.
func (e ErrThingWithAdditionalAttributesByNameAndRangeNullableNotFound) Error() string {
	return "could not find ThingWithAdditionalAttributes"
}

// ScanThingWithAdditionalAttributessByNameAndRangeNullableInput is the input to the ScanThingWithAdditionalAttributessByNameAndRangeNullable method.
type ScanThingWithAdditionalAttributessByNameAndRangeNullableInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithAdditionalAttributes
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingWithAdditionalAttributessByHashNullableAndNameInput is the query input to GetThingWithAdditionalAttributessByHashNullableAndName.
type GetThingWithAdditionalAttributessByHashNullableAndNameInput struct {
	// HashNullable is required
	HashNullable   string
	NameStartingAt *string
	StartingAfter  *models.ThingWithAdditionalAttributes
	Descending     bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithAdditionalAttributesByHashNullableAndNameNotFound is returned when the database fails to find a ThingWithAdditionalAttributes.
type ErrThingWithAdditionalAttributesByHashNullableAndNameNotFound struct {
	HashNullable string
	Name         string
}

var _ error = ErrThingWithAdditionalAttributesByHashNullableAndNameNotFound{}

// Error returns a description of the error.
func (e ErrThingWithAdditionalAttributesByHashNullableAndNameNotFound) Error() string {
	return "could not find ThingWithAdditionalAttributes"
}

// ErrThingWithAdditionalAttributesAlreadyExists is returned when trying to overwrite a ThingWithAdditionalAttributes.
type ErrThingWithAdditionalAttributesAlreadyExists struct {
	Name    string
	Version int64
}

var _ error = ErrThingWithAdditionalAttributesAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithAdditionalAttributesAlreadyExists) Error() string {
	return "ThingWithAdditionalAttributes already exists"
}

// ScanThingWithCompositeAttributessInput is the input to the ScanThingWithCompositeAttributess method.
type ScanThingWithCompositeAttributessInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithCompositeAttributes
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ThingWithCompositeAttributesByNameBranchAndDateFilterableAttribute represents the fields we can apply filters to for queries on this index
type ThingWithCompositeAttributesByNameBranchAndDateFilterableAttribute string

const ThingWithCompositeAttributesVersion ThingWithCompositeAttributesByNameBranchAndDateFilterableAttribute = "version"

// ThingWithCompositeAttributesByNameBranchAndDateFilterValues represents a filter on a particular field to be included in the query
type ThingWithCompositeAttributesByNameBranchAndDateFilterValues struct {
	// AttributeName is the attibute we are attempting to apply the filter to
	AttributeName ThingWithCompositeAttributesByNameBranchAndDateFilterableAttribute
	// AttributeValues is an optional parameter to be used when we want to compare the attibute to a single value or multiple values
	AttributeValues []interface{}
}

// GetThingWithCompositeAttributessByNameBranchAndDateInput is the query input to GetThingWithCompositeAttributessByNameBranchAndDate.
type GetThingWithCompositeAttributessByNameBranchAndDateInput struct {
	// Name is required
	Name string
	// Branch is required
	Branch         string
	DateStartingAt *strfmt.DateTime
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.ThingWithCompositeAttributes
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// FilterValues is an optional array of filters to apply on various table attributes
	FilterValues []ThingWithCompositeAttributesByNameBranchAndDateFilterValues
	// FilterExpression is the filter expression to be applied to our fitlered attributes
	// when referencing an attribute use #ATTRIBUTE_NAME
	// ex: if the attribute is called "created_at" in its wag definition use #CREATED_AT
	// when referencing one of the given values use :{attribute_name}_value0, :{attribute_name}_value1, etc.
	// ex: if the attribute is called "created_at" in its wag definition use :created_at_value0, created_at_value1, etc.
	// see https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.KeyConditionExpressions
	// for guidance on building expressions
	FilterExpression string
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
	// Name is required
	Name string
	// Version is required
	Version        int64
	DateStartingAt *strfmt.DateTime
	StartingAfter  *models.ThingWithCompositeAttributes
	Descending     bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
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

// ScanThingWithCompositeAttributessByNameVersionAndDateInput is the input to the ScanThingWithCompositeAttributessByNameVersionAndDate method.
type ScanThingWithCompositeAttributessByNameVersionAndDateInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithCompositeAttributes
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
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

// ScanThingWithCompositeEnumAttributessInput is the input to the ScanThingWithCompositeEnumAttributess method.
type ScanThingWithCompositeEnumAttributessInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithCompositeEnumAttributes
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingWithCompositeEnumAttributessByNameBranchAndDateInput is the query input to GetThingWithCompositeEnumAttributessByNameBranchAndDate.
type GetThingWithCompositeEnumAttributessByNameBranchAndDateInput struct {
	// Name is required
	Name string
	// BranchID is required
	BranchID       models.Branch
	DateStartingAt *strfmt.DateTime
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.ThingWithCompositeEnumAttributes
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
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

// ScanThingWithDateGSIsInput is the input to the ScanThingWithDateGSIs method.
type ScanThingWithDateGSIsInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithDateGSI
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ErrThingWithDateGSINotFound is returned when the database fails to find a ThingWithDateGSI.
type ErrThingWithDateGSINotFound struct {
	DateH strfmt.Date
}

var _ error = ErrThingWithDateGSINotFound{}

// Error returns a description of the error.
func (e ErrThingWithDateGSINotFound) Error() string {
	return "could not find ThingWithDateGSI"
}

// GetThingWithDateGSIsByIDAndDateRInput is the query input to GetThingWithDateGSIsByIDAndDateR.
type GetThingWithDateGSIsByIDAndDateRInput struct {
	// ID is required
	ID              string
	DateRStartingAt *strfmt.Date
	StartingAfter   *models.ThingWithDateGSI
	Descending      bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithDateGSIByIDAndDateRNotFound is returned when the database fails to find a ThingWithDateGSI.
type ErrThingWithDateGSIByIDAndDateRNotFound struct {
	ID    string
	DateR strfmt.Date
}

var _ error = ErrThingWithDateGSIByIDAndDateRNotFound{}

// Error returns a description of the error.
func (e ErrThingWithDateGSIByIDAndDateRNotFound) Error() string {
	return "could not find ThingWithDateGSI"
}

// GetThingWithDateGSIsByDateHAndIDInput is the query input to GetThingWithDateGSIsByDateHAndID.
type GetThingWithDateGSIsByDateHAndIDInput struct {
	// DateH is required
	DateH         strfmt.Date
	IDStartingAt  *string
	StartingAfter *models.ThingWithDateGSI
	Descending    bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithDateGSIByDateHAndIDNotFound is returned when the database fails to find a ThingWithDateGSI.
type ErrThingWithDateGSIByDateHAndIDNotFound struct {
	DateH strfmt.Date
	ID    string
}

var _ error = ErrThingWithDateGSIByDateHAndIDNotFound{}

// Error returns a description of the error.
func (e ErrThingWithDateGSIByDateHAndIDNotFound) Error() string {
	return "could not find ThingWithDateGSI"
}

// ErrThingWithDateGSIAlreadyExists is returned when trying to overwrite a ThingWithDateGSI.
type ErrThingWithDateGSIAlreadyExists struct {
	DateH strfmt.Date
}

var _ error = ErrThingWithDateGSIAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithDateGSIAlreadyExists) Error() string {
	return "ThingWithDateGSI already exists"
}

// ScanThingWithDateRangesInput is the input to the ScanThingWithDateRanges method.
type ScanThingWithDateRangesInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithDateRange
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingWithDateRangesByNameAndDateInput is the query input to GetThingWithDateRangesByNameAndDate.
type GetThingWithDateRangesByNameAndDateInput struct {
	// Name is required
	Name           string
	DateStartingAt *strfmt.DateTime
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.ThingWithDateRange
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
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

// ScanThingWithDateRangeKeysInput is the input to the ScanThingWithDateRangeKeys method.
type ScanThingWithDateRangeKeysInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithDateRangeKey
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingWithDateRangeKeysByIDAndDateInput is the query input to GetThingWithDateRangeKeysByIDAndDate.
type GetThingWithDateRangeKeysByIDAndDateInput struct {
	// ID is required
	ID             string
	DateStartingAt *strfmt.Date
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.ThingWithDateRangeKey
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithDateRangeKeyNotFound is returned when the database fails to find a ThingWithDateRangeKey.
type ErrThingWithDateRangeKeyNotFound struct {
	ID   string
	Date strfmt.Date
}

var _ error = ErrThingWithDateRangeKeyNotFound{}

// Error returns a description of the error.
func (e ErrThingWithDateRangeKeyNotFound) Error() string {
	return "could not find ThingWithDateRangeKey"
}

// ErrThingWithDateRangeKeyAlreadyExists is returned when trying to overwrite a ThingWithDateRangeKey.
type ErrThingWithDateRangeKeyAlreadyExists struct {
	ID   string
	Date strfmt.Date
}

var _ error = ErrThingWithDateRangeKeyAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithDateRangeKeyAlreadyExists) Error() string {
	return "ThingWithDateRangeKey already exists"
}

// ScanThingWithDateTimeCompositesInput is the input to the ScanThingWithDateTimeComposites method.
type ScanThingWithDateTimeCompositesInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithDateTimeComposite
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput is the query input to GetThingWithDateTimeCompositesByTypeIDAndCreatedResource.
type GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput struct {
	// Type is required
	Type string
	// ID is required
	ID         string
	StartingAt *CreatedResource
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.ThingWithDateTimeComposite
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
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

// CreatedResource struct.
type CreatedResource struct {
	Created  strfmt.DateTime
	Resource string
}

// ScanThingWithDatetimeGSIsInput is the input to the ScanThingWithDatetimeGSIs method.
type ScanThingWithDatetimeGSIsInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithDatetimeGSI
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ErrThingWithDatetimeGSINotFound is returned when the database fails to find a ThingWithDatetimeGSI.
type ErrThingWithDatetimeGSINotFound struct {
	ID string
}

var _ error = ErrThingWithDatetimeGSINotFound{}

// Error returns a description of the error.
func (e ErrThingWithDatetimeGSINotFound) Error() string {
	return "could not find ThingWithDatetimeGSI"
}

// GetThingWithDatetimeGSIsByDatetimeAndIDInput is the query input to GetThingWithDatetimeGSIsByDatetimeAndID.
type GetThingWithDatetimeGSIsByDatetimeAndIDInput struct {
	// Datetime is required
	Datetime      strfmt.DateTime
	IDStartingAt  *string
	StartingAfter *models.ThingWithDatetimeGSI
	Descending    bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithDatetimeGSIByDatetimeAndIDNotFound is returned when the database fails to find a ThingWithDatetimeGSI.
type ErrThingWithDatetimeGSIByDatetimeAndIDNotFound struct {
	Datetime strfmt.DateTime
	ID       string
}

var _ error = ErrThingWithDatetimeGSIByDatetimeAndIDNotFound{}

// Error returns a description of the error.
func (e ErrThingWithDatetimeGSIByDatetimeAndIDNotFound) Error() string {
	return "could not find ThingWithDatetimeGSI"
}

// ScanThingWithDatetimeGSIsByDatetimeAndIDInput is the input to the ScanThingWithDatetimeGSIsByDatetimeAndID method.
type ScanThingWithDatetimeGSIsByDatetimeAndIDInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithDatetimeGSI
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ErrThingWithDatetimeGSIAlreadyExists is returned when trying to overwrite a ThingWithDatetimeGSI.
type ErrThingWithDatetimeGSIAlreadyExists struct {
	ID string
}

var _ error = ErrThingWithDatetimeGSIAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithDatetimeGSIAlreadyExists) Error() string {
	return "ThingWithDatetimeGSI already exists"
}

// ScanThingWithEnumHashKeysInput is the input to the ScanThingWithEnumHashKeys method.
type ScanThingWithEnumHashKeysInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithEnumHashKey
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ThingWithEnumHashKeyByBranchAndDateFilterableAttribute represents the fields we can apply filters to for queries on this index
type ThingWithEnumHashKeyByBranchAndDateFilterableAttribute string

const ThingWithEnumHashKeyDate2 ThingWithEnumHashKeyByBranchAndDateFilterableAttribute = "date2"

// ThingWithEnumHashKeyByBranchAndDateFilterValues represents a filter on a particular field to be included in the query
type ThingWithEnumHashKeyByBranchAndDateFilterValues struct {
	// AttributeName is the attibute we are attempting to apply the filter to
	AttributeName ThingWithEnumHashKeyByBranchAndDateFilterableAttribute
	// AttributeValues is an optional parameter to be used when we want to compare the attibute to a single value or multiple values
	AttributeValues []interface{}
}

// GetThingWithEnumHashKeysByBranchAndDateInput is the query input to GetThingWithEnumHashKeysByBranchAndDate.
type GetThingWithEnumHashKeysByBranchAndDateInput struct {
	// Branch is required
	Branch         models.Branch
	DateStartingAt *strfmt.DateTime
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.ThingWithEnumHashKey
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// FilterValues is an optional array of filters to apply on various table attributes
	FilterValues []ThingWithEnumHashKeyByBranchAndDateFilterValues
	// FilterExpression is the filter expression to be applied to our fitlered attributes
	// when referencing an attribute use #ATTRIBUTE_NAME
	// ex: if the attribute is called "created_at" in its wag definition use #CREATED_AT
	// when referencing one of the given values use :{attribute_name}_value0, :{attribute_name}_value1, etc.
	// ex: if the attribute is called "created_at" in its wag definition use :created_at_value0, created_at_value1, etc.
	// see https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.KeyConditionExpressions
	// for guidance on building expressions
	FilterExpression string
}

// ErrThingWithEnumHashKeyNotFound is returned when the database fails to find a ThingWithEnumHashKey.
type ErrThingWithEnumHashKeyNotFound struct {
	Branch models.Branch
	Date   strfmt.DateTime
}

var _ error = ErrThingWithEnumHashKeyNotFound{}

// Error returns a description of the error.
func (e ErrThingWithEnumHashKeyNotFound) Error() string {
	return "could not find ThingWithEnumHashKey"
}

// GetThingWithEnumHashKeysByBranchAndDate2Input is the query input to GetThingWithEnumHashKeysByBranchAndDate2.
type GetThingWithEnumHashKeysByBranchAndDate2Input struct {
	// Branch is required
	Branch          models.Branch
	Date2StartingAt *strfmt.DateTime
	StartingAfter   *models.ThingWithEnumHashKey
	Descending      bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithEnumHashKeyByBranchAndDate2NotFound is returned when the database fails to find a ThingWithEnumHashKey.
type ErrThingWithEnumHashKeyByBranchAndDate2NotFound struct {
	Branch models.Branch
	Date2  strfmt.DateTime
}

var _ error = ErrThingWithEnumHashKeyByBranchAndDate2NotFound{}

// Error returns a description of the error.
func (e ErrThingWithEnumHashKeyByBranchAndDate2NotFound) Error() string {
	return "could not find ThingWithEnumHashKey"
}

// ScanThingWithEnumHashKeysByBranchAndDate2Input is the input to the ScanThingWithEnumHashKeysByBranchAndDate2 method.
type ScanThingWithEnumHashKeysByBranchAndDate2Input struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithEnumHashKey
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ErrThingWithEnumHashKeyAlreadyExists is returned when trying to overwrite a ThingWithEnumHashKey.
type ErrThingWithEnumHashKeyAlreadyExists struct {
	Branch models.Branch
	Date   strfmt.DateTime
}

var _ error = ErrThingWithEnumHashKeyAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithEnumHashKeyAlreadyExists) Error() string {
	return "ThingWithEnumHashKey already exists"
}

// ScanThingWithMatchingKeyssInput is the input to the ScanThingWithMatchingKeyss method.
type ScanThingWithMatchingKeyssInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithMatchingKeys
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ThingWithMatchingKeysByBearAndAssocTypeIDFilterableAttribute represents the fields we can apply filters to for queries on this index
type ThingWithMatchingKeysByBearAndAssocTypeIDFilterableAttribute string

const ThingWithMatchingKeysCreated ThingWithMatchingKeysByBearAndAssocTypeIDFilterableAttribute = "created"

// ThingWithMatchingKeysByBearAndAssocTypeIDFilterValues represents a filter on a particular field to be included in the query
type ThingWithMatchingKeysByBearAndAssocTypeIDFilterValues struct {
	// AttributeName is the attibute we are attempting to apply the filter to
	AttributeName ThingWithMatchingKeysByBearAndAssocTypeIDFilterableAttribute
	// AttributeValues is an optional parameter to be used when we want to compare the attibute to a single value or multiple values
	AttributeValues []interface{}
}

// GetThingWithMatchingKeyssByBearAndAssocTypeIDInput is the query input to GetThingWithMatchingKeyssByBearAndAssocTypeID.
type GetThingWithMatchingKeyssByBearAndAssocTypeIDInput struct {
	// Bear is required
	Bear       string
	StartingAt *AssocTypeAssocID
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.ThingWithMatchingKeys
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// FilterValues is an optional array of filters to apply on various table attributes
	FilterValues []ThingWithMatchingKeysByBearAndAssocTypeIDFilterValues
	// FilterExpression is the filter expression to be applied to our fitlered attributes
	// when referencing an attribute use #ATTRIBUTE_NAME
	// ex: if the attribute is called "created_at" in its wag definition use #CREATED_AT
	// when referencing one of the given values use :{attribute_name}_value0, :{attribute_name}_value1, etc.
	// ex: if the attribute is called "created_at" in its wag definition use :created_at_value0, created_at_value1, etc.
	// see https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Query.html#Query.KeyConditionExpressions
	// for guidance on building expressions
	FilterExpression string
}

// ErrThingWithMatchingKeysNotFound is returned when the database fails to find a ThingWithMatchingKeys.
type ErrThingWithMatchingKeysNotFound struct {
	Bear      string
	AssocType string
	AssocID   string
}

var _ error = ErrThingWithMatchingKeysNotFound{}

// Error returns a description of the error.
func (e ErrThingWithMatchingKeysNotFound) Error() string {
	return "could not find ThingWithMatchingKeys"
}

// GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput is the query input to GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBear.
type GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput struct {
	// AssocType is required
	AssocType string
	// AssocID is required
	AssocID       string
	StartingAt    *CreatedBear
	StartingAfter *models.ThingWithMatchingKeys
	Descending    bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithMatchingKeysByAssocTypeIDAndCreatedBearNotFound is returned when the database fails to find a ThingWithMatchingKeys.
type ErrThingWithMatchingKeysByAssocTypeIDAndCreatedBearNotFound struct {
	AssocType string
	AssocID   string
	Created   strfmt.DateTime
	Bear      string
}

var _ error = ErrThingWithMatchingKeysByAssocTypeIDAndCreatedBearNotFound{}

// Error returns a description of the error.
func (e ErrThingWithMatchingKeysByAssocTypeIDAndCreatedBearNotFound) Error() string {
	return "could not find ThingWithMatchingKeys"
}

// ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput is the input to the ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear method.
type ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithMatchingKeys
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// AssocTypeAssocID struct.
type AssocTypeAssocID struct {
	AssocType string
	AssocID   string
}

// CreatedBear struct.
type CreatedBear struct {
	Created strfmt.DateTime
	Bear    string
}

// ScanThingWithMultiUseCompositeAttributesInput is the input to the ScanThingWithMultiUseCompositeAttributes method.
type ScanThingWithMultiUseCompositeAttributesInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithMultiUseCompositeAttribute
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ErrThingWithMultiUseCompositeAttributeNotFound is returned when the database fails to find a ThingWithMultiUseCompositeAttribute.
type ErrThingWithMultiUseCompositeAttributeNotFound struct {
	One string
}

var _ error = ErrThingWithMultiUseCompositeAttributeNotFound{}

// Error returns a description of the error.
func (e ErrThingWithMultiUseCompositeAttributeNotFound) Error() string {
	return "could not find ThingWithMultiUseCompositeAttribute"
}

// GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput is the query input to GetThingWithMultiUseCompositeAttributesByThreeAndOneTwo.
type GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput struct {
	// Three is required
	Three         string
	StartingAt    *OneTwo
	StartingAfter *models.ThingWithMultiUseCompositeAttribute
	Descending    bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithMultiUseCompositeAttributeByThreeAndOneTwoNotFound is returned when the database fails to find a ThingWithMultiUseCompositeAttribute.
type ErrThingWithMultiUseCompositeAttributeByThreeAndOneTwoNotFound struct {
	Three string
	One   string
	Two   string
}

var _ error = ErrThingWithMultiUseCompositeAttributeByThreeAndOneTwoNotFound{}

// Error returns a description of the error.
func (e ErrThingWithMultiUseCompositeAttributeByThreeAndOneTwoNotFound) Error() string {
	return "could not find ThingWithMultiUseCompositeAttribute"
}

// ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput is the input to the ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwo method.
type ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithMultiUseCompositeAttribute
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput is the query input to GetThingWithMultiUseCompositeAttributesByFourAndOneTwo.
type GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput struct {
	// Four is required
	Four          string
	StartingAt    *OneTwo
	StartingAfter *models.ThingWithMultiUseCompositeAttribute
	Descending    bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithMultiUseCompositeAttributeByFourAndOneTwoNotFound is returned when the database fails to find a ThingWithMultiUseCompositeAttribute.
type ErrThingWithMultiUseCompositeAttributeByFourAndOneTwoNotFound struct {
	Four string
	One  string
	Two  string
}

var _ error = ErrThingWithMultiUseCompositeAttributeByFourAndOneTwoNotFound{}

// Error returns a description of the error.
func (e ErrThingWithMultiUseCompositeAttributeByFourAndOneTwoNotFound) Error() string {
	return "could not find ThingWithMultiUseCompositeAttribute"
}

// ScanThingWithMultiUseCompositeAttributesByFourAndOneTwoInput is the input to the ScanThingWithMultiUseCompositeAttributesByFourAndOneTwo method.
type ScanThingWithMultiUseCompositeAttributesByFourAndOneTwoInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithMultiUseCompositeAttribute
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// OneTwo struct.
type OneTwo struct {
	One string
	Two string
}

// ScanThingWithRequiredCompositePropertiesAndKeysOnlysInput is the input to the ScanThingWithRequiredCompositePropertiesAndKeysOnlys method.
type ScanThingWithRequiredCompositePropertiesAndKeysOnlysInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithRequiredCompositePropertiesAndKeysOnly
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ErrThingWithRequiredCompositePropertiesAndKeysOnlyNotFound is returned when the database fails to find a ThingWithRequiredCompositePropertiesAndKeysOnly.
type ErrThingWithRequiredCompositePropertiesAndKeysOnlyNotFound struct {
	PropertyThree string
}

var _ error = ErrThingWithRequiredCompositePropertiesAndKeysOnlyNotFound{}

// Error returns a description of the error.
func (e ErrThingWithRequiredCompositePropertiesAndKeysOnlyNotFound) Error() string {
	return "could not find ThingWithRequiredCompositePropertiesAndKeysOnly"
}

// GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput is the query input to GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree.
type GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput struct {
	// PropertyOne is required
	PropertyOne string
	// PropertyTwo is required
	PropertyTwo             string
	PropertyThreeStartingAt *string
	StartingAfter           *models.ThingWithRequiredCompositePropertiesAndKeysOnly
	Descending              bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithRequiredCompositePropertiesAndKeysOnlyByPropertyOneAndTwoAndPropertyThreeNotFound is returned when the database fails to find a ThingWithRequiredCompositePropertiesAndKeysOnly.
type ErrThingWithRequiredCompositePropertiesAndKeysOnlyByPropertyOneAndTwoAndPropertyThreeNotFound struct {
	PropertyOne   string
	PropertyTwo   string
	PropertyThree string
}

var _ error = ErrThingWithRequiredCompositePropertiesAndKeysOnlyByPropertyOneAndTwoAndPropertyThreeNotFound{}

// Error returns a description of the error.
func (e ErrThingWithRequiredCompositePropertiesAndKeysOnlyByPropertyOneAndTwoAndPropertyThreeNotFound) Error() string {
	return "could not find ThingWithRequiredCompositePropertiesAndKeysOnly"
}

// ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput is the input to the ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree method.
type ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithRequiredCompositePropertiesAndKeysOnly
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ScanThingWithRequiredFieldssInput is the input to the ScanThingWithRequiredFieldss method.
type ScanThingWithRequiredFieldssInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithRequiredFields
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
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

// ScanThingWithRequiredFields2sInput is the input to the ScanThingWithRequiredFields2s method.
type ScanThingWithRequiredFields2sInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithRequiredFields2
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// GetThingWithRequiredFields2sByNameAndIDInput is the query input to GetThingWithRequiredFields2sByNameAndID.
type GetThingWithRequiredFields2sByNameAndIDInput struct {
	// Name is required
	Name         string
	IDStartingAt *string
	// StartingAfter is a required specification of an exclusive starting point.
	StartingAfter *models.ThingWithRequiredFields2
	Descending    bool
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithRequiredFields2NotFound is returned when the database fails to find a ThingWithRequiredFields2.
type ErrThingWithRequiredFields2NotFound struct {
	Name string
	ID   string
}

var _ error = ErrThingWithRequiredFields2NotFound{}

// Error returns a description of the error.
func (e ErrThingWithRequiredFields2NotFound) Error() string {
	return "could not find ThingWithRequiredFields2"
}

// ErrThingWithRequiredFields2AlreadyExists is returned when trying to overwrite a ThingWithRequiredFields2.
type ErrThingWithRequiredFields2AlreadyExists struct {
	Name string
	ID   string
}

var _ error = ErrThingWithRequiredFields2AlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithRequiredFields2AlreadyExists) Error() string {
	return "ThingWithRequiredFields2 already exists"
}

// ScanThingWithTransactMultipleGSIsInput is the input to the ScanThingWithTransactMultipleGSIs method.
type ScanThingWithTransactMultipleGSIsInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithTransactMultipleGSI
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ErrThingWithTransactMultipleGSINotFound is returned when the database fails to find a ThingWithTransactMultipleGSI.
type ErrThingWithTransactMultipleGSINotFound struct {
	DateH strfmt.Date
}

var _ error = ErrThingWithTransactMultipleGSINotFound{}

// Error returns a description of the error.
func (e ErrThingWithTransactMultipleGSINotFound) Error() string {
	return "could not find ThingWithTransactMultipleGSI"
}

// GetThingWithTransactMultipleGSIsByIDAndDateRInput is the query input to GetThingWithTransactMultipleGSIsByIDAndDateR.
type GetThingWithTransactMultipleGSIsByIDAndDateRInput struct {
	// ID is required
	ID              string
	DateRStartingAt *strfmt.Date
	StartingAfter   *models.ThingWithTransactMultipleGSI
	Descending      bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithTransactMultipleGSIByIDAndDateRNotFound is returned when the database fails to find a ThingWithTransactMultipleGSI.
type ErrThingWithTransactMultipleGSIByIDAndDateRNotFound struct {
	ID    string
	DateR strfmt.Date
}

var _ error = ErrThingWithTransactMultipleGSIByIDAndDateRNotFound{}

// Error returns a description of the error.
func (e ErrThingWithTransactMultipleGSIByIDAndDateRNotFound) Error() string {
	return "could not find ThingWithTransactMultipleGSI"
}

// GetThingWithTransactMultipleGSIsByDateHAndIDInput is the query input to GetThingWithTransactMultipleGSIsByDateHAndID.
type GetThingWithTransactMultipleGSIsByDateHAndIDInput struct {
	// DateH is required
	DateH         strfmt.Date
	IDStartingAt  *string
	StartingAfter *models.ThingWithTransactMultipleGSI
	Descending    bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
}

// ErrThingWithTransactMultipleGSIByDateHAndIDNotFound is returned when the database fails to find a ThingWithTransactMultipleGSI.
type ErrThingWithTransactMultipleGSIByDateHAndIDNotFound struct {
	DateH strfmt.Date
	ID    string
}

var _ error = ErrThingWithTransactMultipleGSIByDateHAndIDNotFound{}

// Error returns a description of the error.
func (e ErrThingWithTransactMultipleGSIByDateHAndIDNotFound) Error() string {
	return "could not find ThingWithTransactMultipleGSI"
}

// ErrThingWithTransactMultipleGSIAlreadyExists is returned when trying to overwrite a ThingWithTransactMultipleGSI.
type ErrThingWithTransactMultipleGSIAlreadyExists struct {
	DateH strfmt.Date
}

var _ error = ErrThingWithTransactMultipleGSIAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithTransactMultipleGSIAlreadyExists) Error() string {
	return "ThingWithTransactMultipleGSI already exists"
}

// ScanThingWithTransactionsInput is the input to the ScanThingWithTransactions method.
type ScanThingWithTransactionsInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithTransaction
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ErrThingWithTransactionNotFound is returned when the database fails to find a ThingWithTransaction.
type ErrThingWithTransactionNotFound struct {
	Name string
}

var _ error = ErrThingWithTransactionNotFound{}

// Error returns a description of the error.
func (e ErrThingWithTransactionNotFound) Error() string {
	return "could not find ThingWithTransaction"
}

// ErrThingWithTransactionAlreadyExists is returned when trying to overwrite a ThingWithTransaction.
type ErrThingWithTransactionAlreadyExists struct {
	Name string
}

var _ error = ErrThingWithTransactionAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithTransactionAlreadyExists) Error() string {
	return "ThingWithTransaction already exists"
}

// ScanThingWithTransactionWithSimpleThingsInput is the input to the ScanThingWithTransactionWithSimpleThings method.
type ScanThingWithTransactionWithSimpleThingsInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithTransactionWithSimpleThing
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int32
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
}

// ErrThingWithTransactionWithSimpleThingNotFound is returned when the database fails to find a ThingWithTransactionWithSimpleThing.
type ErrThingWithTransactionWithSimpleThingNotFound struct {
	Name string
}

var _ error = ErrThingWithTransactionWithSimpleThingNotFound{}

// Error returns a description of the error.
func (e ErrThingWithTransactionWithSimpleThingNotFound) Error() string {
	return "could not find ThingWithTransactionWithSimpleThing"
}

// ErrThingWithTransactionWithSimpleThingAlreadyExists is returned when trying to overwrite a ThingWithTransactionWithSimpleThing.
type ErrThingWithTransactionWithSimpleThingAlreadyExists struct {
	Name string
}

var _ error = ErrThingWithTransactionWithSimpleThingAlreadyExists{}

// Error returns a description of the error.
func (e ErrThingWithTransactionWithSimpleThingAlreadyExists) Error() string {
	return "ThingWithTransactionWithSimpleThing already exists"
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
