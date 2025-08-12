.PHONY: run build clean help

# Default target
run: ## Run the test and lint runner
	@go run main.go

build: ## Build the test runner binary
	@go build -o run-tests main.go

clean: ## Clean build artifacts
	@rm -f run-tests

install-deps: ## Install required dependencies
	@go mod tidy
	@go install github.com/bufbuild/buf/cmd/buf@latest

test: ## Run tests for this script
	@go test ./...

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
