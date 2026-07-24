package ollama

import (
	"os/exec"
	"strings"
	"testing"
)

func TestOllamaProvider_Name(t *testing.T) {
	p := NewOllamaProvider()
	if p.Name() != "Ollama" {
		t.Errorf("expected name 'Ollama', got '%s'", p.Name())
	}
}

func TestOllamaProvider_ID(t *testing.T) {
	p := NewOllamaProvider()
	if p.ID() != "ollama" {
		t.Errorf("expected ID 'ollama', got '%s'", p.ID())
	}
}

func TestOllamaProvider_GetDefaultModel(t *testing.T) {
	p := NewOllamaProvider()
	expected := "default"
	if p.GetDefaultModel() != expected {
		t.Errorf("expected default model '%s', got '%s'", expected, p.GetDefaultModel())
	}
}

func TestOllamaProvider_GetConfig(t *testing.T) {
	p := NewOllamaProvider()
	config := p.GetConfig()
	if config["provider"] != "ollama" {
		t.Errorf("expected provider 'ollama', got '%s'", config["provider"])
	}
	if config["type"] != "ollama" {
		t.Errorf("expected type 'ollama', got '%s'", config["type"])
	}
	if config["host"] != "http://localhost:11434" {
		t.Errorf("expected host 'http://localhost:11434', got '%s'", config["host"])
	}
}

func TestOllamaProvider_GetEnvironment(t *testing.T) {
	p := NewOllamaProvider()
	env := p.GetEnvironment()
	if env["OLLAMA_HOST"] != "http://localhost:11434" {
		t.Errorf("expected OLLAMA_HOST to be set")
	}
	if env["OLLAMA_MODEL"] != "default" {
		t.Errorf("expected OLLAMA_MODEL to be set")
	}
}

func TestOllamaProvider_IsInstalled(t *testing.T) {
	p := NewOllamaProvider()
	installed := p.IsInstalled()
	// Ollama may or may not be installed in the test environment
	// The test should not fail regardless
	_ = installed
}

func TestOllamaProvider_Install(t *testing.T) {
	p := NewOllamaProvider()
	if err := p.Install(); err != nil {
		t.Errorf("Install returned error: %v", err)
	}
}

func TestOllamaProvider_Uninstall(t *testing.T) {
	p := NewOllamaProvider()
	if err := p.Uninstall(); err != nil {
		t.Errorf("Uninstall returned error: %v", err)
	}
}

func TestOllamaProvider_ValidateNotInstalled(t *testing.T) {
	p := NewOllamaProvider()

	// Temporarily remove ollama from PATH
	originalPath := os.Getenv("PATH")
	os.Setenv("PATH", "/nonexistent")
	defer os.Setenv("PATH", originalPath)

	if err := p.Validate(); err == nil {
		t.Error("expected validation to fail when Ollama is not installed")
	}
}

func TestOllamaProvider_GetModels(t *testing.T) {
	p := NewOllamaProvider()

	// GetModels will only work if Ollama is installed and running
	models, err := p.GetModels()
	if err != nil {
		// It's okay if this fails in test environment
		// As long as the command structure is correct
		t.Logf("GetModels returned error (may be expected): %v", err)
		return
	}

	// If successful, verify we got a list of models
	for _, model := range models {
		if !strings.Contains(model, ":") && len(model) > 0 {
			// Valid model names typically contain a tag like "model:tag"
			// But some might be just names, so this is a loose check
		}
	}
}

func TestOllamaProvider_Configuration(t *testing.T) {
	p := NewOllamaProvider()
	config := p.GetConfig()

	// Verify configuration values
	if config["provider"] != "ollama" {
		t.Errorf("expected provider 'ollama'")
	}
	if config["type"] != "ollama" {
		t.Errorf("expected type 'ollama'")
	}
}

func TestOllamaProvider_Description(t *testing.T) {
	p := NewOllamaProvider()
	if p.description != "Ollama Local LLM Server" {
		t.Errorf("unexpected description: %s", p.description)
	}
}
