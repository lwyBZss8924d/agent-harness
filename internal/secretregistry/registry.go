package secretregistry

import "github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/secretpolicy"

type Entry struct {
	Name           string                    `json:"name"`
	Category       secretpolicy.Category     `json:"category"`
	Backend        string                    `json:"backend"`
	Reference      string                    `json:"reference"`
	RevealPolicy   secretpolicy.RevealPolicy `json:"reveal_policy"`
	UsagePolicy    secretpolicy.UsagePolicy  `json:"usage_policy"`
	AllowedActions []secretpolicy.Action     `json:"allowed_actions,omitempty"`
	Metadata       map[string]string         `json:"metadata,omitempty"`
}

func (e Entry) GetBackend() string {
	return e.Backend
}

func (e Entry) GetReference() string {
	return e.Reference
}
