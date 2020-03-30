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

// ThingWithNullableAttrsInGSITable represents the user-configurable properties of the ThingWithNullableAttrsInGSI table.
type ThingWithNullableAttrsInGSITable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithNullableAttrsInGSIPrimaryKey represents the primary key of a ThingWithNullableAttrsInGSI in DynamoDB.
type ddbThingWithNullableAttrsInGSIPrimaryKey struct {
	PropertyOne string `dynamodbav:"propertyOne"`
}

// ddbThingWithNullableAttrsInGSIGSIByPropertyTwoAndPropertyThree represents the byPropertyTwoAndPropertyThree GSI.
type ddbThingWithNullableAttrsInGSIGSIByPropertyTwoAndPropertyThree struct {
	PropertyTwo   string `dynamodbav:"propertyTwo"`
	PropertyThree string `dynamodbav:"propertyThree"`
}

// ddbThingWithNullableAttrsInGSIGSIByPropertyTwoAndFourAndPropertyThree represents the byPropertyTwoAndFourAndPropertyThree GSI.
type ddbThingWithNullableAttrsInGSIGSIByPropertyTwoAndFourAndPropertyThree struct {
	PropertyTwoAndFour string `dynamodbav:"propertyTwoAndFour"`
	PropertyThree      string `dynamodbav:"propertyThree"`
}

// ddbThingWithNullableAttrsInGSI represents a ThingWithNullableAttrsInGSI as stored in DynamoDB.
type ddbThingWithNullableAttrsInGSI struct {
	models.ThingWithNullableAttrsInGSI
}

func (t ThingWithNullableAttrsInGSITable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-nullable-attrs-in-g-s-is", t.Prefix)
}

