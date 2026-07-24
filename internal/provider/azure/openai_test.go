package azure

import (
	"os"
	"strings"
	"testing"
)

func TestAzureOpenAIProvider_Name(t *testing.T) {
	p := NewAzureOpenAIProvider()
	if p.Name() != "Azure OpenAI" {
		t.Errorf("expected name 'Azure OpenAI', got '%s'", p.Name())
	}
}

func TestAzureOpenAIProvider_ID(t *testing.T) {
	p := NewAzureOpenAIProvider()
	if p.ID() != "azure" {
		t.Errorf("expected ID 'azure', got '%s'", p.ID())
	}
}

func TestAzureOpenAIProvider_GetDefaultModel(t *testing.T) {
	p := NewAzureOpenAIProvider()
	expected := "gpt-4o"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestAzureOpenAIProvider_GetConfig(t *testing.T) {
	p := NewAzureOpenAIProvider()
	config := p.GetConfig()
	if config["provider"] != "azure" {
		t.Errorf("expected provider 'azure', got '%s'", config["provider"])
	}
	if config["type"] != "azure_openai" {
		t.Errorf("expected type 'azure_openai', got '%s'", config["type"])
	}
}

func TestAzureOpenAIProvider_GetEnvironment(t *testing.T) {
	p := NewAzureOpenAIProvider()
	env := p.GetEnvironment()
	if env["AZURE_OPENAI_API_KEY"] != "YOUR_AZURE_OPENAI_API_KEY" {
		t.Errorf("expected AZURE_OPENAI_API_KEY placeholder")
	}
	if env["AZURE_OPENAI_ENDPOINT"] == "" {
		t.Errorf("expected AZURE_OPENAI_ENDPOINT to be set")
	}
}

func TestAzureOpenAIProvider_IsInstalled(t *testing.T) {
	p := NewAzureOpenAIProvider()
	if !p.IsInstalled() {
		t.Error("Azure OpenAI provider should be considered installed")
	}
}

func TestAzureOpenAIProvider_Install(t *testing.T) {
	p := NewAzureOpenAIProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestAzureOpenAIProvider_Validate(t *testing.T) {
	os.Setenv("AZURE_OPENAI_API_KEY", "test-key-12345")
	os.Setenv("AZURE_OPENAI_ENDPOINT", "https://test.openai.azure.com/")
	defer func() {
		os.Unsetenv("AZURE_OPENAI_API_KEY")
		os.Unsetenv("AZURE_OPENAI_ENDPOINT")
	}()

	p := NewAzureOpenAIProvider()
	if err := p.Validate(); err != nil {
		t.Errorf("Validate returned error: %v", err)
	}
}

func TestAzureOpenAIProvider_ValidateInvalidEndpoint(t *testing.T) {
	os.Setenv("AZURE_OPENAI_API_KEY", "test-key-12345")
	os.Setenv("AZURE_OPENAI_ENDPOINT", "http://invalid-endpoint")
	defer func() {
		os.Unsetenv("AZURE_OPENAI_API_KEY")
		os.Unsetenv("AZURE_OPENAI_ENDPOINT")
	}()

	p := NewAzureOpenAIProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail for http endpoint")
	}
	if !strings.Contains(err.Error(), "https://") {
		t.Errorf("expected error about https requirement, got: %v", err)
	}
}

func TestAzureOpenAIProvider_ValidateMissingKey(t *testing.T) {
	os.Unsetenv("AZURE_OPENAI_API_KEY")
	os.Setenv("AZURE_OPENAI_ENDPOINT", "https://test.openai.azure.com/")
	defer os.Unsetenv("AZURE_OPENAI_ENDPOINT")

	p := NewAzureOpenAIProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail for missing API key")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("expected error about missing API key, got: %v", err)
	}
}

func TestAzureOpenAIProvider_ValidateMissingEndpoint(t *testing.T) {
	os.Setenv("AZURE_OPENAI_API_KEY", "test-key-12345")
	os.Unsetenv("AZURE_OPENAI_ENDPOINT")
	defer os.Unsetenv("AZURE_OPENAI_API_KEY")

	p := NewAzureOpenAIProvider()
	err := p.Validate()
	if err == nil {
		t.Error("expected validation to fail for missing endpoint")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("expected error about missing endpoint, got: %v", err)
	}
}

func TestAzureOpenAIProvider_GetModels(t *testing.T) {
	p := NewAzureOpenAIProvider()
	models, err := p.GetModels()
	if err != nil {
		t.Errorf("GetModels returned error: %v", err)
	}

	expectedModels := []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo"}
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
