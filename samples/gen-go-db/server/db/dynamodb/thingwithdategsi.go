package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithDateGSITable represents the user-configurable properties of the ThingWithDateGSI table.
type ThingWithDateGSITable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithDateGSIPrimaryKey represents the primary key of a ThingWithDateGSI in DynamoDB.
type ddbThingWithDateGSIPrimaryKey struct {
	DateH strfmt.Date `dynamodbav:"dateH"`
}

// ddbThingWithDateGSIGSIRangeDate represents the rangeDate GSI.
type ddbThingWithDateGSIGSIRangeDate struct {
	ID    string      `dynamodbav:"id"`
	DateR strfmt.Date `dynamodbav:"dateR"`
}

// ddbThingWithDateGSIGSIHash represents the hash GSI.
type ddbThingWithDateGSIGSIHash struct {
	DateH strfmt.Date `dynamodbav:"dateH"`
	ID    string      `dynamodbav:"id"`
}

// ddbThingWithDateGSI represents a ThingWithDateGSI as stored in DynamoDB.
type ddbThingWithDateGSI struct {
	models.ThingWithDateGSI
}

func (t ThingWithDateGSITable) create(ctx context.Context) error {
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

func (t ThingWithDateGSITable) saveThingWithDateGSI(ctx context.Context, m models.ThingWithDateGSI) error {
	data, err := encodeThingWithDateGSI(m)
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
				return db.ErrThingWithDateGSIAlreadyExists{
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

func (t ThingWithDateGSITable) getThingWithDateGSI(ctx context.Context, dateH strfmt.Date) (*models.ThingWithDateGSI, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithDateGSIPrimaryKey{
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
		return nil, db.ErrThingWithDateGSINotFound{
			DateH: dateH,
		}
	}

	var m models.ThingWithDateGSI
	if err := decodeThingWithDateGSI(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithDateGSITable) scanThingWithDateGSIs(ctx context.Context, input db.ScanThingWithDateGSIsInput, fn func(m *models.ThingWithDateGSI, lastThingWithDateGSI bool) bool) error {
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
		items, err := decodeThingWithDateGSIs(out.Items)
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

func (t ThingWithDateGSITable) deleteThingWithDateGSI(ctx context.Context, dateH strfmt.Date) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithDateGSIPrimaryKey{
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

func (t ThingWithDateGSITable) getThingWithDateGSIsByIDAndDateR(ctx context.Context, input db.GetThingWithDateGSIsByIDAndDateRInput, fn func(m *models.ThingWithDateGSI, lastThingWithDateGSI bool) bool) error {
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
		items, err := decodeThingWithDateGSIs(queryOutput.Items)
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
func (t ThingWithDateGSITable) getThingWithDateGSIsByDateHAndID(ctx context.Context, input db.GetThingWithDateGSIsByDateHAndIDInput, fn func(m *models.ThingWithDateGSI, lastThingWithDateGSI bool) bool) error {
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
		items, err := decodeThingWithDateGSIs(queryOutput.Items)
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

// encodeThingWithDateGSI encodes a ThingWithDateGSI as a DynamoDB map of attribute values.
func encodeThingWithDateGSI(m models.ThingWithDateGSI) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingWithDateGSI{
		ThingWithDateGSI: m,
	})
}

// decodeThingWithDateGSI translates a ThingWithDateGSI stored in DynamoDB to a ThingWithDateGSI struct.
func decodeThingWithDateGSI(m map[string]*dynamodb.AttributeValue, out *models.ThingWithDateGSI) error {
	var ddbThingWithDateGSI ddbThingWithDateGSI
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithDateGSI); err != nil {
		return err
	}
	*out = ddbThingWithDateGSI.ThingWithDateGSI
	return nil
}

// decodeThingWithDateGSIs translates a list of ThingWithDateGSIs stored in DynamoDB to a slice of ThingWithDateGSI structs.
func decodeThingWithDateGSIs(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithDateGSI, error) {
	thingWithDateGSIs := make([]models.ThingWithDateGSI, len(ms))
	for i, m := range ms {
		var thingWithDateGSI models.ThingWithDateGSI
		if err := decodeThingWithDateGSI(m, &thingWithDateGSI); err != nil {
			return nil, err
		}
		thingWithDateGSIs[i] = thingWithDateGSI
	}
	return thingWithDateGSIs, nil
}
