// Package credential handles secure credential storage and management.
// It supports OS credential stores, encrypted config files, and environment variables.
package credential

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/root975638-alt/proxybridge/internal/config"
	"github.com/root975638-alt/proxybridge/internal/logging"
)

const (
	// CredentialStoreFileName is the encrypted credential store file
	CredentialStoreFileName = "credentials.json"

	// CredentialStoreDirName is the credentials directory
	CredentialStoreDirName = "credentials"
)

// CredentialStore represents the in-memory credential storage
type CredentialStore struct {
	Credentials map[string]EncryptedCredential `json:"credentials"`
	Version     string                         `json:"version"`
	Encrypted   bool                           `json:"encrypted"`
}

// EncryptedCredential holds a credential with encryption info
type EncryptedCredential struct {
	Provider    string `json:"provider"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Encrypted   bool   `json:"encrypted"`
	EncryptedAt string `json:"encrypted_at"`
}

// Manager provides credential storage operations
type Manager struct {
	configDir string
	credentials *CredentialStore
	credentialsPath string
	logger    *logging.Logger
	masterKey []byte
}

// NewManager creates a new credential manager
func NewManager(logger *logging.Logger) (*Manager, error) {
	configDir, err := config.GetConfigDirectoryWithDefault()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	credentialsPath := filepath.Join(configDir, CredentialStoreDirName, CredentialStoreFileName)
	manager := &Manager{
		configDir:       configDir,
		credentialsPath: credentialsPath,
		logger:          logger,
		credentials: &CredentialStore{
			Credentials: make(map[string]EncryptedCredential),
			Version:     "1.0.0",
		},
	}

	// Try to load existing credentials
	if err := manager.loadCredentials(); err != nil {
		logger.Warn("Failed to load existing credentials, starting fresh")
	}

	return manager, nil
}

// loadCredentials loads credentials from disk
func (m *Manager) loadCredentials() error {
	data, err := os.ReadFile(m.credentialsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read credentials file: %w", err)
	}

	var store CredentialStore
	if err := json.Unmarshal(data, &store); err != nil {
		return fmt.Errorf("failed to parse credentials file: %w", err)
	}

	m.credentials = &store
	return nil
}

// saveCredentials saves credentials to disk
func (m *Manager) saveCredentials() error {
	dir := filepath.Dir(m.credentialsPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create credentials directory: %w", err)
	}

	data, err := json.MarshalIndent(m.credentials, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Encrypt if we have a master key
	if len(m.masterKey) > 0 {
		encrypted, err := m.encryptData(data)
		if err != nil {
			return fmt.Errorf("failed to encrypt credentials: %w", err)
		}
		data = encrypted
	}

	if err := os.WriteFile(m.credentialsPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

// StoreCredential stores a credential securely
func (m *Manager) StoreCredential(provider, key, value string) error {
	m.logger.Info("Storing credential", "provider", provider, "key", key)

	// Try OS credential store first
	if stored, err := m.storeInOSCredentialStore(provider, key, value); stored || err != nil {
		if err != nil {
			m.logger.Warn("OS credential store failed, falling back to file", "error", err)
		}
		return err
	}

	// Fall back to encrypted config
	encryptedValue, err := m.encryptValue(value)
	if err != nil {
		return fmt.Errorf("failed to encrypt credential: %w", err)
	}

	id := buildCredentialID(provider, key)
	m.credentials.Credentials[id] = EncryptedCredential{
		Provider:    provider,
		Key:         key,
		Value:       encryptedValue,
		Encrypted:   true,
		EncryptedAt: currentTimestamp(),
	}

	if err := m.saveCredentials(); err != nil {
		return err
	}

	m.logger.Info("Credential stored successfully")
	return nil
}

// GetCredential retrieves a credential
func (m *Manager) GetCredential(provider, key string) (string, error) {
	m.logger.Debug("Retrieving credential", "provider", provider, "key", key)

	// Try OS credential store first
	if value, err := m.getFromOSCredentialStore(provider, key); value != "" || err != nil {
		return value, err
	}

	// Fall back to encrypted config
	id := buildCredentialID(provider, key)
	cred, exists := m.credentials.Credentials[id]
	if !exists {
		return "", fmt.Errorf("credential not found: %s/%s", provider, key)
	}

	if !cred.Encrypted {
		return cred.Value, nil
	}

	return m.decryptValue(cred.Value)
}

// DeleteCredential removes a credential
func (m *Manager) DeleteCredential(provider, key string) error {
	m.logger.Info("Deleting credential", "provider", provider, "key", key)

	// Delete from OS credential store
	if err := m.deleteFromOSCredentialStore(provider, key); err != nil {
		m.logger.Warn("Failed to delete from OS credential store", "error", err)
	}

	// Delete from file
	id := buildCredentialID(provider, key)
	delete(m.credentials.Credentials, id)

	return m.saveCredentials()
}

// ValidateCredential validates a credential is working
func (m *Manager) ValidateCredential(provider string) error {
	m.logger.Info("Validating credential", "provider", provider)

	// Get the provider validator
	validator, err := GetProviderValidator(provider)
	if err != nil {
		return fmt.Errorf("unknown provider: %s", provider)
	}

	// Retrieve credentials and validate
	return validator(m)
}

// GetProviderValidator returns the validator function for a provider
func GetProviderValidator(provider string) (func(*Manager) error, error) {
	provider = strings.ToLower(provider)

	switch {
	case strings.Contains(provider, "openai"):
		return validateOpenAI, nil
	case strings.Contains(provider, "anthropic"):
		return validateAnthropic, nil
	case strings.Contains(provider, "aws") || strings.Contains(provider, "bedrock"):
		return validateAWS, nil
	case strings.Contains(provider, "google") || strings.Contains(provider, "gemini"):
		return validateGoogle, nil
	case strings.Contains(provider, "azure"):
		return validateAzure, nil
	case strings.Contains(provider, "openrouter"):
		return validateOpenRouter, nil
	case strings.Contains(provider, "groq"):
		return validateGroq, nil
	case strings.Contains(provider, "deepseek"):
		return validateDeepSeek, nil
	case strings.Contains(provider, "mistral"):
		return validateMistral, nil
	default:
		return nil, fmt.Errorf("no validator available forprovider: %s", provider)
	}
}

// validateOpenAI validates OpenAI credentials
func validateOpenAI(m *Manager) error {
	apiKey, err := m.GetCredential("openai", "api_key")
	if err != nil {
		return fmt.Errorf("missing OpenAI API key: %w", err)
	}
	if len(apiKey) < 10 {
		return fmt.Errorf("invalid OpenAI API key format")
	}
	return nil
}

// validateAnthropic validates Anthropic credentials
func validateAnthropic(m *Manager) error {
	apiKey, err := m.GetCredential("anthropic", "api_key")
	if err != nil {
		return fmt.Errorf("missing Anthropic API key: %w", err)
	}
	if !strings.HasPrefix(apiKey, "sk-ant-") {
		return fmt.Errorf("invalid Anthropic API key format")
	}
	return nil
}

// validateAWS validates AWS credentials
func validateAWS(m *Manager) error {
	accessKey, err := m.GetCredential("aws", "access_key")
	if err != nil {
		return fmt.Errorf("missing AWS access key: %w", err)
	}
	if len(accessKey) != 20 {
		return fmt.Errorf("invalid AWS access key length")
	}

	secretKey, err := m.GetCredential("aws", "secret_key")
	if err != nil {
		return fmt.Errorf("missing AWS secret key: %w", err)
	}
	if len(secretKey) != 40 {
		return fmt.Errorf("invalid AWS secret key length")
	}
	return nil
}

// validateGoogle validates Google credentials
func validateGoogle(m *Manager) error {
	_, err := m.GetCredential("google", "api_key")
	if err != nil {
		// Google uses service account JSON instead
		_, err := m.GetCredential("google", "service_account")
		if err != nil {
			return fmt.Errorf("missing Google API key or service account: %w", err)
		}
	}
	return nil
}

// validateAzure validates Azure credentials
func validateAzure(m *Manager) error {
	apiKey, err := m.GetCredential("azure", "api_key")
	if err != nil {
		return fmt.Errorf("missing Azure API key: %w", err)
	}
	if len(apiKey) < 10 {
		return fmt.Errorf("invalid Azure API key format")
	}
	return nil
}

// validateOpenRouter validates OpenRouter credentials
func validateOpenRouter(m *Manager) error {
	apiKey, err := m.GetCredential("openrouter", "api_key")
	if err != nil {
		return fmt.Errorf("missing OpenRouter API key: %w", err)
	}
	if len(apiKey) < 10 {
		return fmt.Errorf("invalid OpenRouter API key format")
	}
	return nil
}

// validateGroq validates Groq credentials
func validateGroq(m *Manager) error {
	apiKey, err := m.GetCredential("groq", "api_key")
	if err != nil {
		return fmt.Errorf("missing Groq API key: %w", err)
	}
	if len(apiKey) < 10 {
		return fmt.Errorf("invalid Groq API key format")
	}
	return nil
}

// validateDeepSeek validates DeepSeek credentials
func validateDeepSeek(m *Manager) error {
	apiKey, err := m.GetCredential("deepseek", "api_key")
	if err != nil {
		return fmt.Errorf("missing DeepSeek API key: %w", err)
	}
	if len(apiKey) < 10 {
		return fmt.Errorf("invalid DeepSeek API key format")
	}
	return nil
}

// validateMistral validates Mistral credentials
func validateMistral(m *Manager) error {
	apiKey, err := m.GetCredential("mistral", "api_key")
	if err != nil {
		return fmt.Errorf("missing Mistral API key: %w", err)
	}
	if len(apiKey) < 10 {
		return fmt.Errorf("invalid Mistral API key format")
	}
	return nil
}

// StoreEnvironmentVariable stores an environment variable
func (m *Manager) StoreEnvironmentVariable(key, value string) error {
	envPath := filepath.Join(m.configDir, config.EnvironmentFileName)

	// Read existing env file
	var envVars(map[string]string)
	if data, err := os.ReadFile(envPath); err == nil {
		envVars = parseEnvFile(string(data))
	}

	envVars[key] = value

	// Format as env file
	var buf bytes.Buffer
	for k, v := range envVars {
		fmt.Fprintf(&buf, "%s=%s\n", k, v)
	}

	if err := os.WriteFile(envPath, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write environment file: %w", err)
	}

	m.logger.Info("Environment variable stored", "key", key)
	return nil
}

// GetEnvironmentVariable retrieves an environment variable
func (m *Manager) GetEnvironmentVariable(key string) (string, error) {
	envPath := filepath.Join(m.configDir, config.EnvironmentFileName)

	data, err := os.ReadFile(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("environment file not found")
		}
		return "", fmt.Errorf("failed to read environment file: %w", err)
	}

	envVars := parseEnvFile(string(data))
	if value, ok := envVars[key]; ok {
		return value, nil
	}

	return "", fmt.Errorf("environment variable not found: %s", key)
}

// parseEnvFile parses a .env file
func parseEnvFile(content string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}

	return result
}

// encryptValue encrypts a credential value
func (m *Manager) encryptValue(value string) (string, error) {
	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %w", err)
	}

	// Pad the value
	paddedValue := padPKCS7([]byte(value), aes.BlockSize)

	// Create cipher
	block, err := aes.NewCipher(m.masterKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Encrypt
	ciphertext := make([]byte, len(paddedValue))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedValue)

	// Combine IV and ciphertext
	result := append(iv, ciphertext...)

	return base64.StdEncoding.EncodeToString(result), nil
}

// decryptValue decrypts a credential value
func (m *Manager) decryptValue(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("failed to decode value: %w", err)
	}

	if len(data) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext length")
	}

	// Extract IV
	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	// Create cipher
	block, err := aes.NewCipher(m.masterKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Decrypt
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Unpad
	plaintext, err := unpadPKCS7(ciphertext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// encryptData encrypts entire data blob
func (m *Manager) encryptData(data []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(m.masterKey)
	if err != nil {
		return nil, err
	}

	// PKCS7 padding
	padded := padPKCS7(data, aes.BlockSize)

	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	return append(iv, ciphertext...), nil
}

// decryptData decrypts data blob
func (m *Manager) decryptData(data []byte) ([]byte, error) {
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("invalid encrypted data")
	}

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	block, err := aes.NewCipher(m.masterKey)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return unpadPKCS7(ciphertext)
}

// padPKCS7 pads data to block size
func padPKCS7(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// unpadPKCS7 removes PKCS7 padding
func unpadPKCS7(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	padding := int(data[len(data)-1])
	if padding == 0 || padding > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}

	for i := 0; i < padding; i++ {
		if data[len(data)-1-i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padding], nil
}

// buildCredentialID builds a unique ID for a credential
func buildCredentialID(provider, key string) string {
	return fmt.Sprintf("%s/%s", provider, key)
}

// currentTimestamp returns current timestamp in ISO format
func currentTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// storeInOSCredentialStore attempts to store in OS credential store
func (m *Manager) storeInOSCredentialStore(provider, key, value string) (bool, error) {
	// Platform-specific implementation
	// For now, returns (false, nil) to use file-based fallback
	// Implement OS-specific stores in platform-specific files

	switch {
	case isWindows():
		// Use Windows Credential Manager
		// return windowsCredentialStore(provider, key, value)
	case isMacOS():
		// Use macOS Keychain
		// return macosKeychainStore(provider, key, value)
	default:
		// Use Secret Service API or pass
		// return secretServiceStore(provider, key, value)
	}

	return false, nil
}

// getFromOSCredentialStore attempts to retrieve from OS credential store
func (m *Manager) getFromOSCredentialStore(provider, key string) (string, error) {
	switch {
	case isWindows():
		// return windowsCredentialGet(provider, key)
	case isMacOS():
		// return macosKeychainGet(provider, key)
	default:
		// return secretServiceGet(provider, key)
	}
	return "", nil
}

// deleteFromOSCredentialStore attempts to delete from OS credential store
func (m *Manager) deleteFromOSCredentialStore(provider, key string) error {
	switch {
	case isWindows():
		// return windowsCredentialDelete(provider, key)
	case isMacOS():
		// return macosKeychainDelete(provider, key)
	default:
		// return secretServiceDelete(provider, key)
	}
	return nil
}

// isWindows checks if running on Windows
func isWindows() bool {
	return strings.EqualFold(os.Getenv("OS"), "Windows_NT")
}

// isMacOS checks if running on macOS
func isMacOS() bool {
	return strings.EqualFold(runtime.GOOS, "darwin")
}
