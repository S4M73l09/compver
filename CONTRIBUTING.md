# Contributing to Compver

Thank you for your interest in contributing to Compver.

Compver is an open-source project written in Go. Its architecture is designed
around a generic analysis engine, independent adapters, providers, version
logic, and a command-line interface.

## Before starting

For small fixes, documentation changes, and tests, you can open a pull
request directly.

For new adapters, architectural changes, new dependencies, or changes that
affect the public CLI, open an issue first so the proposed direction can be
discussed before implementation.

Please check existing issues and pull requests before starting duplicated work.

## Development workflow

1. Fork the repository and create a focused branch from `pre-release`.
2. Use a descriptive branch name, for example:
   - `feature/terraform-constraints`
   - `fix/cache-timeout`
   - `docs/contributing-guide`
3. Keep each pull request focused on one change.
4. Add or update tests for behavior changes.
5. Update `README.md`, `ARCHITECTURE.md`, or `CHANGELOG.md` when appropriate.
6. Run the local validation commands before opening a pull request.
7. Explain the design and trade-offs in the pull request description.

The `main` branch is intended for stable code. The `pre-release` branch is
used for active development and validation of upcoming changes.

## Local validation

Run these commands from the repository root:

```bash
gofmt -w path/to/changed.go
go test ./...
go vet ./...
```

Do not use `gofmt -w .` if your working tree contains unrelated local changes
that must not be reformatted. In that case, format only the files changed by
your contribution.

## Architecture guidelines

- Keep the generic engine independent from specific tools.
- Put tool-specific behavior inside `internal/adapters/`.
- Keep external version lookups inside providers.
- Reuse `internal/version` for parsing, comparison, selection, and constraints.
- Keep terminal presentation inside `internal/cli/`.
- Prefer the shared cache instead of implementing a tool-specific cache.
- Avoid adding dependencies when the standard library is sufficient.
- Do not include credentials, tokens, private keys, or sensitive project data.

## Pull requests

A pull request should:

- Explain the problem and the proposed solution.
- Include tests for new or changed behavior.
- Keep unrelated refactors out of the change.
- Document user-visible CLI changes.
- Mention known limitations or follow-up work.
- Pass the required repository checks.

Maintainers may request changes, close, or defer a pull request when it is
outside the project scope, duplicates existing work, lacks sufficient tests,
introduces unnecessary complexity, or does not provide a clear benefit to the
project.

## License

By contributing to Compver, you agree that your contribution may be
distributed under the Apache License 2.0 used by this project.


# Thanks^^