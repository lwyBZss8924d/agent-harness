package aih

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/authstatus"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/broker"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/browser"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/devopsstatus"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/discovery"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/doctor"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/facts"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/keychain"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/llmapi"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/paths"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/profilereadiness"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/profilescaffold"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/releasegate"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/runtimeauth"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/runtimeprobe"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretbackendregistry"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretpolicy"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secrets"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/servicestatus"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/version"
)

type App struct {
	Stdout io.Writer
	Stderr io.Writer
}

func New() *App {
	return &App{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func (a *App) Run(args []string) int {
	if len(args) == 0 {
		a.printHelp()
		return 0
	}

	switch args[0] {
	case "version":
		return a.runVersion(args[1:])
	case "paths":
		return a.runPaths(args[1:])
	case "config":
		return a.runConfig(args[1:])
	case "discover":
		return a.runDiscover(args[1:])
	case "architecture":
		return a.runArchitecture(args[1:])
	case "release":
		return a.runRelease(args[1:])
	case "facts":
		return a.runFacts(args[1:])
	case "env":
		return a.runEnv(args[1:])
	case "doctor":
		return a.runDoctor(args[1:])
	case "service":
		return a.runService(args[1:])
	case "auth":
		return a.runAuth(args[1:])
	case "browser":
		return a.runBrowser(args[1:])
	case "devops":
		return a.runDevOps(args[1:])
	case "profile":
		return a.runProfile(args[1:])
	case "secret":
		return a.runSecret(args[1:])
	case "llm":
		return a.runLLM(args[1:])
	case "help", "-h", "--help":
		a.printHelp()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "aih-go: unknown command: %s\n", args[0])
		a.printHelp()
		return 2
	}
}

func (a *App) runVersion(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	payload := map[string]string{
		"version":    version.Version,
		"git_commit": version.GitCommit,
		"build_time": version.BuildTime,
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	fmt.Fprintf(a.Stdout, "aih-go %s\n", version.Version)
	return 0
}

func (a *App) runPaths(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	payload := map[string]string{
		"source_repo_root": paths.SourceRepoRoot(),
		"agents_dir":       paths.AgentsDir(),
		"legacy_runtime":   paths.LegacyRuntimeDir(),
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	fmt.Fprintf(a.Stdout, "source_repo_root=%s\n", payload["source_repo_root"])
	fmt.Fprintf(a.Stdout, "agents_dir=%s\n", payload["agents_dir"])
	fmt.Fprintf(a.Stdout, "legacy_runtime=%s\n", payload["legacy_runtime"])
	return 0
}

func (a *App) runArchitecture(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	payload := map[string]any{
		"mode":                "product",
		"source_repo_root":    cfg.SourceRepoRoot,
		"legacy_runtime_root": cfg.LegacyRuntime,
		"status":              "go_secret_runtime_cutover_complete",
		"secret_service":      cfg.SecretService.Kind,
		"llm_primary_profile": cfg.LLMProfiles.PrimaryProfile,
		"planned_domains": []string{
			"secret",
			"doctor",
			"env",
			"service",
			"auth",
			"browser",
			"facts",
			"system-runtime-network-discovery",
			"dependency-update-checks",
			"docker",
			"k8s",
			"cicd-devops",
		},
		"principles": []string{
			"ai-first-cli",
			"text-and-json-first",
			"opaque-use-over-plaintext-reveal",
			"config-driven-contracts",
			"not-a-transparent-api-gateway",
		},
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	fmt.Fprintln(a.Stdout, "aih-go product architecture")
	fmt.Fprintln(a.Stdout, "source repo:", cfg.SourceRepoRoot)
	fmt.Fprintln(a.Stdout, "legacy runtime:", cfg.LegacyRuntime)
	fmt.Fprintln(a.Stdout, "status: Go `aih secret` runtime is the active reference implementation")
	return 0
}

func (a *App) runConfig(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	if asJSON {
		return writeJSON(a.Stdout, cfg)
	}
	fmt.Fprintf(a.Stdout, "source_repo_root=%s\n", cfg.SourceRepoRoot)
	fmt.Fprintf(a.Stdout, "agents_dir=%s\n", cfg.AgentsDir)
	fmt.Fprintf(a.Stdout, "legacy_runtime=%s\n", cfg.LegacyRuntime)
	fmt.Fprintf(a.Stdout, "compat_backend=%s\n", cfg.CompatBackend)
	fmt.Fprintf(a.Stdout, "broker_socket=%s\n", cfg.Broker.SocketPath)
	fmt.Fprintf(a.Stdout, "secret_service=%s\n", cfg.SecretService.Kind)
	fmt.Fprintf(a.Stdout, "llm_primary_profile=%s\n", cfg.LLMProfiles.PrimaryProfile)
	return 0
}

func (a *App) runDiscover(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	snapshot := discovery.SnapshotFromConfig(cfg)
	if asJSON {
		return writeJSON(a.Stdout, snapshot)
	}
	fmt.Fprintf(a.Stdout, "current_user=%s\n", snapshot.CurrentUser)
	fmt.Fprintf(a.Stdout, "home_dir=%s\n", snapshot.HomeDir)
	fmt.Fprintf(a.Stdout, "os=%s\n", snapshot.OS)
	fmt.Fprintf(a.Stdout, "arch=%s\n", snapshot.Arch)
	fmt.Fprintf(a.Stdout, "compat_backend=%s\n", snapshot.CompatBackend)
	fmt.Fprintf(a.Stdout, "legacy_runtime=%s\n", snapshot.LegacyRuntime)
	fmt.Fprintf(a.Stdout, "broker_socket=%s\n", snapshot.BrokerSocket)
	return 0
}

func (a *App) runSecret(args []string) int {
	if len(args) == 0 {
		a.printSecretHelp()
		return 0
	}

	switch args[0] {
	case "status":
		return a.runSecretStatus(args[1:])
	case "audit":
		return a.runSecretAudit(args[1:])
	case "list":
		return a.runSecretList(args[1:])
	case "read":
		return a.runSecretRead(args[1:])
	case "get":
		return a.runSecretGet(args[1:])
	case "env":
		return a.runSecretEnv(args[1:])
	case "exec":
		return a.runSecretExec(args[1:])
	case "sudo":
		return a.runSecretSudo(args[1:])
	case "bootstrap-token":
		return a.runSecretBootstrapToken(args[1:])
	case "cache":
		return a.runSecretCache(args[1:])
	case "help", "-h", "--help":
		a.printSecretHelp()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "aih-go secret: unknown subcommand: %s\n", args[0])
		a.printSecretHelp()
		return 2
	}
}

func (a *App) runProfile(args []string) int {
	if len(args) == 0 {
		a.printProfileHelp()
		return 0
	}

	switch args[0] {
	case "kinds":
		return writeJSON(a.Stdout, map[string]any{"kinds": profilescaffold.SupportedKinds()})
	case "scaffold":
		return a.runProfileScaffold(args[1:])
	case "help", "-h", "--help":
		a.printProfileHelp()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "aih-go profile: unknown subcommand: %s\n", args[0])
		a.printProfileHelp()
		return 2
	}
}

func (a *App) runProfileScaffold(args []string) int {
	asJSON := false
	req := profilescaffold.Request{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			asJSON = true
		case "--kind":
			i++
			if i >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go profile scaffold: missing value after --kind")
				return 2
			}
			req.Kind = args[i]
		case "--name":
			i++
			if i >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go profile scaffold: missing value after --name")
				return 2
			}
			req.Name = args[i]
		case "--backend":
			i++
			if i >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go profile scaffold: missing value after --backend")
				return 2
			}
			req.Backend = args[i]
		case "--protocol":
			i++
			if i >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go profile scaffold: missing value after --protocol")
				return 2
			}
			req.Protocol = args[i]
		case "--secret-alias":
			i++
			if i >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go profile scaffold: missing value after --secret-alias")
				return 2
			}
			req.SecretAlias = args[i]
		case "--base-url-alias":
			i++
			if i >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go profile scaffold: missing value after --base-url-alias")
				return 2
			}
			req.BaseURLAlias = args[i]
		case "--default-model":
			i++
			if i >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go profile scaffold: missing value after --default-model")
				return 2
			}
			req.DefaultModel = args[i]
		case "--primary":
			req.PrimaryProfile = true
		default:
			fmt.Fprintf(a.Stderr, "aih-go profile scaffold: unknown flag %s\n", args[i])
			return 2
		}
	}

	result, err := profilescaffold.Build(req)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go profile scaffold: %v\n", err)
		return 2
	}

	payload := map[string]any{
		"secret_aliases": result.SecretAliases,
	}
	if len(result.LLMProfiles) > 0 {
		payload["llm_profiles"] = map[string]any{
			"primary_profile": result.PrimaryProfile,
			"profiles":        result.LLMProfiles,
		}
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) printHelp() {
	fmt.Fprintln(a.Stdout, "aih-go")
	fmt.Fprintln(a.Stdout, "")
	fmt.Fprintln(a.Stdout, "Commands:")
	fmt.Fprintln(a.Stdout, "  version [--json]       Print version metadata")
	fmt.Fprintln(a.Stdout, "  paths [--json]         Print source/runtime path layout")
	fmt.Fprintln(a.Stdout, "  config [--json]        Print effective toolkit config")
	fmt.Fprintln(a.Stdout, "  discover [--json]      Print workstation discovery snapshot")
	fmt.Fprintln(a.Stdout, "  architecture [--json]  Print current product architecture status")
	fmt.Fprintln(a.Stdout, "  release [--json]       Print 0.0.1 delivery gate and release status")
	fmt.Fprintln(a.Stdout, "  facts ...              Generate and inspect host facts snapshots")
	fmt.Fprintln(a.Stdout, "  env ...                Print runtime and workstation summaries")
	fmt.Fprintln(a.Stdout, "  doctor [--json]        Run generic workstation doctor checks")
	fmt.Fprintln(a.Stdout, "  service status [--json] Print listening services and common port probes")
	fmt.Fprintln(a.Stdout, "  auth status [--json]   Print auth/config markers for common agent CLIs")
	fmt.Fprintln(a.Stdout, "  browser ...            Inspect and validate local browser automation surfaces")
	fmt.Fprintln(a.Stdout, "  devops status [--json] Print generic Docker/Kubernetes/DevOps runtime status")
	fmt.Fprintln(a.Stdout, "  profile ...            Scaffold and inspect protocol-based profile contracts")
	fmt.Fprintln(a.Stdout, "  secret ...             Native generic secret service and runtime controls")
	fmt.Fprintln(a.Stdout, "  llm ...                Opaque-use LLM profile actions")
}

