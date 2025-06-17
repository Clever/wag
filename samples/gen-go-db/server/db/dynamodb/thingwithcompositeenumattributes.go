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

// ThingWithCompositeEnumAttributesTable represents the user-configurable properties of the ThingWithCompositeEnumAttributes table.
type ThingWithCompositeEnumAttributesTable struct {
	DynamoDBAPI        *dynamodb.Client
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbThingWithCompositeEnumAttributesPrimaryKey represents the primary key of a ThingWithCompositeEnumAttributes in DynamoDB.
type ddbThingWithCompositeEnumAttributesPrimaryKey struct {
	NameBranch string          `dynamodbav:"name_branch"`
	Date       strfmt.DateTime `dynamodbav:"date"`
}

// ddbThingWithCompositeEnumAttributes represents a ThingWithCompositeEnumAttributes as stored in DynamoDB.
type ddbThingWithCompositeEnumAttributes struct {
	models.ThingWithCompositeEnumAttributes
}

func (t ThingWithCompositeEnumAttributesTable) create(ctx context.Context) error {
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

func (t ThingWithCompositeEnumAttributesTable) saveThingWithCompositeEnumAttributes(ctx context.Context, m models.ThingWithCompositeEnumAttributes) error {
	data, err := encodeThingWithCompositeEnumAttributes(m)
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
			return db.ErrThingWithCompositeEnumAttributesAlreadyExists{
				NameBranch: fmt.Sprintf("%s@%s", *m.Name, m.BranchID),
				Date:       *m.Date,
			}
		}
		return err
	}
	return nil
}

