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

// EventTable represents the user-configurable properties of the Event table.
type EventTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbEventPrimaryKey represents the primary key of a Event in DynamoDB.
type ddbEventPrimaryKey struct {
	Pk string `dynamodbav:"pk"`
	Sk string `dynamodbav:"sk"`
}

// ddbEventGSIBySK represents the bySK GSI.
type ddbEventGSIBySK struct {
	Sk   string `dynamodbav:"sk"`
	Data []byte `dynamodbav:"data"`
}

// ddbEvent represents a Event as stored in DynamoDB.
type ddbEvent struct {
	models.Event
}

func (t EventTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("data"),
				AttributeType: types.ScalarAttributeType("B"),
			},
			{
				AttributeName: aws.String("pk"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("sk"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("sk"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("bySK"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("sk"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("data"),
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

func (t EventTable) saveEvent(ctx context.Context, m models.Event) error {
	data, err := encodeEvent(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
	})
	return err
}

func (t EventTable) getEvent(ctx context.Context, pk string, sk string) (*models.Event, error) {
	key, err := attributevalue.MarshalMap(ddbEventPrimaryKey{
		Pk: pk,
		Sk: sk,
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
		return nil, db.ErrEventNotFound{
			Pk: pk,
			Sk: sk,
		}
	}

	var m models.Event
	if err := decodeEvent(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t EventTable) scanEvents(ctx context.Context, input db.ScanEventsInput, fn func(m *models.Event, lastEvent bool) bool) error {
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
			"pk": exclusiveStartKey["pk"],
			"sk": exclusiveStartKey["sk"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeEvents(out.Items)
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

func (t EventTable) getEventsByPkAndSkParseFilters(queryInput *dynamodb.QueryInput, input db.GetEventsByPkAndSkInput) {
	for _, filterValue := range input.FilterValues {
		switch filterValue.AttributeName {
		case db.EventData:
			queryInput.ExpressionAttributeNames["#DATA"] = string(db.EventData)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.EventData), i)] = &types.AttributeValueMemberB{
					Value: attributeValue.([]byte),
				}
			}
		}
	}
}

func (t EventTable) getEventsByPkAndSk(ctx context.Context, input db.GetEventsByPkAndSkInput, fn func(m *models.Event, lastEvent bool) bool) error {
	if input.SkStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.SkStartingAt or input.StartingAfter")
	}
	if input.Pk == "" {
		return fmt.Errorf("Hash key input.Pk cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		ExpressionAttributeNames: map[string]string{
			"#PK": "pk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{
				Value: input.Pk,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.SkStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#PK = :pk")
	} else {
		queryInput.ExpressionAttributeNames["#SK"] = "sk"
		queryInput.ExpressionAttributeValues[":sk"] = &types.AttributeValueMemberS{
			Value: string(*input.SkStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#PK = :pk AND #SK <= :sk")
		} else {
			queryInput.KeyConditionExpression = aws.String("#PK = :pk AND #SK >= :sk")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"sk": &types.AttributeValueMemberS{
				Value: string(input.StartingAfter.Sk),
			},

			"pk": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Pk,
			},
		}
	}
	if len(input.FilterValues) > 0 && input.FilterExpression != "" {
		t.getEventsByPkAndSkParseFilters(queryInput, input)
		queryInput.FilterExpression = aws.String(input.FilterExpression)
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeEvents(queryOutput.Items)
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

func (t EventTable) deleteEvent(ctx context.Context, pk string, sk string) error {

	key, err := attributevalue.MarshalMap(ddbEventPrimaryKey{
		Pk: pk,
		Sk: sk,
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

func (t EventTable) getEventsBySkAndData(ctx context.Context, input db.GetEventsBySkAndDataInput, fn func(m *models.Event, lastEvent bool) bool) error {
	if input.DataStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DataStartingAt or input.StartingAfter")
	}
	if input.Sk == "" {
		return fmt.Errorf("Hash key input.Sk cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("bySK"),
		ExpressionAttributeNames: map[string]string{
			"#SK": "sk",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk": &types.AttributeValueMemberS{
				Value: input.Sk,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.DataStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#SK = :sk")
	} else {
		queryInput.ExpressionAttributeNames["#DATA"] = "data"
		queryInput.ExpressionAttributeValues[":data"] = &types.AttributeValueMemberB{
			Value: input.DataStartingAt,
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#SK = :sk AND #DATA <= :data")
		} else {
			queryInput.KeyConditionExpression = aws.String("#SK = :sk AND #DATA >= :data")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"data": &types.AttributeValueMemberB{
				Value: input.StartingAfter.Data,
			},
			"sk": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Sk,
			},
			"pk": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Pk,
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeEvents(queryOutput.Items)
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
func (t EventTable) scanEventsBySkAndData(ctx context.Context, input db.ScanEventsBySkAndDataInput, fn func(m *models.Event, lastEvent bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	scanInput.IndexName = aws.String("bySK")
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"pk":   exclusiveStartKey["pk"],
			"sk":   exclusiveStartKey["sk"],
			"data": exclusiveStartKey["data"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeEvents(out.Items)
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

// encodeEvent encodes a Event as a DynamoDB map of attribute values.
func encodeEvent(m models.Event) (map[string]types.AttributeValue, error) {
	// First marshal the model to get all fields
	rawVal, err := attributevalue.MarshalMap(m)
	if err != nil {
		return nil, err
	}

	// Create a new map with the correct field names from json tags
	val := make(map[string]types.AttributeValue)

	// Get the type of the Event struct
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

// decodeEvent translates a Event stored in DynamoDB to a Event struct.
func decodeEvent(m map[string]types.AttributeValue, out *models.Event) error {
	var ddbEvent ddbEvent
	if err := attributevalue.UnmarshalMap(m, &ddbEvent); err != nil {
		return err
	}
	*out = ddbEvent.Event
	return nil
}

// decodeEvents translates a list of Events stored in DynamoDB to a slice of Event structs.
func decodeEvents(ms []map[string]types.AttributeValue) ([]models.Event, error) {
	events := make([]models.Event, len(ms))
	for i, m := range ms {
		var event models.Event
		if err := decodeEvent(m, &event); err != nil {
			return nil, err
		}
		events[i] = event
	}
	return events, nil
}
