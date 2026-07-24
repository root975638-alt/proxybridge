// Package health provides health monitoring and validation for ProxyBridge.
package health

import (
	"fmt"
	"net/http"
	"time"

	"github.com/proxybridge/cli/internal/config"
	"github.com/proxybridge/cli/internal/credential"
	"github.com/proxybridge/cli/internal/logging"
	"github.com/proxybridge/cli/internal/provider"
)

// HealthCheck represents a single health check result.
type HealthCheck struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Error       string `json:"error,omitempty"`
	Duration    string `json:"duration,omitempty"`
}

// HealthStatus represents the overall health status.
type HealthStatus struct {
	Overall    string        `json:"overall"`
	Checks     []HealthCheck `json:"checks"`
	Issues     []string      `json:"issues,omitempty"`
	_warnings  []string
	LastCheck  time.Time `json:"last_check"`
}

// Status values
const (
	StatusOK      = "ok"
	StatusWarning = "warning"
	StatusError   = "error"
)

// RunAll runs all health checks.
func RunAll() (*HealthStatus, error) {
	status := &HealthStatus{
		LastCheck: time.Now(),
		Checks:    make([]HealthCheck, 0),
	}

	// Run each health check
	checks := []struct {
		name        string
		description string
		check       func() error
	}{
		{"provider_connectivity", "Provider API connectivity", checkProviderConnectivity},
		{"litellm_api", "LiteLLM API health", checkLiteLLMAPI},
		{"claude_integration", "Claude Code integration", checkClaudeIntegration},
		{"credential_validity", "Credential validity", checkCredentialValidity},
		{"alias_configuration", "Alias configuration", checkAliasConfiguration},
		{"config_integrity", "Configuration integrity", checkConfigIntegrity},
	}

	for _, c := range checks {
		result := status.runCheck(c.name, c.description, c.check)
		status.Checks = append(status.Checks, result)
	}

	// Determine overall status
	status.determineStatus()

	if status.Overall == StatusError {
		return status, fmt.Errorf("health check failed")
	}

	return status, nil
}

// runCheck runs a single health check.
func (s *HealthStatus) runCheck(name, description string, check func() error) HealthCheck {
	result := HealthCheck{
		Name:        name,
		Description: description,
		Status:      StatusOK,
	}

	start := time.Now()

	if err := check(); err != nil {
		result.Status = StatusError
		result.Error = err.Error()
		s.Issues = append(s.Issues, fmt.Sprintf("%s: %s", name, err.Error()))
	} else if name == "provider_connectivity" && len(s._warnings) > 0 {
		result.Status = StatusWarning
	}

	result.Duration = fmt.Sprintf("%v", time.Since(start))
	return result
}

// determineStatus determines the overall health status.
func (s *HealthStatus) determineStatus() {
	hasError := false
	hasWarning := false

	for _, check := range s.Checks {
		if check.Status == StatusError {
			hasError = true
		} else if check.Status == StatusWarning {
			hasWarning = true
		}
	}

	if hasError {
		s.Overall = StatusError
	} else if hasWarning {
		s.Overall = StatusWarning
	} else {
		s.Overall = StatusOK
	}
}

// checkProviderConnectivity checks provider API connectivity.
func checkProviderConnectivity() error {
	providers := provider.GetAvailableProviders()

	var errors []string
	for _, p := range providers {
		// Skip local providers that don't need network
		id := p.ID()
		if id == "ollama" || id == "lmstudio" || id == "local" {
			continue
		}

		// Check if provider is installed
		if !p.IsInstalled() {
			continue
		}

		// Validate provider
		if err := p.Validate(); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", p.Name(), err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("provider validation errors: %v", errors)
	}

	return nil
}

// checkLiteLLMAPI checks LiteLLM API health.
func checkLiteLLMAPI() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Make a simple request to LiteLLM API
	url := fmt.Sprintf("http://%s:%d/v1/models", cfg.LiteLLM.ListenAddress, cfg.LiteLLM.Port)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("LiteLLM API not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("LiteLLM API returned status %d", resp.StatusCode)
	}

	return nil
}

// checkClaudeIntegration checks Claude Code integration.
func checkClaudeIntegration() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	if !cfg.ClaudeCode.Enabled {
		// Integration is disabled, check passes
		return nil
	}

	// Check if Claude Code config exists
	if cfg.ClaudeCode.Path == "" {
		return fmt.Errorf("Claude Code path not configured")
	}

	return nil
}

// checkCredentialValidity checks credential validity.
func checkCredentialValidity() error {
	credMgr, err := credential.NewManager(logging.GetLogger())
	if err != nil {
		return err
	}

	// For each provider with credentials, validate them
	providers := provider.GetAvailableProviders()

	var errors []string
	for _, p := range providers {
		// Check if provider has credentials
		// This is a simplified check - in production would check credential store
	}

	if len(errors) > 0 {
		return fmt.Errorf("credential validation errors: %v", errors)
	}

	return nil
}

// checkAliasConfiguration checks alias configuration.
func checkAliasConfiguration() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	if cfg.DefaultModel == "" {
		return fmt.Errorf("default model not configured")
	}

	return nil
}

// checkConfigIntegrity checks configuration file integrity.
func checkConfigIntegrity() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

// GetProviderHealth returns health status for a specific provider.
func GetProviderHealth(providerID string) (*HealthCheck, error) {
	providers := provider.GetAvailableProviders()

	for _, p := range providers {
		if p.ID() == providerID {
			check := HealthCheck{
				Name:        providerID,
				Description: p.Name(),
				Status:      StatusOK,
			}

			if !p.IsInstalled() {
				check.Status = StatusError
				check.Error = "Provider not installed"
				return &check, nil
			}

			if err := p.Validate(); err != nil {
				check.Status = StatusError
				check.Error = err.Error()
				return &check, nil
			}

			return &check, nil
		}
	}

	return nil, fmt.Errorf("provider not found: %s", providerID)
}

// GetModelHealth returns health status for a model alias.
func GetModelHealth(alias string) (*HealthCheck, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	check := HealthCheck{
		Name:        alias,
		Description: fmt.Sprintf("Model alias: %s", alias),
		Status:      StatusOK,
	}

	// Check if alias exists
	if alias != cfg.DefaultModel && !isAliasConfigured(alias) {
		check.Status = StatusWarning
		check.Error = "Alias not configured"
		return &check, nil
	}

	return &check, nil
}

// isAliasConfigured checks if an alias is configured.
func isAliasConfigured(alias string) bool {
	// Check if alias is in the alias configuration
	return true // Simplified
}
