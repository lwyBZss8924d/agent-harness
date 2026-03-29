package profilescaffold

import (
	"fmt"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/profile"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretpolicy"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretregistry"
)

type Request struct {
	Kind           string
	Name           string
	Backend        string
	Protocol       string
	SecretAlias    string
	BaseURLAlias   string
	DefaultModel   string
	PrimaryProfile bool
}

type Result struct {
	SecretAliases  []secretregistry.Entry `json:"secret_aliases,omitempty"`
	LLMProfiles    []profile.Profile      `json:"llm_profiles,omitempty"`
	PrimaryProfile string                 `json:"primary_profile,omitempty"`
}

func SupportedKinds() []string {
	return []string{
		"llm_api",
		"generic_token",
	}
}

func Build(req Request) (Result, error) {
	switch req.Kind {
	case "llm_api":
		return buildLLMAPI(req), nil
	case "generic_token":
		return buildGenericToken(req), nil
	default:
		return Result{}, fmt.Errorf("unsupported scaffold kind: %s", req.Kind)
	}
}

func buildLLMAPI(req Request) Result {
	name := fallback(req.Name, "default")
	backend := fallback(req.Backend, "1password-service-account")
	secretAlias := fallback(req.SecretAlias, name+"-api-key")
	baseURLAlias := fallback(req.BaseURLAlias, name+"-base-url")
	protocol := fallback(req.Protocol, "openai_chat_completions")
	model := fallback(req.DefaultModel, "gpt-5.4-mini")

	return Result{
		SecretAliases: []secretregistry.Entry{
			{
				Name:         secretAlias,
				Category:     secretpolicy.CategoryToken,
				Backend:      backend,
				Reference:    "<fill-me>",
				RevealPolicy: secretpolicy.RevealNever,
				UsagePolicy:  secretpolicy.UsageOpaqueUse,
			},
			{
				Name:         baseURLAlias,
				Category:     secretpolicy.CategoryCredential,
				Backend:      backend,
				Reference:    "<fill-me>",
				RevealPolicy: secretpolicy.RevealAllowed,
				UsagePolicy:  secretpolicy.UsageReveal,
			},
		},
		LLMProfiles: []profile.Profile{
			{
				Name:           name,
				Category:       secretpolicy.CategoryLLMAPI,
				Backend:        backend,
				SecretAlias:    secretAlias,
				RevealPolicy:   secretpolicy.RevealNever,
				UsagePolicy:    secretpolicy.UsageOpaqueUse,
				AllowedActions: []secretpolicy.Action{secretpolicy.ActionLLMVerify, secretpolicy.ActionLLMRequest, secretpolicy.ActionLLMModels},
				BaseURLAlias:   baseURLAlias,
				DefaultModels:  []string{model},
				Metadata: map[string]string{
					"protocol":    protocol,
					"auth_scheme": "bearer",
				},
			},
		},
		PrimaryProfile: primary(name, req.PrimaryProfile),
	}
}

func buildGenericToken(req Request) Result {
	name := fallback(req.Name, "generic-token")
	backend := fallback(req.Backend, "1password-service-account")
	secretAlias := fallback(req.SecretAlias, name)

	return Result{
		SecretAliases: []secretregistry.Entry{
			{
				Name:         secretAlias,
				Category:     secretpolicy.CategoryToken,
				Backend:      backend,
				Reference:    "<fill-me>",
				RevealPolicy: secretpolicy.RevealNever,
				UsagePolicy:  secretpolicy.UsageOpaqueUse,
			},
		},
	}
}

func fallback(value, def string) string {
	if value != "" {
		return value
	}
	return def
}

func primary(name string, enabled bool) string {
	if enabled {
		return name
	}
	return ""
}
