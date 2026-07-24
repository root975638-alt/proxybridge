// Package alias provides model alias management.
package alias

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/proxybridge/cli/internal/config"
	"gopkg.in/yaml.v3"
)

const (
	aliasesFileName = "aliases.yaml"
)

// Alias represents a model alias mapping.
type Alias struct {
	Name       string   `yaml:"name"`
	Model      string   `yaml:"model"`
	Provider   string   `yaml:"provider"`
	Enabled    bool     `yaml:"enabled"`
	Notes      string   `yaml:"notes,omitempty"`
}

// Manager manages model aliases.
type Manager struct {
	configDir string
}

// NewManager creates a new alias manager.
func NewManager() (*Manager, error) {
	configDir, err := config.GetConfigDirectory()
	if err != nil {
		return nil, err
	}
	return &Manager{configDir: configDir}, nil
}

// GetAliases returns all configured aliases.
func (m *Manager) GetAliases() (map[string]*Alias, error) {
	aliasesPath := filepath.Join(m.configDir, aliasesFileName)

	data, err := os.ReadFile(aliasesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return m.getDefaultAliases(), nil
		}
		return nil, fmt.Errorf("failed to read aliases file: %w", err)
	}

	var aliases map[string]*Alias
	if err := yaml.Unmarshal(data, &aliases); err != nil {
		return nil, fmt.Errorf("failed to parse aliases: %w", err)
	}

	if aliases == nil {
		return m.getDefaultAliases(), nil
	}

	return aliases, nil
}

// SetAlias sets a model alias.
func (m *Manager) SetAlias(name, model, provider string) error {
	aliases, err := m.GetAliases()
	if err != nil {
		return err
	}

	aliases[name] = &Alias{
		Name:     name,
		Model:    model,
		Provider: provider,
		Enabled:  true,
	}

	return m.saveAliases(aliases)
}

// RemoveAlias removes a model alias.
func (m *Manager) RemoveAlias(name string) error {
	aliases, err := m.GetAliases()
	if err != nil {
		return err
	}

	delete(aliases, name)
	return m.saveAliases(aliases)
}

// ResolveAlias resolves an alias to its model and provider.
func (m *Manager) ResolveAlias(alias string) (model, provider string, err error) {
	aliases, err := m.GetAliases()
	if err != nil {
		return "", "", err
	}

	a, ok := aliases[alias]
	if !ok {
		// Check if it's a direct model name
		return alias, "", nil
	}

	if !a.Enabled {
		return "", "", fmt.Errorf("alias is disabled: %s", alias)
	}

	return a.Model, a.Provider, nil
}

// ListAliases returns a list of all enabled aliases.
func (m *Manager) ListAliases() ([]*Alias, error) {
	aliases, err := m.GetAliases()
	if err != nil {
		return nil, err
	}

	result := make([]*Alias, 0, len(aliases))
	for _, a := range aliases {
		if a.Enabled {
			result = append(result, a)
		}
	}

	return result, nil
}

