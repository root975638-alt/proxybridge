// Package config provides the core configuration management for ProxyBridge.
// It handles loading, generating, and validating all configuration files.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	// ConfigVersion is the current configuration schema version
	ConfigVersion = "1.0.0"

	// ConfigFileName is the main config file name
	ConfigFileName = "config.yaml"

	// ProvidersFileName is the providers config file name
	ProvidersFileName = "providers.yaml"

	// AliasesFileName is the aliases config file name
	AliasesFileName = "aliases.yaml"

	// EnvironmentFileName is the environment file name
	EnvironmentFileName = "environment"

	// ConfigDirName is the config directory name
	ConfigDirName = "proxybridge"
)

// Config holds the main ProxyBridge configuration
type Config struct {
	Version         string    `yaml:"version"`
	CreatedAt       time.Time `yaml:"created_at"`
	UpdatedAt       time.Time `yaml:"updated_at"`
	InstallID       string    `yaml:"install_id"`
	InstallationType string   `yaml:"installation_type"`

	// Settings
	Settings Settings `yaml:"settings"`

	// Environment variables (non-sensitive)
	Environment map[string]string `yaml:"environment"`

	// Active provider (alias name)
	ActiveProvider string `yaml:"active_provider"`

	// Default model alias
	DefaultModel string `yaml:"default_model"`

	// LiteLLM configuration
	LiteLLM LiteLLMConfig `yaml:"litellm"`

	// Claude Code configuration
	ClaudeCode ClaudeConfig `yaml:"claude_code"`
}

// Settings contains user-configurable settings
type Settings struct {
	// Logging level: debug, info, warn, error
	LogLevel string `yaml:"log_level"`

	// Enable verbose output
	Verbose bool `yaml:"verbose"`

	// Enable JSON output
	JSONOutput bool `yaml:"json_output"`

	// Skip verification on install
	SkipVerification bool `yaml:"skip_verification"`

	// Enable debug mode
	Debug bool `yaml:"debug"`

	// Auto-start LiteLLM on boot
	AutoStart bool `yaml:"auto_start"`

	// Service type: systemd, launchd, windows, process
	ServiceType string `yaml:"service_type"`

	// Listen address for LiteLLM
	ListenAddress string `yaml:"listen_address"`

	// Port for LiteLLM
	Port int `yaml:"port"`
}

// LiteLLMConfig holds LiteLLM-specific configuration
type LiteLLMConfig struct {
	Path           string `yaml:"path"`
	ConfigPath     string `yaml:"config_path"`
	EnvironmentPath string `yaml:"environment_path"`
	LogPath        string `yaml:"log_path"`
	PIDPath        string `yaml:"pid_path"`
}

// ClaudeConfig holds Claude Code configuration
type ClaudeConfig struct {
	Path        string `yaml:"path"`
	Enabled     bool   `yaml:"enabled"`
	BackupPath  string `yaml:"backup_path"`
	AutoUpdate  bool   `yaml:"auto_update"`
}

// GetConfigDirectory returns the platform-appropriate config directory
func GetConfigDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	switch {
	case isWindows():
		return filepath.Join(home, "AppData", "Roaming", ConfigDirName), nil
	case isMacOS():
		return filepath.Join(home, "Library", "Application Support", ConfigDirName), nil
	default:
		return filepath.Join(home, ".config", ConfigDirName), nil
	}
}

// GetConfigPath returns the full path to the main config file
func GetConfigPath() (string, error) {
	dir, err := GetConfigDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFileName), nil
}

// GetConfigDirectoryWithDefault returns the config directory, creating it if needed
func GetConfigDirectoryWithDefault() (string, error) {
	dir, err := GetConfigDirectory()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return dir, nil
}

// IsConfigExists checks if the config file exists
func IsConfigExists() (bool, error) {
	path, err := GetConfigPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// LoadConfig loads the configuration from disk
func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to disk
func (c *Config) SaveConfig() error {
	dir, err := GetConfigDirectoryWithDefault()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, ConfigFileName)

	c.UpdatedAt = time.Now()

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks the configuration for required fields
func (c *Config) Validate() error {
	if c.Version == "" {
		return fmt.Errorf("config version is required")
	}

	if c.InstallID == "" {
		return fmt.Errorf("install ID is required")
	}

	if c.ActiveProvider == "" {
		c.ActiveProvider = "default"
	}

	if c.DefaultModel == "" {
		c.DefaultModel = "claude-3-5-sonnet"
	}

	if c.Settings.ServiceType == "" {
		c.Settings.ServiceType = "process"
	}

	if c.Settings.ListenAddress == "" {
		c.Settings.ListenAddress = "127.0.0.1"
	}

	if c.Settings.Port == 0 {
		c.Settings.Port = 4000
	}

	if c.LiteLLM.Path == "" {
		c.LiteLLM.Path = "litellm"
	}

	if c.LiteLLM.Port == 0 {
		c.LiteLLM.Port = 4000
	}

	return nil
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	installID, err := uuid.NewUUID()
	if err != nil {
		// Fallback if UUID generation fails
		installID = uuid.Must(uuid.NewRandom())
	}

	return &Config{
		Version: ConfigVersion,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		InstallID: installID.String(),
		InstallationType: "cli",
		Settings: Settings{
			LogLevel:      "info",
			Verbose:       false,
			JSONOutput:    false,
			SkipVerification: false,
			Debug:         false,
			AutoStart:     true,
			ServiceType:   "process",
			ListenAddress: "127.0.0.1",
			Port:          4000,
		},
		ActiveProvider: "default",
		DefaultModel:   "claude-3-5-sonnet",
		LiteLLM: LiteLLMConfig{
			Port: 4000,
		},
		ClaudeCode: ClaudeConfig{
			Enabled:     true,
			AutoUpdate:  true,
		},
		Environment: make(map[string]string),
	}
}

// CreateNewConfig creates a fresh config file
func CreateNewConfig() (*Config, error) {
	config := NewConfig()

	if err := config.SaveConfig(); err != nil {
		return nil, err
	}

	return config, nil
}

// isWindows checks if running on Windows
func isWindows() bool {
	return strings.EqualFold(os.Getenv("OS"), "Windows_NT") ||
		strings.HasPrefix(os.Getenv("PLATFORM"), "mingw") ||
		strings.HasPrefix(os.Getenv("PLATFORM"), "cygwin")
}

// isMacOS checks if running on macOS
func isMacOS() bool {
	return strings.EqualFold(runtime.GOOS, "darwin")
}

// GetLogger returns a configured logger instance
func (c *Config) GetLogger() *Logger {
	return NewLogger(c.Settings.LogLevel, c.Settings.JSONOutput)
}
