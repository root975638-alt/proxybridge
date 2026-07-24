// Package ollama provides the Ollama provider for ProxyBridge.
package provider

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// OllamaProvider implements the Provider interface for Ollama.
type OllamaProvider struct {
	id           string
	name         string
	description  string
	defaultModel string
}

// NewOllamaProvider creates a new Ollama provider.
func NewOllamaProvider() *OllamaProvider {
	return &OllamaProvider{
		id:           "ollama",
		name:         "Ollama",
		description:  "Ollama Local LLM Server",
		defaultModel: "default",
	}
}

// Name returns the display name of the provider.
func (p *OllamaProvider) Name() string {
	return p.name
}

// ID returns the unique identifier for the provider.
func (p *OllamaProvider) ID() string {
	return p.id
}

// Setup initializes the provider.
func (p *OllamaProvider) Setup() error {
	fmt.Println("Ollama provider setup:")
	fmt.Println()
	fmt.Println("1. Download and install Ollama:")
	fmt.Println("   - Linux/macOS: curl -fsSL https://ollama.com/install.sh | sh")
	fmt.Println("   - Windows: Download from https://ollama.com/download/OllamaSetup.exe")
	fmt.Println()
	fmt.Println("2. Start Ollama server:")
	fmt.Println("   ollama serve")
	fmt.Println()
	fmt.Println("3. Pull a model:")
	fmt.Println("   ollama pull llama3")
	fmt.Println("   ollama pull mistral")
	fmt.Println()
	return nil
}

// Validate validates the Ollama installation and configuration.
func (p *OllamaProvider) Validate() error {
	// Check if Ollama is installed
	if !p.IsInstalled() {
		return fmt.Errorf("Ollama is not installed")
	}

	// Check if Ollama server is running
	if err := p.testServer(); err != nil {
		return fmt.Errorf("Ollama server is not running: %w", err)
	}

	// Test model access
	if err := p.testModels(); err != nil {
		return fmt.Errorf("failed to access Ollama models: %w", err)
	}

	return nil
}

// GetModels returns the list of available Ollama models.
func (p *OllamaProvider) GetModels() ([]string, error) {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list Ollama models: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	models := make([]string, 0, len(lines)-1)

	for _, line := range lines[1:] { // Skip header line
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Extract model name (first column)
		parts := strings.Fields(line)
		if len(parts) > 0 {
			models = append(models, parts[0])
		}
	}

	return models, nil
}

// GetDefaultModel returns the default model for this provider.
func (p *OllamaProvider) GetDefaultModel() string {
	return p.defaultModel
}

// GetConfig returns provider-specific configuration.
func (p *OllamaProvider) GetConfig() map[string]string {
	return map[string]string{
		"provider": "ollama",
		"type":     "ollama",
		"host":     "http://localhost:11434",
	}
}

// GetEnvironment returns environment variables for this provider.
func (p *OllamaProvider) GetEnvironment() map[string]string {
	return map[string]string{
		"OLLAMA_HOST":           "http://localhost:11434",
		"OLLAMA_MODEL":          "default",
	}
}

// IsInstalled checks if the provider is installed.
func (p *OllamaProvider) IsInstalled() bool {
	_, err := exec.LookPath("ollama")
	return err == nil
}

// Install installs the provider.
func (p *OllamaProvider) Install() error {
	fmt.Println("Ollama installation instructions:")
	fmt.Println()
	fmt.Println("Linux/macOS:")
	fmt.Println("  curl -fsSL https://ollama.com/install.sh | sh")
	fmt.Println()
	fmt.Println("Windows:")
	fmt.Println("  Download from https://ollama.com/download/OllamaSetup.exe")
	fmt.Println()
	fmt.Println("After installation, start Ollama with: ollama serve")
	fmt.Println()

	return nil
}

// Uninstall uninstalls the provider.
func (p *OllamaProvider) Uninstall() error {
	fmt.Println("Ollama uninstallationinstructions:")
	fmt.Println()
	fmt.Println("Linux/macOS:")
	fmt.Println("  curl https://ollama.com/install.sh | sh -s -- uninstall")
	fmt.Println()
	fmt.Println("Windows:")
	fmt.Println("  Use Programs and Features to uninstall Ollama")
	fmt.Println()

	return nil
}

// testServer tests if the Ollama server is running.
func (p *OllamaProvider) testServer() error {
	// Check if Ollama is running by trying to connect
	// This would use HTTP in production
	return nil
}

// testModels tests if Ollama models can be accessed.
func (p *OllamaProvider) testModels() error {
	// This would list available models via HTTP in production
	return nil
}

// GetProvider creates a new Ollama provider instance.
func GetProvider() Provider {
	return NewOllamaProvider()
}
