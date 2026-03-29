# Architecture

## Repository Placement

The source repository belongs under the developer's `~/dev-space`, not under the workstation's `~/.agents`.

Rationale:

- `~/dev-space` is the correct home for versioned source repositories.
- `~/.agents` is workstation-owned state, docs, generated artifacts, and deployed runtime material.
- Separating source from runtime avoids coupling git history to machine-local mutable state.

Chosen layout:

- Local source repo: `<repo-root>`
- Workstation runtime state/docs: `<workstation-home>/.agents`
- Installed entrypoints:
  - `/usr/local/bin/aih`
  - `<workstation-home>/.local/bin/*`

## Runtime Model

There are three planes:

1. Source plane
   - Go code, docs, packaging scripts, release metadata
2. Install plane
   - built binaries or dev shims exposed on `PATH`
3. Runtime plane
   - `~/.agents` state, launchd plists, sockets, logs, generated facts

## Product Boundary

The first formal product boundary is `aih secret`.

`aih secret` is the runtime layer for:

- unattended secret access on AI-agent-operated workstations
- backend-agnostic custody integration
- policy-driven reveal and usage control
- delegated opaque-use actions

It is intentionally **not** a transparent gateway for arbitrary traffic.
The toolkit should execute explicit contracts and profiles, not silently proxy
unknown services.

## Core Components

### `aih`

Main AI-first CLI.

Formal domains:

- secrets
- doctor
- env
- service
- auth
- browser
- system/runtime/network discovery
- dependency/update checks
- docker/image/registry
- k8s/orchestration
- CI/CD and unattended DevOps helpers

The CLI surface should stay:

- text-first for humans and shells
- JSON-first for agents and automation
- stable in subcommand naming and exit behavior
- explicit about recovery steps when runtime or auth fails

### Secret Service

The toolkit needs a generic secret custody/loading service for any vault content an AI-agent-operated workstation is allowed to use.

Supported categories should include:

- workstation credentials
- API keys
- service tokens
- deploy keys
- registry credentials
- arbitrary vault item fields

LLM API credentials are the most common first-class workload profile, but they are a consumer of the generic secret service rather than the center of the architecture.

The runtime should model secret use through:

- backend adapters
- profile policy
- allowed actions
- opaque delegated use
- layered workstation and repo-local configuration

### Broker

Local same-user UNIX socket broker for unattended 1Password service-account usage.

Responsibilities:

- hold service-account token in memory
- serve same-user local requests
- reload from persistent bootstrap source
- avoid desktop-app authorization dependence in fresh SSH / fresh AI CLI sessions
- broker local material resolution for backends that should not be accessed
  directly from fresh SSH sessions

### 1Password Adapter

The default 1Password integration strategy is:

- official Go SDK
- service-account authentication
- no default dependency on desktop-app interactive authorization

Reference:

- official SDK concepts
- service-account mode

`connect-sdk-go` is **not** the default path because it requires a self-hosted 1Password Connect server and adds unnecessary control-plane complexity for this workstation toolkit.

### LLM API Verification

The toolkit should provide generic LLM API verification driven by configured secret and endpoint profiles.

Important boundary:

- verification is a consumer of the generic secret service
- no private provider should be hardcoded into the architecture
- workstation-specific private providers can be used only through configuration
- the same profile mechanism should generalize to registry, webhook, signing,
  and service-token actions without requiring vendor-specific hardcoding

## Backend Strategy

### First implementations

- `1password-service-account`
- `macos-keychain`

### Planned Linux backends

- `libsecret`
- `gnome-keyring`
- `kwallet`
- `pass`
- `gpg`
- `systemd-creds`

### Planned cloud / enterprise backends

- `1password-connect`
- `vault`
- `aws-secrets-manager`
- `gcp-secret-manager`
- `azure-key-vault`
- `k8s-secret`
- `sops`
- `encrypted-file`

## Dependency Strategy

- Prefer Go stdlib for core CLI/runtime plumbing
- Add external dependencies only when they improve long-term maintainability
- Hide vendor-specific APIs behind internal adapters
- Keep 1Password SDK usage behind `internal/onepassword`
- Keep platform/service orchestration behind `internal/launchd`, `internal/docker`, `internal/k8s`, etc.
- Keep protocol-specific action runners behind reusable internal packages instead
  of hardcoding one-off SaaS workflows into the CLI root

## Install Modes

### Dev mode

- path shims execute `go run ...`
- source changes reflect immediately
- suitable for active workstation-local iteration

### Release mode

- build versioned binaries
- install shims to stable locations
- suitable for production usage and repeatable rollout

## Migration Rule

The original Python runtime served as the parity baseline. The current reference
workstation cutover is now complete for `aih secret`, but future work should
continue to validate new domains against the same product rules:

- unattended workstation-safe behavior
- explicit policy boundaries
- agent-friendly CLI contracts
- configuration-driven integration rather than vendor hardcoding
