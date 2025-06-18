package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Clever/wag/samples/gen-go-db/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db/server/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}
var _ = errors.New("")
var _ = []types.AttributeValue{}

// ThingWithMatchingKeysTable represents the user-configurable properties of the ThingWithMatchingKeys table.
type ThingWithMatchingKeysTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithMatchingKeysPrimaryKey represents the primary key of a ThingWithMatchingKeys in DynamoDB.
type ddbThingWithMatchingKeysPrimaryKey struct {
	Bear        string `dynamodbav:"bear"`
	AssocTypeID string `dynamodbav:"assocTypeID"`
}

// ddbThingWithMatchingKeysGSIByAssoc represents the byAssoc GSI.
type ddbThingWithMatchingKeysGSIByAssoc struct {
	AssocTypeID string `dynamodbav:"assocTypeID"`
	CreatedBear string `dynamodbav:"createdBear"`
}

// ddbThingWithMatchingKeys represents a ThingWithMatchingKeys as stored in DynamoDB.
type ddbThingWithMatchingKeys struct {
	models.ThingWithMatchingKeys
}

func (t ThingWithMatchingKeysTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("assocTypeID"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("bear"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("createdBear"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("bear"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("assocTypeID"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("byAssoc"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("assocTypeID"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("createdBear"),
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

func (t ThingWithMatchingKeysTable) saveThingWithMatchingKeys(ctx context.Context, m models.ThingWithMatchingKeys) error {
	data, err := encodeThingWithMatchingKeys(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
	})
	return err
}

func (t ThingWithMatchingKeysTable) getThingWithMatchingKeys(ctx context.Context, bear string, assocType string, assocID string) (*models.ThingWithMatchingKeys, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithMatchingKeysPrimaryKey{
		Bear:        bear,
		AssocTypeID: fmt.Sprintf("%s^%s", assocType, assocID),
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
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrThingWithMatchingKeysNotFound{
			Bear:      bear,
			AssocType: assocType,
			AssocID:   assocID,
		}
	}

	var m models.ThingWithMatchingKeys
	if err := decodeThingWithMatchingKeys(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithMatchingKeysTable) scanThingWithMatchingKeyss(ctx context.Context, input db.ScanThingWithMatchingKeyssInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"bear": exclusiveStartKey["bear"],
			"assocTypeID": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s^%s", input.StartingAfter.AssocType, input.StartingAfter.AssocID),
			},
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithMatchingKeyss(out.Items)
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

func (t ThingWithMatchingKeysTable) getThingWithMatchingKeyssByBearAndAssocTypeIDParseFilters(queryInput *dynamodb.QueryInput, input db.GetThingWithMatchingKeyssByBearAndAssocTypeIDInput) {
	for _, filterValue := range input.FilterValues {
		switch filterValue.AttributeName {
		case db.ThingWithMatchingKeysCreated:
			queryInput.ExpressionAttributeNames["#CREATED"] = string(db.ThingWithMatchingKeysCreated)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithMatchingKeysCreated), i)] = &types.AttributeValueMemberS{
					Value: datetimeToDynamoTimeString(attributeValue.(strfmt.DateTime)),
				}
			}
		}
	}
}

func (t ThingWithMatchingKeysTable) getThingWithMatchingKeyssByBearAndAssocTypeID(ctx context.Context, input db.GetThingWithMatchingKeyssByBearAndAssocTypeIDInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
	if input.StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of StartingAt or StartingAfter")
	}
	if input.Bear == "" {
		return fmt.Errorf("Hash key input.Bear cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		ExpressionAttributeNames: map[string]string{
			"#BEAR": "bear",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":bear": &types.AttributeValueMemberS{
				Value: input.Bear,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#BEAR = :bear")
	} else {
		queryInput.ExpressionAttributeNames["#ASSOCTYPEID"] = "assocTypeID"
		queryInput.ExpressionAttributeValues[":assocTypeId"] = &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%s^%s", input.StartingAt.AssocType, input.StartingAt.AssocID),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#BEAR = :bear AND #ASSOCTYPEID <= :assocTypeId")
		} else {
			queryInput.KeyConditionExpression = aws.String("#BEAR = :bear AND #ASSOCTYPEID >= :assocTypeId")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"assocTypeID": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s^%s", input.StartingAfter.AssocType, input.StartingAfter.AssocID),
			},

			"bear": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Bear,
			},
		}
	}
	if len(input.FilterValues) > 0 && input.FilterExpression != "" {
		t.getThingWithMatchingKeyssByBearAndAssocTypeIDParseFilters(queryInput, input)
		queryInput.FilterExpression = aws.String(input.FilterExpression)
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithMatchingKeyss(queryOutput.Items)
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

func (t ThingWithMatchingKeysTable) deleteThingWithMatchingKeys(ctx context.Context, bear string, assocType string, assocID string) error {

	key, err := attributevalue.MarshalMap(ddbThingWithMatchingKeysPrimaryKey{
		Bear:        bear,
		AssocTypeID: fmt.Sprintf("%s^%s", assocType, assocID),
	})
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(t.TableName),
	})
	if err != nil {
		return err
	}

	return nil
}

func (t ThingWithMatchingKeysTable) getThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx context.Context, input db.GetThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
	if input.StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.StartingAt or input.StartingAfter")
	}
	if input.AssocType == "" {
		return fmt.Errorf("Hash key input.AssocType cannot be empty")
	}
	if input.AssocID == "" {
		return fmt.Errorf("Hash key input.AssocID cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("byAssoc"),
		ExpressionAttributeNames: map[string]string{
			"#ASSOCTYPEID": "assocTypeID",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":assocTypeId": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s^%s", input.AssocType, input.AssocID),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#ASSOCTYPEID = :assocTypeId")
	} else {
		queryInput.ExpressionAttributeNames["#CREATEDBEAR"] = "createdBear"
		queryInput.ExpressionAttributeValues[":createdBear"] = &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%s^%s", input.StartingAt.Created, input.StartingAt.Bear),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#ASSOCTYPEID = :assocTypeId AND #CREATEDBEAR <= :createdBear")
		} else {
			queryInput.KeyConditionExpression = aws.String("#ASSOCTYPEID = :assocTypeId AND #CREATEDBEAR >= :createdBear")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"createdBear": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s^%s", input.StartingAfter.Created, input.StartingAfter.Bear),
			},
			"assocTypeID": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s^%s", input.StartingAfter.AssocType, input.StartingAfter.AssocID),
			},
			"bear": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Bear,
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithMatchingKeyss(queryOutput.Items)
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
func (t ThingWithMatchingKeysTable) scanThingWithMatchingKeyssByAssocTypeIDAndCreatedBear(ctx context.Context, input db.ScanThingWithMatchingKeyssByAssocTypeIDAndCreatedBearInput, fn func(m *models.ThingWithMatchingKeys, lastThingWithMatchingKeys bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	scanInput.IndexName = aws.String("byAssoc")
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"bear": exclusiveStartKey["bear"],
			"assocTypeID": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s^%s", input.StartingAfter.AssocType, input.StartingAfter.AssocID),
			},
			"createdBear": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s^%s", input.StartingAfter.Created, input.StartingAfter.Bear),
			},
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithMatchingKeyss(out.Items)
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

