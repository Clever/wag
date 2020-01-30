package dynamodb

import (
	"bufio"
	"context"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
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
	// make sure nothing is running 8002
	if _, err := net.DialTimeout("tcp", "localhost:8002", 100*time.Millisecond); err == nil {
		t.Fatal(`zombie ddb local process running. Kill it and try again: pgrep -f "java -jar /tmp/DynamoDBLocal.jar" | xargs kill`)
	}

	// spin up dynamodb local, making sure to kill it when
	// - the test function is finished (defer cancel())
	// - the test is sigkill'd
	testCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Signal(syscall.SIGTERM))
	go func() {
		for range c {
			t.Logf("ctrl-c received")
			cancel()
		}
	}()
	cmd := exec.CommandContext(testCtx, "./dynamodb-local.sh")
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatalf("cmd.Start: %v", err)
	}

	// relay stdout and stderr of the ddblocal process, and
	// check for signs that it didn't start up correctly
	outLines := bufio.NewScanner(stdout)
	errLines := bufio.NewScanner(stderr)
	go func() {
		for outLines.Scan() {
			t.Logf("ddblocal stdout: %s", outLines.Text())
		}
	}()
	go func() {
		for errLines.Scan() {
			txt := errLines.Text()
			t.Logf("ddblocal stderr: %s", txt)
			if txt == "java.net.BindException: Address already in use" {
				t.Fatal(`zombie ddb local process running. Kill it and try again: pgrep -f "java -jar /tmp/DynamoDBLocal.jar" | xargs kill`)
			}
		}
	}()

	// the ddblocal command should not exit with an error
	go func() {
		if err := cmd.Wait(); err != nil {
			t.Fatalf("cmd.Wait: %s", err)
		}
	}()

	// loop for 10s trying to establish a connection
	connected := false
	for start := time.Now(); start.Before(start.Add(10 * time.Second)); time.Sleep(1 * time.Second) {
		if c, err := net.DialTimeout("tcp", "localhost:8002", 100*time.Millisecond); err == nil {
			c.Close()
			connected = true
			break
		} else {
			t.Logf("could not connect to ddb local, will retry: %s", err)
		}
	}
	if connected == false {
		t.Fatal("failed to connect within 60 seconds")
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
