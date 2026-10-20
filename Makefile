# Common project commands. Run `make help` to list them.

VERSION := $(shell tr -d '[:space:]' < VERSION)
IMAGE   ?= docs-drift:local

.DEFAULT_GOAL := help

.PHONY: help build test lint modules vulncheck check validate image clean

help: ## List available targets
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build bin/docs-drift with the version from VERSION
	sh scripts/build.sh bin/docs-drift

test: ## Run tests with the coverage floor
	sh scripts/coverage.sh

lint: ## Check formatting and run golangci-lint
	test -z "$$(gofmt -l .)" || { gofmt -l .; exit 1; }
	golangci-lint run ./...

modules: ## Check go.mod and go.sum are tidy and verified
	go mod tidy -diff
	go mod verify

vulncheck: ## Scan dependencies with govulncheck
	govulncheck ./...

check: lint modules test vulncheck ## Run every quality gate CI runs

validate: ## Validate the example registry
	go run ./cmd/docs-drift registry validate examples/sources.yaml

image: ## Build the container image
	docker build -t $(IMAGE) --build-arg VERSION=$(VERSION) .

clean: ## Remove build and coverage output
	rm -rf bin coverage.out
