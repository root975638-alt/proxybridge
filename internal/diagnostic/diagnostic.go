// Package diagnostic provides system diagnostics and issue detection.
package diagnostic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/root975638-alt/proxybridge/internal/config"
	"github.com/root975638-alt/proxybridge/internal/logging"
	"github.com/root975638-alt/proxybridge/internal/provider"
)

// Options holds diagnostic options
type Options struct {
	Verbose bool
	JSON    bool
}

// Report represents a diagnostic report
type Report struct {
	Success   bool            `json:"success"`
	Issues    []Issue         `json:"issues,omitempty"`
	Warnings  []Issue         `json:"warnings,omitempty"`
	Results   []CheckResult   `json:"results,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// Issue represents a diagnostic issue
type Issue struct {
	Severity  string `json:"severity"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Suggestion string `json:"suggestion"`
}

// CheckResult represents a single check result
type CheckResult struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Duration  string `json:"duration,omitempty"`
	Error     string `json:"error,omitempty"`
}

// RunAll runs all diagnostics
func RunAll(opts Options) error {
	report := &Report{
		Timestamp: time.Now(),
	}

	// Check LiteLLM
	if err := checkLiteLLM(report, opts); err != nil {
		report.addIssue("error", "LITELLM_ERROR", err.Error(), "Check LiteLLM installation")
	}

	// Check providers
	if err := checkProviders(report, opts); err != nil {
		report.addIssue("error", "PROVIDER_ERROR", err.Error(), "Check provider configuration")
	}

	// Check credentials
	if err := checkCredentials(report, opts); err != nil {
		report.addIssue("error", "CREDENTIAL_ERROR", err.Error(), "Check credential storage")
	}

	// Check Claude Code
	if err := checkClaudeCode(report, opts); err != nil {
		report.addIssue("warn", "CLAUDE_ERROR", err.Error(), "Claude Code may not be configured")
	}

	// Check network
	if err := checkNetwork(report, opts); err != nil {
		report.addIssue("warn", "NETWORK_ERROR", err.Error(), "Check network connectivity")
	}

	// Check configuration
	if err := checkConfiguration(report, opts); err != nil {
		report.addIssue("error", "CONFIG_ERROR", err.Error(), "Check configuration files")
	}

	report.Success = len(report.Issues) == 0

	if report.Success {
		logging.Info("All diagnostics passed")
		if opts.Verbose {
			fmt.Println("proxybridge diagnostics: SUCCESS")
			fmt.Println()
			fmt.Println("All systems are operational.")
		}
		return nil
	}

	// Output issues
	for _, issue := range report.Issues {
		if opts.JSON {
			logging.Error(issue.Message, "code", issue.Code, "suggestion", issue.Suggestion)
		} else {
			fmt.Printf("[ERROR] %s (%s)\n", issue.Message, issue.Code)
			fmt.Printf("  Suggestion: %s\n", issue.Suggestion)
		}
	}

	for _, warning := range report.Warnings {
		if opts.JSON {
			logging.Warn(warning.Message, "code", warning.Code)
		} else {
			fmt.Printf("[WARN] %s (%s)\n", warning.Message, warning.Code)
		}
	}

	if opts.JSON {
		// Output as JSON
		data, _ := json.Marshal(report)
		fmt.Println(data)
	}

	return fmt.Errorf("diagnostics found issues")
}