func (a *App) runRelease(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	status := releasegate.Evaluate(version.Version)
	if asJSON {
		return writeJSON(a.Stdout, status)
	}
	fmt.Fprintf(a.Stdout, "release_target=%s\n", status.ReleaseTarget)
	fmt.Fprintf(a.Stdout, "current_version=%s\n", status.CurrentVersion)
	fmt.Fprintf(a.Stdout, "release_scope=%s\n", status.ReleaseScope)
	fmt.Fprintf(a.Stdout, "overall_status=%s\n", status.OverallStatus)
	fmt.Fprintf(a.Stdout, "completed_gates=%d/%d\n", status.CompletedGateCount, status.TotalGateCount)
	fmt.Fprintf(a.Stdout, "install_mode=%s\n", status.Install.Mode)
	fmt.Fprintf(a.Stdout, "executable_path=%s\n", status.Install.ExecutablePath)
	fmt.Fprintln(a.Stdout, "next_tranche:")
	for _, item := range status.NextTranche {
		fmt.Fprintf(a.Stdout, "  - %s\n", item)
	}
	return 0
}

func (a *App) printProfileHelp() {
	fmt.Fprintln(a.Stdout, "aih-go profile")
	fmt.Fprintln(a.Stdout, "")
	fmt.Fprintln(a.Stdout, "Subcommands:")
	fmt.Fprintln(a.Stdout, "  kinds                     List supported scaffold kinds")
	fmt.Fprintln(a.Stdout, "  scaffold [options]        Emit a minimal protocol-based config snippet")
}

func (a *App) printFactsHelp() {
	fmt.Fprintln(a.Stdout, "aih-go facts")
	fmt.Fprintln(a.Stdout, "")
	fmt.Fprintln(a.Stdout, "Subcommands:")
	fmt.Fprintln(a.Stdout, "  refresh [--json] [--quiet]   Refresh generated host facts JSON and markdown")
	fmt.Fprintln(a.Stdout, "  path [--json]                Print generated facts file paths")
}

func (a *App) printEnvHelp() {
	fmt.Fprintln(a.Stdout, "aih-go env")
	fmt.Fprintln(a.Stdout, "")
	fmt.Fprintln(a.Stdout, "Subcommands:")
	fmt.Fprintln(a.Stdout, "  summary [--json]             Print a compact workstation/runtime/session summary")
}

func (a *App) printServiceHelp() {
	fmt.Fprintln(a.Stdout, "aih-go service")
	fmt.Fprintln(a.Stdout, "")
	fmt.Fprintln(a.Stdout, "Subcommands:")
	fmt.Fprintln(a.Stdout, "  status [--json]              Print listening services and common port probes")
}

