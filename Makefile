include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
.PHONY: test build release
PKG := github.com/Clever/wag
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /samples/gen* | grep -v /hardcoded)
VERSION := $(shell head -n 1 VERSION)
EXECUTABLE := wag

$(eval $(call golang-version-check,1.9))

build: hardcoded/hardcoded.go
	go build -o bin/wag

test: build generate $(PKGS) js-tests

js-tests:
	cd test/js && rm -rf node_modules && npm install && npm test

jsdoc2md:
	hash npm 2>/dev/null || (echo "Could not run npm, please install node" && false)
	hash jsdoc2md 2>/dev/null || npm install -g jsdoc-to-markdown@^2.0.0

generate: hardcoded/hardcoded.go jsdoc2md
	./bin/wag -file samples/swagger.yml -go-package $(PKG)/samples/gen-go -js-path $(GOPATH)/src/$(PKG)/samples/gen-js
	(cd $(GOPATH)/src/$(PKG)/samples/gen-js && jsdoc2md index.js types.js > ./README.md)
	go generate $(PKG)/samples/gen-go...
	./bin/wag -file samples/deprecated.yml -go-package $(PKG)/samples/gen-go-deprecated -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-deprecated
	(cd $(GOPATH)/src/$(PKG)/samples/gen-js-deprecated && jsdoc2md index.js types.js > ./README.md)
	go generate ${PKG}/samples/gen-go-deprecated...
	./bin/wag -file samples/errors.yml -go-package $(PKG)/samples/gen-go-errors -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-errors
	(cd $(GOPATH)/src/$(PKG)/samples/gen-js-errors && jsdoc2md index.js types.js > ./README.md)
	go generate ${PKG}/samples/gen-go-errors...
	./bin/wag -file samples/nils.yml -go-package $(PKG)/samples/gen-go-nils -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-nils
	(cd $(GOPATH)/src/$(PKG)/samples/gen-js-nils && jsdoc2md index.js types.js > ./README.md)
	go generate ${PKG}/samples/gen-go-nils...


$(PKGS): golang-test-all-strict-deps
	$(call golang-test-all-strict,$@)

$(GOPATH)/bin/go-bindata:
	go get -u github.com/jteeuwen/go-bindata/...

hardcoded/hardcoded.go: $(GOPATH)/bin/go-bindata _hardcoded/*
	$(GOPATH)/bin/go-bindata -pkg hardcoded -o hardcoded/hardcoded.go _hardcoded/
	# gofmt doesn't like what go-bindata creates
	gofmt -w hardcoded/hardcoded.go

.PHONY: $(GOPATH)/bin/glide
$(GOPATH)/bin/glide:
	@go get -u github.com/Masterminds/glide

release: hardcoded/hardcoded.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o="$@/$(EXECUTABLE)"
	tar -C $@ -zcvf "$@/$(EXECUTABLE)-$(VERSION)-linux-amd64.tar.gz" $(EXECUTABLE)
	@rm "$@/$(EXECUTABLE)"
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o="$@/$(EXECUTABLE)"
	tar -C $@ -zcvf "$@/$(EXECUTABLE)-$(VERSION)-darwin-amd64.tar.gz" $(EXECUTABLE)
	@rm "$@/$(EXECUTABLE)"

install_deps: golang-dep-vendor-deps
	$(call golang-dep-vendor)
	go build -o $(GOPATH)/bin/mockgen ./vendor/github.com/golang/mock/mockgen
