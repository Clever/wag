package dynamodb

import (
	"context"
	"net"
	"os/exec"
	"strings"
	"testing"
	"time"

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

	// restart test db if it gets killed before tests are completed.
	testsDone := make(chan struct{})
	defer close(testsDone)
	go func(doneC chan struct{}) {
		cmd := exec.CommandContext(testCtx, "./dynamodb-local.sh")
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		for {
			select {
			case <-doneC:
				return
			default:
				if err := cmd.Wait(); err != nil {
					cmd = exec.CommandContext(testCtx, "./dynamodb-local.sh")
					cmd.Start()
				}
			}
		}
	}(testsDone)

	var err error
	end := time.Now().Add(60 * time.Second)
	for time.Now().Before(end) {
		var c net.Conn
		c, err = net.Dial("tcp", "localhost:8002")
		if err == nil {
			c.Close()
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		t.Fatal(err)
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
