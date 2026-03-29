package llmapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/profile"
)

type VerifyResult struct {
	Profile         string         `json:"profile"`
	Model           string         `json:"model"`
	BaseURL         string         `json:"base_url"`
	Protocol        string         `json:"protocol"`
	HTTPStatus      int            `json:"http_status"`
	RequestDuration int64          `json:"request_duration_ms"`
	Pass            bool           `json:"pass"`
	ResponsePreview string         `json:"response_preview,omitempty"`
	Usage           map[string]any `json:"usage,omitempty"`
}

type RequestResult struct {
	Profile         string         `json:"profile"`
	Model           string         `json:"model,omitempty"`
	BaseURL         string         `json:"base_url"`
	Protocol        string         `json:"protocol"`
	HTTPStatus      int            `json:"http_status"`
	RequestDuration int64          `json:"request_duration_ms"`
	Response        map[string]any `json:"response"`
}

func VerifyOpaque(profile profile.Profile, apiKey string, marker string) (VerifyResult, error) {
	protocol := profile.Metadata["protocol"]
	if protocol == "" {
		protocol = "openai_chat_completions"
	}

	model := ""
	if len(profile.DefaultModels) > 0 {
		model = profile.DefaultModels[0]
	}
	if model == "" {
		return VerifyResult{}, fmt.Errorf("llm profile %q has no default model configured", profile.Name)
	}

	baseURL := selectBaseURL(profile, model)
	if baseURL == "" {
		return VerifyResult{}, fmt.Errorf("llm profile %q has no usable base URL", profile.Name)
	}

	switch protocol {
	case "openai_chat_completions":
		return verifyOpenAI(profile.Name, baseURL, model, apiKey, marker)
	default:
		return VerifyResult{}, fmt.Errorf("unsupported llm protocol: %s", protocol)
	}
}

func RequestOpaque(profile profile.Profile, apiKey string, body map[string]any) (RequestResult, error) {
	protocol := profile.Metadata["protocol"]
	if protocol == "" {
		protocol = "openai_chat_completions"
	}

	model, _ := body["model"].(string)
	if model == "" && len(profile.DefaultModels) > 0 {
		model = profile.DefaultModels[0]
		body["model"] = model
	}

	baseURL := selectBaseURL(profile, model)
	if baseURL == "" {
		return RequestResult{}, fmt.Errorf("llm profile %q has no usable base URL", profile.Name)
	}

	switch protocol {
	case "openai_chat_completions":
		return requestOpenAI(profile.Name, baseURL, apiKey, body)
	default:
		return RequestResult{}, fmt.Errorf("unsupported llm protocol: %s", protocol)
	}
}

func verifyOpenAI(profileName, baseURL, model, apiKey, marker string) (VerifyResult, error) {
	body := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": "Reply with exactly " + marker},
		},
		"temperature": 0,
		"max_tokens":  16,
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return VerifyResult{}, err
	}

	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(baseURL, "/")+"/chat/completions", bytes.NewReader(encoded))
	if err != nil {
		return VerifyResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return VerifyResult{}, err
	}
	defer resp.Body.Close()

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return VerifyResult{}, err
	}

	previewBytes, _ := json.Marshal(payload)
	preview := string(previewBytes)
	if len(preview) > 800 {
		preview = preview[:800] + "...[truncated]"
	}

	pass := false
	if choices, ok := payload["choices"].([]any); ok && len(choices) > 0 {
		if first, ok := choices[0].(map[string]any); ok {
			if message, ok := first["message"].(map[string]any); ok {
				if content, ok := message["content"].(string); ok && strings.Contains(content, marker) {
					pass = true
				}
			}
		}
	}

	result := VerifyResult{
		Profile:         profileName,
		Model:           model,
		BaseURL:         baseURL,
		Protocol:        "openai_chat_completions",
		HTTPStatus:      resp.StatusCode,
		RequestDuration: time.Since(start).Milliseconds(),
		Pass:            pass && resp.StatusCode == 200,
		ResponsePreview: preview,
	}
	if usage, ok := payload["usage"].(map[string]any); ok {
		result.Usage = usage
	}
	return result, nil
}

func requestOpenAI(profileName, baseURL, apiKey string, body map[string]any) (RequestResult, error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return RequestResult{}, err
	}

	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(baseURL, "/")+"/chat/completions", bytes.NewReader(encoded))
	if err != nil {
		return RequestResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return RequestResult{}, err
	}
	defer resp.Body.Close()

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return RequestResult{}, err
	}

	model, _ := body["model"].(string)
	return RequestResult{
		Profile:         profileName,
		Model:           model,
		BaseURL:         baseURL,
		Protocol:        "openai_chat_completions",
		HTTPStatus:      resp.StatusCode,
		RequestDuration: time.Since(start).Milliseconds(),
		Response:        payload,
	}, nil
}

func selectBaseURL(p profile.Profile, model string) string {
	if strings.HasPrefix(model, "claude-") || strings.HasPrefix(model, "gemini-") {
		if p.ConverterBaseURL != "" {
			return p.ConverterBaseURL
		}
	}
	return p.BaseURL
}
