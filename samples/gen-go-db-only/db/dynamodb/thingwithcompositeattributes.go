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

// ThingWithCompositeAttributesTable represents the user-configurable properties of the ThingWithCompositeAttributes table.
type ThingWithCompositeAttributesTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithCompositeAttributesPrimaryKey represents the primary key of a ThingWithCompositeAttributes in DynamoDB.
type ddbThingWithCompositeAttributesPrimaryKey struct {
	NameBranch string          `dynamodbav:"name_branch"`
	Date       strfmt.DateTime `dynamodbav:"date"`
}

// ddbThingWithCompositeAttributesGSINameVersion represents the nameVersion GSI.
type ddbThingWithCompositeAttributesGSINameVersion struct {
	NameVersion string          `dynamodbav:"name_version"`
	Date        strfmt.DateTime `dynamodbav:"date"`
}

// ddbThingWithCompositeAttributes represents a ThingWithCompositeAttributes as stored in DynamoDB.
type ddbThingWithCompositeAttributes struct {
	models.ThingWithCompositeAttributes `dynamodbav:",inline"`
}

func (t ThingWithCompositeAttributesTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("date"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("name_branch"),
				AttributeType: types.ScalarAttributeType("S"),
			},
			{
				AttributeName: aws.String("name_version"),
				AttributeType: types.ScalarAttributeType("S"),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name_branch"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("date"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("nameVersion"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("name_version"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("date"),
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

func (t ThingWithCompositeAttributesTable) saveThingWithCompositeAttributes(ctx context.Context, m models.ThingWithCompositeAttributes) error {
	data, err := encodeThingWithCompositeAttributes(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#NAME_BRANCH": "name_branch",
			"#DATE":        "date",
		},
		ConditionExpression: aws.String(
			"" +
				"" +
				"attribute_not_exists(#NAME_BRANCH)" +
				"" +
				" AND " +
				"attribute_not_exists(#DATE)" +
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
			return db.ErrThingWithCompositeAttributesAlreadyExists{
				NameBranch: fmt.Sprintf("%s@%s", *m.Name, *m.Branch),
				Date:       *m.Date,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithCompositeAttributesTable) getThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) (*models.ThingWithCompositeAttributes, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branch),
		Date:       date,
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
		return nil, db.ErrThingWithCompositeAttributesNotFound{
			Name:   name,
			Branch: branch,
			Date:   date,
		}
	}

	var m models.ThingWithCompositeAttributes
	if err := decodeThingWithCompositeAttributes(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithCompositeAttributesTable) scanThingWithCompositeAttributess(ctx context.Context, input db.ScanThingWithCompositeAttributessInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error {
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
			"name_branch": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s@%s", *input.StartingAfter.Name, *input.StartingAfter.Branch),
			},
			"date": exclusiveStartKey["date"],
		}
	}
	totalRecordsProcessed := int64(0)

	paginator := dynamodb.NewScanPaginator(t.DynamoDBAPI, scanInput)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error getting next page: %s", err.Error())
		}

		items, err := decodeThingWithCompositeAttributess(out.Items)
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

func (t ThingWithCompositeAttributesTable) getThingWithCompositeAttributessByNameBranchAndDateParseFilters(queryInput *dynamodb.QueryInput, input db.GetThingWithCompositeAttributessByNameBranchAndDateInput) {
	for _, filterValue := range input.FilterValues {
		switch filterValue.AttributeName {
		case db.ThingWithCompositeAttributesVersion:
			queryInput.ExpressionAttributeNames["#VERSION"] = string(db.ThingWithCompositeAttributesVersion)
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.ThingWithCompositeAttributesVersion), i)] = &types.AttributeValueMemberN{
					Value: fmt.Sprint(attributeValue.(int64)),
				}
			}
		}
	}
}

