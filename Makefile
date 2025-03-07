include golang.mk
include node.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
NODE_VERSION := "v18"
$(eval $(call node-version-check,$(NODE_VERSION)))

export PATH := $(PWD)/bin:$(PATH)
MAJOR_VERSION := $(shell head -n 1 VERSION | sed 's/[[:alpha:]|[:space:]]//g' | cut -d. -f1)
PKG := github.com/Clever/wag/v$(MAJOR_VERSION)
PKGS := $(shell go list ./... | grep -v /hardcoded | grep -v /tools | grep -v /gendb)
VERSION := $(shell head -n 1 VERSION)
EXECUTABLE := wag

$(eval $(call golang-version-check,1.24))

.PHONY: test build release js-tests jsdoc2md go-generate generate $(PKGS) install_deps

build: go-generate
	go build -o bin/wag

test: build generate $(PKGS) js-tests
	$(MAKE) -C samples test

js-tests:
	cd samples/gen-js && rm -rf node_modules && npm install
	cd samples/test/js && rm -rf node_modules && npm install && npm test

jsdoc2md:
	hash npm 2>/dev/null || (echo "Could not run npm, please install node" && false)
	hash jsdoc2md 2>/dev/null || npm install -g jsdoc-to-markdown@^4.0.0

go-generate:
	go generate ./hardcoded/
	go generate ./server/gendb/

generate: build jsdoc2md
	$(MAKE) -C samples generate

$(PKGS): golang-test-all-strict-deps
	$(call golang-test-all-strict,$@)

release:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o="$@/$(EXECUTABLE)"
	tar -C $@ -zcvf "$@/$(EXECUTABLE)-$(VERSION)-linux-amd64.tar.gz" $(EXECUTABLE)
	@rm "$@/$(EXECUTABLE)"
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o="$@/$(EXECUTABLE)"
	tar -C $@ -zcvf "$@/$(EXECUTABLE)-$(VERSION)-darwin-amd64.tar.gz" $(EXECUTABLE)
	@rm "$@/$(EXECUTABLE)"

install_deps:
	go mod tidy
	go mod vendor
	go build -o bin/go-bindata ./vendor/github.com/kevinburke/go-bindata/go-bindata
	go build -o bin/mockgen    ./vendor/github.com/golang/mock/mockgen

