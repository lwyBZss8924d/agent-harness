#!/bin/zsh

set -euo pipefail

repo_root="${0:A:h:h}"
src_dir="${repo_root}/build/bin"
build_manifest="${repo_root}/build/release-manifest.json"
dist_root="${repo_root}/dist"
releases_root="${dist_root}/releases"
current_link="${dist_root}/current"
shim_dir="${HOME}/.local/bin"

if [[ ! -x "${src_dir}/aih" ]]; then
  print -u2 "missing built binary: ${src_dir}/aih"
  print -u2 "run ./scripts/build-release.sh first"
  exit 2
fi
if [[ ! -f "${build_manifest}" ]]; then
  print -u2 "missing build manifest: ${build_manifest}"
  print -u2 "run ./scripts/build-release.sh first"
  exit 2
fi

version="$(python3 - <<'PY' "${build_manifest}"
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as fh:
    print(json.load(fh)["version"])
PY
)"

release_root="${releases_root}/${version}"
bin_dir="${release_root}/bin"

mkdir -p "$bin_dir" "$shim_dir" "$releases_root"

cp "${src_dir}/aih" "${bin_dir}/aih"
cp "${src_dir}/op-sa-broker" "${bin_dir}/op-sa-broker"
cp "${src_dir}/op-sa-broker-client" "${bin_dir}/op-sa-broker-client"
cp "${build_manifest}" "${release_root}/manifest.json"

ln -sfn "${release_root}" "${current_link}"

write_shim() {
  local name="$1"
  local target="$2"
  cat > "${shim_dir}/${name}" <<EOF
#!/bin/zsh
set -euo pipefail
exec "${current_link}/bin/${target}" "\$@"
EOF
  chmod 700 "${shim_dir}/${name}"
}

write_shim "aih-go" "aih"
write_shim "op-sa-broker-go" "op-sa-broker"
write_shim "op-sa-broker-client-go" "op-sa-broker-client"

cat > "${dist_root}/install-manifest.json" <<EOF
{
  "channel": "local-release",
  "current_version": "${version}",
  "current_release_root": "${release_root}",
  "current_link": "${current_link}",
  "shim_dir": "${shim_dir}",
  "binaries": [
    "aih-go",
    "op-sa-broker-go",
    "op-sa-broker-client-go"
  ]
}
EOF

echo "installed release binaries under ${release_root}"
