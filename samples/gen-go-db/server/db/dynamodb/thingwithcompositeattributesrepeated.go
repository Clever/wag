package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Clever/wag/v6/samples/gen-go-db/models"
	"github.com/Clever/wag/v6/samples/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithCompositeAttributesRepeatedTable represents the user-configurable properties of the ThingWithCompositeAttributesRepeated table.
type ThingWithCompositeAttributesRepeatedTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithCompositeAttributesRepeatedPrimaryKey represents the primary key of a ThingWithCompositeAttributesRepeated in DynamoDB.
type ddbThingWithCompositeAttributesRepeatedPrimaryKey struct {
	Uno string `dynamodbav:"uno"`
}

// ddbThingWithCompositeAttributesRepeatedGSIThreeIndex represents the threeIndex GSI.
type ddbThingWithCompositeAttributesRepeatedGSIThreeIndex struct {
	Tres   string `dynamodbav:"tres"`
	UnoDos string `dynamodbav:"uno_dos"`
}

// ddbThingWithCompositeAttributesRepeatedGSIFourIndex represents the fourIndex GSI.
type ddbThingWithCompositeAttributesRepeatedGSIFourIndex struct {
	Cuatro string `dynamodbav:"cuatro"`
	Tres   string `dynamodbav:"tres"`
}

// ddbThingWithCompositeAttributesRepeated represents a ThingWithCompositeAttributesRepeated as stored in DynamoDB.
type ddbThingWithCompositeAttributesRepeated struct {
	models.ThingWithCompositeAttributesRepeated
}

func (t ThingWithCompositeAttributesRepeatedTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-composite-attributes-repeateds", t.Prefix)
}

func (t ThingWithCompositeAttributesRepeatedTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("cuatro"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("tres"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("uno"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("uno_dos"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("uno"),
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
						AttributeName: aws.String("tres"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("uno_dos"),
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
						AttributeName: aws.String("cuatro"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("tres"),
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

func (t ThingWithCompositeAttributesRepeatedTable) saveThingWithCompositeAttributesRepeated(ctx context.Context, m models.ThingWithCompositeAttributesRepeated) error {
	data, err := encodeThingWithCompositeAttributesRepeated(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
	})
	return err
}

func (t ThingWithCompositeAttributesRepeatedTable) getThingWithCompositeAttributesRepeated(ctx context.Context, uno string) (*models.ThingWithCompositeAttributesRepeated, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeAttributesRepeatedPrimaryKey{
		Uno: uno,
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
		return nil, db.ErrThingWithCompositeAttributesRepeatedNotFound{
			Uno: uno,
		}
	}

	var m models.ThingWithCompositeAttributesRepeated
	if err := decodeThingWithCompositeAttributesRepeated(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithCompositeAttributesRepeatedTable) scanThingWithCompositeAttributesRepeateds(ctx context.Context, input db.ScanThingWithCompositeAttributesRepeatedsInput, fn func(m *models.ThingWithCompositeAttributesRepeated, lastThingWithCompositeAttributesRepeated bool) bool) error {
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
			"uno": exclusiveStartKey["uno"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithCompositeAttributesRepeateds(out.Items)
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

func (t ThingWithCompositeAttributesRepeatedTable) deleteThingWithCompositeAttributesRepeated(ctx context.Context, uno string) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeAttributesRepeatedPrimaryKey{
		Uno: uno,
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

func (t ThingWithCompositeAttributesRepeatedTable) getThingWithCompositeAttributesRepeatedsByTresAndUnoDos(ctx context.Context, input db.GetThingWithCompositeAttributesRepeatedsByTresAndUnoDosInput, fn func(m *models.ThingWithCompositeAttributesRepeated, lastThingWithCompositeAttributesRepeated bool) bool) error {
	if input.StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.StartingAt or input.StartingAfter")
	}
	if input.Tres == "" {
		return fmt.Errorf("Hash key input.Tres cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("threeIndex"),
		ExpressionAttributeNames: map[string]*string{
			"#TRES": aws.String("tres"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":tres": &dynamodb.AttributeValue{
				S: aws.String(input.Tres),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#TRES = :tres")
	} else {
		queryInput.ExpressionAttributeNames["#UNO_DOS"] = aws.String("uno_dos")
		queryInput.ExpressionAttributeValues[":unoDos"] = &dynamodb.AttributeValue{
			S: aws.String(fmt.Sprintf("%s_%s", input.StartingAt.Uno, input.StartingAt.Dos)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#TRES = :tres AND #UNO_DOS <= :unoDos")
		} else {
			queryInput.KeyConditionExpression = aws.String("#TRES = :tres AND #UNO_DOS >= :unoDos")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"uno_dos": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s_%s", *input.StartingAfter.Uno, *input.StartingAfter.Dos)),
			},
			"tres": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.Tres),
			},
			"uno": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.Uno),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithCompositeAttributesRepeateds(queryOutput.Items)
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
func (t ThingWithCompositeAttributesRepeatedTable) scanThingWithCompositeAttributesRepeatedsByTresAndUnoDos(ctx context.Context, input db.ScanThingWithCompositeAttributesRepeatedsByTresAndUnoDosInput, fn func(m *models.ThingWithCompositeAttributesRepeated, lastThingWithCompositeAttributesRepeated bool) bool) error {
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
			"uno":  exclusiveStartKey["uno"],
			"tres": exclusiveStartKey["tres"],
			"uno_dos": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s_%s", *input.StartingAfter.Uno, *input.StartingAfter.Dos)),
			},
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithCompositeAttributesRepeateds(out.Items)
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

func (t ThingWithCompositeAttributesRepeatedTable) getThingWithCompositeAttributesRepeatedsByCuatroAndTres(ctx context.Context, input db.GetThingWithCompositeAttributesRepeatedsByCuatroAndTresInput, fn func(m *models.ThingWithCompositeAttributesRepeated, lastThingWithCompositeAttributesRepeated bool) bool) error {
	if input.TresStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.TresStartingAt or input.StartingAfter")
	}
	if input.Cuatro == "" {
		return fmt.Errorf("Hash key input.Cuatro cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("fourIndex"),
		ExpressionAttributeNames: map[string]*string{
			"#CUATRO": aws.String("cuatro"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":cuatro": &dynamodb.AttributeValue{
				S: aws.String(input.Cuatro),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.TresStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#CUATRO = :cuatro")
	} else {
		queryInput.ExpressionAttributeNames["#TRES"] = aws.String("tres")
		queryInput.ExpressionAttributeValues[":tres"] = &dynamodb.AttributeValue{
			S: aws.String(*input.TresStartingAt),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#CUATRO = :cuatro AND #TRES <= :tres")
		} else {
			queryInput.KeyConditionExpression = aws.String("#CUATRO = :cuatro AND #TRES >= :tres")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"tres": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.Tres),
			},
			"cuatro": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.Cuatro),
			},
			"uno": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.Uno),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithCompositeAttributesRepeateds(queryOutput.Items)
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
func (t ThingWithCompositeAttributesRepeatedTable) scanThingWithCompositeAttributesRepeatedsByCuatroAndTres(ctx context.Context, input db.ScanThingWithCompositeAttributesRepeatedsByCuatroAndTresInput, fn func(m *models.ThingWithCompositeAttributesRepeated, lastThingWithCompositeAttributesRepeated bool) bool) error {
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
			"uno":    exclusiveStartKey["uno"],
			"cuatro": exclusiveStartKey["cuatro"],
			"tres":   exclusiveStartKey["tres"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithCompositeAttributesRepeateds(out.Items)
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

// encodeThingWithCompositeAttributesRepeated encodes a ThingWithCompositeAttributesRepeated as a DynamoDB map of attribute values.
func encodeThingWithCompositeAttributesRepeated(m models.ThingWithCompositeAttributesRepeated) (map[string]*dynamodb.AttributeValue, error) {
	val, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeAttributesRepeated{
		ThingWithCompositeAttributesRepeated: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(*m.Dos, "_") {
		return nil, fmt.Errorf("dos cannot contain '_': %s", *m.Dos)
	}
	if strings.Contains(*m.Uno, "_") {
		return nil, fmt.Errorf("uno cannot contain '_': %s", *m.Uno)
	}
	// add in composite attributes
	threeIndex, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeAttributesRepeatedGSIThreeIndex{
		Tres:   *m.Tres,
		UnoDos: fmt.Sprintf("%s_%s", *m.Uno, *m.Dos),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range threeIndex {
		val[k] = v
	}
	return val, err
}

// decodeThingWithCompositeAttributesRepeated translates a ThingWithCompositeAttributesRepeated stored in DynamoDB to a ThingWithCompositeAttributesRepeated struct.
func decodeThingWithCompositeAttributesRepeated(m map[string]*dynamodb.AttributeValue, out *models.ThingWithCompositeAttributesRepeated) error {
	var ddbThingWithCompositeAttributesRepeated ddbThingWithCompositeAttributesRepeated
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithCompositeAttributesRepeated); err != nil {
		return err
	}
	*out = ddbThingWithCompositeAttributesRepeated.ThingWithCompositeAttributesRepeated
	return nil
}

// decodeThingWithCompositeAttributesRepeateds translates a list of ThingWithCompositeAttributesRepeateds stored in DynamoDB to a slice of ThingWithCompositeAttributesRepeated structs.
func decodeThingWithCompositeAttributesRepeateds(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithCompositeAttributesRepeated, error) {
	thingWithCompositeAttributesRepeateds := make([]models.ThingWithCompositeAttributesRepeated, len(ms))
	for i, m := range ms {
		var thingWithCompositeAttributesRepeated models.ThingWithCompositeAttributesRepeated
		if err := decodeThingWithCompositeAttributesRepeated(m, &thingWithCompositeAttributesRepeated); err != nil {
			return nil, err
		}
		thingWithCompositeAttributesRepeateds[i] = thingWithCompositeAttributesRepeated
	}
	return thingWithCompositeAttributesRepeateds, nil
}
