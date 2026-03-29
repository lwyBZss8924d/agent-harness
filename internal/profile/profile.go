package profile

import (
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretpolicy"
)

type Profile struct {
	Name             string                    `json:"name"`
	Category         secretpolicy.Category     `json:"category"`
	Backend          string                    `json:"backend"`
	SecretRef        string                    `json:"secret_ref"`
	SecretAlias      string                    `json:"secret_alias,omitempty"`
	RevealPolicy     secretpolicy.RevealPolicy `json:"reveal_policy"`
	UsagePolicy      secretpolicy.UsagePolicy  `json:"usage_policy"`
	AllowedActions   []secretpolicy.Action     `json:"allowed_actions"`
	BaseURL          string                    `json:"base_url,omitempty"`
	ConverterBaseURL string                    `json:"converter_base_url,omitempty"`
	DefaultModels    []string                  `json:"default_models,omitempty"`
	BaseURLAlias     string                    `json:"base_url_alias,omitempty"`
	Metadata         map[string]string         `json:"metadata,omitempty"`
}

func (p Profile) AllowsAction(action secretpolicy.Action) bool {
	for _, candidate := range p.AllowedActions {
		if candidate == action {
			return true
		}
	}
	return false
}

func (p Profile) AllowsReveal() bool {
	return p.RevealPolicy == secretpolicy.RevealAllowed
}
