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

Fork the repository and submit a pull request.
See GitHub's [pull request guide](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request).

This project uses [Devbox](https://github.com/jetify-com/devbox) 0.18.3 for its development environment.
Run `devbox shell` from the repository root to install and activate the dependencies.
You can also install the required dependencies manually.

All the development commands are provided via `Makefile`.
You can run `make help` to see the list of available commands.

Run all checks from the repository root:

```sh
make check
```

To fix formatting errors, run:

```sh
make format
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

## Website

See [website development](website/README.md) to preview the site, update schema
documentation, and edit the specification.

## License

Apache 2.0, see [LICENSE](LICENSE).
