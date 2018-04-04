package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db/models"
	"github.com/Clever/wag/samples/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

// ThingWithDateRangeTable represents the user-configurable properties of the ThingWithDateRange table.
type ThingWithDateRangeTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithDateRangePrimaryKey represents the primary key of a ThingWithDateRange in DynamoDB.
type ddbThingWithDateRangePrimaryKey struct {
	Name string          `dynamodbav:"name"`
	Date strfmt.DateTime `dynamodbav:"date"`
}

// ddbThingWithDateRange represents a ThingWithDateRange as stored in DynamoDB.
type ddbThingWithDateRange struct {
	ddbThingWithDateRangePrimaryKey
	ThingWithDateRange models.ThingWithDateRange `dynamodbav:"thing-with-date-range"`
}

func (t ThingWithDateRangeTable) name() string {
	return fmt.Sprintf("%s-thing-with-date-ranges", t.Prefix)
}

func (t ThingWithDateRangeTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
			{
				AttributeName: aws.String("date"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("date"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
		},
		TableName: aws.String(t.name()),
	}); err != nil {
		return err
	}
	return nil
}

func (t ThingWithDateRangeTable) saveThingWithDateRange(ctx context.Context, m models.ThingWithDateRange) error {
	data, err := encodeThingWithDateRange(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
	})
	return err
}

func (t ThingWithDateRangeTable) getThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) (*models.ThingWithDateRange, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithDateRangePrimaryKey{
		Name: name,
		Date: date,
	})
	if err != nil {
		return nil, err
	}
	res, err := t.DynamoDBAPI.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(t.name()),
	})
	if err != nil {
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrThingWithDateRangeNotFound{
			Name: name,
			Date: date,
		}
	}

	var m models.ThingWithDateRange
	if err := decodeThingWithDateRange(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithDateRangeTable) deleteThingWithDateRange(ctx context.Context, name string, date strfmt.DateTime) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithDateRangePrimaryKey{
		Name: name,
		Date: date,
	})
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(t.name()),
	})
	if err != nil {
		return err
	}
	return nil
}

// encodeThingWithDateRange encodes a ThingWithDateRange as a DynamoDB map of attribute values.
func encodeThingWithDateRange(m models.ThingWithDateRange) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingWithDateRange{
		ddbThingWithDateRangePrimaryKey: ddbThingWithDateRangePrimaryKey{
			Name: m.Name,
			Date: m.Date,
		},
		ThingWithDateRange: m,
	})
}

// decodeThingWithDateRange translates a ThingWithDateRange stored in DynamoDB to a ThingWithDateRange struct.
func decodeThingWithDateRange(m map[string]*dynamodb.AttributeValue, out *models.ThingWithDateRange) error {
	var ddbThingWithDateRange ddbThingWithDateRange
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithDateRange); err != nil {
		return err
	}
	*out = ddbThingWithDateRange.ThingWithDateRange
	return nil
}
