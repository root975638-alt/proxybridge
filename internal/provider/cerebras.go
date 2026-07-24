// Package cerebras provides the Cerebras provider for ProxyBridge.
package provider

import (
	"fmt"
	"os"
)

// CerebrasProvider implements the Provider interface for Cerebras.
type CerebrasProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewCerebrasProvider creates a new Cerebras provider.
func NewCerebrasProvider() *CerebrasProvider {
	return &CerebrasProvider{
		id:           "cerebras",
		name:         "Cerebras",
		description:  "Cerebras Cloud API",
		defaultModel: "cerebras-llama3.1-8b",
	}
}

// Name returns the display name of the provider.
func (p *CerebrasProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *CerebrasProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *CerebrasProvider) Setup() error {
	fmt.Println("Cerebras provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://cloud.cerebras.ai/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add cerebras")
	fmt.Println()
	return nil
}

// Validate validates the API key and configuration.
func (p *CerebrasProvider) Validate() error {
	apiKey := os.Getenv("CEREBRAS_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("CEREBRAS_API_KEY environment variable not set")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid Cerebras API key format")
	}

	// Test the API key
	if err := p.testAPIKey(apiKey); err != nil {
		return fmt.Errorf("Cerebras API key validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available Cerebras models.
func (p *CerebrasProvider) GetModels() ([]string, error) {
	// Cerebras models
	models := []string{
		"cerebras-llama3.1-8b",
		"cerebras-llama3.1-70b",
		"cerebras-llama3.1-8b-exp",
		"cerebras-llama3.1-70b-exp",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *CerebrasProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *CerebrasProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "cerebras",
		"type":     "cerebras",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *CerebrasProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"CEREBRAS_API_KEY":      "YOUR_CEREBRAS_API_KEY",
		"CEREBRAS_BASE_URL":     "https://api.cerebras.ai/v1",
	}
}

// IsInstalled checks if the provider is installed.
func (p *CerebrasProvider) IsInstalled() bool {
	// Cerebras API doesn't require installation
	return true
}

// Install installs the provider.
func (p *CerebrasProvider) Install() error {
	fmt.Println("Cerebras provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://cloud.cerebras.ai/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add cerebras")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *CerebrasProvider) Uninstall() error {
	// Nothing to uninstall for Cerebras
	return nil
}

// testAPIKey tests if the API key is valid.
func (p *CerebrasProvider) testAPIKey(apiKey string) error {
	// In production, would make a test request to Cerebras API
	// For now, just validate format
	return nil
}

