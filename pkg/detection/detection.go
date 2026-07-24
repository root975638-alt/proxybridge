// Package detection provides platform and tool detection capabilities.
package detection

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// PlatformInfo holds detected platform information
type PlatformInfo struct {
	OS               string
	Architecture     string
	Shell            string
	PackageManager   string
	PythonAvailable  bool
	PythonVersion    string
	NodeAvailable    bool
	NodeVersion      string
	ClaudeAvailable  bool
	ClaudePath       string
	LiteLLMAvailable bool
	LiteLLMVersion   string
	DockerAvailable  bool
	DockerVersion    string
	GitAvailable     bool
	GitVersion       string
	WSL              bool
}

// Detector provides platform detection functionality
type Detector struct {
	info *PlatformInfo
}

// NewDetector creates a new detector instance
func NewDetector() *Detector {
	return &Detector{
		info: &PlatformInfo{},
	}
}

// DetectAll runs all detection methods
func (d *Detector) DetectAll() (*PlatformInfo, error) {
	if err := d.detectOS(); err != nil {
		return nil, err
	}
	if err := d.detectShell(); err != nil {
		return nil, err
	}
	if err := d.detectPackageManager(); err != nil {
		return nil, err
	}
	if err := d.detectPython(); err != nil {
		return nil, err
	}
	if err := d.detectNode(); err != nil {
		return nil, err
	}
	if err := d.detectClaude(); err != nil {
		return nil, err
	}
	if err := d.detectLiteLLM(); err != nil {
		return nil, err
	}
	if err := d.detectDocker(); err != nil {
		return nil, err
	}
	if err := d.detectGit(); err != nil {
		return nil, err
	}
	d.detectWSL()

	return d.info, nil
}

// detectOS detects the operating system
func (d *Detector) detectOS() error {
	d.info.OS = runtime.GOOS
	d.info.Architecture = runtime.GOARCH

	// Try to get more specific Linux distribution info
	if d.info.OS == "linux" {
		if distro, err := detectLinuxDistribution(); err == nil {
			d.info.OS = distro
		}
	}

	return nil
}

// detectShell detects the current shell
func (d *Detector) detectShell() error {
	// Check common shell environment variables
	if shell := os.Getenv("SHELL"); shell != "" {
		d.info.Shell = filepath.Base(shell)
		return nil
	}

	// Check for shell process name
	if proc, err := os.ReadFile("/proc/self/comm"); err == nil {
		d.info.Shell = strings.TrimSpace(string(proc))
		return nil
	}

	// Fallback - try to detect from process tree
	if shell := detectShellFromProcess(); shell != "" {
		d.info.Shell = shell
		return nil
	}

	// Default
	d.info.Shell = "unknown"
	return nil
}

// detectPackageManager detects the package manager
func (d *Detector) detectPackageManager() error {
	switch {
	case isCommandAvailable("apt"):
		d.info.PackageManager = "apt"
	case isCommandAvailable("dnf"):
		d.info.PackageManager = "dnf"
	case isCommandAvailable("yum"):
		d.info.PackageManager = "yum"
	case isCommandAvailable("pacman"):
		d.info.PackageManager = "pacman"
	case isCommandAvailable("brew"):
		d.info.PackageManager = "brew"
	case isCommandAvailable("winget"):
		d.info.PackageManager = "winget"
	case isCommandAvailable("choco"):
		d.info.PackageManager = "choco"
	case isCommandAvailable("scoop"):
		d.info.PackageManager = "scoop"
	default:
		d.info.PackageManager = "unknown"
	}

	return nil
}

// detectPython detects Python installation
func (d *Detector) detectPython() error {
	// Check for python3 first
	if path, err := exec.LookPath("python3"); err == nil {
		d.info.PythonAvailable = true
		version, _ := getVersion(path, "--version")
		d.info.PythonVersion = strings.TrimSpace(strings.TrimPrefix(version, "Python "))
		return nil
	}

	// Fallback to python
	if path, err := exec.LookPath("python"); err == nil {
		d.info.PythonAvailable = true
		version, _ := getVersion(path, "--version")
		d.info.PythonVersion = strings.TrimSpace(strings.TrimPrefix(version, "Python "))
		return nil
	}

	d.info.PythonAvailable = false
	d.info.PythonVersion = ""
	return nil
}

// detectNode detects Node.js installation
func (d *Detector) detectNode() error {
	if path, err := exec.LookPath("node"); err == nil {
		d.info.NodeAvailable = true
		version, _ := getVersion(path, "--version")
		d.info.NodeVersion = strings.TrimSpace(strings.TrimPrefix(version, "v"))
		return nil
	}

	d.info.NodeAvailable = false
	d.info.NodeVersion = ""
	return nil
}

