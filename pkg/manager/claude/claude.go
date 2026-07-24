// Package claude provides Claude Code configuration management.
package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/proxybridge/cli/internal/config"
	"github.com/proxybridge/cli/internal/logging"
)

// Manager manages Claude Code configuration.
type Manager struct {
	config *config.Config
	logger *logging.Logger
}

// ClaudeConfig represents Claude Code configuration.
type ClaudeConfig struct {
	Version           string    `json:"version"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Models            []Model   `json:"models"`
	DefaultModel      string    `json:"defaultModel"`
	Provider          string    `json:"provider"`
	ProxyURL          string    `json:"proxyUrl,omitempty"`
	ProxyEnabled      bool      `json:"proxyEnabled,omitempty"`
	AutoUpdateModels  bool      `json:"autoUpdateModels,omitempty"`
}

// Model represents a model configuration.
type Model struct {
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	Model       string `json:"model"`
	APIKey      string `json:"apiKey,omitempty"`
	BaseURL     string `json:"baseUrl,omitempty"`
	MaxTokens   int    `json:"maxTokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

// NewManager creates a new Claude Code manager.
func NewManager(cfg *config.Config, logger *logging.Logger) *Manager {
	return &Manager{
		config: cfg,
		logger: logger,
	}
}

// Detect detects Claude Code installation.
func (m *Manager) Detect() (bool, string) {
	// Check common Claude Code paths
	paths := []string{
		"/usr/local/bin/claude",
		"/opt/claude",
		filepath.Join(os.Getenv("HOME"), ".local/bin/claude"),
		filepath.Join(os.Getenv("HOME"), "bin/claude"),
	}

	for _, path := range paths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return true, path
		}
	}

	// Check PATH
	if path, err := exec.LookPath("claude"); err == nil {
		return true, path
	}

	// Check macOS application
	if isMacOS() {
		appPath := filepath.Join(os.Getenv("HOME"), "Applications/Claude.app")
		if info, err := os.Stat(appPath); err == nil && info.IsDir() {
			return true, appPath
		}
	}

	// Check Windows paths
	if isWindows() {
		paths := []string{
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs/claude/claude.exe"),
			filepath.Join(os.Getenv("ProgramFiles"), "claude/claude.exe"),
		}
		for _, path := range paths {
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				return true, path
			}
		}
	}

	return false, ""
}

// GetConfigPath returns the Claude Code config path.
func (m *Manager) GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch {
	case isWindows():
		return filepath.Join(home, "AppData", "Roaming", "Claude", "claude_config.json"), nil
	case isMacOS():
		return filepath.Join(home, "Library", "Application Support", "Claude", "claude_config.json"), nil
	default:
		return filepath.Join(home, ".config", "Claude", "claude_config.json"), nil
	}
}

// ReadConfig reads the current Claude Code configuration.
func (m *Manager) ReadConfig() (*ClaudeConfig, error) {
	path, err := m.GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return m.createDefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read Claude config: %w", err)
	}

	var config ClaudeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse Claude config: %w", err)
	}

	return &config, nil
}

// WriteConfig writes the Claude Code configuration.
func (m *Manager) WriteConfig(config *ClaudeConfig) error {
	// Create backup first
	if err := m.createBackup(); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	path, err := m.GetConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// SetModelProvider configures a model for a specific provider.
func (m *Manager) SetModelProvider(model, provider string) error {
	config, err := m.ReadConfig()
	if err != nil {
		return err
	}

	// Check if model already exists
	exists := false
	for i, m := range config.Models {
		if m.Name == model {
			config.Models[i].Provider = provider
			exists = true
			break
		}
	}

	if !exists {
		config.Models = append(config.Models, Model{
			Name:      model,
			Provider:  provider,
			Model:     fmt.Sprintf("%s/%s", provider, model),
			MaxTokens: 8192,
		})
	}

	config.DefaultModel = model
	config.Provider = provider
	config.UpdatedAt = time.Now()

	if err := m.WriteConfig(config); err != nil {
		return err
	}

	m.logger.Info("Model provider configured", "model", model, "provider", provider)
	return nil
}

// GetAvailableModels returns available models from Claude Code.
func (m *Manager) GetAvailableModels() ([]string, error) {
	config, err := m.ReadConfig()
	if err != nil {
		return nil, err
	}

	models := make([]string, 0, len(config.Models))
	for _, model := range config.Models {
		models = append(models, model.Name)
	}

	return models, nil
}

// CreateBackup creates a backup of the current config.
func (m *Manager) createBackup() error {
	path, err := m.GetConfigPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// Create backup directory
	backupDir := filepath.Join(filepath.Dir(path), ".backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	// Create timestamped backup
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("claude_config_%s.json", timestamp))

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return err
	}

	m.logger.Info("Configuration backed up", "path", backupPath)
	return nil
}

// RestoreBackup restores a backup.
func (m *Manager) RestoreBackup(timestamp string) error {
	backupDir := filepath.Join(filepath.Dir(m.config.ClaudeCode.BackupPath), ".backups")

	backupPath := filepath.Join(backupDir, fmt.Sprintf("claude_config_%s.json", timestamp))

	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup: %w", err)
	}

	path, err := m.GetConfigPath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	m.logger.Info("Configuration restored", "backup", timestamp)
	return nil
}

// createDefaultConfig creates a default Claude Code configuration.
func (m *Manager) createDefaultConfig() *ClaudeConfig {
	return &ClaudeConfig{
		Version:        "1.0",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		DefaultModel:   "claude-3-5-sonnet",
		Provider:       "anthropic",
		AutoUpdateModels: true,
		Models: []Model{
			{
				Name:        "claude-3-5-sonnet",
				Provider:    "anthropic",
				Model:       "anthropic/claude-3-5-sonnet",
				MaxTokens:   8192,
				Temperature: 0.7,
			},
			{
				Name:        "claude-opus",
				Provider:    "anthropic",
				Model:       "anthropic/claude-opus",
				MaxTokens:   8192,
				Temperature: 0.7,
			},
			{
				Name:        "gpt-4o",
				Provider:    "openai",
				Model:       "openai/gpt-4o",
				MaxTokens:   8192,
				Temperature: 0.7,
			},
		},
	}
}

// isMacOS checks if running on macOS.
func isMacOS() bool {
	return strings.EqualFold(os.Getenv("GOOS"), "darwin")
}

// isWindows checks if running on Windows.
func isWindows() bool {
	return strings.EqualFold(os.Getenv("GOOS"), "windows")
}
