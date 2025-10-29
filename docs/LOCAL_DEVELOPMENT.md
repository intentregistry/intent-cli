# Local Development Guide

This guide covers setting up a local development environment for intent-cli with a mock API server for testing the full publish/install cycle.

## Prerequisites

- Go 1.23+
- Git
- Make (optional, for Makefile commands)

## Setup

### 1. Clone and Build

```bash
git clone https://github.com/intentregistry/intent-cli.git
cd intent-cli

# Build the CLI
go build -o intent ./cmd/intent

# Build the mock API
go build -o mock-api ./cmd/mock-api

# Verify builds
./intent --version
./mock-api --help
```

### 2. Start Mock API

```bash
# In one terminal
./mock-api

# Default port: 8080
# Endpoints:
#   GET    /health
#   POST   /v1/packages/publish
#   GET    /v1/packages/resolve?spec=@scope/name[@version]
#   GET    /v1/packages/search?q=query
#   GET    /v1/packages/tarball/:name/:version.itpkg
```

### 3. Configure CLI for Local API

Option A: Environment Variables
```bash
export INTENT_API_URL=http://localhost:8080
export INTENT_TOKEN=local-dev-token
```

Option B: Create `.env` file in project root
```bash
cat > .env << EOF
INTENT_API_URL=http://localhost:8080
INTENT_TOKEN=local-dev-token
EOF
```

Option C: Use `intent login`
```bash
./intent login
# Enter: http://localhost:8080
# Enter: local-dev-token
```

## Development Workflow

### 1. Create a Test Project

```bash
# Create new project
./intent init test-project --app

cd test-project

# Edit your intent
vim intents/hello.itml
```

### 2. Run Locally

```bash
# Run the intent
../intent run intents/hello.itml --inputs name="Developer"

# With verbose output
../intent run intents/hello.itml --verbose
```

### 3. Package Your Intent

```bash
# Create unsigned package (for development)
../intent package . --unsigned

# Or with signing key
../intent package . --sign-key ~/.ssh/intent_sign_key.hex
```

### 4. Publish to Local Registry

```bash
# Publish the package
../intent publish dist/test-project-0.1.0.itpkg

# Or package and publish
../intent publish .
```

### 5. Test Install Flow

```bash
# Create install test directory
mkdir ../test-install && cd ../test-install

# Install from local registry
../intent install @scope/test-project@0.1.0

# Verify installation
ls intents/
```

## Testing the CLI

### Unit Tests

```bash
# Run all tests
go test -v ./...

# Run specific package tests
go test -v ./internal/cmd/...

# Run with coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Run only integration tests
go test -v -run Integration ./...
```

### Lint

```bash
# Run golangci-lint
golangci-lint run

# Fix issues
golangci-lint run --fix
```

## Testing Full Workflows

### Scenario 1: Simple Publish & Install

```bash
# Terminal 1: Start API
./mock-api

# Terminal 2: Create and publish
./intent init weather-app --app
cd weather-app
../intent package . --unsigned
../intent publish .

# Terminal 3: Install in different location
mkdir weather-consumer
cd weather-consumer
../intent install @scope/weather-app
```

### Scenario 2: Update and Republish

```bash
cd weather-app

# Update version in itpkg.json
sed -i 's/"0.1.0"/"0.1.1"/' itpkg.json

# Re-package and publish
../intent package . --unsigned
../intent publish .

# Install new version
../intent install @scope/weather-app@0.1.1
```

### Scenario 3: Signed Package

```bash
# Generate signing key
./gen_intent_key.sh

# Package with signature
./intent package . --sign-key private_key.hex

# Verify signature
./intent verify dist/weather-app-0.1.0.itpkg

# Publish
./intent publish dist/weather-app-0.1.0.itpkg
```

## Debugging

### Enable Debug Output

```bash
# All commands support --debug
./intent run intents/hello.itml --debug

./intent publish . --debug

./intent install @scope/package --debug
```

### Check Configuration

```bash
# See current config
./intent whoami

# See saved config
cat ~/.intent/config.yaml

# Check environment
echo "API: $INTENT_API_URL"
echo "Token: $INTENT_TOKEN"
```

### Mock API Logs

```bash
# See what packages are in registry
curl http://localhost:8080/v1/packages/search

# Resolve a package
curl "http://localhost:8080/v1/packages/resolve?spec=@scope/weather-app"

# Health check
curl http://localhost:8080/health
```

## Common Tasks

### Generate a New Signing Key

```bash
./gen_intent_key.sh
# Creates: private_key.hex, public_key.hex
```

### Create Multiple Test Projects

```bash
for i in 1 2 3; do
  ./intent init example-$i --app
done
```

### Test All Commands

```bash
# Create project
./intent init test --app && cd test

# Edit intent
echo 'intent "Test"
inputs: []
outputs: []
workflow:
  → log("Hello")
  → return(status="ok")' > intents/test.itml

# Run
../intent run intents/test.itml

# Package
../intent package . --unsigned

# Publish
../intent publish .

# (In another dir) Install
cd ../test-consumer
../../intent install @scope/test
```

## Troubleshooting

### Port 8080 Already in Use

```bash
# Find what's using it
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use different port
PORT=8081 ./mock-api
export INTENT_API_URL=http://localhost:8081
```

### Package Signature Mismatch

```bash
# Ensure you're using same key for signing
export INTENT_SIGN_KEY=./private_key.hex

# Verify package
./intent verify dist/package.itpkg
```

### Package Version Conflicts

```bash
# Clear registry by restarting mock-api
# Mock API stores packages in memory only

pkill mock-api
./mock-api
```

### Config Not Found

```bash
# Ensure config exists
mkdir -p ~/.intent/
touch ~/.intent/config.yaml

# Check permissions
ls -la ~/.intent/
```

## CI/CD Checks

### Before Committing

```bash
# Run linter
golangci-lint run

# Run tests
go test -v ./...

# Build for all platforms
GOOS=linux GOARCH=amd64 go build -o intent-linux ./cmd/intent
GOOS=darwin GOARCH=arm64 go build -o intent-darwin ./cmd/intent
GOOS=windows GOARCH=amd64 go build -o intent-windows.exe ./cmd/intent
```

### Create a Release

```bash
# Tag a version
git tag v0.4.0
git push origin v0.4.0

# GitHub Actions will automatically:
# 1. Run CI tests
# 2. Build for all platforms
# 3. Create GitHub release
# 4. Update Homebrew formula
```

## Performance Testing

### Benchmark Package Creation

```bash
# Create large project
for i in {1..100}; do
  echo "intent \"Intent $i\"
  inputs: []
  workflow:
    → log(\"Test $i\")
    → return(status=\"ok\")" > test/intents/intent-$i.itml
done

# Time package creation
time ./intent package test/ --unsigned
```

### Benchmark Publishing

```bash
# Time publishing
time ./intent publish test/dist/test-0.1.0.itpkg
```

## Additional Resources

- [User Guide](USER_GUIDE.md) - Complete usage guide
- [ITML Format](../README.md#itml-format) - Intent language reference
- [API Documentation](API.md) - Registry API endpoints
- [Package Format](itpkg_definition.md) - .itpkg file specification

