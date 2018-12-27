package dynamodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Clever/wag/samples/gen-go-db/models"
	"github.com/Clever/wag/samples/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

// ThingWithCompositeEnumAttributesTable represents the user-configurable properties of the ThingWithCompositeEnumAttributes table.
type ThingWithCompositeEnumAttributesTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithCompositeEnumAttributesPrimaryKey represents the primary key of a ThingWithCompositeEnumAttributes in DynamoDB.
type ddbThingWithCompositeEnumAttributesPrimaryKey struct {
	NameBranch string          `dynamodbav:"name_branch"`
	Date       strfmt.DateTime `dynamodbav:"date"`
}

// ddbThingWithCompositeEnumAttributes represents a ThingWithCompositeEnumAttributes as stored in DynamoDB.
type ddbThingWithCompositeEnumAttributes struct {
	models.ThingWithCompositeEnumAttributes
}

func (t ThingWithCompositeEnumAttributesTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-composite-enum-attributess", t.Prefix)
}

func (t ThingWithCompositeEnumAttributesTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("date"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("name_branch"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name_branch"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("date"),
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

func (t ThingWithCompositeEnumAttributesTable) saveThingWithCompositeEnumAttributes(ctx context.Context, m models.ThingWithCompositeEnumAttributes) error {
	data, err := encodeThingWithCompositeEnumAttributes(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
		ExpressionAttributeNames: map[string]*string{
			"#NAME_BRANCH": aws.String("name_branch"),
			"#DATE":        aws.String("date"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#NAME_BRANCH) AND attribute_not_exists(#DATE)"),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				return db.ErrThingWithCompositeEnumAttributesAlreadyExists{
					NameBranch: fmt.Sprintf("%s@%s", *m.Name, m.BranchID),
					Date:       *m.Date,
				}
			}
		}
		return err
	}
	return nil
}

func (t ThingWithCompositeEnumAttributesTable) getThingWithCompositeEnumAttributes(ctx context.Context, name string, branchID models.Branch, date strfmt.DateTime) (*models.ThingWithCompositeEnumAttributes, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeEnumAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branchID),
		Date:       date,
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
		return nil, db.ErrThingWithCompositeEnumAttributesNotFound{
			Name:     name,
			BranchID: branchID,
			Date:     date,
		}
	}

	var m models.ThingWithCompositeEnumAttributes
	if err := decodeThingWithCompositeEnumAttributes(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithCompositeEnumAttributesTable) getThingWithCompositeEnumAttributessByNameBranchAndDate(ctx context.Context, input db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput) ([]models.ThingWithCompositeEnumAttributes, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#NAME_BRANCH": aws.String("name_branch"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":nameBranch": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s@%s", input.Name, input.BranchID)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME_BRANCH = :nameBranch")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = aws.String("date")
		queryInput.ExpressionAttributeValues[":date"] = &dynamodb.AttributeValue{
			S: aws.String(time.Time(*input.DateStartingAt).Format(time.RFC3339)), // dynamodb attributevalue only supports RFC3339 resolution
		}
		queryInput.KeyConditionExpression = aws.String("#NAME_BRANCH = :nameBranch AND #DATE >= :date")
	}

	queryOutput, err := t.DynamoDBAPI.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if len(queryOutput.Items) == 0 {
		return []models.ThingWithCompositeEnumAttributes{}, nil
	}

	return decodeThingWithCompositeEnumAttributess(queryOutput.Items)
}

func (t ThingWithCompositeEnumAttributesTable) deleteThingWithCompositeEnumAttributes(ctx context.Context, name string, branchID models.Branch, date strfmt.DateTime) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeEnumAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branchID),
		Date:       date,
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

// encodeThingWithCompositeEnumAttributes encodes a ThingWithCompositeEnumAttributes as a DynamoDB map of attribute values.
func encodeThingWithCompositeEnumAttributes(m models.ThingWithCompositeEnumAttributes) (map[string]*dynamodb.AttributeValue, error) {
	val, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeEnumAttributes{
		ThingWithCompositeEnumAttributes: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(*m.Name, "@") {
		return nil, fmt.Errorf("name cannot contain '@': %s", *m.Name)
	}
	// add in composite attributes
	primaryKey, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeEnumAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", *m.Name, m.BranchID),
		Date:       *m.Date,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range primaryKey {
		val[k] = v
	}
	return val, err
}

// decodeThingWithCompositeEnumAttributes translates a ThingWithCompositeEnumAttributes stored in DynamoDB to a ThingWithCompositeEnumAttributes struct.
func decodeThingWithCompositeEnumAttributes(m map[string]*dynamodb.AttributeValue, out *models.ThingWithCompositeEnumAttributes) error {
	var ddbThingWithCompositeEnumAttributes ddbThingWithCompositeEnumAttributes
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithCompositeEnumAttributes); err != nil {
		return err
	}
	*out = ddbThingWithCompositeEnumAttributes.ThingWithCompositeEnumAttributes
	return nil
}

// decodeThingWithCompositeEnumAttributess translates a list of ThingWithCompositeEnumAttributess stored in DynamoDB to a slice of ThingWithCompositeEnumAttributes structs.
func decodeThingWithCompositeEnumAttributess(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithCompositeEnumAttributes, error) {
	thingWithCompositeEnumAttributess := make([]models.ThingWithCompositeEnumAttributes, len(ms))
	for i, m := range ms {
		var thingWithCompositeEnumAttributes models.ThingWithCompositeEnumAttributes
		if err := decodeThingWithCompositeEnumAttributes(m, &thingWithCompositeEnumAttributes); err != nil {
			return nil, err
		}
		thingWithCompositeEnumAttributess[i] = thingWithCompositeEnumAttributes
	}
	return thingWithCompositeEnumAttributess, nil
}