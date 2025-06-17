package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Clever/wag/samples/gen-go-db-custom-path/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db-custom-path/db"
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

// ThingWithDatetimeGSITable represents the user-configurable properties of the ThingWithDatetimeGSI table.
type ThingWithDatetimeGSITable struct {
	DynamoDBAPI        *dynamodb.Client
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
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("datetime"),
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
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("byDateTime"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("datetime"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("id"),
						KeyType:       types.KeyTypeRange,
					},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
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

func (t ThingWithDatetimeGSITable) saveThingWithDatetimeGSI(ctx context.Context, m models.ThingWithDatetimeGSI) error {
	data, err := encodeThingWithDatetimeGSI(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#ID": "id",
		},
		ConditionExpression: aws.String(
			"" +
				"" +
				"attribute_not_exists(#ID)" +
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
			return db.ErrThingWithDatetimeGSIAlreadyExists{
				ID: m.ID,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithDatetimeGSITable) getThingWithDatetimeGSI(ctx context.Context, id string) (*models.ThingWithDatetimeGSI, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithDatetimeGSIPrimaryKey{
		ID: id,
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
			"id": exclusiveStartKey["id"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithDatetimeGSIs(out.Items)
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

func (t ThingWithDatetimeGSITable) deleteThingWithDatetimeGSI(ctx context.Context, id string) error {

	key, err := attributevalue.MarshalMap(ddbThingWithDatetimeGSIPrimaryKey{
		ID: id,
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
		ExpressionAttributeNames: map[string]string{
			"#DATETIME": "datetime",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":datetime": &types.AttributeValueMemberS{
				Value: datetimeToDynamoTimeString(input.Datetime),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.IDStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#DATETIME = :datetime")
	} else {
		queryInput.ExpressionAttributeNames["#ID"] = "id"
		queryInput.ExpressionAttributeValues[":id"] = &types.AttributeValueMemberS{
			Value: string(*input.IDStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#DATETIME = :datetime AND #ID <= :id")
		} else {
			queryInput.KeyConditionExpression = aws.String("#DATETIME = :datetime AND #ID >= :id")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: string(input.StartingAfter.ID),
			},
			"datetime": &types.AttributeValueMemberS{
				Value: datetimeToDynamoTimeString(input.StartingAfter.Datetime),
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
func (t ThingWithDatetimeGSITable) scanThingWithDatetimeGSIsByDatetimeAndID(ctx context.Context, input db.ScanThingWithDatetimeGSIsByDatetimeAndIDInput, fn func(m *models.ThingWithDatetimeGSI, lastThingWithDatetimeGSI bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	scanInput.IndexName = aws.String("byDateTime")
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"id":       exclusiveStartKey["id"],
			"datetime": exclusiveStartKey["datetime"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithDatetimeGSIs(out.Items)
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

// encodeThingWithDatetimeGSI encodes a ThingWithDatetimeGSI as a DynamoDB map of attribute values.
func encodeThingWithDatetimeGSI(m models.ThingWithDatetimeGSI) (map[string]types.AttributeValue, error) {
	// First marshal the model to get all fields
	rawVal, err := attributevalue.MarshalMap(m)
	if err != nil {
		return nil, err
	}

	// Create a new map with the correct field names from json tags
	val := make(map[string]types.AttributeValue)

	// Get the type of the ThingWithDatetimeGSI struct
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

// decodeThingWithDatetimeGSI translates a ThingWithDatetimeGSI stored in DynamoDB to a ThingWithDatetimeGSI struct.
func decodeThingWithDatetimeGSI(m map[string]types.AttributeValue, out *models.ThingWithDatetimeGSI) error {
	var ddbThingWithDatetimeGSI ddbThingWithDatetimeGSI
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithDatetimeGSI); err != nil {
		return err
	}
	*out = ddbThingWithDatetimeGSI.ThingWithDatetimeGSI
	return nil
}

// decodeThingWithDatetimeGSIs translates a list of ThingWithDatetimeGSIs stored in DynamoDB to a slice of ThingWithDatetimeGSI structs.
func decodeThingWithDatetimeGSIs(ms []map[string]types.AttributeValue) ([]models.ThingWithDatetimeGSI, error) {
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
