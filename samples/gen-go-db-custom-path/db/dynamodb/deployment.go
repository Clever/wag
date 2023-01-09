package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Clever/wag/samples/gen-go-db-custom-path/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-db-custom-path/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/go-openapi/strfmt"
)

var _ = strfmt.DateTime{}

// DeploymentTable represents the user-configurable properties of the Deployment table.
type DeploymentTable struct {
	DynamoDBAPI        dynamodbiface.DynamoDBAPI
	Prefix             string
	TableName          string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// ddbDeploymentPrimaryKey represents the primary key of a Deployment in DynamoDB.
type ddbDeploymentPrimaryKey struct {
	EnvApp  string `dynamodbav:"envApp"`
	Version string `dynamodbav:"version"`
}

// ddbDeploymentGSIByDate represents the byDate GSI.
type ddbDeploymentGSIByDate struct {
	EnvApp string          `dynamodbav:"envApp"`
	Date   strfmt.DateTime `dynamodbav:"date"`
}

// ddbDeploymentGSIByEnvironment represents the byEnvironment GSI.
type ddbDeploymentGSIByEnvironment struct {
	Environment string          `dynamodbav:"environment"`
	Date        strfmt.DateTime `dynamodbav:"date"`
}

// ddbDeploymentGSIByVersion represents the byVersion GSI.
type ddbDeploymentGSIByVersion struct {
	Version string `dynamodbav:"version"`
}

// ddbDeployment represents a Deployment as stored in DynamoDB.
type ddbDeployment struct {
	models.Deployment
}

func (t DeploymentTable) name() string {
	if t.TableName != "" {
		return t.TableName
	}
	return fmt.Sprintf("%s-deployments", t.Prefix)
}

func (t DeploymentTable) create(ctx context.Context) error {
	if _, err := t.DynamoDBAPI.CreateTableWithContext(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("date"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("envApp"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("environment"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("envApp"),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       aws.String(dynamodb.KeyTypeRange),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("byDate"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("envApp"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("date"),
						KeyType:       aws.String(dynamodb.KeyTypeRange),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("byEnvironment"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("environment"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
					},
					{
						AttributeName: aws.String("date"),
						KeyType:       aws.String(dynamodb.KeyTypeRange),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(t.ReadCapacityUnits),
					WriteCapacityUnits: aws.Int64(t.WriteCapacityUnits),
				},
			},
			{
				IndexName: aws.String("byVersion"),
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("version"),
						KeyType:       aws.String(dynamodb.KeyTypeHash),
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

func (t DeploymentTable) saveDeployment(ctx context.Context, m models.Deployment) error {
	data, err := encodeDeployment(m)
	if err != nil {
		return err
	}
	_, err = t.DynamoDBAPI.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.name()),
		Item:      data,
	})
	return err
}

func (t DeploymentTable) getDeployment(ctx context.Context, environment string, application string, version string) (*models.Deployment, error) {
	key, err := dynamodbattribute.MarshalMap(ddbDeploymentPrimaryKey{
		EnvApp:  fmt.Sprintf("%s--%s", environment, application),
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
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, db.ErrDeploymentNotFound{
			Environment: environment,
			Application: application,
			Version:     version,
		}
	}

	var m models.Deployment
	if err := decodeDeployment(res.Item, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func (t DeploymentTable) scanDeployments(ctx context.Context, input db.ScanDeploymentsInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
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
			"envApp": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s--%s", input.StartingAfter.Environment, input.StartingAfter.Application)),
			},
			"version": exclusiveStartKey["version"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeDeployments(out.Items)
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

func (t DeploymentTable) getDeploymentsByEnvAppAndVersionParseFilters(queryInput *dynamodb.QueryInput, input db.GetDeploymentsByEnvAppAndVersionInput) {
	for _, filterValue := range input.FilterValues {
		switch filterValue.AttributeName {
		case db.DeploymentDate:
			queryInput.ExpressionAttributeNames["#DATE"] = aws.String(string(db.DeploymentDate))
			for i, attributeValue := range filterValue.AttributeValues {
				queryInput.ExpressionAttributeValues[fmt.Sprintf(":%s_value%d", string(db.DeploymentDate), i)] = &dynamodb.AttributeValue{
					S: aws.String(toDynamoTimeString(attributeValue.(strfmt.DateTime))),
				}
			}
		}
	}
}

func (t DeploymentTable) getDeploymentsByEnvAppAndVersion(ctx context.Context, input db.GetDeploymentsByEnvAppAndVersionInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
	if input.VersionStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.VersionStartingAt or input.StartingAfter")
	}
	if input.Environment == "" {
		return fmt.Errorf("Hash key input.Environment cannot be empty")
	}
	if input.Application == "" {
		return fmt.Errorf("Hash key input.Application cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		ExpressionAttributeNames: map[string]*string{
			"#ENVAPP": aws.String("envApp"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":envApp": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s--%s", input.Environment, input.Application)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(!input.DisableConsistentRead),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.VersionStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#ENVAPP = :envApp")
	} else {
		queryInput.ExpressionAttributeNames["#VERSION"] = aws.String("version")
		queryInput.ExpressionAttributeValues[":version"] = &dynamodb.AttributeValue{
			S: aws.String(*input.VersionStartingAt),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#ENVAPP = :envApp AND #VERSION <= :version")
		} else {
			queryInput.KeyConditionExpression = aws.String("#ENVAPP = :envApp AND #VERSION >= :version")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"version": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Version),
			},
			"envApp": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s--%s", input.StartingAfter.Environment, input.StartingAfter.Application)),
			},
		}
	}
	if len(input.FilterValues) > 0 && input.FilterExpression != "" {
		t.getDeploymentsByEnvAppAndVersionParseFilters(queryInput, input)
		queryInput.FilterExpression = aws.String(input.FilterExpression)
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		// Only assume an empty page means no more results if there are no filters applied
		if (len(input.FilterValues) == 0 || input.FilterExpression == "") && len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeDeployments(queryOutput.Items)
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
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}

func (t DeploymentTable) deleteDeployment(ctx context.Context, environment string, application string, version string) error {
	key, err := dynamodbattribute.MarshalMap(ddbDeploymentPrimaryKey{
		EnvApp:  fmt.Sprintf("%s--%s", environment, application),
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
		return err
	}

	return nil
}

func (t DeploymentTable) getDeploymentsByEnvAppAndDate(ctx context.Context, input db.GetDeploymentsByEnvAppAndDateInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
	if input.DateStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateStartingAt or input.StartingAfter")
	}
	if input.Environment == "" {
		return fmt.Errorf("Hash key input.Environment cannot be empty")
	}
	if input.Application == "" {
		return fmt.Errorf("Hash key input.Application cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("byDate"),
		ExpressionAttributeNames: map[string]*string{
			"#ENVAPP": aws.String("envApp"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":envApp": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s--%s", input.Environment, input.Application)),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#ENVAPP = :envApp")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = aws.String("date")
		queryInput.ExpressionAttributeValues[":date"] = &dynamodb.AttributeValue{
			S: aws.String(toDynamoTimeString(*input.DateStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#ENVAPP = :envApp AND #DATE <= :date")
		} else {
			queryInput.KeyConditionExpression = aws.String("#ENVAPP = :envApp AND #DATE >= :date")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"date": &dynamodb.AttributeValue{
				S: aws.String(toDynamoTimeString(input.StartingAfter.Date)),
			},
			"envApp": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s--%s", input.StartingAfter.Environment, input.StartingAfter.Application)),
			},
			"version": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Version),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeDeployments(queryOutput.Items)
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
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}
func (t DeploymentTable) scanDeploymentsByEnvAppAndDate(ctx context.Context, input db.ScanDeploymentsByEnvAppAndDateInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("byDate"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"envApp": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s--%s", input.StartingAfter.Environment, input.StartingAfter.Application)),
			},
			"version": exclusiveStartKey["version"],
			"date":    exclusiveStartKey["date"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeDeployments(out.Items)
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
func (t DeploymentTable) getDeploymentsByEnvironmentAndDate(ctx context.Context, input db.GetDeploymentsByEnvironmentAndDateInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
	if input.DateStartingAt != nil && input.StartingAfter != nil {
		return fmt.Errorf("Can specify only one of input.DateStartingAt or input.StartingAfter")
	}
	if input.Environment == "" {
		return fmt.Errorf("Hash key input.Environment cannot be empty")
	}
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("byEnvironment"),
		ExpressionAttributeNames: map[string]*string{
			"#ENVIRONMENT": aws.String("environment"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":environment": &dynamodb.AttributeValue{
				S: aws.String(input.Environment),
			},
		},
		ScanIndexForward: aws.Bool(!input.Descending),
		ConsistentRead:   aws.Bool(false),
	}
	if input.Limit != nil {
		queryInput.Limit = input.Limit
	}
	if input.DateStartingAt == nil {
		queryInput.KeyConditionExpression = aws.String("#ENVIRONMENT = :environment")
	} else {
		queryInput.ExpressionAttributeNames["#DATE"] = aws.String("date")
		queryInput.ExpressionAttributeValues[":date"] = &dynamodb.AttributeValue{
			S: aws.String(toDynamoTimeString(*input.DateStartingAt)),
		}
		if input.Descending {
			queryInput.KeyConditionExpression = aws.String("#ENVIRONMENT = :environment AND #DATE <= :date")
		} else {
			queryInput.KeyConditionExpression = aws.String("#ENVIRONMENT = :environment AND #DATE >= :date")
		}
	}
	if input.StartingAfter != nil {
		queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"date": &dynamodb.AttributeValue{
				S: aws.String(toDynamoTimeString(input.StartingAfter.Date)),
			},
			"environment": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Environment),
			},
			"version": &dynamodb.AttributeValue{
				S: aws.String(input.StartingAfter.Version),
			},
			"envApp": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s--%s", input.StartingAfter.Environment, input.StartingAfter.Application)),
			},
		}
	}

	totalRecordsProcessed := int64(0)
	var pageFnErr error
	pageFn := func(queryOutput *dynamodb.QueryOutput, lastPage bool) bool {
		if len(queryOutput.Items) == 0 {
			return false
		}
		items, err := decodeDeployments(queryOutput.Items)
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
		return err
	}
	if pageFnErr != nil {
		return pageFnErr
	}

	return nil
}
func (t DeploymentTable) getDeploymentByVersion(ctx context.Context, version string) (*models.Deployment, error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(t.name()),
		IndexName: aws.String("byVersion"),
		ExpressionAttributeNames: map[string]*string{
			"#VERSION": aws.String("version"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":version": &dynamodb.AttributeValue{
				S: aws.String(version),
			},
		},
		KeyConditionExpression: aws.String("#VERSION = :version"),
	}

	queryOutput, err := t.DynamoDBAPI.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if len(queryOutput.Items) == 0 {
		return nil, db.ErrDeploymentByVersionNotFound{Version: version}
	}

	var deployment models.Deployment
	if err := decodeDeployment(queryOutput.Items[0], &deployment); err != nil {
		return nil, err
	}
	return &deployment, nil
}
func (t DeploymentTable) scanDeploymentsByVersion(ctx context.Context, input db.ScanDeploymentsByVersionInput, fn func(m *models.Deployment, lastDeployment bool) bool) error {
	scanInput := &dynamodb.ScanInput{
		TableName:      aws.String(t.name()),
		ConsistentRead: aws.Bool(!input.DisableConsistentRead),
		Limit:          input.Limit,
		IndexName:      aws.String("byVersion"),
	}
	if input.StartingAfter != nil {
		exclusiveStartKey, err := dynamodbattribute.MarshalMap(input.StartingAfter)
		if err != nil {
			return fmt.Errorf("error encoding exclusive start key for scan: %s", err.Error())
		}
		// must provide the fields constituting the index and the primary key
		// https://stackoverflow.com/questions/40988397/dynamodb-pagination-with-withexclusivestartkey-on-a-global-secondary-index
		scanInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"envApp": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%s--%s", input.StartingAfter.Environment, input.StartingAfter.Application)),
			},
			"version": exclusiveStartKey["version"],
		}
	}
	totalRecordsProcessed := int64(0)
	var innerErr error
	err := t.DynamoDBAPI.ScanPagesWithContext(ctx, scanInput, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		items, err := decodeDeployments(out.Items)
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

// WTH

// encodeDeployment encodes a Deployment as a DynamoDB map of attribute values.
func encodeDeployment(m models.Deployment) (map[string]*dynamodb.AttributeValue, error) {
	val, err := dynamodbattribute.MarshalMap(ddbDeployment{
		Deployment: m,
	})
	if err != nil {
		return nil, err
	}
	// make sure composite attributes don't contain separator characters
	if strings.Contains(m.Application, "--") {
		return nil, fmt.Errorf("application cannot contain '--': %s", m.Application)
	}
	if strings.Contains(m.Environment, "--") {
		return nil, fmt.Errorf("environment cannot contain '--': %s", m.Environment)
	}
	// add in composite attributes
	primaryKey, err := dynamodbattribute.MarshalMap(ddbDeploymentPrimaryKey{
		EnvApp:  fmt.Sprintf("%s--%s", m.Environment, m.Application),
		Version: m.Version,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range primaryKey {
		val[k] = v
	}
	byDate, err := dynamodbattribute.MarshalMap(ddbDeploymentGSIByDate{
		EnvApp: fmt.Sprintf("%s--%s", m.Environment, m.Application),
		Date:   m.Date,
	})
	if err != nil {
		return nil, err
	}
	for k, v := range byDate {
		val[k] = v
	}
	return val, err
}

// decodeDeployment translates a Deployment stored in DynamoDB to a Deployment struct.
func decodeDeployment(m map[string]*dynamodb.AttributeValue, out *models.Deployment) error {
	var ddbDeployment ddbDeployment
	if err := dynamodbattribute.UnmarshalMap(m, &ddbDeployment); err != nil {
		return err
	}
	*out = ddbDeployment.Deployment
	return nil
}

// decodeDeployments translates a list of Deployments stored in DynamoDB to a slice of Deployment structs.
func decodeDeployments(ms []map[string]*dynamodb.AttributeValue) ([]models.Deployment, error) {
	deployments := make([]models.Deployment, len(ms))
	for i, m := range ms {
		var deployment models.Deployment
		if err := decodeDeployment(m, &deployment); err != nil {
			return nil, err
		}
		deployments[i] = deployment
	}
	return deployments, nil
}
