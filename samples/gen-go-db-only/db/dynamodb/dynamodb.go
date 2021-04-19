package dynamodb

import (
	"context"
	"errors"
	"time"

	"github.com/Clever/wag/v7/samples/gen-go-db-only/db"
	"github.com/Clever/wag/v7/samples/gen-go-db-only/models"
	ddb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
	otaws "github.com/opentracing-contrib/go-aws-sdk"
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
	// DeploymentTable configuration.
	DeploymentTable DeploymentTable
	// EventTable configuration.
	EventTable EventTable
	// NoRangeThingWithCompositeAttributesTable configuration.
	NoRangeThingWithCompositeAttributesTable NoRangeThingWithCompositeAttributesTable
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
	// ThingWithEnumHashKeyTable configuration.
	ThingWithEnumHashKeyTable ThingWithEnumHashKeyTable
	// ThingWithMatchingKeysTable configuration.
	ThingWithMatchingKeysTable ThingWithMatchingKeysTable
	// ThingWithMultiUseCompositeAttributeTable configuration.
	ThingWithMultiUseCompositeAttributeTable ThingWithMultiUseCompositeAttributeTable
	// ThingWithRequiredCompositePropertiesAndKeysOnlyTable configuration.
	ThingWithRequiredCompositePropertiesAndKeysOnlyTable ThingWithRequiredCompositePropertiesAndKeysOnlyTable
	// ThingWithRequiredFieldsTable configuration.
	ThingWithRequiredFieldsTable ThingWithRequiredFieldsTable
	// ThingWithRequiredFields2Table configuration.
	ThingWithRequiredFields2Table ThingWithRequiredFields2Table
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

	// add OpenTracing observability to DDB calls if it's a non-mocked client
	if ddbClient, ok := config.DynamoDBAPI.(*ddb.DynamoDB); ok {
		otaws.AddOTHandlers(ddbClient.Client)
	}

	if config.DefaultWriteCapacityUnits == 0 {
		config.DefaultWriteCapacityUnits = 1
	}
	if config.DefaultReadCapacityUnits == 0 {
		config.DefaultReadCapacityUnits = 1
	}
	// configure Deployment table
	deploymentTable := config.DeploymentTable
	if deploymentTable.DynamoDBAPI == nil {
		deploymentTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if deploymentTable.Prefix == "" {
		deploymentTable.Prefix = config.DefaultPrefix
	}
	if deploymentTable.ReadCapacityUnits == 0 {
		deploymentTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if deploymentTable.WriteCapacityUnits == 0 {
		deploymentTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure Event table
	eventTable := config.EventTable
	if eventTable.DynamoDBAPI == nil {
		eventTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if eventTable.Prefix == "" {
		eventTable.Prefix = config.DefaultPrefix
	}
	if eventTable.ReadCapacityUnits == 0 {
		eventTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if eventTable.WriteCapacityUnits == 0 {
		eventTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure NoRangeThingWithCompositeAttributes table
	noRangeThingWithCompositeAttributesTable := config.NoRangeThingWithCompositeAttributesTable
	if noRangeThingWithCompositeAttributesTable.DynamoDBAPI == nil {
		noRangeThingWithCompositeAttributesTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if noRangeThingWithCompositeAttributesTable.Prefix == "" {
		noRangeThingWithCompositeAttributesTable.Prefix = config.DefaultPrefix
	}
	if noRangeThingWithCompositeAttributesTable.ReadCapacityUnits == 0 {
		noRangeThingWithCompositeAttributesTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if noRangeThingWithCompositeAttributesTable.WriteCapacityUnits == 0 {
		noRangeThingWithCompositeAttributesTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
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
	// configure ThingWithEnumHashKey table
	thingWithEnumHashKeyTable := config.ThingWithEnumHashKeyTable
	if thingWithEnumHashKeyTable.DynamoDBAPI == nil {
		thingWithEnumHashKeyTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingWithEnumHashKeyTable.Prefix == "" {
		thingWithEnumHashKeyTable.Prefix = config.DefaultPrefix
	}
	if thingWithEnumHashKeyTable.ReadCapacityUnits == 0 {
		thingWithEnumHashKeyTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingWithEnumHashKeyTable.WriteCapacityUnits == 0 {
		thingWithEnumHashKeyTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure ThingWithMatchingKeys table
	thingWithMatchingKeysTable := config.ThingWithMatchingKeysTable
	if thingWithMatchingKeysTable.DynamoDBAPI == nil {
		thingWithMatchingKeysTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingWithMatchingKeysTable.Prefix == "" {
		thingWithMatchingKeysTable.Prefix = config.DefaultPrefix
	}
	if thingWithMatchingKeysTable.ReadCapacityUnits == 0 {
		thingWithMatchingKeysTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingWithMatchingKeysTable.WriteCapacityUnits == 0 {
		thingWithMatchingKeysTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure ThingWithMultiUseCompositeAttribute table
	thingWithMultiUseCompositeAttributeTable := config.ThingWithMultiUseCompositeAttributeTable
	if thingWithMultiUseCompositeAttributeTable.DynamoDBAPI == nil {
		thingWithMultiUseCompositeAttributeTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingWithMultiUseCompositeAttributeTable.Prefix == "" {
		thingWithMultiUseCompositeAttributeTable.Prefix = config.DefaultPrefix
	}
	if thingWithMultiUseCompositeAttributeTable.ReadCapacityUnits == 0 {
		thingWithMultiUseCompositeAttributeTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingWithMultiUseCompositeAttributeTable.WriteCapacityUnits == 0 {
		thingWithMultiUseCompositeAttributeTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}
	// configure ThingWithRequiredCompositePropertiesAndKeysOnly table
	thingWithRequiredCompositePropertiesAndKeysOnlyTable := config.ThingWithRequiredCompositePropertiesAndKeysOnlyTable
	if thingWithRequiredCompositePropertiesAndKeysOnlyTable.DynamoDBAPI == nil {
		thingWithRequiredCompositePropertiesAndKeysOnlyTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingWithRequiredCompositePropertiesAndKeysOnlyTable.Prefix == "" {
		thingWithRequiredCompositePropertiesAndKeysOnlyTable.Prefix = config.DefaultPrefix
	}
	if thingWithRequiredCompositePropertiesAndKeysOnlyTable.ReadCapacityUnits == 0 {
		thingWithRequiredCompositePropertiesAndKeysOnlyTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingWithRequiredCompositePropertiesAndKeysOnlyTable.WriteCapacityUnits == 0 {
		thingWithRequiredCompositePropertiesAndKeysOnlyTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
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
	// configure ThingWithRequiredFields2 table
	thingWithRequiredFields2Table := config.ThingWithRequiredFields2Table
	if thingWithRequiredFields2Table.DynamoDBAPI == nil {
		thingWithRequiredFields2Table.DynamoDBAPI = config.DynamoDBAPI
	}
	if thingWithRequiredFields2Table.Prefix == "" {
		thingWithRequiredFields2Table.Prefix = config.DefaultPrefix
	}
	if thingWithRequiredFields2Table.ReadCapacityUnits == 0 {
		thingWithRequiredFields2Table.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if thingWithRequiredFields2Table.WriteCapacityUnits == 0 {
		thingWithRequiredFields2Table.WriteCapacityUnits = config.DefaultWriteCapacityUnits
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
		deploymentTable:                                      deploymentTable,
		eventTable:                                           eventTable,
		noRangeThingWithCompositeAttributesTable:             noRangeThingWithCompositeAttributesTable,
		simpleThingTable:                                     simpleThingTable,
		teacherSharingRuleTable:                              teacherSharingRuleTable,
		thingTable:                                           thingTable,
		thingWithCompositeAttributesTable:                    thingWithCompositeAttributesTable,
		thingWithCompositeEnumAttributesTable:                thingWithCompositeEnumAttributesTable,
		thingWithDateRangeTable:                              thingWithDateRangeTable,
		thingWithDateTimeCompositeTable:                      thingWithDateTimeCompositeTable,
		thingWithEnumHashKeyTable:                            thingWithEnumHashKeyTable,
		thingWithMatchingKeysTable:                           thingWithMatchingKeysTable,
		thingWithMultiUseCompositeAttributeTable:             thingWithMultiUseCompositeAttributeTable,
		thingWithRequiredCompositePropertiesAndKeysOnlyTable: thingWithRequiredCompositePropertiesAndKeysOnlyTable,
		thingWithRequiredFieldsTable:                         thingWithRequiredFieldsTable,
		thingWithRequiredFields2Table:                        thingWithRequiredFields2Table,
		thingWithUnderscoresTable:                            thingWithUnderscoresTable,
	}, nil
}

// DB implements the database interface using DynamoDB to store data.
type DB struct {
	deploymentTable                                      DeploymentTable
	eventTable                                           EventTable
	noRangeThingWithCompositeAttributesTable             NoRangeThingWithCompositeAttributesTable
	simpleThingTable                                     SimpleThingTable
	teacherSharingRuleTable                              TeacherSharingRuleTable
	thingTable                                           ThingTable
	thingWithCompositeAttributesTable                    ThingWithCompositeAttributesTable
	thingWithCompositeEnumAttributesTable                ThingWithCompositeEnumAttributesTable
	thingWithDateRangeTable                              ThingWithDateRangeTable
	thingWithDateTimeCompositeTable                      ThingWithDateTimeCompositeTable
	thingWithEnumHashKeyTable                            ThingWithEnumHashKeyTable
	thingWithMatchingKeysTable                           ThingWithMatchingKeysTable
	thingWithMultiUseCompositeAttributeTable             ThingWithMultiUseCompositeAttributeTable
	thingWithRequiredCompositePropertiesAndKeysOnlyTable ThingWithRequiredCompositePropertiesAndKeysOnlyTable
	thingWithRequiredFieldsTable                         ThingWithRequiredFieldsTable
	thingWithRequiredFields2Table                        ThingWithRequiredFields2Table
	thingWithUnderscoresTable                            ThingWithUnderscoresTable
}

var _ db.Interface = DB{}

// CreateTables creates all tables.
func (d DB) CreateTables(ctx context.Context) error {
	if err := d.deploymentTable.create(ctx); err != nil {
		return err
	}
	if err := d.eventTable.create(ctx); err != nil {
		return err
	}
	if err := d.noRangeThingWithCompositeAttributesTable.create(ctx); err != nil {
		return err
	}
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
	if err := d.thingWithEnumHashKeyTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithMatchingKeysTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithMultiUseCompositeAttributeTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithRequiredCompositePropertiesAndKeysOnlyTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithRequiredFieldsTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithRequiredFields2Table.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithUnderscoresTable.create(ctx); err != nil {
		return err
	}
	return nil
}

// SaveDeployment saves a Deployment to the database.
func (d DB) SaveDeployment(ctx context.Context, m models.Deployment) error {
	return d.deploymentTable.saveDeployment(ctx, m)
}

// GetDeployment retrieves a Deployment from the database.
func (d DB) GetDeployment(ctx context.Context, environment string, application string, version string) (*models.Deployment, error) {
	return d.deploymentTable.getDeployment(ctx, environment, application, version)
}

// ScanDeployments runs a scan on the Deployments table.
func (d DB) ScanDeployments(ctx context.Context, input db.ScanDeploymentsInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
	return d.deploymentTable.scanDeployments(ctx, input, fn)
}

// GetDeploymentsByEnvAppAndVersion retrieves a page of Deployments from the database.
func (d DB) GetDeploymentsByEnvAppAndVersion(ctx context.Context, input db.GetDeploymentsByEnvAppAndVersionInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
	return d.deploymentTable.getDeploymentsByEnvAppAndVersion(ctx, input, fn)
}

// DeleteDeployment deletes a Deployment from the database.
func (d DB) DeleteDeployment(ctx context.Context, environment string, application string, version string) error {
	return d.deploymentTable.deleteDeployment(ctx, environment, application, version)
}

// GetDeploymentsByEnvAppAndDate retrieves a page of Deployments from the database.
func (d DB) GetDeploymentsByEnvAppAndDate(ctx context.Context, input db.GetDeploymentsByEnvAppAndDateInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
	return d.deploymentTable.getDeploymentsByEnvAppAndDate(ctx, input, fn)
}

// ScanDeploymentsByEnvAppAndDate runs a scan on the EnvAppAndDate index.
func (d DB) ScanDeploymentsByEnvAppAndDate(ctx context.Context, input db.ScanDeploymentsByEnvAppAndDateInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
	return d.deploymentTable.scanDeploymentsByEnvAppAndDate(ctx, input, fn)
}

// GetDeploymentsByEnvironmentAndDate retrieves a page of Deployments from the database.
func (d DB) GetDeploymentsByEnvironmentAndDate(ctx context.Context, input db.GetDeploymentsByEnvironmentAndDateInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
	return d.deploymentTable.getDeploymentsByEnvironmentAndDate(ctx, input, fn)
}

// GetDeploymentByVersion retrieves a Deployment from the database.
func (d DB) GetDeploymentByVersion(ctx context.Context, version string) (*models.Deployment, error) {
	return d.deploymentTable.getDeploymentByVersion(ctx, version)
}

// ScanDeploymentsByVersion runs a scan on the Version index.
func (d DB) ScanDeploymentsByVersion(ctx context.Context, input db.ScanDeploymentsByVersionInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
	return d.deploymentTable.scanDeploymentsByVersion(ctx, input, fn)
}

// SaveEvent saves a Event to the database.
func (d DB) SaveEvent(ctx context.Context, m models.Event) error {
	return d.eventTable.saveEvent(ctx, m)
}

// GetEvent retrieves a Event from the database.
func (d DB) GetEvent(ctx context.Context, pk string, sk string) (*models.Event, error) {
	return d.eventTable.getEvent(ctx, pk, sk)
}

// ScanEvents runs a scan on the Events table.
func (d DB) ScanEvents(ctx context.Context, input db.ScanEventsInput, fn func(m *models.Event, lastEvent bool) bool) error {
	return d.eventTable.scanEvents(ctx, input, fn)
}

// GetEventsByPkAndSk retrieves a page of Events from the database.
func (d DB) GetEventsByPkAndSk(ctx context.Context, input db.GetEventsByPkAndSkInput, fn func(m *models.Event, lastEvent bool) bool) error {
	return d.eventTable.getEventsByPkAndSk(ctx, input, fn)
}

// DeleteEvent deletes a Event from the database.
func (d DB) DeleteEvent(ctx context.Context, pk string, sk string) error {
	return d.eventTable.deleteEvent(ctx, pk, sk)
}

// GetEventsBySkAndData retrieves a page of Events from the database.
func (d DB) GetEventsBySkAndData(ctx context.Context, input db.GetEventsBySkAndDataInput, fn func(m *models.Event, lastEvent bool) bool) error {
	return d.eventTable.getEventsBySkAndData(ctx, input, fn)
}

// ScanEventsBySkAndData runs a scan on the SkAndData index.
func (d DB) ScanEventsBySkAndData(ctx context.Context, input db.ScanEventsBySkAndDataInput, fn func(m *models.Event, lastEvent bool) bool) error {
	return d.eventTable.scanEventsBySkAndData(ctx, input, fn)
}

// SaveNoRangeThingWithCompositeAttributes saves a NoRangeThingWithCompositeAttributes to the database.
func (d DB) SaveNoRangeThingWithCompositeAttributes(ctx context.Context, m models.NoRangeThingWithCompositeAttributes) error {
	return d.noRangeThingWithCompositeAttributesTable.saveNoRangeThingWithCompositeAttributes(ctx, m)
}

// GetNoRangeThingWithCompositeAttributes retrieves a NoRangeThingWithCompositeAttributes from the database.
func (d DB) GetNoRangeThingWithCompositeAttributes(ctx context.Context, name string, branch string) (*models.NoRangeThingWithCompositeAttributes, error) {
	return d.noRangeThingWithCompositeAttributesTable.getNoRangeThingWithCompositeAttributes(ctx, name, branch)
}

// ScanNoRangeThingWithCompositeAttributess runs a scan on the NoRangeThingWithCompositeAttributess table.
func (d DB) ScanNoRangeThingWithCompositeAttributess(ctx context.Context, input db.ScanNoRangeThingWithCompositeAttributessInput, fn func(m *models.NoRangeThingWithCompositeAttributes, lastNoRangeThingWithCompositeAttributes bool) bool) error {
	return d.noRangeThingWithCompositeAttributesTable.scanNoRangeThingWithCompositeAttributess(ctx, input, fn)
}

// DeleteNoRangeThingWithCompositeAttributes deletes a NoRangeThingWithCompositeAttributes from the database.
func (d DB) DeleteNoRangeThingWithCompositeAttributes(ctx context.Context, name string, branch string) error {
	return d.noRangeThingWithCompositeAttributesTable.deleteNoRangeThingWithCompositeAttributes(ctx, name, branch)
}

// GetNoRangeThingWithCompositeAttributessByNameVersionAndDate retrieves a page of NoRangeThingWithCompositeAttributess from the database.
func (d DB) GetNoRangeThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input db.GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.NoRangeThingWithCompositeAttributes, lastNoRangeThingWithCompositeAttributes bool) bool) error {
	return d.noRangeThingWithCompositeAttributesTable.getNoRangeThingWithCompositeAttributessByNameVersionAndDate(ctx, input, fn)
}

// ScanNoRangeThingWithCompositeAttributessByNameVersionAndDate runs a scan on the NameVersionAndDate index.
func (d DB) ScanNoRangeThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input db.ScanNoRangeThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.NoRangeThingWithCompositeAttributes, lastNoRangeThingWithCompositeAttributes bool) bool) error {
	return d.noRangeThingWithCompositeAttributesTable.scanNoRangeThingWithCompositeAttributessByNameVersionAndDate(ctx, input, fn)
}

// SaveSimpleThing saves a SimpleThing to the database.
func (d DB) SaveSimpleThing(ctx context.Context, m models.SimpleThing) error {
	return d.simpleThingTable.saveSimpleThing(ctx, m)
}

// GetSimpleThing retrieves a SimpleThing from the database.
func (d DB) GetSimpleThing(ctx context.Context, name string) (*models.SimpleThing, error) {
	return d.simpleThingTable.getSimpleThing(ctx, name)
}

// ScanSimpleThings runs a scan on the SimpleThings table.
func (d DB) ScanSimpleThings(ctx context.Context, input db.ScanSimpleThingsInput, fn func(m *models.SimpleThing, lastSimpleThing bool) bool) error {
	return d.simpleThingTable.scanSimpleThings(ctx, input, fn)
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

// ScanTeacherSharingRules runs a scan on the TeacherSharingRules table.
func (d DB) ScanTeacherSharingRules(ctx context.Context, input db.ScanTeacherSharingRulesInput, fn func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool) error {
	return d.teacherSharingRuleTable.scanTeacherSharingRules(ctx, input, fn)
}

// GetTeacherSharingRulesByTeacherAndSchoolApp retrieves a page of TeacherSharingRules from the database.
func (d DB) GetTeacherSharingRulesByTeacherAndSchoolApp(ctx context.Context, input db.GetTeacherSharingRulesByTeacherAndSchoolAppInput, fn func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool) error {
	return d.teacherSharingRuleTable.getTeacherSharingRulesByTeacherAndSchoolApp(ctx, input, fn)
}

// DeleteTeacherSharingRule deletes a TeacherSharingRule from the database.
func (d DB) DeleteTeacherSharingRule(ctx context.Context, teacher string, school string, app string) error {
	return d.teacherSharingRuleTable.deleteTeacherSharingRule(ctx, teacher, school, app)
}

// GetTeacherSharingRulesByDistrictAndSchoolTeacherApp retrieves a page of TeacherSharingRules from the database.
func (d DB) GetTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx context.Context, input db.GetTeacherSharingRulesByDistrictAndSchoolTeacherAppInput, fn func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool) error {
	return d.teacherSharingRuleTable.getTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx, input, fn)
}

// ScanTeacherSharingRulesByDistrictAndSchoolTeacherApp runs a scan on the DistrictAndSchoolTeacherApp index.
func (d DB) ScanTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx context.Context, input db.ScanTeacherSharingRulesByDistrictAndSchoolTeacherAppInput, fn func(m *models.TeacherSharingRule, lastTeacherSharingRule bool) bool) error {
	return d.teacherSharingRuleTable.scanTeacherSharingRulesByDistrictAndSchoolTeacherApp(ctx, input, fn)
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

// GetThingsByNameAndVersion retrieves a page of Things from the database.
func (d DB) GetThingsByNameAndVersion(ctx context.Context, input db.GetThingsByNameAndVersionInput, fn func(m *models.Thing, lastThing bool) bool) error {
	return d.thingTable.getThingsByNameAndVersion(ctx, input, fn)
}

// DeleteThing deletes a Thing from the database.
func (d DB) DeleteThing(ctx context.Context, name string, version int64) error {
	return d.thingTable.deleteThing(ctx, name, version)
}

// GetThingByID retrieves a Thing from the database.
func (d DB) GetThingByID(ctx context.Context, id string) (*models.Thing, error) {
	return d.thingTable.getThingByID(ctx, id)
}

// ScanThingsByID runs a scan on the ID index.
func (d DB) ScanThingsByID(ctx context.Context, input db.ScanThingsByIDInput, fn func(m *models.Thing, lastThing bool) bool) error {
	return d.thingTable.scanThingsByID(ctx, input, fn)
}

// GetThingsByNameAndCreatedAt retrieves a page of Things from the database.
func (d DB) GetThingsByNameAndCreatedAt(ctx context.Context, input db.GetThingsByNameAndCreatedAtInput, fn func(m *models.Thing, lastThing bool) bool) error {
	return d.thingTable.getThingsByNameAndCreatedAt(ctx, input, fn)
}

// ScanThingsByNameAndCreatedAt runs a scan on the NameAndCreatedAt index.
func (d DB) ScanThingsByNameAndCreatedAt(ctx context.Context, input db.ScanThingsByNameAndCreatedAtInput, fn func(m *models.Thing, lastThing bool) bool) error {
	return d.thingTable.scanThingsByNameAndCreatedAt(ctx, input, fn)
}

// SaveThingWithCompositeAttributes saves a ThingWithCompositeAttributes to the database.
func (d DB) SaveThingWithCompositeAttributes(ctx context.Context, m models.ThingWithCompositeAttributes) error {
	return d.thingWithCompositeAttributesTable.saveThingWithCompositeAttributes(ctx, m)
}

// GetThingWithCompositeAttributes retrieves a ThingWithCompositeAttributes from the database.
func (d DB) GetThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) (*models.ThingWithCompositeAttributes, error) {
	return d.thingWithCompositeAttributesTable.getThingWithCompositeAttributes(ctx, name, branch, date)
}

// ScanThingWithCompositeAttributess runs a scan on the ThingWithCompositeAttributess table.
func (d DB) ScanThingWithCompositeAttributess(ctx context.Context, input db.ScanThingWithCompositeAttributessInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error {
	return d.thingWithCompositeAttributesTable.scanThingWithCompositeAttributess(ctx, input, fn)
}

// GetThingWithCompositeAttributessByNameBranchAndDate retrieves a page of ThingWithCompositeAttributess from the database.
func (d DB) GetThingWithCompositeAttributessByNameBranchAndDate(ctx context.Context, input db.GetThingWithCompositeAttributessByNameBranchAndDateInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error {
	return d.thingWithCompositeAttributesTable.getThingWithCompositeAttributessByNameBranchAndDate(ctx, input, fn)
}

// DeleteThingWithCompositeAttributes deletes a ThingWithCompositeAttributes from the database.
func (d DB) DeleteThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) error {
	return d.thingWithCompositeAttributesTable.deleteThingWithCompositeAttributes(ctx, name, branch, date)
}

// GetThingWithCompositeAttributessByNameVersionAndDate retrieves a page of ThingWithCompositeAttributess from the database.
func (d DB) GetThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input db.GetThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error {
	return d.thingWithCompositeAttributesTable.getThingWithCompositeAttributessByNameVersionAndDate(ctx, input, fn)
}

// ScanThingWithCompositeAttributessByNameVersionAndDate runs a scan on the NameVersionAndDate index.
func (d DB) ScanThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input db.ScanThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error {
	return d.thingWithCompositeAttributesTable.scanThingWithCompositeAttributessByNameVersionAndDate(ctx, input, fn)
}

// SaveThingWithCompositeEnumAttributes saves a ThingWithCompositeEnumAttributes to the database.
func (d DB) SaveThingWithCompositeEnumAttributes(ctx context.Context, m models.ThingWithCompositeEnumAttributes) error {
	return d.thingWithCompositeEnumAttributesTable.saveThingWithCompositeEnumAttributes(ctx, m)
}

// GetThingWithCompositeEnumAttributes retrieves a ThingWithCompositeEnumAttributes from the database.
func (d DB) GetThingWithCompositeEnumAttributes(ctx context.Context, name string, branchID models.Branch, date strfmt.DateTime) (*models.ThingWithCompositeEnumAttributes, error) {
	return d.thingWithCompositeEnumAttributesTable.getThingWithCompositeEnumAttributes(ctx, name, branchID, date)
}

// ScanThingWithCompositeEnumAttributess runs a scan on the ThingWithCompositeEnumAttributess table.
func (d DB) ScanThingWithCompositeEnumAttributess(ctx context.Context, input db.ScanThingWithCompositeEnumAttributessInput, fn func(m *models.ThingWithCompositeEnumAttributes, lastThingWithCompositeEnumAttributes bool) bool) error {
	return d.thingWithCompositeEnumAttributesTable.scanThingWithCompositeEnumAttributess(ctx, input, fn)
}

// GetThingWithCompositeEnumAttributessByNameBranchAndDate retrieves a page of ThingWithCompositeEnumAttributess from the database.
func (d DB) GetThingWithCompositeEnumAttributessByNameBranchAndDate(ctx context.Context, input db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput, fn func(m *models.ThingWithCompositeEnumAttributes, lastThingWithCompositeEnumAttributes bool) bool) error {
	return d.thingWithCompositeEnumAttributesTable.getThingWithCompositeEnumAttributessByNameBranchAndDate(ctx, input, fn)
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

// ScanThingWithDateRanges runs a scan on the ThingWithDateRanges table.
func (d DB) ScanThingWithDateRanges(ctx context.Context, input db.ScanThingWithDateRangesInput, fn func(m *models.ThingWithDateRange, lastThingWithDateRange bool) bool) error {
	return d.thingWithDateRangeTable.scanThingWithDateRanges(ctx, input, fn)
}

// GetThingWithDateRangesByNameAndDate retrieves a page of ThingWithDateRanges from the database.
func (d DB) GetThingWithDateRangesByNameAndDate(ctx context.Context, input db.GetThingWithDateRangesByNameAndDateInput, fn func(m *models.ThingWithDateRange, lastThingWithDateRange bool) bool) error {
	return d.thingWithDateRangeTable.getThingWithDateRangesByNameAndDate(ctx, input, fn)
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

// ScanThingWithDateTimeComposites runs a scan on the ThingWithDateTimeComposites table.
func (d DB) ScanThingWithDateTimeComposites(ctx context.Context, input db.ScanThingWithDateTimeCompositesInput, fn func(m *models.ThingWithDateTimeComposite, lastThingWithDateTimeComposite bool) bool) error {
	return d.thingWithDateTimeCompositeTable.scanThingWithDateTimeComposites(ctx, input, fn)
}

// GetThingWithDateTimeCompositesByTypeIDAndCreatedResource retrieves a page of ThingWithDateTimeComposites from the database.
func (d DB) GetThingWithDateTimeCompositesByTypeIDAndCreatedResource(ctx context.Context, input db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput, fn func(m *models.ThingWithDateTimeComposite, lastThingWithDateTimeComposite bool) bool) error {
	return d.thingWithDateTimeCompositeTable.getThingWithDateTimeCompositesByTypeIDAndCreatedResource(ctx, input, fn)
}

// DeleteThingWithDateTimeComposite deletes a ThingWithDateTimeComposite from the database.
func (d DB) DeleteThingWithDateTimeComposite(ctx context.Context, typeVar string, id string, created strfmt.DateTime, resource string) error {
	return d.thingWithDateTimeCompositeTable.deleteThingWithDateTimeComposite(ctx, typeVar, id, created, resource)
}

// SaveThingWithEnumHashKey saves a ThingWithEnumHashKey to the database.
func (d DB) SaveThingWithEnumHashKey(ctx context.Context, m models.ThingWithEnumHashKey) error {
	return d.thingWithEnumHashKeyTable.saveThingWithEnumHashKey(ctx, m)
}

// GetThingWithEnumHashKey retrieves a ThingWithEnumHashKey from the database.
func (d DB) GetThingWithEnumHashKey(ctx context.Context, branch models.Branch, date strfmt.DateTime) (*models.ThingWithEnumHashKey, error) {
	return d.thingWithEnumHashKeyTable.getThingWithEnumHashKey(ctx, branch, date)
}

// ScanThingWithEnumHashKeys runs a scan on the ThingWithEnumHashKeys table.
func (d DB) ScanThingWithEnumHashKeys(ctx context.Context, input db.ScanThingWithEnumHashKeysInput, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error {
	return d.thingWithEnumHashKeyTable.scanThingWithEnumHashKeys(ctx, input, fn)
}

// GetThingWithEnumHashKeysByBranchAndDate retrieves a page of ThingWithEnumHashKeys from the database.
func (d DB) GetThingWithEnumHashKeysByBranchAndDate(ctx context.Context, input db.GetThingWithEnumHashKeysByBranchAndDateInput, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error {
	return d.thingWithEnumHashKeyTable.getThingWithEnumHashKeysByBranchAndDate(ctx, input, fn)
}

// DeleteThingWithEnumHashKey deletes a ThingWithEnumHashKey from the database.
func (d DB) DeleteThingWithEnumHashKey(ctx context.Context, branch models.Branch, date strfmt.DateTime) error {
	return d.thingWithEnumHashKeyTable.deleteThingWithEnumHashKey(ctx, branch, date)
}

// GetThingWithEnumHashKeysByBranchAndDate2 retrieves a page of ThingWithEnumHashKeys from the database.
func (d DB) GetThingWithEnumHashKeysByBranchAndDate2(ctx context.Context, input db.GetThingWithEnumHashKeysByBranchAndDate2Input, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error {
	return d.thingWithEnumHashKeyTable.getThingWithEnumHashKeysByBranchAndDate2(ctx, input, fn)
}

// ScanThingWithEnumHashKeysByBranchAndDate2 runs a scan on the BranchAndDate2 index.
func (d DB) ScanThingWithEnumHashKeysByBranchAndDate2(ctx context.Context, input db.ScanThingWithEnumHashKeysByBranchAndDate2Input, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error {
	return d.thingWithEnumHashKeyTable.scanThingWithEnumHashKeysByBranchAndDate2(ctx, input, fn)
}

// SaveThingWithMatchingKeys saves a ThingWithMatchingKeys to the database.
func (d DB) SaveThingWithMatchingKeys(ctx context.Context, m models.ThingWithMatchingKeys) error {
	return d.thingWithMatchingKeysTable.saveThingWithMatchingKeys(ctx, m)
}

// GetThingWithMatchingKeys retrieves a ThingWithMatchingKeys from the database.
func (d DB) GetThingWithMatchingKeys(ctx context.Context, bear string, assocType string, assocID string) (*models.ThingWithMatchingKeys, error) {
	return d.thingWithMatchingKeysTable.getThingWithMatchingKeys(ctx, bear, assocType, assocID)
}

// ScanThingWithMatchingKeyss runs a scan on the ThingWithMatchingKeyss table.
func (d DB) ScanThingWithMatchingKeyss(ctx context.Context, input db.ScanThingWithMatchingKeyssInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
	return d.thingWithMatchingKeysTable.scanThingWithMatchingKeyss(ctx, input, fn)
}

// GetThingWithMatchingKeyssByBearAndAssocTypeID retrieves a page of ThingWithMatchingKeyss from the database.
func (d DB) GetThingWithMatchingKeyssByBearAndAssocTypeID(ctx context.Context, input db.GetThingWithMatchingKeyssByBearAndAssocTypeIDInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
	return d.thingWithMatchingKeysTable.getThingWithMatchingKeyssByBearAndAssocTypeID(ctx, input, fn)
}

// DeleteThingWithMatchingKeys deletes a ThingWithMatchingKeys from the database.
func (d DB) DeleteThingWithMatchingKeys(ctx context.Context, bear string, assocType string, assocID string) error {
	return d.thingWithMatchingKeysTable.deleteThingWithMatchingKeys(ctx, bear, assocType, assocID)
}

// GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBear retrieves a page of ThingWithMatchingKeyss from the database.
func (d DB) GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx context.Context, input db.GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
	return d.thingWithMatchingKeysTable.getThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx, input, fn)
}

// ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear runs a scan on the AssocTypeIDAndCreatedBear index.
func (d DB) ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx context.Context, input db.ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
	return d.thingWithMatchingKeysTable.scanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx, input, fn)
}

// SaveThingWithMultiUseCompositeAttribute saves a ThingWithMultiUseCompositeAttribute to the database.
func (d DB) SaveThingWithMultiUseCompositeAttribute(ctx context.Context, m models.ThingWithMultiUseCompositeAttribute) error {
	return d.thingWithMultiUseCompositeAttributeTable.saveThingWithMultiUseCompositeAttribute(ctx, m)
}

// GetThingWithMultiUseCompositeAttribute retrieves a ThingWithMultiUseCompositeAttribute from the database.
func (d DB) GetThingWithMultiUseCompositeAttribute(ctx context.Context, one string) (*models.ThingWithMultiUseCompositeAttribute, error) {
	return d.thingWithMultiUseCompositeAttributeTable.getThingWithMultiUseCompositeAttribute(ctx, one)
}

// ScanThingWithMultiUseCompositeAttributes runs a scan on the ThingWithMultiUseCompositeAttributes table.
func (d DB) ScanThingWithMultiUseCompositeAttributes(ctx context.Context, input db.ScanThingWithMultiUseCompositeAttributesInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	return d.thingWithMultiUseCompositeAttributeTable.scanThingWithMultiUseCompositeAttributes(ctx, input, fn)
}

// DeleteThingWithMultiUseCompositeAttribute deletes a ThingWithMultiUseCompositeAttribute from the database.
func (d DB) DeleteThingWithMultiUseCompositeAttribute(ctx context.Context, one string) error {
	return d.thingWithMultiUseCompositeAttributeTable.deleteThingWithMultiUseCompositeAttribute(ctx, one)
}

// GetThingWithMultiUseCompositeAttributesByThreeAndOneTwo retrieves a page of ThingWithMultiUseCompositeAttributes from the database.
func (d DB) GetThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx context.Context, input db.GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	return d.thingWithMultiUseCompositeAttributeTable.getThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx, input, fn)
}

// ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwo runs a scan on the ThreeAndOneTwo index.
func (d DB) ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx context.Context, input db.ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	return d.thingWithMultiUseCompositeAttributeTable.scanThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx, input, fn)
}

// GetThingWithMultiUseCompositeAttributesByFourAndOneTwo retrieves a page of ThingWithMultiUseCompositeAttributes from the database.
func (d DB) GetThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx context.Context, input db.GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	return d.thingWithMultiUseCompositeAttributeTable.getThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx, input, fn)
}

// ScanThingWithMultiUseCompositeAttributesByFourAndOneTwo runs a scan on the FourAndOneTwo index.
func (d DB) ScanThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx context.Context, input db.ScanThingWithMultiUseCompositeAttributesByFourAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	return d.thingWithMultiUseCompositeAttributeTable.scanThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx, input, fn)
}

// SaveThingWithRequiredCompositePropertiesAndKeysOnly saves a ThingWithRequiredCompositePropertiesAndKeysOnly to the database.
func (d DB) SaveThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, m models.ThingWithRequiredCompositePropertiesAndKeysOnly) error {
	return d.thingWithRequiredCompositePropertiesAndKeysOnlyTable.saveThingWithRequiredCompositePropertiesAndKeysOnly(ctx, m)
}

// GetThingWithRequiredCompositePropertiesAndKeysOnly retrieves a ThingWithRequiredCompositePropertiesAndKeysOnly from the database.
func (d DB) GetThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, propertyThree string) (*models.ThingWithRequiredCompositePropertiesAndKeysOnly, error) {
	return d.thingWithRequiredCompositePropertiesAndKeysOnlyTable.getThingWithRequiredCompositePropertiesAndKeysOnly(ctx, propertyThree)
}

// ScanThingWithRequiredCompositePropertiesAndKeysOnlys runs a scan on the ThingWithRequiredCompositePropertiesAndKeysOnlys table.
func (d DB) ScanThingWithRequiredCompositePropertiesAndKeysOnlys(ctx context.Context, input db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysInput, fn func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, lastThingWithRequiredCompositePropertiesAndKeysOnly bool) bool) error {
	return d.thingWithRequiredCompositePropertiesAndKeysOnlyTable.scanThingWithRequiredCompositePropertiesAndKeysOnlys(ctx, input, fn)
}

// DeleteThingWithRequiredCompositePropertiesAndKeysOnly deletes a ThingWithRequiredCompositePropertiesAndKeysOnly from the database.
func (d DB) DeleteThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, propertyThree string) error {
	return d.thingWithRequiredCompositePropertiesAndKeysOnlyTable.deleteThingWithRequiredCompositePropertiesAndKeysOnly(ctx, propertyThree)
}

// GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree retrieves a page of ThingWithRequiredCompositePropertiesAndKeysOnlys from the database.
func (d DB) GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx context.Context, input db.GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput, fn func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, lastThingWithRequiredCompositePropertiesAndKeysOnly bool) bool) error {
	return d.thingWithRequiredCompositePropertiesAndKeysOnlyTable.getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx, input, fn)
}

// ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree runs a scan on the PropertyOneAndTwoAndPropertyThree index.
func (d DB) ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx context.Context, input db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput, fn func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, lastThingWithRequiredCompositePropertiesAndKeysOnly bool) bool) error {
	return d.thingWithRequiredCompositePropertiesAndKeysOnlyTable.scanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx, input, fn)
}

