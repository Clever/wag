package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-app-service/models"
	"github.com/Clever/wag/samples/gen-go-app-service/server/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// SetupStepTable represents the user-configurable properties of the SetupStep table.
type SetupStepTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbSetupStepPrimaryKey represents the primary key of a SetupStep in DynamoDB.
type ddbSetupStepPrimaryKey struct {
	AppID string `dynamodbav:"app_id"`
	ID    string `dynamodbav:"id"`
}

// ddbSetupStep represents a SetupStep as stored in DynamoDB.
type ddbSetupStep struct {
	models.SetupStep
}

func (t SetupStepTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-setup-steps", t.Prefix)
}

func (t SetupStepTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("app_id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("app_id"),
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

func (t SetupStepTable) saveSetupStep(ctx context.Context, m models.SetupStep) error {
	data, err := encodeSetupStep(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
	})
	return err
}

func (t SetupStepTable) getSetupStep(ctx context.Context, appID string, id string) (*models.SetupStep, error) {
	key, err := dynamodbattribute.MarshalMap(ddbSetupStepPrimaryKey{
		AppID: appID,
		ID:    id,
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
		return nil, db.ErrSetupStepNotFound{
			AppID: appID,
			ID:    id,
		}
	}

	var m models.SetupStep
	if err := decodeSetupStep(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t SetupStepTable) getSetupStepsByAppIDAndID(ctx context.Context, input db.GetSetupStepsByAppIDAndIDInput, fn func(m *models.SetupStep, lastSetupStep bool) bool) error {
	if input.IDStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.IDStartingAt or input.StartingAfter")
	}
	if input.AppID == "" {
		return fmt.Errorf("Hash key input.AppID cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#APP_ID": aws.String("app_id"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":appId": &dynamodb.AttributeValue{
				S: aws.String(input.AppID),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.IDStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#APP_ID = :appId")
	} else {
		queryInput.ExpressionAttributeNames["#ID"] = aws.String("id")
		queryInput.ExpressionAttributeValues[":id"] = &dynamodb.AttributeValue{
			S: aws.String(*input.IDStartingAt),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#APP_ID = :appId AND #ID <= :id")
		} else {
			queryInput.KeyConditionExpression = aws.String("#APP_ID = :appId AND #ID >= :id")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.ID),
			},
			"app_id": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.AppID),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeSetupSteps(queryOutput.Items)
		if err != nil {
			pageFnErr = err
			return false
		}
		hasMore := true
		for i := range items {
			if lastPage == true {
				hasMore = i < len(items)-1
			}
			if !fn(&items[i], !hasMore) {
				return false
			}
			totalRecordsProcessed++
			// if the Limit of records have been passed to fn, don't pass anymore records.
			if input.Limit != nil && totalRecordsProcessed == *input.Limit {
				return false
			}
		}
		return true
	}

	err := t.DynamoDBAPI.QueryPagesWithContext(ctx, queryInput, pageFn)
	if err != nil {
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}

func (t SetupStepTable) deleteSetupStep(ctx context.Context, appID string, id string) error {
	key, err := dynamodbattribute.MarshalMap(ddbSetupStepPrimaryKey{
		AppID: appID,
		ID:    id,
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

// encodeSetupStep encodes a SetupStep as a DynamoDB map of attribute values.
func encodeSetupStep(m models.SetupStep) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbSetupStep{
		SetupStep: m,
	})
}

// decodeSetupStep translates a SetupStep stored in DynamoDB to a SetupStep struct.
func decodeSetupStep(m map[string]*dynamodb.AttributeValue, out *models.SetupStep) error {
	var ddbSetupStep ddbSetupStep
	if err := dynamodbattribute.UnmarshalMap(m, &ddbSetupStep); err != nil {
		return err
	}
	*out = ddbSetupStep.SetupStep
	return nil
}

// decodeSetupSteps translates a list of SetupSteps stored in DynamoDB to a slice of SetupStep structs.
func decodeSetupSteps(ms []map[string]*dynamodb.AttributeValue) ([]models.SetupStep, error) {
	setupSteps := make([]models.SetupStep, len(ms))
	for i, m := range ms {
		var setupStep models.SetupStep
		if err := decodeSetupStep(m, &setupStep); err != nil {
			return nil, err
		}
		setupSteps[i] = setupStep
	}
	return setupSteps, nil
}
