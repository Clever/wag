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

// ThingWithCompositeAttributesTable represents the user-configurable properties of the ThingWithCompositeAttributes table.
type ThingWithCompositeAttributesTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithCompositeAttributesPrimaryKey represents the primary key of a ThingWithCompositeAttributes in DynamoDB.
type ddbThingWithCompositeAttributesPrimaryKey struct {
	NameBranch string          `dynamodbav:"name_branch"`
	Date       strfmt.DateTime `dynamodbav:"date"`
}

// ddbThingWithCompositeAttributesGSINameVersion represents the nameVersion GSI.
type ddbThingWithCompositeAttributesGSINameVersion struct {
	NameVersion string          `dynamodbav:"name_version"`
	Date        strfmt.DateTime `dynamodbav:"date"`
}

// ddbThingWithCompositeAttributes represents a ThingWithCompositeAttributes as stored in DynamoDB.
type ddbThingWithCompositeAttributes struct {
	models.ThingWithCompositeAttributes
}

func (t ThingWithCompositeAttributesTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-composite-attributess", t.Prefix)
}

func (t ThingWithCompositeAttributesTable) create(ctx context.Context) error {
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
			{
				AttributeName: aws.String("name_version"),
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
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("nameVersion"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("name_version"),
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

func (t ThingWithCompositeAttributesTable) saveThingWithCompositeAttributes(ctx context.Context, m models.ThingWithCompositeAttributes) error {
	data, err := encodeThingWithCompositeAttributes(m)
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
				return db.ErrThingWithCompositeAttributesAlreadyExists{
					NameBranch: fmt.Sprintf("%s@%s", *m.Name, *m.Branch),
					Date:       *m.Date,
				}
			}
		}
		return err
	}
	return nil
}

func (t ThingWithCompositeAttributesTable) getThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) (*models.ThingWithCompositeAttributes, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branch),
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
		return nil, db.ErrThingWithCompositeAttributesNotFound{
			Name:   name,
			Branch: branch,
			Date:   date,
		}
	}

	var m models.ThingWithCompositeAttributes
	if err := decodeThingWithCompositeAttributes(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithCompositeAttributesTable) getThingWithCompositeAttributessByNameBranchAndDate(ctx context.Context, input db.GetThingWithCompositeAttributessByNameBranchAndDateInput) ([]models.ThingWithCompositeAttributes, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#NAME_BRANCH": aws.String("name_branch"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":nameBranch": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s@%s", input.Name, input.Branch)),
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
		return []models.ThingWithCompositeAttributes{}, nil
	}

	return decodeThingWithCompositeAttributess(queryOutput.Items)
}

func (t ThingWithCompositeAttributesTable) deleteThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branch),
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

func (t ThingWithCompositeAttributesTable) getThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input db.GetThingWithCompositeAttributessByNameVersionAndDateInput) ([]models.ThingWithCompositeAttributes, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("nameVersion"),
		ExpressionAttributeNames: map[string]*string{
			"#NAME_VERSION": aws.String("name_version"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":nameVersion": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s:%d", input.Name, input.Version)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME_VERSION = :nameVersion")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = aws.String("date")
		queryInput.ExpressionAttributeValues[":date"] = &dynamodb.AttributeValue{
			S: aws.String(time.Time(*input.DateStartingAt).Format(time.RFC3339)), // dynamodb attributevalue only supports RFC3339 resolution
		}
		queryInput.KeyConditionExpression = aws.String("#NAME_VERSION = :nameVersion AND #DATE >= :date")
	}

	queryOutput, err := t.DynamoDBAPI.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if len(queryOutput.Items) == 0 {
		return []models.ThingWithCompositeAttributes{}, nil
	}

	return decodeThingWithCompositeAttributess(queryOutput.Items)
}

// encodeThingWithCompositeAttributes encodes a ThingWithCompositeAttributes as a DynamoDB map of attribute values.
func encodeThingWithCompositeAttributes(m models.ThingWithCompositeAttributes) (map[string]*dynamodb.AttributeValue, error) {
	val, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeAttributes{
		ThingWithCompositeAttributes: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(*m.Name, ":") {
		return nil, fmt.Errorf("name cannot contain ':': %s", *m.Name)
	}
	if strings.Contains(*m.Branch, "@") {
		return nil, fmt.Errorf("branch cannot contain '@': %s", *m.Branch)
	}
	if strings.Contains(*m.Name, "@") {
		return nil, fmt.Errorf("name cannot contain '@': %s", *m.Name)
	}
	// add in composite attributes
	primaryKey, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", *m.Name, *m.Branch),
		Date:       *m.Date,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range primaryKey {
		val[k] = v
	}
	nameVersion, err := dynamodbattribute.MarshalMap(ddbThingWithCompositeAttributesGSINameVersion{
		NameVersion: fmt.Sprintf("%s:%d", *m.Name, m.Version),
		Date:        *m.Date,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range nameVersion {
		val[k] = v
	}
	return val, err
}

// decodeThingWithCompositeAttributes translates a ThingWithCompositeAttributes stored in DynamoDB to a ThingWithCompositeAttributes struct.
func decodeThingWithCompositeAttributes(m map[string]*dynamodb.AttributeValue, out *models.ThingWithCompositeAttributes) error {
	var ddbThingWithCompositeAttributes ddbThingWithCompositeAttributes
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithCompositeAttributes); err != nil {
		return err
	}
	*out = ddbThingWithCompositeAttributes.ThingWithCompositeAttributes
	return nil
}

// decodeThingWithCompositeAttributess translates a list of ThingWithCompositeAttributess stored in DynamoDB to a slice of ThingWithCompositeAttributes structs.
func decodeThingWithCompositeAttributess(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithCompositeAttributes, error) {
	thingWithCompositeAttributess := make([]models.ThingWithCompositeAttributes, len(ms))
	for i, m := range ms {
		var thingWithCompositeAttributes models.ThingWithCompositeAttributes
		if err := decodeThingWithCompositeAttributes(m, &thingWithCompositeAttributes); err != nil {
			return nil, err
		}
		thingWithCompositeAttributess[i] = thingWithCompositeAttributes
	}
	return thingWithCompositeAttributess, nil
}
