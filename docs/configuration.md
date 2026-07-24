# Configuration Guide

This guide covers ProxyBridge configuration in detail.

## Configuration Files

ProxyBridge stores configuration in `~/.config/proxybridge/`:

```
~/.config/proxybridge/
├── config.yaml          # Main configuration
├── providers.yaml       # Provider settings
├── aliases.yaml         # Model aliases
├── environment          # Environment variables
├── litellm.yaml         # LiteLLM configuration
├── .claude.json         # Claude Code config (backup)
├── backups/             # Configuration backups
│   └── 20240101-120000/
└── credentials/         # Credential references
```

## Main Configuration (config.yaml)

```yaml
version: "1.0.0"
created_at: "2024-01-01T00:00:00Z"
updated_at: "2024-01-01T00:00:00Z"
install_id: "uuid-here"
installation_type: "cli"

settings:
  log_level: info
  verbose: false
  json_output: false
  skip_verification: false
  debug: false
  auto_start: true
  service_type: process
  listen_address: 127.0.0.1
  port: 4000

active_provider: "anthropic"
default_model: "claude-3-5-sonnet"

litellm:
  path: "litellm"
  config_path: "/home/user/.config/proxybridge/litellm.yaml"
  environment_path: "/home/user/.config/proxybridge/environment"
  log_path: "/home/user/.config/proxybridge/litellm.log"
  pid_path: "/home/user/.config/proxybridge/litellm.pid"
  port: 4000

claude_code:
  path: "/usr/local/bin/claude"
  enabled: true
  backup_path: "/home/user/.config/proxybrain/.claude.json"
  auto_update: true
```

## Provider Configuration (providers.yaml)

```yaml
openai:
  api_key: "${OPENAI_API_KEY}"
  base_url: "https://api.openai.com/v1"
  models:
    - gpt-4o
    - gpt-4o-mini
    - gpt-4-turbo

anthropic:
  api_key: "${ANTHROPIC_API_KEY}"
  base_url: "https://api.anthropic.com/v1"
  models:
    - claude-3-5-sonnet-20241022
    - claude-3-opus-20240229

aws:
  access_key: "${AWS_ACCESS_KEY_ID}"
  secret_key: "${AWS_SECRET_ACCESS_KEY}"
  region: "us-east-1"
  models:
    - anthropic.claude-3-5-sonnet-20241022-v2:0
```

## Model Aliases (aliases.yaml)

```yaml
claude-3-5-sonnet:
  name: "claude-3-5-sonnet"
  model: "claude-3-5-sonnet-20241022"
  provider: "anthropic"
  enabled: true
  notes: "Anthropic Claude 3.5 Sonnet"

gpt-4o:
  name: "gpt-4o"
  model: "gpt-4o"
  provider: "openai"
  enabled: true
  notes: "OpenAI GPT-4o"

claude:
  name: "claude"
  model: "claude-3-5-sonnet-20241022"
  provider: "anthropic"
  enabled: true
  notes: "Default Claude model"
```

## LiteLLM Configuration (litellm.yaml)

```yaml
model_list:
  - model_name: default
    litellm_params:
      model: openai/default
      api_key: ${LITELLM_API_KEY}

  - model_name: claude-3-5-sonnet
    litellm_params:
      model: anthropic/claude-3-5-sonnet-20241022
      api_key: ${ANTHROPIC_API_KEY}

  - model_name: gpt-4o
    litellm_params:
      model: openai/gpt-4o
      api_key: ${OPENAI_API_KEY}

general_settings:
  master_key: ${LITELLM_MASTER_KEY}
  database_url: "file:/tmp/litellm.db"
  num_workers: 4
  cooldown_seconds: 30
  tracing: false
  alerting: []
  alert_types: []
```

## Environment Variables

```bash
# Provider API Keys
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
GOOGLE_API_KEY=AIza...

# LiteLLM Settings
LITELLM_HOST=127.0.0.1
LITELLM_PORT=4000
LITELLM_API_KEY=sk-...

# Claude Code Integration
ClaudeCodeEnabled=true
ClaudeCodeModel=claude-3-5-sonnet
```

## Configuration Management

### Viewing Configuration

```bash
# View current configuration
proxybridge config show

# View in JSON format
proxybridge config show --json
```

### Editing Configuration

```bash
# Edit configuration file
proxybridge config edit
```

### Export/Import Configuration

```bash
# Export to YAML
proxybridge export config.yaml

# Export to JSON
proxybridge export config.json

# Import configuration
proxybridge import config.yaml
```

## Advanced Configuration

### Custom LiteLLM Settings

```yaml
litellm_settings:
  drop_params: true
  add_function_to_prompt: false
  headers:
    Cache-Control: "no-cache"
    Content-Type: "application/json"
```

### Debug Mode

```bash
# Start with debug logging
proxybridge start --log-level debug

# Or set in config
settings:
  log_level: debug
```

### Change Port

```yaml
settings:
  port: 8000

litellm:
  port: 8000
```

## Environment-Specific Configuration

### Development

```yaml
settings:
  log_level: debug
  verbose: true
```

### Production

```yaml
settings:
  log_level: warn
  skip_verification: true
```

## Troubleshooting Configuration

### Validation Errors

```bash
# Validate configuration
proxybridge validate

# View validation report
proxybridge doctor
```

### Reset Configuration

```bash
# Create backup first
proxybridge export backup_$(date +%Y%m%d).yaml

# Then reinstall
proxybridge uninstall
proxybridge install
```

### Log Files

Logs are stored in `~/.config/proxybridge/logs/`:

```bash
# View LiteLLM logs
proxybridge logs

# View last N lines
proxybridge logs --lines 100
```
