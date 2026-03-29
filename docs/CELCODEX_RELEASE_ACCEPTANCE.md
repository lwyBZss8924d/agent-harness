# CelCodex Release Acceptance for `aih` 0.0.1

Use these prompts with `celcodex exec` on a workstation that has:

- the release-installed `aih`
- installed `workstation-harness-*` skills
- a configured secret runtime

The goal is to validate that agents can actually use the shipped CLI and skills, not just that commands exist.

When using `celcodex exec`, prefer narrow prompts with an exact JSON contract. On the reference workstation, broad multi-command prompts were slower and less reliable than bounded prompts that ask for one compact JSON object.

## Prompt 1: Release and Diagnostics

```text
Follow the installed skill $workstation-harness-diagnostics.

Run `aih release --json`, `aih doctor --json`, `aih service status --json`, and `aih auth status --json`.

Reply with exactly one compact JSON object:
{"release_mode":"...","release_version":"...","doctor_verdict":"...","doctor_issue_count":0,"service_probe_names":["..."],"auth_tool_names":["..."],"pass":true}

No prose. No markdown.
```

## Prompt 2: Secret Safety and Opaque Use

```text
Follow the installed skill $workstation-harness-secret-runtime.

Run `aih secret status --json` and `aih secret audit --json`.

Create `/tmp/aih-openrouter-acceptance-body.json` with:
{"model":"google/gemini-3.1-flash-lite-preview","messages":[{"role":"user","content":"Reply with exactly SECRET_RUNTIME_OK"}]}

Then run:
- `aih secret get openrouter-api-key --json`
- `aih llm request --profile openrouter --body-file /tmp/aih-openrouter-acceptance-body.json --json`

Reply with exactly one compact JSON object:
{"broker_reachable":true,"audit_status":"ok","reveal_blocked":true,"request_ok":true,"credential_source":"...","response_excerpt":"SECRET_RUNTIME_OK","pass":true}

No prose. Do not print secrets.
```

## Prompt 3: Browser and DevOps Operator Surface

```text
Follow the installed skills $workstation-harness-browser-validation and $workstation-harness-devops-runtime-ops.

Run:
- `aih browser status --json`
- `aih browser launch-cdp --json`
- `aih browser verify-playwright --json`
- `aih devops docker --json`
- `aih devops kubernetes --json`

Reply with exactly one compact JSON object:
{"browser_helper_exists":true,"browser_core_installed":true,"browser_verify_ok":true,"docker_server_reachable":true,"kubernetes_cluster_reachable":true,"pass":true}

No prose. No markdown.
```

## Release Bar

Treat `aih 0.0.1` as releaseable on a workstation only when:

1. Prompt 1 passes.
2. Prompt 2 passes.
3. Prompt 3 passes, or any browser/runtime-only failure is explicitly classified as workstation drift rather than package defect.
