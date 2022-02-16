package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Clever/wag/samples/v8/gen-go-db/models"
	"github.com/Clever/wag/samples/v8/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingAllowingBatchWritesWithCompositeAttributesTable represents the user-configurable properties of the ThingAllowingBatchWritesWithCompositeAttributes table.
type ThingAllowingBatchWritesWithCompositeAttributesTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingAllowingBatchWritesWithCompositeAttributesPrimaryKey represents the primary key of a ThingAllowingBatchWritesWithCompositeAttributes in DynamoDB.
type ddbThingAllowingBatchWritesWithCompositeAttributesPrimaryKey struct {
	NameBranch string          `dynamodbav:"name_branch"`
	Date       strfmt.DateTime `dynamodbav:"date"`
}

// ddbThingAllowingBatchWritesWithCompositeAttributes represents a ThingAllowingBatchWritesWithCompositeAttributes as stored in DynamoDB.
type ddbThingAllowingBatchWritesWithCompositeAttributes struct {
	models.ThingAllowingBatchWritesWithCompositeAttributes
}

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-allowing-batch-writes-with-composite-attributess", t.Prefix)
}

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("date"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("name_branch"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name_branch"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("date"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
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

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) saveThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, m models.ThingAllowingBatchWritesWithCompositeAttributes) error {
	data, err := encodeThingAllowingBatchWritesWithCompositeAttributes(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
		ExpressionAttributeNames: map[string]*string{
			"#NAME_BRANCH": aws.String("name_branch"),
			"#DATE":        aws.String("date"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#NAME_BRANCH) AND attribute_not_exists(#DATE)"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return db.ErrThingAllowingBatchWritesWithCompositeAttributesAlreadyExists{
					NameBranch: fmt.Sprintf("%s@%s", *m.Name, *m.Branch),
					Date:       *m.Date,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	return nil
}
func (t ThingAllowingBatchWritesWithCompositeAttributesTable) saveArrayOfThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, ms []models.ThingAllowingBatchWritesWithCompositeAttributes) error {
	if len(ms) > maxDynamoDBBatchItems {
		return fmt.Errorf("saveArrayOfThingAllowingBatchWritesWithCompositeAttributes received %d items to save, which is greater than the maximum of %d", len(ms), maxDynamoDBBatchItems)
	}

	batch := make([]*dynamodb.WriteRequest, len(ms))
	for i := range ms {
		data, err := encodeThingAllowingBatchWritesWithCompositeAttributes(ms[i])
		if err != nil {
			return err
		}
		batch[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: data,
			},
		}
	}
	tname := t.name()
	for {
		if out, err := t.DynamoDBAPI.BatchWriteItemWithContext(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				tname: batch,
			},
		}); err != nil {
			return fmt.Errorf("BatchWriteItem: %v", err)
		} else if out.UnprocessedItems != nil && len(out.UnprocessedItems[tname]) > 0 {
			batch = out.UnprocessedItems[tname]
		} else {
			break
		}
	}
	return nil
}

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) deleteArrayOfThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, ms []models.ThingAllowingBatchWritesWithCompositeAttributes) error {
	if len(ms) > maxDynamoDBBatchItems {
		return fmt.Errorf("deleteArrayOfThingAllowingBatchWritesWithCompositeAttributes received %d items to delete, which is greater than the maximum of %d", len(ms), maxDynamoDBBatchItems)
	}

	batch := make([]*dynamodb.WriteRequest, len(ms))
	for i := range ms {
		key, err := dynamodbattribute.MarshalMap(ddbThingAllowingBatchWritesWithCompositeAttributesPrimaryKey{
			NameBranch: fmt.Sprintf("%s@%s", *ms[i].Name, *ms[i].Branch),
			Date:       *ms[i].Date,
		})
		if err != nil {
			return err
		}

		batch[i] = &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: key,
			},
		}
	}
	tname := t.name()
	for {
		if out, err := t.DynamoDBAPI.BatchWriteItemWithContext(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				tname: batch,
			},
		}); err != nil {
			return fmt.Errorf("BatchWriteItem: %v", err)
		} else if out.UnprocessedItems != nil && len(out.UnprocessedItems[tname]) > 0 {
			batch = out.UnprocessedItems[tname]
		} else {
			break
		}
	}
	return nil
}

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) getThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) (*models.ThingAllowingBatchWritesWithCompositeAttributes, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingAllowingBatchWritesWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branch),
		Date:       date,
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
		return nil, db.ErrThingAllowingBatchWritesWithCompositeAttributesNotFound{
			Name:   name,
			Branch: branch,
			Date:   date,
		}
	}

	var m models.ThingAllowingBatchWritesWithCompositeAttributes
	if err := decodeThingAllowingBatchWritesWithCompositeAttributes(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) scanThingAllowingBatchWritesWithCompositeAttributess(ctx context.Context, input db.ScanThingAllowingBatchWritesWithCompositeAttributessInput, fn func(m *models.ThingAllowingBatchWritesWithCompositeAttributes, lastThingAllowingBatchWritesWithCompositeAttributes bool) bool) error {
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
			"name_branch": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s@%s", *input.StartingAfter.Name, *input.StartingAfter.Branch)),
			},
			"date": exclusiveStartKey["date"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingAllowingBatchWritesWithCompositeAttributess(out.Items)
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

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) getThingAllowingBatchWritesWithCompositeAttributessByNameBranchAndDate(ctx context.Context, input db.GetThingAllowingBatchWritesWithCompositeAttributessByNameBranchAndDateInput, fn func(m *models.ThingAllowingBatchWritesWithCompositeAttributes, lastThingAllowingBatchWritesWithCompositeAttributes bool) bool) error {
	if input.DateStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	if input.Branch == "" {
		return fmt.Errorf("Hash key input.Branch cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#NAME_BRANCH": aws.String("name_branch"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":nameBranch": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s@%s", input.Name, input.Branch)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME_BRANCH = :nameBranch")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = aws.String("date")
		queryInput.ExpressionAttributeValues[":date"] = &dynamodb.AttributeValue{
			S: aws.String(toDynamoTimeString(*input.DateStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME_BRANCH = :nameBranch AND #DATE <= :date")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME_BRANCH = :nameBranch AND #DATE >= :date")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"date": &dynamodb.AttributeValue{
				S: aws.String(toDynamoTimeStringPtr(input.StartingAfter.Date)),
			},
			"name_branch": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s@%s", *input.StartingAfter.Name, *input.StartingAfter.Branch)),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingAllowingBatchWritesWithCompositeAttributess(queryOutput.Items)
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

func (t ThingAllowingBatchWritesWithCompositeAttributesTable) deleteThingAllowingBatchWritesWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingAllowingBatchWritesWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branch),
		Date:       date,
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

// encodeThingAllowingBatchWritesWithCompositeAttributes encodes a ThingAllowingBatchWritesWithCompositeAttributes as a DynamoDB map of attribute values.
func encodeThingAllowingBatchWritesWithCompositeAttributes(m models.ThingAllowingBatchWritesWithCompositeAttributes) (map[string]*dynamodb.AttributeValue, error) {
	val, err := dynamodbattribute.MarshalMap(ddbThingAllowingBatchWritesWithCompositeAttributes{
		ThingAllowingBatchWritesWithCompositeAttributes: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(*m.Branch, "@") {
		return nil, fmt.Errorf("branch cannot contain '@': %s", *m.Branch)
	}
	if strings.Contains(*m.Name, "@") {
		return nil, fmt.Errorf("name cannot contain '@': %s", *m.Name)
	}
	// add in composite attributes
	primaryKey, err := dynamodbattribute.MarshalMap(ddbThingAllowingBatchWritesWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", *m.Name, *m.Branch),
		Date:       *m.Date,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range primaryKey {
		val[k] = v
	}
	return val, err
}

// decodeThingAllowingBatchWritesWithCompositeAttributes translates a ThingAllowingBatchWritesWithCompositeAttributes stored in DynamoDB to a ThingAllowingBatchWritesWithCompositeAttributes struct.
func decodeThingAllowingBatchWritesWithCompositeAttributes(m map[string]*dynamodb.AttributeValue, out *models.ThingAllowingBatchWritesWithCompositeAttributes) error {
	var ddbThingAllowingBatchWritesWithCompositeAttributes ddbThingAllowingBatchWritesWithCompositeAttributes
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingAllowingBatchWritesWithCompositeAttributes); err != nil {
		return err
	}
	*out = ddbThingAllowingBatchWritesWithCompositeAttributes.ThingAllowingBatchWritesWithCompositeAttributes
	return nil
}

// decodeThingAllowingBatchWritesWithCompositeAttributess translates a list of ThingAllowingBatchWritesWithCompositeAttributess stored in DynamoDB to a slice of ThingAllowingBatchWritesWithCompositeAttributes structs.
func decodeThingAllowingBatchWritesWithCompositeAttributess(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingAllowingBatchWritesWithCompositeAttributes, error) {
	thingAllowingBatchWritesWithCompositeAttributess := make([]models.ThingAllowingBatchWritesWithCompositeAttributes, len(ms))
	for i, m := range ms {
		var thingAllowingBatchWritesWithCompositeAttributes models.ThingAllowingBatchWritesWithCompositeAttributes
		if err := decodeThingAllowingBatchWritesWithCompositeAttributes(m, &thingAllowingBatchWritesWithCompositeAttributes); err != nil {
			return nil, err
		}
		thingAllowingBatchWritesWithCompositeAttributess[i] = thingAllowingBatchWritesWithCompositeAttributes
	}
	return thingAllowingBatchWritesWithCompositeAttributess, nil
}
