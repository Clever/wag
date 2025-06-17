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

// ThingWithUnderscoresTable represents the user-configurable properties of the ThingWithUnderscores table.
type ThingWithUnderscoresTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithUnderscoresPrimaryKey represents the primary key of a ThingWithUnderscores in DynamoDB.
type ddbThingWithUnderscoresPrimaryKey struct {
	IDApp string `dynamodbav:"id_app"`
}

// ddbThingWithUnderscores represents a ThingWithUnderscores as stored in DynamoDB.
type ddbThingWithUnderscores struct {
	models.ThingWithUnderscores
}

func (t ThingWithUnderscoresTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id_app"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id_app"),
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

func (t ThingWithUnderscoresTable) saveThingWithUnderscores(ctx context.Context, m models.ThingWithUnderscores) error {
	data, err := encodeThingWithUnderscores(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
	})
	return err
}

func (t ThingWithUnderscoresTable) getThingWithUnderscores(ctx context.Context, iDApp string) (*models.ThingWithUnderscores, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithUnderscoresPrimaryKey{
		IDApp: iDApp,
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
		return nil, db.ErrThingWithUnderscoresNotFound{
			IDApp: iDApp,
		}
	}

	var m models.ThingWithUnderscores
	if err := decodeThingWithUnderscores(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithUnderscoresTable) deleteThingWithUnderscores(ctx context.Context, iDApp string) error {

	key, err := attributevalue.MarshalMap(ddbThingWithUnderscoresPrimaryKey{
		IDApp: iDApp,
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

// encodeThingWithUnderscores encodes a ThingWithUnderscores as a DynamoDB map of attribute values.
func encodeThingWithUnderscores(m models.ThingWithUnderscores) (map[string]types.AttributeValue, error) {
	// First marshal the model to get all fields
	rawVal, err := attributevalue.MarshalMap(m)
	if err != nil {
		return nil, err
	}

	// Create a new map with the correct field names from json tags
	val := make(map[string]types.AttributeValue)

	// Get the type of the ThingWithUnderscores struct
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

// decodeThingWithUnderscores translates a ThingWithUnderscores stored in DynamoDB to a ThingWithUnderscores struct.
func decodeThingWithUnderscores(m map[string]types.AttributeValue, out *models.ThingWithUnderscores) error {
	var ddbThingWithUnderscores ddbThingWithUnderscores
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithUnderscores); err != nil {
		return err
	}
	*out = ddbThingWithUnderscores.ThingWithUnderscores
	return nil
}

// decodeThingWithUnderscoress translates a list of ThingWithUnderscoress stored in DynamoDB to a slice of ThingWithUnderscores structs.
func decodeThingWithUnderscoress(ms []map[string]types.AttributeValue) ([]models.ThingWithUnderscores, error) {
	thingWithUnderscoress := make([]models.ThingWithUnderscores, len(ms))
	for i, m := range ms {
		var thingWithUnderscores models.ThingWithUnderscores
		if err := decodeThingWithUnderscores(m, &thingWithUnderscores); err != nil {
			return nil, err
		}
		thingWithUnderscoress[i] = thingWithUnderscores
	}
	return thingWithUnderscoress, nil
}
