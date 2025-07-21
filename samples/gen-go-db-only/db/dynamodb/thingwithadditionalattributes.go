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
var _ = errors.New("")
var _ = []types.AttributeValue{}

// ThingWithAdditionalAttributesTable represents the user-configurable properties of the ThingWithAdditionalAttributes table.
type ThingWithAdditionalAttributesTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithAdditionalAttributesPrimaryKey represents the primary key of a ThingWithAdditionalAttributes in DynamoDB.
type ddbThingWithAdditionalAttributesPrimaryKey struct {
	Name    string `dynamodbav:"name"`
	Version int64  `dynamodbav:"version"`
}

// ddbThingWithAdditionalAttributesGSIThingID represents the thingID GSI.
type ddbThingWithAdditionalAttributesGSIThingID struct {
	ID string `dynamodbav:"id"`
}

// ddbThingWithAdditionalAttributesGSINameCreatedAt represents the name-createdAt GSI.
type ddbThingWithAdditionalAttributesGSINameCreatedAt struct {
	Name      string          `dynamodbav:"name"`
	CreatedAt strfmt.DateTime `dynamodbav:"createdAt"`
}

// ddbThingWithAdditionalAttributesGSINameRangeNullable represents the name-rangeNullable GSI.
type ddbThingWithAdditionalAttributesGSINameRangeNullable struct {
	Name          string          `dynamodbav:"name"`
	RangeNullable strfmt.DateTime `dynamodbav:"rangeNullable"`
}

// ddbThingWithAdditionalAttributesGSINameHashNullable represents the name-hashNullable GSI.
type ddbThingWithAdditionalAttributesGSINameHashNullable struct {
	HashNullable string `dynamodbav:"hashNullable"`
	Name         string `dynamodbav:"name"`
}

// ddbThingWithAdditionalAttributes represents a ThingWithAdditionalAttributes as stored in DynamoDB.
type ddbThingWithAdditionalAttributes struct {
	models.ThingWithAdditionalAttributes
}

func (t ThingWithAdditionalAttributesTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("createdAt"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("hashNullable"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("name"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("rangeNullable"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: types.ScalarAttributeType("N"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("thingID"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("id"),
						KeyType:       types.KeyTypeHash,
					},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("name-createdAt"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("name"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("createdAt"),
						KeyType:       types.KeyTypeRange,
					},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("name-rangeNullable"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("name"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("rangeNullable"),
						KeyType:       types.KeyTypeRange,
					},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("name-hashNullable"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("hashNullable"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("name"),
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

func (t ThingWithAdditionalAttributesTable) saveThingWithAdditionalAttributes(ctx context.Context, m models.ThingWithAdditionalAttributes) error {
	data, err := encodeThingWithAdditionalAttributes(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#NAME":    "name",
			"#VERSION": "version",
		},
		ConditionExpression: aws.String(
			"" +
				"" +
				"attribute_not_exists(#NAME)" +
				"" +
				" AND " +
				"attribute_not_exists(#VERSION)" +
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
			return db.ErrThingWithAdditionalAttributesAlreadyExists{
				Name:    m.Name,
				Version: m.Version,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributes(ctx context.Context, name string, version int64) (*models.ThingWithAdditionalAttributes, error) {
	key, err := attributevalue.MarshalMapWithOptions(ddbThingWithAdditionalAttributesPrimaryKey{
		Name:    name,
		Version: version,
	}, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
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
		return nil, db.ErrThingWithAdditionalAttributesNotFound{
			Name:    name,
			Version: version,
		}
	}

	var m models.ThingWithAdditionalAttributes
	if err := decodeThingWithAdditionalAttributes(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithAdditionalAttributesTable) scanThingWithAdditionalAttributess(ctx context.Context, input db.ScanThingWithAdditionalAttributessInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMapWithOptions(input.StartingAfter, func(o *attributevalue.EncoderOptions) {
			o.TagKey = "json"
		})
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name":    exclusiveStartKey["name"],
			"version": exclusiveStartKey["version"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithAdditionalAttributess(out.Items)
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

func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributessByNameAndVersionParseFilters(queryInput *dynamodb.QueryInput, input db.GetThingWithAdditionalAttributessByNameAndVersionInput) {
	for _, filterValue := range input.FilterValues {
		switch filterValue.AttributeName {
		case db.ThingWithAdditionalAttributesAdditionalBAttribute:
			queryInput.ExpressionAttributeNames["#ADDITIONALBATTRIBUTE"] = string(db.ThingWithAdditionalAttributesAdditionalBAttribute)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesAdditionalBAttribute), i)] = &types.AttributeValueMemberB{
					Value: attributeValue.([]byte),
				}
			}
		case db.ThingWithAdditionalAttributesAdditionalNAttribute:
			queryInput.ExpressionAttributeNames["#ADDITIONALNATTRIBUTE"] = string(db.ThingWithAdditionalAttributesAdditionalNAttribute)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesAdditionalNAttribute), i)] = &types.AttributeValueMemberN{
					Value: fmt.Sprint(attributeValue.(int64)),
				}
			}
		case db.ThingWithAdditionalAttributesAdditionalSAttribute:
			queryInput.ExpressionAttributeNames["#ADDITIONALSATTRIBUTE"] = string(db.ThingWithAdditionalAttributesAdditionalSAttribute)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesAdditionalSAttribute), i)] = &types.AttributeValueMemberS{
					Value: attributeValue.(string),
				}
			}
		case db.ThingWithAdditionalAttributesCreatedAt:
			queryInput.ExpressionAttributeNames["#CREATEDAT"] = string(db.ThingWithAdditionalAttributesCreatedAt)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesCreatedAt), i)] = &types.AttributeValueMemberS{
					Value: datetimeToDynamoTimeString(attributeValue.(strfmt.DateTime)),
				}
			}
		case db.ThingWithAdditionalAttributesHashNullable:
			queryInput.ExpressionAttributeNames["#HASHNULLABLE"] = string(db.ThingWithAdditionalAttributesHashNullable)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesHashNullable), i)] = &types.AttributeValueMemberS{
					Value: attributeValue.(string),
				}
			}
		case db.ThingWithAdditionalAttributesID:
			queryInput.ExpressionAttributeNames["#ID"] = string(db.ThingWithAdditionalAttributesID)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesID), i)] = &types.AttributeValueMemberS{
					Value: attributeValue.(string),
				}
			}
		case db.ThingWithAdditionalAttributesRangeNullable:
			queryInput.ExpressionAttributeNames["#RANGENULLABLE"] = string(db.ThingWithAdditionalAttributesRangeNullable)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesRangeNullable), i)] = &types.AttributeValueMemberS{
					Value: datetimeToDynamoTimeString(attributeValue.(strfmt.DateTime)),
				}
			}
		}
	}
}