func (a *App) printAuthHelp() {
	fmt.Fprintln(a.Stdout, "aih-go auth")
	fmt.Fprintln(a.Stdout, "")
	fmt.Fprintln(a.Stdout, "Subcommands:")
	fmt.Fprintln(a.Stdout, "  status [--json]              Print auth/config markers for common agent CLIs")
}

func (a *App) printBrowserHelp() {
	fmt.Fprintln(a.Stdout, "aih-go browser")
	fmt.Fprintln(a.Stdout, "")
	fmt.Fprintln(a.Stdout, "Subcommands:")
	fmt.Fprintln(a.Stdout, "  status [--json] [--port N]              Print browser/CDP/Playwright status")
	fmt.Fprintln(a.Stdout, "  launch-cdp [--json] [--port N] [--profile-dir PATH]  Launch Chrome with CDP enabled")
	fmt.Fprintln(a.Stdout, "  verify-playwright [--json] [--port N] [--timeout SEC]  Verify Playwright can connect over CDP")
}

func (a *App) printDevOpsHelp() {
	fmt.Fprintln(a.Stdout, "aih-go devops")
	fmt.Fprintln(a.Stdout, "")
	fmt.Fprintln(a.Stdout, "Subcommands:")
	fmt.Fprintln(a.Stdout, "  status [--json]              Print generic Docker/Kubernetes/DevOps runtime status")
	fmt.Fprintln(a.Stdout, "  tools [--json]               Print configured DevOps tool registry status")
	fmt.Fprintln(a.Stdout, "  docker status [--json]       Print Docker runtime/context status")
	fmt.Fprintln(a.Stdout, "  docker containers [--json]   Print running Docker containers")
	fmt.Fprintln(a.Stdout, "  docker images [--json]       Print local Docker images")
	fmt.Fprintln(a.Stdout, "  docker volumes [--json]      Print Docker volumes")
	fmt.Fprintln(a.Stdout, "  docker compose-projects [--json] Print Docker Compose projects")
	fmt.Fprintln(a.Stdout, "  docker inspect --target ID [--json] Print Docker inspect payload for a target")
	fmt.Fprintln(a.Stdout, "  docker logs --container ID [--tail N] [--json] Print Docker container logs")
	fmt.Fprintln(a.Stdout, "  kubernetes status [--json]   Print Kubernetes client/context status")
	fmt.Fprintln(a.Stdout, "  kubernetes namespaces [--json] Print Kubernetes namespaces")
	fmt.Fprintln(a.Stdout, "  kubernetes contexts [--json] Print Kubernetes contexts")
	fmt.Fprintln(a.Stdout, "  kubernetes nodes [--json]    Print Kubernetes nodes")
	fmt.Fprintln(a.Stdout, "  kubernetes deployments [--namespace NS] [--json] Print Kubernetes deployments")
	fmt.Fprintln(a.Stdout, "  kubernetes events [--namespace NS] [--json] Print Kubernetes events")
	fmt.Fprintln(a.Stdout, "  kubernetes services [--namespace NS] [--json] Print Kubernetes services")
	fmt.Fprintln(a.Stdout, "  kubernetes pods [--namespace NS] [--json] Print Kubernetes pods")
	fmt.Fprintln(a.Stdout, "  kubernetes logs --pod NAME [--namespace NS] [--container NAME] [--tail N] [--json] Print Pod logs")
	fmt.Fprintln(a.Stdout, "  registry status [--json]     Print generic registry auth/config status")
}

func (a *App) printSecretHelp() {
	fmt.Fprintln(a.Stdout, "aih-go secret")
	fmt.Fprintln(a.Stdout, "")
	fmt.Fprintln(a.Stdout, "Subcommands:")
	fmt.Fprintln(a.Stdout, "  status [--json]        Print native runtime status for the generic secret service")
	fmt.Fprintln(a.Stdout, "  audit [--json]         Print native audit findings for the generic secret service")
	fmt.Fprintln(a.Stdout, "  list [--json]          List configured secret aliases and their policies")
	fmt.Fprintln(a.Stdout, "  get <alias>            Resolve a configured alias only when reveal is allowed or explicitly unsafe")
	fmt.Fprintln(a.Stdout, "  read <op://...>        Read a secret reference only in explicit unsafe/admin mode")
	fmt.Fprintln(a.Stdout, "  env                    Print transient service-account env export only in explicit unsafe mode")
	fmt.Fprintln(a.Stdout, "  exec -- <cmd...>       Run a child process with transient service-account env injection")
	fmt.Fprintln(a.Stdout, "  sudo -- <cmd...>       Run a child process through sudo using the configured account-password alias")
	fmt.Fprintln(a.Stdout, "  bootstrap-token        Import a service-account token into the local keychain and restart the broker")
	fmt.Fprintln(a.Stdout, "  cache <alias>          Cache a resolved alias into the macOS keychain for backend handoff")
}

func (a *App) runFacts(args []string) int {
	if len(args) == 0 {
		a.printFactsHelp()
		return 0
	}
	switch args[0] {
	case "refresh":
		return a.runFactsRefresh(args[1:])
	case "path":
		return a.runFactsPath(args[1:])
	case "help", "-h", "--help":
		a.printFactsHelp()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "aih-go facts: unknown subcommand: %s\n", args[0])
		a.printFactsHelp()
		return 2
	}
}

func (a *App) runEnv(args []string) int {
	if len(args) == 0 {
		a.printEnvHelp()
		return 0
	}
	switch args[0] {
	case "summary":
		return a.runEnvSummary(args[1:])
	case "help", "-h", "--help":
		a.printEnvHelp()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "aih-go env: unknown subcommand: %s\n", args[0])
		a.printEnvHelp()
		return 2
	}
}

func (a *App) runDoctor(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	report := doctor.Collect(cfg)
	if asJSON {
		code := writeJSON(a.Stdout, map[string]any{"doctor": report})
		if code != 0 {
			return code
		}
		if report.Status != "ok" {
			return 2
		}
		return 0
	}
	code := writeJSON(a.Stdout, map[string]any{"doctor": report})
	if code != 0 {
		return code
	}
	if report.Status != "ok" {
		return 2
	}
	return 0
}

func (a *App) runService(args []string) int {
	if len(args) == 0 {
		a.printServiceHelp()
		return 0
	}
	switch args[0] {
	case "status":
		return a.runServiceStatus(args[1:])
	case "help", "-h", "--help":
		a.printServiceHelp()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "aih-go service: unknown subcommand: %s\n", args[0])
		a.printServiceHelp()
		return 2
	}
}

