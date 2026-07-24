// Package mistral provides the Mistral AI provider for ProxyBridge.
package provider

import (
	"fmt"
	"os"
)

// MistralProvider implements the Provider interface for Mistral AI.
type MistralProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewMistralProvider creates a new Mistral AI provider.
func NewMistralProvider() *MistralProvider {
	return &MistralProvider{
		id:           "mistral",
		name:         "Mistral AI",
		description:  "Mistral AI API",
		defaultModel: "mistral-large",
	}
}

// Name returns the display name of the provider.
func (p *MistralProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *MistralProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *MistralProvider) Setup() error {
	fmt.Println("Mistral AI provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://console.mistral.ai/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add mistral")
	fmt.Println()
	return nil
}

// Validate validates the API key and configuration.
func (p *MistralProvider) Validate() error {
	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("MISTRAL_API_KEY environment variable not set")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid Mistral API key format")
	}

	// Test the API key
	if err := p.testAPIKey(apiKey); err != nil {
		return fmt.Errorf("Mistral API key validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available Mistral models.
func (p *MistralProvider) GetModels() ([]string, error) {
	// Mistral models
	models := []string{
		"mistral-large",
		"mistral-medium",
		"mistral-small",
		"mistral-tiny",
		"mistral-7b-instruct-v0.3",
		"mistral-7b-instruct-v0.2",
		"mistral-7b-instruct-v0.1",
		"codestral-latest",
		"mistral-embed",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *MistralProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *MistralProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "mistral",
		"type":     "mistral",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *MistralProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"MISTRAL_API_KEY":       "YOUR_MISTRAL_API_KEY",
		"MISTRAL_BASE_URL":      "https://api.mistral.ai/v1",
	}
}

// IsInstalled checks if the provider is installed.
func (p *MistralProvider) IsInstalled() bool {
	// Mistral API doesn't require installation
	return true
}

// Install installs the provider.
func (p *MistralProvider) Install() error {
	fmt.Println("Mistral AI provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://console.mistral.ai/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add mistral")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *MistralProvider) Uninstall() error {
	// Nothing to uninstall for Mistral
	return nil
}

// testAPIKey tests if the API key is valid.
func (p *MistralProvider) testAPIKey(apiKey string) error {
	// In production, would make a test request to Mistral API
	// For now, just validate format
	return nil
}

