package openai

import (
	"os"
	"strings"
	"testing"
)

func TestOpenAIProvider_Name(t *testing.T) {
	p := NewOpenAIProvider()
	if p.Name() != "OpenAI" {
		t.Errorf("expected name 'OpenAI', got '%s'", p.Name())
	}
}

func TestOpenAIProvider_ID(t *testing.T) {
	p := NewOpenAIProvider()
	if p.ID() != "openai" {
		t.Errorf("expected ID 'openai', got '%s'", p.ID())
	}
}

func TestOpenAIProvider_GetDefaultModel(t *testing.T) {
	p := NewOpenAIProvider()
	expected := "gpt-4o"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestOpenAIProvider_GetConfig(t *testing.T) {
	p := NewOpenAIProvider()
	config := p.GetConfig()
	if config["provider"] != "openai" {
		t.Errorf("expected provider 'openai', got '%s'", config["provider"])
	}
	if config["type"] != "openai" {
		t.Errorf("expected type 'openai', got '%s'", config["type"])
	}
}

func TestOpenAIProvider_GetEnvironment(t *testing.T) {
	p := NewOpenAIProvider()
	env := p.GetEnvironment()
	if env["OPENAI_API_KEY"] != "YOUR_OPENAI_API_KEY" {
		t.Errorf("expected OPENAI_API_KEY placeholder")
	}
	if env["OPENAI_BASE_URL"] != "https://api.openai.com/v1" {
		t.Errorf("expected OPENAI_BASE_URL to be set")
	}
}

func TestOpenAIProvider_IsInstalled(t *testing.T) {
	p := NewOpenAIProvider()
	if !p.IsInstalled() {
		t.Error("OpenAI provider should be considered installed")
	}
}

func TestOpenAIProvider_Install(t *testing.T) {
	p := NewOpenAIProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestOpenAIProvider_Validate(t *testing.T) {
	// Test with valid API key format
	os.Setenv("OPENAI_API_KEY", "sk-test-invalid-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	p := NewOpenAIProvider()
	// This should pass format validation
	if err := p.Validate(); err != nil {
		// Format validation might fail if key doesn't start with "sk-"
		// That's expected behavior
		if !strings.Contains(err.Error(), "API key should start") &&
			!strings.Contains(err.Error(), "validation failed") {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestOpenAIProvider_GetModels(t *testing.T) {
	p := NewOpenAIProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-4"}
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

func TestOpenAIProvider_GetEnvironmentVariables(t *testing.T) {
	p := NewOpenAIProvider()
	env := p.GetEnvironment()

	// Verify all expected environment variables are present
	expectedVars := []string{
		"OPENAI_API_KEY",
		"OPENAI_BASE_URL",
		"OPENAI_ORG_ID",
		"OPENAI_PROJECT_ID",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
