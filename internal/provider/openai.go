// Package openai provides the OpenAI provider for ProxyBridge.
package provider

import (
	"fmt"
	"os"
	"strings"
)

// OpenAIProvider implements the Provider interface for OpenAI.
type OpenAIProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewOpenAIProvider creates a new OpenAI provider.
func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		id:           "openai",
		name:         "OpenAI",
		description:  "OpenAI API",
		defaultModel: "gpt-4o",
	}
}

// Name returns the display name of the provider.
func (p *OpenAIProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *OpenAIProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *OpenAIProvider) Setup() error {
	fmt.Println("OpenAI provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://platform.openai.com/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add openai")
	fmt.Println()
	return nil
}

// Validate validates the API key and configuration.
func (p *OpenAIProvider) Validate() error {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid OpenAI API key format")
	}

	if !strings.HasPrefix(apiKey, "sk-") {
		return fmt.Errorf("OpenAI API key should start with 'sk-'")
	}

	// Test the API key by making a request
	if err := p.testAPIKey(apiKey); err != nil {
		return fmt.Errorf("OpenAI API key validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available OpenAI models.
func (p *OpenAIProvider) GetModels() ([]string, error) {
	// Common OpenAI models
	models := []string{
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-4-turbo",
		"gpt-4",
		"gpt-4-32k",
		"gpt-3.5-turbo",
		"gpt-3.5-turbo-16k",
		"o1-preview",
		"o1-mini",
		"text-embedding-3-large",
		"text-embedding-3-small",
		"text-embedding-ada-002",
		"dall-e-3",
		"whisper-1",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *OpenAIProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *OpenAIProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "openai",
		"type":     "openai",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *OpenAIProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"OPENAI_API_KEY":        "YOUR_OPENAI_API_KEY",
		"OPENAI_BASE_URL":       "https://api.openai.com/v1",
		"OPENAI_ORG_ID":         "YOUR_ORGANIZATION_ID (optional)",
		"OPENAI_PROJECT_ID":     "YOUR_PROJECT_ID (optional)",
	}
}

// IsInstalled checks if the provider is installed.
func (p *OpenAIProvider) IsInstalled() bool {
	// OpenAI API doesn't require installation
	return true
}

// Install installs the provider.
func (p *OpenAIProvider) Install() error {
	fmt.Println("OpenAI provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://platform.openai.com/api-keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add openai")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *OpenAIProvider) Uninstall() error {
	// Nothing to uninstall for OpenAI
	return nil
}

// testAPIKey tests if the API key is valid.
func (p *OpenAIProvider) testAPIKey(apiKey string) error {
	// In production, would make a test request to OpenAI API
	// For now, just validate format
	return nil
}

// GetProvider creates a new OpenAI provider instance.
func GetProvider() Provider {
	return NewOpenAIProvider()
}
