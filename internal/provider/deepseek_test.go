package deepseek

import (
	"os"
	"strings"
	"testing"
)

func TestDeepSeekProvider_Name(t *testing.T) {
	p := NewDeepSeekProvider()
	if p.Name() != "DeepSeek" {
		t.Errorf("expected name 'DeepSeek', got '%s'", p.Name())
	}
}

func TestDeepSeekProvider_ID(t *testing.T) {
	p := NewDeepSeekProvider()
	if p.ID() != "deepseek" {
		t.Errorf("expected ID 'deepseek', got '%s'", p.ID())
	}
}

func TestDeepSeekProvider_GetDefaultModel(t *testing.T) {
	p := NewDeepSeekProvider()
	expected := "deepseek-chat"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestDeepSeekProvider_GetConfig(t *testing.T) {
	p := NewDeepSeekProvider()
	config := p.GetConfig()
	if config["provider"] != "deepseek" {
		t.Errorf("expected provider 'deepseek', got '%s'", config["provider"])
	}
	if config["type"] != "deepseek" {
		t.Errorf("expected type 'deepseek', got '%s'", config["type"])
	}
}

func TestDeepSeekProvider_GetEnvironment(t *testing.T) {
	p := NewDeepSeekProvider()
	env := p.GetEnvironment()
	if env["DEEPSEEK_API_KEY"] != "YOUR_DEEPSEEK_API_KEY" {
		t.Errorf("expected DEEPSEEK_API_KEY placeholder")
	}
	if env["DEEPSEEK_BASE_URL"] != "https://api.deepseek.com/v1" {
		t.Errorf("expected DEEPSEEK_BASE_URL to be set")
	}
}

func TestDeepSeekProvider_IsInstalled(t *testing.T) {
	p := NewDeepSeekProvider()
	if !p.IsInstalled() {
		t.Error("DeepSeek provider should be considered installed")
	}
}

func TestDeepSeekProvider_Install(t *testing.T) {
	p := NewDeepSeekProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestDeepSeekProvider_Validate(t *testing.T) {
	os.Setenv("DEEPSEEK_API_KEY", "test-deepseek-key-12345")
	defer os.Unsetenv("DEEPSEEK_API_KEY")

	p := NewDeepSeekProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestDeepSeekProvider_ValidateMissingKey(t *testing.T) {
	os.Unsetenv("DEEPSEEK_API_KEY")

	p := NewDeepSeekProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail when API key is missing")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("expected error about missing API key, got: %v", err)
	}
}

func TestDeepSeekProvider_GetModels(t *testing.T) {
	p := NewDeepSeekProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"deepseek-chat", "deepseek-coder"}
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

func TestDeepSeekProvider_EnvironmentVariables(t *testing.T) {
	p := NewDeepSeekProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"DEEPSEEK_API_KEY",
		"DEEPSEEK_BASE_URL",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
