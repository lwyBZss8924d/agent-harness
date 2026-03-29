package secretbackendregistry

import (
	"context"
	"fmt"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretbackend"
	macoskeychain "github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretbackend/macoskeychain"
	onepasswordbackend "github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretbackend/onepassword"
)

func ByName(cfg config.Config, name string) (secretbackend.Backend, error) {
	switch name {
	case "1password-service-account":
		return onepasswordbackend.Backend{Config: cfg}, nil
	case "macos-keychain":
		return macoskeychain.Backend{Config: cfg}, nil
	default:
		return nil, fmt.Errorf("unknown secret backend: %s", name)
	}
}

func StatusAll(ctx context.Context, cfg config.Config) []secretbackend.Status {
	names := []string{"1password-service-account", "macos-keychain"}
	results := make([]secretbackend.Status, 0, len(names))
	for _, name := range names {
		backend, err := ByName(cfg, name)
		if err != nil {
			results = append(results, secretbackend.Status{Name: name, Available: false, Details: map[string]any{"error": err.Error()}})
			continue
		}
		status, statusErr := backend.Status(ctx)
		if statusErr != nil {
			if status.Details == nil {
				status.Details = map[string]any{}
			}
			status.Details["error"] = statusErr.Error()
		}
		results = append(results, status)
	}
	return results
}
