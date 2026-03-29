package releasegate

import "github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/releaseinstall"

type Status string

const (
	StatusComplete   Status = "complete"
	StatusInProgress Status = "in_progress"
	StatusPlanned    Status = "planned"
)

type Gate struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Status     Status   `json:"status"`
	Scope      string   `json:"scope"`
	Evidence   []string `json:"evidence,omitempty"`
	Remaining  []string `json:"remaining,omitempty"`
	References []string `json:"references,omitempty"`
}

type Domain struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Status     Status   `json:"status"`
	SourceDocs []string `json:"source_docs,omitempty"`
	Projection string   `json:"projection"`
}

type ReleaseStatus struct {
	ReleaseTarget      string                       `json:"release_target"`
	CurrentVersion     string                       `json:"current_version"`
	ReleaseScope       string                       `json:"release_scope"`
	OverallStatus      Status                       `json:"overall_status"`
	Principles         []string                     `json:"principles"`
	CompletedGateCount int                          `json:"completed_gate_count"`
	TotalGateCount     int                          `json:"total_gate_count"`
	Install            releaseinstall.InstallStatus `json:"install"`
	Gates              []Gate                       `json:"gates"`
	NextTranche        []string                     `json:"next_tranche"`
	FutureDomains      []Domain                     `json:"future_domains"`
}

