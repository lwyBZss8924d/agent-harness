---
name: workstation-harness-diagnostics
description: Use this skill when an agent needs workstation diagnostics from `aih doctor`, `aih service`, or `aih auth`, and must interpret verdicts, issue sources, and remediation hints.
---

# Workstation Harness Diagnostics

Use this skill for operator-style workstation diagnostics.

## Use This Skill When

- you need `aih doctor --json`
- you need `aih service status --json`
- you need `aih auth status --json`
- you need to interpret workstation issues by source and severity

## Default Approach

1. Start with `aih doctor --json`.
2. Use `summary.verdict` before drilling into raw issues.
3. Use `source` on each issue to decide which domain to inspect next.
4. Use:
   - `aih service status --json`
   - `aih auth status --json`
   for deeper evidence when needed.

## Rules

- Treat `doctor` as the operator verdict surface, not just a probe dump.
- Prefer config-driven suppression or severity override over ad-hoc code changes for workstation-local noise.

## References

- `../../../docs/DOCTOR_SERVICE_AUTH.md`
- `../../../docs/DELIVERY_GATE.md`

