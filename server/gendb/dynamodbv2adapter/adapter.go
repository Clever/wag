package dynamodbv2adapter

import (
	"context"

	v2aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	v1aws "github.com/aws/aws-sdk-go/aws"
	v1request "github.com/aws/aws-sdk-go/aws/request"
	v1dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
)

// Adapter implements dynamodbiface.DynamoDBAPI using a v2 client
type Adapter struct {
	client *dynamodb.Client
}

// New creates a new adapter for the v2 client
func New(client *dynamodb.Client) *Adapter {
	return &Adapter{
		client: client,
	}
}

// convertV1ToV2AttributeValue converts a v1 AttributeValue to a v2 AttributeValue
func convertV1ToV2AttributeValue(v1 *v1dynamodb.AttributeValue) types.AttributeValue {
	if v1 == nil {
		return nil
	}
	if v1.B != nil {
		return &types.AttributeValueMemberB{Value: v1.B}
	}
	if v1.BOOL != nil {
		return &types.AttributeValueMemberBOOL{Value: *v1.BOOL}
	}
	if v1.BS != nil {
		return &types.AttributeValueMemberBS{Value: v1.BS}
	}
	if v1.L != nil {
		l := make([]types.AttributeValue, len(v1.L))
		for i, v := range v1.L {
			l[i] = convertV1ToV2AttributeValue(v)
		}
		return &types.AttributeValueMemberL{Value: l}
	}
	if v1.M != nil {
		m := make(map[string]types.AttributeValue, len(v1.M))
		for k, v := range v1.M {
			m[k] = convertV1ToV2AttributeValue(v)
		}
		return &types.AttributeValueMemberM{Value: m}
	}
	if v1.N != nil {
		return &types.AttributeValueMemberN{Value: *v1.N}
	}
	if v1.NS != nil {
		ns := make([]string, len(v1.NS))
		for i, v := range v1.NS {
			ns[i] = *v
		}
		return &types.AttributeValueMemberNS{Value: ns}
	}
	if v1.NULL != nil {
		return &types.AttributeValueMemberNULL{Value: *v1.NULL}
	}
	if v1.S != nil {
		return &types.AttributeValueMemberS{Value: *v1.S}
	}
	if v1.SS != nil {
		ss := make([]string, len(v1.SS))
		for i, v := range v1.SS {
			ss[i] = *v
		}
		return &types.AttributeValueMemberSS{Value: ss}
	}
	return nil
}

// convertV2ToV1AttributeValue converts a v2 AttributeValue to a v1 AttributeValue
func convertV2ToV1AttributeValue(v2 types.AttributeValue) *v1dynamodb.AttributeValue {
	if v2 == nil {
		return nil
	}
	v1 := &v1dynamodb.AttributeValue{}
	switch v := v2.(type) {
	case *types.AttributeValueMemberB:
		v1.B = v.Value
	case *types.AttributeValueMemberBOOL:
		v1.BOOL = &v.Value
	case *types.AttributeValueMemberBS:
		v1.BS = v.Value
	case *types.AttributeValueMemberL:
		v1.L = make([]*v1dynamodb.AttributeValue, len(v.Value))
		for i, val := range v.Value {
			v1.L[i] = convertV2ToV1AttributeValue(val)
		}
	case *types.AttributeValueMemberM:
		v1.M = make(map[string]*v1dynamodb.AttributeValue, len(v.Value))
		for k, val := range v.Value {
			v1.M[k] = convertV2ToV1AttributeValue(val)
		}
	case *types.AttributeValueMemberN:
		v1.N = &v.Value
	case *types.AttributeValueMemberNS:
		v1.NS = make([]*string, len(v.Value))
		for i, val := range v.Value {
			v1.NS[i] = &val
		}
	case *types.AttributeValueMemberNULL:
		v1.NULL = &v.Value
	case *types.AttributeValueMemberS:
		v1.S = &v.Value
	case *types.AttributeValueMemberSS:
		v1.SS = make([]*string, len(v.Value))
		for i, val := range v.Value {
			v1.SS[i] = &val
		}
	}
	return v1
}

// convertV1ToV2BatchStatementRequest converts a v1 BatchStatementRequest to a v2 BatchStatementRequest
func convertV1ToV2BatchStatementRequest(v1 *v1dynamodb.BatchStatementRequest) types.BatchStatementRequest {
	v2 := types.BatchStatementRequest{
		Statement: v1.Statement,
	}
	if v1.Parameters != nil {
		v2.Parameters = make([]types.AttributeValue, len(v1.Parameters))
		for i, p := range v1.Parameters {
			v2.Parameters[i] = convertV1ToV2AttributeValue(p)
		}
	}
	return v2
}

// convertV2ToV1BatchStatementResponse converts a v2 BatchStatementResponse to a v1 BatchStatementResponse
func convertV2ToV1BatchStatementResponse(v2 types.BatchStatementResponse) *v1dynamodb.BatchStatementResponse {
	v1 := &v1dynamodb.BatchStatementResponse{
		TableName: v2.TableName,
	}
	if v2.Error != nil {
		code := string(v2.Error.Code)
		v1.Error = &v1dynamodb.BatchStatementError{
			Code:    &code,
			Message: v2.Error.Message,
		}
	}
	if v2.Item != nil {
		v1.Item = make(map[string]*v1dynamodb.AttributeValue, len(v2.Item))
		for k, v := range v2.Item {
			v1.Item[k] = convertV2ToV1AttributeValue(v)
		}
	}
	return v1
}

// BatchExecuteStatement implements dynamodbiface.DynamoDBAPI
func (a *Adapter) BatchExecuteStatement(input *v1dynamodb.BatchExecuteStatementInput) (*v1dynamodb.BatchExecuteStatementOutput, error) {
	return a.BatchExecuteStatementWithContext(context.Background(), input)
}

