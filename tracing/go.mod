module github.com/Clever/wag/tracing

go 1.16

require (
	github.com/Clever/kayvee-go/v7 v7.5.1-0.20220809022937-ea1e87b3d8d5
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.34.0
	go.opentelemetry.io/otel v1.9.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.8.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.8.0
	go.opentelemetry.io/otel/sdk v1.9.0
	go.opentelemetry.io/otel/trace v1.9.0
)
