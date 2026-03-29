# DevOps and Runtime Ops

## Purpose

This domain brings generic Docker, Kubernetes, registry, and infrastructure
status into the AI-first workstation harness without binding the product to one
specific local runtime or one cloud vendor.

## Current Command

```bash
aih devops status --json
```

## Current Surface

The first Go version reports:

- configured DevOps tool registry with availability and version data
- Docker context and daemon status
- Kubernetes client/context/cluster reachability
- registry auth/config state from Docker config
- running Docker containers
- local Docker images
- Kubernetes namespaces
- Kubernetes contexts
- Kubernetes nodes

Current command surface:

```bash
aih devops status --json
aih devops tools --json
aih devops docker --json
aih devops docker containers --json
aih devops docker images --json
aih devops docker volumes --json
aih devops docker compose-projects --json
aih devops docker inspect --target <id> --json
aih devops docker logs --container <id> --tail 50 --json
aih devops kubernetes --json
aih devops kubernetes contexts --json
aih devops kubernetes namespaces --json
aih devops kubernetes nodes --json
aih devops kubernetes deployments --json
aih devops kubernetes events --json
aih devops kubernetes services --json
aih devops kubernetes pods --json
aih devops kubernetes logs --pod <name> --namespace <ns> --tail 50 --json
aih devops registry --json
```

## Design Boundary

This domain should:

- stay configuration-driven
- remain operator-facing and JSON-first
- avoid hardcoding one workstation's historical runtime assumptions as product
  invariants

It should not:

- become a full deployment orchestrator by default
- hardcode one vendor's workflow as the only supported path
- assume one runtime like OrbStack, Docker Desktop, or a specific cloud

## Next Steps

- feed generic DevOps signals back into doctor heuristics
- expand registry/runtime status surfaces
- add richer orchestration and CI/CD helper flows
- keep moving from status-only views toward light operator workflows
