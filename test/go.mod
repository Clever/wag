module github.com/Clever/wag/test/v8

go 1.16

require (
	github.com/Clever/wag/samples/v8 v8.0.0
	github.com/Clever/wag/v8 v8.0.0-00010101000000-000000000000
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/go-openapi/swag v0.19.14
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
	gopkg.in/Clever/kayvee-go.v6 v6.27.0
)

replace github.com/Clever/wag/samples/v8 => ../samples/

replace github.com/Clever/wag/v8 => ../
