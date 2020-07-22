include golang.mk
include node.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
NODE_VERSION := "v7"
$(eval $(call node-version-check,$(NODE_VERSION)))

export PATH := $(PWD)/bin:$(PATH)
PKG := github.com/Clever/wag
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /samples/gen* | grep -v /hardcoded)
VERSION := $(shell head -n 1 VERSION)
EXECUTABLE := wag
PKGS := $(PKGS) $(PKG)/samples/gen-go-db/server/db/dynamodb

$(eval $(call golang-version-check,1.13))

.PHONY: test build release js-tests jsdoc2md go-generate generate $(PKGS) install_deps

build: go-generate
	go build -o bin/wag

test: build generate $(PKGS) js-tests

js-tests:
	cd test/js && rm -rf node_modules && npm install && npm test

jsdoc2md:
	hash npm 2>/dev/null || (echo "Could not run npm, please install node" && false)
	hash jsdoc2md 2>/dev/null || npm install -g jsdoc-to-markdown@^6.0.0

go-generate:
	go generate ./hardcoded/
	go generate ./server/gendb/

generate: build jsdoc2md
	./bin/wag -file samples/swagger.yml -go-package $(PKG)/samples/gen-go -js-path $(GOPATH)/src/$(PKG)/samples/gen-js
	cd $(GOPATH)/src/$(PKG)/samples/gen-js && jsdoc2md index.js types.js > ./README.md
	go generate $(PKG)/samples/gen-go...
	./bin/wag -file samples/deprecated.yml -go-package $(PKG)/samples/gen-go-deprecated -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-deprecated
	cd $(GOPATH)/src/$(PKG)/samples/gen-js-deprecated && jsdoc2md index.js types.js > ./README.md
	go generate ${PKG}/samples/gen-go-deprecated...
	./bin/wag -file samples/errors.yml -go-package $(PKG)/samples/gen-go-errors -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-errors
	cd $(GOPATH)/src/$(PKG)/samples/gen-js-errors && jsdoc2md index.js types.js > ./README.md
	go generate ${PKG}/samples/gen-go-errors...
	./bin/wag -file samples/nils.yml -go-package $(PKG)/samples/gen-go-nils -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-nils
	cd $(GOPATH)/src/$(PKG)/samples/gen-js-nils && jsdoc2md index.js types.js > ./README.md
	go generate ${PKG}/samples/gen-go-nils...
	./bin/wag -file samples/blog.yml -go-package $(PKG)/samples/gen-go-blog -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-blog
	cd $(GOPATH)/src/$(PKG)/samples/gen-js-blog && jsdoc2md index.js types.js > ./README.md
	go generate ${PKG}/samples/gen-go-nils...
	./bin/wag -file samples/db.yml -go-package $(PKG)/samples/gen-go-db -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-db
	cd $(GOPATH)/src/$(PKG)/samples/gen-js-db && jsdoc2md index.js types.js > ./README.md
	go generate ${PKG}/samples/gen-go-db...

$(PKGS): golang-test-all-strict-deps
	$(call golang-test-all-strict,$@)

release:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o="$@/$(EXECUTABLE)"
	tar -C $@ -zcvf "$@/$(EXECUTABLE)-$(VERSION)-linux-amd64.tar.gz" $(EXECUTABLE)
	@rm "$@/$(EXECUTABLE)"
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o="$@/$(EXECUTABLE)"
	tar -C $@ -zcvf "$@/$(EXECUTABLE)-$(VERSION)-darwin-amd64.tar.gz" $(EXECUTABLE)
	@rm "$@/$(EXECUTABLE)"

install_deps: golang-dep-vendor-deps
	$(call golang-dep-vendor)
	go build -o bin/go-bindata ./vendor/github.com/tmthrgd/go-bindata/go-bindata
	go build -o bin/mockgen    ./vendor/github.com/golang/mock/mockgen
	mkdir -p $(GOPATH)/bin
	cp bin/mockgen $(GOPATH)/bin/mockgen
