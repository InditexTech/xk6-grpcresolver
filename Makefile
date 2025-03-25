PROJECT_VERSION := 1.0.0

XK6_VERSION := v0.13.4
XK6_BINARY := $(shell command -v xk6 2> /dev/null)

GOLANGCI_VERSION := v1.64.5
GOLANGCI_BINARY := $(shell command -v golangci-lint 2> /dev/null)

.DEFAULT_GOAL := all

.PHONY: all
all: format lint test build

.PHONY: deps
deps:
	@if [ -z "$(XK6_BINARY)" ]; then \
		echo "Installing xk6..."; \
		go install go.k6.io/xk6/cmd/xk6@$(XK6_VERSION); \
	else \
		echo "xk6 is already installed."; \
	fi

	@if [ -z "$(GOLANGCI_BINARY)" ]; then \
			echo "Installing golangci-lint..."; \
			go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION); \
	else \
		echo "golangci-lint is already installed."; \
	fi

.PHONY: build
build: deps
	@echo "Building k6 with grpcresolver extension..."
	@xk6 build --with github.com/InditexTech/xk6-grpcresolver=.

.PHONY: run
run: build
	@echo "Running example in docker..."
	@docker compose -f docker/docker-compose.yaml up --abort-on-container-exit

.PHONY: verify
verify: format lint test

.PHONY: test
test:
	@echo "Running unit tests..."
	@go clean -testcache && go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: tidy
tidy:
	@echo "Running go mod tidy..."
	@go mod tidy

.PHONY: format
format:
	@echo "Running go fmt..."
	go fmt ./...

.PHONY: lint
lint: deps
	@echo "Running golangci-lint..."
	@golangci-lint run

.PHONY: reuse-deps
reuse-deps:
	@if [ -z "reuse" ]; then \
		echo "Installing reuse tool..."; \
		pip3 install --user reuse ;\
	else \
		echo "reuse is already installed."; \
	fi


.PHONY: add-copyright-headers
reuse-annotate: reuse-deps
	@echo "Adding copyright headers..."
	@reuse annotate --copyright "Industria de Diseño Textil S.A. INDITEX" --license "Apache-2.0" --year "$$(date +%Y)" --merge-copyrights *.go **/**/*.go
	@reuse lint

.PHONY: get-version
get-version:
	@echo $(PROJECT_VERSION)