func (t ThingWithCompositeAttributesTable) getThingWithCompositeAttributessByNameBranchAndDate(ctx context.Context, input db.GetThingWithCompositeAttributessByNameBranchAndDateInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error {
	if input.DateStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	if input.Branch == "" {
		return fmt.Errorf("Hash key input.Branch cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		ExpressionAttributeNames: map[string]string{
			"#NAME_BRANCH": "name_branch",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nameBranch": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s@%s", input.Name, input.Branch),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME_BRANCH = :nameBranch")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = "date"
		queryInput.ExpressionAttributeValues[":date"] = &types.AttributeValueMemberS{
			Value: datetimeToDynamoTimeString(*input.DateStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME_BRANCH = :nameBranch AND #DATE <= :date")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME_BRANCH = :nameBranch AND #DATE >= :date")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"date": &types.AttributeValueMemberS{
				Value: datetimePtrToDynamoTimeString(input.StartingAfter.Date),
			},

			"name_branch": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s@%s", *input.StartingAfter.Name, *input.StartingAfter.Branch),
			},
		}
	}
	if len(input.FilterValues) > 0 && input.FilterExpression != "" {
		t.getThingWithCompositeAttributessByNameBranchAndDateParseFilters(queryInput, input)
		queryInput.FilterExpression = aws.String(input.FilterExpression)
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithCompositeAttributess(queryOutput.Items)
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

func (t ThingWithCompositeAttributesTable) deleteThingWithCompositeAttributes(ctx context.Context, name string, branch string, date strfmt.DateTime) error {

	key, err := attributevalue.MarshalMap(ddbThingWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branch),
		Date:       date,
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

func (t ThingWithCompositeAttributesTable) getThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input db.GetThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error {
	if input.DateStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("nameVersion"),
		ExpressionAttributeNames: map[string]string{
			"#NAME_VERSION": "name_version",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nameVersion": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s:%d", input.Name, input.Version),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#NAME_VERSION = :nameVersion")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = "date"
		queryInput.ExpressionAttributeValues[":date"] = &types.AttributeValueMemberS{
			Value: datetimeToDynamoTimeString(*input.DateStartingAt),
		}

		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#NAME_VERSION = :nameVersion AND #DATE <= :date")
		} else {
			queryInput.KeyConditionExpression = aws.String("#NAME_VERSION = :nameVersion AND #DATE >= :date")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"date": &types.AttributeValueMemberS{
				Value: datetimePtrToDynamoTimeString(input.StartingAfter.Date),
			},
			"name_version": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s:%d", *input.StartingAfter.Name, input.StartingAfter.Version),
			},
			"name_branch": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s@%s", *input.StartingAfter.Name, *input.StartingAfter.Branch),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithCompositeAttributess(queryOutput.Items)
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
func (t ThingWithCompositeAttributesTable) scanThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input db.ScanThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.ThingWithCompositeAttributes, lastThingWithCompositeAttributes bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	scanInput.IndexName = aws.String("nameVersion")
	if input.StartingAfter != nil {
		exclusiveStartKey, err := attributevalue.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name_branch": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s@%s", *input.StartingAfter.Name, *input.StartingAfter.Branch),
			},
			"date": exclusiveStartKey["date"],
			"name_version": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s:%d", *input.StartingAfter.Name, input.StartingAfter.Version),
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

		items, err := decodeThingWithCompositeAttributess(out.Items)
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

// encodeThingWithCompositeAttributes encodes a ThingWithCompositeAttributes as a DynamoDB map of attribute values.
func encodeThingWithCompositeAttributes(m models.ThingWithCompositeAttributes) (map[string]types.AttributeValue, error) {
	val, err := attributevalue.MarshalMap(ddbThingWithCompositeAttributes{
		ThingWithCompositeAttributes: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(*m.Name, ":") {
		return nil, fmt.Errorf("name cannot contain ':': %s", *m.Name)
	}
	if strings.Contains(*m.Branch, "@") {
		return nil, fmt.Errorf("branch cannot contain '@': %s", *m.Branch)
	}
	if strings.Contains(*m.Name, "@") {
		return nil, fmt.Errorf("name cannot contain '@': %s", *m.Name)
	}
	// add in composite attributes
	primaryKey, err := attributevalue.MarshalMap(ddbThingWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", *m.Name, *m.Branch),
		Date:       *m.Date,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range primaryKey {
		val[k] = v
	}
	nameVersion, err := attributevalue.MarshalMap(ddbThingWithCompositeAttributesGSINameVersion{
		NameVersion: fmt.Sprintf("%s:%d", *m.Name, m.Version),
		Date:        *m.Date,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range nameVersion {
		val[k] = v
	}
	return val, err
}

// decodeThingWithCompositeAttributes translates a ThingWithCompositeAttributes stored in DynamoDB to a ThingWithCompositeAttributes struct.
func decodeThingWithCompositeAttributes(m map[string]types.AttributeValue, out *models.ThingWithCompositeAttributes) error {
	var ddbThingWithCompositeAttributes ddbThingWithCompositeAttributes
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithCompositeAttributes); err != nil {
		return err
	}
	*out = ddbThingWithCompositeAttributes.ThingWithCompositeAttributes
	return nil
}

// decodeThingWithCompositeAttributess translates a list of ThingWithCompositeAttributess stored in DynamoDB to a slice of ThingWithCompositeAttributes structs.
func decodeThingWithCompositeAttributess(ms []map[string]types.AttributeValue) ([]models.ThingWithCompositeAttributes, error) {
	thingWithCompositeAttributess := make([]models.ThingWithCompositeAttributes, len(ms))
	for i, m := range ms {
		var thingWithCompositeAttributes models.ThingWithCompositeAttributes
		if err := decodeThingWithCompositeAttributes(m, &thingWithCompositeAttributes); err != nil {
			return nil, err
		}
		thingWithCompositeAttributess[i] = thingWithCompositeAttributes
	}
	return thingWithCompositeAttributess, nil
}
