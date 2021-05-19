module github.com/Clever/wag/v7

go 1.13

require (
	github.com/Clever/go-utils v0.0.0-20150501165843-abc25366fa8e
	github.com/Clever/wag/v7/samples v0.0.0-internal
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/awslabs/goformation/v2 v2.3.1
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/inflect v0.0.0-20130829110746-b1f6470ffb9c // indirect
	github.com/go-openapi/jsonreference v0.0.0-20161105162150-36d33bfe519e
	github.com/go-openapi/loads v0.0.0-20171207192234-2a2b323bab96
	github.com/go-openapi/spec v0.0.0-20180213232550-1de3e0542de6
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.0.0-20180222165948-180bba53b988
	github.com/go-swagger/go-swagger v0.2.1-0.20171112234155-b015bda48dfc
	github.com/go-swagger/scan-repo-boundary v0.0.0-20180623220736-973b3573c013 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/mock v1.5.0
	github.com/hashicorp/hcl v0.0.0-20171017181929-23c074d0eceb // indirect
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/kevinburke/go-bindata v3.22.0+incompatible
	github.com/kr/text v0.2.0 // indirect
	github.com/magiconair/properties v1.7.6 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pelletier/go-toml v1.1.0 // indirect
	github.com/spf13/afero v1.0.2 // indirect
	github.com/spf13/cast v1.2.0 // indirect
	github.com/spf13/jwalterweatherman v0.0.0-20180109140146-7c0cea34c8ec // indirect
	github.com/spf13/pflag v1.0.0 // indirect
	github.com/spf13/viper v1.0.1-0.20171227194143-aafc9e6bc7b7 // indirect
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/sdk v0.15.0
	google.golang.org/grpc v1.36.0 // indirect
	gopkg.in/Clever/kayvee-go.v6 v6.24.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	gopkg.in/yaml.v2 v2.4.0 // indirect
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

replace github.com/Clever/wag/v7/samples v0.0.0-internal => ./samples
