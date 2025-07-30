#!/bin/bash
set -euo pipefail

# Get the microservice directory from the first argument
MICROSERVICE_DIR="${1:-}"

if [[ -z "$MICROSERVICE_DIR" ]]; then
  echo "Usage: $0 <microservice_directory>"
  exit 1
fi

# Dynamic tag prefix based on microservice
TAG_PREFIX="${MICROSERVICE_DIR}/v"

echo "Bumping final version for microservice: $MICROSERVICE_DIR"
echo "Tag prefix: $TAG_PREFIX"

# Run cog bump from the root directory for final release
echo "Running: cog bump --auto --package $MICROSERVICE_DIR"

# First, do a dry run to get the exact tag that would be created
dry_run_output=$(~/.cargo/bin/cog bump --auto --package "$MICROSERVICE_DIR" --dry-run 2>&1) || {
  dry_run_exit_code=$?
  # Check if the dry run error is due to no conventional commits found
  if [[ $dry_run_exit_code -eq 1 ]] && echo "$dry_run_output" | grep -qi "No conventional commits found"; then
    echo "No conventional commits found to bump version - this is expected and considered successful"
    exit 0
  else
    echo "Error: dry run failed"
    echo "$dry_run_output"
    exit 1
  fi
}

# Extract the tag from the dry run output - look for the pattern like "grpc/v1.2.3"
expected_tag=$(echo "$dry_run_output" | grep -E "${MICROSERVICE_DIR}/v[0-9]+\.[0-9]+\.[0-9]+" | tail -n 1 | tr -d '\n')

# Now run the actual bump
output=$(~/.cargo/bin/cog bump --auto --package "$MICROSERVICE_DIR" 2>&1) || {
  exit_code=$?
  # Check if the error is due to no conventional commits found
  if [[ $exit_code -eq 1 ]] && echo "$output" | grep -qi "No conventional commits found"; then
    echo "No conventional commits found to bump version - this is expected and considered successful"
    exit 0
  else
    echo "Error: failed to bump version"
    echo "$output"
    exit 1
  fi
}

echo "Version bumped successfully"
# Output the newly created tag - get the actual tag from the dry run, not from the bump output
echo "NEW_TAG:$expected_tag"