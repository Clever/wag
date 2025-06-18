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

// NoRangeThingWithCompositeAttributesTable represents the user-configurable properties of the NoRangeThingWithCompositeAttributes table.
type NoRangeThingWithCompositeAttributesTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbNoRangeThingWithCompositeAttributesPrimaryKey represents the primary key of a NoRangeThingWithCompositeAttributes in DynamoDB.
type ddbNoRangeThingWithCompositeAttributesPrimaryKey struct {
	NameBranch string `dynamodbav:"name_branch"`
}

// ddbNoRangeThingWithCompositeAttributesGSINameVersion represents the nameVersion GSI.
type ddbNoRangeThingWithCompositeAttributesGSINameVersion struct {
	NameVersion string          `dynamodbav:"name_version"`
	Date        strfmt.DateTime `dynamodbav:"date"`
}

// ddbNoRangeThingWithCompositeAttributesGSINameBranchCommit represents the nameBranchCommit GSI.
type ddbNoRangeThingWithCompositeAttributesGSINameBranchCommit struct {
	NameBranchCommit string `dynamodbav:"name_branch_commit"`
}

// ddbNoRangeThingWithCompositeAttributes represents a NoRangeThingWithCompositeAttributes as stored in DynamoDB.
type ddbNoRangeThingWithCompositeAttributes struct {
	models.NoRangeThingWithCompositeAttributes
}

func (t NoRangeThingWithCompositeAttributesTable) create(ctx context.Context) error {
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
				AttributeName: aws.String("name_branch_commit"),
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
			{
				IndexName: aws.String("nameBranchCommit"),
				Projection: &types.Projection{
					ProjectionType: types.ProjectionType("ALL"),
				},
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("name_branch_commit"),
						KeyType:       types.KeyTypeHash,
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

func (t NoRangeThingWithCompositeAttributesTable) saveNoRangeThingWithCompositeAttributes(ctx context.Context, m models.NoRangeThingWithCompositeAttributes) error {
	data, err := encodeNoRangeThingWithCompositeAttributes(m)
	if err != nil {
		return err
	}

	_, err = t.DynamoDBAPI.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName),
		Item:      data,
		ExpressionAttributeNames: map[string]string{
			"#NAME_BRANCH": "name_branch",
		},
		ConditionExpression: aws.String(
			"" +
				"" +
				"attribute_not_exists(#NAME_BRANCH)" +
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
			return db.ErrNoRangeThingWithCompositeAttributesAlreadyExists{
				NameBranch: fmt.Sprintf("%s@%s", *m.Name, *m.Branch),
			}
		}
		return err
	}
	return nil
}

func (t NoRangeThingWithCompositeAttributesTable) getNoRangeThingWithCompositeAttributes(ctx context.Context, name string, branch string) (*models.NoRangeThingWithCompositeAttributes, error) {
	key, err := attributevalue.MarshalMap(ddbNoRangeThingWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branch),
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
		return nil, db.ErrNoRangeThingWithCompositeAttributesNotFound{
			Name:   name,
			Branch: branch,
		}
	}

	var m models.NoRangeThingWithCompositeAttributes
	if err := decodeNoRangeThingWithCompositeAttributes(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t NoRangeThingWithCompositeAttributesTable) scanNoRangeThingWithCompositeAttributess(ctx context.Context, input db.ScanNoRangeThingWithCompositeAttributessInput, fn func(m *models.NoRangeThingWithCompositeAttributes, lastNoRangeThingWithCompositeAttributes bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.TableName),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		scanInput.Limit = aws.Int32(int32(*input.Limit))
	}
	if input.StartingAfter != nil {
		// must provide only the fields constituting the index
		scanInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"name_branch": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s@%s", *input.StartingAfter.Name, *input.StartingAfter.Branch),
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

		items, err := decodeNoRangeThingWithCompositeAttributess(out.Items)
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

func (t NoRangeThingWithCompositeAttributesTable) deleteNoRangeThingWithCompositeAttributes(ctx context.Context, name string, branch string) error {

	key, err := attributevalue.MarshalMap(ddbNoRangeThingWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branch),
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

func (t NoRangeThingWithCompositeAttributesTable) getNoRangeThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input db.GetNoRangeThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.NoRangeThingWithCompositeAttributes, lastNoRangeThingWithCompositeAttributes bool) bool) error {
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
		items, err := decodeNoRangeThingWithCompositeAttributess(queryOutput.Items)
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
func (t NoRangeThingWithCompositeAttributesTable) scanNoRangeThingWithCompositeAttributessByNameVersionAndDate(ctx context.Context, input db.ScanNoRangeThingWithCompositeAttributessByNameVersionAndDateInput, fn func(m *models.NoRangeThingWithCompositeAttributes, lastNoRangeThingWithCompositeAttributes bool) bool) error {
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
			"name_version": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s:%d", *input.StartingAfter.Name, input.StartingAfter.Version),
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

		items, err := decodeNoRangeThingWithCompositeAttributess(out.Items)
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
func (t NoRangeThingWithCompositeAttributesTable) getNoRangeThingWithCompositeAttributesByNameBranchCommit(ctx context.Context, name string, branch string, commit string) (*models.NoRangeThingWithCompositeAttributes, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		IndexName: aws.String("nameBranchCommit"),
		ExpressionAttributeNames: map[string]string{
			"#NAME_BRANCH_COMMIT": "name_branch_commit",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nameBranchCommit": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s--%s--%s", name, branch, commit),
			},
		},
		KeyConditionExpression: aws.String("#NAME_BRANCH_COMMIT = :nameBranchCommit"),
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
		return nil, db.ErrNoRangeThingWithCompositeAttributesByNameBranchCommitNotFound{
			Name:   name,
			Branch: branch,
			Commit: commit,
		}
	}

	var noRangeThingWithCompositeAttributes models.NoRangeThingWithCompositeAttributes
	if err := decodeNoRangeThingWithCompositeAttributes(queryOutput.Items[0], &noRangeThingWithCompositeAttributes); err != nil {
		return nil, err
	}
	return &noRangeThingWithCompositeAttributes, nil
}

