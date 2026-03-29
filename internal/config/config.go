package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/paths"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/profile"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretpolicy"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretregistry"
)

type Config struct {
	SourceRepoRoot string                 `json:"source_repo_root"`
	AgentsDir      string                 `json:"agents_dir"`
	LegacyRuntime  string                 `json:"legacy_runtime"`
	CompatBackend  string                 `json:"compat_backend"`
	Broker         BrokerConfig           `json:"broker"`
	Doctor         DoctorConfig           `json:"doctor"`
	Service        ServiceConfig          `json:"service"`
	Auth           AuthConfig             `json:"auth"`
	Browser        BrowserConfig          `json:"browser"`
	DevOps         DevOpsConfig           `json:"devops"`
	SecretService  SecretServiceConfig    `json:"secret_service"`
	SecretAliases  []secretregistry.Entry `json:"secret_aliases"`
	LLMProfiles    LLMProfilesConfig      `json:"llm_profiles"`
}

type BrokerConfig struct {
	SocketPath       string `json:"socket_path"`
	LaunchAgentLabel string `json:"launch_agent_label"`
}

type DoctorConfig struct {
	FactsMaxAgeHours         int               `json:"facts_max_age_hours"`
	PathIgnoreEntries        []string          `json:"path_ignore_entries,omitempty"`
	PathIgnorePrefixes       []string          `json:"path_ignore_prefixes,omitempty"`
	PathIgnoreContains       []string          `json:"path_ignore_contains,omitempty"`
	IssueIgnoreKinds         []string          `json:"issue_ignore_kinds,omitempty"`
	IssueIgnoreSources       []string          `json:"issue_ignore_sources,omitempty"`
	IgnoreCurrentSessionOnly bool              `json:"ignore_current_session_only"`
	IssueSeverityOverrides   map[string]string `json:"issue_severity_overrides,omitempty"`
}

type ServiceConfig struct {
	PortProbes []PortProbeConfig `json:"port_probes"`
}

type PortProbeConfig struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Required bool   `json:"required"`
}

type AuthConfig struct {
	Tools []AuthToolConfig `json:"tools"`
}

type AuthToolConfig struct {
	Name      string `json:"name"`
	Binary    string `json:"binary"`
	ConfigDir string `json:"config_dir,omitempty"`
	AuthFile  string `json:"auth_file,omitempty"`
	Required  bool   `json:"required"`
}

type BrowserConfig struct {
	CDPPortDefault       int    `json:"cdp_port_default"`
	ChromeAppPath        string `json:"chrome_app_path"`
	ChromeBinaryPath     string `json:"chrome_binary_path"`
	AutomationProfileDir string `json:"automation_profile_dir"`
	PlaywrightHelperPath string `json:"playwright_helper_path"`
	PlaywrightCoreDir    string `json:"playwright_core_dir"`
	DoctorChecksEnabled  bool   `json:"doctor_checks_enabled"`
}

type DevOpsConfig struct {
	Tools            []DevOpsToolConfig `json:"tools"`
	DockerConfigPath string             `json:"docker_config_path,omitempty"`
}

type DevOpsToolConfig struct {
	Name        string   `json:"name"`
	Binary      string   `json:"binary"`
	VersionArgs []string `json:"version_args,omitempty"`
	Required    bool     `json:"required"`
}

type SecretServiceConfig struct {
	Kind              string `json:"kind"`
	Provider          string `json:"provider"`
	ServiceAccountEnv string `json:"service_account_env"`
	VaultName         string `json:"vault_name"`
	VaultUUID         string `json:"vault_uuid"`
}

type LLMProfilesConfig struct {
	Configured     bool              `json:"configured"`
	PrimaryProfile string            `json:"primary_profile"`
	Profiles       []profile.Profile `json:"profiles"`
}

func Load() Config {
	cfg := defaultConfig()
	if fileCfg, ok := loadFileConfig(); ok {
		cfg = mergeConfig(cfg, fileCfg)
	}
	cfg = applyEnvOverrides(cfg)
	return cfg
}

