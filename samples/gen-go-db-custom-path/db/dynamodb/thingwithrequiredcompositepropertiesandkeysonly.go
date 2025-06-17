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

// ThingWithRequiredCompositePropertiesAndKeysOnlyTable represents the user-configurable properties of the ThingWithRequiredCompositePropertiesAndKeysOnly table.
type ThingWithRequiredCompositePropertiesAndKeysOnlyTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithRequiredCompositePropertiesAndKeysOnlyPrimaryKey represents the primary key of a ThingWithRequiredCompositePropertiesAndKeysOnly in DynamoDB.
type ddbThingWithRequiredCompositePropertiesAndKeysOnlyPrimaryKey struct {
	PropertyThree string `dynamodbav:"propertyThree"`
}

// ddbThingWithRequiredCompositePropertiesAndKeysOnlyGSIPropertyOneAndTwoPropertyThree represents the propertyOneAndTwo_PropertyThree GSI.
type ddbThingWithRequiredCompositePropertiesAndKeysOnlyGSIPropertyOneAndTwoPropertyThree struct {
	PropertyOneAndTwo string `dynamodbav:"propertyOneAndTwo"`
	PropertyThree     string `dynamodbav:"propertyThree"`
}

// ddbThingWithRequiredCompositePropertiesAndKeysOnly represents a ThingWithRequiredCompositePropertiesAndKeysOnly as stored in DynamoDB.
type ddbThingWithRequiredCompositePropertiesAndKeysOnly struct {
	models.ThingWithRequiredCompositePropertiesAndKeysOnly
}

func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("propertyOneAndTwo"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("propertyThree"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("propertyThree"),
				KeyType:       types.KeyTypeHash,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("propertyOneAndTwo_PropertyThree"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("KEYS_ONLY"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("propertyOneAndTwo"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("propertyThree"),
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

func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) saveThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, m models.ThingWithRequiredCompositePropertiesAndKeysOnly) error {
	data, err := encodeThingWithRequiredCompositePropertiesAndKeysOnly(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
	})
	return err
}

func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) getThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, propertyThree string) (*models.ThingWithRequiredCompositePropertiesAndKeysOnly, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithRequiredCompositePropertiesAndKeysOnlyPrimaryKey{
		PropertyThree: propertyThree,
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
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrThingWithRequiredCompositePropertiesAndKeysOnlyNotFound{
			PropertyThree: propertyThree,
		}
	}

	var m models.ThingWithRequiredCompositePropertiesAndKeysOnly
	if err := decodeThingWithRequiredCompositePropertiesAndKeysOnly(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) scanThingWithRequiredCompositePropertiesAndKeysOnlys(ctx context.Context, input db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysInput, fn func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, lastThingWithRequiredCompositePropertiesAndKeysOnly bool) bool) error {
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
			"propertyThree": exclusiveStartKey["propertyThree"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithRequiredCompositePropertiesAndKeysOnlys(out.Items)
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

func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) deleteThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, propertyThree string) error {

	key, err := attributevalue.MarshalMap(ddbThingWithRequiredCompositePropertiesAndKeysOnlyPrimaryKey{
		PropertyThree: propertyThree,
	})
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(t.TableName),
	})
	if err != nil {
		return err
	}

	return nil
}

