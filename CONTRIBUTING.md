# Contributing to Gitwok

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
1. [Project setup](#project-setup)
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

## Testing

One may use [`go-expect`](https://github.com/Netflix/go-expect) for creating test for terminal or console based programs.

## Filing an issue

For `Bug Report` and `Feature Request`, please file a [new issue](https://github.com/Roytangrb/gitwok/issues/new/choose) and fill in the issue templates.

Adding `Labels` to raised issues are also useful for identifying issue scope. For example, for issue related to Windows platform, please add `windows` label to the issue.

## Contributing

Go test and test coverage report are include in the [CI github workflow](https://github.com/Roytangrb/gitwok/tree/main/.github/workflows). While submiting a pull requeset, please add along your test if appropriate.