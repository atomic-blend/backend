#!/bin/bash
set -euo pipefail

# Check if microservice directory is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <microservice-directory>"
    echo "Example: $0 auth"
    exit 1
fi

MICROSERVICE_DIR="$1"

# Check if directory exists
if [ ! -d "$MICROSERVICE_DIR" ]; then
    echo "Error: Directory $MICROSERVICE_DIR does not exist"
    exit 1
fi

# Check if cog.toml exists in the directory
if [ ! -f "$MICROSERVICE_DIR/cog.toml" ]; then
    echo "Error: cog.toml not found in $MICROSERVICE_DIR"
    exit 1
fi

# Extract tag prefix from cog.toml
TAG_PREFIX=$(grep "tag_prefix" "$MICROSERVICE_DIR/cog.toml" | cut -d'"' -f2)

# Patterns
FINAL_TAG_PATTERN="${TAG_PREFIX}[0-9]*.[0-9]*.[0-9]*"
PRE_TAG_PATTERN="${FINAL_TAG_PATTERN}-rc.[0-9]*"

# Safely get the latest final tag (ignore pre-releases)
latest_final=$(git tag --list "$FINAL_TAG_PATTERN" | grep -v -- '-rc' | sort -V | tail -n 1 || true)

# Get latest pre-release (rc.*)
latest_rc=$(git tag --list "$PRE_TAG_PATTERN" | sort -V | tail -n 1 || true)

# Extract base version (remove prefix and pre-release suffix)
extract_base_version() {
  echo "$1" | sed -E "s/^$TAG_PREFIX//" | sed -E 's/-rc\.[0-9]+$//'
}

base_final=""
base_rc=""

if [[ -n "$latest_final" ]]; then
  base_final=$(extract_base_version "$latest_final")
fi

if [[ -n "$latest_rc" ]]; then
  base_rc=$(extract_base_version "$latest_rc")
fi

# Determine next RC
if [[ -z "$latest_rc" ]] || [[ "$base_rc" != "$base_final" ]]; then
  echo "Starting new RC cycle for $MICROSERVICE_DIR: rc.1"
  next_rc="rc.1"
else
  rc_number=$(echo "$latest_rc" | sed -nE 's/.*-rc\.([0-9]+)$/\1/p')
  next_rc="rc.$((rc_number + 1))"
  echo "Incrementing RC for $MICROSERVICE_DIR: $latest_rc → $next_rc"
fi

# Change to microservice directory and run cog bump
echo "Running: cd $MICROSERVICE_DIR && cog bump --pre $next_rc"
cd "$MICROSERVICE_DIR"
~/.cargo/bin/cog bump --pre "$next_rc"