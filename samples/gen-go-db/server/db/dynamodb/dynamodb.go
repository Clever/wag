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
	// ThingWithDateRangeTable configuration.
	ThingWithDateRangeTable ThingWithDateRangeTable
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

	return &DB{
		simpleThingTable:        simpleThingTable,
		thingTable:              thingTable,
		thingWithDateRangeTable: thingWithDateRangeTable,
	}, nil
}

// DB implements the database interface using DynamoDB to store data.
type DB struct {
	simpleThingTable        SimpleThingTable
	thingTable              ThingTable
	thingWithDateRangeTable ThingWithDateRangeTable
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
	if err := d.thingWithDateRangeTable.create(ctx); err != nil {
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

// DeleteThing deletes a Thing from the database.
func (d DB) DeleteThing(ctx context.Context, name string, version int64) error {
	return d.thingTable.deleteThing(ctx, name, version)
}

// SaveThingWithDateRange saves a ThingWithDateRange to the database.
func (d DB) SaveThingWithDateRange(ctx context.Context, m models.ThingWithDateRange) error {
	return d.thingWithDateRangeTable.saveThingWithDateRange(ctx, m)
}

// GetThingWithDateRange retrieves a ThingWithDateRange from the database.
func (d DB) GetThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) (*models.ThingWithDateRange, error) {
	return d.thingWithDateRangeTable.getThingWithDateRange(ctx, name, date)
}

// DeleteThingWithDateRange deletes a ThingWithDateRange from the database.
func (d DB) DeleteThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) error {
	return d.thingWithDateRangeTable.deleteThingWithDateRange(ctx, name, date)
}