func (t ThingWithCompositeEnumAttributesTable) getThingWithCompositeEnumAttributes(ctx context.Context, name string, branchID models.Branch, date strfmt.DateTime) (*models.ThingWithCompositeEnumAttributes, error) {
	key, err := attributevalue.MarshalMap(ddbThingWithCompositeEnumAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branchID),
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
		return nil, db.ErrThingWithCompositeEnumAttributesNotFound{
			Name:     name,
			BranchID: branchID,
			Date:     date,
		}
	}

	var m models.ThingWithCompositeEnumAttributes
	if err := decodeThingWithCompositeEnumAttributes(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t ThingWithCompositeEnumAttributesTable) scanThingWithCompositeEnumAttributess(ctx context.Context, input db.ScanThingWithCompositeEnumAttributessInput, fn func(m *models.ThingWithCompositeEnumAttributes, lastThingWithCompositeEnumAttributes bool) bool) error {
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
				Value: fmt.Sprintf("%s@%s", *input.StartingAfter.Name, input.StartingAfter.BranchID),
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

		items, err := decodeThingWithCompositeEnumAttributess(out.Items)
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

func (t ThingWithCompositeEnumAttributesTable) getThingWithCompositeEnumAttributessByNameBranchAndDate(ctx context.Context, input db.GetThingWithCompositeEnumAttributessByNameBranchAndDateInput, fn func(m *models.ThingWithCompositeEnumAttributes, lastThingWithCompositeEnumAttributes bool) bool) error {
	if input.DateStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateStartingAt or input.StartingAfter")
	}
	if input.Name == "" {
		return fmt.Errorf("Hash key input.Name cannot be empty")
	}
	if input.BranchID == "" {
		return fmt.Errorf("Hash key input.BranchID cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.TableName),
		ExpressionAttributeNames: map[string]string{
			"#NAME_BRANCH": "name_branch",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":nameBranch": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%s@%s", input.Name, input.BranchID),
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
				Value: fmt.Sprintf("%s@%s", *input.StartingAfter.Name, input.StartingAfter.BranchID),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeThingWithCompositeEnumAttributess(queryOutput.Items)
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

func (t ThingWithCompositeEnumAttributesTable) deleteThingWithCompositeEnumAttributes(ctx context.Context, name string, branchID models.Branch, date strfmt.DateTime) error {

	key, err := attributevalue.MarshalMap(ddbThingWithCompositeEnumAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", name, branchID),
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

// encodeThingWithCompositeEnumAttributes encodes a ThingWithCompositeEnumAttributes as a DynamoDB map of attribute values.
func encodeThingWithCompositeEnumAttributes(m models.ThingWithCompositeEnumAttributes) (map[string]types.AttributeValue, error) {
	val, err := attributevalue.MarshalMap(ddbThingWithCompositeEnumAttributes{
		ThingWithCompositeEnumAttributes: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(*m.Name, "@") {
		return nil, fmt.Errorf("name cannot contain '@': %s", *m.Name)
	}
	// add in composite attributes
	primaryKey, err := attributevalue.MarshalMap(ddbThingWithCompositeEnumAttributesPrimaryKey{
		NameBranch: fmt.Sprintf("%s@%s", *m.Name, m.BranchID),
		Date:       *m.Date,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range primaryKey {
		val[k] = v
	}

	// Ensure all attribute names match DynamoDB expectations
	if v, ok := val["NameBranch"]; ok {
		// Convert to the correct attribute value type
		switch av := v.(type) {
		case *types.AttributeValueMemberS:
			val["name_branch"] = &types.AttributeValueMemberS{Value: av.Value}
		case *types.AttributeValueMemberN:
			val["name_branch"] = &types.AttributeValueMemberN{Value: av.Value}
		case *types.AttributeValueMemberB:
			val["name_branch"] = &types.AttributeValueMemberB{Value: av.Value}
		case *types.AttributeValueMemberBOOL:
			val["name_branch"] = &types.AttributeValueMemberBOOL{Value: av.Value}
		case *types.AttributeValueMemberNULL:
			val["name_branch"] = &types.AttributeValueMemberNULL{Value: av.Value}
		case *types.AttributeValueMemberM:
			val["name_branch"] = &types.AttributeValueMemberM{Value: av.Value}
		case *types.AttributeValueMemberL:
			val["name_branch"] = &types.AttributeValueMemberL{Value: av.Value}
		case *types.AttributeValueMemberSS:
			val["name_branch"] = &types.AttributeValueMemberSS{Value: av.Value}
		case *types.AttributeValueMemberNS:
			val["name_branch"] = &types.AttributeValueMemberNS{Value: av.Value}
		case *types.AttributeValueMemberBS:
			val["name_branch"] = &types.AttributeValueMemberBS{Value: av.Value}
		default:
			val["name_branch"] = v
		}
		delete(val, "NameBranch")
	}
	if v, ok := val["Date"]; ok {
		// Convert to the correct attribute value type
		switch av := v.(type) {
		case *types.AttributeValueMemberS:
			val["date"] = &types.AttributeValueMemberS{Value: av.Value}
		case *types.AttributeValueMemberN:
			val["date"] = &types.AttributeValueMemberN{Value: av.Value}
		case *types.AttributeValueMemberB:
			val["date"] = &types.AttributeValueMemberB{Value: av.Value}
		case *types.AttributeValueMemberBOOL:
			val["date"] = &types.AttributeValueMemberBOOL{Value: av.Value}
		case *types.AttributeValueMemberNULL:
			val["date"] = &types.AttributeValueMemberNULL{Value: av.Value}
		case *types.AttributeValueMemberM:
			val["date"] = &types.AttributeValueMemberM{Value: av.Value}
		case *types.AttributeValueMemberL:
			val["date"] = &types.AttributeValueMemberL{Value: av.Value}
		case *types.AttributeValueMemberSS:
			val["date"] = &types.AttributeValueMemberSS{Value: av.Value}
		case *types.AttributeValueMemberNS:
			val["date"] = &types.AttributeValueMemberNS{Value: av.Value}
		case *types.AttributeValueMemberBS:
			val["date"] = &types.AttributeValueMemberBS{Value: av.Value}
		default:
			val["date"] = v
		}
		delete(val, "Date")
	}

	return val, err
}

// decodeThingWithCompositeEnumAttributes translates a ThingWithCompositeEnumAttributes stored in DynamoDB to a ThingWithCompositeEnumAttributes struct.
func decodeThingWithCompositeEnumAttributes(m map[string]types.AttributeValue, out *models.ThingWithCompositeEnumAttributes) error {
	var ddbThingWithCompositeEnumAttributes ddbThingWithCompositeEnumAttributes
	if err := attributevalue.UnmarshalMap(m, &ddbThingWithCompositeEnumAttributes); err != nil {
		return err
	}
	*out = ddbThingWithCompositeEnumAttributes.ThingWithCompositeEnumAttributes
	return nil
}

// decodeThingWithCompositeEnumAttributess translates a list of ThingWithCompositeEnumAttributess stored in DynamoDB to a slice of ThingWithCompositeEnumAttributes structs.
func decodeThingWithCompositeEnumAttributess(ms []map[string]types.AttributeValue) ([]models.ThingWithCompositeEnumAttributes, error) {
	thingWithCompositeEnumAttributess := make([]models.ThingWithCompositeEnumAttributes, len(ms))
	for i, m := range ms {
		var thingWithCompositeEnumAttributes models.ThingWithCompositeEnumAttributes
		if err := decodeThingWithCompositeEnumAttributes(m, &thingWithCompositeEnumAttributes); err != nil {
			return nil, err
		}
		thingWithCompositeEnumAttributess[i] = thingWithCompositeEnumAttributes
	}
	return thingWithCompositeEnumAttributess, nil
}
