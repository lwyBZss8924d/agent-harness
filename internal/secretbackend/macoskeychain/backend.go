package macoskeychain

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/broker"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/keychain"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretbackend"
)

type Backend struct {
	Config config.Config
}

func (b Backend) Name() string {
	return "macos-keychain"
}

func (b Backend) Status(ctx context.Context) (secretbackend.Status, error) {
	status := secretbackend.Status{
		Name:      b.Name(),
		Available: runtime.GOOS == "darwin",
		Capabilities: []secretbackend.Capability{
			secretbackend.CapabilityResolveReference,
			secretbackend.CapabilityHealthCheck,
			secretbackend.CapabilityReveal,
			secretbackend.CapabilityOpaqueUse,
		},
		Details: map[string]any{
			"os": runtime.GOOS,
		},
	}
	if runtime.GOOS != "darwin" {
		return status, fmt.Errorf("macos-keychain backend is only available on darwin")
	}
	return status, nil
}

func (b Backend) ResolveReference(ctx context.Context, request secretbackend.ResolveRequest) (secretbackend.ResolveResult, error) {
	if b.Config.Broker.SocketPath != "" {
		client := broker.Client{SocketPath: b.Config.Broker.SocketPath}
		if response, err := client.GetMaterial(request.Reference); err == nil && response.OK {
			return secretbackend.ResolveResult{
				Value:  response.Material,
				Source: response.MaterialSource,
			}, nil
		}
	}

	service, account, err := keychain.ParseReference(request.Reference)
	if err != nil {
		return secretbackend.ResolveResult{}, err
	}

	output, err := keychain.ReadGenericPassword(ctx, service, account)
	if err != nil {
		return secretbackend.ResolveResult{}, err
	}
	return secretbackend.ResolveResult{
		Value:  strings.TrimSpace(output),
		Source: "macos-keychain:security",
	}, nil
}
