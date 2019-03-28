package dynamodb

import (
	"context"
	"database/sql/driver"
	"errors"
	"log"
	"time"

	"github.com/Clever/wag/samples/gen-go-db/models"
	"github.com/Clever/wag/samples/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
	"github.com/mailru/easyjson/jwriter"
)

// Config is used to create a new DB struct.
type Config struct {
	// DynamoDBAPI is used to communicate with DynamoDB. It is required.
	// It can be overriden on a table-by-table basis.
	DynamoDBAPI dynamodbiface.DynamoDBAPI

	// DefaultPrefix configures a prefix on all table names. It is required.
	// It can be overriden on a table-by-table basis.
	DefaultPrefix string

	// DefaultWriteCapacityUnits configures a default write capacity when creating tables. It defaults to 1.
	// It can be overriden on a table-by-table basis.
	DefaultWriteCapacityUnits int64

	// DefaultReadCapacityUnits configures a default read capacity when creating tables. It defaults to 1.
	// It can be overriden on a table-by-table basis.
	DefaultReadCapacityUnits int64
	// SimpleThingTable configuration.
	SimpleThingTable SimpleThingTable
	// TeacherSharingRuleTable configuration.
	TeacherSharingRuleTable TeacherSharingRuleTable
	// ThingTable configuration.
	ThingTable ThingTable
	// ThingWithCompositeAttributesTable configuration.
	ThingWithCompositeAttributesTable ThingWithCompositeAttributesTable
	// ThingWithCompositeEnumAttributesTable configuration.
	ThingWithCompositeEnumAttributesTable ThingWithCompositeEnumAttributesTable
	// ThingWithDateRangeTable configuration.
	ThingWithDateRangeTable ThingWithDateRangeTable
	// ThingWithDateTimeCompositeTable configuration.
	ThingWithDateTimeCompositeTable ThingWithDateTimeCompositeTable
	// ThingWithRequiredFieldsTable configuration.
	ThingWithRequiredFieldsTable ThingWithRequiredFieldsTable
	// ThingWithUnderscoresTable configuration.
	ThingWithUnderscoresTable ThingWithUnderscoresTable
}

