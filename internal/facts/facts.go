package facts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/discovery"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/paths"
)

type Snapshot struct {
	GeneratedAtUTC string             `json:"generated_at_utc"`
	Host           HostInfo           `json:"host"`
	Path           PathInfo           `json:"path"`
	Runtime        RuntimeInfo        `json:"runtime"`
	SSH            SSHInfo            `json:"ssh"`
	Toolkit        ToolkitInfo        `json:"toolkit"`
	Discovery      discovery.Snapshot `json:"discovery"`
}

type HostInfo struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Kernel   string `json:"kernel"`
	Arch     string `json:"arch"`
	User     string `json:"user"`
	Home     string `json:"home"`
	Shell    string `json:"shell"`
}

type PathInfo struct {
	Head            string   `json:"head,omitempty"`
	Entries         []string `json:"entries"`
	FreshLoginHead  string   `json:"fresh_login_head,omitempty"`
	FreshLoginPaths []string `json:"fresh_login_entries,omitempty"`
	Anomalies       []Issue  `json:"anomalies,omitempty"`
	FreshAnomalies  []Issue  `json:"fresh_login_anomalies,omitempty"`
}

type Issue struct {
	Severity string `json:"severity"`
	Kind     string `json:"kind"`
	Detail   string `json:"detail"`
}

type RuntimeInfo struct {
	Node    string `json:"node,omitempty"`
	NPM     string `json:"npm,omitempty"`
	Python3 string `json:"python3,omitempty"`
	Go      string `json:"go,omitempty"`
	Git     string `json:"git,omitempty"`
}

type SSHInfo struct {
	SSHClient     string `json:"ssh_client,omitempty"`
	SSHConnection string `json:"ssh_connection,omitempty"`
	SSHTTY        string `json:"ssh_tty,omitempty"`
}

type ToolkitInfo struct {
	FactsJSON      string `json:"facts_json"`
	FactsMarkdown  string `json:"facts_markdown"`
	AgentsDir      string `json:"agents_dir"`
	BrokerSocket   string `json:"broker_socket"`
	CompatBackend  string `json:"compat_backend"`
	SourceRepoRoot string `json:"source_repo_root"`
}

type RefreshResult struct {
	GeneratedAtUTC string `json:"generated_at_utc"`
	FactsJSON      string `json:"facts_json"`
	FactsMarkdown  string `json:"facts_markdown"`
}

func Collect(cfg config.Config) Snapshot {
	pathEntries := splitPath(os.Getenv("PATH"))
	freshLogin := freshLoginPathEntries()

	host := HostInfo{
		Hostname: commandOutput("hostname"),
		OS:       hostOS(),
		Kernel:   commandOutput("uname", "-a"),
		Arch:     runtime.GOARCH,
		User:     fallbackString(commandOutput("whoami"), os.Getenv("USER")),
		Home:     paths.HomeDir(),
		Shell:    os.Getenv("SHELL"),
	}

	if runtime.GOOS == "darwin" {
		if arch := commandOutput("uname", "-m"); arch != "" {
			host.Arch = arch
		}
	}

	return Snapshot{
		GeneratedAtUTC: time.Now().UTC().Format(time.RFC3339),
		Host:           host,
		Path: PathInfo{
			Head:            first(pathEntries),
			Entries:         pathEntries,
			FreshLoginHead:  first(freshLogin),
			FreshLoginPaths: freshLogin,
			Anomalies:       pathAnomalies(pathEntries),
			FreshAnomalies:  pathAnomalies(freshLogin),
		},
		Runtime: RuntimeInfo{
			Node:    commandOutput("node", "--version"),
			NPM:     commandOutput("npm", "--version"),
			Python3: commandOutput("python3", "--version"),
			Go:      commandOutput("go", "version"),
			Git:     commandOutput("git", "--version"),
		},
		SSH: SSHInfo{
			SSHClient:     os.Getenv("SSH_CLIENT"),
			SSHConnection: os.Getenv("SSH_CONNECTION"),
			SSHTTY:        os.Getenv("SSH_TTY"),
		},
		Toolkit: ToolkitInfo{
			FactsJSON:      paths.FactsJSONPath(),
			FactsMarkdown:  paths.FactsMarkdownPath(),
			AgentsDir:      cfg.AgentsDir,
			BrokerSocket:   cfg.Broker.SocketPath,
			CompatBackend:  cfg.CompatBackend,
			SourceRepoRoot: cfg.SourceRepoRoot,
		},
		Discovery: discovery.SnapshotFromConfig(cfg),
	}
}

