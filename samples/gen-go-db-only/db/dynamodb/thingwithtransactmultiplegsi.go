package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db-only/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db-only/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}
var _ = errors.New("")
var _ = []types.AttributeValue{}

// ThingWithTransactMultipleGSITable represents the user-configurable properties of the ThingWithTransactMultipleGSI table.
type ThingWithTransactMultipleGSITable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithTransactMultipleGSIPrimaryKey represents the primary key of a ThingWithTransactMultipleGSI in DynamoDB.
type ddbThingWithTransactMultipleGSIPrimaryKey struct {
	DateH strfmt.Date `dynamodbav:"dateH"`
}

// ddbThingWithTransactMultipleGSIGSIRangeDate represents the rangeDate GSI.
type ddbThingWithTransactMultipleGSIGSIRangeDate struct {
	ID    string      `dynamodbav:"id"`
	DateR strfmt.Date `dynamodbav:"dateR"`
}

// ddbThingWithTransactMultipleGSIGSIHash represents the hash GSI.
type ddbThingWithTransactMultipleGSIGSIHash struct {
	DateH strfmt.Date `dynamodbav:"dateH"`
	ID    string      `dynamodbav:"id"`
}

// ddbThingWithTransactMultipleGSI represents a ThingWithTransactMultipleGSI as stored in DynamoDB.
type ddbThingWithTransactMultipleGSI struct {
	models.ThingWithTransactMultipleGSI
}

