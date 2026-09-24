# OpenSLO website

This directory builds the OpenSLO website with MkDocs Material.
The repository's root Devbox environment includes Node.js, Yarn, Python, Ruff, and Go.
See [devbox.json](../devbox.json) for the tool versions.
It installs the website's Python dependencies from `requirements.txt`.

## Development

Run `devbox shell` from the repository root to install and activate the dependencies.
Use the root Makefile:

- `make website/serve` previews the website locally.
- `make website/build` writes the website to `website/site/`.
- `make website/check` lints Python and tests schema generation through MkDocs.
- `make check` runs all repository checks, including YAML example validation.
- `make format` formats Python files.

For one command, use `devbox run -- make website/serve` instead of opening a shell.
Run all commands below from the repository root unless stated otherwise.
Website file paths in this document are relative to `website/`.

Pull requests run these checks and a strict website build in GitHub Actions.
Pushes to `main` publish the checked website to the `gh-pages` branch.
`docs/CNAME` preserves the `openslo.com` domain.

### Publication setup

GitHub Pages must use the `gh-pages` branch and its root directory.
For the initial cutover, retire the publisher in `OpenSLO/openslo.github.io`,
release its `openslo.com` custom domain, and configure that domain in this
repository's Pages settings. Confirm HTTPS and the published routes after
the first deployment. The site remains at `https://openslo.com/`; `website/`
is only the source directory.

## Specification

Edit specification prose and examples directly in `docs/specification.md`.
This page is the source of truth for the specification.
Keep its schema reference links, v2 draft links, and Go SDK validation note
when updating the content.

## Schema reference

The schema reference comes from `api.json`.
The OpenSLO Go SDK combines [govy validation plans](https://github.com/nobl9/govy)
with source comments through [govydoc](https://github.com/nieomylnieja/govydoc).
Its `internal/cmd/objectdoc` generator writes `docs/manifest.json`.

To update the website from an SDK checkout:

1. In the SDK checkout, run `make generate/govydoc` with its required Go toolchain.
2. From the OpenSLO repository root, run `make generate/schema SDK_PATH=/path/to/go-sdk`.
   The default SDK path is `../go-sdk`, beside the OpenSLO checkout.
   This command validates the manifest before it replaces `api.json`.
3. Review the `api.json` diff, then run `make check` and `make website/build`.

Website tests cover rendering, references, links, navigation, and imports.
Field extraction and validation-plan correctness belong to govydoc, govy, and
the SDK generator.

Run `make check/examples` to validate complete YAML examples from the specification
and authored schema pages. This content check uses the SDK's public decoder and
validator. The required Go version and SDK revision are defined in
[`tools/example-check/go.mod`](tools/example-check/go.mod).
Update the SDK pin when the examples must follow a different SDK revision.
CI runs this check before publication.

The MkDocs hook in `main.py` generates pages and navigation for every version and
object in the manifest.
`schema.py` and `templates/` render descriptions, values, examples, and validation
rules.
Builds use the checked-in manifest and require no Go toolchain or SDK checkout.

The SDK generator uses govydoc to register durations as opaque strings.
Their `componentPlans` retain validation for the unit and numeric value.
The website displays these rules inside the duration panel and keeps their
conditions, values, and examples separate from the parent property's rules.
Components do not create child property headings or affect the parent's required
badge.

An authored page under `docs/schema/<version>/<kind>.md` can add examples around
the schema macros.
Otherwise, MkDocs creates the page in memory during the build.
Object filenames use lowercase kind names, such as `alertcondition.md`.

`property-links.json` selects properties that reference a shared definition.
Each version can define:

- `_types`: references keyed by the SDK's `typeInfo.name`.
- `_common`: references keyed by an exact property path, for all object kinds.
- An object kind, such as `SLO`: references keyed by an exact property path.

Object-specific paths take precedence over common paths, then type references.
Each entry contains a `link` and a Jinja `template` for its reference text.
Links are relative to the object page and include the definition's anchor.
For example, v1 `Metadata` points to `../v1.md#metadata`.

Referenced properties have linked headings and a **reference** badge.
They keep their field descriptions and validation rules, but omit their type
descriptions and descendants. The linked page contains that definition.
The generator expands a definition when its reference points to its own heading.
This keeps standalone objects complete while their inline uses link to them.

Inline alert wrappers retain fields such as `conditionRef` and `targetRef`.
Their `spec` properties link to the standalone specifications.

The generator also converts govydoc's `pkg.go.dev` symbol links to website links.
Object names use their generated pages, and shared types use `_types` targets.
Each version can add `_symbols` entries for fields, constants, and other types.
For example, `"SLOObjective.Operator": "slo.md#spec-objectives-items-op"` links
the Go field to its schema property.
Symbol names refer to the version's SDK package.
Use `package/import/path#Symbol` for symbols from another package.
Targets are relative to the version's object directory, as in `_types`.
These entries override inferred targets and do not collapse property definitions.
Links without a website target keep their original URLs.
Code examples remain unchanged.

Do not edit generated HTML under `site/`.
