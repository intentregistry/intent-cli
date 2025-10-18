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

- `INTENT_API_URL` (default: `https://api.intentregistry.com`)
- `INTENT_TOKEN`

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

## Development

To build locally:

```bash
go mod tidy
go run ./cmd/intent --help
make build
make pack-darwin-arm64
make checksum
```

### Development Builds

For development and testing, you can create dev builds that include commit information:

```bash
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
