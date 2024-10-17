package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db-custom-path/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db-custom-path/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithDatetimeGSITable represents the user-configurable properties of the ThingWithDatetimeGSI table.
type ThingWithDatetimeGSITable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithDatetimeGSIPrimaryKey represents the primary key of a ThingWithDatetimeGSI in DynamoDB.
type ddbThingWithDatetimeGSIPrimaryKey struct {
	ID string `dynamodbav:"id"`
}

// ddbThingWithDatetimeGSIGSIByDateTime represents the byDateTime GSI.
type ddbThingWithDatetimeGSIGSIByDateTime struct {
	Datetime strfmt.DateTime `dynamodbav:"datetime"`
	ID       string          `dynamodbav:"id"`
}

// ddbThingWithDatetimeGSI represents a ThingWithDatetimeGSI as stored in DynamoDB.
type ddbThingWithDatetimeGSI struct {
	models.ThingWithDatetimeGSI
}

func (t ThingWithDatetimeGSITable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("datetime"),
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
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("byDateTime"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("datetime"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("id"),
						KeyType:       aws.String(dynamodb.KeyTypeRange),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
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

func (t ThingWithDatetimeGSITable) saveThingWithDatetimeGSI(ctx context.Context, m models.ThingWithDatetimeGSI) error {
	data, err := encodeThingWithDatetimeGSI(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]*string{
			"#ID": aws.String("id"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#ID)"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return db.ErrThingWithDatetimeGSIAlreadyExists{
					ID: m.ID,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
		}
		return err
	}
	return nil
}

func (t ThingWithDatetimeGSITable) getThingWithDatetimeGSI(ctx context.Context, id string) (*models.ThingWithDatetimeGSI, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithDatetimeGSIPrimaryKey{
		ID: id,
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
		return nil, db.ErrThingWithDatetimeGSINotFound{
			ID: id,
		}
	}

	var m models.ThingWithDatetimeGSI
	if err := decodeThingWithDatetimeGSI(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithDatetimeGSITable) scanThingWithDatetimeGSIs(ctx context.Context, input db.ScanThingWithDatetimeGSIsInput, fn func(m *models.ThingWithDatetimeGSI, lastThingWithDatetimeGSI bool) bool) error {
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
			"id": exclusiveStartKey["id"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithDatetimeGSIs(out.Items)
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

func (t ThingWithDatetimeGSITable) deleteThingWithDatetimeGSI(ctx context.Context, id string) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithDatetimeGSIPrimaryKey{
		ID: id,
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

func (t ThingWithDatetimeGSITable) getThingWithDatetimeGSIsByDatetimeAndID(ctx context.Context, input db.GetThingWithDatetimeGSIsByDatetimeAndIDInput, fn func(m *models.ThingWithDatetimeGSI, lastThingWithDatetimeGSI bool) bool) error {
	if input.IDStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.IDStartingAt or input.StartingAfter")
	}
	if datetimeToDynamoTimeString(input.Datetime) == "" {
		return fmt.Errorf("Hash key input.Datetime cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("byDateTime"),
		ExpressionAttributeNames: map[string]*string{
			"#DATETIME": aws.String("datetime"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":datetime": &dynamodb.AttributeValue{
				S: aws.String(datetimeToDynamoTimeString(input.Datetime)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.IDStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#DATETIME = :datetime")
	} else {
		queryInput.ExpressionAttributeNames["#ID"] = aws.String("id")
		queryInput.ExpressionAttributeValues[":id"] = &dynamodb.AttributeValue{
			S: aws.String(string(*input.IDStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#DATETIME = :datetime AND #ID <= :id")
		} else {
			queryInput.KeyConditionExpression = aws.String("#DATETIME = :datetime AND #ID >= :id")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(string(input.StartingAfter.ID)),
			},
			"datetime": &dynamodb.AttributeValue{
				S: aws.String(datetimeToDynamoTimeString(input.StartingAfter.Datetime)),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithDatetimeGSIs(queryOutput.Items)
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
func (t ThingWithDatetimeGSITable) scanThingWithDatetimeGSIsByDatetimeAndID(ctx context.Context, input db.ScanThingWithDatetimeGSIsByDatetimeAndIDInput, fn func(m *models.ThingWithDatetimeGSI, lastThingWithDatetimeGSI bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("byDateTime"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id":       exclusiveStartKey["id"],
			"datetime": exclusiveStartKey["datetime"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithDatetimeGSIs(out.Items)
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

// encodeThingWithDatetimeGSI encodes a ThingWithDatetimeGSI as a DynamoDB map of attribute values.
func encodeThingWithDatetimeGSI(m models.ThingWithDatetimeGSI) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingWithDatetimeGSI{
		ThingWithDatetimeGSI: m,
	})
}

// decodeThingWithDatetimeGSI translates a ThingWithDatetimeGSI stored in DynamoDB to a ThingWithDatetimeGSI struct.
func decodeThingWithDatetimeGSI(m map[string]*dynamodb.AttributeValue, out *models.ThingWithDatetimeGSI) error {
	var ddbThingWithDatetimeGSI ddbThingWithDatetimeGSI
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithDatetimeGSI); err != nil {
		return err
	}
	*out = ddbThingWithDatetimeGSI.ThingWithDatetimeGSI
	return nil
}

// decodeThingWithDatetimeGSIs translates a list of ThingWithDatetimeGSIs stored in DynamoDB to a slice of ThingWithDatetimeGSI structs.
func decodeThingWithDatetimeGSIs(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithDatetimeGSI, error) {
	thingWithDatetimeGSIs := make([]models.ThingWithDatetimeGSI, len(ms))
	for i, m := range ms {
		var thingWithDatetimeGSI models.ThingWithDatetimeGSI
		if err := decodeThingWithDatetimeGSI(m, &thingWithDatetimeGSI); err != nil {
			return nil, err
		}
		thingWithDatetimeGSIs[i] = thingWithDatetimeGSI
	}
	return thingWithDatetimeGSIs, nil
}
