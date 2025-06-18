package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Clever/wag/samples/gen-go-db-only/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db-only/db"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}
var _ = errors.New("")
var _ = []types.AttributeValue{}

// ThingWithMultiUseCompositeAttributeTable represents the user-configurable properties of the ThingWithMultiUseCompositeAttribute table.
type ThingWithMultiUseCompositeAttributeTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithMultiUseCompositeAttributePrimaryKey represents the primary key of a ThingWithMultiUseCompositeAttribute in DynamoDB.
type ddbThingWithMultiUseCompositeAttributePrimaryKey struct {
	One string `dynamodbav:"one"`
}

// ddbThingWithMultiUseCompositeAttributeGSIThreeIndex represents the threeIndex GSI.
type ddbThingWithMultiUseCompositeAttributeGSIThreeIndex struct {
	Three  string `dynamodbav:"three"`
	OneTwo string `dynamodbav:"one_two"`
}

// ddbThingWithMultiUseCompositeAttributeGSIFourIndex represents the fourIndex GSI.
type ddbThingWithMultiUseCompositeAttributeGSIFourIndex struct {
	Four   string `dynamodbav:"four"`
	OneTwo string `dynamodbav:"one_two"`
}

// ddbThingWithMultiUseCompositeAttribute represents a ThingWithMultiUseCompositeAttribute as stored in DynamoDB.
type ddbThingWithMultiUseCompositeAttribute struct {
	models.ThingWithMultiUseCompositeAttribute
}

