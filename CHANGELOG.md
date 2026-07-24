# Changelog

All notable changes to ProxyBridge will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-07-24

### Added
- Initial release of ProxyBridge
- Universal LLM proxy through LiteLLM
- Support for 20+ LLM providers
- Model aliasing system
- Secure credential storage (OS vault + encrypted fallback)
- Cross-platform: Windows, macOS, Linux, WSL2
- Provider plugin architecture
- Built-in diagnostics and health monitoring

### Supported Providers
- AWS Bedrock
- OpenAI (GPT-4, GPT-3.5, O series)
- Azure OpenAI
- Anthropic (Claude models)
- Google Gemini
- OpenRouter
- Groq (Llama models)
- DeepSeek
- Mistral
- Together AI
- Fireworks AI
- Cerebras
- xAI (Grok)
- Ollama (local)
- LM Studio (local)

### CLI Commands
- `install` - Install and configure ProxyBridge
- `uninstall` - Remove ProxyBridge
- `repair` - Repair installation
- `doctor` - Run diagnostics
- `status` - Show system status
- `start/stop/restart` - Manage LiteLLM service
- `models` - List configured models
- `providers` - List available providers
- `test <provider>` - Test provider connection
- And more...

### Repository Information
- Owner: Anand Damdiyal
- Repository: https://github.com/root975638-alt/proxybridge
- License: MIT
