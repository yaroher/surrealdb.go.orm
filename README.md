# surrealdb.go.orm

SurrealDB ORM toolkit for Go: schema definitions, migrations, and query building with a CLI.

## Features

- Query builder for SurrealQL
- ORM code generation from Go source
- Migration management against SurrealDB

## Requirements

- Go 1.25+

## Install

```bash
go install github.com/yaroher/surrealdb.go.orm/cmd/surreal-orm@latest
go install github.com/yaroher/surrealdb.go.orm/cmd/surreal-orm-gen@latest
```

## CLI usage

```bash
surreal-orm generate --dir .

surreal-orm migrate init --dsn <dsn> --ns <namespace> --db <database>
surreal-orm migrate generate --dsn <dsn> --ns <namespace> --db <database> --code . --name init
surreal-orm migrate up --dsn <dsn> --ns <namespace> --db <database>
surreal-orm migrate down --dsn <dsn> --ns <namespace> --db <database> --steps 1
surreal-orm migrate list --dsn <dsn> --ns <namespace> --db <database>
surreal-orm migrate reset --dsn <dsn> --ns <namespace> --db <database>
surreal-orm migrate prune --dsn <dsn> --ns <namespace> --db <database>
```

## Development

```bash
make configure
make tests
make linter
make fmt
make build
```

## Makefile targets

```text
configure        Install tools and Go dependencies
tests            Run unit tests with coverage
tests_integration Run integration tests (tags=integration)
show_cover       Open coverage report
linter           Run golangci-lint
linter_fix       Run golangci-lint with fixes
imports_fix      Run goimports on all Go files
fmt              Run gofmt on all Go files
build            Build CLI binaries into ./bin
install          Install CLI binaries into GOPATH/bin
migrate-status   Show migration status
migrate-new      Generate a migration (CLI)
migrate-up       Apply migrations (CLI)
migrate-down     Roll back migrations (CLI)
migrate-clear    Reset migrations (CLI)
migrate-prune    Prune migrations (CLI)
revision         Create and push a git tag (tag=... required)
```

## License

Apache-2.0. See `LICENSE`.
