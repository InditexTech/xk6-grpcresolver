# Contributing

Thank you for your interest in contributing to this project! We value and appreciate any contributions you can make.
To maintain a collaborative and respectful environment, please consider the following guidelines when contributing to
this project.

## Prerequisites

- Before starting to contribute to the code, you must first sign the
  [Contributor License Agreement (CLA)](https://github.com/InditexTech/foss/blob/main/documents/CLA.pdf).
  Detailed instructions on how to proceed can be found [here](https://github.com/InditexTech/foss/blob/main/CONTRIBUTING.md).

## How to Contribute

1. Open an issue to discuss and gather feedback on the feature or fix you wish to address.
2. Fork the repository and clone it to your local machine.
3. Create a new branch to work on your contribution: `git checkout -b your-branch-name`.
4. Make the necessary changes in your local branch.
5. Ensure that your code follows the established project style and formatting guidelines.
6. Perform testing to ensure your changes do not introduce errors.
7. Make clear and descriptive commits that explain your changes.
8. Push your branch to the remote repository: `git push origin your-branch-name`.
9. Open a pull request describing your changes and linking the corresponding issue.
10. Await comments and discussions on your pull request. Make any necessary modifications based on the received feedback.
11. Once your pull request is approved, your contribution will be merged into the main branch.

## Contribution Guidelines

- All contributors are expected to follow the project's [code of conduct](CODE_OF_CONDUCT.md). Please be respectful and
  considerate towards other contributors.
- Before starting work on a new feature or fix, check existing [issues](../../issues) and [pull requests](../../pulls)
  to avoid duplications and unnecessary discussions.
- If you wish to work on an existing issue, comment on the issue to inform other contributors that you are working on it.
  This will help coordinate efforts and prevent conflicts.
- It is always advisable to discuss and gather feedback from the community before making significant changes to the
  project's structure or architecture.
- Ensure a clean and organized commit history. Divide your changes into logical and descriptive commits.
- Document any new changes or features you add. This will help other contributors and project users understand your work
  and its purpose.
- Be sure to link the corresponding issue in your pull request to maintain proper tracking of contributions.

## Development

Make sure that you have:

- Read the rest of the [`CONTRIBUTING.md`](CONTRIBUTING.md) sections.
- Meet the [prerequisites](#prerequisites).
- [Golang](https://golang.org/doc/install) (version `1.23.5` or higher).
- [GNU Make](https://www.gnu.org/software/make/) (version `4.2.1` or higher).
- [Git](https://git-scm.com/downloads) (version `2.25.1` or higher).

## Testing the application

You can run the tests with the following command:

```sh
make test
```

### Writing tests

We use [stretchr/testify suite package](https://github.com/stretchr/testify#suite-package) for testing when needed. You
can also write regular tests without using the suite package.

## Verifying code integrity

You can verify the code integrity with the following command:
```sh
make verify
```

It's recommended to run this command before pushing your changes to the repository.

## Helpful Resources

- [Project documentation](README.md): Refer to our documentation for more information on the project structure and how
  to contribute.
- [Issues](../../issues): Check open issues and look for opportunities to contribute. Make sure to open an issue before
  starting work on a new feature or fix.

Thank you for your time and contribution! Your work helps to grow and improve this project. If you have any questions,
feel free to reach out to us.
