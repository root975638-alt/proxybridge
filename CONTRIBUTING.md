# Contributing to ProxyBridge

Thank you for your interest in contributing to ProxyBridge! This document provides guidelines and instructions for contributing.

**Repository**: https://github.com/root975638-alt/proxybridge  
**Owner**: Anand Damdiyal  
**License**: MIT

## Code of Conduct

Please be respectful and constructive in all interactions. We aim to be an inclusive, welcoming community.

## How Can I Contribute?

### Reporting Bugs

- Use the GitHub issue tracker
- Include steps to reproduce
- Include expected vs actual behavior
- Include relevant configuration and environment details

### Suggesting Features

- Use the GitHub issue tracker
- Explain the use case
- Describe the proposed implementation
- Consider backward compatibility

### Pull Requests

1. Fork the repository
2. Create a branch from `main`
3. Make your changes
4. Write tests if applicable
5. Update documentation as needed
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.23+
- Python 3.8+
- Git

### Getting Started

```bash
# Clone the repository
git clone https://github.com/root975638-alt/proxybridge.git
cd proxybridge

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o proxybridge ./cmd/proxybridge.go
```

## Project Structure

```
proxybridge/
├── cmd/                    # CLI entry points
├── internal/               # Private application code
│   ├── alias/             # Model alias management
│   ├── config/            # Configuration management
│   ├── credential/        # Credential storage
│   ├── diagnostic/        # Diagnostics engine
│   ├── exportimport/      # Config export/import
│   ├── health/            # Health monitoring
│   ├── installer/         # Installation logic
│   ├── logging/           # Structured logging
│   ├── provider/          # Provider plugins
│   └── util/              # Utility functions
├── pkg/                    # Public packages
│   ├── detection/         # Platform detection
│   ├── manager/           # Service managers
│   └── template/          # Config templates
├── docs/                   # Documentation
└── tests/                  # Test files
```

## Coding Standards

### Go Style

- Follow Go conventions (gofmt)
- Use meaningful names
- Keep functions focused (single responsibility)
- Document public API
- Handle errors explicitly

### Testing

- All new features need tests
- Aim for high coverage
- Use table-driven tests where applicable
- Mock external dependencies

### Documentation

- Update README for user-facing changes
- Update architecture docs for structural changes
- Document new providers
- Include examples where helpful

## Pull Request Process

1. Update documentation as needed
2. Add tests if applicable
3. Ensure tests pass: `go test ./...`
4. Ensure code is formatted: `gofmt -w .`
5. Ensure linting passes: `golangci-lint run`
6. Squash commits when appropriate
7. Wait for review and CI checks

## Commit Messages

Follow the conventional commits format:

```
<type>(<scope>): <description>

[optional body]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test-related changes
- `chore`: Maintenance tasks
- `refactor`: Code refactoring

Scopes:
- `cli`: CLI interface
- `provider`: Provider plugins
- `install`: Installation
- `config`: Configuration
- `diagnostic`: Diagnostics
- `docs`: Documentation

Example:
```
feat(provider): add support for OpenRouter

- Register OpenRouter provider
- Add API key validation
- Include model list
```

## Questions?

- Open a GitHub issue at https://github.com/root975638-alt/proxybridge/issues
- Email maintainers at anand.damdiyal@proxybridge.dev

## Recognition

Contributors are recognized in:
- `AUTHORS` file
- Release notes
- README contributors section
