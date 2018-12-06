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
)

// ThingWithRequiredFieldsTable represents the user-configurable properties of the ThingWithRequiredFields table.
type ThingWithRequiredFieldsTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithRequiredFieldsPrimaryKey represents the primary key of a ThingWithRequiredFields in DynamoDB.
type ddbThingWithRequiredFieldsPrimaryKey struct {
	Name string `dynamodbav:"name"`
}

// ddbThingWithRequiredFields represents a ThingWithRequiredFields as stored in DynamoDB.
type ddbThingWithRequiredFields struct {
	models.ThingWithRequiredFields
}

func (t ThingWithRequiredFieldsTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-required-fieldss", t.Prefix)
}

func (t ThingWithRequiredFieldsTable) create(ctx context.Context) error {
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

func (t ThingWithRequiredFieldsTable) saveThingWithRequiredFields(ctx context.Context, m models.ThingWithRequiredFields) error {
	data, err := encodeThingWithRequiredFields(m)
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
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return db.ErrThingWithRequiredFieldsAlreadyExists{
					Name: *m.Name,
				}
			}
		}
		return err
	}
	return nil
}

func (t ThingWithRequiredFieldsTable) getThingWithRequiredFields(ctx context.Context, name string) (*models.ThingWithRequiredFields, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithRequiredFieldsPrimaryKey{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	res, err := t.DynamoDBAPI.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(t.name()),
	})
	if err != nil {
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrThingWithRequiredFieldsNotFound{
			Name: name,
		}
	}

	var m models.ThingWithRequiredFields
	if err := decodeThingWithRequiredFields(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithRequiredFieldsTable) deleteThingWithRequiredFields(ctx context.Context, name string) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithRequiredFieldsPrimaryKey{
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
		return err
	}
	return nil
}

// encodeThingWithRequiredFields encodes a ThingWithRequiredFields as a DynamoDB map of attribute values.
func encodeThingWithRequiredFields(m models.ThingWithRequiredFields) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingWithRequiredFields{
		ThingWithRequiredFields: m,
	})
}

// decodeThingWithRequiredFields translates a ThingWithRequiredFields stored in DynamoDB to a ThingWithRequiredFields struct.
func decodeThingWithRequiredFields(m map[string]*dynamodb.AttributeValue, out *models.ThingWithRequiredFields) error {
	var ddbThingWithRequiredFields ddbThingWithRequiredFields
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithRequiredFields); err != nil {
		return err
	}
	*out = ddbThingWithRequiredFields.ThingWithRequiredFields
	return nil
}

// decodeThingWithRequiredFieldss translates a list of ThingWithRequiredFieldss stored in DynamoDB to a slice of ThingWithRequiredFields structs.
func decodeThingWithRequiredFieldss(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithRequiredFields, error) {
	thingWithRequiredFieldss := make([]models.ThingWithRequiredFields, len(ms))
	for i, m := range ms {
		var thingWithRequiredFields models.ThingWithRequiredFields
		if err := decodeThingWithRequiredFields(m, &thingWithRequiredFields); err != nil {
			return nil, err
		}
		thingWithRequiredFieldss[i] = thingWithRequiredFields
	}
	return thingWithRequiredFieldss, nil
}