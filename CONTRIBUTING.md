# Contributing

Thanks for your interest in improving this project. This document describes the expected workflow.

## Prerequisites

- Go 1.25+
- Make

## Development setup

```bash
make tidy
make configure
make tests
make linter
```

## Makefile commands

Common targets:

```text
configure        Install tools and Go dependencies
tests            Run unit tests with coverage
tests_integration Run integration tests (tags=integration)
linter           Run golangci-lint
linter_fix       Run golangci-lint with fixes
imports_fix      Run goimports on all Go files
fmt              Run gofmt on all Go files
build            Build CLI binaries into ./bin
install          Install CLI binaries into GOPATH/bin
```

Migration helpers (use CLI flags for connection details):

```text
migrate-status   Show migration status
migrate-new      Generate a migration (CLI)
migrate-up       Apply migrations (CLI)
migrate-down     Roll back migrations (CLI)
migrate-clear    Reset migrations (CLI)
migrate-prune    Prune migrations (CLI)
```

## Coding rules

- Keep changes focused and small.
- Run `gofmt` (or `make fmt`) on Go files.
- Add tests for bug fixes or new behavior.
- Avoid breaking the public API without discussion.

## Submitting changes

1. Fork the repository and create a branch from `main`.
2. Make your changes and add tests.
3. Run `make test` and `make lint`.
4. Open a PR with a clear description and rationale.

## Reporting issues

Use the issue templates and include reproduction steps, expected vs actual behavior, and environment details.
