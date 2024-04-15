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
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithTransactionTable represents the user-configurable properties of the ThingWithTransaction table.
type ThingWithTransactionTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithTransactionPrimaryKey represents the primary key of a ThingWithTransaction in DynamoDB.
type ddbThingWithTransactionPrimaryKey struct {
	Name string `dynamodbav:"name"`
}

// ddbThingWithTransaction represents a ThingWithTransaction as stored in DynamoDB.
type ddbThingWithTransaction struct {
	models.ThingWithTransaction
}

func (t ThingWithTransactionTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
		},
		TableName: aws.String(t.TableName),
	}); err != nil {
		return err
	}
	return nil
}

func (t ThingWithTransactionTable) saveThingWithTransaction(ctx context.Context, m models.ThingWithTransaction) error {
	data, err := encodeThingWithTransaction(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]*string{
			"#NAME": aws.String("name"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#NAME)"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return db.ErrThingWithTransactionAlreadyExists{
					Name: m.Name,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
		}
		return err
	}
	return nil
}

func (t ThingWithTransactionTable) getThingWithTransaction(ctx context.Context, name string) (*models.ThingWithTransaction, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithTransactionPrimaryKey{
		Name: name,
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
		return nil, db.ErrThingWithTransactionNotFound{
			Name: name,
		}
	}

	var m models.ThingWithTransaction
	if err := decodeThingWithTransaction(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithTransactionTable) scanThingWithTransactions(ctx context.Context, input db.ScanThingWithTransactionsInput, fn func(m *models.ThingWithTransaction, lastThingWithTransaction bool) bool) error {
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
			"name": exclusiveStartKey["name"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithTransactions(out.Items)
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

func (t ThingWithTransactionTable) deleteThingWithTransaction(ctx context.Context, name string) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithTransactionPrimaryKey{
		Name: name,
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

func (t ThingWithTransactionTable) transactSaveThingWithTransactionAndThing(ctx context.Context, m1 models.ThingWithTransaction, m1Conditions *expression.ConditionBuilder, m2 models.Thing, m2Conditions *expression.ConditionBuilder) error {
	data1, err := encodeThingWithTransaction(m1)
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

// encodeThingWithTransaction encodes a ThingWithTransaction as a DynamoDB map of attribute values.
func encodeThingWithTransaction(m models.ThingWithTransaction) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingWithTransaction{
		ThingWithTransaction: m,
	})
}

// decodeThingWithTransaction translates a ThingWithTransaction stored in DynamoDB to a ThingWithTransaction struct.
func decodeThingWithTransaction(m map[string]*dynamodb.AttributeValue, out *models.ThingWithTransaction) error {
	var ddbThingWithTransaction ddbThingWithTransaction
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithTransaction); err != nil {
		return err
	}
	*out = ddbThingWithTransaction.ThingWithTransaction
	return nil
}

// decodeThingWithTransactions translates a list of ThingWithTransactions stored in DynamoDB to a slice of ThingWithTransaction structs.
func decodeThingWithTransactions(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithTransaction, error) {
	thingWithTransactions := make([]models.ThingWithTransaction, len(ms))
	for i, m := range ms {
		var thingWithTransaction models.ThingWithTransaction
		if err := decodeThingWithTransaction(m, &thingWithTransaction); err != nil {
			return nil, err
		}
		thingWithTransactions[i] = thingWithTransaction
	}
	return thingWithTransactions, nil
}
