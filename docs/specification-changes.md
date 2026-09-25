# Specification changes

This guide covers changes to OpenSLO object kinds and specification versions.
See [CONTEXT.md](CONTEXT.md) for the terms used here and the relationship between the specification, SDK, and schema manifest.
See [DEVELOPMENT.md](DEVELOPMENT.md) for environment setup and repository checks.

## Adding or modifying an object kind

1. Update the relevant [specification](../website/docs/specification.md) or [draft proposal](../enhancements/), including examples.
2. Update the examples, validation rules, and tests in the [Go SDK repository](https://github.com/OpenSLO/go-sdk).
3. Use the [schema update procedure](../website/README.md#schema-reference) to import the SDK manifest into `website/api.json`.
   Regenerate this manifest through the SDK and import it instead of editing it directly.
4. If the examples require a newer SDK, update the SDK dependency in [the example validator](../website/tools/example-check/go.mod).
5. Run `make check` and `make website/build` from the repository root.

## Adding a specification version

Stable specification versions use names such as `v1`.
Alpha versions use names such as `v1alpha` and `v2alpha`, without a trailing sequence number.
See the [v2 draft](../enhancements/v2alpha.md) for the current proposal.

Coordinate new versions and SDK releases through the [contribution process](../CONTRIBUTING.md#specification-changes).
Import the updated SDK manifest through the same [schema update procedure](../website/README.md#schema-reference).
The website generates schema pages and navigation for every version and object in that manifest.
