# Documentation for ProxyBridge

## Overview

ProxyBridge is a production-grade CLI that enables Claude Code to work with any LLM provider through LiteLLM. This documentation covers installation, configuration, and usage.

## Installation

### Prerequisites

- Python 3.8+: Required for LiteLLM
- Node.js 18+: Optional, for some integrations
- Go 1.23+: Required for building from source

### Quick Install

```bash
curl -fsSL https://proxybridge.com/install.sh | sh
```

### Manual Install

```bash
# Build from source
git clone https://github.com/proxybridge/cli.git
cd cli
go build -o proxybridge

# Install system-wide
sudo cp proxybridge /usr/local/bin/
```

## Configuration

### Default Setup

Run `proxybridge install` to automatically configure everything.

### Manual Configuration

Configuration files are stored in `~/.config/proxybridge/`:

- `config.yaml` - Main ProxyBridge configuration
- `providers.yaml` - Provider API keys and settings
- `aliases.yaml` - Model name aliases
- `environment` - Environment variables
- `litellm.yaml` - LiteLLM server configuration

### Editing Configuration

```bash
# View current configuration
proxybridge config show

# Edit configuration
proxybridge config edit
```

## Supported Providers

### Cloud Providers

| Provider | API Key URL | Default Model |
|----------|-------------|---------------|
| AWS Bedrock | Console | anthropic.claude-3-5-sonnet |
| OpenAI | platform.openai.com | gpt-4o |
| Azure OpenAI | Azure Portal | gpt-4o |
| Anthropic | console.anthropic.com | claude-3-5-sonnet |
| Google Gemini | aistudio.google.com | gemini-1.5-flash |
| OpenRouter | openrouter.ai | openai/gpt-4o |
| Groq | console.groq.com | llama3-8b-8192 |
| DeepSeek | platform.deepseek.com | deepseek-chat |
| Mistral | console.mistral.ai | mistral-large |
| Together AI | cloud.together.ai | togethercomputer/llama-3-70b-chat |
| Fireworks AI | fireworks.ai | accounts/fireworks/models/firefunction-v2 |
| Cerebras | cloud.cerebras.ai | cerebras-llama3.1-8b |
| xAI | console.x.ai | grok-beta |

### Local Providers

| Provider | Setup | Default Model |
|----------|-------|---------------|
| Ollama | ollama.com | default |
| LM Studio | lmstudio.ai | default |
| vLLM | docs.vllm.ai | default |
| TGI | huggingface.co | default |

## Model Aliases

ProxyBridge provides built-in aliases for popular models:

| Alias | Provider | Model |
|-------|----------|-------|
| `claude` | Anthropic | claude-3-5-sonnet |
| `claude-3-5-sonnet` | Anthropic | claude-3-5-sonnet-20241022 |
| `claude-opus` | Anthropic | claude-opus-4-20250502 |
| `claude-3-haiku` | Anthropic | claude-3-haiku-20240307 |
| `gpt-4o` | OpenAI | gpt-4o |
| `gpt-4o-mini` | OpenAI | gpt-4o-mini |
| `gpt-3.5-turbo` | OpenAI | gpt-3.5-turbo |
| `gemini` | Google | gemini-1.5-flash |
| `gemini-1.5-flash` | Google | gemini-1.5-flash |
| `gemini-1.5-pro` | Google | gemini-1.5-pro |
| `mistral` | Mistral | mistral-large |
| `llama` | Groq | llama3-8b-8192 |
| `llama-3` | Groq | llama3-8b-8192 |
| `llama-3.1-70b` | Groq | llama-3.1-70b-versatile |
| `codex` | Any | Claude 3.5 Sonnet |
| `reasoning` | Any | Claude 3.7 Sonnet |

### Adding Custom Aliases

```bash
proxybridge alias add my-model gpt-4o openai
```

### Switching Providers

```bash
# Switch default provider
proxybridge switch anthropic

# Switch specific alias
proxybridge alias switch claude openai
```

## Usage

### Basic Commands

```bash
# Install and configure
proxybridge install

# Check status
proxybridge status

# Run diagnostics
proxybridge doctor

# Start LiteLLM
proxybridge start

# Stop LiteLLM
proxybridge stop

# View logs
proxybridge logs
```

### Provider Management

```bash
# List all providers
proxybridge providers

# Test a provider
proxybridge test openai

# Add credentials
proxybridge credentials add openai

# List credentials
proxybridge credentials list
```

### Model Management

```bash
# List configured models
proxybridge models

# Test a model
proxybridge test-model claude-3-5-sonnet
```

### Configuration Export/Import

```bash
# Export configuration
proxybridge export backup.yaml

# Import configuration
proxybridge import backup.yaml
```

## Troubleshooting

### Common Issues

1. **LiteLLM not starting**
   - Check if Python is installed: `python3 --version`
   - Check LiteLLM installation: `pipx list | grep litellm`
   - View logs: `proxybridge logs`

2. **API key errors**
   - Verify credentials: `proxybridge credentials validate`
   - Re-add credentials: `proxybridge credentials add <provider>`

3. **Claude Code can't connect**
   - Ensure LiteLLM is running: `proxybridge start`
   - Check configuration: `proxybridge config show`
   - Validate installation: `proxybridge doctor`

4. **Model not found**
   - Check available models: `proxybridge models`
   - Verify alias: `proxybridge alias list`

### Getting Help

- Check logs: `proxybridge logs`
- Run diagnostics: `proxybridge doctor`
- Open GitHub issue

## Advanced Configuration

### Environment Variables

ProxyBridge uses the following environment variables:

| Variable | Description |
|----------|-------------|
| `LITELLM_HOST` | LiteLLM server address |
| `LITELLM_PORT` | LiteLLM server port |
| `CONFIG_DIR` | Custom config directory |

### Custom LiteLLM Config

```bash
# Edit LiteLLM config directly
proxybridge config edit
```

## API Reference

### LiteLLM API

ProxyBridge exposes a standard OpenAI-compatible API:

```
http://localhost:4000/v1/models
http://localhost:4000/v1/chat/completions
```

### Provider API Keys

API keys are stored securely in OS credential stores or encrypted config files.

## Security

- API keys use OS credential stores when available
- Encrypted fallback for local storage
- No secrets in logs
- Secure defaults for all configurations

## Uninstall

```bash
proxybridge uninstall
```

This removes all ProxyBridge configuration but preserves OS credential store entries.
