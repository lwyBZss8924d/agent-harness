package doctor

import (
	"os"
	"slices"
	"strings"
	"time"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/authstatus"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/devopsstatus"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/facts"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/runtimeprobe"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/servicestatus"
)

type Issue struct {
	Severity string `json:"severity"`
	Kind     string `json:"kind"`
	Source   string `json:"source"`
	Detail   string `json:"detail"`
}

type Report struct {
	GeneratedAtUTC    string               `json:"generated_at_utc"`
	Status            string               `json:"status"`
	Summary           Summary              `json:"summary"`
	Issues            []Issue              `json:"issues"`
	FactsFile         string               `json:"facts_file"`
	FactsFileAgeHours *float64             `json:"facts_file_age_hours,omitempty"`
	Runtime           runtimeprobe.Probe   `json:"runtime"`
	Services          servicestatus.Status `json:"services"`
	Auth              authstatus.Status    `json:"auth"`
	DevOps            devopsstatus.Status  `json:"devops"`
}

type Summary struct {
	Verdict            string         `json:"verdict"`
	HighestSeverity    string         `json:"highest_severity"`
	SeverityCounts     map[string]int `json:"severity_counts"`
	SourceCounts       map[string]int `json:"source_counts"`
	BlockingSources    []string       `json:"blocking_sources,omitempty"`
	RecommendedActions []string       `json:"recommended_actions,omitempty"`
	IssueCount         int            `json:"issue_count"`
}

func Collect(cfg config.Config) Report {
	snapshot := facts.Collect(cfg)
	runtimeProbe := runtimeprobe.Collect(cfg)
	services := servicestatus.Collect(cfg)
	auth := authstatus.Collect(cfg)
	devops := devopsstatus.Collect(cfg)

	issues := make([]Issue, 0)
	freshIssueDetails := map[string]struct{}{}
	for _, item := range snapshot.Path.FreshAnomalies {
		freshIssueDetails[item.Detail] = struct{}{}
	}

	if cfg.Browser.DoctorChecksEnabled && cfg.Browser.ChromeBinaryPath != "" {
		if _, err := os.Stat(cfg.Browser.ChromeBinaryPath); err == nil {
			if cfg.Browser.PlaywrightHelperPath != "" {
				if _, helperErr := os.Stat(cfg.Browser.PlaywrightHelperPath); helperErr != nil {
					issues = append(issues, Issue{
						Severity: severityFor(cfg, "browser_playwright_helper_missing", "medium"),
						Kind:     "browser_playwright_helper_missing",
						Source:   "browser",
						Detail:   "Browser automation helper is missing",
					})
				}
			}
			if cfg.Browser.PlaywrightCoreDir != "" {
				if _, coreErr := os.Stat(cfg.Browser.PlaywrightCoreDir); coreErr != nil {
					issues = append(issues, Issue{
						Severity: severityFor(cfg, "browser_playwright_core_missing", "medium"),
						Kind:     "browser_playwright_core_missing",
						Source:   "browser",
						Detail:   "Playwright core runtime is missing",
					})
				}
			}
		}
	}
	for _, item := range snapshot.Path.Anomalies {
		if shouldIgnorePathIssue(cfg, item) {
			continue
		}
		severity := severityFor(cfg, item.Kind, item.Severity)
		detail := item.Detail
		if _, ok := freshIssueDetails[item.Detail]; !ok && cfg.Doctor.IgnoreCurrentSessionOnly {
			severity = "info"
			detail = item.Detail + " (current session only; fresh login shell is cleaner)"
		}
		issues = append(issues, Issue{
			Severity: severity,
			Kind:     item.Kind,
			Source:   "facts",
			Detail:   detail,
		})
	}

	factsFile := snapshot.Toolkit.FactsJSON
	var factsAgeHours *float64
	if info, err := os.Stat(factsFile); err == nil {
		ageHours := time.Since(info.ModTime()).Hours()
		factsAgeHours = &ageHours
		if ageHours > 24 {
			issues = append(issues, Issue{
				Severity: severityFor(cfg, "stale_generated_facts", "medium"),
				Kind:     "stale_generated_facts",
				Source:   "facts",
				Detail:   "Generated host facts are older than 24 hours",
			})
		}
	} else {
		issues = append(issues, Issue{
			Severity: severityFor(cfg, "missing_generated_facts", "high"),
			Kind:     "missing_generated_facts",
			Source:   "facts",
			Detail:   "Generated host facts file does not exist yet",
		})
	}

	if secretRuntimeConfigured(cfg) {
		if !runtimeProbe.Broker.SocketExists {
			issues = append(issues, Issue{
				Severity: severityFor(cfg, "broker_socket_missing", "high"),
				Kind:     "broker_socket_missing",
				Source:   "runtime",
				Detail:   "Local broker socket is missing",
			})
		} else if !runtimeProbe.Broker.Reachable {
			issues = append(issues, Issue{
				Severity: severityFor(cfg, "broker_unreachable", "high"),
				Kind:     "broker_unreachable",
				Source:   "runtime",
				Detail:   "Local broker socket exists but is not reachable",
			})
		}
	}

	for _, probe := range services.Probes {
		if probe.Required && !probe.Reachable {
			issues = append(issues, Issue{
				Severity: severityFor(cfg, "required_port_probe_unreachable", "high"),
				Kind:     "required_port_probe_unreachable",
				Source:   "service",
				Detail:   "Required service probe is unreachable: " + probe.Name,
			})
		}
	}

	for _, tool := range auth.Tools {
		if tool.Required && !tool.Available {
			issues = append(issues, Issue{
				Severity: severityFor(cfg, "required_auth_tool_missing", "high"),
				Kind:     "required_auth_tool_missing",
				Source:   "auth",
				Detail:   "Required tool is not available in PATH: " + tool.Name,
			})
		}
	}

	for _, tool := range devops.Tools {
		if tool.Required && !tool.Available {
			issues = append(issues, Issue{
				Severity: severityFor(cfg, "required_devops_tool_missing", "high"),
				Kind:     "required_devops_tool_missing",
				Source:   "devops",
				Detail:   "Required DevOps tool is not available in PATH: " + tool.Name,
			})
		}
	}

	if devops.Docker.Available {
		if devops.Docker.ContextHostInForeignHome {
			issues = append(issues, Issue{
				Severity: severityFor(cfg, "docker_context_points_to_foreign_home", "high"),
				Kind:     "docker_context_points_to_foreign_home",
				Source:   "devops",
				Detail:   "Docker context points at another user's home directory",
			})
		}
		if devops.Docker.ContextHostKind == "unix" && devops.Docker.ContextSocketPath != "" && !devops.Docker.ContextSocketExists {
			issues = append(issues, Issue{
				Severity: severityFor(cfg, "docker_context_socket_missing", "high"),
				Kind:     "docker_context_socket_missing",
				Source:   "devops",
				Detail:   "Docker context expects a missing unix socket",
			})
		}
		if devops.Docker.Context != "" && !devops.Docker.ServerReachable {
			issues = append(issues, Issue{
				Severity: severityFor(cfg, "docker_server_unreachable", "medium"),
				Kind:     "docker_server_unreachable",
				Source:   "devops",
				Detail:   "Docker is installed and has a current context, but the daemon is not reachable",
			})
		}
	}

	if devops.Kubernetes.Available && devops.Kubernetes.CurrentContext != "" && !devops.Kubernetes.ClusterReachable {
		issues = append(issues, Issue{
			Severity: severityFor(cfg, "kubernetes_context_unreachable", "info"),
			Kind:     "kubernetes_context_unreachable",
			Source:   "devops",
			Detail:   "kubectl has a current context but cluster-info is not currently reachable",
		})
	}

	issues = filterIssues(cfg, issues)

	status := "ok"
	for _, issue := range issues {
		switch issue.Severity {
		case "high":
			status = "needs_attention"
			goto sorted
		case "medium":
			if status == "ok" {
				status = "degraded"
			}
		}
	}
sorted:
	sortIssues(issues)
	summary := summarizeIssues(status, issues)

	return Report{
		GeneratedAtUTC:    time.Now().UTC().Format(time.RFC3339),
		Status:            status,
		Summary:           summary,
		Issues:            issues,
		FactsFile:         factsFile,
		FactsFileAgeHours: factsAgeHours,
		Runtime:           runtimeProbe,
		Services:          services,
		Auth:              auth,
		DevOps:            devops,
	}
}

