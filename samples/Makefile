include ../golang.mk

TEST_PKGS := $(shell go list -mod=mod ./test...)
TEST_DB_PKG := "./gen-go-db/server/db/dynamodb"

.PHONY: generate test $(TEST_PKGS) $(TEST_DB_PKG) install_deps

# Function to generate code, run jsdoc2md, and tidy go modules
define generate_code
	../bin/wag -with-tests -file $(1) -output-path $(2) -js-path $(3) $(4)
	cd $(3) && jsdoc2md index.js types.js > ./README.md
	cd $(2)/client && go mod tidy
	cd $(2)/models && go mod tidy
	go generate -mod=mod $(2)...
endef

# Function to generate code without client
define generate_code_no_client
	../bin/wag -with-tests -file $(1) -output-path $(2) $(3)
	cd $(2)/models && go mod tidy
	go generate -mod=mod $(2)...
endef

# Function to generate code without jsdoc2md
define generate_code_no_js
	../bin/wag -with-tests -file $(1) -output-path $(2) $(3)
	cd $(2)/models && go mod tidy
	go generate -mod=mod $(2)...
endef

# Function to generate only a js client with x-unwrapped-fetch set. The spec is copied from
# $(1) with the flag added, so the client can be tested with the same tests as the default client.
define generate_unwrapped_fetch_js
	echo "# Generated from $(1) by 'make generate'. Do not edit." > $(2)
	awk '{print} /x-npm-package:/ {print "  x-unwrapped-fetch: true"}' $(1) >> $(2)
	../bin/wag -file $(2) -output-path $(3) -js-path $(3) -client-only -client-language js
	cd $(3) && jsdoc2md index.js types.js > ./README.md
endef

# These go:generate don't actually care about dependencies, so we use -mod=mod to override -mod=vendor
# If not for that, it's difficult to bootstrap
generate:
	echo  $(TEST_PKGS) $(TEST_DB_PKG)
	$(call generate_code,./swagger.yml,./gen-go-basic,./gen-js)
	$(call generate_code,./deprecated.yml,./gen-go-deprecated,./gen-js-deprecated)
	$(call generate_code,./errors.yml,./gen-go-errors,./gen-js-errors)
	$(call generate_code,./nils.yml,./gen-go-nils,./gen-js-nils)
	$(call generate_code,./strings.yml,./gen-go-strings,./gen-js-strings)
	$(call generate_code,./blog.yml,./gen-go-blog,./gen-js-blog)
	$(call generate_code,./db.yml,./gen-go-db,./gen-js-db)
	$(call generate_code,./swagger.yml,./gen-go-client-only,./gen-js-client-only,--client-only)
	$(call generate_code_no_client,./db.yml,./gen-go-db-only,--dynamo-only)
	$(call generate_code,./db.yml,./gen-go-db-custom-path,./gen-js-db-custom-path,-dynamo-path db)
	$(call generate_unwrapped_fetch_js,./swagger.yml,./unwrapped-fetch.yml,./gen-js-unwrapped-fetch)

	go install -mod=mod golang.org/x/tools/cmd/goimports@v0.24.0
	goimports -w .

$(TEST_PKGS) $(TEST_DB_PKG): golang-test-all-strict-deps
	$(call golang-test-all,$@)

install_deps:
	go mod tidy
	go mod vendor

test: install_deps $(TEST_PKGS) $(TEST_DB_PKG)