// SaveThingWithRequiredFields saves a ThingWithRequiredFields to the database.
func (d DB) SaveThingWithRequiredFields(ctx context.Context, m models.ThingWithRequiredFields) error {
	return d.thingWithRequiredFieldsTable.saveThingWithRequiredFields(ctx, m)
}

// GetThingWithRequiredFields retrieves a ThingWithRequiredFields from the database.
func (d DB) GetThingWithRequiredFields(ctx context.Context, name string) (*models.ThingWithRequiredFields, error) {
	return d.thingWithRequiredFieldsTable.getThingWithRequiredFields(ctx, name)
}

// ScanThingWithRequiredFieldss runs a scan on the ThingWithRequiredFieldss table.
func (d DB) ScanThingWithRequiredFieldss(ctx context.Context, input db.ScanThingWithRequiredFieldssInput, fn func(m *models.ThingWithRequiredFields, lastThingWithRequiredFields bool) bool) error {
	return d.thingWithRequiredFieldsTable.scanThingWithRequiredFieldss(ctx, input, fn)
}

// DeleteThingWithRequiredFields deletes a ThingWithRequiredFields from the database.
func (d DB) DeleteThingWithRequiredFields(ctx context.Context, name string) error {
	return d.thingWithRequiredFieldsTable.deleteThingWithRequiredFields(ctx, name)
}

