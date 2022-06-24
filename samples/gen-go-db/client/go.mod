
module github.com/Clever/wag/samples/v8/gen-go-db/client

go 1.16

require (
	github.com/Clever/wag/samples/v8/gen-go-db/models v0.0.0
	github.com/Clever/discovery-go v1.7.2
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/donovanhide/eventsource v0.0.0-20171031113327-3ed64d21fb0b
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/go-openapi/analysis v0.21.2 // indirect
	github.com/go-openapi/errors v0.20.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/loads v0.21.1 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/strfmt v0.21.2 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-openapi/validate v0.22.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.1-0.20200424115421-065759f9c3d7 // indirect
	go.mongodb.org/mongo-driver v1.7.5 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/Clever/kayvee-go.v6 v6.27.0 // indirect //Still want to remove this once I merge discovery-go PR. 
	gopkg.in/yaml.v2 v2.4.0 // indirect

	replace github.com/Clever/wag/samples/v8/gen-go-db/models => ../models
)
