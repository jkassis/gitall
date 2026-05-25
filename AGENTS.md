# AGENTS.md

## Project Orientation

- `gitall` is a Go CLI for operating on groups of git repositories.
- The main user-facing binary lives in `cmd/` and is built from `./cmd`.
- The CLI uses Cobra for commands, Viper for flags/config binding, go-git for repository operations, go-github for GitHub API access, and keyring for local secret storage.
- This repo vendors dependencies in `vendor/` and currently has generated release/build artifacts under `build/` and `dist/`.

## Important Files

- `README.md`: user-facing description, install notes, and command examples.
- `cmd/main.go`: root Cobra command and command registration.
- `cmd/common.go`: shared git status logic, output formatting, SSH auth setup, and GitHub client creation.
- `cmd/flags.go`: interactive prompts, keyring-backed credentials, and common Cobra/Viper flags.
- `cmd/cmdStatus.go`: `gitall status`.
- `cmd/cmdWhatWhere.go`: `gitall whatwhere`.
- `cmd/cmdUpdateTap.go` and `cmd/brewFormula.go`: Homebrew tap update flow and formula rendering.
- `bin/make.go`: Go-based build/release/package helper invoked with `go run ./bin/make.go <command>`.
- `bin/release.bash` and `bin/buildx.bash`: release/build shell helpers.
- `docs/runbooks/RELEASE.md`: release workflow notes.
- `tasks/TODO.md`: committed task index.
- `tasks/IDEAS.md`: captured but uncommitted future work.
- `tasks/DONE.md`: completed task index.
- `.semver.yaml`: release version state used by release tooling.

## First-Turn Checklist

1. Read this file, then inspect `tasks/TODO.md` and `tasks/IDEAS.md`.
2. Check `git status --short` before editing. The repo may contain local generated files or user work.
3. Review recent history with `git log --oneline -n 8` for release/version churn.
4. Prefer reading first-party files in `cmd/`, `bin/`, `README.md`, `docs/`, and `tasks/` before scanning `vendor/`.

## Build And Test Commands

- Build the CLI:
  ```sh
  go build -o dist/main ./cmd/
  ```
- Equivalent project helper:
  ```sh
  go run ./bin/make.go build
  ```
- Run tests:
  ```sh
  go test ./...
  ```
- Run the local CLI without installing:
  ```sh
  go run ./cmd --help
  go run ./cmd status <repo-dir>...
  go run ./cmd whatwhere <repo-dir>...
  ```
- Cross-platform packaging and release commands in `bin/make.go` may require Docker, GitHub CLI, `dpkg-scanpackages`, adjacent distro repositories, and a clean/in-sync git branch.

## Development Guidance

- Keep command registration in `cmd/main.go` small; put command-specific setup in `CMD<Name>Init` functions.
- Follow the existing command pattern: create a fresh `viper.New()`, bind flags in the command init function, and call a command-specific implementation from Cobra `Run`.
- Keep shared auth, prompting, formatting, and repository traversal logic in `cmd/common.go` or `cmd/flags.go` unless a narrower file is clearly better.
- Avoid adding dependencies unless they materially simplify CLI behavior; this is a small utility with vendored dependencies.
- Prefer explicit errors returned from helper functions. Existing command entrypoints sometimes use `log.Fatal`; do not spread fatal exits deeper into reusable logic.
- Be careful with ANSI color constants and Unicode output in `StatiPrint`; these are part of the CLI user experience.
- Do not mutate release versioning, `dist/`, `build/`, or adjacent distro repos unless the task is explicitly about release/package flow.

## Secrets And External Effects

- SSH key passphrases and GitHub credentials are stored through the local keyring service named `gitall`.
- `status` and `updatetap` fetch from remotes using go-git and SSH auth; they can touch the network.
- `updatetap` writes Homebrew formula files to the user-provided tap repo, creates a commit, and pushes it.
- `bin/make.go release` checks git cleanliness, bumps `.semver.yaml`, commits, tags, pushes, and creates a GitHub release.
- Treat commands that push, tag, release, write neighboring repos, or prompt for credentials as high-impact operations; do not run them casually during validation.

## Generated And Vendored Content

- Avoid broad edits in `vendor/`; update dependencies through Go module tooling only when requested.
- `build/` and `dist/` are generated artifact locations. They may be present locally and should not be reformatted or cleaned unless the task requires it.
- If `go mod tidy` changes `go.sum` or `vendor/`, inspect the diff carefully and explain why dependency metadata changed.

## Git Hygiene

- Preserve user changes. Do not reset, checkout, clean, or remove generated artifacts without an explicit request.
- Keep commits focused. For code changes, include any relevant tests or explain why tests were not added.
- If asked to commit, stage only files related to the task and use a concise non-interactive commit message.
