# BPUF - Bipartite Union Find Library
# Makefile for common development tasks

# Variables
BINARY_NAME=bpuf-clickhouse
BIN_DIR=bin
DIST_DIR=dist

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags="-s -w"

.PHONY: clean
clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)
	rm -f *.xml

.PHONY: deps
deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: test
test: ## Run all tests
	$(GOTEST) -v ./...

.PHONY: test-race
test-race: ## Run tests with race detection
	$(GOTEST) -race -short ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run

.PHONY: lint-fix  
lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --fix

.PHONY: nilaway
nilaway: ## Run nilaway nil pointer analysis
	nilaway ./...

.PHONY: lint-all
lint-all: lint nilaway ## Run all linters (golangci-lint + nilaway)

.PHONY: fmt
fmt: ## Format code
	$(GOCMD) fmt ./...
	goimports -w .

.PHONY: vet
vet: ## Run go vet
	$(GOCMD) vet ./...

.PHONY: bench
bench: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./...

.PHONY: check
check: deps fmt vet lint-all test

.PHONY: ci
ci: check test-race ## Run CI checks locally

.PHONY: pre-commit
pre-commit: fmt vet lint-all test ## Run pre-commit checks

.PHONY: dev-setup
dev-setup: ## Set up development environment
	@echo "Installing development dependencies..."
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint/v2@latest
	$(GOCMD) install golang.org/x/tools/cmd/goimports@latest
	$(GOCMD) install go.uber.org/nilaway/cmd/nilaway@latest
	$(GOMOD) download
	@echo "Development environment setup done"

.PHONY: release
release: ## Create and push a new release tag
	@read -p "Enter version (e.g., v1.0.0): " version; \
	git tag $$version && \
	git push origin $$version && \
	gh release create $$version --generate-notes
