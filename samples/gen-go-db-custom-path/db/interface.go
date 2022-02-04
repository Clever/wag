package db

import (
	"context"

	"github.com/Clever/wag/samples/v8/gen-go-db-custom-path/models"
	"github.com/go-openapi/strfmt"
	"golang.org/x/time/rate"
)

//go:generate mockgen -source=$GOFILE -destination=mock_db.go -package=db

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

	// SaveThingAllowingBatchWrites saves a ThingAllowingBatchWrites to the database.
	SaveThingAllowingBatchWrites(ctx context.Context, m models.ThingAllowingBatchWrites) error
	// SaveArrayOfThingAllowingBatchWrites batch saves all items in []ThingAllowingBatchWrites to the database.
	SaveArrayOfThingAllowingBatchWrites(ctx context.Context, ms []models.ThingAllowingBatchWrites) error
	// GetThingAllowingBatchWrites retrieves a ThingAllowingBatchWrites from the database.
	GetThingAllowingBatchWrites(ctx context.Context, name string, version int64) (*models.ThingAllowingBatchWrites, error)
	// ScanThingAllowingBatchWritess runs a scan on the ThingAllowingBatchWritess table.
	ScanThingAllowingBatchWritess(ctx context.Context, input ScanThingAllowingBatchWritessInput, fn func(m *models.ThingAllowingBatchWrites, lastThingAllowingBatchWrites bool) bool) error
	// GetThingAllowingBatchWritessByNameAndVersion retrieves a page of ThingAllowingBatchWritess from the database.
	GetThingAllowingBatchWritessByNameAndVersion(ctx context.Context, input GetThingAllowingBatchWritessByNameAndVersionInput, fn func(m *models.ThingAllowingBatchWrites, lastThingAllowingBatchWrites bool) bool) error
	// DeleteThingAllowingBatchWrites deletes a ThingAllowingBatchWrites from the database.
	DeleteThingAllowingBatchWrites(ctx context.Context, name string, version int64) error
	// GetThingAllowingBatchWritesByID retrieves a ThingAllowingBatchWrites from the database.
	GetThingAllowingBatchWritesByID(ctx context.Context, id string) (*models.ThingAllowingBatchWrites, error)
	// GetThingAllowingBatchWritessByNameAndCreatedAt retrieves a page of ThingAllowingBatchWritess from the database.
	GetThingAllowingBatchWritessByNameAndCreatedAt(ctx context.Context, input GetThingAllowingBatchWritessByNameAndCreatedAtInput, fn func(m *models.ThingAllowingBatchWrites, lastThingAllowingBatchWrites bool) bool) error

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

// ScanDeploymentsInput is the input to the ScanDeployments method.
type ScanDeploymentsInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.Deployment
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int64
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
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
	Limit *int64
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
	Limit *int64
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
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
	Limit *int64
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
	Limit *int64
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

// ErrThingAllowingBatchWritesByIDNotFound is returned when the database fails to find a ThingAllowingBatchWrites.
type ErrThingAllowingBatchWritesByIDNotFound struct {
	ID string
}

var _ error = ErrThingAllowingBatchWritesByIDNotFound{}

// Error returns a description of the error.
func (e ErrThingAllowingBatchWritesByIDNotFound) Error() string {
	return "could not find ThingAllowingBatchWrites"
}

// GetThingAllowingBatchWritessByNameAndCreatedAtInput is the query input to GetThingAllowingBatchWritessByNameAndCreatedAt.
type GetThingAllowingBatchWritessByNameAndCreatedAtInput struct {
	// Name is required
	Name                string
	CreatedAtStartingAt *strfmt.DateTime
	StartingAfter       *models.ThingAllowingBatchWrites
	Descending          bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int64
}

// ErrThingAllowingBatchWritesByNameAndCreatedAtNotFound is returned when the database fails to find a ThingAllowingBatchWrites.
type ErrThingAllowingBatchWritesByNameAndCreatedAtNotFound struct {
	Name      string
	CreatedAt strfmt.DateTime
}

var _ error = ErrThingAllowingBatchWritesByNameAndCreatedAtNotFound{}

// Error returns a description of the error.
func (e ErrThingAllowingBatchWritesByNameAndCreatedAtNotFound) Error() string {
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

// ScanThingWithCompositeAttributessInput is the input to the ScanThingWithCompositeAttributess method.
type ScanThingWithCompositeAttributessInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithCompositeAttributes
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int64
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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

// ScanThingWithDateRangesInput is the input to the ScanThingWithDateRanges method.
type ScanThingWithDateRangesInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithDateRange
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int64
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
	Limit *int64
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

// ScanThingWithDateTimeCompositesInput is the input to the ScanThingWithDateTimeComposites method.
type ScanThingWithDateTimeCompositesInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithDateTimeComposite
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int64
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
	Limit *int64
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

// ScanThingWithEnumHashKeysInput is the input to the ScanThingWithEnumHashKeys method.
type ScanThingWithEnumHashKeysInput struct {
	// StartingAfter is an optional specification of an (exclusive) starting point.
	StartingAfter *models.ThingWithEnumHashKey
	// DisableConsistentRead turns off the default behavior of running a consistent read.
	DisableConsistentRead bool
	// Limit is an optional limit of how many items to evaluate.
	Limit *int64
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
	// Limiter is an optional limit on how quickly items are scanned.
	Limiter *rate.Limiter
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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
	Limit *int64
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

// ErrThingWithUnderscoresNotFound is returned when the database fails to find a ThingWithUnderscores.
type ErrThingWithUnderscoresNotFound struct {
	IDApp string
}

var _ error = ErrThingWithUnderscoresNotFound{}

// Error returns a description of the error.
func (e ErrThingWithUnderscoresNotFound) Error() string {
	return "could not find ThingWithUnderscores"
}
