package dynamodb

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/Clever/wag/samples/gen-go-db/server/db"
	"github.com/Clever/wag/samples/gen-go-db/server/db/tests"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestDynamoDBStore(t *testing.T) {
	// spin up dynamodb local
	testCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(testCtx, "./dynamodb-local.sh")
	ddbLocalOutputReader, ddbLocalOutputWriter := io.Pipe()
	cmd.Stdout = io.MultiWriter(os.Stdout, ddbLocalOutputWriter)
	cmd.Stderr = io.MultiWriter(os.Stderr, ddbLocalOutputWriter)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// wait for dynamodb local to output the correct startup log
	scanner := bufio.NewScanner(ddbLocalOutputReader)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Initializing DynamoDB Local with the following configuration") {
			break
		}
	}

	dynamoDBAPI := dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String("doesntmatter"),
			Endpoint:    aws.String("http://localhost:8002" /* default dynamodb-local port */),
			Credentials: credentials.NewStaticCredentials("id", "secret", "token"),
		},
	})))

	tests.RunDBTests(t, func() db.Interface {
		prefix := "automated-testing"
		listTablesOutput, err := dynamoDBAPI.ListTablesWithContext(testCtx, &dynamodb.ListTablesInput{})
		if err != nil {
			t.Fatal(err)
		}
		for _, tableName := range listTablesOutput.TableNames {
			if strings.HasPrefix(*tableName, prefix) {
				dynamoDBAPI.DeleteTableWithContext(testCtx, &dynamodb.DeleteTableInput{
					TableName: tableName,
				})
			}
		}
		d, err := New(Config{
			DynamoDBAPI:               dynamoDBAPI,
			DefaultPrefix:             prefix,
			DefaultReadCapacityUnits:  10,
			DefaultWriteCapacityUnits: 10,
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := d.CreateTables(testCtx); err != nil {
			t.Fatal(err)
		}
		return d
	})
}
