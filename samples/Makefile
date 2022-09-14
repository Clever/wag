include ../golang.mk

TEST_PKGS := $(shell go list -mod=mod ./test...)
TEST_DB_PKG := "./gen-go-db/server/db/dynamodb"

.PHONY: generate test $(TEST_PKGS) $(TEST_DB_PKG) install_deps

# These go:generate don't actually care about dependencies, so we use -mod=mod to override -mod=vendor
# If not for that, it's difficult to bootstrap
generate:
	../bin/wag -file ./swagger.yml -output-path ./gen-go-basic -js-path ./gen-js
	cd ./gen-js && jsdoc2md index.js types.js > ./README.md
	go generate -mod=mod ./gen-go...
	../bin/wag -file ./deprecated.yml -output-path ./gen-go-deprecated -js-path ./gen-js-deprecated
	cd ./gen-js-deprecated && jsdoc2md index.js types.js > ./README.md
	go generate -mod=mod ./gen-go-deprecated...
	../bin/wag -file ./errors.yml -output-path ./gen-go-errors -js-path ./gen-js-errors
	cd ./gen-js-errors && jsdoc2md index.js types.js > ./README.md
	go generate -mod=mod ./gen-go-errors...
	../bin/wag -file ./nils.yml -output-path ./gen-go-nils -js-path ./gen-js-nils
	cd ./gen-js-nils && jsdoc2md index.js types.js > ./README.md
	go generate -mod=mod ./gen-go-nils...
	../bin/wag -file ./blog.yml -output-path ./gen-go-blog -js-path ./gen-js-blog
	cd ./gen-js-blog && jsdoc2md index.js types.js > ./README.md
	go generate -mod=mod ./gen-go-nils...
	../bin/wag -file ./db.yml -output-path ./gen-go-db -js-path ./gen-js-db
	cd ./gen-js-db && jsdoc2md index.js types.js > ./README.md
	go generate -mod=mod ./gen-go-db...
	../bin/wag -file ./swagger.yml -output-path ./gen-go-client-only -js-path ./gen-js-client-only --client-only
	cd ./gen-js-client-only && jsdoc2md index.js types.js > ./README.md
	go generate -mod=mod ./gen-go-client-only...
	../bin/wag -file ./db.yml -output-path ./gen-go-db-only -dynamo-only
	go generate -mod=mod ./gen-go-db-only...
	../bin/wag -file ./db.yml -output-path ./gen-go-db-custom-path -js-path ./gen-js-db-custom-path -dynamo-path db
	cd ./gen-js-db-custom-path && jsdoc2md index.js types.js > ./README.md
	go generate -mod=mod ./gen-go-db-custom-path...

$(TEST_PKGS) $(TEST_DB_PKG): golang-test-all-strict-deps
	$(call golang-test-all,$@)

install_deps:
	go mod vendor

test: install_deps $(TEST_PKGS) $(TEST_DB_PKG)