// encodeThingWithMatchingKeys encodes a ThingWithMatchingKeys as a DynamoDB map of attribute values.
func encodeThingWithMatchingKeys(m models.ThingWithMatchingKeys) (map[string]types.AttributeValue, error) {
	// with composite attributes, marshal the model
	val, err := attributevalue.MarshalMapWithOptions(m, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(m.AssocID, "^") {
		return nil, fmt.Errorf("assocID cannot contain '^': %s", m.AssocID)
	}
	if strings.Contains(m.AssocType, "^") {
		return nil, fmt.Errorf("assocType cannot contain '^': %s", m.AssocType)
	}
	if strings.Contains(m.Bear, "^") {
		return nil, fmt.Errorf("bear cannot contain '^': %s", m.Bear)
	}
	// add in composite attributes
	primaryKey, err := attributevalue.MarshalMap(ddbThingWithMatchingKeysPrimaryKey{
		Bear:        m.Bear,
		AssocTypeID: fmt.Sprintf("%s^%s", m.AssocType, m.AssocID),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range primaryKey {
		val[k] = v
	}
	byAssoc, err := attributevalue.MarshalMap(ddbThingWithMatchingKeysGSIByAssoc{
		AssocTypeID: fmt.Sprintf("%s^%s", m.AssocType, m.AssocID),
		CreatedBear: fmt.Sprintf("%s^%s", m.Created, m.Bear),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range byAssoc {
		val[k] = v
	}
	return val, err
}

// decodeThingWithMatchingKeys translates a ThingWithMatchingKeys stored in DynamoDB to a ThingWithMatchingKeys struct.
func decodeThingWithMatchingKeys(m map[string]types.AttributeValue, out *models.ThingWithMatchingKeys) error {
	var ddbThingWithMatchingKeys ddbThingWithMatchingKeys
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithMatchingKeys); err != nil {
		return err
	}
	*out = ddbThingWithMatchingKeys.ThingWithMatchingKeys
	return nil
}

// decodeThingWithMatchingKeyss translates a list of ThingWithMatchingKeyss stored in DynamoDB to a slice of ThingWithMatchingKeys structs.
func decodeThingWithMatchingKeyss(ms []map[string]types.AttributeValue) ([]models.ThingWithMatchingKeys, error) {
	thingWithMatchingKeyss := make([]models.ThingWithMatchingKeys, len(ms))
	for i, m := range ms {
		var thingWithMatchingKeys models.ThingWithMatchingKeys
		if err := decodeThingWithMatchingKeys(m, &thingWithMatchingKeys); err != nil {
			return nil, err
		}
		thingWithMatchingKeyss[i] = thingWithMatchingKeys
	}
	return thingWithMatchingKeyss, nil
}