// New creates a new DB object.
func New(config Config) (*DB, error) {
	if config.DynamoDBAPI == nil {
		return nil, errors.New("must specify DynamoDBAPI")
	}
	if config.DefaultPrefix == "" {
		return nil, errors.New("must specify DefaultPrefix")
	}

	if config.DefaultWriteCapacityUnits == 0 {
		config.DefaultWriteCapacityUnits = 1
	}
	if config.DefaultReadCapacityUnits == 0 {
		config.DefaultReadCapacityUnits = 1
	}
	// configure SimpleThing table
	simpleThingTable := config.SimpleThingTable
	if simpleThingTable.DynamoDBAPI == nil {
		simpleThingTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if simpleThingTable.Prefix == "" {
		simpleThingTable.Prefix = config.DefaultPrefix
	}
	if simpleThingTable.ReadCapacityUnits == 0 {
		simpleThingTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if simpleThingTable.WriteCapacityUnits == 0 {
		simpleThingTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure TeacherSharingRule table
	teacherSharingRuleTable := config.TeacherSharingRuleTable
	if teacherSharingRuleTable.DynamoDBAPI == nil {
		teacherSharingRuleTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if teacherSharingRuleTable.Prefix == "" {
		teacherSharingRuleTable.Prefix = config.DefaultPrefix
	}
	if teacherSharingRuleTable.ReadCapacityUnits == 0 {
		teacherSharingRuleTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if teacherSharingRuleTable.WriteCapacityUnits == 0 {
		teacherSharingRuleTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure Thing table
	thingTable := config.ThingTable
	if thingTable.DynamoDBAPI == nil {
		thingTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingTable.Prefix == "" {
		thingTable.Prefix = config.DefaultPrefix
	}
	if thingTable.ReadCapacityUnits == 0 {
		thingTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingTable.WriteCapacityUnits == 0 {
		thingTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure ThingWithCompositeAttributes table
	thingWithCompositeAttributesTable := config.ThingWithCompositeAttributesTable
	if thingWithCompositeAttributesTable.DynamoDBAPI == nil {
		thingWithCompositeAttributesTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingWithCompositeAttributesTable.Prefix == "" {
		thingWithCompositeAttributesTable.Prefix = config.DefaultPrefix
	}
	if thingWithCompositeAttributesTable.ReadCapacityUnits == 0 {
		thingWithCompositeAttributesTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingWithCompositeAttributesTable.WriteCapacityUnits == 0 {
		thingWithCompositeAttributesTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure ThingWithCompositeEnumAttributes table
	thingWithCompositeEnumAttributesTable := config.ThingWithCompositeEnumAttributesTable
	if thingWithCompositeEnumAttributesTable.DynamoDBAPI == nil {
		thingWithCompositeEnumAttributesTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingWithCompositeEnumAttributesTable.Prefix == "" {
		thingWithCompositeEnumAttributesTable.Prefix = config.DefaultPrefix
	}
	if thingWithCompositeEnumAttributesTable.ReadCapacityUnits == 0 {
		thingWithCompositeEnumAttributesTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingWithCompositeEnumAttributesTable.WriteCapacityUnits == 0 {
		thingWithCompositeEnumAttributesTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure ThingWithDateRange table
	thingWithDateRangeTable := config.ThingWithDateRangeTable
	if thingWithDateRangeTable.DynamoDBAPI == nil {
		thingWithDateRangeTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingWithDateRangeTable.Prefix == "" {
		thingWithDateRangeTable.Prefix = config.DefaultPrefix
	}
	if thingWithDateRangeTable.ReadCapacityUnits == 0 {
		thingWithDateRangeTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingWithDateRangeTable.WriteCapacityUnits == 0 {
		thingWithDateRangeTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure ThingWithDateTimeComposite table
	thingWithDateTimeCompositeTable := config.ThingWithDateTimeCompositeTable
	if thingWithDateTimeCompositeTable.DynamoDBAPI == nil {
		thingWithDateTimeCompositeTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingWithDateTimeCompositeTable.Prefix == "" {
		thingWithDateTimeCompositeTable.Prefix = config.DefaultPrefix
	}
	if thingWithDateTimeCompositeTable.ReadCapacityUnits == 0 {
		thingWithDateTimeCompositeTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingWithDateTimeCompositeTable.WriteCapacityUnits == 0 {
		thingWithDateTimeCompositeTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure ThingWithRequiredFields table
	thingWithRequiredFieldsTable := config.ThingWithRequiredFieldsTable
	if thingWithRequiredFieldsTable.DynamoDBAPI == nil {
		thingWithRequiredFieldsTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingWithRequiredFieldsTable.Prefix == "" {
		thingWithRequiredFieldsTable.Prefix = config.DefaultPrefix
	}
	if thingWithRequiredFieldsTable.ReadCapacityUnits == 0 {
		thingWithRequiredFieldsTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingWithRequiredFieldsTable.WriteCapacityUnits == 0 {
		thingWithRequiredFieldsTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure ThingWithUnderscores table
	thingWithUnderscoresTable := config.ThingWithUnderscoresTable
	if thingWithUnderscoresTable.DynamoDBAPI == nil {
		thingWithUnderscoresTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingWithUnderscoresTable.Prefix == "" {
		thingWithUnderscoresTable.Prefix = config.DefaultPrefix
	}
	if thingWithUnderscoresTable.ReadCapacityUnits == 0 {
		thingWithUnderscoresTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingWithUnderscoresTable.WriteCapacityUnits == 0 {
		thingWithUnderscoresTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}

	return &DB{
		simpleThingTable:                      simpleThingTable,
		teacherSharingRuleTable:               teacherSharingRuleTable,
		thingTable:                            thingTable,
		thingWithCompositeAttributesTable:     thingWithCompositeAttributesTable,
		thingWithCompositeEnumAttributesTable: thingWithCompositeEnumAttributesTable,
		thingWithDateRangeTable:               thingWithDateRangeTable,
		thingWithDateTimeCompositeTable:       thingWithDateTimeCompositeTable,
		thingWithRequiredFieldsTable:          thingWithRequiredFieldsTable,
		thingWithUnderscoresTable:             thingWithUnderscoresTable,
	}, nil
}

// DB implements the database interface using DynamoDB to store data.
type DB struct {
	simpleThingTable                      SimpleThingTable
	teacherSharingRuleTable               TeacherSharingRuleTable
	thingTable                            ThingTable
	thingWithCompositeAttributesTable     ThingWithCompositeAttributesTable
	thingWithCompositeEnumAttributesTable ThingWithCompositeEnumAttributesTable
	thingWithDateRangeTable               ThingWithDateRangeTable
	thingWithDateTimeCompositeTable       ThingWithDateTimeCompositeTable
	thingWithRequiredFieldsTable          ThingWithRequiredFieldsTable
	thingWithUnderscoresTable             ThingWithUnderscoresTable
}

var _ db.Interface = DB{}

// CreateTables creates all tables.
func (d DB) CreateTables(ctx context.Context) error {
	if err := d.simpleThingTable.create(ctx); err != nil {
		return err
	}
	if err := d.teacherSharingRuleTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithCompositeAttributesTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithCompositeEnumAttributesTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithDateRangeTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithDateTimeCompositeTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithRequiredFieldsTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithUnderscoresTable.create(ctx); err != nil {
		return err
	}
	return nil
}

// SaveSimpleThing saves a SimpleThing to the database.
func (d DB) SaveSimpleThing(ctx context.Context, m models.SimpleThing) error {
	return d.simpleThingTable.saveSimpleThing(ctx, m)
}

// GetSimpleThing retrieves a SimpleThing from the database.
func (d DB) GetSimpleThing(ctx context.Context, name string) (*models.SimpleThing, error) {
	return d.simpleThingTable.getSimpleThing(ctx, name)
}

// DeleteSimpleThing deletes a SimpleThing from the database.
func (d DB) DeleteSimpleThing(ctx context.Context, name string) error {
	return d.simpleThingTable.deleteSimpleThing(ctx, name)
}

// SaveTeacherSharingRule saves a TeacherSharingRule to the database.
func (d DB) SaveTeacherSharingRule(ctx context.Context, m models.TeacherSharingRule) error {
	return d.teacherSharingRuleTable.saveTeacherSharingRule(ctx, m)
}

// GetTeacherSharingRule retrieves a TeacherSharingRule from the database.
func (d DB) GetTeacherSharingRule(ctx context.Context, teacher string, school string, app string) (*models.TeacherSharingRule, error) {
	return d.teacherSharingRuleTable.getTeacherSharingRule(ctx, teacher, school, app)
}

// GetTeacherSharingRulesByTeacherAndSchoolApp retrieves a list of TeacherSharingRules from the database.
func (d DB) GetTeacherSharingRulesByTeacherAndSchoolApp(ctx context.Context, input db.GetTeacherSharingRulesByTeacherAndSchoolAppInput) ([]models.TeacherSharingRule, error) {
	return d.teacherSharingRuleTable.getTeacherSharingRulesByTeacherAndSchoolApp(ctx, input)
}

// GetTeacherSharingRulesByTeacherAndSchoolAppPage retrieves a page of TeacherSharingRules from the database.
func (d DB) GetTeacherSharingRulesByTeacherAndSchoolAppPage(ctx context.Context, input db.GetTeacherSharingRulesByTeacherAndSchoolAppPageInput, fn func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool) error {
	return d.teacherSharingRuleTable.getTeacherSharingRulesByTeacherAndSchoolAppPage(ctx, input, fn)
}

// DeleteTeacherSharingRule deletes a TeacherSharingRule from the database.
func (d DB) DeleteTeacherSharingRule(ctx context.Context, teacher string, school string, app string) error {
	return d.teacherSharingRuleTable.deleteTeacherSharingRule(ctx, teacher, school, app)
}

// GetTeacherSharingRulesByDistrictAndSchoolTeacherApp retrieves a list of TeacherSharingRules from the database.
func (d DB) GetTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx context.Context, input db.GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput) ([]models.TeacherSharingRule, error) {
	return d.teacherSharingRuleTable.getTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx, input)
}

// SaveThing saves a Thing to the database.
func (d DB) SaveThing(ctx context.Context, m models.Thing) error {
	return d.thingTable.saveThing(ctx, m)
}

// GetThing retrieves a Thing from the database.
func (d DB) GetThing(ctx context.Context, name string, version int64) (*models.Thing, error) {
	return d.thingTable.getThing(ctx, name, version)
}

// ScanThings runs a scan on the Things table.
func (d DB) ScanThings(ctx context.Context, input db.ScanThingsInput, fn func(m *models.Thing, lastThing bool) bool) error {
	return d.thingTable.scanThings(ctx, input, fn)
}

// GetThingsByNameAndVersion retrieves a list of Things from the database.
func (d DB) GetThingsByNameAndVersion(ctx context.Context, input db.GetThingsByNameAndVersionInput) ([]models.Thing, error) {
	return d.thingTable.getThingsByNameAndVersion(ctx, input)
}

// GetThingsByNameAndVersionPage retrieves a page of Things from the database.
func (d DB) GetThingsByNameAndVersionPage(ctx context.Context, input db.GetThingsByNameAndVersionPageInput, fn func(m *models.Thing, lastThing bool) bool) error {
	return d.thingTable.getThingsByNameAndVersionPage(ctx, input, fn)
}

// DeleteThing deletes a Thing from the database.
func (d DB) DeleteThing(ctx context.Context, name string, version int64) error {
	return d.thingTable.deleteThing(ctx, name, version)
}

// GetThingByID retrieves a Thing from the database.
func (d DB) GetThingByID(ctx context.Context, id string) (*models.Thing, error) {
	return d.thingTable.getThingByID(ctx, id)
}

// GetThingsByNameAndCreatedAt retrieves a list of Things from the database.
func (d DB) GetThingsByNameAndCreatedAt(ctx context.Context, input db.GetThingsByNameAndCreatedAtInput) ([]models.Thing, error) {
	return d.thingTable.getThingsByNameAndCreatedAt(ctx, input)
}

// SaveThingWithCompositeAttributes saves a ThingWithCompositeAttributes to the database.
func (d DB) SaveThingWithCompositeAttributes(ctx context.Context, m models.ThingWithCompositeAttributes) error {
	return d.thingWithCompositeAttributesTable.saveThingWithCompositeAttributes(ctx, m)
}

// GetThingWithCompositeAttributes retrieves a ThingWithCompositeAttributes from the database.
func (d DB) GetThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) (*models.ThingWithCompositeAttributes, error) {
	return d.thingWithCompositeAttributesTable.getThingWithCompositeAttributes(ctx, name, branch, date)
}

// GetThingWithCompositeAttributessByNameBranchAndDate retrieves a list of ThingWithCompositeAttributess from the database.
func (d DB) GetThingWithCompositeAttributessByNameBranchAndDate(ctx context.Context, input db.GetThingWithCompositeAttributessByNameBranchAndDateInput) ([]models.ThingWithCompositeAttributes, error) {
	return d.thingWithCompositeAttributesTable.getThingWithCompositeAttributessByNameBranchAndDate(ctx, input)
}

// GetThingWithCompositeAttributessByNameBranchAndDatePage retrieves a page of ThingWithCompositeAttributess from the database.
func (d DB) GetThingWithCompositeAttributessByNameBranchAndDatePage(ctx context.Context, input db.GetThingWithCompositeAttributessByNameBranchAndDatePageInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error {
	return d.thingWithCompositeAttributesTable.getThingWithCompositeAttributessByNameBranchAndDatePage(ctx, input, fn)
}

// DeleteThingWithCompositeAttributes deletes a ThingWithCompositeAttributes from the database.
func (d DB) DeleteThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) error {
	return d.thingWithCompositeAttributesTable.deleteThingWithCompositeAttributes(ctx, name, branch, date)
}

// GetThingWithCompositeAttributessByNameVersionAndDate retrieves a list of ThingWithCompositeAttributess from the database.
func (d DB) GetThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input db.GetThingWithCompositeAttributessByNameVersionAndDateInput) ([]models.ThingWithCompositeAttributes, error) {
	return d.thingWithCompositeAttributesTable.getThingWithCompositeAttributessByNameVersionAndDate(ctx, input)
}

// SaveThingWithCompositeEnumAttributes saves a ThingWithCompositeEnumAttributes to the database.
func (d DB) SaveThingWithCompositeEnumAttributes(ctx context.Context, m models.ThingWithCompositeEnumAttributes) error {
	return d.thingWithCompositeEnumAttributesTable.saveThingWithCompositeEnumAttributes(ctx, m)
}

// GetThingWithCompositeEnumAttributes retrieves a ThingWithCompositeEnumAttributes from the database.
func (d DB) GetThingWithCompositeEnumAttributes(ctx context.Context, name string, branchID models.Branch, date strfmt.DateTime) (*models.ThingWithCompositeEnumAttributes, error) {
	return d.thingWithCompositeEnumAttributesTable.getThingWithCompositeEnumAttributes(ctx, name, branchID, date)
}

// GetThingWithCompositeEnumAttributessByNameBranchAndDate retrieves a list of ThingWithCompositeEnumAttributess from the database.
func (d DB) GetThingWithCompositeEnumAttributessByNameBranchAndDate(ctx context.Context, input db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput) ([]models.ThingWithCompositeEnumAttributes, error) {
	return d.thingWithCompositeEnumAttributesTable.getThingWithCompositeEnumAttributessByNameBranchAndDate(ctx, input)
}

// GetThingWithCompositeEnumAttributessByNameBranchAndDatePage retrieves a page of ThingWithCompositeEnumAttributess from the database.
func (d DB) GetThingWithCompositeEnumAttributessByNameBranchAndDatePage(ctx context.Context, input db.GetThingWithCompositeEnumAttributessByNameBranchAndDatePageInput, fn func(m *models.ThingWithCompositeEnumAttributes, lastThingWithCompositeEnumAttributes bool) bool) error {
	return d.thingWithCompositeEnumAttributesTable.getThingWithCompositeEnumAttributessByNameBranchAndDatePage(ctx, input, fn)
}

// DeleteThingWithCompositeEnumAttributes deletes a ThingWithCompositeEnumAttributes from the database.
func (d DB) DeleteThingWithCompositeEnumAttributes(ctx context.Context, name string, branchID models.Branch, date strfmt.DateTime) error {
	return d.thingWithCompositeEnumAttributesTable.deleteThingWithCompositeEnumAttributes(ctx, name, branchID, date)
}

// SaveThingWithDateRange saves a ThingWithDateRange to the database.
func (d DB) SaveThingWithDateRange(ctx context.Context, m models.ThingWithDateRange) error {
	return d.thingWithDateRangeTable.saveThingWithDateRange(ctx, m)
}

// GetThingWithDateRange retrieves a ThingWithDateRange from the database.
func (d DB) GetThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) (*models.ThingWithDateRange, error) {
	return d.thingWithDateRangeTable.getThingWithDateRange(ctx, name, date)
}

// GetThingWithDateRangesByNameAndDate retrieves a list of ThingWithDateRanges from the database.
func (d DB) GetThingWithDateRangesByNameAndDate(ctx context.Context, input db.GetThingWithDateRangesByNameAndDateInput) ([]models.ThingWithDateRange, error) {
	return d.thingWithDateRangeTable.getThingWithDateRangesByNameAndDate(ctx, input)
}

// GetThingWithDateRangesByNameAndDatePage retrieves a page of ThingWithDateRanges from the database.
func (d DB) GetThingWithDateRangesByNameAndDatePage(ctx context.Context, input db.GetThingWithDateRangesByNameAndDatePageInput, fn func(m *models.ThingWithDateRange, lastThingWithDateRange bool) bool) error {
	return d.thingWithDateRangeTable.getThingWithDateRangesByNameAndDatePage(ctx, input, fn)
}

// DeleteThingWithDateRange deletes a ThingWithDateRange from the database.
func (d DB) DeleteThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) error {
	return d.thingWithDateRangeTable.deleteThingWithDateRange(ctx, name, date)
}

// SaveThingWithDateTimeComposite saves a ThingWithDateTimeComposite to the database.
func (d DB) SaveThingWithDateTimeComposite(ctx context.Context, m models.ThingWithDateTimeComposite) error {
	return d.thingWithDateTimeCompositeTable.saveThingWithDateTimeComposite(ctx, m)
}

// GetThingWithDateTimeComposite retrieves a ThingWithDateTimeComposite from the database.
func (d DB) GetThingWithDateTimeComposite(ctx context.Context, typeVar string, id string, created strfmt.DateTime, resource string) (*models.ThingWithDateTimeComposite, error) {
	return d.thingWithDateTimeCompositeTable.getThingWithDateTimeComposite(ctx, typeVar, id, created, resource)
}

// GetThingWithDateTimeCompositesByTypeIDAndCreatedResource retrieves a list of ThingWithDateTimeComposites from the database.
func (d DB) GetThingWithDateTimeCompositesByTypeIDAndCreatedResource(ctx context.Context, input db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput) ([]models.ThingWithDateTimeComposite, error) {
	return d.thingWithDateTimeCompositeTable.getThingWithDateTimeCompositesByTypeIDAndCreatedResource(ctx, input)
}

// GetThingWithDateTimeCompositesByTypeIDAndCreatedResourcePage retrieves a page of ThingWithDateTimeComposites from the database.
func (d DB) GetThingWithDateTimeCompositesByTypeIDAndCreatedResourcePage(ctx context.Context, input db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourcePageInput, fn func(m *models.ThingWithDateTimeComposite, lastThingWithDateTimeComposite bool) bool) error {
	return d.thingWithDateTimeCompositeTable.getThingWithDateTimeCompositesByTypeIDAndCreatedResourcePage(ctx, input, fn)
}

// DeleteThingWithDateTimeComposite deletes a ThingWithDateTimeComposite from the database.
func (d DB) DeleteThingWithDateTimeComposite(ctx context.Context, typeVar string, id string, created strfmt.DateTime, resource string) error {
	return d.thingWithDateTimeCompositeTable.deleteThingWithDateTimeComposite(ctx, typeVar, id, created, resource)
}

// SaveThingWithRequiredFields saves a ThingWithRequiredFields to the database.
func (d DB) SaveThingWithRequiredFields(ctx context.Context, m models.ThingWithRequiredFields) error {
	return d.thingWithRequiredFieldsTable.saveThingWithRequiredFields(ctx, m)
}

// GetThingWithRequiredFields retrieves a ThingWithRequiredFields from the database.
func (d DB) GetThingWithRequiredFields(ctx context.Context, name string) (*models.ThingWithRequiredFields, error) {
	return d.thingWithRequiredFieldsTable.getThingWithRequiredFields(ctx, name)
}

// DeleteThingWithRequiredFields deletes a ThingWithRequiredFields from the database.
func (d DB) DeleteThingWithRequiredFields(ctx context.Context, name string) error {
	return d.thingWithRequiredFieldsTable.deleteThingWithRequiredFields(ctx, name)
}

// SaveThingWithUnderscores saves a ThingWithUnderscores to the database.
func (d DB) SaveThingWithUnderscores(ctx context.Context, m models.ThingWithUnderscores) error {
	return d.thingWithUnderscoresTable.saveThingWithUnderscores(ctx, m)
}

// GetThingWithUnderscores retrieves a ThingWithUnderscores from the database.
func (d DB) GetThingWithUnderscores(ctx context.Context, iDApp string) (*models.ThingWithUnderscores, error) {
	return d.thingWithUnderscoresTable.getThingWithUnderscores(ctx, iDApp)
}

// DeleteThingWithUnderscores deletes a ThingWithUnderscores from the database.
func (d DB) DeleteThingWithUnderscores(ctx context.Context, iDApp string) error {
	return d.thingWithUnderscoresTable.deleteThingWithUnderscores(ctx, iDApp)
}

// bad hack for type checking
type strfmtTime interface {
	Value() (driver.Value, error)
	MarshalEasyJSON(w *jwriter.Writer)
}

func toDynamoTimeString(d strfmtTime) string {
	switch v := d.(type) {
	case strfmt.DateTime:
		return time.Time(v).Format(time.RFC3339) // dynamodb attributevalue only supports RFC3339 resolution
	case *strfmt.DateTime:
		return time.Time(*v).Format(time.RFC3339) // dynamodb attributevalue only supports RFC3339 resolution
	default:
		log.Fatal("oops")
		return ""
	}
}
