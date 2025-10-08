#!/bin/bash

# Utility functions for the setup script
# This module contains common functions used across the setup process

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to show help
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Fetch latest non-RC versions of custom Docker images from GitHub Container Registry"
    echo "and optionally update the .env file with new versions."
    echo ""
    echo "Options:"
    echo "  -c, --compose-file FILE    Specify docker-compose file (default: docker-compose.yaml)"
    echo "  -e, --env-file FILE        Specify .env file (default: .env)"
    echo "  -b, --branch BRANCH        Specify branch to download files from (default: main)"
    echo "  -h, --help                 Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Use default files"
    echo "  $0 -c docker-compose-dev.yaml        # Use custom compose file"
    echo "  $0 -e .env.production                # Use custom env file"
    echo "  $0 -b develop                         # Download from develop branch"
    echo "  $0 -c docker-compose-dev.yaml -e .env.dev -b feature/new-setup  # Use custom files and branch"
}

# Function to parse command line arguments
parse_arguments() {
    local docker_compose_file="docker-compose.yaml"
    local env_file=".env"
    local branch="main"
    local github_repo="https://raw.githubusercontent.com/atomic-blend/backend"

    while [[ $# -gt 0 ]]; do
        case $1 in
            -c|--compose-file)
                docker_compose_file="$2"
                shift 2
                ;;
            -e|--env-file)
                env_file="$2"
                shift 2
                ;;
            -b|--branch)
                branch="$2"
                shift 2
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                echo -e "${RED}Error: Unknown option $1${NC}" >&2
                echo "Use -h or --help for usage information." >&2
                exit 1
                ;;
        esac
    done

    # Export variables for use in other modules
    export DOCKER_COMPOSE_FILE="$docker_compose_file"
    export ENV_FILE="$env_file"
    export BRANCH="$branch"
    export GITHUB_REPO="$github_repo"
}

# Function to check if required tools are installed
check_required_tools() {
    command -v curl >/dev/null 2>&1 || { echo -e "${RED}Error: curl is required but not installed.${NC}" >&2; exit 1; }
    command -v jq >/dev/null 2>&1 || { echo -e "${RED}Error: jq is required but not installed.${NC}" >&2; exit 1; }
}

# Function to check for GitHub token
check_github_token() {
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
}

# Function to get current version from .env file
get_current_version() {
    local env_var=$1
    local env_file="${ENV_FILE:-.env}"
    
    # Get the absolute path to the docker directory
    local docker_dir="$(dirname "$(dirname "${BASH_SOURCE[0]}")")"
    
    # Convert relative path to absolute path
    if [[ "$env_file" != /* ]]; then
        env_file="$docker_dir/$env_file"
    fi
    
    local current_version=$(grep "^${env_var}=" "$env_file" 2>/dev/null | cut -d'=' -f2 | tr -d '"' | tr -d "'")
    
    if [ -z "$current_version" ]; then
        echo "latest"
    else
        echo "$current_version"
    fi
}
