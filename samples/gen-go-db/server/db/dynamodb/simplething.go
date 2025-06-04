package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// SimpleThingTable represents the user-configurable properties of the SimpleThing table.
type SimpleThingTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbSimpleThingPrimaryKey represents the primary key of a SimpleThing in DynamoDB.
type ddbSimpleThingPrimaryKey struct {
	Name string `dynamodbav:"name"`
}

// ddbSimpleThing represents a SimpleThing as stored in DynamoDB.
type ddbSimpleThing struct {
	models.SimpleThing
}

func (t SimpleThingTable) create(ctx context.Context) error {
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

func (t SimpleThingTable) saveSimpleThing(ctx context.Context, m models.SimpleThing) error {
	data, err := encodeSimpleThing(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#NAME": "name",
		},
		ConditionExpression: aws.String("attribute_not_exists(#NAME)"),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		var conditionalCheckFailedErr *types.ConditionalCheckFailedException
		if errors.As(err, &resourceNotFoundErr) {
			return fmt.Errorf("table or index not found: %s", t.TableName)
		}
		if errors.As(err, &conditionalCheckFailedErr) {
			return db.ErrSimpleThingAlreadyExists{
				Name: m.Name,
			}
		}
		return err
	}
	return nil
}

func (t SimpleThingTable) getSimpleThing(ctx context.Context, name string) (*models.SimpleThing, error) {
	// swad-get-7
	key, err := attributevalue.MarshalMap(ddbSimpleThingPrimaryKey{
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
		return nil, db.ErrSimpleThingNotFound{
			Name: name,
		}
	}

	var m models.SimpleThing
	if err := decodeSimpleThing(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t SimpleThingTable) scanSimpleThings(ctx context.Context, input db.ScanSimpleThingsInput, fn func(m *models.SimpleThing, lastSimpleThing bool) bool) error {
	// swad-scan-1
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name": exclusiveStartKey["name"],
		}
	}
	totalRecordsProcessed := int32(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeSimpleThings(out.Items)
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

func (t SimpleThingTable) deleteSimpleThing(ctx context.Context, name string) error {

	key, err := attributevalue.MarshalMap(ddbSimpleThingPrimaryKey{
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

// encodeSimpleThing encodes a SimpleThing as a DynamoDB map of attribute values.
func encodeSimpleThing(m models.SimpleThing) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(ddbSimpleThing{
		SimpleThing: m,
	})
}

// decodeSimpleThing translates a SimpleThing stored in DynamoDB to a SimpleThing struct.
func decodeSimpleThing(m map[string]types.AttributeValue, out *models.SimpleThing) error {
	// swad-decode-1
	var ddbSimpleThing ddbSimpleThing
	if err := attributevalue.UnmarshalMap(m, &ddbSimpleThing); err != nil {
		return err
	}
	*out = ddbSimpleThing.SimpleThing
	return nil
}

// decodeSimpleThings translates a list of SimpleThings stored in DynamoDB to a slice of SimpleThing structs.
func decodeSimpleThings(ms []map[string]types.AttributeValue) ([]models.SimpleThing, error) {
	simpleThings := make([]models.SimpleThing, len(ms))
	for i, m := range ms {
		var simpleThing models.SimpleThing
		if err := decodeSimpleThing(m, &simpleThing); err != nil {
			return nil, err
		}
		simpleThings[i] = simpleThing
	}
	return simpleThings, nil
}
