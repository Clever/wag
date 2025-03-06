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
	../bin/wag -with-tests -file $(1) -output-path $(2) -js-path $(3) $(4)
	cd $(3) && jsdoc2md index.js types.js > ./README.md
	cd $(2)/models && go mod tidy
	go generate -mod=mod $(2)...
endef

define generate_code_db_only
	../bin/wag -with-tests -file $(1) -output-path $(2) -dynamo-only
	cd $(2)/models && go mod tidy
	go generate -mod=mod $(2)...
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
	$(call generate_code_no_js,./db.yml,./gen-go-db,./gen-js-db)
	$(call generate_code_no_client,./swagger.yml,./gen-go-client-only,./gen-js-client-only,--client-only)
	$(call generate_code_db_only,./db.yml,./gen-go-db-only,,--dynamo-only)
	$(call generate_code,./db.yml,./gen-go-db-custom-path,./gen-js-db-custom-path,-dynamo-path db)

	go install -mod=mod golang.org/x/tools/cmd/goimports@v0.24.0
	goimports -w .
	
$(TEST_PKGS) $(TEST_DB_PKG): golang-test-all-strict-deps
	$(call golang-test-all,$@)

install_deps:
	go mod tidy
	go mod vendor

test: install_deps $(TEST_PKGS) $(TEST_DB_PKG)
