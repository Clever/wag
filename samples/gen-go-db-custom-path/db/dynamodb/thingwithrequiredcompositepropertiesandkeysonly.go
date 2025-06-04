package dynamodb

import (
	"context"
	"errors"
	"fmt"
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
	// swad-get-7
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
			"propertyThree": exclusiveStartKey["propertyThree"],
		}
	}
	totalRecordsProcessed := int32(0)

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
	// swad-get-33
	if input.PropertyThreeStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.PropertyThreeStartingAt or input.StartingAfter")
	}
	// swad-get-33f
	if input.PropertyOne == "" {
		return fmt.Errorf("Hash key input.PropertyOne cannot be empty")
	}
	// swad-get-33f
	if input.PropertyTwo == "" {
		return fmt.Errorf("Hash key input.PropertyTwo cannot be empty")
	}
	// swad-get-331
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("propertyOneAndTwo_PropertyThree"),
		ExpressionAttributeNames: map[string]string{
			"#PROPERTYONEANDTWO": "propertyOneAndTwo",
		},
		// swad-get-3312
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":propertyOneAndTwo": &types.AttributeValueMemberS{
				// swad-get-33a
				Value: fmt.Sprintf("%s_%s", input.PropertyOne, input.PropertyTwo),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	// swad-get-332
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.PropertyThreeStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#PROPERTYONEANDTWO = :propertyOneAndTwo")
	} else {
		// swad-get-333
		queryInput.ExpressionAttributeNames["#PROPERTYTHREE"] = "propertyThree"

		// swad-get-3331a
		queryInput.ExpressionAttributeValues[":propertyThree"] = &types.AttributeValueMemberS{
			Value: string(*input.PropertyThreeStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#PROPERTYONEANDTWO = :propertyOneAndTwo AND #PROPERTYTHREE <= :propertyThree")
		} else {
			queryInput.KeyConditionExpression = aws.String("#PROPERTYONEANDTWO = :propertyOneAndTwo AND #PROPERTYTHREE >= :propertyThree")
		}
	}
	// swad-get-334
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"propertyThree": &types.AttributeValueMemberS{
				Value: string(*input.StartingAfter.PropertyThree),
			},
			// swad-get-3341
			"propertyOneAndTwo": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s_%s", *input.StartingAfter.PropertyOne, *input.StartingAfter.PropertyTwo),
			},
			// swad-get-3342

			// swad-get-336
		}
	}

	// swad-get-339

	totalRecordsProcessed := int32(0)
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
		Limit:          input.Limit,
		IndexName:      aws.String("propertyOneAndTwo_PropertyThree"),
	}
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
	totalRecordsProcessed := int32(0)

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
	val, err := attributevalue.MarshalMap(ddbThingWithRequiredCompositePropertiesAndKeysOnly{
		ThingWithRequiredCompositePropertiesAndKeysOnly: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(*m.PropertyOne, "_") {
		return nil, fmt.Errorf("propertyOne cannot contain '_': %s", *m.PropertyOne)
	}
	if strings.Contains(*m.PropertyTwo, "_") {
		return nil, fmt.Errorf("propertyTwo cannot contain '_': %s", *m.PropertyTwo)
	}
	// add in composite attributes
	propertyOneAndTwoPropertyThree, err := attributevalue.MarshalMap(ddbThingWithRequiredCompositePropertiesAndKeysOnlyGSIPropertyOneAndTwoPropertyThree{
		PropertyOneAndTwo: fmt.Sprintf("%s_%s", *m.PropertyOne, *m.PropertyTwo),
		PropertyThree:     *m.PropertyThree,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range propertyOneAndTwoPropertyThree {
		val[k] = v
	}
	return val, err
}

// decodeThingWithRequiredCompositePropertiesAndKeysOnly translates a ThingWithRequiredCompositePropertiesAndKeysOnly stored in DynamoDB to a ThingWithRequiredCompositePropertiesAndKeysOnly struct.
func decodeThingWithRequiredCompositePropertiesAndKeysOnly(m map[string]types.AttributeValue, out *models.ThingWithRequiredCompositePropertiesAndKeysOnly) error {
	// swad-decode-1
	var ddbThingWithRequiredCompositePropertiesAndKeysOnly ddbThingWithRequiredCompositePropertiesAndKeysOnly
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithRequiredCompositePropertiesAndKeysOnly); err != nil {
		return err
	}
	*out = ddbThingWithRequiredCompositePropertiesAndKeysOnly.ThingWithRequiredCompositePropertiesAndKeysOnly
	// parse composite attributes from projected secondary indexes and fill
	// in model properties
	// swad-decode-2
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
