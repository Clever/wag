package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Clever/wag/v7/samples/gen-go-db-custom-path/db"
	"github.com/Clever/wag/v7/samples/gen-go-db-custom-path/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithMultiUseCompositeAttributeTable represents the user-configurable properties of the ThingWithMultiUseCompositeAttribute table.
type ThingWithMultiUseCompositeAttributeTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithMultiUseCompositeAttributePrimaryKey represents the primary key of a ThingWithMultiUseCompositeAttribute in DynamoDB.
type ddbThingWithMultiUseCompositeAttributePrimaryKey struct {
	One string `dynamodbav:"one"`
}

// ddbThingWithMultiUseCompositeAttributeGSIThreeIndex represents the threeIndex GSI.
type ddbThingWithMultiUseCompositeAttributeGSIThreeIndex struct {
	Three  string `dynamodbav:"three"`
	OneTwo string `dynamodbav:"one_two"`
}

// ddbThingWithMultiUseCompositeAttributeGSIFourIndex represents the fourIndex GSI.
type ddbThingWithMultiUseCompositeAttributeGSIFourIndex struct {
	Four   string `dynamodbav:"four"`
	OneTwo string `dynamodbav:"one_two"`
}

// ddbThingWithMultiUseCompositeAttribute represents a ThingWithMultiUseCompositeAttribute as stored in DynamoDB.
type ddbThingWithMultiUseCompositeAttribute struct {
	models.ThingWithMultiUseCompositeAttribute
}

func (t ThingWithMultiUseCompositeAttributeTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-multi-use-composite-attributes", t.Prefix)
}

