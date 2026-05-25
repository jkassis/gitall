# Release Runbook

`gitall` uses the GitHub CLI and project build tooling to build, package, tag, and publish releases.

## Build Workflow

The release tooling is centered on `bin/make.go`:

```sh
go run ./bin/make.go build
go run ./bin/make.go buildx
go run ./bin/make.go package
go run ./bin/make.go release
```

Cross-platform builds use Docker and the `jkassis/xgo:1.19.5` image. Packaging uses `nfpm` to produce Linux and Darwin artifacts under `dist/`.

## Preconditions

- GitHub CLI (`gh`) is installed and authenticated.
- Docker is installed for cross-platform builds.
- The worktree is clean before running the release command.
- The current branch is in sync with `origin/<branch>`.
- Release artifacts are already present under `dist/`.

## Release

```sh
go run ./bin/make.go release
```

The release command:

1. Checks GitHub CLI authentication.
2. Verifies the repository has no uncommitted changes.
3. Verifies the current branch matches `origin/<branch>`.
4. Increments the patch version in `.semver.yaml`.
5. Commits the version bump.
6. Tags and pushes the new release tag.
7. Creates a GitHub release with files from `dist/`.

## Historical Workflow

Older release notes described manually pushing a version tag:

```sh
git tag v7.7.7
git push
git push --tags
```

Prefer the `bin/make.go release` workflow unless intentionally bypassing the project release helper.

