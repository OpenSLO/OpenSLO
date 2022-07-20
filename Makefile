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

.PHONY: install/checks/schema-validation
install/checks/schema-validation:
	cd schema-test && go mod download

.PHONY: run/checks/schema-validation
run/checks/schema-validation:
	cd schema-test && go clean -testcache && go test .

.PHONY: gen/v1/go
gen/v1/go:
	cd schemas/v1 && quicktype -l go -t OpenSLO -o openslo.go --package v1 --just-types-and-package -s schema openslo.schema.json $$(find kinds -name '*.schema.json' | xargs -I {} echo -n ' -S {}')  $$(find parts -name '*.schema.json' | xargs -I {} echo -n ' -S {}')
