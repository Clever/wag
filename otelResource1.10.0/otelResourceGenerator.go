package otelResource

import (
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

// NewResource generates a new otel resource with the service name and version passed into it.
// Needed because semconv is part of the whole otel package.
func NewResource(service, version string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceVersionKey.String(version),
		),
	)
	return r
}
