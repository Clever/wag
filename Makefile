include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
.PHONY: test build
PKG := github.com/Clever/wag
PKGS := $(shell go list ./... | grep -v /vendor)
$(eval $(call golang-version-check,1.6))

build:
	# disable CGO and link completely statically (this is to enable us to run in containers that don't use glibc)
	CGO_ENABLED=0 go build -installsuffix cgo -o bin/wag

test: build
	rm -rf generated/*
	./bin/wag -file swagger.yml -package $(PKG)/generated
	cd impl && go build
	# Temporarily run the client here since it isn't used in impl
	cd generated/client && go build

vendor: golang-godep-vendor-deps
	$(call golang-godep-vendor,$(PKGS))
