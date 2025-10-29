# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.7] - 2025-10-29

### Fixed
- Fixed scaffold command to create required directories even if `itpkg.json` already exists
- Improved directory structure validation error messages

## [0.3.6] - 2025-10-29

### Fixed
- Fixed linter errors in test files (os.Chdir error checks)

## [0.3.5] - 2025-10-29

### Added
- **Standalone `intent package` command**: Creates signed `.itpkg` archives independently
- **.itpkg format v0.1**: Flat tar.gz structure with `itpkg.json`, `MANIFEST.sha256`, and `SIGNATURE`
- **ed25519 signing**: Cryptographic signature over MANIFEST.sha256 for package integrity
- **Required manifest**: `itpkg.json` with name, version, entry, policies, capabilities
- **Structure validation**: Validates required directories (`intents/`, `policies/`) and recommended ones
- **Scaffold support**: `--scaffold` flag generates `itpkg.json` and required directories
- **Signing options**: Support for `--sign-key`, `INTENT_SIGN_KEY` env var, or `--unsigned` flag
- **Package types**: Support for both `app` (with entry) and `lib` (without entry) packages
- **Policy validation**: Enforces `policies.security.network` for app packages
- **Key generation script**: `gen_intent_key.sh` for creating ed25519 signing keys

### Changed
- Package command now defaults to `.itpkg` format (removed tar.gz option)
- Package structure changed from nested `payload.tar.gz` to flat tar.gz archive
- Signing changed from HMAC-SHA256 to ed25519 for non-repudiation

### Fixed
- Fixed packaging to skip output archive file when packaging into source directory

## [0.3.4] - 2025-10-29

### Fixed
- Fixed packaging to prevent self-inclusion when output directory is inside source directory

## [0.3.3] - 2025-10-29

### Added
- Standalone `package` command with custom output path support

## [0.3.2] - 2025-10-28

### Added
- Custom ITML DSL format parser (primary format)
- Multi-format support: ITML (primary), JSON (fallback), YAML (fallback)
- Workflow execution: `→ log()` and `→ return()` commands
- Dual template syntax: `{name}` and `{{name}}` support
- Enhanced parameter parsing: `name (type) default="value"` syntax
- Default parameter value support in template processing

### Fixed
- Fixed template processing to use default parameter values when not provided

## [0.3.1] - 2025-10-28

### Added
- YAML parsing support using `gopkg.in/yaml.v3`
- YAML struct tags for all parser structures

## [0.3.0] - 2025-10-18

### Added
- Enhanced shell completion support for zsh, bash, fish, and PowerShell
- Automatic completion installation via Homebrew
- Development build support with commit information
- Comprehensive README with installation and usage instructions
- Enhanced login command that saves both token and api_url to ~/.intent/config.yaml
- Exponential backoff retry strategy for improved network resilience
- Friendlier DNS/network error messages with actionable hints
- User-Agent header for better observability and debugging

### Changed
- Improved project structure following Go best practices
- Moved subcommands to internal package for better organization
- Enhanced versioning system with short and long format support
- Login command now prompts for API URL and saves complete configuration
- Retry strategy now uses exponential backoff (500ms → 1s → 2s → 4s, max 5s)
- Network error handling provides clearer, more actionable error messages

### Fixed
- Fixed tree command functionality
- Improved release and Homebrew integration
- Enhanced completion system reliability
- Tamed noisy Resty retry logs by disabling logger in non-debug mode
- Whoami command now shows helpful "not logged in" message when no token present

## [0.2.9] - 2025-10-18

### Added
- Automatic shell completion installation via Homebrew formula

## [0.2.8] - 2025-10-18

### Added
- Basic shell completion support

## [0.2.7] - 2025-10-18

### Fixed
- Various bug fixes and improvements

## [0.2.6] - 2025-10-18

### Fixed
- Release and Homebrew integration fixes

## [0.2.5] - 2025-10-18

### Fixed
- Release and Homebrew integration fixes

## [0.2.4] - 2025-10-18

### Fixed
- Homebrew formula improvements

## [0.2.3] - 2025-10-18

### Fixed
- Homebrew formula fixes

## [0.2.2] - 2025-10-18

### Fixed
- Tree command fixes

## [0.2.1] - 2025-10-18

### Changed
- Moved subcommands to internal package structure

## [0.2.0] - 2025-10-18

### Added
- Core CLI functionality
- Login, publish, install, search, and whoami commands
- Configuration management
- HTTP client for API communication
- Tar.gz packaging utilities

## [0.1.0] - Initial Release 2025-10-18

### Added
- Initial project setup
- Basic CLI structure
- Go module configuration
- Build system with Makefile
