package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingAllowingBatchWritesTable represents the user-configurable properties of the ThingAllowingBatchWrites table.
type ThingAllowingBatchWritesTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingAllowingBatchWritesPrimaryKey represents the primary key of a ThingAllowingBatchWrites in DynamoDB.
type ddbThingAllowingBatchWritesPrimaryKey struct {
	Name    string `dynamodbav:"name"`
	Version int64  `dynamodbav:"version"`
}

// ddbThingAllowingBatchWrites represents a ThingAllowingBatchWrites as stored in DynamoDB.
type ddbThingAllowingBatchWrites struct {
	models.ThingAllowingBatchWrites
}

func (t ThingAllowingBatchWritesTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-allowing-batch-writess", t.Prefix)
}

func (t ThingAllowingBatchWritesTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("version"),
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

func (t ThingAllowingBatchWritesTable) saveThingAllowingBatchWrites(ctx context.Context, m models.ThingAllowingBatchWrites) error {
	data, err := encodeThingAllowingBatchWrites(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
		ExpressionAttributeNames: map[string]*string{
			"#NAME":    aws.String("name"),
			"#VERSION": aws.String("version"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#NAME) AND attribute_not_exists(#VERSION)"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return db.ErrThingAllowingBatchWritesAlreadyExists{
					Name:    m.Name,
					Version: m.Version,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	return nil
}
func (t ThingAllowingBatchWritesTable) saveArrayOfThingAllowingBatchWrites(ctx context.Context, ms []models.ThingAllowingBatchWrites) error {
	if len(ms) > maxDynamoDBBatchItems {
		return fmt.Errorf("saveArrayOfThingAllowingBatchWrites received %d items to save, which is greater than the maximum of %d", len(ms), maxDynamoDBBatchItems)
	}

	if len(ms) == 0 {
		return nil
	}

	batch := make([]*dynamodb.WriteRequest, len(ms))
	for i := range ms {
		data, err := encodeThingAllowingBatchWrites(ms[i])
		if err != nil {
			return err
		}
		batch[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: data,
			},
		}
	}
	tname := t.name()
	for {
		if out, err := t.DynamoDBAPI.BatchWriteItemWithContext(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				tname: batch,
			},
		}); err != nil {
			return fmt.Errorf("BatchWriteItem: %v", err)
		} else if out.UnprocessedItems != nil && len(out.UnprocessedItems[tname]) > 0 {
			batch = out.UnprocessedItems[tname]
		} else {
			break
		}
	}
	return nil
}

func (t ThingAllowingBatchWritesTable) deleteArrayOfThingAllowingBatchWrites(ctx context.Context, ms []models.ThingAllowingBatchWrites) error {
	if len(ms) > maxDynamoDBBatchItems {
		return fmt.Errorf("deleteArrayOfThingAllowingBatchWrites received %d items to delete, which is greater than the maximum of %d", len(ms), maxDynamoDBBatchItems)
	}

	if len(ms) == 0 {
		return nil
	}

	batch := make([]*dynamodb.WriteRequest, len(ms))
	for i := range ms {
		key, err := dynamodbattribute.MarshalMap(ddbThingAllowingBatchWritesPrimaryKey{
			Name:    ms[i].Name,
			Version: ms[i].Version,
		})
		if err != nil {
			return err
		}

		batch[i] = &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: key,
			},
		}
	}
	tname := t.name()
	for {
		if out, err := t.DynamoDBAPI.BatchWriteItemWithContext(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				tname: batch,
			},
		}); err != nil {
			return fmt.Errorf("BatchWriteItem: %v", err)
		} else if out.UnprocessedItems != nil && len(out.UnprocessedItems[tname]) > 0 {
			batch = out.UnprocessedItems[tname]
		} else {
			break
		}
	}
	return nil
}

func (t ThingAllowingBatchWritesTable) getThingAllowingBatchWrites(ctx context.Context, name string, version int64) (*models.ThingAllowingBatchWrites, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingAllowingBatchWritesPrimaryKey{
		Name:    name,
		Version: version,
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
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return nil, fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrThingAllowingBatchWritesNotFound{
			Name:    name,
			Version: version,
		}
	}

	var m models.ThingAllowingBatchWrites
	if err := decodeThingAllowingBatchWrites(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingAllowingBatchWritesTable) scanThingAllowingBatchWritess(ctx context.Context, input db.ScanThingAllowingBatchWritessInput, fn func(m *models.ThingAllowingBatchWrites, lastThingAllowingBatchWrites bool) bool) error {
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
			"name":    exclusiveStartKey["name"],
			"version": exclusiveStartKey["version"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingAllowingBatchWritess(out.Items)
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

func (t ThingAllowingBatchWritesTable) getThingAllowingBatchWritessByNameAndVersion(ctx context.Context, input db.GetThingAllowingBatchWritessByNameAndVersionInput, fn func(m *models.ThingAllowingBatchWrites, lastThingAllowingBatchWrites bool) bool) error {
	if input.VersionStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.VersionStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
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
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.VersionStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		queryInput.ExpressionAttributeNames["#VERSION"] = aws.String("version")
		queryInput.ExpressionAttributeValues[":version"] = &dynamodb.AttributeValue{
			N: aws.String(fmt.Sprintf("%d", *input.VersionStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #VERSION <= :version")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #VERSION >= :version")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"version": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%d", input.StartingAfter.Version)),
			},
			"name": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Name),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingAllowingBatchWritess(queryOutput.Items)
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
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}

func (t ThingAllowingBatchWritesTable) deleteThingAllowingBatchWrites(ctx context.Context, name string, version int64) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingAllowingBatchWritesPrimaryKey{
		Name:    name,
		Version: version,
	})
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(t.name()),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}

	return nil
}

// encodeThingAllowingBatchWrites encodes a ThingAllowingBatchWrites as a DynamoDB map of attribute values.
func encodeThingAllowingBatchWrites(m models.ThingAllowingBatchWrites) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingAllowingBatchWrites{
		ThingAllowingBatchWrites: m,
	})
}

// decodeThingAllowingBatchWrites translates a ThingAllowingBatchWrites stored in DynamoDB to a ThingAllowingBatchWrites struct.
func decodeThingAllowingBatchWrites(m map[string]*dynamodb.AttributeValue, out *models.ThingAllowingBatchWrites) error {
	var ddbThingAllowingBatchWrites ddbThingAllowingBatchWrites
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingAllowingBatchWrites); err != nil {
		return err
	}
	*out = ddbThingAllowingBatchWrites.ThingAllowingBatchWrites
	return nil
}

// decodeThingAllowingBatchWritess translates a list of ThingAllowingBatchWritess stored in DynamoDB to a slice of ThingAllowingBatchWrites structs.
func decodeThingAllowingBatchWritess(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingAllowingBatchWrites, error) {
	thingAllowingBatchWritess := make([]models.ThingAllowingBatchWrites, len(ms))
	for i, m := range ms {
		var thingAllowingBatchWrites models.ThingAllowingBatchWrites
		if err := decodeThingAllowingBatchWrites(m, &thingAllowingBatchWrites); err != nil {
			return nil, err
		}
		thingAllowingBatchWritess[i] = thingAllowingBatchWrites
	}
	return thingAllowingBatchWritess, nil
}
