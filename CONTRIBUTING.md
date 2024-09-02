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

## License

Apache 2.0, see [LICENSE](LICENSE).
