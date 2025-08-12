#!/bin/bash

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Get the backend workspace root (parent of scripts directory)
WORKSPACE_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# Change to the script directory
cd "$SCRIPT_DIR"

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed or not in PATH"
    exit 1
fi

# Check if buf is available for gRPC linting
if ! command -v buf &> /dev/null; then
    echo "⚠️  Warning: buf CLI tool not found. gRPC linting will be skipped."
    echo "   Install with: go install github.com/bufbuild/buf/cmd/buf@latest"
fi

echo "🚀 Starting Microservice Test and gRPC Lint Runner"
echo "📍 Workspace: $WORKSPACE_ROOT"
echo ""

# Run the Go script
go run main.go

# Capture the exit code
EXIT_CODE=$?

# Return to the original directory
cd "$WORKSPACE_ROOT"

# Exit with the same code
exit $EXIT_CODE
