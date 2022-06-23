module github.com/Clever/wag/v9

go 1.16

require (
	github.com/Clever/go-utils v0.0.0-20180917210021-2dac0ec6f2ac
	github.com/Clever/wag/v8 v8.2.1
	github.com/awslabs/goformation/v2 v2.3.1
	github.com/go-openapi/errors v0.19.4
	github.com/go-openapi/jsonreference v0.20.0
	github.com/go-openapi/loads v0.19.7
	github.com/go-openapi/spec v0.19.8
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.21.1
	github.com/go-openapi/validate v0.19.15
	github.com/go-swagger/go-swagger v0.23.0
	github.com/golang/mock v1.6.0
	github.com/kevinburke/go-bindata v3.23.0+incompatible
	github.com/stretchr/testify v1.7.2
	google.golang.org/grpc v1.37.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

// For validate, something seems to have broken or changed between this pinned version and 0.16.0  (commit 7c1911976134d3a24d0c03127505163c9f16aa3b)
// That version suddenly stops accepting almost anything inline inside the `paths` block.
// TODO remove this pin.
replace github.com/go-openapi/validate => github.com/go-openapi/validate v0.0.0-20180703152151-9a6e517cddf1 // pre-modules tag 0.15.0x

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.0.0