func (t ThingWithMultiUseCompositeAttributeTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("four"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("one"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("one_two"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("three"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("one"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("threeIndex"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("three"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("one_two"),
						KeyType:       aws.String(dynamodb.KeyTypeRange),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("fourIndex"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("four"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("one_two"),
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

func (t ThingWithMultiUseCompositeAttributeTable) saveThingWithMultiUseCompositeAttribute(ctx context.Context, m models.ThingWithMultiUseCompositeAttribute) error {
	data, err := encodeThingWithMultiUseCompositeAttribute(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
	})
	return err
}

func (t ThingWithMultiUseCompositeAttributeTable) getThingWithMultiUseCompositeAttribute(ctx context.Context, one string) (*models.ThingWithMultiUseCompositeAttribute, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithMultiUseCompositeAttributePrimaryKey{
		One: one,
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
		return nil, db.ErrThingWithMultiUseCompositeAttributeNotFound{
			One: one,
		}
	}

	var m models.ThingWithMultiUseCompositeAttribute
	if err := decodeThingWithMultiUseCompositeAttribute(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithMultiUseCompositeAttributeTable) scanThingWithMultiUseCompositeAttributes(ctx context.Context, input db.ScanThingWithMultiUseCompositeAttributesInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
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
			"one": exclusiveStartKey["one"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithMultiUseCompositeAttributes(out.Items)
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

func (t ThingWithMultiUseCompositeAttributeTable) deleteThingWithMultiUseCompositeAttribute(ctx context.Context, one string) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithMultiUseCompositeAttributePrimaryKey{
		One: one,
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

func (t ThingWithMultiUseCompositeAttributeTable) getThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx context.Context, input db.GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	if input.StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.StartingAt or input.StartingAfter")
	}
	if input.Three == "" {
		return fmt.Errorf("Hash key input.Three cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("threeIndex"),
		ExpressionAttributeNames: map[string]*string{
			"#THREE": aws.String("three"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":three": &dynamodb.AttributeValue{
				S: aws.String(input.Three),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#THREE = :three")
	} else {
		queryInput.ExpressionAttributeNames["#ONE_TWO"] = aws.String("one_two")
		queryInput.ExpressionAttributeValues[":oneTwo"] = &dynamodb.AttributeValue{
			S: aws.String(fmt.Sprintf("%s_%s", input.StartingAt.One, input.StartingAt.Two)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#THREE = :three AND #ONE_TWO <= :oneTwo")
		} else {
			queryInput.KeyConditionExpression = aws.String("#THREE = :three AND #ONE_TWO >= :oneTwo")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"one_two": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s_%s", *input.StartingAfter.One, *input.StartingAfter.Two)),
			},
			"three": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.Three),
			},
			"one": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.One),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithMultiUseCompositeAttributes(queryOutput.Items)
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
func (t ThingWithMultiUseCompositeAttributeTable) scanThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx context.Context, input db.ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("threeIndex"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"one":   exclusiveStartKey["one"],
			"three": exclusiveStartKey["three"],
			"one_two": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s_%s", *input.StartingAfter.One, *input.StartingAfter.Two)),
			},
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithMultiUseCompositeAttributes(out.Items)
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

func (t ThingWithMultiUseCompositeAttributeTable) getThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx context.Context, input db.GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	if input.StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.StartingAt or input.StartingAfter")
	}
	if input.Four == "" {
		return fmt.Errorf("Hash key input.Four cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("fourIndex"),
		ExpressionAttributeNames: map[string]*string{
			"#FOUR": aws.String("four"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":four": &dynamodb.AttributeValue{
				S: aws.String(input.Four),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#FOUR = :four")
	} else {
		queryInput.ExpressionAttributeNames["#ONE_TWO"] = aws.String("one_two")
		queryInput.ExpressionAttributeValues[":oneTwo"] = &dynamodb.AttributeValue{
			S: aws.String(fmt.Sprintf("%s_%s", input.StartingAt.One, input.StartingAt.Two)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#FOUR = :four AND #ONE_TWO <= :oneTwo")
		} else {
			queryInput.KeyConditionExpression = aws.String("#FOUR = :four AND #ONE_TWO >= :oneTwo")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"one_two": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s_%s", *input.StartingAfter.One, *input.StartingAfter.Two)),
			},
			"four": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.Four),
			},
			"one": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.One),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithMultiUseCompositeAttributes(queryOutput.Items)
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
func (t ThingWithMultiUseCompositeAttributeTable) scanThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx context.Context, input db.ScanThingWithMultiUseCompositeAttributesByFourAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("fourIndex"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"one":  exclusiveStartKey["one"],
			"four": exclusiveStartKey["four"],
			"one_two": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s_%s", *input.StartingAfter.One, *input.StartingAfter.Two)),
			},
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithMultiUseCompositeAttributes(out.Items)
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

// encodeThingWithMultiUseCompositeAttribute encodes a ThingWithMultiUseCompositeAttribute as a DynamoDB map of attribute values.
func encodeThingWithMultiUseCompositeAttribute(m models.ThingWithMultiUseCompositeAttribute) (map[string]*dynamodb.AttributeValue, error) {
	val, err := dynamodbattribute.MarshalMap(ddbThingWithMultiUseCompositeAttribute{
		ThingWithMultiUseCompositeAttribute: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(*m.One, "_") {
		return nil, fmt.Errorf("one cannot contain '_': %s", *m.One)
	}
	if strings.Contains(*m.Two, "_") {
		return nil, fmt.Errorf("two cannot contain '_': %s", *m.Two)
	}
	// add in composite attributes
	threeIndex, err := dynamodbattribute.MarshalMap(ddbThingWithMultiUseCompositeAttributeGSIThreeIndex{
		Three:  *m.Three,
		OneTwo: fmt.Sprintf("%s_%s", *m.One, *m.Two),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range threeIndex {
		val[k] = v
	}
	fourIndex, err := dynamodbattribute.MarshalMap(ddbThingWithMultiUseCompositeAttributeGSIFourIndex{
		Four:   *m.Four,
		OneTwo: fmt.Sprintf("%s_%s", *m.One, *m.Two),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range fourIndex {
		val[k] = v
	}
	return val, err
}

// decodeThingWithMultiUseCompositeAttribute translates a ThingWithMultiUseCompositeAttribute stored in DynamoDB to a ThingWithMultiUseCompositeAttribute struct.
func decodeThingWithMultiUseCompositeAttribute(m map[string]*dynamodb.AttributeValue, out *models.ThingWithMultiUseCompositeAttribute) error {
	var ddbThingWithMultiUseCompositeAttribute ddbThingWithMultiUseCompositeAttribute
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithMultiUseCompositeAttribute); err != nil {
		return err
	}
	*out = ddbThingWithMultiUseCompositeAttribute.ThingWithMultiUseCompositeAttribute
	return nil
}

// decodeThingWithMultiUseCompositeAttributes translates a list of ThingWithMultiUseCompositeAttributes stored in DynamoDB to a slice of ThingWithMultiUseCompositeAttribute structs.
func decodeThingWithMultiUseCompositeAttributes(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithMultiUseCompositeAttribute, error) {
	thingWithMultiUseCompositeAttributes := make([]models.ThingWithMultiUseCompositeAttribute, len(ms))
	for i, m := range ms {
		var thingWithMultiUseCompositeAttribute models.ThingWithMultiUseCompositeAttribute
		if err := decodeThingWithMultiUseCompositeAttribute(m, &thingWithMultiUseCompositeAttribute); err != nil {
			return nil, err
		}
		thingWithMultiUseCompositeAttributes[i] = thingWithMultiUseCompositeAttribute
	}
	return thingWithMultiUseCompositeAttributes, nil
}
