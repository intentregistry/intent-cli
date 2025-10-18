# Intent CLI

The official command-line interface for publishing and installing
**Intents** on [IntentRegistry](https://intentregistry.com).

## Installation

You can install the CLI using Homebrew:

``` bash
brew tap intentregistry/homebrew-tap
brew install intent-cli
intent --help
```

## Usage

``` bash
intent login
intent publish [path] --private --tag beta --message "first release"
intent install @scope/name[@version] --dest intents
intent search "vector embeddings"
intent whoami
```

## Configuration

Configuration file path:

    ~/.intent/config.yaml

You can also set environment variables:

-   `INTENT_API_URL` (default: `https://api.intentregistry.com`)
-   `INTENT_TOKEN`

## Development

To build locally:

``` bash
go mod tidy
go run ./cmd/intent --help
make build
make pack-darwin-arm64
make checksum
```

### Development Builds

For development and testing, you can create dev builds that include commit information:

``` bash
make build-dev
make pack-dev
```

Dev builds will show version as `dev+<commit-hash>` when you run `intent --version`.

## Project Structure

    intent-cli/
    ├─ cmd/intent/              # main entrypoint
    ├─ internal/cmd/            # subcommands
    ├─ internal/config/         # configuration management
    ├─ internal/httpclient/     # API HTTP client
    ├─ internal/pack/           # tar.gz packaging utilities
    ├─ internal/version/        # version and build metadata
    ├─ .goreleaser.yaml         # release automation
    ├─ go.mod                   # Go module
    ├─ Makefile                 # build tasks
    └─ README.md

## License

MIT License © 2025 [IntentRegistry](https://intentregistry.com)
