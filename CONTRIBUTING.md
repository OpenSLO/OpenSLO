# Contributing to OpenSLO Spec

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

The following is a set of guidelines for contributing to the OpenSLO Spec.
These are mostly guidelines, not rules. Use your best judgment, and feel
free to propose changes to this document in a pull request.

Your pull request will be reviewed by one of the maintainers, and we won't bite.
We encourage and welcome any and all feedback from the community.

## Slack

Use the button `Join our Slack` from the official website [openslo.com](https://openslo.com/).

## Making a pull request

Please make a fork of the repo, and summit a PR from there. More information can
be found [here](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request).
Ensure that checks pass, perform

```bash
make install/checks/spell-and-markdown
```

to install required dependencies to perform the below command

```bash
make run/checks/spell-and-markdown
```

which executes checks for spelling, markdown files and redundant whitespaces. Configuration for them is in the below files:

- [cspell.json](./cspell.json) - missing words can be added to the dictionary (please maintain alphabetical order)

- [.markdownlint.json](./.markdownlint.json) - rules can be adjusted

### Merge Request title

Try to be as more descriptive as you can in your Merge Request title.

## License

Apache 2.0, see [LICENSE](LICENSE).
