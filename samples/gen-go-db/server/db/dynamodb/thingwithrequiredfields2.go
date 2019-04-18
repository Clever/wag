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

// ThingWithRequiredFields2Table represents the user-configurable properties of the ThingWithRequiredFields2 table.
type ThingWithRequiredFields2Table struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithRequiredFields2PrimaryKey represents the primary key of a ThingWithRequiredFields2 in DynamoDB.
type ddbThingWithRequiredFields2PrimaryKey struct {
	Name string `dynamodbav:"name"`
	ID   string `dynamodbav:"id"`
}

// ddbThingWithRequiredFields2 represents a ThingWithRequiredFields2 as stored in DynamoDB.
type ddbThingWithRequiredFields2 struct {
	models.ThingWithRequiredFields2
}

func (t ThingWithRequiredFields2Table) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-required-fields2s", t.Prefix)
}

func (t ThingWithRequiredFields2Table) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
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
			{
				AttributeName: aws.String("id"),
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

func (t ThingWithRequiredFields2Table) saveThingWithRequiredFields2(ctx context.Context, m models.ThingWithRequiredFields2) error {
	data, err := encodeThingWithRequiredFields2(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
		ExpressionAttributeNames: map[string]*string{
			"#NAME": aws.String("name"),
			"#ID":   aws.String("id"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#NAME) AND attribute_not_exists(#ID)"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return db.ErrThingWithRequiredFields2AlreadyExists{
					Name: *m.Name,
					ID:   *m.ID,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	return nil
}

func (t ThingWithRequiredFields2Table) getThingWithRequiredFields2(ctx context.Context, name string, id string) (*models.ThingWithRequiredFields2, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithRequiredFields2PrimaryKey{
		Name: name,
		ID:   id,
	})
	if err != nil {
		return nil, err
	}
	res, err := t.DynamoDBAPI.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(t.name()),
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
		return nil, db.ErrThingWithRequiredFields2NotFound{
			Name: name,
			ID:   id,
		}
	}

	var m models.ThingWithRequiredFields2
	if err := decodeThingWithRequiredFields2(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithRequiredFields2Table) getThingWithRequiredFields2sByNameAndID(ctx context.Context, input db.GetThingWithRequiredFields2sByNameAndIDInput, fn func(m *models.ThingWithRequiredFields2, lastThingWithRequiredFields2 bool) bool) error {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#NAME": aws.String("name"),
			"#ID":   aws.String("id"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAt.Name),
			},
			":id": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAt.ID),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
		Limit:            input.Limit,
	}
	if input.Exclusive {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAt.ID),
			},
			"name": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAt.Name),
			},
		}
	}
	if input.Descending {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #ID <= :id")
	} else {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #ID >= :id")
	}

	queryOutput, err := t.DynamoDBAPI.QueryWithContext(ctx, queryInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	if len(queryOutput.Items) == 0 {
		return nil
	}

	items, err := decodeThingWithRequiredFields2s(queryOutput.Items)
	if err != nil {
		return err
	}

	for i, item := range items {
		hasMore := false
		if len(queryOutput.LastEvaluatedKey) > 0 {
			hasMore = true
		} else {
			hasMore = i < len(items)-1
		}
		if !fn(&item, !hasMore) {
			break
		}
	}

	return nil
}

func (t ThingWithRequiredFields2Table) deleteThingWithRequiredFields2(ctx context.Context, name string, id string) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithRequiredFields2PrimaryKey{
		Name: name,
		ID:   id,
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

// encodeThingWithRequiredFields2 encodes a ThingWithRequiredFields2 as a DynamoDB map of attribute values.
func encodeThingWithRequiredFields2(m models.ThingWithRequiredFields2) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingWithRequiredFields2{
		ThingWithRequiredFields2: m,
	})
}

// decodeThingWithRequiredFields2 translates a ThingWithRequiredFields2 stored in DynamoDB to a ThingWithRequiredFields2 struct.
func decodeThingWithRequiredFields2(m map[string]*dynamodb.AttributeValue, out *models.ThingWithRequiredFields2) error {
	var ddbThingWithRequiredFields2 ddbThingWithRequiredFields2
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithRequiredFields2); err != nil {
		return err
	}
	*out = ddbThingWithRequiredFields2.ThingWithRequiredFields2
	return nil
}

// decodeThingWithRequiredFields2s translates a list of ThingWithRequiredFields2s stored in DynamoDB to a slice of ThingWithRequiredFields2 structs.
func decodeThingWithRequiredFields2s(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithRequiredFields2, error) {
	thingWithRequiredFields2s := make([]models.ThingWithRequiredFields2, len(ms))
	for i, m := range ms {
		var thingWithRequiredFields2 models.ThingWithRequiredFields2
		if err := decodeThingWithRequiredFields2(m, &thingWithRequiredFields2); err != nil {
			return nil, err
		}
		thingWithRequiredFields2s[i] = thingWithRequiredFields2
	}
	return thingWithRequiredFields2s, nil
}
