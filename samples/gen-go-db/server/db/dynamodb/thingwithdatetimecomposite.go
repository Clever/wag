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

// ThingWithDateTimeCompositeTable represents the user-configurable properties of the ThingWithDateTimeComposite table.
type ThingWithDateTimeCompositeTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithDateTimeCompositePrimaryKey represents the primary key of a ThingWithDateTimeComposite in DynamoDB.
type ddbThingWithDateTimeCompositePrimaryKey struct {
	TypeID          string `dynamodbav:"typeID"`
	CreatedResource string `dynamodbav:"createdResource"`
}

// ddbThingWithDateTimeComposite represents a ThingWithDateTimeComposite as stored in DynamoDB.
type ddbThingWithDateTimeComposite struct {
	models.ThingWithDateTimeComposite
}

func (t ThingWithDateTimeCompositeTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("createdResource"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("typeID"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("typeID"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("createdResource"),
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

func (t ThingWithDateTimeCompositeTable) saveThingWithDateTimeComposite(ctx context.Context, m models.ThingWithDateTimeComposite) error {
	data, err := encodeThingWithDateTimeComposite(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
	})
	return err
}

func (t ThingWithDateTimeCompositeTable) getThingWithDateTimeComposite(ctx context.Context, typeVar string, id string, created strfmt.DateTime, resource string) (*models.ThingWithDateTimeComposite, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithDateTimeCompositePrimaryKey{
		TypeID:          fmt.Sprintf("%s|%s", typeVar, id),
		CreatedResource: fmt.Sprintf("%s|%s", created, resource),
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
		return nil, db.ErrThingWithDateTimeCompositeNotFound{
			Type:     typeVar,
			ID:       id,
			Created:  created,
			Resource: resource,
		}
	}

	var m models.ThingWithDateTimeComposite
	if err := decodeThingWithDateTimeComposite(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithDateTimeCompositeTable) scanThingWithDateTimeComposites(ctx context.Context, input db.ScanThingWithDateTimeCompositesInput, fn func(m *models.ThingWithDateTimeComposite, lastThingWithDateTimeComposite bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAfter != nil {
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"typeID": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s|%s", input.StartingAfter.Type, input.StartingAfter.ID),
			},
			"createdResource": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s|%s", input.StartingAfter.Created, input.StartingAfter.Resource),
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

		items, err := decodeThingWithDateTimeComposites(out.Items)
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

func (t ThingWithDateTimeCompositeTable) getThingWithDateTimeCompositesByTypeIDAndCreatedResource(ctx context.Context, input db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput, fn func(m *models.ThingWithDateTimeComposite, lastThingWithDateTimeComposite bool) bool) error {
	if input.StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of StartingAt or StartingAfter")
	}
	if input.Type == "" {
		return fmt.Errorf("Hash key input.Type cannot be empty")
	}
	if input.ID == "" {
		return fmt.Errorf("Hash key input.ID cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		ExpressionAttributeNames: map[string]string{
			"#TYPEID": "typeID",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":typeId": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s|%s", input.Type, input.ID),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#TYPEID = :typeId")
	} else {
		queryInput.ExpressionAttributeNames["#CREATEDRESOURCE"] = "createdResource"
		queryInput.ExpressionAttributeValues[":createdResource"] = &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%s|%s", input.StartingAt.Created, input.StartingAt.Resource),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#TYPEID = :typeId AND #CREATEDRESOURCE <= :createdResource")
		} else {
			queryInput.KeyConditionExpression = aws.String("#TYPEID = :typeId AND #CREATEDRESOURCE >= :createdResource")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"createdResource": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s|%s", input.StartingAfter.Created, input.StartingAfter.Resource),
			},

			"typeID": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s|%s", input.StartingAfter.Type, input.StartingAfter.ID),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithDateTimeComposites(queryOutput.Items)
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

func (t ThingWithDateTimeCompositeTable) deleteThingWithDateTimeComposite(ctx context.Context, typeVar string, id string, created strfmt.DateTime, resource string) error {

	key, err := attributevalue.MarshalMap(ddbThingWithDateTimeCompositePrimaryKey{
		TypeID:          fmt.Sprintf("%s|%s", typeVar, id),
		CreatedResource: fmt.Sprintf("%s|%s", created, resource),
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

// encodeThingWithDateTimeComposite encodes a ThingWithDateTimeComposite as a DynamoDB map of attribute values.
func encodeThingWithDateTimeComposite(m models.ThingWithDateTimeComposite) (map[string]types.AttributeValue, error) {
	// First marshal the model to get all fields
	rawVal, err := attributevalue.MarshalMap(m)
	if err != nil {
		return nil, err
	}

	// Create a new map with the correct field names from json tags
	val := make(map[string]types.AttributeValue)

	// Get the type of the ThingWithDateTimeComposite struct
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

// decodeThingWithDateTimeComposite translates a ThingWithDateTimeComposite stored in DynamoDB to a ThingWithDateTimeComposite struct.
func decodeThingWithDateTimeComposite(m map[string]types.AttributeValue, out *models.ThingWithDateTimeComposite) error {
	var ddbThingWithDateTimeComposite ddbThingWithDateTimeComposite
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithDateTimeComposite); err != nil {
		return err
	}
	*out = ddbThingWithDateTimeComposite.ThingWithDateTimeComposite
	return nil
}

// decodeThingWithDateTimeComposites translates a list of ThingWithDateTimeComposites stored in DynamoDB to a slice of ThingWithDateTimeComposite structs.
func decodeThingWithDateTimeComposites(ms []map[string]types.AttributeValue) ([]models.ThingWithDateTimeComposite, error) {
	thingWithDateTimeComposites := make([]models.ThingWithDateTimeComposite, len(ms))
	for i, m := range ms {
		var thingWithDateTimeComposite models.ThingWithDateTimeComposite
		if err := decodeThingWithDateTimeComposite(m, &thingWithDateTimeComposite); err != nil {
			return nil, err
		}
		thingWithDateTimeComposites[i] = thingWithDateTimeComposite
	}
	return thingWithDateTimeComposites, nil
}
