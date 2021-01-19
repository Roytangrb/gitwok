# Contributing to Gitwok

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
1. [Project setup](#project-setup)
1. [Workflows](#workflows)
1. [Filing an issue](#filing-an-issue)
1. [Contributing](#contributing)

## Code of Conduct

This project and its contibutors are expected to uphold the [Go Community Code of Conduct](https://golang.org/conduct). By participating, you are expected to follow these guidelines.

## Project Setup

First, fork and clone this repo to your local `$GOPATH/src/gitwok`, (with go module, the project can also be outside of your `$GOPATH/src/`).

### Run
```
$ go run main.go <subcommands> <flags> [-v]
```
use `--verboes` or `-v` flag for debug verbose ouput

### Build
```
$ go build .
```
outputs executable binary to the current directory

### Install to `GOBIN`
```
$ go install .
```
installs executable binary to go bin directory

### Test
To test all src packages:
```
$ go test -v ./...
```

To view html test coverage report 
```
$ go test ./... -coverprofile=c.out -covermode=atomic && go tool cover -html=c.out
```

One may use [`go-expect`](https://github.com/Netflix/go-expect) for creating test for terminal or console based programs.

### Debug dependency
1. `go mod vendor` checkout dependency modules to `./vendor` directory (git ignored)
1. `go build -mod=vendor` build with src from `./vendor` directory
1. run the executable `./gitwok`

## Workflows

* `ci.yml`: this is the standard go CI workflow to run all tests, build and uploading test coverage report to `codecov.io`. The detailed report can be view [here](https://codecov.io/gh/Roytangrb/gitwok)

* `release.yml`: the CD workflow is triggered whenever a tag starting with `v` is pushed. It creates a release with the version tag. The default archived source code tar zip url is updated to [homebrew tap repo](https://github.com/Roytangrb/homebrew-gitwok). and version will be bumped. And the corresponding Homebrew Fomula changes will be committed.

## Filing an issue

For `Bug Report` and `Feature Request`, please file a [new issue](https://github.com/Roytangrb/gitwok/issues/new/choose) and fill in the issue templates.

Adding `Labels` to raised issues are also useful for identifying issue scope. For example, for issue related to Windows platform, please add `windows` label to the issue.

## Contributing

While submiting a pull requeset, please add along your test if appropriate.