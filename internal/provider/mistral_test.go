package provider

import (
	"os"
	"strings"
	"testing"
)

func TestMistralProvider_Name(t *testing.T) {
	p := NewMistralProvider()
	if p.Name() != "Mistral AI" {
		t.Errorf("expected name 'Mistral AI', got '%s'", p.Name())
	}
}

func TestMistralProvider_ID(t *testing.T) {
	p := NewMistralProvider()
	if p.ID() != "mistral" {
		t.Errorf("expected ID 'mistral', got '%s'", p.ID())
	}
}

func TestMistralProvider_GetDefaultModel(t *testing.T) {
	p := NewMistralProvider()
	expected := "mistral-large"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestMistralProvider_GetConfig(t *testing.T) {
	p := NewMistralProvider()
	config := p.GetConfig()
	if config["provider"] != "mistral" {
		t.Errorf("expected provider 'mistral', got '%s'", config["provider"])
	}
	if config["type"] != "mistral" {
		t.Errorf("expected type 'mistral', got '%s'", config["type"])
	}
}

func TestMistralProvider_GetEnvironment(t *testing.T) {
	p := NewMistralProvider()
	env := p.GetEnvironment()
	if env["MISTRAL_API_KEY"] != "YOUR_MISTRAL_API_KEY" {
		t.Errorf("expected MISTRAL_API_KEY placeholder")
	}
	if env["MISTRAL_BASE_URL"] != "https://api.mistral.ai/v1" {
		t.Errorf("expected MISTRAL_BASE_URL to be set")
	}
}

func TestMistralProvider_IsInstalled(t *testing.T) {
	p := NewMistralProvider()
	if !p.IsInstalled() {
		t.Error("Mistral provider should be considered installed")
	}
}

func TestMistralProvider_Install(t *testing.T) {
	p := NewMistralProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestMistralProvider_Validate(t *testing.T) {
	os.Setenv("MISTRAL_API_KEY", "test-mistral-key-12345")
	defer os.Unsetenv("MISTRAL_API_KEY")

	p := NewMistralProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestMistralProvider_ValidateMissingKey(t *testing.T) {
	os.Unsetenv("MISTRAL_API_KEY")

	p := NewMistralProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail when API key is missing")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("expected error about missing API key, got: %v", err)
	}
}

func TestMistralProvider_GetModels(t *testing.T) {
	p := NewMistralProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"mistral-large", "mistral-medium", "mistral-small"}
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

func TestMistralProvider_EnvironmentVariables(t *testing.T) {
	p := NewMistralProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"MISTRAL_API_KEY",
		"MISTRAL_BASE_URL",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
