package dynamodb

import (
	"context"
	"errors"

	"github.com/Clever/wag/samples/gen-go-db/models"
	"github.com/Clever/wag/samples/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
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
	// ThingTable configuration.
	ThingTable ThingTable
	// ThingWithCompositeAttributesTable configuration.
	ThingWithCompositeAttributesTable ThingWithCompositeAttributesTable
	// ThingWithDateRangeTable configuration.
	ThingWithDateRangeTable ThingWithDateRangeTable
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
		simpleThingTable:                  simpleThingTable,
		thingTable:                        thingTable,
		thingWithCompositeAttributesTable: thingWithCompositeAttributesTable,
		thingWithDateRangeTable:           thingWithDateRangeTable,
		thingWithUnderscoresTable:         thingWithUnderscoresTable,
	}, nil
}

// DB implements the database interface using DynamoDB to store data.
type DB struct {
	simpleThingTable                  SimpleThingTable
	thingTable                        ThingTable
	thingWithCompositeAttributesTable ThingWithCompositeAttributesTable
	thingWithDateRangeTable           ThingWithDateRangeTable
	thingWithUnderscoresTable         ThingWithUnderscoresTable
}

var _ db.Interface = DB{}

// CreateTables creates all tables.
func (d DB) CreateTables(ctx context.Context) error {
	if err := d.simpleThingTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithCompositeAttributesTable.create(ctx); err != nil {
		return err
	}
	if err := d.thingWithDateRangeTable.create(ctx); err != nil {
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

// SaveThing saves a Thing to the database.
func (d DB) SaveThing(ctx context.Context, m models.Thing) error {
	return d.thingTable.saveThing(ctx, m)
}

// GetThing retrieves a Thing from the database.
func (d DB) GetThing(ctx context.Context, name string, version int64) (*models.Thing, error) {
	return d.thingTable.getThing(ctx, name, version)
}

// GetThingsByNameAndVersion retrieves a list of Things from the database.
func (d DB) GetThingsByNameAndVersion(ctx context.Context, input db.GetThingsByNameAndVersionInput) ([]models.Thing, error) {
	return d.thingTable.getThingsByNameAndVersion(ctx, input)
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

// DeleteThingWithCompositeAttributes deletes a ThingWithCompositeAttributes from the database.
func (d DB) DeleteThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) error {
	return d.thingWithCompositeAttributesTable.deleteThingWithCompositeAttributes(ctx, name, branch, date)
}

// GetThingWithCompositeAttributessByNameVersionAndDate retrieves a list of ThingWithCompositeAttributess from the database.
func (d DB) GetThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input db.GetThingWithCompositeAttributessByNameVersionAndDateInput) ([]models.ThingWithCompositeAttributes, error) {
	return d.thingWithCompositeAttributesTable.getThingWithCompositeAttributessByNameVersionAndDate(ctx, input)
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

// DeleteThingWithDateRange deletes a ThingWithDateRange from the database.
func (d DB) DeleteThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) error {
	return d.thingWithDateRangeTable.deleteThingWithDateRange(ctx, name, date)
}

// SaveThingWithUnderscores saves a ThingWithUnderscores to the database.
func (d DB) SaveThingWithUnderscores(ctx context.Context, m models.ThingWithUnderscores) error {
	return d.thingWithUnderscoresTable.saveThingWithUnderscores(ctx, m)
}

// GetThingWithUnderscores retrieves a ThingWithUnderscores from the database.
func (d DB) GetThingWithUnderscores(ctx context.Context, idApp string) (*models.ThingWithUnderscores, error) {
	return d.thingWithUnderscoresTable.getThingWithUnderscores(ctx, idApp)
}

// DeleteThingWithUnderscores deletes a ThingWithUnderscores from the database.
func (d DB) DeleteThingWithUnderscores(ctx context.Context, idApp string) error {
	return d.thingWithUnderscoresTable.deleteThingWithUnderscores(ctx, idApp)
}