func secretRuntimeConfigured(cfg config.Config) bool {
	if cfg.SecretService.Kind == "" {
		return false
	}
	if cfg.SecretService.VaultName != "" || cfg.SecretService.VaultUUID != "" {
		return true
	}
	if len(cfg.SecretAliases) > 0 {
		return true
	}
	if cfg.LLMProfiles.Configured && len(cfg.LLMProfiles.Profiles) > 0 {
		return true
	}
	return false
}

func summarizeIssues(status string, issues []Issue) Summary {
	severityCounts := map[string]int{}
	sourceCounts := map[string]int{}
	blockingSources := map[string]struct{}{}
	highestSeverity := "none"
	for _, issue := range issues {
		severityCounts[issue.Severity]++
		sourceCounts[issue.Source]++
		switch issue.Severity {
		case "high":
			highestSeverity = "high"
			blockingSources[issue.Source] = struct{}{}
		case "medium":
			if highestSeverity == "none" || highestSeverity == "info" {
				highestSeverity = "medium"
			}
		case "info":
			if highestSeverity == "none" {
				highestSeverity = "info"
			}
		}
	}
	verdict := "healthy"
	switch status {
	case "needs_attention":
		verdict = "blocking"
	case "degraded":
		verdict = "degraded"
	}
	sources := make([]string, 0, len(blockingSources))
	for source := range blockingSources {
		sources = append(sources, source)
	}
	slices.Sort(sources)
	actions := recommendedActions(issues)
	return Summary{
		Verdict:            verdict,
		HighestSeverity:    highestSeverity,
		SeverityCounts:     severityCounts,
		SourceCounts:       sourceCounts,
		BlockingSources:    sources,
		RecommendedActions: actions,
		IssueCount:         len(issues),
	}
}