func (a *App) runAuth(args []string) int {
	if len(args) == 0 {
		a.printAuthHelp()
		return 0
	}
	switch args[0] {
	case "status":
		return a.runAuthStatus(args[1:])
	case "help", "-h", "--help":
		a.printAuthHelp()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "aih-go auth: unknown subcommand: %s\n", args[0])
		a.printAuthHelp()
		return 2
	}
}

func (a *App) runBrowser(args []string) int {
	if len(args) == 0 {
		a.printBrowserHelp()
		return 0
	}
	switch args[0] {
	case "status":
		return a.runBrowserStatus(args[1:])
	case "launch-cdp":
		return a.runBrowserLaunchCDP(args[1:])
	case "verify-playwright":
		return a.runBrowserVerifyPlaywright(args[1:])
	case "help", "-h", "--help":
		a.printBrowserHelp()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "aih-go browser: unknown subcommand: %s\n", args[0])
		a.printBrowserHelp()
		return 2
	}
}

func (a *App) runDevOps(args []string) int {
	if len(args) == 0 {
		a.printDevOpsHelp()
		return 0
	}
	switch args[0] {
	case "status":
		return a.runDevOpsStatus(args[1:])
	case "tools":
		return a.runDevOpsTools(args[1:])
	case "docker":
		return a.runDevOpsDocker(args[1:])
	case "kubernetes":
		return a.runDevOpsKubernetes(args[1:])
	case "registry":
		return a.runDevOpsRegistry(args[1:])
	case "help", "-h", "--help":
		a.printDevOpsHelp()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "aih-go devops: unknown subcommand: %s\n", args[0])
		a.printDevOpsHelp()
		return 2
	}
}

func (a *App) runDevOpsTools(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	payload := map[string]any{
		"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
		"tools":            devopsstatus.CollectTools(cfg),
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runDevOpsDocker(args []string) int {
	asJSON := false
	filtered := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--json" {
			asJSON = true
			continue
		}
		filtered = append(filtered, arg)
	}
	if len(filtered) == 0 || filtered[0] == "status" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"docker":           devopsstatus.CollectDocker(),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "containers" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"containers":       devopsstatus.ListDockerContainers(),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "images" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"images":           devopsstatus.ListDockerImages(),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "volumes" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"volumes":          devopsstatus.ListDockerVolumes(),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "compose-projects" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"projects":         devopsstatus.ListDockerComposeProjects(),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "inspect" {
		target := ""
		for i := 1; i < len(filtered); i++ {
			switch filtered[i] {
			case "--target":
				if i+1 >= len(filtered) {
					fmt.Fprintln(a.Stderr, "aih-go devops docker inspect: missing value after --target")
					return 2
				}
				target = filtered[i+1]
				i++
			default:
				fmt.Fprintf(a.Stderr, "aih-go devops docker inspect: unknown flag %s\n", filtered[i])
				return 2
			}
		}
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"inspect":          devopsstatus.DockerInspect(target),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "logs" {
		container := ""
		tail := 200
		for i := 1; i < len(filtered); i++ {
			switch filtered[i] {
			case "--container":
				if i+1 >= len(filtered) {
					fmt.Fprintln(a.Stderr, "aih-go devops docker logs: missing value after --container")
					return 2
				}
				container = filtered[i+1]
				i++
			case "--tail":
				if i+1 >= len(filtered) {
					fmt.Fprintln(a.Stderr, "aih-go devops docker logs: missing value after --tail")
					return 2
				}
				value, err := strconv.Atoi(filtered[i+1])
				if err != nil {
					fmt.Fprintf(a.Stderr, "aih-go devops docker logs: invalid tail %q\n", filtered[i+1])
					return 2
				}
				tail = value
				i++
			default:
				fmt.Fprintf(a.Stderr, "aih-go devops docker logs: unknown flag %s\n", filtered[i])
				return 2
			}
		}
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"logs":             devopsstatus.DockerLogs(container, tail),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	fmt.Fprintf(a.Stderr, "aih-go devops docker: unknown subcommand: %s\n", filtered[0])
	return 2
}

func (a *App) runDevOpsKubernetes(args []string) int {
	asJSON := false
	filtered := make([]string, 0, len(args))
	namespace := ""
	for _, arg := range args {
		if arg == "--json" {
			asJSON = true
			continue
		}
		filtered = append(filtered, arg)
	}
	processed := make([]string, 0, len(filtered))
	for i := 0; i < len(filtered); i++ {
		if filtered[i] == "--namespace" {
			if i+1 >= len(filtered) {
				fmt.Fprintln(a.Stderr, "aih-go devops kubernetes: missing value after --namespace")
				return 2
			}
			namespace = filtered[i+1]
			i++
			continue
		}
		processed = append(processed, filtered[i])
	}
	filtered = processed
	if len(filtered) == 0 || filtered[0] == "status" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"kubernetes":       devopsstatus.CollectKubernetes(),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "namespaces" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"namespaces":       devopsstatus.ListKubernetesNamespaces(),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "contexts" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"contexts":         devopsstatus.ListKubernetesContexts(),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "nodes" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"nodes":            devopsstatus.ListKubernetesNodes(),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "deployments" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"deployments":      devopsstatus.ListKubernetesDeployments(namespace),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "events" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"events":           devopsstatus.ListKubernetesEvents(namespace),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "services" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"services":         devopsstatus.ListKubernetesServices(namespace),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "pods" {
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"pods":             devopsstatus.ListKubernetesPods(namespace),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	if filtered[0] == "logs" {
		pod := ""
		container := ""
		tail := 200
		for i := 1; i < len(filtered); i++ {
			switch filtered[i] {
			case "--pod":
				if i+1 >= len(filtered) {
					fmt.Fprintln(a.Stderr, "aih-go devops kubernetes logs: missing value after --pod")
					return 2
				}
				pod = filtered[i+1]
				i++
			case "--container":
				if i+1 >= len(filtered) {
					fmt.Fprintln(a.Stderr, "aih-go devops kubernetes logs: missing value after --container")
					return 2
				}
				container = filtered[i+1]
				i++
			case "--tail":
				if i+1 >= len(filtered) {
					fmt.Fprintln(a.Stderr, "aih-go devops kubernetes logs: missing value after --tail")
					return 2
				}
				value, err := strconv.Atoi(filtered[i+1])
				if err != nil {
					fmt.Fprintf(a.Stderr, "aih-go devops kubernetes logs: invalid tail %q\n", filtered[i+1])
					return 2
				}
				tail = value
				i++
			default:
				fmt.Fprintf(a.Stderr, "aih-go devops kubernetes logs: unknown flag %s\n", filtered[i])
				return 2
			}
		}
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"logs":             devopsstatus.KubernetesLogs(pod, namespace, container, tail),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	fmt.Fprintf(a.Stderr, "aih-go devops kubernetes: unknown subcommand: %s\n", filtered[0])
	return 2
}

func (a *App) runDevOpsRegistry(args []string) int {
	asJSON := false
	filtered := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--json" {
			asJSON = true
			continue
		}
		filtered = append(filtered, arg)
	}
	if len(filtered) == 0 || filtered[0] == "status" {
		cfg := config.Load()
		payload := map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"registry":         devopsstatus.CollectRegistry(cfg),
		}
		if asJSON {
			return writeJSON(a.Stdout, payload)
		}
		return writeJSON(a.Stdout, payload)
	}
	fmt.Fprintf(a.Stderr, "aih-go devops registry: unknown subcommand: %s\n", filtered[0])
	return 2
}

func (a *App) runFactsRefresh(args []string) int {
	asJSON := false
	quiet := false
	for _, arg := range args {
		switch arg {
		case "--json":
			asJSON = true
		case "--quiet":
			quiet = true
		default:
			fmt.Fprintf(a.Stderr, "aih-go facts refresh: unknown flag %s\n", arg)
			return 2
		}
	}

	cfg := config.Load()
	result, err := facts.Refresh(cfg)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go facts refresh: %v\n", err)
		return 2
	}
	if asJSON {
		return writeJSON(a.Stdout, map[string]any{"facts_refresh": result})
	}
	if !quiet {
		fmt.Fprintf(a.Stdout, "facts_json=%s\nfacts_md=%s\n", result.FactsJSON, result.FactsMarkdown)
	}
	return 0
}

