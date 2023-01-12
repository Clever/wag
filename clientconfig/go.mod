module github.com/Clever/wag/clientconfig

go 1.16

require (
	github.com/Clever/kayvee-go/v7 v7.7.0
	github.com/Clever/wag/logging/wagclientlogger v0.0.0-20230110184825-edb52117e67a // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.35.0
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.10.0
	go.opentelemetry.io/otel/sdk v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
)
