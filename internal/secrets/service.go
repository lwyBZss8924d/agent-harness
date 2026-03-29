package secrets

import "context"

type Reference struct {
	Alias     string `json:"alias"`
	Reference string `json:"reference"`
	Category  string `json:"category"`
}

type ResolveResult struct {
	Value  string `json:"value"`
	Source string `json:"source"`
}

type Service interface {
	ResolveReference(ctx context.Context, reference string) (ResolveResult, error)
}
