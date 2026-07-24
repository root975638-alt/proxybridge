// Package detection provides Windows-specific platform detection.
//go:build windows
// +build windows

package detection

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// detectWindowsPackageManager detects Windows package managers
func detectWindowsPackageManager() string {
	// Check winget first
	if isCommandAvailable("winget") {
		return "winget"
	}

	// Check choco
	if isCommandAvailable("choco") {
		return "choco"
	}

	// Check scoop
	if isCommandAvailable("scoop") {
		return "scoop"
	}

	return "unknown"
}

// detectWindowsClaudePath detects Claude Code on Windows
func detectWindowsClaudePath() (string, error) {
	// Check AppData\Local\Programs\claude
	localAppData := os.Getenv("LOCALAPPDATA")
	possiblePaths := []string{
		filepath.Join(localAppData, "Programs", "claude", "claude.exe"),
		filepath.Join(localAppData, "claude", "claude.exe"),
		"C:\\Program Files\\claude\\claude.exe",
		"C:\\Program Files (x86)\\claude\\claude.exe",
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

// detectWindowsDockerPath detects Docker on Windows
func detectWindowsDockerPath() (string, error) {
	possiblePaths := []string{
		"C:\\Program Files\\Docker\\Docker\\docker.exe",
		"C:\\ProgramData\\DockerDesktop\\version-cli\\docker.exe",
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

// detectRegistryValue attempts to read a registry value
func detectRegistryValue(root registry.Key, path, name string) (string, bool) {
	key, err := registry.OpenKey(root, path, registry.QUERY_VALUE)
	if err != nil {
		return "", false
	}
	defer key.Close()

	value, _, err := key.GetStringValue(name)
	if err != nil {
		return "", false
	}

	return value, true
}

// detectWSLOnWindows detects if running on WSL
func detectWSLOnWindows() bool {
	// On Windows, check if we're in a WSL environment
	// WSL has /proc/version with "microsoft" in it
	if IsWSL() {
		return true
	}

	// Check for WSL2 by looking at /proc/sys/kernel/osrelease
	if data, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		return strings.Contains(strings.ToLower(string(data)), "microsoft")
	}

	return false
}
