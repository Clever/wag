package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/v7/samples/gen-go-db/models"
	"github.com/Clever/wag/v7/samples/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// EventTable represents the user-configurable properties of the Event table.
type EventTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbEventPrimaryKey represents the primary key of a Event in DynamoDB.
type ddbEventPrimaryKey struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
}

// ddbEventGSIBySK represents the bySK GSI.
type ddbEventGSIBySK struct {
	Sk   string `dynamodbav:"sk"`
	Data []byte `dynamodbav:"data"`
}

// ddbEvent represents a Event as stored in DynamoDB.
type ddbEvent struct {
	models.Event
}

func (t EventTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-events", t.Prefix)
}

func (t EventTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("data"),
				AttributeType: aws.String("B"),
			},
			{
				AttributeName: aws.String("pk"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("sk"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("sk"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("bySK"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("sk"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("data"),
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
		TableName: aws.String(t.name()),
	}); err != nil {
		return err
	}
	return nil
}

func (t EventTable) saveEvent(ctx context.Context, m models.Event) error {
	data, err := encodeEvent(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
	})
	return err
}

func (t EventTable) getEvent(ctx context.Context, pk string, sk string) (*models.Event, error) {
	key, err := dynamodbattribute.MarshalMap(ddbEventPrimaryKey{
		Pk: pk,
		Sk: sk,
	})
	if err != nil {
		return nil, err
	}
	res, err := t.DynamoDBAPI.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrEventNotFound{
			Pk: pk,
			Sk: sk,
		}
	}

	var m models.Event
	if err := decodeEvent(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t EventTable) scanEvents(ctx context.Context, input db.ScanEventsInput, fn func(m *models.Event, lastEvent bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
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
			"pk": exclusiveStartKey["pk"],
			"sk": exclusiveStartKey["sk"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeEvents(out.Items)
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

func (t EventTable) getEventsByPkAndSk(ctx context.Context, input db.GetEventsByPkAndSkInput, fn func(m *models.Event, lastEvent bool) bool) error {
	if input.SkStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.SkStartingAt or input.StartingAfter")
	}
	if input.Pk == "" {
		return fmt.Errorf("Hash key input.Pk cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#PK": aws.String("pk"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": &dynamodb.AttributeValue{
				S: aws.String(input.Pk),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.SkStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#PK = :pk")
	} else {
		queryInput.ExpressionAttributeNames["#SK"] = aws.String("sk")
		queryInput.ExpressionAttributeValues[":sk"] = &dynamodb.AttributeValue{
			S: aws.String(*input.SkStartingAt),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#PK = :pk AND #SK <= :sk")
		} else {
			queryInput.KeyConditionExpression = aws.String("#PK = :pk AND #SK >= :sk")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"sk": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Sk),
			},
			"pk": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Pk),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeEvents(queryOutput.Items)
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
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}

func (t EventTable) deleteEvent(ctx context.Context, pk string, sk string) error {
	key, err := dynamodbattribute.MarshalMap(ddbEventPrimaryKey{
		Pk: pk,
		Sk: sk,
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

func (t EventTable) getEventsBySkAndData(ctx context.Context, input db.GetEventsBySkAndDataInput, fn func(m *models.Event, lastEvent bool) bool) error {
	if input.DataStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DataStartingAt or input.StartingAfter")
	}
	if input.Sk == "" {
		return fmt.Errorf("Hash key input.Sk cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("bySK"),
		ExpressionAttributeNames: map[string]*string{
			"#SK": aws.String("sk"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sk": &dynamodb.AttributeValue{
				S: aws.String(input.Sk),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.DataStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#SK = :sk")
	} else {
		queryInput.ExpressionAttributeNames["#DATA"] = aws.String("data")
		queryInput.ExpressionAttributeValues[":data"] = &dynamodb.AttributeValue{
			B: input.DataStartingAt,
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#SK = :sk AND #DATA <= :data")
		} else {
			queryInput.KeyConditionExpression = aws.String("#SK = :sk AND #DATA >= :data")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"data": &dynamodb.AttributeValue{
				B: input.StartingAfter.Data,
			},
			"sk": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Sk),
			},
			"pk": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Pk),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeEvents(queryOutput.Items)
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
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}
func (t EventTable) scanEventsBySkAndData(ctx context.Context, input db.ScanEventsBySkAndDataInput, fn func(m *models.Event, lastEvent bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("bySK"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"pk":   exclusiveStartKey["pk"],
			"sk":   exclusiveStartKey["sk"],
			"data": exclusiveStartKey["data"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeEvents(out.Items)
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

// encodeEvent encodes a Event as a DynamoDB map of attribute values.
func encodeEvent(m models.Event) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbEvent{
		Event: m,
	})
}

// decodeEvent translates a Event stored in DynamoDB to a Event struct.
func decodeEvent(m map[string]*dynamodb.AttributeValue, out *models.Event) error {
	var ddbEvent ddbEvent
	if err := dynamodbattribute.UnmarshalMap(m, &ddbEvent); err != nil {
		return err
	}
	*out = ddbEvent.Event
	return nil
}

// decodeEvents translates a list of Events stored in DynamoDB to a slice of Event structs.
func decodeEvents(ms []map[string]*dynamodb.AttributeValue) ([]models.Event, error) {
	events := make([]models.Event, len(ms))
	for i, m := range ms {
		var event models.Event
		if err := decodeEvent(m, &event); err != nil {
			return nil, err
		}
		events[i] = event
	}
	return events, nil
}
