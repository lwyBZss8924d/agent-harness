package profilereadiness

import (
	"fmt"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
)

type Readiness struct {
	Name           string   `json:"name"`
	Backend        string   `json:"backend"`
	Ready          bool     `json:"ready"`
	Missing        []string `json:"missing,omitempty"`
	Warnings       []string `json:"warnings,omitempty"`
	AllowedActions []string `json:"allowed_actions,omitempty"`
	RevealPolicy   string   `json:"reveal_policy"`
	UsagePolicy    string   `json:"usage_policy"`
}

func AssessAll(cfg config.Config) []Readiness {
	results := make([]Readiness, 0, len(cfg.LLMProfiles.Profiles))
	for _, prof := range cfg.LLMProfiles.Profiles {
		missing := []string{}
		warnings := []string{}

		if prof.SecretRef == "" && prof.SecretAlias == "" {
			missing = append(missing, "secret_ref_or_secret_alias")
		}
		if prof.SecretAlias != "" {
			entry, err := cfg.FindSecretAlias(prof.SecretAlias)
			if err != nil {
				missing = append(missing, "secret_alias:"+prof.SecretAlias)
			} else if entry.Reference == "" || entry.Reference == "<fill-me>" {
				missing = append(missing, "secret_alias_reference:"+prof.SecretAlias)
			}
		}
		if prof.BaseURL == "" && prof.BaseURLAlias == "" {
			missing = append(missing, "base_url_or_base_url_alias")
		}
		if prof.BaseURLAlias != "" {
			entry, err := cfg.FindSecretAlias(prof.BaseURLAlias)
			if err != nil {
				missing = append(missing, "base_url_alias:"+prof.BaseURLAlias)
			} else if entry.Reference == "" || entry.Reference == "<fill-me>" {
				missing = append(missing, "base_url_alias_reference:"+prof.BaseURLAlias)
			}
		}
		if len(prof.DefaultModels) == 0 {
			warnings = append(warnings, "default_models_empty")
		}
		if len(prof.AllowedActions) == 0 {
			warnings = append(warnings, "allowed_actions_empty")
		}

		actionNames := make([]string, 0, len(prof.AllowedActions))
		for _, action := range prof.AllowedActions {
			actionNames = append(actionNames, string(action))
		}

		results = append(results, Readiness{
			Name:           prof.Name,
			Backend:        prof.Backend,
			Ready:          len(missing) == 0,
			Missing:        missing,
			Warnings:       warnings,
			AllowedActions: actionNames,
			RevealPolicy:   string(prof.RevealPolicy),
			UsagePolicy:    string(prof.UsagePolicy),
		})
	}
	return results
}

func Find(cfg config.Config, name string) (Readiness, error) {
	for _, item := range AssessAll(cfg) {
		if item.Name == name {
			return item, nil
		}
	}
	return Readiness{}, fmt.Errorf("profile readiness not found: %s", name)
}
