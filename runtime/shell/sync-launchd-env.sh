#!/bin/zsh

set -eu

SCRIPT_DIR="${0:A:h}"
BASELINE_SCRIPT="${SCRIPT_DIR}/session-env-baseline.zsh"

if [ ! -r "$BASELINE_SCRIPT" ]; then
  print -u2 "sync-launchd-env: missing baseline script: $BASELINE_SCRIPT"
  exit 1
fi

source "$BASELINE_SCRIPT"

vars=(
  HOME
  USER
  SHELL
  LANG
  LC_CTYPE
  EDITOR
  GOPATH
  VOLTA_HOME
  UV_PYTHON_PREFERENCE
  SSH_AUTH_SOCK
  PATH
)

for name in "${vars[@]}"; do
  value="${(P)name:-}"
  if [ -n "$value" ]; then
    launchctl setenv "$name" "$value"
  fi
done
