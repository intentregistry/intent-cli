# Intent CLI User Guide

Welcome to Intent CLI! This guide covers everything you need to know to create, package, and publish intents.

## Table of Contents

- [Quick Start](#quick-start)
- [Creating Intents](#creating-intents)
- [ITML Format](#itml-format)
- [Running Intents](#running-intents)
- [Packaging](#packaging)
- [Publishing](#publishing)
- [Installing Packages](#installing-packages)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

## Quick Start

### Installation

```bash
# Via Homebrew
brew install intent-cli

# Or build from source
git clone https://github.com/intentregistry/intent-cli.git
cd intent-cli
go build -o intent ./cmd/intent
sudo mv intent /usr/local/bin/
```

### Create Your First Intent

```bash
# Initialize a new project
intent init my-project --app

# Navigate to project
cd my-project

# Edit your intent
vim intents/hello.itml

# Run it
intent run intents/hello.itml --inputs name="World"

# Package it
intent package . --unsigned

# Publish it
intent publish .
```

## Creating Intents

### Project Structure

```
my-project/
├── itpkg.json              # Package manifest
├── project.app.itml        # App entry point (if type=app)
├── intents/
│   ├── hello.itml         # Individual intents
│   └── ...
├── policies/
│   └── security.itml      # Security policies
├── tests/                 # Optional: test files
├── schemas/               # Optional: JSON schemas
└── assets/               # Optional: resources
```

### Scaffold a Project

```bash
# Create new project with scaffolding
intent init my-awesome-intent --app

# Flags:
# --app       Create app package (with project.app.itml entry point)
# --scope     Set package scope (default: "scope")
```

## ITML Format

ITML (Intent Task Markup Language) is a custom DSL for defining intents.

### Basic Syntax

```itml
intent "Intent Name"
description: "What this intent does"
itmlVersion: "0.1"

inputs:
  - name (string) required default="default value"
  - count (integer) default="1"
  - enabled (boolean) default="true"

outputs:
  - result (string)
  - data (array)

workflow:
  → log("Starting intent")
  → log("Hello {name}")
  → return(status="ok")
```

### Supported Types

- `string` - Text data
- `integer` - Whole numbers
- `number` - Decimals
- `boolean` - true/false
- `array` - Lists
- `object` - Key-value pairs

### Workflow Commands

- `log("message")` - Print to output
- `http.get(url)` - Make HTTP request
- `file.read(path)` - Read file
- `file.write(path, content)` - Write file
- `transform(data, mapping)` - Transform data
- `return(status="ok")` - Return results

### Parameters

- `required` - Parameter must be provided
- `default="value"` - Default value if not provided

### Examples

See `examples/` directory for complete examples:
- `hello-world.itml` - Basic "Hello World"
- `weather.itml` - HTTP requests
- `text-analyzer.itml` - Text processing
- `image-processor.itml` - File handling
- `data-transformer.itml` - Format conversion

## Running Intents

### Basic Execution

```bash
intent run path/to/intent.itml
```

### With Input Parameters

```bash
intent run intents/weather.itml \
  --inputs city="New York" \
  --inputs units="fahrenheit"
```

### Save Output

```bash
intent run intents/process.itml \
  --inputs data="value" \
  --output-dir ./results/
```

### Verbose Output

```bash
intent run intents/hello.itml --verbose
```

## Packaging

### Create a Package

```bash
# Unsigned package (for development)
intent package . --unsigned

# Signed package (production)
intent package . --sign-key ~/.ssh/intent_sign_key.hex
```

### Package Output

Creates `.itpkg` file containing:
- `itpkg.json` - Manifest with metadata
- `MANIFEST.sha256` - File checksums
- `SIGNATURE` - Ed25519 digital signature
- All intent, policy, and asset files

### Verification

```bash
# Verify package signature and integrity
intent verify dist/my-package-0.1.0.itpkg
```

## Publishing

### Publish a Package

```bash
# Direct publish (from .itpkg file)
intent publish dist/my-package-0.1.0.itpkg

# Package and publish directory
intent publish .
```

### Publish Options

```bash
intent publish . \
  --tag beta                    # beta or rc release
  --private                     # Publish as private
  --message "Release notes"     # Release message
```

### Authentication

```bash
# Login to registry
intent login

# Check login status
intent whoami

# Set API URL
export INTENT_API_URL=https://api.intentregistry.com

# Set authentication token
export INTENT_TOKEN=your-token-here
```

## Installing Packages

### Install a Package

```bash
# Install latest version
intent install @scope/package-name

# Install specific version
intent install @scope/package-name@1.0.0
```

### Search Registry

```bash
# Search for packages
intent search weather

# Search with filter
intent search "@scope/*"
```

### Install Location

Packages are installed to `intents/` directory by default.

```bash
# Custom destination
intent install @scope/package --dest ./lib/
```

## Testing

### Run Tests

```bash
# Test entire package
intent test .

# Test specific file
intent test tests/test.itml

# Verbose output
intent test . --verbose
```

### Test Format

Tests are `.itml` files in `tests/` directory:

```itml
intent "Test Suite"

inputs:
  - subject (string) default="test"

workflow:
  → assert.equal(1, 1, "basic equality")
  → assert.truthy(true, "truthiness check")
  → return(status="ok")
```

## Troubleshooting

### Intent Won't Run

**Error: `file must have .itml extension`**
- Use `.itml` files only
- Check file path and extension

**Error: `required parameter not provided`**
- Check intent definition for required inputs
- Provide all required parameters with `--inputs`

### Packaging Issues

**Error: `itpkg.json not found`**
- Initialize with: `intent init . --app`
- Or use: `intent package . --scaffold`

**Error: `structure validation failed`**
- Create required directories: `mkdir -p intents/ policies/`
- Or use `--scaffold` flag

### Publishing Fails

**Error: `version already exists`**
- Update version in `itpkg.json`
- Version must be unique per publish

**Error: `No authentication token configured`**
- Run `intent login` and enter token
- Or set `INTENT_TOKEN` environment variable

### API Connection Issues

**Error: `404 Not Found` on resolve**
- Check API URL: `intent whoami`
- Verify registry is running
- Check network connectivity

**Error: `connection refused`**
- Ensure registry is accessible
- Check firewall rules
- Verify `INTENT_API_URL` setting

## Configuration

### Config File

CLI config stored in `~/.intent/config.yaml`:

```yaml
api_url: https://api.intentregistry.com
token: your-token-here
telemetry: false
```

### Environment Variables

```bash
# API configuration
export INTENT_API_URL=http://localhost:8080
export INTENT_TOKEN=your-token

# Package signing
export INTENT_SIGN_KEY=~/.ssh/intent_sign_key.hex

# Enable telemetry
export INTENT_TELEMETRY=true
```

### Project .env File

Create `.env` in project root for local overrides:

```
INTENT_API_URL=http://localhost:8080
INTENT_TOKEN=local-dev-token
```

## Advanced Topics

### Custom Policies

Define security policies in `policies/`:

```itml
intent "Security Policy"

policies:
  security:
    network:
      outbound:
        allow:
          - "*.example.com"
        deny:
          - "*"
```

### Capabilities

Declare what your intent can access:

```json
{
  "capabilities": [
    "http.outbound",
    "file.read",
    "file.write",
    "network.dns"
  ]
}
```

### Multi-file Projects

Organize large projects:

```
intents/
  ├── core/
  │   ├── base.itml
  │   └── utils.itml
  ├── web/
  │   ├── fetch.itml
  │   └── parse.itml
  └── main.itml
```

## Support

- 📖 Documentation: `docs/`
- 💬 Issues: github.com/intentregistry/intent-cli/issues
- 🔗 Registry: https://intentregistry.com
- 📧 Contact: support@intentregistry.com

## License

Apache 2.0
