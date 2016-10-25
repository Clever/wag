include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
.PHONY: test build
PKG := github.com/Clever/wag
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /samples/gen* | grep -v /hardcoded)
$(eval $(call golang-version-check,1.7))

MOCKGEN := $(GOPATH)/bin/mockgen
.PHONY: $(MOCKGEN)
$(MOCKGEN):
	go get -u github.com/golang/mock/mockgen

build: hardcoded/hardcoded.go
	# disable CGO and link completely statically (this is to enable us to run in containers that don't use glibc)
	CGO_ENABLED=0 go build -installsuffix cgo -o bin/wag

test: build generate $(PKGS)

generate: hardcoded/hardcoded.go $(MOCKGEN)
	./bin/wag -file samples/swagger.yml -go-package $(PKG)/samples/gen-go -js-path $(GOPATH)/src/$(PKG)/samples/gen-js
	go generate $(PKG)/samples/gen-go...
	./bin/wag -file samples/nodefinitions.yml -go-package $(PKG)/samples/gen-go-no-definitions -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-no-definitions
	go generate ${PKG}/samples/gen-go-no-definitions...
	./bin/wag -file samples/deprecated.yml -go-package $(PKG)/samples/gen-go-deprecated -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-deprecated
	go generate ${PKG}/samples/gen-go-deprecated...
	./bin/wag -file samples/wag-patch.yml -go-package $(PKG)/samples/gen-go-wag-patch -js-path $(GOPATH)/src/$(PKG)/samples/gen-js-wag-patch
	go generate ${PKG}/samples/gen-go-wag-patch...

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
