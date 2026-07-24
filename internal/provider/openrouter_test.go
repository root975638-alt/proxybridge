package openrouter

import (
	"os"
	"strings"
	"testing"
)

func TestOpenRouterProvider_Name(t *testing.T) {
	p := NewOpenRouterProvider()
	if p.Name() != "OpenRouter" {
		t.Errorf("expected name 'OpenRouter', got '%s'", p.Name())
	}
}

func TestOpenRouterProvider_ID(t *testing.T) {
	p := NewOpenRouterProvider()
	if p.ID() != "openrouter" {
		t.Errorf("expected ID 'openrouter', got '%s'", p.ID())
	}
}

func TestOpenRouterProvider_GetDefaultModel(t *testing.T) {
	p := NewOpenRouterProvider()
	expected := "openai/gpt-4o"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestOpenRouterProvider_GetConfig(t *testing.T) {
	p := NewOpenRouterProvider()
	config := p.GetConfig()
	if config["provider"] != "openrouter" {
		t.Errorf("expected provider 'openrouter', got '%s'", config["provider"])
	}
	if config["type"] != "openrouter" {
		t.Errorf("expected type 'openrouter', got '%s'", config["type"])
	}
}

func TestOpenRouterProvider_GetEnvironment(t *testing.T) {
	p := NewOpenRouterProvider()
	env := p.GetEnvironment()
	if env["OPENROUTER_API_KEY"] != "YOUR_OPENROUTER_API_KEY" {
		t.Errorf("expected OPENROUTER_API_KEY placeholder")
	}
	if env["OPENROUTER_BASE_URL"] != "https://openrouter.ai/api/v1" {
		t.Errorf("expected OPENROUTER_BASE_URL to be set")
	}
}

func TestOpenRouterProvider_IsInstalled(t *testing.T) {
	p := NewOpenRouterProvider()
	if !p.IsInstalled() {
		t.Error("OpenRouter provider should be considered installed")
	}
}

func TestOpenRouterProvider_Install(t *testing.T) {
	p := NewOpenRouterProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestOpenRouterProvider_Validate(t *testing.T) {
	os.Setenv("OPENROUTER_API_KEY", "test-openrouter-key-12345")
	defer os.Unsetenv("OPENROUTER_API_KEY")

	p := NewOpenRouterProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestOpenRouterProvider_ValidateMissingKey(t *testing.T) {
	os.Unsetenv("OPENROUTER_API_KEY")

	p := NewOpenRouterProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail when API key is missing")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("expected error about missing API key, got: %v", err)
	}
}

func TestOpenRouterProvider_GetModels(t *testing.T) {
	p := NewOpenRouterProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"openai/gpt-4o", "anthropic/claude-3-5-sonnet"}
	for _, expected := range expectedModels {
		found := false
		for _, model := range models {
			if model == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected model '%s' in list", expected)
		}
	}
}

func TestOpenRouterProvider_EnvironmentVariables(t *testing.T) {
	p := NewOpenRouterProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"OPENROUTER_API_KEY",
		"OPENROUTER_BASE_URL",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
