# ProxyBridge Architecture

## Overview

ProxyBridge is a production-grade CLI that serves as a universal proxy layer between Claude Code and any LLM provider through LiteLLM. This document describes the system architecture, module responsibilities, and design decisions.

## Core Principles

1. **Modularity**: Every concern is isolated in its own module
2. **Extensibility**: New providers require only adding a new plugin module
3. **Platform-aware**: Automatic OS, shell, and package manager detection
4. **Safe-by-default**: Never overwrite user config without backup
5. **Fail-explicit**: Clear diagnostics with actionable fixes
6. **Secure**: No secrets in config files, use OS credential stores where available
7. **Versioned**: All configuration is versioned with rollback support

## Architecture Layers

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CLI Interface Layer                          │
│  (cmd/proxybridge.go - Root command with all subcommands)           │
└─────────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  Installer      │  │  Diagnostics    │  │  Health Monitor   │
│  Engine         │  │  Engine         │  │  Engine         │
└─────────────────┘  └─────────────────┘  └─────────────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Engine Layer (Internal)                          │
│                                                                       │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌───────────┐  │
│  │ Config Eng   │ │ Credential   │ │ Provider     │ │ Alias     │  │
│  │              │ │ Manager      │ │ Plugin System│ │ System  │  │
│  └──────────────┘ └──────────────┘ └──────────────┘ └───────────┘  │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌───────────┐  │
│  │ LiteLLM Mgr  │ │ Claude Code  │ │ Logging      │ │ Export/   │  │
│  │              │ │ Integration  │ │ Engine       │ │ Import    │  │
│  └──────────────┘ └──────────────┘ └──────────────┘ └───────────┘  │
│                                                                       │
└─────────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ Dependency Mgr  │  │ Detection Engine│  │ Template Engine │
│ (OS/Shell/PM)   │  │ (Platform Info) │  │ (Config Gen)    │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

## Module Details

### 1. CLI Interface (`cmd/`)

**Responsibility**: Command parsing, flag handling, output formatting, exit codes

**Components**:
- `proxybridge.go` - Root command, version, help
- `commands/` - Individual command implementations
  - `install.go` - Install command
  - `uninstall.go` - Uninstall command
  - `repair.go` - Repair command
  - `update.go` - Update command
  - `doctor.go` - Full diagnostic check
  - `status.go` - System status
  - `start.go`, `stop.go`, `restart.go` - LiteLLM lifecycle
  - `logs.go` - Tail LiteLLM logs
  - `validate.go` - Configuration validation
  - `test.go` - Connection tests
  - `models.go` - List configured models
  - `providers.go` - List available providers
  - `credentials.go` - Credential management
  - `export.go`, `import.go` - Config migration
  - `config.go` - Show/edit config
  - `switch.go` - Default provider switch
  - `alias.go` - Model alias management
  - `benchmark.go` - Performance testing
  - `version.go` - Version info
  - `self-update.go` - Auto-updater

**Output Modes**:
- `human` - Colored, friendly output (default)
- `json` - Structured JSON for automation
- `quiet` - Only errors

### 2. Configuration Engine (`internal/config/`)

**Responsibility**: Generate, manage, and validate all configuration files

**Components**:
- `config.go` - Core config structures
- `litellm.go` - LiteLLM configuration generation
- `claude.go` - Claude Code configuration
- `environment.go` - Environment variable management
- `aliases.go` - Model alias definitions
- `routing.go` - Routing rules generation
- `validation.go` - Configuration validation

**Config Files Generated**:
- `~/.config/proxybridge/config.yaml` - Main config
- `~/.config/proxybridge/providers.yaml` - Provider mappings
- `~/.config/proxybridge/aliases.yaml` - Model aliases
- `~/.config/proxybridge/environment` - Environment file
- `~/.config/proxybridge/litellm.yaml` - LiteLLM config
- `~/.config/proxybridge/.claude.json` - Claude Code config backup

### 3. Credential Manager (`internal/credential/`)

**Responsibility**: Secure credential storage and retrieval

**Storage Options**:
- **OS Credential Store** (best effort):
  - Linux: Secret Service API (via libsecret) or pass
  - macOS: Keychain
  - Windows: Credential Manager
  - WSL2: Windows Credential Manager via wsldi
- **Encrypted Config** (fallback):
  - AES-256-GCM encryption
  - Key derived from master password
- **Environment Variables** (explicit choice)

**Functions**:
- `StoreCredential(provider, key, secret)` - Store securely
- `GetCredential(provider, key)` - Retrieve
- `DeleteCredential(provider, key)` - Remove
- `ValidateCredential(provider)` - Check validity
- `RotateCredential(provider)` - Update securely

### 4. Provider Plugin System (`internal/provider/`)

**Interface**: `Provider.go`

