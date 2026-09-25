# Examples guide

This directory holds worked OpenSLO examples that explain real use cases in more depth than the
specification. Each example is a directory containing a `README.md` and one or more `.yaml` files.

## Adding an example

1. Create `examples/<name>/` with a `README.md` and the SLO definitions.
2. Answer three questions in that README: what you measure, why you measure it, and how it helps.
   Links to further reading are welcome.
3. Add a bullet for the new directory to the table of contents in `examples/README.md`, with a one
   sentence description.

## Conventions

- A YAML file describes one SLO, starts with a `#` comment summarizing it, and omits the `---`
  document marker.
- Examples use `apiVersion: openslo/v1` and placeholder metric sources with commented field intent.
- Do not rename existing directories or files; published links may point at them.

## Validation

Nothing validates the YAML in this directory. Only `make check/spell` and `make check/trailing` read it,
so check new definitions with the Oslo CLI or the Go SDK yourself. To get automatic validation, put the
example in the specification or in an authored schema page under `website/docs/`, which
`make check/examples` does cover.
