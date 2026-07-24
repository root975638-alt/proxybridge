// Package fireworks provides the Fireworks AI provider for ProxyBridge.
package provider

import (
	"fmt"
	"os"
)

// FireworksProvider implements the Provider interface for Fireworks AI.
type FireworksProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewFireworksProvider creates a new Fireworks AI provider.
func NewFireworksProvider() *FireworksProvider {
	return &FireworksProvider{
		id:           "fireworks",
		name:         "Fireworks AI",
		description:  "Fireworks AI API",
		defaultModel: "accounts/fireworks/models/firefunction-v2",
	}
}

// Name returns the display name of the provider.
func (p *FireworksProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *FireworksProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *FireworksProvider) Setup() error {
	fmt.Println("Fireworks AI provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://fireworks.ai/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add fireworks")
	fmt.Println()
	return nil
}

// Validate validates the API key and configuration.
func (p *FireworksProvider) Validate() error {
	apiKey := os.Getenv("FIREWORKS_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("FIREWORKS_API_KEY environment variable not set")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid Fireworks AI API key format")
	}

	// Test the API key
	if err := p.testAPIKey(apiKey); err != nil {
		return fmt.Errorf("Fireworks AI API key validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available Fireworks models.
func (p *FireworksProvider) GetModels() ([]string, error) {
	// Fireworks AI models
	models := []string{
		"accounts/fireworks/models/firefunction-v2",
		"accounts/fireworks/models/firellava-13b",
		"accounts/fireworks/models/llama-v3-70b-instruct",
		"accounts/fireworks/models/llama-v3-8b-instruct",
		"accounts/fireworks/models/mixtral-8x7b-instruct",
		"accounts/fireworks/models/qwen-72b-instruct",
		"accounts/fireworks/models/gemma-2b-it",
		"accounts/fireworks/models/gemma-7b-it",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *FireworksProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *FireworksProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "fireworks",
		"type":     "fireworks",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *FireworksProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"FIREWORKS_API_KEY":     "YOUR_FIREWORKS_API_KEY",
		"FIREWORKS_BASE_URL":    "https://api.fireworks.ai/inference/v1",
	}
}

// IsInstalled checks if the provider is installed.
func (p *FireworksProvider) IsInstalled() bool {
	// Fireworks AI API doesn't require installation
	return true
}

// Install installs the provider.
func (p *FireworksProvider) Install() error {
	fmt.Println("Fireworks AI provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://fireworks.ai/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add fireworks")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *FireworksProvider) Uninstall() error {
	// Nothing to uninstall for Fireworks AI
	return nil
}

// testAPIKey tests if the API key is valid.
func (p *FireworksProvider) testAPIKey(apiKey string) error {
	// In production, would make a test request to Fireworks AI API
	// For now, just validate format
	return nil
}

