// Package provider provides a central registry for all providers.
package provider

// Registry holds all registered providers
var Registry = NewProviderRegistry()

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
		NewOpenAIProvider(),
		NewAnthropicProvider(),
		newAzureOpenAIProvider(),
		newGoogleGeminiProvider(),
		NewOpenRouterProvider(),
		NewXAIProvider(),

		// Alternative cloud providers
		NewGroqProvider(),
		NewDeepSeekProvider(),
		NewMistralProvider(),
		NewTogetherAIProvider(),
		NewFireworksProvider(),
		NewCerebrasProvider(),

		// Local/On-premise providers
		NewOllamaProvider(),
		NewLMStudioProvider(),
		NewLocalOpenAIProvider(),

		// AWS-specific
		newAWSBedrockProvider(),
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
	_ = reg.RegisterProvider(NewOpenAIProvider())
	_ = reg.RegisterProvider(NewAnthropicProvider())
	_ = reg.RegisterProvider(newAzureOpenAIProvider())
	_ = reg.RegisterProvider(newGoogleGeminiProvider())
	_ = reg.RegisterProvider(NewOpenRouterProvider())
	_ = reg.RegisterProvider(NewXAIProvider())

	// Alternative cloud providers
	_ = reg.RegisterProvider(NewGroqProvider())
	_ = reg.RegisterProvider(NewDeepSeekProvider())
	_ = reg.RegisterProvider(NewMistralProvider())
	_ = reg.RegisterProvider(NewTogetherAIProvider())
	_ = reg.RegisterProvider(NewFireworksProvider())
	_ = reg.RegisterProvider(NewCerebrasProvider())

	// Local/On-premise providers
	_ = reg.RegisterProvider(NewOllamaProvider())
	_ = reg.RegisterProvider(NewLMStudioProvider())
	_ = reg.RegisterProvider(NewLocalOpenAIProvider())

	// AWS-specific
	_ = reg.RegisterProvider(newAWSBedrockProvider())
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
	return Registry.getAllProviderNames()
}

// ValidateProvider validates a provider by ID
func ValidateProvider(id string) error {
	return Registry.validateProvider(id)
}

// InstallProvider installs a provider by ID
func InstallProvider(id string) error {
	return Registry.installProvider(id)
}

// UninstallProvider uninstalls a provider by ID
func UninstallProvider(id string) error {
	return Registry.uninstallProvider(id)
}
