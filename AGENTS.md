# OpenSLO repository guide

OpenSLO is an open, vendor-neutral specification for describing service level objectives (SLOs) in YAML.
The documentation is published at [openslo.com](https://openslo.com) and built from `website/` in this repository.

## Environment and commands

Everything runs from the root `Makefile` inside the environment defined by `devbox.json`
(Node.js 24, Yarn 1.22, Python 3.14, Ruff, Go 1.27).
Run `devbox shell` once to install dependencies, or use `devbox run -- make <target>` for a single command.
Run `make help` to list all targets.

| Command                | Purpose                                             |
| ---------------------- | --------------------------------------------------- |
| `make check`           | All checks: style, website, and examples            |
| `make check/style`     | Spell check, Markdown lint, and trailing whitespace |
| `make check/examples`  | Validates YAML examples against the pinned Go SDK   |
| `make website/check`   | Lints Python and runs the website unit tests        |
| `make website/serve`   | Local preview of the site                           |
| `make website/build`   | Strict build into `website/site/`                   |
| `make test`            | `python -m unittest discover -s website/tests`      |
| `make format`          | `ruff format website`                               |
| `make generate/schema` | Imports a Go SDK manifest into `website/api.json`   |

`make check` must pass before a pull request is reviewed, and both CI workflows run these targets.

## Layout

| Path                 | Contents                                                                       |
| -------------------- | ------------------------------------------------------------------------------ |
| `website/`           | MkDocs site and the schema reference generator. See `website/AGENTS.md`.       |
| `examples/`          | Worked SLO examples with supporting prose. See `examples/AGENTS.md`.           |
| `enhancements/`      | One draft document per future specification version, for example `v2alpha.md`. |
| `glossary/`          | Terminology for newcomers.                                                     |
| `internal/scripts/`  | Helpers used by the `Makefile`.                                                |
| `.github/workflows/` | Style checks, plus the website build and publication.                          |

## Hard rules

- `website/docs/specification.md` is the source of truth for specification prose.
  `README.md` is only an introduction, so do not add specification text there.
- Never edit generated output: `website/api.json` (exported from the Go SDK) and `website/site/` (build output).
  Keep `website/docs/CNAME`, which serves the site on its custom domain.
- Most schema pages under `website/docs/schema/` are generated during the build.
  Only the pages present on disk are authored, and they may only add prose and examples around the macros.
- No file may contain trailing whitespace, and lines in Markdown outside `website/` must stay under 140 characters.
  Long prose belongs in `website/docs/`, where the line length rule is disabled.

## Changing the specification

An object kind change spans this repository, the Go SDK, and the generated reference.
`CONTRIBUTING.md` documents the process in the section "Adding or modifying an object kind":

1. Update `website/docs/specification.md`, or the relevant document in `enhancements/`, including examples.
2. Update the examples, validation rules, and tests in the `OpenSLO/go-sdk` repository.
3. Import the new manifest with `make generate/schema SDK_PATH=/path/to/go-sdk`.
4. Run `make check` and `make website/build`.

A new specification version needs a document in `enhancements/` and a section in the SDK.
Alpha versions are named `v1alpha` and `v2alpha`, without a trailing sequence number.

## Documentation to consult

- `CONTRIBUTING.md` for the contribution process and development environment.
- `website/README.md` for the website, the schema import, and the reference link configuration.
- `devbox.json` for pinned tool versions and `cspell.json` for the spell checker dictionary.
