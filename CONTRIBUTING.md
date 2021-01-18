# Contributing to Gitwok

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
1. [Project setup](#project-setup)
1. [Workflows](#workflows)
1. [Testing](#testing)
1. [Filing an issue](#filing-an-issue)
1. [Contributing](#contributing)

## Code of Conduct

This project and its contibutors are expected to uphold the [Go Community Code of Conduct](https://golang.org/conduct). By participating, you are expected to follow these guidelines.

## Project Setup

First, fork and clone this repo to your local `$GOPATH/src/gitwok`

### Run
`go run main.go <flags> <subcommands>`, use `--verbose` for verbose ouput

### Build
`go build .` outputs executable binary to the current directory

### Test
`go test -v ./...`

### Install to `GOBIN`
`go install` install executable binary to go bin directory

### Debug dependency
1. `go mod vendor` checkout dependency modules to `./vendor` directory (git ignored)
1. `go build -mod=vendor` build with src from `./vendor` directory
1. run the executable `./gitwok`

## Workflows

* `ci.yml`: this is the standard go CI workflow to run all tests, build and uploading test coverage report to `codecov.io`. The detailed report can be view [here](https://codecov.io/gh/Roytangrb/gitwok)

* `release.yml`: the CD workflow is triggered whenever a tag starting with `v` is pushed. It builds the executable binary, tar zips and creates a release with the version tag. The zipped asset is uploaded as github release asset. And the corresponding Homebrew Fomula changes will be PR to another [Homebrew tap repo](https://github.com/Roytangrb/homebrew-gitwok).

## Testing

One may use [`go-expect`](https://github.com/Netflix/go-expect) for creating test for terminal or console based programs.

## Filing an issue

For `Bug Report` and `Feature Request`, please file a [new issue](https://github.com/Roytangrb/gitwok/issues/new/choose) and fill in the issue templates.

Adding `Labels` to raised issues are also useful for identifying issue scope. For example, for issue related to Windows platform, please add `windows` label to the issue.

## Contributing

While submiting a pull requeset, please add along your test if appropriate.