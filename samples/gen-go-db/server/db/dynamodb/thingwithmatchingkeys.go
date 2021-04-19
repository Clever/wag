package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Clever/wag/v7/samples/gen-go-db/models"
	"github.com/Clever/wag/v7/samples/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithMatchingKeysTable represents the user-configurable properties of the ThingWithMatchingKeys table.
type ThingWithMatchingKeysTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithMatchingKeysPrimaryKey represents the primary key of a ThingWithMatchingKeys in DynamoDB.
type ddbThingWithMatchingKeysPrimaryKey struct {
	Bear        string `dynamodbav:"bear"`
	AssocTypeID string `dynamodbav:"assocTypeID"`
}

// ddbThingWithMatchingKeysGSIByAssoc represents the byAssoc GSI.
type ddbThingWithMatchingKeysGSIByAssoc struct {
	AssocTypeID string `dynamodbav:"assocTypeID"`
	CreatedBear string `dynamodbav:"createdBear"`
}

// ddbThingWithMatchingKeys represents a ThingWithMatchingKeys as stored in DynamoDB.
type ddbThingWithMatchingKeys struct {
	models.ThingWithMatchingKeys
}

func (t ThingWithMatchingKeysTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-matching-keyss", t.Prefix)
}

func (t ThingWithMatchingKeysTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("assocTypeID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("bear"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("createdBear"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("bear"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("assocTypeID"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("byAssoc"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("assocTypeID"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("createdBear"),
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

func (t ThingWithMatchingKeysTable) saveThingWithMatchingKeys(ctx context.Context, m models.ThingWithMatchingKeys) error {
	data, err := encodeThingWithMatchingKeys(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
	})
	return err
}

func (t ThingWithMatchingKeysTable) getThingWithMatchingKeys(ctx context.Context, bear string, assocType string, assocID string) (*models.ThingWithMatchingKeys, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithMatchingKeysPrimaryKey{
		Bear:        bear,
		AssocTypeID: fmt.Sprintf("%s^%s", assocType, assocID),
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
		return nil, db.ErrThingWithMatchingKeysNotFound{
			Bear:      bear,
			AssocType: assocType,
			AssocID:   assocID,
		}
	}

	var m models.ThingWithMatchingKeys
	if err := decodeThingWithMatchingKeys(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithMatchingKeysTable) scanThingWithMatchingKeyss(ctx context.Context, input db.ScanThingWithMatchingKeyssInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
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
			"bear": exclusiveStartKey["bear"],
			"assocTypeID": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s^%s", input.StartingAfter.AssocType, input.StartingAfter.AssocID)),
			},
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithMatchingKeyss(out.Items)
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

func (t ThingWithMatchingKeysTable) getThingWithMatchingKeyssByBearAndAssocTypeID(ctx context.Context, input db.GetThingWithMatchingKeyssByBearAndAssocTypeIDInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
	if input.StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of StartingAt or StartingAfter")
	}
	if input.Bear == "" {
		return fmt.Errorf("Hash key input.Bear cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#BEAR": aws.String("bear"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":bear": &dynamodb.AttributeValue{
				S: aws.String(input.Bear),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#BEAR = :bear")
	} else {
		queryInput.ExpressionAttributeNames["#ASSOCTYPEID"] = aws.String("assocTypeID")
		queryInput.ExpressionAttributeValues[":assocTypeId"] = &dynamodb.AttributeValue{
			S: aws.String(fmt.Sprintf("%s^%s", input.StartingAt.AssocType, input.StartingAt.AssocID)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#BEAR = :bear AND #ASSOCTYPEID <= :assocTypeId")
		} else {
			queryInput.KeyConditionExpression = aws.String("#BEAR = :bear AND #ASSOCTYPEID >= :assocTypeId")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"assocTypeID": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s^%s", input.StartingAfter.AssocType, input.StartingAfter.AssocID)),
			},
			"bear": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Bear),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithMatchingKeyss(queryOutput.Items)
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

func (t ThingWithMatchingKeysTable) deleteThingWithMatchingKeys(ctx context.Context, bear string, assocType string, assocID string) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithMatchingKeysPrimaryKey{
		Bear:        bear,
		AssocTypeID: fmt.Sprintf("%s^%s", assocType, assocID),
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

func (t ThingWithMatchingKeysTable) getThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx context.Context, input db.GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
	if input.StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.StartingAt or input.StartingAfter")
	}
	if input.AssocType == "" {
		return fmt.Errorf("Hash key input.AssocType cannot be empty")
	}
	if input.AssocID == "" {
		return fmt.Errorf("Hash key input.AssocID cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("byAssoc"),
		ExpressionAttributeNames: map[string]*string{
			"#ASSOCTYPEID": aws.String("assocTypeID"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":assocTypeId": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s^%s", input.AssocType, input.AssocID)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#ASSOCTYPEID = :assocTypeId")
	} else {
		queryInput.ExpressionAttributeNames["#CREATEDBEAR"] = aws.String("createdBear")
		queryInput.ExpressionAttributeValues[":createdBear"] = &dynamodb.AttributeValue{
			S: aws.String(fmt.Sprintf("%s^%s", input.StartingAt.Created, input.StartingAt.Bear)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#ASSOCTYPEID = :assocTypeId AND #CREATEDBEAR <= :createdBear")
		} else {
			queryInput.KeyConditionExpression = aws.String("#ASSOCTYPEID = :assocTypeId AND #CREATEDBEAR >= :createdBear")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"createdBear": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s^%s", input.StartingAfter.Created, input.StartingAfter.Bear)),
			},
			"assocTypeID": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s^%s", input.StartingAfter.AssocType, input.StartingAfter.AssocID)),
			},
			"bear": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Bear),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithMatchingKeyss(queryOutput.Items)
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
func (t ThingWithMatchingKeysTable) scanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx context.Context, input db.ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("byAssoc"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"bear": exclusiveStartKey["bear"],
			"assocTypeID": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s^%s", input.StartingAfter.AssocType, input.StartingAfter.AssocID)),
			},
			"createdBear": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s^%s", input.StartingAfter.Created, input.StartingAfter.Bear)),
			},
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithMatchingKeyss(out.Items)
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

// encodeThingWithMatchingKeys encodes a ThingWithMatchingKeys as a DynamoDB map of attribute values.
func encodeThingWithMatchingKeys(m models.ThingWithMatchingKeys) (map[string]*dynamodb.AttributeValue, error) {
	val, err := dynamodbattribute.MarshalMap(ddbThingWithMatchingKeys{
		ThingWithMatchingKeys: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(m.AssocID, "^") {
		return nil, fmt.Errorf("assocID cannot contain '^': %s", m.AssocID)
	}
	if strings.Contains(m.AssocType, "^") {
		return nil, fmt.Errorf("assocType cannot contain '^': %s", m.AssocType)
	}
	if strings.Contains(m.Bear, "^") {
		return nil, fmt.Errorf("bear cannot contain '^': %s", m.Bear)
	}
	// add in composite attributes
	primaryKey, err := dynamodbattribute.MarshalMap(ddbThingWithMatchingKeysPrimaryKey{
		Bear:        m.Bear,
		AssocTypeID: fmt.Sprintf("%s^%s", m.AssocType, m.AssocID),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range primaryKey {
		val[k] = v
	}
	byAssoc, err := dynamodbattribute.MarshalMap(ddbThingWithMatchingKeysGSIByAssoc{
		AssocTypeID: fmt.Sprintf("%s^%s", m.AssocType, m.AssocID),
		CreatedBear: fmt.Sprintf("%s^%s", m.Created, m.Bear),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range byAssoc {
		val[k] = v
	}
	return val, err
}

// decodeThingWithMatchingKeys translates a ThingWithMatchingKeys stored in DynamoDB to a ThingWithMatchingKeys struct.
func decodeThingWithMatchingKeys(m map[string]*dynamodb.AttributeValue, out *models.ThingWithMatchingKeys) error {
	var ddbThingWithMatchingKeys ddbThingWithMatchingKeys
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithMatchingKeys); err != nil {
		return err
	}
	*out = ddbThingWithMatchingKeys.ThingWithMatchingKeys
	return nil
}

// decodeThingWithMatchingKeyss translates a list of ThingWithMatchingKeyss stored in DynamoDB to a slice of ThingWithMatchingKeys structs.
func decodeThingWithMatchingKeyss(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithMatchingKeys, error) {
	thingWithMatchingKeyss := make([]models.ThingWithMatchingKeys, len(ms))
	for i, m := range ms {
		var thingWithMatchingKeys models.ThingWithMatchingKeys
		if err := decodeThingWithMatchingKeys(m, &thingWithMatchingKeys); err != nil {
			return nil, err
		}
		thingWithMatchingKeyss[i] = thingWithMatchingKeys
	}
	return thingWithMatchingKeyss, nil
}
