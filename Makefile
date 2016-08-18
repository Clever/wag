include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
.PHONY: test build hardcoded.go
PKG := github.com/Clever/wag
PKGS := $(shell go list ./... | grep -v /vendor)
$(eval $(call golang-version-check,1.7))

build: hardcoded.go
	# disable CGO and link completely statically (this is to enable us to run in containers that don't use glibc)
	CGO_ENABLED=0 go build -installsuffix cgo -o bin/wag

test: build generate $(PKGS)

generate: hardcoded.go
	rm -rf generated
	./bin/wag -file swagger.yml -package $(PKG)/generated

$(PKGS): golang-test-all-deps
	$(call golang-test-all,$@)

$(GOPATH)/bin/go-bindata:
	go get -u github.com/jteeuwen/go-bindata/...

hardcoded.go: $(GOPATH)/bin/go-bindata hardcoded/*
	$(GOPATH)/bin/go-bindata -o hardcoded.go hardcoded/
	# gofmt doesn't like what go-bindata creates
	gofmt hardcoded.go > temp
	cp temp hardcoded.go
	rm temp

vendor: golang-godep-vendor-deps
	$(call golang-godep-vendor,$(PKGS))
