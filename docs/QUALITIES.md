# Quality Attributes

## Reliability

- Repository checks should isolate failures per target directory so one bad repo does not prevent reporting on the rest.
- Network and host-key errors should be surfaced in the per-repo error list with actionable detail.
- High-impact release and tap update flows should fail before writing remote state when preconditions are not met.

## Safety

- `status` and `whatwhere` inspect local repositories but should not mutate them beyond `status` fetching `origin`.
- `updatetap` has external effects: it writes a neighboring tap repository, creates a commit, and pushes it after confirmation.
- Release helper commands may bump versions, tag, push, publish GitHub releases, and modify generated artifact directories.
- Secrets are kept in the OS keyring service named `gitall`; command-line prompts should avoid echoing secret values.

## Operability

- CLI output is optimized for scanning a fleet of repositories, with stable grouping and colored status markers.
- Build and release operations are centralized in `bin/make.go` so the release runbook can reference one command surface.
- Release artifacts are generated under `build/` and `dist/`; those directories should be treated as generated unless the task is explicitly about packaging or release flow.

## Maintainability

- Command-specific Cobra setup belongs in `CMD<Name>Init` functions.
- Shared auth, prompting, status collection, and output formatting belong in `cmd/common.go` or `cmd/flags.go`.
- Reusable helpers should return explicit errors rather than terminating the process from deep helper code.
- Vendored dependencies should be changed through Go module tooling only when dependency updates are intentional.
