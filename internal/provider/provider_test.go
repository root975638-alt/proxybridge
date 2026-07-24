package provider

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Description method to satisfy interface for provider tests
func (p *OllamaProvider) Description() string {
	return p.description
}

func (p *AWSBedrockProvider) Description() string {
	return p.description
}

func (p *OpenAIProvider) Description() string {
	return p.description
}

func (p *AnthropicProvider) Description() string {
	return p.description
}

func (p *GoogleGeminiProvider) Description() string {
	return p.description
}

func (p *MistralProvider) Description() string {
	return p.description
}

func (p *GroqProvider) Description() string {
	return p.description
}

func (p *DeepSeekProvider) Description() string {
	return p.description
}

func (p *TogetherAIProvider) Description() string {
	return p.description
}

func (p *FireworksProvider) Description() string {
	return p.description
}

func (p *CerebrasProvider) Description() string {
	return p.description
}

func (p *XAIProvider) Description() string {
	return p.description
}

func (p *LMStudioProvider) Description() string {
	return p.description
}

func (p *LocalOpenAIProvider) Description() string {
	return p.description
}

func (p *OpenRouterProvider) Description() string {
	return p.description
}

func (p *AzureOpenAIProvider) Description() string {
	return p.description
}

// TestOllamaProvider tests Ollama provider
func TestOllamaProvider(t *testing.T) {
	p := NewOllamaProvider()

	if p.Name() != "Ollama" {
		t.Errorf("expected name 'Ollama', got '%s'", p.Name())
	}

	if p.ID() != "ollama" {
		t.Errorf("expected ID 'ollama', got '%s'", p.ID())
	}

	expected := "default"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}

	config := p.GetConfig()
	if config["provider"] != "ollama" {
		t.Errorf("expected provider 'ollama', got '%s'", config["provider"])
	}

	if !p.IsInstalled() {
		t.Error("Ollama provider should be considered installed")
	}

	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}

	if err := p.Uninstall(); err != nil {
		t.Errorf("Uninstall returned error: %v", err)
	}
}

