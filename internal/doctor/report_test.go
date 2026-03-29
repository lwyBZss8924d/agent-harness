package doctor

import (
	"testing"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/facts"
)

func TestSummarizeIssues(t *testing.T) {
	summary := summarizeIssues("needs_attention", []Issue{
		{Severity: "high", Source: "runtime"},
		{Severity: "info", Source: "facts"},
		{Severity: "info", Source: "facts"},
	})
	if summary.Verdict != "blocking" {
		t.Fatalf("Verdict = %q", summary.Verdict)
	}
	if summary.SeverityCounts["high"] != 1 || summary.SeverityCounts["info"] != 2 {
		t.Fatalf("SeverityCounts = %#v", summary.SeverityCounts)
	}
	if summary.SourceCounts["facts"] != 2 {
		t.Fatalf("SourceCounts = %#v", summary.SourceCounts)
	}
}

func TestFilterIssues(t *testing.T) {
	cfg := config.Config{
		Doctor: config.DoctorConfig{
			IssueIgnoreKinds:   []string{"ignore_me"},
			IssueIgnoreSources: []string{"browser"},
		},
	}
	issues := filterIssues(cfg, []Issue{
		{Kind: "ignore_me", Source: "runtime"},
		{Kind: "keep_me", Source: "facts"},
		{Kind: "other", Source: "browser"},
	})
	if len(issues) != 1 {
		t.Fatalf("len(issues) = %d, want 1", len(issues))
	}
	if issues[0].Kind != "keep_me" {
		t.Fatalf("issues[0].Kind = %q", issues[0].Kind)
	}
}

func TestShouldIgnorePathIssue(t *testing.T) {
	cfg := config.Config{
		Doctor: config.DoctorConfig{
			PathIgnoreEntries:  []string{"skipme"},
			PathIgnorePrefixes: []string{"/var/run/com.apple.security.cryptexd/"},
			PathIgnoreContains: []string{"anaconda3"},
		},
	}
	if !shouldIgnorePathIssue(cfg, facts.Issue{Kind: "missing_path_entry", Detail: "PATH contains a missing directory: /foo/skipme"}) {
		t.Fatal("expected entry ignore")
	}
	if !shouldIgnorePathIssue(cfg, facts.Issue{Kind: "missing_path_entry", Detail: "PATH contains a missing directory: /var/run/com.apple.security.cryptexd/x"}) {
		t.Fatal("expected prefix ignore")
	}
	if !shouldIgnorePathIssue(cfg, facts.Issue{Kind: "missing_path_entry", Detail: "PATH contains a missing directory: /Users/x/anaconda3/condabin"}) {
		t.Fatal("expected contains ignore")
	}
}

func TestAggregateIssues(t *testing.T) {
	issues := aggregateIssues([]Issue{
		{Severity: "info", Kind: "missing_path_entry", Source: "facts", Detail: "a"},
		{Severity: "medium", Kind: "missing_path_entry", Source: "facts", Detail: "b"},
		{Severity: "high", Kind: "broker_socket_missing", Source: "runtime", Detail: "c"},
	})
	if len(issues) != 2 {
		t.Fatalf("len(issues) = %d, want 2", len(issues))
	}
	found := false
	for _, issue := range issues {
		if issue.Kind == "missing_path_entries" {
			found = true
			if issue.Severity != "medium" {
				t.Fatalf("aggregated severity = %q", issue.Severity)
			}
		}
	}
	if !found {
		t.Fatal("expected aggregated missing_path_entries issue")
	}
}

func TestSeverityOverride(t *testing.T) {
	cfg := config.Config{
		Doctor: config.DoctorConfig{
			IssueSeverityOverrides: map[string]string{
				"broker_socket_missing": "medium",
			},
		},
	}
	if got := severityFor(cfg, "broker_socket_missing", "high"); got != "medium" {
		t.Fatalf("severityFor = %q", got)
	}
}

func TestRecommendedActions(t *testing.T) {
	actions := recommendedActions([]Issue{
		{Kind: "broker_socket_missing"},
		{Kind: "missing_generated_facts"},
		{Kind: "docker_server_unreachable"},
	})
	if len(actions) < 3 {
		t.Fatalf("expected multiple actions, got %#v", actions)
	}
}

func TestShouldIgnoreIssue(t *testing.T) {
	cfg := config.Config{
		Doctor: config.DoctorConfig{
			IssueIgnoreKinds:   []string{"foo"},
			IssueIgnoreSources: []string{"browser"},
		},
	}
	if !shouldIgnoreIssue(cfg, Issue{Kind: "foo", Source: "facts"}) {
		t.Fatal("expected ignore by kind")
	}
	if !shouldIgnoreIssue(cfg, Issue{Kind: "bar", Source: "browser"}) {
		t.Fatal("expected ignore by source")
	}
	if shouldIgnoreIssue(cfg, Issue{Kind: "bar", Source: "facts"}) {
		t.Fatal("did not expect ignore")
	}
}
