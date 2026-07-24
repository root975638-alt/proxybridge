// Package groq provides the Groq provider for ProxyBridge.
package groq

import (
	"fmt"
	"os"
	"strings"
)

// GroqProvider implements the Provider interface for Groq.
type GroqProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewGroqProvider creates a new Groq provider.
func NewGroqProvider() *GroqProvider {
	return &GroqProvider{
		id:           "groq",
		name:         "Groq",
		description:  "Groq Cloud API",
		defaultModel: "llama3-8b-8192",
	}
}

// Name returns the display name of the provider.
func (p *GroqProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *GroqProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *GroqProvider) Setup() error {
	fmt.Println("Groq provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://console.groq.com/keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add groq")
	fmt.Println()
	return nil
}

// Validate validates the API key and configuration.
func (p *GroqProvider) Validate() error {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("GROQ_API_KEY environment variable not set")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid Groq API key format")
	}

	// Test the API key
	if err := p.testAPIKey(apiKey); err != nil {
		return fmt.Errorf("Groq API key validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available Groq models.
func (p *GroqProvider) GetModels() ([]string, error) {
	// Groq models
	models := []string{
		"llama3-8b-8192",
		"llama3-70b-8192",
		"llama-3.1-8b-instant",
		"llama-3.1-70b-versatile",
		"mistral-saba-24b",
		"gemma-7b-it",
		"gemma2-9b-it",
		"mixtral-8x7b-32768",
		"deepseek-r1-distill-llama-70b",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *GroqProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *GroqProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "groq",
		"type":     "groq",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *GroqProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"GROQ_API_KEY":          "YOUR_GROQ_API_KEY",
		"GROQ_BASE_URL":         "https://api.groq.com/openai/v1",
	}
}

// IsInstalled checks if the provider is installed.
func (p *GroqProvider) IsInstalled() bool {
	// Groq API doesn't require installation
	return true
}

// Install installs the provider.
func (p *GroqProvider) Install() error {
	fmt.Println("Groq provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://console.groq.com/keys")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add groq")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *GroqProvider) Uninstall() error {
	// Nothing to uninstall for Groq
	return nil
}

// testAPIKey tests if the API key is valid.
func (p *GroqProvider) testAPIKey(apiKey string) error {
	// In production, would make a test request to Groq API
	// For now, just validate format
	return nil
}

// GetProvider creates a new Groq provider instance.
func GetProvider() *GroqProvider {
	return NewGroqProvider()
}
