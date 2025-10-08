#!/bin/bash

# Main setup script for atomic-blend Docker services
# This script orchestrates the entire setup process by calling modular functions

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$SCRIPT_DIR/scripts"

# Source all the modules
source "$SCRIPTS_DIR/utils.sh"
source "$SCRIPTS_DIR/download.sh"
source "$SCRIPTS_DIR/version-check.sh"
source "$SCRIPTS_DIR/env-update.sh"

# Main execution function
main() {
    echo -e "${BLUE}Fetching latest non-RC versions for Docker images...${NC}"
    echo "=================================================="
    echo -e "${BLUE}Using compose file: ${DOCKER_COMPOSE_FILE:-docker-compose.yaml}${NC}"
    echo -e "${BLUE}Using env file: ${ENV_FILE:-.env}${NC}"
    echo -e "${BLUE}Using branch: ${BRANCH:-main}${NC}"
    echo ""

    # Check and download files if needed
    check_and_download_files

    # Extract images and environment variables from docker-compose.yaml
    extract_images_and_env_vars

    # Display version comparison table
    display_version_table

    # Check for outdated services
    check_outdated_services

    # Display update summary and handle updates
    display_update_summary

    echo ""
    echo -e "${GREEN}Script completed successfully!${NC}"
}

# Parse command line arguments
parse_arguments "$@"

# Check required tools
check_required_tools

# Check for GitHub token
check_github_token

# Run main function
main
