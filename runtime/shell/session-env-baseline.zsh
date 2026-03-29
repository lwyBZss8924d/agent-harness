#!/bin/zsh

# Shared agent-first shell baseline for zsh shells and launchd env sync.
if [[ -n "${__AI_AGENT_ENV_BASELINE_SOURCED:-}" ]]; then
  return 0
fi
typeset -g __AI_AGENT_ENV_BASELINE_SOURCED=1

# Cargo / Rust
[ -f "$HOME/.cargo/env" ] && . "$HOME/.cargo/env"

# Homebrew (Intel: /usr/local)
if [ -x /usr/local/bin/brew ]; then
  eval "$(/usr/local/bin/brew shellenv)"
  export PATH="/usr/local/sbin:$PATH"
fi

# Go
export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"

# uv / pipx user tools
export PATH="$HOME/.local/bin:$PATH"
export UV_PYTHON_PREFERENCE="only-system"

# Volta (Node.js) - must be last to win shim priority
export VOLTA_HOME="$HOME/.volta"
export PATH="$VOLTA_HOME/bin:$PATH"

# Prefer the fixed 1Password SSH agent socket for shells and background jobs.
export AI_AGENT_1PASSWORD_SSH_AUTH_SOCK="$HOME/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock"
if [ -d "${AI_AGENT_1PASSWORD_SSH_AUTH_SOCK:h}" ]; then
  export SSH_AUTH_SOCK="$AI_AGENT_1PASSWORD_SSH_AUTH_SOCK"
fi

function ai_agent_normalize_path() {
  typeset -gaU path
  path=("${(@)path:#\~/.dotnet/tools}")
  path=("${(@)path:#/Applications/VMware Fusion.app/Contents/Public}")
  path=("${(@)path:#/Users/*/.celcodex/tmp/*}")
  path=("${(@)path:#/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/local/bin}")
  path=("${(@)path:#/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/bin}")
  path=("${(@)path:#/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/appleinternal/bin}")
  path=("${(@)path:#/Users/*/.orbstack/bin}")

  if [ -d "$HOME/.dotnet/tools" ]; then
    path+=("$HOME/.dotnet/tools")
  fi
  if [ -d "$HOME/.orbstack/bin" ]; then
    path+=("$HOME/.orbstack/bin")
  fi

  export PATH="${(j/:/)path}"
}

ai_agent_normalize_path
export LANG="en_US.UTF-8"
export LC_CTYPE="en_US.UTF-8"
export EDITOR="nvim"
