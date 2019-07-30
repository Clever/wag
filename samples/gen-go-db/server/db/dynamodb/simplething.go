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

// SimpleThingTable represents the user-configurable properties of the SimpleThing table.
type SimpleThingTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
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

func (t SimpleThingTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-simple-things", t.Prefix)
}

func (t SimpleThingTable) create(ctx context.Context) error {
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
		TableName: aws.String(t.name()),
	}); err != nil {
		return err
	}
	return nil
}

func (t SimpleThingTable) saveSimpleThing(ctx context.Context, m models.SimpleThing) error {
	data, err := encodeSimpleThing(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
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
				return db.ErrSimpleThingAlreadyExists{
					Name: m.Name,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	return nil
}

func (t SimpleThingTable) getSimpleThing(ctx context.Context, name string) (*models.SimpleThing, error) {
	key, err := dynamodbattribute.MarshalMap(ddbSimpleThingPrimaryKey{
		Name: name,
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

func (t SimpleThingTable) deleteSimpleThing(ctx context.Context, name string) error {
	key, err := dynamodbattribute.MarshalMap(ddbSimpleThingPrimaryKey{
		Name: name,
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

// encodeSimpleThing encodes a SimpleThing as a DynamoDB map of attribute values.
func encodeSimpleThing(m models.SimpleThing) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbSimpleThing{
		SimpleThing: m,
	})
}

// decodeSimpleThing translates a SimpleThing stored in DynamoDB to a SimpleThing struct.
func decodeSimpleThing(m map[string]*dynamodb.AttributeValue, out *models.SimpleThing) error {
	var ddbSimpleThing ddbSimpleThing
	if err := dynamodbattribute.UnmarshalMap(m, &ddbSimpleThing); err != nil {
		return err
	}
	*out = ddbSimpleThing.SimpleThing
	return nil
}

// decodeSimpleThings translates a list of SimpleThings stored in DynamoDB to a slice of SimpleThing structs.
func decodeSimpleThings(ms []map[string]*dynamodb.AttributeValue) ([]models.SimpleThing, error) {
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
