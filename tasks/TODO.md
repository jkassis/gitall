# TODO

## Prepare the next release

Status: committed
Priority: medium
Owner: unassigned

### Objective

Make release output and release prerequisites reviewable before invoking the high-impact release flow.

### Tasks

- [x] Record user-visible changes in `changelog.md`.
- [ ] Verify build and package artifacts from a clean checkout.
- [ ] Confirm GitHub CLI, Docker, and publishing prerequisites from `docs/runbooks/RELEASE.md`.
- [ ] Run `go run ./bin/make.go release` only after the worktree and branch checks pass.

### Acceptance criteria

- The changelog describes the release contents.
- Required artifacts exist and are accounted for before release.
- The release completes through the documented helper without manual recovery.
