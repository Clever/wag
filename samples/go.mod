module github.com/Clever/wag/samples/v9

go 1.16

require (
	github.com/Clever/discovery-go v1.8.1
	github.com/Clever/go-process-metrics v0.4.0
	github.com/Clever/launch-gen v0.0.0-20210816170657-06a4d01b3706
	github.com/cespare/reflex v0.3.1
	github.com/davecgh/go-spew v1.1.1
	github.com/get-woke/woke v0.17.1
	github.com/go-errors/errors v1.1.1
	github.com/go-openapi/strfmt v0.21.2
	github.com/go-openapi/swag v0.21.1
	github.com/golang/mock v1.6.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	go.mongodb.org/mongo-driver v1.10.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.34.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.34.0 // indirect
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/jaeger v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.10.0
	go.opentelemetry.io/otel/sdk v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f

)

replace github.com/go-openapi/validate => github.com/go-openapi/validate v0.0.0-20180703152151-9a6e517cddf1 // pre-modules tag 0.15.0x

replace github.com/Clever/wag/samples/gen-go-basic/models/v9 => ./gen-go-basic/models

replace github.com/Clever/wag/samples/gen-go-blog/models/v9 => ./gen-go-blog/models

replace github.com/Clever/wag/samples/gen-go-client-only/models/v9 => ./gen-go-client-only/models

replace github.com/Clever/wag/samples/gen-go-db/models/v9 => ./gen-go-db/models

replace github.com/Clever/wag/samples/gen-go-db-custom-path/models/v9 => ./gen-go-db-custom-path/models

replace github.com/Clever/wag/samples/gen-go-db-only/models/v9 => ./gen-go-db-only/models

replace github.com/Clever/wag/samples/gen-go-deprecated/models/v9 => ./gen-go-deprecated/models

replace github.com/Clever/wag/samples/gen-go-errors/models/v9 => ./gen-go-errors/models

replace github.com/Clever/wag/samples/gen-go-nils/models/v9 => ./gen-go-nils/models

replace github.com/Clever/wag/samples/gen-go-basic/client/v9 => ./gen-go-basic/client

replace github.com/Clever/wag/samples/gen-go-blog/client/v9 => ./gen-go-blog/client

replace github.com/Clever/wag/samples/gen-go-client-only/client/v9 => ./gen-go-client-only/models

replace github.com/Clever/wag/samples/gen-go-db/client/v9 => ./gen-go-db/client

replace github.com/Clever/wag/samples/gen-go-db-custom-path/client/v9 => ./gen-go-db-custom-path/client

replace github.com/Clever/wag/samples/gen-go-db-only/client/v9 => ./gen-go-db-only/client

replace github.com/Clever/wag/samples/gen-go-deprecated/client/v9 => ./gen-go-deprecated/client

replace github.com/Clever/wag/samples/gen-go-errors/client/v9 => ./gen-go-errors/client

replace github.com/Clever/wag/samples/gen-go-nils/client/v9 => ./gen-go-nils/client

replace github.com/go-openapi/swag => github.com/go-openapi/swag v0.0.0-20180102232305-84f4bee7c0a6

replace github.com/Clever/wag/v9 => ../
