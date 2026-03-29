package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
)

type Status struct {
	ChromeAppExists         bool           `json:"chrome_app_exists"`
	ChromeBinary            string         `json:"chrome_binary"`
	ChromeBinaryExists      bool           `json:"chrome_binary_exists"`
	CDPPort                 int            `json:"cdp_port"`
	CDPListening            bool           `json:"cdp_listening"`
	CDPVersionAvailable     bool           `json:"cdp_version_available"`
	CDPVersion              map[string]any `json:"cdp_version,omitempty"`
	PlaywrightHelper        string         `json:"playwright_helper"`
	PlaywrightHelperExists  bool           `json:"playwright_helper_exists"`
	PlaywrightCoreDir       string         `json:"playwright_core_dir"`
	PlaywrightCoreInstalled bool           `json:"playwright_core_installed"`
	AutomationProfileDir    string         `json:"automation_profile_dir"`
	AutomationProfileExists bool           `json:"automation_profile_exists"`
}

func Collect(cfg config.Config, port int) Status {
	if port == 0 {
		port = cfg.Browser.CDPPortDefault
	}

	versionOK, versionPayload := cdpVersion(port)
	return Status{
		ChromeAppExists:         fileExists(cfg.Browser.ChromeAppPath),
		ChromeBinary:            cfg.Browser.ChromeBinaryPath,
		ChromeBinaryExists:      fileExists(cfg.Browser.ChromeBinaryPath),
		CDPPort:                 port,
		CDPListening:            portListening(port),
		CDPVersionAvailable:     versionOK,
		CDPVersion:              versionPayload,
		PlaywrightHelper:        cfg.Browser.PlaywrightHelperPath,
		PlaywrightHelperExists:  fileExists(cfg.Browser.PlaywrightHelperPath),
		PlaywrightCoreDir:       cfg.Browser.PlaywrightCoreDir,
		PlaywrightCoreInstalled: fileExists(cfg.Browser.PlaywrightCoreDir),
		AutomationProfileDir:    cfg.Browser.AutomationProfileDir,
		AutomationProfileExists: fileExists(cfg.Browser.AutomationProfileDir),
	}
}

func LaunchCDP(cfg config.Config, port int, profileDir string) (Status, error) {
	if port == 0 {
		port = cfg.Browser.CDPPortDefault
	}
	if profileDir == "" {
		profileDir = cfg.Browser.AutomationProfileDir
	}
	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		return Status{}, err
	}
	if !fileExists(cfg.Browser.ChromeBinaryPath) {
		return Status{}, fmt.Errorf("chrome binary is not installed at %s", cfg.Browser.ChromeBinaryPath)
	}

	switch runtime.GOOS {
	case "darwin":
		appTarget := strings.TrimSuffix(filepath.Base(cfg.Browser.ChromeAppPath), ".app")
		if appTarget == "" || appTarget == "." {
			appTarget = "Google Chrome"
		}
		cmd := exec.Command(
			"open",
			"-na",
			appTarget,
			"--args",
			fmt.Sprintf("--remote-debugging-port=%d", port),
			fmt.Sprintf("--user-data-dir=%s", profileDir),
			"--no-first-run",
			"--no-default-browser-check",
			"--disable-default-apps",
			"--disable-background-networking",
			"about:blank",
		)
		if output, err := cmd.CombinedOutput(); err != nil {
			return Status{}, fmt.Errorf("launch chrome with cdp: %w (%s)", err, strings.TrimSpace(string(output)))
		}
	default:
		return Status{}, fmt.Errorf("browser launch-cdp is not yet implemented for %s", runtime.GOOS)
	}

	for i := 0; i < 20; i++ {
		if portListening(port) {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	overrideCfg := cfg
	overrideCfg.Browser.AutomationProfileDir = profileDir
	status := Collect(overrideCfg, port)
	if !status.CDPListening {
		return status, fmt.Errorf("chrome did not expose cdp on port %d", port)
	}
	return status, nil
}

func VerifyPlaywright(cfg config.Config, port int, timeoutSeconds int) (map[string]any, int, error) {
	if port == 0 {
		port = cfg.Browser.CDPPortDefault
	}
	if timeoutSeconds <= 0 {
		timeoutSeconds = 20
	}
	if !fileExists(cfg.Browser.PlaywrightHelperPath) {
		return nil, 2, fmt.Errorf("missing helper: %s", cfg.Browser.PlaywrightHelperPath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "node", cfg.Browser.PlaywrightHelperPath)
	cmd.Env = append(os.Environ(), fmt.Sprintf("AIH_CDP_PORT=%d", port))
	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, 2, fmt.Errorf("playwright verification timed out after %ds", timeoutSeconds)
	}

	var payload map[string]any
	if len(output) > 0 {
		_ = json.Unmarshal(output, &payload)
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return payload, exitErr.ExitCode(), nil
		}
		return payload, 2, err
	}
	return payload, 0, nil
}

func cdpVersion(port int) (bool, map[string]any) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/json/version", port), nil)
	if err != nil {
		return false, nil
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}
	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return false, nil
	}
	return true, payload
}

func portListening(port int) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", port)), 500*time.Millisecond)
	if err == nil {
		_ = conn.Close()
		return true
	}
	return false
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}
