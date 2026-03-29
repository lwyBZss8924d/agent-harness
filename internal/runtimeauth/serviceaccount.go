package runtimeauth

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/broker"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/keychain"
)

func ResolveServiceAccountToken(cfg config.Config) (string, string, error) {
	if value := strings.TrimSpace(os.Getenv(cfg.SecretService.ServiceAccountEnv)); TokenLooksValid(value) {
		return value, "env", nil
	}

	if cfg.Broker.SocketPath != "" {
		client := broker.Client{SocketPath: cfg.Broker.SocketPath}
		response, err := client.GetToken()
		if err == nil && response.OK && TokenLooksValid(response.Token) {
			return strings.TrimSpace(response.Token), "broker", nil
		}
	}

	if value, err := launchctlGetenv(cfg.SecretService.ServiceAccountEnv); err == nil && TokenLooksValid(value) {
		return strings.TrimSpace(value), "launchd", nil
	}

	return "", "", errors.New("no valid service account token is available from env, broker, or launchd")
}

func ResolveServiceAccountTokenForBroker(cfg config.Config) (string, string, error) {
	if value := strings.TrimSpace(os.Getenv(cfg.SecretService.ServiceAccountEnv)); TokenLooksValid(value) {
		return value, "env", nil
	}

	if value, err := keychainReadServiceAccountToken(); err == nil && TokenLooksValid(value) {
		return strings.TrimSpace(value), "keychain", nil
	}

	return "", "", errors.New("no valid service account token is available from env or keychain")
}

func TokenLooksValid(value string) bool {
	trimmed := strings.TrimSpace(value)
	return strings.HasPrefix(trimmed, "ops_") && len(trimmed) >= 200
}

func launchctlGetenv(name string) (string, error) {
	output, err := exec.Command("launchctl", "getenv", name).CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func keychainReadServiceAccountToken() (string, error) {
	return keychain.ReadGenericPassword(context.Background(), keychain.ServiceAccountTokenService, keychain.CurrentUser())
}
