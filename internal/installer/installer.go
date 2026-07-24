// Package installer handles the installation and management of ProxyBridge,
// LiteLLM, and connected providers.
package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/proxybridge/cli/internal/config"
	"github.com/proxybridge/cli/internal/credential"
	"github.com/proxybridge/cli/internal/diagnostic"
	"github.com/proxybridge/cli/internal/logging"
	"github.com/proxybridge/cli/internal/provider"
	"github.com/proxybridge/cli/internal/util/detection"
	"github.com/proxybridge/cli/pkg/manager/litellm"
	"github.com/proxybridge/cli/pkg/manager/claude"
	"github.com/proxybridge/cli/pkg/template"
)

// Run performs the full installation
func Run(verbose bool) error {
	logging.Info("Starting ProxyBridge installation...")

	// Detect platform
	detector := detection.NewDetector()
	platformInfo, err := detector.DetectAll()
	if err != nil {
		logging.Error("Failed to detect platform", "error", err)
		return err
	}

	if verbose {
		logging.Info("Platform detected", "platform", platformInfo)
	}

	// Check for existing installation
	if installed, _ := config.IsConfigExists(); installed {
		logging.Warn("ProxyBridge is already installed")
		logging.Info("Use 'proxybridge repair' to fix issues or 'proxybridge uninstall' to remove")
		return fmt.Errorf("ProxyBridge is already installed")
	}

	// Create new config
	cfg, err := config.CreateNewConfig()
	if err != nil {
		return err
	}

	// Step 1: Install dependencies
	logging.Info("Installing dependencies...")
	if err := installDependencies(platformInfo); err != nil {
		logging.Warn("Dependency installation had warnings", "error", err)
	}

	// Step 2: Install LiteLLM
	logging.Info("Installing LiteLLM...")
	lm := litellm.NewManager(cfg, logging.GetLogger())
	if err := lm.Install(); err != nil {
		logging.Error("Failed to install LiteLLM", "error", err)
		return err
	}

	// Step 3: Install providers
	logging.Info("Installing providers...")
	installedProviders, err := installAllProviders(verbose)
	if err != nil {
		logging.Warn("Some provider installations had warnings", "error", err)
	}

	// Step 4: Detect Claude Code
	logging.Info("Detecting Claude Code...")
	cClaude := claude.NewManager(cfg, logging.GetLogger())
	claudeInstalled, claudePath := cClaude.Detect()

	// Update config with Claude info
	if claudeInstalled {
		cfg.ClaudeCode.Path = claudePath
		cfg.ClaudeCode.Enabled = true
		logging.Info("Claude Code detected", "path", claudePath)
	} else {
		logging.Warn("Claude Code not found - will configure on next install")
	}

	// Step 5: Generate configuration
	logging.Info("Generating configuration...")
	if err := generateConfiguration(cfg, installedProviders, verbose); err != nil {
		logging.Error("Failed to generate configuration", "error", err)
		return err
	}

	// Step 6: Set up credentials
	logging.Info("Setting up credentials...")
	if err := setupCredentials(cfg, verbose); err != nil {
		// Don't fail installation if credentials fail
		logging.Warn("Credential setup incomplete", "error", err)
	}

	// Step 7: Start LiteLLM
	logging.Info("Starting LiteLLM...")
	if err := lm.Start(); err != nil {
		logging.Warn("Failed to start LiteLLM", "error", err)
		logging.Info("LiteLLM configuration created but not started")
		logging.Info("Start LiteLLM manually with: proxybridge start")
	}

	// Step 8: Validate installation
	if !cfg.Settings.SkipVerification {
		logging.Info("Validating installation...")
		validationOpts := diagnostic.Options{
			Verbose: verbose,
			JSON:    false,
		}
		if err := diagnostic.RunAll(validationOpts); err != nil {
			logging.Warn("Validation found issues", "error", err)
			logging.Info("Installation complete but with issues - run 'proxybridge doctor' for details")
		} else {
			logging.Info("Installation validated successfully!")
		}
	}

	logging.Info("Installation complete!")

	// Print next steps
	printNextSteps(cfg)

	return nil
}

// installDependencies installs required system dependencies
func installDependencies(info *detection.PlatformInfo) error {
	var missing []string

	// Check for required tools
	if !info.PythonAvailable {
		missing = append(missing, "python3")
	}
	if !info.NodeAvailable {
		missing = append(missing, "node")
	}

	if len(missing) > 0 {
		logging.Info("Installing missing dependencies...")
		for _, pkg := range missing {
			if err := installPackage(info, pkg); err != nil {
				return fmt.Errorf("failed to install %s: %w", pkg, err)
			}
		}
	}

	return nil
}

