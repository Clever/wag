package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}
var _ = errors.New("")
var _ = []types.AttributeValue{}

// ThingWithRequiredFieldsTable represents the user-configurable properties of the ThingWithRequiredFields table.
type ThingWithRequiredFieldsTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithRequiredFieldsPrimaryKey represents the primary key of a ThingWithRequiredFields in DynamoDB.
type ddbThingWithRequiredFieldsPrimaryKey struct {
	Name string `dynamodbav:"name"`
}

// ddbThingWithRequiredFields represents a ThingWithRequiredFields as stored in DynamoDB.
type ddbThingWithRequiredFields struct {
	models.ThingWithRequiredFields
}

func (t ThingWithRequiredFieldsTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       types.KeyTypeHash,
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

func (t ThingWithRequiredFieldsTable) saveThingWithRequiredFields(ctx context.Context, m models.ThingWithRequiredFields) error {
	data, err := encodeThingWithRequiredFields(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#NAME": "name",
		},
		ConditionExpression: aws.String(
			"" +
				"" +
				"attribute_not_exists(#NAME)" +
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
			return db.ErrThingWithRequiredFieldsAlreadyExists{
				Name: *m.Name,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithRequiredFieldsTable) getSliceOfThingWithRequiredFields(ctx context.Context, ms []models.ThingWithRequiredFields) ([]models.ThingWithRequiredFields, error) {
	if len(ms) == 0 {
		return nil, nil
	}

	allKeys := make([]map[string]types.AttributeValue, len(ms))
	for i := range ms {
		key, err := attributevalue.MarshalMap(ddbThingWithRequiredFieldsPrimaryKey{
			Name: *ms[i].Name,
		})
		if err != nil {
			return nil, err
		}
		allKeys[i] = key
	}

	tname := t.TableName
	var items []models.ThingWithRequiredFields
	for len(allKeys) > 0 {
		chunkSize := len(allKeys)
		if chunkSize > maxDynamoDBBatchGetItems {
			chunkSize = maxDynamoDBBatchGetItems
		}
		requestKeys := allKeys[:chunkSize]
		allKeys = allKeys[chunkSize:]
		for attempt := 0; ; attempt++ {
			out, err := t.DynamoDBAPI.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
				RequestItems: map[string]types.KeysAndAttributes{
					tname: {Keys: requestKeys},
				},
			})
			if err != nil {
				return nil, fmt.Errorf("BatchGetItem: %v", err)
			}
			for _, item := range out.Responses[tname] {
				var m models.ThingWithRequiredFields
				if err := decodeThingWithRequiredFields(item, &m); err != nil {
					return nil, err
				}
				items = append(items, m)
			}
			if len(out.UnprocessedKeys[tname].Keys) == 0 {
				break
			}
			requestKeys = out.UnprocessedKeys[tname].Keys
			if attempt >= maxBatchUnprocessedRetries {
				return nil, db.ErrBatchUnprocessedItems{
					Operation:   "BatchGetItem",
					Table:       tname,
					Unprocessed: len(requestKeys),
					Attempts:    attempt + 1,
				}
			}
			// DynamoDB throttled some keys; back off (with jitter) before resubmitting
			// the unprocessed keys to avoid a thundering-herd retry storm.
			if err := waitBatchBackoff(ctx, attempt); err != nil {
				return nil, err
			}
		}
	}
	return items, nil
}

func (t ThingWithRequiredFieldsTable) getThingWithRequiredFields(ctx context.Context, name string) (*models.ThingWithRequiredFields, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithRequiredFieldsPrimaryKey{
		Name: name,
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
		return nil, db.ErrThingWithRequiredFieldsNotFound{
			Name: name,
		}
	}

	var m models.ThingWithRequiredFields
	if err := decodeThingWithRequiredFields(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithRequiredFieldsTable) scanThingWithRequiredFieldss(ctx context.Context, input db.ScanThingWithRequiredFieldssInput, fn func(m *models.ThingWithRequiredFields, lastThingWithRequiredFields bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMapWithOptions(input.StartingAfter, func(o *attributevalue.EncoderOptions) {
			o.TagKey = "json"
		})
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name": exclusiveStartKey["name"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithRequiredFieldss(out.Items)
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

func (t ThingWithRequiredFieldsTable) deleteThingWithRequiredFields(ctx context.Context, name string) error {

	key, err := attributevalue.MarshalMap(ddbThingWithRequiredFieldsPrimaryKey{
		Name: name,
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

// encodeThingWithRequiredFields encodes a ThingWithRequiredFields as a DynamoDB map of attribute values.
func encodeThingWithRequiredFields(m models.ThingWithRequiredFields) (map[string]types.AttributeValue, error) {
	// no composite attributes, marshal the model with the json tag
	val, err := attributevalue.MarshalMapWithOptions(m, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

// decodeThingWithRequiredFields translates a ThingWithRequiredFields stored in DynamoDB to a ThingWithRequiredFields struct.
func decodeThingWithRequiredFields(m map[string]types.AttributeValue, out *models.ThingWithRequiredFields) error {
	var ddbThingWithRequiredFields ddbThingWithRequiredFields
	if err := attributevalue.UnmarshalMapWithOptions(m, &ddbThingWithRequiredFields, func(o *attributevalue.DecoderOptions) {
		o.TagKey = "json"
	}); err != nil {
		return err
	}
	*out = ddbThingWithRequiredFields.ThingWithRequiredFields
	return nil
}

// decodeThingWithRequiredFieldss translates a list of ThingWithRequiredFieldss stored in DynamoDB to a slice of ThingWithRequiredFields structs.
func decodeThingWithRequiredFieldss(ms []map[string]types.AttributeValue) ([]models.ThingWithRequiredFields, error) {
	thingWithRequiredFieldss := make([]models.ThingWithRequiredFields, len(ms))
	for i, m := range ms {
		var thingWithRequiredFields models.ThingWithRequiredFields
		if err := decodeThingWithRequiredFields(m, &thingWithRequiredFields); err != nil {
			return nil, err
		}
		thingWithRequiredFieldss[i] = thingWithRequiredFields
	}
	return thingWithRequiredFieldss, nil
}
