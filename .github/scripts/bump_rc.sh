#!/bin/bash
set -euo pipefail

TAG_PREFIX="v"  # Leave empty if no prefix

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
  echo "Starting new RC cycle: rc.1"
  next_rc="rc.1"
else
  rc_number=$(echo "$latest_rc" | sed -nE 's/.*-rc\.([0-9]+)$/\1/p')
  next_rc="rc.$((rc_number + 1))"
  echo "Incrementing RC: $latest_rc → $next_rc"
fi

echo "Running: cog bump --pre $next_rc"
cog bump --pre "$next_rc"