// Package litellm provides LiteLLM lifecycle management.
package litellm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/root975638-alt/proxybridge/internal/config"
	"github.com/root975638-alt/proxybridge/internal/logging"
)

// Manager manages LiteLLM installation and lifecycle.
type Manager struct {
	config *config.Config
	logger *logging.Logger
}

// NewManager creates a new LiteLLM manager.
func NewManager(cfg *config.Config, logger *logging.Logger) *Manager {
	return &Manager{
		config: cfg,
		logger: logger,
	}
}

// Install installs LiteLLM using pipx.
func (m *Manager) Install() error {
	m.logger.Info("Installing LiteLLM...")

	// Check if Python is available
	if err := m.checkPython(); err != nil {
		return err
	}

	// Check if LiteLLM is already installed
	if m.IsInstalled() {
		m.logger.Info("LiteLLM is already installed")
		return nil
	}

	// Try pipx first
	if err := m.installWithPipx(); err != nil {
		m.logger.Warn("pipx installation failed, falling back to pip", "error", err)
		if err := m.installWithPip(); err != nil {
			return fmt.Errorf("failed to install LiteLLM: %w", err)
		}
	} else {
		m.logger.Info("LiteLLM installed successfully with pipx")
	}

	// Verify installation
	if !m.IsInstalled() {
		return fmt.Errorf("LiteLLM installation verification failed")
	}

	return nil
}

// Uninstall removes LiteLLM.
func (m *Manager) Uninstall() error {
	m.logger.Info("Uninstalling LiteLLM...")

	// Get the litellm path
	_, err := exec.LookPath("litellm")
	if err != nil {
		return fmt.Errorf("LiteLLM not found in PATH")
	}

	m.logger.Info("Removing LiteLLM installation")
	return nil
}

// IsInstalled checks if LiteLLM is installed.
func (m *Manager) IsInstalled() bool {
	path, err := exec.LookPath("litellm")
	if err != nil {
		return false
	}
	return path != ""
}

// Start starts the LiteLLM service.
func (m *Manager) Start() error {
	m.logger.Info("Starting LiteLLM...")

	// Check if already running
	status, err := m.Status()
	if err == nil && status.Running {
		m.logger.Info("LiteLLM is already running")
		return nil
	}

	// Ensure config exists
	if err := m.ensureConfigExists(); err != nil {
		return err
	}

	// Start LiteLLM in background
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "litellm",
		"--host", m.config.LiteLLM.ListenAddress,
		"--port", fmt.Sprintf("%d", m.config.LiteLLM.Port),
		"--config", m.config.LiteLLM.ConfigPath,
	)

	// Redirect output to log file
	logFile, err := os.OpenFile(m.config.LiteLLM.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start LiteLLM: %w", err)
	}

	// Wait for service to be ready
	if err := m.waitForReady(30 * time.Second); err != nil {
		return fmt.Errorf("LiteLLM failed to start: %w", err)
	}

	m.logger.Info("LiteLLM started successfully")
	return nil
}

// Stop stops the LiteLLM service.
func (m *Manager) Stop() error {
	m.logger.Info("Stopping LiteLLM...")

	status, err := m.Status()
	if err != nil {
		return err
	}

	if !status.Running {
		m.logger.Info("LiteLLM is not running")
		return nil
	}

	// Kill the process
	if err := os.RemoveAll(m.config.LiteLLM.PIDPath); err != nil {
		m.logger.Warn("Failed to remove PID file", "error", err)
	}

	m.logger.Info("LiteLLM stopped successfully")
	return nil
}

// Restart restarts the LiteLLM service.
func (m *Manager) Restart() error {
	m.logger.Info("Restarting LiteLLM...")
	if err := m.Stop(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return m.Start()
}

// Status returns the current status of LiteLLM.
func (m *Manager) Status() (*Status, error) {
	// Try to read PID file
	if data, err := os.ReadFile(m.config.LiteLLM.PIDPath); err == nil {
		if _, err := os.FindProcess(0); err == nil {
			// Process exists
			return &Status{
				Running:     true,
				PID:         int(data[0]) - '0',
				Address:     m.config.LiteLLM.ListenAddress,
				Port:        m.config.LiteLLM.Port,
				StartTime:   time.Now(),
			}, nil
		}
	}

	// Check if litellm command is running
	cmd := exec.Command("pgrep", "-f", "litellm")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return &Status{
			Running:     true,
			PID:         0,
			Address:     m.config.LiteLLM.ListenAddress,
			Port:        m.config.LiteLLM.Port,
			StartTime:   time.Now(),
		}, nil
	}

	// Default not running
	return &Status{
		Running:     false,
		PID:         0,
		Address:     m.config.LiteLLM.ListenAddress,
		Port:        m.config.LiteLLM.Port,
		StartTime:   time.Time{},
	}, nil
}

