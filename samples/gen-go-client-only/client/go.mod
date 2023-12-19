
module github.com/Clever/wag/samples/gen-go-client-only/client/v9

go 1.21

require (
	github.com/Clever/discovery-go v1.8.1
	github.com/Clever/wag/logging/wagclientlogger v0.0.0-20221024182247-2bf828ef51be
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/donovanhide/eventsource v0.0.0-20171031113327-3ed64d21fb0b
	github.com/smartystreets/goconvey v1.7.2 // indirect
)
//Replace directives will work locally but mess up imports.
replace github.com/Clever/wag/samples/gen-go-client-only/models/v9 => ../models 