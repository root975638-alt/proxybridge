package anthropic

import (
	"os"
	"strings"
	"testing"
)

func TestAnthropicProvider_Name(t *testing.T) {
	p := NewAnthropicProvider()
	if p.Name() != "Anthropic" {
		t.Errorf("expected name 'Anthropic', got '%s'", p.Name())
	}
}

func TestAnthropicProvider_ID(t *testing.T) {
	p := NewAnthropicProvider()
	if p.ID() != "anthropic" {
		t.Errorf("expected ID 'anthropic', got '%s'", p.ID())
	}
}

func TestAnthropicProvider_GetDefaultModel(t *testing.T) {
	p := NewAnthropicProvider()
	expected := "claude-3-5-sonnet-20241022"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestAnthropicProvider_GetConfig(t *testing.T) {
	p := NewAnthropicProvider()
	config := p.GetConfig()
	if config["provider"] != "anthropic" {
		t.Errorf("expected provider 'anthropic', got '%s'", config["provider"])
	}
	if config["type"] != "anthropic" {
		t.Errorf("expected type 'anthropic', got '%s'", config["type"])
	}
}

func TestAnthropicProvider_GetEnvironment(t *testing.T) {
	p := NewAnthropicProvider()
	env := p.GetEnvironment()
	if env["ANTHROPIC_API_KEY"] != "YOUR_ANTHROPIC_API_KEY" {
		t.Errorf("expected ANTHROPIC_API_KEY placeholder")
	}
	if env["ANTHROPIC_BASE_URL"] != "https://api.anthropic.com/v1" {
		t.Errorf("expected ANTHROPIC_BASE_URL to be set")
	}
}

func TestAnthropicProvider_IsInstalled(t *testing.T) {
	p := NewAnthropicProvider()
	if !p.IsInstalled() {
		t.Error("Anthropic provider should be considered installed")
	}
}

func TestAnthropicProvider_Install(t *testing.T) {
	p := NewAnthropicProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestAnthropicProvider_Validate(t *testing.T) {
	// Test with valid API key format
	os.Setenv("ANTHROPIC_API_KEY", "sk-ant-internal0-test")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	p := NewAnthropicProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestAnthropicProvider_ValidateInvalidKey(t *testing.T) {
	// Test with invalid API key format (doesn't start with sk-ant-)
	os.Setenv("ANTHROPIC_API_KEY", "sk-invalid-key-format")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	p := NewAnthropicProvider()
	if err := p.Validate(); err == nil {
		t.Error("expected validation to fail for invalid key format")
	}
}

func TestAnthropicProvider_GetModels(t *testing.T) {
	p := NewAnthropicProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"claude-3-5-sonnet-20241022", "claude-3-opus-20240229", "claude-3-haiku-20240307"}
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

func TestAnthropicProvider_EnvironmentVariables(t *testing.T) {
	p := NewAnthropicProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"ANTHROPIC_API_KEY",
		"ANTHROPIC_BASE_URL",
		"ANTHROPIC_VERSION",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}

func TestAnthropicProvider_ValidateMissingKey(t *testing.T) {
	// Ensure env var is not set
	os.Unsetenv("ANTHROPIC_API_KEY")

	p := NewAnthropicProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail when API key is missing")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("expected error about missing API key, got: %v", err)
	}
}