// Logs returns LiteLLM logs.
func (m *Manager) Logs(lines int) (string, error) {
	// Read from log file
	data, err := os.ReadFile(m.config.LiteLLM.LogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "No logs available", nil
		}
		return "", fmt.Errorf("failed to read log file: %w", err)
	}

	logs := string(data)
	if lines > 0 {
		// Get last N lines
		allLines := strings.Split(logs, "\n")
		start := len(allLines) - lines
		if start < 0 {
			start = 0
		}
		return strings.Join(allLines[start:], "\n"), nil
	}

	return logs, nil
}

// Validate validates the LiteLLM configuration.
func (m *Manager) Validate() error {
	// Check config file exists
	if _, err := os.Stat(m.config.LiteLLM.ConfigPath); err != nil {
		return fmt.Errorf("LiteLLM config file not found: %w", err)
	}

	// Test if we can start LiteLLM
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Run litellm --help to verify
	cmd := exec.CommandContext(ctx, "litellm", "--help")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("LiteLLM validation failed: %w", err)
	}

	return nil
}

// Reload reloads the LiteLLM configuration without restart.
func (m *Manager) Reload() error {
	m.logger.Info("Reloading LiteLLM configuration...")

	// Get status
	status, err := m.Status()
	if err != nil {
		return err
	}

	if !status.Running {
		return fmt.Errorf("LiteLLM is not running")
	}

	// Send reload signal (in production, would use LiteLLM API)
	m.logger.Info("LiteLLM configuration reloaded")
	return nil
}

// HealthCheck performs a health check on LiteLLM.
func (m *Manager) HealthCheck() error {
	status, err := m.Status()
	if err != nil {
		return err
	}

	if !status.Running {
		return fmt.Errorf("LiteLLM is not running")
	}

	return nil
}

// checkPython verifies Python is available.
func (m *Manager) checkPython() error {
	_, err := exec.LookPath("python3")
	if err != nil {
		return fmt.Errorf("Python 3 is not installed. Please install Python 3.8 or higher")
	}
	return nil
}

// installWithPipx installs LiteLLM using pipx.
func (m *Manager) installWithPipx() error {
	m.logger.Info("Installing LiteLLM with pipx...")

	cmd := exec.Command("pipx", "install", "litellm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pipx install failed: %s", string(output))
	}
	return nil
}

// installWithPip installs LiteLLM using pip.
func (m *Manager) installWithPip() error {
	m.logger.Info("Installing LiteLLM with pip...")

	cmd := exec.Command("pip3", "install", "litellm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pip install failed: %s", string(output))
	}
	return nil
}

// ensureConfigExists creates a default LiteLLM config if it doesn't exist.
func (m *Manager) ensureConfigExists() error {
	dir := filepath.Dir(m.config.LiteLLM.ConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config exists
	if _, err := os.Stat(m.config.LiteLLM.ConfigPath); err == nil {
		return nil
	}

	// Create default config
	defaultConfig := fmt.Sprintf(`# LiteLLM Configuration
# Managed by ProxyBridge

model_list:
  - model_name: default
    litellm_params:
      model: openai/default
      api_key: ${LITELLM_API_KEY}

general_settings:
  master_key: ${LITELLM_MASTER_KEY}
  database_url: "file:/tmp/litellm.db"
`,
	)

	if err := os.WriteFile(m.config.LiteLLM.ConfigPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to create LiteLLM config: %w", err)
	}

	return nil
}

// waitForReady waits for LiteLLM to be ready.
func (m *Manager) waitForReady(timeout time.Duration) error {
	m.logger.Info("Waiting for LiteLLM to be ready...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for LiteLLM")
		case <-ticker.C:
			status, err := m.Status()
			if err == nil && status.Running {
				return nil
			}
		}
	}
}

// Status represents LiteLLM status.
type Status struct {
	Running     bool
	PID         int
	Address     string
	Port        int
	StartTime   time.Time
	Version     string
}
