# Release and Install Model

## Goals

- stable, versioned binary installs
- workstation-local dev mode with immediate source reflection
- predictable rollback

## Dev Install

`scripts/install-dev.sh`

Default behavior:

- installs separate `*-go-dev` entrypoints under `~/.local/bin`
- shims call `go run` against this source tree
- source edits are reflected immediately

Optional future behavior:

- replace the legacy `aih` shim only after parity acceptance

## Release Install

`scripts/install-release.sh`

Release mode:

- builds binaries into `build/bin`
- writes `build/release-manifest.json`
- installs versioned binaries under `dist/releases/<version>/bin`
- updates `dist/current` to point at the active installed release root
- writes `dist/install-manifest.json`
- writes stable shims that execute `dist/current/bin/*`
- supports a formal local release check via `./scripts/check-release.sh`

Installed release status should be visible through:

```bash
aih release --json
```

When running an installed release binary, the output should include:

- install mode
- executable path
- release root
- readable install manifest

## Versioning

The toolkit version should be carried in:

- `internal/version/version.go`
- linker-injected build metadata for release builds

Recommended format:

- dev: `0.0.1-dev`
- preview: `0.0.1-beta.1`
- stable: `0.0.1`

## Auto-Update Direction

There are two separate meanings of “auto update”:

1. dev mode
   - source edits immediately affect installed shims because they execute `go run`
2. release mode
   - explicit upgrade command or install script swaps the installed binary set

The toolkit should eventually gain:

- `aih self update`
- channel selection
- local version reporting
- rollback support

Current `0.0.1` target:

- stable local release installs
- versioned local rollback points under `dist/releases`
- machine-readable release/install status
- repeatable local release checklist
- CI entrypoint that runs unit tests and the local release check
