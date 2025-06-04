package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db-custom-path/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db-custom-path/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithDateRangeKeyTable represents the user-configurable properties of the ThingWithDateRangeKey table.
type ThingWithDateRangeKeyTable struct {
	DynamoDBAPI        *dynamodb.Client
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
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("date"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("date"),
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

func (t ThingWithDateRangeKeyTable) saveThingWithDateRangeKey(ctx context.Context, m models.ThingWithDateRangeKey) error {
	data, err := encodeThingWithDateRangeKey(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#ID":   "id",
			"#DATE": "date",
		},
		ConditionExpression: aws.String("attribute_not_exists(#ID) AND attribute_not_exists(#DATE)"),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		var conditionalCheckFailedErr *types.ConditionalCheckFailedException
		if errors.As(err, &resourceNotFoundErr) {
			return fmt.Errorf("table or index not found: %s", t.TableName)
		}
		if errors.As(err, &conditionalCheckFailedErr) {
			return db.ErrThingWithDateRangeKeyAlreadyExists{
				ID:   m.ID,
				Date: m.Date,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithDateRangeKeyTable) getThingWithDateRangeKey(ctx context.Context, id string, date strfmt.Date) (*models.ThingWithDateRangeKey, error) {
	// swad-get-7
	key, err := attributevalue.MarshalMap(ddbThingWithDateRangeKeyPrimaryKey{
		ID:   id,
		Date: date,
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
	// swad-scan-1
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"id":   exclusiveStartKey["id"],
			"date": exclusiveStartKey["date"],
		}
	}
	totalRecordsProcessed := int32(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithDateRangeKeys(out.Items)
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

func (t ThingWithDateRangeKeyTable) getThingWithDateRangeKeysByIDAndDate(ctx context.Context, input db.GetThingWithDateRangeKeysByIDAndDateInput, fn func(m *models.ThingWithDateRangeKey, lastThingWithDateRangeKey bool) bool) error {
	// swad-get-2
	if input.DateStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateStartingAt or input.StartingAfter")
	}
	if input.ID == "" {
		return fmt.Errorf("Hash key input.ID cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		ExpressionAttributeNames: map[string]string{
			"#ID": "id",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{
				Value: input.ID,
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
		// swad-get-21
		queryInput.ExpressionAttributeNames["#DATE"] = "date"
		queryInput.ExpressionAttributeValues[":date"] = &types.AttributeValueMemberS{
			Value: dateToDynamoTimeString(*input.DateStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#ID = :id AND #DATE <= :date")
		} else {
			queryInput.KeyConditionExpression = aws.String("#ID = :id AND #DATE >= :date")
		}
	}
	// swad-get-22
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"date": &types.AttributeValueMemberS{
				Value: dateToDynamoTimeString(input.StartingAfter.Date),
			},

			// swad-get-223
			"id": &types.AttributeValueMemberS{
				Value: input.StartingAfter.ID,
			},
		}
	}

	totalRecordsProcessed := int32(0)
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

func (t ThingWithDateRangeKeyTable) deleteThingWithDateRangeKey(ctx context.Context, id string, date strfmt.Date) error {

	key, err := attributevalue.MarshalMap(ddbThingWithDateRangeKeyPrimaryKey{
		ID:   id,
		Date: date,
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

// encodeThingWithDateRangeKey encodes a ThingWithDateRangeKey as a DynamoDB map of attribute values.
func encodeThingWithDateRangeKey(m models.ThingWithDateRangeKey) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(ddbThingWithDateRangeKey{
		ThingWithDateRangeKey: m,
	})
}

// decodeThingWithDateRangeKey translates a ThingWithDateRangeKey stored in DynamoDB to a ThingWithDateRangeKey struct.
func decodeThingWithDateRangeKey(m map[string]types.AttributeValue, out *models.ThingWithDateRangeKey) error {
	// swad-decode-1
	var ddbThingWithDateRangeKey ddbThingWithDateRangeKey
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithDateRangeKey); err != nil {
		return err
	}
	*out = ddbThingWithDateRangeKey.ThingWithDateRangeKey
	return nil
}

// decodeThingWithDateRangeKeys translates a list of ThingWithDateRangeKeys stored in DynamoDB to a slice of ThingWithDateRangeKey structs.
func decodeThingWithDateRangeKeys(ms []map[string]types.AttributeValue) ([]models.ThingWithDateRangeKey, error) {
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