```go
type Provider interface {
    Name() string
    ID() string
    Setup() error
    Validate() (bool, error)
    GetModels() ([]string, error)
    GetDefaultModel() string
    GetConfig() map[string]string
    GetEnvironment() map[string]string
    IsInstalled() bool
    Install() error
    Uninstall() error
}
```

**First-Party Providers**:
- `aws/bedrock.go` - AWS Bedrock
- `openai.go` - OpenAI
- `azure/openai.go` - Azure OpenAI
- `anthropic.go` - Anthropic
- `google/gemini.go` - Google Gemini
- `openrouter.go` - OpenRouter
- `ollama.go` - Ollama
- `groq.go` - Groq
- `deepseek.go` - DeepSeek
- `mistral.go` - Mistral
- `together.go` - Together AI
- `fireworks.go` - Fireworks AI
- `cerebras.go` - Cerebras
- `xai/grok.go` - xAI (Grok)
- `lmstudio.go` - LM Studio
- `local.go` - Any OpenAI-compatible server

### 5. LiteLLM Manager (`pkg/manager/litellm.go`)

**Responsibility**: LiteLLM lifecycle management

**Operations**:
- `Install()` - pipx or pip install
- `Uninstall()` - Remove installation
- `Start()` - Start LiteLLM service
- `Stop()` - Stop LiteLLM
- `Restart()` - Restart LiteLLM
- `Status()` - Check running status
- `Logs()` - Read logs (file and stdout)
- `Validate()` - Validate config
- `Reload()` - Reload configuration without restart
- `HealthCheck()` - API health check

**Service Managers**:
- Systemd (Linux)
- Launchd (macOS)
- Windows Service (Windows)
- Plain process (fallback)

### 6. Claude Code Integration (`pkg/manager/claude.go`)

**Responsibility**: Claude Code configuration

**Operations**:
- `Detect()` - Find Claude Code installation
- `GetConfigPath()` - Resolve config location
- `ReadConfig()` - Load current config
- `WriteConfig()` - Save with backup
- `SetModelProvider()` - Configure model
- `GetAvailableModels()` - List configured models
- `CreateBackups()` - Backup before changes
- `RestoreBackup()` - Rollback if needed

### 7. Dependency Manager (`pkg/detection/dependency.go`)

**Responsibility**: Detect and install dependencies

**OS Detection**:
- Linux distros (Ubuntu, Debian, Mint, Fedora, Arch)
- macOS
- Windows
- WSL2

**Package Manager Detection**:
- `apt` - Debian/Ubuntu
- `dnf` - Fedora
- `yum` - RHEL/CentOS
- `pacman` - Arch
- `brew` / `brew cask` - macOS
- `winget` - Windows
- `choco` - Windows
- `scoop` - Windows
- `pip` / `pipx` - Python
- `npm` / `npx` - Node

### 8. Detection Engine (`pkg/detection/`)

**Responsibility**: Platform and tool detection

**Detections**:
- `os.go` - Operating system detection
- `shell.go` - Shell detection (bash, zsh, fish, powershell)
- `package_manager.go` - Package manager detection
- `python.go` - Python detection and version
- `node.go` - Node detection and version
- `claude.go` - Claude Code detection
- `litellm.go` - LiteLLM detection
- `docker.go` - Docker detection
- `git.go` - Git detection

### 9. Template Engine (`pkg/template/`)

**Responsibility**: Configuration file generation

**Templates**:
- `litellm.yaml.tmpl` - LiteLLM config template
- `environment.sh.tmpl` - Environment file
- `service.systemd.tmpl` - Systemd service
- `service.launchd.tmpl` - Launchd service
- `service.windows.tmpl` - Windows service script
- `start.script.tmpl` - Platform start script
- `claude.json.tmpl` - Claude Code config

### 10. Diagnostic Engine (`internal/diagnostic/`)

**Responsibility**: Issue detection and remediation

**Checks**:
- `api_key.go` - API key validation
- `region.go` - Region configuration
- `bedrock.go` - AWS Bedrock access
- `litellm_config.go` - LiteLLM config validation
- `claude_config.go` - Claude Code config
- `port.go` - Port conflict detection
- `firewall.go` - Firewall rules
- `dependencies.go` - Missing dependencies
- `permissions.go` - File permissions
- `proxy.go` - Proxy configuration
- `network.go` - Network connectivity
- `mismatch.go` - Configuration mismatches

### 11. Health Monitor (`internal/health/`)

**Responsibility**: Continuous health verification

**Checks**:
- `provider.go` - Provider connectivity
- `litellm.go` - LiteLLM API health
- `claude.go` - Claude Code integration
- `credential.go` - Credential validity
- `alias.go` - Alias configuration
- `config.go` - Config integrity
- `performance.go` - Performance metrics
- `latency.go` - Latency tracking

### 12. Logging Engine (`internal/logging/`)

**Responsibility**: Structured logging

**Features**:
- `log.go` - Core logging
- `formatter.go` - Human/JSON output
- `rotation.go` - Log rotation
- `export.go` - Log export

