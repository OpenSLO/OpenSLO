.PHONY: build
build:
	go build

.PHONY: install/checks/spell-and-markdow
install/checks/spell-and-markdown:
	yarn

.PHONY: run/checks/spell-and-markdown
run/checks/spell-and-markdown:
	yarn check-trailing-whitespaces
	yarn check-word-lists
	yarn cspell --no-progress '**/**'
	yarn markdownlint --ignore 'node_modules/' '**/*.md'

.PHONY: install/checks/schema-validation
install/checks/schema-validation:
	cd schema-test && go mod download

.PHONY: run/checks/schema-validation
run/checks/schema-validation:
	cd schema-test && go test .