// SaveThingWithRequiredFields2 saves a ThingWithRequiredFields2 to the database.
func (d DB) SaveThingWithRequiredFields2(ctx context.Context, m models.ThingWithRequiredFields2) error {
	return d.thingWithRequiredFields2Table.saveThingWithRequiredFields2(ctx, m)
}

// GetThingWithRequiredFields2 retrieves a ThingWithRequiredFields2 from the database.
func (d DB) GetThingWithRequiredFields2(ctx context.Context, name string, id string) (*models.ThingWithRequiredFields2, error) {
	return d.thingWithRequiredFields2Table.getThingWithRequiredFields2(ctx, name, id)
}

// ScanThingWithRequiredFields2s runs a scan on the ThingWithRequiredFields2s table.
func (d DB) ScanThingWithRequiredFields2s(ctx context.Context, input db.ScanThingWithRequiredFields2sInput, fn func(m *models.ThingWithRequiredFields2, lastThingWithRequiredFields2 bool) bool) error {
	return d.thingWithRequiredFields2Table.scanThingWithRequiredFields2s(ctx, input, fn)
}

// GetThingWithRequiredFields2sByNameAndID retrieves a page of ThingWithRequiredFields2s from the database.
func (d DB) GetThingWithRequiredFields2sByNameAndID(ctx context.Context, input db.GetThingWithRequiredFields2sByNameAndIDInput, fn func(m *models.ThingWithRequiredFields2, lastThingWithRequiredFields2 bool) bool) error {
	return d.thingWithRequiredFields2Table.getThingWithRequiredFields2sByNameAndID(ctx, input, fn)
}

// DeleteThingWithRequiredFields2 deletes a ThingWithRequiredFields2 from the database.
func (d DB) DeleteThingWithRequiredFields2(ctx context.Context, name string, id string) error {
	return d.thingWithRequiredFields2Table.deleteThingWithRequiredFields2(ctx, name, id)
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

func toDynamoTimeString(d strfmt.DateTime) string {
	return time.Time(d).Format(time.RFC3339) // dynamodb attributevalue only supports RFC3339 resolution
}

func toDynamoTimeStringPtr(d *strfmt.DateTime) string {
	return time.Time(*d).Format(time.RFC3339) // dynamodb attributevalue only supports RFC3339 resolution
}
