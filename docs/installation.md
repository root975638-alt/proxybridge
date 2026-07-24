# Installation Guide

This guide covers installing ProxyBridge on all supported platforms.

## Prerequisites

Before installing ProxyBridge, ensure you have:

- Python 3.8+: Required for LiteLLM
- Node.js 18+: Optional, for Claude Code integration
- Administrative access: For system-wide installation

## Quick Installation

### Linux/macOS

```bash
curl -fsSL https://proxybridge.com/install.sh | sh
```

### Windows

```powershell
powershell -ExecutionPolicy Bypass -Command "iwr -useb https://proxybridge.com/install.ps1 | iex"
```

## Manual Installation

### Building from Source

```bash
# Clone the repository
git clone https://github.com/root975638-alt/proxybridge.git
cd cli

# Build the binary
go build -o proxybridge

# Install system-wide
sudo cp proxybridge /usr/local/bin/

# Verify installation
proxybridge version
```

### Platform-Specific Installation

#### Ubuntu/Debian

```bash
# Add repository
wget -qO- https://proxybridge.com/repo.key | sudo gpg --dearmor -o /usr/share/keyrings/proxybridge-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/proxybridge-archive-keyring.gpg] https://proxybridge.com/deb stable main" | sudo tee /etc/apt/sources.list.d/proxybridge.list

# Update and install
sudo apt update
sudo apt install proxybridge
```

#### Fedora/RHEL/CentOS

```bash
# Add repository
sudo dnf config-manager --add-repo https://proxybridge.com/rpm/proxybridge.repo

# Install
sudo dnf install proxybridge
```

#### macOS (Homebrew)

```bash
brew install proxybridge/tap/proxybridge
```

#### Windows (Winget)

```powershell
winget install ProxyBridge.ProxyBridge
```

## Post-Installation Setup

1. Run the installation wizard:
```bash
proxybridge install
```

2. Follow the prompts to:
   - Set up your API credentials
   - Choose default provider
   - Configure Claude Code integration

## Verifying Installation

```bash
# Check version
proxybridge version

# Check status
proxybridge status

# Run diagnostics
proxybridge doctor
```

## Troubleshooting

### Python Not Found

If Python is not found:

```bash
# Ubuntu/Debian
sudo apt install python3

# Fedora/RHEL
sudo dnf install python3

# macOS
brew install python

# Windows
winget install Python.Python.3.12
```

### Permission Denied

If you encounter permission issues:

```bash
# Linux/macOS
sudo ./proxybridge install

# Or change ownership
sudo chown -R $USER:$USER ~/.config/proxybridge
```

### LiteLLM Installation Failed

LiteLLM requires Python pip:

```bash
# Install pip
python3 -m ensurepip --upgrade

# Or use pipx
pipx ensurepath
pipx install litellm
```

## Updating ProxyBridge

```bash
# Using the built-in update command
proxybridge self-update

# Or using package manager
brew upgrade proxybridge
sudo apt upgrade proxybridge
```

## uninstallation

```bash
# Using the uninstall command
proxybridge uninstall
```

This will:
1. Stop LiteLLM
2. Remove ProxyBridge configuration
3. Keep API keys in OS credential store
