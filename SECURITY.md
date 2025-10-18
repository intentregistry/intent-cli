# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security bugs seriously. We appreciate your efforts to responsibly disclose your findings, and will make every effort to acknowledge your contributions.

### How to Report a Security Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **Email**: Send details to [security@intentregistry.com](mailto:security@intentregistry.com)
2. **GitHub Security Advisories**: Use GitHub's private vulnerability reporting feature if available

### What to Include

When reporting a vulnerability, please include:

- **Description**: A clear description of the vulnerability
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Impact**: Potential impact of the vulnerability
- **Environment**: OS, Go version, and any other relevant environment details
- **Suggested Fix**: If you have ideas for how to fix the issue

### What to Expect

- **Acknowledgment**: We will acknowledge receipt of your report within 48 hours
- **Initial Assessment**: We will provide an initial assessment within 5 business days
- **Regular Updates**: We will keep you informed of our progress
- **Resolution**: We will work with you to resolve the issue and coordinate disclosure

### Disclosure Timeline

- **Immediate**: Critical vulnerabilities (remote code execution, authentication bypass)
- **7 days**: High severity vulnerabilities (data exposure, privilege escalation)
- **30 days**: Medium severity vulnerabilities
- **90 days**: Low severity vulnerabilities

## Security Best Practices

### For Users

- **Keep Updated**: Always use the latest version of the CLI
- **Secure Configuration**: Store API tokens securely and use environment variables
- **Network Security**: Use HTTPS endpoints and verify SSL certificates
- **Access Control**: Limit API token permissions to minimum required scope

### For Developers

- **Dependencies**: Keep dependencies updated and scan for vulnerabilities
- **Input Validation**: Validate all user inputs and API responses
- **Error Handling**: Avoid exposing sensitive information in error messages
- **Logging**: Be careful not to log sensitive data (tokens, passwords, etc.)
- **Testing**: Include security testing in your development process

## Security Features

### Authentication & Authorization

- API token-based authentication
- Secure token storage in configuration files
- Environment variable support for sensitive data

### Data Protection

- HTTPS-only communication with API endpoints
- No storage of sensitive user data locally
- Secure handling of temporary files

### Input Validation

- Validation of all command-line arguments
- Sanitization of file paths and URLs
- Protection against path traversal attacks

## Known Security Considerations

### API Token Security

- Tokens are stored in plain text in configuration files
- Consider using environment variables in production environments
- Rotate tokens regularly
- Use tokens with minimal required permissions

### File System Access

- The CLI reads and writes files in user-specified directories
- Be cautious when using with untrusted file paths
- Consider using sandboxed environments for automated usage

### Network Communication

- All API communication uses HTTPS
- Certificate validation is performed by default
- No custom certificate handling (relies on system trust store)

## Security Updates

Security updates will be released as soon as possible after a vulnerability is confirmed and a fix is available. We will:

1. Release a patch version with the security fix
2. Update the changelog with security-related information
3. Notify users through appropriate channels
4. Coordinate with security researchers on disclosure timing

## Security Contact

For security-related questions or concerns:

- **Email**: [security@intentregistry.com](mailto:security@intentregistry.com)
- **Response Time**: We aim to respond within 48 hours
- **PGP Key**: Available upon request for encrypted communication

## Acknowledgments

We would like to thank the security researchers and community members who help keep Intent CLI secure by responsibly disclosing vulnerabilities.

## Legal

This security policy is provided for informational purposes only. By using Intent CLI, you agree to use it at your own risk and in accordance with the project's license terms.
