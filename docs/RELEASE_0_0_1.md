# Release 0.0.1

## Scope

Release `0.0.1` is the first formal software-delivery milestone for `aih`.

Its scope is intentionally narrow:

- `aih secret`
- the local unattended broker runtime
- configuration-driven alias/profile contracts
- `1password-service-account` and `macos-keychain`
- opaque-use LLM action flow as the first generic action consumer

It is not the full workstation product yet. The rest of the toolkit domains are
the next delivery tranche after the first secret-runtime release line is stable.

## Gate Summary

The release gate is exposed through:

```bash
aih release --json
```

That command should remain the machine-readable control plane for release
progress.

The expected `0.0.1` interpretation is:

- completed gates mean the secret-runtime product contract is in place
- in-progress gates identify what still needs to harden before a formal tagged
  release
- future domains describe what belongs after the first release boundary rather
  than inside it

The release/install side of `0.0.1` should also provide:

- versioned local release roots
- a stable `current` pointer
- install metadata that the CLI can report
- deterministic local rollback points
- a repeatable local release checklist and check script

## Minimum Delivery Contract

`0.0.1` should preserve:

- no plaintext reveal by default
- release-installed packages disable generic unsafe reveal/env-injection surfaces
- JSON-first operational status
- fresh SSH usability with no 1Password desktop prompt
- repo-local config support
- no hardcoded private provider assumptions in the architecture

## Next Tranche After 0.0.1

After `aih secret`, the next formal product tranche is:

- facts/discovery
- doctor/diagnostics
- env/service/auth/browser/facts
- runtime/toolchain/network inspection
- DevOps and unattended workstation operations