func Evaluate(currentVersion string) ReleaseStatus {
	gates := []Gate{
		{
			ID:     "scope-boundary",
			Title:  "Formal 0.0.1 Product Boundary",
			Status: StatusComplete,
			Scope:  "`aih secret` is the first formal release boundary and explicitly not a transparent gateway.",
			Evidence: []string{
				"docs/PRODUCT_SCOPE.md",
				"docs/DELIVERY_GATE.md",
				"docs/SECRET_RUNTIME.md",
			},
			References: []string{
				"docs/PRODUCT_SCOPE.md",
				"docs/DELIVERY_GATE.md",
			},
		},
		{
			ID:     "background-runtime",
			Title:  "Unattended Local Secret Runtime",
			Status: StatusComplete,
			Scope:  "A workstation-local Go broker provides same-user secret runtime operations without desktop approval in fresh SSH sessions.",
			Evidence: []string{
				"cmd/op-sa-broker",
				"cmd/op-sa-broker-client",
				"internal/brokerdaemon/server.go",
			},
			References: []string{
				"docs/ARCHITECTURE.md",
				"docs/PARITY.md",
			},
		},
		{
			ID:     "backend-adapters",
			Title:  "Initial Custody Backend Set",
			Status: StatusComplete,
			Scope:  "The runtime supports `1password-service-account` and `macos-keychain` through a shared backend interface.",
			Evidence: []string{
				"internal/secretbackend/onepassword/backend.go",
				"internal/secretbackend/macoskeychain/backend.go",
				"internal/secretbackendregistry/registry.go",
			},
			References: []string{
				"docs/ARCHITECTURE.md",
				"docs/SECRET_RUNTIME.md",
			},
		},
		{
			ID:     "policy-and-opaque-use",
			Title:  "Policy-Driven Opaque Use",
			Status: StatusComplete,
			Scope:  "Profiles, aliases, reveal policy, usage policy, and allowed actions enforce delegated use over plaintext reveal by default.",
			Evidence: []string{
				"internal/secretpolicy/policy.go",
				"internal/profile/profile.go",
				"internal/secrets/material.go",
				"internal/cli/aih/app.go",
			},
			References: []string{
				"docs/SECRET_RUNTIME.md",
				"docs/DELIVERY_GATE.md",
			},
		},
		{
			ID:     "agent-cli-ux",
			Title:  "AI-Friendly Secret CLI Surface",
			Status: StatusComplete,
			Scope:  "The CLI provides stable subcommands, JSON status/audit surfaces, concise failure modes, and explicit unsafe/admin gates.",
			Evidence: []string{
				"internal/cli/aih/app.go",
			},
			References: []string{
				"docs/DELIVERY_GATE.md",
			},
		},
		{
			ID:     "config-driven-contracts",
			Title:  "Global and Repo-Local Contracts",
			Status: StatusComplete,
			Scope:  "Workstation and project operators can define aliases and profiles without hardcoding every provider in toolkit source.",
			Evidence: []string{
				"internal/config/config.go",
				"internal/profilescaffold/scaffold.go",
			},
			References: []string{
				"docs/SECRET_RUNTIME.md",
				"docs/DELIVERY_GATE.md",
			},
		},
		{
			ID:     "release-install-sdlc",
			Title:  "Formal 0.0.1 Release and Install Lifecycle",
			Status: StatusComplete,
			Scope:  "The repo provides versioned local release installs, install metadata, machine-readable release status, and a repeatable local release check for 0.0.1.",
			Evidence: []string{
				"docs/RELEASE.md",
				"docs/RELEASE_CHECKLIST_0_0_1.md",
				"scripts/build-release.sh",
				"scripts/install-dev.sh",
				"scripts/install-release.sh",
				"scripts/check-release.sh",
			},
			References: []string{
				"docs/RELEASE.md",
				"docs/RELEASE_CHECKLIST_0_0_1.md",
				"docs/DEPLOYMENT.md",
			},
		},
	}

	completed := 0
	overall := StatusComplete
	for _, gate := range gates {
		if gate.Status == StatusComplete {
			completed++
			continue
		}
		if gate.Status == StatusInProgress {
			overall = StatusInProgress
		}
	}

	futureDomains := []Domain{
		{
			ID:     "facts-discovery",
			Title:  "Facts and Discovery",
			Status: StatusComplete,
			SourceDocs: []string{
				"/Users/ai-agent-owner/.agents/os-dev-environment/index.llms.txt",
				"/Users/ai-agent-owner/.agents/os-dev-environment/system.llms.txt",
			},
			Projection: "Core `aih facts` and `aih env summary` surfaces now provide host identity, runtime inventory, session-sensitive facts, and generated machine snapshots.",
		},
		{
			ID:     "doctor-diagnostics",
			Title:  "Doctor and Diagnostics",
			Status: StatusInProgress,
			SourceDocs: []string{
				"/Users/ai-agent-owner/.agents/harness-cli.llms.txt",
				"/Users/ai-agent-owner/.agents/os-dev-environment/services-ports.llms.txt",
			},
			Projection: "Initial `doctor`, `service status`, and `auth status` command surfaces now exist with verdict summaries, issue-source attribution, severity overrides, and config-driven noise controls; the next step is to deepen the heuristic quality.",
		},
		{
			ID:     "devops-runtime-ops",
			Title:  "DevOps and Unattended Runtime Ops",
			Status: StatusInProgress,
			SourceDocs: []string{
				"/Users/ai-agent-owner/.agents/os-dev-environment/devops-tools.llms.txt",
				"/Users/ai-agent-owner/.agents/os-dev-environment/package-managers.llms.txt",
			},
			Projection: "Initial `devops` status and read-only/light operator subcommands now report tools, Docker, registry, Kubernetes, containers, images, volumes, namespaces, contexts, nodes, events, and logs; the next step is to deepen this into richer generic workflows.",
		},
		{
			ID:     "browser-validation",
			Title:  "Browser and Validation Harness",
			Status: StatusComplete,
			SourceDocs: []string{
				"/Users/ai-agent-owner/.agents/harness-cli.llms.txt",
				"/Users/ai-agent-owner/.agents/os-dev-environment/ai-agents.llms.txt",
			},
			Projection: "Core browser status, launch-cdp, and Playwright verification commands now exist and have been validated on the reference workstation as a generic local harness domain.",
		},
	}

	return ReleaseStatus{
		ReleaseTarget:  "0.0.1",
		CurrentVersion: currentVersion,
		ReleaseScope:   "`aih secret` with AI-first unattended workstation secret runtime",
		OverallStatus:  overall,
		Principles: []string{
			"ai-first-cli",
			"text-and-json-first",
			"opaque-use-over-plaintext-reveal",
			"config-driven-contracts",
			"not-a-transparent-api-gateway",
		},
		CompletedGateCount: completed,
		TotalGateCount:     len(gates),
		Install:            releaseinstall.Detect(),
		Gates:              gates,
		NextTranche: []string{
			"Deepen doctor/service/auth from initial generic surfaces into richer reusable workstation diagnostics.",
			"Deepen facts/discovery with richer inventory and reusable doctor inputs.",
			"Deepen devops/runtime ops from read-only operator subcommands into richer Docker, registry, and orchestration workflows.",
			"Keep future domains configuration-driven and workstation-agnostic.",
		},
		FutureDomains: futureDomains,
	}
}
