include golang.mk
include node.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
NODE_VERSION := "v7"
$(eval $(call node-version-check,$(NODE_VERSION)))

export PATH := $(PWD)/bin:$(PATH)
MAJOR_VERSION := $(shell head -n 1 VERSION | sed 's/[[:alpha:]|[:space:]]//g' | cut -d. -f1)
PKG := github.com/Clever/wag/v$(MAJOR_VERSION)
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /samples/gen* | grep -v /hardcoded | grep -v /tools)
VERSION := $(shell head -n 1 VERSION)
EXECUTABLE := wag
PKGS := $(PKGS) $(PKG)/samples/gen-go-db/server/db/dynamodb

$(eval $(call golang-version-check,1.16))

.PHONY: test build release js-tests jsdoc2md go-generate generate $(PKGS) install_deps

build: go-generate
	go build -o bin/wag

test: build generate $(PKGS) js-tests

js-tests:
	cd test/js && rm -rf node_modules && npm install && npm test

jsdoc2md:
	hash npm 2>/dev/null || (echo "Could not run npm, please install node" && false)
	hash jsdoc2md 2>/dev/null || npm install -g jsdoc-to-markdown@^4.0.0

go-generate:
	go generate ./hardcoded/
	go generate ./server/gendb/

generate: build jsdoc2md
	./bin/wag -file ./samples/swagger.yml -output-path ./samples/gen-go -js-path ./samples/gen-js
	cd ./samples/gen-js && jsdoc2md index.js types.js > ./README.md
	go generate ./samples/gen-go...
	./bin/wag -file ./samples/deprecated.yml -output-path ./samples/gen-go-deprecated -js-path ./samples/gen-js-deprecated
	cd ./samples/gen-js-deprecated && jsdoc2md index.js types.js > ./README.md
	go generate ./samples/gen-go-deprecated...
	./bin/wag -file ./samples/errors.yml -output-path ./samples/gen-go-errors -js-path ./samples/gen-js-errors
	cd ./samples/gen-js-errors && jsdoc2md index.js types.js > ./README.md
	go generate ./samples/gen-go-errors...
	./bin/wag -file ./samples/nils.yml -output-path ./samples/gen-go-nils -js-path ./samples/gen-js-nils
	cd ./samples/gen-js-nils && jsdoc2md index.js types.js > ./README.md
	go generate ./samples/gen-go-nils...
	./bin/wag -file ./samples/blog.yml -output-path ./samples/gen-go-blog -js-path ./samples/gen-js-blog
	cd ./samples/gen-js-blog && jsdoc2md index.js types.js > ./README.md
	go generate ./samples/gen-go-nils...
	./bin/wag -file ./samples/db.yml -output-path ./samples/gen-go-db -js-path ./samples/gen-js-db
	cd ./samples/gen-js-db && jsdoc2md index.js types.js > ./README.md
	go generate ./samples/gen-go-db...
	./bin/wag -file ./samples/swagger.yml -output-path ./samples/gen-go-client-only -js-path ./samples/gen-js-client-only --client-only
	cd ./samples/gen-js-client-only && jsdoc2md index.js types.js > ./README.md
	go generate ./samples/gen-go-client-only...
	./bin/wag -file ./samples/db.yml -output-path ./samples/gen-go-db-only -dynamo-only
	go generate ./samples/gen-go-db-only...
	./bin/wag -file ./samples/db.yml -output-path ./samples/gen-go-db-custom-path -js-path ./samples/gen-js-db-custom-path -dynamo-path db
	cd ./samples/gen-js-db-custom-path && jsdoc2md index.js types.js > ./README.md
	go generate ./samples/gen-go-db-custom-path...

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
	go mod vendor
	go build -o bin/go-bindata ./vendor/github.com/kevinburke/go-bindata/go-bindata
	go build -o bin/mockgen    ./vendor/github.com/golang/mock/mockgen
	mkdir -p $(GOPATH)/bin
	cp bin/mockgen $(GOPATH)/bin/mockgen
