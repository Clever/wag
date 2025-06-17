package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

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
var _ = reflect.TypeOf(int(0))
var _ = strings.Split("", "")

// ThingAllowingBatchWritesWithCompositeAttributesTable represents the user-configurable properties of the ThingAllowingBatchWritesWithCompositeAttributes table.
type ThingAllowingBatchWritesWithCompositeAttributesTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingAllowingBatchWritesWithCompositeAttributesPrimaryKey represents the primary key of a ThingAllowingBatchWritesWithCompositeAttributes in DynamoDB.
type ddbThingAllowingBatchWritesWithCompositeAttributesPrimaryKey struct {
	NameID string          `dynamodbav:"name_id"`
	Date   strfmt.DateTime `dynamodbav:"date"`
}

// ddbThingAllowingBatchWritesWithCompositeAttributes represents a ThingAllowingBatchWritesWithCompositeAttributes as stored in DynamoDB.
type ddbThingAllowingBatchWritesWithCompositeAttributes struct {
	models.ThingAllowingBatchWritesWithCompositeAttributes
}

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("date"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("name_id"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name_id"),
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

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) saveThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, m models.ThingAllowingBatchWritesWithCompositeAttributes) error {
	data, err := encodeThingAllowingBatchWritesWithCompositeAttributes(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#NAME_ID": "name_id",
			"#DATE":    "date",
		},
		ConditionExpression: aws.String(
			"" +
				"" +
				"attribute_not_exists(#NAME_ID)" +
				"" +
				" AND " +
				"attribute_not_exists(#DATE)" +
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
			return db.ErrThingAllowingBatchWritesWithCompositeAttributesAlreadyExists{
				NameID: fmt.Sprintf("%s@%s", *m.Name, *m.ID),
				Date:   *m.Date,
			}
		}
		return err
	}
	return nil
}
func (t ThingAllowingBatchWritesWithCompositeAttributesTable) saveArrayOfThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, ms []models.ThingAllowingBatchWritesWithCompositeAttributes) error {
	if len(ms) > maxDynamoDBBatchItems {
		return fmt.Errorf("saveArrayOfThingAllowingBatchWritesWithCompositeAttributes received %d items to save, which is greater than the maximum of %d", len(ms), maxDynamoDBBatchItems)
	}

	if len(ms) == 0 {
		return nil
	}

	batch := make([]types.WriteRequest, len(ms))
	for i := range ms {
		data, err := encodeThingAllowingBatchWritesWithCompositeAttributes(ms[i])
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

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) deleteArrayOfThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, ms []models.ThingAllowingBatchWritesWithCompositeAttributes) error {
	if len(ms) > maxDynamoDBBatchItems {
		return fmt.Errorf("deleteArrayOfThingAllowingBatchWritesWithCompositeAttributes received %d items to delete, which is greater than the maximum of %d", len(ms), maxDynamoDBBatchItems)
	}

	if len(ms) == 0 {
		return nil
	}

	batch := make([]types.WriteRequest, len(ms))
	for i := range ms {
		key, err := attributevalue.MarshalMap(ddbThingAllowingBatchWritesWithCompositeAttributesPrimaryKey{
			NameID: fmt.Sprintf("%s@%s", *ms[i].Name, *ms[i].ID),
			Date:   *ms[i].Date,
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

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) getThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, name string, id string, date strfmt.DateTime) (*models.ThingAllowingBatchWritesWithCompositeAttributes, error) {
	key, err := attributevalue.MarshalMap(ddbThingAllowingBatchWritesWithCompositeAttributesPrimaryKey{
		NameID: fmt.Sprintf("%s@%s", name, id),
		Date:   date,
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
		return nil, db.ErrThingAllowingBatchWritesWithCompositeAttributesNotFound{
			Name: name,
			ID:   id,
			Date: date,
		}
	}

	var m models.ThingAllowingBatchWritesWithCompositeAttributes
	if err := decodeThingAllowingBatchWritesWithCompositeAttributes(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) scanThingAllowingBatchWritesWithCompositeAttributess(ctx context.Context, input db.ScanThingAllowingBatchWritesWithCompositeAttributessInput, fn func(m *models.ThingAllowingBatchWritesWithCompositeAttributes, lastThingAllowingBatchWritesWithCompositeAttributes bool) bool) error {
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
			"name_id": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s@%s", *input.StartingAfter.Name, *input.StartingAfter.ID),
			},
			"date": exclusiveStartKey["date"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingAllowingBatchWritesWithCompositeAttributess(out.Items)
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

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) getThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDate(ctx context.Context, input db.GetThingAllowingBatchWritesWithCompositeAttributessByNameIDAndDateInput, fn func(m *models.ThingAllowingBatchWritesWithCompositeAttributes, lastThingAllowingBatchWritesWithCompositeAttributes bool) bool) error {
	if input.DateStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	if input.ID == "" {
		return fmt.Errorf("Hash key input.ID cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		ExpressionAttributeNames: map[string]string{
			"#NAME_ID": "name_id",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nameId": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s@%s", input.Name, input.ID),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME_ID = :nameId")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = "date"
		queryInput.ExpressionAttributeValues[":date"] = &types.AttributeValueMemberS{
			Value: datetimeToDynamoTimeString(*input.DateStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME_ID = :nameId AND #DATE <= :date")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME_ID = :nameId AND #DATE >= :date")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"date": &types.AttributeValueMemberS{
				Value: datetimePtrToDynamoTimeString(input.StartingAfter.Date),
			},

			"name_id": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s@%s", *input.StartingAfter.Name, *input.StartingAfter.ID),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingAllowingBatchWritesWithCompositeAttributess(queryOutput.Items)
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

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) deleteThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, name string, id string, date strfmt.DateTime) error {

	key, err := attributevalue.MarshalMap(ddbThingAllowingBatchWritesWithCompositeAttributesPrimaryKey{
		NameID: fmt.Sprintf("%s@%s", name, id),
		Date:   date,
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

// encodeThingAllowingBatchWritesWithCompositeAttributes encodes a ThingAllowingBatchWritesWithCompositeAttributes as a DynamoDB map of attribute values.
func encodeThingAllowingBatchWritesWithCompositeAttributes(m models.ThingAllowingBatchWritesWithCompositeAttributes) (map[string]types.AttributeValue, error) {
	// First marshal the model to get all fields
	rawVal, err := attributevalue.MarshalMap(m)
	if err != nil {
		return nil, err
	}

	// Create a new map with the correct field names from json tags
	val := make(map[string]types.AttributeValue)

	// Get the type of the ThingAllowingBatchWritesWithCompositeAttributes struct
	t := reflect.TypeOf(m)

	// Create a map of struct field names to their json tags and types
	fieldMap := make(map[string]struct {
		jsonName string
		isMap    bool
	})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			// Handle omitempty in the tag
			jsonTag = strings.Split(jsonTag, ",")[0]
			fieldMap[field.Name] = struct {
				jsonName string
				isMap    bool
			}{
				jsonName: jsonTag,
				isMap:    field.Type.Kind() == reflect.Map || field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Map,
			}
		}
	}

	for k, v := range rawVal {
		// Skip null values
		if _, ok := v.(*types.AttributeValueMemberNULL); ok {
			continue
		}

		// Get the field info from the map
		if fieldInfo, ok := fieldMap[k]; ok {
			// Handle map fields
			if fieldInfo.isMap {
				if memberM, ok := v.(*types.AttributeValueMemberM); ok {
					// Create a new map for the nested structure
					nestedVal := make(map[string]types.AttributeValue)
					for mk, mv := range memberM.Value {
						// Skip null values in nested map
						if _, ok := mv.(*types.AttributeValueMemberNULL); ok {
							continue
						}
						nestedVal[mk] = mv
					}
					val[fieldInfo.jsonName] = &types.AttributeValueMemberM{Value: nestedVal}
				}
				continue
			}

			val[fieldInfo.jsonName] = v
		}
	}

	return val, nil
}

// decodeThingAllowingBatchWritesWithCompositeAttributes translates a ThingAllowingBatchWritesWithCompositeAttributes stored in DynamoDB to a ThingAllowingBatchWritesWithCompositeAttributes struct.
func decodeThingAllowingBatchWritesWithCompositeAttributes(m map[string]types.AttributeValue, out *models.ThingAllowingBatchWritesWithCompositeAttributes) error {
	var ddbThingAllowingBatchWritesWithCompositeAttributes ddbThingAllowingBatchWritesWithCompositeAttributes
	if err := attributevalue.UnmarshalMap(m, &ddbThingAllowingBatchWritesWithCompositeAttributes); err != nil {
		return err
	}
	*out = ddbThingAllowingBatchWritesWithCompositeAttributes.ThingAllowingBatchWritesWithCompositeAttributes
	return nil
}

// decodeThingAllowingBatchWritesWithCompositeAttributess translates a list of ThingAllowingBatchWritesWithCompositeAttributess stored in DynamoDB to a slice of ThingAllowingBatchWritesWithCompositeAttributes structs.
func decodeThingAllowingBatchWritesWithCompositeAttributess(ms []map[string]types.AttributeValue) ([]models.ThingAllowingBatchWritesWithCompositeAttributes, error) {
	thingAllowingBatchWritesWithCompositeAttributess := make([]models.ThingAllowingBatchWritesWithCompositeAttributes, len(ms))
	for i, m := range ms {
		var thingAllowingBatchWritesWithCompositeAttributes models.ThingAllowingBatchWritesWithCompositeAttributes
		if err := decodeThingAllowingBatchWritesWithCompositeAttributes(m, &thingAllowingBatchWritesWithCompositeAttributes); err != nil {
			return nil, err
		}
		thingAllowingBatchWritesWithCompositeAttributess[i] = thingAllowingBatchWritesWithCompositeAttributes
	}
	return thingAllowingBatchWritesWithCompositeAttributess, nil
}
