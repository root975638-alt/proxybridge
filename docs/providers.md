# Provider Guide

This guide covers all supported providers and how to configure them.

## Cloud Providers

### AWS Bedrock

AWS Bedrock provides access to Anthropic Claude, Amazon Titan, and other models.

#### Setup

1. Create AWS IAM user with Bedrock access
2. Configure AWS credentials:
   ```bash
   aws configure
   ```
3. Enable Bedrock access in AWS Console

#### Environment Variables

```bash
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1  # Bedrock available regions
```

#### Models

- `anthropic.claude-3-5-sonnet-20241022-v2:0`
- `anthropic.claude-3-opus-20240229-v1:0`
- `amazon.nova-pro-v1:0`
- `meta.llama3-1-70b-instruct-v1:0`

### OpenAI

Standard OpenAI API access.

#### Setup

1. Get API key from https://platform.openai.com/api-keys
2. Add credentials:
   ```bash
   proxybridge credentials add openai
   ```

#### Environment Variables

```bash
OPENAI_API_KEY=sk-...
OPENAI_ORG_ID=your_org_id  # Optional
```

#### Models

- `gpt-4o`
- `gpt-4o-mini`
- `gpt-4-turbo`
- `gpt-3.5-turbo`
- `o1-preview`
- `o1-mini`

### Azure OpenAI

Azure OpenAI Service provides enterprise-grade AI.

#### Setup

1. Create Azure OpenAI resource in Azure Portal
2. Get endpoint and API key
3. Configure:
   ```bash
   proxybridge credentials add azure
   ```

#### Environment Variables

```bash
AZURE_OPENAI_API_KEY=your_key
AZURE_OPENAI_ENDPOINT=https://your_resource.openai.azure.com/
AZURE_OPENAI_API_VERSION=2024-02-15-preview
```

### Anthropic

Anthropic API for Claude models.

#### Setup

1. Get API key from https://console.anthropic.com/settings/keys
2. Add credentials:
   ```bash
   proxybridge credentials add anthropic
   ```

#### Environment Variables

```bash
ANTHROPIC_API_KEY=sk-ant-...
ANTHROPIC_BASE_URL=https://api.anthropic.com/v1
```

### Google Gemini

Google Gemini API access.

#### Setup

1. Get API key from https://aistudio.google.com/app/api_key
2. Add credentials:
   ```bash
   proxybridge credentials add google
   ```

#### Environment Variables

```bash
GOOGLE_API_KEY=AIza...
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json  # Optional
```

### OpenRouter

OpenRouter provides access to many models through one API.

#### Setup

1. Get API key from https://openrouter.ai/keys
2. Add credentials:
   ```bash
   proxybridge credentials add openrouter
   ```

#### Models

- `openai/gpt-4o`
- `anthropic/claude-3-5-sonnet`
- `google/gemini-1.5-flash`
- `meta-llama/llama-3-70b-instruct`

### Groq

Groq provides fast inference for Llama models.

#### Setup

1. Get API key from https://console.groq.com/keys
2. Add credentials:
   ```bash
   proxybridge credentials add groq
   ```

#### Models

- `llama3-8b-8192`
- `llama3-70b-8192`
- `llama-3.1-8b-instant`
- `llama-3.1-70b-versatile`
- `mixtral-8x7b-32768`

### DeepSeek

DeepSeek API access.

#### Setup

1. Get API key from https://platform.deepseek.com/api_keys
2. Add credentials:
   ```bash
   proxybridge credentials add deepseek
   ```

#### Models

- `deepseek-chat`
- `deepseek-coder`
- `deepseek-reasoner`

### Mistral

Mistral AI API access.

#### Setup

1. Get API key from https://console.mistral.ai/api-keys
2. Add credentials:
   ```bash
   proxybridge credentials add mistral
   ```

#### Models

- `mistral-large`
- `mistral-medium`
- `mistral-small`
- `mistral-tiny`

### Together AI

Together AI provides open-source models.

#### Setup

1. Get API key from https://cloud.together.ai/api-keys
2. Add credentials:
   ```bash
   proxybridge credentials add together
   ```

#### Models

- `togethercomputer/llama-3-70b-chat`
- `togethercomputer/llama-3-8b-chat`
- `togethercomputer/qwen-72b-chat`

### Fireworks AI

Fireworks AI provides fast inference.

#### Setup

1. Get API key from https://fireworks.ai/api-keys
2. Add credentials:
   ```bash
   proxybridge credentials add fireworks
   ```

#### Models

- `accounts/fireworks/models/firefunction-v2`
- `accounts/fireworks/models/llama-v3-70b-instruct`

### Cerebras

Cerebras Cloud API.

#### Setup

1. Get API key from https://cloud.cerebras.ai/api-keys
2. Add credentials:
   ```bash
   proxybridge credentials add cerebras
   ```

#### Models

- `cerebras-llama3.1-8b`
- `cerebras-llama3.1-70b`

### xAI (Grok)

xAI Grok API access.

#### Setup

1. Get API key from https://console.x.ai/api-keys
2. Add credentials:
   ```bash
   proxybridge credentials add xai
   ```

#### Models

- `grok-beta`
- `grok-2-1212`
- `grok-2-vision-1212`

## Local Providers

### Ollama

Run models locally with Ollama.

#### Setup

1. Install Ollama: https://ollama.com/install
2. Pull a model:
   ```bash
   ollama pull llama3
   ollama serve
   ```
3. Add credentials:
   ```bash
   proxybridge credentials add ollama
   ```

### LM Studio

Local LLM server with GUI.

#### Setup

1. Install LM Studio
2. Load a model
3. Start local server (port 1234)
4. Configure ProxyBridge

### vLLM

High-performance inference server.

#### Setup

```bash
pip install vllm
python -m vllm.entrypoints.openai.api_server --model meta-llama/Llama-3.1-8B-Instruct
```

### TGI (Text Generation Inference)

Hugging Face inference server.

#### Setup

```bash
docker run --gpus all --shm-size 1g -p 8080:80 \
  ghcr.io/huggingface/text-generation-inference \
  --model-id meta-llama/Llama-3.1-8B-Instruct
```

## Provider Switching

```bash
# Switch default provider
proxybridge switch anthropic

# Test a provider
proxybridge test openai

# List providers
proxybridge providers

# Add credentials for provider
proxybridge credentials add <provider>
```

## Troubleshooting

### API Key Errors

```bash
# Validate credentials
proxybridge credentials validate

# Add new credentials
proxybridge credentials add openai
```

### Model Not Available

```bash
# List available models for provider
proxybridge models

# Check provider status
proxybridge providers
```

### Connection Issues

```bash
# Check LiteLLM status
proxybridge status

# View logs
proxybridge logs

# Run diagnostics
proxybridge doctor
```
