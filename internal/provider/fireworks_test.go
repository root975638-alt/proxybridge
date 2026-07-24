package provider

import (
	"os"
	"strings"
	"testing"
)

func TestFireworksProvider_Name(t *testing.T) {
	p := NewFireworksProvider()
	if p.Name() != "Fireworks AI" {
		t.Errorf("expected name 'Fireworks AI', got '%s'", p.Name())
	}
}

func TestFireworksProvider_ID(t *testing.T) {
	p := NewFireworksProvider()
	if p.ID() != "fireworks" {
		t.Errorf("expected ID 'fireworks', got '%s'", p.ID())
	}
}

func TestFireworksProvider_GetDefaultModel(t *testing.T) {
	p := NewFireworksProvider()
	expected := "accounts/fireworks/models/firefunction-v2"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestFireworksProvider_GetConfig(t *testing.T) {
	p := NewFireworksProvider()
	config := p.GetConfig()
	if config["provider"] != "fireworks" {
		t.Errorf("expected provider 'fireworks', got '%s'", config["provider"])
	}
	if config["type"] != "fireworks" {
		t.Errorf("expected type 'fireworks', got '%s'", config["type"])
	}
}

func TestFireworksProvider_GetEnvironment(t *testing.T) {
	p := NewFireworksProvider()
	env := p.GetEnvironment()
	if env["FIREWORKS_API_KEY"] != "YOUR_FIREWORKS_API_KEY" {
		t.Errorf("expected FIREWORKS_API_KEY placeholder")
	}
	if env["FIREWORKS_BASE_URL"] != "https://api.fireworks.ai/inference/v1" {
		t.Errorf("expected FIREWORKS_BASE_URL to be set")
	}
}

func TestFireworksProvider_IsInstalled(t *testing.T) {
	p := NewFireworksProvider()
	if !p.IsInstalled() {
		t.Error("Fireworks AI provider should be considered installed")
	}
}

func TestFireworksProvider_Install(t *testing.T) {
	p := NewFireworksProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestFireworksProvider_Validate(t *testing.T) {
	os.Setenv("FIREWORKS_API_KEY", "test-fireworks-key-12345")
	defer os.Unsetenv("FIREWORKS_API_KEY")

	p := NewFireworksProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestFireworksProvider_ValidateMissingKey(t *testing.T) {
	os.Unsetenv("FIREWORKS_API_KEY")

	p := NewFireworksProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail when API key is missing")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("expected error about missing API key, got: %v", err)
	}
}

func TestFireworksProvider_GetModels(t *testing.T) {
	p := NewFireworksProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"accounts/fireworks/models/firefunction-v2", "accounts/fireworks/models/llama-v3-70b-instruct"}
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

func TestFireworksProvider_EnvironmentVariables(t *testing.T) {
	p := NewFireworksProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"FIREWORKS_API_KEY",
		"FIREWORKS_BASE_URL",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
