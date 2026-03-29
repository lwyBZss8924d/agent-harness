package discovery

import (
	"os"
	"os/user"
	"runtime"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
)

type Snapshot struct {
	CurrentUser   string            `json:"current_user"`
	HomeDir       string            `json:"home_dir"`
	OS            string            `json:"os"`
	Arch          string            `json:"arch"`
	CompatBackend string            `json:"compat_backend"`
	AgentsDir     string            `json:"agents_dir"`
	LegacyRuntime string            `json:"legacy_runtime"`
	BrokerSocket  string            `json:"broker_socket"`
	EnvHints      map[string]string `json:"env_hints"`
}

func SnapshotFromConfig(cfg config.Config) Snapshot {
	currentUser := ""
	homeDir := ""
	if u, err := user.Current(); err == nil {
		currentUser = u.Username
		homeDir = u.HomeDir
	} else {
		currentUser = os.Getenv("USER")
		homeDir, _ = os.UserHomeDir()
	}

	return Snapshot{
		CurrentUser:   currentUser,
		HomeDir:       homeDir,
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		CompatBackend: cfg.CompatBackend,
		AgentsDir:     cfg.AgentsDir,
		LegacyRuntime: cfg.LegacyRuntime,
		BrokerSocket:  cfg.Broker.SocketPath,
		EnvHints: map[string]string{
			"AIH_SOURCE_ROOT":            os.Getenv("AIH_SOURCE_ROOT"),
			"AIH_LEGACY_RUNTIME":         os.Getenv("AIH_LEGACY_RUNTIME"),
			"AIH_COMPAT_BACKEND":         os.Getenv("AIH_COMPAT_BACKEND"),
			"AIH_BROKER_SOCKET":          os.Getenv("AIH_BROKER_SOCKET"),
			"AIH_LLM_PROFILE":            os.Getenv("AIH_LLM_PROFILE"),
			"AIH_LLM_BASE_URL":           os.Getenv("AIH_LLM_BASE_URL"),
			"AIH_LLM_CONVERTER_BASE_URL": os.Getenv("AIH_LLM_CONVERTER_BASE_URL"),
			"AIH_LLM_API_KEY_ALIAS":      os.Getenv("AIH_LLM_API_KEY_ALIAS"),
			"AIH_LLM_BASE_URL_ALIAS":     os.Getenv("AIH_LLM_BASE_URL_ALIAS"),
			"AIH_SECRET_BACKEND":         os.Getenv("AIH_SECRET_BACKEND"),
			"AIH_SERVICE_ACCOUNT_ENV":    os.Getenv("AIH_SERVICE_ACCOUNT_ENV"),
			"AIH_VAULT_NAME":             os.Getenv("AIH_VAULT_NAME"),
			"AIH_VAULT_UUID":             os.Getenv("AIH_VAULT_UUID"),
		},
	}
}