func (a *App) runFactsPath(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	payload := facts.FactsPaths()
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	fmt.Fprintf(a.Stdout, "facts_json=%s\n", payload["facts_json"])
	fmt.Fprintf(a.Stdout, "facts_md=%s\n", payload["facts_md"])
	return 0
}

func (a *App) runEnvSummary(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	payload := facts.EnvSummary(cfg)
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runServiceStatus(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	payload := map[string]any{
		"service": servicestatus.Collect(cfg),
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runAuthStatus(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	payload := map[string]any{
		"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
		"auth":             authstatus.Collect(cfg),
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runBrowserStatus(args []string) int {
	asJSON := false
	port := 0
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			asJSON = true
		case "--port":
			if i+1 >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go browser status: missing value after --port")
				return 2
			}
			value, err := strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Fprintf(a.Stderr, "aih-go browser status: invalid port %q\n", args[i+1])
				return 2
			}
			port = value
			i++
		default:
			fmt.Fprintf(a.Stderr, "aih-go browser status: unknown flag %s\n", args[i])
			return 2
		}
	}
	cfg := config.Load()
	payload := map[string]any{
		"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
		"browser":          browser.Collect(cfg, port),
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runBrowserLaunchCDP(args []string) int {
	asJSON := false
	port := 0
	profileDir := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			asJSON = true
		case "--port":
			if i+1 >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go browser launch-cdp: missing value after --port")
				return 2
			}
			value, err := strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Fprintf(a.Stderr, "aih-go browser launch-cdp: invalid port %q\n", args[i+1])
				return 2
			}
			port = value
			i++
		case "--profile-dir":
			if i+1 >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go browser launch-cdp: missing value after --profile-dir")
				return 2
			}
			profileDir = args[i+1]
			i++
		default:
			fmt.Fprintf(a.Stderr, "aih-go browser launch-cdp: unknown flag %s\n", args[i])
			return 2
		}
	}
	cfg := config.Load()
	status, err := browser.LaunchCDP(cfg, port, profileDir)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go browser launch-cdp: %v\n", err)
		return 2
	}
	payload := map[string]any{
		"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
		"browser":          status,
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runBrowserVerifyPlaywright(args []string) int {
	asJSON := false
	port := 0
	timeoutSeconds := 20
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			asJSON = true
		case "--port":
			if i+1 >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go browser verify-playwright: missing value after --port")
				return 2
			}
			value, err := strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Fprintf(a.Stderr, "aih-go browser verify-playwright: invalid port %q\n", args[i+1])
				return 2
			}
			port = value
			i++
		case "--timeout":
			if i+1 >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go browser verify-playwright: missing value after --timeout")
				return 2
			}
			value, err := strconv.Atoi(args[i+1])
			if err != nil {
				fmt.Fprintf(a.Stderr, "aih-go browser verify-playwright: invalid timeout %q\n", args[i+1])
				return 2
			}
			timeoutSeconds = value
			i++
		default:
			fmt.Fprintf(a.Stderr, "aih-go browser verify-playwright: unknown flag %s\n", args[i])
			return 2
		}
	}
	cfg := config.Load()
	payload, exitCode, err := browser.VerifyPlaywright(cfg, port, timeoutSeconds)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go browser verify-playwright: %v\n", err)
		return 2
	}
	if asJSON {
		code := writeJSON(a.Stdout, map[string]any{
			"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
			"browser_verify":   payload,
		})
		if code != 0 {
			return code
		}
		if exitCode != 0 {
			return 2
		}
		return 0
	}
	code := writeJSON(a.Stdout, map[string]any{
		"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
		"browser_verify":   payload,
	})
	if code != 0 {
		return code
	}
	if exitCode != 0 {
		return 2
	}
	return 0
}

