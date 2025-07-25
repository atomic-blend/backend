#!/bin/bash
set -euo pipefail

# Check if directories are provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <microservice-directory1> [microservice-directory2] ..."
    echo "Example: $0 auth productivity"
    exit 1
fi

MICROSERVICE_DIRS=("$@")
GLOBAL_CHANGELOG="GLOBAL_CHANGELOG.md"

echo "# Global Release Changelog" > "$GLOBAL_CHANGELOG"
echo "" >> "$GLOBAL_CHANGELOG"
echo "This changelog contains all changes for the current release across all microservices." >> "$GLOBAL_CHANGELOG"
echo "" >> "$GLOBAL_CHANGELOG"

for dir in "${MICROSERVICE_DIRS[@]}"; do
    if [ ! -d "$dir" ]; then
        echo "Warning: Directory $dir does not exist, skipping..."
        continue
    fi
    
    if [ ! -f "$dir/CHANGELOG.md" ]; then
        echo "Warning: CHANGELOG.md not found in $dir, skipping..."
        continue
    fi
    
    # Get the latest tag for this microservice
    cd "$dir"
    LATEST_TAG=$(git describe --abbrev=0 --tags --match "$(grep "tag_prefix" cog.toml | cut -d'"' -f2)*" 2>/dev/null || echo "")
    cd ..
    
    if [ -z "$LATEST_TAG" ]; then
        echo "Warning: No tags found for $dir, skipping..."
        continue
    fi
    
    echo "## $dir ($LATEST_TAG)" >> "$GLOBAL_CHANGELOG"
    echo "" >> "$GLOBAL_CHANGELOG"
    
    # Extract the changelog for the latest version
    CLEAN_TAG=$(echo "$LATEST_TAG" | sed 's/^[^-]*-v//')
    
    # Try different changelog entry formats
    CHANGELOG=$(sed -n "/^## $LATEST_TAG - /,/^## /p" "$dir/CHANGELOG.md" | sed '1d;$d' || echo "")
    
    if [ -z "$CHANGELOG" ]; then
        CHANGELOG=$(sed -n "/^## $CLEAN_TAG - /,/^## /p" "$dir/CHANGELOG.md" | sed '1d;$d' || echo "")
    fi
    
    if [ -z "$CHANGELOG" ]; then
        CHANGELOG=$(sed -n "/^## \[$CLEAN_TAG\] - /,/^## /p" "$dir/CHANGELOG.md" | sed '1d;$d' || echo "")
    fi
    
    if [ -n "$CHANGELOG" ]; then
        echo "$CHANGELOG" >> "$GLOBAL_CHANGELOG"
    else
        echo "No changelog entries found for this version." >> "$GLOBAL_CHANGELOG"
    fi
    
    echo "" >> "$GLOBAL_CHANGELOG"
    echo "---" >> "$GLOBAL_CHANGELOG"
    echo "" >> "$GLOBAL_CHANGELOG"
done

echo "Global changelog generated at $GLOBAL_CHANGELOG"
