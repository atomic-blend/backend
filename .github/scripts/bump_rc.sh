#!/bin/bash
set -euo pipefail

# Get the microservice directory from the first argument
MICROSERVICE_DIR="${1:-}"

if [[ -z "$MICROSERVICE_DIR" ]]; then
  echo "Usage: $0 <microservice_directory>"
  exit 1
fi

# Get the current commit hash (short version)
COMMIT_HASH=$(git rev-parse --short HEAD)
RC_IDENTIFIER="rc-${COMMIT_HASH}"
echo "Using RC identifier: $RC_IDENTIFIER"

# Run cog bump from the root directory using commit hash as pre-release identifier
echo "Running: cog bump --auto --pre $RC_IDENTIFIER from root directory"

# First, do a dry run to get the exact tag that would be created
dry_run_output=$(~/.cargo/bin/cog bump --auto --pre "$RC_IDENTIFIER" --package "$MICROSERVICE_DIR" --dry-run 2>&1)
dry_run_exit_code=$?

echo "Dry run output: $dry_run_output"
echo "Dry run exit code: $dry_run_exit_code"

# Check if the dry run indicates no conventional commits found (regardless of exit code)
if echo "$dry_run_output" | grep -qi "No conventional commits found"; then
  echo "No conventional commits found to bump version - this is expected and considered successful"
  exit 0
fi

# If dry run failed for other reasons, exit with error
if [[ $dry_run_exit_code -ne 0 ]]; then
  echo "Error: dry run failed"
  echo "$dry_run_output"
  exit 1
fi

# Extract the tag from the dry run output - look for the pattern like "grpc/v1.2.3-rc-abc1234"
expected_tag=$(echo "$dry_run_output" | grep -E "${MICROSERVICE_DIR}/v[0-9]+\.[0-9]+\.[0-9]+-${RC_IDENTIFIER}" | tail -n 1 | tr -d '\n')

# Now run the actual bump
output=$(~/.cargo/bin/cog bump --auto --pre "$RC_IDENTIFIER" --package "$MICROSERVICE_DIR" 2>&1) || {
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
# Output the newly created tag
echo "NEW_TAG:$expected_tag"