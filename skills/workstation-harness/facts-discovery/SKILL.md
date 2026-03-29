---
name: workstation-harness-facts-discovery
description: Use this skill when an agent needs host facts, discovery snapshots, environment summaries, or generated workstation fact files from `aih`.
---

# Workstation Harness Facts and Discovery

Use this skill for host identity, runtime inventory, and generated fact files.

## Use This Skill When

- you need `aih facts refresh`
- you need `aih facts path`
- you need `aih env summary`
- you need to reason about host identity or current workstation/runtime metadata

## Default Approach

1. Refresh facts if freshness matters:
   - `aih facts refresh --json`
2. Use JSON as source of truth:
   - `aih facts path --json`
   - `aih env summary --json`
3. Treat generated markdown as scan-friendly projection, not canonical machine state.

## References

- `../../../docs/FACTS_DISCOVERY.md`
- `../../../docs/ARCHITECTURE.md`

