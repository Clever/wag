package test

import (
	"context"
	"testing"

	"github.com/Clever/wag/v5/samples/gen-go/client"
	"github.com/Clever/wag/v5/samples/gen-go/models"
	"github.com/Clever/wag/v5/samples/gen-go/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

func TestOpenTelemetryInstrumentation(t *testing.T) {
	// client should
	// 1. generate a span
	// 2. send
	//    - the w3c-approved "Traceparent" header
	//    - the aws-xray "X-Amzn-Trace-Id" header
	// server should
	// 1. generate a span
	// 2. log trace and span ids
	ctx := context.Background()
	s, _ := setupServer()
	defer s.Close()
	exporter, err := tracing.SetupGlobalTraceProviderAndExporterForTest()
	if err != nil {
		t.Fatal(err)
	}
	defer exporter.Shutdown(ctx)
	c := client.New(s.URL)
	c.HealthCheck(ctx)
	spans := exporter.GetSpans()
	require.Equal(t, 2, len(spans))
	serverSpan := spans[0]
	clientSpan := spans[1]
	assert.Equal(t, "/v1/health/check", serverSpan.Name)
	assert.Equal(t, true, serverSpan.HasRemoteParent)
	assert.Equal(t, clientSpan.SpanContext.SpanID, serverSpan.ParentSpanID)
	assert.Equal(t, "healthCheck", clientSpan.Name)
	assert.Equal(t, false, clientSpan.HasRemoteParent)

	// check it sets resource to server info
	assert.True(t, hasAttribute(serverSpan.Resource, "deploy_env"))
	assert.True(t, hasAttribute(serverSpan.Resource, "app_name"))

	// test that the client joins a pre-existing span in the ctx
	parentSpanID := serverSpan.SpanContext.SpanID
	ctx = trace.ContextWithRemoteSpanContext(ctx, serverSpan.SpanContext)
	c.GetBookByID(ctx, &models.GetBookByIDInput{BookID: 1})
	spans = exporter.GetSpans()
	require.Equal(t, 4, len(spans))
	serverSpan = spans[2]
	clientSpan = spans[3]
	assert.Equal(t, "/v1/books/{book_id}", serverSpan.Name)
	assert.Equal(t, "getBookByID", clientSpan.Name)
	assert.Equal(t, true, clientSpan.HasRemoteParent)
	assert.Equal(t, parentSpanID, clientSpan.ParentSpanID)
}

func hasAttribute(r *resource.Resource, attr string) bool {
	for _, a := range r.Attributes() {
		if string(a.Key) == attr {
			return true
		}
	}
	return false
}
