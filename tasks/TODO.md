# TODO

## Add continuous integration

Status: committed
Priority: high
Owner: unassigned

### Objective

Run the repository's baseline validation automatically for every change.

### Tasks

- [ ] Add GitHub Actions workflow for `go test ./...`.
- [ ] Add `go vet ./...` and a CLI build to the workflow.
- [ ] Keep the workflow compatible with vendored dependencies.
- [ ] Document the CI check in the repository development guidance.

### Acceptance criteria

- Pull requests run tests, vet, and build on a supported Go version.
- A failing check prevents the change from being considered release-ready.
- The workflow is reproducible locally with commands documented in `AGENTS.md`.

## Harden repository status edge cases

Status: committed
Priority: medium
Owner: unassigned

### Objective

Make multi-repository inspection predictable when repositories have unusual or failing states.

### Tasks

- [ ] Define behavior for detached HEADs, missing `origin`, missing matching remote branches, and fetch failures.
- [ ] Add focused tests for those cases and for partial failure across multiple repositories.
- [ ] Ensure per-repository errors remain actionable without aborting the full scan.
- [ ] Update `docs/INTERFACES.md` and `docs/QUALITIES.md` for any clarified behavior.

### Acceptance criteria

- Each edge case has an explicit result and test coverage.
- One failing repository does not suppress results for other requested repositories.
- CLI output and exit behavior are documented.

## Prepare the next release

Status: committed
Priority: medium
Owner: unassigned

### Objective

Make release output and release prerequisites reviewable before invoking the high-impact release flow.

### Tasks

- [ ] Record user-visible changes in `changelog.md`.
- [ ] Verify build and package artifacts from a clean checkout.
- [ ] Confirm GitHub CLI, Docker, and publishing prerequisites from `docs/runbooks/RELEASE.md`.
- [ ] Run `go run ./bin/make.go release` only after the worktree and branch checks pass.

### Acceptance criteria

- The changelog describes the release contents.
- Required artifacts exist and are accounted for before release.
- The release completes through the documented helper without manual recovery.
