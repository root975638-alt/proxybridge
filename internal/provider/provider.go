// Package provider provides the provider plugin interface and registry.
package provider

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

// Provider is the interface that all providers must implement.
type Provider interface {
	// Name returns the display name of the provider
	Name() string

	// ID returns the unique identifier for the provider
	ID() string

	// Setup initializes the provider
	Setup() error

	// Validate validates the provider configuration and credentials
	Validate() error

	// GetModels returns the list of available models
	GetModels() ([]string, error)

	// GetDefaultModel returns the default model for this provider
	GetDefaultModel() string

	// GetConfig returns provider-specific configuration
	GetConfig() map[string]string

	// GetEnvironment returns environment variables for this provider
	GetEnvironment() map[string]string

	// IsInstalled checks if the provider is installed
	IsInstalled() bool

	// Install installs the provider
	Install() error

	// Uninstall uninstalls the provider
	Uninstall() error
}

// ProviderRegistry manages registered providers
type ProviderRegistry struct {
	providers map[string]Provider
}

// RegisterProvider registers a provider
func (r *ProviderRegistry) RegisterProvider(p Provider) error {
	if r.providers == nil {
		r.providers = make(map[string]Provider)
	}

	id := strings.ToLower(p.ID())
	if _, exists := r.providers[id]; exists {
		return fmt.Errorf("provider %s is already registered", id)
	}

	r.providers[id] = p
	return nil
}

// GetProvider returns a registered provider by ID
func (r *ProviderRegistry) GetProvider(id string) Provider {
	id = strings.ToLower(id)
	return r.providers[id]
}

// GetAllProviders returns all registered providers
func (r *ProviderRegistry) GetAllProviders() []Provider {
	providers := make([]Provider, 0, len(r.providers))
	for _, p := range r.providers {
		providers = append(providers, p)
	}
	return providers
}

// GetAvailableProviders returns all available providers
func GetAvailableProviders() []Provider {
	return registry.GetAllProviders()
}

// Manager provides provider management functionality
type Manager struct {
	registry *ProviderRegistry
}

// NewManager creates a new provider manager
func NewManager() (*Manager, error) {
	mgr := &Manager{
		registry: &ProviderRegistry{
			providers: make(map[string]Provider),
		},
	}

	if err := mgr.registerDefaultProviders(); err != nil {
		return nil, err
	}

	return mgr, nil
}

// registerDefaultProviders registers all default providers
func (m *Manager) registerDefaultProviders() error {
	// List of default providers to register
	defaultProviders := []Provider{
		&AWSBedrockProvider{},
		&OpenAIProvider{},
		&AzureOpenAIProvider{},
		&AnthropicProvider{},
		&GoogleGeminiProvider{},
		&OpenRouterProvider{},
		&OllamaProvider{},
		&GroqProvider{},
		&DeepSeekProvider{},
		&MistralProvider{},
		&TogetherAIProvider{},
		&FireworksProvider{},
		&CerebrasProvider{},
		&XAIProvider{},
		&LMStudioProvider{},
		&LocalOpenAIProvider{},
	}

	for _, p := range defaultProviders {
		if err := m.registry.RegisterProvider(p); err != nil {
			return err
		}
	}

	return nil
}

// GetProvider returns a provider by ID
func (m *Manager) GetProvider(id string) Provider {
	return m.registry.GetProvider(id)
}

// GetProviderNames returns all provider names
func (m *Manager) GetProviderNames() []string {
	providers := m.registry.GetAllProviders()
	names := make([]string, 0, len(providers))
	for _, p := range providers {
		names = append(names, p.Name())
	}
	return names
}

// ValidateProvider validates a provider by ID
func (m *Manager) ValidateProvider(id string) error {
	provider := m.GetProvider(id)
	if provider == nil {
		return fmt.Errorf("provider not found: %s", id)
	}
	return provider.Validate()
}

// installProvider installs a provider by ID
func (m *Manager) InstallProvider(id string) error {
	provider := m.GetProvider(id)
	if provider == nil {
		return fmt.Errorf("provider not found: %s", id)
	}
	return provider.Install()
}

// uninstallProvider uninstalls a provider by ID
func (m *Manager) UninstallProvider(id string) error {
	provider := m.GetProvider(id)
	if provider == nil {
		return fmt.Errorf("provider not found: %s", id)
	}
	return provider.Uninstall()
}

// getDefaultModels returns default model mappings for providers
func getDefaultModels() map[string]string {
	return map[string]string{
		"aws":          "anthropic.claude-3-5-sonnet-20241022-v2:0",
		"openai":       "gpt-4o",
		"azure_openai": "gpt-4o",
		"anthropic":    "claude-3-5-sonnet-20241022",
		"google":       "gemini-1.5-flash",
		"openrouter":   "openai/gpt-4o",
		"ollama":       "default",
		"groq":         "llama3-8b-8192",
		"deepseek":     "deepseek-chat",
		"mistral":      "mistral-large",
		"together":     "togethercomputer/llama-3-70b-chat",
		"fireworks":    "accounts/fireworks/models/firefunction-v2",
		"cerebras":     "cerebras-llama3.1-8b",
		"xai":          "grok-beta",
		"lmstudio":     "default",
		"local":        "default",
	}
}

// modelRegex is used to match model names
var modelRegex = regexp.MustCompile(`^[a-z0-9-/]+$`)

// isValidModelName checks if a model name is valid
func isValidModelName(name string) bool {
	return modelRegex.MatchString(name)
}

// getPlatform gets the current platform identifier
func getPlatform() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

// isPlatformSupported checks if a provider supports the current platform
func isPlatformSupported(supportedPlatforms []string) bool {
	current := getPlatform()
	for _, platform := range supportedPlatforms {
		if strings.EqualFold(platform, current) || strings.EqualFold(platform, "all") {
			return true
		}
	}
	// If no platforms specified, assume all are supported
	return len(supportedPlatforms) == 0
}

// getProviderIDFromName converts a display name to an ID
func getProviderIDFromName(name string) string {
	id := strings.ToLower(name)
	id = strings.ReplaceAll(id, " ", "_")
	id = strings.ReplaceAll(id, "-", "_")
	return id
}

// ProviderStatus represents the status of a provider
type ProviderStatus struct {
	Installed  bool
	Validated  bool
	Models     []string
	DefaultModel string
	Error      string
}

// GetProviderStatus returns the status of a provider
func GetProviderStatus(p Provider) ProviderStatus {
	status := ProviderStatus{
		Installed:    p.IsInstalled(),
		Validated:    false,
		DefaultModel: p.GetDefaultModel(),
	}

	if status.Installed {
		if err := p.Validate(); err == nil {
			status.Validated = true
		} else {
			status.Error = err.Error()
		}

		models, _ := p.GetModels()
		status.Models = models
	}

	return status
}
