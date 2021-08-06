include ../golang.mk

PKGS := $(shell go list ./...)

.PHONY: test install_deps

test: install_deps $(PKGS)

$(PKGS):
	$(call golang-test-all-strict,$@)

install_deps:
	go mod vendor
