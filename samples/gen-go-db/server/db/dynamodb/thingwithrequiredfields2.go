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

// ThingWithRequiredFields2Table represents the user-configurable properties of the ThingWithRequiredFields2 table.
type ThingWithRequiredFields2Table struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithRequiredFields2PrimaryKey represents the primary key of a ThingWithRequiredFields2 in DynamoDB.
type ddbThingWithRequiredFields2PrimaryKey struct {
	Name string `dynamodbav:"name"`
	ID   string `dynamodbav:"id"`
}

// ddbThingWithRequiredFields2 represents a ThingWithRequiredFields2 as stored in DynamoDB.
type ddbThingWithRequiredFields2 struct {
	models.ThingWithRequiredFields2
}

func (t ThingWithRequiredFields2Table) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeType("S"),
			},
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
			{
				AttributeName: aws.String("id"),
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

func (t ThingWithRequiredFields2Table) saveThingWithRequiredFields2(ctx context.Context, m models.ThingWithRequiredFields2) error {
	data, err := encodeThingWithRequiredFields2(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#NAME": "name",
			"#ID":   "id",
		},
		ConditionExpression: aws.String(
			"" +
				"" +
				"attribute_not_exists(#NAME)" +
				"" +
				" AND " +
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
			return db.ErrThingWithRequiredFields2AlreadyExists{
				Name: *m.Name,
				ID:   *m.ID,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithRequiredFields2Table) getThingWithRequiredFields2(ctx context.Context, name string, id string) (*models.ThingWithRequiredFields2, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithRequiredFields2PrimaryKey{
		Name: name,
		ID:   id,
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
		return nil, db.ErrThingWithRequiredFields2NotFound{
			Name: name,
			ID:   id,
		}
	}

	var m models.ThingWithRequiredFields2
	if err := decodeThingWithRequiredFields2(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithRequiredFields2Table) scanThingWithRequiredFields2s(ctx context.Context, input db.ScanThingWithRequiredFields2sInput, fn func(m *models.ThingWithRequiredFields2, lastThingWithRequiredFields2 bool) bool) error {
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
			"name": exclusiveStartKey["name"],
			"id":   exclusiveStartKey["id"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithRequiredFields2s(out.Items)
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

func (t ThingWithRequiredFields2Table) getThingWithRequiredFields2sByNameAndID(ctx context.Context, input db.GetThingWithRequiredFields2sByNameAndIDInput, fn func(m *models.ThingWithRequiredFields2, lastThingWithRequiredFields2 bool) bool) error {
	if input.IDStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.IDStartingAt or input.StartingAfter")
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
	if input.IDStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		queryInput.ExpressionAttributeNames["#ID"] = "id"
		queryInput.ExpressionAttributeValues[":id"] = &types.AttributeValueMemberS{
			Value: string(*input.IDStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #ID <= :id")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #ID >= :id")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: string(*input.StartingAfter.ID),
			},

			"name": &types.AttributeValueMemberS{
				Value: *input.StartingAfter.Name,
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithRequiredFields2s(queryOutput.Items)
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

func (t ThingWithRequiredFields2Table) deleteThingWithRequiredFields2(ctx context.Context, name string, id string) error {

	key, err := attributevalue.MarshalMap(ddbThingWithRequiredFields2PrimaryKey{
		Name: name,
		ID:   id,
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

// encodeThingWithRequiredFields2 encodes a ThingWithRequiredFields2 as a DynamoDB map of attribute values.
func encodeThingWithRequiredFields2(m models.ThingWithRequiredFields2) (map[string]types.AttributeValue, error) {
	// First marshal the model to get all fields
	rawVal, err := attributevalue.MarshalMap(m)
	if err != nil {
		return nil, err
	}

	// Create a new map with the correct field names from json tags
	val := make(map[string]types.AttributeValue)

	// Get the type of the ThingWithRequiredFields2 struct
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

// decodeThingWithRequiredFields2 translates a ThingWithRequiredFields2 stored in DynamoDB to a ThingWithRequiredFields2 struct.
func decodeThingWithRequiredFields2(m map[string]types.AttributeValue, out *models.ThingWithRequiredFields2) error {
	var ddbThingWithRequiredFields2 ddbThingWithRequiredFields2
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithRequiredFields2); err != nil {
		return err
	}
	*out = ddbThingWithRequiredFields2.ThingWithRequiredFields2
	return nil
}

// decodeThingWithRequiredFields2s translates a list of ThingWithRequiredFields2s stored in DynamoDB to a slice of ThingWithRequiredFields2 structs.
func decodeThingWithRequiredFields2s(ms []map[string]types.AttributeValue) ([]models.ThingWithRequiredFields2, error) {
	thingWithRequiredFields2s := make([]models.ThingWithRequiredFields2, len(ms))
	for i, m := range ms {
		var thingWithRequiredFields2 models.ThingWithRequiredFields2
		if err := decodeThingWithRequiredFields2(m, &thingWithRequiredFields2); err != nil {
			return nil, err
		}
		thingWithRequiredFields2s[i] = thingWithRequiredFields2
	}
	return thingWithRequiredFields2s, nil
}
