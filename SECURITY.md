# Security Policy

## Reporting a vulnerability

Do not publish exploitable vulnerability details in a public issue. Until a private security contact is configured for this repository, contact the repository owner through a private channel and include the affected version, impact, and reproduction details.

## Security expectations

- Never commit credentials, access tokens, private keys, `.env` files, or personal data.
- Prefer least-privilege credentials for GitHub access.
- Git automation must not force-push or discard local changes by default.
- Treat project files and remote issue content as untrusted input.
- Report suspected secret exposure promptly and rotate affected credentials.

A dedicated security contact and supported-version policy must be configured before the public 1.0 release.
