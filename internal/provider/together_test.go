package together

import (
	"os"
	"strings"
	"testing"
)

func TestTogetherAIProvider_Name(t *testing.T) {
	p := NewTogetherAIProvider()
	if p.Name() != "Together AI" {
		t.Errorf("expected name 'Together AI', got '%s'", p.Name())
	}
}

func TestTogetherAIProvider_ID(t *testing.T) {
	p := NewTogetherAIProvider()
	if p.ID() != "together" {
		t.Errorf("expected ID 'together', got '%s'", p.ID())
	}
}

func TestTogetherAIProvider_GetDefaultModel(t *testing.T) {
	p := NewTogetherAIProvider()
	expected := "togethercomputer/llama-3-70b-chat"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestTogetherAIProvider_GetConfig(t *testing.T) {
	p := NewTogetherAIProvider()
	config := p.GetConfig()
	if config["provider"] != "together" {
		t.Errorf("expected provider 'together', got '%s'", config["provider"])
	}
	if config["type"] != "together" {
		t.Errorf("expected type 'together', got '%s'", config["type"])
	}
}

func TestTogetherAIProvider_GetEnvironment(t *testing.T) {
	p := NewTogetherAIProvider()
	env := p.GetEnvironment()
	if env["TOGETHER_API_KEY"] != "YOUR_TOGETHER_API_KEY" {
		t.Errorf("expected TOGETHER_API_KEY placeholder")
	}
	if env["TOGETHER_BASE_URL"] != "https://api.together.xyz/v1" {
		t.Errorf("expected TOGETHER_BASE_URL to be set")
	}
}

func TestTogetherAIProvider_IsInstalled(t *testing.T) {
	p := NewTogetherAIProvider()
	if !p.IsInstalled() {
		t.Error("Together AI provider should be considered installed")
	}
}

func TestTogetherAIProvider_Install(t *testing.T) {
	p := NewTogetherAIProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestTogetherAIProvider_Validate(t *testing.T) {
	os.Setenv("TOGETHER_API_KEY", "test-together-key-12345")
	defer os.Unsetenv("TOGETHER_API_KEY")

	p := NewTogetherAIProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestTogetherAIProvider_ValidateMissingKey(t *testing.T) {
	os.Unsetenv("TOGETHER_API_KEY")

	p := NewTogetherAIProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail when API key is missing")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("expected error about missing API key, got: %v", err)
	}
}

func TestTogetherAIProvider_GetModels(t *testing.T) {
	p := NewTogetherAIProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"togethercomputer/llama-3-70b-chat", "togethercomputer/llama-3-8b-chat"}
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

func TestTogetherAIProvider_EnvironmentVariables(t *testing.T) {
	p := NewTogetherAIProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"TOGETHER_API_KEY",
		"TOGETHER_BASE_URL",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
