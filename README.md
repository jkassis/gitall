# gitall

## Description

A CLI to perform git operations on multiple repositories. Right now this generates a one-line sync status for multiple repos.

## Usage

>
```
> [I] jkassis@Jeremys-MacBook-Pro ~/code> gitall [\<repo-dir\>]+
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
