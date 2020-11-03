package dynamodb

import (
	"context"
	"errors"
	"time"

	"github.com/Clever/wag/samples/gen-go-app-service/models"
	"github.com/Clever/wag/samples/gen-go-app-service/server/db"
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
	// SetupStepTable configuration.
	SetupStepTable SetupStepTable
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
	// configure SetupStep table
	setupStepTable := config.SetupStepTable
	if setupStepTable.DynamoDBAPI == nil {
		setupStepTable.DynamoDBAPI = config.DynamoDBAPI
	}
	if setupStepTable.Prefix == "" {
		setupStepTable.Prefix = config.DefaultPrefix
	}
	if setupStepTable.ReadCapacityUnits == 0 {
		setupStepTable.ReadCapacityUnits = config.DefaultReadCapacityUnits
	}
	if setupStepTable.WriteCapacityUnits == 0 {
		setupStepTable.WriteCapacityUnits = config.DefaultWriteCapacityUnits
	}

	return &DB{
		setupStepTable: setupStepTable,
	}, nil
}

// DB implements the database interface using DynamoDB to store data.
type DB struct {
	setupStepTable SetupStepTable
}

var _ db.Interface = DB{}

// CreateTables creates all tables.
func (d DB) CreateTables(ctx context.Context) error {
	if err := d.setupStepTable.create(ctx); err != nil {
		return err
	}
	return nil
}

// SaveSetupStep saves a SetupStep to the database.
func (d DB) SaveSetupStep(ctx context.Context, m models.SetupStep) error {
	return d.setupStepTable.saveSetupStep(ctx, m)
}

// GetSetupStep retrieves a SetupStep from the database.
func (d DB) GetSetupStep(ctx context.Context, appID string, id string) (*models.SetupStep, error) {
	return d.setupStepTable.getSetupStep(ctx, appID, id)
}

// GetSetupStepsByAppIDAndID retrieves a page of SetupSteps from the database.
func (d DB) GetSetupStepsByAppIDAndID(ctx context.Context, input db.GetSetupStepsByAppIDAndIDInput, fn func(m *models.SetupStep, lastSetupStep bool) bool) error {
	return d.setupStepTable.getSetupStepsByAppIDAndID(ctx, input, fn)
}

// DeleteSetupStep deletes a SetupStep from the database.
func (d DB) DeleteSetupStep(ctx context.Context, appID string, id string) error {
	return d.setupStepTable.deleteSetupStep(ctx, appID, id)
}

func toDynamoTimeString(d strfmt.DateTime) string {
	return time.Time(d).Format(time.RFC3339) // dynamodb attributevalue only supports RFC3339 resolution
}

func toDynamoTimeStringPtr(d *strfmt.DateTime) string {
	return time.Time(*d).Format(time.RFC3339) // dynamodb attributevalue only supports RFC3339 resolution
}
