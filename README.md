![GitHub Workflow Status (with branch)](https://img.shields.io/github/actions/workflow/status/kairoaraujo/tufie/tests.yml?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/kairoaraujo/tufie)](https://goreportcard.com/report/github.com/kairoaraujo/tufie)
[![Go Reference](https://pkg.go.dev/badge/github.com/kairoaraujo/tufie.svg)](https://pkg.go.dev/github.com/kairoaraujo/tufie#readme-manage-tuf-artifact-repositories)
[![codecov](https://codecov.io/gh/kairoaraujo/tufie/graph/badge.svg?token=V9WL81Q0VA)](https://codecov.io/gh/kairoaraujo/tufie)


# TUFie: An Open Source generic TUF client

TUFie is a TUF (The Update Framework) command-line client. The TUFie simplifies
the client's high-level usage without building a client from scratch.

This client allows simple use cases such as downloading an artifact from an
existent TUF Repository or scripting and CI/CD.

```console
$ tufie download v1.0.3/demo_package-1.0.3.tar.gz

Artifact v1.0.3/demo_package-1.0.3.tar.gz donwload completed.
```

## Install

## HomeBrew

```shell
brew install kairoaraujo/tap/tufie
```

## Winget (Windows)

```shell
winget install tufie
```

Download the Binary

Download from the releases page or use the install script to download the latest release.

[Releases](https://github.com/kairoaraujo/tufie/releases)

```shell
bash <(curl -s https://raw.githubusercontent.com/kairoaraujo/tufie/main/install.sh)
```


## Usage

TUFie has a simple interface

### Download artifacts

```console
$ tufie download v1.0.3/demo_package-1.0.3.tar.gz

Artifact v1.0.3/demo_package-1.0.3.tar.gz donwload completed.
```

### Manage TUF/Artifact repositories

TUFie supports multiple repositories

```console
$ tufie repository -h
Manage TUF repository configurations

Usage:
  tufie repository [REPOSITORY NAME] [flags]
  tufie repository [command]

Available Commands:
  add         Add a new repository
  list        List all repositories
  remove      Remove a repository
  set         Set the default repository
```

#### Add new repository

```console
Usage:
  tufie repository add [flags]

Flags:
      --artifact-hash         add hash prefix to artifact [default: false]
  -a, --artifact-url string   content artifact base URL
  -d, --default               set repository as default
  -h, --help                  help for add
  -m, --metadata-url string   metadata URL
  -n, --name string           repository name
  -r, --root string           trusted Root metadata

$ tufie repository add --default --artifact-url https://rubygems.org --metadata-url https://metadata.rubygems.org --root rubygems-root.json --name rubygems
Config file used for tuf: /Users/kairoaraujo/.tufie/config.yml

Repository 'rubygems' added.
```

#### List repositories

```console
$ tufie repository list
Config file used for tuf: /Users/kairoaraujo/.tufie/config.yml

Default repository: rubygems.org

Repository: rstuf
Artifact Base URL: https://github.com/kairoaraujo/demo-package/releases/download/
Metadata Base URL: http://metadata.dev.rstuf.org

Repository: rubygems
Artifact Base URL: https://rubygems.org
Metadata Base URL: https://metadata.rubygems.org
```

#### Set repository as the default

```console
$ tufie repository set rstuf
Config file used for tuf: /Users/kairoaraujo/.tufie/config.yml

Updated default repository to 'rstuf'.
```

## Contributing

[Fork](https://docs.github.com/en/get-started/quickstart/fork-a-repo) the
repository on [GitHub](https://github.com/kairoaraujo/tufie) and
[clone](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository)
it to your local machine:

```
git clone git@github.com:YOUR-USERNAME/tufie.git
```

Add a [remote](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/configuring-a-remote-for-a-fork) and
regularly [sync](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/syncing-a-fork) to make sure you stay up-to-date with our repository:

```
git remote add upstream https://github.com/kairoaraujo/tufie
git checkout main
git fetch upstream
git merge upstream/main
```

## Requirements

Install [Go](https://go.dev)

```
go mod tidy
```

## Tests

```
make test
```

## Link

Install [golangci-lint](https://golangci-lint.run/usage/install/)

```
make lint
```