func (a *App) runDevOpsStatus(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	payload := map[string]any{
		"devops": devopsstatus.Collect(cfg),
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) printLLMHelp() {
	fmt.Fprintln(a.Stdout, "aih-go llm")
	fmt.Fprintln(a.Stdout, "")
	fmt.Fprintln(a.Stdout, "Subcommands:")
	fmt.Fprintln(a.Stdout, "  profiles [--json]                  List configured LLM profiles without revealing secrets")
	fmt.Fprintln(a.Stdout, "  verify --profile <name> [--json]   Verify an LLM profile through opaque delegated use")
	fmt.Fprintln(a.Stdout, "  request --profile <name> [--json]  Send a JSON request through an opaque delegated LLM profile")
}

func (a *App) runSecretStatus(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	probe := runtimeprobe.Collect(cfg)
	backends := secretbackendregistry.StatusAll(context.Background(), cfg)
	payload := map[string]any{
		"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
		"secret": map[string]any{
			"mode":               "native-scaffold",
			"compat_backend":     cfg.CompatBackend,
			"secret_service":     cfg.SecretService,
			"available_backends": backends,
			"profile_readiness":  profilereadiness.AssessAll(cfg),
			"broker":             probe.Broker,
			"environment":        probe.Environment,
			"launchd":            probe.Launchd,
			"keychain":           probe.Keychain,
			"legacy_runtime":     probe.LegacyRuntime,
			"llm_profiles":       cfg.LLMProfiles,
		},
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runSecretAudit(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	probe := runtimeprobe.Collect(cfg)
	backends := secretbackendregistry.StatusAll(context.Background(), cfg)

	findings := []map[string]string{}
	status := "ok"

	if !probe.Broker.SocketExists {
		status = "needs_runtime"
		findings = append(findings, map[string]string{
			"severity": "high",
			"kind":     "broker_socket_missing",
			"detail":   "Local broker socket is missing; unattended secret service is not ready.",
		})
	}
	if probe.Broker.SocketExists && !probe.Broker.Reachable {
		status = "degraded"
		findings = append(findings, map[string]string{
			"severity": "high",
			"kind":     "broker_unreachable",
			"detail":   fallbackString(probe.Broker.Error, "Local broker socket exists but did not return a healthy response."),
		})
	}
	if !probe.Environment.ServiceAccountEnvValid {
		if probe.Broker.Status == nil || !probe.Broker.Status.TokenAvailable {
			status = "needs_bootstrap"
			findings = append(findings, map[string]string{
				"severity": "high",
				"kind":     "service_account_runtime_missing",
				"detail":   "No valid service-account token was found in the process environment or broker runtime state.",
			})
		}
	}

	payload := map[string]any{
		"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
		"secret_audit": map[string]any{
			"mode":               "native-scaffold",
			"status":             status,
			"compat_backend":     cfg.CompatBackend,
			"secret_service":     cfg.SecretService,
			"available_backends": backends,
			"profile_readiness":  profilereadiness.AssessAll(cfg),
			"probe":              probe,
			"findings":           findings,
		},
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runSecretRead(args []string) int {
	asJSON := false
	unsafeReveal := false
	trimmed := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--json" {
			asJSON = true
			continue
		}
		if arg == "--unsafe-reveal" {
			unsafeReveal = true
			continue
		}
		trimmed = append(trimmed, arg)
	}
	if len(trimmed) != 1 {
		fmt.Fprintln(a.Stderr, "aih-go secret read: usage: aih secret read <op://...> [--json] [--unsafe-reveal]")
		return 2
	}
	if !unsafeReveal && os.Getenv("AIH_UNSAFE_REVEAL") != "1" {
		fmt.Fprintln(a.Stderr, "aih-go secret read: plaintext reveal is disabled by default; set AIH_UNSAFE_REVEAL=1 or pass --unsafe-reveal in admin mode")
		return 2
	}

	cfg := config.Load()
	result, err := secrets.ResolveReferenceNative(context.Background(), cfg, trimmed[0])
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go secret read: %v\n", err)
		return 2
	}

	if asJSON {
		return writeJSON(a.Stdout, map[string]any{
			"reference": trimmed[0],
			"value":     result.Value,
			"source":    result.Source,
		})
	}
	fmt.Fprint(a.Stdout, result.Value)
	return 0
}

func (a *App) runSecretGet(args []string) int {
	asJSON := false
	unsafeReveal := false
	trimmed := make([]string, 0, len(args))
	for _, arg := range args {
		switch arg {
		case "--json":
			asJSON = true
		case "--unsafe-reveal":
			unsafeReveal = true
		default:
			trimmed = append(trimmed, arg)
		}
	}
	if len(trimmed) != 1 {
		fmt.Fprintln(a.Stderr, "aih-go secret get: usage: aih secret get <alias> [--json] [--unsafe-reveal]")
		return 2
	}

	cfg := config.Load()
	entry, result, err := secrets.ResolveAlias(context.Background(), cfg, trimmed[0])
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go secret get: %v\n", err)
		return 2
	}
	if !unsafeReveal && os.Getenv("AIH_UNSAFE_REVEAL") != "1" && entry.RevealPolicy != secretpolicy.RevealAllowed {
		fmt.Fprintf(a.Stderr, "aih-go secret get: alias %q is not revealable by default (reveal_policy=%s)\n", entry.Name, entry.RevealPolicy)
		return 2
	}

	if asJSON {
		return writeJSON(a.Stdout, map[string]any{
			"alias":         entry.Name,
			"value":         result.Value,
			"source":        result.Source,
			"backend":       entry.Backend,
			"reveal_policy": entry.RevealPolicy,
			"usage_policy":  entry.UsagePolicy,
		})
	}
	fmt.Fprint(a.Stdout, result.Value)
	return 0
}

func (a *App) runSecretList(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	aliases := map[string]map[string]any{}
	for _, entry := range cfg.SecretAliases {
		aliases[entry.Name] = map[string]any{
			"backend":         entry.Backend,
			"reference":       entry.Reference,
			"category":        entry.Category,
			"reveal_policy":   entry.RevealPolicy,
			"usage_policy":    entry.UsagePolicy,
			"allowed_actions": entry.AllowedActions,
			"metadata":        entry.Metadata,
		}
	}
	payload := map[string]any{
		"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
		"secret": map[string]any{
			"aliases": aliases,
		},
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runSecretEnv(args []string) int {
	asJSON := false
	unsafeInject := false
	for _, arg := range args {
		switch arg {
		case "--json":
			asJSON = true
		case "--unsafe-inject":
			unsafeInject = true
		default:
			fmt.Fprintf(a.Stderr, "aih-go secret env: unknown flag %s\n", arg)
			return 2
		}
	}
	if !unsafeInject && os.Getenv("AIH_UNSAFE_INJECT") != "1" {
		fmt.Fprintln(a.Stderr, "aih-go secret env: env injection is disabled by default; set AIH_UNSAFE_INJECT=1 or pass --unsafe-inject in admin mode")
		return 2
	}

	cfg := config.Load()
	token, source, err := runtimeauth.ResolveServiceAccountToken(cfg)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go secret env: %v\n", err)
		return 2
	}
	payload := map[string]any{
		"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
		"secret_env": map[string]any{
			"success": true,
			"source":  source,
			"env": map[string]string{
				cfg.SecretService.ServiceAccountEnv: token,
			},
		},
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	fmt.Fprintf(a.Stdout, "export %s=%q\n", cfg.SecretService.ServiceAccountEnv, token)
	return 0
}

func (a *App) runSecretExec(args []string) int {
	unsafeInject := false
	command := make([]string, 0, len(args))
	for _, arg := range args {
		switch {
		case arg == "--unsafe-inject":
			unsafeInject = true
		default:
			command = append(command, arg)
		}
	}
	if len(command) > 0 && command[0] == "--" {
		command = command[1:]
	}
	if len(command) == 0 {
		fmt.Fprintln(a.Stderr, "aih-go secret exec: provide a command after `aih secret exec -- ...`")
		return 2
	}
	if !unsafeInject && os.Getenv("AIH_UNSAFE_INJECT") != "1" {
		fmt.Fprintln(a.Stderr, "aih-go secret exec: env injection is disabled by default; set AIH_UNSAFE_INJECT=1 or pass --unsafe-inject in admin mode")
		return 2
	}

	cfg := config.Load()
	token, _, err := runtimeauth.ResolveServiceAccountToken(cfg)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go secret exec: %v\n", err)
		return 2
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = a.Stdout
	cmd.Stderr = a.Stderr
	cmd.Env = append(os.Environ(), cfg.SecretService.ServiceAccountEnv+"="+token)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(a.Stderr, "aih-go secret exec: %v\n", err)
		return 1
	}
	return 0
}

func (a *App) runSecretSudo(args []string) int {
	timeout := 60 * time.Second
	command := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--timeout":
			if i+1 >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go secret sudo: missing value after --timeout")
				return 2
			}
			parsed, err := time.ParseDuration(args[i+1] + "s")
			if err != nil {
				fmt.Fprintf(a.Stderr, "aih-go secret sudo: invalid timeout %q\n", args[i+1])
				return 2
			}
			timeout = parsed
			i++
		default:
			command = append(command, args[i])
		}
	}
	if len(command) > 0 && command[0] == "--" {
		command = command[1:]
	}
	if len(command) == 0 {
		fmt.Fprintln(a.Stderr, "aih-go secret sudo: provide a command after `aih secret sudo -- ...`")
		return 2
	}

	cfg := config.Load()
	_, result, err := secrets.ResolveAlias(context.Background(), cfg, "account-password")
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go secret sudo: %v\n", err)
		return 2
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sudo", append([]string{"-S", "-p", "", "-k"}, command...)...)
	cmd.Stdin = strings.NewReader(result.Value + "\n")
	cmd.Stdout = a.Stdout
	cmd.Stderr = a.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Fprintf(a.Stderr, "aih-go secret sudo: command timed out after %s\n", timeout)
			return 2
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(a.Stderr, "aih-go secret sudo: %v\n", err)
		return 1
	}
	return 0
}

