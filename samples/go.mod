module github.com/Clever/wag/v7/samples

go 1.13

require (
	github.com/Clever/discovery-go v1.7.2
	github.com/Clever/go-process-metrics v0.2.0
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/aws/aws-sdk-go v1.38.25
	github.com/donovanhide/eventsource v0.0.0-20171031113327-3ed64d21fb0b
	github.com/go-errors/errors v1.1.1
	github.com/go-openapi/analysis v0.0.0-20180126163718-f59a71f0ece6 // indirect
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.0.0-20161105162150-36d33bfe519e // indirect
	github.com/go-openapi/loads v0.0.0-20171207192234-2a2b323bab96 // indirect
	github.com/go-openapi/runtime v0.0.0-20180131174916-09fac855d850 // indirect
	github.com/go-openapi/spec v0.0.0-20180213232550-1de3e0542de6 // indirect
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.14
	github.com/go-openapi/validate v0.0.0-20180222165948-180bba53b988
	github.com/golang/mock v1.5.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.15.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.15.1
	go.opentelemetry.io/contrib/propagators/aws v0.15.1
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/exporters/otlp v0.15.0
	go.opentelemetry.io/otel/sdk v0.15.0
	golang.org/x/net v0.0.0-20210423184538-5f58ad60dda6
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/Clever/kayvee-go.v6 v6.24.1
)

replace github.com/go-openapi/swag => github.com/go-openapi/swag v0.0.0-20180102232305-84f4bee7c0a6
