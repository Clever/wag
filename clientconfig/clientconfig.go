package clientconfig

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Clever/kayvee-go/v7/logger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func ClientLogger(wagAppName string) *logger.Logger {
	return logger.NewConcreteLogger(fmt.Sprintf("%s-wagclient", wagAppName))
}

func WithoutTracing(wagAppName string) (*logger.Logger, *http.RoundTripper) {
	return ClientLogger(wagAppName), &http.DefaultTransport
}

func WithTracing(ctx context.Context, wagAppName string, exporter sdktrace.SpanExporter) (*logger.Logger, *http.RoundTripper) {
	samplingProbability := determineSampling()
	baseTransport := http.DefaultTransport
	tp := newTracerProvider(exporter, samplingProbability, wagAppName)

	instrumentedTransport := DefaultInstrumentor(ctx, baseTransport, *tp)

	return ClientLogger(wagAppName), &instrumentedTransport
}
