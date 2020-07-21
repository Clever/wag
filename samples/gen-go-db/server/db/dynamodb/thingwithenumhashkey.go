package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db/models"
	"github.com/Clever/wag/samples/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithEnumHashKeyTable represents the user-configurable properties of the ThingWithEnumHashKey table.
type ThingWithEnumHashKeyTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithEnumHashKeyPrimaryKey represents the primary key of a ThingWithEnumHashKey in DynamoDB.
type ddbThingWithEnumHashKeyPrimaryKey struct {
	Branch models.Branch   `dynamodbav:"branch"`
	Date   strfmt.DateTime `dynamodbav:"date"`
}

// ddbThingWithEnumHashKeyGSIByBranch represents the byBranch GSI.
type ddbThingWithEnumHashKeyGSIByBranch struct {
	Branch models.Branch   `dynamodbav:"branch"`
	Date2  strfmt.DateTime `dynamodbav:"date2"`
}

// ddbThingWithEnumHashKey represents a ThingWithEnumHashKey as stored in DynamoDB.
type ddbThingWithEnumHashKey struct {
	models.ThingWithEnumHashKey
}

func (t ThingWithEnumHashKeyTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-enum-hash-keys", t.Prefix)
}

func (t ThingWithEnumHashKeyTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("branch"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("date"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("date2"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("branch"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("date"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("byBranch"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("branch"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("date2"),
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

func (t ThingWithEnumHashKeyTable) saveThingWithEnumHashKey(ctx context.Context, m models.ThingWithEnumHashKey) error {
	data, err := encodeThingWithEnumHashKey(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
		ExpressionAttributeNames: map[string]*string{
			"#BRANCH": aws.String("branch"),
			"#DATE":   aws.String("date"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#BRANCH) AND attribute_not_exists(#DATE)"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return db.ErrThingWithEnumHashKeyAlreadyExists{
					Branch: m.Branch,
					Date:   m.Date,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	return nil
}

func (t ThingWithEnumHashKeyTable) getThingWithEnumHashKey(ctx context.Context, branch models.Branch, date strfmt.DateTime) (*models.ThingWithEnumHashKey, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithEnumHashKeyPrimaryKey{
		Branch: branch,
		Date:   date,
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
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return nil, fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrThingWithEnumHashKeyNotFound{
			Branch: branch,
			Date:   date,
		}
	}

	var m models.ThingWithEnumHashKey
	if err := decodeThingWithEnumHashKey(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithEnumHashKeyTable) getThingWithEnumHashKeysByBranchAndDate(ctx context.Context, input db.GetThingWithEnumHashKeysByBranchAndDateInput, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error {
	if input.DateStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateStartingAt or input.StartingAfter")
	}
	if input.Branch == "" {
		return fmt.Errorf("Hash key input.Branch cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#BRANCH": aws.String("branch"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":branch": &dynamodb.AttributeValue{
				S: aws.String(string(input.Branch)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = aws.String("date")
		queryInput.ExpressionAttributeValues[":date"] = &dynamodb.AttributeValue{
			S: aws.String(toDynamoTimeString(*input.DateStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch AND #DATE <= :date")
		} else {
			queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch AND #DATE >= :date")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"date": &dynamodb.AttributeValue{
				S: aws.String(toDynamoTimeString(input.StartingAfter.Date)),
			},
			"branch": &dynamodb.AttributeValue{
				S: aws.String(string(input.StartingAfter.Branch)),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithEnumHashKeys(queryOutput.Items)
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
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}

func (t ThingWithEnumHashKeyTable) deleteThingWithEnumHashKey(ctx context.Context, branch models.Branch, date strfmt.DateTime) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithEnumHashKeyPrimaryKey{
		Branch: branch,
		Date:   date,
	})
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.DeleteItemWithContext(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(t.name()),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}

	return nil
}

func (t ThingWithEnumHashKeyTable) getThingWithEnumHashKeysByBranchAndDate2(ctx context.Context, input db.GetThingWithEnumHashKeysByBranchAndDate2Input, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error {
	if input.Date2StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.Date2StartingAt or input.StartingAfter")
	}
	if input.Branch == "" {
		return fmt.Errorf("Hash key input.Branch cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("byBranch"),
		ExpressionAttributeNames: map[string]*string{
			"#BRANCH": aws.String("branch"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":branch": &dynamodb.AttributeValue{
				S: aws.String(string(input.Branch)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.Date2StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch")
	} else {
		queryInput.ExpressionAttributeNames["#DATE2"] = aws.String("date2")
		queryInput.ExpressionAttributeValues[":date2"] = &dynamodb.AttributeValue{
			S: aws.String(toDynamoTimeString(*input.Date2StartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch AND #DATE2 <= :date2")
		} else {
			queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch AND #DATE2 >= :date2")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"date2": &dynamodb.AttributeValue{
				S: aws.String(toDynamoTimeString(input.StartingAfter.Date2)),
			},
			"branch": &dynamodb.AttributeValue{
				S: aws.String(string(input.StartingAfter.Branch)),
			},
			"date": &dynamodb.AttributeValue{
				S: aws.String(toDynamoTimeString(input.StartingAfter.Date)),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithEnumHashKeys(queryOutput.Items)
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
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}

// encodeThingWithEnumHashKey encodes a ThingWithEnumHashKey as a DynamoDB map of attribute values.
func encodeThingWithEnumHashKey(m models.ThingWithEnumHashKey) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingWithEnumHashKey{
		ThingWithEnumHashKey: m,
	})
}

// decodeThingWithEnumHashKey translates a ThingWithEnumHashKey stored in DynamoDB to a ThingWithEnumHashKey struct.
func decodeThingWithEnumHashKey(m map[string]*dynamodb.AttributeValue, out *models.ThingWithEnumHashKey) error {
	var ddbThingWithEnumHashKey ddbThingWithEnumHashKey
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithEnumHashKey); err != nil {
		return err
	}
	*out = ddbThingWithEnumHashKey.ThingWithEnumHashKey
	return nil
}

// decodeThingWithEnumHashKeys translates a list of ThingWithEnumHashKeys stored in DynamoDB to a slice of ThingWithEnumHashKey structs.
func decodeThingWithEnumHashKeys(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithEnumHashKey, error) {
	thingWithEnumHashKeys := make([]models.ThingWithEnumHashKey, len(ms))
	for i, m := range ms {
		var thingWithEnumHashKey models.ThingWithEnumHashKey
		if err := decodeThingWithEnumHashKey(m, &thingWithEnumHashKey); err != nil {
			return nil, err
		}
		thingWithEnumHashKeys[i] = thingWithEnumHashKey
	}
	return thingWithEnumHashKeys, nil
}
