#!/bin/zsh

set -euo pipefail

repo_root="${0:A:h:h}"
system_bin_dir="${AIH_SYSTEM_BIN_DIR:-/usr/local/bin}"
launch_agent_dir="${AIH_LAUNCH_AGENT_DIR:-${HOME}/Library/LaunchAgents}"
broker_label="${AIH_BROKER_LABEL:-com.ai-agent-owner.op-sa-broker}"
config_file="${AIH_CONFIG_FILE:-${HOME}/.config/aih/config.json}"
broker_socket="${AIH_BROKER_SOCKET:-${HOME}/.agents/run/op-sa-broker.sock}"
browser_runtime_dir="${repo_root}/dist/current/runtime/browser"
shell_runtime_dir="${repo_root}/dist/current/runtime/shell"
config_runtime_dir="${HOME}/.config/aih/runtime"

"${repo_root}/scripts/build-release.sh"
"${repo_root}/scripts/install-release.sh"

mkdir -p "${system_bin_dir}" "${launch_agent_dir}" "${HOME}/tmp" "${HOME}/.cache/go-build"

cat > "${system_bin_dir}/aih" <<EOF
#!/bin/zsh
set -euo pipefail
exec "${repo_root}/dist/current/bin/aih" "\$@"
EOF
chmod 755 "${system_bin_dir}/aih"

python3 - <<'PY' "${repo_root}" "${launch_agent_dir}" "${broker_label}" "${config_file}" "${broker_socket}"
from pathlib import Path
import sys

repo_root = Path(sys.argv[1])
launch_agent_dir = Path(sys.argv[2])
broker_label = sys.argv[3]
config_file = sys.argv[4]
broker_socket = sys.argv[5]

template = (repo_root / "packaging" / "launchd" / "com.ai-agent-owner.op-sa-broker.plist.tmpl").read_text()
plist = (
    template
    .replace("{{BROKER_LABEL}}", broker_label)
    .replace("{{BROKER_BINARY}}", str(repo_root / "dist" / "current" / "bin" / "op-sa-broker"))
    .replace("{{WORKING_DIRECTORY}}", str(repo_root))
    .replace("{{STDOUT_PATH}}", str(Path.home() / ".agents" / "logs" / "op-sa-broker.stdout.log"))
    .replace("{{STDERR_PATH}}", str(Path.home() / ".agents" / "logs" / "op-sa-broker.stderr.log"))
)
insert = """
  <key>EnvironmentVariables</key>
  <dict>
    <key>AIH_CONFIG_FILE</key>
    <string>{config}</string>
    <key>AIH_BROKER_SOCKET</key>
    <string>{socket}</string>
    <key>TMPDIR</key>
    <string>{tmpdir}</string>
    <key>GOCACHE</key>
    <string>{gocache}</string>
  </dict>
""".format(
    config=config_file,
    socket=broker_socket,
    tmpdir=str(Path.home() / "tmp"),
    gocache=str(Path.home() / ".cache" / "go-build"),
)
plist = plist.replace('  <key>StandardOutPath</key>\n', insert + '  <key>StandardOutPath</key>\n')
(launch_agent_dir / f"{broker_label}.plist").write_text(plist)
PY

echo "installed workstation package:"
echo "  aih shim: ${system_bin_dir}/aih"
echo "  broker plist: ${launch_agent_dir}/${broker_label}.plist"
echo "  release root: ${repo_root}/dist/current"

if [[ -d "${shell_runtime_dir}" ]]; then
  mkdir -p "${config_runtime_dir}/shell"
  cp -R "${shell_runtime_dir}/." "${config_runtime_dir}/shell/"
  chmod 755 "${config_runtime_dir}/shell/"*.sh 2>/dev/null || true
  echo "  shell runtime: ${config_runtime_dir}/shell"
fi

if [[ -f "${browser_runtime_dir}/package.json" ]]; then
  if command -v npm >/dev/null 2>&1; then
    (
      cd "${browser_runtime_dir}"
      npm install --omit=dev --no-audit --no-fund >/dev/null
    )
    echo "  browser runtime: ${browser_runtime_dir}"
  else
    echo "warning: npm not found; browser runtime dependencies not installed under ${browser_runtime_dir}" >&2
  fi
fi
