# Interfaces

## CLI

The primary executable is built from `./cmd` and registered in `cmd/main.go`.

```text
gitall [command]
```

### `status`

```text
gitall status [-k <ssh-key-path>] [-p] <repo-dir>...
```

Checks each target directory as a git repository, fetches `origin`, compares local branches to matching `origin/<branch>` refs, then reports:

- repository errors
- repositories in sync
- repositories with staged or unstaged worktree changes
- repositories whose local branches are missing or divergent from origin branches

### `whatwhere`

```text
gitall whatwhere [-k <ssh-key-path>] [-p] <repo-dir>...
```

Prints the current branch and `origin` URL for each target repository.

### `updatetap`

```text
gitall updatetap -b <brew-tap-repo-path> [-k <ssh-key-path>] [-p] [--ghuser] [--ghpass] <repo-dir>...
```

For repositories that pass the same status check used by `status`, reads the latest GitHub release, renders a Homebrew formula, writes it into the local tap repo, stages the formula, prompts for confirmation, commits, and pushes the tap change.

## Credential Inputs

- `-k, --ssh_key_path`: path to the private SSH key, defaulting to `~/.ssh/id_rsa`.
- `-p, --ssh_key_pass_prompt`: prompt for the SSH key passphrase and store it in the `gitall` keyring service.
- `--ghuser`: prompt for the GitHub username and store it in the keyring.
- `--ghpass`: prompt for the GitHub password/token and store it in the keyring.

## Maintainer Tooling

The build and release helper is run with:

```text
go run ./bin/make.go <command>
```

Main helper commands:

- `build`: compile `./cmd` into `dist/main`.
- `buildx`: run cross-platform builds through Docker and the `jkassis/xgo:1.19.5` image.
- `package`: create distro artifacts under `dist/`.
- `release`: verify release preconditions, bump `.semver.yaml`, tag, push, and create a GitHub release.
- `distro`: distribute artifacts to supported distro repositories.

See [Release Runbook](runbooks/RELEASE.md) for operational sequencing.
