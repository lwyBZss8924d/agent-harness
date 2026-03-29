---
name: workstation-harness-secret-runtime
description: Use this skill when an agent needs to operate `aih secret` safely on an AI-agent-operated workstation, including secret health checks, alias/profile usage, broker/runtime repair context, and non-default unsafe surfaces.
---

# Workstation Harness Secret Runtime

Use this skill for workstation secret operations through `aih`.

## Use This Skill When

- you need `aih secret status` or `aih secret audit`
- you need to inspect alias/profile-driven secret behavior
- you need to use configured secrets through opaque-use flows
- you need to reason about broker/runtime health before using secrets

## Default Approach

1. Start with:
   - `aih secret status --json`
   - `aih secret audit --json`
2. Prefer profile/action-based usage over plaintext reveal.
3. Treat `aih secret get/read/env/exec` as policy-sensitive surfaces.
4. If a workflow can use `aih llm ...` or another delegated action, prefer that.

## High-Value Commands

```bash
aih secret status --json
aih secret audit --json
aih secret list --json
aih secret get <alias> --json
aih secret cache <alias> --json
aih secret sudo -- <command>
aih llm verify --profile <name> --json
aih llm request --profile <name> --body-file <file> --json
```

## Rules

- Prefer `--json`.
- Do not assume a secret should be revealed just because it exists.
- If the runtime is unhealthy, report the broker/runtime failure before trying workaround behavior.
- Keep workstation-global and repo-local contracts distinct.

## References

- `../../../docs/SECRET_RUNTIME.md`
- `../../../docs/PRODUCT_SCOPE.md`
- `../../../docs/DELIVERY_GATE.md`

