package profilescaffold

import "testing"

func TestBuildLLMAPI(t *testing.T) {
	result, err := Build(Request{
		Kind:           "llm_api",
		Name:           "openrouter",
		Backend:        "1password-service-account",
		Protocol:       "openai_chat_completions",
		SecretAlias:    "openrouter-api-key",
		BaseURLAlias:   "openrouter-base-url",
		DefaultModel:   "gpt-5.4-mini",
		PrimaryProfile: true,
	})
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if got := result.PrimaryProfile; got != "openrouter" {
		t.Fatalf("PrimaryProfile = %q, want openrouter", got)
	}
	if len(result.SecretAliases) != 2 {
		t.Fatalf("len(SecretAliases) = %d, want 2", len(result.SecretAliases))
	}
	if len(result.LLMProfiles) != 1 {
		t.Fatalf("len(LLMProfiles) = %d, want 1", len(result.LLMProfiles))
	}
	profile := result.LLMProfiles[0]
	if profile.Name != "openrouter" {
		t.Fatalf("profile.Name = %q, want openrouter", profile.Name)
	}
	if profile.SecretAlias != "openrouter-api-key" {
		t.Fatalf("profile.SecretAlias = %q", profile.SecretAlias)
	}
	if profile.BaseURLAlias != "openrouter-base-url" {
		t.Fatalf("profile.BaseURLAlias = %q", profile.BaseURLAlias)
	}
	if len(profile.DefaultModels) != 1 || profile.DefaultModels[0] != "gpt-5.4-mini" {
		t.Fatalf("profile.DefaultModels = %#v", profile.DefaultModels)
	}
}

func TestBuildUnsupportedKind(t *testing.T) {
	if _, err := Build(Request{Kind: "nope"}); err == nil {
		t.Fatal("expected error for unsupported kind")
	}
}
