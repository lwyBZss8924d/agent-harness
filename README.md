# AIH Toolkit

`aih` is the AI-first workstation harness toolkit for AI-agent-operated owner workstations.

This repository is the **source of truth** for the next-generation Go implementation.
It is intentionally separate from any legacy bespoke runtime that may still exist under `~/.agents`.

## Source vs Runtime

- Development source repo: `<repo-root>`
- Target workstation state and docs: `<workstation-home>/.agents`
- Target workstation legacy bespoke runtime: `<workstation-home>/.agents/retired/...`
- Target workstation installed entrypoints:
  - `/usr/local/bin/aih`
  - `<workstation-home>/.local/bin/*`

## Formal 0.0.1 Boundary

The first formal release boundary is `aih secret`.

`aih secret` in `0.0.1` is expected to provide:

- an unattended workstation-local background service for secret runtime operations
- configuration-driven secret loading for arbitrary declared keys and references
- backend adapters for third-party secret custody systems without exposing plaintext to agents by default
- agent-friendly CLI feedback with stable text and JSON output
- explicit extension points for future backends and action protocols

Important non-goal:

- `aih` is **not** a generic transparent API gateway for arbitrary unknown services
- `aih` should execute declared capability contracts, not silently proxy all traffic

## Design Goals

- Single toolkit repo with versioned source
- Clear separation between source, runtime state, and installed entrypoints
- Dev mode that reflects source changes immediately
- Release mode with versioned binaries and deterministic installation
- AI-first CLI surface with predictable JSON output and explicit failure modes
- pure Unix text/JSON-first agent UX with stable command contracts
- Generic secret custody/loading for any secret stored in an authorized vault
- Broker-backed 1Password service-account access as the first unattended backend
- macOS keychain as a first-class local backend
- LLM API credentials as the most common first-class consumer profile, not the only supported secret class
- config-driven protocols so administrators define service contracts rather than the toolkit hardcoding vendor-specific request formats
- long-lived harness engineering control loops: discovery, diagnosis, verification, recovery, and unattended workstation operations

## Current Status

The Go implementation now owns the reference workstation `aih secret` runtime.

It provides:

- versioned source repo and release/dev install scaffolding
- Go broker daemon and client
- generic secret-service and profile-oriented configuration model
- `1password-service-account` backend
- `macos-keychain` backend
- policy-driven opaque-use profiles and repo-local contract loading
- native `aih secret` operational commands
- real action-based LLM verification/request flows through configured profiles
- native `facts refresh`, `facts path`, and `env summary` commands
- native `doctor`, `service status`, and `auth status` commands
- config-driven auth tool and service probe contracts for diagnostics
- native `browser status`, `browser launch-cdp`, and `browser verify-playwright` commands
- native `devops status` command
- native `devops tools`, `devops docker`, `devops kubernetes`, and `devops registry` subcommands
- deployment/build scripts
- architecture, parity, runtime, and roadmap docs

On the reference workstation, the old Python runtime has been retired to backup-only files.

## Development

```bash
make fmt
make build
make test
make ci
make install-dev
```

The default dev install creates separate `*-go-dev` shims so the Go implementation can evolve without breaking the live Python runtime.

Current regression entrypoints:

- `make test` -> `go test ./...`
- `make ci` -> unit tests plus local release check

Agent skill projection:

- repo skills live under `skills/workstation-harness/`
- local install helper: `./scripts/install-workstation-skills.sh`
- target install destination defaults to `~/.agents/skills`
- workstation package install helper: `./scripts/install-workstation-package.sh`

For remote workstation rollout, add a deployment step rather than treating the repo as the deployed runtime.

Configuration can be layered from:

- global workstation config
- repo-local `.aih/config.json`
- repo-local `.codex/aih/config.json`
- explicit `AIH_CONFIG_FILE`

## Release Strategy

- Dev mode: shim -> `go run ...`
- Release mode: versioned binary build -> installed shim
- Production cutover only after parity acceptance tests pass on a target workstation
- Future release generalization should keep the same product contract while allowing different custody backends and workstation environments

See:

- `docs/ARCHITECTURE.md`
- `docs/DELIVERY_GATE.md`
- `docs/FACTS_DISCOVERY.md`
- `docs/DOCTOR_SERVICE_AUTH.md`
- `docs/BROWSER_VALIDATION.md`
- `docs/DEVOPS_RUNTIME_OPS.md`
- `docs/PRODUCT_SCOPE.md`
- `docs/RELEASE_0_0_1.md`
- `docs/RELEASE_CHECKLIST_0_0_1.md`
- `docs/DEPLOYMENT.md`
- `docs/PARITY.md`
- `docs/RELEASE.md`
- `docs/ROADMAP.md`
- `docs/SECRET_RUNTIME.md`