func (t ThingWithNullableAttrsInGSITable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("propertyOne"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("propertyThree"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("propertyTwo"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("propertyTwoAndFour"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("propertyOne"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("byPropertyTwoAndPropertyThree"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("propertyTwo"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("propertyThree"),
						KeyType:       aws.String(dynamodb.KeyTypeRange),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("byPropertyTwoAndFourAndPropertyThree"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("propertyTwoAndFour"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("propertyThree"),
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

func (t ThingWithNullableAttrsInGSITable) saveThingWithNullableAttrsInGSI(ctx context.Context, m models.ThingWithNullableAttrsInGSI) error {
	data, err := encodeThingWithNullableAttrsInGSI(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
	})
	return err
}

func (t ThingWithNullableAttrsInGSITable) getThingWithNullableAttrsInGSI(ctx context.Context, propertyOne string) (*models.ThingWithNullableAttrsInGSI, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithNullableAttrsInGSIPrimaryKey{
		PropertyOne: propertyOne,
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
		return nil, db.ErrThingWithNullableAttrsInGSINotFound{
			PropertyOne: propertyOne,
		}
	}

	var m models.ThingWithNullableAttrsInGSI
	if err := decodeThingWithNullableAttrsInGSI(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithNullableAttrsInGSITable) deleteThingWithNullableAttrsInGSI(ctx context.Context, propertyOne string) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithNullableAttrsInGSIPrimaryKey{
		PropertyOne: propertyOne,
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

func (t ThingWithNullableAttrsInGSITable) getThingWithNullableAttrsInGSIsByPropertyTwoAndPropertyThree(ctx context.Context, input db.GetThingWithNullableAttrsInGSIsByPropertyTwoAndPropertyThreeInput, fn func(m *models.ThingWithNullableAttrsInGSI, lastThingWithNullableAttrsInGSI bool) bool) error {
	if input.PropertyThreeStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.PropertyThreeStartingAt or input.StartingAfter")
	}
	if input.PropertyTwo == "" {
		return fmt.Errorf("Hash key input.PropertyTwo cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("byPropertyTwoAndPropertyThree"),
		ExpressionAttributeNames: map[string]*string{
			"#PROPERTYTWO": aws.String("propertyTwo"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":propertyTwo": &dynamodb.AttributeValue{
				S: aws.String(input.PropertyTwo),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.PropertyThreeStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#PROPERTYTWO = :propertyTwo")
	} else {
		queryInput.ExpressionAttributeNames["#PROPERTYTHREE"] = aws.String("propertyThree")
		queryInput.ExpressionAttributeValues[":propertyThree"] = &dynamodb.AttributeValue{
			S: aws.String(*input.PropertyThreeStartingAt),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#PROPERTYTWO = :propertyTwo AND #PROPERTYTHREE <= :propertyThree")
		} else {
			queryInput.KeyConditionExpression = aws.String("#PROPERTYTWO = :propertyTwo AND #PROPERTYTHREE >= :propertyThree")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"propertyThree": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.PropertyThree),
			},
			"propertyTwo": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.PropertyTwo),
			},
			"propertyOne": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.PropertyOne),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithNullableAttrsInGSIs(queryOutput.Items)
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

func (t ThingWithNullableAttrsInGSITable) getThingWithNullableAttrsInGSIsByPropertyTwoAndFourAndPropertyThree(ctx context.Context, input db.GetThingWithNullableAttrsInGSIsByPropertyTwoAndFourAndPropertyThreeInput, fn func(m *models.ThingWithNullableAttrsInGSI, lastThingWithNullableAttrsInGSI bool) bool) error {
	if input.PropertyThreeStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.PropertyThreeStartingAt or input.StartingAfter")
	}
	if input.PropertyTwo == "" {
		return fmt.Errorf("Hash key input.PropertyTwo cannot be empty")
	}
	if input.PropertyFour == "" {
		return fmt.Errorf("Hash key input.PropertyFour cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("byPropertyTwoAndFourAndPropertyThree"),
		ExpressionAttributeNames: map[string]*string{
			"#PROPERTYTWOANDFOUR": aws.String("propertyTwoAndFour"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":propertyTwoAndFour": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s_%s", input.PropertyTwo, input.PropertyFour)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.PropertyThreeStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#PROPERTYTWOANDFOUR = :propertyTwoAndFour")
	} else {
		queryInput.ExpressionAttributeNames["#PROPERTYTHREE"] = aws.String("propertyThree")
		queryInput.ExpressionAttributeValues[":propertyThree"] = &dynamodb.AttributeValue{
			S: aws.String(*input.PropertyThreeStartingAt),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#PROPERTYTWOANDFOUR = :propertyTwoAndFour AND #PROPERTYTHREE <= :propertyThree")
		} else {
			queryInput.KeyConditionExpression = aws.String("#PROPERTYTWOANDFOUR = :propertyTwoAndFour AND #PROPERTYTHREE >= :propertyThree")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"propertyThree": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.PropertyThree),
			},
			"propertyTwoAndFour": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s_%s", input.StartingAfter.PropertyTwo, input.StartingAfter.PropertyFour)),
			},
			"propertyOne": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.PropertyOne),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithNullableAttrsInGSIs(queryOutput.Items)
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

// encodeThingWithNullableAttrsInGSI encodes a ThingWithNullableAttrsInGSI as a DynamoDB map of attribute values.
func encodeThingWithNullableAttrsInGSI(m models.ThingWithNullableAttrsInGSI) (map[string]*dynamodb.AttributeValue, error) {
	val, err := dynamodbattribute.MarshalMap(ddbThingWithNullableAttrsInGSI{
		ThingWithNullableAttrsInGSI: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(m.PropertyFour, "_") {
		return nil, fmt.Errorf("propertyFour cannot contain '_': %s", m.PropertyFour)
	}
	if strings.Contains(m.PropertyTwo, "_") {
		return nil, fmt.Errorf("propertyTwo cannot contain '_': %s", m.PropertyTwo)
	}
	// add in composite attributes
	byPropertyTwoAndFourAndPropertyThree, err := dynamodbattribute.MarshalMap(ddbThingWithNullableAttrsInGSIGSIByPropertyTwoAndFourAndPropertyThree{
		PropertyTwoAndFour: fmt.Sprintf("%s_%s", m.PropertyTwo, m.PropertyFour),
		PropertyThree:      *m.PropertyThree,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range byPropertyTwoAndFourAndPropertyThree {
		val[k] = v
	}
	return val, err
}

// decodeThingWithNullableAttrsInGSI translates a ThingWithNullableAttrsInGSI stored in DynamoDB to a ThingWithNullableAttrsInGSI struct.
func decodeThingWithNullableAttrsInGSI(m map[string]*dynamodb.AttributeValue, out *models.ThingWithNullableAttrsInGSI) error {
	var ddbThingWithNullableAttrsInGSI ddbThingWithNullableAttrsInGSI
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithNullableAttrsInGSI); err != nil {
		return err
	}
	*out = ddbThingWithNullableAttrsInGSI.ThingWithNullableAttrsInGSI
	return nil
}

// decodeThingWithNullableAttrsInGSIs translates a list of ThingWithNullableAttrsInGSIs stored in DynamoDB to a slice of ThingWithNullableAttrsInGSI structs.
func decodeThingWithNullableAttrsInGSIs(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithNullableAttrsInGSI, error) {
	thingWithNullableAttrsInGSIs := make([]models.ThingWithNullableAttrsInGSI, len(ms))
	for i, m := range ms {
		var thingWithNullableAttrsInGSI models.ThingWithNullableAttrsInGSI
		if err := decodeThingWithNullableAttrsInGSI(m, &thingWithNullableAttrsInGSI); err != nil {
			return nil, err
		}
		thingWithNullableAttrsInGSIs[i] = thingWithNullableAttrsInGSI
	}
	return thingWithNullableAttrsInGSIs, nil
}
