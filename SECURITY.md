# Security Policy

## Supported Versions

We release security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please report it responsibly.

### How to Report

**Please do not open a public issue for security vulnerabilities.**

Instead, please email us at: **security@coddy.dev** (coming soon)

Or use GitHub's private vulnerability reporting feature:
1. Go to the repository's Security tab
2. Click "Report a vulnerability"
3. Fill out the form with details

### What to Include

Please include as much of the following information as possible:

- **Description**: Clear description of the vulnerability
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Impact**: What could an attacker do with this vulnerability?
- **Affected Versions**: Which versions are affected?
- **Mitigation**: Any suggested fixes or workarounds
- **Proof of Concept**: If applicable, a minimal proof of concept

### Response Timeline

We aim to respond to security reports within:

- **24 hours**: Acknowledgment of receipt
- **72 hours**: Initial assessment
- **7 days**: Plan for fix or mitigation
- **30 days**: Fix released (depending on severity)

### Security Measures

Coddy implements several security measures:

#### Sandbox Isolation
- Docker containers with `--network=none` by default
- Read-only root filesystem
- Resource limits (memory, CPU)
- Execution timeouts
- Non-root user execution

#### Known Limitations

The following are known limitations that users should be aware of:

1. **Subprocess Sandbox**: The subprocess sandbox (used in development) is NOT secure and should never be used in production or with untrusted code.

2. **Docker Sandbox**: While Docker provides isolation, it is not perfect. For high-security environments, consider using gVisor or Firecracker.

3. **Network Access**: Network access can be enabled per-session. This increases risk and should be used carefully.

4. **File System**: Sandboxes have access to a writable `/home/user` directory. Sensitive files should not be mounted into sandboxes.

### Best Practices

When deploying Coddy:

1. **Use Docker Sandbox**: Never use subprocess sandbox in production
2. **Disable Network**: Keep `--network=none` unless specifically needed
3. **Resource Limits**: Set appropriate memory and CPU limits
4. **Timeouts**: Use reasonable execution timeouts
5. **Regular Updates**: Keep Docker images and dependencies updated
6. **Monitoring**: Monitor for unusual resource usage or execution patterns

### Security Updates

Security updates will be released as patch versions (e.g., 0.1.1). We recommend:

- Enabling Dependabot alerts on your fork
- Watching this repository for releases
- Updating promptly when security fixes are released

### Hall of Fame

We thank the following security researchers for responsibly disclosing vulnerabilities:

*None yet — be the first!*

---

Last updated: 2026-04-12
