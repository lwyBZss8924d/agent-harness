---
name: workstation-harness-devops-runtime-ops
description: Use this skill when an agent needs Docker, Kubernetes, registry, or related workstation runtime evidence from `aih devops`, especially for read-only or low-risk operator workflows.
---

# Workstation Harness DevOps Runtime Ops

Use this skill for Docker/Kubernetes/registry status and light operator workflows.

## Use This Skill When

- you need `aih devops status`
- you need Docker context, images, containers, compose projects, inspect, or logs
- you need Kubernetes contexts, namespaces, nodes, pods, deployments, services, or logs
- you need registry auth/config state

## Default Approach

Start narrow instead of jumping to the aggregate view:

- `aih devops docker --json`
- `aih devops registry --json`
- `aih devops kubernetes --json`

Then move to targeted subcommands:

- `aih devops docker containers --json`
- `aih devops docker images --json`
- `aih devops docker compose-projects --json`
- `aih devops docker logs --container <id> --tail 50 --json`
- `aih devops kubernetes contexts --json`
- `aih devops kubernetes namespaces --json`
- `aih devops kubernetes nodes --json`
- `aih devops kubernetes pods --json`
- `aih devops kubernetes deployments --json`
- `aih devops kubernetes services --json`
- `aih devops kubernetes logs --pod <name> --namespace <ns> --tail 50 --json`

## Rules

- Keep workflows read-only or low-risk unless a higher-trust operator flow is explicitly needed.
- Prefer the narrowest subcommand that answers the question.
- Treat runtime unreachability as a state to report, not a cue to guess.

## References

- `../../../docs/DEVOPS_RUNTIME_OPS.md`
- `../../../docs/DOCTOR_SERVICE_AUTH.md`

