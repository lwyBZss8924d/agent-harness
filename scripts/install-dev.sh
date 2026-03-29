#!/bin/zsh

set -euo pipefail

repo_root="${0:A:h:h}"
bin_dir="${HOME}/.local/bin"
mkdir -p "$bin_dir"

write_shim() {
  local name="$1"
  local target="$2"
  cat > "${bin_dir}/${name}" <<EOF
#!/bin/zsh
set -euo pipefail
exec /usr/local/bin/env go run ${repo_root}/${target} "\$@"
EOF
  chmod 700 "${bin_dir}/${name}"
}

write_shim "aih-go-dev" "cmd/aih"
write_shim "op-sa-broker-go-dev" "cmd/op-sa-broker"
write_shim "op-sa-broker-client-go-dev" "cmd/op-sa-broker-client"

echo "installed dev shims to ${bin_dir}"
