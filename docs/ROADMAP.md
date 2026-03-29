# Roadmap

## Phase 0

- establish source repo
- define install model
- preserve current Python runtime

## Phase 1

Completed on the reference workstation:

- implement Go `aih` secret surface
- implement generic secret service interfaces and vault reference resolution
- implement Go broker daemon and client
- implement Go generic LLM API verification consumers
- parity tests against current Python runtime
- add policy-driven profiles with opaque-use defaults
- establish first backend set: `1password-service-account`, `macos-keychain`
- retire the Python runtime to backup-only status on the reference workstation

## Phase 2

- formalize `aih secret 0.0.1` product docs, release metadata, and install/update
  workflow
- deliver the first facts/discovery migration with native `facts refresh`,
  `facts path`, and `env summary`
- deliver the first doctor/service/auth migration with native status commands
- deliver the first browser/validation migration with native browser status,
  launch-cdp, and verify-playwright commands
- deliver the first devops/runtime ops migration with native `devops status`
- continue deepening doctor/env/service/auth/browser/facts surfaces in the Go
  product
- start the first devops/runtime ops migration with generic status commands
- define AI-first CLI UX rules for all user-facing command surfaces
- replace legacy shims with Go entrypoints where parity is complete

## Phase 3

- workstation discovery and inventory for server, VM, container, and cluster
  environments
- docker/image/registry management
- k8s context and cluster utilities
- dependency/update inspections
- system/network service helpers
- unattended SSH / DevOps / CI-oriented helper commands

## Phase 4

- package/release automation
- self-update channels
- multi-workstation discovery/adaptation
- additional custody backends for Linux, cloud, and enterprise environments
