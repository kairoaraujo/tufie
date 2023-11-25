![GitHub Workflow Status (with branch)](https://img.shields.io/github/actions/workflow/status/kairoaraujo/tufie/tests.yml?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/kairoaraujo/tufie)](https://goreportcard.com/report/github.com/kairoaraujo/tufie)
[![Go Reference](https://pkg.go.dev/badge/github.com/kairoaraujo/tufie.svg)](https://pkg.go.dev/github.com/kairoaraujo/tufie#readme-manage-tuf-artifact-repositories)


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

Adding a new repository

```console
$ tufie repository add -d -a https://rubygems.org -m https://metadata.rubygems.org -r rubygems-root.json -n rubygems
Config file used for tuf: /Users/kairoaraujo/.tufie/config.yml

Repository 'rubygems' added.
```

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

```console
$ tufie repository set rubygems
Config file used for tuf: /Users/kairoaraujo/.tufie/config.yml

Updated default repository to 'rubygems'.
```