**Log Levels**:
- `debug` - Detailed debugging
- `info` - General information
- `warn` - Warnings
- `error` - Errors
- `fatal` - Fatal errors

### 13. Alias System (`internal/alias/`)

**Responsibility**: Model alias management

**Features**:
- `alias.go` - Core alias definitions
- `mapping.go` - Alias-to-model mapping
- `switch.go` - Provider switching via aliases
- `validate.go` - Alias validation

**Predefined Aliases**:
- `claude-3-5-sonnet` - Anthropic sonnet
- `claude-opus` - Anthropic opus
- `coder` - Best coding model
- `reasoning` - Best reasoning model
- `deepseek` - DeepSeek V3
- `grok` - xAI grok
- `qwen` - Qwen models
- `gemini` - Gemini models
- `gpt` - GPT-4o

### 14. Export/Import System (`internal/exportimport/`)

**Responsibility**: Configuration migration

**Features**:
- `export.go` - Export configuration
- `import.go` - Import configuration
- `backup.go` - Backup creation/restore
- `migrate.go` - Version migration

## Data Flow

### Installation Flow

```
User runs: proxybridge install

1. Detection Phase
   ├─ Detect OS, shell, package managers
   ├─ Detect existing installations (Claude, LiteLLM, Python, Node)
   └─ Detect available providers

2. Dependency Installation
   ├─ Install LiteLLM
   ├─ Install Python dependencies
   └─ Install system dependencies

3. Configuration Generation
   ├─ Generate LiteLLM config
   ├─ Generate environment file
   ├─ Generate Claude Code config
   └─ Generate aliases

4. Credential Setup
   ├─ Prompt for API keys
   ├─ Store securely (OS credential store first)
   └─ Validate credentials

5. Service Setup
   ├─ Create systemd/launchd service
   ├─ Configure environment
   └─ Start LiteLLM

6. Validation
   ├─ Test LiteLLM connectivity
   ├─ Test Claude Code integration
   └─ Final health check
```

### Request Flow (Claude → LiteLLM → Provider)

```
Claude Code
    │
    │ 1. Request to /v1/chat/completions
    │
    ▼
LiteLLM (proxy)
    │
    │ 2. Parse model name, resolve alias
    │ 3. Look up provider credentials
    │ 4. Route to correct provider
    │
    ▼
Provider API (OpenAI, Anthropic, etc.)
    │
    │ 5. Response
    │
    ▼
LiteLLM (format response)
    │
    │ 6. Standard OpenAI format
    │
    ▼
Claude Code
```

## Error Handling

**Exit Codes**:
- `0` - Success
- `1` - General error
- `2` - Configuration error
- `3` - Invalid credentials
- `4` - Network error
- `5` - Permission denied
- `6` - Service not running
- `7` - Invalid command

**Error Format**:
```json
{
  "error": {
    "code": "ERR_NAME",
    "message": "Human-readable message",
    "details": "Additional details",
    "suggestion": "Actionable fix"
  }
}
```

## Security Considerations

1. **No secrets in logs**: All logging sanitizes secrets
2. **No secrets in config**: Credentials stored in OS vault
3. **Encrypted fallback**: AES-256-GCM for encrypted config
4. **TLS validation**: Always verify certificates
5. **Least privilege**: Run with minimal permissions
6. **Secure defaults**: All defaults are secure

## Configuration Storage

```
~/.config/proxybridge/
├── config.yaml           # Main config (no secrets)
├── providers.yaml        # Provider mappings (no secrets)
├── aliases.yaml          # Model aliases (no secrets)
├── environment           # Environment variables (encrypted if needed)
├── litellm.yaml          # LiteLLM config (generated)
├── .claude.json          # Claude Code config backup
├── backups/              # Configuration backups
│   ├── 20240101-120000/
│   └── ...
├── credentials/          # OS credential references
├── logs/                 # Log files
├── services/             # Service definitions
└── state.json            # Runtime state
```

## Extension Points

### Adding a New Provider

1. Create `internal/provider/newprovider.go`
2. Implement `Provider` interface
3. Register in `internal/provider/registry.go`
4. Add tests in `internal/provider/newprovider_test.go`

### Adding a New CLI Command

1. Create `cmd/commands/newcommand.go`
2. Implement command logic
3. Register in `cmd/proxybridge.go`

### Adding a New Diagnostic Check

1. Create `internal/diagnostic/newcheck.go`
2. Implement check logic
3. Register in `internal/diagnostic/registry.go`

## Testing Strategy

1. **Unit Tests**: Per-package tests
2. **Integration Tests**: End-to-end scenarios
3. **Cross-Platform Tests**: All OSes
4. **Provider Mock Tests**: Simulate provider APIs
5. **CI Automation**: GitHub Actions

## Future Extensibility

- Webhook notifications
- Usage analytics
- Cost tracking
- Rate limit management
- Request retry policies
- Circuit breakers
- Request/response transformation plugins
