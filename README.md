# gitall[![License: CC0-1.0](https://img.shields.io/badge/License-CC0_1.0-lightgrey.svg)](https://spdx.org/licenses/CC0-1.0.html)

## Description

A purego CLI for operations on groups of git repos.

## Distributions
TL/DR: You can install the mac version from [my brew tap](https://github.com/jkassis/dist.brew.pub).

Earlier versions of this utility called os.exec to invoke the git CLI to run git commands. This version uses [go-git](https://github.com/go-git/go-git) to run git commands. Ironically... this "pure go" approach requires [CGO](https://pkg.go.dev/cmd/cgo) and a cross-platform build pipeline to produce native executables.

This cross-platform build chain currently produces... 

• gitall-x.y.z.aarch64.rpm
• gitall-x.y.z.x86_64.rpm
• gitall_x.y.z_aarch64.apk
• gitall_x.y.z_amd64.deb
• gitall_x.y.z_arm64.deb
• gitall_x.y.z_x86_64.apk

Download binaries from this repo's releases or install the mac version from [my brew tap](https://github.com/jkassis/dist.brew.pub).


## Usage

```
A CLI for operations on groups of git repos.

Usage:
  gitall [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  status      Get the status for multiple git repos

Flags:
  -h, --help   help for gitall

Use "gitall [command] --help" for more information about a command.
```

## eg

```
[I] jkassis@Jeremys-MacBook-Pro ~/code> gitall *                                                                                                                                                                    07.08 10:38
 ✔  charge-controller                        in sync (git@github.com:jkassis/charge-controller.git)
 ✔  coding-tests                             in sync (git@github.com:jkassis/coding-tests.git)
 ✔  digits                                   in sync (git@github.com:jkassis/digits.git)
 ✔  escpos                                   in sync (git@github.com:jkassis/escpos.git)
 ✔  gitall                                   in sync (git@github.com:jkassis/gitall.git)
 ✔  hid                                      in sync (git@github.com:jkassis/hid.git)
 ✔  nm2                                      in sync (git@github.com:jkassis/nm2.git)
 ✔  w3af                                     in sync (https://github.com/andresriancho/w3af.git)
<-> jerrie                                   out of sync (git@github.com:jkassis/jerrie.git)
<-> live.shinetribe.media                    out of sync (git@github.com:jkassis/live.shinetribe.media.git)
<-> merchie                                  out of sync (git@github.com:jkassis/merchie.git)
```

## Installation

>
```
> [I] jkassis@Jeremys-MacBook-Pro ~/code>  tap jkassis/keg
> [I] jkassis@Jeremys-MacBook-Pro ~/code>  brew install gitall
```