// installPackage installs a package using the OS package manager
func installPackage(info *detection.PlatformInfo, pkg string) error {
	var cmd *exec.Cmd

	switch info.PackageManager {
	case "apt":
		cmd = exec.Command("sudo", "apt", "install", "-y", pkg)
	case "dnf":
		cmd = exec.Command("sudo", "dnf", "install", "-y", pkg)
	case "yum":
		cmd = exec.Command("sudo", "yum", "install", "-y", pkg)
	case "pacman":
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", pkg)
	case "brew":
		cmd = exec.Command("brew", "install", pkg)
	case "winget":
		cmd = exec.Command("winget", "install", "--silent", pkg)
	default:
		return fmt.Errorf("no package manager available for %s", info.OS)
	}

	logging.Debug("Running command", "command", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("package install failed: %s", string(output))
	}

	return nil
}

// installAllProviders installs all registered providers
func installAllProviders(verbose bool) ([]string, error) {
	var installed []string
	var errs []error

	providers := provider.GetAvailableProviders()

	for _, p := range providers {
		if verbose {
			logging.Info("Processing provider", "provider", p.ID())
		}

		// Check if provider needs installation
		if !p.IsInstalled() {
			if err := p.Install(); err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", p.ID(), err))
				continue
			}
		}

		installed = append(installed, p.ID())
	}

	if len(errs) > 0 {
		return installed, fmt.Errorf("some providers failed: %v", errs)
	}

	return installed, nil
}

// generateConfiguration creates all configuration files
func generateConfiguration(cfg *config.Config, providers []string, verbose bool) error {
	// Generate LiteLLM configuration
	litellmConfig, err := template.GenerateLiteLLMConfig(cfg, providers)
	if err != nil {
		return err
	}

	// Save LiteLLM config
	litellmPath := filepath.Join(
		filepath.Dir(cfg.LiteLLM.ConfigPath),
		"litellm_config.yaml",
	)
	if err := os.WriteFile(litellmPath, []byte(litellmConfig), 0644); err != nil {
		return err
	}

	// Generate environment file
	envConfig, err := template.GenerateEnvironmentConfig(cfg, providers)
	if err != nil {
		return err
	}
	if err := os.WriteFile(cfg.EnvironmentPath, []byte(envConfig), 0644); err != nil {
		return err
	}

	// Save main config
	if err := cfg.SaveConfig(); err != nil {
		return err
	}

	// Generate Claude Code config
	if cfg.ClaudeCode.Enabled {
		claudeConfig := template.GenerateClaudeConfig(cfg)
		claudePath := filepath.Join(filepath.Dir(cfg.ClaudeCode.Path), ".claude.json")
		if err := os.WriteFile(claudePath, []byte(claudeConfig), 0644); err != nil {
			return err
		}
	}

	return nil
}

// setupCredentials handles credential setup
func setupCredentials(cfg *config.Config, verbose bool) error {
	credentialMgr, err := credential.NewManager(logging.GetLogger())
	if err != nil {
		return err
	}

	// Get available providers
	providers := provider.GetAvailableProviders()

	if verbose {
		logging.Info("Available providers", "count", len(providers))
	}

	// Prompt for credentials (interactive mode)
	// In production, this would prompt for each provider
	// For now, we'll log what would be needed
	if len(providers) > 0 {
		logging.Info("To configure providers, run:")
		logging.Info("  proxybridge providers  # List providers")
		logging.Info("  proxybridge credentials add <provider>  # Add credentials")
	}

	return nil
}

