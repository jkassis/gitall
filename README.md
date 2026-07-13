# gitall [![License: CC0-1.0](https://img.shields.io/badge/License-CC0_1.0-lightgrey.svg)](https://spdx.org/licenses/CC0-1.0.html)

`gitall` is a Go CLI for operating on groups of git repositories. It is intended for local repo fleet maintenance: checking whether repositories are clean and in sync, listing branch/origin locations, and updating a Homebrew tap from GitHub releases.

## Commands

- `gitall status <repo-dir>...`: fetches each repo's `origin`, compares local branches with matching remote branches, and reports clean, dirty, out-of-sync, and error states.
- `gitall whatwhere <repo-dir>...`: prints each repo's current branch and `origin` URL.
- `gitall updatetap -b <tap-repo> <repo-dir>...`: updates formula files in a local Homebrew tap for repos that are in sync and have GitHub releases.

## Installation

Apple Silicon quick start:

```sh
brew tap https://github.com/jkassis/dist.brew.pub
brew install gitall
```

Download binaries from this repo's releases or install the macOS version from [the brew tap](https://github.com/jkassis/dist.brew.pub).

Earlier versions invoked the `git` CLI. The current implementation uses [go-git](https://github.com/go-git/go-git) for repository operations. The release pipeline still needs CGO-capable cross-platform builds to produce native artifacts.

The package build currently produces artifacts such as:

- `gitall-x.y.z.aarch64.rpm`
- `gitall-x.y.z.x86_64.rpm`
- `gitall_x.y.z_aarch64.apk`
- `gitall_x.y.z_amd64.deb`
- `gitall_x.y.z_arm64.deb`
- `gitall_x.y.z_x86_64.apk`

## Usage

```sh
gitall status ~/code/*
gitall whatwhere ~/code/*
gitall updatetap -b ~/code/dist.brew.pub ~/code/gitall
```

Example `status` output:

```text
 ✔  charge-controller                        in sync
 ✔  coding-tests                             in sync
 ✔  digits                                   in sync
 ✔  escpos                                   in sync
 ✔  gitall                                   in sync
 ✔  hid                                      in sync
 ✔  nm2                                      in sync
 ✔  w3af                                     in sync
<-> jerrie master                            out of sync with origin
```

## Development

Build the CLI:

```sh
go build -o dist/main ./cmd/
```

Run tests:

```sh
go test ./...
```

Run locally:

```sh
go run ./cmd --help
go run ./cmd status <repo-dir>...
go run ./cmd whatwhere <repo-dir>...
```

The Go-based project helper wraps common build, package, and release workflows:

```sh
go run ./bin/make.go build
go run ./bin/make.go buildx
go run ./bin/make.go package
go run ./bin/make.go release
```

## Documentation

- [System context](docs/CONTEXT.md)
- [Interfaces](docs/INTERFACES.md)
- [Quality attributes](docs/QUALITIES.md)
- [Release runbook](docs/runbooks/RELEASE.md)
