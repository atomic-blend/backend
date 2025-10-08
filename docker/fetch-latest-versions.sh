#!/bin/bash

# Script to fetch latest non-RC versions of custom Docker images from GitHub Container Registry
# This script will output the latest stable versions for each atomic-blend package
# The script will error if it cannot find valid versioned tags (like 0.2.0) for any package

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Fetching latest non-RC versions for Docker images...${NC}"
echo "=================================================="


# Function to get latest version from GitHub Container Registry
get_ghcr_latest() {
    local image=$1
    local org_repo=$(echo $image | sed 's/ghcr\.io\///')
    
    # Get tags from GitHub API with authentication
    local response
    if [ -n "${GITHUB_TOKEN:-}" ]; then
        response=$(curl -s -H "Authorization: Bearer $GITHUB_TOKEN" "https://api.github.com/orgs/$(echo $org_repo | cut -d'/' -f1)/packages/container/$(echo $org_repo | cut -d'/' -f2)/versions")
    else
        response=$(curl -s "https://api.github.com/orgs/$(echo $org_repo | cut -d'/' -f1)/packages/container/$(echo $org_repo | cut -d'/' -f2)/versions")
    fi
    if [ $? -ne 0 ] || [ -z "$response" ]; then
        echo -e "${RED}Error: Failed to fetch versions for $image${NC}" >&2
        exit 1
    fi
    
    # Check for authentication errors
    local error_msg=$(echo "$response" | jq -r '.message' 2>/dev/null)
    if [ "$error_msg" != "null" ] && [ -n "$error_msg" ]; then
        if [[ "$error_msg" == *"Bad credentials"* ]] || [[ "$error_msg" == *"Not Found"* ]]; then
            echo -e "${RED}Error: Authentication failed for $image. Please check your GITHUB_TOKEN${NC}" >&2
        else
            echo -e "${RED}Error: $error_msg for $image${NC}" >&2
        fi
        exit 1
    fi
    
    local latest=$(echo "$response" | \
        jq -r '.[]? | select(.metadata.container.tags[]? | test("^[0-9]+\\.[0-9]+(\\.[0-9]+)?$")) | .metadata.container.tags[]?' 2>/dev/null | \
        grep -E '^[0-9]+\.[0-9]+(\.[0-9]+)?$' | \
        sort -V | tail -1)
    
    if [ -z "$latest" ] || [ "$latest" = "null" ]; then
        echo -e "${RED}Error: No valid version found for $image${NC}" >&2
        exit 1
    else
        echo "$latest"
    fi
}

