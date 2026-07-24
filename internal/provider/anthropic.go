// Package anthropic provides the Anthropic provider for ProxyBridge.
package provider

import (
	"fmt"
	"os"
	"strings"
)

// AnthropicProvider implements the Provider interface for Anthropic.
type AnthropicProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewAnthropicProvider creates a new Anthropic provider.
func NewAnthropicProvider() *AnthropicProvider {
	return &AnthropicProvider{
		id:           "anthropic",
		name:         "Anthropic",
		description:  "Anthropic API (Claude)",
		defaultModel: "claude-3-5-sonnet-20241022",
	}
}

// Name returns the display name of the provider.
func (p *AnthropicProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *AnthropicProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *AnthropicProvider) Setup() error {
	fmt.Println("Anthropic provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://console.anthropic.com/settings/keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add anthropic")
	fmt.Println()
	return nil
}

// Validate validates the API key and configuration.
func (p *AnthropicProvider) Validate() error {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid Anthropic API key format")
	}

	if !strings.HasPrefix(apiKey, "sk-ant-") {
		return fmt.Errorf("Anthropic API key should start with 'sk-ant-'")
	}

	// Test the API key
	if err := p.testAPIKey(apiKey); err != nil {
		return fmt.Errorf("Anthropic API key validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available Anthropic models.
func (p *AnthropicProvider) GetModels() ([]string, error) {
	// Anthropic Claude models
	models := []string{
		"claude-3-5-sonnet-20241022",
		"claude-3-5-sonnet-20240620",
		"claude-3-7-sonnet-20250219",
		"claude-3-5-haiku-20241022",
		"claude-3-opus-20240229",
		"claude-3-haiku-20240307",
		"claude-2.1",
		"claude-2.0",
		"claude-instant-1.2",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *AnthropicProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *AnthropicProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "anthropic",
		"type":     "anthropic",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *AnthropicProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"ANTHROPIC_API_KEY":     "YOUR_ANTHROPIC_API_KEY",
		"ANTHROPIC_BASE_URL":    "https://api.anthropic.com/v1",
		"ANTHROPIC_VERSION":     "2023-06-01",
	}
}

// IsInstalled checks if the provider is installed.
func (p *AnthropicProvider) IsInstalled() bool {
	// Anthropic API doesn't require installation
	return true
}

// Install installs the provider.
func (p *AnthropicProvider) Install() error {
	fmt.Println("Anthropic provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://console.anthropic.com/settings/keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add anthropic")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *AnthropicProvider) Uninstall() error {
	// Nothing to uninstall for Anthropic
	return nil
}

// testAPIKey tests if the API key is valid.
func (p *AnthropicProvider) testAPIKey(apiKey string) error {
	// In production, would make a test request to Anthropic API
	// For now, just validate format
	return nil
}

// GetProvider creates a new Anthropic provider instance.
func GetProvider() Provider {
	return NewAnthropicProvider()
}
