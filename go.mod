module github.com/Clever/wag/v6

go 1.13

require (
	github.com/Clever/discovery-go v1.7.2-0.20180111182807-aec3a7cef89e
	github.com/Clever/go-process-metrics v0.0.0-20171109172046-76790fe7fd86
	github.com/Clever/go-utils v0.0.0-20150501165843-abc25366fa8e
	github.com/PuerkitoBio/purell v1.1.0 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/afex/hystrix-go v0.0.0-20180329224416-4847ceb883b5
	github.com/aws/aws-sdk-go v1.37.8
	github.com/awslabs/goformation/v2 v2.3.1
	github.com/codahale/hdrhistogram v0.9.0 // indirect
	github.com/donovanhide/eventsource v0.0.0-20171031113327-3ed64d21fb0b
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/go-errors/errors v1.1.1
	github.com/go-openapi/analysis v0.0.0-20180126163718-f59a71f0ece6 // indirect
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/inflect v0.0.0-20130829110746-b1f6470ffb9c // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.0.0-20161105162150-36d33bfe519e
	github.com/go-openapi/loads v0.0.0-20171207192234-2a2b323bab96
	github.com/go-openapi/runtime v0.0.0-20180131174916-09fac855d850 // indirect
	github.com/go-openapi/spec v0.0.0-20180213232550-1de3e0542de6
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.14
	github.com/go-openapi/validate v0.0.0-20180222165948-180bba53b988
	github.com/go-swagger/go-swagger v0.2.1-0.20171112234155-b015bda48dfc
	github.com/go-swagger/scan-repo-boundary v0.0.0-20180623220736-973b3573c013 // indirect
	github.com/golang/mock v1.1.1
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/hcl v0.0.0-20171017181929-23c074d0eceb // indirect
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1
	github.com/kevinburke/go-bindata v3.15.0+incompatible
	github.com/kr/text v0.2.0 // indirect
	github.com/magiconair/properties v1.7.6 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opentracing-contrib/go-aws-sdk v0.0.0-20190205132030-9c29407076c0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pelletier/go-toml v1.1.0 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/spf13/afero v1.0.2 // indirect
	github.com/spf13/cast v1.2.0 // indirect
	github.com/spf13/jwalterweatherman v0.0.0-20180109140146-7c0cea34c8ec // indirect
	github.com/spf13/pflag v1.0.0 // indirect
	github.com/spf13/viper v1.0.1-0.20171227194143-aafc9e6bc7b7 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.19.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.19.0
	go.opentelemetry.io/contrib/propagators/aws v0.19.0
	go.opentelemetry.io/otel v0.19.0
	go.opentelemetry.io/otel/exporters/otlp v0.19.0
	go.opentelemetry.io/otel/sdk v0.19.0
	go.opentelemetry.io/otel/trace v0.19.0
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/Clever/kayvee-go.v6 v6.24.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

exclude (
	github.com/codahale/hdrhistogram v1.0.0
	github.com/codahale/hdrhistogram v1.0.1
)

exclude (
	github.com/uber-go/atomic v1.5.0
	github.com/uber-go/atomic v1.5.1
	github.com/uber-go/atomic v1.6.0
	github.com/uber-go/atomic v1.7.0
)

replace github.com/go-openapi/swag => github.com/go-openapi/swag v0.0.0-20180102232305-84f4bee7c0a6

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.0.0
