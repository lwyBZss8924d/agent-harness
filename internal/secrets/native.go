package secrets

import (
	"context"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretbackend"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretbackendregistry"
)

func ResolveReferenceNative(ctx context.Context, cfg config.Config, reference string) (ResolveResult, error) {
	return ResolveReferenceForBackend(ctx, cfg, cfg.SecretService.Kind, reference)
}

func ResolveReferenceForBackend(ctx context.Context, cfg config.Config, backendName string, reference string) (ResolveResult, error) {
	backend, err := secretbackendregistry.ByName(cfg, backendName)
	if err != nil {
		return ResolveResult{}, err
	}
	result, err := backend.ResolveReference(ctx, secretbackend.ResolveRequest{Reference: reference})
	if err != nil {
		return ResolveResult{}, err
	}
	return ResolveResult{
		Value:  result.Value,
		Source: result.Source,
	}, nil
}
