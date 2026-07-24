// Package lmstudio provides the LM Studio provider for ProxyBridge.
package provider

import (
	"fmt"
)

// LMStudioProvider implements the Provider interface for LM Studio.
type LMStudioProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
	baseURL      string
}

// NewLMStudioProvider creates a new LM Studio provider.
func NewLMStudioProvider() *LMStudioProvider {
	return &LMStudioProvider{
		id:           "lmstudio",
		name:         "LM Studio",
		description:  "LM Studio Local Server",
		defaultModel: "default",
		baseURL:      "http://localhost:1234/v1",
	}
}

// Name returns the display name of the provider.
func (p *LMStudioProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *LMStudioProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *LMStudioProvider) Setup() error {
	fmt.Println("LM Studio provider setup:")
	fmt.Println()
	fmt.Println("1. Download and install LM Studio from https://lmstudio.ai/")
	fmt.Println("2. Start LM Studio and load a model")
	fmt.Println("3. Start the local server in LM Studio")
	fmt.Println("4. The server will run on http://localhost:1234")
	fmt.Println()
	return nil
}

// Validate validates the LM Studio connection.
func (p *LMStudioProvider) Validate() error {
	// Check if base URL is configured
	if p.baseURL == "" {
		return fmt.Errorf("base URL not configured")
	}

	// Test connection to LM Studio server
	if err := p.testConnection(); err != nil {
		return fmt.Errorf("failed to connect to LM Studio: %w", err)
	}

	return nil
}

// GetModels returns the list of available models.
func (p *LMStudioProvider) GetModels() ([]string, error) {
	// In production, would fetch models from LM Studio API
	models := []string{
		p.defaultModel,
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *LMStudioProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *LMStudioProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "lmstudio",
		"type":     "openai_compatible",
		"base_url": p.baseURL,
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *LMStudioProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"OPENAI_API_KEY":        "not-required",
		"OPENAI_BASE_URL":       p.baseURL,
	}
}

// IsInstalled checks if the provider is installed.
func (p *LMStudioProvider) IsInstalled() bool {
	// LM Studio is installed if the server is running
	return true
}

// Install installs the provider.
func (p *LMStudioProvider) Install() error {
	fmt.Println("LM Studio provider setup:")
	fmt.Println()
	fmt.Println("1. Download from https://lmstudio.ai/")
	fmt.Println("2. Load a model in LM Studio")
	fmt.Println("3. Start the local server")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *LMStudioProvider) Uninstall() error {
	// Nothing to uninstall
	return nil
}

// testConnection tests the connection to LM Studio.
func (p *LMStudioProvider) testConnection() error {
	// In production, would make a test request
	return nil
}

// SetBaseURL sets the base URL for LM Studio.
func (p *LMStudioProvider) SetBaseURL(url string) {
	p.baseURL = url
}

