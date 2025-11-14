include golang.mk
include node.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
NODE_VERSION := "v24"
$(eval $(call node-version-check,$(NODE_VERSION)))

export PATH := $(PWD)/bin:$(PATH)
MAJOR_VERSION := $(shell head -n 1 VERSION | sed 's/[[:alpha:]|[:space:]]//g' | cut -d. -f1)
PKG := github.com/Clever/wag/v$(MAJOR_VERSION)
PKGS := $(shell go list ./... | grep -v /hardcoded | grep -v /tools | grep -v /gendb)
RAWVERSION :=$(shell head -n 1 VERSION)
VERSION := $(RAWVERSION)$(shell if [[ -z "$(CI)" ]]; then echo "-dev"; fi)
EXECUTABLE := wag

$(eval $(call golang-version-check,1.24))

.PHONY: test build release js-tests jsdoc2md go-generate generate $(PKGS) install_deps

build: go-generate
	go build -ldflags="-X main.version=$(VERSION)" -o bin/wag

test: build generate $(PKGS) js-tests
	$(MAKE) -C samples test
	echo "Currently DB tests are disabled because they are failing and aren't able to prioritize" \
	"invesigating and fixing them. Remove t.Skip() from server/gendb/dynamodb_test.go.tmpl to re-enable" \
	"https://clever.atlassian.net/browse/INFRANG-6880"

js-tests:
	cd samples/gen-js && rm -rf node_modules && npm install
	cd samples/test/js && rm -rf node_modules && npm install && npm test

jsdoc2md:
	hash npm 2>/dev/null || (echo "Could not run npm, please install node" && false)
	@if [ ! -f ./node_modules/.bin/jsdoc2md ]; then \
		npm install jsdoc-to-markdown@9.0.0; \
	elif [ -f ./node_modules/jsdoc-to-markdown/package.json ]; then \
		INSTALLED_VERSION=$$(node -p "require('./node_modules/jsdoc-to-markdown/package.json').version" 2>/dev/null || echo ""); \
		if [ "$$INSTALLED_VERSION" != "9.0.0" ]; then \
			echo "jsdoc-to-markdown version $$INSTALLED_VERSION detected, updating to 9.0.0..."; \
			npm install jsdoc-to-markdown@9.0.0; \
		fi; \
	fi

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
	npm install