func defaultConfig() Config {
	agentsDir := paths.AgentsDir()
	sourceRepoRoot := paths.SourceRepoRoot()
	runtimeBrowserDir := defaultRuntimeBrowserDir(sourceRepoRoot, agentsDir)
	return Config{
		SourceRepoRoot: sourceRepoRoot,
		AgentsDir:      agentsDir,
		LegacyRuntime:  defaultLegacyRuntime(agentsDir),
		CompatBackend:  "legacy-python",
		Broker: BrokerConfig{
			SocketPath:       filepath.Join(agentsDir, "run", "op-sa-broker.sock"),
			LaunchAgentLabel: "com.ai-agent-owner.op-sa-broker",
		},
		Doctor: DoctorConfig{
			FactsMaxAgeHours:         24,
			PathIgnoreEntries:        []string{},
			PathIgnorePrefixes:       []string{"/var/run/com.apple.security.cryptexd/", "/Library/Apple/usr/bin"},
			PathIgnoreContains:       []string{},
			IssueIgnoreKinds:         []string{},
			IssueIgnoreSources:       []string{},
			IgnoreCurrentSessionOnly: true,
			IssueSeverityOverrides:   map[string]string{},
		},
		Service: ServiceConfig{
			PortProbes: []PortProbeConfig{
				{Name: "ssh_22", Host: "127.0.0.1", Port: 22, Required: false},
				{Name: "browser_cdp_9222", Host: "127.0.0.1", Port: 9222, Required: false},
				{Name: "postgres_5432", Host: "127.0.0.1", Port: 5432, Required: false},
			},
		},
		Auth: AuthConfig{
			Tools: []AuthToolConfig{
				{Name: "codex", Binary: "codex", ConfigDir: filepath.Join(paths.HomeDir(), ".codex"), AuthFile: filepath.Join(paths.HomeDir(), ".codex", "auth.json"), Required: false},
				{Name: "claude", Binary: "claude", ConfigDir: filepath.Join(paths.HomeDir(), ".claude"), AuthFile: filepath.Join(paths.HomeDir(), ".claude", "api-key.txt"), Required: false},
				{Name: "gemini", Binary: "gemini", ConfigDir: filepath.Join(paths.HomeDir(), ".config", "gemini"), Required: false},
				{Name: "github_cli", Binary: "gh", ConfigDir: filepath.Join(paths.HomeDir(), ".config", "gh"), AuthFile: filepath.Join(paths.HomeDir(), ".config", "gh", "hosts.yml"), Required: false},
				{Name: "onepassword", Binary: "op", Required: false},
			},
		},
		Browser: BrowserConfig{
			CDPPortDefault:       9222,
			ChromeAppPath:        "/Applications/Google Chrome.app",
			ChromeBinaryPath:     "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			AutomationProfileDir: filepath.Join(agentsDir, "state", "browser", "chrome-cdp-profile"),
			PlaywrightHelperPath: filepath.Join(runtimeBrowserDir, "verify-playwright.mjs"),
			PlaywrightCoreDir:    filepath.Join(runtimeBrowserDir, "node_modules", "playwright-core"),
			DoctorChecksEnabled:  false,
		},
		DevOps: DevOpsConfig{
			Tools: []DevOpsToolConfig{
				{Name: "docker", Binary: "docker", VersionArgs: []string{"--version"}, Required: false},
				{Name: "kubectl", Binary: "kubectl", VersionArgs: []string{"version", "--client=true"}, Required: false},
				{Name: "terraform", Binary: "terraform", VersionArgs: []string{"--version"}, Required: false},
				{Name: "aws", Binary: "aws", VersionArgs: []string{"--version"}, Required: false},
				{Name: "gcloud", Binary: "gcloud", VersionArgs: []string{"--version"}, Required: false},
				{Name: "az", Binary: "az", VersionArgs: []string{"--version"}, Required: false},
				{Name: "helm", Binary: "helm", VersionArgs: []string{"version", "--short"}, Required: false},
				{Name: "cloudflared", Binary: "cloudflared", VersionArgs: []string{"--version"}, Required: false},
			},
			DockerConfigPath: filepath.Join(paths.HomeDir(), ".docker", "config.json"),
		},
		SecretService: SecretServiceConfig{
			Kind:              "1password-service-account",
			Provider:          "1password",
			ServiceAccountEnv: "OP_SERVICE_ACCOUNT_TOKEN",
		},
		SecretAliases: []secretregistry.Entry{},
		LLMProfiles: LLMProfilesConfig{
			Configured:     false,
			PrimaryProfile: "default",
			Profiles: []profile.Profile{
				{
					Name:           "default",
					Category:       secretpolicy.CategoryLLMAPI,
					Backend:        "1password-service-account",
					SecretAlias:    "llm-api-key",
					RevealPolicy:   secretpolicy.RevealNever,
					UsagePolicy:    secretpolicy.UsageOpaqueUse,
					AllowedActions: []secretpolicy.Action{secretpolicy.ActionLLMVerify, secretpolicy.ActionLLMRequest, secretpolicy.ActionLLMModels},
					DefaultModels:  []string{},
					BaseURLAlias:   "llm-base-url",
					Metadata: map[string]string{
						"protocol":    "openai_chat_completions",
						"auth_scheme": "bearer",
					},
				},
			},
		},
	}
}

