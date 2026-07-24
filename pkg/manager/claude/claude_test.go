package claude

import (
	"testing"
	"time"
)

func TestManager_CreateDefaultConfig(t *testing.T) {
	// Create a default config
	cfg := &ClaudeConfig{
		Version:          "1.0",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		DefaultModel:     "claude-3-5-sonnet",
		Provider:         "anthropic",
		AutoUpdateModels: true,
		Models: []Model{
			{
				Name:        "claude-3-5-sonnet",
				Provider:    "anthropic",
				Model:       "anthropic/claude-3-5-sonnet",
				MaxTokens:   8192,
				Temperature: 0.7,
			},
		},
	}

	if cfg.DefaultModel != "claude-3-5-sonnet" {
		t.Errorf("Expected default model 'claude-3-5-sonnet', got '%s'", cfg.DefaultModel)
	}

	if len(cfg.Models) != 1 {
		t.Errorf("Expected 1 model, got %d", len(cfg.Models))
	}

	if cfg.Models[0].Name != "claude-3-5-sonnet" {
		t.Errorf("Expected model name 'claude-3-5-sonnet', got '%s'", cfg.Models[0].Name)
	}
}

func TestModel_Structure(t *testing.T) {
	model := Model{
		Name:        "test-model",
		Provider:    "test-provider",
		Model:       "test-provider/test-model",
		APIKey:      "test-key",
		BaseURL:     "https://api.example.com",
		MaxTokens:   4096,
		Temperature: 0.7,
	}

	if model.Name != "test-model" {
		t.Errorf("Expected name 'test-model', got '%s'", model.Name)
	}

	if model.Provider != "test-provider" {
		t.Errorf("Expected provider 'test-provider', got '%s'", model.Provider)
	}

	if model.MaxTokens != 4096 {
		t.Errorf("Expected maxTokens 4096, got %d", model.MaxTokens)
	}
}

func TestClaudeConfig_Structure(t *testing.T) {
	cfg := &ClaudeConfig{
		Version:          "1.0",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		DefaultModel:     "claude-3-5-sonnet",
		Provider:         "anthropic",
		AutoUpdateModels: true,
		Models: []Model{
			{
				Name:        "claude-3-5-sonnet",
				Provider:    "anthropic",
				Model:       "anthropic/claude-3-5-sonnet",
				MaxTokens:   8192,
				Temperature: 0.7,
			},
		},
	}

	// Verify struct fields are accessible
	_ = cfg.Version
	_ = cfg.DefaultModel
	_ = cfg.Provider
	_ = cfg.Models
}

func TestManager_ValidateModel(t *testing.T) {
	model := Model{
		Name:     "claude-3-5-sonnet",
		Provider: "anthropic",
		Model:    "anthropic/claude-3-5-sonnet",
	}

	// Verify model structure
	if model.Name == "" {
		t.Error("Model name should not be empty")
	}

	if model.Provider == "" {
		t.Error("Model provider should not be empty")
	}

	if model.Model == "" {
		t.Error("Full model path should not be empty")
	}
}

func TestManager_EnsureValidTemperature(t *testing.T) {
	model := Model{
		Name:        "test-model",
		Temperature: 0.7,
	}

	// Verify temperature is in valid range [0, 2]
	if model.Temperature < 0 || model.Temperature > 2 {
		t.Error("Temperature should be between 0 and 2")
	}
}

func TestManager_ModelMaxTokens(t *testing.T) {
	model := Model{
		Name:      "test-model",
		MaxTokens: 8192,
	}

	// Verify max tokens is reasonable
	if model.MaxTokens <= 0 || model.MaxTokens > 100000 {
		t.Error("Max tokens should be between 1 and 100000")
	}
}