func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) getThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx context.Context, input db.GetThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput, fn func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, lastThingWithRequiredCompositePropertiesAndKeysOnly bool) bool) error {
	if input.PropertyThreeStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.PropertyThreeStartingAt or input.StartingAfter")
	}
	if input.PropertyOne == "" {
		return fmt.Errorf("Hash key input.PropertyOne cannot be empty")
	}
	if input.PropertyTwo == "" {
		return fmt.Errorf("Hash key input.PropertyTwo cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("propertyOneAndTwo_PropertyThree"),
		ExpressionAttributeNames: map[string]string{
			"#PROPERTYONEANDTWO": "propertyOneAndTwo",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":propertyOneAndTwo": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s_%s", input.PropertyOne, input.PropertyTwo),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.PropertyThreeStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#PROPERTYONEANDTWO = :propertyOneAndTwo")
	} else {
		queryInput.ExpressionAttributeNames["#PROPERTYTHREE"] = "propertyThree"
		queryInput.ExpressionAttributeValues[":propertyThree"] = &types.AttributeValueMemberS{
			Value: string(*input.PropertyThreeStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#PROPERTYONEANDTWO = :propertyOneAndTwo AND #PROPERTYTHREE <= :propertyThree")
		} else {
			queryInput.KeyConditionExpression = aws.String("#PROPERTYONEANDTWO = :propertyOneAndTwo AND #PROPERTYTHREE >= :propertyThree")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"propertyThree": &types.AttributeValueMemberS{
				Value: string(*input.StartingAfter.PropertyThree),
			},
			"propertyOneAndTwo": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s_%s", *input.StartingAfter.PropertyOne, *input.StartingAfter.PropertyTwo),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithRequiredCompositePropertiesAndKeysOnlys(queryOutput.Items)
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
func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) scanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx context.Context, input db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput, fn func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, lastThingWithRequiredCompositePropertiesAndKeysOnly bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	scanInput.IndexName = aws.String("propertyOneAndTwo_PropertyThree")
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"propertyThree": exclusiveStartKey["propertyThree"],
			"propertyOneAndTwo": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s_%s", *input.StartingAfter.PropertyOne, *input.StartingAfter.PropertyTwo),
			},
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithRequiredCompositePropertiesAndKeysOnlys(out.Items)
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

// encodeThingWithRequiredCompositePropertiesAndKeysOnly encodes a ThingWithRequiredCompositePropertiesAndKeysOnly as a DynamoDB map of attribute values.
func encodeThingWithRequiredCompositePropertiesAndKeysOnly(m models.ThingWithRequiredCompositePropertiesAndKeysOnly) (map[string]types.AttributeValue, error) {
	// First marshal the model to get all fields
	rawVal, err := attributevalue.MarshalMap(m)
	if err != nil {
		return nil, err
	}

	// Create a new map with the correct field names from json tags
	val := make(map[string]types.AttributeValue)

	// Get the type of the ThingWithRequiredCompositePropertiesAndKeysOnly struct
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

// decodeThingWithRequiredCompositePropertiesAndKeysOnly translates a ThingWithRequiredCompositePropertiesAndKeysOnly stored in DynamoDB to a ThingWithRequiredCompositePropertiesAndKeysOnly struct.
func decodeThingWithRequiredCompositePropertiesAndKeysOnly(m map[string]types.AttributeValue, out *models.ThingWithRequiredCompositePropertiesAndKeysOnly) error {
	var ddbThingWithRequiredCompositePropertiesAndKeysOnly ddbThingWithRequiredCompositePropertiesAndKeysOnly
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithRequiredCompositePropertiesAndKeysOnly); err != nil {
		return err
	}
	*out = ddbThingWithRequiredCompositePropertiesAndKeysOnly.ThingWithRequiredCompositePropertiesAndKeysOnly
	// parse composite attributes from projected secondary indexes and fill
	// in model properties
	if v, ok := m["propertyOneAndTwo"]; ok {
		if s, ok := v.(*types.AttributeValueMemberS); ok {
			parts := strings.Split(s.Value, "_")
			if len(parts) != 2 {
				return fmt.Errorf("expected 2 parts: '%s'", s.Value)
			}
			out.PropertyOne = &parts[0]
			out.PropertyTwo = &parts[1]
		}
	}
	return nil
}

// decodeThingWithRequiredCompositePropertiesAndKeysOnlys translates a list of ThingWithRequiredCompositePropertiesAndKeysOnlys stored in DynamoDB to a slice of ThingWithRequiredCompositePropertiesAndKeysOnly structs.
func decodeThingWithRequiredCompositePropertiesAndKeysOnlys(ms []map[string]types.AttributeValue) ([]models.ThingWithRequiredCompositePropertiesAndKeysOnly, error) {
	thingWithRequiredCompositePropertiesAndKeysOnlys := make([]models.ThingWithRequiredCompositePropertiesAndKeysOnly, len(ms))
	for i, m := range ms {
		var thingWithRequiredCompositePropertiesAndKeysOnly models.ThingWithRequiredCompositePropertiesAndKeysOnly
		if err := decodeThingWithRequiredCompositePropertiesAndKeysOnly(m, &thingWithRequiredCompositePropertiesAndKeysOnly); err != nil {
			return nil, err
		}
		thingWithRequiredCompositePropertiesAndKeysOnlys[i] = thingWithRequiredCompositePropertiesAndKeysOnly
	}
	return thingWithRequiredCompositePropertiesAndKeysOnlys, nil
}
