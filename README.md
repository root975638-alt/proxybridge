# ProxyBridge

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.23-green.svg)](https://golang.org)

ProxyBridge is a production-grade CLI that enables Claude Code to work with any LLM provider through LiteLLM, providing a universal interface and model aliasing system.

## Overview

ProxyBridge acts as a universal proxy between Claude Code and any LLM provider. It uses LiteLLM as the underlying proxy server to route requests to the appropriate provider based on the model name or alias.

**Key Features:**
- Universal LLM proxy through LiteLLM
- Support for 20+ LLM providers (OpenAI, Anthropic, AWS Bedrock, Groq, DeepSeek, Mistral, and more)
- Model aliasing system (use memorable names like `claude`, `gpt-4o`, `gemini`)
- Secure credential storage (OS vault + encrypted fallback)
- Cross-platform: Windows, macOS, Linux, WSL2
- Provider plugin architecture
- Built-in diagnostics and health monitoring

## Installation

### Prerequisites

- **Go 1.23+**: Required for building
- **Python 3.8+**: Required for LiteLLM
- **pipx** (recommended): For LiteLLM installation
- **apt/dnf/yum/pacman/brew/winget**: For system package management

### Quick Install

```bash
# Clone the repository
git clone https://github.com/root975638-alt/proxybridge.git
cd proxybridge

# Build the binary
go build -o proxybridge ./cmd/proxybridge.go

# Install to system path (requires sudo)
sudo cp proxybridge /usr/local/bin/
```

### Manual Installation

```bash
# Install LiteLLM
pipx install litellm

# Clone and build
git clone https://github.com/root975638-alt/proxybridge.git
cd proxybridge
go build -o proxybridge ./cmd/proxybridge.go
sudo cp proxybridge /usr/local/bin/
```

## Usage

### Install and Configure

```bash
# Run the installation wizard
proxybridge install

# This will:
# - Install LiteLLM (via pipx)
# - Detect available providers
# - Generate configuration
# - Set up credentials
```

### Commands

| Command | Description |
|---------|-------------|
| `install` | Install ProxyBridge and configure LiteLLM |
| `uninstall` | Remove ProxyBridge from the system |
| `repair` | Repair the ProxyBridge installation |
| `update` | Update ProxyBridge to latest version |
| `doctor` | Run comprehensive diagnostic checks |
| `status` | Show system status |
| `start` | Start LiteLLM service |
| `stop` | Stop LiteLLM service |
| `restart` | Restart LiteLLM service |
| `logs` | View LiteLLM logs |
| `validate` | Validate configuration |
| `test <provider>` | Test provider connection |
| `models` | List configured model aliases |
| `providers` | List available providers |
| `credentials <action>` | Manage credentials |
| `export <file>` | Export configuration |
| `import <file>` | Import configuration |
| `config <action>` | Show or edit configuration |
| `switch <provider>` | Switch default provider |
| `alias <action>` | Manage model aliases |
| `benchmark` | Run benchmark tests |
| `version` | Show version information |
| `self-update` | Update ProxyBridge |

### Examples

```bash
# Start LiteLLM
proxybridge start

# Check system status
proxybridge status

# Run diagnostics
proxybridge doctor

# List available providers
proxybridge providers

# Add credentials for a provider
proxybridge credentials add openai

# Test a provider connection
proxybridge test openai

# Export configuration
proxybridge export backup.yaml

# Switch default provider
proxybridge switch anthropic
```

## Supported Providers

### Cloud Providers

| Provider | Description | API Key URL |
|----------|-------------|-------------|
| AWS Bedrock | Amazon Bedrock models | AWS Console |
| OpenAI | GPT-4, GPT-3.5, O series | platform.openai.com |
| Azure OpenAI | Azure OpenAI Service | Azure Portal |
| Anthropic | Claude models | console.anthropic.com |
| Google Gemini | Gemini Flash/Pro models | aistudio.google.com |
| OpenRouter | Many providers via one API | openrouter.ai |
| Groq | Llama models with fast inference | console.groq.com |
| DeepSeek | DeepSeek Chat/Coder/Reasoner | platform.deepseek.com |
| Mistral | Mistral Large/Small/Tiny | console.mistral.ai |
| Together AI | Open-source models | cloud.together.ai |
| Fireworks AI | Fast inference | fireworks.ai |
| Cerebras | Llama 3.1 models | cloud.cerebras.ai |
| xAI | Grok models | console.x.ai |

### Local Providers

| Provider | Description | Setup |
|----------|-------------|-------|
| Ollama | Local LLM server | ollama.com |
| LM Studio | GUI LLM server | lmstudio.ai |
| vLLM | High-performance inference | docs.vllm.ai |
| TGI | HuggingFace inference | huggingface.co |

## Model Aliases

ProxyBridge includes built-in aliases for popular models:

| Alias | Model | Provider |
|-------|-------|----------|
| `claude` | claude-3-5-sonnet-20241022 | Anthropic |
| `claude-3-5-sonnet` | claude-3-5-sonnet-20241022 | Anthropic |
| `claude-opus` | claude-opus-4-20250502 | Anthropic |
| `claude-3-haiku` | claude-3-haiku-20240307 | Anthropic |
| `gpt-4o` | gpt-4o | OpenAI |
| `gpt-4o-mini` | gpt-4o-mini | OpenAI |
| `gpt-3.5-turbo` | gpt-3.5-turbo | OpenAI |
| `gemini` | gemini-1.5-flash | Google |
| `gemini-1.5-flash` | gemini-1.5-flash | Google |
| `gemini-1.5-pro` | gemini-1.5-pro | Google |
| `mistral` | mistral-large | Mistral |
| `llama` | llama3-8b-8192 | Groq |
| `llama-3.1-70b` | llama-3.1-70b-versatile | Groq |
| `deepseek` | deepseek-chat | DeepSeek |
| `grok` | grok-beta | xAI |
| `coder` | claude-3-5-sonnet-20241022 | Anthropic |
| `reasoning` | claude-3-7-sonnet-20250219 | Anthropic |

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       ProxyBridge CLI                        │
│  (CLI, Diagnostics, Health, Installer, Provider Plugins)    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       LiteLLM Proxy                          │
│  (Routes requests to appropriate providers)                  │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  OpenAI      │     │  Anthropic   │     │   AWS etc.   │
└──────────────┘     └──────────────┘     └──────────────┘
```

## Configuration

Configuration files are stored in `~/.config/proxybridge/`:

```
~/.config/proxybridge/
├── config.yaml          # Main ProxyBridge configuration
├── providers.yaml       # Provider settings
├── aliases.yaml         # Model aliases
├── environment          # Environment variables
├── litellm.yaml         # LiteLLM configuration
├── .claude.json         # Claude Code config (backup)
├── backups/             # Configuration backups
└── credentials/         # Credential references
```

## Security

- **OS Credential Stores**: Uses Keychain (macOS), Credential Manager (Windows), Secret Service (Linux)
- **Encrypted Fallback**: AES-256-GCM encryption when OS vault unavailable
- **No Secrets in Logs**: All logging sanitizes sensitive data
- **Secure Defaults**: TLS verification, least privilege design

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development Setup

```bash
git clone https://github.com/root975638-alt/proxybridge.git
cd proxybridge
go mod download
go test ./...
go build -o proxybridge
```

### Adding a New Provider

1. Create `internal/provider/newprovider.go`
2. Implement the `Provider` interface
3. Add tests in `internal/provider/newprovider_test.go`
4. Register in `internal/provider/registry.go`

See [docs/plugins.md](docs/plugins.md) for detailed Plugin Development Guide.

## License

MIT License

Copyright (c) 2026 Anand Damdiyal

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Acknowledgments

- [LiteLLM](https://github.com/BerriAI/litellm) - The excellent proxy server that powers ProxyBridge
- All contributors and users of ProxyBridge

## Repository

**Repository**: https://github.com/root975638-alt/proxybridge  
**Owner**: Anand Damdiyal  
**Issues**: https://github.com/root975638-alt/proxybridge/issues
