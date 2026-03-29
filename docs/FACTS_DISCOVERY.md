# Facts and Discovery

## Purpose

This domain is the first post-secret product tranche after `0.0.1`.

Its job is to turn workstation identity, runtime inventory, session-sensitive
facts, and generated snapshots into stable CLI surfaces that agents can rely on
without reading workstation-specific manuals first.

## Current Surface

The Go CLI now provides:

```bash
aih facts path --json
aih facts refresh --json
aih env summary --json
```

These commands intentionally focus on a small generic contract:

- generated facts file locations
- host identity
- runtime/toolchain versions
- PATH head and fresh-login PATH head
- SSH session hints
- toolkit/broker context

## Output Model

`facts refresh` writes:

- JSON: `~/.agents/state/facts/host-facts.json`
- Markdown: `~/.agents/state/facts/host-facts.md`

The JSON payload is intended to be the machine-readable source of truth.
The Markdown file is the scan-friendly projection for operators.

## Design Rules

- keep the schema generic across workstations
- do not encode workstation-specific doctor heuristics here
- keep dynamic values command-verifiable
- keep future diagnostics built on top of these facts instead of embedding
  ad-hoc probes directly into every command

## Next Steps

The next step in this domain is to layer:

- doctor diagnostics
- service status
- auth status
- richer inventory and environment scans

on top of the same facts/discovery model instead of rebuilding a separate
workstation-specific contract.