func defaultLegacyRuntime(agentsDir string) string {
	candidate := filepath.Join(agentsDir, "tools", "aih.py")
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return ""
}

func defaultRuntimeBrowserDir(sourceRepoRoot, agentsDir string) string {
	candidates := []string{
		filepath.Join(sourceRepoRoot, "dist", "current", "runtime", "browser"),
		filepath.Join(sourceRepoRoot, "runtime", "browser"),
		filepath.Join(agentsDir, "runtime", "browser"),
		filepath.Join(agentsDir, "tools", "browser"),
	}
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return filepath.Join(sourceRepoRoot, "dist", "current", "runtime", "browser")
}

func applyEnvOverrides(cfg Config) Config {
	cfg.SourceRepoRoot = defaultString(os.Getenv("AIH_SOURCE_ROOT"), cfg.SourceRepoRoot)
	cfg.LegacyRuntime = defaultString(os.Getenv("AIH_LEGACY_RUNTIME"), cfg.LegacyRuntime)
	cfg.CompatBackend = defaultString(os.Getenv("AIH_COMPAT_BACKEND"), cfg.CompatBackend)
	cfg.Broker.SocketPath = defaultString(os.Getenv("AIH_BROKER_SOCKET"), cfg.Broker.SocketPath)
	cfg.Broker.LaunchAgentLabel = defaultString(os.Getenv("AIH_BROKER_LABEL"), cfg.Broker.LaunchAgentLabel)
	cfg.SecretService.Kind = defaultString(os.Getenv("AIH_SECRET_BACKEND"), cfg.SecretService.Kind)
	cfg.SecretService.ServiceAccountEnv = defaultString(os.Getenv("AIH_SERVICE_ACCOUNT_ENV"), cfg.SecretService.ServiceAccountEnv)
	cfg.SecretService.VaultName = defaultString(os.Getenv("AIH_VAULT_NAME"), cfg.SecretService.VaultName)
	cfg.SecretService.VaultUUID = defaultString(os.Getenv("AIH_VAULT_UUID"), cfg.SecretService.VaultUUID)
	cfg.LLMProfiles.PrimaryProfile = defaultString(os.Getenv("AIH_LLM_PROFILE"), cfg.LLMProfiles.PrimaryProfile)
	if len(cfg.LLMProfiles.Profiles) > 0 {
		cfg.LLMProfiles.Profiles[0].Name = defaultString(os.Getenv("AIH_LLM_PROFILE"), cfg.LLMProfiles.Profiles[0].Name)
		cfg.LLMProfiles.Profiles[0].BaseURL = defaultString(os.Getenv("AIH_LLM_BASE_URL"), cfg.LLMProfiles.Profiles[0].BaseURL)
		cfg.LLMProfiles.Profiles[0].ConverterBaseURL = defaultString(os.Getenv("AIH_LLM_CONVERTER_BASE_URL"), cfg.LLMProfiles.Profiles[0].ConverterBaseURL)
		cfg.LLMProfiles.Profiles[0].SecretAlias = defaultString(os.Getenv("AIH_LLM_API_KEY_ALIAS"), cfg.LLMProfiles.Profiles[0].SecretAlias)
		cfg.LLMProfiles.Profiles[0].BaseURLAlias = defaultString(os.Getenv("AIH_LLM_BASE_URL_ALIAS"), cfg.LLMProfiles.Profiles[0].BaseURLAlias)
		if cfg.LLMProfiles.Profiles[0].Metadata == nil {
			cfg.LLMProfiles.Profiles[0].Metadata = map[string]string{}
		}
		cfg.LLMProfiles.Profiles[0].Metadata["protocol"] = defaultString(os.Getenv("AIH_LLM_PROTOCOL"), cfg.LLMProfiles.Profiles[0].Metadata["protocol"])
		cfg.LLMProfiles.Profiles[0].Metadata["auth_scheme"] = defaultString(os.Getenv("AIH_LLM_AUTH_SCHEME"), cfg.LLMProfiles.Profiles[0].Metadata["auth_scheme"])
	}
	return cfg
}

