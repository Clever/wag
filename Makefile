include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
.PHONY: test
PKG := github.com/Clever/inter-service-api-testing/codgen-poc
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /generated)
$(eval $(call golang-version-check,1.6))

test:
	rm generated/controller.go || true
	go run main.go
	cp hardcoded/* generated/
	cd generated && go build main.go middleware.go router.go handlers.go contexts.go controller.go outputs.go

vendor: golang-godep-vendor-deps
	$(call golang-godep-vendor,$(PKGS))