// printNextSteps displays post-installation instructions
func printNextSteps(cfg *config.Config) {
	fmt.Println()
	fmt.Println("Installation complete!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println()
	fmt.Println("  # Start LiteLLM")
	fmt.Printf("  proxybridge start\n")
	fmt.Println()
	fmt.Println("  # Add provider credentials")
	fmt.Printf("  proxybridge credentials add openai  # Follow prompts\n")
	fmt.Println()
	fmt.Println("  # List available providers")
	fmt.Printf("  proxybridge providers\n")
	fmt.Println()
	fmt.Println("  # Validate installation")
	fmt.Printf("  proxybridge doctor\n")
	fmt.Println()
	fmt.Println("  # Configure Claude Code")
	fmt.Printf("  proxybridge config edit\n")
	fmt.Println()
}

// Uninstall removes ProxyBridge
func Uninstall() error {
	logging.Info("Starting uninstall...")

	cfg, err := config.LoadConfig()
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("ProxyBridge not installed")
		}
		return err
	}

	// Step 1: Stop LiteLLM
	lm := litellm.NewManager(cfg, logging.GetLogger())
	if err := lm.Stop(); err != nil {
		logging.Warn("Failed to stop LiteLLM", "error", err)
	}

	// Step 2: Remove LiteLLM
	if err := lm.Uninstall(); err != nil {
		logging.Warn("Failed to uninstall LiteLLM", "error", err)
	}

	// Step 3: Remove config directory
	configDir, err := config.GetConfigDirectory()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(configDir); err != nil {
		return fmt.Errorf("failed to remove config directory: %w", err)
	}

	logging.Info("Uninstall complete!")
	fmt.Println("ProxyBridge has been uninstalled.")
	fmt.Println("Note: API credentials stored in OS credential stores were preserved.")

	return nil
}

// Repair fixes common issues
func Repair() error {
	logging.Info("Starting repair...")

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Detect platform
	detector := detection.NewDetector()
	platformInfo, err := detector.DetectAll()
	if err != nil {
		return err
	}

	// Check and install missing dependencies
	if err := installDependencies(platformInfo); err != nil {
		logging.Warn("Some dependencies could not be installed", "error", err)
	}

	// Verify LiteLLM installation
	lm := litellm.NewManager(cfg, logging.GetLogger())
	if !lm.IsInstalled() {
		logging.Info("LiteLLM not found, reinstalling...")
		if err := lm.Install(); err != nil {
			return err
		}
	}

	// Verify configuration
	if err := lm.Validate(); err != nil {
		logging.Info("LiteLLM config invalid, regenerating...")
		if err := generateConfiguration(cfg, nil, false); err != nil {
			return err
		}
	}

	// Restart LiteLLM
	if err := lm.Restart(); err != nil {
		return err
	}

	logging.Info("Repair complete!")
	fmt.Println("ProxyBridge has been repaired.")

	return nil
}

// Update updates ProxyBridge
func Update() error {
	return SelfUpdate()
}

// Status shows system status
func Status() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("ProxyBridge not installed")
		}
		return err
	}

	// Check LiteLLM
	lm := litellm.NewManager(cfg, logging.GetLogger())
	lmStatus, err := lm.Status()

	// Check providers
	providers := provider.GetAvailableProviders()

	// Print status
	fmt.Println("ProxyBridge Status")
	fmt.Println("==================")
	fmt.Printf("Version: %s\n", config.ConfigVersion)
	fmt.Printf("Config: %s\n", cfg.LiteLLM.ConfigPath)
	fmt.Println()

	fmt.Println("LiteLLM:")
	if err != nil {
		fmt.Printf("  Status: Error - %v\n", err)
	} else if lmStatus.Running {
		fmt.Printf("  Status: Running (PID %d)\n", lmStatus.PID)
		fmt.Printf("  URL: http://%s:%d\n", lmStatus.Address, lmStatus.Port)
	} else {
		fmt.Println("  Status: Not running")
	}
	fmt.Println()

	fmt.Println("Providers:")
	for _, p := range providers {
		status := "installed"
		if !p.IsInstalled() {
			status = "not installed"
		}
		fmt.Printf("  %s: %s\n", p.Name(), status)
	}
	fmt.Println()

	fmt.Println("Active Provider:", cfg.ActiveProvider)
	fmt.Println("Default Model:", cfg.DefaultModel)

	return nil
}

// StartLiteLLM starts LiteLLM
func StartLiteLLM() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	lm := litellm.NewManager(cfg, logging.GetLogger())
	if err := lm.Start(); err != nil {
		return err
	}

	fmt.Println("LiteLLM started successfully!")

	return nil
}

// StopLiteLLM stops LiteLLM
func StopLiteLLM() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	lm := litellm.NewManager(cfg, logging.GetLogger())
	if err := lm.Stop(); err != nil {
		return err
	}

	fmt.Println("LiteLLM stopped successfully!")

	return nil
}

// RestartLiteLLM restarts LiteLLM
func RestartLiteLLM() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	lm := litellm.NewManager(cfg, logging.GetLogger())
	if err := lm.Restart(); err != nil {
		return err
	}

	fmt.Println("LiteLLM restarted successfully!")

	return nil
}

// ShowLogs displays LiteLLM logs
func ShowLogs(lines int) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	lm := litellm.NewManager(cfg, logging.GetLogger())
	logs, err := lm.Logs(lines)
	if err != nil {
		return err
	}

	fmt.Println(logs)

	return nil
}