func loadFileConfig() (Config, bool) {
	loaded := false
	merged := Config{}
	for _, path := range configSearchPaths() {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var cfg Config
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}
		if !loaded {
			merged = cfg
			loaded = true
			continue
		}
		merged = mergeConfig(merged, cfg)
	}
	return merged, loaded
}

func configSearchPaths() []string {
	home := paths.HomeDir()
	candidates := []string{}
	if explicit := os.Getenv("AIH_CONFIG_FILE"); explicit != "" {
		candidates = append(candidates, explicit)
	}
	if home != "" {
		candidates = append(candidates,
			filepath.Join(home, ".config", "aih", "config.json"),
			filepath.Join(home, ".agents", "aih", "config.json"),
		)
	}
	candidates = append(candidates, repoConfigSearchPaths()...)
	return candidates
}

func repoConfigSearchPaths() []string {
	wd := os.Getenv("AIH_CALLER_CWD")
	if wd == "" {
		var err error
		wd, err = os.Getwd()
		if err != nil {
			return nil
		}
	}
	if wd == "" {
		return nil
	}

	var candidates []string
	current := wd
	visited := map[string]struct{}{}
	for {
		if _, ok := visited[current]; ok {
			break
		}
		visited[current] = struct{}{}
		candidates = append(candidates,
			filepath.Join(current, ".aih", "config.json"),
			filepath.Join(current, ".codex", "aih", "config.json"),
		)
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	slices.Reverse(candidates)
	return candidates
}

func mergeConfig(base Config, override Config) Config {
	if override.SourceRepoRoot != "" {
		base.SourceRepoRoot = override.SourceRepoRoot
	}
	if override.AgentsDir != "" {
		base.AgentsDir = override.AgentsDir
	}
	if override.LegacyRuntime != "" {
		base.LegacyRuntime = override.LegacyRuntime
	}
	if override.CompatBackend != "" {
		base.CompatBackend = override.CompatBackend
	}
	if override.Broker.SocketPath != "" {
		base.Broker.SocketPath = override.Broker.SocketPath
	}
	if override.Broker.LaunchAgentLabel != "" {
		base.Broker.LaunchAgentLabel = override.Broker.LaunchAgentLabel
	}
	if override.Doctor.FactsMaxAgeHours != 0 {
		base.Doctor.FactsMaxAgeHours = override.Doctor.FactsMaxAgeHours
	}
	if len(override.Doctor.PathIgnorePrefixes) > 0 {
		base.Doctor.PathIgnorePrefixes = append([]string{}, override.Doctor.PathIgnorePrefixes...)
	}
	if len(override.Doctor.PathIgnoreEntries) > 0 {
		base.Doctor.PathIgnoreEntries = append([]string{}, override.Doctor.PathIgnoreEntries...)
	}
	if len(override.Doctor.PathIgnoreContains) > 0 {
		base.Doctor.PathIgnoreContains = append([]string{}, override.Doctor.PathIgnoreContains...)
	}
	if len(override.Doctor.IssueIgnoreKinds) > 0 {
		base.Doctor.IssueIgnoreKinds = append([]string{}, override.Doctor.IssueIgnoreKinds...)
	}
	if len(override.Doctor.IssueIgnoreSources) > 0 {
		base.Doctor.IssueIgnoreSources = append([]string{}, override.Doctor.IssueIgnoreSources...)
	}
	if len(override.Doctor.IssueSeverityOverrides) > 0 {
		if base.Doctor.IssueSeverityOverrides == nil {
			base.Doctor.IssueSeverityOverrides = map[string]string{}
		}
		for key, value := range override.Doctor.IssueSeverityOverrides {
			base.Doctor.IssueSeverityOverrides[key] = value
		}
	}
	if override.Doctor.IgnoreCurrentSessionOnly {
		base.Doctor.IgnoreCurrentSessionOnly = true
	}
	if len(override.Service.PortProbes) > 0 {
		base.Service.PortProbes = mergePortProbes(base.Service.PortProbes, override.Service.PortProbes)
	}
	if len(override.Auth.Tools) > 0 {
		base.Auth.Tools = mergeAuthTools(base.Auth.Tools, override.Auth.Tools)
	}
	if override.Browser.CDPPortDefault != 0 {
		base.Browser.CDPPortDefault = override.Browser.CDPPortDefault
	}
	if override.Browser.ChromeAppPath != "" {
		base.Browser.ChromeAppPath = override.Browser.ChromeAppPath
	}
	if override.Browser.ChromeBinaryPath != "" {
		base.Browser.ChromeBinaryPath = override.Browser.ChromeBinaryPath
	}
	if override.Browser.AutomationProfileDir != "" {
		base.Browser.AutomationProfileDir = override.Browser.AutomationProfileDir
	}
	if override.Browser.PlaywrightHelperPath != "" {
		base.Browser.PlaywrightHelperPath = override.Browser.PlaywrightHelperPath
	}
	if override.Browser.PlaywrightCoreDir != "" {
		base.Browser.PlaywrightCoreDir = override.Browser.PlaywrightCoreDir
	}
	if override.Browser.DoctorChecksEnabled {
		base.Browser.DoctorChecksEnabled = true
	}
	if len(override.DevOps.Tools) > 0 {
		base.DevOps.Tools = mergeDevOpsTools(base.DevOps.Tools, override.DevOps.Tools)
	}
	if override.DevOps.DockerConfigPath != "" {
		base.DevOps.DockerConfigPath = override.DevOps.DockerConfigPath
	}
	if override.SecretService.Kind != "" {
		base.SecretService.Kind = override.SecretService.Kind
	}
	if override.SecretService.Provider != "" {
		base.SecretService.Provider = override.SecretService.Provider
	}
	if override.SecretService.ServiceAccountEnv != "" {
		base.SecretService.ServiceAccountEnv = override.SecretService.ServiceAccountEnv
	}
	if override.SecretService.VaultName != "" {
		base.SecretService.VaultName = override.SecretService.VaultName
	}
	if override.SecretService.VaultUUID != "" {
		base.SecretService.VaultUUID = override.SecretService.VaultUUID
	}
	if len(override.SecretAliases) > 0 {
		base.SecretAliases = mergeSecretAliases(base.SecretAliases, override.SecretAliases)
	}
	if override.LLMProfiles.PrimaryProfile != "" {
		base.LLMProfiles.PrimaryProfile = override.LLMProfiles.PrimaryProfile
	}
	if len(override.LLMProfiles.Profiles) > 0 {
		base.LLMProfiles.Profiles = mergeProfiles(base.LLMProfiles.Profiles, override.LLMProfiles.Profiles)
		base.LLMProfiles.Configured = true
	}
	return base
}

func mergeSecretAliases(base []secretregistry.Entry, override []secretregistry.Entry) []secretregistry.Entry {
	order := make([]string, 0, len(base)+len(override))
	index := map[string]secretregistry.Entry{}
	for _, entry := range base {
		order = append(order, entry.Name)
		index[entry.Name] = entry
	}
	for _, entry := range override {
		if _, exists := index[entry.Name]; !exists {
			order = append(order, entry.Name)
		}
		index[entry.Name] = entry
	}
	result := make([]secretregistry.Entry, 0, len(order))
	seen := map[string]struct{}{}
	for _, name := range order {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, index[name])
	}
	return result
}

