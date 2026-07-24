# Plugin Development Guide

Extend ProxyBridge with custom providers.

## Provider Interface

All providers must implement the `Provider` interface:

```go
type Provider interface {
    Name() string
    ID() string
    Setup() error
    Validate() error
    GetModels() ([]string, error)
    GetDefaultModel() string
    GetConfig() map[string]string
    GetEnvironment() map[string]string
    IsInstalled() bool
    Install() error
    Uninstall() error
}
```

## Provider Template

```go
package myprovider

import (
    "fmt"
    "os"
)

// MyProvider implements the Provider interface
type MyProvider struct {
    id           string
    name         string
    description  string
    defaultModel string
}

// NewMyProvider creates a new provider instance
func NewMyProvider() *MyProvider {
    return &MyProvider{
        id:           "myprovider",
        name:         "My Provider",
        description:  "My Custom Provider",
        defaultModel: "my-default-model",
    }
}

// Name returns the display name
func (p *MyProvider) Name() string {
    return p.name
}

// ID returns the unique identifier
func (p *MyProvider) ID() string {
    return p.id
}

// Setup performs initial setup
func (p *MyProvider) Setup() error {
    fmt.Println("Setting up My Provider...")
    return nil
}

// Validate validates credentials and configuration
func (p *MyProvider) Validate() error {
    apiKey := os.Getenv("MY_PROVIDER_API_KEY")
    if apiKey == "" {
        return fmt.Errorf("MY_PROVIDER_API_KEY not set")
    }
    return nil
}

// GetModels returns available models
func (p *MyProvider) GetModels() ([]string, error) {
    return []string{"my-model-1", "my-model-2"}, nil
}

// GetDefaultModel returns default model
func (p *MyProvider) GetDefaultModel() string {
    return p.defaultModel
}

// GetConfig returns provider config
func (p *MyProvider) GetConfig() map[string]string {
    return map[string]string{
        "provider": "myprovider",
        "type":     "custom",
    }
}

// GetEnvironment returns environment variables
func (p *MyProvider) GetEnvironment() map[string]string {
    return map[string]string{
        "MY_PROVIDER_API_KEY":   "YOUR_API_KEY",
        "MY_PROVIDER_BASE_URL":  "https://api.myprovider.com/v1",
    }
}

// IsInstalled checks if provider is installed
func (p *MyProvider) IsInstalled() bool {
    return true
}

// Install installs the provider
func (p *MyProvider) Install() error {
    fmt.Println("Install My Provider...")
    return nil
}

// Uninstall uninstalls the provider
func (p *MyProvider) Uninstall() error {
    return nil
}

// GetProvider returns a new provider instance
func GetProvider() *MyProvider {
    return NewMyProvider()
}
```

## Registering the Provider

Add to `internal/provider/registry.go`:

```go
import "github.com/root975638-alt/proxybridge/internal/provider/myprovider"

// In registerProvidersInRegistry():
_ = reg.RegisterProvider(myprovider.GetProvider())
```

## Test Template

```go
package myprovider

import (
    "os"
    "testing"
)

func TestMyProvider_Name(t *testing.T) {
    p := NewMyProvider()
    if p.Name() != "My Provider" {
        t.Errorf("expected name 'My Provider', got '%s'", p.Name())
    }
}

func TestMyProvider_ID(t *testing.T) {
    p := NewMyProvider()
    if p.ID() != "myprovider" {
        t.Errorf("expected ID 'myprovider', got '%s'", p.ID())
    }
}

func TestMyProvider_Validate(t *testing.T) {
    os.Setenv("MY_PROVIDER_API_KEY", "test-key")
    defer os.Unsetenv("MY_PROVIDER_API_KEY")
    
    p := NewMyProvider()
    if err := p.Validate(); err != nil {
        t.Errorf("Validate returned error: %v", err)
    }
}

func TestMyProvider_ValidateMissingKey(t *testing.T) {
    os.Unsetenv("MY_PROVIDER_API_KEY")
    
    p := NewMyProvider()
    err := p.Validate()
    if err == nil {
        t.Error("expected validation to fail when API key is missing")
    }
}
```

## Adding to Provider Registry

Update `internal/provider/provider.go`:

```go
// Register in registerDefaultProviders():
_ = reg.RegisterProvider(myprovider.GetProvider())
```

## Deployment

### Local Development

```bash
# Add your provider to the registry
# Build and test
go build -o proxybridge
go test ./...
```

### Distribution

1. Submit a pull request to the main repository
2. Or distribute as a standalone binary
3. Or package for your platform (deb, rpm, brew, etc.)

## Best Practices

1. **Name Conventions**:
   - Provider ID: lowercase, no spaces
   - Environment variables: `PROVIDER_NAME_API_KEY`

2. **Error Handling**:
   - Always return descriptive errors
   - Include actionable suggestions

3. **Documentation**:
   - Document all public methods
   - Provide example configuration

4. **Testing**:
   - Test all public methods
   - Test with and without API keys
   - Test error conditions

5. **Security**:
   - Never log API keys
   - Use OS credential store when available
   - Validate all user input

## Examples

See existing providers for examples:
- `internal/provider/aws/bedrock.go`
- `internal/provider/openai.go`
- `internal/provider/anthropic.go`
- `internal/provider/ollama.go`

## Support

For plugin-related questions:
- Open GitHub issue
- Check existing provider implementations
- Review test files for examples
