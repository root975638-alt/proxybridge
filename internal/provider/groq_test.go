package provider

import (
	"os"
	"strings"
	"testing"
)

func TestGroqProvider_Name(t *testing.T) {
	p := NewGroqProvider()
	if p.Name() != "Groq" {
		t.Errorf("expected name 'Groq', got '%s'", p.Name())
	}
}

func TestGroqProvider_ID(t *testing.T) {
	p := NewGroqProvider()
	if p.ID() != "groq" {
		t.Errorf("expected ID 'groq', got '%s'", p.ID())
	}
}

func TestGroqProvider_GetDefaultModel(t *testing.T) {
	p := NewGroqProvider()
	expected := "llama3-8b-8192"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestGroqProvider_GetConfig(t *testing.T) {
	p := NewGroqProvider()
	config := p.GetConfig()
	if config["provider"] != "groq" {
		t.Errorf("expected provider 'groq', got '%s'", config["provider"])
	}
	if config["type"] != "groq" {
		t.Errorf("expected type 'groq', got '%s'", config["type"])
	}
}

func TestGroqProvider_GetEnvironment(t *testing.T) {
	p := NewGroqProvider()
	env := p.GetEnvironment()
	if env["GROQ_API_KEY"] != "YOUR_GROQ_API_KEY" {
		t.Errorf("expected GROQ_API_KEY placeholder")
	}
	if env["GROQ_BASE_URL"] != "https://api.groq.com/openai/v1" {
		t.Errorf("expected GROQ_BASE_URL to be set")
	}
}

func TestGroqProvider_IsInstalled(t *testing.T) {
	p := NewGroqProvider()
	if !p.IsInstalled() {
		t.Error("Groq provider should be considered installed")
	}
}

func TestGroqProvider_Install(t *testing.T) {
	p := NewGroqProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestGroqProvider_Validate(t *testing.T) {
	os.Setenv("GROQ_API_KEY", "test-groq-key-12345")
	defer os.Unsetenv("GROQ_API_KEY")

	p := NewGroqProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestGroqProvider_ValidateMissingKey(t *testing.T) {
	os.Unsetenv("GROQ_API_KEY")

	p := NewGroqProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail when API key is missing")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("expected error about missing API key, got: %v", err)
	}
}

func TestGroqProvider_GetModels(t *testing.T) {
	p := NewGroqProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"llama3-8b-8192", "llama3-70b-8192", "mixtral-8x7b-32768"}
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

func TestGroqProvider_EnvironmentVariables(t *testing.T) {
	p := NewGroqProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"GROQ_API_KEY",
		"GROQ_BASE_URL",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
