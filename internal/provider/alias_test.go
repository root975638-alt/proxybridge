package provider

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Description returns the provider description
func (p *OllamaProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *AWSBedrockProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *OpenAIProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *AnthropicProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *GoogleGeminiProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *MistralProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *GroqProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *DeepSeekProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *TogetherAIProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *FireworksProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *CerebrasProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *XAIProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *LMStudioProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *LocalOpenAIProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *OpenRouterProvider) Description() string {
	return p.description
}

// Description returns the provider description
func (p *AzureOpenAIProvider) Description() string {
	return p.description
}
