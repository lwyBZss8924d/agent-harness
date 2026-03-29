---
name: workstation-harness-browser-validation
description: Use this skill when an agent needs local browser harness checks through `aih browser`, including CDP status, launching Chrome for automation, and Playwright-over-CDP verification.
---

# Workstation Harness Browser Validation

Use this skill for local browser automation readiness and validation loops.

## Use This Skill When

- you need `aih browser status`
- you need `aih browser launch-cdp`
- you need `aih browser verify-playwright`

## Default Approach

1. Inspect readiness:
   - `aih browser status --json`
2. If CDP is not up, launch it:
   - `aih browser launch-cdp --json`
3. Verify Playwright connectivity:
   - `aih browser verify-playwright --json`

## Rules

- Prefer the workstation harness commands over ad-hoc browser launch logic.
- Treat browser verification as a local validation surface, not as a general-purpose browser testing framework.

## References

- `../../../docs/BROWSER_VALIDATION.md`

