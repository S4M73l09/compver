# Security Policy

## Supported versions

Security fixes are prioritized for the latest code in the `pre-release` branch
and for the latest stable release available from `main`.

Older versions may not receive security fixes. Users should upgrade to the
latest supported version before reporting whether an issue is still present.

## Reporting a vulnerability

Please do not report security vulnerabilities through public GitHub issues,
pull requests, or discussions.

Use GitHub's private vulnerability reporting feature when it is available:

<https://github.com/S4M73l09/compver/security/advisories/new>

If private vulnerability reporting is unavailable, contact the repository
maintainer privately through GitHub before disclosing the issue publicly.

Please include, when possible:

- A concise description of the vulnerability.
- The affected version, commit, or branch.
- Reproduction steps or a minimal proof of concept.
- The expected and observed behavior.
- The potential impact.
- Any suggested mitigation or fix.

Do not include secrets, credentials, private keys, or sensitive infrastructure
data in the report.

## Response process

The maintainer will attempt to acknowledge a valid report, investigate the
impact, coordinate a fix, and publish an advisory when appropriate.

Please allow time for investigation and coordinated disclosure before making a
security issue public.

## Scope

Security reports may include issues involving:

- Unsafe filesystem traversal or file handling.
- Unexpected network requests or insecure provider communication.
- Cache poisoning or unsafe cache contents.
- Credential, token, or secret exposure.
- Dependency vulnerabilities introduced by Compver.
- Command execution or untrusted input handling.
