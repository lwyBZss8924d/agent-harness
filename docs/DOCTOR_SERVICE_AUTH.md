# Doctor, Service, and Auth

## Purpose

This domain is the first generic diagnostics layer built on top of
facts/discovery.

It translates the old bespoke workstation harness ideas into reusable product
surfaces:

- `aih doctor --json`
- `aih service status --json`
- `aih auth status --json`

## Design Boundary

These commands should remain generic.

They should report:

- workstation/runtime drift
- missing or stale generated facts
- broker/runtime availability
- listening TCP services
- common auth/config markers for AI-facing CLIs

They should **not** hardcode one workstation's historical operational quirks as
permanent product rules.

Instead, this domain should increasingly use configuration-driven contracts:

- auth tool registry
- service probe registry
- generic issue severity rules

## Current State

The first Go version now provides:

- doctor report with machine-readable issues
- issue source attribution (`facts`, `runtime`, `browser`, `service`, `auth`, `devops`)
- severity override hooks by issue kind
- issue ignore hooks by kind and source
- listening TCP socket inventory
- auth/config markers from a configurable auth tool registry
- port probes from a configurable service probe registry
- operator-facing verdict summary with issue counts by severity and source
- browser checks can be enabled or disabled for doctor independently of the browser CLI itself

## Next Steps

The next iterations should improve:

- generic issue severity rules
- richer network/runtime/service probes
- cleaner path-drift heuristics
- stronger integration between facts, doctor, and service status
- configurable ignore/override policy so workstation-local noise can be tamed without code edits
