#!/bin/bash
set -euo pipefail

# Get the microservice directory from the first argument
MICROSERVICE_DIR="${1:-}"

if [[ -z "$MICROSERVICE_DIR" ]]; then
  echo "Usage: $0 <microservice_directory>"
  exit 1
fi

# Dynamic tag prefix based on microservice
TAG_PREFIX="${MICROSERVICE_DIR}-v"

# Patterns
FINAL_TAG_PATTERN="${TAG_PREFIX}[0-9]*.[0-9]*.[0-9]*"
PRE_TAG_PATTERN="${FINAL_TAG_PATTERN}-rc.[0-9]*"

# Safely get the latest final tag (ignore pre-releases)
latest_final=$(git tag --list "$FINAL_TAG_PATTERN" | grep -v -- '-rc' | sort -V | tail -n 1 || true)

# Get latest pre-release (rc.*)
latest_rc=$(git tag --list "$PRE_TAG_PATTERN" | sort -V | tail -n 1 || true)

echo "Latest final tag: $latest_final"
echo "Latest RC tag: $latest_rc"

# Extract base version (remove prefix and pre-release suffix)
extract_base_version() {
  echo "$1" | sed -E "s/^${TAG_PREFIX//\//\\/}//" | sed -E 's/-rc\.[0-9]+$//'
}

base_final=""
base_rc=""

if [[ -n "$latest_final" ]]; then
  base_final=$(extract_base_version "$latest_final")
fi

if [[ -n "$latest_rc" ]]; then
  base_rc=$(extract_base_version "$latest_rc")
fi

echo "Base final version: $base_final"
echo "Base RC version: $base_rc"

# Determine next RC
if [[ -z "$latest_rc" ]]; then
  echo "No RC found, starting new RC cycle: rc.1"
  next_rc="rc.1"
else
  # Check what the new version will be by doing a dry run
  temp_dry_run=$(~/.cargo/bin/cog bump --auto --package "$MICROSERVICE_DIR" --dry-run 2>&1) || {
    # If dry run fails due to no conventional commits, use current version
    temp_dry_run=""
  }
  
  if [[ -n "$temp_dry_run" ]]; then
    # Extract the new version that would be created
    new_version=$(echo "$temp_dry_run" | tail -n 1 | sed -E 's/.*\/v([0-9]+\.[0-9]+\.[0-9]+).*/\1/')
    echo "New version that would be created: $new_version"
    echo "Current RC base version: $base_rc"
    
    # Compare with current RC base version using sort for proper version comparison
    if [[ -n "$base_rc" ]] && [[ "$(echo -e "$new_version\n$base_rc" | sort -V | tail -n 1)" == "$new_version" ]] && [[ "$new_version" != "$base_rc" ]]; then
      # New version is higher than current RC version, reset to rc.1
      echo "New version $new_version is higher than current RC version $base_rc, starting new RC cycle: rc.1"
      next_rc="rc.1"
    else
      # Same version or lower version, increment RC
      rc_number=$(echo "$latest_rc" | sed -nE 's/.*-rc\.([0-9]+)$/\1/p')
      next_rc="rc.$((rc_number + 1))"
      echo "Same version or lower version, incrementing RC: $latest_rc → $next_rc"
    fi
  else
    # No conventional commits, increment RC
    rc_number=$(echo "$latest_rc" | sed -nE 's/.*-rc\.([0-9]+)$/\1/p')
    next_rc="rc.$((rc_number + 1))"
    echo "No conventional commits, incrementing RC: $latest_rc → $next_rc"
  fi
fi

# Run cog bump from the root directory
echo "Running: cog bump --auto --pre $next_rc from root directory"

# First, do a dry run to get the exact tag that would be created
dry_run_output=$(~/.cargo/bin/cog bump --auto --pre "$next_rc" --package "$MICROSERVICE_DIR" --dry-run 2>&1) || {
  dry_run_exit_code=$?
  # Check if the dry run error is due to no conventional commits found
  if [[ $dry_run_exit_code -eq 1 ]] && echo "$dry_run_output" | grep -q "No conventional commit found to bump current version"; then
    echo "No conventional commits found to bump version - this is expected and considered successful"
    exit 0
  else
    echo "Error: dry run failed"
    echo "$dry_run_output"
    exit 1
  fi
}

expected_tag=$(echo "$dry_run_output" | tail -n 1 | tr -d '\n')

# Now run the actual bump
output=$(~/.cargo/bin/cog bump --auto --pre "$next_rc" --package "$MICROSERVICE_DIR" 2>&1) || {
  exit_code=$?
  # Check if the error is due to no conventional commits found
  if [[ $exit_code -eq 1 ]] && echo "$output" | grep -q "No conventional commit found to bump current version"; then
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