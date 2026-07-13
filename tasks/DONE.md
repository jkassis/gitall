# Done

## Add continuous integration

Completed: 2026-07-12

- Added `.github/workflows/ci.yml` for pull requests and pushes to `master`.
- CI runs tests, vet, and a CLI build with vendored dependencies.
- Documented equivalent local validation commands in `AGENTS.md`.
- Validation passed locally with `go test -mod=vendor ./...`, `go vet -mod=vendor ./...`, and `go build -mod=vendor -o /tmp/gitall ./cmd/`.
