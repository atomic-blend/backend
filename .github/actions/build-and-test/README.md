# Build and Test Backend Services Action

This GitHub Action builds and tests Go backend services with code coverage reporting and email template compilation.

## Features

- ✅ Automatic Go and Node.js setup
- 🧹 Go code linting with golint
- 🧪 Test execution with code coverage
- 📦 Application compilation
- 📧 Maizzle email template compilation
- 📊 Coverage reports upload to Codecov (optional)
- 🚀 Build artifacts upload

## Usage

### Basic usage

```yaml
steps:
  - name: Checkout code
    uses: actions/checkout@v3
    
  - name: Build and Test
    uses: ./.github/actions/build-and-test
```

### Usage with custom parameters

```yaml
steps:
  - name: Checkout code
    uses: actions/checkout@v3
    
  - name: Build and Test
    uses: ./.github/actions/build-and-test
    with:
      go-version: '1.22'
      node-version: '20'
      working-directory: './backend'
      maizzle-directory: './backend/maizzle'
      codecov-token: ${{ secrets.CODECOV_TOKEN }}
      upload-artifacts: 'true'
```

## Input Parameters

| Parameter | Description | Required | Default |
|-----------|-------------|----------|---------|
| `go-version` | Go version to use | No | `1.21` |
| `node-version` | Node.js version to use | No | `23.9` |
| `working-directory` | Working directory for the action | No | `.` |
| `maizzle-directory` | Directory containing Maizzle templates | No | `./maizzle` |
| `codecov-token` | Codecov token for coverage upload | No | - |
| `upload-artifacts` | Whether to upload build artifacts | No | `true` |

## Outputs

| Output | Description |
|--------|-------------|
| `coverage-file` | Path to the code coverage output file |
| `binary-path` | Path to the compiled binary |

## Complete workflow example

```yaml
name: CI/CD Pipeline

on:
  pull_request:
    branches: [main, dev]
  push:
    branches: [main]

jobs:
  backend-tests:
    runs-on: ubuntu-22.04
    strategy:
      matrix:
        service: [auth, productivity]
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
        
      - name: Build and Test ${{ matrix.service }}
        uses: ./.github/actions/build-and-test
        with:
          working-directory: ./${{ matrix.service }}
          maizzle-directory: ./${{ matrix.service }}/maizzle
          codecov-token: ${{ secrets.CODECOV_TOKEN }}
        
      - name: Use outputs
        run: |
          echo "Coverage file: ${{ steps.build-test.outputs.coverage-file }}"
          echo "Binary path: ${{ steps.build-test.outputs.binary-path }}"
```

## Prerequisites

- The project must contain a `go.mod` file in the working directory
- The Maizzle directory must contain a `package.json` with a `build` script
- Go tests must be present in the project

## Generated Artifacts

- **coverage-report**: HTML coverage report
- **backend-binary**: Compiled application binary

## Notes

- The action uses Go modules and automatically caches them
- Node.js also caches npm dependencies
- Code coverage is uploaded to Codecov only if a token is provided
- Artifacts are uploaded only if `upload-artifacts` is set to `true`
