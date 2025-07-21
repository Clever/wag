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

// ThingWithTransactionWithSimpleThingTable represents the user-configurable properties of the ThingWithTransactionWithSimpleThing table.
type ThingWithTransactionWithSimpleThingTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithTransactionWithSimpleThingPrimaryKey represents the primary key of a ThingWithTransactionWithSimpleThing in DynamoDB.
type ddbThingWithTransactionWithSimpleThingPrimaryKey struct {
	Name string `dynamodbav:"name"`
}

// ddbThingWithTransactionWithSimpleThing represents a ThingWithTransactionWithSimpleThing as stored in DynamoDB.
type ddbThingWithTransactionWithSimpleThing struct {
	models.ThingWithTransactionWithSimpleThing
}

func (t ThingWithTransactionWithSimpleThingTable) create(ctx context.Context) error {
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

func (t ThingWithTransactionWithSimpleThingTable) saveThingWithTransactionWithSimpleThing(ctx context.Context, m models.ThingWithTransactionWithSimpleThing) error {
	data, err := encodeThingWithTransactionWithSimpleThing(m)
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
			return db.ErrThingWithTransactionWithSimpleThingAlreadyExists{
				Name: m.Name,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithTransactionWithSimpleThingTable) getThingWithTransactionWithSimpleThing(ctx context.Context, name string) (*models.ThingWithTransactionWithSimpleThing, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithTransactionWithSimpleThingPrimaryKey{
		Name: name,
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
		return nil, db.ErrThingWithTransactionWithSimpleThingNotFound{
			Name: name,
		}
	}

	var m models.ThingWithTransactionWithSimpleThing
	if err := decodeThingWithTransactionWithSimpleThing(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithTransactionWithSimpleThingTable) scanThingWithTransactionWithSimpleThings(ctx context.Context, input db.ScanThingWithTransactionWithSimpleThingsInput, fn func(m *models.ThingWithTransactionWithSimpleThing, lastThingWithTransactionWithSimpleThing bool) bool) error {
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

		items, err := decodeThingWithTransactionWithSimpleThings(out.Items)
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

func (t ThingWithTransactionWithSimpleThingTable) deleteThingWithTransactionWithSimpleThing(ctx context.Context, name string) error {

	key, err := attributevalue.MarshalMap(ddbThingWithTransactionWithSimpleThingPrimaryKey{
		Name: name,
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

func (t ThingWithTransactionWithSimpleThingTable) transactSaveThingWithTransactionWithSimpleThingAndSimpleThing(ctx context.Context, m1 models.ThingWithTransactionWithSimpleThing, m1Conditions *expression.ConditionBuilder, m2 models.SimpleThing, m2Conditions *expression.ConditionBuilder) error {
	data1, err := encodeThingWithTransactionWithSimpleThing(m1)
	if err != nil {
		return err
	}

	m1CondExpr, m1ExprVals, m1ExprNames, err := buildCondExpr(m1Conditions)
	if err != nil {
		return err
	}

	data2, err := encodeSimpleThing(m2)
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
					TableName:                 aws.String(fmt.Sprintf("%s-SimpleThings", t.Prefix)),
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

// encodeThingWithTransactionWithSimpleThing encodes a ThingWithTransactionWithSimpleThing as a DynamoDB map of attribute values.
func encodeThingWithTransactionWithSimpleThing(m models.ThingWithTransactionWithSimpleThing) (map[string]types.AttributeValue, error) {
	// no composite attributes, marshal the model with the json tag
	val, err := attributevalue.MarshalMapWithOptions(m, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

// decodeThingWithTransactionWithSimpleThing translates a ThingWithTransactionWithSimpleThing stored in DynamoDB to a ThingWithTransactionWithSimpleThing struct.
func decodeThingWithTransactionWithSimpleThing(m map[string]types.AttributeValue, out *models.ThingWithTransactionWithSimpleThing) error {
	var ddbThingWithTransactionWithSimpleThing ddbThingWithTransactionWithSimpleThing
	if err := attributevalue.UnmarshalMapWithOptions(m, &ddbThingWithTransactionWithSimpleThing, func(o *attributevalue.DecoderOptions) {
		o.TagKey = "json"
	}); err != nil {
		return err
	}
	*out = ddbThingWithTransactionWithSimpleThing.ThingWithTransactionWithSimpleThing
	return nil
}

// decodeThingWithTransactionWithSimpleThings translates a list of ThingWithTransactionWithSimpleThings stored in DynamoDB to a slice of ThingWithTransactionWithSimpleThing structs.
func decodeThingWithTransactionWithSimpleThings(ms []map[string]types.AttributeValue) ([]models.ThingWithTransactionWithSimpleThing, error) {
	thingWithTransactionWithSimpleThings := make([]models.ThingWithTransactionWithSimpleThing, len(ms))
	for i, m := range ms {
		var thingWithTransactionWithSimpleThing models.ThingWithTransactionWithSimpleThing
		if err := decodeThingWithTransactionWithSimpleThing(m, &thingWithTransactionWithSimpleThing); err != nil {
			return nil, err
		}
		thingWithTransactionWithSimpleThings[i] = thingWithTransactionWithSimpleThing
	}
	return thingWithTransactionWithSimpleThings, nil
}