func (t ThingWithMultiUseCompositeAttributeTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("four"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("one"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("one_two"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("three"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("one"),
				KeyType:       types.KeyTypeHash,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("threeIndex"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("three"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("one_two"),
						KeyType:       types.KeyTypeRange,
					},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("fourIndex"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("four"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("one_two"),
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

func (t ThingWithMultiUseCompositeAttributeTable) saveThingWithMultiUseCompositeAttribute(ctx context.Context, m models.ThingWithMultiUseCompositeAttribute) error {
	data, err := encodeThingWithMultiUseCompositeAttribute(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
	})
	return err
}

func (t ThingWithMultiUseCompositeAttributeTable) getThingWithMultiUseCompositeAttribute(ctx context.Context, one string) (*models.ThingWithMultiUseCompositeAttribute, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithMultiUseCompositeAttributePrimaryKey{
		One: one,
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
		return nil, db.ErrThingWithMultiUseCompositeAttributeNotFound{
			One: one,
		}
	}

	var m models.ThingWithMultiUseCompositeAttribute
	if err := decodeThingWithMultiUseCompositeAttribute(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithMultiUseCompositeAttributeTable) scanThingWithMultiUseCompositeAttributes(ctx context.Context, input db.ScanThingWithMultiUseCompositeAttributesInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
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
			"one": exclusiveStartKey["one"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithMultiUseCompositeAttributes(out.Items)
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

func (t ThingWithMultiUseCompositeAttributeTable) deleteThingWithMultiUseCompositeAttribute(ctx context.Context, one string) error {

	key, err := attributevalue.MarshalMap(ddbThingWithMultiUseCompositeAttributePrimaryKey{
		One: one,
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

func (t ThingWithMultiUseCompositeAttributeTable) getThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx context.Context, input db.GetThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	if input.StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.StartingAt or input.StartingAfter")
	}
	if input.Three == "" {
		return fmt.Errorf("Hash key input.Three cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("threeIndex"),
		ExpressionAttributeNames: map[string]string{
			"#THREE": "three",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":three": &types.AttributeValueMemberS{
				Value: input.Three,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#THREE = :three")
	} else {
		queryInput.ExpressionAttributeNames["#ONE_TWO"] = "one_two"
		queryInput.ExpressionAttributeValues[":oneTwo"] = &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%s_%s", input.StartingAt.One, input.StartingAt.Two),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#THREE = :three AND #ONE_TWO <= :oneTwo")
		} else {
			queryInput.KeyConditionExpression = aws.String("#THREE = :three AND #ONE_TWO >= :oneTwo")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"one_two": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s_%s", *input.StartingAfter.One, *input.StartingAfter.Two),
			},
			"three": &types.AttributeValueMemberS{
				Value: *input.StartingAfter.Three,
			},
			"one": &types.AttributeValueMemberS{
				Value: *input.StartingAfter.One,
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithMultiUseCompositeAttributes(queryOutput.Items)
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
func (t ThingWithMultiUseCompositeAttributeTable) scanThingWithMultiUseCompositeAttributesByThreeAndOneTwo(ctx context.Context, input db.ScanThingWithMultiUseCompositeAttributesByThreeAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	scanInput.IndexName = aws.String("threeIndex")
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"one":   exclusiveStartKey["one"],
			"three": exclusiveStartKey["three"],
			"one_two": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s_%s", *input.StartingAfter.One, *input.StartingAfter.Two),
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

		items, err := decodeThingWithMultiUseCompositeAttributes(out.Items)
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
func (t ThingWithMultiUseCompositeAttributeTable) getThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx context.Context, input db.GetThingWithMultiUseCompositeAttributesByFourAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	if input.StartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.StartingAt or input.StartingAfter")
	}
	if input.Four == "" {
		return fmt.Errorf("Hash key input.Four cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("fourIndex"),
		ExpressionAttributeNames: map[string]string{
			"#FOUR": "four",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":four": &types.AttributeValueMemberS{
				Value: input.Four,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#FOUR = :four")
	} else {
		queryInput.ExpressionAttributeNames["#ONE_TWO"] = "one_two"
		queryInput.ExpressionAttributeValues[":oneTwo"] = &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%s_%s", input.StartingAt.One, input.StartingAt.Two),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#FOUR = :four AND #ONE_TWO <= :oneTwo")
		} else {
			queryInput.KeyConditionExpression = aws.String("#FOUR = :four AND #ONE_TWO >= :oneTwo")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"one_two": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s_%s", *input.StartingAfter.One, *input.StartingAfter.Two),
			},
			"four": &types.AttributeValueMemberS{
				Value: *input.StartingAfter.Four,
			},
			"one": &types.AttributeValueMemberS{
				Value: *input.StartingAfter.One,
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithMultiUseCompositeAttributes(queryOutput.Items)
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
func (t ThingWithMultiUseCompositeAttributeTable) scanThingWithMultiUseCompositeAttributesByFourAndOneTwo(ctx context.Context, input db.ScanThingWithMultiUseCompositeAttributesByFourAndOneTwoInput, fn func(m *models.ThingWithMultiUseCompositeAttribute, lastThingWithMultiUseCompositeAttribute bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	scanInput.IndexName = aws.String("fourIndex")
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"one":  exclusiveStartKey["one"],
			"four": exclusiveStartKey["four"],
			"one_two": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s_%s", *input.StartingAfter.One, *input.StartingAfter.Two),
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

		items, err := decodeThingWithMultiUseCompositeAttributes(out.Items)
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

// encodeThingWithMultiUseCompositeAttribute encodes a ThingWithMultiUseCompositeAttribute as a DynamoDB map of attribute values.
func encodeThingWithMultiUseCompositeAttribute(m models.ThingWithMultiUseCompositeAttribute) (map[string]types.AttributeValue, error) {
	// with composite attributes, marshal the model
	val, err := attributevalue.MarshalMapWithOptions(m, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(*m.One, "_") {
		return nil, fmt.Errorf("one cannot contain '_': %s", *m.One)
	}
	if strings.Contains(*m.Two, "_") {
		return nil, fmt.Errorf("two cannot contain '_': %s", *m.Two)
	}
	// add in composite attributes
	threeIndex, err := attributevalue.MarshalMap(ddbThingWithMultiUseCompositeAttributeGSIThreeIndex{
		Three:  *m.Three,
		OneTwo: fmt.Sprintf("%s_%s", *m.One, *m.Two),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range threeIndex {
		val[k] = v
	}
	fourIndex, err := attributevalue.MarshalMap(ddbThingWithMultiUseCompositeAttributeGSIFourIndex{
		Four:   *m.Four,
		OneTwo: fmt.Sprintf("%s_%s", *m.One, *m.Two),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range fourIndex {
		val[k] = v
	}
	return val, err
}

// decodeThingWithMultiUseCompositeAttribute translates a ThingWithMultiUseCompositeAttribute stored in DynamoDB to a ThingWithMultiUseCompositeAttribute struct.
func decodeThingWithMultiUseCompositeAttribute(m map[string]types.AttributeValue, out *models.ThingWithMultiUseCompositeAttribute) error {
	var ddbThingWithMultiUseCompositeAttribute ddbThingWithMultiUseCompositeAttribute
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithMultiUseCompositeAttribute); err != nil {
		return err
	}
	*out = ddbThingWithMultiUseCompositeAttribute.ThingWithMultiUseCompositeAttribute
	return nil
}

// decodeThingWithMultiUseCompositeAttributes translates a list of ThingWithMultiUseCompositeAttributes stored in DynamoDB to a slice of ThingWithMultiUseCompositeAttribute structs.
func decodeThingWithMultiUseCompositeAttributes(ms []map[string]types.AttributeValue) ([]models.ThingWithMultiUseCompositeAttribute, error) {
	thingWithMultiUseCompositeAttributes := make([]models.ThingWithMultiUseCompositeAttribute, len(ms))
	for i, m := range ms {
		var thingWithMultiUseCompositeAttribute models.ThingWithMultiUseCompositeAttribute
		if err := decodeThingWithMultiUseCompositeAttribute(m, &thingWithMultiUseCompositeAttribute); err != nil {
			return nil, err
		}
		thingWithMultiUseCompositeAttributes[i] = thingWithMultiUseCompositeAttribute
	}
	return thingWithMultiUseCompositeAttributes, nil
}
