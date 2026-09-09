# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in StellarYard CLI, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please email: **security@stellaryard.dev**

Include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

## Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial assessment**: Within 1 week
- **Fix or mitigation**: Depends on severity, typically within 2 weeks

## Scope

This security policy applies to:
- The `stellaryard-cli` Go application
- API client communication
- Output formatting and data handling

## Out of Scope

- stellaryard-core vulnerabilities — report to the core repo
- Third-party dependencies — report to their respective maintainers

## Key Security Considerations

### API Communication

The CLI communicates with stellaryard-core over HTTP/WS. Ensure the core URL is correct and that the connection is to a trusted local instance.

### Exit Codes

Exit codes are part of the security contract. Misclassified exit codes can cause CI/CD pipelines to take incorrect actions.

### Credential Handling

The CLI does not store credentials. All authentication is delegated to core.

## Disclosure Policy

We follow responsible disclosure. We will:
- Credit reporters (unless they prefer anonymity)
- Not pursue legal action for good-faith security research
- Work with reporters on disclosure timing
