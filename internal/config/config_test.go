package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromExplicitAndRepoConfig(t *testing.T) {
	home := t.TempDir()
	repo := filepath.Join(home, "repo", "subdir")
	if err := os.MkdirAll(filepath.Join(home, ".config", "aih"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(home, "repo", ".aih"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}

	globalPath := filepath.Join(home, ".config", "aih", "config.json")
	repoPath := filepath.Join(home, "repo", ".aih", "config.json")
	explicitPath := filepath.Join(home, "explicit.json")

	if err := os.WriteFile(globalPath, []byte(`{
	  "broker": {"socket_path": "/global.sock"},
	  "llm_profiles": {"configured": true, "profiles": [{"name":"global","backend":"1password-service-account"}]}
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(repoPath, []byte(`{
	  "broker": {"socket_path": "/repo.sock"},
	  "secret_aliases": [{"name":"repo-secret","backend":"1password-service-account","reference":"op://x","reveal_policy":"never","usage_policy":"opaque_use"}]
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(explicitPath, []byte(`{
	  "compat_backend": "native-go"
	}`), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("HOME", home)
	t.Setenv("AIH_CONFIG_FILE", explicitPath)
	t.Setenv("AIH_CALLER_CWD", repo)

	cfg := Load()
	if cfg.Broker.SocketPath != "/repo.sock" {
		t.Fatalf("Broker.SocketPath = %q", cfg.Broker.SocketPath)
	}
	if cfg.CompatBackend != "native-go" {
		t.Fatalf("CompatBackend = %q", cfg.CompatBackend)
	}
	if _, err := cfg.FindSecretAlias("repo-secret"); err != nil {
		t.Fatalf("FindSecretAlias(repo-secret): %v", err)
	}
	if _, err := cfg.FindLLMProfile("global"); err != nil {
		t.Fatalf("FindLLMProfile(global): %v", err)
	}
}
