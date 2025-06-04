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

// ThingTable represents the user-configurable properties of the Thing table.
type ThingTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingPrimaryKey represents the primary key of a Thing in DynamoDB.
type ddbThingPrimaryKey struct {
	Name    string `dynamodbav:"name"`
	Version int64  `dynamodbav:"version"`
}

// ddbThingGSIThingID represents the thingID GSI.
type ddbThingGSIThingID struct {
	ID string `dynamodbav:"id"`
}

// ddbThingGSINameCreatedAt represents the name-createdAt GSI.
type ddbThingGSINameCreatedAt struct {
	Name      string          `dynamodbav:"name"`
	CreatedAt strfmt.DateTime `dynamodbav:"createdAt"`
}

// ddbThingGSINameRangeNullable represents the name-rangeNullable GSI.
type ddbThingGSINameRangeNullable struct {
	Name          string          `dynamodbav:"name"`
	RangeNullable strfmt.DateTime `dynamodbav:"rangeNullable"`
}

// ddbThingGSINameHashNullable represents the name-hashNullable GSI.
type ddbThingGSINameHashNullable struct {
	HashNullable string `dynamodbav:"hashNullable"`
	Name         string `dynamodbav:"name"`
}

// ddbThing represents a Thing as stored in DynamoDB.
type ddbThing struct {
	models.Thing
}

func (t ThingTable) create(ctx context.Context) error {
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

func (t ThingTable) saveThing(ctx context.Context, m models.Thing) error {
	data, err := encodeThing(m)
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
		ConditionExpression: aws.String("attribute_not_exists(#NAME) AND attribute_not_exists(#VERSION)"),
	})
	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException
		var conditionalCheckFailedErr *types.ConditionalCheckFailedException
		if errors.As(err, &resourceNotFoundErr) {
			return fmt.Errorf("table or index not found: %s", t.TableName)
		}
		if errors.As(err, &conditionalCheckFailedErr) {
			return db.ErrThingAlreadyExists{
				Name:    m.Name,
				Version: m.Version,
			}
		}
		return err
	}
	return nil
}

func (t ThingTable) getThing(ctx context.Context, name string, version int64) (*models.Thing, error) {
	// swad-get-7
	key, err := attributevalue.MarshalMap(ddbThingPrimaryKey{
		Name:    name,
		Version: version,
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
		return nil, db.ErrThingNotFound{
			Name:    name,
			Version: version,
		}
	}

	var m models.Thing
	if err := decodeThing(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingTable) scanThings(ctx context.Context, input db.ScanThingsInput, fn func(m *models.Thing, lastThing bool) bool) error {
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
			"name":    exclusiveStartKey["name"],
			"version": exclusiveStartKey["version"],
		}
	}
	totalRecordsProcessed := int32(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThings(out.Items)
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

func (t ThingTable) getThingsByNameAndVersionParseFilters(queryInput *dynamodb.QueryInput, input db.GetThingsByNameAndVersionInput) {
	// swad-get-11
	for _, filterValue := range input.FilterValues {
		switch filterValue.AttributeName {
		case db.ThingCreatedAt:
			queryInput.ExpressionAttributeNames["#CREATEDAT"] = string(db.ThingCreatedAt)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingCreatedAt), i)] = &types.AttributeValueMemberS{
					Value: datetimeToDynamoTimeString(attributeValue.(strfmt.DateTime)),
				}
			}
		case db.ThingHashNullable:
			queryInput.ExpressionAttributeNames["#HASHNULLABLE"] = string(db.ThingHashNullable)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingHashNullable), i)] = &types.AttributeValueMemberS{
					Value: attributeValue.(string),
				}
			}
		case db.ThingID:
			queryInput.ExpressionAttributeNames["#ID"] = string(db.ThingID)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingID), i)] = &types.AttributeValueMemberS{
					Value: attributeValue.(string),
				}
			}
		case db.ThingRangeNullable:
			queryInput.ExpressionAttributeNames["#RANGENULLABLE"] = string(db.ThingRangeNullable)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingRangeNullable), i)] = &types.AttributeValueMemberS{
					Value: datetimeToDynamoTimeString(attributeValue.(strfmt.DateTime)),
				}
			}
		}
	}
}

