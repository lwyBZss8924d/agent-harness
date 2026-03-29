# Product Scope

## Release 0.0.1

The first formal release boundary for `aih` is `aih secret`.

`aih secret` exists to give an AI-agent-operated workstation a safe, unattended,
configuration-driven secret runtime. The formal `0.0.1` scope is:

1. Provide a background local service for workstation secret operations.
2. Load configured secrets for declared contracts without requiring human
   interaction in fresh SSH or fresh AI CLI sessions.
3. Support third-party custody systems such as 1Password service account and
   macOS Keychain behind a consistent runtime and policy layer.
4. Prefer safe delegated use and opaque-use actions over plaintext reveal.
5. Provide AI-friendly command feedback:
   - stable subcommands
   - machine-readable JSON
   - short actionable errors
   - non-release engineering mode may expose explicit break-glass gates, but release-installed workstation packages do not
6. Leave extensibility hooks for future backends and action protocols.

## Non-Goals

`aih` is not intended to become:

- a transparent network gateway for arbitrary unknown services
- a generic plaintext secret dump tool
- a vendor-specific wrapper that hardcodes every third-party request format

Instead, administrators define contracts and profiles, and `aih` executes those
contracts through a stable runtime.

## Contract Model

The toolkit should hardcode only:

- policy semantics
- profile and alias schemas
- runtime orchestration
- protocol adapters that are generic enough to reuse

The toolkit should not hardcode:

- private provider endpoints
- project-specific key names
- service-specific request bodies for every SaaS

Project or workstation operators provide those via configuration.

## Post-Secret Product Surface

After `aih secret`, the formal product evolves into a general-purpose
AI Agent owner workstation self harness-engineering CLI.

Its next domains are:

- workstation discovery and environment inspection
- system, runtime, and network diagnostics
- dev environment dependency management and update checks
- container/image/registry utilities
- orchestration and cluster utilities
- unattended SSH, CI/CD, and DevOps-oriented helper flows
- repository-local and workstation-local contract discovery

The product should remain:

- AI-first
- CLI-first
- Unix text/JSON-first
- configuration-driven
- workstation-focused rather than dashboard-first
