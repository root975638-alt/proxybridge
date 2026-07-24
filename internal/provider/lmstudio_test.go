package provider

import (
	"testing"
)

func TestLMStudioProvider_Name(t *testing.T) {
	p := NewLMStudioProvider()
	if p.Name() != "LM Studio" {
		t.Errorf("expected name 'LM Studio', got '%s'", p.Name())
	}
}

func TestLMStudioProvider_ID(t *testing.T) {
	p := NewLMStudioProvider()
	if p.ID() != "lmstudio" {
		t.Errorf("expected ID 'lmstudio', got '%s'", p.ID())
	}
}

func TestLMStudioProvider_GetDefaultModel(t *testing.T) {
	p := NewLMStudioProvider()
	expected := "default"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestLMStudioProvider_GetConfig(t *testing.T) {
	p := NewLMStudioProvider()
	config := p.GetConfig()
	if config["provider"] != "lmstudio" {
		t.Errorf("expected provider 'lmstudio', got '%s'", config["provider"])
	}
	if config["type"] != "openai_compatible" {
		t.Errorf("expected type 'openai_compatible', got '%s'", config["type"])
	}
}

func TestLMStudioProvider_GetEnvironment(t *testing.T) {
	p := NewLMStudioProvider()
	env := p.GetEnvironment()
	if env["OPENAI_BASE_URL"] != "http://localhost:1234/v1" {
		t.Errorf("expected OPENAI_BASE_URL to be set")
	}
}

func TestLMStudioProvider_IsInstalled(t *testing.T) {
	p := NewLMStudioProvider()
	if !p.IsInstalled() {
		t.Error("LM Studio provider should be considered installed")
	}
}

func TestLMStudioProvider_Install(t *testing.T) {
	p := NewLMStudioProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestLMStudioProvider_Uninstall(t *testing.T) {
	p := NewLMStudioProvider()
	if err := p.Uninstall(); err != nil {
		t.Errorf("Uninstall returned error: %v", err)
	}
}

func TestLMStudioProvider_GetModels(t *testing.T) {
	p := NewLMStudioProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	if len(models) == 0 {
		t.Error("expected at least one model")
	}
}

func TestLMStudioProvider_EnvironmentVariables(t *testing.T) {
	p := NewLMStudioProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"OPENAI_API_KEY",
		"OPENAI_BASE_URL",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