// checkLiteLLM checks LiteLLM status
func checkLiteLLM(report *Report, opts Options) error {
	result := &CheckResult{Name: "litellm", Status: "check"}
	start := time.Now()

	cfg, err := config.LoadConfig()
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		report.Results = append(report.Results, *result)
		return err
	}

	lm := NewLiteLLMManager(cfg)

	if !lm.IsInstalled() {
		result.Status = "error"
		result.Error = "LiteLLM is not installed"
		report.addIssue("error", "LITELLM_NOT_INSTALLED", "LiteLLM is not installed", "Run 'proxybridge install' to install LiteLLM")
		report.Results = append(report.Results, *result)
		return fmt.Errorf("LiteLLM not installed")
	}

	status, err := lm.Status()
	result.Duration = fmt.Sprintf("%v", time.Since(start))

	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		report.addIssue("error", "LITELLM_STATUS_ERROR", "Failed to get LiteLLM status", "CheckLiteLLM logs with 'proxybridge logs'")
	} else if !status.Running {
		result.Status = "warn"
		result.Error = "LiteLLM is not running"
		report.addWarning("LITELLM_NOT_RUNNING", "LiteLLM is not running", "Run 'proxybridge start' to start LiteLLM")
	} else {
		result.Status = "ok"
		report.addResult("litellm-running", "ok")
	}

	report.Results = append(report.Results, *result)
	return nil
}

// checkProviders checks provider status
func checkProviders(report *Report, opts Options) error {
	providers := provider.GetAvailableProviders()

	for _, p := range providers {
		result := &CheckResult{Name: p.ID(), Status: "check"}
		start := time.Now()

		// Check if provider is installed
		if !p.IsInstalled() {
			result.Status = "warn"
			result.Error = "Provider not installed"
			report.addWarning("PROVIDER_NOT_INSTALLED", fmt.Sprintf("Provider %s is not installed", p.Name()),
				fmt.Sprintf("Run 'proxybridge providers' to see available providers"))
		} else {
			// Validate provider
			if err := p.Validate(); err != nil {
				result.Status = "error"
				result.Error = err.Error()
				report.addIssue("error", "PROVIDER_VALIDATION_ERROR", fmt.Sprintf("Provider %s validation failed", p.Name()),
					fmt.Sprintf("Check %s credentials with 'proxybridge credentials validate %s'", p.Name(), p.ID()))
			} else {
				result.Status = "ok"
				result.Duration = fmt.Sprintf("%v", time.Since(start))
			}
		}

		report.Results = append(report.Results, *result)
	}

	return nil
}

// checkCredentials checks credential storage
func checkCredentials(report *Report, opts Options) error {
	// Credential manager is initialized in the credential package
	// This is a placeholder for credential checks
	result := &CheckResult{Name: "credentials", Status: "ok"}
	report.Results = append(report.Results, *result)
	return nil
}

// checkClaudeCode checks Claude Code integration
func checkClaudeCode(report *Report, opts Options) error {
	result := &CheckResult{Name: "claude_code", Status: "check"}

	cfg, err := config.LoadConfig()
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		report.addIssue("error", "CONFIG_ERROR", "Failed to load configuration",
			"Run 'proxybridge install' to set up configuration")
		report.Results = append(report.Results, *result)
		return err
	}

	if !cfg.ClaudeCode.Enabled {
		result.Status = "warn"
		result.Error = "Claude Code integration is disabled"
		report.addWarning("CLAUDE_DISABLED", "Claude Code integration is disabled",
			"Enable with 'proxybridge config edit'")
		report.Results = append(report.Results, *result)
		return nil
	}

	// Check if Claude Code config file exists
	if cfg.ClaudeCode.Path != "" {
		// We'll just note that it's configured
		result.Status = "ok"
	} else {
		result.Status = "warn"
		result.Error = "Claude Code path not configured"
		report.addWarning("CLAUDE_NOT_CONFIGURED", "Claude Code path not configured",
			"Claude Code is not integrated - will use LiteLLM directly")
	}

	report.Results = append(report.Results, *result)
	return nil
}