func mergeProfiles(base []profile.Profile, override []profile.Profile) []profile.Profile {
	order := make([]string, 0, len(base)+len(override))
	index := map[string]profile.Profile{}
	for _, entry := range base {
		order = append(order, entry.Name)
		index[entry.Name] = entry
	}
	for _, entry := range override {
		if _, exists := index[entry.Name]; !exists {
			order = append(order, entry.Name)
		}
		index[entry.Name] = entry
	}
	result := make([]profile.Profile, 0, len(order))
	seen := map[string]struct{}{}
	for _, name := range order {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, index[name])
	}
	return result
}

func mergePortProbes(base []PortProbeConfig, override []PortProbeConfig) []PortProbeConfig {
	order := make([]string, 0, len(base)+len(override))
	index := map[string]PortProbeConfig{}
	for _, entry := range base {
		order = append(order, entry.Name)
		index[entry.Name] = entry
	}
	for _, entry := range override {
		if _, exists := index[entry.Name]; !exists {
			order = append(order, entry.Name)
		}
		index[entry.Name] = entry
	}
	result := make([]PortProbeConfig, 0, len(order))
	seen := map[string]struct{}{}
	for _, name := range order {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, index[name])
	}
	return result
}

func mergeAuthTools(base []AuthToolConfig, override []AuthToolConfig) []AuthToolConfig {
	order := make([]string, 0, len(base)+len(override))
	index := map[string]AuthToolConfig{}
	for _, entry := range base {
		order = append(order, entry.Name)
		index[entry.Name] = entry
	}
	for _, entry := range override {
		if _, exists := index[entry.Name]; !exists {
			order = append(order, entry.Name)
		}
		index[entry.Name] = entry
	}
	result := make([]AuthToolConfig, 0, len(order))
	seen := map[string]struct{}{}
	for _, name := range order {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, index[name])
	}
	return result
}

