# Troubleshooting Guide

Common issues and their solutions for ProxyBridge.

## Installation Issues

### Python Not Found

**Error**: `Python 3 is not installed`

**Solution**:
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

### LiteLLM Installation Failed

**Error**: `Failed to install LiteLLM`

**Solution**:
```bash
# Check pip is available
python3 -m pip --version

# Install pip if needed
python3 -m ensurepip --upgrade

# Install LiteLLM manually
pipx install litellm

# Or use pip
pip3 install litellm
```

### Permission Denied

**Error**: `Permission denied writing to config directory`

**Solution**:
```bash
# Fix directory permissions
sudo chown -R $USER:$USER ~/.config/proxybridge
chmod -R 755 ~/.config/proxybridge

# Or run with sudo
sudo proxybridge install
```

## Runtime Issues

### LiteLLM Not Starting

**Error**: `LiteLLM failed to start`

**Solutions**:
1. Check if LiteLLM is installed:
   ```bash
   litellm --version
   ```

2. Check for port conflicts:
   ```bash
   # Check if port 4000 is in use
   lsof -i :4000
   netstat -tulpn | grep :4000
   ```

3. View logs:
   ```bash
   proxybridge logs
   ```

4. Try with verbose logging:
   ```bash
   proxybridge start --verbose
   ```

### API Key Errors

**Error**: `API key validation failed`

**Solutions**:
1. Verify credentials:
   ```bash
   proxybridge credentials validate
   ```

2. Re-add credentials:
   ```bash
   proxybridge credentials add openai
   ```

3. Check environment variable:
   ```bash
   echo $OPENAI_API_KEY
   ```

### Claude Code Can't Connect

**Error**: `Claude Code connection failed`

**Solutions**:
1. Ensure LiteLLM is running:
   ```bash
   proxybridge start
   ```

2. Verify LiteLLM is accessible:
   ```bash
   curl http://localhost:4000/v1/models
   ```

3. Check Claude Code config:
   ```bash
   proxybridge config show
   ```

4. Restart both:
   ```bash
   proxybridge restart
   ```

## Configuration Issues

### Invalid Configuration

**Error**: `Configuration validation failed`

**Solution**:
```bash
# Validate configuration
proxybridge validate

# View diagnostic report
proxybridge doctor

# Export current config
proxybridge export backup.yaml

# Reinstall
proxybridge uninstall
proxybridge install
```

### Provider Not Configured

**Error**: `Provider not configured`

**Solution**:
```bash
# List providers
proxybridge providers

# Add credentials
proxybridge credentials add <provider>

# Test provider
proxybridge test <provider>
```

##Diagnostic Issues

### Run Diagnostics

```bash
proxybridge doctor
```

This will check:
- LiteLLM status
- Provider connectivity
- Credential validity
- Claude Code integration
- Configuration integrity

### View Detailed Logs

```bash
# View LiteLLM logs
proxybridge logs

# View last N lines
proxybridge logs --lines 100

# Follow logs
tail -f ~/.config/proxybridge/logs/litellm.log
```

## Common Commands

| Problem | Command |
|---------|---------|
| Check status | `proxybridge status` |
| Run diagnostics | `proxybridge doctor` |
| Start LiteLLM | `proxybridge start` |
| Stop LiteLLM | `proxybridge stop` |
| Restart LiteLLM | `proxybridge restart` |
| View logs | `proxybridge logs` |
| Validate config | `proxybridge validate` |
| Export config | `proxybridge export file.yaml` |
| Import config | `proxybridge import file.yaml` |

## Getting Help

If you can't resolve the issue:

1. Run diagnostics:
   ```bash
   proxybridge doctor --json > diagnostics.json
   ```

2. Collect logs:
   ```bash
   proxybridge logs > logs.txt
   ```

3. Open GitHub issue with:
   - Output from `proxybridge version`
   - Output from `proxybridge status`
   - Diagnostics file
   - Log file
   - Steps to reproduce

## Platform-Specific Issues

### Windows

**Issue**: Path too long

**Solution**:
```powershell
# Enable long paths
reg add HKLM\SYSTEM\CurrentControlSet\Control\FileSystem /v LongPathsEnabled /t REG_DWORD /d 1

# Or use short paths
proxybridge --config C:\pb\config
```

### macOS

**Issue**: Security restrictions

**Solution**:
```bash
# Allow from unidentified developer
xattr -d com.apple.quarantine /usr/local/bin/proxybridge

# Or allow in System Preferences > Security
```

### WSL2

**Issue**: Port binding

**Solution**:
```bash
# Use 127.0.0.1 instead of localhost
proxybridge start --host 127.0.0.1

# Or configure in config.yaml
settings:
  listen_address: 127.0.0.1
```
