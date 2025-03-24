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

.PHONY: add-copyright-headers
add-copyright-headers:
	@bash -c ' \
		SPDX1="// SPDX-FileCopyrightText: © 2025 Industria de Diseño Textil S.A. INDITEX"; \
		SPDX2="// SPDX-License-Identifier: Apache-2.0"; \
		find . -type f -name "*.go" | while read file; do \
			line1=$$(sed -n "1p" $$file); \
			line2=$$(sed -n "2p" $$file); \
			if [[ "$$line1" =~ ^//\ SPDX-FileCopyrightText: && "$$line2" =~ ^//\ SPDX-License-Identifier: ]]; then \
				sed -i "1s|.*|$$SPDX1|" $$file; \
				sed -i "2s|.*|$$SPDX2|" $$file; \
			else \
				{ echo "$$SPDX1"; echo "$$SPDX2"; cat $$file; } > $$file.tmp && mv $$file.tmp $$file; \
			fi \
		done'


.PHONY: get-version
get-version:
	@echo $(PROJECT_VERSION)
