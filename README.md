# gitall

## Description

A purego CLI for operations on groups of git repos.

## Status

Ironically... because this is pure go, it relies on building against OS native libraries... ie those that implement the os package commands.

Because of that, cross-platform builds are an f'ing nightmatre and we currently only support this for i386.


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
