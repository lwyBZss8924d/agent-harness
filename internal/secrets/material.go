package secrets

import (
	"context"
	"fmt"
	"strings"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/compat/legacy"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/profile"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretregistry"
)

func ResolveProfileMaterial(ctx context.Context, cfg config.Config, p profile.Profile) (ResolveResult, error) {
	if p.SecretRef != "" {
		return ResolveReferenceForBackend(ctx, cfg, p.Backend, p.SecretRef)
	}

	if p.SecretAlias != "" {
		if entry, err := cfg.FindSecretAlias(p.SecretAlias); err == nil {
			if result, resolveErr := resolveAliasNative(ctx, cfg, entry); resolveErr == nil {
				return result, nil
			}
			return resolveAliasCompatEntry(cfg, entry)
		}
		return resolveAliasCompat(cfg, p.SecretAlias)
	}

	return ResolveResult{}, fmt.Errorf("profile %q has neither secret_ref nor secret_alias", p.Name)
}

func ResolveProfileBaseURL(ctx context.Context, cfg config.Config, p profile.Profile) (string, error) {
	if strings.TrimSpace(p.BaseURL) != "" {
		return strings.TrimSpace(p.BaseURL), nil
	}

	if p.BaseURLAlias != "" {
		if entry, err := cfg.FindSecretAlias(p.BaseURLAlias); err == nil {
			if result, resolveErr := resolveAliasNative(ctx, cfg, entry); resolveErr == nil {
				return strings.TrimSpace(result.Value), nil
			}
			result, resolveErr := resolveAliasCompatEntry(cfg, entry)
			if resolveErr != nil {
				return "", resolveErr
			}
			return strings.TrimSpace(result.Value), nil
		}
		result, err := resolveAliasCompat(cfg, p.BaseURLAlias)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(result.Value), nil
	}

	return "", fmt.Errorf("profile %q has no base_url or base_url_alias", p.Name)
}

func resolveAliasNative(ctx context.Context, cfg config.Config, entry secretregistry.Entry) (ResolveResult, error) {
	if entry.Reference == "" {
		return ResolveResult{}, fmt.Errorf("secret alias %q has no reference configured", entry.Name)
	}
	return ResolveReferenceForBackend(ctx, cfg, entry.Backend, entry.Reference)
}

func resolveAliasCompat(cfg config.Config, alias string) (ResolveResult, error) {
	if cfg.CompatBackend != "legacy-python" || cfg.LegacyRuntime == "" {
		return ResolveResult{}, fmt.Errorf("compat alias resolution is unavailable for alias %q", alias)
	}

	output, exitCode, err := legacy.New(cfg.LegacyRuntime).Output([]string{"secret", "get", alias})
	if err != nil {
		return ResolveResult{}, err
	}
	if exitCode != 0 {
		return ResolveResult{}, fmt.Errorf("compat alias resolution failed for %q with exit code %d", alias, exitCode)
	}

	return ResolveResult{
		Value:  strings.TrimSpace(output),
		Source: "compat:legacy-python:alias",
	}, nil
}

func resolveAliasCompatEntry(cfg config.Config, entry secretregistry.Entry) (ResolveResult, error) {
	if cfg.CompatBackend != "legacy-python" || cfg.LegacyRuntime == "" {
		return ResolveResult{}, fmt.Errorf("compat alias resolution is unavailable for alias %q", entry.Name)
	}
	if entry.Reference == "" {
		return ResolveResult{}, fmt.Errorf("compat alias resolution is unavailable for alias %q without reference", entry.Name)
	}

	output, exitCode, err := legacy.New(cfg.LegacyRuntime).Output([]string{"secret", "read", entry.Reference})
	if err != nil {
		return ResolveResult{}, err
	}
	if exitCode != 0 {
		return ResolveResult{}, fmt.Errorf("compat alias resolution by reference failed for %q with exit code %d", entry.Name, exitCode)
	}
	return ResolveResult{
		Value:  strings.TrimSpace(output),
		Source: "compat:legacy-python:reference",
	}, nil
}
