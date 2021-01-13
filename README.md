<p align="center">
  <img alt="GitWok Logo" src="" width="140" height="140" />
  <h3 align="center">GitWok</h3>
  <p align="center">Configurable CLI with conventional commits, changelog, git hooks all in one</p>
</p>

<p>
  <a href="https://github.com/Roytangrb/gitwok/actions">
    <img alt="Actions Status" src="https://github.com/Roytangrb/gitwok/workflows/Go/badge.svg" />
  </a>
  <a href="https://codecov.io/gh/Roytangrb/gitwok">
    <img alt="codecov" src="https://codecov.io/gh/Roytangrb/gitwok/branch/main/graph/badge.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/Roytangrb/gitwok">
    <img alt="goreport" src="https://goreportcard.com/badge/github.com/Roytangrb/gitwok" />
  </a>
  <a href="https://conventionalcommits.org">
    <img alt="Conventional Commits" src="https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg" />
  </a>
</p>

## Table of Contents
<details>
<summary>Install</summary>

- [Go get](#go-get)
- [Homebrew](#homebrew)

</details>

<details>
<summary>Usage</summary>

- [`commit` command](#commit-command)

</details>

<details>
<summary>Configuration</summary>

- [Conventional commits](#conventional-commits)
- [Conventional changelog](#conventional-changelog)

</details>

## Overview

## Install

### Go get
If you have `Go` setup on your machine, get and install by:
```
$ go get -u github.com/Roytangrb/gitwok
```
`gitwok` executable should be available if `$GOPATH/bin` is already in your `PATH`, otherwise put the binary in one of your `PATH` directories

### Homebrew

> coming soon

## Usage

### `commit` command
The `gitwok commit` command is used for building the commit message following [conventional commit v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) specification, followed by executing `git commit -m <msg>`.

#### `flags` mode
You can build the commit message using flags for subcommand `commit`, example: 
```
$ gitwok commit -t docs -s readme.md -d "commit command usage"
```
which commits with a simple and valid message:
```
$ docs(readme.md): commit command usage
```
You can check all flags by `gitwok commit --help`

#### `interactive` mode
You may also build the commit message using interactively by running 
```
$ gitwok commit
````
with no flag or argument. You will be prompted for selecting/entering each conventional commit message component.

> coming soon

## Configuration

### Conventional commits

> coming soon

### Conventional changelog

> coming soon

## Reference
* [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/)