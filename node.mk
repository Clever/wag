# This is the default Clever Node Makefile.
# It is stored in the dev-handbook repo, github.com/Clever/dev-handbook
# Please do not alter this file directly.
NODE_MK_VERSION := 0.3.3
SHELL := /bin/bash

# This block checks and confirms that the proper node version is installed.
# arg1: node version. e.g. v6
define node-version-check
_ := $(if \
	$(shell node -v | grep $(1)), \
	@echo "", \
	$(error "Node $(1) is required, use nvm to install / use it"))
endef

# This block checks and confirms the number of coffeescript files in the repo. The function must be
# called inside a target or it will result in a syntax error.
# arg1: number of coffeescript files. e.g. 10
define node-coffeescript-file-count-check
@echo -e "\nChecking count of coffeescript files"
@git ls-files '*.coffee' | wc -l | tr -d ' ' > /tmp/node-coffee-file-count
@if [ "`cat /tmp/node-coffee-file-count`" -eq "$(1)" ]; then \
	echo -e "\033[0;32m✓ No change in file count.\033[0m\n"; \
elif [ "`cat /tmp/node-coffee-file-count`" -gt "$(1)" ]; then \
	echo -e "\033[0;31m✖ Found new coffeescript file(s). All new modules should be written in ES6.\033[0m\n"; \
	exit 1; \
else \
	echo -e "\033[0;31m✖ Congrats! You have reduced the file count to `cat /tmp/node-coffee-file-count`. Please lower the expected count in the Makefile.\033[0m\n"; \
	exit 1; \
fi
endef

# node-guarded-lint runs a lint with an expected number of errors (and exits
# non-zero if there number of errors differs).
# arg1: linter name (for human readability)
# arg2: path to linter executable
# arg3: files to lint
# arg4: expected number of problems
define node-guarded-lint
@$(2) $(3) > /tmp/lint-output.txt || true
@# look for eslint/tslint-esque problem count, fall back to counting number of ✖'s
@sed -n 's/^✖ \(.*\) problems.*/\1/p' /tmp/lint-output.txt > /tmp/lint-problem-count
@if [[ `cat /tmp/lint-problem-count` = "" ]]; then  \
	cat /tmp/lint-output.txt | sed -e 's/\(.\)/\1\n/g' | grep '✖' | wc -l > /tmp/lint-problem-count;  \
fi
@if [ "`cat /tmp/lint-problem-count`" -gt "$(4)" ]; then  \
	cat /tmp/lint-output.txt;  \
	echo -e "\033[0;31m✖  Added $(1) errors. Please don't introduce new lint errors! If you are decaffeinating a file, you may increase the $(1) problem count in the Makefile.\033[0m"; \
	echo "  Expected: $(4)"; \
	echo -e "  Actual: \033[0;31m`cat /tmp/lint-problem-count`\033[0m\n"; \
	exit 1;  \
elif [ "`cat /tmp/lint-problem-count`" -lt "$(4)" ]; then  \
	echo -e "\033[0;31m✖  Congrats! You have decreased the $(1) problem count. Please lower the count in the Makefile.\033[0m"; \
	echo "  Expected: $(4)"; \
	echo -e "  Actual: \033[0;32m`cat /tmp/lint-problem-count`\033[0m\n"; \
	exit 1;  \
else  \
	echo "✓  No new $(1) errors found.";  \
fi
endef

# node-coffee-edit-check exits non-zero if this branch adds any lines to
# coffeescript files. NOTE: Uses the merge-base, as otherwise we would get
# changes in master that were merged after this branch diverged from master.
.PHONY: node-coffee-edit-check
node-coffee-edit-check:
	@coffee_edits=$$(git diff $(git merge-base master HEAD) --stat | grep '\.coffee.*+');  \
	if [[ $$? != 1 ]]; then  \
		echo "You're not allowed to edit coffeescript files. Please decaffeinate the following files: ";  \
		echo "$$coffee_edits" | awk '{ print "-", $$1 }';  \
		echo;  \
		echo "Don't know how to do this? Read https://clever.atlassian.net/wiki/display/ENG/ES6+at+Clever#ES6atClever-decaffeinate.";  \
		echo;  \
		exit 1;  \
	fi

# node-update-makefile downloads latest version of node.mk
.PHONY: node-update-makefile
node-update-makefile:
	@wget https://raw.githubusercontent.com/Clever/dev-handbook/master/make/node.mk -O /tmp/node.mk 2>/dev/null
	@if ! grep -q $(NODE_MK_VERSION) /tmp/node.mk; then cp /tmp/node.mk node.mk && echo "node.mk updated"; else echo "node.mk is up-to-date"; fi

##########################
## Linting / Formatting ##
##########################
# Match .js, .ts, .jsx, and .tsx files. Ignore directories contaning installed dependencies or auto-generated files.
NODE_ALL_FILES := $(shell find . -type f \( -name "*.ts" -or -name "*.tsx" -name "*.js" -or -name "*.jsx" \) -not -path "./node_modules/*" -a -not -path "./gen-js/*" -a -not -path "./build/*" -a -not -path "./__build/*")

.PHONY: node-format
node-format: node-prettier-format node-tslint-fix

.PHONY: node-lint
node-lint: node-prettier-lint node-tslint

.PHONY: node-prettier-lint
node-prettier-lint:
	@echo "Running prettier lint..."
	@./node_modules/.bin/prettier -l $(NODE_ALL_FILES)

.PHONY: node-tslint
node-tslint:
	@echo "Running tslint..."
	@./node_modules/.bin/tslint --project tsconfig.json -t verbose $(NODE_ALL_FILES)

.PHONY: node-prettier-format
node-prettier-format:
	@echo "Running prettier..."
	@./node_modules/.bin/prettier --write $(NODE_ALL_FILES)

.PHONY: node-tslint-fix
node-tslint-fix:
	@echo "Running tslint fix..."
	@./node_modules/.bin/tslint --fix --project tsconfig.json -t verbose $(NODE_ALL_FILES)