func (a *App) runSecretBootstrapToken(args []string) int {
	asJSON := false
	noVerify := false
	for _, arg := range args {
		switch arg {
		case "--json":
			asJSON = true
		case "--no-verify":
			noVerify = true
		default:
			fmt.Fprintf(a.Stderr, "aih-go secret bootstrap-token: unknown flag %s\n", arg)
			return 2
		}
	}

	rawToken := strings.TrimSpace(os.Getenv("OP_SERVICE_ACCOUNT_TOKEN"))
	source := ""
	if rawToken != "" {
		source = "env"
	} else if stat, _ := os.Stdin.Stat(); stat != nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err == nil {
			rawToken = strings.TrimSpace(string(data))
			if rawToken != "" {
				source = "stdin"
			}
		}
	}
	if !runtimeauth.TokenLooksValid(rawToken) {
		fmt.Fprintln(a.Stderr, "aih-go secret bootstrap-token: provide a full 1Password service-account token via OP_SERVICE_ACCOUNT_TOKEN or stdin")
		return 2
	}

	if err := keychain.UpsertGenericPassword(context.Background(), keychain.ServiceAccountTokenService, keychain.CurrentUser(), rawToken); err != nil {
		fmt.Fprintf(a.Stderr, "aih-go secret bootstrap-token: %v\n", err)
		return 2
	}

	cfg := config.Load()
	restartErr := restartBrokerLaunchAgent(cfg.Broker.LaunchAgentLabel)
	payload := map[string]any{
		"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
		"secret_bootstrap_token": map[string]any{
			"success":          restartErr == nil,
			"source":           source,
			"keychain_service": keychain.ServiceAccountTokenService,
			"broker_label":     cfg.Broker.LaunchAgentLabel,
			"token_length":     len(rawToken),
			"error":            errorString(restartErr),
		},
	}
	if !noVerify {
		payload["verification"] = runtimeprobe.Collect(cfg)
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runSecretCache(args []string) int {
	asJSON := false
	trimmed := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--json" {
			asJSON = true
			continue
		}
		trimmed = append(trimmed, arg)
	}
	if len(trimmed) != 1 {
		fmt.Fprintln(a.Stderr, "aih-go secret cache: usage: aih secret cache <alias> [--json]")
		return 2
	}

	cfg := config.Load()
	entry, result, err := secrets.ResolveAlias(context.Background(), cfg, trimmed[0])
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go secret cache: %v\n", err)
		return 2
	}
	service := defaultString(entry.Metadata["keychain_service"], "aih."+entry.Name)
	account := defaultString(entry.Metadata["keychain_account"], keychain.CurrentUser())
	source := result.Source
	if cfg.Broker.SocketPath != "" {
		client := broker.Client{SocketPath: cfg.Broker.SocketPath}
		response, cacheErr := client.CacheMaterial(entry.Reference, service, account)
		if cacheErr == nil && response.OK {
			source = response.MaterialSource
		} else {
			if err := keychain.UpsertGenericPassword(context.Background(), service, account, result.Value); err != nil {
				fmt.Fprintf(a.Stderr, "aih-go secret cache: %v\n", err)
				return 2
			}
		}
	} else {
		if err := keychain.UpsertGenericPassword(context.Background(), service, account, result.Value); err != nil {
			fmt.Fprintf(a.Stderr, "aih-go secret cache: %v\n", err)
			return 2
		}
	}

	payload := map[string]any{
		"generated_at_utc": time.Now().UTC().Format(time.RFC3339),
		"secret_cache": map[string]any{
			"alias":            entry.Name,
			"success":          true,
			"backend":          entry.Backend,
			"source":           source,
			"keychain_service": service,
			"keychain_account": account,
			"reference_hint":   fmt.Sprintf("keychain://%s?account=%s", service, account),
			"value_length":     len(result.Value),
		},
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runLLM(args []string) int {
	if len(args) == 0 {
		a.printLLMHelp()
		return 0
	}

	switch args[0] {
	case "profiles":
		return a.runLLMProfiles(args[1:])
	case "verify":
		return a.runLLMVerify(args[1:])
	case "request":
		return a.runLLMRequest(args[1:])
	case "help", "-h", "--help":
		a.printLLMHelp()
		return 0
	default:
		fmt.Fprintf(a.Stderr, "aih-go llm: unknown subcommand: %s\n", args[0])
		a.printLLMHelp()
		return 2
	}
}

func (a *App) runLLMProfiles(args []string) int {
	asJSON := len(args) > 0 && args[0] == "--json"
	cfg := config.Load()
	payload := map[string]any{
		"configured":      cfg.LLMProfiles.Configured,
		"primary_profile": cfg.LLMProfiles.PrimaryProfile,
		"profiles":        cfg.LLMProfiles.Profiles,
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runLLMVerify(args []string) int {
	asJSON := false
	profileName := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			asJSON = true
		case "--profile":
			if i+1 >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go llm verify: missing value after --profile")
				return 2
			}
			profileName = args[i+1]
			i++
		default:
			if strings.HasPrefix(args[i], "--") {
				fmt.Fprintf(a.Stderr, "aih-go llm verify: unknown flag %s\n", args[i])
				return 2
			}
		}
	}

	cfg := config.Load()
	prof, err := cfg.FindLLMProfile(profileName)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go llm verify: %v\n", err)
		return 2
	}
	readiness, readinessErr := profilereadiness.Find(cfg, prof.Name)
	if readinessErr == nil && !readiness.Ready {
		fmt.Fprintf(a.Stderr, "aih-go llm verify: profile %q is not ready: %v\n", prof.Name, readiness.Missing)
		return 2
	}
	if !prof.AllowsAction("llm.verify") {
		fmt.Fprintf(a.Stderr, "aih-go llm verify: profile %q does not allow action %q\n", prof.Name, "llm.verify")
		return 2
	}
	if string(prof.UsagePolicy) != "opaque_use" {
		fmt.Fprintf(a.Stderr, "aih-go llm verify: profile %q is not configured for opaque-use verification\n", prof.Name)
		return 2
	}

	secretMaterial, err := secrets.ResolveProfileMaterial(context.Background(), cfg, prof)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go llm verify: failed to resolve profile credential: %v\n", err)
		return 2
	}

	baseURL, err := secrets.ResolveProfileBaseURL(context.Background(), cfg, prof)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go llm verify: failed to resolve profile base URL: %v\n", err)
		return 2
	}
	prof.BaseURL = baseURL

	marker := "AIH_GO_LLM_VERIFY_OK"
	result, err := llmapi.VerifyOpaque(prof, secretMaterial.Value, marker)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go llm verify: %v\n", err)
		return 2
	}

	payload := map[string]any{
		"profile":           prof.Name,
		"category":          prof.Category,
		"reveal_policy":     prof.RevealPolicy,
		"usage_policy":      prof.UsagePolicy,
		"allowed_actions":   prof.AllowedActions,
		"credential_source": secretMaterial.Source,
		"result":            result,
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func (a *App) runLLMRequest(args []string) int {
	asJSON := false
	profileName := ""
	bodyFile := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			asJSON = true
		case "--profile":
			if i+1 >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go llm request: missing value after --profile")
				return 2
			}
			profileName = args[i+1]
			i++
		case "--body-file":
			if i+1 >= len(args) {
				fmt.Fprintln(a.Stderr, "aih-go llm request: missing value after --body-file")
				return 2
			}
			bodyFile = args[i+1]
			i++
		default:
			if strings.HasPrefix(args[i], "--") {
				fmt.Fprintf(a.Stderr, "aih-go llm request: unknown flag %s\n", args[i])
				return 2
			}
		}
	}
	if bodyFile == "" {
		fmt.Fprintln(a.Stderr, "aih-go llm request: --body-file is required")
		return 2
	}

	cfg := config.Load()
	prof, err := cfg.FindLLMProfile(profileName)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go llm request: %v\n", err)
		return 2
	}
	readiness, readinessErr := profilereadiness.Find(cfg, prof.Name)
	if readinessErr == nil && !readiness.Ready {
		fmt.Fprintf(a.Stderr, "aih-go llm request: profile %q is not ready: %v\n", prof.Name, readiness.Missing)
		return 2
	}
	if !prof.AllowsAction("llm.request") {
		fmt.Fprintf(a.Stderr, "aih-go llm request: profile %q does not allow action %q\n", prof.Name, "llm.request")
		return 2
	}
	if string(prof.UsagePolicy) != "opaque_use" && string(prof.UsagePolicy) != "inject_env" {
		fmt.Fprintf(a.Stderr, "aih-go llm request: profile %q is not configured for delegated request usage\n", prof.Name)
		return 2
	}

	secretMaterial, err := secrets.ResolveProfileMaterial(context.Background(), cfg, prof)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go llm request: failed to resolve profile credential: %v\n", err)
		return 2
	}
	baseURL, err := secrets.ResolveProfileBaseURL(context.Background(), cfg, prof)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go llm request: failed to resolve profile base URL: %v\n", err)
		return 2
	}
	prof.BaseURL = baseURL

	data, err := os.ReadFile(bodyFile)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go llm request: read body file: %v\n", err)
		return 2
	}
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		fmt.Fprintf(a.Stderr, "aih-go llm request: parse body file: %v\n", err)
		return 2
	}

	result, err := llmapi.RequestOpaque(prof, secretMaterial.Value, body)
	if err != nil {
		fmt.Fprintf(a.Stderr, "aih-go llm request: %v\n", err)
		return 2
	}
	payload := map[string]any{
		"profile":           prof.Name,
		"category":          prof.Category,
		"reveal_policy":     prof.RevealPolicy,
		"usage_policy":      prof.UsagePolicy,
		"allowed_actions":   prof.AllowedActions,
		"credential_source": secretMaterial.Source,
		"result":            result,
	}
	if asJSON {
		return writeJSON(a.Stdout, payload)
	}
	return writeJSON(a.Stdout, payload)
}

