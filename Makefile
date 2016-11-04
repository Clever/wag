include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
.PHONY: test build release
PKG := github.com/Clever/wag
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /samples/gen* | grep -v /hardcoded)
VERSION := $(shell head -n 1 VERSION)
EXECUTABLE := wag

$(eval $(call golang-version-check,1.7))

MOCKGEN := $(GOPATH)/bin/mockgen
.PHONY: $(MOCKGEN)
$(MOCKGEN):
	go get -u github.com/golang/mock/mockgen

build: hardcoded/hardcoded.go
	go build -o bin/wag

test: build generate $(PKGS) js-tests

js-tests:
	cd test/js && npm install && npm test

generate: hardcoded/hardcoded.go $(MOCKGEN)
	./bin/wag -file samples/swagger.yml -go-package $(PKG)/samples/gen-go -js-path $(GOPATH)/src/$(PKG)/samples/gen-js
	go generate $(PKG)/samples/gen-go...
	./bin/wag -file samples/deprecated.yml -go-package $(PKG)/samples/gen-go-deprecated -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-deprecated
	go generate ${PKG}/samples/gen-go-deprecated...
	./bin/wag -file samples/errors.yml -go-package $(PKG)/samples/gen-go-errors -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-errors
	go generate ${PKG}/samples/gen-go-errors...

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

install_deps: $(GOPATH)/bin/glide
	$(GOPATH)/bin/glide install -v

release: hardcoded/hardcoded.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o="$@/$(EXECUTABLE)"
	tar -C $@ -zcvf "$@/$(EXECUTABLE)-v$(VERSION)-linux-amd64.tar.gz" $(EXECUTABLE)
	@rm "$@/$(EXECUTABLE)"
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o="$@/$(EXECUTABLE)"
	tar -C $@ -zcvf "$@/$(EXECUTABLE)-v$(VERSION)-darwin-amd64.tar.gz" $(EXECUTABLE)
	@rm "$@/$(EXECUTABLE)"
