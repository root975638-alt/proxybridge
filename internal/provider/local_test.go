package provider

import (
	"testing"
)

func TestLocalOpenAIProvider_Name(t *testing.T) {
	p := NewLocalOpenAIProvider()
	if p.Name() != "Local / Self-Hosted" {
		t.Errorf("expected name 'Local / Self-Hosted', got '%s'", p.Name())
	}
}

func TestLocalOpenAIProvider_ID(t *testing.T) {
	p := NewLocalOpenAIProvider()
	if p.ID() != "local" {
		t.Errorf("expected ID 'local', got '%s'", p.ID())
	}
}

func TestLocalOpenAIProvider_GetDefaultModel(t *testing.T) {
	p := NewLocalOpenAIProvider()
	expected := "default"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestLocalOpenAIProvider_GetConfig(t *testing.T) {
	p := NewLocalOpenAIProvider()
	config := p.GetConfig()
	if config["provider"] != "local" {
		t.Errorf("expected provider 'local', got '%s'", config["provider"])
	}
	if config["type"] != "openai_compatible" {
		t.Errorf("expected type 'openai_compatible', got '%s'", config["type"])
	}
}

func TestLocalOpenAIProvider_GetEnvironment(t *testing.T) {
	p := NewLocalOpenAIProvider()
	env := p.GetEnvironment()
	if env["OPENAI_BASE_URL"] == "" {
		t.Errorf("expected OPENAI_BASE_URL to be set")
	}
}

func TestLocalOpenAIProvider_IsInstalled(t *testing.T) {
	p := NewLocalOpenAIProvider()
	if !p.IsInstalled() {
		t.Error("Local provider should be considered installed")
	}
}

func TestLocalOpenAIProvider_SetBaseURL(t *testing.T) {
	p := NewLocalOpenAIProvider()
	p.SetBaseURL("http://localhost:8000/v1")
}

func TestLocalOpenAIProvider_SetModelName(t *testing.T) {
	p := NewLocalOpenAIProvider()
	p.SetModelName("mistral-7b")
}

func TestLocalOpenAIProvider_Install(t *testing.T) {
	p := NewLocalOpenAIProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestLocalOpenAIProvider_Uninstall(t *testing.T) {
	p := NewLocalOpenAIProvider()
	if err := p.Uninstall(); err != nil {
		t.Errorf("Uninstall returned error: %v", err)
	}
}

func TestLocalOpenAIProvider_Configuration(t *testing.T) {
	p := NewLocalOpenAIProvider()
	config := p.GetConfig()

	// Verify configuration values
	if config["provider"] != "local" {
		t.Errorf("expected provider 'local'")
	}
	if config["type"] != "openai_compatible" {
		t.Errorf("expected type 'openai_compatible'")
	}
}

func TestLocalOpenAIProvider_EnvironmentVariables(t *testing.T) {
	p := NewLocalOpenAIProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"OPENAI_API_KEY",
		"OPENAI_BASE_URL",
		"LOCAL_MODEL_NAME",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
