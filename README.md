<p align="center">
  <img alt="GitWok Logo" src="" width="140" height="140" />
  <h3 align="center">GitWok</h3>
  <p align="center">Configurable CLI with conventional commits, changelog, git hooks all in one</p>
</p>

<p>
  <a href="https://github.com/Roytangrb/gitwok/actions">
    <img alt="Actions Status" src="https://github.com/Roytangrb/gitwok/workflows/CI/badge.svg" />
  </a>
  <a href="https://codecov.io/gh/Roytangrb/gitwok">
    <img alt="codecov" src="https://codecov.io/gh/Roytangrb/gitwok/branch/main/graph/badge.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/Roytangrb/gitwok">
    <img alt="goreport" src="https://goreportcard.com/badge/github.com/Roytangrb/gitwok" />
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

- [`add` command](#add-command)
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

### `add` command

The add subcommand prompts for selecting unstaged changes of the current directory to be added for commiting.
```
$ gitwok add
```
![add command capture](docs/images/add.png)

### `commit` command

The `commit` subcommand is used for building the commit message following <a href="https://www.conventionalcommits.org/en/v1.0.0/" target="_blank"><img alt="Conventional Commits" src="https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg" /></a> specification, and execute `git commit -m <msg>`.

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

You may also build the commit message interactively by running :
```
$ gitwok commit
````
You will be prompted for selecting/entering each commit message component.

![commit command capture](docs/images/commit.png)

## Configuration

### Conventional commits

> coming soon

### Conventional changelog

> coming soon

## Reference
* [Conventional Commits 1.0.0](https://www.conventionalcommits.org/en/v1.0.0/)
* [Cobra](https://github.com/spf13/cobra)
* [Survey](https://github.com/AlecAivazis/survey)
* [Carbon](https://carbon.now.sh/)