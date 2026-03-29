package onepasswordbackend

import (
	"context"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/broker"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	opsdk "github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/onepassword"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/runtimeauth"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretbackend"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/version"
)

type Backend struct {
	Config config.Config
}

func (b Backend) Name() string {
	return "1password-service-account"
}

func (b Backend) Status(ctx context.Context) (secretbackend.Status, error) {
	status := secretbackend.Status{
		Name:      b.Name(),
		Available: false,
		Capabilities: []secretbackend.Capability{
			secretbackend.CapabilityResolveReference,
			secretbackend.CapabilityHealthCheck,
			secretbackend.CapabilityReveal,
			secretbackend.CapabilityOpaqueUse,
		},
		Details: map[string]any{
			"provider":            b.Config.SecretService.Provider,
			"service_account_env": b.Config.SecretService.ServiceAccountEnv,
			"vault_name":          b.Config.SecretService.VaultName,
			"vault_uuid":          b.Config.SecretService.VaultUUID,
		},
	}

	_, source, err := runtimeauth.ResolveServiceAccountToken(b.Config)
	if err != nil {
		status.Details["token_source"] = ""
		status.Details["error"] = err.Error()
		return status, err
	}
	status.Available = true
	status.Details["token_source"] = source
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

	token, source, err := runtimeauth.ResolveServiceAccountToken(b.Config)
	if err != nil {
		return secretbackend.ResolveResult{}, err
	}

	client := opsdk.ServiceAccountClient{
		IntegrationName:    "aih-go",
		IntegrationVersion: version.Version,
		TokenEnv:           b.Config.SecretService.ServiceAccountEnv,
		Token:              token,
	}
	value, err := client.ResolveReference(ctx, request.Reference)
	if err != nil {
		return secretbackend.ResolveResult{}, err
	}
	return secretbackend.ResolveResult{
		Value:  value,
		Source: "onepassword:service_account:" + source,
	}, nil
}
