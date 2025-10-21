#!/bin/bash

# Version checking functionality for setup script
# This module handles fetching latest versions from GitHub Container Registry

# Source utils for colors and common functions
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/utils.sh"

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

# Function to extract images and env vars from docker-compose.yaml
extract_images_and_env_vars() {
    local docker_compose_file="${DOCKER_COMPOSE_FILE:-docker-compose.yaml}"
    local images=()
    local env_vars=()

    # Get the absolute path to the docker-compose file
    # If it's a relative path, make it relative to the docker directory
    if [[ "$docker_compose_file" != /* ]]; then
        # Get the parent directory of the scripts directory (which is the docker directory)
        local docker_dir="$(dirname "$SCRIPT_DIR")"
        docker_compose_file="$docker_dir/$docker_compose_file"
    fi

    # Check if the file exists
    if [ ! -f "$docker_compose_file" ]; then
        echo -e "${RED}Error: $docker_compose_file not found${NC}" >&2
        exit 1
    fi


    # Parse each image line to extract both the base image and env var name
    while IFS= read -r line; do
        # Extract service name and env var using sed for more reliable parsing
        local service_name=$(echo "$line" | sed -n 's/.*ghcr\.io\/atomic-blend\/\([^:]*\):.*/\1/p')
        local env_var=$(echo "$line" | sed -n 's/.*\${\([^}]*\):-latest}.*/\1/p')
        
        if [ -n "$service_name" ] && [ -n "$env_var" ]; then
            local base_image="ghcr.io/atomic-blend/${service_name}"
            
            images+=("$base_image")
            env_vars+=("$env_var")
        fi
    done < <(grep -E "image:\s*ghcr\.io/atomic-blend/" "$docker_compose_file")

    # Check if we found any images
    if [ ${#images[@]} -eq 0 ]; then
        echo -e "${RED}Error: No atomic-blend images found in $docker_compose_file${NC}" >&2
        exit 1
    fi

    # Set global arrays for use in other modules
    IMAGES=("${images[@]}")
    ENV_VARS=("${env_vars[@]}")
    
}

# Function to display version comparison table
display_version_table() {
    echo -e "${YELLOW}Checking for latest versions...${NC}"
    echo ""

    # Create table header
    printf "%-40s %-15s %-15s %-10s\n" "IMAGE" "CURRENT" "LATEST" "STATUS"
    printf "%-40s %-15s %-15s %-10s\n" "----------------------------------------" "---------------" "---------------" "----------"

    # Fetch and display latest versions in table format
    for i in "${!IMAGES[@]}"; do
        local base_image="${IMAGES[$i]}"
        local env_var="${ENV_VARS[$i]}"
        local current_version=$(get_current_version "$env_var")
        local latest_version=$(get_latest_version "$base_image")
        
        # Extract just the service name for cleaner display
        local service_name=$(echo $base_image | sed 's/ghcr\.io\/atomic-blend\///')
        
        # Determine status and print with colors
        if [ "$latest_version" = "$current_version" ]; then
            printf "%-40s %-15s %-15s " "$service_name" "$current_version" "$latest_version"
            echo -e "${GREEN}UP TO DATE${NC}"
        else
            printf "%-40s %-15s %-15s " "$service_name" "$current_version" "$latest_version"
            echo -e "${YELLOW}OUTDATED${NC}"
        fi
    done
}
