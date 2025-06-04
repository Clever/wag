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

// ThingWithEnumHashKeyTable represents the user-configurable properties of the ThingWithEnumHashKey table.
type ThingWithEnumHashKeyTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithEnumHashKeyPrimaryKey represents the primary key of a ThingWithEnumHashKey in DynamoDB.
type ddbThingWithEnumHashKeyPrimaryKey struct {
	Branch models.Branch   `dynamodbav:"branch"`
	Date   strfmt.DateTime `dynamodbav:"date"`
}

// ddbThingWithEnumHashKeyGSIByBranch represents the byBranch GSI.
type ddbThingWithEnumHashKeyGSIByBranch struct {
	Branch models.Branch   `dynamodbav:"branch"`
	Date2  strfmt.DateTime `dynamodbav:"date2"`
}

// ddbThingWithEnumHashKey represents a ThingWithEnumHashKey as stored in DynamoDB.
type ddbThingWithEnumHashKey struct {
	models.ThingWithEnumHashKey
}

func (t ThingWithEnumHashKeyTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("branch"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("date"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("date2"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("branch"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("date"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("byBranch"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("branch"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("date2"),
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

func (t ThingWithEnumHashKeyTable) saveThingWithEnumHashKey(ctx context.Context, m models.ThingWithEnumHashKey) error {
	data, err := encodeThingWithEnumHashKey(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#BRANCH": "branch",
			"#DATE":   "date",
		},
		ConditionExpression: aws.String("attribute_not_exists(#BRANCH) AND attribute_not_exists(#DATE)"),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		var conditionalCheckFailedErr *types.ConditionalCheckFailedException
		if errors.As(err, &resourceNotFoundErr) {
			return fmt.Errorf("table or index not found: %s", t.TableName)
		}
		if errors.As(err, &conditionalCheckFailedErr) {
			return db.ErrThingWithEnumHashKeyAlreadyExists{
				Branch: m.Branch,
				Date:   m.Date,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithEnumHashKeyTable) getThingWithEnumHashKey(ctx context.Context, branch models.Branch, date strfmt.DateTime) (*models.ThingWithEnumHashKey, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithEnumHashKeyPrimaryKey{
		Branch: branch,
		Date:   date,
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
		return nil, db.ErrThingWithEnumHashKeyNotFound{
			Branch: branch,
			Date:   date,
		}
	}

	var m models.ThingWithEnumHashKey
	if err := decodeThingWithEnumHashKey(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithEnumHashKeyTable) scanThingWithEnumHashKeys(ctx context.Context, input db.ScanThingWithEnumHashKeysInput, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error {
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
			"branch": exclusiveStartKey["branch"],
			"date":   exclusiveStartKey["date"],
		}
	}
	totalRecordsProcessed := int32(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithEnumHashKeys(out.Items)
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

func (t ThingWithEnumHashKeyTable) getThingWithEnumHashKeysByBranchAndDateParseFilters(queryInput *dynamodb.QueryInput, input db.GetThingWithEnumHashKeysByBranchAndDateInput) {
	for _, filterValue := range input.FilterValues {
		switch filterValue.AttributeName {
		case db.ThingWithEnumHashKeyDate2:
			queryInput.ExpressionAttributeNames["#DATE2"] = string(db.ThingWithEnumHashKeyDate2)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithEnumHashKeyDate2), i)] = &types.AttributeValueMemberS{
					Value: datetimeToDynamoTimeString(attributeValue.(strfmt.DateTime)),
				}
			}
		}
	}
}

func (t ThingWithEnumHashKeyTable) getThingWithEnumHashKeysByBranchAndDate(ctx context.Context, input db.GetThingWithEnumHashKeysByBranchAndDateInput, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error {
	if input.DateStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateStartingAt or input.StartingAfter")
	}
	if input.Branch == "" {
		return fmt.Errorf("Hash key input.Branch cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		ExpressionAttributeNames: map[string]string{
			"#BRANCH": "branch",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":branch": &types.AttributeValueMemberS{
				Value: string(input.Branch),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = "date"
		queryInput.ExpressionAttributeValues[":date"] = &types.AttributeValueMemberS{
			Value: datetimeToDynamoTimeString(*input.DateStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch AND #DATE <= :date")
		} else {
			queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch AND #DATE >= :date")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"date": &types.AttributeValueMemberS{
				Value: datetimeToDynamoTimeString(input.StartingAfter.Date),
			},

			"branch": &types.AttributeValueMemberS{
				Value: string(input.StartingAfter.Branch),
			},
		}
	}
	if len(input.FilterValues) > 0 && input.FilterExpression != "" {
		t.getThingWithEnumHashKeysByBranchAndDateParseFilters(queryInput, input)
		queryInput.FilterExpression = aws.String(input.FilterExpression)
	}

	totalRecordsProcessed := int32(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithEnumHashKeys(queryOutput.Items)
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

func (t ThingWithEnumHashKeyTable) deleteThingWithEnumHashKey(ctx context.Context, branch models.Branch, date strfmt.DateTime) error {

	key, err := attributevalue.MarshalMap(ddbThingWithEnumHashKeyPrimaryKey{
		Branch: branch,
		Date:   date,
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

func (t ThingWithEnumHashKeyTable) getThingWithEnumHashKeysByBranchAndDate2(ctx context.Context, input db.GetThingWithEnumHashKeysByBranchAndDate2Input, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error {
	if input.Date2StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.Date2StartingAt or input.StartingAfter")
	}
	if input.Branch == "" {
		return fmt.Errorf("Hash key input.Branch cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("byBranch"),
		ExpressionAttributeNames: map[string]string{
			"#BRANCH": "branch",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":branch": &types.AttributeValueMemberS{
				Value: string(input.Branch),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.Date2StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch")
	} else {
		queryInput.ExpressionAttributeNames["#DATE2"] = "date2"
		queryInput.ExpressionAttributeValues[":date2"] = &types.AttributeValueMemberS{
			Value: datetimeToDynamoTimeString(*input.Date2StartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch AND #DATE2 <= :date2")
		} else {
			queryInput.KeyConditionExpression = aws.String("#BRANCH = :branch AND #DATE2 >= :date2")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"date2": &types.AttributeValueMemberS{
				Value: datetimeToDynamoTimeString(input.StartingAfter.Date2),
			},
			"branch": &types.AttributeValueMemberS{
				Value: string(input.StartingAfter.Branch),
			},
			"date": &types.AttributeValueMemberS{
				Value: datetimeToDynamoTimeString(input.StartingAfter.Date),
			},
		}
	}

	totalRecordsProcessed := int32(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithEnumHashKeys(queryOutput.Items)
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
func (t ThingWithEnumHashKeyTable) scanThingWithEnumHashKeysByBranchAndDate2(ctx context.Context, input db.ScanThingWithEnumHashKeysByBranchAndDate2Input, fn func(m *models.ThingWithEnumHashKey, lastThingWithEnumHashKey bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("byBranch"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"branch": exclusiveStartKey["branch"],
			"date":   exclusiveStartKey["date"],
			"date2":  exclusiveStartKey["date2"],
		}
	}
	totalRecordsProcessed := int32(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithEnumHashKeys(out.Items)
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

// encodeThingWithEnumHashKey encodes a ThingWithEnumHashKey as a DynamoDB map of attribute values.
func encodeThingWithEnumHashKey(m models.ThingWithEnumHashKey) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(ddbThingWithEnumHashKey{
		ThingWithEnumHashKey: m,
	})
}

// decodeThingWithEnumHashKey translates a ThingWithEnumHashKey stored in DynamoDB to a ThingWithEnumHashKey struct.
func decodeThingWithEnumHashKey(m map[string]types.AttributeValue, out *models.ThingWithEnumHashKey) error {
	var ddbThingWithEnumHashKey ddbThingWithEnumHashKey
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithEnumHashKey); err != nil {
		return err
	}
	*out = ddbThingWithEnumHashKey.ThingWithEnumHashKey
	return nil
}

// decodeThingWithEnumHashKeys translates a list of ThingWithEnumHashKeys stored in DynamoDB to a slice of ThingWithEnumHashKey structs.
func decodeThingWithEnumHashKeys(ms []map[string]types.AttributeValue) ([]models.ThingWithEnumHashKey, error) {
	thingWithEnumHashKeys := make([]models.ThingWithEnumHashKey, len(ms))
	for i, m := range ms {
		var thingWithEnumHashKey models.ThingWithEnumHashKey
		if err := decodeThingWithEnumHashKey(m, &thingWithEnumHashKey); err != nil {
			return nil, err
		}
		thingWithEnumHashKeys[i] = thingWithEnumHashKey
	}
	return thingWithEnumHashKeys, nil
}
