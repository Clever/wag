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

var _ = strfmt.DateTime{}

// ThingWithDateRangeTable represents the user-configurable properties of the ThingWithDateRange table.
type ThingWithDateRangeTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
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
	models.ThingWithDateRange
}

func (t ThingWithDateRangeTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-date-ranges", t.Prefix)
}

func (t ThingWithDateRangeTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("date"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
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

func (t ThingWithDateRangeTable) getThingWithDateRangesByNameAndDate(ctx context.Context, input db.GetThingWithDateRangesByNameAndDateInput) ([]models.ThingWithDateRange, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#NAME": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": &dynamodb.AttributeValue{
				S: aws.String(input.Name),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = aws.String("date")
		queryInput.ExpressionAttributeValues[":date"] = &dynamodb.AttributeValue{
			S: aws.String(toDynamoTimeString(*input.DateStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #DATE <= :date")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #DATE >= :date")
		}
	}

	queryOutput, err := t.DynamoDBAPI.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if len(queryOutput.Items) == 0 {
		return []models.ThingWithDateRange{}, nil
	}

	return decodeThingWithDateRanges(queryOutput.Items)
}

func (t ThingWithDateRangeTable) getThingWithDateRangesByNameAndDatePage(ctx context.Context, input db.GetThingWithDateRangesByNameAndDatePageInput, fn func(m *models.ThingWithDateRange, lastThingWithDateRange bool) bool) error {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#NAME": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Name),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
		Limit:            input.Limit,
		ExclusiveStartKey: map[string]*dynamodb.AttributeValue{
			"date": &dynamodb.AttributeValue{
				S: aws.String(toDynamoTimeString(input.StartingAfter.Date)),
			},
			"name": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Name),
			},
		},
	}
	queryInput.ExpressionAttributeNames["#DATE"] = aws.String("date")
	queryInput.ExpressionAttributeValues[":date"] = &dynamodb.AttributeValue{
		S: aws.String(toDynamoTimeString(input.StartingAfter.Date)),
	}
	if input.Descending {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #DATE <= :date")
	} else {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #DATE >= :date")
	}

	queryOutput, err := t.DynamoDBAPI.QueryWithContext(ctx, queryInput)
	if err != nil {
		return err
	}
	if len(queryOutput.Items) == 0 {
		fn(nil, false)
		return nil
	}

	items, err := decodeThingWithDateRanges(queryOutput.Items)
	if err != nil {
		return err
	}

	for i, item := range items {
		hasMore := false
		if len(queryOutput.LastEvaluatedKey) > 0 {
			hasMore = true
		} else {
			hasMore = i < len(items)-1
		}
		if !fn(&item, !hasMore) {
			break
		}
	}

	return nil
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

// decodeThingWithDateRanges translates a list of ThingWithDateRanges stored in DynamoDB to a slice of ThingWithDateRange structs.
func decodeThingWithDateRanges(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithDateRange, error) {
	thingWithDateRanges := make([]models.ThingWithDateRange, len(ms))
	for i, m := range ms {
		var thingWithDateRange models.ThingWithDateRange
		if err := decodeThingWithDateRange(m, &thingWithDateRange); err != nil {
			return nil, err
		}
		thingWithDateRanges[i] = thingWithDateRange
	}
	return thingWithDateRanges, nil
}