// detectClaude detects Claude Code installation
func (d *Detector) detectClaude() error {
	// Check various Claude Code installation paths
	prefixes := []string{
		"/usr/local/bin/claude",
		"/opt/claude",
		filepath.Join(os.Getenv("HOME"), ".local/bin/claude"),
	}

	for _, prefix := range prefixes {
		if info, err := os.Stat(prefix); err == nil && !info.IsDir() {
			d.info.ClaudeAvailable = true
			d.info.ClaudePath = prefix
			return nil
		}
	}

	// Check PATH
	if path, err := exec.LookPath("claude"); err == nil {
		d.info.ClaudeAvailable = true
		d.info.ClaudePath = path
		return nil
	}

	d.info.ClaudeAvailable = false
	d.info.ClaudePath = ""
	return nil
}

// detectLiteLLM detects LiteLLM installation
func (d *Detector) detectLiteLLM() error {
	// Check if pipx has litellm installed
	if path, err := exec.LookPath("litellm"); err == nil {
		d.info.LiteLLMAvailable = true
		version, _ := getVersion(path, "--version")
		d.info.LiteLLMVersion = strings.TrimSpace(version)
		return nil
	}

	// Check if litellm is installed in Python
	if path, err := exec.LookPath("python3"); err == nil {
		cmd := exec.Command(path, "-m", "litellm", "--version")
		if output, err := cmd.CombinedOutput(); err == nil {
			d.info.LiteLLMAvailable = true
			d.info.LiteLLMVersion = strings.TrimSpace(string(output))
			return nil
		}
	}

	d.info.LiteLLMAvailable = false
	d.info.LiteLLMVersion = ""
	return nil
}

// detectDocker detects Docker installation
func (d *Detector) detectDocker() error {
	if path, err := exec.LookPath("docker"); err == nil {
		d.info.DockerAvailable = true
		version, _ := getVersion(path, "--version")
		d.info.DockerVersion = strings.TrimSpace(extractDockerVersion(version))
		return nil
	}

	d.info.DockerAvailable = false
	d.info.DockerVersion = ""
	return nil
}

// detectGit detects Git installation
func (d *Detector) detectGit() error {
	if path, err := exec.LookPath("git"); err == nil {
		d.info.GitAvailable = true
		version, _ := getVersion(path, "--version")
		d.info.GitVersion = strings.TrimSpace(strings.TrimPrefix(version, "git version "))
		return nil
	}

	d.info.GitAvailable = false
	d.info.GitVersion = ""
	return nil
}

// detectWSL detects if running on WSL
func (d *Detector) detectWSL() {
	if !IsLinux() {
		d.info.WSL = false
		return
	}

	// Check for WSL kernel
	if data, err := os.ReadFile("/proc/version"); err == nil {
		if strings.Contains(strings.ToLower(string(data)), "microsoft") {
			d.info.WSL = true
			return
		}
	}

	// Check for WSL hostname pattern
	if hostname, err := os.Hostname(); err == nil {
		if strings.Contains(strings.ToLower(hostname), "wsl") {
			d.info.WSL = true
			return
		}
	}

	d.info.WSL = false
}

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// getVersion runs a command and extracts version information
func getVersion(cmd, flag string) (string, error) {
	fullCmd := exec.Command(cmd, flag)
	output, err := fullCmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// extractDockerVersion extracts version from docker --version output
func extractDockerVersion(output string) string {
	re := regexp.MustCompile(`Docker version ([\d.]+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1]
	}
	return output
}

// detectLinuxDistribution detects the Linux distribution
func detectLinuxDistribution() (string, error) {
	// Check /etc/os-release
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "ID=") {
				return strings.ToLower(strings.TrimSpace(strings.TrimPrefix(line, "ID="))), nil
			}
		}
	}

	return "linux", nil
}

// detectShellFromProcess detects shell from process information
func detectShellFromProcess() string {
	// Check shell environment
	if shell := os.Getenv("BASH_VERSION"); shell != "" {
		return "bash"
	}
	if shell := os.Getenv("ZSH_VERSION"); shell != "" {
		return "zsh"
	}
	if shell := os.Getenv("FISH_VERSION"); shell != "" {
		return "fish"
	}
	return ""
}

// IsLinux checks if running on Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsMacOS checks if running on macOS
func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// IsWindows checks if running on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// GetConfigDirectory gets the config directory path
func GetConfigDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch {
	case IsWindows():
		return filepath.Join(home, "AppData", "Roaming", "proxybridge"), nil
	case IsMacOS():
		return filepath.Join(home, "Library", "Application Support", "proxybridge"), nil
	default:
		return filepath.Join(home, ".config", "proxybridge"), nil
	}
}

// GetConfigPath gets the full path to the config file
func GetConfigPath() (string, error) {
	dir, err := GetConfigDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// GetLogDirectory gets the log directory path
func GetLogDirectory() (string, error) {
	dir, err := GetConfigDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "logs"), nil
}
