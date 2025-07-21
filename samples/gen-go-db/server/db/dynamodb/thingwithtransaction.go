package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}
var _ = errors.New("")
var _ = []types.AttributeValue{}

// ThingWithTransactionTable represents the user-configurable properties of the ThingWithTransaction table.
type ThingWithTransactionTable struct {
	DynamoDBAPI        *dynamodb.Client
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
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       types.KeyTypeHash,
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

func (t ThingWithTransactionTable) saveThingWithTransaction(ctx context.Context, m models.ThingWithTransaction) error {
	data, err := encodeThingWithTransaction(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#NAME": "name",
		},
		ConditionExpression: aws.String(
			"" +
				"" +
				"attribute_not_exists(#NAME)" +
				"",
		),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		var conditionalCheckFailedErr *types.ConditionalCheckFailedException
		if errors.As(err, &resourceNotFoundErr) {
			return fmt.Errorf("table or index not found: %s", t.TableName)
		}
		if errors.As(err, &conditionalCheckFailedErr) {
			return db.ErrThingWithTransactionAlreadyExists{
				Name: m.Name,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithTransactionTable) getThingWithTransaction(ctx context.Context, name string) (*models.ThingWithTransaction, error) {
	key, err := attributevalue.MarshalMapWithOptions(ddbThingWithTransactionPrimaryKey{
		Name: name,
	}, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
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
		var resourceNotFoundErr *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundErr) {
			return nil, fmt.Errorf("table or index not found: %s", t.TableName)
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
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMapWithOptions(input.StartingAfter, func(o *attributevalue.EncoderOptions) {
			o.TagKey = "json"
		})
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name": exclusiveStartKey["name"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithTransactions(out.Items)
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

func (t ThingWithTransactionTable) deleteThingWithTransaction(ctx context.Context, name string) error {

	key, err := attributevalue.MarshalMapWithOptions(ddbThingWithTransactionPrimaryKey{
		Name: name,
	}, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
	})
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(t.TableName),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundErr) {
			return fmt.Errorf("table or index not found: %s", t.TableName)
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

	// Convert map[string]*string to map[string]string for ExpressionAttributeNames
	toStringMap := func(in map[string]*string) map[string]string {
		if in == nil {
			return nil
		}
		out := make(map[string]string, len(in))
		for k, v := range in {
			if v != nil {
				out[k] = *v
			}
		}
		return out
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName:                 aws.String(t.TableName),
					Item:                      data1,
					ConditionExpression:       m1CondExpr,
					ExpressionAttributeValues: m1ExprVals,
					ExpressionAttributeNames:  toStringMap(m1ExprNames),
				},
			},
			{
				Put: &types.Put{
					TableName:                 aws.String(fmt.Sprintf("%s-Things", t.Prefix)),
					Item:                      data2,
					ConditionExpression:       m2CondExpr,
					ExpressionAttributeValues: m2ExprVals,
					ExpressionAttributeNames:  toStringMap(m2ExprNames),
				},
			},
		},
	}
	_, err = t.DynamoDBAPI.TransactWriteItems(ctx, input)

	return err
}

// encodeThingWithTransaction encodes a ThingWithTransaction as a DynamoDB map of attribute values.
func encodeThingWithTransaction(m models.ThingWithTransaction) (map[string]types.AttributeValue, error) {
	// no composite attributes, marshal the model with the json tag
	val, err := attributevalue.MarshalMapWithOptions(m, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

// decodeThingWithTransaction translates a ThingWithTransaction stored in DynamoDB to a ThingWithTransaction struct.
func decodeThingWithTransaction(m map[string]types.AttributeValue, out *models.ThingWithTransaction) error {
	var ddbThingWithTransaction ddbThingWithTransaction
	if err := attributevalue.UnmarshalMapWithOptions(m, &ddbThingWithTransaction, func(o *attributevalue.DecoderOptions) {
		o.TagKey = "json"
	}); err != nil {
		return err
	}
	*out = ddbThingWithTransaction.ThingWithTransaction
	return nil
}

// decodeThingWithTransactions translates a list of ThingWithTransactions stored in DynamoDB to a slice of ThingWithTransaction structs.
func decodeThingWithTransactions(ms []map[string]types.AttributeValue) ([]models.ThingWithTransaction, error) {
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
