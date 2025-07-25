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
~/.cargo/bin/cog bump --auto --skip-untracked --package "$MICROSERVICE_DIR"