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

// WithoutTracing is now Deprecated and so is WithTracing, but both are still provided for backwards compatibility.
// Use Default instead.
func WithoutTracing(wagAppName string) (*logger.Logger, *http.RoundTripper) {
	return Default(wagAppName)
}

// Default returns a logger and a transport to use in client requests.
// It is meant as a convenience function for initiating Wag clients
func Default(wagAppName string) (logger.Logger, *http.RoundTripper) {
	baseTransport := http.DefaultTransport
	instrumentedTransport := DefaultInstrumentor(baseTransport, wagAppName)
	return ClientLogger(wagAppName), &instrumentedTransport
}

// WithTracing is now Deprecated and so is WithoutTracing, but both are still provided for backwards compatibility.
// Use Default instad.
func WithTracing(wagAppName string, exporter sdktrace.SpanExporter) (*logger.Logger, *http.RoundTripper) {
	return Default(wagAppName)
}
