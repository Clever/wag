include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile
.PHONY: test build
PKG := github.com/Clever/wag
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /generated)
$(eval $(call golang-version-check,1.7))

MOCKGEN := $(GOPATH)/bin/mockgen
.PHONY: $(MOCKGEN)
$(MOCKGEN):
	go get -u github.com/golang/mock/mockgen

build: hardcoded.go
	# disable CGO and link completely statically (this is to enable us to run in containers that don't use glibc)
	CGO_ENABLED=0 go build -installsuffix cgo -o bin/wag

test: build generate $(PKGS)

generate: hardcoded.go $(MOCKGEN)
	rm -rf generated
	./bin/wag -file swagger.yml -package $(PKG)/generated
	go generate $(PKG)/generated...

$(PKGS): golang-test-all-deps
	$(call golang-test-all,$@)

$(GOPATH)/bin/go-bindata:
	go get -u github.com/jteeuwen/go-bindata/...

hardcoded.go: $(GOPATH)/bin/go-bindata hardcoded/*
	$(GOPATH)/bin/go-bindata -o hardcoded.go hardcoded/
	# gofmt doesn't like what go-bindata creates
	gofmt -w hardcoded.go

.PHONY: $(GOPATH)/bin/glide
$(GOPATH)/bin/glide:
	@go get github.com/Masterminds/glide
	
install_deps: $(GOPATH)/bin/glide
	@# stash dependencies into vendor directory (-v removes nested vendor dirs)
	$(GOPATH)/bin/glide install -v
