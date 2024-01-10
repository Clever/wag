module github.com/Clever/wag/v9

go 1.21

require (
	github.com/Clever/go-utils v0.0.0-20180917210021-2dac0ec6f2ac
	github.com/awslabs/goformation/v2 v2.3.1
	github.com/go-openapi/errors v0.19.4
	github.com/go-openapi/jsonreference v0.20.2
	github.com/go-openapi/loads v0.19.7
	github.com/go-openapi/spec v0.19.8
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.22.3
	github.com/go-openapi/validate v0.19.15
	github.com/go-swagger/go-swagger v0.23.0
	github.com/golang/mock v1.6.0
	github.com/kevinburke/go-bindata v3.24.0+incompatible
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/asaskevich/govalidator v0.0.0-20200108200545-475eaeb16496 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-openapi/analysis v0.19.10 // indirect
	github.com/go-openapi/inflect v0.19.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/runtime v0.19.12 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/pelletier/go-toml v1.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.6.2 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	go.mongodb.org/mongo-driver v1.3.1 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sys v0.0.0-20220829200755-d48e67d00261 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.13-0.20220908144252-ce397412b6a4 // indirect
	gopkg.in/ini.v1 v1.54.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// For validate, something seems to have broken or changed between this pinned version and 0.16.0  (commit 7c1911976134d3a24d0c03127505163c9f16aa3b)
// That version suddenly stops accepting almost anything inline inside the `paths` block.
// TODO remove this pin.
replace github.com/go-openapi/validate => github.com/go-openapi/validate v0.0.0-20180703152151-9a6e517cddf1 // pre-modules tag 0.15.0x

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.0.0

// Newer versions of swag use a newer version of yaml which fix a "bug" where the yaml spec was not fully
// addhered  to. Unfortunately this was a breaking change in yaml parsing, so pinning swag for now.
replace github.com/go-openapi/swag => github.com/go-openapi/swag v0.21.1
