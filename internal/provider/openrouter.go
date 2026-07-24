// Package openrouter provides the OpenRouter provider for ProxyBridge.
package provider

import (
	"fmt"
	"os"
)

// OpenRouterProvider implements the Provider interface for OpenRouter.
type OpenRouterProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewOpenRouterProvider creates a new OpenRouter provider.
func NewOpenRouterProvider() *OpenRouterProvider {
	return &OpenRouterProvider{
		id:           "openrouter",
		name:         "OpenRouter",
		description:  "OpenRouter API",
		defaultModel: "openai/gpt-4o",
	}
}

// Name returns the display name of the provider.
func (p *OpenRouterProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *OpenRouterProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *OpenRouterProvider) Setup() error {
	fmt.Println("OpenRouter provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://openrouter.ai/keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add openrouter")
	fmt.Println()
	return nil
}

// Validate validates the API key and configuration.
func (p *OpenRouterProvider) Validate() error {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("OPENROUTER_API_KEY environment variable not set")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid OpenRouter API key format")
	}

	// Test the API key
	if err := p.testAPIKey(apiKey); err != nil {
		return fmt.Errorf("OpenRouter API key validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available OpenRouter models.
func (p *OpenRouterProvider) GetModels() ([]string, error) {
	// OpenRouter supports many models
	// Here are the most popular ones
	models := []string{
		"openai/gpt-4o",
		"openai/gpt-4o-mini",
		"openai/gpt-4-turbo",
		"openai/gpt-3.5-turbo",
		"anthropic/claude-3-5-sonnet",
		"anthropic/claude-3-opus",
		"anthropic/claude-3-haiku",
		"google/gemini-1.5-flash",
		"google/gemini-1.5-pro",
		"meta-llama/llama-3-70b-instruct",
		"meta-llama/llama-3-8b-instruct",
		"mistralai/mistral-7b-instruct",
		"mistralai/mixtral-8x7b-instruct",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *OpenRouterProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *OpenRouterProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "openrouter",
		"type":     "openrouter",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *OpenRouterProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"OPENROUTER_API_KEY":    "YOUR_OPENROUTER_API_KEY",
		"OPENROUTER_BASE_URL":   "https://openrouter.ai/api/v1",
	}
}

// IsInstalled checks if the provider is installed.
func (p *OpenRouterProvider) IsInstalled() bool {
	// OpenRouter API doesn't require installation
	return true
}

// Install installs the provider.
func (p *OpenRouterProvider) Install() error {
	fmt.Println("OpenRouter provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://openrouter.ai/keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add openrouter")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *OpenRouterProvider) Uninstall() error {
	// Nothing to uninstall for OpenRouter
	return nil
}

// testAPIKey tests if the API key is valid.
func (p *OpenRouterProvider) testAPIKey(apiKey string) error {
	// In production, would make a test request to OpenRouter API
	// For now, just validate format
	return nil
}

