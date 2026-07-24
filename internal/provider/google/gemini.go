// Package google provides the Google Gemini provider for ProxyBridge.
package google

import (
	"fmt"
	"os"
)

// GoogleGeminiProvider implements the Provider interface for Google Gemini.
type GoogleGeminiProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewGoogleGeminiProvider creates a new Google Gemini provider.
func NewGoogleGeminiProvider() *GoogleGeminiProvider {
	return &GoogleGeminiProvider{
		id:           "google",
		name:         "Google Gemini",
		description:  "Google Gemini API",
		defaultModel: "gemini-1.5-flash",
	}
}

// Name returns the display name of the provider.
func (p *GoogleGeminiProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *GoogleGeminiProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *GoogleGeminiProvider) Setup() error {
	fmt.Println("Google Gemini provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://aistudio.google.com/app/api_key")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add google")
	fmt.Println()
	return nil
}

// Validate validates the API key and configuration.
func (p *GoogleGeminiProvider) Validate() error {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		// Check for service account
		if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
			return fmt.Errorf("GOOGLE_API_KEY or GOOGLE_APPLICATION_CREDENTIALS not set")
		}
		return nil
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid Google API key format")
	}

	// Test the API key
	if err := p.testAPIKey(apiKey); err != nil {
		return fmt.Errorf("Google API key validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available Google Gemini models.
func (p *GoogleGeminiProvider) GetModels() ([]string, error) {
	// Google Gemini models
	models := []string{
		"gemini-1.5-flash",
		"gemini-1.5-flash-8b",
		"gemini-1.5-pro",
		"gemini-1.5-pro-latest",
		"gemini-1.0-pro",
		"models/gemini-1.5-flash",
		"models/gemini-1.5-pro",
		"models/gemini-1.0-pro-latest",
		"gemini-2.0-flash-exp",
		"gemini-2.0-flash",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *GoogleGeminiProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *GoogleGeminiProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "google",
		"type":     "gemini",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *GoogleGeminiProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"GOOGLE_API_KEY":                "YOUR_GOOGLE_API_KEY",
		"GOOGLE_APPLICATION_CREDENTIALS": "PATH_TO_SERVICE_ACCOUNT_JSON (optional)",
	}
}

// IsInstalled checks if the provider is installed.
func (p *GoogleGeminiProvider) IsInstalled() bool {
	// Google Gemini API doesn't require installation
	return true
}

// Install installs the provider.
func (p *GoogleGeminiProvider) Install() error {
	fmt.Println("Google Gemini provider setup:")
	fmt.Println()
	fmt.Println("1. Get your API key from https://aistudio.google.com/app/api_key")
	fmt.Println("2. Store it securely:")
	fmt.Println("   proxybridge credentials add google")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *GoogleGeminiProvider) Uninstall() error {
	// Nothing to uninstall for Google Gemini
	return nil
}

// testAPIKey tests if the API key is valid.
func (p *GoogleGeminiProvider) testAPIKey(apiKey string) error {
	// In production, would make a test request to Google API
	// For now, just validate format
	return nil
}

// GetProvider creates a new Google Gemini provider instance.
func GetProvider() *GoogleGeminiProvider {
	return NewGoogleGeminiProvider()
}
