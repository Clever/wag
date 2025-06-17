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
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
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
	// First marshal the model to get all fields
	rawVal, err := attributevalue.MarshalMap(m)
	if err != nil {
		return nil, err
	}

	// Create a new map with the correct field names from json tags
	val := make(map[string]types.AttributeValue)

	// Get the type of the ThingWithRequiredFields struct
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

// decodeThingWithRequiredFields translates a ThingWithRequiredFields stored in DynamoDB to a ThingWithRequiredFields struct.
func decodeThingWithRequiredFields(m map[string]types.AttributeValue, out *models.ThingWithRequiredFields) error {
	var ddbThingWithRequiredFields ddbThingWithRequiredFields
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithRequiredFields); err != nil {
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
