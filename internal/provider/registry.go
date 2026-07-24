// Package provider provides a central registry for all providers.
package provider

import (
	// Import all provider implementations
	_ "github.com/proxybridge/cli/internal/provider/aws"
	_ "github.com/proxybridge/cli/internal/provider/google"
	_ "github.com/proxybridge/cli/internal/provider/azure"
	"github.com/proxybridge/cli/internal/provider"
	"github.com/proxybridge/cli/internal/provider/anthropic"
	"github.com/proxybridge/cli/internal/provider/deepseek"
	"github.com/proxybridge/cli/internal/provider/fireworks"
	"github.com/proxybridge/cli/internal/provider/groq"
	"github.com/proxybridge/cli/internal/provider/local"
	"github.com/proxybridge/cli/internal/provider/lmstudio"
	"github.com/proxybridge/cli/internal/provider/mistral"
	"github.com/proxybridge/cli/internal/provider/ollama"
	"github.com/proxybridge/cli/internal/provider/openai"
	"github.com/proxybridge/cli/internal/provider/openrouter"
	"github.com/proxybridge/cli/internal/provider/together"
	"github.com/proxybridge/cli/internal/provider/xai"
	"github.com/proxybridge/cli/internal/provider/cerebras"
)

// Registry holds all registered providers
var Registry = provider.NewProviderRegistry()

func init() {
	// Register all providers
	registerProviders()
}

// registerProviders registers all providers at init time
func registerProviders() {
	registry := &ProviderRegistry{
		providers: make(map[string]Provider),
	}

	providers := []Provider{
		// Cloud providers
		openai.GetProvider(),
		anthropic.GetProvider(),
		azure.GetProvider(),
		google.GetProvider(),
		openrouter.GetProvider(),
		xai.GetProvider(),

		// Alternative cloud providers
		groq.GetProvider(),
		deepseek.GetProvider(),
		mistral.GetProvider(),
		together.GetProvider(),
		fireworks.GetProvider(),
		cerebras.GetProvider(),

		// Local/On-premise providers
		ollama.GetProvider(),
		lmstudio.GetProvider(),
		local.GetProvider(),

		// AWS-specific
		aws.GetProvider(),
	}

	for _, p := range providers {
		registry.RegisterProvider(p)
	}
}

// NewProviderRegistry creates a new provider registry with all providers registered
func NewProviderRegistry() *ProviderRegistry {
	reg := &ProviderRegistry{
		providers: make(map[string]Provider),
	}
	registerProvidersInRegistry(reg)
	return reg
}

// registerProvidersInRegistry registers all providers into a given registry
func registerProvidersInRegistry(reg *ProviderRegistry) {
	// Cloud providers
	_ = reg.RegisterProvider(openai.GetProvider())
	_ = reg.RegisterProvider(anthropic.GetProvider())
	_ = reg.RegisterProvider(azure.GetProvider())
	_ = reg.RegisterProvider(google.GetProvider())
	_ = reg.RegisterProvider(openrouter.GetProvider())
	_ = reg.RegisterProvider(xai.GetProvider())

	// Alternative cloud providers
	_ = reg.RegisterProvider(groq.GetProvider())
	_ = reg.RegisterProvider(deepseek.GetProvider())
	_ = reg.RegisterProvider(mistral.GetProvider())
	_ = reg.RegisterProvider(together.GetProvider())
	_ = reg.RegisterProvider(fireworks.GetProvider())
	_ = reg.RegisterProvider(cerebras.GetProvider())

	// Local/On-premise providers
	_ = reg.RegisterProvider(ollama.GetProvider())
	_ = reg.RegisterProvider(lmstudio.GetProvider())
	_ = reg.RegisterProvider(local.GetProvider())

	// AWS-specific
	_ = reg.RegisterProvider(aws.GetProvider())
}

// GetProvider returns a provider by ID
func GetProvider(id string) Provider {
	return Registry.GetProvider(id)
}

// GetAllProviders returns all registered providers
func GetAllProviders() []Provider {
	return Registry.GetAllProviders()
}

// GetProviderNames returns all provider names
func GetProviderNames() []string {
	return Registry.GetProviderNames()
}

// ValidateProvider validates a provider by ID
func ValidateProvider(id string) error {
	return Registry.ValidateProvider(id)
}

// InstallProvider installs a provider by ID
func InstallProvider(id string) error {
	return Registry.InstallProvider(id)
}

// UninstallProvider uninstalls a provider by ID
func UninstallProvider(id string) error {
	return Registry.UninstallProvider(id)
}
