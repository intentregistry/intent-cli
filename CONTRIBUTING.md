# Contributing to Intent CLI

Thank you for your interest in contributing to the Intent CLI! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)
- [Code Style Guidelines](#code-style-guidelines)
- [Documentation](#documentation)

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/). By participating, you agree to uphold this code.

## Getting Started

### Prerequisites

- Go 1.23.0 or later
- Git
- Make (for build automation)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/intent-cli.git
   cd intent-cli
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/intentregistry/intent-cli.git
   ```

## Development Setup

### Install Dependencies

```bash
go mod tidy
```

### Build the Project

```bash
# Build the binary
make build

# Build and package for macOS ARM64
make pack-darwin-arm64

# Generate checksums
make checksum
```

### Run Locally

```bash
# Run the CLI directly
go run ./cmd/intent --help

# Test specific commands
go run ./cmd/intent login
go run ./cmd/intent whoami
```

## Project Structure

```
intent-cli/
├── cmd/intent/              # Main application entrypoint
├── internal/
│   ├── cmd/                 # CLI subcommands (login, publish, install, etc.)
│   ├── config/              # Configuration management
│   ├── httpclient/          # HTTP client for API communication
│   ├── pack/                # Tar.gz packaging utilities
│   └── version/             # Version and build metadata
├── dist/                    # Build artifacts
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
├── Makefile                 # Build automation
└── README.md               # Project documentation
```

### Key Components

- **`cmd/intent/main.go`**: Application entrypoint that registers all subcommands
- **`internal/cmd/`**: Contains all CLI subcommands using Cobra framework
- **`internal/config/`**: Handles configuration file and environment variable management
- **`internal/httpclient/`**: HTTP client wrapper for API interactions
- **`internal/pack/`**: Utilities for creating tar.gz packages

## Making Changes

### Branch Strategy

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the [Code Style Guidelines](#code-style-guidelines)

3. Test your changes thoroughly

4. Commit your changes with clear, descriptive messages

### Adding New Commands

To add a new CLI command:

1. Create a new file in `internal/cmd/` (e.g., `newcommand.go`)
2. Implement the command using the Cobra framework
3. Register the command in `cmd/intent/main.go`
4. Add appropriate tests
5. Update documentation

Example command structure:
```go
package cmd

import (
    "github.com/spf13/cobra"
)

func NewCommandCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "newcommand",
        Short: "Brief description",
        Long:  "Detailed description",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Command implementation
            return nil
        },
    }
}
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/cmd/
```

### Test Guidelines

- Write unit tests for new functionality
- Aim for good test coverage
- Use table-driven tests where appropriate
- Mock external dependencies (HTTP clients, file system, etc.)

## Submitting Changes

### Pull Request Process

1. Ensure your branch is up to date with `main`:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

3. Create a Pull Request on GitHub with:
   - Clear title and description
   - Reference any related issues
   - Include screenshots for UI changes
   - Ensure all CI checks pass

### Pull Request Guidelines

- **Title**: Use clear, descriptive titles
- **Description**: Explain what changes you made and why
- **Size**: Keep PRs focused and reasonably sized
- **Tests**: Include tests for new functionality
- **Documentation**: Update documentation as needed

### Review Process

- All PRs require review from maintainers
- Address feedback promptly
- Keep discussions constructive and focused
- Be patient - maintainers are volunteers

## Release Process

Releases are automated and managed by maintainers. The process includes:

1. Version bumping
2. Building binaries for multiple platforms
3. Creating GitHub releases
4. Publishing to package managers

## Code Style Guidelines

### Go Code Style

- Follow standard Go formatting (`gofmt`)
- Use `golint` and `go vet` for code quality
- Follow Go naming conventions
- Use meaningful variable and function names
- Add comments for exported functions and types

### Error Handling

- Always handle errors explicitly
- Use `fmt.Errorf` for error wrapping
- Provide meaningful error messages
- Log errors appropriately

### Configuration

- Use environment variables for configuration
- Support configuration files
- Provide sensible defaults
- Document all configuration options

### HTTP Client Usage

- Use the provided HTTP client wrapper
- Handle HTTP errors appropriately
- Include proper headers and authentication
- Implement retry logic where appropriate

## Documentation

### Code Documentation

- Document all exported functions and types
- Use Go doc comments format
- Include examples for complex functions
- Keep documentation up to date

### User Documentation

- Update README.md for user-facing changes
- Document new CLI commands and options
- Provide usage examples
- Keep installation instructions current

## Getting Help

- **Issues**: Use GitHub Issues for bug reports and feature requests
- **Discussions**: Use GitHub Discussions for questions and general discussion
- **Email**: Contact maintainers directly for sensitive issues

## Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes (for significant contributions)
- Project documentation

Thank you for contributing to Intent CLI! 🚀
