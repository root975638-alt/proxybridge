// Package exportimport provides configuration export and import functionality.
package exportimport

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/root975638-alt/proxybridge/internal/config"
	"github.com/root975638-alt/proxybridge/internal/logging"
	"gopkg.in/yaml.v3"
)

// ExportConfig exports ProxyBridge configuration to a file.
func ExportConfig(outputPath string) error {
	logging.Info("Exporting configuration...")

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create export data
	exportData := map[string]interface{}{
		"version":      config.ConfigVersion,
		"exported_at":  time.Now().UTC().Format(time.RFC3339),
		"install_id":   cfg.InstallID,
		"config":       cfg,
	}

	// Determine format based on file extension
	ext := filepath.Ext(outputPath)
	var data []byte

	switch ext {
	case ".json":
		data, err = json.MarshalIndent(exportData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
	case ".yaml", ".yml":
		data, err = yaml.Marshal(exportData)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}
	default:
		// Default to YAML
		data, err = yaml.Marshal(exportData)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	logging.Info("Configuration exported successfully", "path", outputPath)
	fmt.Printf("Configuration exported to: %s\n", outputPath)

	return nil
}

// ImportConfig imports ProxyBridge configuration from a file.
func ImportConfig(inputPath string) error {
	logging.Info("Importing configuration...")

	// Check if file exists
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("config file not found: %w", err)
	}

	// Read file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Determine format and unmarshal
	var exportData map[string]interface{}

	ext := filepath.Ext(inputPath)
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &exportData); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &exportData); err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
	default:
		// Try JSON first, then YAML
		if err := json.Unmarshal(data, &exportData); err == nil {
			break
		}
		if err := yaml.Unmarshal(data, &exportData); err != nil {
			return fmt.Errorf("failed to parse config (not valid JSON or YAML)")
		}
	}

	// Extract config from export data
	rawConfig, ok := exportData["config"]
	if !ok {
		return fmt.Errorf("export file missing 'config' field")
	}

	// Convert to config struct
	configData, err := json.Marshal(rawConfig)
	if err != nil {
		return fmt.Errorf("failed to process config: %w", err)
	}

	var cfg config.Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create backup of existing config
	exists, _ := config.IsConfigExists()
	if exists {
		if err := createBackup(); err != nil {
			logging.Warn("Failed to create backup before import", "error", err)
		}
	}

	// Save new config
	if err := cfg.SaveConfig(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	logging.Info("Configuration imported successfully")
	fmt.Println("Configuration imported successfully!")

	return nil
}

// createBackup creates a backup of the current configuration.
func createBackup() error {
	cfgDir, err := config.GetConfigDirectory()
	if err != nil {
		return err
	}

	// Create backups directory
	backupsDir := filepath.Join(cfgDir, "backups")
	if err := os.MkdirAll(backupsDir, 0755); err != nil {
		return fmt.Errorf("failed to create backups directory: %w", err)
	}

	// Read current config
	configPath := filepath.Join(cfgDir, config.ConfigFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Create timestamped backup
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupsDir, fmt.Sprintf("config_%s.yaml", timestamp))

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	logging.Info("Configuration backed up", "path", backupPath)
	return nil
}

// ExportConfiguration exports only the configuration without export metadata.
func ExportConfiguration(outputPath string) error {
	logging.Info("Exporting configuration data...")

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Marshal config
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	logging.Info("Configuration exported successfully", "path", outputPath)
	return nil
}

// ImportConfiguration imports configuration data.
func ImportConfiguration(inputPath string) error {
	logging.Info("Importing configuration data...")

	// Read file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Save
	if err := cfg.SaveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	logging.Info("Configuration imported successfully")
	return nil
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// GetExportFormats returns available export formats.
func GetExportFormats() []string {
	return []string{"yaml", "json", "yml"}
}
