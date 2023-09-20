package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Clever/wag/samples/gen-go-db/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithRequiredCompositePropertiesAndKeysOnlyTable represents the user-configurable properties of the ThingWithRequiredCompositePropertiesAndKeysOnly table.
type ThingWithRequiredCompositePropertiesAndKeysOnlyTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
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

func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-required-composite-properties-and-keys-onlys", t.Prefix)
}

func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("propertyOneAndTwo"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("propertyThree"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("propertyThree"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("propertyOneAndTwo_PropertyThree"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("KEYS_ONLY"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("propertyOneAndTwo"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("propertyThree"),
						KeyType:       aws.String(dynamodb.KeyTypeRange),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
		},
		TableName: aws.String(t.name()),
	}); err != nil {
		return err
	}
	return nil
}

func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) saveThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, m models.ThingWithRequiredCompositePropertiesAndKeysOnly) error {
	data, err := encodeThingWithRequiredCompositePropertiesAndKeysOnly(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
	})
	return err
}

func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) getThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, propertyThree string) (*models.ThingWithRequiredCompositePropertiesAndKeysOnly, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithRequiredCompositePropertiesAndKeysOnlyPrimaryKey{
		PropertyThree: propertyThree,
	})
	if err != nil {
		return nil, err
	}
	res, err := t.DynamoDBAPI.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(t.name()),
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
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"propertyThree": exclusiveStartKey["propertyThree"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithRequiredCompositePropertiesAndKeysOnlys(out.Items)
		if err != nil {
			innerErr = fmt.Errorf("error decoding %s", err.Error())
			return false
		}
		for i := range items {
			if input.Limiter != nil {
				if err := input.Limiter.Wait(ctx); err != nil {
					innerErr = err
					return false
				}
			}
			isLastModel := lastPage && i == len(items)-1
			if shouldContinue := fn(&items[i], isLastModel); !shouldContinue {
				return false
			}
			totalRecordsProcessed++
			// if the Limit of records have been passed to fn, don't pass anymore records.
			if input.Limit != nil && totalRecordsProcessed == *input.Limit {
				return false
			}
		}
		return true
	})
	if innerErr != nil {
		return innerErr
	}
	return err
}

func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) deleteThingWithRequiredCompositePropertiesAndKeysOnly(ctx context.Context, propertyThree string) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithRequiredCompositePropertiesAndKeysOnlyPrimaryKey{
		PropertyThree: propertyThree,
	})
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(t.name()),
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
		TableName: aws.String(t.name()),
		IndexName: aws.String("propertyOneAndTwo_PropertyThree"),
		ExpressionAttributeNames: map[string]*string{
			"#PROPERTYONEANDTWO": aws.String("propertyOneAndTwo"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":propertyOneAndTwo": {
				S: aws.String(fmt.Sprintf("%s_%s", input.PropertyOne, input.PropertyTwo)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.PropertyThreeStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#PROPERTYONEANDTWO = :propertyOneAndTwo")
	} else {
		queryInput.ExpressionAttributeNames["#PROPERTYTHREE"] = aws.String("propertyThree")
		queryInput.ExpressionAttributeValues[":propertyThree"] = &dynamodb.AttributeValue{
			S: aws.String(*input.PropertyThreeStartingAt),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#PROPERTYONEANDTWO = :propertyOneAndTwo AND #PROPERTYTHREE <= :propertyThree")
		} else {
			queryInput.KeyConditionExpression = aws.String("#PROPERTYONEANDTWO = :propertyOneAndTwo AND #PROPERTYTHREE >= :propertyThree")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"propertyThree": {
				S: aws.String(*input.StartingAfter.PropertyThree),
			},
			"propertyOneAndTwo": {
				S: aws.String(fmt.Sprintf("%s_%s", *input.StartingAfter.PropertyOne, *input.StartingAfter.PropertyTwo)),
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

	err := t.DynamoDBAPI.QueryPagesWithContext(ctx, queryInput, pageFn)
	if err != nil {
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}
func (t ThingWithRequiredCompositePropertiesAndKeysOnlyTable) scanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThree(ctx context.Context, input db.ScanThingWithRequiredCompositePropertiesAndKeysOnlysByPropertyOneAndTwoAndPropertyThreeInput, fn func(m *models.ThingWithRequiredCompositePropertiesAndKeysOnly, lastThingWithRequiredCompositePropertiesAndKeysOnly bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("propertyOneAndTwo_PropertyThree"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"propertyThree": exclusiveStartKey["propertyThree"],
			"propertyOneAndTwo": {
				S: aws.String(fmt.Sprintf("%s_%s", *input.StartingAfter.PropertyOne, *input.StartingAfter.PropertyTwo)),
			},
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithRequiredCompositePropertiesAndKeysOnlys(out.Items)
		if err != nil {
			innerErr = fmt.Errorf("error decoding %s", err.Error())
			return false
		}
		for i := range items {
			if input.Limiter != nil {
				if err := input.Limiter.Wait(ctx); err != nil {
					innerErr = err
					return false
				}
			}
			isLastModel := lastPage && i == len(items)-1
			if shouldContinue := fn(&items[i], isLastModel); !shouldContinue {
				return false
			}
			totalRecordsProcessed++
			// if the Limit of records have been passed to fn, don't pass anymore records.
			if input.Limit != nil && totalRecordsProcessed == *input.Limit {
				return false
			}
		}
		return true
	})
	if innerErr != nil {
		return innerErr
	}
	return err
}

// encodeThingWithRequiredCompositePropertiesAndKeysOnly encodes a ThingWithRequiredCompositePropertiesAndKeysOnly as a DynamoDB map of attribute values.
func encodeThingWithRequiredCompositePropertiesAndKeysOnly(m models.ThingWithRequiredCompositePropertiesAndKeysOnly) (map[string]*dynamodb.AttributeValue, error) {
	val, err := dynamodbattribute.MarshalMap(ddbThingWithRequiredCompositePropertiesAndKeysOnly{
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
	propertyOneAndTwoPropertyThree, err := dynamodbattribute.MarshalMap(ddbThingWithRequiredCompositePropertiesAndKeysOnlyGSIPropertyOneAndTwoPropertyThree{
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
func decodeThingWithRequiredCompositePropertiesAndKeysOnly(m map[string]*dynamodb.AttributeValue, out *models.ThingWithRequiredCompositePropertiesAndKeysOnly) error {
	var ddbThingWithRequiredCompositePropertiesAndKeysOnly ddbThingWithRequiredCompositePropertiesAndKeysOnly
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithRequiredCompositePropertiesAndKeysOnly); err != nil {
		return err
	}
	*out = ddbThingWithRequiredCompositePropertiesAndKeysOnly.ThingWithRequiredCompositePropertiesAndKeysOnly
	// parse composite attributes from projected secondary indexes and fill
	// in model properties
	if v, ok := m["propertyOneAndTwo"]; ok && v.S != nil {
		parts := strings.Split(*v.S, "_")
		if len(parts) != 2 {
			return fmt.Errorf("expected 2 parts: '%s'", *v.S)
		}
		out.PropertyOne = &parts[0]
		out.PropertyTwo = &parts[1]
	}
	return nil
}

// decodeThingWithRequiredCompositePropertiesAndKeysOnlys translates a list of ThingWithRequiredCompositePropertiesAndKeysOnlys stored in DynamoDB to a slice of ThingWithRequiredCompositePropertiesAndKeysOnly structs.
func decodeThingWithRequiredCompositePropertiesAndKeysOnlys(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithRequiredCompositePropertiesAndKeysOnly, error) {
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
