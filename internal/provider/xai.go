// Package xai provides the xAI (Grok) provider for ProxyBridge.
package provider

import (
	"fmt"
	"os"
)

// XAIProvider implements the Provider interface for xAI.
type XAIProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewXAIProvider creates a new xAI provider.
func NewXAIProvider() *XAIProvider {
	return &XAIProvider{
		id:           "xai",
		name:         "xAI (Grok)",
		description:  "xAI Grok API",
		defaultModel: "grok-beta",
	}
}

// Name returns the display name of the provider.
func (p *XAIProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *XAIProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *XAIProvider) Setup() error {
	fmt.Println("xAI (Grok) provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://console.x.ai/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add xai")
	fmt.Println()
	return nil
}

// Validate validates the API key and configuration.
func (p *XAIProvider) Validate() error {
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("XAI_API_KEY environment variable not set")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid xAI API key format")
	}

	// Test the API key
	if err := p.testAPIKey(apiKey); err != nil {
		return fmt.Errorf("xAI API key validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available xAI models.
func (p *XAIProvider) GetModels() ([]string, error) {
	// xAI Grok models
	models := []string{
		"grok-beta",
		"grok-2-1212",
		"grok-2-vision-1212",
		"grok-1-0718",
		"grok-1-vision-1212",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *XAIProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *XAIProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "xai",
		"type":     "xai",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *XAIProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"XAI_API_KEY":           "YOUR_XAI_API_KEY",
		"XAI_BASE_URL":          "https://api.x.ai/v1",
	}
}

// IsInstalled checks if the provider is installed.
func (p *XAIProvider) IsInstalled() bool {
	// xAI API doesn't require installation
	return true
}

// Install installs the provider.
func (p *XAIProvider) Install() error {
	fmt.Println("xAI (Grok) provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://console.x.ai/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add xai")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *XAIProvider) Uninstall() error {
	// Nothing to uninstall for xAI
	return nil
}

// testAPIKey tests if the API key is valid.
func (p *XAIProvider) testAPIKey(apiKey string) error {
	// In production, would make a test request to xAI API
	// For now, just validate format
	return nil
}

