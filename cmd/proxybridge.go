// Package cmd implements the ProxyBridge CLI.
//
// ProxyBridge is a production-grade CLI that integrates Claude Code with
// LiteLLM as a universal LLM proxy, supporting any LLM provider through
// a plugin-based architecture.
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/root975638-alt/proxybridge/internal/config"
	"github.com/root975638-alt/proxybridge/internal/diagnostic"
	"github.com/root975638-alt/proxybridge/internal/exportimport"
	"github.com/root975638-alt/proxybridge/internal/installer"
	"github.com/root975638-alt/proxybridge/internal/logging"
	"github.com/root975638-alt/proxybridge/internal/util"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "development"

	// BuildDate is set at build time
	BuildDate = "unknown"

	// Commit is set at build time
	Commit = "unknown"
)

var (
	// Global flags
	rootArgs struct {
		verbose    bool
		jsonOutput bool
		logLevel   string
		configPath string
	}

	// Root command
	rootCmd = &cobra.Command{
		Use:           "proxybridge",
		Short:         "ProxyBridge - Universal LLM Proxy for Claude Code",
		Long:          LongHelp(),
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Initialize logging
			logLevel := "info"
			if rootArgs.verbose {
				logLevel = "debug"
			}
			if rootArgs.logLevel != "" {
				logLevel = rootArgs.logLevel
			}

			logging.Init(logLevel, rootArgs.jsonOutput)

			// Set config path if specified
			if rootArgs.configPath != "" {
				// _, err := config.SetCustomConfigPath(rootArgs.configPath)
				// if err != nil {
				// 	logging.Error("Failed to set custom config path", "error", err)
				// 	os.Exit(1)
				// }
				logging.Warn("Custom config path not currently supported", "path", rootArgs.configPath)
			}

			// Validate config directory exists
			if _, err := config.GetConfigDirectoryWithDefault(); err != nil {
				logging.Error("Failed to access config directory", "error", err)
				os.Exit(1)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

// Execute runs the CLI
func Execute() {
	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logging.Error("Command failed", "error", err)
		os.Exit(getExitCode(err))
	}
}

func getExitCode(err error) int {
	// Map errors to exit codes
	switch {
	case strings.Contains(err.Error(), "permission"):
		return 5
	case strings.Contains(err.Error(), "not found"):
		return 2
	case strings.Contains(err.Error(), "invalid"):
		return 3
	case strings.Contains(err.Error(), "connection"):
		return 4
	default:
		return 1
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&rootArgs.verbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&rootArgs.jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringVar(&rootArgs.logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&rootArgs.configPath, "config", "", "Path to config file (overrides default)")

	// Add subcommands
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(repairCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(modelsCmd)
	rootCmd.AddCommand(providersCmd)
	rootCmd.AddCommand(credentialsCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(benchmarkCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(selfUpdateCmd)
}

// installCmd handles the install command
var installCmd = &cobra.Command{
	Use:     "install",
	Short:   "Install ProxyBridge and configure LiteLLM",
	Long:    `Install ProxyBridge and set up LiteLLM for use with Claude Code.`,
	Example: `  proxybridge install  # Interactive installation`,
	RunE:    runInstall,
}

func runInstall(cmd *cobra.Command, args []string) error {
	return installer.Run(rootArgs.verbose)
}

// uninstallCmd handles the uninstall command
var uninstallCmd = &cobra.Command{
	Use:     "uninstall",
	Short:   "Uninstall ProxyBridge",
	Long:    `Remove ProxyBridge and its configuration from the system.`,
	RunE:    runUninstall,
}

func runUninstall(cmd *cobra.Command, args []string) error {
	return installer.Uninstall()
}

// repairCmd handles the repair command
var repairCmd = &cobra.Command{
	Use:     "repair",
	Short:   "Repair ProxyBridge installation",
	Long:    `Repair ProxyBridge and fix common configuration issues.`,
	RunE:    runRepair,
}

func runRepair(cmd *cobra.Command, args []string) error {
	return installer.Repair()
}

// updateCmd handles the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update ProxyBridge",
	Long:    `Update ProxyBridge to the latest version.`,
	RunE:    runUpdate,
}

func runUpdate(cmd *cobra.Command, args []string) error {
	return installer.Update()
}

// doctorCmd handles the doctor command
var doctorCmd = &cobra.Command{
	Use:     "doctor",
	Short:   "Run diagnostic checks",
	Long:    `Run comprehensive diagnostic checks on the ProxyBridge installation.`,
	RunE:    runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) error {
	opts := diagnostic.Options{
		Verbose: rootArgs.verbose,
		JSON:    rootArgs.jsonOutput,
	}
	return diagnostic.RunAll(opts)
}

// statusCmd handles the status command
var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Show system status",
	Long:    `Display the current status of ProxyBridge, LiteLLM, and connected providers.`,
	RunE:    runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	return installer.Status()
}

// startCmd handles the start command
var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start LiteLLM",
	Long:    `Start the LiteLLM service.`,
	RunE:    runStart,
}

func runStart(cmd *cobra.Command, args []string) error {
	return installer.StartLiteLLM()
}

// stopCmd handles the stop command
var stopCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stop LiteLLM",
	Long:    `Stop the LiteLLM service.`,
	RunE:    runStop,
}

func runStop(cmd *cobra.Command, args []string) error {
	return installer.StopLiteLLM()
}

// restartCmd handles the restart command
var restartCmd = &cobra.Command{
	Use:     "restart",
	Short:   "Restart LiteLLM",
	Long:    `Restart the LiteLLM service.`,
	RunE:    runRestart,
}

func runRestart(cmd *cobra.Command, args []string) error {
	return installer.RestartLiteLLM()
}

// logsCmd handles the logs command
var logsCmd = &cobra.Command{
	Use:     "logs",
	Short:   "View LiteLLM logs",
	Long:    `Display logs from the LiteLLM service.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runLogs,
}

func runLogs(cmd *cobra.Command, args []string) error {
	lines := 100
	if len(args) > 0 {
		lines = 100 // Keep simple for now
	}
	return installer.ShowLogs(lines)
}

// validateCmd handles the validate command
var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate configuration",
	Long:    `Validate the ProxyBridge configuration.`,
	RunE:    runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	return installer.Validate()
}

// testCmd handles the test command
var testCmd = &cobra.Command{
	Use:     "test",
	Short:   "Test provider connection",
	Long:    `Test connection to a specific provider.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runTest,
}

func runTest(cmd *cobra.Command, args []string) error {
	providerName := args[0]
	return installer.TestProvider(providerName)
}

// modelsCmd handles the models command
var modelsCmd = &cobra.Command{
	Use:     "models",
	Short:   "List configured models",
	Long:    `List all configured model aliases and their mappings.`,
	RunE:    runModels,
}

func runModels(cmd *cobra.Command, args []string) error {
	return installer.ListModels()
}

// providersCmd handles the providers command
var providersCmd = &cobra.Command{
	Use:     "providers",
	Short:   "List available providers",
	Long:    `List all supported providers and their status.`,
	RunE:    runProviders,
}

func runProviders(cmd *cobra.Command, args []string) error {
	return installer.ListProviders()
}

// credentialsCmd handles the credentials command
var credentialsCmd = &cobra.Command{
	Use:     "credentials",
	Short:   "Manage credentials",
	Long:    `Manage API credentials for providers.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runCredentials,
}

func runCredentials(cmd *cobra.Command, args []string) error {
	return installer.ManageCredentials(args[0])
}

// exportCmd handles the export command
var exportCmd = &cobra.Command{
	Use:     "export",
	Short:   "Export configuration",
	Long:    `Export ProxyBridge configuration to a file.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runExport,
}

func runExport(cmd *cobra.Command, args []string) error {
	return exportimport.ExportConfig(args[0])
}

// importCmd handles the import command
var importCmd = &cobra.Command{
	Use:     "import",
	Short:   "Import configuration",
	Long:    `Import ProxyBridge configuration from a file.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runImport,
}

func runImport(cmd *cobra.Command, args []string) error {
	return exportimport.ImportConfig(args[0])
}

// configCmd handles the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Show or edit configuration",
	Long:    `Display or modify the ProxyBridge configuration.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runConfig,
}

func runConfig(cmd *cobra.Command, args []string) error {
	option := args[0]
	switch option {
	case "show":
		return installer.ShowConfig()
	case "edit":
		return installer.EditConfig()
	default:
		return fmt.Errorf("unknown config option: %s", option)
	}
}

// switchCmd handles the switch command
var switchCmd = &cobra.Command{
	Use:     "switch",
	Short:   "Switch default provider",
	Long:    `Switch the default provider for model aliases.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runSwitch,
}

func runSwitch(cmd *cobra.Command, args []string) error {
	return installer.SwitchProvider(args[0])
}

// aliasCmd handles the alias command
var aliasCmd = &cobra.Command{
	Use:     "alias",
	Short:   "Manage model aliases",
	Long:    `Create, list, or remove model aliases.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runAlias,
}

func runAlias(cmd *cobra.Command, args []string) error {
	return installer.ManageAliases(args[0])
}

// benchmarkCmd handles the benchmark command
var benchmarkCmd = &cobra.Command{
	Use:     "benchmark",
	Short:   "Run benchmark tests",
	Long:    `Run benchmark tests to measure provider performance.`,
	RunE:    runBenchmark,
}

func runBenchmark(cmd *cobra.Command, args []string) error {
	return installer.RunBenchmark()
}

// versionCmd handles the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Show version information",
	Long:    `Display ProxyBridge version and build information.`,
	Run:     runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("ProxyBridge version %s (commit %s, built %s)\n", Version, Commit, BuildDate)
	fmt.Printf("CLI version: %s\n", Version)
	fmt.Printf("Go version: %s\n", util.GetGoVersion())
	fmt.Printf("Platform: %s/%s\n", util.GetOS(), util.GetArch())
	fmt.Printf("Config dir: %s\n", filepath.Join(util.GetHomeDir(), ".config", "proxybridge"))
}

// selfUpdateCmd handles the self-update command
var selfUpdateCmd = &cobra.Command{
	Use:     "self-update",
	Short:   "Update ProxyBridge",
	Long:    `Update ProxyBridge to the latest version using GitHub releases.`,
	RunE:    runSelfUpdate,
}

func runSelfUpdate(cmd *cobra.Command, args []string) error {
	return installer.SelfUpdate()
}

// LongHelp returns extended help text
func LongHelp() string {
	return `ProxyBridge - Universal LLM Proxy for Claude Code

ProxyBridge enables Claude Code to work with any LLM provider through LiteLLM,
providing a universal interface and model aliasing system.

Usage:
  proxybridge [command]

Examples:
  proxybridge install              # Install and configure everything
  proxybridge doctor               # Run diagnostic checks
  proxybridge status               # Show system status
  proxybridge start                # Start LiteLLM
  proxybridge stop                 # Stop LiteLLM
  proxybridge models               # List configured models
  proxybridge providers            # List available providers

Available Commands:
  install       Install ProxyBridge and configure LiteLLM
  uninstall     Remove ProxyBridge from the system
  repair        Repair the ProxyBridge installation
  update        Update ProxyBridge
  doctor        Run diagnostic checks
  status        Show system status
  start         Start LiteLLM
  stop          Stop LiteLLM
  restart       Restart LiteLLM
  logs          View LiteLLM logs
  validate      Validate configuration
  test          Test provider connection
  models        List configured models
  providers     List available providers
  credentials   Manage API credentials
  export        Export configuration
  import        Import configuration
  config        Show or edit configuration
  switch        Switch default provider
  alias         Manage model aliases
  benchmark     Run benchmark tests
  version       Show version information
  self-update   Update ProxyBridge

Flags:
  -v, --verbose    Enable verbose output
  -j, --json       Output in JSON format
  -h, --help       Help for proxybridge
`
}