func (t ThingWithTransactMultipleGSITable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("dateH"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("dateR"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("dateH"),
				KeyType:       types.KeyTypeHash,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("rangeDate"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("id"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("dateR"),
						KeyType:       types.KeyTypeRange,
					},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("hash"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("dateH"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("id"),
						KeyType:       types.KeyTypeRange,
					},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
		},
		TableName: aws.String(t.TableName),
	}); err != nil {
		return fmt.Errorf("failed to create table %s: %w", t.TableName, err)
	}
	return nil
}

func (t ThingWithTransactMultipleGSITable) saveThingWithTransactMultipleGSI(ctx context.Context, m models.ThingWithTransactMultipleGSI) error {
	data, err := encodeThingWithTransactMultipleGSI(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#DATEH": "dateH",
		},
		ConditionExpression: aws.String(
			"" +
				"" +
				"attribute_not_exists(#DATEH)" +
				"",
		),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		var conditionalCheckFailedErr *types.ConditionalCheckFailedException
		if errors.As(err, &resourceNotFoundErr) {
			return fmt.Errorf("table or index not found: %s", t.TableName)
		}
		if errors.As(err, &conditionalCheckFailedErr) {
			return db.ErrThingWithTransactMultipleGSIAlreadyExists{
				DateH: m.DateH,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithTransactMultipleGSITable) getThingWithTransactMultipleGSI(ctx context.Context, dateH strfmt.Date) (*models.ThingWithTransactMultipleGSI, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithTransactMultipleGSIPrimaryKey{
		DateH: dateH,
	})
	if err != nil {
		return nil, err
	}
	res, err := t.DynamoDBAPI.GetItem(ctx, &dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundErr) {
			return nil, fmt.Errorf("table or index not found: %s", t.TableName)
		}
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrThingWithTransactMultipleGSINotFound{
			DateH: dateH,
		}
	}

	var m models.ThingWithTransactMultipleGSI
	if err := decodeThingWithTransactMultipleGSI(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithTransactMultipleGSITable) scanThingWithTransactMultipleGSIs(ctx context.Context, input db.ScanThingWithTransactMultipleGSIsInput, fn func(m *models.ThingWithTransactMultipleGSI, lastThingWithTransactMultipleGSI bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          aws.Int32(int32(*input.Limit)),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"dateH": exclusiveStartKey["dateH"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithTransactMultipleGSIs(out.Items)
		if err != nil {
			return fmt.Errorf("error decoding items: %s", err.Error())
		}

		for i := range items {
			if input.Limiter != nil {
				if err := input.Limiter.Wait(ctx); err != nil {
					return err
				}
			}

			isLastModel := !paginator.HasMorePages() && i == len(items)-1
			if shouldContinue := fn(&items[i], isLastModel); !shouldContinue {
				return nil
			}

			totalRecordsProcessed++
			if input.Limit != nil && totalRecordsProcessed == *input.Limit {
				return nil
			}
		}
	}

	return nil
}

func (t ThingWithTransactMultipleGSITable) deleteThingWithTransactMultipleGSI(ctx context.Context, dateH strfmt.Date) error {

	key, err := attributevalue.MarshalMap(ddbThingWithTransactMultipleGSIPrimaryKey{
		DateH: dateH,
	})
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(t.TableName),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundErr) {
			return fmt.Errorf("table or index not found: %s", t.TableName)
		}
		return err
	}

	return nil
}

func (t ThingWithTransactMultipleGSITable) getThingWithTransactMultipleGSIsByIDAndDateR(ctx context.Context, input db.GetThingWithTransactMultipleGSIsByIDAndDateRInput, fn func(m *models.ThingWithTransactMultipleGSI, lastThingWithTransactMultipleGSI bool) bool) error {
	if input.DateRStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateRStartingAt or input.StartingAfter")
	}
	if input.ID == "" {
		return fmt.Errorf("Hash key input.ID cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("rangeDate"),
		ExpressionAttributeNames: map[string]string{
			"#ID": "id",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{
				Value: input.ID,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.DateRStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#ID = :id")
	} else {
		queryInput.ExpressionAttributeNames["#DATER"] = "dateR"
		queryInput.ExpressionAttributeValues[":dateR"] = &types.AttributeValueMemberS{
			Value: dateToDynamoTimeString(*input.DateRStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#ID = :id AND #DATER <= :dateR")
		} else {
			queryInput.KeyConditionExpression = aws.String("#ID = :id AND #DATER >= :dateR")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"dateR": &types.AttributeValueMemberS{
				Value: dateToDynamoTimeString(input.StartingAfter.DateR),
			},
			"id": &types.AttributeValueMemberS{
				Value: input.StartingAfter.ID,
			},
			"dateH": &types.AttributeValueMemberS{
				Value: dateToDynamoTimeString(input.StartingAfter.DateH),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithTransactMultipleGSIs(queryOutput.Items)
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

	paginator := dynamodb.NewQueryPaginator(t.DynamoDBAPI, queryInput)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			var resourceNotFoundErr *types.ResourceNotFoundException
			if errors.As(err, &resourceNotFoundErr) {
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
			return err
		}
		if !pageFn(output, !paginator.HasMorePages()) {
			break
		}
	}

	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}
func (t ThingWithTransactMultipleGSITable) getThingWithTransactMultipleGSIsByDateHAndID(ctx context.Context, input db.GetThingWithTransactMultipleGSIsByDateHAndIDInput, fn func(m *models.ThingWithTransactMultipleGSI, lastThingWithTransactMultipleGSI bool) bool) error {
	if input.IDStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.IDStartingAt or input.StartingAfter")
	}
	if dateToDynamoTimeString(input.DateH) == "" {
		return fmt.Errorf("Hash key input.DateH cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("hash"),
		ExpressionAttributeNames: map[string]string{
			"#DATEH": "dateH",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":dateH": &types.AttributeValueMemberS{
				Value: dateToDynamoTimeString(input.DateH),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.IDStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#DATEH = :dateH")
	} else {
		queryInput.ExpressionAttributeNames["#ID"] = "id"
		queryInput.ExpressionAttributeValues[":id"] = &types.AttributeValueMemberS{
			Value: string(*input.IDStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#DATEH = :dateH AND #ID <= :id")
		} else {
			queryInput.KeyConditionExpression = aws.String("#DATEH = :dateH AND #ID >= :id")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: string(input.StartingAfter.ID),
			},
			"dateH": &types.AttributeValueMemberS{
				Value: dateToDynamoTimeString(input.StartingAfter.DateH),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithTransactMultipleGSIs(queryOutput.Items)
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

	paginator := dynamodb.NewQueryPaginator(t.DynamoDBAPI, queryInput)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			var resourceNotFoundErr *types.ResourceNotFoundException
			if errors.As(err, &resourceNotFoundErr) {
				return fmt.Errorf("table or index not found: %s", t.TableName)
			}
			return err
		}
		if !pageFn(output, !paginator.HasMorePages()) {
			break
		}
	}

	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}
func (t ThingWithTransactMultipleGSITable) transactSaveThingWithTransactMultipleGSIAndThing(ctx context.Context, m1 models.ThingWithTransactMultipleGSI, m1Conditions *expression.ConditionBuilder, m2 models.Thing, m2Conditions *expression.ConditionBuilder) error {
	data1, err := encodeThingWithTransactMultipleGSI(m1)
	if err != nil {
		return err
	}

	m1CondExpr, m1ExprVals, m1ExprNames, err := buildCondExpr(m1Conditions)
	if err != nil {
		return err
	}

	data2, err := encodeThing(m2)
	if err != nil {
		return err
	}

	m2CondExpr, m2ExprVals, m2ExprNames, err := buildCondExpr(m2Conditions)
	if err != nil {
		return err
	}

	// Convert map[string]*string to map[string]string for ExpressionAttributeNames
	toStringMap := func(in map[string]*string) map[string]string {
		if in == nil {
			return nil
		}
		out := make(map[string]string, len(in))
		for k, v := range in {
			if v != nil {
				out[k] = *v
			}
		}
		return out
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName:                 aws.String(t.TableName),
					Item:                      data1,
					ConditionExpression:       m1CondExpr,
					ExpressionAttributeValues: m1ExprVals,
					ExpressionAttributeNames:  toStringMap(m1ExprNames),
				},
			},
			{
				Put: &types.Put{
					TableName:                 aws.String(fmt.Sprintf("%s-Things", t.Prefix)),
					Item:                      data2,
					ConditionExpression:       m2CondExpr,
					ExpressionAttributeValues: m2ExprVals,
					ExpressionAttributeNames:  toStringMap(m2ExprNames),
				},
			},
		},
	}
	_, err = t.DynamoDBAPI.TransactWriteItems(ctx, input)

	return err
}

// encodeThingWithTransactMultipleGSI encodes a ThingWithTransactMultipleGSI as a DynamoDB map of attribute values.
func encodeThingWithTransactMultipleGSI(m models.ThingWithTransactMultipleGSI) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(ddbThingWithTransactMultipleGSI{
		ThingWithTransactMultipleGSI: m,
	})
}

// decodeThingWithTransactMultipleGSI translates a ThingWithTransactMultipleGSI stored in DynamoDB to a ThingWithTransactMultipleGSI struct.
func decodeThingWithTransactMultipleGSI(m map[string]types.AttributeValue, out *models.ThingWithTransactMultipleGSI) error {
	var ddbThingWithTransactMultipleGSI ddbThingWithTransactMultipleGSI
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithTransactMultipleGSI); err != nil {
		return err
	}
	*out = ddbThingWithTransactMultipleGSI.ThingWithTransactMultipleGSI
	return nil
}

// decodeThingWithTransactMultipleGSIs translates a list of ThingWithTransactMultipleGSIs stored in DynamoDB to a slice of ThingWithTransactMultipleGSI structs.
func decodeThingWithTransactMultipleGSIs(ms []map[string]types.AttributeValue) ([]models.ThingWithTransactMultipleGSI, error) {
	thingWithTransactMultipleGSIs := make([]models.ThingWithTransactMultipleGSI, len(ms))
	for i, m := range ms {
		var thingWithTransactMultipleGSI models.ThingWithTransactMultipleGSI
		if err := decodeThingWithTransactMultipleGSI(m, &thingWithTransactMultipleGSI); err != nil {
			return nil, err
		}
		thingWithTransactMultipleGSIs[i] = thingWithTransactMultipleGSI
	}
	return thingWithTransactMultipleGSIs, nil
}
