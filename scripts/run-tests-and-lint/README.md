# Microservice Test and gRPC Lint Runner

This Go script automatically runs test suites and gRPC linting for all microservices in the backend workspace, then displays the results in a clear, readable format.

## Features

- 🧪 **Automatic Test Detection**: Automatically detects which services have test files
- 🔍 **gRPC Lint Detection**: Automatically detects which services have gRPC/proto files
- ⏱️ **Timing Information**: Shows how long each test suite and lint run takes
- 📊 **Clear Results**: Displays results with emojis and clear formatting
- 🚨 **Failure Details**: Shows detailed error information for failed tests/lint
- 🔄 **Exit Codes**: Returns appropriate exit codes for CI/CD integration
- 📈 **Progress Bars**: Visual progress indicators for overall completion
- ⏳ **Spinner Animation**: Animated spinners during test execution
- 📊 **Test Statistics**: Detailed counts of passed, failed, and skipped tests
- 🎯 **JSON Test Output**: Uses Go's JSON test format for better data parsing

## Supported Services

The script automatically checks these microservices:
- `auth` - Authentication service
- `mail` - Mail service
- `mail-server` - Mail server service
- `productivity` - Productivity service
- `grpc` - gRPC definitions

## Prerequisites

- Go 1.21 or later
- `buf` CLI tool (for gRPC linting) - install with `go install github.com/bufbuild/buf/cmd/buf@latest`

## Important

**This script must be run from the backend directory**. The script will automatically detect if it's running from the wrong directory and provide a helpful error message.

## Usage

### From the backend directory:

```bash
# Run directly
go run scripts/run-tests-and-lint/main.go

# Or navigate to the script directory first
cd scripts/run-tests-and-lint
go run main.go
```

### Build and run:

```bash
cd scripts/run-tests-and-lint
go build -o run-tests main.go
./run-tests
```

### From anywhere in the backend workspace:

```bash
go run ./scripts/run-tests-and-lint
```

## Output Example

```
🚀 Running Microservice Tests and gRPC Linting
==================================================

🔍 Processing auth...
  🧪 Running tests for auth...
  ✅ Tests passed for auth (2.34s)
  🔍 Running gRPC lint for auth...
  ✅ gRPC lint passed for auth (0.45s)

🔍 Processing mail...
  🧪 Running tests for mail...
  ✅ Tests passed for mail (1.87s)
  🔍 Running gRPC lint for mail...
  ✅ gRPC lint passed for mail (0.32s)

==================================================
📊 RESULTS SUMMARY
==================================================

🏗️  AUTH
   ---
   ✅ Tests: PASSED (2.34s)
   ✅ gRPC Lint: PASSED (0.45s)

🏗️  MAIL
   ---
   ✅ Tests: PASSED (1.87s)
   ✅ gRPC Lint: PASSED (0.32s)

==================================================
📈 SUMMARY
==================================================
Total Services: 2
Tests Passed: 2
Tests Failed: 0
gRPC Lint Passed: 2
gRPC Lint Failed: 0
Total Duration: 4.98s

✅ All checks passed successfully!
```

## Exit Codes

- `0` - All tests and linting passed successfully
- `1` - One or more tests or linting checks failed

## Configuration

The script automatically:
- Detects services based on directory structure
- Identifies test files (`*_test.go`)
- Identifies gRPC files (`.proto` files or `grpc/` directories)
- Sets appropriate timeouts (5 minutes for tests, 2 minutes for linting)

## Integration

This script is perfect for:
- Local development workflow
- CI/CD pipelines
- Pre-commit hooks
- Team development standards enforcement