// Validate validates the installation
func Validate() error {
	opts := diagnostic.Options{
		Verbose: false,
		JSON:    false,
	}
	return diagnostic.RunAll(opts)
}

// TestProvider tests a specific provider
func TestProvider(providerName string) error {
	providerMgr, err := provider.NewManager()
	if err != nil {
		return err
	}

	prov := providerMgr.GetProvider(providerName)
	if prov == nil {
		return fmt.Errorf("provider not found: %s", providerName)
	}

	// Get credential manager
	credentialMgr, err := credential.NewManager(logging.GetLogger())
	if err != nil {
		return err
	}

	// Validate credential
	if err := credentialMgr.ValidateCredential(prov.ID()); err != nil {
		return fmt.Errorf("credential validation failed: %w", err)
	}

	// Test connection
	if err := prov.Validate(); err != nil {
		return fmt.Errorf("provider test failed: %w", err)
	}

	fmt.Printf("Provider %s is working correctly!\n", prov.Name())

	return nil
}

// ListModels lists configured models
func ListModels() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	aliases, err := loadAliases()
	if err != nil {
		return err
	}

	fmt.Println("Configured Models (Aliases):")
	fmt.Println("============================")
	for name, mapping := range aliases {
		fmt.Printf("  %-20s -> %s (provider: %s)\n", name, mapping.Model, mapping.Provider)
	}

	return nil
}

// loadAliases loads alias configuration
func loadAliases() (map[string]*Aliases, error) {
	// Default aliases
	defaultAliases := map[string]*Aliases{
		"claude-3-5-sonnet":  {Model: "claude-3-5-sonnet-20241022", Provider: "anthropic"},
		"claude-opus":        {Model: "claude-opus-4-20250502", Provider: "anthropic"},
		"claude-3-haiku":     {Model: "claude-3-haiku-20240307", Provider: "anthropic"},
		"claude-3-7-sonnet":  {Model: "claude-3-7-sonnet-20250219", Provider: "anthropic"},
		"claude-3-5-haiku":   {Model: "claude-3-5-haiku-20241022", Provider: "anthropic"},
		"claude-3-opus":      {Model: "claude-3-opus-20240229", Provider: "anthropic"},
		"claude-sonnet":      {Model: "claude-3-5-sonnet-20241022", Provider: "anthropic"},
		"claude":             {Model: "claude-3-5-sonnet-20241022", Provider: "anthropic"},
		"coder":              {Model: "claude-3-5-sonnet-20241022", Provider: "anthropic"},
		"reasoning":          {Model: "claude-3-7-sonnet-20250219", Provider: "anthropic"},
		"deepseek":           {Model: "deepseek-chat", Provider: "deepseek"},
		"grok":               {Model: "grok-beta", Provider: "xai"},
		"grok-2":             {Model: "grok-2-1212", Provider: "xai"},
		"qwen":               {Model: "qwen-plus", Provider: "alibaba"},
		"gemini":             {Model: "gemini-1.5-flash", Provider: "google"},
		"gemini-1.5-flash":   {Model: "gemini-1.5-flash", Provider: "google"},
		"gemini-1.5-pro":     {Model: "gemini-1.5-pro", Provider: "google"},
		"gemini-2.0-flash":   {Model: "gemini-2.0-flash-exp", Provider: "google"},
		"gpt-4o":             {Model: "gpt-4o", Provider: "openai"},
		"gpt-4o-mini":        {Model: "gpt-4o-mini", Provider: "openai"},
		"gpt-4-turbo":        {Model: "gpt-4-turbo", Provider: "openai"},
		"gpt-4":              {Model: "gpt-4", Provider: "openai"},
		"gpt-3.5-turbo":      {Model: "gpt-3.5-turbo", Provider: "openai"},
		"gpt":                {Model: "gpt-4o", Provider: "openai"},
		"mistral":            {Model: "mistral-large", Provider: "mistral"},
		"mistral-small":      {Model: "mistral-small", Provider: "mistral"},
		"llama-3.1-70b":      {Model: "llama-3.1-70b-versatile", Provider: "groq"},
		"llama-3.1-8b":       {Model: "llama-3.1-8b-instant", Provider: "groq"},
		"llama-3":            {Model: "llama3-8b-8192", Provider: "groq"},
		"llama-2":            {Model: "llama2-70b-4096", Provider: "groq"},
		"mixtral":            {Model: "mixtral-8x7b-32768", Provider: "groq"},
		"llama":              {Model: "llama3-8b-8192", Provider: "groq"},
		"openai":             {Model: "gpt-4o", Provider: "openai"},
		"anthropic":          {Model: "claude-3-5-sonnet-20241022", Provider: "anthropic"},
		"bedrock":            {Model: "anthropic.claude-3-5-sonnet-20241022-v2:0", Provider: "aws"},
		"azure":              {Model: "gpt-4o", Provider: "azure"},
		"vertex":             {Model: "gemini-1.5-flash", Provider: "google"},
		"openrouter":         {Model: "openai/gpt-4o", Provider: "openrouter"},
		"fireworks":          {Model: "accounts/fireworks/models/firefunction-v2", Provider: "fireworks"},
		"together":           {Model: "togethercomputer/llama-3-70b-chat", Provider: "together"},
		"groq":               {Model: "llama3-8b-8192", Provider: "groq"},
		"deepseek":           {Model: "deepseek-chat", Provider: "deepseek"},
		"deepseek-coder":     {Model: "deepseek-coder", Provider: "deepseek"},
		"mistral":            {Model: "mistral-large", Provider: "mistral"},
		"together":           {Model: "togethercomputer/llama-3-70b-chat", Provider: "together"},
		"fireworks":          {Model: "accounts/fireworks/models/firefunction-v2", Provider: "fireworks"},
		"cerebras":           {Model: "cerebras-llama3.1-8b", Provider: "cerebras"},
		"grok":               {Model: "grok-beta", Provider: "xai"},
		"lm-studio":          {Model: "default", Provider: "lmstudio"},
		"ollama":             {Model: "default", Provider: "ollama"},
	}

	return defaultAliases, nil
}

