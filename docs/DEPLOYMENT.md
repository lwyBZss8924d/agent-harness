# Deployment Model

## Source of Truth

The source of truth lives in the local development repository:

- `<repo-root>`

This repository should be versioned, reviewed, and eventually published from here.

## Target Workstation

The target workstation keeps:

- runtime state under `~/.agents`
- installed entrypoints under `/usr/local/bin` and `~/.local/bin`
- logs, sockets, and launchd artifacts under the target user's home directory

The workstation should **not** be treated as the primary source repository.

## Deployment Phases

### Phase 1

- Go repo exists locally
- Python runtime may remain active as parity baseline
- deployment is explicit and selective

### Phase 2

- Go binaries are built locally
- release artifacts are copied to the workstation
- shims on the workstation point to installed Go binaries
- the primary broker can be cut over to the Go runtime after parity validation
- workstation package install can be standardized through `scripts/install-workstation-package.sh`

### Phase 2a

Reference workstation state after cutover:

- primary `aih` entrypoint points to the Go runtime
- primary broker LaunchAgent points to the Go broker binary
- Python runtime files are retained only as timestamped backup artifacts

### Phase 3

- workstation installs are channel-aware
- self-update and rollback are available
- release installs expose machine-readable install metadata

## Current Rule

For a new workstation rollout, do not overwrite the live runtime until the Go
implementation passes parity and reboot validation. The reference workstation
has already completed this cutover for `aih secret`.