// BatchExecuteStatementWithContext implements dynamodbiface.DynamoDBAPI
func (a *Adapter) BatchExecuteStatementWithContext(ctx context.Context, input *v1dynamodb.BatchExecuteStatementInput, opts ...v1request.Option) (*v1dynamodb.BatchExecuteStatementOutput, error) {
	v2Input := &dynamodb.BatchExecuteStatementInput{
		Statements: make([]types.BatchStatementRequest, len(input.Statements)),
	}

	for i, stmt := range input.Statements {
		v2Input.Statements[i] = types.BatchStatementRequest{
			Statement:  v2aws.String(*stmt.Statement),
			Parameters: convertV1ToV2AttributeValues(stmt.Parameters),
		}
	}

	output, err := a.client.BatchExecuteStatement(ctx, v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.BatchExecuteStatementOutput{
		Responses: make([]*v1dynamodb.BatchStatementResponse, len(output.Responses)),
	}

	for i, resp := range output.Responses {
		v1Output.Responses[i] = &v1dynamodb.BatchStatementResponse{
			Error: &v1dynamodb.BatchStatementError{
				Code:    v1aws.String(string(resp.Error.Code)),
				Message: v1aws.String(*resp.Error.Message),
			},
			Item: convertV2ToV1AttributeValues(resp.Item),
		}
	}

	return v1Output, nil
}

// BatchExecuteStatementRequest implements dynamodbiface.DynamoDBAPI
func (a *Adapter) BatchExecuteStatementRequest(input *v1dynamodb.BatchExecuteStatementInput) (*v1request.Request, *v1dynamodb.BatchExecuteStatementOutput) {
	// Since v2 doesn't have a direct equivalent to Request objects, we'll create a new request
	// and execute it immediately
	req := v1request.New(v1aws.Config{}, v1aws.ClientInfo{}, v1request.Handlers{}, nil, &v1request.Operation{
		Name:       "BatchExecuteStatement",
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}, input, nil)

	output, err := a.BatchExecuteStatement(input)
	if err != nil {
		req.Error = err
	}

	return req, output
}

// BatchGetItem implements dynamodbiface.DynamoDBAPI
func (a *Adapter) BatchGetItem(input *v1dynamodb.BatchGetItemInput) (*v1dynamodb.BatchGetItemOutput, error) {
	v2Input := &dynamodb.BatchGetItemInput{
		RequestItems: make(map[string]types.KeysAndAttributes, len(input.RequestItems)),
	}

	// Convert request items
	for k, v := range input.RequestItems {
		v2Input.RequestItems[k] = convertV1ToV2KeysAndAttributes(v)
	}

	output, err := a.client.BatchGetItem(context.Background(), v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.BatchGetItemOutput{
		ConsumedCapacity: make([]*v1dynamodb.ConsumedCapacity, len(output.ConsumedCapacity)),
		Responses:        make(map[string][]map[string]*v1dynamodb.AttributeValue, len(output.Responses)),
		UnprocessedKeys:  make(map[string]*v1dynamodb.KeysAndAttributes, len(output.UnprocessedKeys)),
	}

	// Convert consumed capacity
	for i, cap := range output.ConsumedCapacity {
		v1Output.ConsumedCapacity[i] = &v1dynamodb.ConsumedCapacity{
			CapacityUnits:          cap.CapacityUnits,
			GlobalSecondaryIndexes: nil, // Conversion needed if used
			LocalSecondaryIndexes:  nil, // Conversion needed if used
			ReadCapacityUnits:      cap.ReadCapacityUnits,
			Table:                  nil, // Conversion needed if used
			TableName:              cap.TableName,
			WriteCapacityUnits:     cap.WriteCapacityUnits,
		}
	}

	// Convert responses
	for k, v := range output.Responses {
		v1Output.Responses[k] = make([]map[string]*v1dynamodb.AttributeValue, len(v))
		for i, item := range v {
			v1Output.Responses[k][i] = make(map[string]*v1dynamodb.AttributeValue, len(item))
			for key, val := range item {
				v1Output.Responses[k][i][key] = convertV2ToV1AttributeValue(val)
			}
		}
	}

	// Convert unprocessed keys
	for k, v := range output.UnprocessedKeys {
		v1Output.UnprocessedKeys[k] = convertV2ToV1KeysAndAttributes(v)
	}

	return v1Output, nil
}

// BatchWriteItem implements dynamodbiface.DynamoDBAPI
func (a *Adapter) BatchWriteItem(input *v1dynamodb.BatchWriteItemInput) (*v1dynamodb.BatchWriteItemOutput, error) {
	v2Input := &dynamodb.BatchWriteItemInput{
		RequestItems: make(map[string][]types.WriteRequest, len(input.RequestItems)),
	}

	// Convert request items
	for k, v := range input.RequestItems {
		v2Input.RequestItems[k] = make([]types.WriteRequest, len(v))
		for i, req := range v {
			v2Input.RequestItems[k][i] = convertV1ToV2WriteRequest(req)
		}
	}

	output, err := a.client.BatchWriteItem(context.Background(), v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.BatchWriteItemOutput{
		ConsumedCapacity:      make([]*v1dynamodb.ConsumedCapacity, len(output.ConsumedCapacity)),
		ItemCollectionMetrics: make(map[string][]*v1dynamodb.ItemCollectionMetrics, len(output.ItemCollectionMetrics)),
		UnprocessedItems:      make(map[string][]*v1dynamodb.WriteRequest, len(output.UnprocessedItems)),
	}

	// Convert consumed capacity
	for i, cap := range output.ConsumedCapacity {
		v1Output.ConsumedCapacity[i] = &v1dynamodb.ConsumedCapacity{
			CapacityUnits:          cap.CapacityUnits,
			GlobalSecondaryIndexes: nil, // Conversion needed if used
			LocalSecondaryIndexes:  nil, // Conversion needed if used
			ReadCapacityUnits:      cap.ReadCapacityUnits,
			Table:                  nil, // Conversion needed if used
			TableName:              cap.TableName,
			WriteCapacityUnits:     cap.WriteCapacityUnits,
		}
	}

	// Convert item collection metrics
	for k, v := range output.ItemCollectionMetrics {
		v1Output.ItemCollectionMetrics[k] = make([]*v1dynamodb.ItemCollectionMetrics, len(v))
		for i := range v {
			v1Output.ItemCollectionMetrics[k][i] = &v1dynamodb.ItemCollectionMetrics{
				ItemCollectionKey:   nil, // Conversion needed if used
				SizeEstimateRangeGB: nil, // Conversion needed if used
			}
		}
	}

	// Convert unprocessed items
	for k, v := range output.UnprocessedItems {
		v1Output.UnprocessedItems[k] = make([]*v1dynamodb.WriteRequest, len(v))
		for i, req := range v {
			v1Output.UnprocessedItems[k][i] = convertV2ToV1WriteRequest(req)
		}
	}

	return v1Output, nil
}

// CreateTable implements dynamodbiface.DynamoDBAPI
func (a *Adapter) CreateTable(input *v1dynamodb.CreateTableInput) (*v1dynamodb.CreateTableOutput, error) {
	// TODO: Implement and use conversion helpers for each field above
	return nil, nil // Placeholder until conversion helpers are implemented
}

// DeleteItem implements dynamodbiface.DynamoDBAPI
func (a *Adapter) DeleteItem(input *v1dynamodb.DeleteItemInput) (*v1dynamodb.DeleteItemOutput, error) {
	v2Input := &dynamodb.DeleteItemInput{
		TableName: input.TableName,
	}
	if input.ConditionExpression != nil {
		v2Input.ConditionExpression = input.ConditionExpression
	}
	if input.ExpressionAttributeNames != nil {
		v2Input.ExpressionAttributeNames = make(map[string]string, len(input.ExpressionAttributeNames))
		for k, v := range input.ExpressionAttributeNames {
			v2Input.ExpressionAttributeNames[k] = *v
		}
	}
	if input.ExpressionAttributeValues != nil {
		v2Input.ExpressionAttributeValues = make(map[string]types.AttributeValue, len(input.ExpressionAttributeValues))
		for k, v := range input.ExpressionAttributeValues {
			v2Input.ExpressionAttributeValues[k] = convertV1ToV2AttributeValue(v)
		}
	}
	if input.Key != nil {
		v2Input.Key = make(map[string]types.AttributeValue, len(input.Key))
		for k, v := range input.Key {
			v2Input.Key[k] = convertV1ToV2AttributeValue(v)
		}
	}
	if input.ReturnConsumedCapacity != nil {
		v2Input.ReturnConsumedCapacity = types.ReturnConsumedCapacity(*input.ReturnConsumedCapacity)
	}
	if input.ReturnItemCollectionMetrics != nil {
		v2Input.ReturnItemCollectionMetrics = types.ReturnItemCollectionMetrics(*input.ReturnItemCollectionMetrics)
	}
	if input.ReturnValues != nil {
		v2Input.ReturnValues = types.ReturnValue(*input.ReturnValues)
	}

	output, err := a.client.DeleteItem(context.Background(), v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.DeleteItemOutput{}
	if output.Attributes != nil {
		v1Output.Attributes = make(map[string]*v1dynamodb.AttributeValue, len(output.Attributes))
		for k, v := range output.Attributes {
			v1Output.Attributes[k] = convertV2ToV1AttributeValue(v)
		}
	}
	if output.ConsumedCapacity != nil {
		v1Output.ConsumedCapacity = &v1dynamodb.ConsumedCapacity{
			CapacityUnits:      output.ConsumedCapacity.CapacityUnits,
			ReadCapacityUnits:  output.ConsumedCapacity.ReadCapacityUnits,
			WriteCapacityUnits: output.ConsumedCapacity.WriteCapacityUnits,
			TableName:          output.ConsumedCapacity.TableName,
		}
	}
	if output.ItemCollectionMetrics != nil {
		v1Output.ItemCollectionMetrics = &v1dynamodb.ItemCollectionMetrics{
			ItemCollectionKey:   nil, // TODO: Implement conversion if needed
			SizeEstimateRangeGB: nil, // TODO: Implement conversion if needed
		}
	}

	return v1Output, nil
}

// DeleteTable implements dynamodbiface.DynamoDBAPI
func (a *Adapter) DeleteTable(input *v1dynamodb.DeleteTableInput) (*v1dynamodb.DeleteTableOutput, error) {
	v2Input := &dynamodb.DeleteTableInput{
		TableName: input.TableName,
	}

	output, err := a.client.DeleteTable(context.Background(), v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.DeleteTableOutput{
		TableDescription: &v1dynamodb.TableDescription{
			ArchivalSummary:           output.TableDescription.ArchivalSummary,
			AttributeDefinitions:      make([]*v1dynamodb.AttributeDefinition, len(output.TableDescription.AttributeDefinitions)),
			BillingModeSummary:        output.TableDescription.BillingModeSummary,
			CreationDateTime:          output.TableDescription.CreationDateTime,
			DeletionProtectionEnabled: output.TableDescription.DeletionProtectionEnabled,
			GlobalSecondaryIndexes:    make([]*v1dynamodb.GlobalSecondaryIndexDescription, len(output.TableDescription.GlobalSecondaryIndexes)),
			GlobalTableVersion:        output.TableDescription.GlobalTableVersion,
			ItemCount:                 output.TableDescription.ItemCount,
			KeySchema:                 make([]*v1dynamodb.KeySchemaElement, len(output.TableDescription.KeySchema)),
			LatestStreamArn:           output.TableDescription.LatestStreamArn,
			LatestStreamLabel:         output.TableDescription.LatestStreamLabel,
			LocalSecondaryIndexes:     make([]*v1dynamodb.LocalSecondaryIndexDescription, len(output.TableDescription.LocalSecondaryIndexes)),
			ProvisionedThroughput:     output.TableDescription.ProvisionedThroughput,
			Replicas:                  make([]*v1dynamodb.ReplicaDescription, len(output.TableDescription.Replicas)),
			RestoreSummary:            output.TableDescription.RestoreSummary,
			SSEDescription:            output.TableDescription.SSEDescription,
			StreamSpecification:       output.TableDescription.StreamSpecification,
			TableArn:                  output.TableDescription.TableArn,
			TableId:                   output.TableDescription.TableId,
			TableName:                 output.TableDescription.TableName,
			TableSizeBytes:            output.TableDescription.TableSizeBytes,
			TableStatus:               output.TableDescription.TableStatus,
		},
	}

	// Convert attribute definitions
	for i, def := range output.TableDescription.AttributeDefinitions {
		v1Output.TableDescription.AttributeDefinitions[i] = &v1dynamodb.AttributeDefinition{
			AttributeName: def.AttributeName,
			AttributeType: def.AttributeType,
		}
	}

	// Convert global secondary indexes
	for i, gsi := range output.TableDescription.GlobalSecondaryIndexes {
		v1Output.TableDescription.GlobalSecondaryIndexes[i] = &v1dynamodb.GlobalSecondaryIndexDescription{
			Backfilling:           gsi.Backfilling,
			IndexArn:              gsi.IndexArn,
			IndexName:             gsi.IndexName,
			IndexSizeBytes:        gsi.IndexSizeBytes,
			IndexStatus:           gsi.IndexStatus,
			ItemCount:             gsi.ItemCount,
			KeySchema:             make([]*v1dynamodb.KeySchemaElement, len(gsi.KeySchema)),
			Projection:            gsi.Projection,
			ProvisionedThroughput: gsi.ProvisionedThroughput,
		}
		for j, key := range gsi.KeySchema {
			v1Output.TableDescription.GlobalSecondaryIndexes[i].KeySchema[j] = &v1dynamodb.KeySchemaElement{
				AttributeName: key.AttributeName,
				KeyType:       key.KeyType,
			}
		}
	}

	// Convert key schema
	for i, key := range output.TableDescription.KeySchema {
		v1Output.TableDescription.KeySchema[i] = &v1dynamodb.KeySchemaElement{
			AttributeName: key.AttributeName,
			KeyType:       key.KeyType,
		}
	}

	// Convert local secondary indexes
	for i, lsi := range output.TableDescription.LocalSecondaryIndexes {
		v1Output.TableDescription.LocalSecondaryIndexes[i] = &v1dynamodb.LocalSecondaryIndexDescription{
			IndexArn:       lsi.IndexArn,
			IndexName:      lsi.IndexName,
			IndexSizeBytes: lsi.IndexSizeBytes,
			ItemCount:      lsi.ItemCount,
			KeySchema:      make([]*v1dynamodb.KeySchemaElement, len(lsi.KeySchema)),
			Projection:     lsi.Projection,
		}
		for j, key := range lsi.KeySchema {
			v1Output.TableDescription.LocalSecondaryIndexes[i].KeySchema[j] = &v1dynamodb.KeySchemaElement{
				AttributeName: key.AttributeName,
				KeyType:       key.KeyType,
			}
		}
	}

	return v1Output, nil
}

// DescribeTable implements dynamodbiface.DynamoDBAPI
func (a *Adapter) DescribeTable(input *v1dynamodb.DescribeTableInput) (*v1dynamodb.DescribeTableOutput, error) {
	v2Input := &dynamodb.DescribeTableInput{
		TableName: input.TableName,
	}

	output, err := a.client.DescribeTable(context.Background(), v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.DescribeTableOutput{
		Table: &v1dynamodb.TableDescription{
			ArchivalSummary:           output.Table.ArchivalSummary,
			AttributeDefinitions:      make([]*v1dynamodb.AttributeDefinition, len(output.Table.AttributeDefinitions)),
			BillingModeSummary:        output.Table.BillingModeSummary,
			CreationDateTime:          output.Table.CreationDateTime,
			DeletionProtectionEnabled: output.Table.DeletionProtectionEnabled,
			GlobalSecondaryIndexes:    make([]*v1dynamodb.GlobalSecondaryIndexDescription, len(output.Table.GlobalSecondaryIndexes)),
			GlobalTableVersion:        output.Table.GlobalTableVersion,
			ItemCount:                 output.Table.ItemCount,
			KeySchema:                 make([]*v1dynamodb.KeySchemaElement, len(output.Table.KeySchema)),
			LatestStreamArn:           output.Table.LatestStreamArn,
			LatestStreamLabel:         output.Table.LatestStreamLabel,
			LocalSecondaryIndexes:     make([]*v1dynamodb.LocalSecondaryIndexDescription, len(output.Table.LocalSecondaryIndexes)),
			ProvisionedThroughput:     output.Table.ProvisionedThroughput,
			Replicas:                  make([]*v1dynamodb.ReplicaDescription, len(output.Table.Replicas)),
			RestoreSummary:            output.Table.RestoreSummary,
			SSEDescription:            output.Table.SSEDescription,
			StreamSpecification:       output.Table.StreamSpecification,
			TableArn:                  output.Table.TableArn,
			TableId:                   output.Table.TableId,
			TableName:                 output.Table.TableName,
			TableSizeBytes:            output.Table.TableSizeBytes,
			TableStatus:               output.Table.TableStatus,
		},
	}

	// Convert attribute definitions
	for i, def := range output.Table.AttributeDefinitions {
		v1Output.Table.AttributeDefinitions[i] = &v1dynamodb.AttributeDefinition{
			AttributeName: def.AttributeName,
			AttributeType: def.AttributeType,
		}
	}

	// Convert global secondary indexes
	for i, gsi := range output.Table.GlobalSecondaryIndexes {
		v1Output.Table.GlobalSecondaryIndexes[i] = &v1dynamodb.GlobalSecondaryIndexDescription{
			Backfilling:           gsi.Backfilling,
			IndexArn:              gsi.IndexArn,
			IndexName:             gsi.IndexName,
			IndexSizeBytes:        gsi.IndexSizeBytes,
			IndexStatus:           gsi.IndexStatus,
			ItemCount:             gsi.ItemCount,
			KeySchema:             make([]*v1dynamodb.KeySchemaElement, len(gsi.KeySchema)),
			Projection:            gsi.Projection,
			ProvisionedThroughput: gsi.ProvisionedThroughput,
		}
		for j, key := range gsi.KeySchema {
			v1Output.Table.GlobalSecondaryIndexes[i].KeySchema[j] = &v1dynamodb.KeySchemaElement{
				AttributeName: key.AttributeName,
				KeyType:       key.KeyType,
			}
		}
	}

	// Convert key schema
	for i, key := range output.Table.KeySchema {
		v1Output.Table.KeySchema[i] = &v1dynamodb.KeySchemaElement{
			AttributeName: key.AttributeName,
			KeyType:       key.KeyType,
		}
	}

	// Convert local secondary indexes
	for i, lsi := range output.Table.LocalSecondaryIndexes {
		v1Output.Table.LocalSecondaryIndexes[i] = &v1dynamodb.LocalSecondaryIndexDescription{
			IndexArn:       lsi.IndexArn,
			IndexName:      lsi.IndexName,
			IndexSizeBytes: lsi.IndexSizeBytes,
			ItemCount:      lsi.ItemCount,
			KeySchema:      make([]*v1dynamodb.KeySchemaElement, len(lsi.KeySchema)),
			Projection:     lsi.Projection,
		}
		for j, key := range lsi.KeySchema {
			v1Output.Table.LocalSecondaryIndexes[i].KeySchema[j] = &v1dynamodb.KeySchemaElement{
				AttributeName: key.AttributeName,
				KeyType:       key.KeyType,
			}
		}
	}

	return v1Output, nil
}

// GetItem implements dynamodbiface.DynamoDBAPI
func (a *Adapter) GetItem(input *v1dynamodb.GetItemInput) (*v1dynamodb.GetItemOutput, error) {
	v2Input := &dynamodb.GetItemInput{
		AttributesToGet:          input.AttributesToGet,
		ConsistentRead:           input.ConsistentRead,
		ExpressionAttributeNames: input.ExpressionAttributeNames,
		Key:                      input.Key,
		ProjectionExpression:     input.ProjectionExpression,
		ReturnConsumedCapacity:   input.ReturnConsumedCapacity,
		TableName:                input.TableName,
	}

	output, err := a.client.GetItem(context.Background(), v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.GetItemOutput{
		ConsumedCapacity: output.ConsumedCapacity,
		Item:             output.Item,
	}

	return v1Output, nil
}

// PutItem implements dynamodbiface.DynamoDBAPI
func (a *Adapter) PutItem(input *v1dynamodb.PutItemInput) (*v1dynamodb.PutItemOutput, error) {
	v2Input := &dynamodb.PutItemInput{
		ConditionExpression:         input.ConditionExpression,
		ConditionalOperator:         input.ConditionalOperator,
		Expected:                    input.Expected,
		ExpressionAttributeNames:    input.ExpressionAttributeNames,
		ExpressionAttributeValues:   input.ExpressionAttributeValues,
		Item:                        input.Item,
		ReturnConsumedCapacity:      input.ReturnConsumedCapacity,
		ReturnItemCollectionMetrics: input.ReturnItemCollectionMetrics,
		ReturnValues:                input.ReturnValues,
		TableName:                   input.TableName,
	}

	output, err := a.client.PutItem(context.Background(), v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.PutItemOutput{
		Attributes:            output.Attributes,
		ConsumedCapacity:      output.ConsumedCapacity,
		ItemCollectionMetrics: output.ItemCollectionMetrics,
	}

	return v1Output, nil
}

// Query implements dynamodbiface.DynamoDBAPI
func (a *Adapter) Query(input *v1dynamodb.QueryInput) (*v1dynamodb.QueryOutput, error) {
	v2Input := &dynamodb.QueryInput{
		AttributesToGet:           input.AttributesToGet,
		ConditionalOperator:       input.ConditionalOperator,
		ConsistentRead:            input.ConsistentRead,
		ExclusiveStartKey:         input.ExclusiveStartKey,
		ExpressionAttributeNames:  input.ExpressionAttributeNames,
		ExpressionAttributeValues: input.ExpressionAttributeValues,
		FilterExpression:          input.FilterExpression,
		IndexName:                 input.IndexName,
		KeyConditionExpression:    input.KeyConditionExpression,
		KeyConditions:             input.KeyConditions,
		Limit:                     input.Limit,
		ProjectionExpression:      input.ProjectionExpression,
		QueryFilter:               input.QueryFilter,
		ReturnConsumedCapacity:    input.ReturnConsumedCapacity,
		ScanIndexForward:          input.ScanIndexForward,
		Select:                    input.Select,
		TableName:                 input.TableName,
	}

	output, err := a.client.Query(context.Background(), v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.QueryOutput{
		ConsumedCapacity: output.ConsumedCapacity,
		Count:            output.Count,
		Items:            output.Items,
		LastEvaluatedKey: output.LastEvaluatedKey,
		ScannedCount:     output.ScannedCount,
	}

	return v1Output, nil
}

// Scan implements dynamodbiface.DynamoDBAPI
func (a *Adapter) Scan(input *v1dynamodb.ScanInput) (*v1dynamodb.ScanOutput, error) {
	v2Input := &dynamodb.ScanInput{
		AttributesToGet:           input.AttributesToGet,
		ConditionalOperator:       input.ConditionalOperator,
		ConsistentRead:            input.ConsistentRead,
		ExclusiveStartKey:         input.ExclusiveStartKey,
		ExpressionAttributeNames:  input.ExpressionAttributeNames,
		ExpressionAttributeValues: input.ExpressionAttributeValues,
		FilterExpression:          input.FilterExpression,
		IndexName:                 input.IndexName,
		Limit:                     input.Limit,
		ProjectionExpression:      input.ProjectionExpression,
		ReturnConsumedCapacity:    input.ReturnConsumedCapacity,
		ScanFilter:                input.ScanFilter,
		Segment:                   input.Segment,
		Select:                    input.Select,
		TableName:                 input.TableName,
		TotalSegments:             input.TotalSegments,
	}

	output, err := a.client.Scan(context.Background(), v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.ScanOutput{
		ConsumedCapacity: output.ConsumedCapacity,
		Count:            output.Count,
		Items:            output.Items,
		LastEvaluatedKey: output.LastEvaluatedKey,
		ScannedCount:     output.ScannedCount,
	}

	return v1Output, nil
}

// UpdateItem implements dynamodbiface.DynamoDBAPI
func (a *Adapter) UpdateItem(input *v1dynamodb.UpdateItemInput) (*v1dynamodb.UpdateItemOutput, error) {
	v2Input := &dynamodb.UpdateItemInput{
		AttributeUpdates:            input.AttributeUpdates,
		ConditionExpression:         input.ConditionExpression,
		ConditionalOperator:         input.ConditionalOperator,
		Expected:                    input.Expected,
		ExpressionAttributeNames:    input.ExpressionAttributeNames,
		ExpressionAttributeValues:   input.ExpressionAttributeValues,
		Key:                         input.Key,
		ReturnConsumedCapacity:      input.ReturnConsumedCapacity,
		ReturnItemCollectionMetrics: input.ReturnItemCollectionMetrics,
		ReturnValues:                input.ReturnValues,
		TableName:                   input.TableName,
		UpdateExpression:            input.UpdateExpression,
	}

	output, err := a.client.UpdateItem(context.Background(), v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.UpdateItemOutput{
		Attributes:            output.Attributes,
		ConsumedCapacity:      output.ConsumedCapacity,
		ItemCollectionMetrics: output.ItemCollectionMetrics,
	}

	return v1Output, nil
}

// UpdateTable implements dynamodbiface.DynamoDBAPI
func (a *Adapter) UpdateTable(input *v1dynamodb.UpdateTableInput) (*v1dynamodb.UpdateTableOutput, error) {
	v2Input := &dynamodb.UpdateTableInput{
		AttributeDefinitions:        input.AttributeDefinitions,
		BillingMode:                 input.BillingMode,
		DeletionProtectionEnabled:   input.DeletionProtectionEnabled,
		GlobalSecondaryIndexUpdates: input.GlobalSecondaryIndexUpdates,
		ProvisionedThroughput:       input.ProvisionedThroughput,
		SSESpecification:            input.SSESpecification,
		StreamSpecification:         input.StreamSpecification,
		TableName:                   input.TableName,
	}

	output, err := a.client.UpdateTable(context.Background(), v2Input)
	if err != nil {
		return nil, err
	}

	v1Output := &v1dynamodb.UpdateTableOutput{
		TableDescription: &v1dynamodb.TableDescription{
			ArchivalSummary:           output.TableDescription.ArchivalSummary,
			AttributeDefinitions:      make([]*v1dynamodb.AttributeDefinition, len(output.TableDescription.AttributeDefinitions)),
			BillingModeSummary:        output.TableDescription.BillingModeSummary,
			CreationDateTime:          output.TableDescription.CreationDateTime,
			DeletionProtectionEnabled: output.TableDescription.DeletionProtectionEnabled,
			GlobalSecondaryIndexes:    make([]*v1dynamodb.GlobalSecondaryIndexDescription, len(output.TableDescription.GlobalSecondaryIndexes)),
			GlobalTableVersion:        output.TableDescription.GlobalTableVersion,
			ItemCount:                 output.TableDescription.ItemCount,
			KeySchema:                 make([]*v1dynamodb.KeySchemaElement, len(output.TableDescription.KeySchema)),
			LatestStreamArn:           output.TableDescription.LatestStreamArn,
			LatestStreamLabel:         output.TableDescription.LatestStreamLabel,
			LocalSecondaryIndexes:     make([]*v1dynamodb.LocalSecondaryIndexDescription, len(output.TableDescription.LocalSecondaryIndexes)),
			ProvisionedThroughput:     output.TableDescription.ProvisionedThroughput,
			Replicas:                  make([]*v1dynamodb.ReplicaDescription, len(output.TableDescription.Replicas)),
			RestoreSummary:            output.TableDescription.RestoreSummary,
			SSEDescription:            output.TableDescription.SSEDescription,
			StreamSpecification:       output.TableDescription.StreamSpecification,
			TableArn:                  output.TableDescription.TableArn,
			TableId:                   output.TableDescription.TableId,
			TableName:                 output.TableDescription.TableName,
			TableSizeBytes:            output.TableDescription.TableSizeBytes,
			TableStatus:               output.TableDescription.TableStatus,
		},
	}

	// Convert attribute definitions
	for i, def := range output.TableDescription.AttributeDefinitions {
		v1Output.TableDescription.AttributeDefinitions[i] = &v1dynamodb.AttributeDefinition{
			AttributeName: def.AttributeName,
			AttributeType: def.AttributeType,
		}
	}

	// Convert global secondary indexes
	for i, gsi := range output.TableDescription.GlobalSecondaryIndexes {
		v1Output.TableDescription.GlobalSecondaryIndexes[i] = &v1dynamodb.GlobalSecondaryIndexDescription{
			Backfilling:           gsi.Backfilling,
			IndexArn:              gsi.IndexArn,
			IndexName:             gsi.IndexName,
			IndexSizeBytes:        gsi.IndexSizeBytes,
			IndexStatus:           gsi.IndexStatus,
			ItemCount:             gsi.ItemCount,
			KeySchema:             make([]*v1dynamodb.KeySchemaElement, len(gsi.KeySchema)),
			Projection:            gsi.Projection,
			ProvisionedThroughput: gsi.ProvisionedThroughput,
		}
		for j, key := range gsi.KeySchema {
			v1Output.TableDescription.GlobalSecondaryIndexes[i].KeySchema[j] = &v1dynamodb.KeySchemaElement{
				AttributeName: key.AttributeName,
				KeyType:       key.KeyType,
			}
		}
	}

	// Convert key schema
	for i, key := range output.TableDescription.KeySchema {
		v1Output.TableDescription.KeySchema[i] = &v1dynamodb.KeySchemaElement{
			AttributeName: key.AttributeName,
			KeyType:       key.KeyType,
		}
	}

	// Convert local secondary indexes
	for i, lsi := range output.TableDescription.LocalSecondaryIndexes {
		v1Output.TableDescription.LocalSecondaryIndexes[i] = &v1dynamodb.LocalSecondaryIndexDescription{
			IndexArn:       lsi.IndexArn,
			IndexName:      lsi.IndexName,
			IndexSizeBytes: lsi.IndexSizeBytes,
			ItemCount:      lsi.ItemCount,
			KeySchema:      make([]*v1dynamodb.KeySchemaElement, len(lsi.KeySchema)),
			Projection:     lsi.Projection,
		}
		for j, key := range lsi.KeySchema {
			v1Output.TableDescription.LocalSecondaryIndexes[i].KeySchema[j] = &v1dynamodb.KeySchemaElement{
				AttributeName: key.AttributeName,
				KeyType:       key.KeyType,
			}
		}
	}

	return v1Output, nil
}

// WaitUntilTableExists implements dynamodbiface.DynamoDBAPI
func (a *Adapter) WaitUntilTableExists(input *v1dynamodb.DescribeTableInput) error {
	return a.client.WaitUntilTableExists(context.Background(), &dynamodb.DescribeTableInput{
		TableName: input.TableName,
	})
}

// WaitUntilTableNotExists implements dynamodbiface.DynamoDBAPI
func (a *Adapter) WaitUntilTableNotExists(input *v1dynamodb.DescribeTableInput) error {
	return a.client.WaitUntilTableNotExists(context.Background(), &dynamodb.DescribeTableInput{
		TableName: input.TableName,
	})
}

// convertV1ToV2KeysAndAttributes converts a v1 KeysAndAttributes to a v2 KeysAndAttributes
func convertV1ToV2KeysAndAttributes(v1 *v1dynamodb.KeysAndAttributes) types.KeysAndAttributes {
	v2 := types.KeysAndAttributes{
		ConsistentRead: v1.ConsistentRead,
	}
	if v1.AttributesToGet != nil {
		v2.AttributesToGet = make([]string, len(v1.AttributesToGet))
		for i, v := range v1.AttributesToGet {
			v2.AttributesToGet[i] = *v
		}
	}
	if v1.ExpressionAttributeNames != nil {
		v2.ExpressionAttributeNames = make(map[string]string, len(v1.ExpressionAttributeNames))
		for k, v := range v1.ExpressionAttributeNames {
			v2.ExpressionAttributeNames[k] = *v
		}
	}
	if v1.Keys != nil {
		v2.Keys = make([]map[string]types.AttributeValue, len(v1.Keys))
		for i, k := range v1.Keys {
			v2.Keys[i] = make(map[string]types.AttributeValue, len(k))
			for k2, v := range k {
				v2.Keys[i][k2] = convertV1ToV2AttributeValue(v)
			}
		}
	}
	if v1.ProjectionExpression != nil {
		v2.ProjectionExpression = v1.ProjectionExpression
	}
	return v2
}

// convertV2ToV1KeysAndAttributes converts a v2 KeysAndAttributes to a v1 KeysAndAttributes
func convertV2ToV1KeysAndAttributes(v2 types.KeysAndAttributes) *v1dynamodb.KeysAndAttributes {
	v1 := &v1dynamodb.KeysAndAttributes{
		ConsistentRead: v2.ConsistentRead,
	}
	if v2.AttributesToGet != nil {
		v1.AttributesToGet = make([]*string, len(v2.AttributesToGet))
		for i, v := range v2.AttributesToGet {
			v1.AttributesToGet[i] = &v
		}
	}
	if v2.ExpressionAttributeNames != nil {
		v1.ExpressionAttributeNames = make(map[string]*string, len(v2.ExpressionAttributeNames))
		for k, v := range v2.ExpressionAttributeNames {
			v1.ExpressionAttributeNames[k] = &v
		}
	}
	if v2.Keys != nil {
		v1.Keys = make([]map[string]*v1dynamodb.AttributeValue, len(v2.Keys))
		for i, k := range v2.Keys {
			v1.Keys[i] = make(map[string]*v1dynamodb.AttributeValue, len(k))
			for k2, v := range k {
				v1.Keys[i][k2] = convertV2ToV1AttributeValue(v)
			}
		}
	}
	if v2.ProjectionExpression != nil {
		v1.ProjectionExpression = v2.ProjectionExpression
	}
	return v1
}

// convertV1ToV2Capacity converts a v1 Capacity to a v2 Capacity
func convertV1ToV2Capacity(v1 *v1dynamodb.Capacity) types.Capacity {
	v2 := types.Capacity{
		CapacityUnits:      v1.CapacityUnits,
		ReadCapacityUnits:  v1.ReadCapacityUnits,
		WriteCapacityUnits: v1.WriteCapacityUnits,
	}
	if v1.GlobalSecondaryIndexes != nil {
		v2.GlobalSecondaryIndexes = make(map[string]types.Capacity, len(v1.GlobalSecondaryIndexes))
		for k, v := range v1.GlobalSecondaryIndexes {
			v2.GlobalSecondaryIndexes[k] = convertV1ToV2Capacity(v)
		}
	}
	if v1.LocalSecondaryIndexes != nil {
		v2.LocalSecondaryIndexes = make(map[string]types.Capacity, len(v1.LocalSecondaryIndexes))
		for k, v := range v1.LocalSecondaryIndexes {
			v2.LocalSecondaryIndexes[k] = convertV1ToV2Capacity(v)
		}
	}
	if v1.Table != nil {
		v2.Table = &convertV1ToV2Capacity(v1.Table)
	}
	return v2
}

// convertV2ToV1Capacity converts a v2 Capacity to a v1 Capacity
func convertV2ToV1Capacity(v2 types.Capacity) *v1dynamodb.Capacity {
	v1 := &v1dynamodb.Capacity{
		CapacityUnits:      v2.CapacityUnits,
		ReadCapacityUnits:  v2.ReadCapacityUnits,
		WriteCapacityUnits: v2.WriteCapacityUnits,
	}
	if v2.GlobalSecondaryIndexes != nil {
		v1.GlobalSecondaryIndexes = make(map[string]*v1dynamodb.Capacity, len(v2.GlobalSecondaryIndexes))
		for k, v := range v2.GlobalSecondaryIndexes {
			v1.GlobalSecondaryIndexes[k] = &convertV2ToV1Capacity(v)
		}
	}
	if v2.LocalSecondaryIndexes != nil {
		v1.LocalSecondaryIndexes = make(map[string]*v1dynamodb.Capacity, len(v2.LocalSecondaryIndexes))
		for k, v := range v2.LocalSecondaryIndexes {
			v1.LocalSecondaryIndexes[k] = &convertV2ToV1Capacity(v)
		}
	}
	if v2.Table != nil {
		v1.Table = &convertV2ToV1Capacity(*v2.Table)
	}
	return v1
}

// convertV1ToV2WriteRequest converts a v1 WriteRequest to a v2 WriteRequest
func convertV1ToV2WriteRequest(v1 *v1dynamodb.WriteRequest) types.WriteRequest {
	v2 := types.WriteRequest{}
	if v1.DeleteRequest != nil {
		v2.DeleteRequest = &types.DeleteRequest{
			Key: make(map[string]types.AttributeValue, len(v1.DeleteRequest.Key)),
		}
		for k, v := range v1.DeleteRequest.Key {
			v2.DeleteRequest.Key[k] = convertV1ToV2AttributeValue(v)
		}
	}
	if v1.PutRequest != nil {
		v2.PutRequest = &types.PutRequest{
			Item: make(map[string]types.AttributeValue, len(v1.PutRequest.Item)),
		}
		for k, v := range v1.PutRequest.Item {
			v2.PutRequest.Item[k] = convertV1ToV2AttributeValue(v)
		}
	}
	return v2
}

// convertV2ToV1WriteRequest converts a v2 WriteRequest to a v1 WriteRequest
func convertV2ToV1WriteRequest(v2 types.WriteRequest) *v1dynamodb.WriteRequest {
	v1 := &v1dynamodb.WriteRequest{}
	if v2.DeleteRequest != nil {
		v1.DeleteRequest = &v1dynamodb.DeleteRequest{
			Key: make(map[string]*v1dynamodb.AttributeValue, len(v2.DeleteRequest.Key)),
		}
		for k, v := range v2.DeleteRequest.Key {
			v1.DeleteRequest.Key[k] = convertV2ToV1AttributeValue(v)
		}
	}
	if v2.PutRequest != nil {
		v1.PutRequest = &v1dynamodb.PutRequest{
			Item: make(map[string]*v1dynamodb.AttributeValue, len(v2.PutRequest.Item)),
		}
		for k, v := range v2.PutRequest.Item {
			v1.PutRequest.Item[k] = convertV2ToV1AttributeValue(v)
		}
	}
	return v1
}

// convertV1ToV2TableDescription converts a v1 TableDescription to v2 TableDescription
func convertV1ToV2TableDescription(v1 *v1dynamodb.TableDescription) *types.TableDescription {
	if v1 == nil {
		return nil
	}
	return &types.TableDescription{
		ArchivalSummary:        convertV1ToV2ArchivalSummary(v1.ArchivalSummary),
		AttributeDefinitions:   convertV1ToV2AttributeDefinitions(v1.AttributeDefinitions),
		BillingModeSummary:     convertV1ToV2BillingModeSummary(v1.BillingModeSummary),
		CreationDateTime:       v1.CreationDateTime,
		GlobalSecondaryIndexes: convertV1ToV2GlobalSecondaryIndexes(v1.GlobalSecondaryIndexes),
		GlobalTableVersion:     v1.GlobalTableVersion,
		ItemCount:              v1.ItemCount,
		KeySchema:              convertV1ToV2KeySchema(v1.KeySchema),
		LatestStreamArn:        v1.LatestStreamArn,
		LatestStreamLabel:      v1.LatestStreamLabel,
		LocalSecondaryIndexes:  convertV1ToV2LocalSecondaryIndexes(v1.LocalSecondaryIndexes),
		ProvisionedThroughput:  convertV1ToV2ProvisionedThroughputDescription(v1.ProvisionedThroughput),
		Replicas:               convertV1ToV2Replicas(v1.Replicas),
		RestoreSummary:         convertV1ToV2RestoreSummary(v1.RestoreSummary),
		SSEDescription:         convertV1ToV2SSEDescription(v1.SSEDescription),
		StreamSpecification:    convertV1ToV2StreamSpecification(v1.StreamSpecification),
		TableArn:               v1.TableArn,
		TableId:                v1.TableId,
		TableName:              v1.TableName,
		TableSizeBytes:         v1.TableSizeBytes,
		TableStatus:            types.TableStatus(*v1.TableStatus),
	}
}

// convertV2ToV1TableDescription converts a v2 TableDescription to v1 TableDescription
func convertV2ToV1TableDescription(v2 *types.TableDescription) *v1dynamodb.TableDescription {
	if v2 == nil {
		return nil
	}
	status := string(v2.TableStatus)
	return &v1dynamodb.TableDescription{
		ArchivalSummary:        convertV2ToV1ArchivalSummary(v2.ArchivalSummary),
		AttributeDefinitions:   convertV2ToV1AttributeDefinitions(v2.AttributeDefinitions),
		BillingModeSummary:     convertV2ToV1BillingModeSummary(v2.BillingModeSummary),
		CreationDateTime:       v2.CreationDateTime,
		GlobalSecondaryIndexes: convertV2ToV1GlobalSecondaryIndexes(v2.GlobalSecondaryIndexes),
		GlobalTableVersion:     v2.GlobalTableVersion,
		ItemCount:              v2.ItemCount,
		KeySchema:              convertV2ToV1KeySchema(v2.KeySchema),
		LatestStreamArn:        v2.LatestStreamArn,
		LatestStreamLabel:      v2.LatestStreamLabel,
		LocalSecondaryIndexes:  convertV2ToV1LocalSecondaryIndexes(v2.LocalSecondaryIndexes),
		ProvisionedThroughput:  convertV2ToV1ProvisionedThroughputDescription(v2.ProvisionedThroughput),
		Replicas:               convertV2ToV1Replicas(v2.Replicas),
		RestoreSummary:         convertV2ToV1RestoreSummary(v2.RestoreSummary),
		SSEDescription:         convertV2ToV1SSEDescription(v2.SSEDescription),
		StreamSpecification:    convertV2ToV1StreamSpecification(v2.StreamSpecification),
		TableArn:               v2.TableArn,
		TableId:                v2.TableId,
		TableName:              v2.TableName,
		TableSizeBytes:         v2.TableSizeBytes,
		TableStatus:            &status,
	}
}

// convertV1ToV2AttributeDefinitions converts v1 AttributeDefinitions to v2 AttributeDefinitions
func convertV1ToV2AttributeDefinitions(v1 []*v1dynamodb.AttributeDefinition) []types.AttributeDefinition {
	if v1 == nil {
		return nil
	}
	v2 := make([]types.AttributeDefinition, len(v1))
	for i, def := range v1 {
		v2[i] = types.AttributeDefinition{
			AttributeName: v2aws.String(*def.AttributeName),
			AttributeType: types.ScalarAttributeType(*def.AttributeType),
		}
	}
	return v2
}

// convertV2ToV1AttributeDefinitions converts v2 AttributeDefinitions to v1 AttributeDefinitions
func convertV2ToV1AttributeDefinitions(v2 []types.AttributeDefinition) []*v1dynamodb.AttributeDefinition {
	if v2 == nil {
		return nil
	}
	v1 := make([]*v1dynamodb.AttributeDefinition, len(v2))
	for i, def := range v2 {
		attrType := string(def.AttributeType)
		v1[i] = &v1dynamodb.AttributeDefinition{
			AttributeName: v1aws.String(*def.AttributeName),
			AttributeType: &attrType,
		}
	}
	return v1
}

// convertV1ToV2GlobalSecondaryIndexes converts v1 GlobalSecondaryIndexes to v2 GlobalSecondaryIndexes
func convertV1ToV2GlobalSecondaryIndexes(v1 []*v1dynamodb.GlobalSecondaryIndex) []types.GlobalSecondaryIndex {
	if v1 == nil {
		return nil
	}
	v2 := make([]types.GlobalSecondaryIndex, len(v1))
	for i, gsi := range v1 {
		v2[i] = types.GlobalSecondaryIndex{
			IndexName:             v2aws.String(*gsi.IndexName),
			KeySchema:             convertV1ToV2KeySchema(gsi.KeySchema),
			Projection:            convertV1ToV2Projection(gsi.Projection),
			ProvisionedThroughput: convertV1ToV2ProvisionedThroughput(gsi.ProvisionedThroughput),
		}
	}
	return v2
}

// convertV2ToV1GlobalSecondaryIndexes converts v2 GlobalSecondaryIndexes to v1 GlobalSecondaryIndexes
func convertV2ToV1GlobalSecondaryIndexes(v2 []types.GlobalSecondaryIndex) []*v1dynamodb.GlobalSecondaryIndex {
	if v2 == nil {
		return nil
	}
	v1 := make([]*v1dynamodb.GlobalSecondaryIndex, len(v2))
	for i, gsi := range v2 {
		status := string(gsi.IndexStatus)
		v1[i] = &v1dynamodb.GlobalSecondaryIndex{
			IndexName:             v1aws.String(*gsi.IndexName),
			KeySchema:             convertV2ToV1KeySchema(gsi.KeySchema),
			Projection:            convertV2ToV1Projection(gsi.Projection),
			ProvisionedThroughput: convertV2ToV1ProvisionedThroughput(gsi.ProvisionedThroughput),
			IndexStatus:           &status,
		}
	}
	return v1
}

// convertV1ToV2ArchivalSummary converts v1 ArchivalSummary to v2 ArchivalSummary
func convertV1ToV2ArchivalSummary(v1 *v1dynamodb.ArchivalSummary) types.ArchivalSummary {
	if v1 == nil {
		return nil
	}
	return &types.ArchivalSummary{
		ArchivalDateTime: v1.ArchivalDateTime,
		ArchivalReason:   v1.ArchivalReason,
	}
}

// convertV2ToV1ArchivalSummary converts v2 ArchivalSummary to v1 ArchivalSummary
func convertV2ToV1ArchivalSummary(v2 types.ArchivalSummary) *v1dynamodb.ArchivalSummary {
	if v2 == nil {
		return nil
	}
	return &v1dynamodb.ArchivalSummary{
		ArchivalDateTime: v2.ArchivalDateTime,
		ArchivalReason:   v2.ArchivalReason,
	}
}

// convertV1ToV2BillingModeSummary converts v1 BillingModeSummary to v2 BillingModeSummary
func convertV1ToV2BillingModeSummary(v1 *v1dynamodb.BillingModeSummary) types.BillingModeSummary {
	if v1 == nil {
		return nil
	}
	return &types.BillingModeSummary{
		LastUpdateToPayPerRequestDateTime: v1.LastUpdateToPayPerRequestDateTime,
		LastUpdateToProvisionedDateTime:   v1.LastUpdateToProvisionedDateTime,
		Period:                            v1.Period,
		State:                             types.BillingMode(*v1.State),
	}
}

// convertV2ToV1BillingModeSummary converts v2 BillingModeSummary to v1 BillingModeSummary
func convertV2ToV1BillingModeSummary(v2 types.BillingModeSummary) *v1dynamodb.BillingModeSummary {
	if v2 == nil {
		return nil
	}
	return &v1dynamodb.BillingModeSummary{
		LastUpdateToPayPerRequestDateTime: v2.LastUpdateToPayPerRequestDateTime,
		LastUpdateToProvisionedDateTime:   v2.LastUpdateToProvisionedDateTime,
		Period:                            v2.Period,
		State:                             v1aws.String(string(v2.State)),
	}
}

// convertV1ToV2KeySchema converts v1 KeySchema to v2 KeySchema
func convertV1ToV2KeySchema(v1 []*v1dynamodb.KeySchemaElement) []types.KeySchemaElement {
	if v1 == nil {
		return nil
	}
	v2 := make([]types.KeySchemaElement, len(v1))
	for i, key := range v1 {
		v2[i] = types.KeySchemaElement{
			AttributeName: key.AttributeName,
			KeyType:       types.KeyType(*key.KeyType),
		}
	}
	return v2
}

// convertV2ToV1KeySchema converts v2 KeySchema to v1 KeySchema
func convertV2ToV1KeySchema(v2 []types.KeySchemaElement) []*v1dynamodb.KeySchemaElement {
	if v2 == nil {
		return nil
	}
	v1 := make([]*v1dynamodb.KeySchemaElement, len(v2))
	for i, key := range v2 {
		v1[i] = &v1dynamodb.KeySchemaElement{
			AttributeName: key.AttributeName,
			KeyType:       v1aws.String(string(key.KeyType)),
		}
	}
	return v1
}

// convertV1ToV2LocalSecondaryIndexes converts v1 LocalSecondaryIndexes to v2 LocalSecondaryIndexes
func convertV1ToV2LocalSecondaryIndexes(v1 []*v1dynamodb.LocalSecondaryIndex) []types.LocalSecondaryIndex {
	if v1 == nil {
		return nil
	}
	v2 := make([]types.LocalSecondaryIndex, len(v1))
	for i, lsi := range v1 {
		v2[i] = types.LocalSecondaryIndex{
			IndexArn:       lsi.IndexArn,
			IndexName:      lsi.IndexName,
			IndexSizeBytes: lsi.IndexSizeBytes,
			ItemCount:      lsi.ItemCount,
			KeySchema:      convertV1ToV2KeySchema(lsi.KeySchema),
			Projection:     convertV1ToV2Projection(lsi.Projection),
		}
	}
	return v2
}

// convertV2ToV1LocalSecondaryIndexes converts v2 LocalSecondaryIndexes to v1 LocalSecondaryIndexes
func convertV2ToV1LocalSecondaryIndexes(v2 []types.LocalSecondaryIndex) []*v1dynamodb.LocalSecondaryIndex {
	if v2 == nil {
		return nil
	}
	v1 := make([]*v1dynamodb.LocalSecondaryIndex, len(v2))
	for i, lsi := range v2 {
		v1[i] = &v1dynamodb.LocalSecondaryIndex{
			IndexArn:       lsi.IndexArn,
			IndexName:      lsi.IndexName,
			IndexSizeBytes: lsi.IndexSizeBytes,
			ItemCount:      lsi.ItemCount,
			KeySchema:      convertV2ToV1KeySchema(lsi.KeySchema),
			Projection:     convertV2ToV1Projection(lsi.Projection),
		}
	}
	return v1
}

// convertV1ToV2Projection converts v1 Projection to v2 Projection
func convertV1ToV2Projection(v1 *v1dynamodb.Projection) types.Projection {
	if v1 == nil {
		return nil
	}
	return &types.Projection{
		NonKeyAttributes: v1.NonKeyAttributes,
		ProjectionType:   types.ProjectionType(*v1.ProjectionType),
	}
}

// convertV2ToV1Projection converts v2 Projection to v1 Projection
func convertV2ToV1Projection(v2 types.Projection) *v1dynamodb.Projection {
	if v2 == nil {
		return nil
	}
	return &v1dynamodb.Projection{
		NonKeyAttributes: v2.NonKeyAttributes,
		ProjectionType:   v1aws.String(string(v2.ProjectionType)),
	}
}

// convertV1ToV2ProvisionedThroughput converts v1 ProvisionedThroughput to v2 ProvisionedThroughput
func convertV1ToV2ProvisionedThroughput(v1 *v1dynamodb.ProvisionedThroughput) types.ProvisionedThroughput {
	if v1 == nil {
		return nil
	}
	return &types.ProvisionedThroughput{
		ReadCapacityUnits:  v1.ReadCapacityUnits,
		WriteCapacityUnits: v1.WriteCapacityUnits,
	}
}

// convertV2ToV1ProvisionedThroughput converts v2 ProvisionedThroughput to v1 ProvisionedThroughput
func convertV2ToV1ProvisionedThroughput(v2 types.ProvisionedThroughput) *v1dynamodb.ProvisionedThroughput {
	if v2 == nil {
		return nil
	}
	return &v1dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  v2.ReadCapacityUnits,
		WriteCapacityUnits: v2.WriteCapacityUnits,
	}
}

// convertV1ToV2Replicas converts v1 Replicas to v2 Replicas
func convertV1ToV2Replicas(v1 []*v1dynamodb.ReplicaDescription) []types.ReplicaDescription {
	if v1 == nil {
		return nil
	}
	v2 := make([]types.ReplicaDescription, len(v1))
	for i, replica := range v1 {
		v2[i] = types.ReplicaDescription{
			RegionName: replica.RegionName,
		}
	}
	return v2
}

// convertV2ToV1Replicas converts v2 Replicas to v1 Replicas
func convertV2ToV1Replicas(v2 []types.ReplicaDescription) []*v1dynamodb.ReplicaDescription {
	if v2 == nil {
		return nil
	}
	v1 := make([]*v1dynamodb.ReplicaDescription, len(v2))
	for i, replica := range v2 {
		v1[i] = &v1dynamodb.ReplicaDescription{
			RegionName: replica.RegionName,
		}
	}
	return v1
}

// convertV1ToV2RestoreSummary converts v1 RestoreSummary to v2 RestoreSummary
func convertV1ToV2RestoreSummary(v1 *v1dynamodb.RestoreSummary) types.RestoreSummary {
	if v1 == nil {
		return nil
	}
	return &types.RestoreSummary{
		RestoreDateTime:   v1.RestoreDateTime,
		RestoreInProgress: v1.RestoreInProgress,
		SourceBackupArn:   v1.SourceBackupArn,
		SourceTableArn:    v1.SourceTableArn,
	}
}

// convertV2ToV1RestoreSummary converts v2 RestoreSummary to v1 RestoreSummary
func convertV2ToV1RestoreSummary(v2 types.RestoreSummary) *v1dynamodb.RestoreSummary {
	if v2 == nil {
		return nil
	}
	return &v1dynamodb.RestoreSummary{
		RestoreDateTime:   v2.RestoreDateTime,
		RestoreInProgress: v2.RestoreInProgress,
		SourceBackupArn:   v2.SourceBackupArn,
		SourceTableArn:    v2.SourceTableArn,
	}
}

// convertV1ToV2SSEDescription converts v1 SSEDescription to v2 SSEDescription
func convertV1ToV2SSEDescription(v1 *v1dynamodb.SSEDescription) types.SSEDescription {
	if v1 == nil {
		return nil
	}
	return &types.SSEDescription{
		KMSMasterKeyArn: v1.KMSMasterKeyArn,
		Status:          types.SSEStatus(*v1.Status),
	}
}

// convertV2ToV1SSEDescription converts v2 SSEDescription to v1 SSEDescription
func convertV2ToV1SSEDescription(v2 types.SSEDescription) *v1dynamodb.SSEDescription {
	if v2 == nil {
		return nil
	}
	return &v1dynamodb.SSEDescription{
		KMSMasterKeyArn: v2.KMSMasterKeyArn,
		Status:          v1aws.String(string(v2.Status)),
	}
}

// convertV1ToV2StreamSpecification converts v1 StreamSpecification to v2 StreamSpecification
func convertV1ToV2StreamSpecification(v1 *v1dynamodb.StreamSpecification) types.StreamSpecification {
	if v1 == nil {
		return nil
	}
	return &types.StreamSpecification{
		StreamEnabled:  v1.StreamEnabled,
		StreamViewType: types.StreamViewType(*v1.StreamViewType),
	}
}

// convertV2ToV1StreamSpecification converts v2 StreamSpecification to v1 StreamSpecification
func convertV2ToV1StreamSpecification(v2 types.StreamSpecification) *v1dynamodb.StreamSpecification {
	if v2 == nil {
		return nil
	}
	return &v1dynamodb.StreamSpecification{
		StreamEnabled:  v2.StreamEnabled,
		StreamViewType: v1aws.String(string(v2.StreamViewType)),
	}
}

// convertV1ToV2AttributeValues converts a slice of v1 AttributeValues to v2 AttributeValues
func convertV1ToV2AttributeValues(v1Values []*v1dynamodb.AttributeValue) []types.AttributeValue {
	if v1Values == nil {
		return nil
	}
	v2Values := make([]types.AttributeValue, len(v1Values))
	for i, v := range v1Values {
		v2Values[i] = convertV1ToV2AttributeValue(v)
	}
	return v2Values
}

// convertV2ToV1AttributeValues converts a map of v2 AttributeValues to v1 AttributeValues
func convertV2ToV1AttributeValues(v2Values map[string]types.AttributeValue) map[string]*v1dynamodb.AttributeValue {
	if v2Values == nil {
		return nil
	}
	v1Values := make(map[string]*v1dynamodb.AttributeValue, len(v2Values))
	for k, v := range v2Values {
		v1Values[k] = convertV2ToV1AttributeValue(v)
	}
	return v1Values
}

// convertV1ToV2ProvisionedThroughputDescription converts v1 ProvisionedThroughputDescription to v2 ProvisionedThroughputDescription
func convertV1ToV2ProvisionedThroughputDescription(v1 *v1dynamodb.ProvisionedThroughputDescription) *types.ProvisionedThroughputDescription {
	if v1 == nil {
		return nil
	}
	return &types.ProvisionedThroughputDescription{
		LastDecreaseDateTime:   v1.LastDecreaseDateTime,
		LastIncreaseDateTime:   v1.LastIncreaseDateTime,
		NumberOfDecreasesToday: v1.NumberOfDecreasesToday,
		ReadCapacityUnits:      v1.ReadCapacityUnits,
		WriteCapacityUnits:     v1.WriteCapacityUnits,
	}
}

// convertV2ToV1ProvisionedThroughputDescription converts v2 ProvisionedThroughputDescription to v1 ProvisionedThroughputDescription
func convertV2ToV1ProvisionedThroughputDescription(v2 *types.ProvisionedThroughputDescription) *v1dynamodb.ProvisionedThroughputDescription {
	if v2 == nil {
		return nil
	}
	return &v1dynamodb.ProvisionedThroughputDescription{
		LastDecreaseDateTime:   v2.LastDecreaseDateTime,
		LastIncreaseDateTime:   v2.LastIncreaseDateTime,
		NumberOfDecreasesToday: v2.NumberOfDecreasesToday,
		ReadCapacityUnits:      v2.ReadCapacityUnits,
		WriteCapacityUnits:     v2.WriteCapacityUnits,
	}
}
