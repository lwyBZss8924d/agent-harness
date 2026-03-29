# Browser and Validation Harness

## Purpose

This domain brings local browser validation into the AI-first workstation
harness product without turning `aih` into a generic browser automation
framework.

It provides workstation operators with a stable control surface for:

- checking browser automation readiness
- launching a local CDP-enabled browser session
- verifying Playwright-over-CDP connectivity

## Current Commands

```bash
aih browser status --json
aih browser launch-cdp --json
aih browser verify-playwright --json
```

## Design Boundary

This domain should remain:

- local
- explicit
- config-driven
- validation-oriented

It should not become:

- a generic browser test runner
- a hidden background browser daemon
- a product-specific UI script bundle

## Current Implementation

The first Go version supports:

- Chrome binary/app path from config
- CDP port from config or flag
- automation profile dir from config or flag
- Playwright helper path from config
- machine-readable status and verification output

## Next Steps

- improve browser-specific doctor integration
- support richer validation metadata
- keep provider- and project-specific validation logic outside the CLI root
