// Package deepseek provides the DeepSeek provider for ProxyBridge.
package provider

import (
	"fmt"
	"os"
)

// DeepSeekProvider implements the Provider interface for DeepSeek.
type DeepSeekProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewDeepSeekProvider creates a new DeepSeek provider.
func NewDeepSeekProvider() *DeepSeekProvider {
	return &DeepSeekProvider{
		id:           "deepseek",
		name:         "DeepSeek",
		description:  "DeepSeek API",
		defaultModel: "deepseek-chat",
	}
}

// Name returns the display name of the provider.
func (p *DeepSeekProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *DeepSeekProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *DeepSeekProvider) Setup() error {
	fmt.Println("DeepSeek provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://platform.deepseek.com/api_keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add deepseek")
	fmt.Println()
	return nil
}

// Validate validates the API key and configuration.
func (p *DeepSeekProvider) Validate() error {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("DEEPSEEK_API_KEY environment variable not set")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid DeepSeek API key format")
	}

	// Test the API key
	if err := p.testAPIKey(apiKey); err != nil {
		return fmt.Errorf("DeepSeek API key validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available DeepSeek models.
func (p *DeepSeekProvider) GetModels() ([]string, error) {
	// DeepSeek models
	models := []string{
		"deepseek-chat",
		"deepseek-coder",
		"deepseek-reasoner",
		"deepseek-chat-v3",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *DeepSeekProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *DeepSeekProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "deepseek",
		"type":     "deepseek",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *DeepSeekProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"DEEPSEEK_API_KEY":      "YOUR_DEEPSEEK_API_KEY",
		"DEEPSEEK_BASE_URL":     "https://api.deepseek.com/v1",
	}
}

// IsInstalled checks if the provider is installed.
func (p *DeepSeekProvider) IsInstalled() bool {
	// DeepSeek API doesn't require installation
	return true
}

// Install installs the provider.
func (p *DeepSeekProvider) Install() error {
	fmt.Println("DeepSeek provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://platform.deepseek.com/api_keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add deepseek")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *DeepSeekProvider) Uninstall() error {
	// Nothing to uninstall for DeepSeek
	return nil
}

// testAPIKey tests if the API key is valid.
func (p *DeepSeekProvider) testAPIKey(apiKey string) error {
	// In production, would make a test request to DeepSeek API
	// For now, just validate format
	return nil
}

