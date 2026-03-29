#!/bin/zsh

set -euo pipefail

repo_root="${0:A:h:h}"

cd "${repo_root}"

./scripts/build-release.sh
./scripts/install-release.sh

release_json="$(~/.local/bin/aih-go release --json)"

python3 - <<'PY' "${release_json}"
import json
import sys

payload = json.loads(sys.argv[1])

errors = []

if payload.get("release_target") != "0.0.1":
    errors.append("release_target must be 0.0.1")

install = payload.get("install") or {}
if install.get("mode") != "release-installed":
    errors.append("install.mode must be release-installed")
if not install.get("manifest_readable"):
    errors.append("install manifest must be readable")

gates = {gate["id"]: gate for gate in payload.get("gates", [])}
required_gate = gates.get("release-install-sdlc")
if not required_gate:
    errors.append("release-install-sdlc gate missing")
elif required_gate.get("status") != "complete":
    errors.append("release-install-sdlc gate must be complete")

if errors:
    for item in errors:
        print(item, file=sys.stderr)
    sys.exit(2)

print(json.dumps({
    "ok": True,
    "release_target": payload["release_target"],
    "current_version": payload["current_version"],
    "install_mode": install["mode"],
    "manifest_path": install.get("manifest_path"),
}, indent=2))
PY
