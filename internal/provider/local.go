// Package local provides the Local OpenAI-compatible provider for ProxyBridge.
package provider

import (
	"fmt"
)

// LocalOpenAIProvider implements the Provider interface for OpenAI-compatible servers.
type LocalOpenAIProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
	baseURL      string
	modelName    string
}

// NewLocalOpenAIProvider creates a new local OpenAI-compatible provider.
func NewLocalOpenAIProvider() *LocalOpenAIProvider {
	return &LocalOpenAIProvider{
		id:           "local",
		name:         "Local / Self-Hosted",
		description:  "Any OpenAI-compatible server",
		defaultModel: "default",
		baseURL:      "http://localhost:8000/v1",
		modelName:    "default",
	}
}

// Name returns the display name of the provider.
func (p *LocalOpenAIProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *LocalOpenAIProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *LocalOpenAIProvider) Setup() error {
	fmt.Println("Local OpenAI-compatible provider setup:")
	fmt.Println()
	fmt.Println("Supported servers:")
	fmt.Println("  - vLLM")
	fmt.Println("  - TGI (Text Generation Inference)")
	fmt.Println("  - LM Studio")
	fmt.Println("  - Any OpenAI-compatible server")
	fmt.Println()
	fmt.Println("1. Start your local server")
	fmt.Println("2. Configure the base URL and model:")
	fmt.Println("   proxybridge config edit")
	fmt.Println()
	return nil
}

// Validate validates the local server connection.
func (p *LocalOpenAIProvider) Validate() error {
	// In production, would make a test request to the local server
	// Check if base URL is configured
	if p.baseURL == "" || p.baseURL == "http://localhost:8000/v1" {
		// User needs to configure this
		return fmt.Errorf("base URL not configured. Edit with: proxybridge config edit")
	}

	// Test connection to local server
	if err := p.testConnection(); err != nil {
		return fmt.Errorf("failed to connect to local server: %w", err)
	}

	return nil
}

// GetModels returns the list of available models from the local server.
func (p *LocalOpenAIProvider) GetModels() ([]string, error) {
	// In production, would fetch models from the local server API
	models := []string{
		p.modelName,
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *LocalOpenAIProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *LocalOpenAIProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "local",
		"type":     "openai_compatible",
		"base_url": p.baseURL,
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *LocalOpenAIProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"OPENAI_API_KEY":        "YOUR_API_KEY (optional)",
		"OPENAI_BASE_URL":       p.baseURL,
		"LOCAL_MODEL_NAME":      p.modelName,
	}
}

// IsInstalled checks if the provider is installed.
func (p *LocalOpenAIProvider) IsInstalled() bool {
	// No installation required for local servers
	return true
}

// Install installs the provider.
func (p *LocalOpenAIProvider) Install() error {
	fmt.Println("Local OpenAI-compatible provider setup:")
	fmt.Println()
	fmt.Println("Supported servers:")
	fmt.Println("  - vLLM: https://github.com/vllm-project/vllm")
	fmt.Println("  - TGI: https://github.com/huggingface/text-generation-inference")
	fmt.Println("  - LM Studio: https://lmstudio.ai/")
	fmt.Println()
	fmt.Println("1. Install and start your local server")
	fmt.Println("2. Configure with: proxybridge config edit")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *LocalOpenAIProvider) Uninstall() error {
	// Nothing to uninstall
	return nil
}

// testConnection tests the connection to the local server.
func (p *LocalOpenAIProvider) testConnection() error {
	// In production, would make a test request
	return nil
}

// SetBaseURL sets the base URL for the local server.
func (p *LocalOpenAIProvider) SetBaseURL(url string) {
	p.baseURL = url
}

// SetModelName sets the model name.
func (p *LocalOpenAIProvider) SetModelName(name string) {
	p.modelName = name
}

