#!/bin/zsh

set -euo pipefail

repo_root="${0:A:h:h}"
target_host="${1:-}"
target_root="${2:-/Users/ai-agent-owner/dev-space/aih-toolkit}"

if [[ -z "$target_host" ]]; then
  print -u2 "usage: scripts/deploy-remote.sh <host> [target-root]"
  exit 2
fi

rsync -az --delete \
  --exclude '.git' \
  --exclude 'build' \
  --exclude 'dist' \
  "${repo_root}/" "${target_host}:${target_root}/"

print "deployed source repo to ${target_host}:${target_root}"
