package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db-only/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db-only/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}
var _ = errors.New("")
var _ = []types.AttributeValue{}

// ThingAllowingBatchWritesTable represents the user-configurable properties of the ThingAllowingBatchWrites table.
type ThingAllowingBatchWritesTable struct {
	DynamoDBAPI        *dynamodb.Client
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
	models.ThingAllowingBatchWrites `dynamodbav:",inline"`
}

func (t ThingAllowingBatchWritesTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: types.ScalarAttributeType("N"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
		},
		TableName: aws.String(t.TableName),
	}); err != nil {
		return fmt.Errorf("failed to create table %s: %w", t.TableName, err)
	}
	return nil
}

func (t ThingAllowingBatchWritesTable) saveThingAllowingBatchWrites(ctx context.Context, m models.ThingAllowingBatchWrites) error {
	data, err := encodeThingAllowingBatchWrites(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#NAME":    "name",
			"#VERSION": "version",
		},
		ConditionExpression: aws.String(
			"" +
				"" +
				"attribute_not_exists(#NAME)" +
				"" +
				" AND " +
				"attribute_not_exists(#VERSION)" +
				"",
		),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		var conditionalCheckFailedErr *types.ConditionalCheckFailedException
		if errors.As(err, &resourceNotFoundErr) {
			return fmt.Errorf("table or index not found: %s", t.TableName)
		}
		if errors.As(err, &conditionalCheckFailedErr) {
			return db.ErrThingAllowingBatchWritesAlreadyExists{
				Name:    m.Name,
				Version: m.Version,
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

	batch := make([]types.WriteRequest, len(ms))
	for i := range ms {
		data, err := encodeThingAllowingBatchWrites(ms[i])
		if err != nil {
			return err
		}
		batch[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: data,
			},
		}
	}
	tname := t.TableName
	for {
		if out, err := t.DynamoDBAPI.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
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

	batch := make([]types.WriteRequest, len(ms))
	for i := range ms {
		key, err := attributevalue.MarshalMap(ddbThingAllowingBatchWritesPrimaryKey{
			Name:    ms[i].Name,
			Version: ms[i].Version,
		})
		if err != nil {
			return err
		}

		batch[i] = types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: key,
			},
		}
	}
	tname := t.TableName
	for {
		if out, err := t.DynamoDBAPI.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
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
	key, err := attributevalue.MarshalMap(ddbThingAllowingBatchWritesPrimaryKey{
		Name:    name,
		Version: version,
	})
	if err != nil {
		return nil, err
	}
	res, err := t.DynamoDBAPI.GetItem(ctx, &dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundErr) {
			return nil, fmt.Errorf("table or index not found: %s", t.TableName)
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
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name":    exclusiveStartKey["name"],
			"version": exclusiveStartKey["version"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingAllowingBatchWritess(out.Items)
		if err != nil {
			return fmt.Errorf("error decoding items: %s", err.Error())
		}

		for i := range items {
			if input.Limiter != nil {
				if err := input.Limiter.Wait(ctx); err != nil {
					return err
				}
			}

			isLastModel := !paginator.HasMorePages() && i == len(items)-1
			if shouldContinue := fn(&items[i], isLastModel); !shouldContinue {
				return nil
			}

			totalRecordsProcessed++
			if input.Limit != nil && totalRecordsProcessed == *input.Limit {
				return nil
			}
		}
	}

	return nil
}

func (t ThingAllowingBatchWritesTable) getThingAllowingBatchWritessByNameAndVersion(ctx context.Context, input db.GetThingAllowingBatchWritessByNameAndVersionInput, fn func(m *models.ThingAllowingBatchWrites, lastThingAllowingBatchWrites bool) bool) error {
	if input.VersionStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.VersionStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		ExpressionAttributeNames: map[string]string{
			"#NAME": "name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{
				Value: input.Name,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.VersionStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		queryInput.ExpressionAttributeNames["#VERSION"] = "version"
		queryInput.ExpressionAttributeValues[":version"] = &types.AttributeValueMemberN{
			Value: fmt.Sprintf("%d", *input.VersionStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #VERSION <= :version")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #VERSION >= :version")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"version": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", input.StartingAfter.Version),
			},

			"name": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Name,
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

	paginator := dynamodb.NewQueryPaginator(t.DynamoDBAPI, queryInput)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			var resourceNotFoundErr *types.ResourceNotFoundException
			if errors.As(err, &resourceNotFoundErr) {
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
			return err
		}
		if !pageFn(output, !paginator.HasMorePages()) {
			break
		}
	}

	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}

func (t ThingAllowingBatchWritesTable) deleteThingAllowingBatchWrites(ctx context.Context, name string, version int64) error {

	key, err := attributevalue.MarshalMap(ddbThingAllowingBatchWritesPrimaryKey{
		Name:    name,
		Version: version,
	})
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(t.TableName),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundErr) {
			return fmt.Errorf("table or index not found: %s", t.TableName)
		}
		return err
	}

	return nil
}

// encodeThingAllowingBatchWrites encodes a ThingAllowingBatchWrites as a DynamoDB map of attribute values.
func encodeThingAllowingBatchWrites(m models.ThingAllowingBatchWrites) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(ddbThingAllowingBatchWrites{
		ThingAllowingBatchWrites: m,
	})
}

// decodeThingAllowingBatchWrites translates a ThingAllowingBatchWrites stored in DynamoDB to a ThingAllowingBatchWrites struct.
func decodeThingAllowingBatchWrites(m map[string]types.AttributeValue, out *models.ThingAllowingBatchWrites) error {
	var ddbThingAllowingBatchWrites ddbThingAllowingBatchWrites
	if err := attributevalue.UnmarshalMap(m, &ddbThingAllowingBatchWrites); err != nil {
		return err
	}
	*out = ddbThingAllowingBatchWrites.ThingAllowingBatchWrites
	return nil
}

// decodeThingAllowingBatchWritess translates a list of ThingAllowingBatchWritess stored in DynamoDB to a slice of ThingAllowingBatchWrites structs.
func decodeThingAllowingBatchWritess(ms []map[string]types.AttributeValue) ([]models.ThingAllowingBatchWrites, error) {
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
