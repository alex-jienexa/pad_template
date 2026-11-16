# Contributing to PAD Core

Thank you for contributing. PAD is intended to remain platform-independent, predictable, and safe around project data.

## Before opening an issue or pull request

- Search existing issues and pull requests.
- For substantial changes, discuss the proposal before implementation.
- Keep changes focused and link them to a task or requirement when available.

## Development setup

1. Install Go 1.23 or newer and Git.
2. Clone the repository.
3. Run `go test ./...` and `go vet ./...`.

## Change guidelines

- Follow standard Go formatting (`gofmt`).
- Prefer the standard library unless a dependency has a clear benefit.
- Keep packages cohesive and avoid unnecessary global state.
- Treat file operations and Git operations as potentially destructive.
- Provide tests for new behavior and regression tests for bug fixes.
- Update documentation and examples when behavior or configuration changes.
- Never include credentials, tokens, private keys, or real user data.

## Pull requests

A pull request should explain the problem, the chosen solution, compatibility implications, and verification performed. CI must pass before merging. Breaking changes require explicit documentation and a versioning decision.

## Commit messages

Use Conventional Commits where practical, for example:

```text
feat(cli): add project initialization
fix(git): refuse backup with divergent history
docs(methodology): clarify recovery procedure
```

## License

The repository license is pending owner selection. Contribution licensing terms must be finalized before accepting external contributions.
