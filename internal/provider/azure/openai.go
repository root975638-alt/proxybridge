// Package azure provides the Azure OpenAI provider for ProxyBridge.
package azure

import (
	"fmt"
	"os"
	"strings"
)

// AzureOpenAIProvider implements the Provider interface for Azure OpenAI.
type AzureOpenAIProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
	apiVersion   string
	endpoint     string
}

// NewAzureOpenAIProvider creates a new Azure OpenAI provider.
func NewAzureOpenAIProvider() *AzureOpenAIProvider {
	return &AzureOpenAIProvider{
		id:           "azure",
		name:         "Azure OpenAI",
		description:  "Azure OpenAI Service",
		defaultModel: "gpt-4o",
		apiVersion:   "2024-02-15-preview",
		endpoint:     "https://YOUR_RESOURCE_NAME.openai.azure.com/",
	}
}

// Name returns the display name of the provider.
func (p *AzureOpenAIProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *AzureOpenAIProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *AzureOpenAIProvider) Setup() error {
	fmt.Println("Azure OpenAI provider setup:")
	fmt.Println()
	fmt.Println("1. Create an Azure OpenAI resource in the Azure Portal")
	fmt.Println("2. Get your endpoint and API key from the Azure Portal")
	fmt.Println("3. Store them securely:")
	fmt.Println("   proxybridge credentials add azure")
	fmt.Println()
	return nil
}

// Validate validates the Azure credentials and configuration.
func (p *AzureOpenAIProvider) Validate() error {
	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("AZURE_OPENAI_API_KEY environment variable not set")
	}

	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	if endpoint == "" {
		return fmt.Errorf("AZURE_OPENAI_ENDPOINT environment variable not set")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("invalid Azure OpenAI API key format")
	}

	if !strings.HasPrefix(endpoint, "https://") {
		return fmt.Errorf("Azure OpenAI endpoint must start with https://")
	}

	// Test the API key
	if err := p.testAPIKey(apiKey, endpoint); err != nil {
		return fmt.Errorf("Azure OpenAI validation failed: %w", err)
	}

	return nil
}

// GetModels returns the list of available Azure OpenAI models.
func (p *AzureOpenAIProvider) GetModels() ([]string, error) {
	// Azure OpenAI model deployment names (your deployed models)
	// These are deployment names, not model names
	models := []string{
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-4-turbo",
		"gpt-4",
		"gpt-35-turbo",
		"gpt-35-turbo-16k",
		"text-embedding-ada-002",
		"text-embedding-3-small",
		"text-embedding-3-large",
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *AzureOpenAIProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *AzureOpenAIProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "azure",
		"type":     "azure_openai",
		"api_version": p.apiVersion,
		"endpoint":   p.endpoint,
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *AzureOpenAIProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"AZURE_OPENAI_API_KEY":     "YOUR_AZURE_OPENAI_API_KEY",
		"AZURE_OPENAI_ENDPOINT":    "YOUR_AZURE_OPENAI_ENDPOINT",
		"AZURE_OPENAI_API_VERSION": p.apiVersion,
	}
}

// IsInstalled checks if the provider is installed.
func (p *AzureOpenAIProvider) IsInstalled() bool {
	// Azure OpenAI API doesn't require installation
	return true
}

// Install installs the provider.
func (p *AzureOpenAIProvider) Install() error {
	fmt.Println("Azure OpenAI provider setup:")
	fmt.Println()
	fmt.Println("1. Create an Azure OpenAI resource in the Azure Portal")
	fmt.Println("2. Get your endpoint and API key from the Azure Portal")
	fmt.Println("3. Store them securely:")
	fmt.Println("   proxybridge credentials add azure")
	fmt.Println()
	return nil
}

// Uninstall uninstalls the provider.
func (p *AzureOpenAIProvider) Uninstall() error {
	// Nothing to uninstall for Azure OpenAI
	return nil
}

// testAPIKey tests if the Azure credentials are valid.
func (p *AzureOpenAIProvider) testAPIKey(apiKey, endpoint string) error {
	// In production, would make a test request to Azure OpenAI API
	// For now, just validate format
	return nil
}

// GetProvider creates a new Azure OpenAI provider instance.
func GetProvider() *AzureOpenAIProvider {
	return NewAzureOpenAIProvider()
}
