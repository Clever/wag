module github.com/Clever/wag/samples/v9

go 1.21

toolchain go1.21.0

require (
	github.com/Clever/go-process-metrics v0.4.0
	github.com/Clever/kayvee-go/v7 v7.6.0
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
	gopkg.in/Clever/kayvee-go.v6 v6.27.0
)

require (
	cloud.google.com/go v0.102.0 // indirect
	cloud.google.com/go/bigquery v1.8.0 // indirect
	cloud.google.com/go/compute v1.7.0 // indirect
	cloud.google.com/go/datastore v1.1.0 // indirect
	cloud.google.com/go/iam v0.3.0 // indirect
	cloud.google.com/go/pubsub v1.3.1 // indirect
	cloud.google.com/go/storage v1.22.1 // indirect
	dmitri.shuralyov.com/gpu/mtl v0.0.0-20190408044501-666a987793e9 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802 // indirect
	github.com/Clever/discovery-go v1.8.1 // indirect
	github.com/Clever/go-utils v0.0.0-20180917210021-2dac0ec6f2ac // indirect
	github.com/Clever/wag/logging/wagclientlogger v0.0.0-20221024182247-2bf828ef51be // indirect
	github.com/OneOfOne/xxhash v1.2.2 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc // indirect
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf // indirect
	github.com/antihax/optional v1.0.0 // indirect
	github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/awslabs/goformation/v2 v2.3.1 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.0 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/census-instrumentation/opencensus-proto v0.2.1 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/client9/misspell v0.3.4 // indirect
	github.com/cncf/udpa/go v0.0.0-20210930031921-04548b0d99d4 // indirect
	github.com/cncf/xds/go v0.0.0-20211011173535-cb28da3451f1 // indirect
	github.com/corbym/gocrest v1.0.3 // indirect
	github.com/coreos/bbolt v1.3.2 // indirect
	github.com/coreos/etcd v3.3.10+incompatible // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/coreos/go-semver v0.2.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/creack/pty v1.1.9 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dgryski/go-sip13 v0.0.0-20181026042036-e10d5fee7954 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/donovanhide/eventsource v0.0.0-20171031113327-3ed64d21fb0b // indirect
	github.com/eapache/go-resiliency v1.3.0 // indirect
	github.com/envoyproxy/go-control-plane v0.10.2-0.20220325020618-49ff273808a1 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/globalsign/mgo v0.0.0-20180905125535-1ca0a4f7cbcb // indirect
	github.com/go-gl/glfw v0.0.0-20190409004039-e6da0acd62b1 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20200222043503-6f7a984d4dc4 // indirect
	github.com/go-kit/kit v0.8.0 // indirect
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.20.2 // indirect
	github.com/go-openapi/inflect v0.19.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/runtime v0.24.1 // indirect
	github.com/go-openapi/spec v0.20.7 // indirect
	github.com/go-openapi/validate v0.22.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/go-swagger/go-swagger v0.23.0 // indirect
	github.com/go-swagger/scan-repo-boundary v0.0.0-20180623220736-973b3573c013 // indirect
	github.com/gobuffalo/attrs v0.0.0-20190224210810-a9411de4debd // indirect
	github.com/gobuffalo/depgen v0.1.0 // indirect
	github.com/gobuffalo/envy v1.7.0 // indirect
	github.com/gobuffalo/flect v0.1.3 // indirect
	github.com/gobuffalo/genny v0.1.1 // indirect
	github.com/gobuffalo/gitgen v0.0.0-20190315122116-cc086187d211 // indirect
	github.com/gobuffalo/gogen v0.1.1 // indirect
	github.com/gobuffalo/logger v0.0.0-20190315122211-86e12af44bc2 // indirect
	github.com/gobuffalo/mapi v1.0.2 // indirect
	github.com/gobuffalo/packd v0.1.0 // indirect
	github.com/gobuffalo/packr/v2 v2.2.0 // indirect
	github.com/gobuffalo/syncx v0.0.0-20190224160051-33c29581e754 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/martian v2.1.0+incompatible // indirect
	github.com/google/martian/v3 v3.2.1 // indirect
	github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1 // indirect
	github.com/google/renameio v0.1.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.0.0-20220520183353-fd19c99a87aa // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/googleapis/go-type-adapters v1.0.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20200824232613-28f6c0f3b639 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jmespath/go-jmespath/internal/testify v1.5.1 // indirect
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/julienschmidt/httprouter v1.2.0 // indirect
	github.com/karrick/godirwalk v1.10.3 // indirect
	github.com/kevinburke/go-bindata v3.24.0+incompatible // indirect
	github.com/kisielk/errcheck v1.2.0 // indirect
	github.com/kisielk/gotool v1.0.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/kr/pty v1.1.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/markbates/oncer v0.0.0-20181203154359-bf2de49a0be2 // indirect
	github.com/markbates/safe v1.0.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/onsi/ginkgo v1.5.0 // indirect
	github.com/onsi/gomega v1.2.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pelletier/go-toml v1.9.1 // indirect
	github.com/pkg/diff v0.0.0-20210226163009-20ebb0f2a09e // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/prometheus/client_golang v0.9.3 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/prometheus/common v0.4.0 // indirect
	github.com/prometheus/procfs v0.0.0-20190507164030-5867b95ac084 // indirect
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/rogpeppe/fastuuid v1.2.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/sanathkr/go-yaml v0.0.0-20170819195128-ed9d249f429b // indirect
	github.com/sanathkr/yaml v0.0.0-20170819201035-0056894fa522 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/spaolacci/murmur3 v0.0.0-20180118202830-f09979ecbc72 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v0.0.3 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.6.2 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/toqueteos/webbrowser v1.2.0 // indirect
	github.com/ugorji/go v1.1.4 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v0.0.0-20180714160509-73f8eece6fdc // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yuin/goldmark v1.4.13 // indirect
	go.etcd.io/bbolt v1.3.2 // indirect
	go.mongodb.org/mongo-driver v1.10.1 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.31.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.31.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.26.0 // indirect
	go.opentelemetry.io/otel/metric v0.31.0 // indirect
	go.opentelemetry.io/otel/sdk/export/metric v0.26.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.31.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6 // indirect
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mobile v0.0.0-20190719004257-d2bd2a29d028 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20220826154423-83b083e8dc8b // indirect
	golang.org/x/oauth2 v0.0.0-20220822191816-0ebed06d0094 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.0.0-20220829200755-d48e67d00261 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.13-0.20220908144252-ce397412b6a4 // indirect
	google.golang.org/api v0.84.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220829175752-36a9c930ecbf // indirect
	google.golang.org/grpc v1.49.0 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/errgo.v2 v2.1.0 // indirect
	gopkg.in/ini.v1 v1.54.0 // indirect
	gopkg.in/resty.v1 v1.12.0 // indirect
	gopkg.in/square/go-jose.v2 v2.4.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
	rsc.io/binaryregexp v0.2.0 // indirect
	rsc.io/quote/v3 v3.1.0 // indirect
	rsc.io/sampler v1.3.0 // indirect
)

replace github.com/go-openapi/validate => github.com/go-openapi/validate v0.0.0-20180703152151-9a6e517cddf1 // pre-modules tag 0.15.0x

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