# Function to get latest version for GitHub Container Registry images
get_latest_version() {
    local image=$1
    local base_image=$(echo $image | cut -d':' -f1)
    
    case $base_image in
        ghcr.io/atomic-blend/*)
            get_ghcr_latest "$base_image"
            ;;
        *)
            echo -e "${RED}Error: Unsupported image type: $base_image${NC}" >&2
            exit 1
            ;;
    esac
}

# Function to get current version from .env file
get_current_version() {
    local env_var=$1
    local current_version=$(grep "^${env_var}=" "$env_file" 2>/dev/null | cut -d'=' -f2 | tr -d '"' | tr -d "'")
    
    if [ -z "$current_version" ]; then
        echo "latest"
    else
        echo "$current_version"
    fi
}

# Extract atomic-blend images from docker-compose.yaml
docker_compose_file="docker-compose.yaml"
env_file=".env"

if [ ! -f "$docker_compose_file" ]; then
    echo -e "${RED}Error: $docker_compose_file not found in current directory${NC}" >&2
    exit 1
fi

if [ ! -f "$env_file" ]; then
    echo -e "${RED}Error: $env_file not found in current directory${NC}" >&2
    exit 1
fi

# Extract atomic-blend images and their env vars from docker-compose.yaml
images=()
env_vars=()

# Parse each image line to extract both the base image and env var name
while IFS= read -r line; do
    # Extract the env var name from ${ENV_VAR:-latest} pattern
    # Account for YAML indentation (spaces before "image:")
    if [[ $line =~ ^[[:space:]]*image:[[:space:]]*ghcr\.io/atomic-blend/([^:]+):\$\{([^}]+):-latest\} ]]; then
        service_name="${BASH_REMATCH[1]}"
        env_var="${BASH_REMATCH[2]}"
        base_image="ghcr.io/atomic-blend/${service_name}"
        
        images+=("$base_image")
        env_vars+=("$env_var")
    fi
done < <(grep -E "image:\s*ghcr\.io/atomic-blend/" "$docker_compose_file")

# Debug: Check if we found any images
if [ ${#images[@]} -eq 0 ]; then
    echo -e "${RED}Error: No atomic-blend images found in $docker_compose_file${NC}" >&2
    echo -e "${YELLOW}Debug: Looking for lines matching 'image: ghcr.io/atomic-blend/'${NC}" >&2
    grep -E "image:\s*ghcr\.io/atomic-blend/" "$docker_compose_file" >&2 || echo "No matching lines found" >&2
    exit 1
fi

# Check if required tools are installed
command -v curl >/dev/null 2>&1 || { echo -e "${RED}Error: curl is required but not installed.${NC}" >&2; exit 1; }
command -v jq >/dev/null 2>&1 || { echo -e "${RED}Error: jq is required but not installed.${NC}" >&2; exit 1; }

# Check for GitHub token
if [ -z "${GITHUB_TOKEN:-}" ]; then
    echo -e "${YELLOW}GitHub token is required to fetch package versions.${NC}"
    echo -e "${BLUE}You can either:${NC}"
    echo "1. Set the GITHUB_TOKEN environment variable"
    echo "2. Enter your token now (it will be used for this session only)"
    echo ""
    read -s -p "Enter your GitHub token: " GITHUB_TOKEN
    echo ""
    
    if [ -z "$GITHUB_TOKEN" ]; then
        echo -e "${RED}Error: GitHub token is required to continue.${NC}" >&2
        exit 1
    fi
    
    echo -e "${GREEN}GitHub token provided. Continuing...${NC}"
    echo ""
fi

echo -e "${YELLOW}Checking for latest versions...${NC}"
echo ""

# Create table header
printf "%-40s %-15s %-15s %-10s\n" "IMAGE" "CURRENT" "LATEST" "STATUS"
printf "%-40s %-15s %-15s %-10s\n" "----------------------------------------" "---------------" "---------------" "----------"

# Fetch and display latest versions in table format
for i in "${!images[@]}"; do
    base_image="${images[$i]}"
    env_var="${env_vars[$i]}"
    current_version=$(get_current_version "$env_var")
    latest_version=$(get_latest_version "$base_image")
    
    # Extract just the service name for cleaner display
    service_name=$(echo $base_image | sed 's/ghcr\.io\/atomic-blend\///')
    
    # Determine status and print with colors
    if [ "$latest_version" = "$current_version" ]; then
        printf "%-40s %-15s %-15s " "$service_name" "$current_version" "$latest_version"
        echo -e "${GREEN}UP TO DATE${NC}"
    else
        printf "%-40s %-15s %-15s " "$service_name" "$current_version" "$latest_version"
        echo -e "${YELLOW}OUTDATED${NC}"
    fi
done

echo ""
echo -e "${BLUE}Summary:${NC}"
echo "=========="

# Check if there are any outdated services
outdated_count=0
outdated_services=()

for i in "${!images[@]}"; do
    base_image="${images[$i]}"
    env_var="${env_vars[$i]}"
    current_version=$(get_current_version "$env_var")
    latest_version=$(get_latest_version "$base_image")
    
    # Count outdated services
    if [ "$latest_version" != "$current_version" ]; then
        outdated_count=$((outdated_count + 1))
        outdated_services+=("${env_var}=$latest_version")
    fi
done

# Only show update section if there are outdated services
if [ $outdated_count -gt 0 ]; then
    echo -e "${YELLOW}Add these to your .env.versions file:${NC}"
    echo ""
    
    for service in "${outdated_services[@]}"; do
        echo "$service"
    done
    
    echo ""
    echo -e "${BLUE}Would you like to automatically update the .env file with these new versions?${NC}"
    read -p "Enter 'y' or 'yes' to proceed: " -r
    echo ""
    
    if [[ $REPLY =~ ^[Yy]([Ee][Ss])?$ ]]; then
        echo -e "${YELLOW}Updating .env file...${NC}"
        
        # Create backup of .env file
        cp "$env_file" "${env_file}.backup.$(date +%Y%m%d_%H%M%S)"
        echo -e "${BLUE}Backup created: ${env_file}.backup.$(date +%Y%m%d_%H%M%S)${NC}"
        
        # Update each outdated service in .env file
        for service in "${outdated_services[@]}"; do
            env_var=$(echo "$service" | cut -d'=' -f1)
            new_version=$(echo "$service" | cut -d'=' -f2)
            
            # Check if the variable exists in .env file
            if grep -q "^${env_var}=" "$env_file"; then
                # Get the original line to check for quotes
                original_line=$(grep "^${env_var}=" "$env_file")
                
                # Check if the original value was quoted (double quotes)
                if [[ $original_line =~ ^${env_var}=\".*\"$ ]]; then
                    # Replace with double quotes
                    sed -i.tmp "s/^${env_var}=.*/${env_var}=\"${new_version}\"/" "$env_file"
                # Check if the original value was quoted (single quotes)
                elif [[ $original_line =~ ^${env_var}=\'.*\'$ ]]; then
                    # Replace with single quotes
                    sed -i.tmp "s/^${env_var}=.*/${env_var}='${new_version}'/" "$env_file"
                else
                    # Replace without quotes
                    sed -i.tmp "s/^${env_var}=.*/${env_var}=${new_version}/" "$env_file"
                fi
                echo -e "${GREEN}Updated ${env_var} to ${new_version}${NC}"
            else
                # Add new variable with quotes (default format)
                echo "${env_var}=\"${new_version}\"" >> "$env_file"
                echo -e "${GREEN}Added ${env_var}=\"${new_version}\"${NC}"
            fi
        done
        
        # Remove temporary file created by sed
        rm -f "${env_file}.tmp"
        
        echo -e "${GREEN}Successfully updated .env file!${NC}"
    else
        echo -e "${YELLOW}Update cancelled. You can manually update the .env file with the values shown above.${NC}"
    fi
else
    echo -e "${GREEN}All services are up to date!${NC}"
fi

echo ""
echo -e "${GREEN}Script completed successfully!${NC}"