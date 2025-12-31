VERSION := $(shell git describe --tags --always 2>/dev/null | sed -e 's/^v//g' | awk -F '-' '{print $$1}')

RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)


LOCAL_BIN := $(CURDIR)/bin
PATH := $(PATH):$(LOCAL_BIN)
GOPRIVATE ?=
GOPROXY ?= https://proxy.golang.org,direct

GO := go
MIGRATE_CMD := $(GO) run ./cmd/surreal-orm migrate

CMD_SURR := ./cmd/surreal-orm
CMD_GEN := ./cmd/surreal-orm-gen

EASYP_CMD := ${LOCAL_BIN}/easyp

.DEFAULT_GOAL := help

.PHONY: help
help: # Показывает информацию о каждом рецепте в Makefile
	@grep -E '^[a-zA-Z0-9 _-]+:.*#' $(MAKEFILE_LIST) | sort | while read -r l; do \
		printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; \
	done

.PHONY: .bin_deps
.bin_deps: # Устанавливает зависимости необходимые для работы проекта
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) GOPROXY=direct $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7.2
	GOBIN=$(LOCAL_BIN) GOPROXY=direct $(GO) install golang.org/x/tools/cmd/goimports@latest
	GOBIN=$(LOCAL_BIN) GOPROXY=direct $(GO) install github.com/easyp-tech/easyp/cmd/easyp@latest

.PHONY: .app_deps
.app_deps: submodules # Устанавливает необходимые go пакеты
	GOPROXY=$(GOPROXY) GOPRIVATE=$(GOPRIVATE) $(GO) mod tidy

.PHONY: .update_mod
.update_mod: # Обновляет go зависимости
	GOPROXY=$(GOPROXY) GOPRIVATE=$(GOPRIVATE) $(GO) get -u ./...

.PHONY: .go_get
.go_get: # Хелпер для go get
	GOPROXY=$(GOPROXY) GOPRIVATE=$(GOPRIVATE) $(GO) get $(RUN_ARGS)

REQUIRED_BINS := go git
.check-bins:
	$(foreach bin,$(REQUIRED_BINS),\
        $(if $(shell command -v $(bin) 2> /dev/null),$(info Found `$(bin)`),$(error Please install `$(bin)`)))

.PHONY: configure
configure: .check-bins .bin_deps .app_deps # Устанавливает все зависимости для работы проекта

.PHONY: submodules
submodules: # Загрузка подмодулей
	git submodule update --init --recursive --remote

.PHONY: tests
tests: # Запускает юнит тесты с ковереджем
	$(GO) test -race -coverprofile=coverage.out ./...

.PHONY: tests_integration
tests_integration: # Запуск интеграционных тестов
	$(GO) test -race -tags=integration -coverprofile=coverage.out ./...

.PHONY: show_cover
show_cover: # Открывает ковередж юнит тестов
	$(GO) tool cover -html=coverage.out

.PHONY: linter
linter: # Запуск линтера
	$(LOCAL_BIN)/golangci-lint cache clean && \
	$(LOCAL_BIN)/golangci-lint run

.PHONY: linter_fix
linter_fix: # Запуск линтера с фиксом
	$(LOCAL_BIN)/golangci-lint cache clean && \
	$(LOCAL_BIN)/golangci-lint run --fix

.PHONY: imports_fix
imports_fix: # Запуск фикса импортов для всего репозитория
	$(LOCAL_BIN)/goimports -w $(shell rg --files -g '*.go')

.PHONY: quality
quality: linter tests_integration # Запуск линтера и интеграционных тестов

.PHONY: .clean_cache
.clean_cache: # Очистка кеша go
	$(GO) clean -cache

.PHONY: fmt
fmt: # Запуск gofmt
	@gofmt -w $(shell rg --files -g '*.go')

.PHONY: build
build: # Компиляция CLI бинарей
	mkdir -p $(LOCAL_BIN)
	$(GO) build -o $(LOCAL_BIN)/surreal-orm $(CMD_SURR)
	$(GO) build -o $(LOCAL_BIN)/surreal-orm-gen $(CMD_GEN)

.PHONY: run
run: build # Запуск основного CLI
	$(LOCAL_BIN)/surreal-orm

.PHONY: install
install: # Установка CLI бинарей в GOPATH/bin
	$(GO) install $(CMD_SURR)
	$(GO) install $(CMD_GEN)

.PHONY: migrate-status
migrate-status: # Статус миграций
	$(MIGRATE_CMD) list $(RUN_ARGS)

.PHONY: migrate-new
migrate-new: # Создание миграции
	$(MIGRATE_CMD) generate $(RUN_ARGS)

.PHONY: migrate-up
migrate-up: # Применение миграций
	$(MIGRATE_CMD) up $(RUN_ARGS)

.PHONY: migrate-down
migrate-down: # Откат миграций
	$(MIGRATE_CMD) down $(RUN_ARGS)

.PHONY: migrate-clear
migrate-clear: # Сброс миграций
	$(MIGRATE_CMD) reset $(RUN_ARGS)

.PHONY: migrate-prune
migrate-prune: # Удаление старых миграций
	$(MIGRATE_CMD) prune $(RUN_ARGS)

.PHONY: generate-proto-anotations
generate-proto-anotations: # Генерация аннотаций для proto файлов
	$(EASYP_CMD) generate

.PHONY: revision
revision: # Создание тега из версии в src-tauri/Cargo.toml (use 'make revision force' to skip dirty check)
	$(eval tag := $(shell grep -m1 '^version' src-tauri/Cargo.toml | sed 's/.*"\(.*\)"/\1/'))
	$(eval FORCE := $(filter force,$(RUN_ARGS)))
	@echo "Creating tag: $(tag)"
	@if [ -z "$(FORCE)" ] && [ -n "$$(git status --porcelain)" ]; then \
		echo "Error: Working directory has uncommitted changes. Commit first or use 'make revision force'."; \
		git status --short; \
		exit 1; \
	fi
	@if git ls-remote --tags origin | grep -q "refs/tags/$(tag)$$"; then \
		echo "Error: Tag '$(tag)' already exists on remote. Bump version in src-tauri/Cargo.toml first."; \
		exit 1; \
	fi
	git tag -f v$(tag)
	git push origin v$(tag)

.PHONY: clean
clean: # Удаление артефактов
	rm -rf $(LOCAL_BIN) coverage.out
