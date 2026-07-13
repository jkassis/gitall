# Done

## Add continuous integration

Completed: 2026-07-12

- Added `.github/workflows/ci.yml` for pull requests and pushes to `master`.
- CI runs tests, vet, and a CLI build with vendored dependencies.
- Documented equivalent local validation commands in `AGENTS.md`.
- Validation passed locally with `go test -mod=vendor ./...`, `go vet -mod=vendor ./...`, and `go build -mod=vendor -o /tmp/gitall ./cmd/`.

## Harden repository status edge cases

Completed: 2026-07-12

- Fetch failures and missing `origin` remotes are reported as per-repository errors.
- Detached HEADs are reported explicitly by `status` and `whatwhere`.
- Local branches without matching origin branches remain sync-needed.
- Added regression coverage for these cases and partial-failure isolation.
- Updated interface and quality documentation.
