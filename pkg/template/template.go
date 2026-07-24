// Package template provides configuration template generation.
package template

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/root975638-alt/proxybridge/internal/config"
	"github.com/root975638-alt/proxybridge/internal/provider"
)

// GenerateLiteLLMConfig generates the LiteLLM configuration
func GenerateLiteLLMConfig(cfg *config.Config, providers []string) (string, error) {
	tmpl := `# LiteLLM Configuration
# Managed by ProxyBridge
# Generated: {{.Timestamp}}
# Version: {{.Version}}

model_list:
{{- range $index, $provider := .Providers }}
{{- range $modelName := $provider.Models }}
  - model_name: {{ $modelName }}
    litellm_params:
      model: {{ $provider.Provider }}/{{ $modelName }}
      api_key: "{{ $provider.ID }}_API_KEY"
{{- end }}
{{- end }}

general_settings:
  master_key: "{{ .MasterKey }}"
  database_url: "file:/tmp/litellm.db"
  num_workers: 4
  cooldown_seconds: 30
  tracing: false
  alerting: []
  alert_types: []

# Environment Variables
litellm_settings:
  drop_params: true
  add_function_to_prompt: false
  headers:
    Cache-Control: "no-cache"
    Content-Type: "application/json"
`
	return renderTemplate(tmpl, cfg)
}

// GenerateEnvironmentConfig generates the environment configuration file
func GenerateEnvironmentConfig(cfg *config.Config, providers []string) (string, error) {
	tmpl := `# Environment Configuration
# Managed by ProxyBridge
# Generated: {{.Timestamp}}
# Version: {{.Version}}

# LiteLLM Settings
LITELLM_PORT: {{ .LiteLLM.Port }}
LITELLM_HOST: {{ .LiteLLM.ListenAddress }}

# Provider API Keys (Reference)
# These should be set via the ProxyBridge credential manager
{{- range $provider := .Providers }}
# {{ $provider.ID }}_API_KEY={{ $provider.ID | upper }}_API_KEY_VALUE
{{- end }}

# Claude Code Integration
ClaudeCodeEnabled: "{{ .ClaudeCode.Enabled }}"
ClaudeCodeModel: "{{ .DefaultModel }}"
`
	return renderTemplate(tmpl, cfg)
}

// GenerateClaudeConfig generates the Claude Code configuration
func GenerateClaudeConfig(cfg *config.Config) (string, error) {
	tmpl := `{
  "models": [
    {
      "name": "{{ .DefaultModel }}",
      "provider": "{{ .ActiveProvider }}",
      "model": "{{ .ActiveProvider }}/{{ .DefaultModel }}",
      "apiKey": "{{ .ActiveProvider }}_API_KEY",
      "maxTokens": {{ .Settings.MaxTokens }}
    }
  ],
  "defaultModel": "{{ .DefaultModel }}",
  "provider": "{{ .ActiveProvider }}",
  "proxyUrl": "http://{{ .LiteLLM.ListenAddress }}:{{ .LiteLLM.Port }}",
  "proxyEnabled": true,
  "version": "1.0"
}`
	return renderTemplate(tmpl, cfg)
}

// GenerateSystemdService generates a systemd service file
func GenerateSystemdService(cfg *config.Config) (string, error) {
	tmpl := `[Unit]
Description=LiteLLM - Universal LLM Proxy
Documentation=https://docs.litellm.ai
After=network.target

[Service]
Type=simple
User={{ .User }}
Group={{ .Group }}
Environment="LITELLM_CONFIG={{ .Config.LiteLLM.ConfigPath }}"
Environment="LITELLM_PORT={{ .LiteLLM.Port }}"
Environment="LITELLM_HOST={{ .LiteLLM.ListenAddress }}"
ExecStart={{ .LiteLLM.Path }} --host {{ .LiteLLM.ListenAddress }} --port {{ .LiteLLM.Port }} --config {{ .LiteLLM.ConfigPath }}
Restart=on-failure
RestartSec=5
StartLimitIntervalSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
`
	return renderTemplate(tmpl, cfg)
}

// GenerateLaunchdService generates a macOS launchd plist
func GenerateLaunchdService(cfg *config.Config) (string, error) {
	tmpl := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.proxybridge.litellm</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{ .LiteLLM.Path }}</string>
        <string>--host</string>
        <string>{{ .LiteLLM.ListenAddress }}</string>
        <string>--port</string>
        <string>{{ .LiteLLM.Port }}</string>
        <string>--config</string>
        <string>{{ .LiteLLM.ConfigPath }}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>{{ .LiteLLM.LogPath }}/litellm.log</string>
    <key>StandardErrorPath</key>
    <string>{{ .LiteLLM.LogPath }}/litellm.error.log</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>LITELLM_CONFIG</key>
        <string>{{ .LiteLLM.ConfigPath }}</string>
        <key>LITELLM_PORT</key>
        <string>{{ .LiteLLM.Port }}</string>
    </dict>
