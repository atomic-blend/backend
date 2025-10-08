#!/bin/bash

# Download functionality for setup script
# This module handles downloading files from GitHub repository

# Source utils for colors and common functions
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/utils.sh"

# Function to download file from GitHub repository
download_file() {
    local file_path=$1
    local local_file=$2
    local url="${GITHUB_REPO}/${BRANCH}/${file_path}"
    
    echo -e "${YELLOW}Downloading $local_file from GitHub repository...${NC}"
    echo -e "${BLUE}URL: $url${NC}"
    
    if curl -s -f -o "$local_file" "$url"; then
        echo -e "${GREEN}Successfully downloaded $local_file${NC}"
        return 0
    else
        echo -e "${RED}Failed to download $local_file from $url${NC}" >&2
        return 1
    fi
}

# Function to check file status and handle downloads
check_and_download_files() {
    local docker_compose_file="${DOCKER_COMPOSE_FILE:-docker-compose.yaml}"
    local env_file="${ENV_FILE:-.env}"
    
    # Get the absolute path to the docker directory
    local docker_dir="$(dirname "$SCRIPT_DIR")"
    
    # Convert relative paths to absolute paths
    if [[ "$docker_compose_file" != /* ]]; then
        docker_compose_file="$docker_dir/$docker_compose_file"
    fi
    if [[ "$env_file" != /* ]]; then
        env_file="$docker_dir/$env_file"
    fi
    
    # Check which files exist and which need to be downloaded
    local missing_files=()
    local existing_files=()

    if [ ! -f "$docker_compose_file" ]; then
        missing_files+=("docker-compose.yaml → $docker_compose_file")
    else
        existing_files+=("$docker_compose_file")
    fi

    if [ ! -f "$env_file" ]; then
        missing_files+=(".env.example → $env_file")
    else
        existing_files+=("$env_file")
    fi

    # Show status of files
    if [ ${#existing_files[@]} -gt 0 ]; then
        echo -e "${GREEN}Found existing files:${NC}"
        for file in "${existing_files[@]}"; do
            echo -e "${GREEN}  ✓ $file${NC}"
        done
        echo ""
    fi

    # Ask for confirmation only if files need to be downloaded
    if [ ${#missing_files[@]} -gt 0 ]; then
        echo -e "${YELLOW}Missing files that need to be downloaded:${NC}"
        for file in "${missing_files[@]}"; do
            echo -e "${YELLOW}  - $file${NC}"
        done
        echo ""
        echo -e "${BLUE}Repository: atomic-blend/backend${NC}"
        echo -e "${BLUE}Branch: $BRANCH${NC}"
        echo ""
        read -p "Enter 'y' or 'yes' to download missing files: " -r
        echo ""
        
        if [[ ! $REPLY =~ ^[Yy]([Ee][Ss])?$ ]]; then
            echo -e "${YELLOW}Download cancelled. Using existing files only.${NC}"
            echo ""
        else
            echo -e "${BLUE}Downloading missing files...${NC}"
            echo ""
            
            # Download docker-compose.yaml if missing
            if [ ! -f "$docker_compose_file" ]; then
                echo -e "${YELLOW}Downloading docker-compose.yaml → $docker_compose_file...${NC}"
                local github_path="docker/docker-compose.yaml"
                
                if ! download_file "$github_path" "$docker_compose_file"; then
                    echo -e "${RED}Error: Failed to download $github_path${NC}" >&2
                    exit 1
                fi
            fi
            
            # Download .env.example if env file is missing
            if [ ! -f "$env_file" ]; then
                echo -e "${YELLOW}Downloading .env.example → $env_file...${NC}"
                local github_path="docker/.env.example"
                local temp_file=".env.example"
                
                if ! download_file "$github_path" "$temp_file"; then
                    echo -e "${RED}Error: Failed to download $github_path${NC}" >&2
                    exit 1
                fi
                
                # Rename .env.example to the target env file
                mv "$temp_file" "$env_file"
                echo -e "${GREEN}Renamed .env.example to $env_file${NC}"
            fi
        fi
    else
        echo -e "${GREEN}All required files already exist. No download needed.${NC}"
        echo ""
    fi
}
