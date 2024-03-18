# Определение переменных
BINARY_NAME = filmLibApp
GOLANGCI_LINT_VERSION = v1.43.0

# Цель по умолчанию (выполняется при вызове 'make' без конкретной цели)
all: build lint test

# Цель построения (компилирует код Go)
build:
	go build -o $(BINARY_NAME) ./cmd

# Цель проверки кода (lint) (проверяет код на соответствие стилю и лучшим практикам)
lint:
	wget https://github.com/golangci/golangci-lint/releases/download/$(GOLANGCI_LINT_VERSION)/golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64.tar.gz
	tar -xvf golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64.tar.gz
	./golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64/golangci-lint run
	rm -rf golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64.tar.gz golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64

# Цель тестирования (запускает тесты Go)
test:
	go test -v ./...

# Цель запуска (запускает скомпилированный двоичный файл)
run:
	go run ./cmd/main.go

# Цель форматирования (форматирует код Go с помощью gofmt)
format:
	gofmt -s -w .

# Цель очистки (удаляет скомпилированный двоичный файл)
clean:
	rm -f $(BINARY_NAME)

.PHONY: all build lint test run format clean