func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributessByNameAndVersion(ctx context.Context, input db.GetThingWithAdditionalAttributessByNameAndVersionInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	if input.VersionStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.VersionStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		ExpressionAttributeNames: map[string]string{
			"#NAME": "name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{
				Value: input.Name,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.VersionStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		queryInput.ExpressionAttributeNames["#VERSION"] = "version"
		queryInput.ExpressionAttributeValues[":version"] = &types.AttributeValueMemberN{
			Value: fmt.Sprintf("%d", *input.VersionStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #VERSION <= :version")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #VERSION >= :version")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"version": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", input.StartingAfter.Version),
			},

			"name": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Name,
			},
		}
	}
	if len(input.FilterValues) > 0 && input.FilterExpression != "" {
		t.getThingWithAdditionalAttributessByNameAndVersionParseFilters(queryInput, input)
		queryInput.FilterExpression = aws.String(input.FilterExpression)
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithAdditionalAttributess(queryOutput.Items)
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

func (t ThingWithAdditionalAttributesTable) deleteThingWithAdditionalAttributes(ctx context.Context, name string, version int64) error {

	key, err := attributevalue.MarshalMapWithOptions(ddbThingWithAdditionalAttributesPrimaryKey{
		Name:    name,
		Version: version,
	}, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
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

func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributesByID(ctx context.Context, id string) (*models.ThingWithAdditionalAttributes, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("thingID"),
		ExpressionAttributeNames: map[string]string{
			"#ID": "id",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
		KeyConditionExpression: aws.String("#ID = :id"),
	}

	queryOutput, err := t.DynamoDBAPI.Query(ctx, queryInput)
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundErr) {
			return nil, fmt.Errorf("table or index not found: %s", t.TableName)
		}
		return nil, err
	}
	if len(queryOutput.Items) == 0 {
		return nil, db.ErrThingWithAdditionalAttributesByIDNotFound{
			ID: id,
		}
	}

	var thingWithAdditionalAttributes models.ThingWithAdditionalAttributes
	if err := decodeThingWithAdditionalAttributes(queryOutput.Items[0], &thingWithAdditionalAttributes); err != nil {
		return nil, err
	}
	return &thingWithAdditionalAttributes, nil
}
func (t ThingWithAdditionalAttributesTable) scanThingWithAdditionalAttributessByID(ctx context.Context, input db.ScanThingWithAdditionalAttributessByIDInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	scanInput.IndexName = aws.String("thingID")
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMapWithOptions(input.StartingAfter, func(o *attributevalue.EncoderOptions) {
			o.TagKey = "json"
		})
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name":    exclusiveStartKey["name"],
			"version": exclusiveStartKey["version"],
			"id":      exclusiveStartKey["id"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithAdditionalAttributess(out.Items)
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
func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributessByNameAndCreatedAt(ctx context.Context, input db.GetThingWithAdditionalAttributessByNameAndCreatedAtInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	if input.CreatedAtStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.CreatedAtStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("name-createdAt"),
		ExpressionAttributeNames: map[string]string{
			"#NAME": "name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{
				Value: input.Name,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.CreatedAtStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		queryInput.ExpressionAttributeNames["#CREATEDAT"] = "createdAt"
		queryInput.ExpressionAttributeValues[":createdAt"] = &types.AttributeValueMemberS{
			Value: datetimeToDynamoTimeString(*input.CreatedAtStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #CREATEDAT <= :createdAt")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #CREATEDAT >= :createdAt")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"createdAt": &types.AttributeValueMemberS{
				Value: datetimeToDynamoTimeString(input.StartingAfter.CreatedAt),
			},
			"name": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Name,
			},
			"version": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", input.StartingAfter.Version),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithAdditionalAttributess(queryOutput.Items)
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
func (t ThingWithAdditionalAttributesTable) scanThingWithAdditionalAttributessByNameAndCreatedAt(ctx context.Context, input db.ScanThingWithAdditionalAttributessByNameAndCreatedAtInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	scanInput.IndexName = aws.String("name-createdAt")
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMapWithOptions(input.StartingAfter, func(o *attributevalue.EncoderOptions) {
			o.TagKey = "json"
		})
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name":      exclusiveStartKey["name"],
			"version":   exclusiveStartKey["version"],
			"createdAt": exclusiveStartKey["createdAt"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithAdditionalAttributess(out.Items)
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
func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributessByNameAndRangeNullable(ctx context.Context, input db.GetThingWithAdditionalAttributessByNameAndRangeNullableInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	if input.RangeNullableStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.RangeNullableStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("name-rangeNullable"),
		ExpressionAttributeNames: map[string]string{
			"#NAME": "name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{
				Value: input.Name,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.RangeNullableStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		queryInput.ExpressionAttributeNames["#RANGENULLABLE"] = "rangeNullable"
		queryInput.ExpressionAttributeValues[":rangeNullable"] = &types.AttributeValueMemberS{
			Value: datetimeToDynamoTimeString(*input.RangeNullableStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #RANGENULLABLE <= :rangeNullable")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #RANGENULLABLE >= :rangeNullable")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"rangeNullable": &types.AttributeValueMemberS{
				Value: datetimePtrToDynamoTimeString(input.StartingAfter.RangeNullable),
			},
			"name": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Name,
			},
			"version": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", input.StartingAfter.Version),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithAdditionalAttributess(queryOutput.Items)
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
func (t ThingWithAdditionalAttributesTable) scanThingWithAdditionalAttributessByNameAndRangeNullable(ctx context.Context, input db.ScanThingWithAdditionalAttributessByNameAndRangeNullableInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	scanInput.IndexName = aws.String("name-rangeNullable")
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMapWithOptions(input.StartingAfter, func(o *attributevalue.EncoderOptions) {
			o.TagKey = "json"
		})
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name":          exclusiveStartKey["name"],
			"version":       exclusiveStartKey["version"],
			"rangeNullable": exclusiveStartKey["rangeNullable"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithAdditionalAttributess(out.Items)
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
func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributessByHashNullableAndName(ctx context.Context, input db.GetThingWithAdditionalAttributessByHashNullableAndNameInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	if input.NameStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.NameStartingAt or input.StartingAfter")
	}
	if input.HashNullable == "" {
		return fmt.Errorf("Hash key input.HashNullable cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("name-hashNullable"),
		ExpressionAttributeNames: map[string]string{
			"#HASHNULLABLE": "hashNullable",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":hashNullable": &types.AttributeValueMemberS{
				Value: input.HashNullable,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.NameStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#HASHNULLABLE = :hashNullable")
	} else {
		queryInput.ExpressionAttributeNames["#NAME"] = "name"
		queryInput.ExpressionAttributeValues[":name"] = &types.AttributeValueMemberS{
			Value: string(*input.NameStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#HASHNULLABLE = :hashNullable AND #NAME <= :name")
		} else {
			queryInput.KeyConditionExpression = aws.String("#HASHNULLABLE = :hashNullable AND #NAME >= :name")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{
				Value: string(input.StartingAfter.Name),
			},
			"hashNullable": &types.AttributeValueMemberS{
				Value: *input.StartingAfter.HashNullable,
			},
			"version": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", input.StartingAfter.Version),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithAdditionalAttributess(queryOutput.Items)
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

// encodeThingWithAdditionalAttributes encodes a ThingWithAdditionalAttributes as a DynamoDB map of attribute values.
func encodeThingWithAdditionalAttributes(m models.ThingWithAdditionalAttributes) (map[string]types.AttributeValue, error) {
	// no composite attributes, marshal the model with the json tag
	val, err := attributevalue.MarshalMapWithOptions(m, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

// decodeThingWithAdditionalAttributes translates a ThingWithAdditionalAttributes stored in DynamoDB to a ThingWithAdditionalAttributes struct.
func decodeThingWithAdditionalAttributes(m map[string]types.AttributeValue, out *models.ThingWithAdditionalAttributes) error {
	var ddbThingWithAdditionalAttributes ddbThingWithAdditionalAttributes
	if err := attributevalue.UnmarshalMapWithOptions(m, &ddbThingWithAdditionalAttributes, func(o *attributevalue.DecoderOptions) {
		o.TagKey = "json"
	}); err != nil {
		return err
	}
	*out = ddbThingWithAdditionalAttributes.ThingWithAdditionalAttributes
	return nil
}

// decodeThingWithAdditionalAttributess translates a list of ThingWithAdditionalAttributess stored in DynamoDB to a slice of ThingWithAdditionalAttributes structs.
func decodeThingWithAdditionalAttributess(ms []map[string]types.AttributeValue) ([]models.ThingWithAdditionalAttributes, error) {
	thingWithAdditionalAttributess := make([]models.ThingWithAdditionalAttributes, len(ms))
	for i, m := range ms {
		var thingWithAdditionalAttributes models.ThingWithAdditionalAttributes
		if err := decodeThingWithAdditionalAttributes(m, &thingWithAdditionalAttributes); err != nil {
			return nil, err
		}
		thingWithAdditionalAttributess[i] = thingWithAdditionalAttributes
	}
	return thingWithAdditionalAttributess, nil
}
