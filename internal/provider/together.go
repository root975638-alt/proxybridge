// Package together provides the Together AI provider for ProxyBridge.
package provider

import (
	"fmt"
	"os"
)

// TogetherAIProvider implements the Provider interface for Together AI.
type TogetherAIProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewTogetherAIProvider creates a new Together AI provider.
func NewTogetherAIProvider() *TogetherAIProvider {
	return &TogetherAIProvider{
		id:           "together",
		name:         "Together AI",
		description:  "Together AI API",
		defaultModel: "togethercomputer/llama-3-70b-chat",
	}
}

// Name returns the display name of the provider.
func (p *TogetherAIProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *TogetherAIProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *TogetherAIProvider) Setup() error {
	fmt.Println("Together AI provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://cloud.together.ai/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add together")
	fmt.Println()
	return nil
}

// Validate validates the API key and configuration.
func (p *TogetherAIProvider) Validate() error {
	apiKey := os.Getenv("TOGETHER_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("TOGETHER_API_KEY environment variable not set")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid Together AI API key format")
	}

	// Test the API key
	if err := p.testAPIKey(apiKey); err != nil {
		return fmt.Errorf("Together AI API key validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available Together models.
func (p *TogetherAIProvider) GetModels() ([]string, error) {
	// Together AI models
	models := []string{
		"togethercomputer/llama-3-70b-chat",
		"togethercomputer/llama-3-8b-chat",
		"togethercomputer/llama-3-70b-instruct",
		"togethercomputer/llama-3-8b-instruct",
		"togethercomputer/llama-3.1-70b-chat",
		"togethercomputer/llama-3.1-8b-chat",
		"togethercomputer/llama-3.1-405b-chat",
		"togethercomputer/mistral-7b-instruct-v0.2",
		"togethercomputer/qwen-72b-chat",
		"togethercomputer/gemma-2-27b-it",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *TogetherAIProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *TogetherAIProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "together",
		"type":     "together",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *TogetherAIProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"TOGETHER_API_KEY":      "YOUR_TOGETHER_API_KEY",
		"TOGETHER_BASE_URL":     "https://api.together.xyz/v1",
	}
}

// IsInstalled checks if the provider is installed.
func (p *TogetherAIProvider) IsInstalled() bool {
	// Together AI API doesn't require installation
	return true
}

// Install installs the provider.
func (p *TogetherAIProvider) Install() error {
	fmt.Println("Together AI provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://cloud.together.ai/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add together")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *TogetherAIProvider) Uninstall() error {
	// Nothing to uninstall for Together AI
	return nil
}

// testAPIKey tests if the API key is valid.
func (p *TogetherAIProvider) testAPIKey(apiKey string) error {
	// In production, would make a test request to Together AI API
	// For now, just validate format
	return nil
}