// TestAWSBedrockProvider tests AWS Bedrock provider
func TestAWSBedrockProvider(t *testing.T) {
	p := NewAWSBedrockProvider()

	if p.Name() != "AWS Bedrock" {
		t.Errorf("expected name 'AWS Bedrock', got '%s'", p.Name())
	}

	if p.ID() != "aws" {
		t.Errorf("expected ID 'aws', got '%s'", p.ID())
	}

	expected := "anthropic.claude-3-5-sonnet-20241022-v2:0"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestOpenAIProvider tests OpenAI provider
func TestOpenAIProvider(t *testing.T) {
	p := NewOpenAIProvider()

	if p.Name() != "OpenAI" {
		t.Errorf("expected name 'OpenAI', got '%s'", p.Name())
	}

	if p.ID() != "openai" {
		t.Errorf("expected ID 'openai', got '%s'", p.ID())
	}

	expected := "gpt-4o"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestAnthropicProvider tests Anthropic provider
func TestAnthropicProvider(t *testing.T) {
	p := NewAnthropicProvider()

	if p.Name() != "Anthropic" {
		t.Errorf("expected name 'Anthropic', got '%s'", p.Name())
	}

	if p.ID() != "anthropic" {
		t.Errorf("expected ID 'anthropic', got '%s'", p.ID())
	}

	expected := "claude-3-5-sonnet-20241022"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestGoogleGeminiProvider tests Google Gemini provider
func TestGoogleGeminiProvider(t *testing.T) {
	p := NewGoogleGeminiProvider()

	if p.Name() != "Google Gemini" {
		t.Errorf("expected name 'Google Gemini', got '%s'", p.Name())
	}

	if p.ID() != "google" {
		t.Errorf("expected ID 'google', got '%s'", p.ID())
	}

	expected := "gemini-1.5-flash"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestMistralProvider tests Mistral provider
func TestMistralProvider(t *testing.T) {
	p := NewMistralProvider()

	if p.Name() != "Mistral AI" {
		t.Errorf("expected name 'Mistral AI', got '%s'", p.Name())
	}

	if p.ID() != "mistral" {
		t.Errorf("expected ID 'mistral', got '%s'", p.ID())
	}

	expected := "mistral-large"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestGroqProvider tests Groq provider
func TestGroqProvider(t *testing.T) {
	p := NewGroqProvider()

	if p.Name() != "Groq" {
		t.Errorf("expected name 'Groq', got '%s'", p.Name())
	}

	if p.ID() != "groq" {
		t.Errorf("expected ID 'groq', got '%s'", p.ID())
	}

	expected := "llama3-8b-8192"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestDeepSeekProvider tests DeepSeek provider
func TestDeepSeekProvider(t *testing.T) {
	p := NewDeepSeekProvider()

	if p.Name() != "DeepSeek" {
		t.Errorf("expected name 'DeepSeek', got '%s'", p.Name())
	}

	if p.ID() != "deepseek" {
		t.Errorf("expected ID 'deepseek', got '%s'", p.ID())
	}

	expected := "deepseek-chat"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestTogetherAIProvider tests Together AI provider
func TestTogetherAIProvider(t *testing.T) {
	p := NewTogetherAIProvider()

	if p.Name() != "Together AI" {
		t.Errorf("expected name 'Together AI', got '%s'", p.Name())
	}

	if p.ID() != "together" {
		t.Errorf("expected ID 'together', got '%s'", p.ID())
	}

	expected := "togethercomputer/llama-3-70b-chat"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestFireworksProvider tests Fireworks AI provider
func TestFireworksProvider(t *testing.T) {
	p := NewFireworksProvider()

	if p.Name() != "Fireworks AI" {
		t.Errorf("expected name 'Fireworks AI', got '%s'", p.Name())
	}

	if p.ID() != "fireworks" {
		t.Errorf("expected ID 'fireworks', got '%s'", p.ID())
	}

	expected := "accounts/fireworks/models/firefunction-v2"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestCerebrasProvider tests Cerebras provider
func TestCerebrasProvider(t *testing.T) {
	p := NewCerebrasProvider()

	if p.Name() != "Cerebras" {
		t.Errorf("expected name 'Cerebras', got '%s'", p.Name())
	}

	if p.ID() != "cerebras" {
		t.Errorf("expected ID 'cerebras', got '%s'", p.ID())
	}

	expected := "cerebras-llama3.1-8b"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestXAIProvider tests xAI (Grok) provider
func TestXAIProvider(t *testing.T) {
	p := NewXAIProvider()

	if p.Name() != "xAI (Grok)" {
		t.Errorf("expected name 'xAI (Grok)', got '%s'", p.Name())
	}

	if p.ID() != "xai" {
		t.Errorf("expected ID 'xai', got '%s'", p.ID())
	}

	expected := "grok-beta"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestLMStudioProvider tests LM Studio provider
func TestLMStudioProvider(t *testing.T) {
	p := NewLMStudioProvider()

	if p.Name() != "LM Studio" {
		t.Errorf("expected name 'LM Studio', got '%s'", p.Name())
	}

	if p.ID() != "lmstudio" {
		t.Errorf("expected ID 'lmstudio', got '%s'", p.ID())
	}

	expected := "default"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestLocalOpenAIProvider tests local OpenAI-compatible provider
func TestLocalOpenAIProvider(t *testing.T) {
	p := NewLocalOpenAIProvider()

	if p.Name() != "Local / Self-Hosted" {
		t.Errorf("expected name 'Local / Self-Hosted', got '%s'", p.Name())
	}

	if p.ID() != "local" {
		t.Errorf("expected ID 'local', got '%s'", p.ID())
	}

	expected := "default"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestOpenRouterProvider tests OpenRouter provider
func TestOpenRouterProvider(t *testing.T) {
	p := NewOpenRouterProvider()

	if p.Name() != "OpenRouter" {
		t.Errorf("expected name 'OpenRouter', got '%s'", p.Name())
	}

	if p.ID() != "openrouter" {
		t.Errorf("expected ID 'openrouter', got '%s'", p.ID())
	}

	expected := "openai/gpt-4o"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestAzureOpenAIProvider tests Azure OpenAI provider
func TestAzureOpenAIProvider(t *testing.T) {
	p := NewAzureOpenAIProvider()

	if p.Name() != "Azure OpenAI" {
		t.Errorf("expected name 'Azure OpenAI', got '%s'", p.Name())
	}

	if p.ID() != "azure" {
		t.Errorf("expected ID 'azure', got '%s'", p.ID())
	}

	expected := "gpt-4o"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

// TestGetAvailableProviders tests the GetAvailableProviders function
func TestGetAvailableProviders(t *testing.T) {
	providers := GetAvailableProviders()

	if len(providers) == 0 {
		t.Error("expected at least one provider")
	}

	// Check for expected providers
	expected := []string{"openai", "anthropic", "aws", "google", "ollama"}
	for _, id := range expected {
		found := false
		for _, p := range providers {
			if p.ID() == id {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected provider '%s' to be registered", id)
		}
	}
}

// TestProviderRegistry tests provider registry functionality
func TestProviderRegistry(t *testing.T) {
	registry := NewProviderRegistry()
	if registry == nil {
		t.Error("expected non-nil registry")
	}

	openai := registry.GetProvider("openai")
	if openai == nil {
		t.Error("expected openai provider to be registered")
	} else {
		if openai.ID() != "openai" {
			t.Errorf("expected openai provider ID 'openai', got '%s'", openai.ID())
		}
	}

	all := registry.GetAllProviders()
	if len(all) == 0 {
		t.Error("expected at least one provider to be registered")
	}

	names := registry.GetProviderNames()
	if len(names) == 0 {
		t.Error("expected at least one provider name")
	}
}

// TestProviderInstall tests provider installation
func TestProviderInstall(t *testing.T) {
	registry := NewProviderRegistry()

	err := registry.InstallProvider("ollama")
	if err != nil {
		t.Logf("Install returned error (may be expected): %v", err)
	}

	err = registry.UninstallProvider("ollama")
	if err != nil {
		t.Logf("Uninstall returned error (may be expected): %v", err)
	}
}

// TestProviderValidation tests provider validation
func TestProviderValidation(t *testing.T) {
	providers := GetAvailableProviders()

	for _, p := range providers {
		// Skip checks for local providers that may not be installed
		id := p.ID()
		if id == "ollama" || id == "lmstudio" || id == "local" {
			continue
		}

		if !p.IsInstalled() {
			t.Logf("Provider %s not installed, skipping validation", p.Name())
			continue
		}

		err := p.Validate()
		if err != nil {
			t.Logf("Provider %s validation warning: %v", p.Name(), err)
		}
	}
}

// TestProviderGetConfig tests provider configuration
func TestProviderGetConfig(t *testing.T) {
	providers := GetAvailableProviders()

	for _, p := range providers {
		config := p.GetConfig()

		if config["provider"] != p.ID() {
			t.Errorf("Provider %s: expected provider '%s', got '%s'", p.Name(), p.ID(), config["provider"])
		}

		if config["type"] == "" {
			t.Errorf("Provider %s: type should not be empty", p.Name())
		}
	}
}

// TestProviderGetEnvironment tests provider environment variables
func TestProviderGetEnvironment(t *testing.T) {
	providers := GetAvailableProviders()

	for _, p := range providers {
		env := p.GetEnvironment()

		apiKey := p.ID() + "_API_KEY"
		if _, ok := env[apiKey]; !ok {
			t.Logf("Provider %s: API key environment variable not found", p.Name())
		}
	}
}

// TestProviderGetModels tests provider model listing
func TestProviderGetModels(t *testing.T) {
	providers := GetAvailableProviders()

	for _, p := range providers {
		models, err := p.GetModels()
		if err != nil {
			t.Logf("Provider %s: GetModels error: %v", p.Name(), err)
			continue
		}

		if len(models) == 0 {
			t.Logf("Provider %s: no models returned", p.Name())
			continue
		}

		t.Logf("Provider %s: %d models", p.Name(), len(models))
	}
}
