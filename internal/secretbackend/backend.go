package secretbackend

import "context"

type Capability string

const (
	CapabilityResolveReference Capability = "resolve_reference"
	CapabilityHealthCheck      Capability = "health_check"
	CapabilityOpaqueUse        Capability = "opaque_use"
	CapabilityReveal           Capability = "reveal"
)

type Status struct {
	Name         string         `json:"name"`
	Available    bool           `json:"available"`
	Capabilities []Capability   `json:"capabilities"`
	Details      map[string]any `json:"details,omitempty"`
}

type ResolveRequest struct {
	Reference string `json:"reference"`
}

type ResolveResult struct {
	Value  string `json:"value"`
	Source string `json:"source"`
}

type Backend interface {
	Name() string
	Status(ctx context.Context) (Status, error)
	ResolveReference(ctx context.Context, request ResolveRequest) (ResolveResult, error)
}