func Refresh(cfg config.Config) (RefreshResult, error) {
	snapshot := Collect(cfg)
	generatedDir := paths.GeneratedFactsDir()
	if err := os.MkdirAll(generatedDir, 0o755); err != nil {
		return RefreshResult{}, err
	}

	jsonPayload := map[string]any{
		"facts": snapshot,
	}
	encoded, err := json.MarshalIndent(jsonPayload, "", "  ")
	if err != nil {
		return RefreshResult{}, err
	}
	if err := os.WriteFile(paths.FactsJSONPath(), append(encoded, '\n'), 0o644); err != nil {
		return RefreshResult{}, err
	}
	if err := os.WriteFile(paths.FactsMarkdownPath(), []byte(markdownSummary(snapshot)), 0o644); err != nil {
		return RefreshResult{}, err
	}

	return RefreshResult{
		GeneratedAtUTC: snapshot.GeneratedAtUTC,
		FactsJSON:      paths.FactsJSONPath(),
		FactsMarkdown:  paths.FactsMarkdownPath(),
	}, nil
}

func EnvSummary(cfg config.Config) map[string]any {
	snapshot := Collect(cfg)
	return map[string]any{
		"generated_at_utc": snapshot.GeneratedAtUTC,
		"host":             snapshot.Host,
		"runtime":          snapshot.Runtime,
		"path_head":        snapshot.Path.Head,
		"broker_socket":    snapshot.Toolkit.BrokerSocket,
		"compat_backend":   snapshot.Toolkit.CompatBackend,
		"ssh":              snapshot.SSH,
	}
}

func markdownSummary(snapshot Snapshot) string {
	lines := []string{
		fmt.Sprintf("# Host Facts - %s", fallbackString(snapshot.Host.Hostname, "unknown-host")),
		"",
		fmt.Sprintf("Generated at: %s", snapshot.GeneratedAtUTC),
		"",
		"## Summary",
		"",
		fmt.Sprintf("- Arch: %s", fallbackString(snapshot.Host.Arch, "unknown")),
		fmt.Sprintf("- User: %s", fallbackString(snapshot.Host.User, "unknown")),
		fmt.Sprintf("- Home: %s", fallbackString(snapshot.Host.Home, "unknown")),
		fmt.Sprintf("- Shell: %s", fallbackString(snapshot.Host.Shell, "unknown")),
		fmt.Sprintf("- Node: %s", fallbackString(snapshot.Runtime.Node, "unavailable")),
		fmt.Sprintf("- Python: %s", fallbackString(snapshot.Runtime.Python3, "unavailable")),
		fmt.Sprintf("- Go: %s", fallbackString(snapshot.Runtime.Go, "unavailable")),
		fmt.Sprintf("- Git: %s", fallbackString(snapshot.Runtime.Git, "unavailable")),
		fmt.Sprintf("- PATH head: %s", fallbackString(snapshot.Path.Head, "unknown")),
		fmt.Sprintf("- Fresh login PATH head: %s", fallbackString(snapshot.Path.FreshLoginHead, "unknown")),
		fmt.Sprintf("- Broker socket: %s", fallbackString(snapshot.Toolkit.BrokerSocket, "unset")),
		"",
		"## Files",
		"",
		fmt.Sprintf("- Facts JSON: %s", snapshot.Toolkit.FactsJSON),
		fmt.Sprintf("- Facts Markdown: %s", snapshot.Toolkit.FactsMarkdown),
	}
	return strings.Join(lines, "\n") + "\n"
}

func splitPath(raw string) []string {
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, string(os.PathListSeparator))
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func pathAnomalies(entries []string) []Issue {
	issues := make([]Issue, 0)
	seenMissing := map[string]struct{}{}
	for _, entry := range entries {
		switch {
		case strings.Contains(entry, "~"):
			issues = append(issues, Issue{
				Severity: "medium",
				Kind:     "literal_tilde_segment",
				Detail:   fmt.Sprintf("PATH contains a non-expanded segment: %s", entry),
			})
		default:
			if _, err := os.Stat(entry); err != nil {
				if _, exists := seenMissing[entry]; exists {
					continue
				}
				seenMissing[entry] = struct{}{}
				issues = append(issues, Issue{
					Severity: "info",
					Kind:     "missing_path_entry",
					Detail:   fmt.Sprintf("PATH contains a missing directory: %s", entry),
				})
			}
		}
	}
	return issues
}

func freshLoginPathEntries() []string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/zsh"
	}
	output, err := exec.Command(shell, "-lic", "printf '%s' \"$PATH\"").CombinedOutput()
	if err != nil {
		return []string{}
	}
	return splitPath(string(output))
}

func commandOutput(name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	output, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func hostOS() string {
	if runtime.GOOS == "darwin" {
		if value := commandOutput("sw_vers"); value != "" {
			return value
		}
	}
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

func first(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[0]
}

func fallbackString(value string, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func FactsPaths() map[string]string {
	return map[string]string{
		"facts_json": paths.FactsJSONPath(),
		"facts_md":   paths.FactsMarkdownPath(),
	}
}

func RelativeFactsPaths() map[string]string {
	return map[string]string{
		"facts_json": filepath.ToSlash(paths.FactsJSONPath()),
		"facts_md":   filepath.ToSlash(paths.FactsMarkdownPath()),
	}
}