// getDefaultAliases returns the default alias mappings.
func (m *Manager) getDefaultAliases() map[string]*Alias {
	return map[string]*Alias{
		"claude-3-5-sonnet": {
			Name:     "claude-3-5-sonnet",
			Model:    "claude-3-5-sonnet-20241022",
			Provider: "anthropic",
			Enabled:  true,
			Notes:    "Anthropic Claude 3.5 Sonnet",
		},
		"claude-opus": {
			Name:     "claude-opus",
			Model:    "claude-opus-4-20250502",
			Provider: "anthropic",
			Enabled:  true,
			Notes:    "Anthropic Claude Opus",
		},
		"claude-3-haiku": {
			Name:     "claude-3-haiku",
			Model:    "claude-3-haiku-20240307",
			Provider: "anthropic",
			Enabled:  true,
			Notes:    "Anthropic Claude 3 Haiku",
		},
		"claude-3-7-sonnet": {
			Name:     "claude-3-7-sonnet",
			Model:    "claude-3-7-sonnet-20250219",
			Provider: "anthropic",
			Enabled:  true,
			Notes:    "Anthropic Claude 3.7 Sonnet",
		},
		"claude-3-5-haiku": {
			Name:     "claude-3-5-haiku",
			Model:    "claude-3-5-haiku-20241022",
			Provider: "anthropic",
			Enabled:  true,
			Notes:    "Anthropic Claude 3.5 Haiku",
		},
		"claude-3-opus": {
			Name:     "claude-3-opus",
			Model:    "claude-3-opus-20240229",
			Provider: "anthropic",
			Enabled:  true,
			Notes:    "Anthropic Claude 3 Opus",
		},
		"claude-sonnet": {
			Name:     "claude-sonnet",
			Model:    "claude-3-5-sonnet-20241022",
			Provider: "anthropic",
			Enabled:  true,
			Notes:    "Anthropic Claude 3.5 Sonnet (alias)",
		},
		"claude": {
			Name:     "claude",
			Model:    "claude-3-5-sonnet-20241022",
			Provider: "anthropic",
			Enabled:  true,
			Notes:    "Default Claude model",
		},
		"coder": {
			Name:     "coder",
			Model:    "claude-3-5-sonnet-20241022",
			Provider: "anthropic",
			Enabled:  true,
			Notes:    "Best coding model",
		},
		"reasoning": {
			Name:     "reasoning",
			Model:    "claude-3-7-sonnet-20250219",
			Provider: "anthropic",
			Enabled:  true,
			Notes:    "Best reasoning model",
		},
		"deepseek": {
			Name:     "deepseek",
			Model:    "deepseek-chat",
			Provider: "deepseek",
			Enabled:  true,
			Notes:    "DeepSeek Chat",
		},
		"grok": {
			Name:     "grok",
			Model:    "grok-beta",
			Provider: "xai",
			Enabled:  true,
			Notes:    "xAI Grok Beta",
		},
		"grok-2": {
			Name:     "grok-2",
			Model:    "grok-2-1212",
			Provider: "xai",
			Enabled:  true,
			Notes:    "xAI Grok-2",
		},
		"qwen": {
			Name:     "qwen",
			Model:    "qwen-plus",
			Provider: "alibaba",
			Enabled:  true,
			Notes:    "Alibaba Qwen Plus",
		},
		"gemini": {
			Name:     "gemini",
			Model:    "gemini-1.5-flash",
			Provider: "google",
			Enabled:  true,
			Notes:    "Google Gemini Flash",
		},
		"gemini-1.5-flash": {
			Name:     "gemini-1.5-flash",
			Model:    "gemini-1.5-flash",
			Provider: "google",
			Enabled:  true,
			Notes:    "Google Gemini 1.5 Flash",
		},
		"gemini-1.5-pro": {
			Name:     "gemini-1.5-pro",
			Model:    "gemini-1.5-pro",
			Provider: "google",
			Enabled:  true,
			Notes:    "Google Gemini 1.5 Pro",
		},
		"gemini-2.0-flash": {
			Name:     "gemini-2.0-flash",
			Model:    "gemini-2.0-flash-exp",
			Provider: "google",
			Enabled:  true,
			Notes:    "Google Gemini 2.0 Flash",
		},
		"gpt-4o": {
			Name:     "gpt-4o",
			Model:    "gpt-4o",
			Provider: "openai",
			Enabled:  true,
			Notes:    "OpenAI GPT-4o",
		},
		"gpt-4o-mini": {
			Name:     "gpt-4o-mini",
			Model:    "gpt-4o-mini",
			Provider: "openai",
			Enabled:  true,
			Notes:    "OpenAI GPT-4o Mini",
		},
		"gpt-4-turbo": {
			Name:     "gpt-4-turbo",
			Model:    "gpt-4-turbo",
			Provider: "openai",
			Enabled:  true,
			Notes:    "OpenAI GPT-4 Turbo",
		},
		"gpt-4": {
			Name:     "gpt-4",
			Model:    "gpt-4",
			Provider: "openai",
			Enabled:  true,
			Notes:    "OpenAI GPT-4",
		},
		"gpt-3.5-turbo": {
			Name:     "gpt-3.5-turbo",
			Model:    "gpt-3.5-turbo",
			Provider: "openai",
			Enabled:  true,
			Notes:    "OpenAI GPT-3.5 Turbo",
		},
		"gpt": {
			Name:     "gpt",
			Model:    "gpt-4o",
			Provider: "openai",
			Enabled:  true,
			Notes:    "Default OpenAI model",
		},
		"mistral": {
			Name:     "mistral",
			Model:    "mistral-large",
			Provider: "mistral",
			Enabled:  true,
			Notes:    "Mistral Large",
		},
		"mistral-small": {
			Name:     "mistral-small",
			Model:    "mistral-small",
			Provider: "mistral",
			Enabled:  true,
			Notes:    "Mistral Small",
		},
		"llama-3.1-70b": {
			Name:     "llama-3.1-70b",
			Model:    "llama-3.1-70b-versatile",
			Provider: "groq",
			Enabled:  true,
			Notes:    "Llama 3.1 70B on Groq",
		},
		"llama-3.1-8b": {
			Name:     "llama-3.1-8b",
			Model:    "llama-3.1-8b-instant",
			Provider: "groq",
			Enabled:  true,
			Notes:    "Llama 3.1 8B on Groq",
		},
		"llama-3": {
			Name:     "llama-3",
			Model:    "llama3-8b-8192",
			Provider: "groq",
			Enabled:  true,
			Notes:    "Llama 3 8B on Groq",
		},
		"llama-2": {
			Name:     "llama-2",
			Model:    "llama2-70b-4096",
			Provider: "groq",
			Enabled:  true,
			Notes:    "Llama 2 70B on Groq",
		},
		"mixtral": {
			Name:     "mixtral",
			Model:    "mixtral-8x7b-32768",
			Provider: "groq",
			Enabled:  true,
			Notes:    "Mixtral 8x7B on Groq",
		},
		"llama": {
			Name:     "llama",
			Model:    "llama3-8b-8192",
			Provider: "groq",
			Enabled:  true,
			Notes:    "Llama 3 on Groq",
		},
		"openai": {
			Name:     "openai",
			Model:    "gpt-4o",
			Provider: "openai",
			Enabled:  true,
			Notes:    "OpenAI default",
		},
		"anthropic": {
			Name:     "anthropic",
			Model:    "claude-3-5-sonnet-20241022",
			Provider: "anthropic",
			Enabled:  true,
			Notes:    "Anthropic default",
		},
		"bedrock": {
			Name:     "bedrock",
			Model:    "anthropic.claude-3-5-sonnet-20241022-v2:0",
			Provider: "aws",
			Enabled:  true,
			Notes:    "AWS Bedrock default",
		},
		"azure": {
			Name:     "azure",
			Model:    "gpt-4o",
			Provider: "azure",
			Enabled:  true,
			Notes:    "Azure OpenAI default",
		},
		"vertex": {
			Name:     "vertex",
			Model:    "gemini-1.5-flash",
			Provider: "google",
			Enabled:  true,
			Notes:    "Google Vertex AI default",
		},
		"openrouter": {
			Name:     "openrouter",
			Model:    "openai/gpt-4o",
			Provider: "openrouter",
			Enabled:  true,
			Notes:    "OpenRouter default",
		},
		"fireworks": {
			Name:     "fireworks",
			Model:    "accounts/fireworks/models/firefunction-v2",
			Provider: "fireworks",
			Enabled:  true,
			Notes:    "Fireworks AI default",
		},
		"together": {
			Name:     "together",
			Model:    "togethercomputer/llama-3-70b-chat",
			Provider: "together",
			Enabled:  true,
			Notes:    "Together AI default",
		},
		"groq": {
			Name:     "groq",
			Model:    "llama3-8b-8192",
			Provider: "groq",
			Enabled:  true,
			Notes:    "Groq default",
		},
		"deepseek-coder": {
			Name:     "deepseek-coder",
			Model:    "deepseek-coder",
			Provider: "deepseek",
			Enabled:  true,
			Notes:    "DeepSeek Coder",
		},
		"together": {
			Name:     "together",
			Model:    "togethercomputer/llama-3-70b-chat",
			Provider: "together",
			Enabled:  true,
			Notes:    "Together AI Chat",
		},
		"fireworks": {
			Name:     "fireworks",
			Model:    "accounts/fireworks/models/firefunction-v2",
			Provider: "fireworks",
			Enabled:  true,
			Notes:    "Fireworks AI Function",
		},
		"cerebras": {
			Name:     "cerebras",
			Model:    "cerebras-llama3.1-8b",
			Provider: "cerebras",
			Enabled:  true,
			Notes:    "Cerebras Llama 3.1 8B",
		},
		"lm-studio": {
			Name:     "lm-studio",
			Model:    "default",
			Provider: "lmstudio",
			Enabled:  true,
			Notes:    "LM Studio default",
		},
		"ollama": {
			Name:     "ollama",
			Model:    "default",
			Provider: "ollama",
			Enabled:  true,
			Notes:    "Ollama default",
		},
	}
}

// saveAliases saves the aliases to file.
func (m *Manager) saveAliases(aliases map[string]*Alias) error {
	aliasesPath := filepath.Join(m.configDir, aliasesFileName)

	data, err := yaml.Marshal(aliases)
	if err != nil {
		return fmt.Errorf("failed to marshal aliases: %w", err)
	}

	if err := os.WriteFile(aliasesPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write aliases file: %w", err)
	}

	return nil
}

// ResolveModel resolves a model/alias to its full specification.
func ResolveModel(alias string) (model, provider string, err error) {
	mgr, err := NewManager()
	if err != nil {
		return "", "", err
	}

	return mgr.ResolveAlias(alias)
}

// SwitchProvider switches the default provider for all aliases.
func SwitchProvider(provider string) error {
	mgr, err := NewManager()
	if err != nil {
		return err
	}

	aliases, err := mgr.GetAliases()
	if err != nil {
		return err
	}

	// Update all aliases to use new provider
	for name, alias := range aliases {
		alias.Provider = provider
		aliases[name] = alias
	}

	return mgr.saveAliases(aliases)
}
