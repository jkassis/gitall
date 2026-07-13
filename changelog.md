[Unreleased]

### Added

- Added GitHub Actions validation for tests, vet, and CLI builds.
- Added system context, interface, and quality documentation.

### Changed

- `status` now reports fetch failures, missing `origin` remotes, and detached HEADs per repository.
- `whatwhere` now reports detached HEADs without assuming a branch name.
- Added regression coverage for repository status edge cases and partial failures.
