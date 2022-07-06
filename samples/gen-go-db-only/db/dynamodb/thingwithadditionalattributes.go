package dynamodb

import (
	"context"
	"fmt"

	"github.com/Clever/wag/samples/v8/gen-go-db-only/db"
	"github.com/Clever/wag/samples/v8/gen-go-db-only/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// ThingWithAdditionalAttributesTable represents the user-configurable properties of the ThingWithAdditionalAttributes table.
type ThingWithAdditionalAttributesTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
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

func (t ThingWithAdditionalAttributesTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-thing-with-additional-attributess", t.Prefix)
}

func (t ThingWithAdditionalAttributesTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("createdAt"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("hashNullable"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("rangeNullable"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("thingID"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("id"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("name-createdAt"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("name"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("createdAt"),
						KeyType:       aws.String(dynamodb.KeyTypeRange),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("name-rangeNullable"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("name"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("rangeNullable"),
						KeyType:       aws.String(dynamodb.KeyTypeRange),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("name-hashNullable"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("hashNullable"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("name"),
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

func (t ThingWithAdditionalAttributesTable) saveThingWithAdditionalAttributes(ctx context.Context, m models.ThingWithAdditionalAttributes) error {
	data, err := encodeThingWithAdditionalAttributes(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
		ExpressionAttributeNames: map[string]*string{
			"#NAME":    aws.String("name"),
			"#VERSION": aws.String("version"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#NAME) AND attribute_not_exists(#VERSION)"),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return db.ErrThingWithAdditionalAttributesAlreadyExists{
					Name:    m.Name,
					Version: m.Version,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	return nil
}

func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributes(ctx context.Context, name string, version int64) (*models.ThingWithAdditionalAttributes, error) {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithAdditionalAttributesPrimaryKey{
		Name:    name,
		Version: version,
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
			"name":    exclusiveStartKey["name"],
			"version": exclusiveStartKey["version"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithAdditionalAttributess(out.Items)
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

func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributessByNameAndVersionParseFilters(queryInput *dynamodb.QueryInput, input db.GetThingWithAdditionalAttributessByNameAndVersionInput) {
	for _, filterValue := range input.FilterValues {
		switch filterValue.AttributeName {
		case db.ThingWithAdditionalAttributesAdditionalBAttribute:
			queryInput.ExpressionAttributeNames["#ADDITIONALBATTRIBUTE"] = aws.String(string(db.ThingWithAdditionalAttributesAdditionalBAttribute))
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesAdditionalBAttribute), i)] = &dynamodb.AttributeValue{
					B: attributeValue.([]byte),
				}
			}
		case db.ThingWithAdditionalAttributesAdditionalNAttribute:
			queryInput.ExpressionAttributeNames["#ADDITIONALNATTRIBUTE"] = aws.String(string(db.ThingWithAdditionalAttributesAdditionalNAttribute))
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesAdditionalNAttribute), i)] = &dynamodb.AttributeValue{
					N: aws.String(fmt.Sprint(attributeValue.(int64))),
				}
			}
		case db.ThingWithAdditionalAttributesAdditionalSAttribute:
			queryInput.ExpressionAttributeNames["#ADDITIONALSATTRIBUTE"] = aws.String(string(db.ThingWithAdditionalAttributesAdditionalSAttribute))
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesAdditionalSAttribute), i)] = &dynamodb.AttributeValue{
					S: aws.String(attributeValue.(string)),
				}
			}
		case db.ThingWithAdditionalAttributesCreatedAt:
			queryInput.ExpressionAttributeNames["#CREATEDAT"] = aws.String(string(db.ThingWithAdditionalAttributesCreatedAt))
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesCreatedAt), i)] = &dynamodb.AttributeValue{
					S: aws.String(toDynamoTimeString(attributeValue.(strfmt.DateTime))),
				}
			}
		case db.ThingWithAdditionalAttributesHashNullable:
			queryInput.ExpressionAttributeNames["#HASHNULLABLE"] = aws.String(string(db.ThingWithAdditionalAttributesHashNullable))
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesHashNullable), i)] = &dynamodb.AttributeValue{
					S: aws.String(attributeValue.(string)),
				}
			}
		case db.ThingWithAdditionalAttributesID:
			queryInput.ExpressionAttributeNames["#ID"] = aws.String(string(db.ThingWithAdditionalAttributesID))
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesID), i)] = &dynamodb.AttributeValue{
					S: aws.String(attributeValue.(string)),
				}
			}
		case db.ThingWithAdditionalAttributesRangeNullable:
			queryInput.ExpressionAttributeNames["#RANGENULLABLE"] = aws.String(string(db.ThingWithAdditionalAttributesRangeNullable))
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithAdditionalAttributesRangeNullable), i)] = &dynamodb.AttributeValue{
					S: aws.String(toDynamoTimeString(attributeValue.(strfmt.DateTime))),
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
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#NAME": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": &dynamodb.AttributeValue{
				S: aws.String(input.Name),
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
		queryInput.ExpressionAttributeNames["#VERSION"] = aws.String("version")
		queryInput.ExpressionAttributeValues[":version"] = &dynamodb.AttributeValue{
			N: aws.String(fmt.Sprintf("%d", *input.VersionStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #VERSION <= :version")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #VERSION >= :version")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"version": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%d", input.StartingAfter.Version)),
			},
			"name": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Name),
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
		// Only assume an empty page means no more results if there are no filters applied
		if (len(input.FilterValues) == 0 || input.FilterExpression == "") && len(queryOutput.Items) == 0 {
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

	err := t.DynamoDBAPI.QueryPagesWithContext(ctx, queryInput, pageFn)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}

func (t ThingWithAdditionalAttributesTable) deleteThingWithAdditionalAttributes(ctx context.Context, name string, version int64) error {
	key, err := dynamodbattribute.MarshalMap(ddbThingWithAdditionalAttributesPrimaryKey{
		Name:    name,
		Version: version,
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

func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributesByID(ctx context.Context, id string) (*models.ThingWithAdditionalAttributes, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("thingID"),
		ExpressionAttributeNames: map[string]*string{
			"#ID": aws.String("id"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": &dynamodb.AttributeValue{
				S: aws.String(id),
			},
		},
		KeyConditionExpression: aws.String("#ID = :id"),
	}

	queryOutput, err := t.DynamoDBAPI.QueryWithContext(ctx, queryInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return nil, fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return nil, err
	}
	if len(queryOutput.Items) == 0 {
		return nil, db.ErrThingWithAdditionalAttributesByIDNotFound{ID: id}
	}

	var thingWithAdditionalAttributes models.ThingWithAdditionalAttributes
	if err := decodeThingWithAdditionalAttributes(queryOutput.Items[0], &thingWithAdditionalAttributes); err != nil {
		return nil, err
	}
	return &thingWithAdditionalAttributes, nil
}
func (t ThingWithAdditionalAttributesTable) scanThingWithAdditionalAttributessByID(ctx context.Context, input db.ScanThingWithAdditionalAttributessByIDInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("thingID"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"name":    exclusiveStartKey["name"],
			"version": exclusiveStartKey["version"],
			"id":      exclusiveStartKey["id"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithAdditionalAttributess(out.Items)
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
func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributessByNameAndCreatedAt(ctx context.Context, input db.GetThingWithAdditionalAttributessByNameAndCreatedAtInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	if input.CreatedAtStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.CreatedAtStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("name-createdAt"),
		ExpressionAttributeNames: map[string]*string{
			"#NAME": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": &dynamodb.AttributeValue{
				S: aws.String(input.Name),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.CreatedAtStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		queryInput.ExpressionAttributeNames["#CREATEDAT"] = aws.String("createdAt")
		queryInput.ExpressionAttributeValues[":createdAt"] = &dynamodb.AttributeValue{
			S: aws.String(toDynamoTimeString(*input.CreatedAtStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #CREATEDAT <= :createdAt")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #CREATEDAT >= :createdAt")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"createdAt": &dynamodb.AttributeValue{
				S: aws.String(toDynamoTimeString(input.StartingAfter.CreatedAt)),
			},
			"name": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Name),
			},
			"version": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%d", input.StartingAfter.Version)),
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

	err := t.DynamoDBAPI.QueryPagesWithContext(ctx, queryInput, pageFn)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}
func (t ThingWithAdditionalAttributesTable) scanThingWithAdditionalAttributessByNameAndCreatedAt(ctx context.Context, input db.ScanThingWithAdditionalAttributessByNameAndCreatedAtInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("name-createdAt"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"name":      exclusiveStartKey["name"],
			"version":   exclusiveStartKey["version"],
			"createdAt": exclusiveStartKey["createdAt"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithAdditionalAttributess(out.Items)
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
func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributessByNameAndRangeNullable(ctx context.Context, input db.GetThingWithAdditionalAttributessByNameAndRangeNullableInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	if input.RangeNullableStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.RangeNullableStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("name-rangeNullable"),
		ExpressionAttributeNames: map[string]*string{
			"#NAME": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": &dynamodb.AttributeValue{
				S: aws.String(input.Name),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.RangeNullableStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME = :name")
	} else {
		queryInput.ExpressionAttributeNames["#RANGENULLABLE"] = aws.String("rangeNullable")
		queryInput.ExpressionAttributeValues[":rangeNullable"] = &dynamodb.AttributeValue{
			S: aws.String(toDynamoTimeString(*input.RangeNullableStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #RANGENULLABLE <= :rangeNullable")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME = :name AND #RANGENULLABLE >= :rangeNullable")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"rangeNullable": &dynamodb.AttributeValue{
				S: aws.String(toDynamoTimeStringPtr(input.StartingAfter.RangeNullable)),
			},
			"name": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Name),
			},
			"version": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%d", input.StartingAfter.Version)),
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

	err := t.DynamoDBAPI.QueryPagesWithContext(ctx, queryInput, pageFn)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}
func (t ThingWithAdditionalAttributesTable) scanThingWithAdditionalAttributessByNameAndRangeNullable(ctx context.Context, input db.ScanThingWithAdditionalAttributessByNameAndRangeNullableInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("name-rangeNullable"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"name":          exclusiveStartKey["name"],
			"version":       exclusiveStartKey["version"],
			"rangeNullable": exclusiveStartKey["rangeNullable"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeThingWithAdditionalAttributess(out.Items)
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
func (t ThingWithAdditionalAttributesTable) getThingWithAdditionalAttributessByHashNullableAndName(ctx context.Context, input db.GetThingWithAdditionalAttributessByHashNullableAndNameInput, fn func(m *models.ThingWithAdditionalAttributes, lastThingWithAdditionalAttributes bool) bool) error {
	if input.NameStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.NameStartingAt or input.StartingAfter")
	}
	if input.HashNullable == "" {
		return fmt.Errorf("Hash key input.HashNullable cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("name-hashNullable"),
		ExpressionAttributeNames: map[string]*string{
			"#HASHNULLABLE": aws.String("hashNullable"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":hashNullable": &dynamodb.AttributeValue{
				S: aws.String(input.HashNullable),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.NameStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#HASHNULLABLE = :hashNullable")
	} else {
		queryInput.ExpressionAttributeNames["#NAME"] = aws.String("name")
		queryInput.ExpressionAttributeValues[":name"] = &dynamodb.AttributeValue{
			S: aws.String(*input.NameStartingAt),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#HASHNULLABLE = :hashNullable AND #NAME <= :name")
		} else {
			queryInput.KeyConditionExpression = aws.String("#HASHNULLABLE = :hashNullable AND #NAME >= :name")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"name": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Name),
			},
			"hashNullable": &dynamodb.AttributeValue{
				S: aws.String(*input.StartingAfter.HashNullable),
			},
			"version": &dynamodb.AttributeValue{
				N: aws.String(fmt.Sprintf("%d", input.StartingAfter.Version)),
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

	err := t.DynamoDBAPI.QueryPagesWithContext(ctx, queryInput, pageFn)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("table or index not found: %s", t.name())
			}
		}
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}

// encodeThingWithAdditionalAttributes encodes a ThingWithAdditionalAttributes as a DynamoDB map of attribute values.
func encodeThingWithAdditionalAttributes(m models.ThingWithAdditionalAttributes) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(ddbThingWithAdditionalAttributes{
		ThingWithAdditionalAttributes: m,
	})
}

// decodeThingWithAdditionalAttributes translates a ThingWithAdditionalAttributes stored in DynamoDB to a ThingWithAdditionalAttributes struct.
func decodeThingWithAdditionalAttributes(m map[string]*dynamodb.AttributeValue, out *models.ThingWithAdditionalAttributes) error {
	var ddbThingWithAdditionalAttributes ddbThingWithAdditionalAttributes
	if err := dynamodbattribute.UnmarshalMap(m, &ddbThingWithAdditionalAttributes); err != nil {
		return err
	}
	*out = ddbThingWithAdditionalAttributes.ThingWithAdditionalAttributes
	return nil
}

// decodeThingWithAdditionalAttributess translates a list of ThingWithAdditionalAttributess stored in DynamoDB to a slice of ThingWithAdditionalAttributes structs.
func decodeThingWithAdditionalAttributess(ms []map[string]*dynamodb.AttributeValue) ([]models.ThingWithAdditionalAttributes, error) {
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