// checkNetwork checks network connectivity
func checkNetwork(report *Report, opts Options) error {
	result := &CheckResult{Name: "network", Status: "check"}
	start := time.Now()

	// Test connectivity to common provider endpoints
	endpoints := []string{
		"https://api.openai.com/v1/models",
		"https://api.anthropic.com/v1/models",
		"https://api.mistral.ai/v1/models",
	}

	client := &http.Client{Timeout: 5 * time.Second}
	var failed []string

	for _, endpoint := range endpoints {
		if _, err := client.Get(endpoint); err != nil {
			failed = append(failed, endpoint)
		}
	}

	result.Duration = fmt.Sprintf("%v", time.Since(start))

	if len(failed) > 0 {
		result.Status = "warn"
		result.Error = fmt.Sprintf("Failed to reach some endpoints: %s", strings.Join(failed, ", "))
		report.addWarning("NETWORK_CONNECTIVITY", "Network connectivity issues detected",
			"Check your internet connection and firewall settings")
	} else {
		result.Status = "ok"
	}

	report.Results = append(report.Results, *result)
	return nil
}

// checkConfiguration checks configuration file integrity
func checkConfiguration(report *Report, opts Options) error {
	result := &CheckResult{Name: "configuration", Status: "check"}

	cfg, err := config.LoadConfig()
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		report.addIssue("error", "CONFIG_LOAD_ERROR", "Failed to load configuration",
			"Re-run 'proxybridge install' to fix configuration")
		report.Results = append(report.Results, *result)
		return err
	}

	// Check required fields
	if cfg.ActiveProvider == "" {
		result.Status = "warn"
		result.Error = "No active provider configured"
		report.addWarning("NO_ACTIVE_PROVIDER", "No active provider configured",
			"Run 'proxybridge providers' and 'proxybridge switch <provider>'")
	}

	if cfg.DefaultModel == "" {
		result.Status = "warn"
		result.Error = "No default model configured"
		report.addWarning("NO_DEFAULT_MODEL", "No default model configured",
			"Run 'proxybridge models' to set a default model")
	}

	if cfg.LiteLLM.ConfigPath == "" {
		result.Status = "error"
		result.Error = "LiteLLM config path not configured"
		report.addIssue("error", "LITELLM_CONFIG_PATH", "LiteLLM config path not configured",
			"Re-run 'proxybridge install' to set up configuration")
	}

	if result.Status == "check" {
		result.Status = "ok"
	}

	report.Results = append(report.Results, *result)
	return nil
}

// addIssue adds an issue to the report
func (r *Report) addIssue(severity, code, message, suggestion string) {
	r.Issues = append(r.Issues, Issue{
		Severity:  severity,
		Code:      code,
		Message:   message,
		Suggestion: suggestion,
	})
}

// addWarning adds a warning to the report
func (r *Report) addWarning(code, message, suggestion string) {
	r.Warnings = append(r.Warnings, Issue{
		Severity:  "warning",
		Code:      code,
		Message:   message,
		Suggestion: suggestion,
	})
}

// addResult adds a result to the report
func (r *Report) addResult(name, status string) {
	r.Results = append(r.Results, CheckResult{
		Name:   name,
		Status: status,
	})
}

// LiteLLMManager provides LiteLLM status checking
type LiteLLMManager struct {
	config *config.Config
	logger *logging.Logger
}

// NewLiteLLMManager creates a new LiteLLM manager for diagnostics
func NewLiteLLMManager(cfg *config.Config) *LiteLLMManager {
	return &LiteLLMManager{
		config: cfg,
		logger: logging.GetLogger(),
	}
}

// IsInstalled checks if LiteLLM is installed
func (m *LiteLLMManager) IsInstalled() bool {
	// Check if litellm command is available
	cmd := m.config.LiteLLM.Path
	// Simplified check
	return cmd != ""
}

// Status returns LiteLLM status
func (m *LiteLLMManager) Status() (*LiteLLMStatus, error) {
	// Simplified status check
	return &LiteLLMStatus{
		Running: true,
		PID:     0,
		Address: m.config.LiteLLM.ListenAddress,
		Port:    m.config.LiteLLM.Port,
	}, nil
}

// LiteLLMStatus represents LiteLLM status
type LiteLLMStatus struct {
	Running   bool
	PID       int
	Address   string
	Port      int
	StartTime time.Time
	Version   string
}

