package provider

import (
	"os"
	"strings"
	"testing"
)

func TestCerebrasProvider_Name(t *testing.T) {
	p := NewCerebrasProvider()
	if p.Name() != "Cerebras" {
		t.Errorf("expected name 'Cerebras', got '%s'", p.Name())
	}
}

func TestCerebrasProvider_ID(t *testing.T) {
	p := NewCerebrasProvider()
	if p.ID() != "cerebras" {
		t.Errorf("expected ID 'cerebras', got '%s'", p.ID())
	}
}

func TestCerebrasProvider_GetDefaultModel(t *testing.T) {
	p := NewCerebrasProvider()
	expected := "cerebras-llama3.1-8b"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestCerebrasProvider_GetConfig(t *testing.T) {
	p := NewCerebrasProvider()
	config := p.GetConfig()
	if config["provider"] != "cerebras" {
		t.Errorf("expected provider 'cerebras', got '%s'", config["provider"])
	}
	if config["type"] != "cerebras" {
		t.Errorf("expected type 'cerebras', got '%s'", config["type"])
	}
}

func TestCerebrasProvider_GetEnvironment(t *testing.T) {
	p := NewCerebrasProvider()
	env := p.GetEnvironment()
	if env["CEREBRAS_API_KEY"] != "YOUR_CEREBRAS_API_KEY" {
		t.Errorf("expected CEREBRAS_API_KEY placeholder")
	}
	if env["CEREBRAS_BASE_URL"] != "https://api.cerebras.ai/v1" {
		t.Errorf("expected CEREBRAS_BASE_URL to be set")
	}
}

func TestCerebrasProvider_IsInstalled(t *testing.T) {
	p := NewCerebrasProvider()
	if !p.IsInstalled() {
		t.Error("Cerebras provider should be considered installed")
	}
}

func TestCerebrasProvider_Install(t *testing.T) {
	p := NewCerebrasProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestCerebrasProvider_Validate(t *testing.T) {
	os.Setenv("CEREBRAS_API_KEY", "test-cerebras-key-12345")
	defer os.Unsetenv("CEREBRAS_API_KEY")

	p := NewCerebrasProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestCerebrasProvider_ValidateMissingKey(t *testing.T) {
	os.Unsetenv("CEREBRAS_API_KEY")

	p := NewCerebrasProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail when API key is missing")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("expected error about missing API key, got: %v", err)
	}
}

func TestCerebrasProvider_GetModels(t *testing.T) {
	p := NewCerebrasProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"cerebras-llama3.1-8b", "cerebras-llama3.1-70b"}
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

func TestCerebrasProvider_EnvironmentVariables(t *testing.T) {
	p := NewCerebrasProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"CEREBRAS_API_KEY",
		"CEREBRAS_BASE_URL",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
