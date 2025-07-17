# Test gRPC Action

This GitHub Action installs the Buf CLI and runs comprehensive validation on Protocol Buffer files, including building, linting, and formatting checks.

## Description

The `test-grpc` action is designed to validate Protocol Buffer files in your repository by:
- Installing the Buf CLI tool
- Running `buf build` to compile and validate protobuf files
- Performing explicit linting checks
- Verifying code formatting compliance

The action will fail if any of these validation steps encounter errors, ensuring your protobuf files meet quality standards.

## Inputs

### `buf-version`
- **Description**: Version of Buf CLI to install
- **Required**: No
- **Default**: `1.28.1`

### `working-directory`
- **Description**: Working directory where buf commands will be executed
- **Required**: No
- **Default**: `./grpc`

## Prerequisites

- The target directory must contain a `buf.yaml` configuration file
- Protocol Buffer files should be properly structured according to Buf standards

## Usage

### Basic Usage

```yaml
- name: Test gRPC
  uses: ./.github/actions/test-grpc
```

### With Custom Parameters

```yaml
- name: Test gRPC with custom settings
  uses: ./.github/actions/test-grpc
  with:
    buf-version: '1.30.0'
    working-directory: './proto'
```

### Complete Workflow Example

```yaml
name: gRPC Validation

on:
  push:
    branches: [main, dev]
  pull_request:
    branches: [main]

jobs:
  test-grpc:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Test gRPC files
        uses: ./.github/actions/test-grpc
        with:
          buf-version: '1.28.1'
          working-directory: './grpc'
```

## What It Does

### 1. Setup Buf CLI
- Automatically detects the runner's OS and architecture
- Downloads and installs the specified version of Buf CLI
- Supports Linux and macOS on x86_64 and ARM64 architectures
- Verifies successful installation

### 2. Build and Lint Validation
- Checks for the presence of `buf.yaml` configuration file
- Runs `buf build` to compile and validate protobuf files
- Performs comprehensive linting according to Buf rules
- Fails the action if any validation errors are found

### 3. Explicit Lint Check
- Runs `buf lint` as a separate step for detailed linting feedback
- Provides clear error messages for any linting violations

### 4. Format Validation
- Checks that all protobuf files are properly formatted
- Uses `buf format --diff --exit-code` to ensure consistent formatting
- Fails if any files need formatting changes

## Error Handling

The action will fail and provide detailed error messages in the following scenarios:
- `buf.yaml` file is missing from the working directory
- Protobuf files contain syntax errors
- Linting rules are violated
- Files are not properly formatted
- Buf CLI installation fails
- Unsupported OS or architecture

## Supported Platforms

- **Operating Systems**: Linux, macOS
- **Architectures**: x86_64, ARM64 (aarch64)

## Dependencies

- `curl` (for downloading Buf CLI)
- `sudo` access (for installing Buf CLI to system path)

## Best Practices

1. **Configuration**: Ensure your `buf.yaml` file is properly configured with appropriate linting rules
2. **Integration**: Use this action in your CI/CD pipeline to catch protobuf issues early
3. **Formatting**: Run `buf format` locally before committing to avoid formatting failures
4. **Version Pinning**: Consider pinning to a specific Buf version for consistent behavior

## Troubleshooting

### Common Issues

1. **Missing buf.yaml**: Ensure the configuration file exists in your working directory
2. **Permission Errors**: The action requires sudo access to install Buf CLI
3. **Network Issues**: Verify that the runner can access GitHub releases for downloading Buf CLI
4. **Format Failures**: Run `buf format` locally to fix formatting issues before pushing

### Debug Tips

- Check the action logs for detailed error messages
- Verify your `buf.yaml` configuration is valid
- Test locally with the same Buf version used in the action
- Ensure all protobuf files are properly structured
