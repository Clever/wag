package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Clever/wag/samples/gen-go-db-only/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db-only/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithDateGSITable represents the user-configurable properties of the ThingWithDateGSI table.
type ThingWithDateGSITable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithDateGSIPrimaryKey represents the primary key of a ThingWithDateGSI in DynamoDB.
type ddbThingWithDateGSIPrimaryKey struct {
	DateH strfmt.Date `dynamodbav:"dateH"`
}

// ddbThingWithDateGSIGSIRangeDate represents the rangeDate GSI.
type ddbThingWithDateGSIGSIRangeDate struct {
	ID    string      `dynamodbav:"id"`
	DateR strfmt.Date `dynamodbav:"dateR"`
}

// ddbThingWithDateGSIGSIHash represents the hash GSI.
type ddbThingWithDateGSIGSIHash struct {
	DateH strfmt.Date `dynamodbav:"dateH"`
	ID    string      `dynamodbav:"id"`
}

// ddbThingWithDateGSI represents a ThingWithDateGSI as stored in DynamoDB.
type ddbThingWithDateGSI struct {
	models.ThingWithDateGSI
}

func (t ThingWithDateGSITable) create(ctx context.Context) error {
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

func (t ThingWithDateGSITable) saveThingWithDateGSI(ctx context.Context, m models.ThingWithDateGSI) error {
	data, err := encodeThingWithDateGSI(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#DATEH": "dateH",
		},
		ConditionExpression: aws.String("attribute_not_exists(#DATEH)"),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		var conditionalCheckFailedErr *types.ConditionalCheckFailedException
		if errors.As(err, &resourceNotFoundErr) {
			return fmt.Errorf("table or index not found: %s", t.TableName)
		}
		if errors.As(err, &conditionalCheckFailedErr) {
			return db.ErrThingWithDateGSIAlreadyExists{
				DateH: m.DateH,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithDateGSITable) getThingWithDateGSI(ctx context.Context, dateH strfmt.Date) (*models.ThingWithDateGSI, error) {
	// swad-get-7
	key, err := attributevalue.MarshalMap(ddbThingWithDateGSIPrimaryKey{
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
		return nil, db.ErrThingWithDateGSINotFound{
			DateH: dateH,
		}
	}

	var m models.ThingWithDateGSI
	if err := decodeThingWithDateGSI(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithDateGSITable) scanThingWithDateGSIs(ctx context.Context, input db.ScanThingWithDateGSIsInput, fn func(m *models.ThingWithDateGSI, lastThingWithDateGSI bool) bool) error {
	// swad-scan-1
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
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
	totalRecordsProcessed := int32(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithDateGSIs(out.Items)
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

func (t ThingWithDateGSITable) deleteThingWithDateGSI(ctx context.Context, dateH strfmt.Date) error {

	key, err := attributevalue.MarshalMap(ddbThingWithDateGSIPrimaryKey{
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

func (t ThingWithDateGSITable) getThingWithDateGSIsByIDAndDateR(ctx context.Context, input db.GetThingWithDateGSIsByIDAndDateRInput, fn func(m *models.ThingWithDateGSI, lastThingWithDateGSI bool) bool) error {
	// swad-get-33
	if input.DateRStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateRStartingAt or input.StartingAfter")
	}
	// swad-get-33f
	if input.ID == "" {
		return fmt.Errorf("Hash key input.ID cannot be empty")
	}
	// swad-get-331
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("rangeDate"),
		ExpressionAttributeNames: map[string]string{
			"#ID": "id",
		},
		// swad-get-3312
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{
				// swad-get-33e
				Value: input.ID,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	// swad-get-332
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.DateRStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#ID = :id")
	} else {
		// swad-get-333
		queryInput.ExpressionAttributeNames["#DATER"] = "dateR"

		// swad-get-3331a
		queryInput.ExpressionAttributeValues[":dateR"] = &types.AttributeValueMemberS{
			Value: dateToDynamoTimeString(*input.DateRStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#ID = :id AND #DATER <= :dateR")
		} else {
			queryInput.KeyConditionExpression = aws.String("#ID = :id AND #DATER >= :dateR")
		}
	}
	// swad-get-334
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"dateR": &types.AttributeValueMemberS{
				Value: dateToDynamoTimeString(input.StartingAfter.DateR),
			},
			// swad-get-3341
			"id": &types.AttributeValueMemberS{
				Value: input.StartingAfter.ID,
			},
			// swad-get-3342

			// swad-get-336
			"dateH": &types.AttributeValueMemberS{
				Value: dateToDynamoTimeString(input.StartingAfter.DateH),
			},
		}
	}

	// swad-get-339

	totalRecordsProcessed := int32(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithDateGSIs(queryOutput.Items)
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
func (t ThingWithDateGSITable) getThingWithDateGSIsByDateHAndID(ctx context.Context, input db.GetThingWithDateGSIsByDateHAndIDInput, fn func(m *models.ThingWithDateGSI, lastThingWithDateGSI bool) bool) error {
	// swad-get-33
	if input.IDStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.IDStartingAt or input.StartingAfter")
	}
	if dateToDynamoTimeString(input.DateH) == "" {
		return fmt.Errorf("Hash key input.DateH cannot be empty")
	}
	// swad-get-331
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("hash"),
		ExpressionAttributeNames: map[string]string{
			"#DATEH": "dateH",
		},
		// swad-get-3312
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":dateH": &types.AttributeValueMemberS{
				// swad-get-33c
				Value: dateToDynamoTimeString(input.DateH),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	// swad-get-332
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.IDStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#DATEH = :dateH")
	} else {
		// swad-get-333
		queryInput.ExpressionAttributeNames["#ID"] = "id"

		// swad-get-3331a
		queryInput.ExpressionAttributeValues[":id"] = &types.AttributeValueMemberS{
			Value: string(*input.IDStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#DATEH = :dateH AND #ID <= :id")
		} else {
			queryInput.KeyConditionExpression = aws.String("#DATEH = :dateH AND #ID >= :id")
		}
	}
	// swad-get-334
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: string(input.StartingAfter.ID),
			},
			// swad-get-3341
			"dateH": &types.AttributeValueMemberS{
				Value: dateToDynamoTimeString(input.StartingAfter.DateH),
			},
			// swad-get-3342

			// swad-get-336
		}
	}

	// swad-get-339

	totalRecordsProcessed := int32(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithDateGSIs(queryOutput.Items)
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

// encodeThingWithDateGSI encodes a ThingWithDateGSI as a DynamoDB map of attribute values.
func encodeThingWithDateGSI(m models.ThingWithDateGSI) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(ddbThingWithDateGSI{
		ThingWithDateGSI: m,
	})
}

// decodeThingWithDateGSI translates a ThingWithDateGSI stored in DynamoDB to a ThingWithDateGSI struct.
func decodeThingWithDateGSI(m map[string]types.AttributeValue, out *models.ThingWithDateGSI) error {
	// swad-decode-1
	var ddbThingWithDateGSI ddbThingWithDateGSI
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithDateGSI); err != nil {
		return err
	}
	*out = ddbThingWithDateGSI.ThingWithDateGSI
	return nil
}

// decodeThingWithDateGSIs translates a list of ThingWithDateGSIs stored in DynamoDB to a slice of ThingWithDateGSI structs.
func decodeThingWithDateGSIs(ms []map[string]types.AttributeValue) ([]models.ThingWithDateGSI, error) {
	thingWithDateGSIs := make([]models.ThingWithDateGSI, len(ms))
	for i, m := range ms {
		var thingWithDateGSI models.ThingWithDateGSI
		if err := decodeThingWithDateGSI(m, &thingWithDateGSI); err != nil {
			return nil, err
		}
		thingWithDateGSIs[i] = thingWithDateGSI
	}
	return thingWithDateGSIs, nil
}
