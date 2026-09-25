# Development

This guide covers local development of the OpenSLO specification and website.
See [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution and review guidelines.
See [CONTEXT.md](CONTEXT.md) for project terminology and the sources behind the schema reference.
For changes to object kinds or specification versions, see the [specification change workflow](specification-changes.md).
Run commands from the repository root unless a step names another checkout.

## Repository layout

| Path | Purpose |
| --- | --- |
| [docs/](../docs/) | Domain context and contributor guides for development and specification changes. |
| [website/docs/specification.md](../website/docs/specification.md) | Specification prose and YAML examples. |
| [enhancements/](../enhancements/) | Draft specification proposals. |
| [examples/](../examples/) | Usage examples with explanations. |
| [website/docs/](../website/docs/) | Website content, including authored schema pages and blog posts. |
| [website/api.json](../website/api.json) | Schema manifest imported from the Go SDK. |
| [website/README.md](../website/README.md) | Website rendering, schema imports, and publication details. |
| [website/tests/](../website/tests/) | Python tests for rendering, navigation, links, and schema imports. |
| [website/tools/example-check/](../website/tools/example-check/) | Go module that validates complete YAML examples with the SDK. |

## Development environment

Install Git and Make before you start.
The [Devbox configuration](../devbox.json) defines the versions of Node.js, Yarn, Python, pip, Ruff, and Go.
JavaScript dependencies are in [package.json](../package.json).
Website Python dependencies are in [website/requirements.txt](../website/requirements.txt).

### Devbox

1. If Devbox is not installed, run:

   ```sh
   make install/devbox
   ```

   This target downloads and runs the Devbox installer.

2. Start the development environment:

   ```sh
   devbox shell
   ```

   `make activate` runs the same command.
   The shell hook activates a Python virtual environment and runs `make install` to install JavaScript and Python dependencies.

To run a single command without opening a shell, use:

```sh
devbox run -- make check
```

For automatic environment activation, install direnv and complete these steps:

1. Generate `.envrc`:

   ```sh
   make direnv
   ```

2. Review `.envrc`, then allow direnv to load it:

   ```sh
   direnv allow
   ```

### Manual setup

1. Install the tools at the versions in [devbox.json](../devbox.json), including Ruff and Yarn Classic.
   The example validator's minimum Go version is in [its module definition](../website/tools/example-check/go.mod).

2. Create a Python virtual environment:

   ```sh
   python3 -m venv .venv
   ```

3. Activate the virtual environment:

   ```sh
   . .venv/bin/activate
   ```

4. Install JavaScript and Python dependencies:

   ```sh
   make install
   ```

`make install` does not install the toolchain or Ruff.
Activate the virtual environment again when you open a new shell.

## Checks and formatting

Run these commands before submitting a pull request:

```sh
make check
make website/build
```

`make check` runs style checks, Python linting, website tests, and YAML example validation.
`make website/build` runs a strict MkDocs build.
Together, these commands cover the checks in the [style workflow](../.github/workflows/checks.yml)
and [website workflow](../.github/workflows/build-publish-website.yml).

Use individual targets for a narrower check:

| Command | Purpose |
| --- | --- |
| `make check/style` | Run spelling, trailing whitespace, and Markdown checks. |
| `make check/spell` | Check spelling with [cspell.json](../cspell.json). |
| `make check/trailing` | Check files tracked in the current commit for trailing spaces. |
| `make check/markdown` | Lint repository and website Markdown with their respective configurations. |
| `make check/python` | Lint website Python files with Ruff. |
| `make test` | Run Python tests for website rendering, navigation, links, and schema imports. |
| `make website/check` | Run Python linting and website tests. |
| `make check/examples` | Validate complete YAML examples with the pinned Go SDK. |
| `make format` | Format Python files with Ruff. |

`make format` only formats Python files.
Correct Markdown, spelling, and whitespace errors separately.
Markdown rules are in [.markdownlint.json](../.markdownlint.json) and [website/.markdownlint.json](../website/.markdownlint.json).
For a valid technical term that fails the spelling check, update the appropriate dictionary in `cspell.json`.

`make check/examples` reads YAML code blocks in `website/docs/specification.md` and `website/docs/schema/*/*.md`.
It skips blocks without `apiVersion:` and schematic blocks that contain `name: string`.
It does not validate files under `examples/` or `enhancements/`.
The [Go module](../website/tools/example-check/go.mod) pins the SDK revision for this check.

Run `make help` to list all targets in the [Makefile](../Makefile).

## Website development

Start the local preview:

```sh
make website/serve
```

`make website/build` writes generated HTML to `website/site/`.
Do not edit that output.
Website builds use the imported schema manifest and need no Go toolchain or SDK checkout.
The YAML example check requires Go.

See the [website guide](../website/README.md) for schema rendering, property references, and publication details.
