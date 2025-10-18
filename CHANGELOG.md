# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