func fallbackString(value string, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return fallback
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func restartBrokerLaunchAgent(label string) error {
	if runtime.GOOS != "darwin" {
		return nil
	}
	if label == "" {
		return fmt.Errorf("broker launch agent label is empty")
	}

	uid := os.Getuid()
	target := fmt.Sprintf("gui/%d/%s", uid, label)
	plistPath := filepath.Join(paths.HomeDir(), "Library", "LaunchAgents", label+".plist")

	bootstrap := exec.Command("launchctl", "bootstrap", fmt.Sprintf("gui/%d", uid), plistPath)
	if output, err := bootstrap.CombinedOutput(); err != nil {
		text := strings.ToLower(string(output))
		if !strings.Contains(text, "already bootstrapped") {
			return fmt.Errorf("launchctl bootstrap %s: %w", label, err)
		}
	}
	kickstart := exec.Command("launchctl", "kickstart", "-k", target)
	if output, err := kickstart.CombinedOutput(); err != nil {
		return fmt.Errorf("launchctl kickstart %s: %w (%s)", label, err, strings.TrimSpace(string(output)))
	}
	return nil
}

func writeJSON(w io.Writer, payload any) int {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(payload); err != nil {
		fmt.Fprintf(os.Stderr, "aih-go: failed to encode JSON: %v\n", err)
		return 1
	}
	return 0
}