func recommendedActions(issues []Issue) []string {
	set := map[string]struct{}{}
	for _, issue := range issues {
		switch issue.Kind {
		case "broker_socket_missing", "broker_unreachable":
			set["Repair or start the local broker before relying on unattended secret access."] = struct{}{}
		case "missing_generated_facts", "stale_generated_facts":
			set["Refresh generated host facts with `aih facts refresh --json`."] = struct{}{}
		case "missing_path_entries", "missing_path_entry", "literal_tilde_segment":
			set["Review PATH baseline configuration or add path ignore rules for non-actionable session noise."] = struct{}{}
		case "docker_context_points_to_foreign_home", "docker_context_socket_missing", "docker_server_unreachable":
			set["Inspect Docker context and daemon reachability with `aih devops docker --json`."] = struct{}{}
		case "kubernetes_context_unreachable":
			set["Inspect Kubernetes reachability with `aih devops kubernetes --json` and current cluster context."] = struct{}{}
		case "browser_playwright_helper_missing", "browser_playwright_core_missing":
			set["Repair browser automation dependencies or disable browser doctor checks for this workstation profile."] = struct{}{}
		case "required_auth_tool_missing":
			set["Install or re-expose the missing required auth tool in PATH."] = struct{}{}
		case "required_devops_tool_missing":
			set["Install or re-expose the missing required DevOps tool in PATH."] = struct{}{}
		case "required_port_probe_unreachable":
			set["Inspect the required service probe and verify whether the service should be running on this workstation."] = struct{}{}
		}
	}
	if len(set) == 0 {
		return nil
	}
	result := make([]string, 0, len(set))
	for item := range set {
		result = append(result, item)
	}
	slices.Sort(result)
	return result
}

func shouldIgnorePathIssue(cfg config.Config, issue facts.Issue) bool {
	if issue.Kind != "missing_path_entry" && issue.Kind != "literal_tilde_segment" {
		return false
	}
	for _, entry := range cfg.Doctor.PathIgnoreEntries {
		if entry != "" && strings.HasSuffix(issue.Detail, entry) {
			return true
		}
	}
	for _, prefix := range cfg.Doctor.PathIgnorePrefixes {
		if prefix != "" && strings.Contains(issue.Detail, prefix) {
			return true
		}
	}
	for _, fragment := range cfg.Doctor.PathIgnoreContains {
		if fragment != "" && strings.Contains(issue.Detail, fragment) {
			return true
		}
	}
	return false
}

func filterIssues(cfg config.Config, issues []Issue) []Issue {
	filtered := make([]Issue, 0, len(issues))
	for _, issue := range issues {
		if shouldIgnoreIssue(cfg, issue) {
			continue
		}
		filtered = append(filtered, issue)
	}
	return aggregateIssues(filtered)
}

func shouldIgnoreIssue(cfg config.Config, issue Issue) bool {
	for _, kind := range cfg.Doctor.IssueIgnoreKinds {
		if kind == issue.Kind {
			return true
		}
	}
	for _, source := range cfg.Doctor.IssueIgnoreSources {
		if source == issue.Source {
			return true
		}
	}
	return false
}

func aggregateIssues(issues []Issue) []Issue {
	missingPath := make([]Issue, 0)
	others := make([]Issue, 0, len(issues))
	for _, issue := range issues {
		if issue.Kind == "missing_path_entry" && issue.Source == "facts" {
			missingPath = append(missingPath, issue)
			continue
		}
		others = append(others, issue)
	}
	if len(missingPath) == 0 {
		return others
	}
	entries := make([]string, 0, len(missingPath))
	highest := "info"
	for _, issue := range missingPath {
		entries = append(entries, issue.Detail)
		if issue.Severity == "medium" && highest == "info" {
			highest = "medium"
		}
		if issue.Severity == "high" {
			highest = "high"
		}
	}
	others = append(others, Issue{
		Severity: highest,
		Kind:     "missing_path_entries",
		Source:   "facts",
		Detail:   strings.Join(entries, " | "),
	})
	return others
}

func severityFor(cfg config.Config, kind string, defaultSeverity string) string {
	if cfg.Doctor.IssueSeverityOverrides != nil {
		if value, ok := cfg.Doctor.IssueSeverityOverrides[kind]; ok && value != "" {
			return value
		}
	}
	return defaultSeverity
}

func sortIssues(issues []Issue) {
	severityRank := map[string]int{
		"high":   0,
		"medium": 1,
		"info":   2,
	}
	slices.SortStableFunc(issues, func(a, b Issue) int {
		if severityRank[a.Severity] < severityRank[b.Severity] {
			return -1
		}
		if severityRank[a.Severity] > severityRank[b.Severity] {
			return 1
		}
		if a.Kind < b.Kind {
			return -1
		}
		if a.Kind > b.Kind {
			return 1
		}
		if a.Source < b.Source {
			return -1
		}
		if a.Source > b.Source {
			return 1
		}
		return 0
	})
}