func mergeDevOpsTools(base []DevOpsToolConfig, override []DevOpsToolConfig) []DevOpsToolConfig {
	order := make([]string, 0, len(base)+len(override))
	index := map[string]DevOpsToolConfig{}
	for _, entry := range base {
		order = append(order, entry.Name)
		index[entry.Name] = entry
	}
	for _, entry := range override {
		if _, exists := index[entry.Name]; !exists {
			order = append(order, entry.Name)
		}
		index[entry.Name] = entry
	}
	result := make([]DevOpsToolConfig, 0, len(order))
	seen := map[string]struct{}{}
	for _, name := range order {
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, index[name])
	}
	return result
}

func defaultString(value string, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func (cfg Config) FindLLMProfile(name string) (profile.Profile, error) {
	target := name
	if target == "" {
		target = cfg.LLMProfiles.PrimaryProfile
	}
	for _, entry := range cfg.LLMProfiles.Profiles {
		if entry.Name == target {
			return entry, nil
		}
	}
	return profile.Profile{}, fmt.Errorf("llm profile not found: %s", target)
}

func (cfg Config) FindSecretAlias(name string) (secretregistry.Entry, error) {
	for _, entry := range cfg.SecretAliases {
		if entry.Name == name {
			return entry, nil
		}
	}
	return secretregistry.Entry{}, fmt.Errorf("secret alias not found: %s", name)
}
