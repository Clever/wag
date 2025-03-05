module github.com/Clever/wag/samples/v9

go 1.21

toolchain go1.21.0

require (
	github.com/Clever/go-process-metrics v0.4.0
	github.com/Clever/kayvee-go/v7 v7.7.0
	github.com/Clever/wag/logging/wagclientlogger v0.0.0-20230110184825-edb52117e67a
	github.com/Clever/wag/samples/gen-go-basic/client/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/samples/gen-go-basic/models/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/samples/gen-go-blog/models/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/samples/gen-go-db-custom-path/models/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/samples/gen-go-db-only/models/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/samples/gen-go-db/models/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/samples/gen-go-deprecated/models/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/samples/gen-go-errors/client/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/samples/gen-go-errors/models/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/samples/gen-go-nils/client/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/samples/gen-go-nils/models/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/samples/gen-go-strings/models/v9 v9.0.0-00010101000000-000000000000
	github.com/Clever/wag/v9 v9.0.0-00010101000000-000000000000
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/aws/aws-sdk-go v1.44.89
	github.com/go-errors/errors v1.1.1
	github.com/go-openapi/strfmt v0.21.3
	github.com/go-openapi/swag v0.22.3
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.1.2
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/stretchr/testify v1.8.2
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.34.0
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.10.0
	go.opentelemetry.io/otel/sdk v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f
)

require (
	github.com/Clever/discovery-go v1.8.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.20.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/runtime v0.24.1 // indirect
	github.com/go-openapi/spec v0.20.7 // indirect
	github.com/go-openapi/validate v0.22.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/smartystreets/goconvey v1.8.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	go.mongodb.org/mongo-driver v1.10.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.31.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.31.0 // indirect
	go.opentelemetry.io/otel/metric v0.31.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.31.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	golang.org/x/net v0.0.0-20220826154423-83b083e8dc8b // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220829175752-36a9c930ecbf // indirect
	google.golang.org/grpc v1.49.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-openapi/validate => github.com/go-openapi/validate v0.0.0-20180703152151-9a6e517cddf1 // pre-modules tag 0.15.0x

replace github.com/go-openapi/errors => github.com/go-openapi/errors v0.0.0-20180515155515-b2b2befaf267 // pre-modules tag 0.15.0x

replace github.com/Clever/wag/samples/gen-go-strings/models/v9 => ./gen-go-strings/models

replace github.com/Clever/wag/samples/gen-go-basic/models/v9 => ./gen-go-basic/models

replace github.com/Clever/wag/samples/gen-go-blog/models/v9 => ./gen-go-blog/models

replace github.com/Clever/wag/samples/gen-go-client-only/models/v9 => ./gen-go-client-only/models

replace github.com/Clever/wag/samples/gen-go-db/models/v9 => ./gen-go-db/models

replace github.com/Clever/wag/samples/gen-go-db-custom-path/models/v9 => ./gen-go-db-custom-path/models

replace github.com/Clever/wag/samples/gen-go-db-only/models/v9 => ./gen-go-db-only/models

replace github.com/Clever/wag/samples/gen-go-deprecated/models/v9 => ./gen-go-deprecated/models

replace github.com/Clever/wag/samples/gen-go-errors/models/v9 => ./gen-go-errors/models

replace github.com/Clever/wag/samples/gen-go-nils/models/v9 => ./gen-go-nils/models

replace github.com/Clever/wag/samples/gen-go-strings/client/v9 => ./gen-go-strings/client

replace github.com/Clever/wag/samples/gen-go-basic/client/v9 => ./gen-go-basic/client

replace github.com/Clever/wag/samples/gen-go-blog/client/v9 => ./gen-go-blog/client

replace github.com/Clever/wag/samples/gen-go-client-only/client/v9 => ./gen-go-client-only/models

replace github.com/Clever/wag/samples/gen-go-db/client/v9 => ./gen-go-db/client

replace github.com/Clever/wag/samples/gen-go-db-custom-path/client/v9 => ./gen-go-db-custom-path/client

replace github.com/Clever/wag/samples/gen-go-db-only/client/v9 => ./gen-go-db-only/client

replace github.com/Clever/wag/samples/gen-go-deprecated/client/v9 => ./gen-go-deprecated/client

replace github.com/Clever/wag/samples/gen-go-errors/client/v9 => ./gen-go-errors/client

replace github.com/Clever/wag/samples/gen-go-nils/client/v9 => ./gen-go-nils/client

replace github.com/go-openapi/swag => github.com/go-openapi/swag v0.21.1

replace github.com/Clever/wag/v9 => ../