</dict>
</plist>
`
	return renderTemplate(tmpl, cfg)
}

// GenerateStartScript generates a cross-platform start script
func GenerateStartScript(cfg *config.Config) (string, error) {
	tmpl := `#!/bin/bash
# ProxyBridge LiteLLM Start Script
# Generated: {{ .Timestamp }}

set -e

LITELLM_PATH="{{ .LiteLLM.Path }}"
CONFIG_PATH="{{ .LiteLLM.ConfigPath }}"
PORT="{{ .LiteLLM.Port }}"
HOST="{{ .LiteLLM.ListenAddress }}"

# Check if LiteLLM is installed
if ! command -v "$LITELLM_PATH" &> /dev/null; then
    echo "Error: LiteLLM not found at $LITELLM_PATH"
    exit 1
fi

# Check if config exists
if [ ! -f "$CONFIG_PATH" ]; then
    echo "Error: Config file not found at $CONFIG_PATH"
    exit 1
fi

# Start LiteLLM
echo "Starting LiteLLM on $HOST:$PORT..."
"$LITELLM_PATH" --host "$HOST" --port "$PORT" --config "$CONFIG_PATH"
`
	return renderTemplate(tmpl, cfg)
}

// renderTemplate renders a template with the config
func renderTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("config").Funcs(template.FuncMap{
		"upper": strings.ToUpper,
	}).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ProviderConfig holds provider configuration for templates
type ProviderConfig struct {
	Name              string
	ID                string
	DefaultModel      string
	APIKeyPlaceholder string
	Settings          map[string]string
}

// GenerateProviderConfig generates a provider configuration
func GenerateProviderConfig(p provider.Provider) (string, error) {
	cfg := &ProviderConfig{
		Name:              p.Name(),
		ID:                p.ID(),
		DefaultModel:      p.GetDefaultModel(),
		APIKeyPlaceholder: fmt.Sprintf("%s_API_KEY", strings.ToUpper(p.ID())),
		Settings:          p.GetConfig(),
	}

	tmpl := `# Provider: {{ .Name }}
# ID: {{ .ID }}
# Default Model: {{ .DefaultModel }}

{{ .ID | upper }}_API_KEY={{ .APIKeyPlaceholder }}

# Provider-Specific Settings
{{- range $key, $value := .Settings }}
{{ $key }}={{ $value }}
{{- end }}
`
	return renderTemplate(tmpl, cfg)
}

// GenerateModelAliasConfig generates the model alias configuration
func GenerateModelAliasConfig(cfg *config.Config, aliases map[string]string) (string, error) {
	tmpl := `# Model Alias Configuration
# Managed by ProxyBridge

# Default Mappings
{{- range $alias, $model := .Aliases }}
{{ $alias }}: {{ $model }}
{{- end }}

# Provider Mappings
{{- range $provider, $model := .ProviderMappings }}
{{ $provider }}: {{ $model }}
{{- end }}
`
	return renderTemplate(tmpl, map[string]interface{}{
		"Aliases":        aliases,
		"ProviderMappings": map[string]string{"default": cfg.DefaultModel},
	})
}

// GenerateConfigReport generates a validation report
func GenerateConfigReport(issues []string, warnings []string) (string, error) {
	tmpl := `# ProxyBridge Configuration Report
# Generated: {{ .Timestamp }}

## Status: {{ .Status }}

{{- if .Issues }}
### Issues Found ({{ len .Issues }})
{{- range $index, $issue := .Issues }}
{{ $index | add 1 }}. {{ $issue }}
{{- end }}
{{- end }}

{{- if .Warnings }}
### Warnings ({{ len .Warnings }})
{{- range $index, $warning := .Warnings }}
{{ $index | add 1 }}. {{ $warning }}
{{- end }}
{{- end }}

{{- if and (not .Issues) (not .Warnings) }}
No issues or warnings found. Configuration is valid.
{{- end }}
`
	return renderTemplate(tmpl, map[string]interface{}{
		"Timestamp":  "TODO",
		"Status":     "VALID",
		"Issues":     issues,
		"Warnings":   warnings,
	})
}
