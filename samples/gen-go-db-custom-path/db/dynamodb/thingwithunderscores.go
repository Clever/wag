package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/v8/samples/gen-go-db-custom-path/db"
	"github.com/Clever/wag/v8/samples/gen-go-db-custom-path/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithUnderscoresTable represents the user-configurable properties of the ThingWithUnderscores table.
type ThingWithUnderscoresTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
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

func (t ThingWithUnderscoresTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-underscoress", t.Prefix)
}

func (t ThingWithUnderscoresTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id_app"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id_app"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
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

func (t ThingWithUnderscoresTable) saveThingWithUnderscores(ctx context.Context, m models.ThingWithUnderscores) error {
	data, err := encodeThingWithUnderscores(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
	})
	return err
}

func (t ThingWithUnderscoresTable) getThingWithUnderscores(ctx context.Context, iDApp string) (*models.ThingWithUnderscores, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithUnderscoresPrimaryKey{
		IDApp: iDApp,
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
	key, err := dynamodbattribute.MarshalMap(ddbThingWithUnderscoresPrimaryKey{
		IDApp: iDApp,
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

// encodeThingWithUnderscores encodes a ThingWithUnderscores as a DynamoDB map of attribute values.
func encodeThingWithUnderscores(m models.ThingWithUnderscores) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingWithUnderscores{
		ThingWithUnderscores: m,
	})
}

// decodeThingWithUnderscores translates a ThingWithUnderscores stored in DynamoDB to a ThingWithUnderscores struct.
func decodeThingWithUnderscores(m map[string]*dynamodb.AttributeValue, out *models.ThingWithUnderscores) error {
	var ddbThingWithUnderscores ddbThingWithUnderscores
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithUnderscores); err != nil {
		return err
	}
	*out = ddbThingWithUnderscores.ThingWithUnderscores
	return nil
}

// decodeThingWithUnderscoress translates a list of ThingWithUnderscoress stored in DynamoDB to a slice of ThingWithUnderscores structs.
func decodeThingWithUnderscoress(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithUnderscores, error) {
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
