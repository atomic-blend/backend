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
elif [[ -z "$latest_final" ]] || [[ "$base_rc" > "$base_final" ]]; then
  # RC version is ahead of final version, increment RC
  rc_number=$(echo "$latest_rc" | sed -nE 's/.*-rc\.([0-9]+)$/\1/p')
  next_rc="rc.$((rc_number + 1))"
  echo "RC version ahead of final version, incrementing RC: $latest_rc → $next_rc"
elif [[ "$base_rc" == "$base_final" ]]; then
  # RC version matches final version, increment RC
  rc_number=$(echo "$latest_rc" | sed -nE 's/.*-rc\.([0-9]+)$/\1/p')
  next_rc="rc.$((rc_number + 1))"
  echo "RC version matches final version, incrementing RC: $latest_rc → $next_rc"
else
  # RC version is behind final version, start new RC cycle
  echo "RC version behind final version, starting new RC cycle: rc.1"
  next_rc="rc.1"
fi

# Run cog bump from the root directory
echo "Running: cog bump --auto --pre $next_rc from root directory"
~/.cargo/bin/cog bump --auto --pre "$next_rc" --skip-untracked --package "$MICROSERVICE_DIR"