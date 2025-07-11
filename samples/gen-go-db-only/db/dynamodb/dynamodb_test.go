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

	"github.com/Clever/wag/samples/v9/gen-go-db-only/db"
	"github.com/Clever/wag/samples/v9/gen-go-db-only/db/tests"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

/* FAILING_TEST */
func TestDynamoDBStore(t *testing.T) {
	t.Skip()
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
				t.Errorf(`zombie ddb local process running. Kill it and try again: pgrep -f "java -jar /tmp/DynamoDBLocal.jar" | xargs kill`)
				return
			}
		}
	}()

	// the ddblocal command should not exit with an error before the test is finished
	go func() {
		if err := cmd.Wait(); err != nil && testCtx.Err() == nil {
			t.Errorf("cmd.Wait: %s", err)
			return
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

	cfg, err := config.LoadDefaultConfig(testCtx)
	if err != nil {
		t.Fatal(err)
	}
	dynamoDBAPI := dynamodb.NewFromConfig(cfg)

	tests.RunDBTests(t, func() db.Interface {
		prefix := "automated-testing"
		listTablesOutput, err := dynamoDBAPI.ListTables(testCtx, &dynamodb.ListTablesInput{})
		if err != nil {
			t.Fatal(err)
		}
		for _, tableName := range listTablesOutput.TableNames {
			if strings.HasPrefix(tableName, prefix) {
				dynamoDBAPI.DeleteTable(testCtx, &dynamodb.DeleteTableInput{
					TableName: aws.String(tableName),
				})
			}
		}
		d, err := New(Config{
			DynamoDBAPI:               dynamoDBAPI,
			DefaultPrefix:             prefix,
			DefaultReadCapacityUnits:  10,
			DefaultWriteCapacityUnits: 10,
			DeploymentTable: DeploymentTable{
				TableName: "automated-testing-Deployment",
			},
			EventTable: EventTable{
				TableName: "automated-testing-Event",
			},
			NoRangeThingWithCompositeAttributesTable: NoRangeThingWithCompositeAttributesTable{
				TableName: "automated-testing-NoRangeThingWithCompositeAttributes",
			},
			SimpleThingTable: SimpleThingTable{
				TableName: "automated-testing-SimpleThing",
			},
			TeacherSharingRuleTable: TeacherSharingRuleTable{
				TableName: "automated-testing-TeacherSharingRule",
			},
			ThingTable: ThingTable{
				TableName: "automated-testing-Thing",
			},
			ThingAllowingBatchWritesTable: ThingAllowingBatchWritesTable{
				TableName: "automated-testing-ThingAllowingBatchWrites",
			},
			ThingAllowingBatchWritesWithCompositeAttributesTable: ThingAllowingBatchWritesWithCompositeAttributesTable{
				TableName: "automated-testing-ThingAllowingBatchWritesWithCompositeAttributes",
			},
			ThingWithAdditionalAttributesTable: ThingWithAdditionalAttributesTable{
				TableName: "automated-testing-ThingWithAdditionalAttributes",
			},
			ThingWithCompositeAttributesTable: ThingWithCompositeAttributesTable{
				TableName: "automated-testing-ThingWithCompositeAttributes",
			},
			ThingWithCompositeEnumAttributesTable: ThingWithCompositeEnumAttributesTable{
				TableName: "automated-testing-ThingWithCompositeEnumAttributes",
			},
			ThingWithDateGSITable: ThingWithDateGSITable{
				TableName: "automated-testing-ThingWithDateGSI",
			},
			ThingWithDateRangeTable: ThingWithDateRangeTable{
				TableName: "automated-testing-ThingWithDateRange",
			},
			ThingWithDateRangeKeyTable: ThingWithDateRangeKeyTable{
				TableName: "automated-testing-ThingWithDateRangeKey",
			},
			ThingWithDateTimeCompositeTable: ThingWithDateTimeCompositeTable{
				TableName: "automated-testing-ThingWithDateTimeComposite",
			},
			ThingWithDatetimeGSITable: ThingWithDatetimeGSITable{
				TableName: "automated-testing-ThingWithDatetimeGSI",
			},
			ThingWithEnumHashKeyTable: ThingWithEnumHashKeyTable{
				TableName: "automated-testing-ThingWithEnumHashKey",
			},
			ThingWithMatchingKeysTable: ThingWithMatchingKeysTable{
				TableName: "automated-testing-ThingWithMatchingKeys",
			},
			ThingWithMultiUseCompositeAttributeTable: ThingWithMultiUseCompositeAttributeTable{
				TableName: "automated-testing-ThingWithMultiUseCompositeAttribute",
			},
			ThingWithRequiredCompositePropertiesAndKeysOnlyTable: ThingWithRequiredCompositePropertiesAndKeysOnlyTable{
				TableName: "automated-testing-ThingWithRequiredCompositePropertiesAndKeysOnly",
			},
			ThingWithRequiredFieldsTable: ThingWithRequiredFieldsTable{
				TableName: "automated-testing-ThingWithRequiredFields",
			},
			ThingWithRequiredFields2Table: ThingWithRequiredFields2Table{
				TableName: "automated-testing-ThingWithRequiredFields2",
			},
			ThingWithTransactMultipleGSITable: ThingWithTransactMultipleGSITable{
				TableName: "automated-testing-ThingWithTransactMultipleGSI",
			},
			ThingWithTransactionTable: ThingWithTransactionTable{
				TableName: "automated-testing-ThingWithTransaction",
			},
			ThingWithTransactionWithSimpleThingTable: ThingWithTransactionWithSimpleThingTable{
				TableName: "automated-testing-ThingWithTransactionWithSimpleThing",
			},
			ThingWithUnderscoresTable: ThingWithUnderscoresTable{
				TableName: "automated-testing-ThingWithUnderscores",
			},
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
