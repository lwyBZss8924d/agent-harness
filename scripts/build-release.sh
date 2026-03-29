#!/bin/zsh

set -euo pipefail

repo_root="${0:A:h:h}"
build_dir="${repo_root}/build/bin"
build_meta_dir="${repo_root}/build"
mkdir -p "$build_dir"
mkdir -p "$build_meta_dir"

cd "${repo_root}"

version="${AIH_VERSION:-0.0.1-dev.1}"
commit="${AIH_GIT_COMMIT:-unknown}"
built_at="${AIH_BUILD_TIME:-$(date -u '+%Y-%m-%dT%H:%M:%SZ')}"
ldflags="-X github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/version.Version=${version} -X github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/version.GitCommit=${commit} -X github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/version.BuildTime=${built_at}"

go build -ldflags "$ldflags" -o "${build_dir}/aih" ./cmd/aih
go build -ldflags "$ldflags" -o "${build_dir}/op-sa-broker" ./cmd/op-sa-broker
go build -ldflags "$ldflags" -o "${build_dir}/op-sa-broker-client" ./cmd/op-sa-broker-client

cat > "${build_meta_dir}/release-manifest.json" <<EOF
{
  "version": "${version}",
  "git_commit": "${commit}",
  "build_time": "${built_at}",
  "binaries": [
    "aih",
    "op-sa-broker",
    "op-sa-broker-client"
  ]
}
EOF

echo "built release binaries under ${build_dir}"