type Aliases struct {
	Model    string `json:"model"`
	Provider string `json:"provider"`
}

// ListProviders lists available providers
func ListProviders() error {
	providers := provider.GetAvailableProviders()

	fmt.Println("Available Providers:")
	fmt.Println("====================")
	fmt.Println()

	for _, p := range providers {
		status := "installed"
		if !p.IsInstalled() {
			status = "not installed"
		}
		fmt.Printf("%-20s %s\n", p.Name(), status)
		fmt.Printf("  ID: %s\n", p.ID())
		fmt.Printf("  Models: %s\n", strings.Join(p.GetModels(), ", "))
		fmt.Println()
	}

	return nil
}

// ManageCredentials manages credentials
func ManageCredentials(action string) error {
	credentialMgr, err := credential.NewManager(logging.GetLogger())
	if err != nil {
		return err
	}

	switch action {
	case "list":
		return credentialMgr.ListCredentials()
	case "add":
		return credentialMgr.AddCredential()
	case "remove":
		return credentialMgr.RemoveCredential()
	case "validate":
		return credentialMgr.ValidateAll()
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

// ManageAliases manages model aliases
func ManageAliases(action string) error {
	switch action {
	case "list":
		return ListModels()
	case "add":
		return AddAlias()
	case "remove":
		return RemoveAlias()
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

// AddAlias adds a new model alias
func AddAlias() error {
	// Implementation would prompt for alias details
	fmt.Println("Adding model alias - interactive mode not yet implemented")
	return nil
}

// RemoveAlias removes a model alias
func RemoveAlias() error {
	// Implementation would remove alias
	fmt.Println("Removing model alias - interactive mode not yet implemented")
	return nil
}

// RunBenchmark runs benchmark tests
func RunBenchmark() error {
	fmt.Println("Benchmark tests - implementation not yet complete")
	return nil
}

// ShowConfig shows the current configuration
func ShowConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	data, err := config.PrettyPrintConfig(cfg)
	if err != nil {
		return err
	}

	fmt.Println(data)

	return nil
}

// EditConfig opens the config file for editing
func EditConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}

	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// SwitchProvider switches the default provider
func SwitchProvider(providerName string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Verify provider exists
	providers := provider.GetAvailableProviders()
	providerFound := false
	for _, p := range providers {
		if strings.EqualFold(p.ID(), providerName) {
			providerFound = true
			break
		}
	}

	if !providerFound {
		return fmt.Errorf("provider not found: %s", providerName)
	}

	// Update config
	cfg.ActiveProvider = providerName

	if err := cfg.SaveConfig(); err != nil {
		return err
	}

	fmt.Printf("Default provider switched to: %s\n", providerName)

	return nil
}

// SelfUpdate update ProxyBridge
func SelfUpdate() error {
	fmt.Println("Self-update not yet implemented")
	return nil
}
