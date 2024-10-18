package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db-custom-path/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db-custom-path/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithTransactMultipleGSITable represents the user-configurable properties of the ThingWithTransactMultipleGSI table.
type ThingWithTransactMultipleGSITable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithTransactMultipleGSIPrimaryKey represents the primary key of a ThingWithTransactMultipleGSI in DynamoDB.
type ddbThingWithTransactMultipleGSIPrimaryKey struct {
	DateH strfmt.Date `dynamodbav:"dateH"`
}

// ddbThingWithTransactMultipleGSIGSIRangeDate represents the rangeDate GSI.
type ddbThingWithTransactMultipleGSIGSIRangeDate struct {
	ID    string      `dynamodbav:"id"`
	DateR strfmt.Date `dynamodbav:"dateR"`
}

// ddbThingWithTransactMultipleGSIGSIHash represents the hash GSI.
type ddbThingWithTransactMultipleGSIGSIHash struct {
	DateH strfmt.Date `dynamodbav:"dateH"`
	ID    string      `dynamodbav:"id"`
}

// ddbThingWithTransactMultipleGSI represents a ThingWithTransactMultipleGSI as stored in DynamoDB.
type ddbThingWithTransactMultipleGSI struct {
	models.ThingWithTransactMultipleGSI
}

func (t ThingWithTransactMultipleGSITable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("dateH"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("dateR"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("dateH"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("rangeDate"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("id"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("dateR"),
						KeyType:       aws.String(dynamodb.KeyTypeRange),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("hash"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("dateH"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("id"),
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
		TableName: aws.String(t.TableName),
	}); err != nil {
		return fmt.Errorf("failed to create table %s: %w", t.TableName, err)
	}
	return nil
}

func (t ThingWithTransactMultipleGSITable) saveThingWithTransactMultipleGSI(ctx context.Context, m models.ThingWithTransactMultipleGSI) error {
	data, err := encodeThingWithTransactMultipleGSI(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]*string{
			"#DATEH": aws.String("dateH"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#DATEH)"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return db.ErrThingWithTransactMultipleGSIAlreadyExists{
					DateH: m.DateH,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
		}
		return err
	}
	return nil
}

func (t ThingWithTransactMultipleGSITable) getThingWithTransactMultipleGSI(ctx context.Context, dateH strfmt.Date) (*models.ThingWithTransactMultipleGSI, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithTransactMultipleGSIPrimaryKey{
		DateH: dateH,
	})
	if err != nil {
		return nil, err
	}
	res, err := t.DynamoDBAPI.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return nil, fmt.Errorf("table or index not found: %s", t.TableName)
			}
		}
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrThingWithTransactMultipleGSINotFound{
			DateH: dateH,
		}
	}

	var m models.ThingWithTransactMultipleGSI
	if err := decodeThingWithTransactMultipleGSI(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithTransactMultipleGSITable) scanThingWithTransactMultipleGSIs(ctx context.Context, input db.ScanThingWithTransactMultipleGSIsInput, fn func(m *models.ThingWithTransactMultipleGSI, lastThingWithTransactMultipleGSI bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
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
			"dateH": exclusiveStartKey["dateH"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithTransactMultipleGSIs(out.Items)
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

func (t ThingWithTransactMultipleGSITable) deleteThingWithTransactMultipleGSI(ctx context.Context, dateH strfmt.Date) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithTransactMultipleGSIPrimaryKey{
		DateH: dateH,
	})
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(t.TableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
		}
		return err
	}

	return nil
}

func (t ThingWithTransactMultipleGSITable) getThingWithTransactMultipleGSIsByIDAndDateR(ctx context.Context, input db.GetThingWithTransactMultipleGSIsByIDAndDateRInput, fn func(m *models.ThingWithTransactMultipleGSI, lastThingWithTransactMultipleGSI bool) bool) error {
	if input.DateRStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateRStartingAt or input.StartingAfter")
	}
	if input.ID == "" {
		return fmt.Errorf("Hash key input.ID cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("rangeDate"),
		ExpressionAttributeNames: map[string]*string{
			"#ID": aws.String("id"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": &dynamodb.AttributeValue{
				S: aws.String(input.ID),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.DateRStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#ID = :id")
	} else {
		queryInput.ExpressionAttributeNames["#DATER"] = aws.String("dateR")
		queryInput.ExpressionAttributeValues[":dateR"] = &dynamodb.AttributeValue{
			S: aws.String(dateToDynamoTimeString(*input.DateRStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#ID = :id AND #DATER <= :dateR")
		} else {
			queryInput.KeyConditionExpression = aws.String("#ID = :id AND #DATER >= :dateR")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"dateR": &dynamodb.AttributeValue{
				S: aws.String(dateToDynamoTimeString(input.StartingAfter.DateR)),
			},
			"id": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.ID),
			},
			"dateH": &dynamodb.AttributeValue{
				S: aws.String(dateToDynamoTimeString(input.StartingAfter.DateH)),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithTransactMultipleGSIs(queryOutput.Items)
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
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
		}
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}
func (t ThingWithTransactMultipleGSITable) getThingWithTransactMultipleGSIsByDateHAndID(ctx context.Context, input db.GetThingWithTransactMultipleGSIsByDateHAndIDInput, fn func(m *models.ThingWithTransactMultipleGSI, lastThingWithTransactMultipleGSI bool) bool) error {
	if input.IDStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.IDStartingAt or input.StartingAfter")
	}
	if dateToDynamoTimeString(input.DateH) == "" {
		return fmt.Errorf("Hash key input.DateH cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("hash"),
		ExpressionAttributeNames: map[string]*string{
			"#DATEH": aws.String("dateH"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":dateH": &dynamodb.AttributeValue{
				S: aws.String(dateToDynamoTimeString(input.DateH)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.IDStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#DATEH = :dateH")
	} else {
		queryInput.ExpressionAttributeNames["#ID"] = aws.String("id")
		queryInput.ExpressionAttributeValues[":id"] = &dynamodb.AttributeValue{
			S: aws.String(string(*input.IDStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#DATEH = :dateH AND #ID <= :id")
		} else {
			queryInput.KeyConditionExpression = aws.String("#DATEH = :dateH AND #ID >= :id")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(string(input.StartingAfter.ID)),
			},
			"dateH": &dynamodb.AttributeValue{
				S: aws.String(dateToDynamoTimeString(input.StartingAfter.DateH)),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithTransactMultipleGSIs(queryOutput.Items)
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
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
		}
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}
func (t ThingWithTransactMultipleGSITable) transactSaveThingWithTransactMultipleGSIAndThing(ctx context.Context, m1 models.ThingWithTransactMultipleGSI, m1Conditions *expression.ConditionBuilder, m2 models.Thing, m2Conditions *expression.ConditionBuilder) error {
	data1, err := encodeThingWithTransactMultipleGSI(m1)
	if err != nil {
		return err
	}

	m1CondExpr, m1ExprVals, m1ExprNames, err := buildCondExpr(m1Conditions)
	if err != nil {
		return err
	}

	data2, err := encodeThing(m2)
	if err != nil {
		return err
	}

	m2CondExpr, m2ExprVals, m2ExprNames, err := buildCondExpr(m2Conditions)
	if err != nil {
		return err
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					TableName:                 aws.String(t.TableName),
					Item:                      data1,
					ConditionExpression:       m1CondExpr,
					ExpressionAttributeValues: m1ExprVals,
					ExpressionAttributeNames:  m1ExprNames,
				},
			},
			{
				Put: &dynamodb.Put{
					TableName:                 aws.String(fmt.Sprintf("%s-Things", t.Prefix)),
					Item:                      data2,
					ConditionExpression:       m2CondExpr,
					ExpressionAttributeValues: m2ExprVals,
					ExpressionAttributeNames:  m2ExprNames,
				},
			},
		},
	}
	_, err = t.DynamoDBAPI.TransactWriteItemsWithContext(ctx, input)

	return err
}

// encodeThingWithTransactMultipleGSI encodes a ThingWithTransactMultipleGSI as a DynamoDB map of attribute values.
func encodeThingWithTransactMultipleGSI(m models.ThingWithTransactMultipleGSI) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingWithTransactMultipleGSI{
		ThingWithTransactMultipleGSI: m,
	})
}

// decodeThingWithTransactMultipleGSI translates a ThingWithTransactMultipleGSI stored in DynamoDB to a ThingWithTransactMultipleGSI struct.
func decodeThingWithTransactMultipleGSI(m map[string]*dynamodb.AttributeValue, out *models.ThingWithTransactMultipleGSI) error {
	var ddbThingWithTransactMultipleGSI ddbThingWithTransactMultipleGSI
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithTransactMultipleGSI); err != nil {
		return err
	}
	*out = ddbThingWithTransactMultipleGSI.ThingWithTransactMultipleGSI
	return nil
}

// decodeThingWithTransactMultipleGSIs translates a list of ThingWithTransactMultipleGSIs stored in DynamoDB to a slice of ThingWithTransactMultipleGSI structs.
func decodeThingWithTransactMultipleGSIs(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithTransactMultipleGSI, error) {
	thingWithTransactMultipleGSIs := make([]models.ThingWithTransactMultipleGSI, len(ms))
	for i, m := range ms {
		var thingWithTransactMultipleGSI models.ThingWithTransactMultipleGSI
		if err := decodeThingWithTransactMultipleGSI(m, &thingWithTransactMultipleGSI); err != nil {
			return nil, err
		}
		thingWithTransactMultipleGSIs[i] = thingWithTransactMultipleGSI
	}
	return thingWithTransactMultipleGSIs, nil
}
