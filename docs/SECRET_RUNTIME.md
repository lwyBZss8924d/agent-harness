# Secret Runtime Model

## Design Principle

The toolkit follows an **opaque-use** model inspired by wallet/signing systems such as MetaMask:

- the backend secret store is responsible for custody, persistence, backup, sync, and scoped access
- `aih` is responsible for runtime policy, capability control, injection, and delegated use
- agent-facing flows should prefer **actions** over plaintext secret reveal

The toolkit is **not** a generic plaintext secret dumping interface.
It is also **not** intended to become a transparent universal API gateway for
arbitrary unknown service traffic.

## Responsibility Split

### Secret backend

Backends are responsible for:

- durable storage
- sync/backup/recovery
- identity and backend-native access control
- retrieving raw secret material when permitted

### `aih` secret runtime

The runtime is responsible for:

- profile and alias registry
- reveal policy
- usage policy
- allowed actions
- broker/cache/runtime state
- audit and health checks
- delegated opaque-use actions

## Policies

### Reveal policy

- `never`
- `admin_only`
- `allowed`

### Usage policy

- `opaque_use`
- `inject_env`
- `reveal`

### Allowed actions

Examples:

- `llm.verify`
- `llm.request`
- `llm.models`
- `registry.login`
- `webhook.send`
- `k8s.auth`
- `sign.payload`

## Profile Model

A profile binds:

- backend
- secret reference
- optional endpoint metadata
- action policy
- reveal policy

Example shape:

```json
{
  "name": "default-llm",
  "category": "llm_api",
  "backend": "1password-service-account",
  "secret_ref": "op://AIAGENTS/PIPELLM_API_KEY/credential",
  "reveal_policy": "never",
  "usage_policy": "opaque_use",
  "allowed_actions": ["llm.verify", "llm.request", "llm.models"]
}
```

## First Backends

### Priority 1

- `1password-service-account`
- `macos-keychain`

### Future backends

#### Linux

- `libsecret`
- `gnome-keyring`
- `kwallet`
- `pass`
- `gpg`
- `systemd-creds`

#### Cloud / enterprise / server

- `1password-connect`
- `vault`
- `aws-secrets-manager`
- `gcp-secret-manager`
- `azure-key-vault`
- `k8s-secret`
- `sops`
- `encrypted-file`

## CLI Direction

Preferred agent-facing commands:

- `aih llm verify --profile <name>`
- `aih llm request --profile <name> ...`
- `aih registry login --profile <name>`
- `aih webhook send --profile <name> ...`

Restricted/default-denied commands:

- plaintext secret reveal
- broad secret enumeration
- arbitrary backend passthrough

## Configuration Layers

The runtime should support layered configuration:

- global workstation config
- repo-local config
- explicit session config file

This allows project-specific secret and profile contracts without hardcoding them into the toolkit source.

## Contract Philosophy

The toolkit should hardcode:

- profile and alias schemas
- policy semantics
- runtime orchestration
- reusable protocol adapters

The toolkit should not hardcode:

- every SaaS-specific request body
- every private endpoint
- every project-specific key name

Administrators provide those through config so the same runtime can serve
workstation-global, repo-local, and session-local contracts.
