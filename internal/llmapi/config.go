package llmapi

type VerifierConfig struct {
	ProfileName      string   `json:"profile_name"`
	BaseURL          string   `json:"base_url"`
	ConverterBaseURL string   `json:"converter_base_url"`
	DefaultModels    []string `json:"default_models"`
	APIKeyAlias      string   `json:"api_key_alias"`
	BaseURLAlias     string   `json:"base_url_alias"`
}
