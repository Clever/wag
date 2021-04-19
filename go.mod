module github.com/Clever/wag/v7

go 1.16

require (
	github.com/Clever/discovery-go v1.7.2-0.20180111182807-aec3a7cef89e
	github.com/Clever/go-process-metrics v0.0.0-20171109172046-76790fe7fd86
	github.com/Clever/go-utils v0.0.0-20150501165843-abc25366fa8e
	github.com/GeertJohan/fgt v0.0.0-20160120143236-262f7b11eec0
	github.com/afex/hystrix-go v0.0.0-20180329224416-4847ceb883b5
	github.com/aws/aws-sdk-go v1.37.8
	github.com/awslabs/goformation/v2 v2.3.1
	github.com/codahale/hdrhistogram v0.9.0 // indirect
	github.com/donovanhide/eventsource v0.0.0-20171031113327-3ed64d21fb0b
	github.com/go-errors/errors v1.1.1
	github.com/go-openapi/errors v0.19.4
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.3
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/spec v0.19.7
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.14
	github.com/go-openapi/validate v0.19.7
	github.com/go-swagger/go-swagger v0.23.0
	github.com/golang/mock v1.1.1
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1
	github.com/kevinburke/go-bindata v3.15.0+incompatible
	github.com/opentracing-contrib/go-aws-sdk v0.0.0-20190205132030-9c29407076c0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.7.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools v0.0.0-20210106214847-113979e3529a // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/Clever/kayvee-go.v6 v6.24.1
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	gopkg.in/yaml.v2 v2.4.0
)

// For validate, something seems to have broken or changed between this pinned version and 0.16.0  (commit 7c1911976134d3a24d0c03127505163c9f16aa3b)
// That version suddenly stops accepting almost anything inline inside the `paths` block.
// TODO remove this pin.
replace github.com/go-openapi/validate => github.com/go-openapi/validate v0.0.0-20180703152151-9a6e517cddf1 // pre-modules tag 0.15.0

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.0.0
