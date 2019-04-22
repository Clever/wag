package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Clever/wag/samples/gen-go-db/models"
	"github.com/Clever/wag/samples/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithDateTimeCompositeTable represents the user-configurable properties of the ThingWithDateTimeComposite table.
type ThingWithDateTimeCompositeTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithDateTimeCompositePrimaryKey represents the primary key of a ThingWithDateTimeComposite in DynamoDB.
type ddbThingWithDateTimeCompositePrimaryKey struct {
	TypeID          string `dynamodbav:"typeID"`
	CreatedResource string `dynamodbav:"createdResource"`
}

// ddbThingWithDateTimeComposite represents a ThingWithDateTimeComposite as stored in DynamoDB.
type ddbThingWithDateTimeComposite struct {
	models.ThingWithDateTimeComposite
}

func (t ThingWithDateTimeCompositeTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-date-time-composites", t.Prefix)
}

func (t ThingWithDateTimeCompositeTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("createdResource"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("typeID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("typeID"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("createdResource"),
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

func (t ThingWithDateTimeCompositeTable) saveThingWithDateTimeComposite(ctx context.Context, m models.ThingWithDateTimeComposite) error {
	data, err := encodeThingWithDateTimeComposite(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
	})
	return err
}

func (t ThingWithDateTimeCompositeTable) getThingWithDateTimeComposite(ctx context.Context, typeVar string, id string, created strfmt.DateTime, resource string) (*models.ThingWithDateTimeComposite, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithDateTimeCompositePrimaryKey{
		TypeID:          fmt.Sprintf("%s|%s", typeVar, id),
		CreatedResource: fmt.Sprintf("%s|%s", created, resource),
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
		return nil, db.ErrThingWithDateTimeCompositeNotFound{
			Type:     typeVar,
			ID:       id,
			Created:  created,
			Resource: resource,
		}
	}

	var m models.ThingWithDateTimeComposite
	if err := decodeThingWithDateTimeComposite(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithDateTimeCompositeTable) getThingWithDateTimeCompositesByTypeIDAndCreatedResource(ctx context.Context, input db.GetThingWithDateTimeCompositesByTypeIDAndCreatedResourceInput, fn func(m *models.ThingWithDateTimeComposite, lastThingWithDateTimeComposite bool) bool) error {
	if input.StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of StartingAt or StartingAfter")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#TYPEID": aws.String("typeID"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":typeId": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s|%s", input.Type, input.ID)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#TYPEID = :typeId")
	} else {
		queryInput.ExpressionAttributeNames["#CREATEDRESOURCE"] = aws.String("createdResource")
		queryInput.ExpressionAttributeValues[":createdResource"] = &dynamodb.AttributeValue{
			S: aws.String(fmt.Sprintf("%s|%s", input.StartingAt.Created, input.StartingAt.Resource)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#TYPEID = :typeId AND #CREATEDRESOURCE <= :createdResource")
		} else {
			queryInput.KeyConditionExpression = aws.String("#TYPEID = :typeId AND #CREATEDRESOURCE >= :createdResource")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"createdResource": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s|%s", input.StartingAfter.Created, input.StartingAfter.Resource)),
			},
			"typeID": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s|%s", input.StartingAfter.Type, input.StartingAfter.ID)),
			},
		}
	}

	var decodeErr error
	var items []models.ThingWithDateTimeComposite
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, decodeErr = decodeThingWithDateTimeComposites(queryOutput.Items)
		if decodeErr != nil {
			return false
		}
		hasMore := true
		for i, item := range items {
			if lastPage == true {
				hasMore = i < len(items)-1
			}
			if !fn(&item, !hasMore) {
				return false
			}
		}
		return true
	}

	err := t.DynamoDBAPI.QueryPagesWithContext(ctx, queryInput, pageFn)
	if err != nil {
		return err
	}
	if decodeErr != nil {
		return decodeErr
	}

	return nil
}

func (t ThingWithDateTimeCompositeTable) deleteThingWithDateTimeComposite(ctx context.Context, typeVar string, id string, created strfmt.DateTime, resource string) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithDateTimeCompositePrimaryKey{
		TypeID:          fmt.Sprintf("%s|%s", typeVar, id),
		CreatedResource: fmt.Sprintf("%s|%s", created, resource),
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

// encodeThingWithDateTimeComposite encodes a ThingWithDateTimeComposite as a DynamoDB map of attribute values.
func encodeThingWithDateTimeComposite(m models.ThingWithDateTimeComposite) (map[string]*dynamodb.AttributeValue, error) {
	val, err := dynamodbattribute.MarshalMap(ddbThingWithDateTimeComposite{
		ThingWithDateTimeComposite: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(m.ID, "|") {
		return nil, fmt.Errorf("id cannot contain '|': %s", m.ID)
	}
	if strings.Contains(m.Resource, "|") {
		return nil, fmt.Errorf("resource cannot contain '|': %s", m.Resource)
	}
	if strings.Contains(m.Type, "|") {
		return nil, fmt.Errorf("type cannot contain '|': %s", m.Type)
	}
	// add in composite attributes
	primaryKey, err := dynamodbattribute.MarshalMap(ddbThingWithDateTimeCompositePrimaryKey{
		TypeID:          fmt.Sprintf("%s|%s", m.Type, m.ID),
		CreatedResource: fmt.Sprintf("%s|%s", m.Created, m.Resource),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range primaryKey {
		val[k] = v
	}
	return val, err
}

// decodeThingWithDateTimeComposite translates a ThingWithDateTimeComposite stored in DynamoDB to a ThingWithDateTimeComposite struct.
func decodeThingWithDateTimeComposite(m map[string]*dynamodb.AttributeValue, out *models.ThingWithDateTimeComposite) error {
	var ddbThingWithDateTimeComposite ddbThingWithDateTimeComposite
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithDateTimeComposite); err != nil {
		return err
	}
	*out = ddbThingWithDateTimeComposite.ThingWithDateTimeComposite
	return nil
}

// decodeThingWithDateTimeComposites translates a list of ThingWithDateTimeComposites stored in DynamoDB to a slice of ThingWithDateTimeComposite structs.
func decodeThingWithDateTimeComposites(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithDateTimeComposite, error) {
	thingWithDateTimeComposites := make([]models.ThingWithDateTimeComposite, len(ms))
	for i, m := range ms {
		var thingWithDateTimeComposite models.ThingWithDateTimeComposite
		if err := decodeThingWithDateTimeComposite(m, &thingWithDateTimeComposite); err != nil {
			return nil, err
		}
		thingWithDateTimeComposites[i] = thingWithDateTimeComposite
	}
	return thingWithDateTimeComposites, nil
}
