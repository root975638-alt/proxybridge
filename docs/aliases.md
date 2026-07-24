# Model Aliases Guide

Model aliases allow you to use memorable names that point to any backend provider.

## Built-In Aliases

ProxyBridge comes with many pre-configured aliases:

### Anthropic Claude

| Alias | Model | Provider |
|-------|-------|----------|
| `claude` | claude-3-5-sonnet-20241022 | Anthropic |
| `claude-3-5-sonnet` | claude-3-5-sonnet-20241022 | Anthropic |
| `claude-opus` | claude-opus-4-20250502 | Anthropic |
| `claude-3-haiku` | claude-3-haiku-20240307 | Anthropic |
| `claude-3-7-sonnet` | claude-3-7-sonnet-20250219 | Anthropic |
| `claude-3-5-haiku` | claude-3-5-haiku-20241022 | Anthropic |

### OpenAI GPT

| Alias | Model | Provider |
|-------|-------|----------|
| `gpt-4o` | gpt-4o | OpenAI |
| `gpt-4o-mini` | gpt-4o-mini | OpenAI |
| `gpt-4-turbo` | gpt-4-turbo | OpenAI |
| `gpt-4` | gpt-4 | OpenAI |
| `gpt-3.5-turbo` | gpt-3.5-turbo | OpenAI |
| `gpt` | gpt-4o | OpenAI |

### Google Gemini

| Alias | Model | Provider |
|-------|-------|----------|
| `gemini` | gemini-1.5-flash | Google |
| `gemini-1.5-flash` | gemini-1.5-flash | Google |
| `gemini-1.5-pro` | gemini-1.5-pro | Google |
| `gemini-2.0-flash` | gemini-2.0-flash-exp | Google |

### Groq LLaMA

| Alias | Model | Provider |
|-------|-------|----------|
| `llama` | llama3-8b-8192 | Groq |
| `llama-3` | llama3-8b-8192 | Groq |
| `llama-3.1-70b` | llama-3.1-70b-versatile | Groq |
| `llama-3.1-8b` | llama-3.1-8b-instant | Groq |
| `llama-2` | llama2-70b-4096 | Groq |
| `mixtral` | mixtral-8x7b-32768 | Groq |

### Other Providers

| Alias | Model | Provider |
|-------|-------|----------|
| `deepseek` | deepseek-chat | DeepSeek |
| `grok` | grok-beta | xAI |
| `grok-2` | grok-2-1212 | xAI |
| `qwen` | qwen-plus | Alibaba |
| `mistral` | mistral-large | Mistral |
| `bedrock` | anthropic.claude-3-5-sonnet-v2:0 | AWS |
| `azure` | gpt-4o | Azure |
| `vertex` | gemini-1.5-flash | Google |
| `openrouter` | openai/gpt-4o | OpenRouter |

## Custom Aliases

### Adding an Alias

```bash
# Interactive mode
proxybridge alias add my-model

# Direct add
proxybridge alias add deep-learning gpt-4o openai
```

### Listing Aliases

```bash
# List all aliases
proxybridge models

# Or use alias list
proxybridge alias list
```

### Switching a Provider for an Alias

```bash
# Switch alias to a different provider
proxybridge alias switch claude openai
```

### Removing an Alias

```bash
proxybridge alias remove old-alias
```

## Provider Switching

Switch all aliases to a new provider:

```bash
# Change default provider for all aliases
proxybridge switch anthropic
```

## Use Cases

### Development vs Production

```bash
# Use local models for development
proxybridge alias switch claude ollama

# Swap back to cloud for production
proxybridge alias switch claude anthropic
```

### Testing Different Models

```bash
# Create alias for testing
proxybridge alias add test-model gpt-4o-mini openai

# Use in Claude Code
proxybridge switch test-model
```

### Cost Optimization

```bash
# Use cheaper models for certain tasks
proxybridge alias add cheap-task gpt-4o-mini openai
proxybridge alias add expensive-task gpt-4o openai
```

## Troubleshooting

### Alias Not Found

```bash
# Check if alias exists
proxybridge models

# Alias list
proxybridge alias list
```

### Wrong Provider

```bash
# Check current mapping
proxybridge alias show claude

# Switch provider
proxybridge alias switch claude anthropic
```

### Command Not Found

If `proxybridge alias` commands fail, use the full path:

```bash
proxybridge alias <command>
```
