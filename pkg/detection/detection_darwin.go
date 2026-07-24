// Package detection provides macOS-specific platform detection.
//go:build darwin
// +build darwin

package detection

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// detectMacOSPackageManager detects macOS package managers
func detectMacOSPackageManager() string {
	if isCommandAvailable("brew") {
		return "brew"
	}
	return "unknown"
}

// detectMacOSClaudePath detects Claude Code on macOS
func detectMacOSClaudePath() (string, error) {
	possiblePaths := []string{
		"/usr/local/bin/claude",
		"/opt/homebrew/bin/claude",
		filepath.Join(os.Getenv("HOME"), "Applications", "Claude.app", "Contents", "MacOS", "Claude"),
		filepath.Join(os.Getenv("HOME"), ".local", "bin", "claude"),
	}

	for _, path := range possiblePaths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, nil
		}
	}

	// Check PATH
	if path, err := exec.LookPath("claude"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("Claude Code not found")
}

// detectMacOSDockerPath detects Docker on macOS
func detectMacOSDockerPath() (string, error) {
	possiblePaths := []string{
		"/usr/local/bin/docker",
		"/opt/homebrew/bin/docker",
		"/Applications/Docker.app/Contents/Resources/bin/docker",
	}

	for _, path := range possiblePaths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, nil
		}
	}

	if path, err := exec.LookPath("docker"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("Docker not found")
}

// detectMacOSAppPath checks if an application is installed
func detectMacOSAppPath(appName string) string {
	possiblePaths := []string{
		filepath.Join("/Applications", appName),
		filepath.Join(os.Getenv("HOME"), "Applications", appName),
	}

	for _, path := range possiblePaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path
		}
	}
	return ""
}

// detectWSLOnMacOS detects if running on WSL from macOS
// This would only be relevant in cross-remote scenarios
func detectWSLOnMacOS() bool {
	// Check for WSL by looking at /proc/version
	if data, err := os.ReadFile("/proc/version"); err == nil {
		return strings.Contains(strings.ToLower(string(data)), "microsoft")
	}
	return false
}
