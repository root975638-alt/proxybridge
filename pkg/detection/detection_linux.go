// Package detection provides Linux-specific platform detection.
//go:build linux
// +build linux

package detection

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// detectLinuxPackageManager detects Linux package managers
func detectLinuxPackageManager() string {
	if isCommandAvailable("apt") {
		return "apt"
	}
	if isCommandAvailable("dnf") {
		return "dnf"
	}
	if isCommandAvailable("yum") {
		return "yum"
	}
	if isCommandAvailable("pacman") {
		return "pacman"
	}
	return "unknown"
}

// detectLinuxClaudePath detects Claude Code on Linux
func detectLinuxClaudePath() (string, error) {
	possiblePaths := []string{
		"/usr/local/bin/claude",
		"/usr/bin/claude",
		filepath.Join(os.Getenv("HOME"), ".local", "bin", "claude"),
		filepath.Join(os.Getenv("HOME"), "bin", "claude"),
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

// detectLinuxDockerPath detects Docker on Linux
func detectLinuxDockerPath() (string, error) {
	possiblePaths := []string{
		"/usr/bin/docker",
		"/usr/local/bin/docker",
		filepath.Join(os.Getenv("HOME"), "bin", "docker"),
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

// detectLinuxServiceManager detects the service manager
func detectLinuxServiceManager() string {
	// Check for systemd
	if isCommandAvailable("systemctl") {
		return "systemd"
	}
	// Check for openrc
	if isCommandAvailable("rc-service") {
		return "openrc"
	}
	// Check for runit
	if isCommandAvailable("sv") {
		return "runit"
	}
	return "unknown"
}

// detectWSLFromLinux detects WSL from Linux
func detectWSLFromLinux() bool {
	// Check /proc/version for Microsoft kernel
	if data, err := os.ReadFile("/proc/version"); err == nil {
		if strings.Contains(strings.ToLower(string(data)), "microsoft") {
			return true
		}
	}

	// Check for WSL-specific files
	wslFiles := []string{
		"/proc/sys/kernel/osrelease",
		"/proc/version_signature",
	}

	for _, file := range wslFiles {
		if data, err := os.ReadFile(file); err == nil {
			if strings.Contains(strings.ToLower(string(data)), "microsoft") {
				return true
			}
		}
	}

	// Check hostname for WSL pattern
	if hostname, err := os.Hostname(); err == nil {
		if strings.Contains(strings.ToLower(hostname), "wsl") {
			return true
		}
	}

	return false
}

// detectUbuntuVersion returns the Ubuntu version
func detectUbuntuVersion() string {
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "VERSION_ID=") {
				return strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
			}
		}
	}
	return "unknown"
}

// detectDebianVersion returns the Debian version
func detectDebianVersion() string {
	if data, err := os.ReadFile("/etc/debian_version"); err == nil {
		return strings.TrimSpace(string(data))
	}
	return "unknown"
}

// detectFedoraVersion returns the Fedora version
func detectFedoraVersion() string {
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "VERSION_ID=") {
				return strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
			}
		}
	}
	return "unknown"
}

// detectArchVersion returns the Arch version
func detectArchVersion() string {
	if _, err := os.ReadFile("/etc/arch-release"); err == nil {
		return "rolling"
	}
	return "unknown"
}
