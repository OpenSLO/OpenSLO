# Agent instructions

## Required context

- Before implementation, read [README.md](README.md) for project scope and [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.
- Before environment setup or verification, read the [development guide](docs/DEVELOPMENT.md).
  Use its documented check scopes when reporting results.
- Before specification or SDK dependency changes, read the [specification workflow](docs/specification-changes.md).
  Read the relevant [specification](website/docs/specification.md) or [draft proposal](enhancements/) before editing.
- Before website work, read the [website guide](website/README.md).
  For `website/api.json`, `website/site/`, or generated schema pages, read the [schema reference](website/README.md#schema-reference).
  Use that guidance to identify the source files to change.
- Before editing usage examples, read the [example contribution guidelines](examples/README.md).
  Check the [validation coverage](docs/DEVELOPMENT.md#checks-and-formatting) before reporting that examples pass validation.

## Scope

- Apply changes to an SDK checkout only when the user's task includes that repository.
  Otherwise, report required SDK changes as dependencies of the local work.
- Publish the website or change GitHub Pages settings only when the user's request authorizes those actions.

## Maintaining these instructions

Keep shared setup instructions, workflows, and conventions in the linked guides.
Use this file for agent-specific constraints and pointers to those documents.
