.DEFAULT_GOAL := help
MAKEFLAGS += --silent --no-print-directory

BIN_DIR := ./.bin
SCRIPTS_DIR := ./internal/scripts

# Print Makefile target step description for check.
# Only print 'check' steps this way, and not dependent steps, like 'install'.
# ${1} - step description
define _print_step
	printf -- '------\n%s...\n' "${1}"
endef

## Activate developer environment using devbox. Run `make install/devbox` first If you don't have devbox installed.
activate:
	devbox shell

## Install devbox binary.
install/devbox:
	curl -fsSL https://get.jetpack.io/devbox | bash

## Automatically load devbox environment, requires direnv.
direnv:
	devbox generate direnv

.PHONY: check check/spell check/trailing check/markdown check/format
## Run all checks.
check: check/spell check/trailing check/markdown check/format

## Check spelling, rules are defined in cspell.json.
check/spell:
	$(call _print_step,Verifying spelling)
	yarn --silent cspell --no-progress '**/**'

## Check for trailing whitespaces in any of the projects' files.
check/trailing:
	$(call _print_step,Looking for trailing whitespaces)
	$(SCRIPTS_DIR)/check-trailing-whitespaces.bash

## Check markdown files for potential issues with markdownlint.
check/markdown:
	$(call _print_step,Verifying Markdown files)
	yarn --silent markdownlint '**/*.md' --ignore 'node_modules'

## Verify if the files are formatted.
## You must first commit the changes, otherwise it won't detect the diffs.
check/format:
	$(call _print_step,Checking if files are formatted)
	$(SCRIPTS_DIR)/check-formatting.sh

.PHONY: format format/cspell
## Format files.
format: format/cspell

## Format cspell config file.
format/cspell:
	echo "Formatting cspell.yaml configuration (words list)..."
	yarn --silent format-cspell-config

.PHONY: install
## Install all dev dependencies.
install: install/yarn

## Install JS dependencies with yarn.
install/yarn:
	echo "Installing yarn dependencies..."
	yarn --silent install
	
.PHONY: help
## Print this help message.
help:
	$(SCRIPTS_DIR)/makefile-help.awk $(MAKEFILE_LIST)
