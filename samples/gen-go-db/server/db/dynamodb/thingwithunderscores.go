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
var _ = errors.New("")
var _ = []types.AttributeValue{}

// ThingWithUnderscoresTable represents the user-configurable properties of the ThingWithUnderscores table.
type ThingWithUnderscoresTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithUnderscoresPrimaryKey represents the primary key of a ThingWithUnderscores in DynamoDB.
type ddbThingWithUnderscoresPrimaryKey struct {
	IDApp string `dynamodbav:"id_app"`
}

// ddbThingWithUnderscores represents a ThingWithUnderscores as stored in DynamoDB.
type ddbThingWithUnderscores struct {
	models.ThingWithUnderscores
}

func (t ThingWithUnderscoresTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id_app"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id_app"),
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

func (t ThingWithUnderscoresTable) saveThingWithUnderscores(ctx context.Context, m models.ThingWithUnderscores) error {
	data, err := encodeThingWithUnderscores(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
	})
	return err
}

func (t ThingWithUnderscoresTable) getThingWithUnderscores(ctx context.Context, iDApp string) (*models.ThingWithUnderscores, error) {
	key, err := attributevalue.MarshalMapWithOptions(ddbThingWithUnderscoresPrimaryKey{
		IDApp: iDApp,
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
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrThingWithUnderscoresNotFound{
			IDApp: iDApp,
		}
	}

	var m models.ThingWithUnderscores
	if err := decodeThingWithUnderscores(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithUnderscoresTable) deleteThingWithUnderscores(ctx context.Context, iDApp string) error {

	key, err := attributevalue.MarshalMapWithOptions(ddbThingWithUnderscoresPrimaryKey{
		IDApp: iDApp,
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
		return err
	}

	return nil
}

// encodeThingWithUnderscores encodes a ThingWithUnderscores as a DynamoDB map of attribute values.
func encodeThingWithUnderscores(m models.ThingWithUnderscores) (map[string]types.AttributeValue, error) {
	// no composite attributes, marshal the model with the json tag
	val, err := attributevalue.MarshalMapWithOptions(m, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

// decodeThingWithUnderscores translates a ThingWithUnderscores stored in DynamoDB to a ThingWithUnderscores struct.
func decodeThingWithUnderscores(m map[string]types.AttributeValue, out *models.ThingWithUnderscores) error {
	var ddbThingWithUnderscores ddbThingWithUnderscores
	if err := attributevalue.UnmarshalMapWithOptions(m, &ddbThingWithUnderscores, func(o *attributevalue.DecoderOptions) {
		o.TagKey = "json"
	}); err != nil {
		return err
	}
	*out = ddbThingWithUnderscores.ThingWithUnderscores
	return nil
}

// decodeThingWithUnderscoress translates a list of ThingWithUnderscoress stored in DynamoDB to a slice of ThingWithUnderscores structs.
func decodeThingWithUnderscoress(ms []map[string]types.AttributeValue) ([]models.ThingWithUnderscores, error) {
	thingWithUnderscoress := make([]models.ThingWithUnderscores, len(ms))
	for i, m := range ms {
		var thingWithUnderscores models.ThingWithUnderscores
		if err := decodeThingWithUnderscores(m, &thingWithUnderscores); err != nil {
			return nil, err
		}
		thingWithUnderscoress[i] = thingWithUnderscores
	}
	return thingWithUnderscoress, nil
}
