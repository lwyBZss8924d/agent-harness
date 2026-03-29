package aih

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRunReleaseJSON(t *testing.T) {
	app := New()
	var out bytes.Buffer
	var errOut bytes.Buffer
	app.Stdout = &out
	app.Stderr = &errOut

	code := app.Run([]string{"release", "--json"})
	if code != 0 {
		t.Fatalf("Run returned %d stderr=%s", code, errOut.String())
	}

	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if payload["release_target"] != "0.0.1" {
		t.Fatalf("release_target = %#v", payload["release_target"])
	}
}

func TestRunFactsPathJSON(t *testing.T) {
	app := New()
	var out bytes.Buffer
	var errOut bytes.Buffer
	app.Stdout = &out
	app.Stderr = &errOut

	code := app.Run([]string{"facts", "path", "--json"})
	if code != 0 {
		t.Fatalf("Run returned %d stderr=%s", code, errOut.String())
	}
	var payload map[string]string
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if payload["facts_json"] == "" || payload["facts_md"] == "" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestSecretReadDisabledInReleaseInstalledMode(t *testing.T) {
	t.Setenv("AIH_UNSAFE_REVEAL", "1")
	restore := detectInstallMode
	detectInstallMode = func() string { return "release-installed" }
	defer func() { detectInstallMode = restore }()

	app := New()
	var errOut bytes.Buffer
	app.Stderr = &errOut

	code := app.Run([]string{"secret", "read", "op://AIAGENTS/OPENROUTER_API_KEY/credential"})
	if code != 2 {
		t.Fatalf("Run returned %d stderr=%s", code, errOut.String())
	}
	if !bytes.Contains(errOut.Bytes(), []byte("disabled in release-installed mode")) {
		t.Fatalf("stderr = %s", errOut.String())
	}
}

func TestSecretGetRevealNeverDeniedEvenWithUnsafeFlag(t *testing.T) {
	configPath := writeTempConfig(t, `{
  "secret_aliases": [
    {
      "name": "openrouter-api-key",
      "backend": "1password-service-account",
      "reference": "op://AIAGENTS/OPENROUTER_API_KEY/credential",
      "reveal_policy": "never",
      "usage_policy": "opaque_use"
    }
  ]
}`)
	t.Setenv("AIH_CONFIG_FILE", configPath)
	t.Setenv("AIH_UNSAFE_REVEAL", "1")
	restore := detectInstallMode
	detectInstallMode = func() string { return "dev-go-run" }
	defer func() { detectInstallMode = restore }()

	app := New()
	var errOut bytes.Buffer
	app.Stderr = &errOut

	code := app.Run([]string{"secret", "get", "openrouter-api-key"})
	if code != 2 {
		t.Fatalf("Run returned %d stderr=%s", code, errOut.String())
	}
	if !bytes.Contains(errOut.Bytes(), []byte("is not revealable")) {
		t.Fatalf("stderr = %s", errOut.String())
	}
}

func TestSecretGetAdminOnlyDeniedInReleaseInstalledMode(t *testing.T) {
	configPath := writeTempConfig(t, `{
  "secret_aliases": [
    {
      "name": "account-password",
      "backend": "1password-service-account",
      "reference": "op://AIAGENTS/mac-server-01/wijfsih7ufa5k6eczzfv2bnqai",
      "reveal_policy": "admin_only",
      "usage_policy": "opaque_use"
    }
  ]
}`)
	t.Setenv("AIH_CONFIG_FILE", configPath)
	t.Setenv("AIH_UNSAFE_REVEAL", "1")
	restore := detectInstallMode
	detectInstallMode = func() string { return "release-installed" }
	defer func() { detectInstallMode = restore }()

	app := New()
	var errOut bytes.Buffer
	app.Stderr = &errOut

	code := app.Run([]string{"secret", "get", "account-password"})
	if code != 2 {
		t.Fatalf("Run returned %d stderr=%s", code, errOut.String())
	}
	if !bytes.Contains(errOut.Bytes(), []byte("requires non-release admin mode")) {
		t.Fatalf("stderr = %s", errOut.String())
	}
}

func TestSecretEnvDisabledInReleaseInstalledMode(t *testing.T) {
	t.Setenv("AIH_UNSAFE_INJECT", "1")
	restore := detectInstallMode
	detectInstallMode = func() string { return "release-installed" }
	defer func() { detectInstallMode = restore }()

	app := New()
	var errOut bytes.Buffer
	app.Stderr = &errOut

	code := app.Run([]string{"secret", "env"})
	if code != 2 {
		t.Fatalf("Run returned %d stderr=%s", code, errOut.String())
	}
	if !bytes.Contains(errOut.Bytes(), []byte("disabled in release-installed mode")) {
		t.Fatalf("stderr = %s", errOut.String())
	}
}

func TestSecretExecDisabledInReleaseInstalledMode(t *testing.T) {
	t.Setenv("AIH_UNSAFE_INJECT", "1")
	restore := detectInstallMode
	detectInstallMode = func() string { return "release-installed" }
	defer func() { detectInstallMode = restore }()

	app := New()
	var errOut bytes.Buffer
	app.Stderr = &errOut

	code := app.Run([]string{"secret", "exec", "--", "env"})
	if code != 2 {
		t.Fatalf("Run returned %d stderr=%s", code, errOut.String())
	}
	if !bytes.Contains(errOut.Bytes(), []byte("disabled in release-installed mode")) {
		t.Fatalf("stderr = %s", errOut.String())
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("os.WriteFile: %v", err)
	}
	return path
}