func (t ThingTable) getThingsByNameAndVersion(ctx context.Context, input db.GetThingsByNameAndVersionInput, fn func(m *models.Thing, lastThing bool) bool) error {
	// swad-get-2
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
		queryInput.Limit = input.Limit
	}
	if input.VersionStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		// swad-get-21
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
	// swad-get-22
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"version": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", input.StartingAfter.Version),
			},

			// swad-get-223
			"name": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Name,
			},
		}
	}
	if len(input.FilterValues) > 0 && input.FilterExpression != "" {
		t.getThingsByNameAndVersionParseFilters(queryInput, input)
		queryInput.FilterExpression = aws.String(input.FilterExpression)
	}

	totalRecordsProcessed := int32(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThings(queryOutput.Items)
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

func (t ThingTable) deleteThing(ctx context.Context, name string, version int64) error {

	key, err := attributevalue.MarshalMap(ddbThingPrimaryKey{
		Name:    name,
		Version: version,
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

func (t ThingTable) getThingByID(ctx context.Context, id string) (*models.Thing, error) {
	// swad-get-8
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
		return nil, db.ErrThingByIDNotFound{
			ID: id,
		}
	}

	var thing models.Thing
	if err := decodeThing(queryOutput.Items[0], &thing); err != nil {
		return nil, err
	}
	return &thing, nil
}
func (t ThingTable) scanThingsByID(ctx context.Context, input db.ScanThingsByIDInput, fn func(m *models.Thing, lastThing bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("thingID"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
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
	totalRecordsProcessed := int32(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThings(out.Items)
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
func (t ThingTable) getThingsByNameAndCreatedAt(ctx context.Context, input db.GetThingsByNameAndCreatedAtInput, fn func(m *models.Thing, lastThing bool) bool) error {
	// swad-get-33
	if input.CreatedAtStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.CreatedAtStartingAt or input.StartingAfter")
	}
	// swad-get-33f
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	// swad-get-331
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("name-createdAt"),
		ExpressionAttributeNames: map[string]string{
			"#NAME": "name",
		},
		// swad-get-3312
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{
				// swad-get-33e
				Value: input.Name,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	// swad-get-332
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.CreatedAtStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		// swad-get-333
		queryInput.ExpressionAttributeNames["#CREATEDAT"] = "createdAt"

		// swad-get-3331a
		queryInput.ExpressionAttributeValues[":createdAt"] = &types.AttributeValueMemberS{
			Value: datetimeToDynamoTimeString(*input.CreatedAtStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #CREATEDAT <= :createdAt")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #CREATEDAT >= :createdAt")
		}
	}
	// swad-get-334
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"createdAt": &types.AttributeValueMemberS{
				Value: datetimeToDynamoTimeString(input.StartingAfter.CreatedAt),
			},
			// swad-get-3341
			"name": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Name,
			},
			// swad-get-3342
			"version": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", input.StartingAfter.Version),
			},

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
		items, err := decodeThings(queryOutput.Items)
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
func (t ThingTable) scanThingsByNameAndCreatedAt(ctx context.Context, input db.ScanThingsByNameAndCreatedAtInput, fn func(m *models.Thing, lastThing bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("name-createdAt"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
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
	totalRecordsProcessed := int32(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThings(out.Items)
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
func (t ThingTable) getThingsByNameAndRangeNullable(ctx context.Context, input db.GetThingsByNameAndRangeNullableInput, fn func(m *models.Thing, lastThing bool) bool) error {
	// swad-get-33
	if input.RangeNullableStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.RangeNullableStartingAt or input.StartingAfter")
	}
	// swad-get-33f
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	// swad-get-331
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("name-rangeNullable"),
		ExpressionAttributeNames: map[string]string{
			"#NAME": "name",
		},
		// swad-get-3312
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{
				// swad-get-33e
				Value: input.Name,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	// swad-get-332
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.RangeNullableStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		// swad-get-333
		queryInput.ExpressionAttributeNames["#RANGENULLABLE"] = "rangeNullable"

		// swad-get-3331a
		queryInput.ExpressionAttributeValues[":rangeNullable"] = &types.AttributeValueMemberS{
			Value: datetimeToDynamoTimeString(*input.RangeNullableStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #RANGENULLABLE <= :rangeNullable")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #RANGENULLABLE >= :rangeNullable")
		}
	}
	// swad-get-334
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"rangeNullable": &types.AttributeValueMemberS{
				Value: datetimePtrToDynamoTimeString(input.StartingAfter.RangeNullable),
			},
			// swad-get-3341
			"name": &types.AttributeValueMemberS{
				Value: input.StartingAfter.Name,
			},
			// swad-get-3342
			"version": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", input.StartingAfter.Version),
			},

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
		items, err := decodeThings(queryOutput.Items)
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
func (t ThingTable) scanThingsByNameAndRangeNullable(ctx context.Context, input db.ScanThingsByNameAndRangeNullableInput, fn func(m *models.Thing, lastThing bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("name-rangeNullable"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
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
	totalRecordsProcessed := int32(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThings(out.Items)
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
func (t ThingTable) getThingsByHashNullableAndName(ctx context.Context, input db.GetThingsByHashNullableAndNameInput, fn func(m *models.Thing, lastThing bool) bool) error {
	// swad-get-33
	if input.NameStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.NameStartingAt or input.StartingAfter")
	}
	// swad-get-33f
	if input.HashNullable == "" {
		return fmt.Errorf("Hash key input.HashNullable cannot be empty")
	}
	// swad-get-331
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("name-hashNullable"),
		ExpressionAttributeNames: map[string]string{
			"#HASHNULLABLE": "hashNullable",
		},
		// swad-get-3312
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":hashNullable": &types.AttributeValueMemberS{
				// swad-get-33e
				Value: input.HashNullable,
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	// swad-get-332
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.NameStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#HASHNULLABLE = :hashNullable")
	} else {
		// swad-get-333
		queryInput.ExpressionAttributeNames["#NAME"] = "name"

		// swad-get-3331a
		queryInput.ExpressionAttributeValues[":name"] = &types.AttributeValueMemberS{
			Value: string(*input.NameStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#HASHNULLABLE = :hashNullable AND #NAME <= :name")
		} else {
			queryInput.KeyConditionExpression = aws.String("#HASHNULLABLE = :hashNullable AND #NAME >= :name")
		}
	}
	// swad-get-334
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{
				Value: string(input.StartingAfter.Name),
			},
			// swad-get-3341
			"hashNullable": &types.AttributeValueMemberS{
				Value: *input.StartingAfter.HashNullable,
			},
			// swad-get-3342
			"version": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", input.StartingAfter.Version),
			},

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
		items, err := decodeThings(queryOutput.Items)
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

// encodeThing encodes a Thing as a DynamoDB map of attribute values.
func encodeThing(m models.Thing) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(ddbThing{
		Thing: m,
	})
}

// decodeThing translates a Thing stored in DynamoDB to a Thing struct.
func decodeThing(m map[string]types.AttributeValue, out *models.Thing) error {
	// swad-decode-1
	var ddbThing ddbThing
	if err := attributevalue.UnmarshalMap(m, &ddbThing); err != nil {
		return err
	}
	*out = ddbThing.Thing
	return nil
}

// decodeThings translates a list of Things stored in DynamoDB to a slice of Thing structs.
func decodeThings(ms []map[string]types.AttributeValue) ([]models.Thing, error) {
	things := make([]models.Thing, len(ms))
	for i, m := range ms {
		var thing models.Thing
		if err := decodeThing(m, &thing); err != nil {
			return nil, err
		}
		things[i] = thing
	}
	return things, nil
}
