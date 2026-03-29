# Parity Baseline

The original rule in this document was the cutover gate from Python to Go.
That cutover is complete on the reference workstation for `aih secret`.

This document now serves two purposes:

1. record the parity boundary that was required for the first cutover
2. define the minimum functional contract future releases must preserve

## Required Functional Parity

### Secrets

- generic secret service must support arbitrary 1Password references and alias-backed lookups, not only LLM API keys
- `aih secret status --json`
- `aih secret audit --json`
- `aih secret get <alias>`
- `aih secret read <op://...>`
- `aih secret sudo -- <command>`
- `aih secret bootstrap-token`
- `aih secret env`
- `aih secret exec -- <command>`
- `aih secret cache <alias>`

### Broker

- LaunchAgent auto-start in Aqua session
- automatic restart after process termination
- broker-backed service-account auth after reboot
- no desktop-app authorization required for fresh SSH sessions

### Canary

The Go implementation must pass the same real-world generic LLM API verification class used by the Python runtime:

- provider/model selection must come from configuration
- no private endpoint or private API conventions may be hardcoded into the toolkit architecture
- the workstation-local production baseline may use a private provider, but the Go implementation must model this as a configurable LLM API target

## Required Operational Parity

- no secret value printed by default
- JSON-first status/audit surfaces
- explicit non-zero exit codes for actionable failures
- same vault: `AIAGENTS`
- same vault UUID: `ui2hktsvdem6rq66s5rlm7gwme`

## Cutover Rule

For a new workstation or a future major migration, do not replace the live
runtime until:

1. direct canary passes
2. reboot validation passes
3. fresh SSH validation passes
4. at least two AI CLI TUI flows pass with no 1Password desktop prompt
