package runtimeprobe

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/broker"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
)

type Probe struct {
	SecretService SecretServiceProbe `json:"secret_service"`
	Broker        BrokerProbe        `json:"broker"`
	Environment   EnvironmentProbe   `json:"environment"`
	Launchd       LaunchdProbe       `json:"launchd"`
	Keychain      KeychainProbe      `json:"keychain"`
	LegacyRuntime LegacyRuntimeProbe `json:"legacy_runtime"`
}

type SecretServiceProbe struct {
	Kind              string `json:"kind"`
	Provider          string `json:"provider"`
	ServiceAccountEnv string `json:"service_account_env"`
	VaultName         string `json:"vault_name"`
	VaultUUID         string `json:"vault_uuid"`
}

type BrokerProbe struct {
	SocketPath   string                 `json:"socket_path"`
	SocketExists bool                   `json:"socket_exists"`
	Reachable    bool                   `json:"reachable"`
	Error        string                 `json:"error,omitempty"`
	Status       *broker.StatusResponse `json:"status,omitempty"`
}

type EnvironmentProbe struct {
	ServiceAccountEnvPresent bool `json:"service_account_env_present"`
	ServiceAccountEnvValid   bool `json:"service_account_env_valid"`
}

type LaunchdProbe struct {
	Available                bool   `json:"available"`
	ServiceAccountEnvPresent bool   `json:"service_account_env_present"`
	ServiceAccountEnvValid   bool   `json:"service_account_env_valid"`
	Error                    string `json:"error,omitempty"`
}

type KeychainProbe struct {
	Available    bool   `json:"available"`
	ItemPresent  bool   `json:"item_present"`
	ItemReadable bool   `json:"item_readable"`
	Unlocked     *bool  `json:"unlocked,omitempty"`
	Error        string `json:"error,omitempty"`
}

type LegacyRuntimeProbe struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
}

func Collect(cfg config.Config) Probe {
	envToken := os.Getenv(cfg.SecretService.ServiceAccountEnv)
	launchdValue, launchdErr := launchctlGetenv(cfg.SecretService.ServiceAccountEnv)
	keychainPresent, keychainReadable, keychainUnlocked, keychainErr := keychainProbe()

	result := Probe{
		SecretService: SecretServiceProbe{
			Kind:              cfg.SecretService.Kind,
			Provider:          cfg.SecretService.Provider,
			ServiceAccountEnv: cfg.SecretService.ServiceAccountEnv,
			VaultName:         cfg.SecretService.VaultName,
			VaultUUID:         cfg.SecretService.VaultUUID,
		},
		Broker: BrokerProbe{
			SocketPath:   cfg.Broker.SocketPath,
			SocketExists: fileExists(cfg.Broker.SocketPath),
		},
		Environment: EnvironmentProbe{
			ServiceAccountEnvPresent: envToken != "",
			ServiceAccountEnvValid:   tokenLooksValid(envToken),
		},
		Launchd: LaunchdProbe{
			Available:                runtime.GOOS == "darwin",
			ServiceAccountEnvPresent: launchdValue != "",
			ServiceAccountEnvValid:   tokenLooksValid(launchdValue),
		},
		Keychain: KeychainProbe{
			Available:    runtime.GOOS == "darwin",
			ItemPresent:  keychainPresent,
			ItemReadable: keychainReadable,
			Unlocked:     keychainUnlocked,
		},
		LegacyRuntime: LegacyRuntimeProbe{
			Path:   cfg.LegacyRuntime,
			Exists: fileExists(cfg.LegacyRuntime),
		},
	}

	if launchdErr != nil {
		result.Launchd.Error = launchdErr.Error()
	}
	if keychainErr != nil {
		result.Keychain.Error = keychainErr.Error()
	}

	if result.Broker.SocketExists {
		client := broker.Client{SocketPath: cfg.Broker.SocketPath}
		status, err := client.Status()
		if err != nil {
			result.Broker.Error = err.Error()
		} else {
			result.Broker.Reachable = status.OK
			result.Broker.Status = &status
		}
	}

	return result
}

func tokenLooksValid(value string) bool {
	return strings.HasPrefix(strings.TrimSpace(value), "ops_") && len(strings.TrimSpace(value)) >= 200
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func launchctlGetenv(name string) (string, error) {
	if runtime.GOOS != "darwin" {
		return "", nil
	}
	output, err := exec.Command("launchctl", "getenv", name).CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func keychainProbe() (present bool, readable bool, unlocked *bool, err error) {
	if runtime.GOOS != "darwin" {
		return false, false, nil, nil
	}

	serviceName := "aiagents/op-service-account-token"
	findCmd := exec.Command("security", "find-generic-password", "-a", os.Getenv("USER"), "-s", serviceName)
	if output, findErr := findCmd.CombinedOutput(); findErr == nil {
		present = true
		if strings.Contains(string(output), "\"svce\"<blob>=\""+serviceName+"\"") || len(output) > 0 {
			present = true
		}
	} else if len(output) > 0 {
		if strings.Contains(string(output), "\"svce\"<blob>=\""+serviceName+"\"") {
			present = true
		}
	}

	readCmd := exec.Command("security", "find-generic-password", "-w", "-a", os.Getenv("USER"), "-s", serviceName)
	if output, readErr := readCmd.CombinedOutput(); readErr == nil {
		readable = tokenLooksValid(string(output))
	} else if strings.Contains(strings.ToLower(string(output)), "user interaction is not allowed") {
		value := false
		unlocked = &value
	} else if len(output) > 0 && strings.Contains(strings.ToLower(string(output)), "could not be found") {
		// no-op
	} else if readErr != nil {
		err = readErr
	}

	infoCmd := exec.Command("security", "show-keychain-info", os.ExpandEnv("$HOME/Library/Keychains/login.keychain-db"))
	if output, infoErr := infoCmd.CombinedOutput(); infoErr == nil {
		value := true
		unlocked = &value
		if len(output) == 0 {
			// keep as true
		}
	} else if strings.Contains(strings.ToLower(string(output)), "user interaction is not allowed") {
		value := false
		unlocked = &value
	}

	return present, readable, unlocked, err
}
