include ../golang.mk

TEST_PKGS := $(shell go list -mod=mod ./test...)
TEST_DB_PKG := "./gen-go-db/server/db/dynamodb"

.PHONY: generate test $(TEST_PKGS) $(TEST_DB_PKG) install_deps

# These go:generate don't actually care about dependencies, so we use -mod=mod to override -mod=vendor
# If not for that, it's difficult to bootstrap
generate:
	@echo  $(TEST_PKGS) $(TEST_DB_PKG)
	
	../bin/wag -with-tests -file ./swagger.yml -output-path ./gen-go-basic -js-path ./gen-js
	cd ./gen-js && jsdoc2md index.js types.js > ./README.md    
	cd ./gen-go-basic/client && go mod tidy
	cd ./gen-go-basic/models && go mod tidy
	go generate -mod=mod ./gen-go-basic...


	../bin/wag -with-tests -file ./deprecated.yml -output-path ./gen-go-deprecated -js-path ./gen-js-deprecated
	cd ./gen-js-deprecated && jsdoc2md index.js types.js > ./README.md
	cd ./gen-go-deprecated/client && go mod tidy
	cd ./gen-go-deprecated/models && go mod tidy
	go generate -mod=mod ./gen-go-deprecated...

	../bin/wag -with-tests -file ./errors.yml -output-path ./gen-go-errors -js-path ./gen-js-errors
	cd ./gen-js-errors && jsdoc2md index.js types.js > ./README.md
	cd ./gen-go-errors/client && go mod tidy
	cd ./gen-go-errors/models && go mod tidy
	go generate -mod=mod ./gen-go-errors...

	../bin/wag -with-tests -file ./nils.yml -output-path ./gen-go-nils -js-path ./gen-js-nils
	cd ./gen-js-nils && jsdoc2md index.js types.js > ./README.md
	cd ./gen-go-nils/client && go mod tidy
	cd ./gen-go-nils/models && go mod tidy
	go generate -mod=mod ./gen-go-nils...

	../bin/wag -with-tests -file ./strings.yml -output-path ./gen-go-strings -js-path ./gen-js-strings
	cd ./gen-js-strings && jsdoc2md index.js types.js > ./README.md
	cd ./gen-go-strings/client && go mod tidy
	cd ./gen-go-strings/models && go mod tidy
	go generate -mod=mod ./gen-go-strings...

	../bin/wag -with-tests -file ./blog.yml -output-path ./gen-go-blog -js-path ./gen-js-blog
	cd ./gen-js-blog && jsdoc2md index.js types.js > ./README.md
	cd ./gen-go-blog/client && go mod tidy
	cd ./gen-go-blog/models && go mod tidy
	go generate -mod=mod ./gen-go-nils...

	../bin/wag -with-tests -file ./db.yml -output-path ./gen-go-db -js-path ./gen-js-db
	cd ./gen-js-db && jsdoc2md index.js types.js > ./README.md
	cd ./gen-go-db/client && go mod tidy
	cd ./gen-go-db/models && go mod tidy
	go generate -mod=mod ./gen-go-db...

	../bin/wag -with-tests -file ./swagger.yml -output-path ./gen-go-client-only -js-path ./gen-js-client-only --client-only
	cd ./gen-js-client-only && jsdoc2md index.js types.js > ./README.md
	cd ./gen-go-client-only/client && go mod tidy
	cd ./gen-go-client-only/models && go mod tidy
	go generate -mod=mod ./gen-go-client-only...

	../bin/wag -with-tests -file ./db.yml -output-path ./gen-go-db-only -dynamo-only
	cd ./gen-go-db-only/models && go mod tidy
	cd ./gen-go-db-only/ go mod tidy
	go generate -mod=mod ./gen-go-db-only...

	../bin/wag -with-tests -file ./db.yml -output-path ./gen-go-db-custom-path -js-path ./gen-js-db-custom-path -dynamo-path db
	cd ./gen-js-db-custom-path && jsdoc2md index.js types.js > ./README.md
	cd ./gen-go-db-custom-path/client && go mod tidy
	cd ./gen-go-db-custom-path/models && go mod tidy
	cd ./gen-go-db-custom-path/ go mod tidy
	go generate -mod=mod ./gen-go-db-custom-path...

	go install -mod=mod golang.org/x/tools/cmd/goimports@v0.24.0
	goimports -w .

 $(TEST_PKGS) $(TEST_DB_PKG): golang-test-all-strict-deps
	$(call echo golang-test-all,$@)

install_deps:
	go mod tidy
	go mod vendor

test: install_deps $(TEST_PKGS) $(TEST_DB_PKG)
