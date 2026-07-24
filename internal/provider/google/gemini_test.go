package google

import (
	"os"
	"strings"
	"testing"
)

func TestGoogleGeminiProvider_Name(t *testing.T) {
	p := NewGoogleGeminiProvider()
	if p.Name() != "Google Gemini" {
		t.Errorf("expected name 'Google Gemini', got '%s'", p.Name())
	}
}

func TestGoogleGeminiProvider_ID(t *testing.T) {
	p := NewGoogleGeminiProvider()
	if p.ID() != "google" {
		t.Errorf("expected ID 'google', got '%s'", p.ID())
	}
}

func TestGoogleGeminiProvider_GetDefaultModel(t *testing.T) {
	p := NewGoogleGeminiProvider()
	expected := "gemini-1.5-flash"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestGoogleGeminiProvider_GetConfig(t *testing.T) {
	p := NewGoogleGeminiProvider()
	config := p.GetConfig()
	if config["provider"] != "google" {
		t.Errorf("expected provider 'google', got '%s'", config["provider"])
	}
	if config["type"] != "gemini" {
		t.Errorf("expected type 'gemini', got '%s'", config["type"])
	}
}

func TestGoogleGeminiProvider_GetEnvironment(t *testing.T) {
	p := NewGoogleGeminiProvider()
	env := p.GetEnvironment()
	if env["GOOGLE_API_KEY"] != "YOUR_GOOGLE_API_KEY" {
		t.Errorf("expected GOOGLE_API_KEY placeholder")
	}
}

func TestGoogleGeminiProvider_IsInstalled(t *testing.T) {
	p := NewGoogleGeminiProvider()
	if !p.IsInstalled() {
		t.Error("Google Gemini provider should be considered installed")
	}
}

func TestGoogleGeminiProvider_Install(t *testing.T) {
	p := NewGoogleGeminiProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestGoogleGeminiProvider_Validate(t *testing.T) {
	// Test with valid API key format
	os.Setenv("GOOGLE_API_KEY", "AIza-test-invalid-key")
	defer os.Unsetenv("GOOGLE_API_KEY")

	p := NewGoogleGeminiProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestGoogleGeminiProvider_ValidateWithServiceAccount(t *testing.T) {
	// Test with service account credentials
	os.Setenv("GOOGLE_API_KEY", "")
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/path/to/service-account.json")
	defer func() {
		os.Unsetenv("GOOGLE_API_KEY")
		os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
	}()

	p := NewGoogleGeminiProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestGoogleGeminiProvider_ValidateMissingKey(t *testing.T) {
	// Ensure env vars are not set
	os.Unsetenv("GOOGLE_API_KEY")
	os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")

	p := NewGoogleGeminiProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail when API key is missing")
	}
	if !strings.Contains(err.Error(), "not set") &&
		!strings.Contains(err.Error(), "GOOGLE_API_KEY") &&
		!strings.Contains(err.Error(), "GOOGLE_APPLICATION_CREDENTIALS") {
		t.Errorf("expected error about missing credentials, got: %v", err)
	}
}

func TestGoogleGeminiProvider_GetModels(t *testing.T) {
	p := NewGoogleGeminiProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"gemini-1.5-flash", "gemini-1.5-pro", "gemini-2.0-flash"}
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

func TestGoogleGeminiProvider_EnvironmentVariables(t *testing.T) {
	p := NewGoogleGeminiProvider()
	env := p.GetEnvironment()

	expectedVars := []string{
		"GOOGLE_API_KEY",
		"GOOGLE_APPLICATION_CREDENTIALS",
	}

	for _, varName := range expectedVars {
		if _, ok := env[varName]; !ok {
			t.Errorf("expected environment variable '%s' to be present", varName)
		}
	}
}