// encodeNoRangeThingWithCompositeAttributes encodes a NoRangeThingWithCompositeAttributes as a DynamoDB map of attribute values.
func encodeNoRangeThingWithCompositeAttributes(m models.NoRangeThingWithCompositeAttributes) (map[string]types.AttributeValue, error) {
	// with composite attributes, marshal the model
	val, err := attributevalue.MarshalMapWithOptions(m, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(*m.Branch, "--") {
		return nil, fmt.Errorf("branch cannot contain '--': %s", *m.Branch)
	}
	if strings.Contains(*m.Commit, "--") {
		return nil, fmt.Errorf("commit cannot contain '--': %s", *m.Commit)
	}
	if strings.Contains(*m.Name, "--") {
		return nil, fmt.Errorf("name cannot contain '--': %s", *m.Name)
	}
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
	primaryKey, err := attributevalue.MarshalMap(ddbNoRangeThingWithCompositeAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", *m.Name, *m.Branch),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range primaryKey {
		val[k] = v
	}
	nameVersion, err := attributevalue.MarshalMap(ddbNoRangeThingWithCompositeAttributesGSINameVersion{
		NameVersion: fmt.Sprintf("%s:%d", *m.Name, m.Version),
		Date:        *m.Date,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range nameVersion {
		val[k] = v
	}
	nameBranchCommit, err := attributevalue.MarshalMap(ddbNoRangeThingWithCompositeAttributesGSINameBranchCommit{
		NameBranchCommit: fmt.Sprintf("%s--%s--%s", *m.Name, *m.Branch, *m.Commit),
	})
	if err != nil {
		return nil, err
	}
	for k, v := range nameBranchCommit {
		val[k] = v
	}
	return val, err
}

// decodeNoRangeThingWithCompositeAttributes translates a NoRangeThingWithCompositeAttributes stored in DynamoDB to a NoRangeThingWithCompositeAttributes struct.
func decodeNoRangeThingWithCompositeAttributes(m map[string]types.AttributeValue, out *models.NoRangeThingWithCompositeAttributes) error {
	var ddbNoRangeThingWithCompositeAttributes ddbNoRangeThingWithCompositeAttributes
	if err := attributevalue.UnmarshalMap(m, &ddbNoRangeThingWithCompositeAttributes); err != nil {
		return err
	}
	*out = ddbNoRangeThingWithCompositeAttributes.NoRangeThingWithCompositeAttributes
	return nil
}

// decodeNoRangeThingWithCompositeAttributess translates a list of NoRangeThingWithCompositeAttributess stored in DynamoDB to a slice of NoRangeThingWithCompositeAttributes structs.
func decodeNoRangeThingWithCompositeAttributess(ms []map[string]types.AttributeValue) ([]models.NoRangeThingWithCompositeAttributes, error) {
	noRangeThingWithCompositeAttributess := make([]models.NoRangeThingWithCompositeAttributes, len(ms))
	for i, m := range ms {
		var noRangeThingWithCompositeAttributes models.NoRangeThingWithCompositeAttributes
		if err := decodeNoRangeThingWithCompositeAttributes(m, &noRangeThingWithCompositeAttributes); err != nil {
			return nil, err
		}
		noRangeThingWithCompositeAttributess[i] = noRangeThingWithCompositeAttributes
	}
	return noRangeThingWithCompositeAttributess, nil
}
