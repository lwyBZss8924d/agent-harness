#!/bin/zsh

set -euo pipefail

repo_root="${0:A:h:h}"
src_root="${repo_root}/skills/workstation-harness"
dest_root="${1:-${HOME}/.agents/skills}"

mkdir -p "${dest_root}"

for skill_dir in "${src_root}"/*; do
  [[ -d "${skill_dir}" ]] || continue
  skill_name="$(basename "${skill_dir}")"
  dest_dir="${dest_root}/workstation-harness-${skill_name}"
  rm -rf "${dest_dir}"
  mkdir -p "${dest_dir}"
  rsync -a "${skill_dir}/" "${dest_dir}/"
done

echo "installed workstation harness skills to ${dest_root}"
