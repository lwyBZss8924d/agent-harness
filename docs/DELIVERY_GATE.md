# Delivery Gate

## Purpose

This document projects the original workstation-specific harness requirements
from:

- `<legacy-workstation-docs>/os-dev-environment/*.llms.txt`
- `<legacy-workstation-docs>/harness-cli.llms.txt`
- `<legacy-workstation-docs>/AGENTS.md`

into the formal Go product boundary for `aih-toolkit`.

It also aligns that projection with the core AI-first CLI and harness
engineering principles from:

- `harness-engineering-openai.llm.txt`
- `cli-is-the-new-api-and-mcp.llm.txt`

The point is not to preserve every workstation-specific detail. The point is to
define what must survive as a product contract when converting a bespoke machine
harness into a publishable, reusable AI-first CLI toolkit.

## Product Framing

`aih` is the control plane for AI Agent owner workstation self harness
engineering.

Formal product principles:

- AI-first CLI, not dashboard-first
- Unix text and JSON-first
- explicit command contracts and non-zero exits
- repository-local and workstation-local discoverability
- delegated actions over plaintext secret reveal
- capability-oriented runtime, not transparent proxy/gateway behavior

## Release 0.0.1 Gate: `aih secret`

The first formal release boundary is `aih secret`.

`aih secret` passes the release gate when all of the following are true:

1. It provides a workstation-local unattended background service.
2. It can load declared secrets for configured contracts in fresh SSH and fresh
   AI CLI sessions without human desktop approval.
3. It supports at least:
   - `1password-service-account`
   - `macos-keychain`
4. It enforces agent-safe defaults:
   - no plaintext reveal by default
   - explicit unsafe/admin gates
   - stable JSON status/audit
5. It supports configuration-driven profiles and aliases:
   - workstation-global config
   - repo-local config
   - optional session-local config
6. It can perform opaque-use actions with configured secrets instead of only
   returning raw values.

`aih secret` does **not** need to hardcode every future API provider, request
shape, or key name. It only needs to provide the runtime, policy, and contract
machinery.

## Mapping From The Original `.agents` Harness

### Source: `harness-cli.llms.txt`

Original workstation harness intent:

- machine-level harness CLI
- stable subcommands
- `--json`
- explicit exit codes
- doctor/facts/env/service/auth/secret/browser control surface
- noninteractive secrets access
- browser automation verification

Formal Go projection:

- keep `aih` as the single machine-facing entrypoint
- preserve agent-friendly command ergonomics
- preserve JSON-first operational surfaces
- preserve long-running unattended runtime behavior
- preserve browser/observability/runtime verification as later domains after
  `secret`

### Source: `AGENTS.md`

Original workstation doctrine:

- `aih` is the machine control surface
- AI agents are the primary operators
- docs act as routing and progressive disclosure
- dynamic facts must remain verifiable
- secrets should go through `aih secret`

Formal Go projection:

- `aih` remains the operator-facing CLI root
- repo docs should act as routing and system-of-record context for the agent
- dynamic facts should stay machine-verifiable and command-backed
- no project should bypass the runtime with hardcoded secrets or ad-hoc wrappers

### Source: `os-dev-environment/*.llms.txt`

These files split the machine into reusable operator domains:

- system
- languages/runtimes
- package managers
- devops tools
- services/ports
- cli tools
- environment vars and secrets context
- AI agent tools

Formal Go projection:

- these are not one-off notes; they are the domain model for future `aih`
- each file projects into a future product surface rather than a copied manual

## Post-Secret Product Gate

After `aih secret`, the remaining product boundary is:

- **discovery and facts**
  - host identity
  - session-sensitive facts
  - generated snapshots
- **doctor and diagnostics**
  - workstation health
  - auth/runtime drift
  - stale sockets, proxies, ports, contexts
- **env and service**
  - runtime/toolchain summaries
  - service status
  - port and daemon inspection
- **browser and validation**
  - Chrome/CDP/Playwright checks
  - validation loops for agent work
- **devops and unattended ops**
  - Docker/image/registry helpers
  - k8s/orchestration helpers
  - SSH/CI/CD/operator-friendly actions
- **dependency and tooling hygiene**
  - runtimes
  - package managers
  - upgrade checks
  - drift detection

## What Must Be Generic

The publishable toolkit should generalize:

- the runtime model
- the command ergonomics
- the policy model
- the config schema
- the backend adapter interfaces
- the workstation discovery model

The publishable toolkit should **not** generalize by hardcoding:

- private workstation endpoints
- one customer's service contracts
- arbitrary third-party request formats
- every possible secret item name

Those belong in configuration and repo-local contracts.

## Agent UX Gate

Derived from the CLI and harness references, every formal `aih` domain should
meet these UX rules:

- one obvious command path per capability
- `--help` explains intent, not only syntax
- JSON exists for every non-trivial operational surface
- non-zero exits are actionable
- auth/runtime errors include the next repair step
- output is concise enough for agents to compose
- behavior is stable enough to be treated as a command contract

## Must-Not Boundary

The formal product must not drift into:

- a universal transparent gateway for unknown services
- a generic plaintext secret dump interface
- a pile of workstation-specific glue commands with no contract
- a monolithic manual that agents cannot navigate

The product should remain a reusable harness engineering toolkit, not a custom
wrapper around one workstation's history.
