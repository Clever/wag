package test

import (
	"context"
	"testing"

	clientconfig "github.com/Clever/wag/clientconfig/v9"
	"github.com/Clever/wag/samples/gen-go-basic/client/v9"
	"github.com/Clever/wag/samples/gen-go-basic/models/v9"
	"github.com/Clever/wag/samples/v9/gen-go-basic/servertracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

func TestOpenTelemetryInstrumentation(t *testing.T) {
	// client should
	// 1. generate a span
	// 2. send the w3c-approved "Traceparent" header
	// server should
	// 1. generate a span
	// 2. log trace and span ids
	ctx := context.Background()
	exporter, provider, err := servertracing.SetupGlobalTraceProviderAndExporterForTest()
	if err != nil {
		t.Fatal(err)
	}
	defer exporter.Shutdown(ctx)
	defer provider.Shutdown(ctx)
	s, _ := setupServer()
	defer s.Close()
	logger, transport := clientconfig.WithTracing("service", exporter)

	c := client.New(s.URL, logger, transport)
	c.HealthCheck(ctx)
	spans := exporter.GetSpans()
	require.Equal(t, 2, len(spans))
	serverSpan := spans[0]
	clientSpan := spans[1]
	assert.Equal(t, "/v1/health/check", serverSpan.Name)
	assert.Equal(t, true, serverSpan.Parent.IsRemote())
	assert.Equal(t, clientSpan.SpanContext.SpanID(), serverSpan.Parent.SpanID())
	assert.Equal(t, "healthCheck", clientSpan.Name)
	assert.Equal(t, false, clientSpan.Parent.IsRemote())

	// test that the client joins a pre-existing span in the ctx
	parentSpanID := serverSpan.SpanContext.SpanID()
	ctx = trace.ContextWithRemoteSpanContext(ctx, serverSpan.SpanContext)
	c.GetBookByID(ctx, &models.GetBookByIDInput{BookID: 1})
	spans = exporter.GetSpans()
	require.Equal(t, 4, len(spans))
	serverSpan = spans[2]
	clientSpan = spans[3]
	assert.Equal(t, "/v1/books/{book_id}", serverSpan.Name)
	assert.Equal(t, "getBookByID", clientSpan.Name)
	assert.Equal(t, true, clientSpan.Parent.IsRemote())
	assert.Equal(t, parentSpanID, clientSpan.Parent.SpanID())
}

func hasAttribute(r *resource.Resource, attr string) bool {
	for _, a := range r.Attributes() {
		if string(a.Key) == attr {
			return true
		}
	}
	return false
}
