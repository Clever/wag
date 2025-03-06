package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db-only/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db-only/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithDateRangeKeyTable represents the user-configurable properties of the ThingWithDateRangeKey table.
type ThingWithDateRangeKeyTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithDateRangeKeyPrimaryKey represents the primary key of a ThingWithDateRangeKey in DynamoDB.
type ddbThingWithDateRangeKeyPrimaryKey struct {
	ID   string      `dynamodbav:"id"`
	Date strfmt.Date `dynamodbav:"date"`
}

// ddbThingWithDateRangeKey represents a ThingWithDateRangeKey as stored in DynamoDB.
type ddbThingWithDateRangeKey struct {
	models.ThingWithDateRangeKey
}

func (t ThingWithDateRangeKeyTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("date"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
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
		TableName: aws.String(t.TableName),
	}); err != nil {
		return fmt.Errorf("failed to create table %s: %w", t.TableName, err)
	}
	return nil
}

func (t ThingWithDateRangeKeyTable) saveThingWithDateRangeKey(ctx context.Context, m models.ThingWithDateRangeKey) error {
	data, err := encodeThingWithDateRangeKey(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]*string{
			"#ID":   aws.String("id"),
			"#DATE": aws.String("date"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#ID) AND attribute_not_exists(#DATE)"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return db.ErrThingWithDateRangeKeyAlreadyExists{
					ID:   m.ID,
					Date: m.Date,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
		}
		return err
	}
	return nil
}

func (t ThingWithDateRangeKeyTable) getThingWithDateRangeKey(ctx context.Context, id string, date strfmt.Date) (*models.ThingWithDateRangeKey, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithDateRangeKeyPrimaryKey{
		ID:   id,
		Date: date,
	})
	if err != nil {
		return nil, err
	}
	res, err := t.DynamoDBAPI.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return nil, fmt.Errorf("table or index not found: %s", t.TableName)
			}
		}
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrThingWithDateRangeKeyNotFound{
			ID:   id,
			Date: date,
		}
	}

	var m models.ThingWithDateRangeKey
	if err := decodeThingWithDateRangeKey(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithDateRangeKeyTable) scanThingWithDateRangeKeys(ctx context.Context, input db.ScanThingWithDateRangeKeysInput, fn func(m *models.ThingWithDateRangeKey, lastThingWithDateRangeKey bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id":   exclusiveStartKey["id"],
			"date": exclusiveStartKey["date"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithDateRangeKeys(out.Items)
		if err != nil {
			innerErr = fmt.Errorf("error decoding %s", err.Error())
			return false
		}
		for i := range items {
			if input.Limiter != nil {
				if err := input.Limiter.Wait(ctx); err != nil {
					innerErr = err
					return false
				}
			}
			isLastModel := lastPage && i == len(items)-1
			if shouldContinue := fn(&items[i], isLastModel); !shouldContinue {
				return false
			}
			totalRecordsProcessed++
			// if the Limit of records have been passed to fn, don't pass anymore records.
			if input.Limit != nil && totalRecordsProcessed == *input.Limit {
				return false
			}
		}
		return true
	})
	if innerErr != nil {
		return innerErr
	}
	return err
}

func (t ThingWithDateRangeKeyTable) getThingWithDateRangeKeysByIDAndDate(ctx context.Context, input db.GetThingWithDateRangeKeysByIDAndDateInput, fn func(m *models.ThingWithDateRangeKey, lastThingWithDateRangeKey bool) bool) error {
	if input.DateStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateStartingAt or input.StartingAfter")
	}
	if input.ID == "" {
		return fmt.Errorf("Hash key input.ID cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		ExpressionAttributeNames: map[string]*string{
			"#ID": aws.String("id"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": &dynamodb.AttributeValue{
				S: aws.String(input.ID),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#ID = :id")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = aws.String("date")
		queryInput.ExpressionAttributeValues[":date"] = &dynamodb.AttributeValue{
			S: aws.String(dateToDynamoTimeString(*input.DateStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#ID = :id AND #DATE <= :date")
		} else {
			queryInput.KeyConditionExpression = aws.String("#ID = :id AND #DATE >= :date")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"date": &dynamodb.AttributeValue{
				S: aws.String(dateToDynamoTimeString(input.StartingAfter.Date)),
			},
			"id": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.ID),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithDateRangeKeys(queryOutput.Items)
		if err != nil {
			pageFnErr = err
			return false
		}
		hasMore := true
		for i := range items {
			if lastPage == true {
				hasMore = i < len(items)-1
			}
			if !fn(&items[i], !hasMore) {
				return false
			}
			totalRecordsProcessed++
			// if the Limit of records have been passed to fn, don't pass anymore records.
			if input.Limit != nil && totalRecordsProcessed == *input.Limit {
				return false
			}
		}
		return true
	}

	err := t.DynamoDBAPI.QueryPagesWithContext(ctx, queryInput, pageFn)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
		}
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}

func (t ThingWithDateRangeKeyTable) deleteThingWithDateRangeKey(ctx context.Context, id string, date strfmt.Date) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithDateRangeKeyPrimaryKey{
		ID:   id,
		Date: date,
	})
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(t.TableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
		}
		return err
	}

	return nil
}

// encodeThingWithDateRangeKey encodes a ThingWithDateRangeKey as a DynamoDB map of attribute values.
func encodeThingWithDateRangeKey(m models.ThingWithDateRangeKey) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingWithDateRangeKey{
		ThingWithDateRangeKey: m,
	})
}

// decodeThingWithDateRangeKey translates a ThingWithDateRangeKey stored in DynamoDB to a ThingWithDateRangeKey struct.
func decodeThingWithDateRangeKey(m map[string]*dynamodb.AttributeValue, out *models.ThingWithDateRangeKey) error {
	var ddbThingWithDateRangeKey ddbThingWithDateRangeKey
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithDateRangeKey); err != nil {
		return err
	}
	*out = ddbThingWithDateRangeKey.ThingWithDateRangeKey
	return nil
}

// decodeThingWithDateRangeKeys translates a list of ThingWithDateRangeKeys stored in DynamoDB to a slice of ThingWithDateRangeKey structs.
func decodeThingWithDateRangeKeys(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithDateRangeKey, error) {
	thingWithDateRangeKeys := make([]models.ThingWithDateRangeKey, len(ms))
	for i, m := range ms {
		var thingWithDateRangeKey models.ThingWithDateRangeKey
		if err := decodeThingWithDateRangeKey(m, &thingWithDateRangeKey); err != nil {
			return nil, err
		}
		thingWithDateRangeKeys[i] = thingWithDateRangeKey
	}
	return thingWithDateRangeKeys, nil
}
