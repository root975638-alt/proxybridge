package provider

import (
	"os"
	"strings"
	"testing"
)

func TestXAIProvider_Name(t *testing.T) {
	p := NewXAIProvider()
	if p.Name() != "xAI (Grok)" {
		t.Errorf("expected name 'xAI (Grok)', got '%s'", p.Name())
	}
}

func TestXAIProvider_ID(t *testing.T) {
	p := NewXAIProvider()
	if p.ID() != "xai" {
		t.Errorf("expected ID 'xai', got '%s'", p.ID())
	}
}

func TestXAIProvider_GetDefaultModel(t *testing.T) {
	p := NewXAIProvider()
	expected := "grok-beta"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestXAIProvider_GetConfig(t *testing.T) {
	p := NewXAIProvider()
	config := p.GetConfig()
	if config["provider"] != "xai" {
		t.Errorf("expected provider 'xai', got '%s'", config["provider"])
	}
	if config["type"] != "xai" {
		t.Errorf("expected type 'xai', got '%s'", config["type"])
	}
}

func TestXAIProvider_GetEnvironment(t *testing.T) {
	p := NewXAIProvider()
	env := p.GetEnvironment()
	if env["XAI_API_KEY"] != "YOUR_XAI_API_KEY" {
		t.Errorf("expected XAI_API_KEY placeholder")
	}
	if env["XAI_BASE_URL"] != "https://api.x.ai/v1" {
		t.Errorf("expected XAI_BASE_URL to be set")
	}
}

func TestXAIProvider_IsInstalled(t *testing.T) {
	p := NewXAIProvider()
	if !p.IsInstalled() {
		t.Error("xAI provider should be considered installed")
	}
}

func TestXAIProvider_Install(t *testing.T) {
	p := NewXAIProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestXAIProvider_Validate(t *testing.T) {
	os.Setenv("XAI_API_KEY", "test-xai-key-12345")
	defer os.Unsetenv("XAI_API_KEY")

	p := NewXAIProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestXAIProvider_ValidateMissingKey(t *testing.T) {
	os.Unsetenv("XAI_API_KEY")

	p := NewXAIProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail when API key is missing")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("expected error about missing API key, got: %v", err)
	}
}

func TestXAIProvider_GetModels(t *testing.T) {
	p := NewXAIProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"grok-beta", "grok-2-1212"}
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

func TestXAIProvider_EnvironmentVariables(t *testing.T) {
	p := NewXAIProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"XAI_API_KEY",
		"XAI_BASE_URL",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
