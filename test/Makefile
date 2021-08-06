include ../golang.mk

PKGS := $(shell go list ./...)

.PHONY: test install_deps debug $(PKGS)

debug:
	go env
	@echo $(PKGS)

test: debug install_deps $(PKGS)

$(PKGS):
	$(call golang-test-all-strict,$@)

install_deps:
	go mod vendor
