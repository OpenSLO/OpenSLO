# Contributing to OpenSLO Spec

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

The following is a set of guidelines for contributing to the OpenSLO Spec.
These are mostly suggestions, not strict rules. Use your best judgment, and feel
free to propose changes to this document in a pull request.

Your pull requests will be reviewed by one of the maintainers, and we won't bite.
We encourage and welcome any and all feedback from the community.

## Slack

Use the button `Join our Slack` from the official website [openslo.com](https://openslo.com/).

## Making a pull request

Please make a fork of the repo, and summit a PR from there. More information can
be found [here](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request).

This project utilizes [devbox](https://github.com/jetify-com/devbox) in order
to provide a consistent and reliable development environment.
You can however install the required dependencies manually.

All the development commands are provided via `Makefile`.
You can run `make help` to see the list of available commands.

Checks which are run as part of the CI pipeline can be run locally wth:

```sh
make check
```

If you see formatting or code generation errors you can fix them with:

```sh
make format && make generate
```

If you have devbox installed, you can initialize the environment with:

```sh
make activate
```

Devbox can be easily installed with:

```sh
make install/devbox
```

Furthermore, you can utilize devbox's direnv integration to automatically
activate the environment when you enter the project's directory, to do so run:

```sh
make direnv
```

It will generate an `.envrc` file which is scanned by direnv when you
enter or leave the directory. You might need to run `direnv allow` in order
to whitelist the project's `.envrc` file.

### Merge Request title

Try to be as descriptive as you can in your Merge Request title.

## Adding or modifying new object kind

1. Update main [README.md](./README.md) file with the description of the object.
2. Update SDK code.
    1. Add or modify examples.
    2. Add or modify validation rules.
    3. Add or modify unit tests.

## Adding new version

*TO BE DISCUSSED AND CHANGED*

OpenSLO follows a similar version lifecycle to the k8s API.

- Official versions can be described by the following regular expression: `v[0-9]`, e.g. `v1`.
- Versions still in active development or in experimentation phase can be described by the following regular expression: `v[0-9](alpha|beta)[0-9]`, e.g. `v2alpha1`, `v3beta2`.

Alpha versions are promoted to stable versions following a community consensus.
Once a version is promoted:

- Tag should be created in the [main OpenSLO repository](https://github.com/OpenSLO/OpenSLO).
- The tag should list ALL the changes between the new version and its predecessor.
- The tag should follow [semantic versioning](https://semver.org/), this is of
  particular importance to the OpenSLO SDK.

## License

Apache 2.0, see [LICENSE](LICENSE).
