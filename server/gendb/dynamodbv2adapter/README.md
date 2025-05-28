# DynamoDB v2 Adapter

This package provides an adapter that implements the AWS SDK v1 DynamoDB interface (`dynamodbiface.DynamoDBAPI`) using an AWS SDK v2 DynamoDB client. This allows code that expects the v1 interface to work with the v2 client.

## Usage

```go
import (
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/clever/wag/gendb/dynamodbv2adapter"
)

func main() {
    // Create a v2 client
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        panic(err)
    }
    v2Client := dynamodb.NewFromConfig(cfg)

    // Create the adapter
    adapter := dynamodbv2adapter.New(v2Client)

    // Use the adapter where a v1 client is expected
    config := dynamodbgen.Config{
        DynamoDBAPI: adapter,
        DefaultPrefix: "my-prefix",
        // ... other config options
    }
    db, err := dynamodbgen.New(config)
    if err != nil {
        panic(err)
    }
}
```

## Features

- Implements all methods from the v1 `dynamodbiface.DynamoDBAPI` interface
- Handles conversion between v1 and v2 input/output types
- Maintains context handling from v2 client
- Preserves all DynamoDB functionality while using the v2 client

## Requirements

- Go 1.21 or later
- AWS SDK v2 for Go
- AWS SDK v1 for Go (for interface definitions)

## License

MIT 