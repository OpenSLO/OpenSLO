.PHONY: build
build:
	go build

.PHONY: install/checks/spell-and-markdown
install/checks/spell-and-markdown:
	yarn

.PHONY: run/checks/spell-and-markdown
run/checks/spell-and-markdown:
	yarn check-trailing-whitespaces
	yarn check-word-lists
	yarn cspell --no-progress '**/**'
	yarn markdownlint --ignore 'node_modules/' '**/*.md'
