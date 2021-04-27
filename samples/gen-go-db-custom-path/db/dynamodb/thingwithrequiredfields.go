package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/v7/samples/gen-go-db-custom-path/db"
	"github.com/Clever/wag/v7/samples/gen-go-db-custom-path/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

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
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return db.ErrThingWithRequiredFieldsAlreadyExists{
					Name: *m.Name,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
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

func (t ThingWithRequiredFieldsTable) scanThingWithRequiredFieldss(ctx context.Context, input db.ScanThingWithRequiredFieldssInput, fn func(m *models.ThingWithRequiredFields, lastThingWithRequiredFields bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"name": exclusiveStartKey["name"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithRequiredFieldss(out.Items)
		if err != nil {
			innerErr = fmt.Errorf("error decoding %s", err.Error())
			return false
		}
		for i := range items {
			if input.Limiter != nil {
				if err := input.Limiter.Wait(ctx); err != nil {
					innerErr = err
					return false
				}
			}
			isLastModel := lastPage && i == len(items)-1
			if shouldContinue := fn(&items[i], isLastModel); !shouldContinue {
				return false
			}
			totalRecordsProcessed++
			// if the Limit of records have been passed to fn, don't pass anymore records.
			if input.Limit != nil && totalRecordsProcessed == *input.Limit {
				return false
			}
		}
		return true
	})
	if innerErr != nil {
		return innerErr
	}
	return err
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
