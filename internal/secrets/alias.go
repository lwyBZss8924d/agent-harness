package secrets

import (
	"context"
	"fmt"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretregistry"
)

func ResolveAlias(ctx context.Context, cfg config.Config, alias string) (secretregistry.Entry, ResolveResult, error) {
	entry, err := cfg.FindSecretAlias(alias)
	if err != nil {
		return secretregistry.Entry{}, ResolveResult{}, err
	}

	if result, resolveErr := resolveAliasNative(ctx, cfg, entry); resolveErr == nil {
		return entry, result, nil
	}

	result, resolveErr := resolveAliasCompatEntry(cfg, entry)
	if resolveErr != nil {
		return entry, ResolveResult{}, fmt.Errorf("resolve alias %q: %w", alias, resolveErr)
	}
	return entry, result, nil
}
