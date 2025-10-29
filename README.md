# Intent CLI

The official command-line interface for publishing and installing
**Intents** on [IntentRegistry](https://intentregistry.com).

## Installation

You can install the CLI using Homebrew:

```bash
brew tap intentregistry/homebrew-tap
brew install intent-cli
intent --help
```

## Usage

```bash
# Authentication
intent login

# Package creation (creates signed .itpkg archives)
intent package [path] --scaffold --unsigned  # Development/testing
intent package [path] --sign-key ~/.ssh/intent_key  # Production signing

# Publishing
intent publish [path] --private --tag beta --message "first release"

# Installation
intent install @scope/name[@version] --dest intents

# Execution
intent run FILE.itml --inputs name=World

# Discovery
intent search "vector embeddings"

# Testing
intent test [path] --format json

# Utilities
intent whoami
intent doctor
```

## Configuration

Configuration file path:

    ~/.intent/config.yaml

You can also set environment variables:

- `INTENT_API_URL` (default: `https://api.intentregistry.com`)
- `INTENT_TOKEN`
- `INTENT_SIGN_KEY` (path to ed25519 private key for package signing)

## Shell completion

The CLI can generate completion scripts for **zsh**, **bash**, **fish**, and **PowerShell**.

### Homebrew (recommended)
If you installed via Homebrew, completions are installed automatically by the formula.

Check it’s present:
```bash
brew cat intentregistry/tap/intent-cli | grep -q generate_completions_from_executable && echo "Completions enabled"
```

### Manual setup

Generate the script for your shell and place it in the standard location:

#### zsh (macOS default)
```bash
# one-time: ensure completion system is enabled
echo 'autoload -Uz compinit; compinit' >> ~/.zshrc

# install completion
sudo mkdir -p /opt/homebrew/share/zsh/site-functions
intent completion zsh | sudo tee /opt/homebrew/share/zsh/site-functions/_intent > /dev/null

# reload shell or run:
autoload -Uz compinit; compinit
```

#### bash
```bash
# macOS (Homebrew bash-completion dir)
sudo mkdir -p /opt/homebrew/etc/bash_completion.d
intent completion bash | sudo tee /opt/homebrew/etc/bash_completion.d/intent > /dev/null

# Linux (system-wide)
sudo mkdir -p /etc/bash_completion.d
intent completion bash | sudo tee /etc/bash_completion.d/intent > /dev/null

# current shell only:
source <(intent completion bash)
```

#### fish
```bash
mkdir -p ~/.config/fish/completions
intent completion fish > ~/.config/fish/completions/intent.fish
```

#### PowerShell
```powershell
# current session
intent completion powershell | Out-String | Invoke-Expression

# persist for future sessions (adjust path to your profile)
$OutPath = "$HOME\Documents\PowerShell\Scripts\intent.ps1"
intent completion powershell > $OutPath
Add-Content $PROFILE "`n. $OutPath"
```

> Tip: You can preview the script without installing it:
> ```bash
> intent completion zsh | head
> ```

## Package Format (.itpkg)

The `.itpkg` format is a signed, versioned Intent package container:

- **Structure**: Flat tar.gz archive with `itpkg.json`, `MANIFEST.sha256`, `SIGNATURE`, and project files
- **Signing**: ed25519 signature over MANIFEST.sha256 for integrity verification
- **Manifest**: Required `itpkg.json` with name, version, policies, and capabilities
- **Validation**: Directory structure validation (requires `intents/` and `policies/` directories)

### Quick Start

```bash
# Generate signing key (one-time setup)
./gen_intent_key.sh  # Creates private_key.hex and public_key.hex

# Package with scaffold (creates itpkg.json and required directories)
intent package . --scaffold --unsigned

# Package with signing
export INTENT_SIGN_KEY=~/.ssh/private_key.hex
intent package . --out dist/
```

See [docs/itpkg_definition.md](docs/itpkg_definition.md) for complete specification.

## Development

To build locally:

```bash
go mod tidy
go run ./cmd/intent --help
make build
make pack-darwin-arm64
make checksum
```

### Local Testing with Mock API

For complete local testing, you can run the mock API server:

```bash
# Build mock API
go build -o mock-api ./cmd/mock-api

# Run mock API (starts on http://localhost:8080)
./mock-api

# Configure CLI for local API
export INTENT_API_URL=http://localhost:8080
export INTENT_TOKEN=local-dev-token

# Test full workflow
./intent init test-project --app
cd test-project
../intent run intents/hello.itml
../intent package . --unsigned
../intent publish .
```

See [docs/LOCAL_DEVELOPMENT.md](docs/LOCAL_DEVELOPMENT.md) for complete testing guide.

### Documentation

- **[USER_GUIDE.md](docs/USER_GUIDE.md)** - Complete guide covering all commands, ITML format, troubleshooting
- **[LOCAL_DEVELOPMENT.md](docs/LOCAL_DEVELOPMENT.md)** - Setup, testing, debugging, and CI/CD workflows
- **[itpkg_definition.md](docs/itpkg_definition.md)** - Detailed .itpkg package format specification

### Example Intents

See `examples/` directory for production-ready intents:
- `hello-world.itml` - Basic "Hello World" example
- `weather.itml` - HTTP API integration
- `text-analyzer.itml` - Text processing
- `image-processor.itml` - File handling
- `data-transformer.itml` - Format conversion

## Project Structure

    intent-cli/
    ├─ cmd/intent/              # main entrypoint
    ├─ internal/cmd/            # subcommands
    ├─ internal/config/         # configuration management
    ├─ internal/httpclient/     # API HTTP client
    ├─ internal/pack/           # packaging utilities (.itpkg format + tar.gz)
    ├─ internal/parser/         # ITML format parser (DSL, JSON, YAML)
    ├─ internal/executor/       # intent execution engine
    ├─ internal/version/        # version and build metadata
    ├─ .goreleaser.yaml         # release automation
    ├─ go.mod                   # Go module
    ├─ Makefile                 # build tasks
    └─ README.md

## License

MIT License © 2025 [IntentRegistry](https://intentregistry.com)
