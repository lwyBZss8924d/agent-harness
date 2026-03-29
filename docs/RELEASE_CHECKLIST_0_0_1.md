# Release Checklist 0.0.1

This is the formal local release checklist for `aih` `0.0.1`.

## Scope

This checklist covers the first formal release boundary:

- `aih secret`
- local broker runtime
- versioned local release install
- machine-readable release status

It does not cover the later product domains such as doctor, facts, browser,
DevOps, or runtime discovery.

## Required Checks

1. Build release binaries.

```bash
./scripts/build-release.sh
```

2. Install a versioned local release under `dist/releases/<version>`.

```bash
./scripts/install-release.sh
```

3. Verify the installed release reports itself correctly.

```bash
~/.local/bin/aih-go release --json
```

Expected:

- `release_target = 0.0.1`
- `install.mode = release-installed`
- `install.manifest_readable = true`

4. Run the automated local release check.

```bash
./scripts/check-release.sh
```

5. For workstation release candidates, also verify:

- broker is running
- `aih secret status --json`
- `aih secret audit --json`
- a configured opaque-use profile succeeds

## Exit Rule

`0.0.1` is ready for a tagged release candidate when:

- local release check passes
- the reference workstation cutover remains healthy
- no default plaintext secret reveal regression exists
