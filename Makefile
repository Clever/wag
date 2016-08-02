include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
.PHONY: test
PKG := github.com/Clever/inter-service-api-testing/codegen-poc
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /generated)
$(eval $(call golang-version-check,1.6))

test:
	rm generated/controller.go || true
	cp hardcoded/* generated/
	go run main.go genclients.go -file swagger.yml -package $(PKG)/generated
	cd impl && go build

vendor: golang-godep-vendor-deps
	$(call golang-godep-vendor,$(PKGS))
