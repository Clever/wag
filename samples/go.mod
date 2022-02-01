module github.com/Clever/wag/samples/v8

go 1.16

require (
	github.com/Clever/discovery-go v1.7.3
	github.com/Clever/go-process-metrics v0.4.0
	github.com/Clever/kayvee-go/v7 v7.4.0
	github.com/Clever/wag/v8 v8.0.0-00010101000000-000000000000
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/aws/aws-sdk-go v1.42.44
	github.com/donovanhide/eventsource v0.0.0-20171031113327-3ed64d21fb0b
	github.com/go-errors/errors v1.1.1
	github.com/go-openapi/errors v0.19.9
	github.com/go-openapi/strfmt v0.19.11
	github.com/go-openapi/swag v0.21.1
	github.com/go-openapi/validate v0.19.15
	github.com/golang/mock v1.6.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/stretchr/testify v1.7.0
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.28.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.28.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.3.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.opentelemetry.io/otel/trace v1.3.0
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace github.com/go-openapi/swag => github.com/go-openapi/swag v0.0.0-20180102232305-84f4bee7c0a6

replace github.com/Clever/wag/v8 => ../
