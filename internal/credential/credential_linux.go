// Package credential provides OS-specific credential storage implementations.
// This file contains Linux-specific credential store functions using Secret Service API.
//go:build !windows && !darwin
// +build !windows,!darwin

package credential

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// secretServiceStore stores a credential using Secret Service API (libsecret)
func secretServiceStore(provider, key, value string) (bool, error) {
	// Try to use secret-tool if available
	_, err := exec.LookPath("secret-tool")
	if err != nil {
		return false, nil // Fall back to file storage
	}

	label := fmt.Sprintf("ProxyBridge: %s %s", provider, key)
	cmd := exec.Command(
		"secret-tool",
		"store",
		"--label", label,
		"proxybridge", provider,
		"key", key,
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return false, err
	}

	go func() {
		defer stdin.Close()
		stdin.Write([]byte(value + "\n"))
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("secret-tool failed: %s", string(output))
	}

	return true, nil
}

// secretServiceGet retrieves a credential using Secret Service API
func secretServiceGet(provider, key string) (string, error) {
	_, err := exec.LookPath("secret-tool")
	if err != nil {
		return "", nil // Fall back to file storage
	}

	cmd := exec.Command(
		"secret-tool",
		"lookup",
		"proxybridge", provider,
		"key", key,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("secret-tool lookup failed: %s", string(output))
	}

	return strings.TrimSpace(string(output)), nil
}

// secretServiceDelete deletes a credential using Secret Service API
func secretServiceDelete(provider, key string) error {
	_, err := exec.LookPath("secret-tool")
	if err != nil {
		return nil // Fall back to file storage
	}

	cmd := exec.Command(
		"secret-tool",
		"clear",
		"proxybridge", provider,
		"key", key,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("secret-tool clear failed: %s", string(output))
	}

	return nil
}

// passStoreStore stores a credential using the pass password manager
func passStoreStore(provider, key, value string) (bool, error) {
	_, err := exec.LookPath("pass")
	if err != nil {
		return false, nil
	}

	// Check if pass init has been done
	home, _ := os.UserHomeDir()
	passDir := home + "/.password-store"
	if _, err := os.Stat(passDir); os.IsNotExist(err) {
		return false, nil
	}

	secretPath := fmt.Sprintf("proxybridge/%s/%s", provider, key)

	// Create directory structure if needed
	_ = fmt.Sprintf("proxybridge/%s", provider)
	cmd := exec.Command("pass", "init", "-p", provider)
	cmd.Dir = passDir
	if err := cmd.Run(); err != nil && !strings.Contains(string(err.Error()), "already exists") {
		return false, err
	}

	// Store the credential
	cmd = exec.Command("pass", "insert", "-m", "-f", secretPath)
	cmd.Stdin = strings.NewReader(value + "\n")
	cmd.Dir = passDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("pass insert failed: %s", string(output))
	}

	return true, nil
}

// passStoreGet retrieves a credential using the pass password manager
func passStoreGet(provider, key string) (string, error) {
	_, err := exec.LookPath("pass")
	if err != nil {
		return "", nil
	}

	home, _ := os.UserHomeDir()
	passDir := home + "/.password-store"
	secretPath := fmt.Sprintf("proxybridge/%s/%s", provider, key)

	cmd := exec.Command("pass", "show", secretPath)
	cmd.Dir = passDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pass show failed: %s", string(output))
	}

	return strings.TrimSpace(string(output)), nil
}

// passStoreDelete deletes a credential using the pass password manager
func passStoreDelete(provider, key string) error {
	_, err := exec.LookPath("pass")
	if err != nil {
		return nil
	}

	home, _ := os.UserHomeDir()
	passDir := home + "/.password-store"
	secretPath := fmt.Sprintf("proxybridge/%s/%s", provider, key)

	cmd := exec.Command("pass", "rm", "-f", secretPath)
	cmd.Dir = passDir
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "not exist") {
		return fmt.Errorf("pass rm failed: %s", string(output))
	}

	return nil
}
